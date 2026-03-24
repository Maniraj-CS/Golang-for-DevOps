package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// ================================
// STRUCTS FOR NESTED JSON
// ================================

// This is what GitHub API returns for a repo.
// We only define the fields we care about - Go ignores the rest!

type GithubRepo struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	HTMLURL     string `json:"html_url"`
	Language    string `json:"language"`
	Stars       int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
	Owner       Owner  `json:"owner"`
	PushedAt    string `json:"pushed_at"`
}

type Owner struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
}

func main() {
	// ================================
	// 1. PARSING NESTED JSON WITH STRUCTS
	// ================================

	fmt.Println("=== Parsing Nested JSON with Structs ===")

	// This is a sample GitHub API response (simplified)
	sampleJSON := `[
		{
			"name": "kubernetes-tutorial",
			"full_name": "abhishekkumar/kubernetes-tutorial",
			"description": "Learn K8s step by step",
			"private": false,
			"html_url": "https://github.com/abhishekkumar/kubernetes-tutorial",
			"language": "Go",
			"stargazers_count": 145,
			"forks_count": 23,
			"owner": {
				"login": "abhishekkumar",
				"avatar_url": "https://avatars.githubusercontent.com/u/12345"
			},
			"pushed_at": "2024-01-10T15:30:00Z"
		},
		{
			"name": "devops-scripts",
			"full_name": "abhishekkumar/devops-scripts",
			"description": "Useful scripts for daily DevOps work",
			"private": false,
			"html_url": "https://github.com/abhishekkumar/devops-scripts",
			"language": "Shell",
			"stargazers_count": 89,
			"forks_count": 12,
			"owner": {
				"login": "abhishekkumar",
				"avatar_url": "https://avatars.githubusercontent.com/u/12345"
			},
			"pushed_at": "2024-01-08T10:00:00Z"
		}
	]`

	var repos []GithubRepo
	err := json.Unmarshal([]byte(sampleJSON), &repos)
	if err != nil {
		log.Fatal("Cannot parse JSON:", err)
	}

	for _, repo := range repos {
		fmt.Printf("Repo: %s\n", repo.FullName)
		fmt.Printf("  Owner: %s\n", repo.Owner.Login)
		fmt.Printf("  Language: %s | Stars: %d | Forks: %d\n", repo.Language, repo.Stars, repo.Forks)
		fmt.Printf("  URL: %s\n", repo.HTMLURL)
		fmt.Printf("  Last Push: %s\n\n", repo.PushedAt)
	}

	// ================================
	// 2. PARSING WITH map[string]any (unknown structure)
	// ================================

	fmt.Println("=== Parsing with map[string]any ===")

	// Sometimes you don't know the JSON structure ahead of time.
	// Use map[string]any to handle it dynamically.

	unknownJSON := `{
		"server": "web-01",
		"status": "running",
		"metrics": {
			"cpu": 45.2,
			"memory": 78.1,
			"connections": 1234
		},
		"tags": ["production", "us-east-1", "web"]
	}`

	var data map[string]any
	err = json.Unmarshal([]byte(unknownJSON), &data)
	if err != nil {
		log.Fatal("Cannot parse:", err)
	}

	// Access top-level fields
	fmt.Println("Server:", data["server"])
	fmt.Println("Status:", data["status"])

	// Access nested fields - need type assertion!
	metrics := data["metrics"].(map[string]any)
	fmt.Printf("CPU: %.1f%%\n", metrics["cpu"].(float64))
	fmt.Printf("Memory: %.1f%%\n", metrics["memory"].(float64))
	fmt.Printf("Connections: %.0f\n", metrics["connections"].(float64))

	// Access array fields
	tags := data["tags"].([]any)
	fmt.Print("Tags: ")
	for i, tag := range tags {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(tag.(string))
	}
	fmt.Println()

	// ================================
	// 3. SAFE TYPE ASSERTION
	// ================================

	fmt.Println("\n=== Safe Type Assertion ===")

	// The above code will PANIC if a field is missing or wrong type.
	// Here's the safe way:

	if metrics, ok := data["metrics"].(map[string]any); ok {
		if cpu, ok := metrics["cpu"].(float64); ok {
			fmt.Printf("CPU (safe): %.1f%%\n", cpu)
		} else {
			fmt.Println("CPU field is not a number")
		}
	} else {
		fmt.Println("No metrics found")
	}

	// Check for a field that doesn't exist
	if _, ok := data["nonexistent"].(string); ok {
		fmt.Println("This won't print")
	} else {
		fmt.Println("Field 'nonexistent' not found (expected)")
	}

	// ================================
	// 4. CONVERTING STRUCT TO JSON (Marshal)
	// ================================

	fmt.Println("\n=== Struct to JSON ===")

	newRepo := GithubRepo{
		Name:     "my-new-tool",
		FullName: "abhishekkumar/my-new-tool",
		Language: "Go",
		Stars:    0,
		Owner: Owner{
			Login: "abhishekkumar",
		},
	}

	// Pretty print with indent
	jsonBytes, err := json.MarshalIndent(newRepo, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonBytes))

	// ================================
	// KEY TAKEAWAYS
	// ================================

	fmt.Println("\n=== Key Takeaways ===")
	fmt.Println("1. Use structs when you know the JSON structure (recommended)")
	fmt.Println("2. Use map[string]any when JSON structure is unknown")
	fmt.Println("3. Always use safe type assertions (value, ok := ...) in production")
	fmt.Println("4. Go ignores JSON fields that don't match your struct")
	fmt.Println("5. json.MarshalIndent() for pretty printing")
}
