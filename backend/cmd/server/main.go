package main

import (
	"astroolog/backend/internal/routes"
	"astroolog/backend/internal/web"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"astroolog/backend/internal/config"
	"astroolog/backend/internal/database"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db := database.Connect(cfg)

	if err := database.Migrate(db); err != nil {
		panic(err)
	}
	router := gin.Default()
	_ = router.SetTrustedProxies(nil)
	routes.Register(router, db, cfg)

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

	if directory := os.Getenv("SITE_DIST_DIR"); directory != "" {
		if err := web.Register(router, directory, db); err != nil {
			log.Fatal(err)
		}
	}

	address := fmt.Sprintf(":%s", cfg.AppPort)

	server := &http.Server{Addr: address, Handler: router, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
