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
		cmd.HelpFunc()(cmd, args) // Hiển thị help nếu không có subcommand
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
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
