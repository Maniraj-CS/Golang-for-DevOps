package main

// ================================
// EXAMPLE 01 - Context and Timeout
// ================================
// This shows how context.Context works for timeouts and cancellation.
// Run: go run main.go

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// ================================
	// 1. BASIC TIMEOUT
	// ================================

	fmt.Println("=== Context with Timeout ===")

	// Create a context that cancels after 2 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // always call cancel!

	// Simulate a fast operation (finishes in time)
	result, err := doWork(ctx, 1*time.Second)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}

	// ================================
	// 2. TIMEOUT EXCEEDED
	// ================================

	fmt.Println("\n=== Timeout Exceeded ===")

	ctx2, cancel2 := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel2()

	// This operation takes 3 seconds - it will timeout!
	result, err = doWork(ctx2, 3*time.Second)
	if err != nil {
		if ctx2.Err() == context.DeadlineExceeded {
			fmt.Println("Operation timed out! (expected)")
		} else {
			fmt.Println("Error:", err)
		}
	} else {
		fmt.Println("Result:", result)
	}

	// ================================
	// 3. MANUAL CANCELLATION
	// ================================

	fmt.Println("\n=== Manual Cancellation ===")

	ctx3, cancel3 := context.WithCancel(context.Background())

	// Start work in a goroutine
	done := make(chan string)
	go func() {
		result, err := doWork(ctx3, 5*time.Second)
		if err != nil {
			done <- fmt.Sprintf("cancelled: %v", err)
			return
		}
		done <- result
	}()

	// Cancel after 500ms
	time.Sleep(500 * time.Millisecond)
	cancel3() // cancel the context!

	fmt.Println("Result:", <-done)

	// ================================
	// 4. CONTEXT WITH VALUE (passing data)
	// ================================

	fmt.Println("\n=== Context with Value ===")

	// You can attach values to context (useful for request IDs, user info)
	ctx4 := context.WithValue(context.Background(), "requestID", "req-12345")
	ctx4 = context.WithValue(ctx4, "user", "abhishek")

	processRequest(ctx4)

	// ================================
	// 5. REAL-WORLD PATTERN
	// ================================

	fmt.Println("\n=== Real-World Pattern ===")
	fmt.Println("In DevOps tools, you use context for:")
	fmt.Println("  - Docker API calls (timeout if daemon is slow)")
	fmt.Println("  - K8s API calls (timeout if cluster is unreachable)")
	fmt.Println("  - HTTP requests (we already used this in Day-07)")
	fmt.Println("  - Database queries")
	fmt.Println("  - Any operation that might hang")
}

// doWork simulates an operation that takes some time
func doWork(ctx context.Context, duration time.Duration) (string, error) {
	select {
	case <-time.After(duration):
		// Work completed
		return "done!", nil
	case <-ctx.Done():
		// Context was cancelled or timed out
		return "", ctx.Err()
	}
}

// processRequest shows how to read values from context
func processRequest(ctx context.Context) {
	reqID := ctx.Value("requestID").(string)
	user := ctx.Value("user").(string)
	fmt.Printf("Processing request %s for user %s\n", reqID, user)
}
