package display

// ================================
// JSON Output for kubectl-status
// ================================
// When the user passes `-o json`, we output clean JSON.
// This is super useful for piping to jq, or for automation.

import (
	"encoding/json"
	"fmt"
	"os"

	"kubectl-status/kube"
)

// PrintJSON prints the cluster status as formatted JSON
func PrintJSON(status *kube.ClusterStatus) {
	// MarshalIndent gives us pretty-printed JSON
	data, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to marshal JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(data))
}
