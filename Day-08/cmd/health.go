package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// HealthResult holds the result of one health check
type HealthResult struct {
	URL        string `json:"url"`
	Status     string `json:"status"`
	StatusCode int    `json:"status_code,omitempty"`
	TimeMs     int64  `json:"time_ms"`
	Error      string `json:"error,omitempty"`
}

var healthCmd = &cobra.Command{
	Use:   "health [urls...]",
	Short: "Check health of one or more URLs",
	Long: `Check if URLs are up or down.
Sends a GET request to each URL and reports the status.

Examples:
  devctl health https://google.com https://github.com
  devctl health https://google.com --timeout 5
  devctl health https://google.com -o json`,
	Args: cobra.MinimumNArgs(1),
	Run:  runHealthCheck,
}

var healthTimeout int

func init() {
	rootCmd.AddCommand(healthCmd)
	healthCmd.Flags().IntVar(&healthTimeout, "timeout", 0, "Request timeout in seconds (default from config or 10)")
}

func runHealthCheck(cmd *cobra.Command, args []string) {
	// Get timeout - flag > config > default
	timeout := healthTimeout
	if timeout == 0 {
		timeout = viper.GetInt("health.timeout")
	}
	if timeout == 0 {
		timeout = 10
	}

	concurrent := viper.GetInt("health.concurrent")
	if concurrent == 0 {
		concurrent = 5
	}

	if Verbose {
		fmt.Printf("[verbose] Timeout: %ds | Concurrency: %d\n", timeout, concurrent)
	}

	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Check all URLs concurrently
	var results []HealthResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Use a semaphore to limit concurrency
	sem := make(chan struct{}, concurrent)

	for _, url := range args {
		// Make sure URL has a scheme
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "https://" + url
		}

		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			sem <- struct{}{}        // acquire
			defer func() { <-sem }() // release

			result := checkURL(client, url)

			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(url)
	}

	wg.Wait()

	// Display results
	if OutputFormat == "json" {
		printHealthJSON(results)
	} else {
		printHealthTable(results)
	}
}

func checkURL(client *http.Client, url string) HealthResult {
	result := HealthResult{URL: url}

	start := time.Now()
	resp, err := client.Get(url)
	elapsed := time.Since(start)
	result.TimeMs = elapsed.Milliseconds()

	if err != nil {
		result.Status = "DOWN"
		// Simplify error message
		errMsg := err.Error()
		if strings.Contains(errMsg, "no such host") {
			errMsg = "DNS lookup failed"
		} else if strings.Contains(errMsg, "connection refused") {
			errMsg = "connection refused"
		} else if strings.Contains(errMsg, "timeout") {
			errMsg = "timeout"
		}
		result.Error = errMsg
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

func printHealthTable(results []HealthResult) {
	fmt.Println("\n=== Health Check ===\n")

	// Find max URL length
	maxURL := 3
	for _, r := range results {
		if len(r.URL) > maxURL {
			maxURL = len(r.URL)
		}
	}
	if maxURL > 50 {
		maxURL = 50
	}

	// Header
	fmt.Printf("%-*s  %-6s  %4s  %s\n", maxURL, "URL", "Status", "Code", "Time")
	fmt.Printf("%-*s  %-6s  %4s  %s\n", maxURL, strings.Repeat("-", maxURL), "------", "----", "----")

	up := 0
	down := 0

	for _, r := range results {
		url := r.URL
		if len(url) > 50 {
			url = url[:47] + "..."
		}

		icon := "✅"
		codeStr := fmt.Sprintf("%d", r.StatusCode)
		extra := ""

		if r.Status == "DOWN" {
			icon = "❌"
			down++
			if r.StatusCode == 0 {
				codeStr = "-"
			}
			if r.Error != "" {
				extra = fmt.Sprintf(" (%s)", r.Error)
			}
		} else {
			up++
		}

		fmt.Printf("%-*s  %s %-4s  %4s  %dms%s\n", maxURL, url, icon, r.Status, codeStr, r.TimeMs, extra)
	}

	fmt.Printf("\nChecked: %d | Up: %d | Down: %d\n", len(results), up, down)
}

func printHealthJSON(results []HealthResult) {
	data, _ := json.MarshalIndent(results, "", "  ")
	fmt.Println(string(data))
}
