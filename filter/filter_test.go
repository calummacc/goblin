package filter

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestSimpleError tests that basic errors are handled
func TestSimpleError(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router
	r := gin.New()

	// Create filter manager
	manager := NewExceptionFilterManager()

	// Apply middleware
	r.Use(manager.GetExceptionHandlerMiddleware())

	// Add a route that throws an error
	r.GET("/error", func(c *gin.Context) {
		err := errors.New("test error")
		_ = c.Error(err)
	})

	// Test the route
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/error", nil)
	r.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Internal Server Error")
}

// TestHttpException tests handling of HttpException
func TestHttpException(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router
	r := gin.New()

	// Create filter manager
	manager := NewExceptionFilterManager()

	// Apply middleware
	r.Use(manager.GetExceptionHandlerMiddleware())

	// Add a route that throws an HttpException
	r.GET("/forbidden", func(c *gin.Context) {
		panic(NewHttpException(http.StatusForbidden, "You don't have permission to access this resource"))
	})

	// Test the route
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/forbidden", nil)
	r.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "You don't have permission")
	assert.Contains(t, w.Body.String(), "Forbidden")
}

// CustomException is a test-specific exception
type CustomException struct {
	Code    string
	Message string
}

func (e *CustomException) Error() string {
	return e.Message
}

// CustomExceptionFilter handles CustomException
type CustomExceptionFilter struct {
	BaseExceptionFilter
}

func (f *CustomExceptionFilter) Catch(exception error, ctx *ExceptionContext) {
	if customErr, ok := exception.(*CustomException); ok {
		ctx.GinContext.JSON(http.StatusBadRequest, gin.H{
			"code":    customErr.Code,
			"message": customErr.Message,
			"custom":  true,
		})
	}
}

func (f *CustomExceptionFilter) CanHandle(exception error) bool {
	_, ok := exception.(*CustomException)
	return ok
}

// TestCustomExceptionFilter tests custom exception filters
func TestCustomExceptionFilter(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router
	r := gin.New()

	// Create filter manager
	manager := NewExceptionFilterManager()

	// Register custom filter
	manager.RegisterGlobalFilter(&CustomExceptionFilter{})

	// Apply middleware
	r.Use(manager.GetExceptionHandlerMiddleware())

	// Add a route that throws a CustomException
	r.GET("/custom-error", func(c *gin.Context) {
		panic(&CustomException{
			Code:    "CUSTOM_ERROR",
			Message: "This is a custom error",
		})
	})

	// Test the route
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/custom-error", nil)
	r.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "This is a custom error")
	assert.Contains(t, w.Body.String(), "CUSTOM_ERROR")
	assert.Contains(t, w.Body.String(), "true") // custom field
}

// TestFilterPriority tests that controller-specific filters take precedence over global ones
func TestFilterPriority(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create filter manager
	manager := NewExceptionFilterManager()

	// Create a controller
	controller := &struct{}{}

	// Create test filters
	globalFilter := &testFilter{name: "global"}
	controllerFilter := &testFilter{name: "controller"}

	// Register filters
	manager.RegisterGlobalFilter(globalFilter)
	manager.RegisterFilter(controller, controllerFilter)

	// Create filter context
	ginCtx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx := &ExceptionContext{
		GinContext: ginCtx,
		Host:       controller,
		Exception:  errors.New("test error"),
	}

	// Test filter execution
	manager.handleException(ctx.Exception, ginCtx, controller, nil)

	// Assert that controller filter was called, not global filter
	assert.True(t, controllerFilter.called)
	assert.False(t, globalFilter.called)
}

// testFilter is a helper for testing filter priority
type testFilter struct {
	BaseExceptionFilter
	name   string
	called bool
}

func (f *testFilter) Catch(exception error, ctx *ExceptionContext) {
	f.called = true
	ctx.GinContext.JSON(http.StatusTeapot, gin.H{"filter": f.name})
}

func (f *testFilter) CanHandle(exception error) bool {
	return true
}
