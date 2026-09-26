package web

import (
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"os"
	"path/filepath"
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
