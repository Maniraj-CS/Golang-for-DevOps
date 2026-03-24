package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// ================================
// TYPES
// ================================

type PodInfo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Status    string `json:"status"`
	Ready     string `json:"ready"`
	Restarts  int32  `json:"restarts"`
	Age       string `json:"age"`
	Node      string `json:"node"`
}

type DeploymentInfo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Ready     string `json:"ready"`
	UpToDate  int32  `json:"up_to_date"`
	Available int32  `json:"available"`
	Age       string `json:"age"`
}

type ServiceInfo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Type      string `json:"type"`
	ClusterIP string `json:"cluster_ip"`
	Ports     string `json:"ports"`
	Age       string `json:"age"`
}

// ================================
// K8S OPERATIONS
// ================================

func runK8s(config *Config) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Find kubeconfig
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error: Cannot find home directory:", err)
		return
	}

	kubeconfig := filepath.Join(home, ".kube", "config")

	// Check if kubeconfig exists
	if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
		fmt.Println("No kubeconfig found at", kubeconfig)
		fmt.Println("You need a K8s cluster. Try: minikube start")
		return
	}

	// Build config
	k8sConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Println("Error: Cannot load kubeconfig:", err)
		return
	}

	// Create clientset
	clientset, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		fmt.Println("Error: Cannot create K8s client:", err)
		return
	}

	// Route to subcommand
	switch config.SubCommand {
	case "pods":
		listPods(ctx, clientset, config)
	case "deployments", "deploy":
		listDeployments(ctx, clientset, config)
	case "services", "svc":
		listServices(ctx, clientset, config)
	default:
		fmt.Printf("Unknown K8s resource: %s\n", config.SubCommand)
		fmt.Println("Supported: pods, deployments, services")
	}
}

func listPods(ctx context.Context, clientset *kubernetes.Clientset, config *Config) {
	pods, err := clientset.CoreV1().Pods(config.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error listing pods:", err)
		return
	}

	var infos []PodInfo
	for _, pod := range pods.Items {
		info := PodInfo{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Status:    string(pod.Status.Phase),
			Node:      pod.Spec.NodeName,
		}

		// Ready count
		readyCount := 0
		totalCount := len(pod.Status.ContainerStatuses)
		for _, cs := range pod.Status.ContainerStatuses {
			if cs.Ready {
				readyCount++
			}
			info.Restarts += cs.RestartCount
		}
		info.Ready = fmt.Sprintf("%d/%d", readyCount, totalCount)

		// Age
		if pod.Status.StartTime != nil {
			info.Age = formatAge(time.Since(pod.Status.StartTime.Time))
		}

		infos = append(infos, info)
	}

	if config.Output == "json" {
		data, _ := json.MarshalIndent(infos, "", "  ")
		fmt.Println(string(data))
	} else {
		printPodTable(infos, config.Namespace)
	}
}

func listDeployments(ctx context.Context, clientset *kubernetes.Clientset, config *Config) {
	deployments, err := clientset.AppsV1().Deployments(config.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error listing deployments:", err)
		return
	}

	var infos []DeploymentInfo
	for _, d := range deployments.Items {
		replicas := int32(0)
		if d.Spec.Replicas != nil {
			replicas = *d.Spec.Replicas
		}

		info := DeploymentInfo{
			Name:      d.Name,
			Namespace: d.Namespace,
			Ready:     fmt.Sprintf("%d/%d", d.Status.ReadyReplicas, replicas),
			UpToDate:  d.Status.UpdatedReplicas,
			Available: d.Status.AvailableReplicas,
		}

		if !d.CreationTimestamp.IsZero() {
			info.Age = formatAge(time.Since(d.CreationTimestamp.Time))
		}

		infos = append(infos, info)
	}

	if config.Output == "json" {
		data, _ := json.MarshalIndent(infos, "", "  ")
		fmt.Println(string(data))
	} else {
		printDeploymentTable(infos, config.Namespace)
	}
}

func listServices(ctx context.Context, clientset *kubernetes.Clientset, config *Config) {
	services, err := clientset.CoreV1().Services(config.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error listing services:", err)
		return
	}

	var infos []ServiceInfo
	for _, svc := range services.Items {
		ports := ""
		for i, p := range svc.Spec.Ports {
			if i > 0 {
				ports += ", "
			}
			if p.NodePort != 0 {
				ports += fmt.Sprintf("%d:%d/%s", p.Port, p.NodePort, p.Protocol)
			} else {
				ports += fmt.Sprintf("%d/%s", p.Port, p.Protocol)
			}
		}

		info := ServiceInfo{
			Name:      svc.Name,
			Namespace: svc.Namespace,
			Type:      string(svc.Spec.Type),
			ClusterIP: svc.Spec.ClusterIP,
			Ports:     ports,
		}

		if !svc.CreationTimestamp.IsZero() {
			info.Age = formatAge(time.Since(svc.CreationTimestamp.Time))
		}

		infos = append(infos, info)
	}

	if config.Output == "json" {
		data, _ := json.MarshalIndent(infos, "", "  ")
		fmt.Println(string(data))
	} else {
		printServiceTable(infos, config.Namespace)
	}
}

// ================================
// DISPLAY
// ================================

func printPodTable(infos []PodInfo, namespace string) {
	ns := "all namespaces"
	if namespace != "" {
		ns = namespace
	}
	fmt.Printf("\n=== Kubernetes Pods (%s) ===\n\n", ns)

	if len(infos) == 0 {
		fmt.Println("No pods found.")
		return
	}

	maxName := 4
	for _, info := range infos {
		if len(info.Name) > maxName {
			maxName = len(info.Name)
		}
	}
	if maxName > 45 {
		maxName = 45
	}

	fmt.Printf("%-*s  %-12s  %-10s  %-5s  %8s  %s\n",
		maxName, "Name", "Namespace", "Status", "Ready", "Restarts", "Age")
	fmt.Printf("%-*s  %-12s  %-10s  %-5s  %8s  %s\n",
		maxName, strings.Repeat("-", maxName), "------------", "----------", "-----", "--------", "---")

	for _, info := range infos {
		name := info.Name
		if len(name) > 45 {
			name = name[:42] + "..."
		}

		icon := "🟢"
		if info.Status != "Running" && info.Status != "Succeeded" {
			icon = "🟡"
		}
		if info.Status == "Failed" || info.Status == "CrashLoopBackOff" {
			icon = "🔴"
		}

		fmt.Printf("%s %-*s  %-12s  %-10s  %-5s  %8d  %s\n",
			icon, maxName, name, info.Namespace, info.Status, info.Ready, info.Restarts, info.Age)
	}

	fmt.Printf("\nTotal: %d pods\n", len(infos))
}

func printDeploymentTable(infos []DeploymentInfo, namespace string) {
	ns := "all namespaces"
	if namespace != "" {
		ns = namespace
	}
	fmt.Printf("\n=== Kubernetes Deployments (%s) ===\n\n", ns)

	if len(infos) == 0 {
		fmt.Println("No deployments found.")
		return
	}

	maxName := 4
	for _, info := range infos {
		if len(info.Name) > maxName {
			maxName = len(info.Name)
		}
	}
	if maxName > 40 {
		maxName = 40
	}

	fmt.Printf("%-*s  %-12s  %-7s  %8s  %9s  %s\n",
		maxName, "Name", "Namespace", "Ready", "UpToDate", "Available", "Age")
	fmt.Printf("%-*s  %-12s  %-7s  %8s  %9s  %s\n",
		maxName, strings.Repeat("-", maxName), "------------", "-------", "--------", "---------", "---")

	for _, info := range infos {
		name := info.Name
		if len(name) > 40 {
			name = name[:37] + "..."
		}
		fmt.Printf("%-*s  %-12s  %-7s  %8d  %9d  %s\n",
			maxName, name, info.Namespace, info.Ready, info.UpToDate, info.Available, info.Age)
	}

	fmt.Printf("\nTotal: %d deployments\n", len(infos))
}

func printServiceTable(infos []ServiceInfo, namespace string) {
	ns := "all namespaces"
	if namespace != "" {
		ns = namespace
	}
	fmt.Printf("\n=== Kubernetes Services (%s) ===\n\n", ns)

	if len(infos) == 0 {
		fmt.Println("No services found.")
		return
	}

	maxName := 4
	for _, info := range infos {
		if len(info.Name) > maxName {
			maxName = len(info.Name)
		}
	}
	if maxName > 35 {
		maxName = 35
	}

	fmt.Printf("%-*s  %-12s  %-12s  %-16s  %-20s  %s\n",
		maxName, "Name", "Namespace", "Type", "ClusterIP", "Ports", "Age")
	fmt.Printf("%-*s  %-12s  %-12s  %-16s  %-20s  %s\n",
		maxName, strings.Repeat("-", maxName), "------------", "------------", "----------------", "--------------------", "---")

	for _, info := range infos {
		name := info.Name
		if len(name) > 35 {
			name = name[:32] + "..."
		}
		ports := info.Ports
		if len(ports) > 20 {
			ports = ports[:17] + "..."
		}
		fmt.Printf("%-*s  %-12s  %-12s  %-16s  %-20s  %s\n",
			maxName, name, info.Namespace, info.Type, info.ClusterIP, ports, info.Age)
	}

	fmt.Printf("\nTotal: %d services\n", len(infos))
}

// ================================
// HELPERS
// ================================

func formatAge(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	days := int(d.Hours()) / 24
	if days < 30 {
		return fmt.Sprintf("%dd", days)
	}
	months := days / 30
	return fmt.Sprintf("%dmo", months)
}
