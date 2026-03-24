# Day 04 - File I/O, Error Handling and Packages

## What are we doing today?

Today we will:
1. Learn proper file reading and writing using `os` and `bufio`
2. Master error handling - the Go way
3. Create custom errors and use error wrapping
4. Use `errors.Is`, `errors.As`, and `errors.AsType`
5. Create our own packages and import them
6. Understand the `init()` function
7. Build a log parser that filters logs by severity and writes to a new file

By the end of this episode, you will have a real log parser tool that reads large log files, filters by level (INFO/WARN/ERROR), and writes filtered output to a new file. You will handle every edge case properly - file not found, permission denied, bad format etc.

---

## 1. File I/O with os Package

### Reading a file (entire content at once)

```go
// os.ReadFile - reads entire file into memory
// Good for small config files, bad for huge log files
data, err := os.ReadFile("config.json")
if err != nil {
    log.Fatal(err)
}
fmt.Println(string(data))
```

### Writing a file

```go
// os.WriteFile - creates file and writes all data at once
// 0644 is the file permission (owner can read/write, others can read)
content := []byte("Hello from Go!\n")
err := os.WriteFile("output.txt", content, 0644)
if err != nil {
    log.Fatal(err)
}
```

### Appending to a file

```go
// os.OpenFile with flags - more control over how file is opened
file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
if err != nil {
    log.Fatal(err)
}
defer file.Close()

file.WriteString("2024-01-15 10:30:00 INFO Server started\n")
```

File open flags:
- `os.O_RDONLY` - read only
- `os.O_WRONLY` - write only
- `os.O_RDWR` - read and write
- `os.O_CREATE` - create if not exists
- `os.O_APPEND` - append to end
- `os.O_TRUNC` - empty the file first

### Reading a file line by line with bufio

```go
// bufio.Scanner - perfect for reading log files line by line
// Does NOT load the entire file into memory!
file, err := os.Open("server.log")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

scanner := bufio.NewScanner(file)
lineNum := 0

for scanner.Scan() {
    lineNum++
    line := scanner.Text()
    fmt.Printf("%d: %s\n", lineNum, line)
}

if err := scanner.Err(); err != nil {
    log.Fatal("error reading file:", err)
}
```

### Checking if file or directory exists

```go
// os.Stat gives you file info
info, err := os.Stat("config.json")
if err != nil {
    if os.IsNotExist(err) {
        fmt.Println("File does not exist!")
    } else {
        fmt.Println("Error checking file:", err)
    }
} else {
    fmt.Println("File exists! Size:", info.Size(), "bytes")
}
```

---

## 2. Error Handling - The Go Way

Go does not have try-catch-finally like Java or Python. Instead, functions return an error as the last return value, and you check it immediately.

### The basic pattern

```go
file, err := os.Open("config.json")
if err != nil {
    // handle the error
    fmt.Println("Cannot open file:", err)
    return
}
// use file...
```

This might feel repetitive at first, but it makes error handling explicit. You always know exactly where an error can happen. No hidden exceptions.

### Error wrapping with fmt.Errorf

When you get an error from a function, you can wrap it with extra context using `%w`:

```go
func readConfig(path string) ([]byte, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("readConfig: cannot read %s: %w", path, err)
    }
    return data, nil
}
```

The `%w` verb wraps the original error inside your new error. This keeps the error chain intact - you can still check for the original error later.

### Custom errors

```go
// Simple custom error type
type FileError struct {
    Path    string
    Message string
}

func (e *FileError) Error() string {
    return fmt.Sprintf("%s: %s", e.Path, e.Message)
}

// Using it
func validateFile(path string) error {
    info, err := os.Stat(path)
    if err != nil {
        return &FileError{Path: path, Message: "file not found"}
    }
    if info.Size() == 0 {
        return &FileError{Path: path, Message: "file is empty"}
    }
    return nil
}
```

### errors.Is - checking for a specific error

```go
import "errors"

_, err := os.Open("missing.txt")
if errors.Is(err, os.ErrNotExist) {
    fmt.Println("File not found!")
}

// Works even through wrapped errors
wrapped := fmt.Errorf("config failed: %w", os.ErrNotExist)
if errors.Is(wrapped, os.ErrNotExist) {
    fmt.Println("Still finds it through the wrapping!")
}
```

### errors.As - extracting a specific error type

```go
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    fmt.Println("Path:", pathErr.Path)
    fmt.Println("Op:", pathErr.Op)
}
```

### errors.AsType - the cleaner way (new in recent Go versions)

```go
// errors.AsType is a generic version of errors.As
// Returns the error of that type and a bool, no need for a separate variable
if pathErr, ok := errors.AsType[*os.PathError](err); ok {
    fmt.Println("Path:", pathErr.Path)
    fmt.Println("Op:", pathErr.Op)
}
```

This is much cleaner than `errors.As` because you don't need to declare a variable first.

---

## 3. Creating Your Own Packages

So far we have been writing everything in `package main` in one file. In real projects, you split your code into packages.

### Project structure for this project

```
Day-04/
├── go.mod
├── main.go              <- main package (entry point)
├── parser/
│   └── parser.go        <- parser package (log parsing logic)
├── errors/
│   └── errors.go        <- our custom error types
├── examples/
│   ├── 01-file-io/main.go
│   ├── 02-error-handling/main.go
│   └── 03-packages/
│       ├── main.go
│       └── utils/
│           └── utils.go
└── sample.log
```

### How packages work

```go
// File: utils/utils.go
package utils

import "strings"

// FormatLevel makes the log level consistent
// Exported function - starts with UPPERCASE
func FormatLevel(level string) string {
    return strings.ToUpper(strings.TrimSpace(level))
}

// unexportedHelper - starts with lowercase, only visible inside this package
func unexportedHelper() string {
    return "you can't call me from outside"
}
```

```go
// File: main.go
package main

import (
    "fmt"
    "myproject/utils"  // import your package
)

func main() {
    level := utils.FormatLevel("  error  ")
    fmt.Println(level) // ERROR
}
```

Rules:
- Folder name = package name (convention)
- Exported names start with **Uppercase** (accessible from other packages)
- Unexported names start with **lowercase** (private to the package)
- Import using `modulename/packagename`

### init() function

The `init()` function runs automatically before `main()`. It's useful for setup tasks.

```go
package main

import "fmt"

func init() {
    fmt.Println("init runs first!")
    // Good for: loading config, setting defaults, validating environment
}

func main() {
    fmt.Println("main runs second!")
}
```

- Every package can have an `init()` function
- It runs before anything else in that package
- You cannot call `init()` manually
- Don't put heavy logic here - keep it simple

---

## 4. Project: Log Parser Tool

We build a proper log parser that:
- Reads a log file line by line (memory efficient)
- Filters by severity level (INFO, WARN, ERROR)
- Writes filtered lines to a new output file
- Has proper custom errors with wrapping
- Uses separate packages for clean code
- Handles all edge cases gracefully

### How to run

```bash
# Parse all lines (default)
go run main.go sample.log

# Filter only ERROR lines and write to file
go run main.go sample.log --level ERROR --output errors.txt

# Filter WARN and above
go run main.go sample.log --level WARN

# Build and run
go build -o log-parser
./log-parser sample.log --level ERROR --output errors.txt
```

### Sample output

```
=== Log Parser ===
File: sample.log

Parsing... done! (150 lines processed)

Summary:
  INFO   : 95
  WARN   : 35
  ERROR  : 20

Filter: ERROR only

Filtered lines written to: errors.txt (20 lines)
```

---

## Mini Challenges

1. **Easy**: Add a `--count` flag that only prints the count of each level without showing individual lines
2. **Medium**: Add a `--after` and `--before` flag to filter logs by timestamp (e.g., `--after "2024-01-15 10:00:00"`)
3. **Hard**: Add support for reading gzipped log files (`.log.gz`) using the `compress/gzip` package

---

## What We Learned Today

- File reading - `os.ReadFile` for small files, `bufio.Scanner` for large files
- File writing - `os.WriteFile`, `os.OpenFile` with flags
- Appending to files with `os.O_APPEND`
- Error handling - `if err != nil` pattern
- Error wrapping with `fmt.Errorf` and `%w`
- Custom error types
- `errors.Is` for checking specific errors
- `errors.As` and `errors.AsType` for extracting error types
- Creating your own packages with exported/unexported names
- `init()` function

---

Next up: **Day 05 - HTTP Calls, Health Checks and Goroutines** where we will build a concurrent website health checker!
