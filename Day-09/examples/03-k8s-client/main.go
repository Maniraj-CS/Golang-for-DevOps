package main

// ================================
// EXAMPLE 03 - Kubernetes client-go Basics
// ================================
// This shows how to connect to a K8s cluster and list resources.
// Run: go run main.go
// Needs a running K8s cluster (minikube, kind, or real cluster)!

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// ================================
	// 1. CONNECT TO K8S CLUSTER
	// ================================

	fmt.Println("=== Connecting to Kubernetes ===")

	// Find kubeconfig file
	home, _ := os.UserHomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")

	// Check if kubeconfig exists
	if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
		fmt.Println("No kubeconfig found at", kubeconfig)
		fmt.Println("You need a K8s cluster to run this example.")
		fmt.Println("Try: minikube start")
		os.Exit(0)
	}

	// Build config from kubeconfig file
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal("Cannot load kubeconfig:", err)
	}

	// Create the clientset (this is your K8s API client)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal("Cannot create K8s client:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("Connected to cluster:", config.Host)

	// ================================
	// 2. LIST PODS (all namespaces)
	// ================================

	fmt.Println("\n=== Pods ===")

	pods, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Println("Cannot list pods:", err)
	} else {
		if len(pods.Items) == 0 {
			fmt.Println("No pods found.")
		}
		for _, pod := range pods.Items {
			age := ""
			if pod.Status.StartTime != nil {
				age = time.Since(pod.Status.StartTime.Time).Round(time.Second).String()
			}
			fmt.Printf("  %-40s  %-12s  %-10s  %s\n",
				pod.Name, pod.Namespace, pod.Status.Phase, age)
		}
		fmt.Printf("\nTotal: %d pods\n", len(pods.Items))
	}

	// ================================
	// 3. LIST DEPLOYMENTS
	// ================================

	fmt.Println("\n=== Deployments ===")

	deployments, err := clientset.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Println("Cannot list deployments:", err)
	} else {
		if len(deployments.Items) == 0 {
			fmt.Println("No deployments found.")
		}
		for _, d := range deployments.Items {
			replicas := int32(0)
			if d.Spec.Replicas != nil {
				replicas = *d.Spec.Replicas
			}
			fmt.Printf("  %-40s  %-12s  %d/%d ready\n",
				d.Name, d.Namespace, d.Status.ReadyReplicas, replicas)
		}
		fmt.Printf("\nTotal: %d deployments\n", len(deployments.Items))
	}

	// ================================
	// 4. LIST SERVICES
	// ================================

	fmt.Println("\n=== Services ===")

	services, err := clientset.CoreV1().Services("").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Println("Cannot list services:", err)
	} else {
		if len(services.Items) == 0 {
			fmt.Println("No services found.")
		}
		for _, svc := range services.Items {
			ports := ""
			for i, p := range svc.Spec.Ports {
				if i > 0 {
					ports += ", "
				}
				ports += fmt.Sprintf("%d/%s", p.Port, p.Protocol)
			}
			fmt.Printf("  %-30s  %-12s  %-12s  %s\n",
				svc.Name, svc.Namespace, svc.Spec.Type, ports)
		}
		fmt.Printf("\nTotal: %d services\n", len(services.Items))
	}

	// ================================
	// KEY TAKEAWAYS
	// ================================

	fmt.Println("\n=== Key Takeaways ===")
	fmt.Println("1. Load kubeconfig from ~/.kube/config")
	fmt.Println("2. Use kubernetes.NewForConfig() to create a client")
	fmt.Println("3. CoreV1() for pods, services; AppsV1() for deployments")
	fmt.Println("4. Use \"\" for all namespaces, or specify one like \"default\"")
	fmt.Println("5. Always use context with timeout")
}
