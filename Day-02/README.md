# Day 02 - Conditionals, Loops and Functions

## What are we doing today?

Today we will:
1. Learn if-else and switch
2. Learn the `for` loop (Go has only ONE loop!)
3. Work with slices, arrays and maps
4. Write functions with multiple return values
5. Understand pointers basics
6. Build a log file analyzer tool

By the end of this episode, you will have a tool that reads a server log file, counts how many ERROR, WARN, and INFO lines are there, and prints a nice summary. Every DevOps engineer has done this with `grep` and `awk` - today we do it in Go.

---

## 1. If-Else

```go
port := 8080

if port == 80 {
    fmt.Println("HTTP")
} else if port == 443 {
    fmt.Println("HTTPS")
} else {
    fmt.Printf("Custom port: %d\n", port)
}
```

Go has a nice trick - you can declare a variable inside the if statement:

```go
if status := checkServer(); status == "down" {
    fmt.Println("Server is down! Alert!")
}
// status is NOT available here - it only lives inside the if block
```

This is very common in Go. You will see this pattern everywhere, especially with error handling:

```go
if err := doSomething(); err != nil {
    fmt.Println("Error:", err)
}
```

---

## 2. Switch Case

```go
env := "production"

switch env {
case "development":
    fmt.Println("Dev mode - debug logging ON")
case "staging":
    fmt.Println("Staging mode - testing")
case "production":
    fmt.Println("Prod mode - be careful!")
default:
    fmt.Println("Unknown environment")
}
```

Go's switch is cleaner than other languages - no need for `break`, it stops automatically after a match.

You can also do switch without a condition (like multiple if-else):

```go
hour := 14

switch {
case hour < 12:
    fmt.Println("Good morning!")
case hour < 17:
    fmt.Println("Good afternoon!")
default:
    fmt.Println("Good evening!")
}
```

---

## 3. For Loop

Go has only ONE looping keyword: `for`. But it can do everything.

```go
// Classic for loop
for i := 0; i < 5; i++ {
    fmt.Println(i)
}

// While loop style
count := 0
for count < 5 {
    fmt.Println(count)
    count++
}

// Infinite loop (useful for servers, watchers)
for {
    fmt.Println("Running forever...")
    break // remove this to actually run forever
}

// Range over integers
for i := range 5 {
    fmt.Println(i) // 0, 1, 2, 3, 4
}
```

---

## 4. Arrays and Slices

### Arrays (fixed size - rarely used)

```go
var ports [3]int
ports[0] = 80
ports[1] = 443
ports[2] = 8080
```

### Slices (dynamic size - used everywhere)

```go
// Create a slice
servers := []string{"web-01", "web-02", "db-01"}

// Add to slice
servers = append(servers, "cache-01")

// Length
fmt.Println("Total servers:", len(servers))

// Access by index
fmt.Println("First server:", servers[0])

// Loop over slice
for index, server := range servers {
    fmt.Printf("%d: %s\n", index, server)
}

// If you don't need the index, use _
for _, server := range servers {
    fmt.Println(server)
}
```

Think of slices like Python lists. You will use them all the time.

---

## 5. Maps

Maps are key-value pairs. Like Python dictionaries or JS objects.

```go
// Create a map
serverStatus := map[string]string{
    "web-01":   "running",
    "web-02":   "stopped",
    "db-01":    "running",
    "cache-01": "running",
}

// Access a value
fmt.Println("web-01:", serverStatus["web-01"])

// Add a new key
serverStatus["lb-01"] = "running"

// Delete a key
delete(serverStatus, "web-02")

// Check if key exists
status, exists := serverStatus["web-03"]
if exists {
    fmt.Println("web-03:", status)
} else {
    fmt.Println("web-03 not found")
}

// Loop over map
for server, status := range serverStatus {
    fmt.Printf("%s -> %s\n", server, status)
}
```

---

## 6. Functions

### Basic function

```go
func greet(name string) {
    fmt.Printf("Hello %s!\n", name)
}
```

### Function with return value

```go
func add(a int, b int) int {
    return a + b
}
```

### Multiple return values (very Go thing)

```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("cannot divide by zero")
    }
    return a / b, nil
}

// Usage
result, err := divide(10, 3)
if err != nil {
    fmt.Println("Error:", err)
} else {
    fmt.Println("Result:", result)
}
```

This pattern of returning `(result, error)` is the bread and butter of Go. You will see it in every Go project.

### Named return values

```go
func getServerInfo() (name string, port int) {
    name = "web-server-01"
    port = 8080
    return // naked return - returns the named values
}
```

---

## 7. Pointers Basics

A pointer holds the memory address of a variable. That's it.

```go
name := "Abhishek"
ptr := &name  // & gives the address

fmt.Println(name)   // Abhishek
fmt.Println(ptr)    // 0xc000010230 (some memory address)
fmt.Println(*ptr)   // Abhishek (* gives the value at that address)
```

Why do we need pointers? Two main reasons:

1. **To modify a variable inside a function:**

```go
func changePort(p *int) {
    *p = 9090
}

port := 8080
changePort(&port)
fmt.Println(port) // 9090
```

2. **To avoid copying large data** - passing a pointer is cheaper than copying a whole struct.

Don't overthink pointers right now. Just remember:
- `&` = give me the address
- `*` = give me the value at this address

We will use them more in later episodes.

---

## 8. Project: Log File Analyzer

We will build a tool that:
- Reads a server log file
- Counts how many ERROR, WARN, INFO lines are there
- Shows the summary

This is something every DevOps engineer does daily - just with different tools.

### How to run

First, check out the sample log file `sample.log` in this folder.

```bash
# Run the analyzer on the sample log
go run main.go sample.log

# Build and run
go build -o log-analyzer
./log-analyzer sample.log
```

### Sample output

```
=== Log File Analysis ===
File: sample.log
Total lines: 20

Level       Count
-----       -----
INFO        10
ERROR       6
WARN        4

Top 3 ERROR messages:
  1. Connection refused to database on port 5432
  2. Disk space critical on /dev/sda1
  3. Failed to authenticate user admin
```

---

## Mini Challenges

1. **Easy**: Add a count for DEBUG level logs too
2. **Medium**: Accept a filter flag - show only ERROR lines when user passes `--level ERROR`
3. **Hard**: Show a time-based breakdown - how many errors per hour

---

## What We Learned Today

- if-else with inline variable declaration
- switch case (no break needed in Go)
- for loop - the only loop you need
- range over slices, maps, and integers
- Slices (dynamic arrays) and maps (key-value pairs)
- Functions with multiple return values
- Pointers basics - `&` and `*`
- Built a real log file analyzer

---

Next up: **Day 03 - Structs, Interfaces and JSON** where we will build a config file reader!
