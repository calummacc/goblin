package goblin

import (
	"fmt"

	"github.com/calummacc/goblin/internal/common/enums"
	"github.com/gin-gonic/gin"
)

// RouteConfig defines a single route
type RouteConfig struct {
	Method      enums.RequestMethod
	Path        string
	Handler     interface{}
	Middlewares []gin.HandlerFunc
}

// Route defines a group route and prefix
type Route struct {
	Prefix      string
	RouteConfig []RouteConfig
}

// Controller interface
type Controller interface {
	Routes() []Route
}

// Registers all route of app
func RegisterRoutes(engine *gin.Engine, controllers []Controller) {
	for _, ctrl := range controllers {
		routes := ctrl.Routes()
		for _, route := range routes {
			group := engine.Group(route.Prefix)
			for _, routeConfig := range route.RouteConfig {
				handler := routeConfig.Handler.(func(*gin.Context))
				fullPath := fmt.Sprintf("%s%s", route.Prefix, routeConfig.Path)
				switch routeConfig.Method {
				case enums.GET:
					group.GET(routeConfig.Path, append(routeConfig.Middlewares, handler)...)
				case enums.POST:
					group.POST(routeConfig.Path, append(routeConfig.Middlewares, handler)...)
				case enums.PUT:
					group.PUT(routeConfig.Path, append(routeConfig.Middlewares, handler)...)
				case enums.DELETE:
					group.DELETE(routeConfig.Path, append(routeConfig.Middlewares, handler)...)
				case enums.PATCH:
					group.PATCH(routeConfig.Path, append(routeConfig.Middlewares, handler)...)
				case enums.ALL:
					group.Any(routeConfig.Path, append(routeConfig.Middlewares, handler)...)
				case enums.OPTIONS:
					group.OPTIONS(routeConfig.Path, append(routeConfig.Middlewares, handler)...)
				case enums.HEAD:
					group.HEAD(routeConfig.Path, append(routeConfig.Middlewares, handler)...)
				case enums.SEARCH:
					group.Handle("SEARCH", routeConfig.Path, append(routeConfig.Middlewares, handler)...)
				case enums.PROPFIND:
					group.Handle("PROPFIND", routeConfig.Path, append(routeConfig.Middlewares, handler)...)
				case enums.PROPPATCH:
					group.Handle("PROPPATCH", routeConfig.Path, append(routeConfig.Middlewares, handler)...)
				case enums.MKCOL:
					group.Handle("MKCOL", routeConfig.Path, append(routeConfig.Middlewares, handler)...)
				case enums.COPY:
					group.Handle("COPY", routeConfig.Path, append(routeConfig.Middlewares, handler)...)
				case enums.MOVE:
					group.Handle("MOVE", routeConfig.Path, append(routeConfig.Middlewares, handler)...)
				case enums.LOCK:
					group.Handle("LOCK", routeConfig.Path, append(routeConfig.Middlewares, handler)...)
				case enums.UNLOCK:
					group.Handle("UNLOCK", routeConfig.Path, append(routeConfig.Middlewares, handler)...)
				default:
					panic(fmt.Sprintf("Unsupported HTTP method: %s", routeConfig.Method))
				}
				fmt.Printf("Register route %s %s \n", routeConfig.Method, fullPath)
			}
		}
	}
}
