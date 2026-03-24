package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

// ================================
// TYPES
// ================================

// CheckResult holds the result of checking one URL
type CheckResult struct {
	URL        string
	StatusCode int
	Status     string // "UP" or "DOWN"
	Duration   time.Duration
	Error      string
}

// Config holds CLI configuration
type Config struct {
	URLs    []string
	Timeout time.Duration
	File    string
}

// ================================
// CLI PARSING
// ================================

func parseArgs() (*Config, error) {
	args := os.Args[1:]

	if len(args) == 0 {
		return nil, fmt.Errorf("no URLs specified")
	}

	config := &Config{
		Timeout: 5 * time.Second, // default timeout
	}

	i := 0
	for i < len(args) {
		switch args[i] {
		case "--timeout":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--timeout requires a value (in seconds)")
			}
			i++
			var seconds int
			_, err := fmt.Sscanf(args[i], "%d", &seconds)
			if err != nil || seconds <= 0 {
				return nil, fmt.Errorf("--timeout must be a positive number, got: %s", args[i])
			}
			config.Timeout = time.Duration(seconds) * time.Second

		case "--file":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--file requires a file path")
			}
			i++
			config.File = args[i]

		default:
			// Treat as a URL
			url := args[i]
			if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
				url = "https://" + url // add https if missing
			}
			config.URLs = append(config.URLs, url)
		}
		i++
	}

	// If a file was given, read URLs from it
	if config.File != "" {
		fileURLs, err := readURLsFromFile(config.File)
		if err != nil {
			return nil, fmt.Errorf("cannot read URL file: %w", err)
		}
		config.URLs = append(config.URLs, fileURLs...)
	}

	if len(config.URLs) == 0 {
		return nil, fmt.Errorf("no URLs to check")
	}

	return config, nil
}

// readURLsFromFile reads URLs from a file (one per line)
func readURLsFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue // skip empty lines and comments
		}
		if !strings.HasPrefix(line, "http://") && !strings.HasPrefix(line, "https://") {
			line = "https://" + line
		}
		urls = append(urls, line)
	}

	return urls, scanner.Err()
}

func printUsage() {
	fmt.Println("Usage: health-checker [options] <url1> <url2> ...")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --timeout N    Timeout in seconds (default: 5)")
	fmt.Println("  --file PATH    Read URLs from a file (one per line)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run main.go https://google.com https://github.com")
	fmt.Println("  go run main.go --file urls.txt")
	fmt.Println("  go run main.go --timeout 10 google.com github.com")
}

// ================================
// HEALTH CHECK LOGIC
// ================================

// checkURL pings a single URL and returns the result
func checkURL(url string, timeout time.Duration) CheckResult {
	result := CheckResult{URL: url}

	client := &http.Client{Timeout: timeout}

	start := time.Now()
	resp, err := client.Get(url)
	result.Duration = time.Since(start)

	if err != nil {
		result.Status = "DOWN"
		result.Error = simplifyError(err)
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		result.Status = "UP"
	} else {
		result.Status = "DOWN"
	}

	return result
}

// simplifyError makes error messages shorter and more readable
func simplifyError(err error) string {
	errStr := err.Error()

	switch {
	case strings.Contains(errStr, "timeout") || strings.Contains(errStr, "deadline"):
		return "timeout"
	case strings.Contains(errStr, "no such host"):
		return "DNS lookup failed"
	case strings.Contains(errStr, "connection refused"):
		return "connection refused"
	case strings.Contains(errStr, "certificate"):
		return "TLS/SSL error"
	default:
		// Return a shortened version
		if len(errStr) > 50 {
			return errStr[:50] + "..."
		}
		return errStr
	}
}

// checkAllURLs runs health checks concurrently using goroutines
func checkAllURLs(urls []string, timeout time.Duration) []CheckResult {
	results := make(chan CheckResult, len(urls))

	// Launch a goroutine for each URL using WaitGroup.Go
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Go(func() {
			results <- checkURL(url, timeout)
		})
	}

	// Close channel when all goroutines finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var allResults []CheckResult
	for r := range results {
		allResults = append(allResults, r)
	}

	// Sort by URL for consistent output
	sort.Slice(allResults, func(i, j int) bool {
		return allResults[i].URL < allResults[j].URL
	})

	return allResults
}

// ================================
// DISPLAY
// ================================

func printResults(results []CheckResult) {
	fmt.Printf("\n%-8s %-5s %-10s %s\n", "Status", "Code", "Time", "URL")
	fmt.Printf("%-8s %-5s %-10s %s\n", "------", "----", "----", "---")

	for _, r := range results {
		icon := "✅ UP  "
		if r.Status == "DOWN" {
			icon = "❌ DOWN"
		}

		code := fmt.Sprintf("%d", r.StatusCode)
		if r.StatusCode == 0 {
			code = "---"
		}

		duration := fmt.Sprintf("%dms", r.Duration.Milliseconds())

		extra := ""
		if r.Error != "" {
			extra = fmt.Sprintf(" (%s)", r.Error)
		}

		fmt.Printf("%-8s %-5s %-10s %s%s\n", icon, code, duration, r.URL, extra)
	}
}

func printSummary(results []CheckResult) {
	upCount := 0
	downCount := 0
	var totalDuration time.Duration

	for _, r := range results {
		totalDuration += r.Duration
		if r.Status == "UP" {
			upCount++
		} else {
			downCount++
		}
	}

	avg := totalDuration / time.Duration(len(results))

	fmt.Println("\nSummary:")
	fmt.Printf("  Total  : %d\n", len(results))
	fmt.Printf("  Up     : %d\n", upCount)
	fmt.Printf("  Down   : %d\n", downCount)
	fmt.Printf("  Avg    : %dms\n", avg.Milliseconds())
}

// ================================
// MAIN
// ================================

func main() {
	config, err := parseArgs()
	if err != nil {
		fmt.Printf("Error: %s\n\n", err)
		printUsage()
		os.Exit(1)
	}

	fmt.Println("=== Website Health Checker ===")
	fmt.Printf("Timeout: %s\n", config.Timeout)
	fmt.Printf("Checking %d URLs concurrently...\n", len(config.URLs))

	start := time.Now()
	results := checkAllURLs(config.URLs, config.Timeout)
	elapsed := time.Since(start)

	printResults(results)
	printSummary(results)
	fmt.Printf("  Time   : %s\n\n", elapsed.Round(time.Millisecond))
}
