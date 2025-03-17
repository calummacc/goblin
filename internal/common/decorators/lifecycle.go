package decorators

import (
	"fmt"
	"reflect"
)

// LifecycleHookType represents the type of lifecycle hook
type LifecycleHookType string

// Lifecycle Hook Types
const (
	OnModuleInitHook           LifecycleHookType = "onModuleInit"
	OnModuleDestroyHook        LifecycleHookType = "onModuleDestroy"
	OnApplicationBootstrapHook LifecycleHookType = "onApplicationBootstrap"
	OnApplicationShutdownHook  LifecycleHookType = "onApplicationShutdown"
)

// LifecycleMetadata stores metadata for lifecycle hooks
type LifecycleMetadata struct {
	HookType LifecycleHookType
}

// LifecycleHookDecorator base
func LifecycleHookDecorator(hookType LifecycleHookType) func(interface{}) {
	return func(handler interface{}) {
		SetMetadata(handler, LifecycleMetadataKey, LifecycleMetadata{
			HookType: hookType,
		})
	}
}

func OnModuleInit() func(interface{}) {
	return LifecycleHookDecorator(OnModuleInitHook)
}

func OnModuleDestroy() func(interface{}) {
	return LifecycleHookDecorator(OnModuleDestroyHook)
}

func OnApplicationBootstrap() func(interface{}) {
	return LifecycleHookDecorator(OnApplicationBootstrapHook)
}

func OnApplicationShutdown() func(interface{}) {
	return LifecycleHookDecorator(OnApplicationShutdownHook)
}

func ExtractLifecycleMetadata(target interface{}, methodName string) LifecycleMetadata {
	methodType := reflect.TypeOf(target)
	method, ok := methodType.MethodByName(methodName)
	if !ok {
		panic(fmt.Sprintf("Method %s not found on type %s", methodName, methodType))
	}
	methodInterface := method.Func.Interface()

	if metadata, ok := GetMetadata(methodInterface, LifecycleMetadataKey); ok {
		return metadata.(LifecycleMetadata)
	}
	return LifecycleMetadata{}
}
