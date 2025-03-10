package main

import (
	"github.com/calummacc/goblin/examples/middlewares"
	"github.com/calummacc/goblin/examples/user"
	goblin "github.com/calummacc/goblin/internal/core"
)

func main() {
	//Create module
	module := user.NewModule()
	app := goblin.New(
		goblin.WithPort(":8081"),
		module.Provide(), // Provide user module to fx.Options
	)
	app.AddGlobalMiddleware(
		middlewares.RequestIDMiddleware(),
		middlewares.CORSMiddleware(),
		middlewares.LoggerMiddleware(),
	)
	app.Run()
}
