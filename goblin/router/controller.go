// goblin/router/controller.go
package router

import (
	"github.com/gin-gonic/gin"
)

// ControllerMetadata represents metadata for a controller
type ControllerMetadata struct {
	BasePath string
	Routes   []Route
}

// ControllerBuilder is a builder for controllers
type ControllerBuilder struct {
	basePath    string
	routes      []Route
	middlewares []gin.HandlerFunc
}

// NewControllerBuilder creates a new controller builder
func NewControllerBuilder(basePath string) *ControllerBuilder {
	return &ControllerBuilder{
		basePath: basePath,
		routes:   []Route{},
	}
}

// AddRoute adds a route to the controller
func (b *ControllerBuilder) AddRoute(method, path string, handler gin.HandlerFunc, middlewares ...gin.HandlerFunc) *ControllerBuilder {
	b.routes = append(b.routes, Route{
		Method:      method,
		Path:        path,
		Handler:     handler,
		Middlewares: append(b.middlewares, middlewares...),
	})
	return b
}

// Get adds a GET route
func (b *ControllerBuilder) Get(path string, handler gin.HandlerFunc, middlewares ...gin.HandlerFunc) *ControllerBuilder {
	return b.AddRoute("GET", path, handler, middlewares...)
}

// Post adds a POST route
func (b *ControllerBuilder) Post(path string, handler gin.HandlerFunc, middlewares ...gin.HandlerFunc) *ControllerBuilder {
	return b.AddRoute("POST", path, handler, middlewares...)
}

// Put adds a PUT route
func (b *ControllerBuilder) Put(path string, handler gin.HandlerFunc, middlewares ...gin.HandlerFunc) *ControllerBuilder {
	return b.AddRoute("PUT", path, handler, middlewares...)
}

// Delete adds a DELETE route
func (b *ControllerBuilder) Delete(path string, handler gin.HandlerFunc, middlewares ...gin.HandlerFunc) *ControllerBuilder {
	return b.AddRoute("DELETE", path, handler, middlewares...)
}

// WithMiddleware adds middleware to the controller
func (b *ControllerBuilder) WithMiddleware(middlewares ...gin.HandlerFunc) *ControllerBuilder {
	b.middlewares = append(b.middlewares, middlewares...)
	return b
}

// Build builds the controller metadata
func (b *ControllerBuilder) Build() ControllerMetadata {
	return ControllerMetadata{
		BasePath: b.basePath,
		Routes:   b.routes,
	}
}
