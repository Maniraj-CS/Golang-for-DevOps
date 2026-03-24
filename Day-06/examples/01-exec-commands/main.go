package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// ================================
	// 1. BASIC COMMAND - Get output
	// ================================

	fmt.Println("=== Basic Command ===")

	out, err := exec.Command("hostname").Output()
	if err != nil {
		log.Fatal("Cannot get hostname:", err)
	}
	fmt.Println("Hostname:", strings.TrimSpace(string(out)))

	// ================================
	// 2. COMMAND WITH ARGUMENTS
	// ================================

	fmt.Println("\n=== Command with Arguments ===")

	// Each argument is separate - NOT one big string!
	out, err = exec.Command("uname", "-a").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("System:", strings.TrimSpace(string(out)))

	// List files
	out, err = exec.Command("ls", "-la", "/tmp").Output()
	if err != nil {
		log.Fatal(err)
	}
	// Show only first 3 lines
	lines := strings.Split(string(out), "\n")
	for i, line := range lines {
		if i >= 3 {
			fmt.Printf("  ... and %d more lines\n", len(lines)-3)
			break
		}
		fmt.Println(" ", line)
	}

	// ================================
	// 3. CAPTURING STDOUT AND STDERR SEPARATELY
	// ================================

	fmt.Println("\n=== Stdout vs Stderr ===")

	// This command will fail - /nonexistent doesn't exist
	cmd := exec.Command("ls", "/nonexistent-path-12345")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println("Command failed (expected):")
		fmt.Println("  Stdout:", stdout.String())
		fmt.Println("  Stderr:", strings.TrimSpace(stderr.String()))

		// Get exit code
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Println("  Exit code:", exitErr.ExitCode())
		}
	}

	// ================================
	// 4. SHELL COMMANDS (pipes, redirects)
	// ================================

	fmt.Println("\n=== Shell Commands (pipes) ===")

	// Use sh -c for shell features like pipes
	out, err = exec.Command("sh", "-c", "ps aux | head -5").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Top processes:")
	fmt.Println(string(out))

	// Count Go processes
	out, err = exec.Command("sh", "-c", "ps aux | grep -c '[g]o'").Output()
	if err != nil {
		// grep returns exit code 1 if no matches, that's okay
		fmt.Println("  No Go processes found")
	} else {
		fmt.Println("Go processes:", strings.TrimSpace(string(out)))
	}

	// ================================
	// 5. ENVIRONMENT VARIABLES FOR COMMANDS
	// ================================

	fmt.Println("\n=== Environment Variables for Commands ===")

	cmd = exec.Command("sh", "-c", "echo \"App: $MY_APP, Env: $MY_ENV\"")
	// Start with current env, then add our own
	cmd.Env = append(os.Environ(),
		"MY_APP=sysinfo-collector",
		"MY_ENV=production",
	)
	out, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.TrimSpace(string(out)))

	// ================================
	// 6. CHECKING IF A COMMAND EXISTS
	// ================================

	fmt.Println("\n=== Checking if Command Exists ===")

	commands := []string{"go", "docker", "kubectl", "terraform", "nonexistent-tool"}

	for _, name := range commands {
		path, err := exec.LookPath(name)
		if err != nil {
			fmt.Printf("  ❌ %-20s not found\n", name)
		} else {
			fmt.Printf("  ✅ %-20s -> %s\n", name, path)
		}
	}

	// ================================
	// 7. RUNNING WITH CombinedOutput
	// ================================

	fmt.Println("\n=== CombinedOutput (stdout + stderr together) ===")

	// CombinedOutput merges stdout and stderr into one
	out, err = exec.Command("sh", "-c", "echo 'hello stdout'; echo 'hello stderr' >&2; echo 'bye stdout'").CombinedOutput()
	// err might be non-nil but we still get the output
	fmt.Println("Combined output:")
	fmt.Println(string(out))
}
