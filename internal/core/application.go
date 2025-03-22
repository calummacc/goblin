package core

import (
	"context"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type ApplicationOptions struct {
	Port    int    // Port to run the server on
	Host    string // Host to run the server on
	GinMode string // Gin mode (debug, release, test)
}

// Default options
var defaultOptions = ApplicationOptions{
	Port:    8080,
	Host:    "localhost",
	GinMode: gin.DebugMode,
}

type Application struct {
	mu        sync.RWMutex
	container *Container
	engine    *gin.Engine
	modules   []Module
	options   []fx.Option
	config    ApplicationOptions
}

// Option functions for configuration
func WithPort(port int) func(*ApplicationOptions) {
	return func(opts *ApplicationOptions) {
		opts.Port = port
	}
}

func WithHost(host string) func(*ApplicationOptions) {
	return func(opts *ApplicationOptions) {
		opts.Host = host
	}
}

func WithGinMode(mode string) func(*ApplicationOptions) {
	return func(opts *ApplicationOptions) {
		opts.GinMode = mode
	}
}

func NewGoblinApplication(opts ...func(*ApplicationOptions)) *Application {
	// Start with default options
	config := defaultOptions

	// Apply any provided options
	for _, opt := range opts {
		opt(&config)
	}

	// Set Gin mode
	gin.SetMode(config.GinMode)

	return &Application{
		container: NewContainer(),
		engine:    gin.Default(),
		modules:   make([]Module, 0),
		options:   make([]fx.Option, 0),
		config:    config,
	}
}

func (app *Application) AddModule(module Module) {
	app.mu.Lock()
	defer app.mu.Unlock()
	app.modules = append(app.modules, module)
}

func (app *Application) Configure() {
	app.mu.Lock()
	defer app.mu.Unlock()

	// Configure all modules
	for _, module := range app.modules {
		// Initialize module
		module.Configure(app.container)

		// Add module's fx options if available
		if fxModule, ok := module.(FxModule); ok {
			app.options = append(app.options, fxModule.ProvideDependencies())
		}

		// Call lifecycle hooks if available
		if lifecycleModule, ok := module.(LifecycleModule); ok {
			if err := lifecycleModule.OnInit(); err != nil {
				panic(err)
			}
		}
	}

	// Configure Fx
	app.options = append(app.options,
		fx.Provide(
			func() *gin.Engine { return app.engine },
			func() *Container { return app.container },
		),
		fx.Invoke(app.registerRoutes),
	)
}

func (app *Application) Run(ctx context.Context) error {
	// Create Fx application with all options
	fxApp := fx.New(app.options...)

	// Start the application
	if err := fxApp.Start(ctx); err != nil {
		return err
	}

	// Create a channel for server errors
	errChan := make(chan error, 1)

	// Start HTTP server in a goroutine
	go func() {
		addr := fmt.Sprintf("%s:%d", app.config.Host, app.config.Port)
		if err := app.engine.Run(addr); err != nil {
			errChan <- err
		}
	}()

	// Wait for either context cancellation or server error
	select {
	case <-ctx.Done():
		return app.cleanup()
	case err := <-errChan:
		return err
	}
}

func (app *Application) cleanup() error {
	app.mu.Lock()
	defer app.mu.Unlock()

	// Call OnDestroy for all modules that implement LifecycleModule
	for _, module := range app.modules {
		if lifecycleModule, ok := module.(LifecycleModule); ok {
			if err := lifecycleModule.OnDestroy(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (app *Application) registerRoutes() {
	app.mu.RLock()
	defer app.mu.RUnlock()

	for _, module := range app.modules {
		if routeModule, ok := module.(RouteModule); ok {
			routeModule.RegisterRoutes(app.engine.Group(""))
		}
	}
}

// GetEngine returns the underlying Gin engine
func (app *Application) GetEngine() *gin.Engine {
	return app.engine
}

// GetContainer returns the dependency container
func (app *Application) GetContainer() *Container {
	return app.container
}

// GetConfig returns the application configuration
func (app *Application) GetConfig() ApplicationOptions {
	return app.config
}
