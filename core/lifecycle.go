package core

import (
	"context"
	"log"
	"reflect"
	"sync"
)

// LifecycleState defines the possible states in the application lifecycle.
// The application transitions through these states in a specific order:
// NotStarted -> ModuleInit -> AppBootstrap -> Running -> AppShutdown -> ModuleDestroy
type LifecycleState int

const (
	// StateNotStarted indicates the application has not yet started.
	// This is the initial state when the application is first created.
	StateNotStarted LifecycleState = iota
	// StateModuleInit indicates modules are being initialized.
	// During this state, all modules perform their initialization tasks.
	StateModuleInit
	// StateAppBootstrap indicates the application is being bootstrapped.
	// This state allows for application-wide setup before accepting requests.
	StateAppBootstrap
	// StateRunning indicates the application is running and ready to handle requests.
	// This is the normal operational state of the application.
	StateRunning
	// StateAppShutdown indicates the application is shutting down.
	// During this state, the application stops accepting new requests and begins cleanup.
	StateAppShutdown
	// StateModuleDestroy indicates modules are being destroyed.
	// This is the final state where all resources are cleaned up.
	StateModuleDestroy
)

// OnModuleInit is called after a module is initialized.
// This interface allows modules to perform any necessary setup,
// such as initializing connections, loading configurations, or preparing resources.
// The context provided can be used for cancellation and timeout control.
type OnModuleInit interface {
	// OnModuleInit is called with a context for the initialization operation.
	// It should return any error that occurred during initialization.
	// If an error is returned, the application startup will be aborted.
	OnModuleInit(ctx context.Context) error
}

// OnApplicationBootstrap is called after all modules are initialized,
// before the application is ready to handle requests.
// This interface is useful for application-wide setup that depends on
// all modules being initialized, such as starting servers or initializing
// global services.
type OnApplicationBootstrap interface {
	// OnApplicationBootstrap is called with a context for the bootstrap operation.
	// It should return any error that occurred during bootstrap.
	// If an error is returned, the application startup will be aborted.
	OnApplicationBootstrap(ctx context.Context) error
}

// OnApplicationShutdown is called when the application begins its shutdown process.
// This interface allows components to gracefully stop accepting new requests
// and begin cleanup of resources. It's called before module destruction.
type OnApplicationShutdown interface {
	// OnApplicationShutdown is called with a context for the shutdown operation.
	// It should return any error that occurred during shutdown.
	// Errors during shutdown are logged but don't prevent the shutdown process.
	OnApplicationShutdown(ctx context.Context) error
}

// OnModuleDestroy is called when a module is being destroyed.
// This interface allows modules to clean up their resources, such as
// closing connections, stopping background tasks, or releasing system resources.
// Modules are destroyed in reverse order of initialization.
type OnModuleDestroy interface {
	// OnModuleDestroy is called with a context for the destruction operation.
	// It should return any error that occurred during destruction.
	// Errors during destruction are logged but don't prevent the destruction process.
	OnModuleDestroy(ctx context.Context) error
}

// LifecycleManager manages the application lifecycle.
// It coordinates the initialization and shutdown of modules and providers,
// ensuring proper order of operations and error handling.
// The manager is thread-safe and can be used concurrently.
type LifecycleManager struct {
	// state represents the current lifecycle state
	state LifecycleState
	// modules contains all registered modules
	modules []Module
	// providers contains all registered providers that implement lifecycle hooks
	providers []interface{}
	// shutdownHooks contains functions to be called during application shutdown
	shutdownHooks []func(ctx context.Context) error
	// stateMutex protects concurrent access to state and collections
	stateMutex sync.RWMutex
	// shutdownSignals is used to coordinate shutdown operations
	shutdownSignals chan struct{}
}

/// ... existing code ...

// NewLifecycleManager creates a new LifecycleManager.
// It initializes all fields with empty slices and default state.
// The manager starts in StateNotStarted and is ready to register modules and providers.
//
// Returns:
//   - *LifecycleManager: A new lifecycle manager instance
func NewLifecycleManager() *LifecycleManager {
	return &LifecycleManager{
		state:           StateNotStarted,
		modules:         make([]Module, 0),
		providers:       make([]interface{}, 0),
		shutdownHooks:   make([]func(ctx context.Context) error, 0),
		shutdownSignals: make(chan struct{}),
	}
}

// RegisterModules registers modules with the lifecycle manager.
// This method is thread-safe and can be called concurrently.
// Modules will be initialized in the order they are registered.
//
// Parameters:
//   - modules: The modules to register
func (m *LifecycleManager) RegisterModules(modules []Module) {
	m.stateMutex.Lock()
	defer m.stateMutex.Unlock()

	m.modules = append(m.modules, modules...)
}

// RegisterProviders registers providers with the lifecycle manager.
// This method is thread-safe and can be called concurrently.
// Providers that implement lifecycle hooks will be called during the appropriate lifecycle events.
//
// Parameters:
//   - providers: The providers to register
func (m *LifecycleManager) RegisterProviders(providers []interface{}) {
	m.stateMutex.Lock()
	defer m.stateMutex.Unlock()

	m.providers = append(m.providers, providers...)
}

// RunModuleInit calls OnModuleInit hooks for all modules and providers.
// This method transitions the application to StateModuleInit and calls
// initialization hooks in the order modules and providers were registered.
// If any hook returns an error, the initialization process is aborted.
//
// Parameters:
//   - ctx: The context for the initialization operation
//
// Returns:
//   - error: Any error that occurred during initialization
func (m *LifecycleManager) RunModuleInit(ctx context.Context) error {
	m.stateMutex.Lock()
	m.state = StateModuleInit
	m.stateMutex.Unlock()

	log.Println("Running module initialization hooks...")

	// Call hooks for modules
	for _, module := range m.modules {
		if hook, ok := module.(OnModuleInit); ok {
			if err := hook.OnModuleInit(ctx); err != nil {
				return err
			}
		}
	}

	// Call hooks for providers
	for _, provider := range m.providers {
		if hook, ok := provider.(OnModuleInit); ok {
			if err := hook.OnModuleInit(ctx); err != nil {
				return err
			}
		}
	}

	return nil
}

// RunAppBootstrap calls OnApplicationBootstrap hooks for all modules and providers.
// This method transitions the application to StateAppBootstrap and calls
// bootstrap hooks in the order modules and providers were registered.
// After successful bootstrap, the application transitions to StateRunning.
// If any hook returns an error, the bootstrap process is aborted.
//
// Parameters:
//   - ctx: The context for the bootstrap operation
//
// Returns:
//   - error: Any error that occurred during bootstrap
func (m *LifecycleManager) RunAppBootstrap(ctx context.Context) error {
	m.stateMutex.Lock()
	m.state = StateAppBootstrap
	m.stateMutex.Unlock()

	log.Println("Running application bootstrap hooks...")

	// Call hooks for modules
	for _, module := range m.modules {
		if hook, ok := module.(OnApplicationBootstrap); ok {
			if err := hook.OnApplicationBootstrap(ctx); err != nil {
				return err
			}
		}
	}

	// Call hooks for providers
	for _, provider := range m.providers {
		if hook, ok := provider.(OnApplicationBootstrap); ok {
			if err := hook.OnApplicationBootstrap(ctx); err != nil {
				return err
			}
		}
	}

	m.stateMutex.Lock()
	m.state = StateRunning
	m.stateMutex.Unlock()

	return nil
}

// RunAppShutdown calls OnApplicationShutdown hooks for all modules and providers.
// This method transitions the application to StateAppShutdown and calls
// shutdown hooks in the order modules and providers were registered.
// After this, any registered shutdown hooks are called.
// Errors during shutdown are logged but don't prevent the shutdown process.
//
// Parameters:
//   - ctx: The context for the shutdown operation
//
// Returns:
//   - error: Any error that occurred during shutdown
func (m *LifecycleManager) RunAppShutdown(ctx context.Context) error {
	m.stateMutex.Lock()
	m.state = StateAppShutdown
	m.stateMutex.Unlock()

	log.Println("Running application shutdown hooks...")

	// Call hooks for modules
	for _, module := range m.modules {
		if hook, ok := module.(OnApplicationShutdown); ok {
			if err := hook.OnApplicationShutdown(ctx); err != nil {
				log.Printf("Error during module shutdown: %v\n", err)
			}
		}
	}

	// Call hooks for providers
	for _, provider := range m.providers {
		if hook, ok := provider.(OnApplicationShutdown); ok {
			if err := hook.OnApplicationShutdown(ctx); err != nil {
				log.Printf("Error during provider shutdown: %v\n", err)
			}
		}
	}

	// Call registered shutdown hooks
	for _, hook := range m.shutdownHooks {
		if err := hook(ctx); err != nil {
			log.Printf("Error during shutdown hook: %v\n", err)
		}
	}

	return nil
}

// RunModuleDestroy calls OnModuleDestroy hooks for all modules and providers.
// This method transitions the application to StateModuleDestroy and calls
// destroy hooks in reverse order of initialization (last registered, first destroyed).
// Errors during destruction are logged but don't prevent the destruction process.
//
// Parameters:
//   - ctx: The context for the destruction operation
//
// Returns:
//   - error: Any error that occurred during destruction
func (m *LifecycleManager) RunModuleDestroy(ctx context.Context) error {
	m.stateMutex.Lock()
	m.state = StateModuleDestroy
	m.stateMutex.Unlock()

	log.Println("Running module destroy hooks...")

	// Call hooks for modules (in reverse initialization order)
	for i := len(m.modules) - 1; i >= 0; i-- {
		if hook, ok := m.modules[i].(OnModuleDestroy); ok {
			if err := hook.OnModuleDestroy(ctx); err != nil {
				log.Printf("Error during module destroy: %v\n", err)
			}
		}
	}

	// Call hooks for providers (in reverse initialization order)
	for i := len(m.providers) - 1; i >= 0; i-- {
		if hook, ok := m.providers[i].(OnModuleDestroy); ok {
			if err := hook.OnModuleDestroy(ctx); err != nil {
				log.Printf("Error during provider destroy: %v\n", err)
			}
		}
	}

	return nil
}

// GetState returns the current lifecycle state.
// This method is thread-safe and can be called concurrently.
//
// Returns:
//   - LifecycleState: The current state of the lifecycle
func (m *LifecycleManager) GetState() LifecycleState {
	m.stateMutex.RLock()
	defer m.stateMutex.RUnlock()
	return m.state
}

// RegisterShutdownHook registers a function to be called during application shutdown.
// This method is thread-safe and can be called concurrently.
// The registered hooks will be called after all module and provider shutdown hooks.
//
// Parameters:
//   - hook: The function to call during shutdown
func (m *LifecycleManager) RegisterShutdownHook(hook func(ctx context.Context) error) {
	m.stateMutex.Lock()
	defer m.stateMutex.Unlock()
	m.shutdownHooks = append(m.shutdownHooks, hook)
}

// ExtractLifecycleHooks extracts all providers that implement lifecycle hooks.
// This method filters out function providers (factories) and adds the rest to the manager.
// It's used to automatically discover and register providers that implement
// any of the lifecycle hook interfaces.
//
// Parameters:
//   - providers: The providers to check for lifecycle hooks
func (m *LifecycleManager) ExtractLifecycleHooks(providers []interface{}) {
	for _, provider := range providers {
		providerType := reflect.TypeOf(provider)
		if providerType.Kind() == reflect.Func {
			// Skip function providers (factories)
			continue
		}

		m.providers = append(m.providers, provider)
	}
}
