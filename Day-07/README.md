# Day 07 - Working with REST APIs

## What are we doing today?

Today we will:
1. Make POST, PUT, DELETE requests (not just GET)
2. Set headers and auth tokens
3. Parse nested JSON responses
4. Build our own HTTP server with `http.ServeMux`
5. Serve JSON responses from our server
6. Use environment variables for secrets
7. Build a tool that talks to GitHub API and also runs a local API server

By the end of this episode, you will have a tool that lists your GitHub repos, checks their last commit date, and also exposes a local `/repos` API endpoint that other tools can consume. This is a very common pattern in DevOps - build a tool that both uses APIs and serves APIs.

---

## 1. Making POST, PUT, DELETE Requests

In Day-05 we made GET requests. Now let's learn the rest.

### POST request with JSON body

```go
import (
    "bytes"
    "encoding/json"
    "net/http"
)

// Create the request body
body := map[string]string{
    "title": "Bug: Server crash on deploy",
    "body":  "The server crashes when deploying v2.3.1",
}

jsonBody, err := json.Marshal(body)
if err != nil {
    log.Fatal(err)
}

// Make the POST request
resp, err := http.Post(
    "https://api.example.com/issues",
    "application/json",
    bytes.NewBuffer(jsonBody),
)
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

fmt.Println("Status:", resp.StatusCode) // 201 Created
```

### Using http.NewRequest (more control)

```go
// Create request
req, err := http.NewRequest("POST", "https://api.example.com/issues", bytes.NewBuffer(jsonBody))
if err != nil {
    log.Fatal(err)
}

// Set headers
req.Header.Set("Content-Type", "application/json")
req.Header.Set("Authorization", "Bearer your-token-here")
req.Header.Set("Accept", "application/json")

// Send it
client := &http.Client{Timeout: 10 * time.Second}
resp, err := client.Do(req)
```

### PUT request (update something)

```go
req, err := http.NewRequest("PUT", "https://api.example.com/issues/42", bytes.NewBuffer(jsonBody))
req.Header.Set("Content-Type", "application/json")
req.Header.Set("Authorization", "Bearer "+token)

resp, err := client.Do(req)
```

### DELETE request

```go
req, err := http.NewRequest("DELETE", "https://api.example.com/issues/42", nil)
req.Header.Set("Authorization", "Bearer "+token)

resp, err := client.Do(req)
```

---

## 2. Parsing Nested JSON

Real API responses are deeply nested. Here's how to handle them.

### With structs (recommended)

```go
type GithubRepo struct {
    Name        string `json:"name"`
    FullName    string `json:"full_name"`
    Description string `json:"description"`
    Private     bool   `json:"private"`
    HTMLURL     string `json:"html_url"`
    Language    string `json:"language"`
    Owner       struct {
        Login     string `json:"login"`
        AvatarURL string `json:"avatar_url"`
    } `json:"owner"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
    PushedAt  string `json:"pushed_at"`
}
```

### With map[string]any (when you don't know the structure)

```go
var result map[string]any
json.Unmarshal(data, &result)

// Access nested fields with type assertions
owner := result["owner"].(map[string]any)
login := owner["login"].(string)
```

The struct approach is safer and gives you autocomplete. Use `map[string]any` only when the JSON structure is truly unknown.

---

## 3. Building an HTTP Server

### Basic server with ServeMux

```go
mux := http.NewServeMux()

mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "status": "ok",
    })
})

mux.HandleFunc("GET /api/servers", func(w http.ResponseWriter, r *http.Request) {
    servers := []map[string]string{
        {"name": "web-01", "status": "running"},
        {"name": "db-01", "status": "running"},
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(servers)
})

fmt.Println("Server starting on :8080")
http.ListenAndServe(":8080", mux)
```

### ServeMux patterns (method + path)

In recent Go versions, ServeMux supports method-based routing:

```go
mux.HandleFunc("GET /api/users", listUsers)       // GET only
mux.HandleFunc("POST /api/users", createUser)     // POST only
mux.HandleFunc("GET /api/users/{id}", getUser)    // path parameter!
mux.HandleFunc("DELETE /api/users/{id}", deleteUser)
```

### Reading path parameters

```go
mux.HandleFunc("GET /api/users/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    fmt.Fprintf(w, "User ID: %s", id)
})
```

### Reading query parameters

```go
// URL: /api/repos?sort=updated&per_page=10
sort := r.URL.Query().Get("sort")       // "updated"
perPage := r.URL.Query().Get("per_page") // "10"
```

---

## 4. Environment Variables for Secrets

Never hardcode tokens or passwords! Use environment variables.

```go
token := os.Getenv("GITHUB_TOKEN")
if token == "" {
    log.Fatal("GITHUB_TOKEN environment variable is not set")
}
```

Set it before running:
```bash
export GITHUB_TOKEN="ghp_your_token_here"
go run main.go
```

---

## 5. Project: GitHub API Tool + Local Server

We build a tool that:
- Lists your GitHub repos (public, no auth needed for public repos)
- Shows repo name, language, stars, last push date
- Exposes a local `/repos` endpoint that serves the same data as JSON
- Has a `/health` endpoint
- Supports `GITHUB_TOKEN` env var for accessing private repos

### How to run

```bash
# List repos for a user (no auth needed for public repos)
go run main.go repos abhishekkumar

# List repos in table format
go run main.go repos abhishekkumar --sort updated

# Start the local API server
go run main.go serve

# Start on a custom port
go run main.go serve --port 9090

# Then hit the API:
# curl http://localhost:8080/health
# curl http://localhost:8080/api/repos/abhishekkumar
```

### Sample output - repos command

```
=== GitHub Repos for abhishekkumar ===

Name                     Language       Stars  Last Push
----                     --------       -----  ---------
kubernetes-tutorial      Go             145    2024-01-10
devops-scripts           Shell          89     2024-01-08
terraform-modules        HCL            56     2023-12-20
docker-compose-apps      Dockerfile     34     2023-11-15
ansible-playbooks        YAML           23     2023-10-01

Total: 5 public repos
```

### Sample output - serve command

```
=== GitHub Tool - API Server ===
Endpoints:
  GET /health                    - Health check
  GET /api/repos/{username}      - List GitHub repos for a user

Server starting on :8080
```

---

## Mini Challenges

1. **Easy**: Add a `--language` flag to filter repos by programming language
2. **Medium**: Add a `POST /api/issues/{owner}/{repo}` endpoint that creates a GitHub issue (needs GITHUB_TOKEN)
3. **Hard**: Add a `/api/repos/{username}/stats` endpoint that returns aggregated stats - total stars, most used language, most recently updated repo

---

## What We Learned Today

- Making POST, PUT, DELETE requests with `http.NewRequest`
- Setting headers and auth tokens
- Parsing nested JSON into structs
- Building HTTP servers with `http.ServeMux`
- Method-based routing (`GET /path`, `POST /path`)
- Path parameters with `{id}` and `r.PathValue()`
- Query parameters with `r.URL.Query()`
- Serving JSON responses
- Using environment variables for secrets
- The pattern of building tools that both consume and serve APIs

---

Next up: **Day 08 - Building a Full CLI App with Cobra** where we will build a proper CLI tool with subcommands, flags, and config files!
