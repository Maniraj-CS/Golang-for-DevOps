package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	// ================================
	// 1. SIMPLE GET REQUEST
	// ================================

	fmt.Println("=== Simple GET Request ===")

	resp, err := http.Get("https://httpbin.org/get")
	if err != nil {
		log.Fatal("Request failed:", err)
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Content-Type:", resp.Header.Get("Content-Type"))

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Cannot read body:", err)
	}
	// Show just first 200 chars to keep output clean
	preview := string(body)
	if len(preview) > 200 {
		preview = preview[:200] + "..."
	}
	fmt.Println("Body (preview):", preview)

	// ================================
	// 2. CHECKING STATUS CODES
	// ================================

	fmt.Println("\n=== Checking Status Codes ===")

	urls := []string{
		"https://httpbin.org/status/200",
		"https://httpbin.org/status/301",
		"https://httpbin.org/status/404",
		"https://httpbin.org/status/500",
	}

	// Use a client that does NOT follow redirects, so we see the actual status
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // don't follow redirects
		},
	}

	for _, url := range urls {
		resp, err := client.Get(url)
		if err != nil {
			fmt.Printf("  ❌ %s - Error: %s\n", url, err)
			continue
		}
		resp.Body.Close()

		status := "❓"
		switch {
		case resp.StatusCode >= 200 && resp.StatusCode < 300:
			status = "✅"
		case resp.StatusCode >= 300 && resp.StatusCode < 400:
			status = "↗️ "
		case resp.StatusCode >= 400 && resp.StatusCode < 500:
			status = "❌"
		case resp.StatusCode >= 500:
			status = "💥"
		}

		fmt.Printf("  %s %d - %s\n", status, resp.StatusCode, url)
	}

	// ================================
	// 3. TIMEOUT
	// ================================

	fmt.Println("\n=== Timeout Demo ===")

	fastClient := &http.Client{
		Timeout: 2 * time.Second, // very short timeout
	}

	fmt.Println("Requesting a slow endpoint (10s delay) with 2s timeout...")
	start := time.Now()
	_, err = fastClient.Get("https://httpbin.org/delay/10")
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("  ❌ Failed after %s: %s\n", elapsed.Round(time.Millisecond), "timeout exceeded")
	} else {
		fmt.Printf("  ✅ Succeeded in %s\n", elapsed.Round(time.Millisecond))
	}

	// ================================
	// 4. CUSTOM REQUEST WITH HEADERS
	// ================================

	fmt.Println("\n=== Custom Request with Headers ===")

	req, err := http.NewRequest("GET", "https://httpbin.org/headers", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "GoDevOps-HealthChecker/1.0")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Custom-Header", "hello-from-go")

	normalClient := &http.Client{Timeout: 10 * time.Second}
	resp, err = normalClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Println("Response (shows our headers reflected back):")
	fmt.Println(string(body))

	// ================================
	// 5. MEASURING RESPONSE TIME
	// ================================

	fmt.Println("=== Measuring Response Time ===")

	testURLs := []string{
		"https://google.com",
		"https://github.com",
		"https://go.dev",
	}

	measureClient := &http.Client{Timeout: 10 * time.Second}

	for _, url := range testURLs {
		start := time.Now()
		resp, err := measureClient.Get(url)
		elapsed := time.Since(start)

		if err != nil {
			fmt.Printf("  ❌ %-25s Error: %s\n", url, err)
			continue
		}
		resp.Body.Close()

		fmt.Printf("  ✅ %-25s %d  %s\n", url, resp.StatusCode, elapsed.Round(time.Millisecond))
	}
}
