package cmd

// ================================
// Status Command Logic
// ================================
// This is the main brain of the plugin.
// It connects to K8s, fetches data, and displays output.
// All the pieces from the course come together here!

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"kubectl-status/display"
	"kubectl-status/kube"
)

// runStatus is the main function that runs when the user executes the command
func runStatus(cmd *cobra.Command, args []string) error {
	// Step 1: Validate flags
	if outputFmt != "table" && outputFmt != "json" {
		return fmt.Errorf("invalid output format %q: must be 'table' or 'json'", outputFmt)
	}

	// Step 2: Create Kubernetes client
	client, err := kube.NewClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n%s❌ Cannot connect to Kubernetes cluster%s\n\n", display.BoldRed, display.Reset)
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		fmt.Fprintf(os.Stderr, "Make sure:\n")
		fmt.Fprintf(os.Stderr, "  • kubectl is configured and working\n")
		fmt.Fprintf(os.Stderr, "  • Your kubeconfig is at ~/.kube/config or KUBECONFIG is set\n")
		fmt.Fprintf(os.Stderr, "  • Your cluster is reachable\n\n")
		fmt.Fprintf(os.Stderr, "Try running: kubectl cluster-info\n\n")
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	// Step 3: Fetch cluster status (pods, deployments, services - all concurrently!)
	status, err := client.FetchClusterStatus(namespace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n%s❌ Failed to fetch cluster status%s\n\n", display.BoldRed, display.Reset)
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		fmt.Fprintf(os.Stderr, "The cluster might be unreachable or you might not have permission.\n")
		fmt.Fprintf(os.Stderr, "Try running: kubectl get pods\n\n")
		return fmt.Errorf("failed to fetch cluster status: %w", err)
	}

	// Step 4: Display output in the requested format
	switch outputFmt {
	case "json":
		display.PrintJSON(status)
	default:
		display.PrintTable(status)
	}

	return nil
}
