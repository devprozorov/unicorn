package profile

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/oklog/ulid/v2"

	"unicorn-auth/internal/http/httputil"
	"unicorn-auth/internal/http/middleware"
	"unicorn-auth/internal/models"
	"unicorn-auth/internal/repo"
	"unicorn-auth/internal/security"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type upsertReq struct {
	About    string   `json:"about"`
	Location string   `json:"location,omitempty"`
	Links    []string `json:"links,omitempty"`
	Industry string   `json:"industry,omitempty"`
	Website  string   `json:"website,omitempty"`
}

func Register(r *gin.Engine, sec *security.Security, users *repo.UserRepo, profiles *repo.ProfileRepo) {
	api := r.Group("/api")

	// public profile
	api.GET("/profile/:userId", func(c *gin.Context) {
		p, err := profiles.GetByUserID(c.Request.Context(), c.Param("userId"))
		if err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		if p == nil {
			c.JSON(404, gin.H{"ok": false, "error": "not_found"})
			return
		}
		c.JSON(200, gin.H{"ok": true, "profile": p})
	})

	protected := api.Group("")
	protected.Use(middleware.RequireAuth(sec))

	protected.GET("/profile/me", func(c *gin.Context) {
		uid := c.GetString(middleware.CtxUserID)
		p, err := profiles.GetByUserID(c.Request.Context(), uid)
		if err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		c.JSON(200, gin.H{"ok": true, "profile": p})
	})

	protected.POST("/profile/me", func(c *gin.Context) {
		uid := c.GetString(middleware.CtxUserID)
		u, err := users.FindByUserID(c.Request.Context(), uid)
		if err != nil || u == nil || u.Status.Deleted || u.Status.Blocked {
			c.JSON(401, gin.H{"ok": false, "error": "unauthorized"})
			return
		}
		if existing, _ := profiles.GetByUserID(c.Request.Context(), uid); existing != nil {
			c.JSON(http.StatusConflict, gin.H{"ok": false, "error": "conflict"})
			return
		}
		var req upsertReq
		if !httputil.BindJSONStrict(c, &req, 64<<10) {
			return
		}
		p := &models.Profile{
			UserID:      uid,
			Type:        u.Type,
			DisplayName: u.DisplayName,
			About:       strings.TrimSpace(req.About),
			Location:    strings.TrimSpace(req.Location),
			Links:       req.Links,
			Industry:    strings.TrimSpace(req.Industry),
			Website:     strings.TrimSpace(req.Website),
		}
		if err := profiles.Create(c.Request.Context(), p); err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		c.JSON(200, gin.H{"ok": true})
	})

	protected.PATCH("/profile/me", func(c *gin.Context) {
		uid := c.GetString(middleware.CtxUserID)

		// читаем как map, чтобы:
		// - порядок полей не имел значения (в JSON он и так не имеет)
		// - обновлять только присланные поля
		var raw map[string]any
		if !httputil.BindJSONStrict(c, &raw, 64<<10) {
			return
		}

		setProfile := bson.M{}
		setUser := bson.M{}

		// helper: string field
		setStr := func(key string, target bson.M, max int) bool {
			v, ok := raw[key]
			if !ok {
				return true
			}
			s, ok := v.(string)
			if !ok {
				c.JSON(400, gin.H{"ok": false, "error": "bad_request"})
				return false
			}
			s = strings.TrimSpace(s)
			if max > 0 && len(s) > max {
				c.JSON(400, gin.H{"ok": false, "error": "bad_request"})
				return false
			}
			target[key] = s
			return true
		}

		// helper: []string field
		setStrSlice := func(key string, target bson.M, maxItems int, maxLen int) bool {
			v, ok := raw[key]
			if !ok {
				return true
			}
			arr, ok := v.([]any)
			if !ok {
				c.JSON(400, gin.H{"ok": false, "error": "bad_request"})
				return false
			}
			out := make([]string, 0, len(arr))
			for _, item := range arr {
				s, ok := item.(string)
				if !ok {
					c.JSON(400, gin.H{"ok": false, "error": "bad_request"})
					return false
				}
				s = strings.TrimSpace(s)
				if s == "" {
					continue
				}
				if maxLen > 0 && len(s) > maxLen {
					c.JSON(400, gin.H{"ok": false, "error": "bad_request"})
					return false
				}
				out = append(out, s)
				if maxItems > 0 && len(out) > maxItems {
					c.JSON(400, gin.H{"ok": false, "error": "bad_request"})
					return false
				}
			}
			target[key] = out
			return true
		}

		// profile fields
		if !setStr("about", setProfile, 2000) {
			return
		}
		if !setStr("location", setProfile, 128) {
			return
		}
		if !setStr("industry", setProfile, 128) {
			return
		}
		if !setStr("website", setProfile, 256) {
			return
		}
		if !setStrSlice("links", setProfile, 20, 256) {
			return
		}

		// displayName — синхронно в user и profile
		if v, ok := raw["displayName"]; ok {
			s, ok := v.(string)
			if !ok {
				c.JSON(400, gin.H{"ok": false, "error": "bad_request"})
				return
			}
			s = strings.TrimSpace(s)
			if s == "" || len(s) > 64 {
				c.JSON(400, gin.H{"ok": false, "error": "bad_request"})
				return
			}
			setProfile["displayName"] = s
			setUser["displayName"] = s
		}

		// если ничего не прислали — можно вернуть ok
		if len(setProfile) == 0 && len(setUser) == 0 {
			c.JSON(200, gin.H{"ok": true})
			return
		}

		// обновляем профиль
		if len(setProfile) > 0 {
			if err := profiles.UpdateByUserID(c.Request.Context(), uid, setProfile); err != nil {
				c.JSON(500, gin.H{"ok": false, "error": "server_error"})
				return
			}
		}

		// обновляем user displayName
		if len(setUser) > 0 {
			if err := users.UpdateByUserID(c.Request.Context(), uid, setUser); err != nil {
				c.JSON(500, gin.H{"ok": false, "error": "server_error"})
				return
			}
		}

		c.JSON(200, gin.H{"ok": true})
	})

	// ✅ Upload avatar to /uploads/avatars (inside container)
	protected.POST("/profile/me/avatar", func(c *gin.Context) {
		uid := c.GetString(middleware.CtxUserID)

		// 1) лимит размера тела запроса (8MB)
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 8<<20)

		// 2) multipart file
		fh, err := c.FormFile("avatar")
		if err != nil {
			c.JSON(400, gin.H{"ok": false, "error": "bad_request"})
			return
		}

		// 3) открыть файл и проверить content-type
		f, err := fh.Open()
		if err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		defer f.Close()

		head := make([]byte, 512)
		n, _ := io.ReadFull(f, head)
		head = head[:n]
		ct := http.DetectContentType(head)

		ext := ""
		switch ct {
		case "image/jpeg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/webp":
			ext = ".webp"
		default:
			c.JSON(400, gin.H{"ok": false, "error": "unsupported_file_type"})
			return
		}

		// 4) заново открыть файл (потому что мы прочитали первые байты)
		_ = f.Close()
		f, err = fh.Open()
		if err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		defer f.Close()

		// ✅ 5) директория внутри контейнера
		dir := "/uploads/avatars"
		if err := os.MkdirAll(dir, 0755); err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}

		// 6) имя файла
		name := uid + "-" + ulid.Make().String() + ext
		dstPath := filepath.Join(dir, name)

		// 7) сохранить
		out, err := os.Create(dstPath)
		if err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		defer out.Close()

		if _, err := io.Copy(out, f); err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}

		// ✅ 8) url, который будет отдавать nginx
		avatarURL := "/uploads/avatars/" + name

		// 9) записать в профиль
		if err := profiles.UpdateByUserID(c.Request.Context(), uid, bson.M{
			"avatarUrl": avatarURL,
		}); err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}

		c.JSON(200, gin.H{"ok": true, "avatarUrl": avatarURL})
	})

}
