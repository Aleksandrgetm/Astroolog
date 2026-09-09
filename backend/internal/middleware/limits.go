package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func BodyLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 16*1024)
		c.Next()
	}
}
