package admin

import (
	"astroolog/backend/internal/models"
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (h Handler) Publish(c *gin.Context) {
	directory := os.Getenv("CMS_FRONTEND_DIR")
	if directory == "" {
		directory = "../frontend"
	}
	directory, _ = filepath.Abs(directory)
	site := os.Getenv("VITE_SITE_URL")
	if site == "" {
		site = os.Getenv("APP_ORIGIN")
	}
	if !strings.HasPrefix(site, "https://") {
		c.JSON(503, gin.H{"error": "Для публикации настройте HTTPS-домен VITE_SITE_URL на сервере."})
		return
	}
	publication := models.Publication{UserID: uid(c), Status: "running"}
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if e := tx.Exec("SELECT pg_advisory_xact_lock(20260927)").Error; e != nil {
			return e
		}
		var n int64
		if e := tx.Model(&models.Publication{}).Where("status='running' AND created_at>?", time.Now().Add(-15*time.Minute)).Count(&n).Error; e != nil {
			return e
		}
		if n > 0 {
			return gorm.ErrInvalidTransaction
		}
		if e := tx.Model(&models.Publication{}).Where("status='running' AND created_at<=?", time.Now().Add(-15*time.Minute)).Updates(map[string]any{"status": "failed", "error": "Процесс публикации был прерван. Повторите публикацию."}).Error; e != nil {
			return e
		}
		if e := tx.Create(&publication).Error; e != nil {
			return e
		}
		return audit(tx, uid(c), "publish", "site", strconv.Itoa(int(publication.ID)))
	})
	if err != nil {
		c.JSON(409, gin.H{"error": "Публикация уже выполняется. Обновите Dashboard позже."})
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		releases := os.Getenv("CMS_RELEASE_ROOT")
		if releases == "" {
			releases = filepath.Join(directory, ".releases")
		}
		releases, _ = filepath.Abs(releases)
		target := filepath.Join(releases, strconv.Itoa(int(publication.ID)))
		cmd := exec.CommandContext(ctx, "npm", "run", "build")
		cmd.Dir = directory
		cmd.Env = append(os.Environ(), "VITE_SITE_URL="+site, "SITE_BUILD_DIR="+target)
		// Build output is deliberately not exposed through the admin API; it may include environment diagnostics.
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err == nil {
			link := filepath.Join(releases, ".next-"+strconv.Itoa(int(publication.ID)))
			err = os.Symlink(target, link)
			if err == nil {
				err = os.Rename(link, filepath.Join(releases, "current"))
			}
		}
		status, message := "published", ""
		if err != nil {
			log.Printf("CMS publication %d failed: %v", publication.ID, err)
			status = "failed"
			message = "Публикация не завершена. Текущая версия сайта сохранена; проверьте конфигурацию и доступность API."
		}
		h.DB.Model(&publication).Updates(map[string]any{"status": status, "error": message, "finished_at": time.Now()})
	}()
	c.JSON(202, publication)
}
