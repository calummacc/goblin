package main

import (
	goblin "github.com/calummacc/goblin/internal/core"
)

func main() {
	app := goblin.NewGoblinApp(
		goblin.WithPort(":8081"),
	)

	app.AddGlobalMiddleware()
	app.Run()
}
