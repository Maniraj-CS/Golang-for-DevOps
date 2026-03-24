package main

import "fmt"

func main() {
	// ================================
	// 1. VARIABLES AND DATA TYPES
	// ================================

	// Using var keyword (the full way)
	var name string = "Abhishek"
	var age int = 25
	var salary float64 = 50000.50
	var isDevOps bool = true

	fmt.Println("=== Using var keyword ===")
	fmt.Println("Name:", name)
	fmt.Println("Age:", age)
	fmt.Println("Salary:", salary)
	fmt.Println("Is DevOps:", isDevOps)

	// Using short declaration (the Go way - most common)
	city := "Bangalore"
	experience := 3
	rating := 4.5
	isRemote := false

	fmt.Println("\n=== Using := short declaration ===")
	fmt.Println("City:", city)
	fmt.Println("Experience:", experience)
	fmt.Println("Rating:", rating)
	fmt.Println("Is Remote:", isRemote)

	// Constants - values that never change
	const maxRetries = 3
	const appName = "port-checker"
	const version = "1.0.0"

	fmt.Println("\n=== Constants ===")
	fmt.Println("App:", appName)
	fmt.Println("Version:", version)
	fmt.Println("Max Retries:", maxRetries)

	// You CANNOT change a constant
	// maxRetries = 5  // This will give error: cannot assign to maxRetries

	// ================================
	// 2. ZERO VALUES
	// ================================

	// If you declare a variable without giving value, Go gives a default value
	var emptyString string
	var zeroInt int
	var zeroFloat float64
	var falseBool bool

	fmt.Println("\n=== Zero Values (defaults) ===")
	fmt.Printf("Empty string: '%s'\n", emptyString) // ""
	fmt.Println("Zero int:", zeroInt)                // 0
	fmt.Println("Zero float:", zeroFloat)            // 0
	fmt.Println("Zero bool:", falseBool)             // false

	// ================================
	// 3. TYPE CONVERSION
	// ================================

	// Go does NOT do automatic type conversion. You have to do it yourself.
	myInt := 42
	myFloat := float64(myInt) // convert int to float64
	backToInt := int(myFloat) // convert float64 to int

	fmt.Println("\n=== Type Conversion ===")
	fmt.Println("Int:", myInt)
	fmt.Println("Float:", myFloat)
	fmt.Println("Back to Int:", backToInt)
}
