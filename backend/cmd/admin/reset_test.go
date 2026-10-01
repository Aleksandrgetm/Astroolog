package main

import (
	"astroolog/backend/internal/admin"
	"astroolog/backend/internal/config"
	"astroolog/backend/internal/database"
	"astroolog/backend/internal/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestPasswordRequirements(t *testing.T) {
	for _, tc := range []struct {
		name, password, repeat string
		valid                  bool
	}{
		{"minimum", strings.Repeat("a", 12), strings.Repeat("a", 12), true},
		{"maximum", strings.Repeat("a", 72), strings.Repeat("a", 72), true},
		{"short", "short", "short", false},
		{"long", strings.Repeat("a", 73), strings.Repeat("a", 73), false},
		{"byte limit", strings.Repeat("я", 37), strings.Repeat("я", 37), false},
		{"mismatch", "new-password-123", "other-password-123", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			digest, err := passwordDigest([]byte(tc.password), []byte(tc.repeat))
			if (err == nil) != tc.valid {
				t.Fatalf("unexpected validation result: %v", err)
			}
			if tc.valid {
				cost, _ := bcrypt.Cost(digest)
				if cost != 12 || bcrypt.CompareHashAndPassword(digest, []byte(tc.password)) != nil {
					t.Fatal("create-compatible bcrypt required")
				}
			}
		})
	}
}
func TestResetPasswordIntegration(t *testing.T) {
	name := os.Getenv("TEST_DB_NAME")
	if name == "" {
		t.Skip("isolated TEST_DB_NAME required")
	}
	if !strings.HasSuffix(name, "_test") {
		t.Fatal("isolated database only")
	}
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	cfg.DBName = name
	cfg.AdminAllowedOrigins = []string{"http://localhost:5174"}
	db := database.Connect(cfg)
	if err := database.Migrate(db); err != nil {
		t.Fatal(err)
	}
	unique := fmt.Sprint(time.Now().UnixNano())
	oldPassword := "old-test-password-123"
	newPassword := "new-test-password-456"
	digest, _ := bcrypt.GenerateFromPassword([]byte(oldPassword), 12)
	lastLogin := time.Now().UTC().Truncate(time.Second)
	user := models.User{Email: "reset-" + unique + "@example.invalid", Role: "admin", IsActive: true, PasswordHash: string(digest), LastLoginAt: &lastLogin}
	if err := db.Create(&user).Error; err != nil {
		t.Fatal(err)
	}
	defer func() {
		db.Where("user_id = ?", user.ID).Delete(&models.AdminSession{})
		db.Where("user_id = ?", user.ID).Delete(&models.AdminAuditLog{})
		db.Delete(&user)
	}()
	oldUpdated := user.UpdatedAt
	session := models.AdminSession{ID: "reset-test-" + unique, UserID: user.ID, ExpiresAt: time.Now().Add(time.Hour)}
	if err := db.Create(&session).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := findAdmin(db, "unknown-"+unique+"@example.invalid"); err != errAdminNotFound {
		t.Fatal("unknown email accepted", err)
	}
	for _, pair := range [][2]string{{"short", "short"}, {newPassword, "mismatch-password"}} {
		if err := resetPassword(db, user, []byte(pair[0]), []byte(pair[1])); err == nil {
			t.Fatal("invalid password accepted")
		}
	}
	var unchanged models.User
	db.First(&unchanged, user.ID)
	if unchanged.PasswordHash != user.PasswordHash {
		t.Fatal("failed validation changed password")
	}
	var output bytes.Buffer
	calls := 0
	if err := resetInteractive(db, strings.NewReader("  "+strings.ToUpper(user.Email)+"  \n"), &output, true, func() ([]byte, error) { calls++; return []byte(newPassword), nil }); err != nil {
		t.Fatal(err)
	}
	if calls != 2 || !strings.Contains(output.String(), "Password updated successfully.") {
		t.Fatal("interactive confirmation failed")
	}
	var updated models.User
	db.First(&updated, user.ID)
	if updated.PasswordHash == user.PasswordHash || bcrypt.CompareHashAndPassword([]byte(updated.PasswordHash), []byte(newPassword)) != nil || bcrypt.CompareHashAndPassword([]byte(updated.PasswordHash), []byte(oldPassword)) == nil {
		t.Fatal("password not replaced")
	}
	if updated.Email != user.Email || updated.Role != user.Role || updated.IsActive != user.IsActive || !updated.CreatedAt.Equal(user.CreatedAt) || !updated.LastLoginAt.Equal(lastLogin) || !updated.UpdatedAt.After(oldUpdated) {
		t.Fatal("unexpected account changes")
	}
	if strings.Contains(output.String(), newPassword) || strings.Contains(output.String(), updated.PasswordHash) {
		t.Fatal("credential leaked to console")
	}
	var n int64
	db.Model(&models.AdminSession{}).Where("user_id = ?", user.ID).Count(&n)
	if n != 0 {
		t.Fatal("sessions survived reset")
	}
	var logs []models.AdminAuditLog
	db.Where("user_id = ? AND action = ?", user.ID, "password_reset").Find(&logs)
	if len(logs) != 1 || logs[0].Metadata["source"] != "local_cli" || len(logs[0].Metadata) != 1 {
		t.Fatal("invalid audit")
	}
	if err := resetInteractive(db, strings.NewReader(user.Email+"\n"), &output, false, func() ([]byte, error) { t.Fatal("password read without TTY"); return nil, nil }); err == nil {
		t.Fatal("noninteractive reset allowed")
	}
	router := gin.New()
	admin.Register(router, db, cfg)
	for _, tc := range []struct {
		password string
		status   int
	}{{oldPassword, 401}, {newPassword, 200}} {
		body, _ := json.Marshal(map[string]string{"email": user.Email, "password": tc.password})
		req := httptest.NewRequest("POST", "/api/admin/auth/login", bytes.NewReader(body))
		req.Header.Set("Origin", "http://localhost:5174")
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != tc.status {
			t.Fatalf("login status %d, expected %d", rec.Code, tc.status)
		}
	}
	// If auditing fails, both the password update and session deletion must roll back.
	var beforeSessions int64
	db.Model(&models.AdminSession{}).Where("user_id = ?", user.ID).Count(&beforeSessions)
	if beforeSessions == 0 {
		t.Fatal("login did not establish a session")
	}
	callback := "test:reject-password-reset-audit"
	if err := db.Callback().Create().Before("gorm:create").Register(callback, func(tx *gorm.DB) {
		if tx.Statement.Table == "admin_audit_logs" {
			tx.AddError(errors.New("audit unavailable"))
		}
	}); err != nil {
		t.Fatal(err)
	}
	defer db.Callback().Create().Remove(callback)
	if err := resetPassword(db, user, []byte("another-password-789"), []byte("another-password-789")); err == nil {
		t.Fatal("reset succeeded without audit")
	}
	var afterRollback models.User
	db.First(&afterRollback, user.ID)
	db.Model(&models.AdminSession{}).Where("user_id = ?", user.ID).Count(&n)
	if afterRollback.PasswordHash != updated.PasswordHash || n != beforeSessions {
		t.Fatal("reset did not roll back atomically")
	}
	// A non-admin account must never be promoted or reset by this command.
	if err := db.Model(&user).Update("role", "client").Error; err != nil {
		t.Fatal(err)
	}
	if _, err := findAdmin(db, user.Email); err != errAdminNotFound {
		t.Fatal("non-admin accepted")
	}
}
