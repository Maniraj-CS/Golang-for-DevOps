package main

import "fmt"

// ================================
// STRUCTS - Grouping related data
// ================================

// Basic struct - represents a server in our infra
type Server struct {
	Name   string
	IP     string
	Port   int
	Status string
}

// Nested struct - health check config inside a server
type HealthCheck struct {
	Endpoint string
	Interval int // in seconds
}

type MonitoredServer struct {
	Name        string
	IP          string
	Port        int
	HealthCheck HealthCheck // struct inside a struct
}

func main() {
	// ================================
	// 1. CREATING STRUCTS
	// ================================

	fmt.Println("=== Creating Structs ===")

	// Method 1: Using field names (recommended - always use this)
	web01 := Server{
		Name:   "web-01",
		IP:     "10.0.1.10",
		Port:   8080,
		Status: "running",
	}
	fmt.Printf("Server: %+v\n", web01)

	// Method 2: Without field names (not recommended - hard to read)
	web02 := Server{"web-02", "10.0.1.11", 8080, "stopped"}
	fmt.Printf("Server: %+v\n", web02)

	// Method 3: Zero value then set fields one by one
	var db01 Server
	db01.Name = "db-01"
	db01.IP = "10.0.2.10"
	db01.Port = 5432
	db01.Status = "running"
	fmt.Printf("Server: %+v\n", db01)

	// ================================
	// 2. ACCESSING AND MODIFYING FIELDS
	// ================================

	fmt.Println("\n=== Accessing Fields ===")
	fmt.Println("Name:", web01.Name)
	fmt.Println("IP:", web01.IP)
	fmt.Println("Address:", web01.IP, ":", web01.Port)

	// Modify a field
	web01.Status = "stopped"
	fmt.Println("Updated status:", web01.Status)

	// ================================
	// 3. SLICE OF STRUCTS
	// ================================

	fmt.Println("\n=== Slice of Structs ===")
	servers := []Server{
		{Name: "web-01", IP: "10.0.1.10", Port: 8080, Status: "running"},
		{Name: "web-02", IP: "10.0.1.11", Port: 8080, Status: "stopped"},
		{Name: "api-01", IP: "10.0.1.20", Port: 3000, Status: "running"},
		{Name: "db-01", IP: "10.0.2.10", Port: 5432, Status: "running"},
	}

	for _, s := range servers {
		status := "✅"
		if s.Status == "stopped" {
			status = "❌"
		}
		fmt.Printf("  %s %s (%s:%d)\n", status, s.Name, s.IP, s.Port)
	}

	// ================================
	// 4. NESTED STRUCTS
	// ================================

	fmt.Println("\n=== Nested Structs ===")
	monitored := MonitoredServer{
		Name: "web-01",
		IP:   "10.0.1.10",
		Port: 8080,
		HealthCheck: HealthCheck{
			Endpoint: "/health",
			Interval: 30,
		},
	}

	fmt.Printf("Server: %s\n", monitored.Name)
	fmt.Printf("Health endpoint: %s\n", monitored.HealthCheck.Endpoint)
	fmt.Printf("Check every: %d seconds\n", monitored.HealthCheck.Interval)

	// ================================
	// 5. COMPARING STRUCTS
	// ================================

	fmt.Println("\n=== Comparing Structs ===")
	s1 := Server{Name: "web-01", IP: "10.0.1.10", Port: 8080, Status: "running"}
	s2 := Server{Name: "web-01", IP: "10.0.1.10", Port: 8080, Status: "running"}
	s3 := Server{Name: "web-02", IP: "10.0.1.11", Port: 8080, Status: "running"}

	fmt.Println("s1 == s2:", s1 == s2) // true - all fields match
	fmt.Println("s1 == s3:", s1 == s3) // false - different name and IP
}
