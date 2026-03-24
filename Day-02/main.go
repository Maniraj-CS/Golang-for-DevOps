package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	// Check if user gave a file path
	if len(os.Args) < 2 {
		fmt.Println("Usage: log-analyzer <logfile>")
		fmt.Println("Example: log-analyzer sample.log")
		os.Exit(1)
	}

	filePath := os.Args[1]

	// Open the log file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err)
		os.Exit(1)
	}
	defer file.Close() // This runs when the function exits. Always close files!

	// Maps to store our counts
	levelCount := map[string]int{
		"INFO":  0,
		"WARN":  0,
		"ERROR": 0,
	}

	// Store error messages to show top errors later
	errorMessages := []string{}

	totalLines := 0

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		totalLines++

		// Check which log level this line has
		// We check for each level in the line
		switch {
		case strings.Contains(line, "ERROR"):
			levelCount["ERROR"]++
			// Extract just the message part (after the log level)
			msg := extractMessage(line, "ERROR")
			errorMessages = append(errorMessages, msg)

		case strings.Contains(line, "WARN"):
			levelCount["WARN"]++

		case strings.Contains(line, "INFO"):
			levelCount["INFO"]++
		}
	}

	// Check if there was an error while reading
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}

	// Print the summary
	printSummary(filePath, totalLines, levelCount, errorMessages)
}

// extractMessage pulls out the message part after the log level
func extractMessage(line string, level string) string {
	// Find where the level is in the line
	index := strings.Index(line, level)
	if index == -1 {
		return line
	}
	// Get everything after the level keyword
	msg := strings.TrimSpace(line[index+len(level):])
	return msg
}

// printSummary prints a nice summary of the log analysis
func printSummary(filePath string, totalLines int, levelCount map[string]int, errorMessages []string) {
	fmt.Println("\n=== Log File Analysis ===")
	fmt.Printf("File: %s\n", filePath)
	fmt.Printf("Total lines: %d\n\n", totalLines)

	// Print level counts
	fmt.Printf("%-12s %s\n", "Level", "Count")
	fmt.Printf("%-12s %s\n", "-----", "-----")

	// Sort the levels so output is consistent
	levels := []string{"INFO", "WARN", "ERROR"}
	for _, level := range levels {
		count := levelCount[level]
		fmt.Printf("%-12s %d\n", level, count)
	}

	// Print top error messages (unique ones)
	if len(errorMessages) > 0 {
		fmt.Println("\nTop ERROR messages:")
		uniqueErrors := getUniqueErrors(errorMessages)

		// Show max 5 errors
		limit := 5
		if len(uniqueErrors) < limit {
			limit = len(uniqueErrors)
		}

		for i := 0; i < limit; i++ {
			fmt.Printf("  %d. %s\n", i+1, uniqueErrors[i])
		}
	}

	fmt.Println()
}

// getUniqueErrors returns unique error messages sorted by how often they appear
func getUniqueErrors(messages []string) []string {
	// Count each unique message
	counts := map[string]int{}
	for _, msg := range messages {
		counts[msg]++
	}

	// Get unique messages
	unique := []string{}
	for msg := range counts {
		unique = append(unique, msg)
	}

	// Sort by count (most frequent first)
	sort.Slice(unique, func(i, j int) bool {
		return counts[unique[i]] > counts[unique[j]]
	})

	return unique
}
