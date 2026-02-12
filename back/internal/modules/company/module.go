package company

import (
	"net/http"
	"regexp"
	"strings"

	"unicorn-auth/internal/models"
	"unicorn-auth/internal/repo"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func Register(r *gin.Engine, profiles *repo.ProfileRepo) {
	api := r.Group("/api")

	api.GET("/companies/search", func(c *gin.Context) {
		q := strings.TrimSpace(c.Query("q"))
		location := strings.TrimSpace(c.Query("location"))
		industry := strings.TrimSpace(c.Query("industry"))

		filter := bson.M{"type": models.UserTypeCompany}

		if q != "" {
			if len(q) > 64 {
				c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "bad_request"})
				return
			}
			safe := regexp.QuoteMeta(q)
			filter["displayName"] = bson.M{"$regex": safe, "$options": "i"}
		}
		if location != "" {
			filter["location"] = location
		}
		if industry != "" {
			filter["industry"] = industry
		}

		items, err := profiles.SearchCompanies(c.Request.Context(), filter, 50)
		if err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		c.JSON(200, gin.H{"ok": true, "items": items})
	})
}
