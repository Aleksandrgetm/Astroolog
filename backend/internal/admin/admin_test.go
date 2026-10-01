package admin

import (
	"astroolog/backend/internal/config"
	"astroolog/backend/internal/database"
	"astroolog/backend/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gen2brain/avif"
	"github.com/gen2brain/webp"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMediaValidation(t *testing.T) {
	for _, b := range [][]byte{[]byte("<svg><script>alert(1)</script></svg>"), []byte("<html>bad</html>"), []byte("not an image")} {
		if _, e := decodeImage(b); e == nil {
			t.Fatal("unsafe file accepted")
		}
	}
	var b bytes.Buffer
	png.Encode(&b, image.NewNRGBA(image.Rect(0, 0, 20, 10)))
	img, e := decodeImage(b.Bytes())
	if e != nil || img.Bounds().Dx() != 20 {
		t.Fatal(e)
	}

	for _, format := range []string{"jpeg", "webp", "avif"} {
		var encoded bytes.Buffer
		sample := image.NewNRGBA(image.Rect(0, 0, 20, 10))
		var err error
		switch format {
		case "jpeg":
			err = jpeg.Encode(&encoded, sample, nil)
		case "webp":
			err = webp.Encode(&encoded, sample)
		case "avif":
			err = avif.Encode(&encoded, sample, avif.Options{Speed: 10, Quality: 70})
		}
		if err != nil {
			t.Fatal(format, err)
		}
		decoded, err := decodeImage(encoded.Bytes())
		if err != nil || decoded.Bounds().Dx() != 20 {
			t.Fatal(format, err)
		}
	}
	if safeImage("/media/../../test") {
		t.Fatal("path traversal")
	}
}
func TestAdminIntegration(t *testing.T) {
	name := os.Getenv("TEST_DB_NAME")
	if name == "" {
		t.Skip("isolated TEST_DB_NAME required")
	}
	if !strings.HasSuffix(name, "_test") {
		t.Fatal("isolated database only")
	}
	cfg, e := config.Load()
	if e != nil {
		t.Fatal(e)
	}
	cfg.DBName = name
	db := database.Connect(cfg)
	if e = database.Migrate(db); e != nil {
		t.Fatal(e)
	}
	t.Setenv("APP_ORIGIN", "http://example.test")
	t.Setenv("APP_ENV", "development")
	digest, _ := bcrypt.GenerateFromPassword([]byte("cms-test-only-password-2026"), 12)
	user := models.User{Email: "cms-test@example.invalid", PasswordHash: string(digest), Role: "admin", IsActive: true}
	db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "email"}}, DoUpdates: clause.Assignments(map[string]any{"password_hash": string(digest), "role": "admin", "is_active": true})}).Create(&user)
	db.First(&user, "email=?", user.Email)
	db.Where("key IN ?", []string{hash("ip:192.0.2.1"), hash("email:" + user.Email)}).Delete(&models.LoginAttempt{})
	gin.SetMode(gin.TestMode)
	r := gin.New()
	t.Setenv("MEDIA_ROOT", t.TempDir())
	cfg.AdminAllowedOrigins = []string{"http://example.test"}
	Register(r, db, cfg)
	request := func(method, path string, data any, cookie *http.Cookie, csrf string) *httptest.ResponseRecorder {
		var body bytes.Buffer
		if data != nil {
			json.NewEncoder(&body).Encode(data)
		}
		req := httptest.NewRequest(method, path, &body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Origin", "http://example.test")
		if cookie != nil {
			req.AddCookie(cookie)
		}
		if csrf != "" {
			req.Header.Set("X-CSRF-Token", csrf)
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}
	if w := request("GET", "/api/admin/dashboard", nil, nil, ""); w.Code != 401 {
		t.Fatal("unauthenticated access", w.Code)
	}
	login := request("POST", "/api/admin/auth/login", map[string]string{"email": user.Email, "password": "cms-test-only-password-2026"}, nil, "")
	if login.Code != 200 {
		t.Fatal(login.Code, login.Body.String())
	}
	cookie := login.Result().Cookies()[0]
	if !cookie.HttpOnly || cookie.SameSite != http.SameSiteStrictMode || cookie.Path != "/" {
		t.Fatal("unsafe cookie")
	}
	var response map[string]any
	json.Unmarshal(login.Body.Bytes(), &response)
	csrf := response["csrf"].(string)
	if w := request("POST", "/api/admin/publish", nil, cookie, ""); w.Code != 403 {
		t.Fatal("CSRF bypass", w.Code)
	}
	if w := request("POST", "/api/admin/publish", nil, cookie, "wrong"); w.Code != 403 {
		t.Fatal("CSRF bypass", w.Code)
	}
	if w := request("GET", "/api/admin/dashboard", nil, cookie, ""); w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	list := request("GET", "/api/admin/entities/sections?limit=100", nil, cookie, "")
	json.Unmarshal(list.Body.Bytes(), &response)
	items := response["items"].([]any)
	var section map[string]any
	for _, item := range items {
		m := item.(map[string]any)
		if m["section_key"] == "home.numerology" {
			section = m
		}
	}
	if section == nil {
		t.Fatal("missing seeded section")
	}
	translations := section["translations"].([]any)
	original := translations[0].(map[string]any)["fields"].(map[string]any)["home.aToolForSelfreflectionNotScientificAssessment"]
	translations[0].(map[string]any)["fields"].(map[string]any)["home.aToolForSelfreflectionNotScientificAssessment"] = "CMS integration text"
	path := fmt.Sprintf("/api/admin/entities/sections/%.0f", section["id"])
	saved := request("PUT", path, section, cookie, csrf)
	if saved.Code != 200 {
		t.Fatal(saved.Code, saved.Body.String())
	}
	if w := request("PUT", path, section, cookie, csrf); w.Code != 400 {
		t.Fatal("concurrent edit not rejected")
	}
	json.Unmarshal(saved.Body.Bytes(), &section)
	section["translations"].([]any)[0].(map[string]any)["fields"].(map[string]any)["home.aToolForSelfreflectionNotScientificAssessment"] = original
	if w := request("PUT", path, section, cookie, csrf); w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	// Media is re-encoded, referenced files cannot be deleted, and originals remain immutable.
	var uploadBody bytes.Buffer
	form := multipart.NewWriter(&uploadBody)
	part, _ := form.CreateFormFile("file", "photo.png")
	png.Encode(part, image.NewNRGBA(image.Rect(0, 0, 32, 16)))
	form.Close()
	req := httptest.NewRequest("POST", "/api/admin/media", &uploadBody)
	req.Header.Set("Content-Type", form.FormDataContentType())
	req.Header.Set("Origin", "http://example.test")
	req.Header.Set("X-CSRF-Token", csrf)
	req.AddCookie(cookie)
	uploaded := httptest.NewRecorder()
	r.ServeHTTP(uploaded, req)
	if uploaded.Code != 201 {
		t.Fatal("upload", uploaded.Code, uploaded.Body.String())
	}
	var media models.MediaAsset
	json.Unmarshal(uploaded.Body.Bytes(), &media)
	if media.MimeType != "image/webp" || media.Width != 32 || strings.Contains(media.StoragePath, "photo") {
		t.Fatal("unsafe media metadata")
	}
	if w := request("GET", media.StoragePath, nil, nil, ""); w.Code != 200 || w.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatal("media unavailable")
	}
	var service models.Service
	db.First(&service)
	originalImage := service.Image
	db.Model(&service).Update("image", media.StoragePath)
	mediaPath := fmt.Sprintf("/api/admin/media/%d", media.ID)
	if w := request("DELETE", mediaPath, nil, cookie, csrf); w.Code != 400 {
		t.Fatal("referenced image deleted")
	}
	db.Model(&service).Update("image", originalImage)
	if w := request("DELETE", mediaPath, nil, cookie, csrf); w.Code != 204 {
		t.Fatal("unused image deletion", w.Code, w.Body.String())
	}
	var booking models.Booking
	if db.First(&booking).Error == nil {
		before, _ := json.Marshal(booking)
		w := request("PATCH", fmt.Sprintf("/api/admin/requests/bookings/%d", booking.ID), map[string]any{"status": "in_progress", "internal_note": "integration note", "message": "must not overwrite original", "price": 1}, cookie, csrf)
		if w.Code != 204 {
			t.Fatal("request update", w.Code)
		}
		var after models.Booking
		db.First(&after, booking.ID)
		if after.Message != booking.Message || after.Price != nil && booking.Price != nil && *after.Price != *booking.Price || after.OptionTitle != booking.OptionTitle {
			t.Fatal("original snapshot changed", string(before))
		}
		db.Model(&booking).Updates(map[string]any{"workflow_status": booking.WorkflowStatus, "internal_note": booking.InternalNote, "updated_at": booking.UpdatedAt})
	}
	// No front-end guard can bypass server roles or activity checks.
	db.Model(&user).Update("role", "client")
	if w := request("GET", "/api/admin/dashboard", nil, cookie, ""); w.Code != 403 {
		t.Fatal("RBAC bypass")
	}
	db.Model(&user).Update("role", "admin")
	db.Model(&models.AdminSession{}).Where("id=?", hash(cookie.Value)).Update("expires_at", time.Now().Add(-time.Minute))
	if w := request("GET", "/api/admin/auth/me", nil, cookie, ""); w.Code != 401 {
		t.Fatal("expired session accepted")
	}
	login = request("POST", "/api/admin/auth/login", map[string]string{"email": user.Email, "password": "cms-test-only-password-2026"}, nil, "")
	json.Unmarshal(login.Body.Bytes(), &response)
	cookie = login.Result().Cookies()[0]
	csrf = response["csrf"].(string)
	if w := request("POST", "/api/admin/auth/logout", nil, cookie, csrf); w.Code != 204 {
		t.Fatal("logout failed")
	}
	if w := request("GET", "/api/admin/auth/me", nil, cookie, ""); w.Code != 401 {
		t.Fatal("logout did not invalidate")
	}
	for i := 0; i < 9; i++ {
		w := request("POST", "/api/admin/auth/login", map[string]string{"email": "missing@example.invalid", "password": "wrong"}, nil, "")
		if i == 8 && w.Code != 429 {
			t.Fatal("throttling failed", w.Code)
		}
	}
	var n int64
	db.Model(&models.Review{}).Count(&n)
	if n != 8 {
		t.Fatal("reviews not migrated", n)
	}
	for _, review := range database.CMSDefinitions().Reviews {
		for _, tr := range review.Translations {
			if tr.FullText == "" {
				t.Fatal("lost original review text")
			}
		}
	}
}

func TestSEOPageAcceptsExistingReadOnlyRoute(t *testing.T) {
	for _, path := range []string{"/", "/about", "/reviews"} {
		page := &models.Page{Key: "home", Slug: path, Translations: []models.PageTranslation{{Language: "ru", Title: "Главная", MetaTitle: "Главная страница", MetaDescription: "Описание страницы"}}}
		if err := (Handler{}).validate(nil, "pages", page, map[string]any{"slug": path}); err != nil {
			t.Fatalf("read-only route %s rejected: %v", path, err)
		}
	}
}
