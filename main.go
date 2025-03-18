// goblin/main.go
// Package main provides the entry point for the Goblin Framework application.
// It demonstrates how to set up and run a Goblin application with core modules
// and feature modules.
package main

import (
	"context"
	"log"

	"goblin/core"
	"goblin/database"
	"goblin/events"
	"goblin/examples/user_module"
	"goblin/middleware"

	"github.com/gin-gonic/gin"
)

// main initializes and starts the Goblin application.
// It creates a new GoblinApp instance with core and feature modules,
// then starts the server on port 8080.
func main() {
	// Create a new Goblin app
	app := core.NewGoblinApp(core.GoblinAppOptions{
		Debug: true,
		Modules: []core.Module{
			// Core modules
			NewCoreModule(),
			// Feature modules
			user_module.NewUserModule(),
		},
	})

	// Start the app
	if err := app.Start(":8080"); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}

// CoreModule represents the core module of the application.
// It provides essential services and configurations for the application.
type CoreModule struct {
	*core.BaseModule
}

// NewCoreModule creates and initializes a new core module with essential providers,
// controllers, and exports. It sets up:
// - Event bus for event handling
// - Database ORM and repository
// - Transaction management
// - Basic middleware (logging and recovery)
func NewCoreModule() *CoreModule {
	module := &CoreModule{}
	module.BaseModule = core.NewBaseModule(core.ModuleMetadata{
		Providers: []interface{}{
			events.NewEventBus,
			func() (*database.ORM, error) {
				return database.NewORM()
			},
			database.NewRepository,
			database.NewTransactionManager,
		},
		Controllers: []interface{}{
			func(engine *gin.Engine) {
				engine.Use(middleware.Logger())
				engine.Use(middleware.Recovery())
				// Add a simple ping endpoint for health checks
				engine.GET("/ping", func(c *gin.Context) {
					c.JSON(200, gin.H{
						"message": "pong",
					})
				})
			},
		},
		// Export essential providers for use in other modules
		Exports: []interface{}{
			events.NewEventBus,
			database.NewRepository,
			database.NewTransactionManager,
		},
	})
	return module
}

// OnModuleInit initializes the core module and its resources.
// This method is called when the application starts up.
//
// Parameters:
//   - ctx: The context for the initialization process
//
// Returns:
//   - error: Any error that occurred during initialization
func (m *CoreModule) OnModuleInit(ctx context.Context) error {
	// Initialize core resources (e.g., database connections)
	return nil
}

// OnModuleDestroy performs cleanup operations for the core module.
// This method is called when the application is shutting down.
//
// Parameters:
//   - ctx: The context for the cleanup process
//
// Returns:
//   - error: Any error that occurred during cleanup
func (m *CoreModule) OnModuleDestroy(ctx context.Context) error {
	// Cleanup core resources
	return nil
}
