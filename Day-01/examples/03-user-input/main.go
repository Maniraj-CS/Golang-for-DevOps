package main

import (
	"fmt"
	"os"
)

func main() {
	// ================================
	// METHOD 1: os.Args (Command Line Arguments)
	// ================================

	// os.Args is a slice (like an array) that holds all command line arguments
	// os.Args[0] = program name
	// os.Args[1] = first argument
	// os.Args[2] = second argument ... and so on

	fmt.Println("=== os.Args ===")
	fmt.Println("Total arguments:", len(os.Args))
	fmt.Println("Program name:", os.Args[0])

	// Print all arguments
	for i, arg := range os.Args {
		fmt.Printf("  os.Args[%d] = %s\n", i, arg)
	}

	// Check if user gave an argument
	if len(os.Args) >= 2 {
		fmt.Printf("\nHello %s! Welcome to Go!\n", os.Args[1])
	} else {
		fmt.Println("\nNo name given. Usage: go run main.go <your-name>")
	}

	// ================================
	// METHOD 2: fmt.Scan (Interactive Input)
	// ================================

	fmt.Println("\n=== fmt.Scan ===")

	var city string
	fmt.Print("Enter your city: ")
	fmt.Scan(&city) // & means we are passing the address of the variable
	fmt.Printf("You are from %s!\n", city)

	var port int
	fmt.Print("Enter a port number: ")
	fmt.Scan(&port)
	fmt.Printf("You entered port: %d\n", port)
}

// Run this file:
// go run main.go Abhishek
// go run main.go Abhishek DevOps
