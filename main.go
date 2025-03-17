// goblin/main.go
package main

import (
	"log"

	"github.com/calummacc/goblin/core"
	"github.com/calummacc/goblin/database"
	"github.com/calummacc/goblin/events"
	"github.com/calummacc/goblin/examples/user_module"
	"github.com/calummacc/goblin/middleware"
	"github.com/calummacc/goblin/router"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func main() {
	// Create a new Goblin app
	app := core.NewGoblinApp(core.GoblinAppOptions{
		Debug: true,
		Modules: []core.GoblinModule{
			// Core modules
			CoreModule(),
			// Feature modules
			user_module.Module(),
		},
	})

	// Start the app
	if err := app.Start(":8080"); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}

// CoreModule creates the core module
func CoreModule() core.GoblinModule {
	return core.NewModule("CoreModule", fx.Options(
		// Provide core services
		fx.Provide(func() *gin.Engine {
			engine := gin.Default()
			return engine
		}),
		fx.Provide(events.NewEventBus),
		fx.Provide(func() (*database.ORM, error) {
			return database.NewORM(database.Config{
				Driver:   "sqlite",
				Database: "goblin.db",
			})
		}),
		fx.Provide(database.NewRepository),
		fx.Provide(database.NewTransactionManager),
		fx.Provide(router.NewRouter),

		// Set up global middleware
		fx.Invoke(func(engine *gin.Engine) {
			engine.Use(middleware.Logger())
			engine.Use(middleware.Recovery())
		}),

		// Auto-migrate database models
		fx.Invoke(func(orm *database.ORM) {
			if err := orm.AutoMigrate(&user_module.User{}); err != nil {
				log.Fatalf("Failed to migrate database: %v", err)
			}
		}),

		// Register controllers
		fx.Invoke(func(router *router.RouterRegistry, userController *user_module.UserController) {
			for _, route := range userController.Routes() {
				router.RegisterController(userController.BasePath(), route)
			}
		}),
	))
}
