# Day 09 - Docker and Kubernetes from Go

## What are we doing today?

Today we will:
1. Use the Docker SDK to talk to the Docker daemon from Go
2. List containers, images, and pull new images
3. Use the Kubernetes client-go to talk to a K8s cluster
4. List pods, deployments, and services
5. Understand context and timeout in Go
6. Learn about container-aware `GOMAXPROCS`
7. Build a tool that shows Docker containers + K8s pods in one clean view

By the end of this episode, you will have a tool called `infraview` that connects to your local Docker daemon and optionally to a Kubernetes cluster, and shows you what's running in a nice table or JSON output. This is the kind of tool that SREs and DevOps engineers use every day.

**Note:** You need Docker installed to run the Docker examples. K8s examples need a running cluster (minikube, kind, or a real cluster). If you don't have them, the tool handles it gracefully and just shows an error message.

---

## 1. Context and Timeout

Before we talk to Docker and K8s, we need to understand `context.Context`. It's how Go handles timeouts and cancellation.

### Why context?

When you make an API call to Docker or K8s, what if the server is slow? You don't want to wait forever. Context lets you say "give up after 10 seconds".

### Basic usage

```go
import (
    "context"
    "time"
)

// Create a context with 5 second timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel() // always call cancel!

// Pass ctx to any function that accepts it
result, err := someAPICall(ctx)
```

### Rules of context
1. Always pass context as the **first parameter**
2. Always call `cancel()` with defer
3. Never store context in a struct - pass it to functions

```go
// GOOD
func getContainers(ctx context.Context) ([]Container, error) { ... }

// BAD - don't store context
type Client struct {
    ctx context.Context  // don't do this!
}
```

### Checking if context timed out

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

result, err := slowOperation(ctx)
if err != nil {
    if ctx.Err() == context.DeadlineExceeded {
        fmt.Println("Operation timed out!")
    } else {
        fmt.Println("Other error:", err)
    }
}
```

---

## 2. Docker SDK for Go

### Installing the SDK

```bash
go get github.com/docker/docker/client
go get github.com/docker/docker/api/types
```

### Connecting to Docker

```go
import (
    "github.com/docker/docker/client"
)

// Connect to local Docker daemon
cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
if err != nil {
    log.Fatal("Cannot connect to Docker:", err)
}
defer cli.Close()
```

### Listing containers

```go
containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
if err != nil {
    log.Fatal(err)
}

for _, c := range containers {
    fmt.Printf("ID: %s  Image: %s  Status: %s\n",
        c.ID[:12], c.Image, c.Status)
}
```

### Listing images

```go
images, err := cli.ImageList(ctx, image.ListOptions{})
for _, img := range images {
    // img.RepoTags is a slice like ["nginx:latest", "nginx:1.25"]
    tags := strings.Join(img.RepoTags, ", ")
    size := float64(img.Size) / 1024 / 1024 // convert to MB
    fmt.Printf("%-40s  %.0f MB\n", tags, size)
}
```

### Pulling an image

```go
reader, err := cli.ImagePull(ctx, "nginx:latest", image.PullOptions{})
if err != nil {
    log.Fatal(err)
}
defer reader.Close()

// Print pull progress
io.Copy(os.Stdout, reader)
```

### Getting container stats (CPU, memory)

```go
stats, err := cli.ContainerStats(ctx, containerID, false) // false = one snapshot
if err != nil {
    log.Fatal(err)
}
defer stats.Body.Close()

var statsJSON container.StatsResponse
json.NewDecoder(stats.Body).Decode(&statsJSON)

// Calculate CPU percentage
cpuDelta := float64(statsJSON.CPUStats.CPUUsage.TotalUsage - statsJSON.PreCPUStats.CPUUsage.TotalUsage)
systemDelta := float64(statsJSON.CPUStats.SystemUsage - statsJSON.PreCPUStats.SystemUsage)
cpuPercent := (cpuDelta / systemDelta) * float64(len(statsJSON.CPUStats.CPUUsage.PercpuUsage)) * 100.0

// Memory usage
memUsage := float64(statsJSON.MemoryStats.Usage) / 1024 / 1024 // MB
memLimit := float64(statsJSON.MemoryStats.Limit) / 1024 / 1024 // MB
```

---

## 3. Kubernetes client-go

### Installing client-go

```bash
go get k8s.io/client-go@latest
go get k8s.io/apimachinery@latest
```

### Connecting to a cluster

```go
import (
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/util/homedir"
)

// Load kubeconfig from default location (~/.kube/config)
home := homedir.HomeDir()
kubeconfig := filepath.Join(home, ".kube", "config")

config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
if err != nil {
    log.Fatal("Cannot load kubeconfig:", err)
}

clientset, err := kubernetes.NewForConfig(config)
if err != nil {
    log.Fatal("Cannot create K8s client:", err)
}
```

### Listing pods

```go
pods, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
// "" means all namespaces. Use "default" for a specific namespace.

for _, pod := range pods.Items {
    fmt.Printf("%-20s %-15s %-10s %s\n",
        pod.Name,
        pod.Namespace,
        pod.Status.Phase,
        pod.Status.StartTime)
}
```

### Listing deployments

```go
deployments, err := clientset.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})

for _, d := range deployments.Items {
    fmt.Printf("%-30s %-15s %d/%d replicas\n",
        d.Name,
        d.Namespace,
        d.Status.ReadyReplicas,
        *d.Spec.Replicas)
}
```

### Listing services

```go
services, err := clientset.CoreV1().Services("").List(ctx, metav1.ListOptions{})

for _, svc := range services.Items {
    fmt.Printf("%-30s %-15s %-12s %v\n",
        svc.Name,
        svc.Namespace,
        svc.Spec.Type,
        svc.Spec.Ports)
}
```

---

## 4. Container-Aware GOMAXPROCS

When Go runs inside a Docker container, `runtime.NumCPU()` returns the **host** CPU count, not the container's CPU limit. This means Go may create too many threads and waste resources.

### The problem

```go
// On a 64-core host with a container limited to 2 CPUs:
runtime.NumCPU()    // returns 64 (wrong!)
runtime.GOMAXPROCS(0)  // defaults to 64 (way too many!)
```

### The fix - automaxprocs

```bash
go get go.uber.org/automaxprocs
```

```go
import _ "go.uber.org/automaxprocs"

func main() {
    // GOMAXPROCS is now automatically set to the container's CPU limit
    // If running outside a container, it uses the actual CPU count
    fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0))
}
```

Just import it with `_` (blank import) and it does the rest. Many production Go services use this.

---

## 5. Project: `infraview` - Infrastructure Viewer

We build a tool that:
- Connects to the local Docker daemon and lists containers with status and resource info
- Optionally connects to a K8s cluster and lists pods
- Shows everything in a nice table or JSON format
- Handles errors gracefully (Docker not running, K8s not configured, etc.)

### How to run

```bash
# Show Docker containers
go run . docker

# Show Docker containers with stats (CPU/memory)
go run . docker --stats

# Show Docker images
go run . docker --images

# Show K8s pods (needs a running cluster)
go run . k8s pods
go run . k8s pods --namespace default

# Show K8s deployments
go run . k8s deployments

# Show K8s services
go run . k8s services

# Show everything
go run . all

# JSON output
go run . docker -o json
go run . k8s pods -o json
```

### Sample output - docker command

```
=== Docker Containers ===

ID            Image              Status          Ports
----------    -----------------  --------------  -----
a1b2c3d4e5f6  nginx:latest       Up 2 hours      0.0.0.0:80->80/tcp
f6e5d4c3b2a1  redis:7-alpine     Up 5 hours      6379/tcp
1a2b3c4d5e6f  postgres:16        Up 5 hours      5432/tcp
e6d5c4b3a2f1  myapp:latest       Exited (0) 1h

Running: 3 | Stopped: 1 | Total: 4
```

### Sample output - k8s pods command

```
=== Kubernetes Pods ===

Name                              Namespace    Status   Age
----                              ---------    ------   ---
nginx-deployment-abc123-xyz       default      Running  2d
redis-master-0                    default      Running  5d
api-gateway-def456-uvw            production   Running  1d
worker-ghi789-rst                 production   Pending  5m

Total: 4 pods
```

### Project structure

```
Day-09/
  go.mod
  go.sum
  main.go
  docker.go      <- Docker operations
  k8s.go         <- Kubernetes operations
  display.go     <- Table/JSON output helpers
```

---

## Mini Challenges

1. **Easy**: Add a `--filter` flag to show only running or stopped containers
2. **Medium**: Add a `docker pull <image>` subcommand that pulls an image and shows progress
3. **Hard**: Add a `k8s logs <pod-name>` subcommand that streams pod logs in real time

---

## What We Learned Today

- `context.Context` for timeouts and cancellation
- Docker SDK - connecting, listing containers/images, getting stats
- Kubernetes client-go - connecting, listing pods/deployments/services
- Container-aware `GOMAXPROCS` with `automaxprocs`
- Graceful error handling when services are not available
- Building tools that combine multiple infrastructure data sources

---

Next up: **Day 10 - Capstone Project: Infrastructure Monitoring Agent** where we bring everything together into a production-ready monitoring tool!
