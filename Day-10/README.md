# Day 10 - Capstone Project: Build Your Own Kubectl Plugin

## What are we building today?

This is the capstone - we bring everything together and build a **real kubectl plugin** called `kubectl-status`.

This plugin gives you a single-command overview of your Kubernetes cluster - pods, deployments, and services - all in one clean, color-coded view. No more running three separate commands!

### What makes this a kubectl plugin?

A kubectl plugin is just a binary on your PATH that starts with `kubectl-`. That's it. No special SDK, no registration. If you have a binary called `kubectl-status`, you can run:

```bash
kubectl status
```

Kubernetes automatically finds it. This is one of the simplest and most powerful plugin systems in the DevOps world.

---

## What concepts are we using?

This project uses everything we learned in the course:

| Day | Concept | How we use it here |
|-----|---------|-------------------|
| Day-01 | Variables, fmt, CLI args | Basic program structure |
| Day-02 | Loops, slices, maps, functions | Processing K8s resources |
| Day-03 | Structs, interfaces, JSON | Resource types, JSON output |
| Day-04 | Error handling, packages | Graceful errors, cmd/ package |
| Day-05 | Goroutines, WaitGroup | Fetching pods/deployments/services concurrently |
| Day-06 | String manipulation | Formatting output, parsing |
| Day-07 | HTTP/JSON, environment vars | K8s API communication |
| Day-08 | Cobra, flags, cross-compilation | CLI structure, building the plugin |
| Day-09 | Kubernetes client-go, context | Talking to the cluster |

---

## How kubectl plugins work

### The rules

1. The binary name must start with `kubectl-`
2. It must be on your `PATH`
3. That's it. Seriously.

### Naming

```
Binary name         ->  kubectl command
kubectl-status      ->  kubectl status
kubectl-ns          ->  kubectl ns
kubectl-pod-restart ->  kubectl pod-restart
```

### Installing a plugin

```bash
# Build the binary
go build -o kubectl-status .

# Copy to somewhere on your PATH
cp kubectl-status /usr/local/bin/

# Now you can use it!
kubectl status
kubectl status -n default
kubectl status -o json
```

### Check if kubectl finds it

```bash
kubectl plugin list
```

---

## Project: `kubectl-status`

### Features

- Shows pods, deployments, and services in one view
- Color-coded status (green for healthy, red for problems, yellow for pending)
- Highlights pods in CrashLoopBackOff or with many restarts
- Supports namespace filtering (`-n default`)
- Table output (default) and JSON output (`-o json`)
- Fetches all resources concurrently using goroutines (fast!)
- Graceful error handling when cluster is unreachable
- Works standalone (`./kubectl-status`) or as a plugin (`kubectl status`)

### How to run

```bash
# Run directly (during development)
go run . 
go run . -n default
go run . -o json
go run . --help

# Build the plugin binary
go build -o kubectl-status .

# Install it
cp kubectl-status /usr/local/bin/

# Use as kubectl plugin
kubectl status
kubectl status -n kube-system
kubectl status -o json
```

### Sample output

```
=== Cluster Status ===
Context: minikube

--- Pods (namespace: all) ---

  Name                                Namespace     Status     Ready  Restarts  Age
  ----                                ---------     ------     -----  --------  ---
🟢 nginx-deployment-abc123-xyz        default       Running    1/1    0         2d
🟢 redis-master-0                     default       Running    1/1    0         5d
🔴 api-server-def456-uvw              production    CrashLoop  0/1    42        1h
🟡 worker-ghi789-rst                  production    Pending    0/1    0         5m

--- Deployments (namespace: all) ---

  Name                    Namespace     Ready   UpToDate  Available  Age
  ----                    ---------     -----   --------  ---------  ---
✅ nginx-deployment        default       3/3     3         3          2d
⚠️  api-server              production    0/2     2         0          1h
✅ redis                   default       1/1     1         1          5d

--- Services (namespace: all) ---

  Name                    Namespace     Type          ClusterIP       Ports         Age
  ----                    ---------     ----          ---------       -----         ---
  kubernetes              default       ClusterIP     10.96.0.1       443/TCP       30d
  nginx-service           default       LoadBalancer  10.96.100.50    80:30080/TCP  2d
  redis                   default       ClusterIP     10.96.200.10    6379/TCP      5d

=== Summary ===
Pods:        4 total (2 running, 1 crash-loop, 1 pending)
Deployments: 3 total (2 healthy, 1 unhealthy)
Services:    3 total
```

### Project structure

```
Day-10/
  go.mod
  go.sum
  main.go          <- entry point
  cmd/
    root.go        <- Cobra root command, flags, config
    status.go      <- main status logic (fetch + display)
  kube/
    client.go      <- K8s client setup, kubeconfig loading
    resources.go   <- fetch pods, deployments, services concurrently
    types.go       <- our custom types for clean data
  display/
    table.go       <- table output with colors
    json.go        <- JSON output
    colors.go      <- terminal color helpers
```

---

## Building and Cross-Compiling

```bash
# Build for your current OS
go build -o kubectl-status .

# Build for Linux (for deploying to servers)
GOOS=linux GOARCH=amd64 go build -o kubectl-status-linux-amd64 .

# Build for Linux ARM (Raspberry Pi, AWS Graviton)
GOOS=linux GOARCH=arm64 go build -o kubectl-status-linux-arm64 .

# Build for macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o kubectl-status-darwin-arm64 .

# Build with version info
go build -ldflags "-X kubectl-status/cmd.Version=1.0.0 -X kubectl-status/cmd.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o kubectl-status .
```

---

## Mini Challenges

1. **Easy**: Add a `--pods-only` flag that shows only pods (skip deployments and services)
2. **Medium**: Add a `--watch` flag that refreshes the output every 5 seconds (like `kubectl get pods -w`)
3. **Hard**: Add a `--problems-only` flag that shows only unhealthy resources (crash-looping pods, unavailable deployments)

---

## What We Learned in This Course

Looking back at what we built over 10 days:

1. **Day-01**: Go setup, variables, fmt, CLI basics - built a port checker
2. **Day-02**: Loops, slices, maps, functions - built a log analyzer
3. **Day-03**: Structs, interfaces, JSON - built a config reader
4. **Day-04**: File I/O, error handling, packages - built a log parser
5. **Day-05**: HTTP, goroutines, channels - built a health checker
6. **Day-06**: System commands, regex - built a system info collector
7. **Day-07**: REST APIs, HTTP server - built a GitHub API tool
8. **Day-08**: Cobra, Viper, cross-compilation - built `devctl` swiss knife
9. **Day-09**: Docker SDK, K8s client-go - built `infraview`
10. **Day-10**: Capstone - built a real kubectl plugin

You now have the skills to build real DevOps tools in Go. Go build something cool!
