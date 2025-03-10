package goblin

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

// Application core
type Application struct {
	app    *fx.App
	Engine *gin.Engine
}

// Creates a new Goblin application
func New(opts ...fx.Option) *Application {
	engine := gin.Default()

	// Base dependencies (Gin engine + Fx logger)
	baseOptions := fx.Options(
		fx.Provide(
			func() *gin.Engine { return engine },
			func() string {
				return ":8080" // Default port
			}),
		fx.WithLogger(func() fxevent.Logger {
			return &fxevent.ConsoleLogger{}
		}),
	)

	// Combine all options
	allOptions := fx.Options(append([]fx.Option{baseOptions}, opts...)...)

	// Build Fx app
	fxApp := fx.New(
		allOptions,
		fx.Invoke(registerLifecycleHooks),
		fx.Invoke(RegisterRoutes), // Use RegisterRoutes in router.go
	)

	return &Application{
		app:    fxApp,
		Engine: engine,
	}
}

// Registers server lifecycle hooks
func registerLifecycleHooks(lc fx.Lifecycle, engine *gin.Engine, port string) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := engine.Run(port); err != nil {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Add graceful shutdown logic here
			return nil
		},
	})
}

// WithPort allows custom port configuration
func WithPort(port string) fx.Option {
	return fx.Provide(func() string { return port })
}

// Starts the application
func (a *Application) Run() {
	a.app.Run()
}
