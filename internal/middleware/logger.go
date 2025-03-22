package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		log.Printf("[GIN] %3d | %13v | %15s | %-7s %s",
			status,
			latency,
			c.ClientIP(),
			c.Request.Method,
			path,
		)
	}
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := generateUUID()
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Set("RequestID", requestID)
		c.Next()
	}
}

func generateUUID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(6)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
