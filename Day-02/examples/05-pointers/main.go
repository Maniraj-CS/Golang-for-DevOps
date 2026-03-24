package main

import "fmt"

func main() {
	// ================================
	// WHAT IS A POINTER?
	// ================================
	// A pointer holds the memory address of a variable.
	// & = address of (give me the address)
	// * = value at   (give me the value at this address)

	fmt.Println("=== Pointer Basics ===")

	name := "Abhishek"
	ptr := &name // ptr now holds the address of name

	fmt.Println("Value:", name)         // Abhishek
	fmt.Println("Address:", ptr)        // 0xc000010230 (some address)
	fmt.Println("Value via ptr:", *ptr) // Abhishek

	// Change value using pointer
	*ptr = "DevOps Engineer"
	fmt.Println("After change:", name) // DevOps Engineer

	// ================================
	// WHY POINTERS? - MODIFYING VALUES IN FUNCTIONS
	// ================================
	fmt.Println("\n=== Without Pointer (value is copied) ===")

	port := 8080
	changePortCopy(port) // this does NOT change the original
	fmt.Println("Port after changePortCopy:", port) // still 8080

	fmt.Println("\n=== With Pointer (original is changed) ===")

	changePortReal(&port) // this CHANGES the original
	fmt.Println("Port after changePortReal:", port) // 9090

	// ================================
	// PRACTICAL EXAMPLE
	// ================================
	fmt.Println("\n=== Practical: Increment Counter ===")

	counter := 0
	increment(&counter)
	increment(&counter)
	increment(&counter)
	fmt.Println("Counter:", counter) // 3
}

// This gets a COPY of the value. Changes don't affect original.
func changePortCopy(port int) {
	port = 9090
	fmt.Println("Inside function:", port) // 9090
}

// This gets the ADDRESS. Changes affect the original.
func changePortReal(port *int) {
	*port = 9090
	fmt.Println("Inside function:", *port) // 9090
}

// Increment a counter using pointer
func increment(count *int) {
	*count++
}
