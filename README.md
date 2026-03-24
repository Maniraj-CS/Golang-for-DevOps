# Golang for DevOps 

## Who is this for?

- DevOps engineers who write Bash/Python scripts and want to level up
- SREs who want to build reliable tooling
- Anyone who wants to automate infra stuff using Go

## What we will NOT do

- We won't cover every Go concept out there
- No web frameworks, no frontend stuff
- No DSA, no competitive programming

## What we WILL do

- Learn Go by building real DevOps tools
- Every episode = one concept + one hands-on project
- By the end, you will be mass confident writing Go for your day-to-day DevOps work

---

## Course Syllabus

### Day-01 - Go Setup, Data Types and Your First CLI Tool

**Concepts:**
- Installing Go 1.26, understanding `go env`
- Go project structure, go modules (`go mod init`)
- Data types - string, int, float64, bool
- Variables - `var`, short declaration `:=`, constants
- fmt package - Println, Printf, Sprintf
- Taking user input with `fmt.Scan` and `os.Args`

**Project:** Build a simple CLI tool that checks if a given port is open or not on localhost

---

### Day-02 - Conditionals, Loops and Functions

**Concepts:**
- if-else, switch case
- for loop (the only loop in Go!), range
- Range over integers (`for i := range 10`)
- Slices and arrays - declaration, append, iterate
- Maps - key-value pairs, checking if key exists
- Functions - multiple return values, named returns
- Pointers basics - what is `&` and `*`, why it matters

**Project:** Build a log file analyzer that reads a server log file, counts ERROR vs WARN vs INFO lines, and prints a summary report

---

### Day-03 - Structs, Interfaces and JSON

**Concepts:**
- Structs - defining, creating, nested structs
- Methods on structs (receiver functions)
- Interfaces - what they are, why they matter
- JSON marshal/unmarshal, struct tags (`json:"name"`, `omitempty`, `omitzero`)
- Reading and writing JSON config files

**Project:** Build a tool that reads a JSON config file (list of servers with IP, role, environment) and prints a nicely formatted table. Use interfaces to support both JSON and YAML config formats

---

### Day-04 - File I/O, Error Handling and Packages

**Concepts:**
- os package - reading, writing, appending files
- bufio for reading files line by line
- Error handling - the Go way (`if err != nil`)
- Custom errors, error wrapping with `fmt.Errorf`
- `errors.As`, `errors.Is`, and `errors.AsType`
- Creating your own packages, importing them
- init() function

**Project:** Build a server log parser that reads large log files, filters by severity level, and writes the output to a new file. Handle all edge cases like file not found, permission denied, etc.

---

### Day-05 - HTTP Calls, Health Checks and Goroutines

**Concepts:**
- net/http package - GET requests, reading response body
- Status codes, timeouts, headers
- Goroutines - what are they, how to launch them
- sync.WaitGroup - waiting for goroutines to finish
- `sync.WaitGroup.Go` method - cleaner way to launch goroutines
- Channels basics - sending and receiving data
- Why concurrency is a superpower for DevOps

**Project:** Build a website health checker that pings a list of URLs concurrently using goroutines and reports which ones are up/down with response time

---

### Day-06 - Running System Commands from Go

**Concepts:**
- os/exec package - running shell commands
- Capturing stdout, stderr
- Setting environment variables from code
- Working with filepath package
- Defer, panic, recover
- String manipulation - strings package, regexp basics

**Project:** Build a system info collector that runs commands like `df -h`, `free -m`, `uptime`, parses the output, and gives a clean JSON/table report

---

### Day-07 - Working with REST APIs

**Concepts:**
- Making POST, PUT, DELETE requests
- Setting headers, auth tokens
- Parsing nested JSON responses
- Building a simple HTTP server with `http.HandleFunc` and `http.ServeMux` patterns
- Serving JSON responses
- Environment variables with `os.Getenv`

**Project:** Build a tool that talks to GitHub API - list your repos, check last commit date, create a new issue from terminal. Also expose a simple `/repos` API endpoint locally

---

### Day-08 - Building a Full CLI App with Cobra

**Concepts:**
- Cobra library - why every big Go tool uses it
- Subcommands, flags (string, bool, int flags)
- Persistent flags vs local flags
- Reading from config file using Viper
- `go tool` directive in go.mod - tracking tool dependencies
- Building and distributing Go binaries (cross-compilation with `GOOS` and `GOARCH`)

**Project:** Build a DevOps swiss knife CLI called `devctl` with subcommands like `devctl health`, `devctl sysinfo`, `devctl github` - combining what we built in earlier episodes into one proper tool

---

### Day-09 - Docker and Kubernetes from Go

**Concepts:**
- Docker SDK for Go - connecting to Docker daemon
- Listing containers, images, pulling images
- Kubernetes client-go basics
- Listing pods, deployments, services
- Context and timeout in Go
- Container-aware `GOMAXPROCS` - how Go behaves inside containers

**Project:** Build a tool that lists running Docker containers with their resource usage, and optionally connects to a K8s cluster to list pods and their status

---

### Day-10 - Capstone Project: Build Your Own Kubectl Plugin (Beginner Friendly!)

**Concepts:** Bringing it ALL together - structs, interfaces, goroutines, JSON, Cobra CLI, error handling, Kubernetes client-go, color-coded terminal output, cross-compilation

**Project:** Build a real kubectl plugin called `kubectl-status` that:
- Shows pods, deployments, and services in one clean, color-coded view (🟢🟡🔴)
- Highlights unhealthy pods, CrashLoopBackOff, and restart loops
- Fetches all resources concurrently using goroutines (fast!)
- Supports namespace filtering with `-n` flag
- Outputs as table (default) or JSON (`-o json`)
- Works as a standalone binary AND as a kubectl plugin (`kubectl status`)
- Installable by just copying the binary to your PATH — that's how kubectl plugins work!
- Uses concepts from every single day of this course

---
