package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
)

// MiddlewareFunc is a function that implements middleware logic
type MiddlewareFunc = gin.HandlerFunc

// Middleware interface defines the contract for middleware implementations
type Middleware interface {
	// Handle executes the middleware logic
	Handle(ctx *gin.Context)
}

// LifecycleMiddleware interface defines middleware with lifecycle hooks
type LifecycleMiddleware interface {
	Middleware
	// OnRegister is called when the middleware is registered
	OnRegister(ctx context.Context) error
	// OnShutdown is called when the application is shutting down
	OnShutdown(ctx context.Context) error
}

// Group represents a group of middleware that can be applied together
type Group struct {
	Name        string
	Middlewares []interface{} // Can be MiddlewareFunc, Middleware, or LifecycleMiddleware
}

// Options configures how middleware is applied
type Options struct {
	// Global indicates if the middleware should be applied globally
	Global bool
	// Priority determines the order of middleware execution (lower runs first)
	Priority int
	// Path specifies the route path pattern to match (empty means all paths)
	Path string
	// Methods specifies HTTP methods to match (empty means all methods)
	Methods []string
}

// Config holds configuration for a middleware
type Config struct {
	// Name is a unique identifier for the middleware
	Name string
	// Options configures how the middleware is applied
	Options Options
	// Middleware is the actual middleware implementation
	Middleware interface{} // Can be MiddlewareFunc, Middleware, or LifecycleMiddleware
}

// BaseMiddleware provides a base implementation of LifecycleMiddleware
type BaseMiddleware struct{}

// OnRegister implements LifecycleMiddleware
func (m *BaseMiddleware) OnRegister(ctx context.Context) error {
	return nil
}

// OnShutdown implements LifecycleMiddleware
func (m *BaseMiddleware) OnShutdown(ctx context.Context) error {
	return nil
}
