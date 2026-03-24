package main

// ================================
// EXAMPLE 01 - Basic Cobra CLI
// ================================
// This shows how Cobra works at the simplest level.
// Run: go run main.go
// Run: go run main.go greet
// Run: go run main.go greet --name Abhishek
// Run: go run main.go greet --name Abhishek --shout

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	// ================================
	// 1. ROOT COMMAND
	// ================================
	// This is what runs when you type just the app name

	rootCmd := &cobra.Command{
		Use:   "myapp",
		Short: "A simple demo CLI app",
		Long:  "This is a demo app to show how Cobra works.\nIt has subcommands, flags, and help text - all auto-generated!",
		// If someone runs just "myapp" with no subcommand:
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Welcome to myapp! Use --help to see available commands.")
		},
	}

	// ================================
	// 2. GREET SUBCOMMAND
	// ================================

	var name string
	var shout bool

	greetCmd := &cobra.Command{
		Use:   "greet",
		Short: "Greet someone",
		Long:  "Print a greeting message. Use --name to specify who to greet.",
		Run: func(cmd *cobra.Command, args []string) {
			msg := fmt.Sprintf("Hello, %s! Welcome to DevOps with Go.", name)
			if shout {
				msg = strings.ToUpper(msg)
			}
			fmt.Println(msg)
		},
	}

	// Add flags to greet command
	greetCmd.Flags().StringVar(&name, "name", "World", "Name of the person to greet")
	greetCmd.Flags().BoolVar(&shout, "shout", false, "SHOUT the greeting")

	// ================================
	// 3. INFO SUBCOMMAND
	// ================================

	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Show app info",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("App: myapp")
			fmt.Println("Version: 0.1.0")
			fmt.Println("Built with Cobra")
		},
	}

	// ================================
	// 4. ECHO SUBCOMMAND (takes arguments)
	// ================================

	echoCmd := &cobra.Command{
		Use:   "echo [message...]",
		Short: "Echo back whatever you type",
		Args:  cobra.MinimumNArgs(1), // at least 1 argument required
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("You said:", strings.Join(args, " "))
		},
	}

	// ================================
	// 5. WIRE IT ALL TOGETHER
	// ================================

	rootCmd.AddCommand(greetCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(echoCmd)

	// Execute!
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
