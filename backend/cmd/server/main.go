package main

import (
	"astroolog/backend/internal/routes"
	"fmt"
	"net/http"
	"os"

	"astroolog/backend/internal/config"
	"astroolog/backend/internal/database"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db := database.Connect(cfg)

	if err := database.Migrate(db); err != nil {
		panic(err)
	}
	if os.Getenv("SEED_DEMO") == "true" {
		if err := database.SeedDemo(db); err != nil {
			panic(err)
		}
	}
	router := gin.Default()
	_ = router.SetTrustedProxies(nil)
	routes.Register(router, db)

	router.GET("/api/health", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":   "error",
				"database": "unavailable",
			})
			return
		}

		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":   "error",
				"database": "unavailable",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"service":  "astroolog-api",
			"database": "connected",
		})
	})

	address := fmt.Sprintf(":%s", cfg.AppPort)

	if err := router.Run(address); err != nil {
		panic(err)
	}
}
