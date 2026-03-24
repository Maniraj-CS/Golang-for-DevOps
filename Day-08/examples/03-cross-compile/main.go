package main

// ================================
// EXAMPLE 03 - Cross-Compilation
// ================================
// This program prints build info. Build it for different platforms:
//
//   go build -o myapp .
//   GOOS=linux GOARCH=amd64 go build -o myapp-linux .
//   GOOS=windows GOARCH=amd64 go build -o myapp.exe .
//
// With version info:
//   go build -ldflags "-X main.Version=1.0.0 -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o myapp .

import (
	"fmt"
	"runtime"
)

// These variables are set at build time using -ldflags
// If not set, they use these defaults
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	fmt.Println("=== Build Info ===")
	fmt.Println()
	fmt.Printf("Version:    %s\n", Version)
	fmt.Printf("Build Time: %s\n", BuildTime)
	fmt.Printf("Git Commit: %s\n", GitCommit)
	fmt.Printf("Go Version: %s\n", runtime.Version())
	fmt.Printf("OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Compiler:   %s\n", runtime.Compiler)

	fmt.Println()
	fmt.Println("=== How to Set Build Info ===")
	fmt.Println()
	fmt.Println("Build with version info:")
	fmt.Println(`  go build -ldflags "-X main.Version=1.0.0 -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o myapp .`)
	fmt.Println()
	fmt.Println("Cross-compile for Linux:")
	fmt.Println(`  GOOS=linux GOARCH=amd64 go build -o myapp-linux .`)
	fmt.Println()
	fmt.Println("Cross-compile for Windows:")
	fmt.Println(`  GOOS=windows GOARCH=amd64 go build -o myapp.exe .`)
	fmt.Println()
	fmt.Println("Supported GOOS values:  linux, darwin, windows, freebsd")
	fmt.Println("Supported GOARCH values: amd64, arm64, arm, 386")
}
