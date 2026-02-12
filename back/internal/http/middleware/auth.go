package middleware

import (
	"net/http"
	"strings"

	"unicorn-auth/internal/security"

	"github.com/gin-gonic/gin"
)

const (
	CtxUserID   = "userId"
	CtxUserType = "userType"
	CtxAMR      = "amr"
)

func RequireAuth(sec *security.Security) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"ok": false, "error": "unauthorized"})
			return
		}
		token := strings.TrimSpace(strings.TrimPrefix(h, "Bearer "))
		claims, err := sec.Tokens.ParseAccess(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"ok": false, "error": "unauthorized"})
			return
		}
		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxUserType, claims.Type)
		c.Set(CtxAMR, claims.AMR)
		c.Next()
	}
}

func RequireType(want string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString(CtxUserType) != want {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"ok": false, "error": "forbidden"})
			return
		}
		c.Next()
	}
}
