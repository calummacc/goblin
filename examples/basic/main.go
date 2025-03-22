package main

import (
	"context"
	"log"

	"github.com/calummacc/goblin/internal/core"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create new application with custom configuration
	app := core.NewGoblinApplication(
		core.WithPort(3000),
		core.WithHost("0.0.0.0"),
		core.WithGinMode(gin.ReleaseMode),
	)

	// Add modules
	appModule := NewAppModule()
	app.AddModule(appModule)

	// Configure application
	app.Configure()

	// Run application
	if err := app.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
