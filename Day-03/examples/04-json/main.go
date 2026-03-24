package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// ================================
// JSON - Marshal and Unmarshal
// ================================

// Struct tags tell the json package what field names to use
type Server struct {
	Name     string `json:"name"`
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Status   string `json:"status,omitempty"` // skip if empty string
	Password string `json:"-"`               // NEVER include in JSON
}

func main() {
	// ================================
	// 1. STRUCT TO JSON (Marshal)
	// ================================

	fmt.Println("=== Struct to JSON (Marshal) ===")

	server := Server{
		Name:     "web-01",
		IP:       "10.0.1.10",
		Port:     8080,
		Status:   "running",
		Password: "super-secret-123", // this will NOT appear in JSON
	}

	// Marshal - converts struct to JSON bytes
	data, err := json.Marshal(server)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("JSON:", string(data))

	// MarshalIndent - pretty printed JSON (great for config files)
	prettyData, err := json.MarshalIndent(server, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nPretty JSON:")
	fmt.Println(string(prettyData))

	// ================================
	// 2. JSON TO STRUCT (Unmarshal)
	// ================================

	fmt.Println("\n=== JSON to Struct (Unmarshal) ===")

	jsonStr := `{"name":"api-01","ip":"10.0.1.20","port":3000,"status":"running"}`

	var parsed Server
	err = json.Unmarshal([]byte(jsonStr), &parsed)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Name: %s\n", parsed.Name)
	fmt.Printf("IP: %s\n", parsed.IP)
	fmt.Printf("Port: %d\n", parsed.Port)
	fmt.Printf("Status: %s\n", parsed.Status)

	// ================================
	// 3. SLICE OF STRUCTS TO JSON
	// ================================

	fmt.Println("\n=== Slice of Structs to JSON ===")

	servers := []Server{
		{Name: "web-01", IP: "10.0.1.10", Port: 8080, Status: "running"},
		{Name: "web-02", IP: "10.0.1.11", Port: 8080, Status: "stopped"},
		{Name: "db-01", IP: "10.0.2.10", Port: 5432, Status: "running"},
	}

	output, err := json.MarshalIndent(servers, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(output))

	// ================================
	// 4. JSON ARRAY TO SLICE OF STRUCTS
	// ================================

	fmt.Println("\n=== JSON Array to Slice of Structs ===")

	jsonArray := `[
		{"name":"cache-01","ip":"10.0.3.10","port":6379,"status":"running"},
		{"name":"cache-02","ip":"10.0.3.11","port":6379}
	]`

	var parsedServers []Server
	err = json.Unmarshal([]byte(jsonArray), &parsedServers)
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range parsedServers {
		status := s.Status
		if status == "" {
			status = "(no status)" // omitempty means this field was missing in JSON
		}
		fmt.Printf("  %s -> %s:%d [%s]\n", s.Name, s.IP, s.Port, status)
	}

	// ================================
	// 5. OMITEMPTY DEMO
	// ================================

	fmt.Println("\n=== omitempty Demo ===")

	// Status is empty - it will be skipped in JSON because of omitempty tag
	noStatus := Server{
		Name: "test-01",
		IP:   "10.0.4.10",
		Port: 8080,
	}

	data, _ = json.MarshalIndent(noStatus, "", "  ")
	fmt.Println(string(data))
	// Notice: "status" field is not in the output because it's empty

	// ================================
	// 6. WORKING WITH MAP (dynamic JSON)
	// ================================

	fmt.Println("\n=== Dynamic JSON with map ===")

	// When you don't know the JSON structure, use map[string]any
	dynamicJSON := `{"name":"unknown","cpu":85.5,"healthy":true,"tags":["prod","us-east"]}`

	var result map[string]any
	err = json.Unmarshal([]byte(dynamicJSON), &result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Name:", result["name"])
	fmt.Println("CPU:", result["cpu"])
	fmt.Println("Healthy:", result["healthy"])
	fmt.Println("Tags:", result["tags"])
}
