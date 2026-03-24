package main

import "fmt"

func main() {
	// ================================
	// 1. DEFER - runs when function exits
	// ================================

	fmt.Println("=== Defer Basics ===")
	deferDemo()

	// ================================
	// 2. DEFER IS LIFO (last in, first out)
	// ================================

	fmt.Println("\n=== Defer Order (LIFO) ===")
	deferOrder()

	// ================================
	// 3. DEFER WITH LOOPS (common DevOps pattern)
	// ================================

	fmt.Println("\n=== Defer in Cleanup ===")
	cleanupDemo()

	// ================================
	// 4. PANIC - unrecoverable errors
	// ================================

	fmt.Println("\n=== Panic and Recover ===")
	panicDemo()

	// ================================
	// 5. SAFE RUNNER (recover pattern)
	// ================================

	fmt.Println("\n=== Safe Runner Pattern ===")
	safeRunnerDemo()

	fmt.Println("\n✅ Program completed normally (panics were recovered!)")
}

func deferDemo() {
	fmt.Println("  Step 1: Opening resource")
	defer fmt.Println("  Step 4: Closing resource (deferred - runs last)")
	fmt.Println("  Step 2: Using resource")
	fmt.Println("  Step 3: Done using resource")
	// Step 4 runs here, when function returns
}

func deferOrder() {
	// Defers run in LIFO order - like a stack
	defer fmt.Println("  3rd defer (runs FIRST)")
	defer fmt.Println("  2nd defer (runs SECOND)")
	defer fmt.Println("  1st defer (runs THIRD)")
	fmt.Println("  Main code runs first")
}

func cleanupDemo() {
	// Common pattern: track resources to clean up
	resources := []string{"db-connection", "file-handle", "cache-client"}

	for _, r := range resources {
		r := r // capture for defer
		defer fmt.Printf("  🧹 Cleaned up: %s\n", r)
		fmt.Printf("  📦 Opened: %s\n", r)
	}

	fmt.Println("  (doing work with all resources...)")
	// When function exits, all defers run in reverse order
}

func panicDemo() {
	// recover must be in a defer to catch panics
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("  🛡️  Recovered from panic: %v\n", r)
		}
	}()

	fmt.Println("  About to panic...")
	panic("something went terribly wrong!")

	// This line never executes
	// fmt.Println("This won't print")
}

func safeRunnerDemo() {
	// This pattern is useful for running risky operations
	// without crashing the whole program

	safeRun := func(name string, fn func()) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("  ❌ Task %q panicked: %v\n", name, r)
			}
		}()

		fmt.Printf("  ⏳ Running: %s\n", name)
		fn()
		fmt.Printf("  ✅ Completed: %s\n", name)
	}

	// Task 1 - works fine
	safeRun("check disk", func() {
		// normal operation
	})

	// Task 2 - panics but gets caught
	safeRun("check memory", func() {
		panic("out of memory!")
	})

	// Task 3 - still runs! The panic didn't crash everything
	safeRun("check network", func() {
		// normal operation
	})
}
