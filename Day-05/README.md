# Day 05 - HTTP Calls, Health Checks and Goroutines

## What are we doing today?

Today we will:
1. Learn how to make HTTP requests using `net/http`
2. Handle status codes, timeouts, and headers
3. Understand goroutines - Go's lightweight threads
4. Use `sync.WaitGroup` to wait for goroutines to finish
5. Use `sync.WaitGroup.Go` - the cleaner way to launch goroutines
6. Learn channels - how goroutines talk to each other
7. Build a concurrent website health checker

By the end of this episode, you will have a tool that pings a list of URLs concurrently (at the same time!) and reports which ones are up, which are down, and how fast each one responds. This is a tool every DevOps engineer needs.

---

## 1. Making HTTP Requests

### Simple GET request

```go
resp, err := http.Get("https://httpbin.org/get")
if err != nil {
    log.Fatal("Request failed:", err)
}
defer resp.Body.Close() // ALWAYS close the body!

fmt.Println("Status:", resp.Status)           // "200 OK"
fmt.Println("Status Code:", resp.StatusCode)  // 200
```

### Reading the response body

```go
body, err := io.ReadAll(resp.Body)
if err != nil {
    log.Fatal("Cannot read body:", err)
}
fmt.Println("Body:", string(body))
```

### Setting a timeout

By default, `http.Get` has no timeout - it will wait forever! In DevOps, you always want a timeout.

```go
client := &http.Client{
    Timeout: 5 * time.Second, // give up after 5 seconds
}

resp, err := client.Get("https://httpbin.org/delay/10")
if err != nil {
    fmt.Println("Timed out or failed:", err)
    return
}
defer resp.Body.Close()
```

### Custom request with headers

```go
req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
if err != nil {
    log.Fatal(err)
}

// Add headers
req.Header.Set("Authorization", "Bearer your-token-here")
req.Header.Set("Accept", "application/json")

client := &http.Client{Timeout: 10 * time.Second}
resp, err := client.Do(req)
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()
```

### Checking status codes

```go
switch {
case resp.StatusCode >= 200 && resp.StatusCode < 300:
    fmt.Println("✅ Success")
case resp.StatusCode >= 300 && resp.StatusCode < 400:
    fmt.Println("↗️  Redirect")
case resp.StatusCode >= 400 && resp.StatusCode < 500:
    fmt.Println("❌ Client error")
case resp.StatusCode >= 500:
    fmt.Println("💥 Server error")
}
```

---

## 2. Goroutines

A goroutine is a lightweight thread managed by Go. It's like saying "hey Go, run this function in the background while I do other stuff".

### Launching a goroutine

```go
// Normal function call - runs and blocks until done
doWork()

// Goroutine - runs in the background, does NOT block
go doWork()
```

That's it! Just add `go` before the function call. But there is a catch - if `main()` finishes before the goroutine, the program exits and the goroutine is killed.

```go
func main() {
    go fmt.Println("Hello from goroutine!") // might not print!
    // main exits immediately, goroutine has no time to run
}
```

This is where WaitGroup comes in.

---

## 3. sync.WaitGroup

A WaitGroup waits for a collection of goroutines to finish.

### The classic way

```go
var wg sync.WaitGroup

wg.Add(1) // tell WaitGroup: 1 goroutine is starting
go func() {
    defer wg.Done() // tell WaitGroup: this goroutine is done
    fmt.Println("Task 1 done")
}()

wg.Add(1)
go func() {
    defer wg.Done()
    fmt.Println("Task 2 done")
}()

wg.Wait() // block until all goroutines call Done()
fmt.Println("All tasks complete!")
```

### The new way - WaitGroup.Go (cleaner!)

In recent Go versions, `sync.WaitGroup` has a `.Go()` method that combines `Add(1)` + launching the goroutine + `Done()` in one call:

```go
var wg sync.WaitGroup

wg.Go(func() {
    fmt.Println("Task 1 done")
})

wg.Go(func() {
    fmt.Println("Task 2 done")
})

wg.Wait()
fmt.Println("All tasks complete!")
```

Much cleaner! No more `Add(1)` and `defer wg.Done()` boilerplate.

---

## 4. Channels

Channels are how goroutines communicate with each other. Think of a channel as a pipe - one goroutine sends data in, another goroutine receives it.

### Basic channel

```go
ch := make(chan string) // create a channel that carries strings

// Send data into channel (in a goroutine)
go func() {
    ch <- "hello from goroutine!"
}()

// Receive data from channel (blocks until data arrives)
msg := <-ch
fmt.Println(msg) // "hello from goroutine!"
```

### Buffered channels

```go
// Buffered channel - can hold 3 values without blocking
ch := make(chan string, 3)

ch <- "first"
ch <- "second"
ch <- "third"
// ch <- "fourth" // this would block because buffer is full!

fmt.Println(<-ch) // "first"
fmt.Println(<-ch) // "second"
fmt.Println(<-ch) // "third"
```

### Channel with range

```go
ch := make(chan int, 5)

// Send values
go func() {
    for i := 1; i <= 5; i++ {
        ch <- i
    }
    close(ch) // close when done sending - important!
}()

// Receive until channel is closed
for val := range ch {
    fmt.Println(val)
}
```

### Why channels matter for DevOps

- Health check 10 servers concurrently, collect results through a channel
- Run 50 API calls in parallel, gather all responses
- Process log files in parallel, merge results

---

## 5. Putting It Together - Concurrent Pattern

This is the pattern you will use all the time in DevOps:

```go
// 1. Create a channel for results
results := make(chan Result, len(urls))

// 2. Launch goroutines for each task
var wg sync.WaitGroup
for _, url := range urls {
    wg.Go(func() {
        result := checkURL(url)
        results <- result
    })
}

// 3. Close channel when all goroutines are done
go func() {
    wg.Wait()
    close(results)
}()

// 4. Collect results
for result := range results {
    fmt.Println(result)
}
```

---

## 6. Project: Website Health Checker

We build a tool that:
- Takes a list of URLs (from CLI args or a file)
- Pings all of them concurrently using goroutines
- Reports status (up/down), HTTP status code, and response time
- Has a configurable timeout
- Prints a nice summary table

### How to run

```bash
# Check specific URLs
go run main.go https://google.com https://github.com https://httpbin.org/status/500

# Check URLs from a file (one URL per line)
go run main.go --file urls.txt

# Set custom timeout (default is 5 seconds)
go run main.go --timeout 10 https://google.com https://github.com

# Build and run
go build -o health-checker
./health-checker https://google.com https://github.com
```

### Sample output

```
=== Website Health Checker ===
Timeout: 5s
Checking 5 URLs concurrently...

Status  Code  Time      URL
------  ----  ----      ---
✅ UP    200    245ms   https://google.com
✅ UP    200    312ms   https://github.com
✅ UP    200    189ms   https://go.dev
❌ DOWN  500    156ms   https://httpbin.org/status/500
❌ DOWN  ---   5001ms   https://httpbin.org/delay/10 (timeout)

Summary:
  Total  : 5
  Up     : 3
  Down   : 2
  Avg    : 981ms
```

---

## Mini Challenges

1. **Easy**: Add a `--verbose` flag that shows response headers for each URL
2. **Medium**: Add retry logic - if a URL fails, retry it 2 more times before marking it as down
3. **Hard**: Add a `--watch` flag that runs the health check every N seconds and shows a live dashboard in the terminal

---

## What We Learned Today

- Making HTTP requests with `net/http` package
- Setting timeouts (critical for DevOps tools!)
- Custom requests with headers
- Goroutines - lightweight concurrent execution
- `sync.WaitGroup` - waiting for goroutines to finish
- `sync.WaitGroup.Go` - the cleaner way to launch goroutines
- Channels - how goroutines send data to each other
- Buffered vs unbuffered channels
- The concurrent fan-out/fan-in pattern
- Building a real health checker tool

---

Next up: **Day 06 - Running System Commands from Go** where we will build a system info collector that runs shell commands and parses their output!
