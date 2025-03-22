package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				log.Printf("PANIC: %v\n%s", err, string(stack))

				requestID, exists := c.Get("RequestID")
				errorID := "unknown"
				if exists {
					errorID = requestID.(string)
				}

				c.JSON(http.StatusInternalServerError, gin.H{
					"error":    "Internal Server Error",
					"error_id": errorID,
					"message":  fmt.Sprintf("Recovered from panic: %v", err),
				})

				c.Abort()
			}
		}()

		c.Next()
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			requestID, exists := c.Get("RequestID")
			errorID := "unknown"
			if exists {
				errorID = requestID.(string)
			}

			log.Printf("Error: %v", err.Err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error":    "Internal Server Error",
				"error_id": errorID,
				"message":  err.Error(),
			})
		}
	}
}
