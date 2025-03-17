package decorators

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
)

// ParamType represents the type of parameter decorator
type ParamType string

// Parameter Decorator Types
const (
	ReqParam     ParamType = "req"
	ResParam     ParamType = "res"
	BodyParam    ParamType = "body"
	QueryParam   ParamType = "query"
	PathParam    ParamType = "path"
	HeadersParam ParamType = "headers"
	SessionParam ParamType = "session"
	IpParam      ParamType = "ip"
	MessageBodyParam ParamType = "messageBody"
	ConnectedSocketParam ParamType = "connectedSocket"
)

// ParamMetadata stores metadata for parameter decorators
type ParamMetadata struct {
	ParamType ParamType
	Index     int
	Data      string // Optional data for @Query, @Param, etc.
}

// ParamDecorator base
func ParamDecorator(paramType ParamType, data string) func(interface{}, string, int) {
	return func(target interface{}, methodName string, index int) {
		SetMethodMetadata(target, methodName, ParamMetadataKey, ParamMetadata{
			ParamType: paramType,
			Index:     index,
			Data:      data,
		})
	}
}

func Req() func(interface{}, string, int) {
	return ParamDecorator(ReqParam, "")
}

func Res() func(interface{}, string, int) {
	return ParamDecorator(ResParam, "")
}

func Body(data string) func(interface{}, string, int) {
	return ParamDecorator(BodyParam, data)
}

func Query(data string) func(interface{}, string, int) {
	return ParamDecorator(QueryParam, data)
}

func Param(data string) func(interface{}, string, int) {
	return ParamDecorator(PathParam, data)
}

func Headers(data string) func(interface{}, string, int) {
	return ParamDecorator(HeadersParam, data)
}

func Session() func(interface{}, string, int) {
	return ParamDecorator(SessionParam, "")
}

func Ip() func(interface{}, string, int) {
	return ParamDecorator(IpParam, "")
}

func MessageBody() func(interface{}, string, int){
	return ParamDecorator(MessageBodyParam, "")
}

func ConnectedSocket() func(interface{}, string, int){
	return ParamDecorator(ConnectedSocketParam, "")
}

func ExtractParameterMetadata(target interface{}, methodName string) []ParamMetadata {
	methodType := reflect.TypeOf(target)
	method, ok := methodType.MethodByName(methodName)
	if !ok {
		panic(fmt.Sprintf("Method %s not found on type %s", methodName, methodType))
	}
	methodInterface := method.Func.Interface()

	metadata, ok := GetMetadata(methodInterface, ParamMetadataKey)
	if !ok {
		return []ParamMetadata{}
	}
	var params []ParamMetadata
	if paramValue, ok := metadata.(ParamMetadata); ok {
		params = append(params, paramValue)
	} else if paramsValue, ok := metadata.([]ParamMetadata); ok {
		params = paramsValue
	} else {
		panic("Invalid ParamMetadata")
	}

	return params
}

func ExtractRequestParameters(c *gin.Context, params []ParamMetadata, conn interface{}, message interface{}) []interface{} {
	var requestParams []interface{}
	for _, param := range params {
		switch param.ParamType {
		case ReqParam:
			requestParams = append(requestParams, c.Request)
		case ResParam:
			requestParams = append(requestParams, c.Writer)
		case BodyParam:
			if param.Data != "" {
				// Bind body to data type
				paramType := reflect.TypeOf(c.Handler())
				if paramType.Kind() == reflect.Func {
					for i := 0; i < paramType.NumIn(); i++ {
						fieldType := paramType.In(i)
						if fieldType.String() == param.Data {
							v := reflect.New(fieldType).Interface()
							if err := c.ShouldBindBodyWith(&v, nil); err != nil {
								panic(err)
							}
							requestParams = append(requestParams, reflect.ValueOf(v).Elem().Interface())
						}
					}
				}
			} else {
				//Bind body to map[string]interface{}
				var v map[string]interface{}
				if err := c.ShouldBindBodyWith(&v, nil); err != nil {
					panic(err)
				}
				requestParams = append(requestParams, v)
			}
		case QueryParam:
			requestParams = append(requestParams, c.Query(param.Data))
		case PathParam:
			requestParams = append(requestParams, c.Param(param.Data))
		case HeadersParam:
			requestParams = append(requestParams, c.GetHeader(param.Data))
		case SessionParam:
			// get session
		case IpParam:
			requestParams = append(requestParams, c.ClientIP())
		case MessageBodyParam:
			requestParams = append(requestParams, message)
		case ConnectedSocketParam:
			requestParams = append(requestParams, conn)
		default:
			panic("Invalid ParamType")
		}
	}
	return requestParams
}
