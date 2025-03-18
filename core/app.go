// goblin/core/app.go
// Package core provides the core functionality for the Goblin Framework.
// It implements the main application structure and lifecycle management,
// integrating with Gin for HTTP handling and Fx for dependency injection.
package core

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// GoblinApp represents a Goblin application instance.
// It embeds the Gin Engine and manages the application lifecycle,
// including module management and dependency injection.
type GoblinApp struct {
	// Engine is the embedded Gin HTTP engine
	*gin.Engine
	// app is the Fx application instance for dependency injection
	app *fx.App
	// moduleManager handles module registration and lifecycle
	moduleManager *ModuleManager
	// lifecycleManager manages application lifecycle hooks
	lifecycleManager *LifecycleManager
}

// GoblinAppOptions configures a new Goblin application.
// It allows customization of modules and debug settings.
type GoblinAppOptions struct {
	// Modules specifies the list of modules to initialize
	Modules []Module
	// Debug enables debug mode when true
	Debug bool
}

// NewGoblinApp creates a new Goblin application instance.
// It initializes the Gin engine, sets up module management,
// and configures dependency injection.
//
// Parameters:
//   - opts: Optional configuration for the application
//
// Returns:
//   - *GoblinApp: A new Goblin application instance
func NewGoblinApp(opts ...GoblinAppOptions) *GoblinApp {
	var options GoblinAppOptions
	if len(opts) > 0 {
		options = opts[0]
	}

	// Set default mode to production unless debug is enabled
	if !options.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	if options.Debug {
		engine.Use(gin.Logger())
	}

	// Initialize module manager
	moduleManager := NewModuleManager()

	// Initialize lifecycle manager
	lifecycleManager := NewLifecycleManager()

	// Collect all providers and controllers from modules
	var providers []interface{}
	var controllers []interface{}

	// Register and collect providers/controllers from all modules
	for _, module := range options.Modules {
		if err := moduleManager.RegisterModule(module); err != nil {
			log.Printf("Warning: failed to register module: %v", err)
			continue
		}

		// Collect providers
		moduleProviders := moduleManager.GetModuleProviders(module)
		providers = append(providers, moduleProviders...)

		// Collect controllers
		moduleControllers := moduleManager.GetModuleControllers(module)
		controllers = append(controllers, moduleControllers...)
	}

	// Register modules with lifecycle manager
	lifecycleManager.RegisterModules(options.Modules)

	// Register providers with lifecycle manager
	lifecycleManager.ExtractLifecycleHooks(providers)

	// Add engine to providers
	providers = append(providers, func() *gin.Engine { return engine })

	// Add lifecycle manager to providers
	providers = append(providers, func() *LifecycleManager { return lifecycleManager })

	// Create Fx app with providers and controllers
	app := fx.New(
		fx.Provide(providers...),
		fx.Invoke(controllers...),
	)

	return &GoblinApp{
		Engine:           engine,
		app:              app,
		moduleManager:    moduleManager,
		lifecycleManager: lifecycleManager,
	}
}

// Start initializes and starts the Goblin application.
// It sets up signal handling for graceful shutdown,
// initializes modules, and starts the HTTP server.
//
// Parameters:
//   - port: The port number to listen on (e.g., ":8080")
//
// Returns:
//   - error: Any error that occurred during startup
func (g *GoblinApp) Start(port string) error {
	// Create context for the entire application
	ctx := context.Background()

	// Set up signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		log.Println("Received termination signal, shutting down gracefully...")
		g.Stop()
	}()

	// Initialize modules
	if err := g.moduleManager.InitializeModules(ctx); err != nil {
		return fmt.Errorf("failed to initialize modules: %w", err)
	}

	// Run module initialization hooks
	if err := g.lifecycleManager.RunModuleInit(ctx); err != nil {
		return fmt.Errorf("failed to run module initialization hooks: %w", err)
	}

	// Start Fx app
	startCtx, cancel := context.WithTimeout(context.Background(), fx.DefaultTimeout)
	defer cancel()

	if err := g.app.Start(startCtx); err != nil {
		return fmt.Errorf("failed to start application: %w", err)
	}

	// Run application bootstrap hooks
	if err := g.lifecycleManager.RunAppBootstrap(ctx); err != nil {
		return fmt.Errorf("failed to run application bootstrap hooks: %w", err)
	}

	// Start Gin server
	log.Printf("Goblin app started on port %s", port)
	return g.Engine.Run(port)
}

// Stop gracefully shuts down the Goblin application.
// It runs shutdown hooks, destroys modules, and stops the Fx app.
//
// Returns:
//   - error: Any error that occurred during shutdown
func (g *GoblinApp) Stop() error {
	ctx := context.Background()

	// Run application shutdown hooks
	if err := g.lifecycleManager.RunAppShutdown(ctx); err != nil {
		log.Printf("Warning: failed to run application shutdown hooks: %v", err)
	}

	// Destroy modules
	if err := g.moduleManager.DestroyModules(ctx); err != nil {
		log.Printf("Warning: failed to destroy modules: %v", err)
	}

	// Run module destroy hooks
	if err := g.lifecycleManager.RunModuleDestroy(ctx); err != nil {
		log.Printf("Warning: failed to run module destroy hooks: %v", err)
	}

	// Stop Fx app
	stopCtx, cancel := context.WithTimeout(ctx, fx.DefaultTimeout)
	defer cancel()

	return g.app.Stop(stopCtx)
}

// RegisterModules registers additional modules with the application.
// It initializes the modules and updates the dependency injection container.
//
// Parameters:
//   - modules: The modules to register
//
// Returns:
//   - error: Any error that occurred during module registration
func (g *GoblinApp) RegisterModules(modules ...Module) error {
	ctx := context.Background()

	// Register modules with lifecycle manager
	g.lifecycleManager.RegisterModules(modules)

	for _, module := range modules {
		// Register module
		if err := g.moduleManager.RegisterModule(module); err != nil {
			return fmt.Errorf("failed to register module: %w", err)
		}

		// Collect providers and controllers
		providers := g.moduleManager.GetModuleProviders(module)
		controllers := g.moduleManager.GetModuleControllers(module)

		// Register providers with lifecycle manager
		g.lifecycleManager.ExtractLifecycleHooks(providers)

		// Create new Fx options
		options := fx.Options(
			fx.Provide(providers...),
			fx.Invoke(controllers...),
		)

		// Initialize module
		if err := module.OnModuleInit(ctx); err != nil {
			return fmt.Errorf("failed to initialize module: %w", err)
		}

		// Update Fx app
		g.app = fx.New(options)
	}

	return nil
}

// GetApp returns the underlying Fx application instance.
//
// Returns:
//   - *fx.App: The Fx application instance
func (g *GoblinApp) GetApp() *fx.App {
	return g.app
}

// GetModuleManager returns the module manager instance.
//
// Returns:
//   - *ModuleManager: The module manager instance
func (g *GoblinApp) GetModuleManager() *ModuleManager {
	return g.moduleManager
}

// GetLifecycleManager returns the lifecycle manager instance.
//
// Returns:
//   - *LifecycleManager: The lifecycle manager instance
func (g *GoblinApp) GetLifecycleManager() *LifecycleManager {
	return g.lifecycleManager
}

// RegisterShutdownHook registers a function to be called during application shutdown.
//
// Parameters:
//   - hook: The function to call during shutdown
func (g *GoblinApp) RegisterShutdownHook(hook func(ctx context.Context) error) {
	g.lifecycleManager.RegisterShutdownHook(hook)
}
