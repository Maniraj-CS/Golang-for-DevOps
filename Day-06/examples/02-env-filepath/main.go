package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// ================================
	// 1. ENVIRONMENT VARIABLES
	// ================================

	fmt.Println("=== Reading Environment Variables ===")

	// Get specific env vars
	home := os.Getenv("HOME")
	user := os.Getenv("USER")
	shell := os.Getenv("SHELL")
	gopath := os.Getenv("GOPATH")

	fmt.Println("  HOME  :", home)
	fmt.Println("  USER  :", user)
	fmt.Println("  SHELL :", shell)
	fmt.Println("  GOPATH:", gopath)

	// LookupEnv - check if var exists (Getenv returns "" for both unset and empty)
	fmt.Println("\n=== LookupEnv (check existence) ===")

	vars := []string{"HOME", "KUBERNETES_SERVICE_HOST", "MY_SECRET"}
	for _, v := range vars {
		val, exists := os.LookupEnv(v)
		if exists {
			// Show only first 30 chars for long values
			if len(val) > 30 {
				val = val[:30] + "..."
			}
			fmt.Printf("  ✅ %-30s = %s\n", v, val)
		} else {
			fmt.Printf("  ❌ %-30s (not set)\n", v)
		}
	}

	// Set and unset
	fmt.Println("\n=== Set and Unset ===")

	os.Setenv("MY_APP_NAME", "sysinfo-collector")
	os.Setenv("MY_APP_VERSION", "1.0.0")
	fmt.Println("  MY_APP_NAME   :", os.Getenv("MY_APP_NAME"))
	fmt.Println("  MY_APP_VERSION:", os.Getenv("MY_APP_VERSION"))

	os.Unsetenv("MY_APP_NAME")
	fmt.Println("  After unset   :", os.Getenv("MY_APP_NAME"), "(empty)")

	// ================================
	// 2. FILEPATH PACKAGE
	// ================================

	fmt.Println("\n=== filepath.Join (safe path building) ===")

	paths := [][]string{
		{"home", "user", "config.json"},
		{"/var", "log", "server.log"},
		{"/tmp", "..", "etc", "passwd"},
	}

	for _, parts := range paths {
		joined := filepath.Join(parts...)
		fmt.Printf("  %v -> %s\n", parts, joined)
	}

	fmt.Println("\n=== filepath Parts ===")

	testPaths := []string{
		"/var/log/server.log",
		"/home/user/config.yaml",
		"server.log.gz",
		"/tmp/backup/db-2024.sql.gz",
	}

	fmt.Printf("  %-35s %-15s %-20s %s\n", "Path", "Base", "Dir", "Ext")
	fmt.Printf("  %-35s %-15s %-20s %s\n", "----", "----", "---", "---")
	for _, p := range testPaths {
		fmt.Printf("  %-35s %-15s %-20s %s\n",
			p,
			filepath.Base(p),
			filepath.Dir(p),
			filepath.Ext(p),
		)
	}

	// ================================
	// 3. FILEPATH.WALK - traverse directory tree
	// ================================

	fmt.Println("\n=== filepath.Walk (directory tree) ===")

	// Walk the current directory (just 1 level deep for demo)
	depth := 0
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip errors
		}

		// Limit depth
		depth = strings.Count(path, string(filepath.Separator))
		if depth > 1 {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		icon := "📄"
		if info.IsDir() {
			icon = "📁"
		}

		fmt.Printf("  %s %-30s %d bytes\n", icon, path, info.Size())
		return nil
	})
	if err != nil {
		fmt.Println("Walk error:", err)
	}

	// ================================
	// 4. FILEPATH.GLOB - find files by pattern
	// ================================

	fmt.Println("\n=== filepath.Glob (pattern matching) ===")

	patterns := []string{"*.go", "examples/*/*.go", "*.mod"}
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			fmt.Printf("  Pattern: %-20s Error: %s\n", pattern, err)
			continue
		}
		if len(matches) == 0 {
			fmt.Printf("  Pattern: %-20s (no matches)\n", pattern)
		} else {
			fmt.Printf("  Pattern: %-20s -> %v\n", pattern, matches)
		}
	}
}
