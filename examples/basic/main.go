package main

import (
	"context"
	"log"

	"github.com/calummacc/goblin/internal/core"
)

func main() {

	// Create new application
	app := core.NewApplication()

	// Add root module
	appModule := NewAppModule()
	app.AddModule(appModule)

	// Configure application
	app.Configure()

	// Run application
	if err := app.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
