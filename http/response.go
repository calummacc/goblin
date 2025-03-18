// goblin/http/response.go
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents an HTTP response
type Response struct {
	Context *gin.Context
}

// NewResponse creates a new response from a Gin context
func NewResponse(c *gin.Context) *Response {
	return &Response{
		Context: c,
	}
}

// JSON sends a JSON response
func (r *Response) JSON(status int, data interface{}) {
	r.Context.JSON(status, data)
}

// Error sends an error response
func (r *Response) Error(status int, message string) {
	r.Context.JSON(status, gin.H{
		"error": message,
	})
}

// Success sends a success response
func (r *Response) Success(data interface{}) {
	r.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

// Created sends a created response
func (r *Response) Created(data interface{}) {
	r.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    data,
	})
}

// NoContent sends a no content response
func (r *Response) NoContent() {
	r.Context.Status(http.StatusNoContent)
}

// NotFound sends a not found response
func (r *Response) NotFound(message string) {
	r.Error(http.StatusNotFound, message)
}

// BadRequest sends a bad request response
func (r *Response) BadRequest(message string) {
	r.Error(http.StatusBadRequest, message)
}

// Unauthorized sends an unauthorized response
func (r *Response) Unauthorized(message string) {
	r.Error(http.StatusUnauthorized, message)
}

// Forbidden sends a forbidden response
func (r *Response) Forbidden(message string) {
	r.Error(http.StatusForbidden, message)
}

// InternalServerError sends an internal server error response
func (r *Response) InternalServerError(message string) {
	r.Error(http.StatusInternalServerError, message)
}
