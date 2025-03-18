package interceptor

import (
	"context"

	"github.com/gin-gonic/gin"
)

// HandlerFunc represents a route handler function
type HandlerFunc = gin.HandlerFunc

// ExecutionContext holds information about the current request/response cycle
type ExecutionContext struct {
	// Original Gin context
	GinContext *gin.Context
	// Handler is the original route handler
	Handler HandlerFunc
	// Path is the route path
	Path string
	// Method is the HTTP method
	Method string
	// Custom data that can be shared between interceptors
	Data map[string]interface{}
}

// Interceptor defines the contract for request/response interceptors
type Interceptor interface {
	// Before is called before the route handler
	Before(ctx *ExecutionContext) error
	// After is called after the route handler
	After(ctx *ExecutionContext) error
}

// LifecycleInterceptor extends Interceptor with lifecycle hooks
type LifecycleInterceptor interface {
	Interceptor
	// OnRegister is called when the interceptor is registered
	OnRegister(ctx context.Context) error
	// OnShutdown is called when the application is shutting down
	OnShutdown(ctx context.Context) error
}

// BaseInterceptor provides a base implementation of LifecycleInterceptor
type BaseInterceptor struct{}

// Before implements Interceptor
func (i *BaseInterceptor) Before(ctx *ExecutionContext) error {
	return nil
}

// After implements Interceptor
func (i *BaseInterceptor) After(ctx *ExecutionContext) error {
	return nil
}

// OnRegister implements LifecycleInterceptor
func (i *BaseInterceptor) OnRegister(ctx context.Context) error {
	return nil
}

// OnShutdown implements LifecycleInterceptor
func (i *BaseInterceptor) OnShutdown(ctx context.Context) error {
	return nil
}

// Config holds configuration for an interceptor
type Config struct {
	// Name is a unique identifier for the interceptor
	Name string
	// Priority determines the order of interceptor execution (lower runs first)
	Priority int
	// Path specifies the route path pattern to match (empty means all paths)
	Path string
	// Methods specifies HTTP methods to match (empty means all methods)
	Methods []string
	// Interceptor is the actual interceptor implementation
	Interceptor interface{} // Can be Interceptor or LifecycleInterceptor
}
