package admin

import (
	"astroolog/backend/internal/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"os"
)

func Register(r *gin.Engine, db *gorm.DB, cfg config.Config) {
	root := os.Getenv("MEDIA_ROOT")
	if root == "" {
		root = "./storage/media"
	}
	h := Handler{DB: db, AllowedOrigins: cfg.AdminAllowedOrigins, Store: LocalStorage{Root: root}}
	r.GET("/api/content", h.Public)
	r.GET("/media/:name", h.ServeMedia)
	group := r.Group("/api/admin")
	group.Use(func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("X-Robots-Tag", "noindex, nofollow")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Next()
	})
	group.POST("/auth/login", h.Login)
	group.Use(h.RequireAuth(), RequireRole("admin"), h.CSRF())
	group.GET("/auth/me", h.Me)
	group.POST("/auth/logout", h.Logout)
	group.GET("/dashboard", h.Dashboard)
	group.GET("/schema", h.Schema)
	group.GET("/entities/:kind", h.List)
	group.GET("/entities/:kind/:id", h.Get)
	group.POST("/entities/:kind", h.Save)
	group.PUT("/entities/:kind/:id", h.Save)
	group.GET("/requests", h.Requests)
	group.PATCH("/requests/:kind/:id", h.RequestUpdate)
	group.GET("/audit", h.Audit)
	group.GET("/media", h.Media)
	group.POST("/media", h.Upload)
	group.DELETE("/media/:id", h.DeleteMedia)
	group.POST("/publish", h.Publish)
}
