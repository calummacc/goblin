package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Root command
var rootCmd = &cobra.Command{
	Use:   "goblin",
	Short: "CLI tool for Goblin Framework",
	Run: func(cmd *cobra.Command, args []string) {
		// Hiển thị help nếu không có subcommand
		cmd.Help()
	},
}

// Command "new"
var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new Goblin project",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[1]
		fmt.Printf("Creating project %s...\n", projectName)
		// Thêm logic tạo project
	},
}

// Command "run"
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the Goblin application",
	Run: func(cmd *cobra.Command, args []string) {
		app := core.NewApp()
		app.Run()
	},
}

func init() {
	rootCmd.AddCommand(newCmd, runCmd) // Thêm các subcommand
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
