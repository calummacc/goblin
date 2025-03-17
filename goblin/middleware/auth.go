// goblin/middleware/auth.go
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthOptions configures the auth middleware
type AuthOptions struct {
	TokenHeader string
	TokenField  string
}

// Auth creates an authentication middleware
func Auth(options ...AuthOptions) gin.HandlerFunc {
	opts := AuthOptions{
		TokenHeader: "Authorization",
		TokenField:  "token",
	}

	if len(options) > 0 {
		opts = options[0]
	}

	return func(c *gin.Context) {
		token := c.GetHeader(opts.TokenHeader)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		// Store the token in the context
		c.Set(opts.TokenField, token)
		c.Next()
	}
}

// Logger creates a logger middleware
func Logger() gin.HandlerFunc {
	return gin.Logger()
}

// Recovery creates a recovery middleware
func Recovery() gin.HandlerFunc {
	return gin.Recovery()
}
