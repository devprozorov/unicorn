package middleware

import (
	"net/http"

	"unicorn-auth/internal/repo"
	"unicorn-auth/internal/security"

	"github.com/gin-gonic/gin"
)

func RequireMFAEnabled(_ *security.Security, users *repo.UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetString(CtxUserID)
		u, err := users.FindByUserID(c.Request.Context(), uid)
		if err != nil || u == nil || u.Status.Blocked || u.Status.Deleted {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"ok": false, "error": "unauthorized"})
			return
		}
		if !u.MFA.TOTP.Enabled {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"ok": false, "error": "mfa_required"})
			return
		}
		c.Next()
	}
}
