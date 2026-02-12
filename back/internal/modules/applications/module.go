package applications

import (
	"net/http"
	"strings"

	"unicorn-auth/internal/http/httputil"
	"unicorn-auth/internal/http/middleware"
	"unicorn-auth/internal/models"
	"unicorn-auth/internal/repo"
	"unicorn-auth/internal/security"

	"github.com/gin-gonic/gin"
)

type applyReq struct {
	ApplicationID string `json:"applicationId"`
	VacancyID     string `json:"vacancyId"`
	CompanyID     string `json:"companyId"`
	ResumeID      string `json:"resumeId"`
	Status        string `json:"status"`
	Message       string `json:"message,omitempty"`

	// удобные поля для фронта
	VacancyTitle string `json:"vacancyTitle,omitempty"`
	CompanyName  string `json:"companyName,omitempty"`
}

// апдейт для отлкиков со стороны компании чтоб приходили не только айди а реальные показатели имени и названия вакансии
type inboxItem struct {
	ApplicationID   string `json:"applicationId"`
	VacancyID       string `json:"vacancyId"`
	VacancyTitle    string `json:"vacancyTitle,omitempty"`
	ResumeID        string `json:"resumeId"`
	ResumeTitle     string `json:"resumeTitle,omitempty"`
	UserID          string `json:"userId"`
	UserDisplayName string `json:"userDisplayName,omitempty"`
	UserIsPremium   bool   `json:"userIsPremium"`
	CompanyID       string `json:"companyId"`
	Status          string `json:"status"`
	Message         string `json:"message,omitempty"`
	Viewed          bool   `json:"viewed"`
}

type myAppItem struct {
	ApplicationID string `json:"applicationId"`
	VacancyID     string `json:"vacancyId"`
	CompanyID     string `json:"companyId"`
	ResumeID      string `json:"resumeId"`
	Status        string `json:"status"`
	Message       string `json:"message,omitempty"`

	VacancyTitle       string `json:"vacancyTitle,omitempty"`
	CompanyDisplayName string `json:"companyDisplayName,omitempty"`
}

func Register(r *gin.Engine, sec *security.Security, users *repo.UserRepo,
	vac *repo.VacancyRepo, resumes *repo.ResumeRepo, apps *repo.ApplicationRepo) {

	api := r.Group("/api")
	protected := api.Group("")
	protected.Use(middleware.RequireAuth(sec))
	protected.Use(middleware.RequireMFAEnabled(sec, users))

	// user apply
	protected.POST("/applications", middleware.RequireType("user"), func(c *gin.Context) {
		uid := c.GetString(middleware.CtxUserID)

		var req applyReq
		if !httputil.BindJSONStrict(c, &req, 64<<10) {
			return
		}
		req.Message = strings.TrimSpace(req.Message)

		v, err := vac.GetByID(c.Request.Context(), req.VacancyID)
		if err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		if v == nil || v.Status != "active" {
			c.JSON(404, gin.H{"ok": false, "error": "not_found"})
			return
		}
		rr, err := resumes.GetByID(c.Request.Context(), req.ResumeID)
		if err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		if rr == nil || rr.UserID != uid {
			c.JSON(403, gin.H{"ok": false, "error": "forbidden"})
			return
		}

		a := &models.Application{
			VacancyID: v.VacancyID,
			ResumeID:  rr.ResumeID,
			UserID:    uid,
			CompanyID: v.CompanyID,
			Message:   req.Message,
		}
		if err := apps.Create(c.Request.Context(), a); err != nil {
			c.JSON(http.StatusConflict, gin.H{"ok": false, "error": "conflict"})
			return
		}
		c.JSON(200, gin.H{"ok": true, "applicationId": a.ApplicationID, "status": a.Status})
	})
	// user: my applications
	protected.GET("/applications/my", middleware.RequireType("user"), func(c *gin.Context) {
		uid := c.GetString(middleware.CtxUserID)
		status := c.Query("status") // optional

		items, err := apps.ListByUser(c.Request.Context(), uid, status, 50, 0)
		if err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}

		out := make([]myAppItem, 0, len(items))
		cName := map[string]string{} // cache companyId -> displayName
		for _, a := range items {
			it := myAppItem{
				ApplicationID: a.ApplicationID,
				VacancyID:     a.VacancyID,
				CompanyID:     a.CompanyID,
				ResumeID:      a.ResumeID,
				Status:        a.Status,
				Message:       a.Message,
			}
			// company displayName
			if n, ok := cName[a.CompanyID]; ok {
				it.CompanyDisplayName = n
			} else {
				if cu, _ := users.FindByUserID(c.Request.Context(), a.CompanyID); cu != nil {
					cName[a.CompanyID] = cu.DisplayName
					it.CompanyDisplayName = cu.DisplayName
				}
			}
			// enrich vacancy title (optional)
			if v, _ := vac.GetByID(c.Request.Context(), a.VacancyID); v != nil {
				it.VacancyTitle = v.Title
			}

			out = append(out, it)
		}

		c.JSON(200, gin.H{"ok": true, "items": out})
	})

	// новый инбокс соотвктвующий inboxItem
	// company inbox
	protected.GET("/applications/inbox", middleware.RequireType("company"), func(c *gin.Context) {
		companyID := c.GetString(middleware.CtxUserID)
		status := c.Query("status")

		items, err := apps.ListInbox(c.Request.Context(), companyID, status, 50, 0)
		if err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}

		vTitle := map[string]string{}
		rTitle := map[string]string{}
		uName := map[string]string{}
		uPremium := map[string]bool{}

		out := make([]inboxItem, 0, len(items))

		for _, a := range items {
			it := inboxItem{
				ApplicationID: a.ApplicationID,
				VacancyID:     a.VacancyID,
				ResumeID:      a.ResumeID,
				UserID:        a.UserID,
				CompanyID:     a.CompanyID,
				Status:        a.Status,
				Message:       a.Message,
				Viewed:        a.Viewed,
			}

			// vacancy title
			if t, ok := vTitle[a.VacancyID]; ok {
				it.VacancyTitle = t
			} else {
				if v, _ := vac.GetByID(c.Request.Context(), a.VacancyID); v != nil {
					vTitle[a.VacancyID] = v.Title
					it.VacancyTitle = v.Title
				}
			}

			// resume title
			if t, ok := rTitle[a.ResumeID]; ok {
				it.ResumeTitle = t
			} else {
				if rr, _ := resumes.GetByID(c.Request.Context(), a.ResumeID); rr != nil {
					rTitle[a.ResumeID] = rr.Title
					it.ResumeTitle = rr.Title
				}
			}

			// user displayName and premium status
			if n, ok := uName[a.UserID]; ok {
				it.UserDisplayName = n
				it.UserIsPremium = uPremium[a.UserID]
			} else {
				if u, _ := users.FindByUserID(c.Request.Context(), a.UserID); u != nil {
					uName[a.UserID] = u.DisplayName
					it.UserDisplayName = u.DisplayName
					uPremium[a.UserID] = u.Subscription.Active
					it.UserIsPremium = u.Subscription.Active
				}
			}

			out = append(out, it)
		}

		// ✅ ВАЖНО: отдаём out, НЕ items
		c.JSON(200, gin.H{"ok": true, "items": out})
	})

	protected.POST("/applications/:id/accept", middleware.RequireType("company"), func(c *gin.Context) {
		companyID := c.GetString(middleware.CtxUserID)
		if err := apps.UpdateStatus(c.Request.Context(), c.Param("id"), companyID, "accepted"); err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		c.JSON(200, gin.H{"ok": true, "status": "accepted"})
	})

	protected.POST("/applications/:id/reject", middleware.RequireType("company"), func(c *gin.Context) {
		companyID := c.GetString(middleware.CtxUserID)
		if err := apps.UpdateStatus(c.Request.Context(), c.Param("id"), companyID, "rejected"); err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		c.JSON(200, gin.H{"ok": true, "status": "rejected"})
	})
	protected.POST("/applications/:id/viewed", middleware.RequireType("company"), func(c *gin.Context) {
		companyID := c.GetString(middleware.CtxUserID)

		if err := apps.MarkViewed(c.Request.Context(), c.Param("id"), companyID); err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}

		c.JSON(200, gin.H{"ok": true, "viewed": true})
	})

	// Скрыть отклик для пользователя или компании
	protected.POST("/applications/:id/hide", func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		userType := c.GetString(middleware.CtxUserType)

		isCompany := userType == "company"

		if err := apps.Hide(c.Request.Context(), c.Param("id"), userID, isCompany); err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}

		c.JSON(200, gin.H{"ok": true})
	})
}
