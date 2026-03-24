package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// ================================
	// 1. READING ENVIRONMENT VARIABLES
	// ================================

	fmt.Println("=== Reading Environment Variables ===")

	// os.Getenv returns empty string if not set
	home := os.Getenv("HOME")
	fmt.Println("HOME:", home)

	user := os.Getenv("USER")
	fmt.Println("USER:", user)

	shell := os.Getenv("SHELL")
	fmt.Println("SHELL:", shell)

	// ================================
	// 2. CHECKING IF A VARIABLE IS SET
	// ================================

	fmt.Println("\n=== Checking if Variable is Set ===")

	// os.Getenv returns "" for both "not set" and "set to empty"
	// Use os.LookupEnv to tell the difference!

	token, exists := os.LookupEnv("GITHUB_TOKEN")
	if !exists {
		fmt.Println("GITHUB_TOKEN is not set (this is expected in this example)")
	} else if token == "" {
		fmt.Println("GITHUB_TOKEN is set but empty")
	} else {
		// Don't print the full token! Just show first few chars
		masked := token[:4] + strings.Repeat("*", len(token)-4)
		fmt.Println("GITHUB_TOKEN:", masked)
	}

	// ================================
	// 3. GETENV WITH DEFAULT VALUE (helper pattern)
	// ================================

	fmt.Println("\n=== Getenv with Default ===")

	port := getEnvOrDefault("APP_PORT", "8080")
	env := getEnvOrDefault("APP_ENV", "development")
	logLevel := getEnvOrDefault("LOG_LEVEL", "info")

	fmt.Println("Port:", port)
	fmt.Println("Environment:", env)
	fmt.Println("Log Level:", logLevel)

	// ================================
	// 4. REQUIRED ENV VARS (fail fast)
	// ================================

	fmt.Println("\n=== Required Environment Variables ===")

	// In real apps, some env vars are required. Fail fast if missing!
	// Uncommenting the line below will stop the program:
	// dbHost := requireEnv("DATABASE_HOST")

	fmt.Println("Skipping requireEnv() demo - it would stop the program")
	fmt.Println("In real code, use it for required vars like DATABASE_HOST, API_KEY, etc.")

	// ================================
	// 5. SETTING ENVIRONMENT VARIABLES
	// ================================

	fmt.Println("\n=== Setting Environment Variables ===")

	// You can set env vars from your code too
	os.Setenv("MY_APP_VERSION", "1.2.3")
	os.Setenv("MY_APP_MODE", "debug")

	fmt.Println("MY_APP_VERSION:", os.Getenv("MY_APP_VERSION"))
	fmt.Println("MY_APP_MODE:", os.Getenv("MY_APP_MODE"))

	// Note: These are only set for the current process and its children.
	// They don't change your terminal's environment!

	// ================================
	// 6. LISTING ALL ENV VARS
	// ================================

	fmt.Println("\n=== All Environment Variables (first 5) ===")

	allEnv := os.Environ() // returns []string like "KEY=VALUE"
	count := 5
	if len(allEnv) < count {
		count = len(allEnv)
	}
	for i := 0; i < count; i++ {
		parts := strings.SplitN(allEnv[i], "=", 2)
		key := parts[0]
		value := parts[1]
		// Truncate long values
		if len(value) > 50 {
			value = value[:50] + "..."
		}
		fmt.Printf("  %s = %s\n", key, value)
	}
	fmt.Printf("  ... and %d more\n", len(allEnv)-count)

	// ================================
	// KEY TAKEAWAYS
	// ================================

	fmt.Println("\n=== Key Takeaways ===")
	fmt.Println("1. Use os.Getenv() for simple reads")
	fmt.Println("2. Use os.LookupEnv() to check if a var exists")
	fmt.Println("3. Create helper functions like getEnvOrDefault()")
	fmt.Println("4. NEVER hardcode secrets - always use env vars!")
	fmt.Println("5. In DevOps, env vars are used for: tokens, database URLs, ports, feature flags")
}

// getEnvOrDefault returns the env var value, or the default if not set
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// requireEnv returns the env var value, or exits the program if not set
func requireEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("ERROR: Required environment variable %s is not set\n", key)
	}
	return value
}
