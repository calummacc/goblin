// goblin/router/router.go
package router

import (
	"github.com/gin-gonic/gin"
)

// Controller is a marker interface for controllers
type Controller interface {
	BasePath() string
}

// Route represents a route in the application
type Route struct {
	Method      string
	Path        string
	Handler     gin.HandlerFunc
	Middlewares []gin.HandlerFunc
}

// RouterRegistry manages the routes in the application
type RouterRegistry struct {
	routes map[string][]Route
}

// NewRouterRegistry creates a new router registry
func NewRouterRegistry() *RouterRegistry {
	return &RouterRegistry{
		routes: make(map[string][]Route),
	}
}

// RegisterController registers a controller with the router
func (r *RouterRegistry) RegisterController(basePath string, controller interface{}) {
	// Retrieve routes from controller
	routes := r.getRoutesFromController(basePath, controller)
	r.routes[basePath] = routes
}

// getRoutesFromController extracts routes from a controller
func (r *RouterRegistry) getRoutesFromController(basePath string, controller interface{}) []Route {
	// This is a simplified implementation
	// In a real implementation, you'd use reflection to extract routes
	// For now, we'll assume controllers have a Routes() method

	// This is a placeholder for the actual implementation
	return []Route{}
}

// ConfigureRoutes configures the routes with the Gin engine
func (r *RouterRegistry) ConfigureRoutes(engine *gin.Engine) {
	for basePath, routes := range r.routes {
		for _, route := range routes {
			fullPath := basePath + route.Path
			engine.Handle(route.Method, fullPath, append(route.Middlewares, route.Handler)...)
		}
	}
}

// NewRouter creates a new router to be used in the DI container
func NewRouter(engine *gin.Engine) *RouterRegistry {
	registry := NewRouterRegistry()
	registry.ConfigureRoutes(engine)
	return registry
}
