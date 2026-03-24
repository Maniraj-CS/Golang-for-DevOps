# Day 06 - Running System Commands from Go

## What are we doing today?

Today we will:
1. Run shell commands from Go using `os/exec`
2. Capture stdout and stderr separately
3. Work with environment variables
4. Use the `filepath` package for safe path handling
5. Understand `defer`, `panic`, and `recover`
6. Use `strings` and `regexp` packages for text parsing
7. Build a system info collector that runs system commands and gives a clean report

By the end of this episode, you will have a tool that collects disk usage, memory, CPU, hostname, uptime, and network info from your machine by running system commands, parsing the output, and presenting it as a clean JSON or table report. This is exactly what monitoring agents like node_exporter do.

---

## 1. Running Commands with os/exec

### Basic command execution

```go
import "os/exec"

// Run a command and get the output
out, err := exec.Command("hostname").Output()
if err != nil {
    log.Fatal(err)
}
fmt.Println("Hostname:", string(out))
```

### Command with arguments

```go
// Each argument is a separate string - NOT one big string
out, err := exec.Command("ls", "-la", "/tmp").Output()
```

Common mistake:
```go
// WRONG - this tries to run a binary called "ls -la /tmp"
exec.Command("ls -la /tmp")

// CORRECT - command and args are separate
exec.Command("ls", "-la", "/tmp")
```

### Capturing stdout and stderr separately

```go
cmd := exec.Command("ls", "/nonexistent")

var stdout, stderr bytes.Buffer
cmd.Stdout = &stdout
cmd.Stderr = &stderr

err := cmd.Run()
if err != nil {
    fmt.Println("Error:", err)
    fmt.Println("Stderr:", stderr.String())
} else {
    fmt.Println("Output:", stdout.String())
}
```

### Running a shell command (with pipes, etc.)

```go
// If you need shell features like pipes, use sh -c
out, err := exec.Command("sh", "-c", "ps aux | grep go | wc -l").Output()
```

### Getting exit code

```go
cmd := exec.Command("ls", "/nonexistent")
err := cmd.Run()
if err != nil {
    if exitErr, ok := err.(*exec.ExitError); ok {
        fmt.Println("Exit code:", exitErr.ExitCode())
    }
}
```

---

## 2. Environment Variables

### Reading env vars

```go
// Get a specific env var
home := os.Getenv("HOME")
fmt.Println("Home:", home)

// Check if an env var exists
val, exists := os.LookupEnv("MY_SECRET")
if !exists {
    fmt.Println("MY_SECRET is not set")
}
```

### Setting env vars for a command

```go
cmd := exec.Command("sh", "-c", "echo $MY_APP_ENV")
cmd.Env = append(os.Environ(), "MY_APP_ENV=production")
out, _ := cmd.Output()
fmt.Println(string(out)) // production
```

### Setting env vars in current process

```go
os.Setenv("MY_VAR", "hello")
fmt.Println(os.Getenv("MY_VAR")) // hello
os.Unsetenv("MY_VAR")
```

---

## 3. filepath Package

The `filepath` package handles file paths safely across different operating systems.

```go
import "path/filepath"

// Join paths safely (handles slashes)
path := filepath.Join("home", "user", "config.json")
// Linux/Mac: "home/user/config.json"
// Windows:   "home\\user\\config.json"

// Get file extension
ext := filepath.Ext("server.log.gz")  // ".gz"

// Get base name (just the filename)
base := filepath.Base("/var/log/server.log")  // "server.log"

// Get directory
dir := filepath.Dir("/var/log/server.log")  // "/var/log"

// Walk a directory tree
filepath.Walk("/tmp", func(path string, info os.FileInfo, err error) error {
    if err != nil {
        return err
    }
    fmt.Println(path)
    return nil
})
```

---

## 4. Defer, Panic, Recover

### defer - run something when function exits

```go
func readFile(path string) {
    file, err := os.Open(path)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close() // runs when readFile() returns, no matter what

    // ... use file ...
    // file.Close() is called automatically, even if there's an error
}
```

`defer` is LIFO (last in, first out):
```go
defer fmt.Println("first")
defer fmt.Println("second")
defer fmt.Println("third")
// Output: third, second, first
```

### panic - something really bad happened

```go
// panic stops the program immediately
panic("something terrible happened!")

// Use it for truly unrecoverable situations
// For normal errors, use error returns
```

### recover - catch a panic (like catch in try-catch)

```go
func safeRun(name string, fn func()) {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Recovered from panic in %s: %v\n", name, r)
        }
    }()
    fn()
}

safeRun("risky task", func() {
    panic("oops!")
})
fmt.Println("Program continues!") // this still runs
```

Rule of thumb:
- Use `error` returns for expected failures (file not found, network error)
- Use `panic` only for bugs/impossible situations
- Use `recover` in top-level handlers to prevent crashes

---

## 5. Strings and Regexp

### strings package (you'll use this all the time)

```go
import "strings"

strings.Contains("ERROR: disk full", "ERROR")       // true
strings.HasPrefix("/api/v1/users", "/api")           // true
strings.HasSuffix("server.log.gz", ".gz")            // true
strings.Split("key=value", "=")                      // ["key", "value"]
strings.Join([]string{"a", "b", "c"}, ", ")          // "a, b, c"
strings.TrimSpace("  hello  ")                       // "hello"
strings.ReplaceAll("foo-bar-baz", "-", "_")          // "foo_bar_baz"
strings.ToUpper("error")                              // "ERROR"
strings.Fields("  hello   world  ")                   // ["hello", "world"] (splits by whitespace)
strings.Repeat("-", 40)                               // "----------------------------------------"
```

### regexp package (when strings package is not enough)

```go
import "regexp"

// Compile a pattern
re := regexp.MustCompile(`\d+\.\d+%`)  // matches "78.5%" etc.

// Find first match
match := re.FindString("CPU: 78.5% used")  // "78.5%"

// Find all matches
matches := re.FindAllString("CPU: 78.5%, Mem: 64.2%", -1)  // ["78.5%", "64.2%"]

// Named capture groups
re2 := regexp.MustCompile(`(?P<key>\w+)=(?P<value>\S+)`)
match2 := re2.FindStringSubmatch("port=8080")
// match2[0] = "port=8080", match2[1] = "port", match2[2] = "8080"
```

---

## 6. Project: System Info Collector

We build a tool that:
- Collects system info by running shell commands
- Parses the output using strings and regexp
- Displays as a clean table or JSON report
- Works on macOS and Linux
- Handles errors gracefully when a command is not available

### How to run

```bash
# Show system info as a table (default)
go run main.go

# Show as JSON
go run main.go --json

# Show only specific sections
go run main.go --section disk
go run main.go --section memory
go run main.go --section all

# Build and run
go build -o sysinfo
./sysinfo --json
```

### Sample output

```
=== System Info Collector ===
Collected at: 2024-01-15 10:30:00

--- Host ---
  Hostname     : web-server-01
  OS           : darwin (arm64)
  Uptime       : up 5 days, 12:30

--- CPU ---
  Cores        : 8
  Load Avg     : 1.25, 0.98, 0.87

--- Memory ---
  Total        : 16.00 GB
  Used         : 10.24 GB
  Available    : 5.76 GB
  Usage        : 64.0%

--- Disk ---
  Filesystem   Size    Used    Avail   Use%   Mounted
  /dev/disk1   228G    156G    72G     68%    /
  /dev/disk2   500G    350G    150G    70%    /data

--- Network ---
  Interface    IP Address
  en0          192.168.1.100
  lo0          127.0.0.1
```

---

## Mini Challenges

1. **Easy**: Add a `--hostname` flag that shows only the hostname (useful for scripts)
2. **Medium**: Add a `--watch` flag that refreshes system info every N seconds (use goroutines + time.Ticker)
3. **Hard**: Add `--alert` flag that checks if any metric crosses a threshold (disk > 80%, memory > 90%) and prints warnings

---

## What We Learned Today

- Running shell commands with `os/exec`
- Capturing stdout and stderr separately
- Command arguments (each one separate, not one big string!)
- Environment variables - reading, setting, passing to commands
- `filepath` package for safe path handling
- `defer` - cleanup that always runs
- `panic` and `recover` - emergency error handling
- `strings` package - Split, Fields, Contains, TrimSpace, etc.
- `regexp` package - patterns, capture groups
- Parsing command output into structured data

---

Next up: **Day 07 - Working with REST APIs** where we will talk to GitHub API, build a simple HTTP server, and more!
