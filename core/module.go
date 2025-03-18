// goblin/core/module.go
// Package core provides the core functionality for the Goblin Framework.
// It implements the module system that enables modular application development
// with dependency injection and lifecycle management.
package core

import (
	"context"
	"fmt"
	"reflect"
)

// ModuleMetadata represents the metadata of a module.
// It contains information about the module's dependencies, exports,
// providers, and controllers.
type ModuleMetadata struct {
	// Imports represents other modules that this module depends on.
	// These modules must be initialized before this module.
	Imports []Module
	// Exports represents providers that other modules can use.
	// These providers are made available to dependent modules.
	Exports []interface{}
	// Providers represents the module's providers.
	// These providers are used for dependency injection within the module.
	Providers []interface{}
	// Controllers represents the module's controllers.
	// These controllers handle HTTP requests and implement business logic.
	Controllers []interface{}
}

// Module represents a module in the Goblin Framework.
// A module is a self-contained unit of functionality that can be
// composed with other modules to build a complete application.
type Module interface {
	// GetMetadata returns the module's metadata, including its
	// dependencies, exports, providers, and controllers.
	GetMetadata() ModuleMetadata
	// OnModuleInit is called after the module is initialized.
	// This is where the module can perform any necessary setup.
	OnModuleInit(ctx context.Context) error
	// OnModuleDestroy is called before the module is destroyed.
	// This is where the module can perform any necessary cleanup.
	OnModuleDestroy(ctx context.Context) error
}

// BaseModule provides a base implementation of Module.
// It can be embedded in custom modules to reduce boilerplate code.
type BaseModule struct {
	// metadata contains the module's configuration and dependencies
	metadata ModuleMetadata
}

// NewBaseModule creates a new BaseModule with the given metadata.
//
// Parameters:
//   - metadata: The module's metadata configuration
//
// Returns:
//   - *BaseModule: A new base module instance
func NewBaseModule(metadata ModuleMetadata) *BaseModule {
	return &BaseModule{
		metadata: metadata,
	}
}

// GetMetadata returns the module's metadata.
//
// Returns:
//   - ModuleMetadata: The module's configuration and dependencies
func (m *BaseModule) GetMetadata() ModuleMetadata {
	return m.metadata
}

// OnModuleInit is called after the module is initialized.
// The base implementation does nothing and returns nil.
//
// Parameters:
//   - ctx: The context for the initialization operation
//
// Returns:
//   - error: Any error that occurred during initialization
func (m *BaseModule) OnModuleInit(ctx context.Context) error {
	return nil
}

// OnModuleDestroy is called before the module is destroyed.
// The base implementation does nothing and returns nil.
//
// Parameters:
//   - ctx: The context for the destruction operation
//
// Returns:
//   - error: Any error that occurred during destruction
func (m *BaseModule) OnModuleDestroy(ctx context.Context) error {
	return nil
}

// ModuleBuilder helps to build a module with its dependencies.
// It provides a fluent interface for configuring module metadata.
type ModuleBuilder struct {
	// metadata contains the module's configuration being built
	metadata ModuleMetadata
}

// NewModuleBuilder creates a new ModuleBuilder.
// It initializes empty slices for imports, exports, providers, and controllers.
//
// Returns:
//   - *ModuleBuilder: A new module builder instance
func NewModuleBuilder() *ModuleBuilder {
	return &ModuleBuilder{
		metadata: ModuleMetadata{
			Imports:     make([]Module, 0),
			Exports:     make([]interface{}, 0),
			Providers:   make([]interface{}, 0),
			Controllers: make([]interface{}, 0),
		},
	}
}

// Import adds a module dependency.
//
// Parameters:
//   - module: The module to import
//
// Returns:
//   - *ModuleBuilder: The builder instance for method chaining
func (b *ModuleBuilder) Import(module Module) *ModuleBuilder {
	b.metadata.Imports = append(b.metadata.Imports, module)
	return b
}

// Export adds a provider that other modules can use.
//
// Parameters:
//   - provider: The provider to export
//
// Returns:
//   - *ModuleBuilder: The builder instance for method chaining
func (b *ModuleBuilder) Export(provider interface{}) *ModuleBuilder {
	b.metadata.Exports = append(b.metadata.Exports, provider)
	return b
}

// Provide adds a provider to the module.
//
// Parameters:
//   - provider: The provider to add
//
// Returns:
//   - *ModuleBuilder: The builder instance for method chaining
func (b *ModuleBuilder) Provide(provider interface{}) *ModuleBuilder {
	b.metadata.Providers = append(b.metadata.Providers, provider)
	return b
}

// Controller adds a controller to the module.
//
// Parameters:
//   - controller: The controller to add
//
// Returns:
//   - *ModuleBuilder: The builder instance for method chaining
func (b *ModuleBuilder) Controller(controller interface{}) *ModuleBuilder {
	b.metadata.Controllers = append(b.metadata.Controllers, controller)
	return b
}

// Build creates a new module with the configured metadata.
//
// Returns:
//   - Module: A new module instance with the configured metadata
func (b *ModuleBuilder) Build() Module {
	return NewBaseModule(b.metadata)
}

// ModuleManager manages the lifecycle of modules.
// It handles module registration, initialization, and destruction.
type ModuleManager struct {
	// modules maps module names to their instances
	modules map[string]Module
}

// NewModuleManager creates a new ModuleManager.
//
// Returns:
//   - *ModuleManager: A new module manager instance
func NewModuleManager() *ModuleManager {
	return &ModuleManager{
		modules: make(map[string]Module),
	}
}

// RegisterModule registers a module with the manager.
//
// Parameters:
//   - module: The module to register
//
// Returns:
//   - error: Any error that occurred during registration
func (m *ModuleManager) RegisterModule(module Module) error {
	moduleType := reflect.TypeOf(module)
	moduleName := moduleType.Elem().Name()

	if _, exists := m.modules[moduleName]; exists {
		return fmt.Errorf("module %s is already registered", moduleName)
	}

	m.modules[moduleName] = module
	return nil
}

// GetModule retrieves a module by name.
//
// Parameters:
//   - name: The name of the module to retrieve
//
// Returns:
//   - Module: The module instance if found
//   - bool: Whether the module was found
func (m *ModuleManager) GetModule(name string) (Module, bool) {
	module, exists := m.modules[name]
	return module, exists
}

// InitializeModules initializes all registered modules.
//
// Parameters:
//   - ctx: The context for the initialization operation
//
// Returns:
//   - error: Any error that occurred during initialization
func (m *ModuleManager) InitializeModules(ctx context.Context) error {
	for _, module := range m.modules {
		if err := module.OnModuleInit(ctx); err != nil {
			return fmt.Errorf("failed to initialize module %s: %w", reflect.TypeOf(module).Elem().Name(), err)
		}
	}
	return nil
}

// DestroyModules destroys all registered modules.
//
// Parameters:
//   - ctx: The context for the destruction operation
//
// Returns:
//   - error: Any error that occurred during destruction
func (m *ModuleManager) DestroyModules(ctx context.Context) error {
	for _, module := range m.modules {
		if err := module.OnModuleDestroy(ctx); err != nil {
			return fmt.Errorf("failed to destroy module %s: %w", reflect.TypeOf(module).Elem().Name(), err)
		}
	}
	return nil
}

// GetModuleDependencies returns all dependencies for a module.
// This includes both direct imports and transitive dependencies.
//
// Parameters:
//   - module: The module to get dependencies for
//
// Returns:
//   - []Module: A list of all dependencies
func (m *ModuleManager) GetModuleDependencies(module Module) []Module {
	metadata := module.GetMetadata()
	deps := make([]Module, 0)

	// Add direct imports
	deps = append(deps, metadata.Imports...)

	// Add dependencies of imported modules
	for _, imp := range metadata.Imports {
		deps = append(deps, m.GetModuleDependencies(imp)...)
	}

	return deps
}

// GetModuleProviders returns all providers for a module and its dependencies.
// This includes the module's own providers and exported providers from dependencies.
//
// Parameters:
//   - module: The module to get providers for
//
// Returns:
//   - []interface{}: A list of all providers
func (m *ModuleManager) GetModuleProviders(module Module) []interface{} {
	metadata := module.GetMetadata()
	providers := make([]interface{}, 0)

	// Add module's own providers
	providers = append(providers, metadata.Providers...)

	// Add exported providers from dependencies
	for _, imp := range metadata.Imports {
		impMetadata := imp.GetMetadata()
		providers = append(providers, impMetadata.Exports...)
	}

	return providers
}

// GetModuleControllers returns all controllers for a module and its dependencies.
// This includes the module's own controllers and controllers from dependencies.
//
// Parameters:
//   - module: The module to get controllers for
//
// Returns:
//   - []interface{}: A list of all controllers
func (m *ModuleManager) GetModuleControllers(module Module) []interface{} {
	metadata := module.GetMetadata()
	controllers := make([]interface{}, 0)

	// Add module's own controllers
	controllers = append(controllers, metadata.Controllers...)

	// Add controllers from dependencies
	for _, imp := range metadata.Imports {
		impMetadata := imp.GetMetadata()
		controllers = append(controllers, impMetadata.Controllers...)
	}

	return controllers
}
