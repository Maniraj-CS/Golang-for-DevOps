# Day 08 - Building a Full CLI App with Cobra

## What are we doing today?

Today we will:
1. Learn why Cobra is used in almost every big Go CLI tool (kubectl, docker, gh)
2. Build a CLI app with subcommands, flags, and help text
3. Use persistent flags vs local flags
4. Read config from a file using Viper
5. Understand `go tool` directive for tracking tool dependencies
6. Build and distribute Go binaries with cross-compilation
7. Build a DevOps swiss knife CLI called `devctl`

By the end of this episode, you will have a proper CLI tool called `devctl` with subcommands like `devctl health`, `devctl sysinfo`, `devctl github` - combining what we built in earlier episodes into one clean, professional tool. This is exactly how real DevOps CLI tools are built.

---

## 1. Why Cobra?

Every big Go CLI tool uses Cobra:
- `kubectl` (Kubernetes)
- `docker`
- `gh` (GitHub CLI)
- `terraform`
- `hugo`

Why? Because it gives you:
- Subcommands (`devctl health`, `devctl sysinfo`)
- Flags (`--port 8080`, `--json`)
- Auto-generated help text
- Bash/Zsh completion
- Config file support (with Viper)

### Installing Cobra

```bash
go get github.com/spf13/cobra@latest
```

---

## 2. Basic Cobra App

### The structure

A Cobra app has a **root command** and **sub-commands**.

```
devctl                 -> root command (shows help)
devctl health          -> subcommand
devctl sysinfo         -> subcommand
devctl github repos    -> nested subcommand
```

### Root command

```go
package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "devctl",
    Short: "DevOps swiss knife CLI",
    Long:  "A CLI tool that combines health checks, system info, and GitHub operations.",
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
```

### Adding a subcommand

```go
var healthCmd = &cobra.Command{
    Use:   "health [urls...]",
    Short: "Check health of URLs",
    Args:  cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        // args contains the URLs
        for _, url := range args {
            fmt.Printf("Checking %s...\n", url)
        }
    },
}

func init() {
    rootCmd.AddCommand(healthCmd)
}
```

---

## 3. Flags - Local vs Persistent

### Local flags (only for one command)

```go
var port int

func init() {
    healthCmd.Flags().IntVar(&port, "port", 8080, "Port to check")
    // --port only works with "devctl health"
}
```

### Persistent flags (available for all subcommands)

```go
var verbose bool

func init() {
    rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Verbose output")
    // --verbose works with ALL subcommands
}
```

### Flag types

```go
cmd.Flags().StringVar(&name, "name", "default", "Description")
cmd.Flags().IntVar(&count, "count", 10, "Description")
cmd.Flags().BoolVar(&force, "force", false, "Description")
cmd.Flags().StringSliceVar(&tags, "tag", []string{}, "Tags (can repeat)")
// --tag prod --tag staging
```

### Short flags

```go
cmd.Flags().StringVarP(&output, "output", "o", "table", "Output format")
// -o json  OR  --output json
```

---

## 4. Reading Config Files with Viper

Viper is Cobra's best friend. It reads config from:
- Config files (YAML, JSON, TOML)
- Environment variables
- Flags

### Basic Viper setup

```go
import "github.com/spf13/viper"

func initConfig() {
    viper.SetConfigName("devctl")     // config file name (no extension)
    viper.SetConfigType("yaml")       // file type
    viper.AddConfigPath(".")          // look in current directory
    viper.AddConfigPath("$HOME/.devctl") // also look in home dir

    viper.AutomaticEnv() // read environment variables too

    if err := viper.ReadInConfig(); err != nil {
        // Config file not found is OK - we have defaults
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            fmt.Println("Error reading config:", err)
        }
    }
}
```

### Sample config file (devctl.yaml)

```yaml
github:
  token: "ghp_your_token_here"
  default_user: "abhishekkumar"

health:
  timeout: 10
  concurrent: 5

output: "table"
verbose: false
```

### Reading config values

```go
token := viper.GetString("github.token")
timeout := viper.GetInt("health.timeout")
verbose := viper.GetBool("verbose")
```

### Binding flags to Viper

```go
// Flag value takes priority over config file
viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
```

Priority order (highest to lowest):
1. CLI flags
2. Environment variables
3. Config file
4. Default values

---

## 5. Cross-Compilation

One of Go's superpowers - build a binary for any OS from any OS!

### Build for current OS

```bash
go build -o devctl .
```

### Build for different OS/architecture

```bash
# Build for Linux (great for servers and Docker)
GOOS=linux GOARCH=amd64 go build -o devctl-linux-amd64 .

# Build for Linux ARM (Raspberry Pi, AWS Graviton)
GOOS=linux GOARCH=arm64 go build -o devctl-linux-arm64 .

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o devctl-windows-amd64.exe .

# Build for macOS Intel
GOOS=darwin GOARCH=amd64 go build -o devctl-darwin-amd64 .

# Build for macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o devctl-darwin-arm64 .
```

### Build all platforms at once (simple script)

```bash
#!/bin/bash
APP="devctl"
PLATFORMS=("linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64" "windows/amd64")

for PLATFORM in "${PLATFORMS[@]}"; do
    OS="${PLATFORM%/*}"
    ARCH="${PLATFORM#*/}"
    OUTPUT="${APP}-${OS}-${ARCH}"
    if [ "$OS" = "windows" ]; then
        OUTPUT="${OUTPUT}.exe"
    fi
    echo "Building $OUTPUT..."
    GOOS=$OS GOARCH=$ARCH go build -o "build/$OUTPUT" .
done
echo "Done!"
```

### Adding version info at build time

```go
// In your code
var Version = "dev"
var BuildTime = "unknown"

// Build with:
// go build -ldflags "-X main.Version=1.0.0 -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o devctl .
```

---

## 6. `go tool` Directive

In Go, you can track CLI tool dependencies in your `go.mod` using the `tool` directive:

```go
// go.mod
module devctl

go 1.26

tool github.com/spf13/cobra-cli

require (
    github.com/spf13/cobra v1.9.1
    github.com/spf13/viper v1.20.0
)
```

Then you can run the tool with:
```bash
go tool cobra-cli init
go tool cobra-cli add health
```

This ensures everyone on your team uses the same version of the tool.

---

## 7. Project: `devctl` - DevOps Swiss Knife CLI

We build a tool that combines what we built in earlier episodes:
- `devctl health <url1> <url2>` - Check if URLs are up or down (from Day-05)
- `devctl sysinfo` - Show system information (from Day-06)
- `devctl github repos <username>` - List GitHub repos (from Day-07)
- `devctl version` - Show version info

### How to run

```bash
# Build the tool
go build -o devctl .

# Or just run directly
go run . health https://google.com https://github.com

# System info
go run . sysinfo
go run . sysinfo --json
go run . sysinfo --section disk

# GitHub repos
go run . github repos golang
go run . github repos golang --sort stars
go run . github repos torvalds --language c

# Version
go run . version

# Global flags (work with all commands)
go run . --verbose health https://google.com
go run . --output json sysinfo
```

### Sample output - health command

```
=== Health Check ===

URL                              Status  Code  Time
---                              ------  ----  ----
https://google.com               UP      200   120ms
https://github.com               UP      200   245ms
https://nonexistent.example.com  DOWN    -     0ms (connection refused)

Checked: 3 | Up: 2 | Down: 1
```

### Sample output - sysinfo command

```
=== System Info ===

Host
  Hostname:   macbook-pro.local
  OS:         darwin
  Arch:       arm64
  Uptime:     5 days, 3:22

CPU
  Cores:      10
  Load Avg:   2.15, 1.98, 1.76

Memory
  Total:      16 GB
  Used:       11 GB
  Available:  5 GB
  Usage:      68.7%

Disk
  /          466 GB total, 234 GB used, 232 GB free (50%)
```

### Sample output - github command

```
=== GitHub Repos for golang ===

Name                Language     Stars  Forks  Last Push
----                --------     -----  -----  ---------
go                  Go           133109 18873  2026-03-20
groupcache          Go           13339  1395   2024-11-29
dep                 Go           12755  1030   2020-09-05

Total: 30 repos (sorted by: stars)
```

### Project structure

```
Day-08/
  go.mod
  go.sum
  main.go            <- entry point (just calls cmd.Execute())
  cmd/
    root.go          <- root command + global flags
    health.go        <- health subcommand
    sysinfo.go       <- sysinfo subcommand
    github.go        <- github subcommand with "repos" sub-subcommand
    version.go       <- version subcommand
  devctl.yaml        <- sample config file
```

---

## Mini Challenges

1. **Easy**: Add a `devctl github repos <user> --top 5` flag that shows only top N repos
2. **Medium**: Add a `devctl config show` subcommand that prints the current config (from file + env + defaults)
3. **Hard**: Add bash/zsh completion generation with `devctl completion bash > /etc/bash_completion.d/devctl`

---

## What We Learned Today

- Cobra for building professional CLI tools with subcommands and flags
- Local flags vs persistent flags
- Viper for reading config from files, env vars, and flags
- `go tool` directive for tracking tool dependencies in go.mod
- Cross-compilation with `GOOS` and `GOARCH`
- Building binaries with version info using `-ldflags`
- The pattern of combining multiple tools into one CLI (swiss knife approach)

---

Next up: **Day 09 - Docker and Kubernetes from Go** where we will talk to the Docker daemon and K8s API from Go code!
