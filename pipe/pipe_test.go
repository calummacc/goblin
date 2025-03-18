package pipe

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationPipe(t *testing.T) {
	type TestUser struct {
		ID       int    `json:"id" validate:"required,gt=0"`
		Name     string `json:"name" validate:"required,min=3,max=50"`
		Email    string `json:"email" validate:"required,email"`
		Age      int    `json:"age" validate:"gte=18"`
		IsActive bool   `json:"is_active"`
	}

	validUser := TestUser{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Age:      25,
		IsActive: true,
	}

	invalidUser := TestUser{
		ID:       0,
		Name:     "Jo",
		Email:    "invalid-email",
		Age:      15,
		IsActive: false,
	}

	t.Run("ValidInput", func(t *testing.T) {
		pipe := NewValidationPipe()
		ctx := &TransformContext{
			Value: validUser,
			Type:  reflect.TypeOf(validUser),
		}

		result, err := pipe.Transform(ctx)
		assert.NoError(t, err)
		assert.Equal(t, validUser, result)
	})

	t.Run("InvalidInput", func(t *testing.T) {
		pipe := NewValidationPipe()
		ctx := &TransformContext{
			Value: invalidUser,
			Type:  reflect.TypeOf(invalidUser),
		}

		result, err := pipe.Transform(ctx)
		assert.Error(t, err)
		assert.Nil(t, result)

		// Kiểm tra lỗi validation
		validationErrors, ok := err.(ValidationErrors)
		assert.True(t, ok, "Error should be ValidationErrors")
		assert.GreaterOrEqual(t, len(validationErrors), 3, "Should have at least 3 validation errors")
	})

	t.Run("NilInput", func(t *testing.T) {
		pipe := NewValidationPipe()
		ctx := &TransformContext{
			Value: nil,
			Type:  reflect.TypeOf(TestUser{}),
		}

		result, err := pipe.Transform(ctx)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrInvalidInput, err)
	})

	t.Run("WithTransform", func(t *testing.T) {
		pipe := NewValidationPipe(ValidationOptions{
			Transform: true,
		})

		// Map đầu vào
		input := map[string]interface{}{
			"id":        1,
			"name":      "John Doe",
			"email":     "john@example.com",
			"age":       25,
			"is_active": true,
		}

		ctx := &TransformContext{
			Value: input,
			Type:  reflect.TypeOf(TestUser{}),
		}

		result, err := pipe.Transform(ctx)
		assert.NoError(t, err)

		// Kiểm tra kết quả chuyển đổi
		userPtr, ok := result.(*TestUser)
		assert.True(t, ok, "Result should be *TestUser")
		assert.Equal(t, 1, userPtr.ID)
		assert.Equal(t, "John Doe", userPtr.Name)
	})

	t.Run("WithStopAtFirstError", func(t *testing.T) {
		pipe := NewValidationPipe(ValidationOptions{
			StopAtFirstError: true,
		})
		ctx := &TransformContext{
			Value: invalidUser,
			Type:  reflect.TypeOf(invalidUser),
		}

		_, err := pipe.Transform(ctx)
		assert.Error(t, err)

		// Kiểm tra chỉ có 1 lỗi
		validationErrors, ok := err.(ValidationErrors)
		assert.True(t, ok, "Error should be ValidationErrors")
		assert.Equal(t, 1, len(validationErrors), "Should have exactly 1 validation error")
	})
}

func TestParsePipes(t *testing.T) {
	t.Run("ParseIntPipe", func(t *testing.T) {
		pipe := NewParseIntPipe(10)
		ctx := &TransformContext{
			Value: "123",
			Type:  reflect.TypeOf(int(0)),
		}

		result, err := pipe.Transform(ctx)
		assert.NoError(t, err)
		assert.Equal(t, int64(123), result)

		// Kiểm tra lỗi với input không hợp lệ
		ctx.Value = "abc"
		result, err = pipe.Transform(ctx)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("ParseFloatPipe", func(t *testing.T) {
		pipe := NewParseFloatPipe(64)
		ctx := &TransformContext{
			Value: "123.45",
			Type:  reflect.TypeOf(float64(0)),
		}

		result, err := pipe.Transform(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 123.45, result)

		// Kiểm tra lỗi với input không hợp lệ
		ctx.Value = "abc"
		result, err = pipe.Transform(ctx)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("ParseBoolPipe", func(t *testing.T) {
		pipe := NewParseBoolPipe()

		// Test true values
		trueValues := []string{"true", "True", "1", "yes", "Yes", "y", "Y"}
		for _, value := range trueValues {
			ctx := &TransformContext{
				Value: value,
				Type:  reflect.TypeOf(bool(false)),
			}
			result, err := pipe.Transform(ctx)
			assert.NoError(t, err)
			assert.Equal(t, true, result)
		}

		// Test false values
		falseValues := []string{"false", "False", "0", "no", "No", "n", "N"}
		for _, value := range falseValues {
			ctx := &TransformContext{
				Value: value,
				Type:  reflect.TypeOf(bool(false)),
			}
			result, err := pipe.Transform(ctx)
			assert.NoError(t, err)
			assert.Equal(t, false, result)
		}

		// Kiểm tra lỗi với input không hợp lệ
		ctx := &TransformContext{
			Value: "invalid",
			Type:  reflect.TypeOf(bool(false)),
		}
		result, err := pipe.Transform(ctx)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestTransformPipes(t *testing.T) {
	t.Run("DefaultValuePipe", func(t *testing.T) {
		defaultValue := "default"
		pipe := NewDefaultValuePipe(defaultValue)

		// Kiểm tra nil value
		ctx := &TransformContext{
			Value: nil,
			Type:  reflect.TypeOf(""),
		}
		result, err := pipe.Transform(ctx)
		assert.NoError(t, err)
		assert.Equal(t, defaultValue, result)

		// Kiểm tra non-nil value
		ctx.Value = "actual"
		result, err = pipe.Transform(ctx)
		assert.NoError(t, err)
		assert.Equal(t, "actual", result)
	})

	t.Run("TrimPipe", func(t *testing.T) {
		pipe := NewTrimPipe()
		ctx := &TransformContext{
			Value: "  hello world  ",
			Type:  reflect.TypeOf(""),
		}

		result, err := pipe.Transform(ctx)
		assert.NoError(t, err)
		assert.Equal(t, "hello world", result)
	})

	t.Run("LowerCasePipe", func(t *testing.T) {
		pipe := NewLowerCasePipe()
		ctx := &TransformContext{
			Value: "Hello World",
			Type:  reflect.TypeOf(""),
		}

		result, err := pipe.Transform(ctx)
		assert.NoError(t, err)
		assert.Equal(t, "hello world", result)
	})

	t.Run("UpperCasePipe", func(t *testing.T) {
		pipe := NewUpperCasePipe()
		ctx := &TransformContext{
			Value: "Hello World",
			Type:  reflect.TypeOf(""),
		}

		result, err := pipe.Transform(ctx)
		assert.NoError(t, err)
		assert.Equal(t, "HELLO WORLD", result)
	})
}

func TestCompositePipe(t *testing.T) {
	// Tạo một chuỗi pipe: Trim -> LowerCase -> DefaultValue
	trimPipe := NewTrimPipe()
	lowerPipe := NewLowerCasePipe()
	defaultPipe := NewDefaultValuePipe("default")

	compositePipe := NewCompositePipe(DefaultPipeOptions(), trimPipe, lowerPipe, defaultPipe)

	t.Run("MultipleTransforms", func(t *testing.T) {
		ctx := &TransformContext{
			Value: "  HELLO WORLD  ",
			Type:  reflect.TypeOf(""),
		}

		result, err := compositePipe.Transform(ctx)
		assert.NoError(t, err)
		assert.Equal(t, "hello world", result)
	})

	t.Run("StopOnError", func(t *testing.T) {
		// Tạo một pipe sẽ luôn trả về lỗi
		errorPipe := &TestErrorPipe{err: assert.AnError}

		// Tạo composite pipe với StopOnError = true
		compositePipe := NewCompositePipe(
			PipeOptions{StopOnError: true},
			trimPipe, errorPipe, lowerPipe,
		)

		ctx := &TransformContext{
			Value: "  TEST  ",
			Type:  reflect.TypeOf(""),
		}

		result, err := compositePipe.Transform(ctx)
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		assert.Nil(t, result)
	})

	t.Run("ContinueOnError", func(t *testing.T) {
		// Tạo một pipe sẽ luôn trả về lỗi
		errorPipe := &TestErrorPipe{err: assert.AnError}

		// Tạo composite pipe với StopOnError = false
		compositePipe := NewCompositePipe(
			PipeOptions{StopOnError: false},
			trimPipe, errorPipe, lowerPipe,
		)

		ctx := &TransformContext{
			Value: "  TEST  ",
			Type:  reflect.TypeOf(""),
		}

		// Vẫn tiếp tục xử lý dù có lỗi
		result, err := compositePipe.Transform(ctx)
		assert.NoError(t, err)
		assert.Equal(t, "test", result)
	})
}

// TestErrorPipe là một pipe luôn trả về lỗi, dùng để kiểm thử
type TestErrorPipe struct {
	BasePipe
	err error
}

func (p *TestErrorPipe) Transform(ctx *TransformContext) (interface{}, error) {
	return nil, p.err
}
