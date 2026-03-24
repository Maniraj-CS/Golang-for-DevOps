package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	// Check if user gave us a port number
	if len(os.Args) < 2 {
		fmt.Println("Usage: port-checker <port>")
		fmt.Println("Example: port-checker 8080")
		os.Exit(1)
	}

	// Get the port from command line
	port := os.Args[1]

	// Tell user what we are checking
	fmt.Printf("Checking port %s on localhost...\n", port)

	// Build the address - localhost:port
	address := "localhost:" + port

	// Try to connect with a 3 second timeout
	// net.DialTimeout tries to make a TCP connection
	// If it connects, port is open. If it fails, port is closed.
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)

	// Check if connection failed
	if err != nil {
		fmt.Printf("Port %s is CLOSED\n", port)
		fmt.Printf("Reason: %s\n", err)
		return
	}

	// If we reach here, connection was successful
	// Always close the connection when done
	conn.Close()
	fmt.Printf("Port %s is OPEN\n", port)
}
