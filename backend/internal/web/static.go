// Package web serves only explicitly prerendered public routes; unknown routes are real 404s.
package web

import (
	"astroolog/backend/internal/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"os"
	urlpath "path"
	"path/filepath"
	"strings"
)

type Manifest struct {
	Site  string   `json:"site"`
	Paths []string `json:"paths"`
}

func Register(r *gin.Engine, directory string, db *gorm.DB) error {
	root, err := filepath.Abs(directory)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(filepath.Join(root, "routes.json"))
	if err != nil {
		return fmt.Errorf("prerender manifest: %w", err)
	}
	var manifest Manifest
	if err = json.Unmarshal(data, &manifest); err != nil {
		return err
	}
	known := map[string]bool{}
	for _, p := range manifest.Paths {
		known[p] = true
	}
	r.NoRoute(func(c *gin.Context) {
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			c.Status(404)
			return
		}
		p := c.Request.URL.Path
		if strings.HasPrefix(p, "/api/") {
			c.JSON(404, gin.H{"code": "not_found"})
			return
		}
		locale := ""
		base := p
		for _, lang := range []string{"lv", "en"} {
			if p == "/"+lang || strings.HasPrefix(p, "/"+lang+"/") {
				locale = "/" + lang
				base = strings.TrimPrefix(p, locale)
				break
			}
		}
		query := c.Request.URL.Query()
		if values, ok := query["lang"]; ok {
			locale = ""
			if len(values) > 0 && (values[0] == "lv" || values[0] == "en") {
				locale = "/" + values[0]
			}
			query.Del("lang")
		}
		base = urlpath.Clean("/" + strings.TrimLeft(base, "/"))
		if base == "" {
			base = "/"
		}
		if base == "/services/growth-point" {
			base = "/services"
		}
		normalized := locale + base
		if normalized != p || c.Request.URL.Query().Has("lang") {
			if q := query.Encode(); q != "" {
				normalized += "?" + q
			}
			c.Redirect(301, normalized)
			return
		}
		allowed := known[p]
		if allowed && strings.HasPrefix(base, "/services/") {
			var count int64
			if err := db.Model(&models.Service{}).Where("slug=? AND is_active=true", strings.TrimPrefix(base, "/services/")).Count(&count).Error; err != nil {
				c.Status(503)
				return
			}
			allowed = count == 1
		}
		c.Header("X-Content-Type-Options", "nosniff")
		if allowed {
			c.Header("Cache-Control", "no-cache")
			c.File(filepath.Join(root, strings.TrimPrefix(p, "/"), "index.html"))
			return
		}
		// Only asset trees and generated public discovery files are served directly.
		asset := strings.HasPrefix(p, "/assets/") || strings.HasPrefix(p, "/images/") || strings.HasPrefix(p, "/fonts/") || p == "/robots.txt" || p == "/sitemap.xml" || p == "/favicon.ico"
		clean := filepath.Clean(filepath.Join(root, p))
		if asset && strings.HasPrefix(clean, root+string(os.PathSeparator)) && !strings.Contains(p, "/.") {
			if info, e := os.Stat(clean); e == nil && !info.IsDir() {
				c.File(clean)
				return
			}
		}
		c.Header("X-Robots-Tag", "noindex, follow")
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Status(404)
		body, e := os.ReadFile(filepath.Join(root, strings.TrimPrefix(locale, "/"), "404.html"))
		if e == nil {
			_, _ = c.Writer.Write(body)
		}
	})
	return nil
}
