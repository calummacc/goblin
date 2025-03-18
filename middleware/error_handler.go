package middleware

import (
	"goblin/errors"

	"github.com/gin-gonic/gin"
)

// ErrorHandler is a middleware that handles errors globally
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Handle AppError
			if appErr, ok := err.(*errors.AppError); ok {
				c.JSON(appErr.Status, errors.ErrorResponse{
					Type:    appErr.Type,
					Message: appErr.Message,
					Status:  appErr.Status,
				})
				return
			}

			// Handle other errors as internal server error
			c.JSON(500, errors.ErrorResponse{
				Type:    errors.ErrorTypeInternal,
				Message: "Internal Server Error",
				Status:  500,
			})
		}
	}
}
