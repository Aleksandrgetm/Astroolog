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
		// A completed publication switches atomically; retain the previous release during builds.
		activeRoot := root
		releases := os.Getenv("CMS_RELEASE_ROOT")
		if releases == "" {
			releases = filepath.Join(filepath.Dir(root), ".releases")
		}
		if _, e := os.Stat(filepath.Join(releases, "current", "routes.json")); e == nil {
			activeRoot = filepath.Join(releases, "current")
		}
		currentKnown := known
		if activeRoot != root {
			if b, e := os.ReadFile(filepath.Join(activeRoot, "routes.json")); e == nil {
				var latest Manifest
				if json.Unmarshal(b, &latest) == nil {
					currentKnown = map[string]bool{}
					for _, p := range latest.Paths {
						currentKnown[p] = true
					}
				}
			}
		}
		root := activeRoot
		p := c.Request.URL.Path
		if p == "/admin" || strings.HasPrefix(p, "/admin/") {
			c.Header("X-Robots-Tag", "noindex, nofollow")
			c.Header("Cache-Control", "no-store")
			c.Header("X-Frame-Options", "DENY")
			c.File(filepath.Join(root, "admin.html"))
			return
		}
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
		allowed := currentKnown[p]
		if allowed && strings.HasPrefix(base, "/services/") {
			var count int64
			if err := db.Model(&models.Service{}).Where("slug=? AND is_active=true", strings.TrimPrefix(base, "/services/")).Count(&count).Error; err != nil {
				c.Status(503)
				return
			}
			allowed = count == 1
		}
		if allowed && strings.HasPrefix(base, "/directions/") {
			var count int64
			if err := db.Model(&models.Direction{}).Where("slug=? AND is_active=true", strings.TrimPrefix(base, "/directions/")).Count(&count).Error; err != nil {
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
		// Existing open pages may still request a content-hashed chunk from an older release.
		if strings.HasPrefix(p, "/assets/") && !strings.Contains(p, "/.") && filepath.Base(p) == strings.TrimPrefix(p, "/assets/") {
			candidates := []string{filepath.Join(directory, p)}
			entries, _ := os.ReadDir(releases)
			for _, entry := range entries {
				if entry.IsDir() {
					candidates = append(candidates, filepath.Join(releases, entry.Name(), p))
				}
			}
			for _, candidate := range candidates {
				if info, e := os.Stat(candidate); e == nil && !info.IsDir() {
					c.Header("Cache-Control", "public,max-age=31536000,immutable")
					c.File(candidate)
					return
				}
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
