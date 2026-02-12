package vacancies

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
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Location    string   `json:"location,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

func Register(r *gin.Engine, sec *security.Security, users *repo.UserRepo, vac *repo.VacancyRepo) {
	api := r.Group("/api")

	api.GET("/vacancies", func(c *gin.Context) {
		items, err := vac.ListPublic(c.Request.Context(), 50, 0)
		if err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		c.JSON(200, gin.H{"ok": true, "items": items})
	})

	api.GET("/vacancies/:id", func(c *gin.Context) {
		v, err := vac.GetByID(c.Request.Context(), c.Param("id"))
		if err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		if v == nil {
			c.JSON(404, gin.H{"ok": false, "error": "not_found"})
			return
		}
		c.JSON(200, gin.H{"ok": true, "vacancy": v})
	})

	protected := api.Group("")
	protected.Use(middleware.RequireAuth(sec))
	protected.Use(middleware.RequireType("company"))
	protected.Use(middleware.RequireMFAEnabled(sec, users))

	// GET /api/vacancies/my - получить свои вакансии
	protected.GET("/vacancies/my", func(c *gin.Context) {
		uid := c.GetString(middleware.CtxUserID)
		items, err := vac.ListByCompanyID(c.Request.Context(), uid)
		if err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		c.JSON(200, gin.H{"ok": true, "items": items})
	})

	protected.POST("/vacancies", func(c *gin.Context) {
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

		cnt, err := vac.CountActiveByCompany(c.Request.Context(), uid)
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
		v := &models.Vacancy{
			CompanyID:   uid,
			Title:       strings.TrimSpace(req.Title),
			Description: strings.TrimSpace(req.Description),
			Location:    strings.TrimSpace(req.Location),
			Tags:        req.Tags,
			IsPremium:   u.Subscription.Active,
			ColorCode:   "",
		}

		// Если подписка активна, добавляем цветовой код
		if u.Subscription.Active {
			v.ColorCode = "#FFD700" // Gold color for premium
		}

		if v.Title == "" || v.Description == "" {
			c.JSON(400, gin.H{"ok": false, "error": "bad_request"})
			return
		}
		if err := vac.Create(c.Request.Context(), v); err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		c.JSON(200, gin.H{"ok": true, "vacancyId": v.VacancyID})
	})

	protected.PATCH("/vacancies/:id", func(c *gin.Context) {
		uid := c.GetString(middleware.CtxUserID)
		var req createReq
		if !httputil.BindJSONStrict(c, &req, 64<<10) {
			return
		}
		set := bson.M{
			"title":       strings.TrimSpace(req.Title),
			"description": strings.TrimSpace(req.Description),
			"location":    strings.TrimSpace(req.Location),
			"tags":        req.Tags,
		}
		if err := vac.Update(c.Request.Context(), c.Param("id"), uid, set); err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		c.JSON(200, gin.H{"ok": true})
	})

	protected.DELETE("/vacancies/:id", func(c *gin.Context) {
		uid := c.GetString(middleware.CtxUserID)
		if err := vac.Delete(c.Request.Context(), c.Param("id"), uid); err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		c.JSON(200, gin.H{"ok": true})
	})
}
