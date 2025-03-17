package decorators

import (
	"fmt"
	"reflect"
)

// Injectable decorator
func Injectable() func(interface{}) {
	return func(target interface{}) {
		SetMetadata(target, InjectableMetadataKey, struct{}{})
		// Auto Register
		RegisterInjectable(target)
	}
}

// Inject decorator
func Inject(dependency interface{}) func(interface{}, string) {
	return func(target interface{}, methodName string) {
		SetMethodMetadata(target, methodName, InjectMetadataKey, dependency)
		//Auto register
		RegisterDependency(target)
	}
}

// Optional decorator
func Optional() func(interface{}, string, int) {
	return func(target interface{}, methodName string, index int) {
		targetType := reflect.TypeOf(target)
		method, ok := targetType.MethodByName(methodName)
		if !ok {
			panic(fmt.Sprintf("Method %s not found on type %s", methodName, targetType))
		}

		methodInterface := method.Func.Interface()
		methodValue := reflect.ValueOf(methodInterface)
		for i := 0; i < methodValue.Type().NumIn(); i++ {
			if i == index {
				SetMethodMetadata(target, methodName, InjectMetadataKey, "optional")
			}
		}
	}
}

func ExtractInjectMetadata(target interface{}, methodName string) interface{} {
	methodType := reflect.TypeOf(target)
	method, ok := methodType.MethodByName(methodName)
	if !ok {
		panic(fmt.Sprintf("Method %s not found on type %s", methodName, methodType))
	}
	methodInterface := method.Func.Interface()

	if metadata, ok := GetMetadata(methodInterface, InjectMetadataKey); ok {
		return metadata
	}

	return nil
}
