package admin

import (
	"astroolog/backend/internal/models"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Handler struct {
	DB    *gorm.DB
	Store Storage
}

func token() string {
	b := make([]byte, 32)
	if _, e := rand.Read(b); e != nil {
		panic(e)
	}
	return hex.EncodeToString(b)
}
func hash(s string) string { b := sha256.Sum256([]byte(s)); return hex.EncodeToString(b[:]) }
func cookieName() string {
	if os.Getenv("APP_ENV") == "production" {
		return "__Host-astroolog_admin"
	}
	return "astroolog_admin"
}
func cookie(c *gin.Context, value string, age int) {
	http.SetCookie(c.Writer, &http.Cookie{Name: cookieName(), Value: value, Path: "/", HttpOnly: true, Secure: os.Getenv("APP_ENV") == "production", SameSite: http.SameSiteStrictMode, MaxAge: age})
}
func audit(db *gorm.DB, user uint, action, kind, id string) error {
	return db.Create(&models.AdminAuditLog{UserID: user, Action: action, EntityType: kind, EntityID: id, Metadata: models.Document{}}).Error
}
func uid(c *gin.Context) uint { return c.MustGet("admin").(models.User).ID }
func originOK(c *gin.Context) bool {
	origin := c.GetHeader("Origin")
	u, e := url.Parse(origin)
	if e != nil || origin == "" || u.User != nil {
		return false
	}
	allowed := os.Getenv("APP_ORIGIN")
	if allowed != "" {
		return origin == strings.TrimRight(allowed, "/")
	}
	return os.Getenv("APP_ENV") != "production" && (u.Host == c.Request.Host || origin == "http://localhost:5173" || origin == "http://127.0.0.1:5173") && (u.Scheme == "http" || u.Scheme == "https")
}
func (h Handler) throttle(key string, limit int) bool {
	allowed := false
	e := h.DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		row := models.LoginAttempt{Key: hash(key), WindowStart: now}
		if e := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&row).Error; e != nil {
			return e
		}
		if e := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&row, "key=?", hash(key)).Error; e != nil {
			return e
		}
		if now.Sub(row.WindowStart) > 15*time.Minute {
			row.Count = 0
			row.WindowStart = now
		}
		row.Count++
		allowed = row.Count <= limit
		return tx.Save(&row).Error
	})
	return e == nil && allowed
}

var dummyPassword, _ = bcrypt.GenerateFromPassword([]byte("dummy-password-not-an-account"), 12)

func (h Handler) Login(c *gin.Context) {
	if !originOK(c) {
		c.AbortWithStatusJSON(403, gin.H{"error": "Недопустимый источник запроса."})
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 4096)
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if c.ShouldBindJSON(&input) != nil || len(input.Email) > 254 || len(input.Password) > 72 {
		c.JSON(400, gin.H{"error": "Неверный email или пароль."})
		return
	}
	email := strings.ToLower(strings.TrimSpace(input.Email))
	if !h.throttle("ip:"+c.ClientIP(), 30) || !h.throttle("email:"+email, 8) {
		c.Header("Retry-After", "900")
		c.JSON(429, gin.H{"error": "Слишком много попыток. Попробуйте через 15 минут."})
		return
	}
	var user models.User
	e := h.DB.First(&user, "email=?", email).Error
	password := dummyPassword
	if e == nil {
		password = []byte(user.PasswordHash)
	}
	if bcrypt.CompareHashAndPassword(password, []byte(input.Password)) != nil || e != nil || !user.IsActive || user.Role != "admin" {
		c.JSON(401, gin.H{"error": "Неверный email или пароль."})
		return
	}
	raw := token()
	csrf := hash("csrf:" + raw)
	session := models.AdminSession{ID: hash(raw), UserID: user.ID, CSRFHash: hash(csrf), ExpiresAt: time.Now().Add(12 * time.Hour)}
	e = h.DB.Transaction(func(tx *gorm.DB) error {
		if old, e := c.Cookie(cookieName()); e == nil {
			if e := tx.Delete(&models.AdminSession{}, "id=?", hash(old)).Error; e != nil {
				return e
			}
		}
		if e := tx.Create(&session).Error; e != nil {
			return e
		}
		if e := tx.Model(&user).Update("last_login_at", time.Now()).Error; e != nil {
			return e
		}
		return audit(tx, user.ID, "login", "session", "")
	})
	if e != nil {
		c.JSON(500, gin.H{"error": "Не удалось войти."})
		return
	}
	cookie(c, raw, 12*3600)
	c.JSON(200, gin.H{"user": user, "csrf": csrf})
}
func (h Handler) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, e := c.Cookie(cookieName())
		if e != nil || len(raw) != 64 {
			c.AbortWithStatusJSON(401, gin.H{"error": "Войдите в административную панель."})
			return
		}
		var session models.AdminSession
		var user models.User
		if h.DB.First(&session, "id=? AND expires_at>?", hash(raw), time.Now()).Error != nil || h.DB.First(&user, "id=? AND is_active=true", session.UserID).Error != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Сессия истекла. Войдите снова."})
			return
		}
		c.Set("admin", user)
		c.Set("session", session)
		c.Set("csrf", hash("csrf:"+raw))
		c.Next()
	}
}
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.MustGet("admin").(models.User).Role != role {
			c.AbortWithStatus(403)
			return
		}
		c.Next()
	}
}
func CSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodHead {
			c.Next()
			return
		}
		session := c.MustGet("session").(models.AdminSession)
		if !originOK(c) || subtle.ConstantTimeCompare([]byte(hash(c.GetHeader("X-CSRF-Token"))), []byte(session.CSRFHash)) != 1 {
			c.AbortWithStatusJSON(403, gin.H{"error": "Обновите страницу: проверка безопасности не пройдена."})
			return
		}
		c.Next()
	}
}
func (h Handler) Me(c *gin.Context) {
	c.JSON(200, gin.H{"user": c.MustGet("admin"), "csrf": c.MustGet("csrf")})
}
func (h Handler) Logout(c *gin.Context) {
	s := c.MustGet("session").(models.AdminSession)
	if e := h.DB.Transaction(func(tx *gorm.DB) error {
		if e := tx.Delete(&s).Error; e != nil {
			return e
		}
		return audit(tx, uid(c), "logout", "session", "")
	}); e != nil {
		c.JSON(500, gin.H{"error": "Не удалось завершить сессию."})
		return
	}
	cookie(c, "", -1)
	c.Status(204)
}
