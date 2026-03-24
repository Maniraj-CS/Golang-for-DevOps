package cmd

// ================================
// Root Command for kubectl-status
// ================================
// This sets up the Cobra root command, flags, and version info.
// Cobra gives us automatic --help, flag parsing, and subcommand support.
// Same thing we learned in Day-08!

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version info - set during build with -ldflags
var (
	Version   = "dev"
	BuildTime = "unknown"
)

// Global flags
var (
	namespace  string
	outputFmt  string
)

// rootCmd is the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubectl-status",
	Short: "Get a quick overview of your Kubernetes cluster",
	Long: `kubectl-status gives you a single-command overview of your Kubernetes cluster.

It shows pods, deployments, and services in one clean, color-coded view.
No more running three separate kubectl commands!

Works as a standalone tool or as a kubectl plugin:
  kubectl status            # as a plugin
  kubectl-status            # standalone
  kubectl-status -n default # filter by namespace`,

	// Run the status command by default (no subcommands needed)
	RunE: runStatus,
}

// Execute runs the root command - called from main.go
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Add flags to the root command
	rootCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Kubernetes namespace to filter (default: all namespaces)")
	rootCmd.Flags().StringVarP(&outputFmt, "output", "o", "table", "Output format: table or json")

	// Version flag
	rootCmd.Version = fmt.Sprintf("%s (built: %s)", Version, BuildTime)
}
