package core

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type Application struct {
	mu        sync.RWMutex
	container *Container
	engine    *gin.Engine
	modules   []Module
	options   []fx.Option
}

func NewApplication() *Application {
	return &Application{
		container: NewContainer(),
		engine:    gin.Default(),
		modules:   make([]Module, 0),
		options:   make([]fx.Option, 0),
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
		if err := app.engine.Run(":8080"); err != nil {
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
