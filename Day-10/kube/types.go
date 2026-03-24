package kube

// ================================
// Types for clean data representation
// ================================
// These are our own types - simpler than the K8s API types.
// We convert K8s types to these for easy display and JSON output.

// ClusterStatus holds everything we want to show
type ClusterStatus struct {
	Context     string           `json:"context"`
	Namespace   string           `json:"namespace"`
	Pods        []PodInfo        `json:"pods"`
	Deployments []DeploymentInfo `json:"deployments"`
	Services    []ServiceInfo    `json:"services"`
	Summary     Summary          `json:"summary"`
}

// PodInfo is our simplified pod representation
type PodInfo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Status    string `json:"status"`
	Ready     string `json:"ready"`
	Restarts  int32  `json:"restarts"`
	Age       string `json:"age"`
	Node      string `json:"node,omitempty"`
}

// DeploymentInfo is our simplified deployment representation
type DeploymentInfo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Ready     string `json:"ready"`
	UpToDate  int32  `json:"up_to_date"`
	Available int32  `json:"available"`
	Age       string `json:"age"`
	Healthy   bool   `json:"healthy"`
}

// ServiceInfo is our simplified service representation
type ServiceInfo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Type      string `json:"type"`
	ClusterIP string `json:"cluster_ip"`
	Ports     string `json:"ports"`
	Age       string `json:"age"`
}

// Summary gives a quick overview of the cluster health
type Summary struct {
	TotalPods        int `json:"total_pods"`
	RunningPods      int `json:"running_pods"`
	PendingPods      int `json:"pending_pods"`
	FailedPods       int `json:"failed_pods"`
	CrashLoopPods    int `json:"crashloop_pods"`
	TotalDeployments int `json:"total_deployments"`
	HealthyDeploys   int `json:"healthy_deployments"`
	UnhealthyDeploys int `json:"unhealthy_deployments"`
	TotalServices    int `json:"total_services"`
}

// PodHealthStatus returns a simple health indicator
func PodHealthStatus(status string, restarts int32) string {
	switch {
	case status == "Running" && restarts < 5:
		return "healthy"
	case status == "Running" && restarts >= 5:
		return "warning" // running but has restarted a lot
	case status == "Succeeded":
		return "healthy"
	case status == "Pending":
		return "warning"
	case status == "CrashLoopBackOff" || status == "Error" || status == "Failed":
		return "critical"
	case status == "Terminating":
		return "warning"
	default:
		return "unknown"
	}
}

// DeploymentHealthStatus checks if a deployment is healthy
func DeploymentHealthStatus(ready string, available, desired int32) string {
	if desired == 0 {
		return "warning"
	}
	if available == desired {
		return "healthy"
	}
	if available == 0 {
		return "critical"
	}
	return "warning"
}
