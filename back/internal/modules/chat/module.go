package chat

import (
	"strings"

	"unicorn-auth/internal/http/httputil"
	"unicorn-auth/internal/http/middleware"
	"unicorn-auth/internal/models"
	"unicorn-auth/internal/repo"
	"unicorn-auth/internal/security"

	"github.com/gin-gonic/gin"
)

type sendReq struct {
	Text string `json:"text"`
}

type chatItem struct {
	ApplicationID string `json:"applicationId"`
	VacancyID     string `json:"vacancyId"`
	VacancyTitle  string `json:"vacancyTitle,omitempty"`
	ResumeID      string `json:"resumeId"`
	UserID        string `json:"userId"`
	CompanyID     string `json:"companyId"`
	CompanyName   string `json:"companyName,omitempty"`
	Status        string `json:"status"`
	Message       string `json:"message,omitempty"`
	Viewed        bool   `json:"viewed"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

func Register(r *gin.Engine, sec *security.Security, users *repo.UserRepo, apps *repo.ApplicationRepo, chatRepo *repo.ChatRepo, vac *repo.VacancyRepo, profiles *repo.ProfileRepo) {
	api := r.Group("/api")
	protected := api.Group("")
	protected.Use(middleware.RequireAuth(sec))
	protected.Use(middleware.RequireMFAEnabled(sec, users))

	canAccess := func(c *gin.Context, a *models.Application) bool {
		uid := c.GetString(middleware.CtxUserID)
		ut := c.GetString(middleware.CtxUserType)
		if ut == "user" && a.UserID == uid {
			return true
		}
		if ut == "company" && a.CompanyID == uid {
			return true
		}
		return false
	}

	// GET /api/chats/my - получить список чатов (applications) пользователя
	protected.GET("/chats/my", func(c *gin.Context) {
		uid := c.GetString(middleware.CtxUserID)
		ut := c.GetString(middleware.CtxUserType)

		var items []models.Application
		var err error

		if ut == "user" {
			items, err = apps.ListByUserID(c.Request.Context(), uid)
		} else if ut == "company" {
			items, err = apps.ListByCompanyID(c.Request.Context(), uid)
		} else {
			c.JSON(403, gin.H{"ok": false, "error": "forbidden"})
			return
		}

		if err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}

		c.JSON(200, gin.H{"ok": true, "items": items})
	})

	// GET /api/user/chats - алиас для /api/chats/my
	protected.GET("/user/chats", func(c *gin.Context) {
		uid := c.GetString(middleware.CtxUserID)
		ut := c.GetString(middleware.CtxUserType)

		var items []models.Application
		var err error

		if ut == "user" {
			items, err = apps.ListByUserID(c.Request.Context(), uid)
		} else if ut == "company" {
			items, err = apps.ListByCompanyID(c.Request.Context(), uid)
		} else {
			c.JSON(403, gin.H{"ok": false, "error": "forbidden"})
			return
		}

		if err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}

		// Обогащаем данные названиями вакансий и компаний
		vTitle := map[string]string{}
		cName := map[string]string{}

		out := make([]chatItem, 0, len(items))

		for _, a := range items {
			it := chatItem{
				ApplicationID: a.ApplicationID,
				VacancyID:     a.VacancyID,
				ResumeID:      a.ResumeID,
				UserID:        a.UserID,
				CompanyID:     a.CompanyID,
				Status:        a.Status,
				Message:       a.Message,
				Viewed:        a.Viewed,
				CreatedAt:     a.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
				UpdatedAt:     a.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			}

			// Название вакансии
			if t, ok := vTitle[a.VacancyID]; ok {
				it.VacancyTitle = t
			} else {
				if v, _ := vac.GetByID(c.Request.Context(), a.VacancyID); v != nil {
					vTitle[a.VacancyID] = v.Title
					it.VacancyTitle = v.Title
				}
			}

			// Название компании
			if n, ok := cName[a.CompanyID]; ok {
				it.CompanyName = n
			} else {
				if p, _ := profiles.GetByUserID(c.Request.Context(), a.CompanyID); p != nil {
					cName[a.CompanyID] = p.DisplayName
					it.CompanyName = p.DisplayName
				}
			}

			out = append(out, it)
		}

		c.JSON(200, gin.H{"ok": true, "items": out})
	})

	protected.GET("/chat/:applicationId/messages", func(c *gin.Context) {
		a, err := apps.GetByID(c.Request.Context(), c.Param("applicationId"))
		if err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		if a == nil {
			c.JSON(404, gin.H{"ok": false, "error": "not_found"})
			return
		}
		if !canAccess(c, a) {
			c.JSON(403, gin.H{"ok": false, "error": "forbidden"})
			return
		}
		items, err := chatRepo.List(c.Request.Context(), a.ApplicationID, 100)
		if err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		c.JSON(200, gin.H{"ok": true, "items": items})
	})

	protected.POST("/chat/:applicationId/messages", func(c *gin.Context) {
		a, err := apps.GetByID(c.Request.Context(), c.Param("applicationId"))
		if err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		if a == nil {
			c.JSON(404, gin.H{"ok": false, "error": "not_found"})
			return
		}
		if !canAccess(c, a) {
			c.JSON(403, gin.H{"ok": false, "error": "forbidden"})
			return
		}

		var req sendReq
		if !httputil.BindJSONStrict(c, &req, 32<<10) {
			return
		}
		txt := strings.TrimSpace(req.Text)
		if txt == "" || len(txt) > 2000 {
			c.JSON(400, gin.H{"ok": false, "error": "bad_request"})
			return
		}

		m := &models.ChatMessage{
			ApplicationID: a.ApplicationID,
			SenderID:      c.GetString(middleware.CtxUserID),
			SenderType:    c.GetString(middleware.CtxUserType),
			Text:          txt,
		}
		if err := chatRepo.Create(c.Request.Context(), m); err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}

		// При отправке сообщения снимаем скрытие
		if err := apps.UnhideOnNewMessage(c.Request.Context(), a.ApplicationID); err != nil {
			// Логируем ошибку, но не прерываем процесс
		}

		c.JSON(200, gin.H{"ok": true, "messageId": m.MessageID, "createdAt": m.CreatedAt})
	})
}
