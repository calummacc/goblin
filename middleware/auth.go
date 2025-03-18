// goblin/middleware/auth.go
package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware implements LifecycleMiddleware for authentication
type AuthMiddleware struct {
	BaseMiddleware
	options AuthOptions
}

// AuthOptions configures the auth middleware
type AuthOptions struct {
	TokenHeader string
	TokenField  string
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(options ...AuthOptions) *AuthMiddleware {
	opts := AuthOptions{
		TokenHeader: "Authorization",
		TokenField:  "token",
	}

	if len(options) > 0 {
		opts = options[0]
	}

	return &AuthMiddleware{
		options: opts,
	}
}

// Handle implements Middleware interface
func (m *AuthMiddleware) Handle(c *gin.Context) {
	token := c.GetHeader(m.options.TokenHeader)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		c.Abort()
		return
	}

	// Store the token in the context
	c.Set(m.options.TokenField, token)
	c.Next()
}

// OnRegister implements LifecycleMiddleware
func (m *AuthMiddleware) OnRegister(ctx context.Context) error {
	// Add any initialization logic here
	return nil
}

// OnShutdown implements LifecycleMiddleware
func (m *AuthMiddleware) OnShutdown(ctx context.Context) error {
	// Add any cleanup logic here
	return nil
}

// Auth creates an authentication middleware (function style)
func Auth(options ...AuthOptions) gin.HandlerFunc {
	mw := NewAuthMiddleware(options...)
	return mw.Handle
}

// LoggerMiddleware implements Middleware for logging
type LoggerMiddleware struct {
	BaseMiddleware
}

// NewLoggerMiddleware creates a new logger middleware
func NewLoggerMiddleware() *LoggerMiddleware {
	return &LoggerMiddleware{}
}

// Handle implements Middleware interface
func (m *LoggerMiddleware) Handle(c *gin.Context) {
	gin.Logger()(c)
}

// Logger creates a logger middleware (function style)
func Logger() gin.HandlerFunc {
	mw := NewLoggerMiddleware()
	return mw.Handle
}

// RecoveryMiddleware implements Middleware for panic recovery
type RecoveryMiddleware struct {
	BaseMiddleware
}

// NewRecoveryMiddleware creates a new recovery middleware
func NewRecoveryMiddleware() *RecoveryMiddleware {
	return &RecoveryMiddleware{}
}

// Handle implements Middleware interface
func (m *RecoveryMiddleware) Handle(c *gin.Context) {
	gin.Recovery()(c)
}

// Recovery creates a recovery middleware (function style)
func Recovery() gin.HandlerFunc {
	mw := NewRecoveryMiddleware()
	return mw.Handle
}
