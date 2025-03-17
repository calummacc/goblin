// goblin/http/request.go
package http

import (
	"github.com/gin-gonic/gin"
)

// Request represents an HTTP request
type Request struct {
	Context *gin.Context
}

// NewRequest creates a new request from a Gin context
func NewRequest(c *gin.Context) *Request {
	return &Request{
		Context: c,
	}
}

// GetBody extracts the request body
func (r *Request) GetBody(into interface{}) error {
	return r.Context.ShouldBindJSON(into)
}

// GetQuery extracts a query parameter
func (r *Request) GetQuery(name string) string {
	return r.Context.Query(name)
}

// GetParam extracts a URL parameter
func (r *Request) GetParam(name string) string {
	return r.Context.Param(name)
}

// GetHeader extracts a request header
func (r *Request) GetHeader(name string) string {
	return r.Context.GetHeader(name)
}

// GetMethod returns the HTTP method of the request
func (r *Request) GetMethod() string {
	return r.Context.Request.Method
}
