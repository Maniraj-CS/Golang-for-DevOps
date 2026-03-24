# Day 03 - Structs, Interfaces and JSON

## What are we doing today?

Today we will:
1. Learn structs - Go's way of grouping related data
2. Add methods to structs
3. Understand interfaces - the most powerful concept in Go
4. Work with JSON - marshal, unmarshal, struct tags
5. Build a config reader that supports both JSON and YAML

By the end of this episode, you will have a tool that reads your server inventory from a config file (JSON or YAML) and prints a nice formatted table. This is something you do all the time in DevOps - reading config files, parsing them, doing something with the data.

---

## 1. Structs

A struct is a collection of fields. Think of it like a blueprint for your data.

In DevOps terms - if you want to represent a server, you need its name, IP, port, status etc. Struct lets you group all of that together.

```go
type Server struct {
    Name   string
    IP     string
    Port   int
    Status string
}
```

### Creating a struct

```go
// Method 1: All fields
s1 := Server{
    Name:   "web-01",
    IP:     "10.0.1.10",
    Port:   8080,
    Status: "running",
}

// Method 2: Without field names (not recommended - hard to read)
s2 := Server{"web-02", "10.0.1.11", 8080, "stopped"}

// Method 3: Zero value (all fields get default values)
var s3 Server
s3.Name = "db-01"
s3.IP = "10.0.2.10"
```

### Accessing fields

```go
fmt.Println(s1.Name)   // web-01
fmt.Println(s1.IP)     // 10.0.1.10
s1.Status = "stopped"  // change a field
```

### Nested structs

```go
type HealthCheck struct {
    Endpoint string
    Interval int // in seconds
}

type Server struct {
    Name        string
    IP          string
    Port        int
    HealthCheck HealthCheck // struct inside a struct
}
```

---

## 2. Methods on Structs

A method is just a function attached to a struct. It's like saying "hey Server, you can do these things".

```go
type Server struct {
    Name   string
    IP     string
    Port   int
    Status string
}

// Method with value receiver - cannot modify the struct
func (s Server) Address() string {
    return fmt.Sprintf("%s:%d", s.IP, s.Port)
}

// Method with pointer receiver - CAN modify the struct
func (s *Server) Stop() {
    s.Status = "stopped"
    fmt.Printf("%s has been stopped\n", s.Name)
}

// Usage
server := Server{Name: "web-01", IP: "10.0.1.10", Port: 8080, Status: "running"}
fmt.Println(server.Address()) // 10.0.1.10:8080
server.Stop()                 // web-01 has been stopped
```

Rule of thumb:
- Use **pointer receiver** (`*Server`) when you want to modify the struct
- Use **value receiver** (`Server`) when you just want to read data

---

## 3. Interfaces

This is where Go gets really powerful. An interface defines a set of methods. Any struct that has those methods automatically satisfies the interface. No `implements` keyword needed.

```go
type ConfigReader interface {
    Read(path string) ([]Server, error)
}
```

Now any struct that has a `Read(path string) ([]Server, error)` method is a ConfigReader. That's it!

```go
type JSONReader struct{}

func (j JSONReader) Read(path string) ([]Server, error) {
    // read JSON file and return servers
}

type YAMLReader struct{}

func (y YAMLReader) Read(path string) ([]Server, error) {
    // read YAML file and return servers
}
```

Both `JSONReader` and `YAMLReader` satisfy the `ConfigReader` interface because they both have the `Read` method with the right signature.

Why is this useful? Because now you can write one function that works with ANY config reader:

```go
func loadServers(reader ConfigReader, path string) {
    servers, err := reader.Read(path)
    // works with JSON, YAML, TOML, whatever!
}
```

This is how real Go tools are built. Docker, Kubernetes, Terraform - all of them use interfaces heavily.

---

## 4. JSON - Marshal and Unmarshal

### Struct to JSON (Marshal)

```go
type Server struct {
    Name   string `json:"name"`
    IP     string `json:"ip"`
    Port   int    `json:"port"`
    Status string `json:"status"`
}

server := Server{Name: "web-01", IP: "10.0.1.10", Port: 8080, Status: "running"}

// Convert struct to JSON bytes
data, err := json.Marshal(server)
fmt.Println(string(data))
// {"name":"web-01","ip":"10.0.1.10","port":8080,"status":"running"}

// Pretty print with indent
data, err = json.MarshalIndent(server, "", "  ")
fmt.Println(string(data))
```

### JSON to Struct (Unmarshal)

```go
jsonStr := `{"name":"web-01","ip":"10.0.1.10","port":8080,"status":"running"}`

var server Server
err := json.Unmarshal([]byte(jsonStr), &server)
fmt.Println(server.Name) // web-01
```

### Struct Tags

Those backtick things after each field? Those are struct tags. They tell the JSON package what to name the fields.

```go
type Server struct {
    Name     string `json:"name"`
    IP       string `json:"ip_address"`
    Port     int    `json:"port"`
    Status   string `json:"status,omitempty"`   // skip if empty string
    Internal bool   `json:"internal,omitzero"`  // skip if false (zero value)
    Password string `json:"-"`                  // never include in JSON
}
```

- `omitempty` - skip the field if it has an empty value ("", 0, false, nil)
- `omitzero` - skip the field if it has its zero value (cleaner than omitempty)
- `"-"` - never include this field in JSON output

---

## 5. Reading and Writing JSON Files

```go
// Reading a JSON file
data, err := os.ReadFile("config.json")
if err != nil {
    log.Fatal(err)
}

var servers []Server
err = json.Unmarshal(data, &servers)

// Writing a JSON file
output, err := json.MarshalIndent(servers, "", "  ")
err = os.WriteFile("output.json", output, 0644)
```

---

## 6. Project: Config Reader Tool

We build a tool that:
- Reads a server inventory config file
- Supports both JSON and YAML formats (using interfaces!)
- Prints a nicely formatted table
- Shows a summary of servers by environment and status

### Config files

Check `servers.json` and `servers.yaml` in this folder - both have the same data in different formats.

### How to run

```bash
# Read JSON config
go run main.go servers.json

# Read YAML config
go run main.go servers.yaml

# Build and run
go build -o config-reader
./config-reader servers.json
```

### Sample output

```
=== Server Inventory ===

Name          IP             Port   Env          Status     
----          --             ----   ---          ------     
web-01        10.0.1.10      8080   production   running    
web-02        10.0.1.11      8080   production   running    
api-01        10.0.1.20      3000   production   running    
db-01         10.0.2.10      5432   production   running    
db-02         10.0.2.11      5432   staging      stopped    
cache-01      10.0.3.10      6379   production   running    
dev-01        10.0.4.10      8080   development  running    
test-01       10.0.4.20      8080   staging      running    

Total: 8 servers

By Environment:
  development  : 1
  production   : 5
  staging      : 2

By Status:
  running      : 7
  stopped      : 1
```

---

## Mini Challenges

1. **Easy**: Add a `Region` field to the server struct and update the config files
2. **Medium**: Add a `--env` flag to filter servers by environment (e.g., `./config-reader servers.json --env production`)
3. **Hard**: Add TOML support as another config reader (you will need the `github.com/BurntSushi/toml` package)

---

## What We Learned Today

- Structs - grouping related data together
- Methods on structs - value vs pointer receiver
- Interfaces - the Go way of polymorphism (no `implements` keyword!)
- JSON marshal/unmarshal and struct tags
- `omitempty` and `omitzero` for cleaner JSON
- Reading config files in different formats
- Using interfaces to make code flexible

---

Next up: **Day 04 - File I/O, Error Handling and Packages** where we will build a proper log parser with custom error handling!
