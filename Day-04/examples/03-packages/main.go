package main

import (
	"fmt"

	// Import our own package using module path + folder name
	"log-parser/examples/03-packages/utils"
)

// init runs before main - good for setup
func init() {
	fmt.Println("=== init() function ran first ===")
	fmt.Println("(This is where you'd load config, check environment, etc.)")
	fmt.Println()
}

func main() {
	// ================================
	// 1. USING OUR PACKAGE
	// ================================

	fmt.Println("=== Using Our utils Package ===")

	// FormatLevel is exported (starts with uppercase) - we can use it!
	levels := []string{"  error  ", "warn", "INFO", "debug", "CRITICAL"}

	for _, level := range levels {
		formatted := utils.FormatLevel(level)
		valid := utils.IsValidLevel(level)

		status := "✅"
		if !valid {
			status = "❌"
		}

		fmt.Printf("  %s Input: %-12q -> Formatted: %-6s Valid: %v\n", status, level, formatted, valid)
	}

	// ================================
	// 2. FORMAT BYTES
	// ================================

	fmt.Println("\n=== FormatBytes ===")

	sizes := []int64{512, 2048, 1048576, 5368709120}
	for _, size := range sizes {
		fmt.Printf("  %d bytes -> %s\n", size, utils.FormatBytes(size))
	}

	// ================================
	// 3. USING SEPARATOR
	// ================================

	fmt.Println("\n=== PrintSeparator (uses unexported helper internally) ===")
	utils.PrintSeparator(40)
	fmt.Println("This is between two separators")
	utils.PrintSeparator(40)

	// This would NOT work - repeatChar is unexported:
	// utils.repeatChar("-", 10)  // compile error: unexported function

	fmt.Println("\n✅ Packages working! Notice how:")
	fmt.Println("  - FormatLevel, IsValidLevel, FormatBytes are EXPORTED (uppercase)")
	fmt.Println("  - repeatChar is unexported (lowercase) - we can't call it from here")
	fmt.Println("  - init() ran automatically before main()")
}
