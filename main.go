package main

import (
	goblin "github.com/calummacc/goblin/internal/core"
)

func main() {
	app := goblin.New(
		goblin.WithPort(":8081"),
	)
	app.Run()
}
