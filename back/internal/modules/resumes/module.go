package resumes

import (
	"strings"

	"unicorn-auth/internal/http/httputil"
	"unicorn-auth/internal/http/middleware"
	"unicorn-auth/internal/models"
	"unicorn-auth/internal/repo"
	"unicorn-auth/internal/security"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type createReq struct {
	Title  string   `json:"title"`
	About  string   `json:"about"`
	Skills []string `json:"skills,omitempty"`
	Links  []string `json:"links,omitempty"`
}

func Register(r *gin.Engine, sec *security.Security, users *repo.UserRepo, resumes *repo.ResumeRepo, apps *repo.ApplicationRepo) {
	api := r.Group("/api")

	// shared auth group for both user/company (but MFA required)
	shared := api.Group("")
	shared.Use(middleware.RequireAuth(sec))
	shared.Use(middleware.RequireMFAEnabled(sec, users))

	// ✅ GET /api/resumes/:id
	// user: только своё
	// company: только если есть application на эту company с таким resumeId
	shared.GET("/resumes/:id", func(c *gin.Context) {
		uid := c.GetString(middleware.CtxUserID)
		ut := c.GetString(middleware.CtxUserType)
		id := c.Param("id")

		rr, err := resumes.GetByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		if rr == nil {
			c.JSON(404, gin.H{"ok": false, "error": "not_found"})
			return
		}

		if ut == "user" {
			if rr.UserID != uid {
				c.JSON(403, gin.H{"ok": false, "error": "forbidden"})
				return
			}
			c.JSON(200, gin.H{"ok": true, "resume": rr})
			return
		}

		if ut == "company" {
			ok, err := apps.ExistsCompanyResume(c.Request.Context(), uid, id)
			if err != nil {
				c.JSON(500, gin.H{"ok": false, "error": "server_error"})
				return
			}
			if !ok {
				c.JSON(403, gin.H{"ok": false, "error": "forbidden"})
				return
			}
			c.JSON(200, gin.H{"ok": true, "resume": rr})
			return
		}

		c.JSON(403, gin.H{"ok": false, "error": "forbidden"})
	})

	// my resumes
	protected := api.Group("")
	protected.Use(middleware.RequireAuth(sec))
	protected.Use(middleware.RequireType("user"))
	protected.Use(middleware.RequireMFAEnabled(sec, users))

	protected.GET("/resumes/my", func(c *gin.Context) {
		uid := c.GetString(middleware.CtxUserID)
		items, err := resumes.ListMine(c.Request.Context(), uid)
		if err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		c.JSON(200, gin.H{"ok": true, "items": items})
	})

	protected.POST("/resumes", func(c *gin.Context) {
		uid := c.GetString(middleware.CtxUserID)

		u, err := users.FindByUserID(c.Request.Context(), uid)
		if err != nil || u == nil {
			c.JSON(401, gin.H{"ok": false, "error": "unauthorized"})
			return
		}

		// Проверяем лимиты: 16 для премиум, 2 для обычных
		maxLimit := int64(2)
		if u.Subscription.Active {
			maxLimit = 16
		}

		cnt, err := resumes.CountActiveByUser(c.Request.Context(), uid)
		if err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		if cnt >= maxLimit {
			c.JSON(403, gin.H{"ok": false, "error": "limit_reached", "limit": maxLimit})
			return
		}

		var req createReq
		if !httputil.BindJSONStrict(c, &req, 64<<10) {
			return
		}
		rr := &models.Resume{
			UserID:    uid,
			Title:     strings.TrimSpace(req.Title),
			About:     strings.TrimSpace(req.About),
			Skills:    req.Skills,
			Links:     req.Links,
			IsPremium: u.Subscription.Active,
			ColorCode: "",
		}

		// Если подписка активна, добавляем цветовой код
		if u.Subscription.Active {
			rr.ColorCode = "#FFD700" // Gold color for premium
		}

		if rr.Title == "" || rr.About == "" {
			c.JSON(400, gin.H{"ok": false, "error": "bad_request"})
			return
		}
		if err := resumes.Create(c.Request.Context(), rr); err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		c.JSON(200, gin.H{"ok": true, "resumeId": rr.ResumeID})
	})

	protected.PATCH("/resumes/:id", func(c *gin.Context) {
		uid := c.GetString(middleware.CtxUserID)
		var req createReq
		if !httputil.BindJSONStrict(c, &req, 64<<10) {
			return
		}
		set := bson.M{
			"title":  strings.TrimSpace(req.Title),
			"about":  strings.TrimSpace(req.About),
			"skills": req.Skills,
			"links":  req.Links,
		}
		if err := resumes.Update(c.Request.Context(), c.Param("id"), uid, set); err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		c.JSON(200, gin.H{"ok": true})
	})

	protected.DELETE("/resumes/:id", func(c *gin.Context) {
		uid := c.GetString(middleware.CtxUserID)
		if err := resumes.Delete(c.Request.Context(), c.Param("id"), uid); err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		c.JSON(200, gin.H{"ok": true})
	})
}
