package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	// We use httpbin.org - it's a free service that echoes back whatever we send.
	// Very useful for testing HTTP requests!

	client := &http.Client{Timeout: 10 * time.Second}

	// ================================
	// 1. POST REQUEST - Send JSON data
	// ================================

	fmt.Println("=== POST Request ===")

	postBody := map[string]string{
		"title": "Bug: Server crash on deploy",
		"body":  "App crashes when deploying v2.3.1",
	}

	jsonBody, err := json.Marshal(postBody)
	if err != nil {
		log.Fatal("Cannot marshal JSON:", err)
	}

	// Create a POST request
	req, err := http.NewRequest("POST", "https://httpbin.org/post", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Fatal("Cannot create request:", err)
	}

	// Set headers - always set Content-Type when sending JSON
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Custom-Header", "devops-tool")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("POST failed:", err)
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)

	body, _ := io.ReadAll(resp.Body)
	// httpbin echoes back what we sent - let's parse it
	var result map[string]any
	json.Unmarshal(body, &result)

	// Show the JSON we sent (httpbin puts it in "json" field)
	if data, ok := result["json"]; ok {
		fmt.Println("Server received our JSON:", data)
	}

	// Show our custom headers (httpbin puts them in "headers" field)
	if headers, ok := result["headers"].(map[string]any); ok {
		fmt.Println("Custom header received:", headers["X-Custom-Header"])
	}

	// ================================
	// 2. PUT REQUEST - Update something
	// ================================

	fmt.Println("\n=== PUT Request ===")

	updateBody := map[string]string{
		"title":  "Bug: Server crash on deploy (FIXED)",
		"status": "closed",
	}

	jsonBody, _ = json.Marshal(updateBody)
	req, err = http.NewRequest("PUT", "https://httpbin.org/put", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Fatal("Cannot create PUT request:", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer fake-token-12345")

	resp, err = client.Do(req)
	if err != nil {
		log.Fatal("PUT failed:", err)
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)

	body, _ = io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	if data, ok := result["json"]; ok {
		fmt.Println("Server received our update:", data)
	}

	// ================================
	// 3. DELETE REQUEST
	// ================================

	fmt.Println("\n=== DELETE Request ===")

	// DELETE usually has no body, just headers
	req, err = http.NewRequest("DELETE", "https://httpbin.org/delete", nil)
	if err != nil {
		log.Fatal("Cannot create DELETE request:", err)
	}
	req.Header.Set("Authorization", "Bearer fake-token-12345")

	resp, err = client.Do(req)
	if err != nil {
		log.Fatal("DELETE failed:", err)
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
	fmt.Println("DELETE request sent successfully!")

	// ================================
	// 4. SIMPLE POST (shortcut)
	// ================================

	fmt.Println("\n=== Simple POST (shortcut) ===")

	// http.Post is a shortcut when you don't need custom headers
	simpleBody := []byte(`{"name": "web-server-01", "action": "restart"}`)
	resp, err = http.Post(
		"https://httpbin.org/post",
		"application/json",
		bytes.NewBuffer(simpleBody),
	)
	if err != nil {
		log.Fatal("Simple POST failed:", err)
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)

	body, _ = io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	if data, ok := result["json"]; ok {
		fmt.Println("Server received:", data)
	}

	// ================================
	// KEY TAKEAWAYS
	// ================================

	fmt.Println("\n=== Key Takeaways ===")
	fmt.Println("1. Use http.NewRequest() when you need custom headers or auth tokens")
	fmt.Println("2. Use http.Post() as a shortcut for simple POST requests")
	fmt.Println("3. Always set Content-Type header when sending JSON")
	fmt.Println("4. Always close response body with defer resp.Body.Close()")
	fmt.Println("5. For PUT/DELETE, you must use http.NewRequest() - no shortcuts exist")
}
