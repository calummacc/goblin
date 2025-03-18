package pipe

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidationPipe là một pipe kiểm tra tính hợp lệ của dữ liệu
type ValidationPipe struct {
	BasePipe
	validator *validator.Validate
	options   ValidationOptions
}

// ValidationOptions chứa các tùy chọn cấu hình cho ValidationPipe
type ValidationOptions struct {
	// DisableErrorMessages ngăn chặn hiển thị thông báo lỗi chi tiết
	DisableErrorMessages bool
	// ValidateCustomDecorators bật kiểm tra decorator tùy chỉnh
	ValidateCustomDecorators bool
	// ForbidUnknownValues từ chối các trường không được định nghĩa trong DTO
	ForbidUnknownValues bool
	// Transform tự động chuyển đổi dữ liệu sang kiểu mong muốn
	Transform bool
	// StopAtFirstError dừng kiểm tra khi gặp lỗi đầu tiên
	StopAtFirstError bool
}

// DefaultValidationOptions trả về các tùy chọn mặc định
func DefaultValidationOptions() ValidationOptions {
	return ValidationOptions{
		DisableErrorMessages:     false,
		ValidateCustomDecorators: true,
		ForbidUnknownValues:      false,
		Transform:                true,
		StopAtFirstError:         false,
	}
}

// NewValidationPipe tạo một ValidationPipe mới
func NewValidationPipe(options ...ValidationOptions) *ValidationPipe {
	var opts ValidationOptions
	if len(options) > 0 {
		opts = options[0]
	} else {
		opts = DefaultValidationOptions()
	}

	return &ValidationPipe{
		validator: validator.New(),
		options:   opts,
	}
}

// Transform thực hiện validate dữ liệu đầu vào
func (p *ValidationPipe) Transform(ctx *TransformContext) (interface{}, error) {
	if ctx.Value == nil {
		return nil, ErrInvalidInput
	}

	// Nếu dữ liệu đã có kiểu đúng, chỉ cần validate
	valueType := reflect.TypeOf(ctx.Value)
	if valueType == ctx.Type {
		return p.validateValue(ctx.Value)
	}

	// Nếu cần chuyển đổi kiểu dữ liệu
	if p.options.Transform {
		// Tạo instance mới của kiểu mong muốn
		targetValue := reflect.New(ctx.Type).Interface()

		// Thực hiện chuyển đổi (trong trường hợp thực tế, có thể sử dụng thư viện như mapstructure)
		// Đơn giản hóa: giả sử dữ liệu đầu vào là map[string]interface{}
		if err := p.mapToStruct(ctx.Value, targetValue); err != nil {
			return nil, err
		}

		// Validate dữ liệu sau khi chuyển đổi
		return p.validateValue(targetValue)
	}

	return nil, fmt.Errorf("%w: expected %s but got %s", ErrIncompatibleType, ctx.Type, valueType)
}

// validateValue thực hiện kiểm tra tính hợp lệ của một giá trị
func (p *ValidationPipe) validateValue(value interface{}) (interface{}, error) {
	// Thực hiện validate bằng validator
	if err := p.validator.Struct(value); err != nil {
		// Xử lý lỗi validation
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return nil, p.formatValidationErrors(validationErrors)
		}
		return nil, err
	}

	return value, nil
}

// formatValidationErrors chuyển đổi lỗi từ validator sang ValidationErrors
func (p *ValidationPipe) formatValidationErrors(validationErrors validator.ValidationErrors) error {
	if p.options.DisableErrorMessages {
		return ErrValidationFailed
	}

	var errors ValidationErrors
	for _, err := range validationErrors {
		errors = append(errors, ValidationError{
			Field:   p.formatFieldName(err.Field()),
			Message: p.getErrorMessage(err),
			Tag:     err.Tag(),
			Value:   fmt.Sprintf("%v", err.Value()),
		})

		if p.options.StopAtFirstError {
			break
		}
	}

	return errors
}

// formatFieldName định dạng tên trường
func (p *ValidationPipe) formatFieldName(field string) string {
	return strings.ToLower(field[:1]) + field[1:]
}

// getErrorMessage tạo thông báo lỗi cho một lỗi validation
func (p *ValidationPipe) getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	case "min":
		return fmt.Sprintf("must be at least %s", err.Param())
	case "max":
		return fmt.Sprintf("must be at most %s", err.Param())
	case "len":
		return fmt.Sprintf("must be %s characters long", err.Param())
	case "eq":
		return fmt.Sprintf("must be equal to %s", err.Param())
	case "ne":
		return fmt.Sprintf("must not be equal to %s", err.Param())
	case "oneof":
		return fmt.Sprintf("must be one of %s", err.Param())
	}
	return fmt.Sprintf("failed on '%s' validation", err.Tag())
}

// mapToStruct ánh xạ dữ liệu từ map sang struct
// Đây là một triển khai đơn giản cho ví dụ, trong thực tế có thể cần phức tạp hơn
func (p *ValidationPipe) mapToStruct(input, output interface{}) error {
	inputValue := reflect.ValueOf(input)
	outputValue := reflect.ValueOf(output)

	// Đảm bảo output là một con trỏ
	if outputValue.Kind() != reflect.Ptr {
		return fmt.Errorf("output must be a pointer")
	}

	// Đảm bảo output trỏ đến một struct
	outputElem := outputValue.Elem()
	if outputElem.Kind() != reflect.Struct {
		return fmt.Errorf("output must be a pointer to a struct")
	}

	// Nếu input là map
	if inputValue.Kind() == reflect.Map {
		return p.mapToStructFromMap(inputValue, outputElem)
	}

	// Nếu input là struct
	if inputValue.Kind() == reflect.Struct {
		return p.mapToStructFromStruct(inputValue, outputElem)
	}

	return fmt.Errorf("input must be a map or a struct")
}

// mapToStructFromMap ánh xạ dữ liệu từ map sang struct
func (p *ValidationPipe) mapToStructFromMap(inputMap, outputStruct reflect.Value) error {
	// Duyệt qua các trường của struct đích
	for i := 0; i < outputStruct.NumField(); i++ {
		field := outputStruct.Type().Field(i)
		fieldValue := outputStruct.Field(i)

		// Lấy tên trường trong JSON tag hoặc tên trường
		jsonTag := field.Tag.Get("json")
		name := field.Name
		if jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "-" {
				name = parts[0]
			}
		}

		// Lấy giá trị từ map
		mapValue := inputMap.MapIndex(reflect.ValueOf(name))
		if !mapValue.IsValid() {
			continue // Không tìm thấy giá trị trong map
		}

		// Thử thiết lập giá trị
		if fieldValue.CanSet() {
			// Chuyển đổi kiểu dữ liệu nếu cần
			if mapValue.Type().AssignableTo(fieldValue.Type()) {
				fieldValue.Set(mapValue)
			} else if mapValue.Type().ConvertibleTo(fieldValue.Type()) {
				fieldValue.Set(mapValue.Convert(fieldValue.Type()))
			} else {
				return fmt.Errorf("cannot assign %s to field %s of type %s", mapValue.Type(), field.Name, fieldValue.Type())
			}
		}
	}

	return nil
}

// mapToStructFromStruct ánh xạ dữ liệu từ struct sang struct
func (p *ValidationPipe) mapToStructFromStruct(inputStruct, outputStruct reflect.Value) error {
	// Duyệt qua các trường của struct đích
	for i := 0; i < outputStruct.NumField(); i++ {
		outputField := outputStruct.Type().Field(i)
		outputFieldValue := outputStruct.Field(i)

		// Tìm trường tương ứng trong struct nguồn
		inputFieldValue := inputStruct.FieldByName(outputField.Name)
		if !inputFieldValue.IsValid() {
			continue // Không tìm thấy trường trong struct nguồn
		}

		// Thử thiết lập giá trị
		if outputFieldValue.CanSet() {
			// Chuyển đổi kiểu dữ liệu nếu cần
			if inputFieldValue.Type().AssignableTo(outputFieldValue.Type()) {
				outputFieldValue.Set(inputFieldValue)
			} else if inputFieldValue.Type().ConvertibleTo(outputFieldValue.Type()) {
				outputFieldValue.Set(inputFieldValue.Convert(outputFieldValue.Type()))
			} else {
				return fmt.Errorf("cannot assign %s to field %s of type %s", inputFieldValue.Type(), outputField.Name, outputFieldValue.Type())
			}
		}
	}

	return nil
}
