package kube

// ================================
// Kubernetes Client Setup
// ================================
// This file handles connecting to a Kubernetes cluster.
// It uses the same kubeconfig file that kubectl uses (~/.kube/config).
// So if kubectl works on your machine, this will also work!

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client wraps the Kubernetes clientset and config info
type Client struct {
	Clientset  *kubernetes.Clientset
	RestConfig *rest.Config
	Context    string // current context name (like "minikube", "prod-cluster", etc.)
}

// NewClient creates a new Kubernetes client.
// It loads kubeconfig the same way kubectl does:
//  1. If KUBECONFIG env var is set, use that
//  2. Otherwise, use ~/.kube/config
//
// This is the standard way to connect to K8s from outside the cluster.
func NewClient() (*Client, error) {
	// Step 1: Find the kubeconfig file
	kubeconfigPath := getKubeconfigPath()

	// Step 2: Build the config loading rules
	loadingRules := &clientcmd.ClientConfigLoadingRules{
		ExplicitPath: kubeconfigPath,
	}
	configOverrides := &clientcmd.ConfigOverrides{}

	// Step 3: Load the kubeconfig
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		loadingRules,
		configOverrides,
	)

	// Step 4: Get the current context name (e.g., "minikube", "prod-cluster")
	rawConfig, err := kubeConfig.RawConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read kubeconfig: %w", err)
	}
	currentContext := rawConfig.CurrentContext

	// Step 5: Get the REST config for making API calls
	restConfig, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to build rest config: %w", err)
	}

	// Step 6: Create the clientset (this is what we use to make API calls)
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	return &Client{
		Clientset:  clientset,
		RestConfig: restConfig,
		Context:    currentContext,
	}, nil
}

// getKubeconfigPath returns the path to the kubeconfig file.
// Checks KUBECONFIG env var first, then falls back to ~/.kube/config.
func getKubeconfigPath() string {
	// Check if KUBECONFIG environment variable is set
	if kubeconfig := os.Getenv("KUBECONFIG"); kubeconfig != "" {
		return kubeconfig
	}

	// Default: ~/.kube/config
	home, err := os.UserHomeDir()
	if err != nil {
		// If we can't get home dir, return the default path anyway
		return filepath.Join("~", ".kube", "config")
	}
	return filepath.Join(home, ".kube", "config")
}
