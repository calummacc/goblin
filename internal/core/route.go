package goblin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RouteConfig struct
type RouteConfig struct {
	Method      string            // HTTP method
	Path        string            // Route path
	Handler     interface{}       // Route handler function
	Middlewares []gin.HandlerFunc // Middleware functions for this route
}

// Route struct
type Route struct {
	Prefix      string            // Prefix for the route group
	Middlewares []gin.HandlerFunc // Middleware functions for this route group
	RouteConfig []RouteConfig     // List of routes in this group
}

// Registers all routes of the application. Applies module middleware.
func RegisterRoutes(engine *gin.Engine, controllers []Controller) {
	fmt.Println("Registering routes...")
	for _, ctrl := range controllers {
		routes := ctrl.Routes()
		for _, route := range routes {
			group := engine.Group(route.Prefix, route.Middlewares...)
			for _, routeConfig := range route.RouteConfig {
				handler := routeConfig.Handler.(func(*gin.Context))
				fullPath := fmt.Sprintf("%s%s", route.Prefix, routeConfig.Path)
				allHandlers := append(routeConfig.Middlewares, handler)
				switch routeConfig.Method {
				case http.MethodGet:
					group.GET(routeConfig.Path, allHandlers...)
				case http.MethodPost:
					group.POST(routeConfig.Path, allHandlers...)
				case http.MethodPut:
					group.PUT(routeConfig.Path, allHandlers...)
				case http.MethodDelete:
					group.DELETE(routeConfig.Path, allHandlers...)
				case http.MethodPatch:
					group.PATCH(routeConfig.Path, allHandlers...)
				case http.MethodOptions:
					group.OPTIONS(routeConfig.Path, allHandlers...)
				case http.MethodHead:
					group.HEAD(routeConfig.Path, allHandlers...)
				default:
					engine.NoRoute(func(c *gin.Context) {
						c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "Method Not Allowed"})
					})
				}
				fmt.Printf("Registered route: %s %s\n", routeConfig.Method, fullPath)
			}
		}
	}
}

// registerMiddleware registers global and module middleware.
func registerMiddleware(engine *gin.Engine, globalMiddleware *globalMiddleware) {
	engine.Use(globalMiddleware.getMiddlewares()...)
}
