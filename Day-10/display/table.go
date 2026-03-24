package display

// ================================
// Table Output for kubectl-status
// ================================
// This formats the cluster status as a nice, color-coded table.
// We use fmt.Printf with fixed widths for alignment.
// No external table library needed - keep it simple!

import (
	"fmt"

	"kubectl-status/kube"
)

// PrintTable prints the full cluster status as a formatted table
func PrintTable(status *kube.ClusterStatus) {
	// Main header
	PrintMainHeader(status.Context)

	// Pods section
	printPodsTable(status.Pods, status.Namespace)

	// Deployments section
	printDeploymentsTable(status.Deployments, status.Namespace)

	// Services section
	printServicesTable(status.Services, status.Namespace)

	// Summary section
	printSummary(status.Summary)
}

// printPodsTable prints the pods in a formatted table
func printPodsTable(pods []kube.PodInfo, namespace string) {
	PrintHeader(fmt.Sprintf("Pods (namespace: %s)", namespace))

	if len(pods) == 0 {
		fmt.Printf("  %sNo pods found%s\n", Gray, Reset)
		return
	}

	// Header row
	fmt.Printf("  %-2s %-44s %-14s %-12s %-6s %-10s %s\n",
		"", "Name", "Namespace", "Status", "Ready", "Restarts", "Age")
	fmt.Printf("  %-2s %-44s %-14s %-12s %-6s %-10s %s\n",
		"", "----", "---------", "------", "-----", "--------", "---")

	for _, p := range pods {
		// Get health status and icon
		health := kube.PodHealthStatus(p.Status, p.Restarts)
		icon := StatusIcon(health)

		// Color the status text based on health
		var statusText string
		switch health {
		case "healthy":
			statusText = Colorize(Green, p.Status)
		case "warning":
			statusText = Colorize(Yellow, p.Status)
		case "critical":
			statusText = Colorize(BoldRed, p.Status)
		default:
			statusText = p.Status
		}

		// Color restarts if too many
		restartText := fmt.Sprintf("%d", p.Restarts)
		if p.Restarts >= 10 {
			restartText = Colorize(BoldRed, restartText)
		} else if p.Restarts >= 5 {
			restartText = Colorize(Yellow, restartText)
		}

		// Truncate long pod names
		name := p.Name
		if len(name) > 42 {
			name = name[:39] + "..."
		}

		fmt.Printf("%s %-44s %-14s %-12s %-6s %-10s %s\n",
			icon, name, p.Namespace, statusText, p.Ready, restartText, p.Age)
	}
}

// printDeploymentsTable prints deployments in a formatted table
func printDeploymentsTable(deployments []kube.DeploymentInfo, namespace string) {
	PrintHeader(fmt.Sprintf("Deployments (namespace: %s)", namespace))

	if len(deployments) == 0 {
		fmt.Printf("  %sNo deployments found%s\n", Gray, Reset)
		return
	}

	// Header row
	fmt.Printf("%-3s %-24s %-14s %-8s %-10s %-10s %s\n",
		"", "Name", "Namespace", "Ready", "UpToDate", "Available", "Age")
	fmt.Printf("%-3s %-24s %-14s %-8s %-10s %-10s %s\n",
		"", "----", "---------", "-----", "--------", "---------", "---")

	for _, d := range deployments {
		icon := DeploymentIcon(d.Healthy)

		// Color the ready text
		var readyText string
		if d.Healthy {
			readyText = Colorize(Green, d.Ready)
		} else {
			readyText = Colorize(BoldRed, d.Ready)
		}

		// Truncate long names
		name := d.Name
		if len(name) > 22 {
			name = name[:19] + "..."
		}

		fmt.Printf("%s %-24s %-14s %-8s %-10d %-10d %s\n",
			icon, name, d.Namespace, readyText, d.UpToDate, d.Available, d.Age)
	}
}

// printServicesTable prints services in a formatted table
func printServicesTable(services []kube.ServiceInfo, namespace string) {
	PrintHeader(fmt.Sprintf("Services (namespace: %s)", namespace))

	if len(services) == 0 {
		fmt.Printf("  %sNo services found%s\n", Gray, Reset)
		return
	}

	// Header row
	fmt.Printf("  %-24s %-14s %-14s %-16s %-18s %s\n",
		"Name", "Namespace", "Type", "ClusterIP", "Ports", "Age")
	fmt.Printf("  %-24s %-14s %-14s %-16s %-18s %s\n",
		"----", "---------", "----", "---------", "-----", "---")

	for _, s := range services {
		// Truncate long names
		name := s.Name
		if len(name) > 22 {
			name = name[:19] + "..."
		}

		// Color service type
		var typeText string
		switch s.Type {
		case "LoadBalancer":
			typeText = Colorize(Cyan, s.Type)
		case "NodePort":
			typeText = Colorize(Yellow, s.Type)
		default:
			typeText = s.Type
		}

		fmt.Printf("  %-24s %-14s %-14s %-16s %-18s %s\n",
			name, s.Namespace, typeText, s.ClusterIP, s.Ports, s.Age)
	}
}

// printSummary prints the cluster summary
func printSummary(summary kube.Summary) {
	fmt.Printf("\n%s=== Summary ===%s\n", BoldWhite, Reset)

	// Pods summary
	podDetails := fmt.Sprintf("%d total", summary.TotalPods)
	if summary.RunningPods > 0 {
		podDetails += fmt.Sprintf(" (%s%d running%s", Green, summary.RunningPods, Reset)
	} else {
		podDetails += " (0 running"
	}
	if summary.CrashLoopPods > 0 {
		podDetails += fmt.Sprintf(", %s%d crash-loop%s", BoldRed, summary.CrashLoopPods, Reset)
	}
	if summary.PendingPods > 0 {
		podDetails += fmt.Sprintf(", %s%d pending%s", Yellow, summary.PendingPods, Reset)
	}
	if summary.FailedPods > 0 {
		podDetails += fmt.Sprintf(", %s%d failed%s", Red, summary.FailedPods, Reset)
	}
	podDetails += ")"
	fmt.Printf("Pods:        %s\n", podDetails)

	// Deployments summary
	depDetails := fmt.Sprintf("%d total", summary.TotalDeployments)
	if summary.TotalDeployments > 0 {
		depDetails += fmt.Sprintf(" (%s%d healthy%s", Green, summary.HealthyDeploys, Reset)
		if summary.UnhealthyDeploys > 0 {
			depDetails += fmt.Sprintf(", %s%d unhealthy%s", BoldRed, summary.UnhealthyDeploys, Reset)
		}
		depDetails += ")"
	}
	fmt.Printf("Deployments: %s\n", depDetails)

	// Services summary
	fmt.Printf("Services:    %d total\n", summary.TotalServices)

	PrintSeparator()
}
