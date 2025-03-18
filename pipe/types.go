package pipe

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	// ErrInvalidInput chỉ ra rằng dữ liệu đầu vào không hợp lệ
	ErrInvalidInput = errors.New("invalid input data")
	// ErrValidationFailed chỉ ra rằng dữ liệu không vượt qua kiểm tra validation
	ErrValidationFailed = errors.New("validation failed")
	// ErrIncompatibleType chỉ ra rằng dữ liệu không thể chuyển đổi sang kiểu mong muốn
	ErrIncompatibleType = errors.New("incompatible type")
)

// ValidationError đại diện cho một lỗi validation cụ thể
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag,omitempty"`
	Value   string `json:"value,omitempty"`
}

// ValidationErrors là một slice của ValidationError
type ValidationErrors []ValidationError

// Error trả về thông báo lỗi cho ValidationErrors
func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return "validation failed"
	}

	if len(ve) == 1 {
		return fmt.Sprintf("validation failed: %s %s", ve[0].Field, ve[0].Message)
	}

	return fmt.Sprintf("validation failed with %d errors", len(ve))
}

// TransformContext chứa thông tin về ngữ cảnh chuyển đổi
type TransformContext struct {
	// Value là giá trị đầu vào cần chuyển đổi
	Value interface{}
	// Type là kiểu mong muốn của giá trị đầu ra
	Type reflect.Type
	// Metadata chứa thông tin bổ sung cho quá trình chuyển đổi
	Metadata map[string]interface{}
}

// PipeTransform là interface cốt lõi cho tất cả pipes
type PipeTransform interface {
	// Transform chuyển đổi dữ liệu đầu vào sang định dạng mong muốn
	// hoặc trả về lỗi nếu dữ liệu không hợp lệ
	Transform(ctx *TransformContext) (interface{}, error)
}

// BasePipe cung cấp implementation cơ bản của PipeTransform
type BasePipe struct{}

// Transform implementation mặc định chỉ trả về giá trị đầu vào mà không thay đổi
func (p *BasePipe) Transform(ctx *TransformContext) (interface{}, error) {
	return ctx.Value, nil
}

// PipeOptions chứa các tùy chọn cấu hình cho pipes
type PipeOptions struct {
	// StopOnError dừng chuỗi pipe khi gặp lỗi
	StopOnError bool
}

// DefaultPipeOptions trả về các tùy chọn mặc định
func DefaultPipeOptions() PipeOptions {
	return PipeOptions{
		StopOnError: true,
	}
}

// CompositePipe kết hợp nhiều pipe thành một chuỗi
type CompositePipe struct {
	pipes   []PipeTransform
	options PipeOptions
}

// NewCompositePipe tạo một CompositePipe mới
func NewCompositePipe(options PipeOptions, pipes ...PipeTransform) *CompositePipe {
	return &CompositePipe{
		pipes:   pipes,
		options: options,
	}
}

// Transform chạy tất cả các pipe theo thứ tự và trả về kết quả cuối cùng
func (p *CompositePipe) Transform(ctx *TransformContext) (interface{}, error) {
	value := ctx.Value

	for _, pipe := range p.pipes {
		var err error
		pipeCtx := &TransformContext{
			Value:    value,
			Type:     ctx.Type,
			Metadata: ctx.Metadata,
		}

		value, err = pipe.Transform(pipeCtx)
		if err != nil && p.options.StopOnError {
			return nil, err
		}
	}

	return value, nil
}

// PipeManager quản lý việc đăng ký và sử dụng pipes
type PipeManager struct {
	globalPipes []PipeTransform
}

// NewPipeManager tạo một PipeManager mới
func NewPipeManager() *PipeManager {
	return &PipeManager{
		globalPipes: make([]PipeTransform, 0),
	}
}

// RegisterGlobalPipe đăng ký một pipe toàn cục
func (m *PipeManager) RegisterGlobalPipe(pipe PipeTransform) {
	m.globalPipes = append(m.globalPipes, pipe)
}

// CreateParamDecorator tạo một decorator để áp dụng pipes cho một tham số
func CreateParamDecorator(paramIndex int, pipes ...PipeTransform) func(handler interface{}) interface{} {
	return func(handler interface{}) interface{} {
		handlerType := reflect.TypeOf(handler)
		if handlerType.Kind() != reflect.Func {
			panic("Handler must be a function")
		}

		return func(args ...interface{}) ([]reflect.Value, error) {
			if paramIndex >= len(args) {
				return nil, fmt.Errorf("param index %d out of bounds", paramIndex)
			}

			// Lấy giá trị tham số
			value := args[paramIndex]
			paramType := handlerType.In(paramIndex)

			// Áp dụng pipes
			ctx := &TransformContext{
				Value:    value,
				Type:     paramType,
				Metadata: map[string]interface{}{},
			}

			compositePipe := NewCompositePipe(DefaultPipeOptions(), pipes...)
			result, err := compositePipe.Transform(ctx)
			if err != nil {
				return nil, err
			}

			// Cập nhật giá trị tham số
			args[paramIndex] = result

			// Gọi handler gốc
			handlerValue := reflect.ValueOf(handler)
			inputValues := make([]reflect.Value, len(args))
			for i, arg := range args {
				inputValues[i] = reflect.ValueOf(arg)
			}

			outputValues := handlerValue.Call(inputValues)
			return outputValues, nil
		}
	}
}

// ApplyPipesToHandler áp dụng pipes cho một handler
func ApplyPipesToHandler(handler interface{}, paramPipes map[int][]PipeTransform) interface{} {
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func {
		panic("Handler must be a function")
	}

	return func(args ...interface{}) ([]reflect.Value, error) {
		// Kiểm tra số lượng tham số
		if len(args) != handlerType.NumIn() {
			return nil, fmt.Errorf("expected %d arguments, got %d", handlerType.NumIn(), len(args))
		}

		// Áp dụng pipes cho từng tham số
		for paramIndex, pipes := range paramPipes {
			if paramIndex >= len(args) {
				return nil, fmt.Errorf("param index %d out of bounds", paramIndex)
			}

			// Lấy giá trị tham số
			value := args[paramIndex]
			paramType := handlerType.In(paramIndex)

			// Áp dụng pipes
			ctx := &TransformContext{
				Value:    value,
				Type:     paramType,
				Metadata: map[string]interface{}{},
			}

			compositePipe := NewCompositePipe(DefaultPipeOptions(), pipes...)
			result, err := compositePipe.Transform(ctx)
			if err != nil {
				return nil, err
			}

			// Cập nhật giá trị tham số
			args[paramIndex] = result
		}

		// Gọi handler gốc
		handlerValue := reflect.ValueOf(handler)
		inputValues := make([]reflect.Value, len(args))
		for i, arg := range args {
			if arg == nil {
				// Nếu arg là nil, tạo một zero value cho kiểu tương ứng
				inputValues[i] = reflect.Zero(handlerType.In(i))
			} else {
				inputValues[i] = reflect.ValueOf(arg)
			}
		}

		outputValues := handlerValue.Call(inputValues)
		return outputValues, nil
	}
}

// UsePipes là một decorator để áp dụng pipes cho một method
func UsePipes(paramIndex int, pipes ...PipeTransform) func(target interface{}, methodName string) {
	return func(target interface{}, methodName string) {
		targetType := reflect.TypeOf(target)
		method, ok := targetType.MethodByName(methodName)
		if !ok {
			panic(fmt.Sprintf("Method %s not found in %s", methodName, targetType.Name()))
		}

		// Lưu handler gốc
		originalHandler := method.Func.Interface()

		// Áp dụng pipes
		paramPipes := map[int][]PipeTransform{
			paramIndex: pipes,
		}

		// Tạo handler mới với pipes
		ApplyPipesToHandler(originalHandler, paramPipes)

		// Lưu handler mới (cần triển khai cơ chế để lưu trữ và sử dụng handler mới)
		// Đây chỉ là một ví dụ, cần triển khai chi tiết hơn
		// controllerRegistry.RegisterHandler(target, methodName, newHandler)
	}
}
