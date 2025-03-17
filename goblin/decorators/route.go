// goblin/decorators/route.go
package decorators

import (
	"goblin/router"

	"github.com/gin-gonic/gin"
)

// RouteDecorator is a decorator for routes
type RouteDecorator struct {
	method      string
	path        string
	handler     gin.HandlerFunc
	middlewares []gin.HandlerFunc
}

// NewRouteDecorator creates a new route decorator
func NewRouteDecorator(method, path string, handler gin.HandlerFunc, middlewares ...gin.HandlerFunc) *RouteDecorator {
	return &RouteDecorator{
		method:      method,
		path:        path,
		handler:     handler,
		middlewares: middlewares,
	}
}

// Get creates a GET route decorator
func Get(path string, handler gin.HandlerFunc, middlewares ...gin.HandlerFunc) *RouteDecorator {
	return NewRouteDecorator("GET", path, handler, middlewares...)
}

// Post creates a POST route decorator
func Post(path string, handler gin.HandlerFunc, middlewares ...gin.HandlerFunc) *RouteDecorator {
	return NewRouteDecorator("POST", path, handler, middlewares...)
}

// Put creates a PUT route decorator
func Put(path string, handler gin.HandlerFunc, middlewares ...gin.HandlerFunc) *RouteDecorator {
	return NewRouteDecorator("PUT", path, handler, middlewares...)
}

// Delete creates a DELETE route decorator
func Delete(path string, handler gin.HandlerFunc, middlewares ...gin.HandlerFunc) *RouteDecorator {
	return NewRouteDecorator("DELETE", path, handler, middlewares...)
}

// ToRoute converts a route decorator to a router.Route
func (d *RouteDecorator) ToRoute() router.Route {
	return router.Route{
		Method:      d.method,
		Path:        d.path,
		Handler:     d.handler,
		Middlewares: d.middlewares,
	}
}
