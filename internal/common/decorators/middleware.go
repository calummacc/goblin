package decorators

import (
	"fmt"
	"reflect"
)

// UseInterceptors decorator
func UseInterceptors(interceptors ...interface{}) func(interface{}) {
	return func(target interface{}) {
		SetMetadata(target, MiddlewareMetadataKey, struct {
			Interceptors []interface{}
		}{
			Interceptors: interceptors,
		})
	}
}

// UseGuards decorator
func UseGuards(guards ...interface{}) func(interface{}) {
	return func(target interface{}) {
		SetMetadata(target, MiddlewareMetadataKey, struct {
			Guards []interface{}
		}{
			Guards: guards,
		})
	}
}

// UsePipes decorator
func UsePipes(pipes ...interface{}) func(interface{}) {
	return func(target interface{}) {
		SetMetadata(target, MiddlewareMetadataKey, struct {
			Pipes []interface{}
		}{
			Pipes: pipes,
		})
	}
}

// UseFilters decorator
func UseFilters(filters ...interface{}) func(interface{}) {
	return func(target interface{}) {
		SetMetadata(target, MiddlewareMetadataKey, struct {
			Filters []interface{}
		}{
			Filters: filters,
		})
	}
}

// ExtractMiddlewareMetadata extracts and sets middleware metadata.
func ExtractMiddlewareMetadata(target interface{}, methodName string) map[string][]interface{} {
	middlewareMetadata := make(map[string][]interface{})
	method, ok := reflect.TypeOf(target).MethodByName(methodName)
	if !ok {
		panic(fmt.Sprintf("Method %s not found on type %s", methodName, reflect.TypeOf(target)))
	}

	methodInterface := method.Func.Interface()

	if metadata, ok := GetMetadata(methodInterface, MiddlewareMetadataKey); ok {
		switch m := metadata.(type) {
		case struct {
			Interceptors []interface{}
		}:
			middlewareMetadata["Interceptors"] = m.Interceptors
		case struct {
			Guards []interface{}
		}:
			middlewareMetadata["Guards"] = m.Guards
		case struct {
			Pipes []interface{}
		}:
			middlewareMetadata["Pipes"] = m.Pipes
		case struct {
			Filters []interface{}
		}:
			middlewareMetadata["Filters"] = m.Filters
		}
	}
	return middlewareMetadata
}
