package admin

import (
	"strings"
	"time"

	"unicorn-auth/internal/http/httputil"
	"unicorn-auth/internal/repo"
	"unicorn-auth/internal/security"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type loginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func Register(r *gin.Engine, sec *security.Security, admins *repo.AdminRepo, users *repo.UserRepo) {
	api := r.Group("/api/admin")

	api.POST("/login", func(c *gin.Context) {
		var req loginReq
		if !httputil.BindJSONStrict(c, &req, 16<<10) {
			return
		}
		loginNorm := strings.ToLower(strings.TrimSpace(req.Login))
		a, err := admins.FindByLoginNorm(c.Request.Context(), loginNorm)
		if err != nil || a == nil {
			c.JSON(401, gin.H{"ok": false, "error": "invalid_credentials"})
			return
		}
		ok, _ := security.VerifyPassword(req.Password, a.PasswordHash)
		if !ok {
			c.JSON(401, gin.H{"ok": false, "error": "invalid_credentials"})
			return
		}
		tok, err := sec.Tokens.NewAccessToken(a.AdminID, "admin", []string{"pwd"}, 30*time.Minute)
		if err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		c.JSON(200, gin.H{"ok": true, "accessToken": tok})
	})

	requireAdmin := func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{"ok": false, "error": "unauthorized"})
			return
		}
		claims, err := sec.Tokens.ParseAccess(strings.TrimSpace(strings.TrimPrefix(h, "Bearer ")))
		if err != nil || claims.Type != "admin" {
			c.AbortWithStatusJSON(401, gin.H{"ok": false, "error": "unauthorized"})
			return
		}
		c.Set("adminId", claims.UserID)
		c.Next()
	}

	api.GET("/users", requireAdmin, func(c *gin.Context) {
		typ := strings.TrimSpace(c.Query("type"))
		search := strings.TrimSpace(c.Query("search"))
		blocked := c.Query("blocked")
		deleted := c.Query("deleted")
		premium := c.Query("premium")

		items, err := users.ListWithFilters(c.Request.Context(), typ, search, blocked, deleted, premium, 50, 0)
		if err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		c.JSON(200, gin.H{"ok": true, "items": items})
	})

	api.DELETE("/users/:userId", requireAdmin, func(c *gin.Context) {
		if err := users.SoftDelete(c.Request.Context(), c.Param("userId")); err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		c.JSON(200, gin.H{"ok": true})
	})

	api.POST("/users/:userId/block", requireAdmin, func(c *gin.Context) {
		if err := users.Block(c.Request.Context(), c.Param("userId"), true); err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		c.JSON(200, gin.H{"ok": true})
	})

	api.POST("/users/:userId/unblock", requireAdmin, func(c *gin.Context) {
		if err := users.Block(c.Request.Context(), c.Param("userId"), false); err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		c.JSON(200, gin.H{"ok": true})
	})

	// Получить детальную информацию о пользователе
	api.GET("/users/:userId", requireAdmin, func(c *gin.Context) {
		user, err := users.FindByUserID(c.Request.Context(), c.Param("userId"))
		if err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		if user == nil {
			c.JSON(404, gin.H{"ok": false, "error": "user_not_found"})
			return
		}

		// Формируем детальный ответ
		response := gin.H{
			"ok":          true,
			"userId":      user.UserID,
			"login":       user.Login,
			"displayName": user.DisplayName,
			"type":        user.Type,
			"status": gin.H{
				"deleted": user.Status.Deleted,
				"blocked": user.Status.Blocked,
			},
			"subscription": gin.H{
				"active": user.Subscription.Active,
				"until":  user.Subscription.Until,
			},
			"mfa": gin.H{
				"totpEnabled": user.MFA.TOTP.Enabled,
			},
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		}
		c.JSON(200, response)
	})

	// Редактирование данных пользователя
	api.PATCH("/users/:userId", requireAdmin, func(c *gin.Context) {
		type patchReq struct {
			DisplayName *string `json:"displayName"`
			Login       *string `json:"login"`
		}
		var req patchReq
		if !httputil.BindJSONStrict(c, &req, 16<<10) {
			return
		}

		// Проверяем существование пользователя
		user, err := users.FindByUserID(c.Request.Context(), c.Param("userId"))
		if err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		if user == nil {
			c.JSON(404, gin.H{"ok": false, "error": "user_not_found"})
			return
		}

		update := bson.M{}
		if req.DisplayName != nil {
			update["displayName"] = strings.TrimSpace(*req.DisplayName)
		}
		if req.Login != nil {
			loginTrim := strings.TrimSpace(*req.Login)
			// Проверяем, не занят ли новый логин
			if loginTrim != user.Login {
				existing, _ := users.FindByLoginNorm(c.Request.Context(), normLogin(loginTrim))
				if existing != nil {
					c.JSON(400, gin.H{"ok": false, "error": "login_already_exists"})
					return
				}
				update["login"] = loginTrim
				update["loginNorm"] = normLogin(loginTrim)
			}
		}

		if len(update) == 0 {
			c.JSON(400, gin.H{"ok": false, "error": "no_fields_to_update"})
			return
		}

		if err := users.UpdateByUserID(c.Request.Context(), c.Param("userId"), update); err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}

		c.JSON(200, gin.H{"ok": true})
	})

	// Активировать подписку пользователя
	api.POST("/users/:userId/subscription/activate", requireAdmin, func(c *gin.Context) {
		type activateReq struct {
			Days int `json:"days"` // на сколько дней активировать
		}
		var req activateReq
		if !httputil.BindJSONStrict(c, &req, 16<<10) {
			return
		}

		if req.Days <= 0 || req.Days > 3650 { // максимум 10 лет
			c.JSON(400, gin.H{"ok": false, "error": "invalid_days_value"})
			return
		}

		until := time.Now().UTC().AddDate(0, 0, req.Days)
		update := bson.M{
			"subscription.active": true,
			"subscription.until":  until,
		}

		if err := users.UpdateByUserID(c.Request.Context(), c.Param("userId"), update); err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}

		c.JSON(200, gin.H{"ok": true, "until": until})
	})

	// Деактивировать подписку пользователя
	api.POST("/users/:userId/subscription/deactivate", requireAdmin, func(c *gin.Context) {
		update := bson.M{
			"subscription.active": false,
		}

		if err := users.UpdateByUserID(c.Request.Context(), c.Param("userId"), update); err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}

		c.JSON(200, gin.H{"ok": true})
	})
}

func normLogin(login string) string {
	return strings.ToLower(strings.TrimSpace(login))
}
