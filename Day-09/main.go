package main

// ================================
// infraview - Infrastructure Viewer
// ================================
// A tool that shows Docker containers and K8s pods in one clean view.
//
// Usage:
//   go run . docker                    List Docker containers
//   go run . docker --stats            Show CPU/memory stats
//   go run . docker --images           List Docker images
//   go run . k8s pods                  List K8s pods
//   go run . k8s deployments           List K8s deployments
//   go run . k8s services              List K8s services
//   go run . all                       Show everything
//   go run . docker -o json            JSON output

import (
	"fmt"
	"os"
	"strings"

	_ "go.uber.org/automaxprocs/maxprocs"
)

// Config holds CLI options
type Config struct {
	Command    string // "docker", "k8s", "all"
	SubCommand string // "pods", "deployments", "services", "images"
	Output     string // "table", "json"
	Namespace  string // k8s namespace filter
	ShowStats  bool   // show container CPU/memory
	ShowImages bool   // show images instead of containers
	Filter     string // filter by status
}

func main() {
	config := parseArgs()

	switch config.Command {
	case "docker":
		runDocker(config)
	case "k8s":
		runK8s(config)
	case "all":
		runDocker(config)
		fmt.Println()
		runK8s(config)
	default:
		printUsage()
	}
}

func parseArgs() *Config {
	config := &Config{
		Output:    "table",
		Namespace: "",
	}

	args := os.Args[1:]
	if len(args) == 0 {
		printUsage()
		os.Exit(0)
	}

	config.Command = args[0]

	// Parse rest of args
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "-o", "--output":
			if i+1 < len(args) {
				i++
				config.Output = args[i]
			}
		case "-n", "--namespace":
			if i+1 < len(args) {
				i++
				config.Namespace = args[i]
			}
		case "--stats":
			config.ShowStats = true
		case "--images":
			config.ShowImages = true
		case "--filter":
			if i+1 < len(args) {
				i++
				config.Filter = strings.ToLower(args[i])
			}
		case "--help", "-h":
			printUsage()
			os.Exit(0)
		default:
			// Could be a subcommand for k8s
			if config.Command == "k8s" && config.SubCommand == "" {
				config.SubCommand = args[i]
			}
		}
	}

	// Default subcommand for k8s
	if config.Command == "k8s" && config.SubCommand == "" {
		config.SubCommand = "pods"
	}

	return config
}

func printUsage() {
	fmt.Println("infraview - Infrastructure Viewer")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run . docker [flags]              List Docker containers")
	fmt.Println("  go run . docker --stats              Show CPU/memory stats")
	fmt.Println("  go run . docker --images             List Docker images")
	fmt.Println("  go run . k8s pods [flags]            List K8s pods")
	fmt.Println("  go run . k8s deployments [flags]     List K8s deployments")
	fmt.Println("  go run . k8s services [flags]        List K8s services")
	fmt.Println("  go run . all                         Show everything")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -o, --output table|json     Output format (default: table)")
	fmt.Println("  -n, --namespace NAME        K8s namespace (default: all)")
	fmt.Println("  --stats                     Show container CPU/memory")
	fmt.Println("  --images                    Show images instead of containers")
	fmt.Println("  --filter running|stopped    Filter containers by state")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run . docker")
	fmt.Println("  go run . docker --stats -o json")
	fmt.Println("  go run . k8s pods -n default")
	fmt.Println("  go run . k8s deployments")
	fmt.Println("  go run . all")
}
