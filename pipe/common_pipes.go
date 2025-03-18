package pipe

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ParseIntPipe chuyển đổi string thành int
type ParseIntPipe struct {
	BasePipe
	radix int
}

// NewParseIntPipe tạo một ParseIntPipe mới
func NewParseIntPipe(radix int) *ParseIntPipe {
	if radix <= 0 {
		radix = 10
	}
	return &ParseIntPipe{
		radix: radix,
	}
}

// Transform chuyển đổi string thành int
func (p *ParseIntPipe) Transform(ctx *TransformContext) (interface{}, error) {
	if ctx.Value == nil {
		return nil, ErrInvalidInput
	}

	var strValue string
	switch v := ctx.Value.(type) {
	case string:
		strValue = v
	case []byte:
		strValue = string(v)
	case fmt.Stringer:
		strValue = v.String()
	default:
		return nil, fmt.Errorf("%w: expected string but got %T", ErrIncompatibleType, ctx.Value)
	}

	// Chuyển đổi sang int với radix
	intValue, err := strconv.ParseInt(strValue, p.radix, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse int: %w", err)
	}

	// Chuyển đổi kiểu nếu cần
	switch ctx.Type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intValue, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if intValue < 0 {
			return nil, fmt.Errorf("cannot convert negative value to unsigned int")
		}
		return uint64(intValue), nil
	default:
		return nil, fmt.Errorf("%w: expected integer type but got %s", ErrIncompatibleType, ctx.Type)
	}
}

// ParseFloatPipe chuyển đổi string thành float
type ParseFloatPipe struct {
	BasePipe
	bitSize int
}

// NewParseFloatPipe tạo một ParseFloatPipe mới
func NewParseFloatPipe(bitSize int) *ParseFloatPipe {
	if bitSize != 32 && bitSize != 64 {
		bitSize = 64
	}
	return &ParseFloatPipe{
		bitSize: bitSize,
	}
}

// Transform chuyển đổi string thành float
func (p *ParseFloatPipe) Transform(ctx *TransformContext) (interface{}, error) {
	if ctx.Value == nil {
		return nil, ErrInvalidInput
	}

	var strValue string
	switch v := ctx.Value.(type) {
	case string:
		strValue = v
	case []byte:
		strValue = string(v)
	case fmt.Stringer:
		strValue = v.String()
	default:
		return nil, fmt.Errorf("%w: expected string but got %T", ErrIncompatibleType, ctx.Value)
	}

	// Chuyển đổi sang float
	floatValue, err := strconv.ParseFloat(strValue, p.bitSize)
	if err != nil {
		return nil, fmt.Errorf("failed to parse float: %w", err)
	}

	// Chuyển đổi kiểu nếu cần
	switch ctx.Type.Kind() {
	case reflect.Float32:
		return float32(floatValue), nil
	case reflect.Float64:
		return floatValue, nil
	default:
		return nil, fmt.Errorf("%w: expected float type but got %s", ErrIncompatibleType, ctx.Type)
	}
}

// ParseBoolPipe chuyển đổi string thành bool
type ParseBoolPipe struct {
	BasePipe
}

// NewParseBoolPipe tạo một ParseBoolPipe mới
func NewParseBoolPipe() *ParseBoolPipe {
	return &ParseBoolPipe{}
}

// Transform chuyển đổi string thành bool
func (p *ParseBoolPipe) Transform(ctx *TransformContext) (interface{}, error) {
	if ctx.Value == nil {
		return nil, ErrInvalidInput
	}

	var strValue string
	switch v := ctx.Value.(type) {
	case string:
		strValue = v
	case []byte:
		strValue = string(v)
	case fmt.Stringer:
		strValue = v.String()
	case bool:
		return v, nil
	default:
		return nil, fmt.Errorf("%w: expected string but got %T", ErrIncompatibleType, ctx.Value)
	}

	// Chuyển đổi các giá trị phổ biến
	strValue = strings.ToLower(strings.TrimSpace(strValue))
	switch strValue {
	case "true", "1", "yes", "y":
		return true, nil
	case "false", "0", "no", "n":
		return false, nil
	default:
		return nil, fmt.Errorf("cannot parse %q as bool", strValue)
	}
}

// DefaultValuePipe thiết lập giá trị mặc định nếu đầu vào là nil
type DefaultValuePipe struct {
	BasePipe
	defaultValue interface{}
}

// NewDefaultValuePipe tạo một DefaultValuePipe mới
func NewDefaultValuePipe(defaultValue interface{}) *DefaultValuePipe {
	return &DefaultValuePipe{
		defaultValue: defaultValue,
	}
}

// Transform trả về giá trị mặc định nếu đầu vào là nil
func (p *DefaultValuePipe) Transform(ctx *TransformContext) (interface{}, error) {
	if ctx.Value == nil {
		return p.defaultValue, nil
	}
	return ctx.Value, nil
}

// TrimPipe loại bỏ khoảng trắng ở đầu và cuối string
type TrimPipe struct {
	BasePipe
}

// NewTrimPipe tạo một TrimPipe mới
func NewTrimPipe() *TrimPipe {
	return &TrimPipe{}
}

// Transform loại bỏ khoảng trắng ở đầu và cuối string
func (p *TrimPipe) Transform(ctx *TransformContext) (interface{}, error) {
	if ctx.Value == nil {
		return nil, ErrInvalidInput
	}

	var strValue string
	switch v := ctx.Value.(type) {
	case string:
		strValue = v
	case []byte:
		strValue = string(v)
	case fmt.Stringer:
		strValue = v.String()
	default:
		return nil, fmt.Errorf("%w: expected string but got %T", ErrIncompatibleType, ctx.Value)
	}

	return strings.TrimSpace(strValue), nil
}

// LowerCasePipe chuyển đổi string sang chữ thường
type LowerCasePipe struct {
	BasePipe
}

// NewLowerCasePipe tạo một LowerCasePipe mới
func NewLowerCasePipe() *LowerCasePipe {
	return &LowerCasePipe{}
}

// Transform chuyển đổi string sang chữ thường
func (p *LowerCasePipe) Transform(ctx *TransformContext) (interface{}, error) {
	if ctx.Value == nil {
		return nil, ErrInvalidInput
	}

	var strValue string
	switch v := ctx.Value.(type) {
	case string:
		strValue = v
	case []byte:
		strValue = string(v)
	case fmt.Stringer:
		strValue = v.String()
	default:
		return nil, fmt.Errorf("%w: expected string but got %T", ErrIncompatibleType, ctx.Value)
	}

	return strings.ToLower(strValue), nil
}

// UpperCasePipe chuyển đổi string sang chữ hoa
type UpperCasePipe struct {
	BasePipe
}

// NewUpperCasePipe tạo một UpperCasePipe mới
func NewUpperCasePipe() *UpperCasePipe {
	return &UpperCasePipe{}
}

// Transform chuyển đổi string sang chữ hoa
func (p *UpperCasePipe) Transform(ctx *TransformContext) (interface{}, error) {
	if ctx.Value == nil {
		return nil, ErrInvalidInput
	}

	var strValue string
	switch v := ctx.Value.(type) {
	case string:
		strValue = v
	case []byte:
		strValue = string(v)
	case fmt.Stringer:
		strValue = v.String()
	default:
		return nil, fmt.Errorf("%w: expected string but got %T", ErrIncompatibleType, ctx.Value)
	}

	return strings.ToUpper(strValue), nil
}
