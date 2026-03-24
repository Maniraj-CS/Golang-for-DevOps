package main

// ================================
// EXAMPLE 02 - Docker SDK Basics
// ================================
// This shows how to connect to Docker and list containers/images.
// Run: go run main.go
// Needs Docker running on your machine!

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

func main() {
	// ================================
	// 1. CONNECT TO DOCKER
	// ================================

	fmt.Println("=== Connecting to Docker ===")

	// Create a context with timeout - don't wait forever
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to the local Docker daemon
	// FromEnv reads DOCKER_HOST, DOCKER_TLS_VERIFY, etc.
	// WithAPIVersionNegotiation automatically uses the right API version
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal("Cannot create Docker client:", err)
	}
	defer cli.Close()

	// Test connection by pinging Docker
	ping, err := cli.Ping(ctx)
	if err != nil {
		log.Fatal("Cannot connect to Docker. Is Docker running?\nError:", err)
	}
	fmt.Println("Connected to Docker!")
	fmt.Println("API Version:", ping.APIVersion)

	// ================================
	// 2. LIST CONTAINERS
	// ================================

	fmt.Println("\n=== Docker Containers ===")

	containers, err := cli.ContainerList(ctx, container.ListOptions{
		All: true, // include stopped containers too
	})
	if err != nil {
		log.Fatal("Cannot list containers:", err)
	}

	if len(containers) == 0 {
		fmt.Println("No containers found. Try running: docker run -d nginx")
	}

	running := 0
	stopped := 0

	for _, c := range containers {
		// Container ID is long - show first 12 chars
		id := c.ID[:12]

		// Image name
		img := c.Image

		// Status
		status := c.Status

		// Ports
		ports := formatPorts(c.Ports)

		// Names (Docker adds a leading /)
		name := ""
		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		}

		if c.State == "running" {
			running++
		} else {
			stopped++
		}

		fmt.Printf("  %s  %-25s  %-20s  %-15s  %s\n", id, img, status, ports, name)
	}

	fmt.Printf("\nRunning: %d | Stopped: %d | Total: %d\n", running, stopped, len(containers))

	// ================================
	// 3. LIST IMAGES
	// ================================

	fmt.Println("\n=== Docker Images ===")

	images, err := cli.ImageList(ctx, image.ListOptions{})
	if err != nil {
		log.Fatal("Cannot list images:", err)
	}

	if len(images) == 0 {
		fmt.Println("No images found.")
	}

	for _, img := range images {
		tags := "<none>"
		if len(img.RepoTags) > 0 {
			tags = strings.Join(img.RepoTags, ", ")
		}

		sizeMB := float64(img.Size) / 1024 / 1024

		fmt.Printf("  %-45s  %6.0f MB\n", tags, sizeMB)
	}

	fmt.Printf("\nTotal: %d images\n", len(images))

	// ================================
	// KEY TAKEAWAYS
	// ================================

	fmt.Println("\n=== Key Takeaways ===")
	fmt.Println("1. Use client.NewClientWithOpts(client.FromEnv) to connect")
	fmt.Println("2. Always use context with timeout for API calls")
	fmt.Println("3. ContainerList with All:true includes stopped containers")
	fmt.Println("4. The Docker SDK mirrors the Docker CLI commands")
}

func formatPorts(ports []types.Port) string {
	if len(ports) == 0 {
		return "-"
	}

	var parts []string
	for _, p := range ports {
		if p.PublicPort != 0 {
			parts = append(parts, fmt.Sprintf("%d->%d/%s", p.PublicPort, p.PrivatePort, p.Type))
		} else {
			parts = append(parts, fmt.Sprintf("%d/%s", p.PrivatePort, p.Type))
		}
	}
	return strings.Join(parts, ", ")
}
