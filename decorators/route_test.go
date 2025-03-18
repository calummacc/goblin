package decorators

import (
	"goblin/core"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestGuard for testing
type TestGuard struct {
	core.BaseGuard
}

// TestController for testing decorators
type TestController struct {
	core.BaseController
}

// Simple route handlers for testing
func (c *TestController) GetMethod(ctx *gin.Context)     {}
func (c *TestController) PostMethod(ctx *gin.Context)    {}
func (c *TestController) PutMethod(ctx *gin.Context)     {}
func (c *TestController) DeleteMethod(ctx *gin.Context)  {}
func (c *TestController) PatchMethod(ctx *gin.Context)   {}
func (c *TestController) OptionsMethod(ctx *gin.Context) {}
func (c *TestController) HeadMethod(ctx *gin.Context)    {}
func (c *TestController) AnyMethod(ctx *gin.Context)     {}

// TestHTTPMethodDecorators tests the HTTP method decorators
func TestHTTPMethodDecorators(t *testing.T) {
	// Create a controller
	controller := &TestController{}

	// Apply the controller decorator
	Controller("/test")(controller)

	// Apply HTTP method decorators
	Get("/get")(controller, "GetMethod")
	Post("/post")(controller, "PostMethod")
	Put("/put")(controller, "PutMethod")
	Delete("/delete")(controller, "DeleteMethod")
	Patch("/patch")(controller, "PatchMethod")
	Options("/options")(controller, "OptionsMethod")
	Head("/head")(controller, "HeadMethod")
	All("/all")(controller, "AnyMethod")

	// Parse metadata
	routes := ParseRouteMetadata(controller)

	// Verify the routes
	assert.Equal(t, 8, len(routes))

	// Verify each route
	assert.Equal(t, "GET", routes["GetMethod"].Method)
	assert.Equal(t, "/get", routes["GetMethod"].Path)

	assert.Equal(t, "POST", routes["PostMethod"].Method)
	assert.Equal(t, "/post", routes["PostMethod"].Path)

	assert.Equal(t, "PUT", routes["PutMethod"].Method)
	assert.Equal(t, "/put", routes["PutMethod"].Path)

	assert.Equal(t, "DELETE", routes["DeleteMethod"].Method)
	assert.Equal(t, "/delete", routes["DeleteMethod"].Path)

	assert.Equal(t, "PATCH", routes["PatchMethod"].Method)
	assert.Equal(t, "/patch", routes["PatchMethod"].Path)

	assert.Equal(t, "OPTIONS", routes["OptionsMethod"].Method)
	assert.Equal(t, "/options", routes["OptionsMethod"].Path)

	assert.Equal(t, "HEAD", routes["HeadMethod"].Method)
	assert.Equal(t, "/head", routes["HeadMethod"].Path)

	assert.Equal(t, "ANY", routes["AnyMethod"].Method)
	assert.Equal(t, "/all", routes["AnyMethod"].Path)
}

// TestMetadataDecorators tests the metadata decorators
func TestMetadataDecorators(t *testing.T) {
	// Create a controller
	controller := &TestController{}

	// Apply the controller decorator
	Controller("/api")(controller)

	// Apply method decorators with metadata
	Get("/resource")(controller, "GetMethod")

	// Apply description
	descDecorator, ok := Description("Get a resource").(struct {
		ControllerDecorator func(interface{})
		MethodDecorator     func(interface{}, string)
	})
	assert.True(t, ok)
	descDecorator.MethodDecorator(controller, "GetMethod")

	// Apply tags
	tagsDecorator, ok := Tags("api", "resource").(struct {
		ControllerDecorator func(interface{})
		MethodDecorator     func(interface{}, string)
	})
	assert.True(t, ok)
	tagsDecorator.MethodDecorator(controller, "GetMethod")

	// Apply public flag
	Public()(controller, "GetMethod")

	// Apply deprecated flag
	Deprecated()(controller, "GetMethod")

	// Parse metadata
	routes := ParseRouteMetadata(controller)

	// Verify metadata
	routeInfo := routes["GetMethod"]
	assert.Equal(t, "Get a resource", routeInfo.Description)
	assert.Equal(t, []string{"api", "resource"}, routeInfo.Tags)
	assert.True(t, routeInfo.IsPublic)
	assert.True(t, routeInfo.Deprecated)
}

// TestGuardDecorator tests the guard decorators
func TestGuardDecorator(t *testing.T) {
	// Create a controller
	controller := &TestController{}

	// Apply the controller decorator
	Controller("/secure")(controller)

	// Apply method with guard
	Get("/protected")(controller, "GetMethod")
	UseGuards(&TestGuard{})(controller, "GetMethod")

	// Parse metadata
	routes := ParseRouteMetadata(controller)

	// Verify guards
	routeInfo := routes["GetMethod"]
	assert.Equal(t, 1, len(routeInfo.Guards))
	_, ok := routeInfo.Guards[0].(*TestGuard)
	assert.True(t, ok)
}

// TestMiddlewareDecorator tests the middleware decorators
func TestMiddlewareDecorator(t *testing.T) {
	// Create a controller
	controller := &TestController{}

	// Apply the controller decorator
	Controller("/middleware")(controller)

	// Create test middleware
	testMiddleware := func(ctx *gin.Context) {}

	// Apply method with middleware
	Get("/with-middleware")(controller, "GetMethod")
	UseMiddleware(testMiddleware)(controller, "GetMethod")

	// Parse metadata
	routes := ParseRouteMetadata(controller)

	// Verify middleware
	routeInfo := routes["GetMethod"]
	assert.Equal(t, 1, len(routeInfo.Middleware))
}

// TestRouteConversion tests converting route info to core.Route
func TestRouteConversion(t *testing.T) {
	// Create a controller
	controller := &TestController{}

	// Apply decorators
	Controller("/api")(controller)
	Get("/users")(controller, "GetMethod")

	// Parse metadata
	routes := ParseRouteMetadata(controller)

	// Convert to core.Route
	routerRoutes := ConvertToRoutes("/api", routes)

	// Verify conversion
	assert.Equal(t, 1, len(routerRoutes))
	assert.Equal(t, "GET", routerRoutes[0].Method)
	assert.Equal(t, "/api/users", routerRoutes[0].Path)
	assert.NotNil(t, routerRoutes[0].Handler)
}
