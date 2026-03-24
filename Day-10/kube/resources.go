package kube

// ================================
// Fetch Kubernetes Resources Concurrently
// ================================
// This is where goroutines + WaitGroup shine!
// Instead of fetching pods, then deployments, then services one by one,
// we fetch ALL of them at the same time. This makes the plugin fast.
// This is exactly what we learned in Day-05.

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FetchClusterStatus fetches pods, deployments, and services concurrently.
// If namespace is empty (""), it fetches from ALL namespaces.
func (c *Client) FetchClusterStatus(namespace string) (*ClusterStatus, error) {
	status := &ClusterStatus{
		Context:   c.Context,
		Namespace: namespace,
	}

	// If no namespace specified, show "all"
	if namespace == "" {
		status.Namespace = "all"
	}

	// We use a mutex to safely write to the status struct from multiple goroutines
	var mu sync.Mutex
	// Collect errors from goroutines
	var fetchErrors []string

	// Create a context with timeout - don't hang forever if cluster is unreachable
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// ================================
	// Launch 3 goroutines concurrently!
	// ================================
	var wg sync.WaitGroup

	// Goroutine 1: Fetch Pods
	wg.Add(1)
	go func() {
		defer wg.Done()

		pods, err := c.fetchPods(ctx, namespace)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			fetchErrors = append(fetchErrors, fmt.Sprintf("pods: %v", err))
			return
		}
		status.Pods = pods
	}()

	// Goroutine 2: Fetch Deployments
	wg.Add(1)
	go func() {
		defer wg.Done()

		deployments, err := c.fetchDeployments(ctx, namespace)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			fetchErrors = append(fetchErrors, fmt.Sprintf("deployments: %v", err))
			return
		}
		status.Deployments = deployments
	}()

	// Goroutine 3: Fetch Services
	wg.Add(1)
	go func() {
		defer wg.Done()

		services, err := c.fetchServices(ctx, namespace)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			fetchErrors = append(fetchErrors, fmt.Sprintf("services: %v", err))
			return
		}
		status.Services = services
	}()

	// Wait for all goroutines to finish
	wg.Wait()

	// If ALL fetches failed, return error
	if len(fetchErrors) == 3 {
		return nil, fmt.Errorf("failed to fetch any resources: %s", strings.Join(fetchErrors, "; "))
	}

	// Calculate summary
	status.Summary = calculateSummary(status)

	return status, nil
}

// fetchPods gets all pods in a namespace (or all namespaces if namespace is "")
func (c *Client) fetchPods(ctx context.Context, namespace string) ([]PodInfo, error) {
	podList, err := c.Clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	var pods []PodInfo
	for _, p := range podList.Items {
		// Figure out the pod status
		podStatus := string(p.Status.Phase)

		// Check for more specific statuses from container statuses
		for _, cs := range p.Status.ContainerStatuses {
			if cs.State.Waiting != nil {
				podStatus = cs.State.Waiting.Reason // e.g., "CrashLoopBackOff", "ImagePullBackOff"
			}
			if cs.State.Terminated != nil && cs.State.Terminated.Reason != "" {
				podStatus = cs.State.Terminated.Reason
			}
		}

		// Check if pod is being deleted
		if p.DeletionTimestamp != nil {
			podStatus = "Terminating"
		}

		// Calculate ready containers (e.g., "1/2")
		readyCount := int32(0)
		totalContainers := int32(len(p.Spec.Containers))
		for _, cs := range p.Status.ContainerStatuses {
			if cs.Ready {
				readyCount++
			}
		}
		ready := fmt.Sprintf("%d/%d", readyCount, totalContainers)

		// Total restarts across all containers
		var restarts int32
		for _, cs := range p.Status.ContainerStatuses {
			restarts += cs.RestartCount
		}

		pods = append(pods, PodInfo{
			Name:      p.Name,
			Namespace: p.Namespace,
			Status:    podStatus,
			Ready:     ready,
			Restarts:  restarts,
			Age:       formatAge(p.CreationTimestamp.Time),
			Node:      p.Spec.NodeName,
		})
	}

	return pods, nil
}

// fetchDeployments gets all deployments in a namespace
func (c *Client) fetchDeployments(ctx context.Context, namespace string) ([]DeploymentInfo, error) {
	depList, err := c.Clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}

	var deployments []DeploymentInfo
	for _, d := range depList.Items {
		desired := *d.Spec.Replicas
		available := d.Status.AvailableReplicas
		upToDate := d.Status.UpdatedReplicas

		ready := fmt.Sprintf("%d/%d", available, desired)
		healthy := available == desired

		deployments = append(deployments, DeploymentInfo{
			Name:      d.Name,
			Namespace: d.Namespace,
			Ready:     ready,
			UpToDate:  upToDate,
			Available: available,
			Age:       formatAge(d.CreationTimestamp.Time),
			Healthy:   healthy,
		})
	}

	return deployments, nil
}

// fetchServices gets all services in a namespace
func (c *Client) fetchServices(ctx context.Context, namespace string) ([]ServiceInfo, error) {
	svcList, err := c.Clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	var services []ServiceInfo
	for _, s := range svcList.Items {
		// Format ports like "80/TCP" or "80:30080/TCP"
		var portStrs []string
		for _, port := range s.Spec.Ports {
			if port.NodePort > 0 {
				portStrs = append(portStrs, fmt.Sprintf("%d:%d/%s", port.Port, port.NodePort, port.Protocol))
			} else {
				portStrs = append(portStrs, fmt.Sprintf("%d/%s", port.Port, port.Protocol))
			}
		}
		ports := strings.Join(portStrs, ", ")
		if ports == "" {
			ports = "<none>"
		}

		services = append(services, ServiceInfo{
			Name:      s.Name,
			Namespace: s.Namespace,
			Type:      string(s.Spec.Type),
			ClusterIP: s.Spec.ClusterIP,
			Ports:     ports,
			Age:       formatAge(s.CreationTimestamp.Time),
		})
	}

	return services, nil
}

// calculateSummary builds a quick summary from the fetched data
func calculateSummary(status *ClusterStatus) Summary {
	s := Summary{
		TotalPods:        len(status.Pods),
		TotalDeployments: len(status.Deployments),
		TotalServices:    len(status.Services),
	}

	for _, p := range status.Pods {
		switch {
		case p.Status == "Running":
			s.RunningPods++
		case p.Status == "Pending":
			s.PendingPods++
		case p.Status == "CrashLoopBackOff" || p.Status == "Error" || p.Status == "Failed":
			s.CrashLoopPods++
		case p.Status == "Succeeded":
			// completed pods (like Jobs)
		default:
			s.FailedPods++
		}
	}

	for _, d := range status.Deployments {
		if d.Healthy {
			s.HealthyDeploys++
		} else {
			s.UnhealthyDeploys++
		}
	}

	return s
}

// formatAge converts a time to a human-readable age string like "2d", "5h", "3m"
func formatAge(created time.Time) string {
	duration := time.Since(created)

	switch {
	case duration.Hours() >= 24*365:
		years := int(duration.Hours() / (24 * 365))
		return fmt.Sprintf("%dy", years)
	case duration.Hours() >= 24:
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%dd", days)
	case duration.Hours() >= 1:
		hours := int(duration.Hours())
		return fmt.Sprintf("%dh", hours)
	case duration.Minutes() >= 1:
		mins := int(duration.Minutes())
		return fmt.Sprintf("%dm", mins)
	default:
		secs := int(duration.Seconds())
		return fmt.Sprintf("%ds", secs)
	}
}
