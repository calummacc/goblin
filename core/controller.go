// Package core provides core functionality for the Goblin Framework.
// It includes controller management, routing, and middleware support.
package core

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

// ControllerRegistry stores metadata for all registered controllers
var ControllerRegistry = make(map[reflect.Type]*ControllerMetadata)

// Route represents a single HTTP route in the application.
// It defines the HTTP method, path, handler function, and any middleware
// that should be applied to the route.
type Route struct {
	// Method specifies the HTTP method (GET, POST, PUT, DELETE, etc.)
	Method string
	// Path specifies the URL path for this route
	Path string
	// Handler is the function that handles the HTTP request
	Handler gin.HandlerFunc
	// Middlewares is a list of middleware functions to be applied to this route
	Middlewares []gin.HandlerFunc
}

// RouteMetadata contains information about a route in the Goblin Framework.
// It defines the path, HTTP method, handler function, and associated middleware and guards.
type RouteMetadata struct {
	// Path is the URL path pattern for the route
	Path string
	// Method is the HTTP method (GET, POST, PUT, DELETE, etc.)
	Method string
	// Handler is the function that handles the route
	Handler interface{}
	// Middleware contains the list of middleware functions to be executed before the handler
	Middleware []gin.HandlerFunc
	// Guards contains the list of guards that protect this route
	Guards []Guard
}

// ControllerMetadata contains information about a controller in the Goblin Framework.
// It defines the route prefix, routes, and associated middleware and guards.
type ControllerMetadata struct {
	// Prefix is the URL prefix for all routes in this controller
	Prefix string
	// Routes contains the list of routes defined in this controller
	Routes []RouteMetadata
	// Guards contains the list of guards that protect all routes in this controller
	Guards []Guard
	// Middleware contains the list of middleware functions to be executed for all routes
	Middleware []gin.HandlerFunc
}

// Controller defines the interface for controllers in the Goblin Framework.
// Controllers are responsible for handling HTTP requests and defining routes.
type Controller interface {
	// GetMetadata returns the metadata for this controller
	GetMetadata() *ControllerMetadata
}

// BaseController provides a default implementation of the Controller interface.
// It implements basic metadata management functionality that can be extended
// by custom controllers.
type BaseController struct {
	metadata *ControllerMetadata
}

// NewBaseController creates a new base controller
func NewBaseController() *BaseController {
	return &BaseController{
		metadata: &ControllerMetadata{},
	}
}

// GetMetadata returns the metadata for this controller.
//
// Returns:
//   - *ControllerMetadata: The controller's metadata
func (c *BaseController) GetMetadata() *ControllerMetadata {
	return c.metadata
}

// SetMetadata sets the metadata for this controller.
//
// Parameters:
//   - metadata: The metadata to set
func (c *BaseController) SetMetadata(metadata *ControllerMetadata) {
	c.metadata = metadata
}

// UseMiddleware adds middleware to this controller.
// The middleware will be executed for all routes in this controller.
//
// Parameters:
//   - middleware: The middleware function to add
func (c *BaseController) UseMiddleware(middleware gin.HandlerFunc) {
	if c.metadata == nil {
		c.metadata = &ControllerMetadata{
			Routes:     make([]RouteMetadata, 0),
			Guards:     make([]Guard, 0),
			Middleware: make([]gin.HandlerFunc, 0),
		}
	}
	c.metadata.Middleware = append(c.metadata.Middleware, middleware)
}

// UseGuards adds guards to this controller.
// The guards will protect all routes in this controller.
//
// Parameters:
//   - guard: The guard to add
func (c *BaseController) UseGuards(guard Guard) {
	if c.metadata == nil {
		c.metadata = &ControllerMetadata{
			Routes:     make([]RouteMetadata, 0),
			Guards:     make([]Guard, 0),
			Middleware: make([]gin.HandlerFunc, 0),
		}
	}
	c.metadata.Guards = append(c.metadata.Guards, guard)
}

// Controller decorators

// SetController is a decorator that registers a prefix path for a controller.
// It creates metadata for the controller and stores it in the registry.
//
// Parameters:
//   - prefix: The URL prefix for all routes in this controller
//
// Returns:
//   - func(target interface{}): A decorator function that can be used to decorate controllers
func SetController(prefix string) func(target interface{}) {
	return func(target interface{}) {
		t := reflect.TypeOf(target)
		metadata := &ControllerMetadata{
			Prefix:     prefix,
			Routes:     make([]RouteMetadata, 0),
			Guards:     make([]Guard, 0),
			Middleware: make([]gin.HandlerFunc, 0),
		}
		ControllerRegistry[t] = metadata

		// If target implements Controller interface, set metadata
		if ctrl, ok := target.(Controller); ok {
			if bc, ok := ctrl.(*BaseController); ok {
				bc.SetMetadata(metadata)
			}
		}
	}
}

// Route decorators

// Get is a decorator that registers a GET route.
//
// Parameters:
//   - path: The URL path pattern
//   - middleware: Optional middleware functions to be executed before the handler
//
// Returns:
//   - func(target interface{}, methodName string): A decorator function for the route handler
func Get(path string, middleware ...gin.HandlerFunc) func(target interface{}, methodName string) {
	return route("GET", path, middleware...)
}

// Post is a decorator that registers a POST route.
//
// Parameters:
//   - path: The URL path pattern
//   - middleware: Optional middleware functions to be executed before the handler
//
// Returns:
//   - func(target interface{}, methodName string): A decorator function for the route handler
func Post(path string, middleware ...gin.HandlerFunc) func(target interface{}, methodName string) {
	return route("POST", path, middleware...)
}

// Put is a decorator that registers a PUT route.
//
// Parameters:
//   - path: The URL path pattern
//   - middleware: Optional middleware functions to be executed before the handler
//
// Returns:
//   - func(target interface{}, methodName string): A decorator function for the route handler
func Put(path string, middleware ...gin.HandlerFunc) func(target interface{}, methodName string) {
	return route("PUT", path, middleware...)
}

// Delete is a decorator that registers a DELETE route.
//
// Parameters:
//   - path: The URL path pattern
//   - middleware: Optional middleware functions to be executed before the handler
//
// Returns:
//   - func(target interface{}, methodName string): A decorator function for the route handler
func Delete(path string, middleware ...gin.HandlerFunc) func(target interface{}, methodName string) {
	return route("DELETE", path, middleware...)
}

// route is a helper function that creates route decorators.
// It registers route metadata in the controller registry.
//
// Parameters:
//   - method: The HTTP method
//   - path: The URL path pattern
//   - middleware: Optional middleware functions
//
// Returns:
//   - func(target interface{}, methodName string): A decorator function for the route handler
func route(method string, path string, middleware ...gin.HandlerFunc) func(target interface{}, methodName string) {
	return func(target interface{}, methodName string) {
		t := reflect.TypeOf(target)
		metadata, ok := ControllerRegistry[t]
		if !ok {
			metadata = &ControllerMetadata{
				Routes:     make([]RouteMetadata, 0),
				Guards:     make([]Guard, 0),
				Middleware: make([]gin.HandlerFunc, 0),
			}
			ControllerRegistry[t] = metadata
		}

		// Get method from target
		m, ok := t.MethodByName(methodName)
		if !ok {
			panic(fmt.Sprintf("Method %s not found in controller %s", methodName, t.Name()))
		}

		// Add route to metadata
		metadata.Routes = append(metadata.Routes, RouteMetadata{
			Path:       path,
			Method:     method,
			Handler:    m.Func.Interface(),
			Middleware: middleware,
		})
	}
}

// Guard decorators

// UseGuards is a decorator that adds guards to a controller or route.
// It can be applied to either the entire controller or a specific route.
//
// Parameters:
//   - guards: The guards to add
//
// Returns:
//   - func(target interface{}, methodName string): A decorator function
func UseGuards(guards ...Guard) func(target interface{}, methodName string) {
	return func(target interface{}, methodName string) {
		t := reflect.TypeOf(target)
		metadata, ok := ControllerRegistry[t]
		if !ok {
			return
		}

		if methodName == "" {
			// Apply to entire controller
			metadata.Guards = append(metadata.Guards, guards...)
		} else {
			// Apply to specific route
			for i, route := range metadata.Routes {
				if strings.EqualFold(route.Method, methodName) {
					metadata.Routes[i].Guards = append(metadata.Routes[i].Guards, guards...)
				}
			}
		}
	}
}

// Middleware decorators

// UseMiddleware is a decorator that adds middleware to a controller or route.
// It can be applied to either the entire controller or a specific route.
//
// Parameters:
//   - middleware: The middleware functions to add
//
// Returns:
//   - func(target interface{}, methodName string): A decorator function
func UseMiddleware(middleware ...gin.HandlerFunc) func(target interface{}, methodName string) {
	return func(target interface{}, methodName string) {
		t := reflect.TypeOf(target)
		metadata, ok := ControllerRegistry[t]
		if !ok {
			return
		}

		if methodName == "" {
			// Apply to entire controller
			metadata.Middleware = append(metadata.Middleware, middleware...)
		} else {
			// Apply to specific route
			for i, route := range metadata.Routes {
				if strings.EqualFold(route.Method, methodName) {
					metadata.Routes[i].Middleware = append(metadata.Routes[i].Middleware, middleware...)
				}
			}
		}
	}
}

// ControllerManager manages the registration and initialization of controllers.
// It handles the integration of controllers with the Gin web framework.
type ControllerManager struct {
	// engine is the Gin engine instance
	engine *gin.Engine
}

// NewControllerManager creates a new ControllerManager instance.
//
// Parameters:
//   - engine: The Gin engine to use
//
// Returns:
//   - *ControllerManager: A new controller manager
func NewControllerManager(engine *gin.Engine) *ControllerManager {
	return &ControllerManager{
		engine: engine,
	}
}

// RegisterController registers a controller with the Gin engine.
// It sets up all routes, middleware, and guards defined in the controller.
// The controller must have metadata registered through decorators.
//
// The registration process:
// 1. Creates a router group with the controller's prefix
// 2. Applies controller-level middleware
// 3. For each route:
//   - Adds controller-level guards
//   - Adds route-specific guards
//   - Adds route-specific middleware
//   - Adds the main handler
//
// 4. Registers the route with Gin
//
// Parameters:
//   - controller: The controller to register
//
// Returns:
//   - error: Any error that occurred during registration
func (m *ControllerManager) RegisterController(controller interface{}) error {
	t := reflect.TypeOf(controller)
	metadata, ok := ControllerRegistry[t]
	if !ok {
		return fmt.Errorf("no metadata found for controller %s", t.Name())
	}

	// Create router group with prefix
	group := m.engine.Group(metadata.Prefix)

	// Apply middleware for entire controller
	group.Use(metadata.Middleware...)

	// Register each route
	for _, route := range metadata.Routes {
		handlers := make([]gin.HandlerFunc, 0)

		// Add controller guards
		for _, guard := range metadata.Guards {
			handlers = append(handlers, createGuardMiddleware(guard))
		}

		// Add route guards
		for _, guard := range route.Guards {
			handlers = append(handlers, createGuardMiddleware(guard))
		}

		// Add route middleware
		handlers = append(handlers, route.Middleware...)

		// Add main handler
		handlers = append(handlers, createRouteHandler(controller, route.Handler))

		// Register route with Gin
		group.Handle(route.Method, route.Path, handlers...)
	}

	return nil
}

// createGuardMiddleware creates a Gin middleware function from a Guard.
// It checks if the request should be allowed to proceed and aborts with
// a 403 status code if the guard denies access.
//
// Parameters:
//   - guard: The guard to create middleware from
//
// Returns:
//   - gin.HandlerFunc: A middleware function that implements the guard's logic
func createGuardMiddleware(guard Guard) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Context{
			GinContext: c,
		}
		if ok, err := guard.CanActivate(ctx); !ok {
			if err != nil {
				c.AbortWithError(403, err)
			} else {
				c.AbortWithStatus(403)
			}
			return
		}
		c.Next()
	}
}

// createRouteHandler creates a handler function for a route.
// It uses reflection to call the handler method with the controller instance
// and Gin context.
//
// Parameters:
//   - controller: The controller instance
//   - handler: The handler function to call
//
// Returns:
//   - gin.HandlerFunc: A handler function that calls the route handler
func createRouteHandler(controller interface{}, handler interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Call handler with context
		reflect.ValueOf(handler).Call([]reflect.Value{
			reflect.ValueOf(controller),
			reflect.ValueOf(c),
		})
	}
}
