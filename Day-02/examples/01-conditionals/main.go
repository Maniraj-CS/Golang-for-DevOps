package main

import "fmt"

func main() {
	// ================================
	// IF-ELSE
	// ================================
	fmt.Println("=== If-Else ===")

	port := 8080

	if port == 80 {
		fmt.Println("HTTP")
	} else if port == 443 {
		fmt.Println("HTTPS")
	} else {
		fmt.Printf("Custom port: %d\n", port)
	}

	// Go special: declare variable inside if
	// The variable only exists inside the if-else block
	if serverPort := 3000; serverPort < 1024 {
		fmt.Println("Privileged port - need root")
	} else {
		fmt.Printf("Port %d - no root needed\n", serverPort)
	}
	// serverPort is NOT available here

	// ================================
	// SWITCH CASE
	// ================================
	fmt.Println("\n=== Switch Case ===")

	env := "production"

	switch env {
	case "development":
		fmt.Println("Dev mode - debug logging ON")
	case "staging":
		fmt.Println("Staging mode - testing")
	case "production":
		fmt.Println("Prod mode - be careful!")
	default:
		fmt.Println("Unknown environment")
	}

	// Switch without a condition (like multiple if-else)
	hour := 14

	switch {
	case hour < 12:
		fmt.Println("Good morning!")
	case hour < 17:
		fmt.Println("Good afternoon!")
	default:
		fmt.Println("Good evening!")
	}

	// Switch with multiple values in one case
	day := "Saturday"

	switch day {
	case "Monday", "Tuesday", "Wednesday", "Thursday", "Friday":
		fmt.Println("Weekday - deploy with caution")
	case "Saturday", "Sunday":
		fmt.Println("Weekend - no deploys please!")
	}
}
