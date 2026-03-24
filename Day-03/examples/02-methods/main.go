package main

import "fmt"

// ================================
// METHODS ON STRUCTS
// ================================

type Server struct {
	Name   string
	IP     string
	Port   int
	Status string
}

// VALUE RECEIVER - just reads data, does NOT modify the struct
func (s Server) Address() string {
	return fmt.Sprintf("%s:%d", s.IP, s.Port)
}

// VALUE RECEIVER - returns a nice display string
func (s Server) Summary() string {
	return fmt.Sprintf("[%s] %s - %s (%s)", s.Status, s.Name, s.Address(), s.Status)
}

// POINTER RECEIVER - CAN modify the struct
func (s *Server) Start() {
	s.Status = "running"
	fmt.Printf("  ✅ %s started\n", s.Name)
}

// POINTER RECEIVER - CAN modify the struct
func (s *Server) Stop() {
	s.Status = "stopped"
	fmt.Printf("  ❌ %s stopped\n", s.Name)
}

// POINTER RECEIVER - restart = stop then start
func (s *Server) Restart() {
	fmt.Printf("  🔄 Restarting %s...\n", s.Name)
	s.Stop()
	s.Start()
}

// Value receiver that checks something
func (s Server) IsRunning() bool {
	return s.Status == "running"
}

func main() {
	// ================================
	// 1. USING METHODS
	// ================================

	fmt.Println("=== Using Methods ===")

	server := Server{
		Name:   "web-01",
		IP:     "10.0.1.10",
		Port:   8080,
		Status: "running",
	}

	fmt.Println("Address:", server.Address())
	fmt.Println("Running?", server.IsRunning())

	// ================================
	// 2. POINTER RECEIVER - MODIFYING STRUCT
	// ================================

	fmt.Println("\n=== Pointer Receiver (modifies struct) ===")

	fmt.Println("Before stop - Status:", server.Status)
	server.Stop()
	fmt.Println("After stop - Status:", server.Status)
	fmt.Println("Running?", server.IsRunning())

	server.Start()
	fmt.Println("After start - Status:", server.Status)

	// ================================
	// 3. RESTART
	// ================================

	fmt.Println("\n=== Restart ===")
	server.Restart()
	fmt.Println("Status after restart:", server.Status)

	// ================================
	// 4. METHODS ON A SLICE OF SERVERS
	// ================================

	fmt.Println("\n=== Managing Multiple Servers ===")
	servers := []Server{
		{Name: "web-01", IP: "10.0.1.10", Port: 8080, Status: "running"},
		{Name: "web-02", IP: "10.0.1.11", Port: 8080, Status: "stopped"},
		{Name: "api-01", IP: "10.0.1.20", Port: 3000, Status: "running"},
		{Name: "db-01", IP: "10.0.2.10", Port: 5432, Status: "stopped"},
	}

	fmt.Println("\nStarting all stopped servers:")
	for i := range servers {
		if !servers[i].IsRunning() {
			servers[i].Start() // must use index, not range copy!
		}
	}

	fmt.Println("\nAll servers after starting:")
	for _, s := range servers {
		status := "✅"
		if !s.IsRunning() {
			status = "❌"
		}
		fmt.Printf("  %s %s -> %s\n", status, s.Name, s.Address())
	}
}
