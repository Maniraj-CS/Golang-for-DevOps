package main

import "fmt"

func main() {
	// ================================
	// ARRAYS (fixed size - rarely used in real code)
	// ================================
	fmt.Println("=== Arrays ===")

	var ports [3]int
	ports[0] = 80
	ports[1] = 443
	ports[2] = 8080
	fmt.Println("Ports:", ports)
	fmt.Println("Length:", len(ports))

	// You CANNOT add more items to an array. Size is fixed.
	// ports[3] = 9090 // This will crash - index out of range

	// ================================
	// SLICES (dynamic size - use this everywhere)
	// ================================
	fmt.Println("\n=== Slices ===")

	// Create a slice
	servers := []string{"web-01", "web-02", "db-01"}
	fmt.Println("Servers:", servers)

	// Append (add item)
	servers = append(servers, "cache-01")
	servers = append(servers, "lb-01", "lb-02") // add multiple at once
	fmt.Println("After append:", servers)
	fmt.Println("Length:", len(servers))

	// Access by index
	fmt.Println("First:", servers[0])
	fmt.Println("Last:", servers[len(servers)-1])

	// Slicing (get a portion)
	fmt.Println("First 3:", servers[0:3]) // index 0, 1, 2
	fmt.Println("From 2:", servers[2:])   // from index 2 to end

	// Create empty slice and build it up
	var alerts []string
	fmt.Println("Empty slice:", alerts)         // []
	fmt.Println("Is nil?", alerts == nil)       // true
	fmt.Println("Length:", len(alerts))          // 0

	alerts = append(alerts, "CPU high")
	alerts = append(alerts, "Disk full")
	fmt.Println("Alerts:", alerts)

	// Make a slice with initial capacity (good for performance)
	bigSlice := make([]int, 0, 100) // length 0, but room for 100 items
	fmt.Println("Big slice length:", len(bigSlice))
	fmt.Println("Big slice capacity:", cap(bigSlice))

	// ================================
	// MAPS (key-value pairs)
	// ================================
	fmt.Println("\n=== Maps ===")

	// Create a map
	serverStatus := map[string]string{
		"web-01":   "running",
		"web-02":   "stopped",
		"db-01":    "running",
		"cache-01": "running",
	}
	fmt.Println("Status map:", serverStatus)

	// Access a value
	fmt.Println("web-01:", serverStatus["web-01"])

	// Add a new key
	serverStatus["lb-01"] = "running"

	// Update existing key
	serverStatus["web-02"] = "running"

	// Delete a key
	delete(serverStatus, "cache-01")
	fmt.Println("After changes:", serverStatus)

	// Check if key exists (very important!)
	status, exists := serverStatus["web-03"]
	if exists {
		fmt.Println("web-03:", status)
	} else {
		fmt.Println("web-03 not found in the map")
	}

	// If you just access a missing key, you get the zero value
	fmt.Println("Missing key:", serverStatus["does-not-exist"]) // "" (empty string)

	// Count items in map
	fmt.Println("Total servers:", len(serverStatus))

	// Create an empty map
	errorCounts := make(map[string]int)
	errorCounts["timeout"] = 5
	errorCounts["disk_full"] = 2
	errorCounts["auth_failed"] = 8
	fmt.Println("Error counts:", errorCounts)

	// ================================
	// LOOP OVER MAP
	// ================================
	fmt.Println("\n=== Loop Over Map ===")
	for server, status := range serverStatus {
		fmt.Printf("%-10s -> %s\n", server, status)
	}
}
