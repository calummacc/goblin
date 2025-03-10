package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LoggerMiddleware logs request information.
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		requestID := c.GetString("RequestID")

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		log.Printf("| %3d | %13v | %15s | %s  %s | RequestID: %s",
			statusCode,
			latency,
			clientIP,
			method,
			path,
			requestID,
		)
	}
}

// AuthMiddleware checks for a valid authentication token.  Replace with your actual authentication logic.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		// Replace this with your actual token validation logic
		if token != "secret" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		c.Next()
	}
}

// CORSMiddleware adds CORS headers.
func CORSMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	return cors.New(config)
}

// RateLimitMiddleware limits the number of requests per IP address.
func RateLimitMiddleware(limit int, duration time.Duration) gin.HandlerFunc {
	type ipRecord struct {
		count     int
		timestamp time.Time
	}

	var (
		mu       sync.Mutex
		requests = make(map[string]ipRecord)
	)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		mu.Lock()
		defer mu.Unlock()

		record, ok := requests[ip]
		if !ok || now.Sub(record.timestamp) > duration {
			requests[ip] = ipRecord{count: 1, timestamp: now}
		} else {
			if record.count >= limit {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests"})
				return
			}
			requests[ip] = ipRecord{count: record.count + 1, timestamp: now}
		}

		c.Next()
	}
}

// RequestIDMiddleware adds a unique request ID to each request.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		c.Set("RequestID", requestID)
		fmt.Printf("RequestID: %s\n", requestID)
		c.Next()
	}
}
