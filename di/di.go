package di

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// Scope defines the lifetime of a dependency
type Scope int

const (
	// Singleton - one instance for entire application
	Singleton Scope = iota
	// Transient - new instance created each time
	Transient
	// RequestScoped - one instance per request
	RequestScoped
)

// Provider represents a dependency provider
type Provider struct {
	Type        reflect.Type
	Constructor interface{}
	Scope       Scope
	instance    interface{}
	mutex       sync.RWMutex
}

// Container manages dependencies
type Container struct {
	providers map[reflect.Type]*Provider
	mutex     sync.RWMutex
	app       *fx.App
}

// NewContainer creates a new DI container
func NewContainer() *Container {
	return &Container{
		providers: make(map[reflect.Type]*Provider),
	}
}

// Register adds a provider to the container
func (c *Container) Register(constructor interface{}, scope Scope) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	constructorType := reflect.TypeOf(constructor)
	if constructorType.Kind() != reflect.Func {
		return fmt.Errorf("constructor must be a function")
	}

	// Get the return type
	if constructorType.NumOut() != 1 {
		return fmt.Errorf("constructor must return exactly one value")
	}
	returnType := constructorType.Out(0)

	// Check for circular dependencies
	if err := c.checkCircularDependencies(constructorType, make(map[reflect.Type]bool)); err != nil {
		return err
	}

	c.providers[returnType] = &Provider{
		Type:        returnType,
		Constructor: constructor,
		Scope:       scope,
	}

	return nil
}

// checkCircularDependencies checks for circular dependencies in the constructor
func (c *Container) checkCircularDependencies(constructorType reflect.Type, visited map[reflect.Type]bool) error {
	for i := 0; i < constructorType.NumIn(); i++ {
		dependencyType := constructorType.In(i)

		if visited[dependencyType] {
			return fmt.Errorf("circular dependency detected for type %v", dependencyType)
		}

		visited[dependencyType] = true

		if provider, exists := c.providers[dependencyType]; exists {
			if err := c.checkCircularDependencies(reflect.TypeOf(provider.Constructor), visited); err != nil {
				return err
			}
		}

		delete(visited, dependencyType)
	}

	return nil
}

// Resolve gets an instance of the requested type
func (c *Container) Resolve(t reflect.Type, ctx *gin.Context) (interface{}, error) {
	c.mutex.RLock()
	provider, exists := c.providers[t]
	c.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no provider registered for type %v", t)
	}

	return c.resolveProvider(provider, ctx)
}

// resolveProvider creates or returns an instance based on scope
func (c *Container) resolveProvider(provider *Provider, ctx *gin.Context) (interface{}, error) {
	switch provider.Scope {
	case Singleton:
		return c.resolveSingleton(provider)
	case Transient:
		return c.resolveTransient(provider, ctx)
	case RequestScoped:
		return c.resolveRequestScoped(provider, ctx)
	default:
		return nil, fmt.Errorf("unknown scope: %v", provider.Scope)
	}
}

// resolveSingleton returns or creates a singleton instance
func (c *Container) resolveSingleton(provider *Provider) (interface{}, error) {
	provider.mutex.RLock()
	if provider.instance != nil {
		instance := provider.instance
		provider.mutex.RUnlock()
		return instance, nil
	}
	provider.mutex.RUnlock()

	provider.mutex.Lock()
	defer provider.mutex.Unlock()

	// Double-check after acquiring write lock
	if provider.instance != nil {
		return provider.instance, nil
	}

	instance, err := c.createInstance(provider, nil)
	if err != nil {
		return nil, err
	}

	provider.instance = instance
	return instance, nil
}

// resolveTransient creates a new instance each time
func (c *Container) resolveTransient(provider *Provider, ctx *gin.Context) (interface{}, error) {
	return c.createInstance(provider, ctx)
}

// resolveRequestScoped returns or creates an instance per request
func (c *Container) resolveRequestScoped(provider *Provider, ctx *gin.Context) (interface{}, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context required for request-scoped dependency")
	}

	// Use context to store request-scoped instances
	key := fmt.Sprintf("di:%v", provider.Type)
	if instance, exists := ctx.Get(key); exists {
		return instance, nil
	}

	instance, err := c.createInstance(provider, ctx)
	if err != nil {
		return nil, err
	}

	ctx.Set(key, instance)
	return instance, nil
}

// createInstance creates a new instance using the constructor
func (c *Container) createInstance(provider *Provider, ctx *gin.Context) (interface{}, error) {
	constructorType := reflect.TypeOf(provider.Constructor)
	constructorValue := reflect.ValueOf(provider.Constructor)

	// Prepare arguments for the constructor
	args := make([]reflect.Value, constructorType.NumIn())
	for i := 0; i < constructorType.NumIn(); i++ {
		argType := constructorType.In(i)

		// Special handling for gin.Context
		if argType == reflect.TypeOf(&gin.Context{}) {
			if ctx == nil {
				return nil, fmt.Errorf("gin.Context required but not provided")
			}
			args[i] = reflect.ValueOf(ctx)
			continue
		}

		// Resolve dependencies recursively
		arg, err := c.Resolve(argType, ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve argument %v: %v", argType, err)
		}
		args[i] = reflect.ValueOf(arg)
	}

	// Call the constructor
	results := constructorValue.Call(args)
	if len(results) != 1 {
		return nil, fmt.Errorf("constructor must return exactly one value")
	}

	return results[0].Interface(), nil
}

// BuildFxOptions creates fx.Option for the container
func (c *Container) BuildFxOptions() fx.Option {
	var options []fx.Option

	for _, provider := range c.providers {
		if provider.Scope == Singleton {
			options = append(options, fx.Provide(provider.Constructor))
		}
	}

	return fx.Options(options...)
}

// Inject is a decorator that injects dependencies into struct fields
func Inject(target interface{}) error {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to struct")
	}

	targetElem := targetValue.Elem()
	targetType := targetElem.Type()

	for i := 0; i < targetElem.NumField(); i++ {
		field := targetElem.Field(i)
		fieldType := targetType.Field(i)

		// Check if field should be injected
		if _, ok := fieldType.Tag.Lookup("inject"); !ok {
			continue
		}

		// Get the container instance (you'll need to store this globally or pass it somehow)
		container := NewContainer() // This is a placeholder - you'll need to implement proper container access

		// Resolve the dependency
		instance, err := container.Resolve(field.Type(), nil)
		if err != nil {
			return fmt.Errorf("failed to inject field %s: %v", fieldType.Name, err)
		}

		// Set the field value
		field.Set(reflect.ValueOf(instance))
	}

	return nil
}
