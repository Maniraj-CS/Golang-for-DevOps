package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// ================================
// STRUCTS
// ================================

// Server represents one server in our inventory
type Server struct {
	Name   string `json:"name" yaml:"name"`
	IP     string `json:"ip" yaml:"ip"`
	Port   int    `json:"port" yaml:"port"`
	Env    string `json:"env" yaml:"env"`
	Status string `json:"status" yaml:"status"`
}

// ================================
// INTERFACE
// ================================

// ConfigReader - any struct that can read a file and return servers is a ConfigReader
type ConfigReader interface {
	Read(path string) ([]Server, error)
}

// ================================
// JSON READER
// ================================

type JSONReader struct{}

func (j JSONReader) Read(path string) ([]Server, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}

	var servers []Server
	err = json.Unmarshal(data, &servers)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	return servers, nil
}

// ================================
// YAML READER
// ================================

type YAMLReader struct{}

func (y YAMLReader) Read(path string) ([]Server, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}

	var servers []Server
	err = yaml.Unmarshal(data, &servers)
	if err != nil {
		return nil, fmt.Errorf("invalid YAML: %w", err)
	}

	return servers, nil
}

// ================================
// HELPER FUNCTIONS
// ================================

// getReader returns the right config reader based on file extension
// This is where the interface magic happens!
func getReader(path string) (ConfigReader, error) {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".json":
		return JSONReader{}, nil
	case ".yaml", ".yml":
		return YAMLReader{}, nil
	default:
		return nil, fmt.Errorf("unsupported file format: %s (use .json or .yaml)", ext)
	}
}

// printTable prints a nicely formatted server table
func printTable(servers []Server) {
	// Header
	fmt.Printf("\n%-14s %-15s %-6s %-13s %-10s\n", "Name", "IP", "Port", "Env", "Status")
	fmt.Printf("%-14s %-15s %-6s %-13s %-10s\n", "----", "--", "----", "---", "------")

	// Rows
	for _, s := range servers {
		fmt.Printf("%-14s %-15s %-6d %-13s %-10s\n", s.Name, s.IP, s.Port, s.Env, s.Status)
	}
}

// printSummary shows breakdown by environment and status
func printSummary(servers []Server) {
	envCount := map[string]int{}
	statusCount := map[string]int{}

	for _, s := range servers {
		envCount[s.Env]++
		statusCount[s.Status]++
	}

	fmt.Printf("\nTotal: %d servers\n", len(servers))

	// By Environment - sorted for consistent output
	fmt.Println("\nBy Environment:")
	envKeys := sortedKeys(envCount)
	for _, env := range envKeys {
		fmt.Printf("  %-13s: %d\n", env, envCount[env])
	}

	// By Status - sorted for consistent output
	fmt.Println("\nBy Status:")
	statusKeys := sortedKeys(statusCount)
	for _, status := range statusKeys {
		fmt.Printf("  %-13s: %d\n", status, statusCount[status])
	}
}

// sortedKeys returns map keys in sorted order
func sortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// ================================
// MAIN
// ================================

func main() {
	// Check if user gave a file path
	if len(os.Args) < 2 {
		fmt.Println("Usage: config-reader <config-file>")
		fmt.Println("Example:")
		fmt.Println("  go run main.go servers.json")
		fmt.Println("  go run main.go servers.yaml")
		os.Exit(1)
	}

	configPath := os.Args[1]

	// Get the right reader based on file extension (interface in action!)
	reader, err := getReader(configPath)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	// Read the config file - same code works for JSON and YAML!
	servers, err := reader.Read(configPath)
	if err != nil {
		fmt.Printf("Error reading config: %s\n", err)
		os.Exit(1)
	}

	// Display everything
	fmt.Println("=== Server Inventory ===")
	printTable(servers)
	printSummary(servers)
}
