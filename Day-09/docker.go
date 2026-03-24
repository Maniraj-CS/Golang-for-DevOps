package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

// ================================
// TYPES
// ================================

type ContainerInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Image   string `json:"image"`
	Status  string `json:"status"`
	State   string `json:"state"`
	Ports   string `json:"ports"`
	Created string `json:"created"`
	// Stats (optional)
	CPUPercent float64 `json:"cpu_percent,omitempty"`
	MemUsageMB float64 `json:"mem_usage_mb,omitempty"`
	MemLimitMB float64 `json:"mem_limit_mb,omitempty"`
}

type ImageInfo struct {
	ID     string  `json:"id"`
	Tags   string  `json:"tags"`
	SizeMB float64 `json:"size_mb"`
}

// ================================
// DOCKER OPERATIONS
// ================================

func runDocker(config *Config) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Connect to Docker
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println("Error: Cannot create Docker client:", err)
		fmt.Println("Is Docker installed?")
		return
	}
	defer cli.Close()

	// Test connection
	_, err = cli.Ping(ctx)
	if err != nil {
		fmt.Println("Error: Cannot connect to Docker daemon:", err)
		fmt.Println("Is Docker running? Try: docker info")
		return
	}

	if config.ShowImages {
		listImages(ctx, cli, config)
	} else {
		listContainers(ctx, cli, config)
	}
}

func listContainers(ctx context.Context, cli *client.Client, config *Config) {
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		fmt.Println("Error listing containers:", err)
		return
	}

	if len(containers) == 0 {
		fmt.Println("\n=== Docker Containers ===\n")
		fmt.Println("No containers found.")
		fmt.Println("Try: docker run -d --name test-nginx nginx:latest")
		return
	}

	// Convert to our type
	var infos []ContainerInfo
	running := 0
	stopped := 0

	for _, c := range containers {
		info := ContainerInfo{
			ID:     c.ID[:12],
			Image:  c.Image,
			Status: c.Status,
			State:  c.State,
			Ports:  formatContainerPorts(c.Ports),
		}

		if len(c.Names) > 0 {
			info.Name = strings.TrimPrefix(c.Names[0], "/")
		}

		// Created time
		info.Created = time.Unix(c.Created, 0).Format("2006-01-02 15:04")

		// Filter
		if config.Filter != "" {
			if config.Filter == "running" && c.State != "running" {
				continue
			}
			if config.Filter == "stopped" && c.State == "running" {
				continue
			}
		}

		// Get stats if requested
		if config.ShowStats && c.State == "running" {
			getCPUMemStats(ctx, cli, c.ID, &info)
		}

		if c.State == "running" {
			running++
		} else {
			stopped++
		}

		infos = append(infos, info)
	}

	// Display
	if config.Output == "json" {
		data, _ := json.MarshalIndent(infos, "", "  ")
		fmt.Println(string(data))
	} else {
		printContainerTable(infos, config.ShowStats, running, stopped)
	}
}

func getCPUMemStats(ctx context.Context, cli *client.Client, containerID string, info *ContainerInfo) {
	statsCtx, statsCancel := context.WithTimeout(ctx, 5*time.Second)
	defer statsCancel()

	stats, err := cli.ContainerStats(statsCtx, containerID, false)
	if err != nil {
		return
	}
	defer stats.Body.Close()

	var statsJSON container.StatsResponse
	if err := json.NewDecoder(stats.Body).Decode(&statsJSON); err != nil {
		return
	}

	// CPU percentage
	cpuDelta := float64(statsJSON.CPUStats.CPUUsage.TotalUsage - statsJSON.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(statsJSON.CPUStats.SystemUsage - statsJSON.PreCPUStats.SystemUsage)
	if systemDelta > 0 && cpuDelta > 0 {
		cpuCount := float64(statsJSON.CPUStats.OnlineCPUs)
		if cpuCount == 0 {
			cpuCount = float64(len(statsJSON.CPUStats.CPUUsage.PercpuUsage))
		}
		info.CPUPercent = (cpuDelta / systemDelta) * cpuCount * 100.0
	}

	// Memory
	info.MemUsageMB = float64(statsJSON.MemoryStats.Usage) / 1024 / 1024
	info.MemLimitMB = float64(statsJSON.MemoryStats.Limit) / 1024 / 1024
}

func listImages(ctx context.Context, cli *client.Client, config *Config) {
	images, err := cli.ImageList(ctx, image.ListOptions{})
	if err != nil {
		fmt.Println("Error listing images:", err)
		return
	}

	if len(images) == 0 {
		fmt.Println("\n=== Docker Images ===\n")
		fmt.Println("No images found.")
		return
	}

	var infos []ImageInfo
	for _, img := range images {
		tags := "<none>"
		if len(img.RepoTags) > 0 {
			tags = strings.Join(img.RepoTags, ", ")
		}

		infos = append(infos, ImageInfo{
			ID:     img.ID[7:19], // skip "sha256:" prefix, show 12 chars
			Tags:   tags,
			SizeMB: float64(img.Size) / 1024 / 1024,
		})
	}

	if config.Output == "json" {
		data, _ := json.MarshalIndent(infos, "", "  ")
		fmt.Println(string(data))
	} else {
		printImageTable(infos)
	}
}

// ================================
// DISPLAY
// ================================

func printContainerTable(infos []ContainerInfo, showStats bool, running, stopped int) {
	fmt.Println("\n=== Docker Containers ===\n")

	if len(infos) == 0 {
		fmt.Println("No containers match the filter.")
		return
	}

	// Find max widths
	maxName := 4
	maxImage := 5
	for _, info := range infos {
		if len(info.Name) > maxName {
			maxName = len(info.Name)
		}
		if len(info.Image) > maxImage {
			maxImage = len(info.Image)
		}
	}
	if maxName > 30 {
		maxName = 30
	}
	if maxImage > 30 {
		maxImage = 30
	}

	// Header
	header := fmt.Sprintf("%-12s  %-*s  %-*s  %-20s  %s",
		"ID", maxName, "Name", maxImage, "Image", "Status", "Ports")
	if showStats {
		header += fmt.Sprintf("  %8s  %10s", "CPU%", "Memory")
	}
	fmt.Println(header)
	fmt.Println(strings.Repeat("-", len(header)))

	for _, info := range infos {
		name := info.Name
		if len(name) > 30 {
			name = name[:27] + "..."
		}
		img := info.Image
		if len(img) > 30 {
			img = img[:27] + "..."
		}

		icon := "🟢"
		if info.State != "running" {
			icon = "🔴"
		}

		line := fmt.Sprintf("%s %-12s  %-*s  %-*s  %-20s  %s",
			icon, info.ID, maxName, name, maxImage, img, info.Status, info.Ports)

		if showStats && info.State == "running" {
			memStr := fmt.Sprintf("%.0f/%.0fMB", info.MemUsageMB, info.MemLimitMB)
			line += fmt.Sprintf("  %7.1f%%  %10s", info.CPUPercent, memStr)
		}

		fmt.Println(line)
	}

	fmt.Printf("\nRunning: %d | Stopped: %d | Total: %d\n", running, stopped, len(infos))
}

func printImageTable(infos []ImageInfo) {
	fmt.Println("\n=== Docker Images ===\n")

	maxTags := 4
	for _, info := range infos {
		if len(info.Tags) > maxTags {
			maxTags = len(info.Tags)
		}
	}
	if maxTags > 50 {
		maxTags = 50
	}

	fmt.Printf("%-12s  %-*s  %8s\n", "ID", maxTags, "Tags", "Size")
	fmt.Printf("%-12s  %-*s  %8s\n", "------------", maxTags, strings.Repeat("-", maxTags), "--------")

	var totalSize float64
	for _, info := range infos {
		tags := info.Tags
		if len(tags) > 50 {
			tags = tags[:47] + "..."
		}
		fmt.Printf("%-12s  %-*s  %6.0f MB\n", info.ID, maxTags, tags, info.SizeMB)
		totalSize += info.SizeMB
	}

	fmt.Printf("\nTotal: %d images (%.0f MB)\n", len(infos), totalSize)
}

func formatContainerPorts(ports []types.Port) string {
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
	result := strings.Join(parts, ", ")
	if len(result) > 30 {
		result = result[:27] + "..."
	}
	return result
}
