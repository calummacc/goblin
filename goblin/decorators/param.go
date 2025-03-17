// goblin/decorators/param.go
package decorators

import (
	"github.com/gin-gonic/gin"
)

// ParamDecoratorFunc is a function that extracts a parameter from a request
type ParamDecoratorFunc func(c *gin.Context, name string, into interface{}) error

// Body extracts a parameter from the request body
func Body(name string) ParamDecoratorFunc {
	return func(c *gin.Context, _ string, into interface{}) error {
		return c.ShouldBindJSON(into)
	}
}

// Query extracts a parameter from the query string
func Query(name string) ParamDecoratorFunc {
	return func(c *gin.Context, paramName string, into interface{}) error {
		value := c.Query(paramName)
		if strValue, ok := into.(*string); ok {
			*strValue = value
			return nil
		}
		return c.ShouldBindQuery(into)
	}
}

// Param extracts a parameter from the URL path
func Param(name string) ParamDecoratorFunc {
	return func(c *gin.Context, paramName string, into interface{}) error {
		value := c.Param(paramName)
		if strValue, ok := into.(*string); ok {
			*strValue = value
			return nil
		}
		// For other types, we would need to implement type conversion
		// This is simplified for now
		return nil
	}
}

// Header extracts a parameter from the request headers
func Header(name string) ParamDecoratorFunc {
	return func(c *gin.Context, paramName string, into interface{}) error {
		value := c.GetHeader(paramName)
		if strValue, ok := into.(*string); ok {
			*strValue = value
			return nil
		}
		return nil
	}
}

// ExtractParam is a helper function to extract a parameter using a decorator
func ExtractParam(c *gin.Context, decorator ParamDecoratorFunc, name string, into interface{}) error {
	return decorator(c, name, into)
}
