package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

// ================================
// TYPES
// ================================

// GithubRepo holds the fields we care about from GitHub API
type GithubRepo struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	HTMLURL     string `json:"html_url"`
	Language    string `json:"language"`
	Stars       int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
	OpenIssues  int    `json:"open_issues_count"`
	PushedAt    string `json:"pushed_at"`
}

// Config holds CLI options
type Config struct {
	Command  string // "repos" or "serve"
	Username string
	Sort     string // "name", "stars", "updated"
	Language string // filter by language
	Port     string // server port
}

// ================================
// CLI PARSING
// ================================

func parseArgs() *Config {
	config := &Config{
		Sort: "name",
		Port: "8080",
	}

	args := os.Args[1:]
	if len(args) == 0 {
		printUsage()
		os.Exit(0)
	}

	config.Command = args[0]

	switch config.Command {
	case "repos":
		if len(args) < 2 {
			fmt.Println("Error: username is required")
			fmt.Println("Usage: go run main.go repos <username> [--sort name|stars|updated] [--language go]")
			os.Exit(1)
		}
		config.Username = args[1]

		// Parse optional flags
		for i := 2; i < len(args); i++ {
			switch args[i] {
			case "--sort":
				if i+1 < len(args) {
					i++
					config.Sort = strings.ToLower(args[i])
				}
			case "--language":
				if i+1 < len(args) {
					i++
					config.Language = strings.ToLower(args[i])
				}
			}
		}

	case "serve":
		for i := 1; i < len(args); i++ {
			if args[i] == "--port" && i+1 < len(args) {
				i++
				config.Port = args[i]
			}
		}

	case "--help", "-h", "help":
		printUsage()
		os.Exit(0)

	default:
		fmt.Printf("Unknown command: %s\n", config.Command)
		printUsage()
		os.Exit(1)
	}

	return config
}

func printUsage() {
	fmt.Println("GitHub Tool - List repos and serve API")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run main.go repos <username> [flags]    List GitHub repos")
	fmt.Println("  go run main.go serve [flags]               Start API server")
	fmt.Println()
	fmt.Println("Repos flags:")
	fmt.Println("  --sort name|stars|updated    Sort repos (default: name)")
	fmt.Println("  --language go|python|...     Filter by language")
	fmt.Println()
	fmt.Println("Serve flags:")
	fmt.Println("  --port 8080                  Server port (default: 8080)")
	fmt.Println()
	fmt.Println("Environment:")
	fmt.Println("  GITHUB_TOKEN    Optional. Set for private repos or higher API limits")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run main.go repos golang")
	fmt.Println("  go run main.go repos torvalds --sort stars")
	fmt.Println("  go run main.go repos google --language go")
	fmt.Println("  go run main.go serve")
	fmt.Println("  go run main.go serve --port 9090")
}

// ================================
// GITHUB API CLIENT
// ================================

// fetchRepos calls GitHub API and returns list of repos for a user
func fetchRepos(username string) ([]GithubRepo, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?per_page=100&type=owner", username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "golang-devops-tool")

	// If GITHUB_TOKEN is set, use it for auth (higher rate limit + private repos)
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("user '%s' not found on GitHub", username)
	}
	if resp.StatusCode == 403 {
		return nil, fmt.Errorf("rate limited by GitHub. Set GITHUB_TOKEN env var to increase limit")
	}
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API returned %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response: %w", err)
	}

	var repos []GithubRepo
	err = json.Unmarshal(body, &repos)
	if err != nil {
		return nil, fmt.Errorf("cannot parse JSON: %w", err)
	}

	return repos, nil
}

// ================================
// SORTING AND FILTERING
// ================================

func filterRepos(repos []GithubRepo, language string) []GithubRepo {
	if language == "" {
		return repos
	}

	var filtered []GithubRepo
	for _, repo := range repos {
		if strings.EqualFold(repo.Language, language) {
			filtered = append(filtered, repo)
		}
	}
	return filtered
}

func sortRepos(repos []GithubRepo, sortBy string) {
	switch sortBy {
	case "stars":
		sort.Slice(repos, func(i, j int) bool {
			return repos[i].Stars > repos[j].Stars // highest first
		})
	case "updated":
		sort.Slice(repos, func(i, j int) bool {
			return repos[i].PushedAt > repos[j].PushedAt // most recent first
		})
	default: // "name"
		sort.Slice(repos, func(i, j int) bool {
			return strings.ToLower(repos[i].Name) < strings.ToLower(repos[j].Name)
		})
	}
}

// ================================
// DISPLAY (TABLE FORMAT)
// ================================

func printReposTable(username string, repos []GithubRepo) {
	fmt.Printf("\n=== GitHub Repos for %s ===\n\n", username)

	if len(repos) == 0 {
		fmt.Println("No repos found.")
		return
	}

	// Find max widths for nice formatting
	maxName := 4 // "Name" header
	maxLang := 8 // "Language" header
	for _, repo := range repos {
		if len(repo.Name) > maxName {
			maxName = len(repo.Name)
		}
		lang := repo.Language
		if lang == "" {
			lang = "-"
		}
		if len(lang) > maxLang {
			maxLang = len(lang)
		}
	}
	if maxName > 35 {
		maxName = 35
	}

	// Print header
	fmt.Printf("%-*s  %-*s  %5s  %5s  %s\n", maxName, "Name", maxLang, "Language", "Stars", "Forks", "Last Push")
	fmt.Printf("%-*s  %-*s  %5s  %5s  %s\n", maxName, strings.Repeat("-", maxName), maxLang, strings.Repeat("-", maxLang), "-----", "-----", "---------")

	// Print repos
	for _, repo := range repos {
		name := repo.Name
		if len(name) > 35 {
			name = name[:32] + "..."
		}

		lang := repo.Language
		if lang == "" {
			lang = "-"
		}

		lastPush := formatDate(repo.PushedAt)

		fmt.Printf("%-*s  %-*s  %5d  %5d  %s\n", maxName, name, maxLang, lang, repo.Stars, repo.Forks, lastPush)
	}

	fmt.Printf("\nTotal: %d repos\n", len(repos))
}

// formatDate converts "2024-01-10T15:30:00Z" to "2024-01-10"
func formatDate(isoDate string) string {
	t, err := time.Parse(time.RFC3339, isoDate)
	if err != nil {
		return isoDate
	}
	return t.Format("2006-01-02")
}

// ================================
// REPOS COMMAND
// ================================

func runReposCommand(config *Config) {
	fmt.Printf("Fetching repos for '%s'...\n", config.Username)

	repos, err := fetchRepos(config.Username)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	// Filter by language
	repos = filterRepos(repos, config.Language)

	// Sort
	sortRepos(repos, config.Sort)

	// Display
	printReposTable(config.Username, repos)

	if config.Language != "" {
		fmt.Printf("(filtered by language: %s)\n", config.Language)
	}
	fmt.Printf("(sorted by: %s)\n", config.Sort)
}

// ================================
// SERVE COMMAND (HTTP SERVER)
// ================================

func runServeCommand(config *Config) {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// List repos for a user
	mux.HandleFunc("GET /api/repos/{username}", func(w http.ResponseWriter, r *http.Request) {
		username := r.PathValue("username")

		// Read query params
		sortBy := r.URL.Query().Get("sort")
		if sortBy == "" {
			sortBy = "name"
		}
		language := r.URL.Query().Get("language")

		// Fetch from GitHub
		repos, err := fetchRepos(username)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		// Filter and sort
		repos = filterRepos(repos, language)
		sortRepos(repos, sortBy)

		// Log the request
		log.Printf("GET /api/repos/%s - %d repos returned", username, len(repos))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"username":   username,
			"repo_count": len(repos),
			"sort":       sortBy,
			"repos":      repos,
		})
	})

	// Root endpoint - show help
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"service": "GitHub Tool API",
			"endpoints": map[string]string{
				"GET /health":                          "Health check",
				"GET /api/repos/{username}":            "List repos for a GitHub user",
				"GET /api/repos/{username}?sort=stars": "Sort by: name, stars, updated",
				"GET /api/repos/{username}?language=go": "Filter by programming language",
			},
		})
	})

	// Print startup info
	addr := ":" + config.Port
	fmt.Println("=== GitHub Tool - API Server ===")
	fmt.Println()
	fmt.Println("Endpoints:")
	fmt.Println("  GET /health                    - Health check")
	fmt.Println("  GET /api/repos/{username}      - List GitHub repos for a user")
	fmt.Println()
	fmt.Println("Query parameters:")
	fmt.Println("  ?sort=name|stars|updated        - Sort repos")
	fmt.Println("  ?language=go|python|...         - Filter by language")
	fmt.Println()
	fmt.Println("Try:")
	fmt.Printf("  curl http://localhost:%s/health\n", config.Port)
	fmt.Printf("  curl http://localhost:%s/api/repos/golang\n", config.Port)
	fmt.Printf("  curl http://localhost:%s/api/repos/torvalds?sort=stars\n", config.Port)
	fmt.Println()

	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		fmt.Println("Auth: GITHUB_TOKEN is set (private repos + higher rate limit)")
	} else {
		fmt.Println("Auth: No GITHUB_TOKEN set (public repos only, 60 requests/hour)")
	}

	fmt.Printf("\nServer starting on %s ...\n", addr)

	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}

// ================================
// MAIN
// ================================

func main() {
	config := parseArgs()

	switch config.Command {
	case "repos":
		runReposCommand(config)
	case "serve":
		runServeCommand(config)
	}
}
