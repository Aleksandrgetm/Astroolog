package web

import (
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPublicRouting(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, "routes.json"), []byte(`{"paths":["/","/lv/services"],"site":"https://example.invalid"}`), 0600)
	os.WriteFile(filepath.Join(root, "index.html"), []byte("HOME"), 0600)
	os.WriteFile(filepath.Join(root, "404.html"), []byte("MISSING"), 0600)
	os.MkdirAll(filepath.Join(root, "lv", "services"), 0700)
	os.WriteFile(filepath.Join(root, "lv", "services", "index.html"), []byte("LV SERVICES"), 0600)
	os.MkdirAll(filepath.Join(root, "lv"), 0700)
	os.WriteFile(filepath.Join(root, "lv", "404.html"), []byte("LV MISSING"), 0600)
	r := gin.New()
	if err := Register(r, root, nil); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		path     string
		status   int
		location string
	}{{"//example.invalid/", 301, "/example.invalid"}, {"/", 200, ""}, {"/missing", 404, ""}, {"/lv/missing", 404, ""}, {"/services?lang=lv&option=personality", 301, "/lv/services?option=personality"}, {"/en/services/growth-point", 301, "/en/services"}, {"/lv/services/", 301, "/lv/services"}, {"/api/testimonials", 404, ""}, {"/routes.json", 404, ""}} {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest("GET", tc.path, nil))
		if w.Code != tc.status {
			t.Errorf("%s: %d", tc.path, w.Code)
		}
		if tc.location != "" && w.Header().Get("Location") != tc.location {
			t.Errorf("%s: redirect %s", tc.path, w.Header().Get("Location"))
		}
	}
}

// Exercise the production asset allowlist with the real generated icons.
func TestBrandIcons(t *testing.T) {
	root := t.TempDir()
	t.Setenv("CMS_RELEASE_ROOT", t.TempDir())
	os.WriteFile(filepath.Join(root, "routes.json"), []byte(`{"paths":[]}`), 0600)
	icons := map[string]string{"favicon.svg": "image/svg+xml", "favicon.ico": "image/", "favicon-32x32.png": "image/png", "favicon-192x192.png": "image/png", "apple-touch-icon.png": "image/png"}
	for name := range icons {
		data, err := os.ReadFile(filepath.Join("../../../frontend/public", name))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, name), data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	r := gin.New()
	if err := Register(r, root, nil); err != nil {
		t.Fatal(err)
	}
	for name, contentType := range icons {
		for _, method := range []string{"GET", "HEAD"} {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(method, "/"+name+"?v=ez2", nil))
			if w.Code != 200 || !strings.HasPrefix(w.Header().Get("Content-Type"), contentType) {
				t.Errorf("%s %s: %d %s", method, name, w.Code, w.Header().Get("Content-Type"))
			}
		}
	}
}
