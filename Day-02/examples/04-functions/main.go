package main

import (
	"fmt"
	"math"
)

func main() {
	// ================================
	// BASIC FUNCTION
	// ================================
	fmt.Println("=== Basic Function ===")
	greet("Abhishek")
	greet("DevOps Team")

	// ================================
	// FUNCTION WITH RETURN VALUE
	// ================================
	fmt.Println("\n=== Function with Return ===")
	sum := add(10, 20)
	fmt.Println("10 + 20 =", sum)

	// ================================
	// MULTIPLE RETURN VALUES
	// ================================
	fmt.Println("\n=== Multiple Return Values ===")

	result, err := divide(10, 3)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("10 / 3 = %.2f\n", result)
	}

	// Try dividing by zero
	result, err = divide(10, 0)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// ================================
	// NAMED RETURN VALUES
	// ================================
	fmt.Println("\n=== Named Return Values ===")
	name, port := getServerInfo()
	fmt.Printf("Server: %s, Port: %d\n", name, port)

	// ================================
	// FUNCTION TAKING A SLICE
	// ================================
	fmt.Println("\n=== Function with Slice ===")
	ports := []int{80, 443, 8080, 3000, 9090}
	avg := average(ports)
	fmt.Printf("Average port: %.2f\n", avg)

	// ================================
	// FUNCTION TAKING A MAP
	// ================================
	fmt.Println("\n=== Function with Map ===")
	status := map[string]string{
		"web-01": "running",
		"web-02": "stopped",
		"db-01":  "running",
		"web-03": "stopped",
	}
	running, stopped := countStatus(status)
	fmt.Printf("Running: %d, Stopped: %d\n", running, stopped)
}

// Basic function - no return value
func greet(name string) {
	fmt.Printf("Hello %s! Welcome to the DevOps team.\n", name)
}

// Function with return value
func add(a int, b int) int {
	return a + b
}

// Function with multiple return values
// This is the most common pattern in Go - (result, error)
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("cannot divide by zero")
	}
	result := a / b
	result = math.Round(result*100) / 100 // round to 2 decimal places
	return result, nil
}

// Function with named return values
func getServerInfo() (name string, port int) {
	name = "web-server-01"
	port = 8080
	return // naked return - returns the named values
}

// Function that takes a slice
func average(numbers []int) float64 {
	if len(numbers) == 0 {
		return 0
	}
	total := 0
	for _, n := range numbers {
		total += n
	}
	return float64(total) / float64(len(numbers))
}

// Function that takes a map and returns multiple values
func countStatus(servers map[string]string) (int, int) {
	running := 0
	stopped := 0
	for _, status := range servers {
		switch status {
		case "running":
			running++
		case "stopped":
			stopped++
		}
	}
	return running, stopped
}
