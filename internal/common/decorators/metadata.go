package decorators

import (
	"fmt"
	"reflect"
)

// MetadataKey type for unique metadata keys
type MetadataKey string

// Define Metadata Keys
const (
	ControllerMetadataKey MetadataKey = "controller:metadata"
	RouteMetadataKey      MetadataKey = "route:metadata"
	MiddlewareMetadataKey MetadataKey = "middleware:metadata"
	ParamMetadataKey      MetadataKey = "param:metadata"
	InjectableMetadataKey MetadataKey = "injectable:metadata"
	InjectMetadataKey     MetadataKey = "inject:metadata"
	LifecycleMetadataKey  MetadataKey = "lifecycle:metadata"
	WsMetadataKey         MetadataKey = "ws:metadata"
)

// MetadataStore stores metadata for different types
var MetadataStore = make(map[interface{}]map[MetadataKey]interface{})

// SetMetadata stores metadata for a target
func SetMetadata(target interface{}, key MetadataKey, value interface{}) {
	if _, ok := MetadataStore[target]; !ok {
		MetadataStore[target] = make(map[MetadataKey]interface{})
	}
	MetadataStore[target][key] = value
}

// GetMetadata retrieves metadata for a target
func GetMetadata(target interface{}, key MetadataKey) (interface{}, bool) {
	metadata, ok := MetadataStore[target]
	if !ok {
		return nil, false
	}
	value, ok := metadata[key]
	return value, ok
}

func SetMethodMetadata(target interface{}, methodName string, key MetadataKey, value interface{}) {
	targetType := reflect.TypeOf(target)
	method, ok := targetType.MethodByName(methodName)
	if !ok {
		panic(fmt.Sprintf("Method %s not found on type %s", methodName, targetType))
	}
	methodInterface := method.Func.Interface()
	SetMetadata(methodInterface, key, value)
}

func GetMethodMetadata(target interface{}, methodName string, key MetadataKey) (interface{}, bool) {
	targetType := reflect.TypeOf(target)
	method, ok := targetType.MethodByName(methodName)
	if !ok {
		panic(fmt.Sprintf("Method %s not found on type %s", methodName, targetType))
	}
	methodInterface := method.Func.Interface()
	return GetMetadata(methodInterface, key)
}

// GetControllerRegistry return Controller metadata
func GetControllerRegistry() map[interface{}]ControllerMetadata {
	controllerRegistry := make(map[interface{}]ControllerMetadata)
	for target, metadataMap := range MetadataStore {
		if controllerMetadata, ok := metadataMap[ControllerMetadataKey].(ControllerMetadata); ok {
			controllerRegistry[target] = controllerMetadata
		}
	}
	return controllerRegistry
}

// GetRouteRegistry return Route metadata
func GetRouteRegistry() map[interface{}]RouteMetadata {
	routeRegistry := make(map[interface{}]RouteMetadata)
	for target, metadataMap := range MetadataStore {
		if routeMetadata, ok := metadataMap[RouteMetadataKey].(RouteMetadata); ok {
			routeRegistry[target] = routeMetadata
		}
	}
	return routeRegistry
}

// GetInjectableRegistry return Injectable metadata
func GetInjectableRegistry() []interface{} {
	var injectableRegistry []interface{}
	for target, metadataMap := range MetadataStore {
		if _, ok := metadataMap[InjectableMetadataKey]; ok {
			injectableRegistry = append(injectableRegistry, target)
		}
	}
	return injectableRegistry
}

func GetInjectRegistry() map[string]interface{} {
	var injectRegistry map[string]interface{} = make(map[string]interface{})

	for target, _ := range MetadataStore {
		targetType := reflect.TypeOf(target)
		for i := 0; i < targetType.NumMethod(); i++ {
			method := targetType.Method(i)
			metadata, ok := GetMethodMetadata(target, method.Name, InjectMetadataKey)

			if ok {
				injectRegistry[method.Name] = metadata
			}
		}
	}
	return injectRegistry
}

func GetLifecycleRegistry() map[string]LifecycleMetadata {
	var lifecycleRegistry map[string]LifecycleMetadata = make(map[string]LifecycleMetadata)

	for target, _ := range MetadataStore {
		targetType := reflect.TypeOf(target)
		for i := 0; i < targetType.NumMethod(); i++ {
			method := targetType.Method(i)
			metadata, ok := GetMethodMetadata(target, method.Name, LifecycleMetadataKey)

			if ok {
				lifecycleRegistry[method.Name] = metadata.(LifecycleMetadata)
			}
		}
	}
	return lifecycleRegistry
}

func GetWsRegistry() map[string]WsMetadata {
	var wsRegistry map[string]WsMetadata = make(map[string]WsMetadata)

	for target, _ := range MetadataStore {
		targetType := reflect.TypeOf(target)
		for i := 0; i < targetType.NumMethod(); i++ {
			method := targetType.Method(i)
			metadata, ok := GetMethodMetadata(target, method.Name, WsMetadataKey)

			if ok {
				wsRegistry[method.Name] = metadata.(WsMetadata)
			}
		}
	}
	return wsRegistry
}
