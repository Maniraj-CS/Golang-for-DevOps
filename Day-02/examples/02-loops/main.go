package main

import "fmt"

func main() {
	// ================================
	// CLASSIC FOR LOOP
	// ================================
	fmt.Println("=== Classic For Loop ===")
	for i := 0; i < 5; i++ {
		fmt.Printf("Iteration %d\n", i)
	}

	// ================================
	// WHILE STYLE (for with just condition)
	// ================================
	fmt.Println("\n=== While Style Loop ===")
	retries := 0
	maxRetries := 3
	for retries < maxRetries {
		fmt.Printf("Retry attempt %d of %d\n", retries+1, maxRetries)
		retries++
	}

	// ================================
	// RANGE OVER INTEGERS
	// ================================
	fmt.Println("\n=== Range Over Integers ===")
	for i := range 5 {
		fmt.Printf("Server-%d\n", i+1)
	}

	// ================================
	// INFINITE LOOP WITH BREAK
	// ================================
	fmt.Println("\n=== Infinite Loop with Break ===")
	count := 0
	for {
		if count >= 3 {
			fmt.Println("Breaking out!")
			break
		}
		fmt.Printf("Count: %d\n", count)
		count++
	}

	// ================================
	// CONTINUE - skip current iteration
	// ================================
	fmt.Println("\n=== Continue (skip even numbers) ===")
	for i := range 10 {
		if i%2 == 0 {
			continue // skip even numbers
		}
		fmt.Printf("%d is odd\n", i)
	}

	// ================================
	// LOOP OVER A SLICE
	// ================================
	fmt.Println("\n=== Loop Over Slice ===")
	servers := []string{"web-01", "web-02", "db-01", "cache-01"}

	// With index and value
	for index, server := range servers {
		fmt.Printf("[%d] %s\n", index, server)
	}

	// Without index (use _ to ignore it)
	fmt.Println("\nJust values:")
	for _, server := range servers {
		fmt.Println(server)
	}

	// ================================
	// LOOP OVER A MAP
	// ================================
	fmt.Println("\n=== Loop Over Map ===")
	statusMap := map[string]string{
		"web-01": "running",
		"web-02": "stopped",
		"db-01":  "running",
	}

	for server, status := range statusMap {
		fmt.Printf("%s -> %s\n", server, status)
	}
}
