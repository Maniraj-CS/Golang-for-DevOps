package main

import "fmt"

func main() {
	// ================================
	// fmt.Println - simple print with newline
	// ================================
	fmt.Println("=== Println ===")
	fmt.Println("Hello", "DevOps", "Engineers")
	// Output: Hello DevOps Engineers
	// It adds spaces between values and a newline at the end

	// ================================
	// fmt.Printf - formatted print (like printf in C)
	// ================================
	fmt.Println("\n=== Printf ===")

	name := "Abhishek"
	age := 25
	salary := 75000.50
	isDevOps := true

	// %s = string
	fmt.Printf("Name: %s\n", name)

	// %d = integer (decimal)
	fmt.Printf("Age: %d\n", age)

	// %f = float (decimal number)
	fmt.Printf("Salary: %f\n", salary)

	// %.2f = float with 2 decimal places
	fmt.Printf("Salary (2 decimal): %.2f\n", salary)

	// %t = boolean
	fmt.Printf("Is DevOps: %t\n", isDevOps)

	// %v = any value (Go figures out the type)
	fmt.Printf("Name: %v, Age: %v, DevOps: %v\n", name, age, isDevOps)

	// %T = print the TYPE of a variable
	fmt.Printf("Type of name: %T\n", name)
	fmt.Printf("Type of age: %T\n", age)
	fmt.Printf("Type of salary: %T\n", salary)

	// ================================
	// fmt.Sprintf - like Printf but returns a string instead of printing
	// ================================
	fmt.Println("\n=== Sprintf ===")

	server := "web-server-01"
	port := 8080
	status := "running"

	// This is super useful for building messages, log lines, etc.
	msg := fmt.Sprintf("Server %s is %s on port %d", server, status, port)
	fmt.Println(msg)

	// Another example - building a URL
	host := "api.github.com"
	path := "/repos"
	url := fmt.Sprintf("https://%s%s", host, path)
	fmt.Println("URL:", url)

	// ================================
	// fmt.Print - print without newline
	// ================================
	fmt.Println("\n=== Print (no newline) ===")
	fmt.Print("Hello ")
	fmt.Print("World")
	fmt.Print("\n") // we have to add newline ourselves
}
