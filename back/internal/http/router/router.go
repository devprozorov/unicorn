package router

import (
	"net/http"

	"unicorn-auth/internal/config"
	"unicorn-auth/internal/http/handlers"
	"unicorn-auth/internal/http/middleware"
	"unicorn-auth/internal/repo"
	"unicorn-auth/internal/security"

	"github.com/gin-gonic/gin"
)

func New(cfg config.Config, sec *security.Security, users *repo.UserRepo, sessions *repo.SessionRepo, resumes *repo.ResumeRepo, vacancies *repo.VacancyRepo) *gin.Engine {
	if cfg.AppEnv == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	//AVAupload
	r.MaxMultipartMemory = 8 << 20 // 8MB (можешь поменять) ограничение на аватарку
	r.Static("/uploads", "./uploads")

	// simple health
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	// CORS (simple)
	r.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		for _, allow := range cfg.CorsOrigins {
			if origin == allow {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Vary", "Origin")
				c.Header("Access-Control-Allow-Credentials", "true")
				c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
				c.Header("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE,OPTIONS")
				break
			}
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	rl := middleware.NewRateLimiter(5, 10)
	r.Use(rl.Middleware())

	ah := handlers.NewAuthHandler(cfg, sec, users, sessions)
	hh := handlers.NewHomeHandler(users, resumes, vacancies)

	api := r.Group("/api")
	{
		// Public endpoints
		api.GET("/home/stats", hh.GetStats)

		api.POST("/auth/register", ah.Register)
		api.POST("/auth/login", ah.Login)
		api.POST("/auth/totp/verify", ah.VerifyTOTP)
		api.POST("/auth/refresh", ah.Refresh)
		api.POST("/auth/logout", ah.Logout)

		protected := api.Group("")
		protected.Use(middleware.RequireAuth(sec))
		{
			protected.GET("/auth/me", ah.Me)
			protected.POST("/auth/change-password", ah.ChangePassword)
			protected.POST("/auth/totp/enroll", ah.TotpEnroll)
			protected.POST("/auth/totp/enable", ah.TotpEnable)
		}
	}

	return r
}
