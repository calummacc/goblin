package main

import (
	"os"

	"github.com/calummacc/goblin/cmd"
	"github.com/calummacc/goblin/internal/core"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "new" {
		cmd.Execute()
	} else {
		app := core.NewApp()
		app.Run()
	}
}
