package goblin

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

// Config struct
type Config struct {
	Port string
}

// Application struct
type Application struct {
	app              *fx.App
	done             chan struct{}
	globalMiddleware *globalMiddleware
	moduleRegistry   *ModuleRegistry
}

// New function to create a new Goblin application
func NewGoblinApp(opts ...fx.Option) *Application {
	app := &Application{
		done:             make(chan struct{}),
		globalMiddleware: newGlobalMiddleware(),
	}
	baseOptions := fx.Options(
		fx.Provide(
			func() *gin.Engine { return gin.Default() },
			func() Config { return Config{Port: ":8080"} },
			NewModuleRegistry, // Provide the function to create ModuleRegistry
		),
		fx.WithLogger(func() fxevent.Logger {
			return &fxevent.ConsoleLogger{W: os.Stdout}
		}),
	)

	allOptions := fx.Options(append([]fx.Option{baseOptions}, opts...)...)

	fxApp := fx.New(
		allOptions,
		fx.Invoke(app.registerLifecycleHooks),
	)
	app.app = fxApp
	return app
}

// registerLifecycleHooks function to register lifecycle hooks
func (app *Application) registerLifecycleHooks(lc fx.Lifecycle, engine *gin.Engine, moduleRegistry *ModuleRegistry, cfg Config) {
	srv := &http.Server{
		Addr:         cfg.Port,
		Handler:      engine,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := moduleRegistry.LoadModules(); err != nil {
				return fmt.Errorf("failed to load modules: %w", err)
			}
			registerMiddleware(engine, app.globalMiddleware)
			moduleRegistry.RegisterRoutesForModules(engine)
			log.Printf("Server is listening on %s", cfg.Port)
			go func() {
				if err := srv.ListenAndServe(); err != http.ErrServerClosed {
					log.Printf("failed to start server: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping server...")
			return srv.Shutdown(ctx)
		},
	})
}

// WithPort function to override the default port
func WithPort(port string) fx.Option {
	return fx.Replace(Config{Port: port})
}

// Run function to start the application
func (a *Application) Run() {
	if err := a.app.Start(context.Background()); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
	<-a.done
}

// Stop function to stop the application
func (a *Application) Stop() {
	a.app.Stop(context.Background())
}

// AddGlobalMiddleware function
func (a *Application) AddGlobalMiddleware(middlewares ...gin.HandlerFunc) {
	a.globalMiddleware.addMiddlewares(middlewares...)
}
