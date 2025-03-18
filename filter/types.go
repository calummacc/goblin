package filter

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

// HttpException represents an HTTP exception with status code and message
type HttpException struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	ErrorType  string      `json:"error,omitempty"`
	Details    interface{} `json:"details,omitempty"`
}

// Error implements the error interface
func (e *HttpException) Error() string {
	return fmt.Sprintf("%d %s: %s", e.StatusCode, e.ErrorType, e.Message)
}

// NewHttpException creates a new HttpException
func NewHttpException(statusCode int, message string, details ...interface{}) *HttpException {
	var detailsValue interface{}
	if len(details) > 0 {
		detailsValue = details[0]
	}

	return &HttpException{
		StatusCode: statusCode,
		Message:    message,
		ErrorType:  http.StatusText(statusCode),
		Details:    detailsValue,
	}
}

// ExceptionContext provides context for exception handling
type ExceptionContext struct {
	// Original Gin context
	GinContext *gin.Context
	// Exception is the error that was thrown
	Exception error
	// Host is the controller/handler where the exception was thrown
	Host interface{}
	// Handler is the route handler where the exception was thrown
	Handler interface{}
	// Path is the route path
	Path string
	// Method is the HTTP method
	Method string
}

// ExceptionFilter defines the contract for exception filters
type ExceptionFilter interface {
	// Catch handles an exception
	Catch(exception error, ctx *ExceptionContext)
	// CanHandle checks if this filter can handle the given exception
	CanHandle(exception error) bool
}

// LifecycleExceptionFilter extends ExceptionFilter with lifecycle hooks
type LifecycleExceptionFilter interface {
	ExceptionFilter
	// OnRegister is called when the filter is registered
	OnRegister(ctx context.Context) error
	// OnShutdown is called when the application is shutting down
	OnShutdown(ctx context.Context) error
}

// BaseExceptionFilter provides a base implementation of ExceptionFilter
type BaseExceptionFilter struct{}

// Catch provides default implementation that does nothing
func (f *BaseExceptionFilter) Catch(exception error, ctx *ExceptionContext) {
	// Default implementation does nothing
}

// CanHandle returns true if this filter can handle the given exception
func (f *BaseExceptionFilter) CanHandle(exception error) bool {
	// Default implementation handles all exceptions
	return true
}

// ExceptionFilterManager manages exception filters
type ExceptionFilterManager struct {
	globalFilters []ExceptionFilter
	registry      map[reflect.Type][]ExceptionFilter
}

// NewExceptionFilterManager creates a new ExceptionFilterManager
func NewExceptionFilterManager() *ExceptionFilterManager {
	return &ExceptionFilterManager{
		globalFilters: make([]ExceptionFilter, 0),
		registry:      make(map[reflect.Type][]ExceptionFilter),
	}
}

// RegisterGlobalFilter registers a filter to be applied globally
func (m *ExceptionFilterManager) RegisterGlobalFilter(filter ExceptionFilter) {
	m.globalFilters = append(m.globalFilters, filter)
}

// RegisterFilter registers a filter for a specific controller or handler
func (m *ExceptionFilterManager) RegisterFilter(target interface{}, filter ExceptionFilter) {
	targetType := reflect.TypeOf(target)
	if _, ok := m.registry[targetType]; !ok {
		m.registry[targetType] = make([]ExceptionFilter, 0)
	}
	m.registry[targetType] = append(m.registry[targetType], filter)
}

// GetExceptionHandlerMiddleware returns a middleware that handles exceptions
func (m *ExceptionFilterManager) GetExceptionHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Defer error handling
		defer func() {
			if r := recover(); r != nil {
				var err error
				switch v := r.(type) {
				case error:
					err = v
				default:
					err = fmt.Errorf("%v", v)
				}

				m.handleException(err, c, nil, nil)
				// Abort the request
				c.Abort()
			}
		}()

		// Continue with the next middleware/handler
		c.Next()

		// Check for errors after executing the handler
		if len(c.Errors) > 0 {
			// Handle the last error
			err := c.Errors.Last().Err
			m.handleException(err, c, nil, nil)
		}
	}
}

// handleException processes an exception through registered filters
func (m *ExceptionFilterManager) handleException(exception error, c *gin.Context, host, handler interface{}) {
	// Create exception context
	ctx := &ExceptionContext{
		GinContext: c,
		Exception:  exception,
		Host:       host,
		Handler:    handler,
		Path:       c.FullPath(),
		Method:     c.Request.Method,
	}

	// Check if any specific filters can handle this exception
	if host != nil {
		hostType := reflect.TypeOf(host)
		if filters, ok := m.registry[hostType]; ok {
			for _, filter := range filters {
				if filter.CanHandle(exception) {
					filter.Catch(exception, ctx)
					return
				}
			}
		}
	}

	// Try global filters
	for _, filter := range m.globalFilters {
		if filter.CanHandle(exception) {
			filter.Catch(exception, ctx)
			return
		}
	}

	// Default error handling if no filter caught the exception
	defaultErrorHandling(exception, c)
}

// defaultErrorHandling handles errors when no filter is available
func defaultErrorHandling(err error, c *gin.Context) {
	statusCode := http.StatusInternalServerError
	message := "Internal Server Error"

	// Check if it's an HttpException
	if httpErr, ok := err.(*HttpException); ok {
		statusCode = httpErr.StatusCode
		message = httpErr.Message
		c.JSON(statusCode, httpErr)
		return
	}

	// Return generic error
	c.JSON(statusCode, gin.H{
		"statusCode": statusCode,
		"message":    message,
		"error":      http.StatusText(statusCode),
	})
}

// Decorators for use with controllers and handlers

// UseFilters is a decorator to apply exception filters to a controller or method
func UseFilters(filters ...ExceptionFilter) func(target interface{}, methodName string) {
	return func(target interface{}, methodName string) {
		// This would be implemented to register the filters with the filter manager
		// A global reference to the filter manager would be needed here
		// Or a registry to store decorations that are later processed

		// This is a placeholder implementation
		// In a real implementation, you would register these filters with the manager
	}
}

// GlobalFilters is a decorator to apply exception filters globally
func GlobalFilters(filters ...ExceptionFilter) func(module interface{}) {
	return func(module interface{}) {
		// This would be implemented to register global filters with the filter manager
		// A global reference to the filter manager would be needed here
		// Or a registry to store decorations that are later processed

		// This is a placeholder implementation
		// In a real implementation, you would register these filters with the manager
	}
}
