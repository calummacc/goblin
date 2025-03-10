package main

import (
	"github.com/calummacc/goblin/examples/user"
	goblin "github.com/calummacc/goblin/internal/core"
)

func main() {
	app := goblin.New(
		goblin.WithPort(":8081"),
		user.Module(), // Include user module
	)
	app.Run()
}
