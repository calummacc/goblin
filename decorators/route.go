// goblin/decorators/route.go
package decorators

import (
	"fmt"
	"goblin/core"
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// RouteInfo represents metadata for a route
type RouteInfo struct {
	Path        string
	Method      string
	Handler     interface{}
	HandlerName string
	Middleware  []gin.HandlerFunc
	Guards      []core.Guard
	// Additional metadata fields
	Description string
	Tags        []string
	IsPublic    bool
	Deprecated  bool
	Parameters  []ParameterInfo
	Responses   map[int]ResponseInfo
}

// ParameterInfo represents metadata for a route parameter
type ParameterInfo struct {
	Name        string
	Description string
	Type        string
	Required    bool
	In          string // "path", "query", "header", "body"
}

// ResponseInfo represents metadata for a route response
type ResponseInfo struct {
	Description string
	Type        interface{}
	Headers     map[string]string
}

// Controller metadata registry
var (
	controllerRegistry     = make(map[reflect.Type]*core.ControllerMetadata)
	routeRegistry          = make(map[string]map[string]*RouteInfo) // controllerName -> methodName -> RouteInfo
	registryMutex          sync.RWMutex
	decoratorsRegistry     = make(map[reflect.Type]map[string][]interface{}) // controllerType -> methodName -> []decorators
	decoratorsRegistryLock sync.RWMutex
)

// RegisterDecorator registers a decorator for a controller method
func RegisterDecorator(target interface{}, methodName string, decorator interface{}) {
	decoratorsRegistryLock.Lock()
	defer decoratorsRegistryLock.Unlock()

	t := reflect.TypeOf(target)
	if decoratorsRegistry[t] == nil {
		decoratorsRegistry[t] = make(map[string][]interface{})
	}

	decoratorsRegistry[t][methodName] = append(decoratorsRegistry[t][methodName], decorator)
}

// GetDecorators returns all decorators for a controller method
func GetDecorators(target interface{}, methodName string) []interface{} {
	decoratorsRegistryLock.RLock()
	defer decoratorsRegistryLock.RUnlock()

	t := reflect.TypeOf(target)
	if decoratorsRegistry[t] == nil {
		return nil
	}

	return decoratorsRegistry[t][methodName]
}

// RouteDecorator is a decorator for routes
type RouteDecorator struct {
	Method      string
	Path        string
	Description string
	Tags        []string
	IsPublic    bool
	Deprecated  bool
}

// NewRouteDecorator creates a new route decorator
func NewRouteDecorator(method, path string) *RouteDecorator {
	return &RouteDecorator{
		Method:     method,
		Path:       path,
		IsPublic:   false,
		Deprecated: false,
		Tags:       []string{},
	}
}

// WithDescription adds a description to the route
func (d *RouteDecorator) WithDescription(description string) *RouteDecorator {
	d.Description = description
	return d
}

// WithTags adds tags to the route
func (d *RouteDecorator) WithTags(tags ...string) *RouteDecorator {
	d.Tags = append(d.Tags, tags...)
	return d
}

// AsPublic marks the route as public (not requiring authentication)
func (d *RouteDecorator) AsPublic() *RouteDecorator {
	d.IsPublic = true
	return d
}

// AsDeprecated marks the route as deprecated
func (d *RouteDecorator) AsDeprecated() *RouteDecorator {
	d.Deprecated = true
	return d
}

// Apply applies the decorator to a target method
func (d *RouteDecorator) Apply(target interface{}, methodName string) {
	RegisterDecorator(target, methodName, d)

	// Also register with controller registry for backward compatibility
	registerRouteWithController(d.Method, d.Path, target, methodName)
}

// HTTP Method decorators

// Get creates a decorator for GET requests
func Get(path string) func(target interface{}, methodName string) {
	decorator := NewRouteDecorator("GET", path)

	return func(target interface{}, methodName string) {
		decorator.Apply(target, methodName)
	}
}

// Post creates a decorator for POST requests
func Post(path string) func(target interface{}, methodName string) {
	decorator := NewRouteDecorator("POST", path)

	return func(target interface{}, methodName string) {
		decorator.Apply(target, methodName)
	}
}

// Put creates a decorator for PUT requests
func Put(path string) func(target interface{}, methodName string) {
	decorator := NewRouteDecorator("PUT", path)

	return func(target interface{}, methodName string) {
		decorator.Apply(target, methodName)
	}
}

// Delete creates a decorator for DELETE requests
func Delete(path string) func(target interface{}, methodName string) {
	decorator := NewRouteDecorator("DELETE", path)

	return func(target interface{}, methodName string) {
		decorator.Apply(target, methodName)
	}
}

// Patch creates a decorator for PATCH requests
func Patch(path string) func(target interface{}, methodName string) {
	decorator := NewRouteDecorator("PATCH", path)

	return func(target interface{}, methodName string) {
		decorator.Apply(target, methodName)
	}
}

// Options creates a decorator for OPTIONS requests
func Options(path string) func(target interface{}, methodName string) {
	decorator := NewRouteDecorator("OPTIONS", path)

	return func(target interface{}, methodName string) {
		decorator.Apply(target, methodName)
	}
}

// Head creates a decorator for HEAD requests
func Head(path string) func(target interface{}, methodName string) {
	decorator := NewRouteDecorator("HEAD", path)

	return func(target interface{}, methodName string) {
		decorator.Apply(target, methodName)
	}
}

// All creates a decorator for all HTTP methods
func All(path string) func(target interface{}, methodName string) {
	decorator := NewRouteDecorator("ANY", path)

	return func(target interface{}, methodName string) {
		decorator.Apply(target, methodName)
	}
}

// Guard decorators

// UseGuards creates a decorator that applies guards to a route
func UseGuards(guards ...core.Guard) func(target interface{}, methodName string) {
	return func(target interface{}, methodName string) {
		RegisterDecorator(target, methodName, &GuardDecorator{Guards: guards})

		// Also register with controller registry for backward compatibility
		t := reflect.TypeOf(target)
		metadata, ok := controllerRegistry[t]
		if !ok {
			return
		}

		// Apply guards to specific route
		for i, route := range metadata.Routes {
			if strings.EqualFold(route.Method, methodName) {
				metadata.Routes[i].Guards = append(metadata.Routes[i].Guards, guards...)
			}
		}
	}
}

// GuardDecorator is a decorator for applying guards
type GuardDecorator struct {
	Guards []core.Guard
}

// Middleware decorators

// UseMiddleware creates a decorator that applies middleware to a route
func UseMiddleware(middleware ...gin.HandlerFunc) func(target interface{}, methodName string) {
	return func(target interface{}, methodName string) {
		RegisterDecorator(target, methodName, &MiddlewareDecorator{Middleware: middleware})

		// Also register with controller registry for backward compatibility
		t := reflect.TypeOf(target)
		metadata, ok := controllerRegistry[t]
		if !ok {
			return
		}

		// Apply middleware to specific route
		for i, route := range metadata.Routes {
			if strings.EqualFold(route.Method, methodName) {
				metadata.Routes[i].Middleware = append(metadata.Routes[i].Middleware, middleware...)
			}
		}
	}
}

// MiddlewareDecorator is a decorator for applying middleware
type MiddlewareDecorator struct {
	Middleware []gin.HandlerFunc
}

// Controller decorators

// Controller creates a decorator that registers a controller with a base path
func Controller(basePath string) func(target interface{}) {
	return func(target interface{}) {
		t := reflect.TypeOf(target)
		metadata := &core.ControllerMetadata{
			Prefix:     basePath,
			Routes:     make([]core.RouteMetadata, 0),
			Guards:     make([]core.Guard, 0),
			Middleware: make([]gin.HandlerFunc, 0),
		}

		registryMutex.Lock()
		controllerRegistry[t] = metadata
		registryMutex.Unlock()

		// If target implements Controller interface, set metadata
		if ctrl, ok := target.(core.Controller); ok {
			if bc, ok := ctrl.(*core.BaseController); ok {
				bc.SetMetadata(metadata)
			}
		}
	}
}

// Description adds a description to a controller or route
func Description(description string) interface{} {
	// When applied to a controller
	controllerDecorator := func(target interface{}) {
		// Store description in controller metadata
		RegisterDecorator(target, "", &DescriptionDecorator{Description: description})
	}

	// When applied to a method
	methodDecorator := func(target interface{}, methodName string) {
		// Store description in route metadata
		RegisterDecorator(target, methodName, &DescriptionDecorator{Description: description})
	}

	// Return both types - the caller will use the appropriate one
	return struct {
		ControllerDecorator func(interface{})
		MethodDecorator     func(interface{}, string)
	}{
		ControllerDecorator: controllerDecorator,
		MethodDecorator:     methodDecorator,
	}
}

// DescriptionDecorator adds a description to a controller or route
type DescriptionDecorator struct {
	Description string
}

// Tags adds tags to a controller or route
func Tags(tags ...string) interface{} {
	// When applied to a controller
	controllerDecorator := func(target interface{}) {
		// Store tags in controller metadata
		RegisterDecorator(target, "", &TagsDecorator{Tags: tags})
	}

	// When applied to a method
	methodDecorator := func(target interface{}, methodName string) {
		// Store tags in route metadata
		RegisterDecorator(target, methodName, &TagsDecorator{Tags: tags})
	}

	// Return both types - the caller will use the appropriate one
	return struct {
		ControllerDecorator func(interface{})
		MethodDecorator     func(interface{}, string)
	}{
		ControllerDecorator: controllerDecorator,
		MethodDecorator:     methodDecorator,
	}
}

// TagsDecorator adds tags to a controller or route
type TagsDecorator struct {
	Tags []string
}

// Deprecated marks a route as deprecated
func Deprecated() func(target interface{}, methodName string) {
	return func(target interface{}, methodName string) {
		RegisterDecorator(target, methodName, &DeprecatedDecorator{})
	}
}

// DeprecatedDecorator marks a route as deprecated
type DeprecatedDecorator struct{}

// Public marks a route as public (not requiring authentication)
func Public() func(target interface{}, methodName string) {
	return func(target interface{}, methodName string) {
		RegisterDecorator(target, methodName, &PublicDecorator{})
	}
}

// PublicDecorator marks a route as public
type PublicDecorator struct{}

// parseRouteMetadata parses all decorators for a controller and returns route metadata
func ParseRouteMetadata(controller interface{}) map[string]*RouteInfo {
	controllerType := reflect.TypeOf(controller)
	if controllerType.Kind() == reflect.Ptr {
		controllerType = controllerType.Elem()
	}

	result := make(map[string]*RouteInfo)

	// Get controller-level decorators
	controllerDecorators := GetDecorators(controller, "")

	// Process each method
	for i := 0; i < controllerType.NumMethod(); i++ {
		method := controllerType.Method(i)
		methodName := method.Name

		// Get method decorators
		methodDecorators := GetDecorators(controller, methodName)
		if len(methodDecorators) == 0 {
			continue // Skip methods without decorators
		}

		// Initialize route info
		routeInfo := &RouteInfo{
			HandlerName: methodName,
			Handler:     method.Func.Interface(),
			Middleware:  make([]gin.HandlerFunc, 0),
			Guards:      make([]core.Guard, 0),
			Tags:        make([]string, 0),
			Responses:   make(map[int]ResponseInfo),
		}

		// Process controller-level decorators first
		for _, decorator := range controllerDecorators {
			applyDecoratorToRouteInfo(decorator, routeInfo)
		}

		// Then process method-level decorators
		for _, decorator := range methodDecorators {
			applyDecoratorToRouteInfo(decorator, routeInfo)
		}

		// Only add if it has a valid HTTP method decorator
		if routeInfo.Method != "" {
			result[methodName] = routeInfo
		}
	}

	return result
}

// applyDecoratorToRouteInfo applies a decorator to route info
func applyDecoratorToRouteInfo(decorator interface{}, routeInfo *RouteInfo) {
	switch d := decorator.(type) {
	case *RouteDecorator:
		routeInfo.Method = d.Method
		routeInfo.Path = d.Path
		routeInfo.Description = d.Description
		routeInfo.Tags = append(routeInfo.Tags, d.Tags...)
		routeInfo.IsPublic = d.IsPublic
		routeInfo.Deprecated = d.Deprecated
	case *GuardDecorator:
		routeInfo.Guards = append(routeInfo.Guards, d.Guards...)
	case *MiddlewareDecorator:
		routeInfo.Middleware = append(routeInfo.Middleware, d.Middleware...)
	case *DescriptionDecorator:
		routeInfo.Description = d.Description
	case *TagsDecorator:
		routeInfo.Tags = append(routeInfo.Tags, d.Tags...)
	case *DeprecatedDecorator:
		routeInfo.Deprecated = true
	case *PublicDecorator:
		routeInfo.IsPublic = true
	}
}

// ConvertToRoutes converts route info to core.Route
func ConvertToRoutes(controllerPath string, routesInfo map[string]*RouteInfo) []core.Route {
	result := make([]core.Route, 0, len(routesInfo))

	for _, info := range routesInfo {
		fullPath := controllerPath + info.Path

		// Convert the handler to gin.HandlerFunc
		// In a real implementation, you'd need to adapt the handler signature
		handler := func(c *gin.Context) {
			// This is a simplified implementation
			// You'd need to call the actual handler with the correct arguments
			fmt.Printf("Handler for %s %s called\n", info.Method, info.Path)
		}

		result = append(result, core.Route{
			Method:      info.Method,
			Path:        fullPath,
			Handler:     handler,
			Middlewares: info.Middleware,
		})
	}

	return result
}

// Helper function for backward compatibility
func registerRouteWithController(method, path string, target interface{}, methodName string) {
	t := reflect.TypeOf(target)
	metadata, ok := controllerRegistry[t]
	if !ok {
		metadata = &core.ControllerMetadata{
			Routes:     make([]core.RouteMetadata, 0),
			Guards:     make([]core.Guard, 0),
			Middleware: make([]gin.HandlerFunc, 0),
		}
		controllerRegistry[t] = metadata
	}

	// Get method from target
	m, ok := t.MethodByName(methodName)
	if !ok {
		panic(fmt.Sprintf("Method %s not found in controller %s", methodName, t.Name()))
	}

	// Add route to metadata
	metadata.Routes = append(metadata.Routes, core.RouteMetadata{
		Path:    path,
		Method:  method,
		Handler: m.Func.Interface(),
	})
}
