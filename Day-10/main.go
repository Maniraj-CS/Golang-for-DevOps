package main

// ================================
// kubectl-status: Entry Point
// ================================
// This is the main.go - the starting point of our kubectl plugin.
// It just calls cmd.Execute() which starts Cobra.
// Simple, clean, the Go way!

import "kubectl-status/cmd"

func main() {
	cmd.Execute()
}
