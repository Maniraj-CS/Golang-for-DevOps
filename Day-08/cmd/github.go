package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// ================================
// TYPES
// ================================

type GithubRepo struct {
	Name       string `json:"name"`
	FullName   string `json:"full_name"`
	Language   string `json:"language"`
	Stars      int    `json:"stargazers_count"`
	Forks      int    `json:"forks_count"`
	OpenIssues int    `json:"open_issues_count"`
	PushedAt   string `json:"pushed_at"`
	HTMLURL    string `json:"html_url"`
}

// ================================
// COBRA COMMANDS
// ================================

var githubCmd = &cobra.Command{
	Use:   "github",
	Short: "GitHub operations",
	Long:  "Interact with GitHub API. Use subcommands like 'repos' to list repositories.",
}

var ghSortBy string
var ghLanguage string

var githubReposCmd = &cobra.Command{
	Use:   "repos <username>",
	Short: "List GitHub repos for a user",
	Long: `List public repositories for a GitHub user.

Examples:
  devctl github repos golang
  devctl github repos torvalds --sort stars
  devctl github repos google --language go
  devctl github repos golang --sort updated -o json`,
	Args: cobra.ExactArgs(1),
	Run:  runGithubRepos,
}

func init() {
	rootCmd.AddCommand(githubCmd)
	githubCmd.AddCommand(githubReposCmd)

	githubReposCmd.Flags().StringVar(&ghSortBy, "sort", "name", "Sort by: name, stars, updated")
	githubReposCmd.Flags().StringVar(&ghLanguage, "language", "", "Filter by programming language")
}

func runGithubRepos(cmd *cobra.Command, args []string) {
	username := args[0]

	if Verbose {
		fmt.Printf("[verbose] Fetching repos for %s...\n", username)
	}

	repos, err := fetchGithubRepos(username)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	// Filter by language
	if ghLanguage != "" {
		var filtered []GithubRepo
		for _, r := range repos {
			if strings.EqualFold(r.Language, ghLanguage) {
				filtered = append(filtered, r)
			}
		}
		repos = filtered
	}

	// Sort
	sortGithubRepos(repos, ghSortBy)

	// Display
	if OutputFormat == "json" {
		printGithubJSON(username, repos)
	} else {
		printGithubTable(username, repos)
	}
}

// ================================
// GITHUB API
// ================================

func fetchGithubRepos(username string) ([]GithubRepo, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?per_page=100&type=owner", username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "devctl-tool")

	// Use token if available
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
		if Verbose {
			fmt.Println("[verbose] Using GITHUB_TOKEN for auth")
		}
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("user '%s' not found", username)
	}
	if resp.StatusCode == 403 {
		return nil, fmt.Errorf("rate limited. Set GITHUB_TOKEN env var")
	}
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API returned %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response: %w", err)
	}

	var repos []GithubRepo
	if err := json.Unmarshal(body, &repos); err != nil {
		return nil, fmt.Errorf("cannot parse JSON: %w", err)
	}

	return repos, nil
}

func sortGithubRepos(repos []GithubRepo, by string) {
	switch by {
	case "stars":
		sort.Slice(repos, func(i, j int) bool { return repos[i].Stars > repos[j].Stars })
	case "updated":
		sort.Slice(repos, func(i, j int) bool { return repos[i].PushedAt > repos[j].PushedAt })
	default:
		sort.Slice(repos, func(i, j int) bool {
			return strings.ToLower(repos[i].Name) < strings.ToLower(repos[j].Name)
		})
	}
}

// ================================
// DISPLAY
// ================================

func printGithubTable(username string, repos []GithubRepo) {
	fmt.Printf("\n=== GitHub Repos for %s ===\n\n", username)

	if len(repos) == 0 {
		fmt.Println("No repos found.")
		return
	}

	maxName := 4
	maxLang := 8
	for _, r := range repos {
		if len(r.Name) > maxName {
			maxName = len(r.Name)
		}
		lang := r.Language
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

	fmt.Printf("%-*s  %-*s  %5s  %5s  %s\n", maxName, "Name", maxLang, "Language", "Stars", "Forks", "Last Push")
	fmt.Printf("%-*s  %-*s  %5s  %5s  %s\n", maxName, strings.Repeat("-", maxName), maxLang, strings.Repeat("-", maxLang), "-----", "-----", "---------")

	for _, r := range repos {
		name := r.Name
		if len(name) > 35 {
			name = name[:32] + "..."
		}
		lang := r.Language
		if lang == "" {
			lang = "-"
		}
		lastPush := formatGHDate(r.PushedAt)
		fmt.Printf("%-*s  %-*s  %5d  %5d  %s\n", maxName, name, maxLang, lang, r.Stars, r.Forks, lastPush)
	}

	fmt.Printf("\nTotal: %d repos", len(repos))
	if ghLanguage != "" {
		fmt.Printf(" (language: %s)", ghLanguage)
	}
	fmt.Printf(" (sorted by: %s)\n", ghSortBy)
}

func printGithubJSON(username string, repos []GithubRepo) {
	output := map[string]any{
		"username":   username,
		"repo_count": len(repos),
		"sort":       ghSortBy,
		"repos":      repos,
	}
	data, _ := json.MarshalIndent(output, "", "  ")
	fmt.Println(string(data))
}

func formatGHDate(isoDate string) string {
	t, err := time.Parse(time.RFC3339, isoDate)
	if err != nil {
		return isoDate
	}
	return t.Format("2006-01-02")
}
