package main

import (
	"os"

	"github.com/calummacc/goblin/cmd"
	"github.com/calummacc/goblin/internal/core"
)

func main() {
	// Kiểm tra xem có phải lệnh CLI không
	if len(os.Args) >= 2 {
		cmd.Execute() // Chạy CLI (ví dụ: --help, new, ...)
	} else {
		// Chạy app nếu không có tham số
		app := core.NewApp()
		app.Run()
	}
}
