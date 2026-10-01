package admin

import (
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"testing"
)

func TestStrictAdminOrigin(t *testing.T) {
	for _, tc := range []struct {
		name, origin string
		allowed      []string
		want         bool
	}{
		{"localhost", "http://localhost:5174", []string{"http://localhost:5174", "http://127.0.0.1:5174"}, true},
		{"loopback", "http://127.0.0.1:5174", []string{"http://localhost:5174", "http://127.0.0.1:5174"}, true},
		{"old port", "http://localhost:5173", []string{"http://localhost:5174"}, false},
		{"similar port", "http://localhost:51740", []string{"http://localhost:5174"}, false},
		{"similar host", "http://evil-localhost:5174", []string{"http://localhost:5174"}, false},
		{"evil", "http://evil.example", []string{"http://localhost:5174"}, false},
		{"malformed", "://", []string{"http://localhost:5174"}, false},
		{"path", "http://localhost:5174/", []string{"http://localhost:5174"}, false},
		{"credentials", "http://attacker@localhost:5174", []string{"http://localhost:5174"}, false},
		{"query", "http://localhost:5174?", []string{"http://localhost:5174"}, false},
		{"fragment", "http://localhost:5174#", []string{"http://localhost:5174"}, false},
		{"empty", "", []string{"http://localhost:5174"}, false},
		{"production HTTPS", "https://admin.example.invalid", []string{"https://admin.example.invalid"}, true},
		{"production localhost", "http://localhost:5174", []string{"https://admin.example.invalid"}, false},
		{"production unknown", "https://evil.example", []string{"https://admin.example.invalid"}, false},
		{"production unconfigured", "https://admin.example.invalid", nil, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest("POST", "http://localhost:5173/api/admin/auth/login", nil)
			c.Request.Header.Set("Origin", tc.origin)
			if got := (Handler{AllowedOrigins: tc.allowed}).originOK(c); got != tc.want {
				t.Fatalf("allowed = %v", got)
			}
		})
	}
}
