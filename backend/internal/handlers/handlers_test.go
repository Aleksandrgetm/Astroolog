package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRejectInvalidRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cases := []struct {
		name, body string
		booking    bool
	}{
		{"missing consent", `{"name":"Анна","email":"anna@example.com","message":"Вопрос"}`, false},
		{"invalid email", `{"name":"Анна","email":"bad","message":"Вопрос","consent":true}`, false},
		{"whitespace name", `{"name":"  ","email":"anna@example.com","message":"Вопрос","consent":true}`, false},
		{"empty contact message", `{"name":"Анна","email":"anna@example.com","message":"  ","consent":true}`, false},
		{"invalid phone", `{"name":"Анна","email":"anna@example.com","phone":"abcdef","message":"Вопрос","consent":true}`, false},
		{"missing service type", `{"name":"Анна","email":"anna@example.com","consent":true,"service_id":1}`, true},
		{"missing service", `{"name":"Анна","email":"anna@example.com","consent":true}`, true},
		{"malformed JSON", `{`, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.New()
			h := Handler{}
			r.POST("/", func(c *gin.Context) { h.Submit(c, tc.booking) })
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			if w.Code != 400 {
				t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
			}
		})
	}
}
