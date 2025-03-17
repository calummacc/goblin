package decorators

import (
	"fmt"
	"reflect"
)

// WsMetadata stores metadata for websocket
type WsMetadata struct{
	EventName string
}

// SubscribeMessage decorator
func SubscribeMessage(eventName string) func(interface{}) {
	return func(handler interface{}) {
		SetMetadata(handler, WsMetadataKey, WsMetadata{
			EventName: eventName,
		})
	}
}

func ExtractWsMetadata(target interface{}, methodName string) WsMetadata {
	methodType := reflect.TypeOf(target)
	method, ok := methodType.MethodByName(methodName)
	if !ok {
		panic(fmt.Sprintf("Method %s not found on type %s", methodName, methodType))
	}
	methodInterface := method.Func.Interface()

	if metadata, ok := GetMetadata(methodInterface, WsMetadataKey); ok {
		return metadata.(WsMetadata)
	}

	return WsMetadata{}
}
