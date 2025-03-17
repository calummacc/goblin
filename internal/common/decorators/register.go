package decorators

import (
	"fmt"
	"reflect"
)

// RegisterController registers a controller and its routes
func RegisterController(controller interface{}) {
	controllerMetadata := ExtractControllerMetadata(controller)
	for _, route := range controllerMetadata.Routes {
		// fmt.Println("Register route:", controllerMetadata.Prefix+route.Path)
		SetMethodMetadata(controller, getMethodName(route.Handler), RouteMetadataKey, route)
	}
	// Auto call RegisterDependency and RegisterLifecycleHook
	RegisterDependency(controller)
	RegisterLifecycleHook()
}

// RegisterDependency registers dependencies using decorators
func RegisterDependency(controller interface{}) {
	controllerType := reflect.TypeOf(controller)
	for i := 0; i < controllerType.NumMethod(); i++ {
		method := controllerType.Method(i)
		metadata := ExtractInjectMetadata(controller, method.Name)
		if metadata != nil {
			// Check if the method has a matching dependency
			if metadata == "optional" {
				// Handle optional dependency
				fmt.Println("Optional Dependency:", method.Name)
				continue
			}

			// Handle the required dependency.
			if reflect.TypeOf(metadata).Kind() == reflect.Ptr {
				methodValue := reflect.ValueOf(controller).Method(i)
				if !methodValue.IsValid() || !methodValue.CanSet() {
					panic("Invalid method for dependency injection")
				}
				injectableType := reflect.TypeOf(metadata)
				var injectableValue reflect.Value
				for _, injectable := range GetInjectableRegistry() {
					if reflect.TypeOf(injectable) == injectableType {
						injectableValue = reflect.ValueOf(injectable)
						break
					}
				}
				// Call the method with the injected dependency
				methodValue.Call([]reflect.Value{injectableValue})
			}
		}
	}
}

func RegisterInjectable(injectable interface{}) {
	SetMetadata(injectable, InjectableMetadataKey, struct{}{})
}

// RegisterLifecycleHook calls lifecycle hooks
func RegisterLifecycleHook() {
	lifecycleRegistry := GetLifecycleRegistry()
	for methodName, metadata := range lifecycleRegistry {
		switch metadata.HookType {
		case OnModuleInitHook:
			fmt.Println("OnModuleInit", methodName)
			// Call function
			for target, _ := range MetadataStore {
				targetType := reflect.TypeOf(target)
				if method, ok := targetType.MethodByName(methodName); ok{
					methodValue := method.Func
					methodValue.Call(nil)
				}
			}
		case OnModuleDestroyHook:
			fmt.Println("OnModuleDestroy", methodName)
		case OnApplicationBootstrapHook:
			fmt.Println("OnApplicationBootstrap", methodName)
		case OnApplicationShutdownHook:
			fmt.Println("OnApplicationShutdown", methodName)
		}
	}
}

// RegisterWsHandler call ws handler
func RegisterWsHandler(conn interface{}, message interface{}, handler interface{}) {
	handlerType := reflect.TypeOf(handler)

	for i := 0; i < handlerType.NumMethod(); i++ {
		method := handlerType.Method(i)
		metadata := ExtractWsMetadata(handler, method.Name)
		if metadata.EventName != "" {
			methodValue := reflect.ValueOf(handler).Method(i)
			paramMetadata := ExtractParameterMetadata(handler, method.Name)
			requestParams := ExtractRequestParameters(nil, paramMetadata, conn, message)

			if !methodValue.IsValid() || !methodValue.CanSet() {
				panic("Invalid method for ws handler")
			}

			var args []reflect.Value
			for _, p := range requestParams {
				args = append(args, reflect.ValueOf(p))
			}
			methodValue.Call(args)
		}
	}
}
