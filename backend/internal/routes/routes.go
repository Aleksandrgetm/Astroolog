package routes

import (
	"astroolog/backend/internal/admin"
	"astroolog/backend/internal/config"
	"astroolog/backend/internal/handlers"
	"astroolog/backend/internal/middleware"
	"astroolog/backend/internal/repositories"
	"astroolog/backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(r *gin.Engine, db *gorm.DB, cfg config.Config) {
	admin.Register(r, db, cfg)
	repo := repositories.Repository{DB: db}
	h := handlers.Handler{Repo: repo, Requests: services.Requests{Repo: repo}}
	api := r.Group("/api")
	api.GET("/booking-options", h.BookingOptions)
	api.GET("/services", h.ListServices)
	api.GET("/services/:slug", h.Service)
	api.POST("/bookings", middleware.BodyLimit(), func(c *gin.Context) { h.Submit(c, true) })
	api.POST("/contact", middleware.BodyLimit(), func(c *gin.Context) { h.Submit(c, false) })
}
