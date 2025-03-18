package core

import (
	"github.com/gin-gonic/gin"
)

// Context represents the request context in the Goblin Framework.
// It wraps the Gin context and provides additional information about
// the current request handler and class.
type Context struct {
	// GinContext is the underlying Gin context that contains the HTTP request/response data
	GinContext *gin.Context
	// Handler is the current request handler being executed
	Handler interface{}
	// Class is the controller or service class that contains the handler
	Class interface{}
}

// Guard defines an interface for request authorization in the Goblin Framework.
// Guards are used to protect routes and determine if a request should be allowed
// to proceed to the handler. They can implement custom authorization logic
// such as authentication checks, role-based access control, or other security rules.
type Guard interface {
	// CanActivate determines if the current request should be allowed to proceed.
	// It receives the current request context and returns:
	// - bool: true if the request should proceed, false if it should be blocked
	// - error: any error that occurred during the authorization check
	CanActivate(ctx *Context) (bool, error)
}

// BaseGuard provides a default implementation of the Guard interface.
// It serves as a base class for custom guards and implements a permissive
// authorization policy that allows all requests by default.
type BaseGuard struct{}

// CanActivate implements the default authorization logic that allows all requests.
// This is a permissive implementation that should be overridden by custom guards
// to implement specific authorization rules.
//
// Parameters:
//   - ctx: The current request context
//
// Returns:
//   - bool: Always returns true to allow all requests
//   - error: Always returns nil as no errors can occur in this implementation
func (g *BaseGuard) CanActivate(ctx *Context) (bool, error) {
	return true, nil
}
