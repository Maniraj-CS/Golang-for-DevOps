package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// ================================
// TYPES
// ================================

type Server struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	IP     string `json:"ip"`
}

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// In-memory data (in real apps, this would come from a database)
var servers = []Server{
	{Name: "web-01", Status: "running", IP: "10.0.1.10"},
	{Name: "web-02", Status: "running", IP: "10.0.1.11"},
	{Name: "db-01", Status: "running", IP: "10.0.2.10"},
	{Name: "cache-01", Status: "stopped", IP: "10.0.3.10"},
}

func main() {
	// ================================
	// CREATING A SERVER WITH ServeMux
	// ================================

	// ServeMux is Go's built-in HTTP router.
	// In recent Go versions, it supports method-based routing!

	mux := http.NewServeMux()

	// ================================
	// 1. HEALTH CHECK ENDPOINT
	// ================================

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(HealthResponse{
			Status:    "ok",
			Timestamp: time.Now().Format(time.RFC3339),
		})
	})

	// ================================
	// 2. LIST ALL SERVERS
	// ================================

	mux.HandleFunc("GET /api/servers", func(w http.ResponseWriter, r *http.Request) {
		// Read query parameter for filtering
		status := r.URL.Query().Get("status")

		var result []Server
		if status != "" {
			// Filter by status
			for _, s := range servers {
				if s.Status == status {
					result = append(result, s)
				}
			}
		} else {
			result = servers
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	// ================================
	// 3. GET SINGLE SERVER BY NAME (path parameter!)
	// ================================

	mux.HandleFunc("GET /api/servers/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name") // This is the new Go way to read path params!

		for _, s := range servers {
			if s.Name == name {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(s)
				return
			}
		}

		// Server not found
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("server '%s' not found", name),
		})
	})

	// ================================
	// 4. ADD A NEW SERVER (POST)
	// ================================

	mux.HandleFunc("POST /api/servers", func(w http.ResponseWriter, r *http.Request) {
		var newServer Server
		err := json.NewDecoder(r.Body).Decode(&newServer)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid JSON body",
			})
			return
		}

		// Basic validation
		if newServer.Name == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "name is required",
			})
			return
		}

		servers = append(servers, newServer)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newServer)
	})

	// ================================
	// 5. DELETE A SERVER
	// ================================

	mux.HandleFunc("DELETE /api/servers/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")

		for i, s := range servers {
			if s.Name == name {
				servers = append(servers[:i], servers[i+1:]...)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]string{
					"message": fmt.Sprintf("server '%s' deleted", name),
				})
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("server '%s' not found", name),
		})
	})

	// ================================
	// START THE SERVER
	// ================================

	port := ":8080"
	fmt.Println("=== DevOps Server API ===")
	fmt.Println()
	fmt.Println("Endpoints:")
	fmt.Println("  GET    /health                  - Health check")
	fmt.Println("  GET    /api/servers              - List all servers")
	fmt.Println("  GET    /api/servers?status=running - Filter by status")
	fmt.Println("  GET    /api/servers/{name}       - Get one server")
	fmt.Println("  POST   /api/servers              - Add a server")
	fmt.Println("  DELETE /api/servers/{name}       - Delete a server")
	fmt.Println()
	fmt.Println("Try these commands in another terminal:")
	fmt.Println("  curl http://localhost:8080/health")
	fmt.Println("  curl http://localhost:8080/api/servers")
	fmt.Println("  curl http://localhost:8080/api/servers/web-01")
	fmt.Println("  curl http://localhost:8080/api/servers?status=running")
	fmt.Println(`  curl -X POST http://localhost:8080/api/servers -d '{"name":"app-01","status":"running","ip":"10.0.4.10"}'`)
	fmt.Println(`  curl -X DELETE http://localhost:8080/api/servers/cache-01`)
	fmt.Println()
	fmt.Printf("Server starting on %s ...\n", port)

	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
