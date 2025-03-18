package di

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

// Test types
type (
	TestSingleton struct {
		Value int
	}

	TestTransient struct {
		Value int
	}

	TestRequestScoped struct {
		Value int
	}

	TestService struct {
		Singleton     *TestSingleton     `inject:""`
		Transient     *TestTransient     `inject:""`
		RequestScoped *TestRequestScoped `inject:""`
	}
)

// Constructors
func NewTestSingleton() *TestSingleton {
	return &TestSingleton{Value: 1}
}

func NewTestTransient() *TestTransient {
	return &TestTransient{Value: 1}
}

func NewTestRequestScoped(c *gin.Context) *TestRequestScoped {
	return &TestRequestScoped{Value: 1}
}

// TestSingletonScope tests singleton scope behavior
func TestSingletonScope(t *testing.T) {
	container := NewContainer()
	err := container.Register(NewTestSingleton, Singleton)
	assert.NoError(t, err)

	// Get two instances
	instance1, err := container.Resolve(reflect.TypeOf(&TestSingleton{}), nil)
	assert.NoError(t, err)
	singleton1 := instance1.(*TestSingleton)

	instance2, err := container.Resolve(reflect.TypeOf(&TestSingleton{}), nil)
	assert.NoError(t, err)
	singleton2 := instance2.(*TestSingleton)

	// Should be the same instance
	assert.Equal(t, singleton1, singleton2)

	// Modify first instance
	singleton1.Value = 2

	// Second instance should reflect the change
	assert.Equal(t, 2, singleton2.Value)
}

// TestTransientScope tests transient scope behavior
func TestTransientScope(t *testing.T) {
	container := NewContainer()
	err := container.Register(NewTestTransient, Transient)
	assert.NoError(t, err)

	// Get two instances
	instance1, err := container.Resolve(reflect.TypeOf(&TestTransient{}), nil)
	assert.NoError(t, err)
	transient1 := instance1.(*TestTransient)

	instance2, err := container.Resolve(reflect.TypeOf(&TestTransient{}), nil)
	assert.NoError(t, err)
	transient2 := instance2.(*TestTransient)

	// Should be different instances
	assert.NotEqual(t, transient1, transient2)

	// Modify first instance
	transient1.Value = 2

	// Second instance should not reflect the change
	assert.Equal(t, 1, transient2.Value)
}

// TestRequestScope tests request-scoped behavior
func TestRequestScope(t *testing.T) {
	container := NewContainer()
	err := container.Register(NewTestRequestScoped, RequestScoped)
	assert.NoError(t, err)

	// Create two different contexts
	ctx1 := &gin.Context{}
	ctx2 := &gin.Context{}

	// Get instances for first context
	instance1, err := container.Resolve(reflect.TypeOf(&TestRequestScoped{}), ctx1)
	assert.NoError(t, err)
	request1 := instance1.(*TestRequestScoped)

	instance2, err := container.Resolve(reflect.TypeOf(&TestRequestScoped{}), ctx1)
	assert.NoError(t, err)
	request2 := instance2.(*TestRequestScoped)

	// Should be the same instance within the same context
	assert.Equal(t, request1, request2)

	// Get instance for second context
	instance3, err := container.Resolve(reflect.TypeOf(&TestRequestScoped{}), ctx2)
	assert.NoError(t, err)
	request3 := instance3.(*TestRequestScoped)

	// Should be different instance for different context
	assert.NotEqual(t, request1, request3)
}

// TestCircularDependency tests circular dependency detection
func TestCircularDependency(t *testing.T) {
	container := NewContainer()

	// Create types with circular dependency
	type A struct{}
	type B struct{}

	NewA := func(b *B) *A { return &A{} }
	NewB := func(a *A) *B { return &B{} }

	// Register A
	err := container.Register(NewA, Singleton)
	assert.NoError(t, err)

	// Register B - should detect circular dependency
	err = container.Register(NewB, Singleton)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circular dependency")
}

// TestInject tests the Inject decorator
func TestInject(t *testing.T) {
	container := NewContainer()

	// Register dependencies
	err := container.Register(NewTestSingleton, Singleton)
	assert.NoError(t, err)

	err = container.Register(NewTestTransient, Transient)
	assert.NoError(t, err)

	err = container.Register(NewTestRequestScoped, RequestScoped)
	assert.NoError(t, err)

	// Create service instance
	service := &TestService{}

	// Inject dependencies
	err = Inject(service)
	assert.NoError(t, err)

	// Verify injection
	assert.NotNil(t, service.Singleton)
	assert.NotNil(t, service.Transient)
	assert.NotNil(t, service.RequestScoped)
}

// TestInvalidRegistration tests invalid provider registration
func TestInvalidRegistration(t *testing.T) {
	container := NewContainer()

	// Try to register non-function
	err := container.Register("not a function", Singleton)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be a function")

	// Try to register function with multiple return values
	multiReturn := func() (*TestSingleton, error) { return &TestSingleton{}, nil }
	err = container.Register(multiReturn, Singleton)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must return exactly one value")
}

// TestRequestScopedWithoutContext tests request-scoped resolution without context
func TestRequestScopedWithoutContext(t *testing.T) {
	container := NewContainer()
	err := container.Register(NewTestRequestScoped, RequestScoped)
	assert.NoError(t, err)

	// Try to resolve without context
	_, err = container.Resolve(reflect.TypeOf(&TestRequestScoped{}), nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context required")
}

// TestFxIntegration tests integration with uber-fx
func TestFxIntegration(t *testing.T) {
	container := NewContainer()

	// Register singleton
	err := container.Register(NewTestSingleton, Singleton)
	assert.NoError(t, err)

	// Get fx options
	options := container.BuildFxOptions()
	assert.NotNil(t, options)

	// Create fx app
	app := fx.New(options)
	assert.NotNil(t, app)
}
