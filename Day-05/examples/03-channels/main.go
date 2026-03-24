package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// ================================
	// 1. BASIC CHANNEL
	// ================================

	fmt.Println("=== Basic Channel ===")

	ch := make(chan string) // unbuffered channel

	// Send data in a goroutine
	go func() {
		time.Sleep(100 * time.Millisecond) // simulate some work
		ch <- "Hello from goroutine!"      // send into channel
	}()

	// Receive data (blocks until something is sent)
	msg := <-ch
	fmt.Println("Received:", msg)

	// ================================
	// 2. BUFFERED CHANNEL
	// ================================

	fmt.Println("\n=== Buffered Channel ===")

	// Buffered channel can hold 3 values without blocking
	buffered := make(chan string, 3)

	// These don't block because buffer has space
	buffered <- "first"
	buffered <- "second"
	buffered <- "third"

	fmt.Println(<-buffered) // first
	fmt.Println(<-buffered) // second
	fmt.Println(<-buffered) // third

	// ================================
	// 3. CHANNEL WITH RANGE (receive until closed)
	// ================================

	fmt.Println("\n=== Channel with Range ===")

	numbers := make(chan int, 5)

	go func() {
		for i := 1; i <= 5; i++ {
			numbers <- i * 10
		}
		close(numbers) // MUST close when done sending!
	}()

	// range automatically stops when channel is closed
	for num := range numbers {
		fmt.Printf("  Got: %d\n", num)
	}

	// ================================
	// 4. COLLECTING RESULTS FROM GOROUTINES
	// ================================

	fmt.Println("\n=== Collecting Results (DevOps Pattern) ===")

	type PingResult struct {
		Server   string
		Latency  time.Duration
		IsAlive  bool
	}

	servers := []string{"web-01", "web-02", "api-01", "db-01", "cache-01"}
	results := make(chan PingResult, len(servers))

	// Launch goroutines to ping all servers
	var wg sync.WaitGroup
	for _, server := range servers {
		wg.Go(func() {
			// Simulate pinging with random-ish latency
			latency := time.Duration(len(server)*50) * time.Millisecond
			time.Sleep(latency)

			results <- PingResult{
				Server:  server,
				Latency: latency,
				IsAlive: server != "db-01", // pretend db-01 is down
			}
		})
	}

	// Close channel when all goroutines finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect and print results as they come in
	for r := range results {
		status := "✅ UP  "
		if !r.IsAlive {
			status = "❌ DOWN"
		}
		fmt.Printf("  %s  %-10s  %s\n", status, r.Server, r.Latency)
	}

	// ================================
	// 5. CHANNEL DIRECTION (send-only, receive-only)
	// ================================

	fmt.Println("\n=== Channel Direction ===")

	ch2 := make(chan string, 1)

	// This function can ONLY send to the channel
	sendOnly := func(ch chan<- string, msg string) {
		ch <- msg
		// <-ch  // compile error: cannot receive from send-only channel
	}

	// This function can ONLY receive from the channel
	receiveOnly := func(ch <-chan string) string {
		return <-ch
		// ch <- "nope"  // compile error: cannot send to receive-only channel
	}

	sendOnly(ch2, "type-safe channels!")
	fmt.Println("Received:", receiveOnly(ch2))

	fmt.Println("\n  Channel directions prevent bugs - you can't accidentally")
	fmt.Println("  send when you should receive, or vice versa.")

	// ================================
	// 6. SELECT - LISTENING ON MULTIPLE CHANNELS
	// ================================

	fmt.Println("\n=== Select (multiple channels) ===")

	fast := make(chan string, 1)
	slow := make(chan string, 1)

	go func() {
		time.Sleep(100 * time.Millisecond)
		fast <- "fast result"
	}()

	go func() {
		time.Sleep(500 * time.Millisecond)
		slow <- "slow result"
	}()

	// select picks whichever channel is ready first
	for i := 0; i < 2; i++ {
		select {
		case msg := <-fast:
			fmt.Println("  From fast:", msg)
		case msg := <-slow:
			fmt.Println("  From slow:", msg)
		case <-time.After(1 * time.Second):
			fmt.Println("  Timeout!")
		}
	}

	fmt.Println("\n  select is like a switch for channels - perfect for")
	fmt.Println("  handling timeouts and multiple data sources.")
}
