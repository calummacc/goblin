package decorators

import (
	"fmt"
	"reflect"

	"net/http"

	"github.com/gin-gonic/gin"
)

// RouteHandler is the handler function for a route.
type RouteHandler func(*gin.Context)

// ControllerMetadata stores metadata about a controller.
type ControllerMetadata struct {
	Prefix      string
	Middlewares []interface{} // Middleware can be functions, types, etc.
	Routes      []RouteMetadata
}

// RouteMetadata stores metadata about a route.
type RouteMetadata struct {
	Path        string
	Method      string
	Handler     interface{}
	Middlewares []interface{} // Local route middleware
}

// Controller decorator
func Controller(prefix string, middlewares ...interface{}) func(interface{}) {
	return func(target interface{}) {
		SetMetadata(target, ControllerMetadataKey, ControllerMetadata{
			Prefix:      prefix,
			Middlewares: middlewares,
			Routes:      []RouteMetadata{},
		})
		// Auto register
		RegisterController(target)
	}
}

// Route Decorator
func Route(method, path string, middlewares ...interface{}) func(interface{}) {
	return func(handler interface{}) {
		// Route metadata already added in @Get, @Post, etc.
		if _, ok := GetMetadata(handler, RouteMetadataKey); !ok {
			SetMetadata(handler, RouteMetadataKey, RouteMetadata{
				Path:        path,
				Method:      method,
				Handler:     handler,
				Middlewares: middlewares,
			})
		} else {
			metadata, _ := GetMetadata(handler, RouteMetadataKey)
			routeMetadata := metadata.(RouteMetadata)
			routeMetadata.Middlewares = append(routeMetadata.Middlewares, middlewares...)
			SetMetadata(handler, RouteMetadataKey, routeMetadata)
		}
	}
}

func Get(path string, middlewares ...interface{}) func(interface{}, string) {
	return func(handler interface{}, methodName string) {
		Route(http.MethodGet, path, middlewares...)(getMethod(handler, methodName))
	}
}

func Post(path string, middlewares ...interface{}) func(interface{}, string) {
	return func(handler interface{}, methodName string) {
		Route(http.MethodPost, path, middlewares...)(getMethod(handler, methodName))
	}
}

func Put(path string, middlewares ...interface{}) func(interface{}, string) {
	return func(handler interface{}, methodName string) {
		Route(http.MethodPut, path, middlewares...)(getMethod(handler, methodName))
	}
}

func Delete(path string, middlewares ...interface{}) func(interface{}, string) {
	return func(handler interface{}, methodName string) {
		Route(http.MethodDelete, path, middlewares...)(getMethod(handler, methodName))
	}
}

func Patch(path string, middlewares ...interface{}) func(interface{}, string) {
	return func(handler interface{}, methodName string) {
		Route(http.MethodPatch, path, middlewares...)(getMethod(handler, methodName))
	}
}

func Options(path string, middlewares ...interface{}) func(interface{}, string) {
	return func(handler interface{}, methodName string) {
		Route(http.MethodOptions, path, middlewares...)(getMethod(handler, methodName))
	}
}

func Head(path string, middlewares ...interface{}) func(interface{}, string) {
	return func(handler interface{}, methodName string) {
		Route(http.MethodHead, path, middlewares...)(getMethod(handler, methodName))
	}
}

// ExtractControllerMetadata extracts and sets route metadata.
func ExtractControllerMetadata(controller interface{}) ControllerMetadata {
	metadata, ok := GetMetadata(controller, ControllerMetadataKey)
	if !ok {
		panic("Controller metadata not found")
	}
	controllerMetadata := metadata.(ControllerMetadata)

	t := reflect.TypeOf(controller)
	v := reflect.ValueOf(controller)

	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		methodValue := v.Method(i)
		methodInterface := methodValue.Interface()

		if routeMetadataValue, ok := GetMetadata(methodInterface, RouteMetadataKey); ok {
			routeMetadata := routeMetadataValue.(RouteMetadata)
			routeMetadata.Handler = methodInterface
			controllerMetadata.Routes = append(controllerMetadata.Routes, routeMetadata)
		}
	}

	SetMetadata(controller, ControllerMetadataKey, controllerMetadata)
	return controllerMetadata
}

func getMethod(target interface{}, methodName string) interface{} {
	targetType := reflect.TypeOf(target)
	method, ok := targetType.MethodByName(methodName)
	if !ok {
		panic(fmt.Sprintf("Method %s not found on type %s", methodName, targetType))
	}
	methodInterface := method.Func.Interface()
	return methodInterface
}
