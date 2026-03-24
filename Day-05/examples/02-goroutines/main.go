package main

import (
	"fmt"
	"sync"
	"time"
)

// simulateWork pretends to do some work (like calling an API)
func simulateWork(task string, duration time.Duration) {
	fmt.Printf("  ⏳ [%s] Starting... (takes %s)\n", task, duration)
	time.Sleep(duration)
	fmt.Printf("  ✅ [%s] Done!\n", task)
}

func main() {
	// ================================
	// 1. WITHOUT GOROUTINES (sequential)
	// ================================

	fmt.Println("=== Without Goroutines (Sequential) ===")
	start := time.Now()

	simulateWork("Check server-1", 500*time.Millisecond)
	simulateWork("Check server-2", 300*time.Millisecond)
	simulateWork("Check server-3", 400*time.Millisecond)

	fmt.Printf("Total time: %s\n", time.Since(start).Round(time.Millisecond))
	// ~1200ms because they run one after another

	// ================================
	// 2. WITH GOROUTINES + WAITGROUP (classic way)
	// ================================

	fmt.Println("\n=== With Goroutines + WaitGroup (Classic) ===")
	start = time.Now()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		simulateWork("Check server-1", 500*time.Millisecond)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		simulateWork("Check server-2", 300*time.Millisecond)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		simulateWork("Check server-3", 400*time.Millisecond)
	}()

	wg.Wait() // block until all 3 are done
	fmt.Printf("Total time: %s\n", time.Since(start).Round(time.Millisecond))
	// ~500ms because they all run at the same time!

	// ================================
	// 3. WITH WAITGROUP.Go (the new clean way)
	// ================================

	fmt.Println("\n=== With WaitGroup.Go (Clean Way) ===")
	start = time.Now()

	var wg2 sync.WaitGroup

	// Much cleaner - no Add(1) or defer Done()!
	wg2.Go(func() {
		simulateWork("Check server-1", 500*time.Millisecond)
	})

	wg2.Go(func() {
		simulateWork("Check server-2", 300*time.Millisecond)
	})

	wg2.Go(func() {
		simulateWork("Check server-3", 400*time.Millisecond)
	})

	wg2.Wait()
	fmt.Printf("Total time: %s\n", time.Since(start).Round(time.Millisecond))

	// ================================
	// 4. LAUNCHING GOROUTINES IN A LOOP
	// ================================

	fmt.Println("\n=== Goroutines in a Loop ===")
	start = time.Now()

	servers := []string{"web-01", "web-02", "api-01", "db-01", "cache-01"}

	var wg3 sync.WaitGroup
	for _, server := range servers {
		wg3.Go(func() {
			simulateWork("Ping "+server, 200*time.Millisecond)
		})
	}

	wg3.Wait()
	fmt.Printf("Pinged %d servers in: %s\n", len(servers), time.Since(start).Round(time.Millisecond))
	// ~200ms for all 5 - they all run at the same time!

	// ================================
	// 5. HOW MANY GOROUTINES CAN YOU RUN?
	// ================================

	fmt.Println("\n=== Goroutine Scaling ===")

	counts := []int{10, 100, 1000, 10000}

	for _, count := range counts {
		start = time.Now()
		var wg4 sync.WaitGroup

		for range count {
			wg4.Go(func() {
				time.Sleep(10 * time.Millisecond) // tiny bit of work
			})
		}

		wg4.Wait()
		fmt.Printf("  %6d goroutines completed in: %s\n", count, time.Since(start).Round(time.Millisecond))
	}

	fmt.Println("\n  Goroutines are super cheap! Each one uses only ~2KB of stack.")
	fmt.Println("  You can easily run thousands of them. Try doing that with OS threads!")
}
