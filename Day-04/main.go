package main

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"log-parser/apperrors"
	"log-parser/parser"
)

// ================================
// CLI FLAG PARSING (simple, no external library)
// ================================

type Config struct {
	FilePath   string
	Level      string // filter level: INFO, WARN, ERROR, ALL
	OutputFile string // optional: write filtered lines to this file
}

func parseArgs() (*Config, error) {
	args := os.Args[1:]

	if len(args) == 0 {
		return nil, fmt.Errorf("no log file specified")
	}

	config := &Config{
		FilePath: args[0],
		Level:    "ALL", // default: show all levels
	}

	// Parse optional flags
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--level":
			if i+1 >= len(args) {
				return nil, &apperrors.ConfigError{
					Flag: "level", Value: "", Message: "missing value",
				}
			}
			i++
			config.Level = strings.ToUpper(args[i])

		case "--output":
			if i+1 >= len(args) {
				return nil, &apperrors.ConfigError{
					Flag: "output", Value: "", Message: "missing value",
				}
			}
			i++
			config.OutputFile = args[i]

		default:
			return nil, &apperrors.ConfigError{
				Flag: args[i], Value: "", Message: "unknown flag",
			}
		}
	}

	return config, nil
}

func printUsage() {
	fmt.Println("Usage: log-parser <logfile> [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --level LEVEL    Filter by level: INFO, WARN, ERROR, DEBUG, ALL (default: ALL)")
	fmt.Println("  --output FILE    Write filtered lines to a file")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run main.go sample.log")
	fmt.Println("  go run main.go sample.log --level ERROR")
	fmt.Println("  go run main.go sample.log --level ERROR --output errors.txt")
}

// ================================
// DISPLAY FUNCTIONS
// ================================

func printSummary(result *parser.ParseResult) {
	fmt.Printf("Parsing... done! (%d lines processed)\n\n", result.TotalLines)

	fmt.Println("Summary:")

	// Sort levels for consistent output
	levels := make([]string, 0, len(result.LevelCount))
	for level := range result.LevelCount {
		levels = append(levels, level)
	}
	sort.Strings(levels)

	for _, level := range levels {
		fmt.Printf("  %-8s: %d\n", level, result.LevelCount[level])
	}
}

func printFilteredEntries(entries []parser.LogEntry, level string) {
	if level == "ALL" {
		return // don't print all lines if no filter
	}

	fmt.Printf("\nFilter: %s only\n", level)

	if len(entries) == 0 {
		fmt.Println("  (no matching lines)")
		return
	}

	// Show max 10 lines in terminal
	limit := 10
	if len(entries) < limit {
		limit = len(entries)
	}

	fmt.Println()
	for i := 0; i < limit; i++ {
		fmt.Printf("  %s\n", entries[i].Raw)
	}

	if len(entries) > limit {
		fmt.Printf("  ... and %d more lines\n", len(entries)-limit)
	}
}

// ================================
// MAIN
// ================================

func main() {
	// Parse CLI arguments
	config, err := parseArgs()
	if err != nil {
		fmt.Printf("Error: %s\n\n", err)
		printUsage()
		os.Exit(1)
	}

	// Validate the log file
	err = parser.ValidateFile(config.FilePath)
	if err != nil {
		handleError(err)
		os.Exit(1)
	}

	// Validate the filter level
	err = parser.ValidateLevel(config.Level)
	if err != nil {
		handleError(err)
		os.Exit(1)
	}

	// Parse the log file
	fmt.Println("=== Log Parser ===")
	fmt.Printf("File: %s\n\n", config.FilePath)

	result, err := parser.ParseFile(config.FilePath)
	if err != nil {
		handleError(err)
		os.Exit(1)
	}

	// Print summary
	printSummary(result)

	// Filter entries
	filtered := parser.FilterByLevel(result.Entries, config.Level)

	// Show filtered entries in terminal
	printFilteredEntries(filtered, config.Level)

	// Write to output file if requested
	if config.OutputFile != "" {
		err = parser.WriteToFile(filtered, config.OutputFile)
		if err != nil {
			handleError(err)
			os.Exit(1)
		}

		levelLabel := config.Level
		if levelLabel == "ALL" {
			levelLabel = "all levels"
		}
		fmt.Printf("\nFiltered lines written to: %s (%d lines, filter: %s)\n", config.OutputFile, len(filtered), levelLabel)
	}

	fmt.Println()
}

// handleError prints user-friendly error messages based on error type
func handleError(err error) {
	// Try to extract our custom error types for better messages
	var fileErr *apperrors.FileError
	var configErr *apperrors.ConfigError
	var parseErr *apperrors.ParseError

	switch {
	case errors.As(err, &fileErr):
		fmt.Printf("❌ File error: %s\n", fileErr)
		switch fileErr.Message {
		case "file not found":
			fmt.Printf("   Hint: check if the file path is correct\n")
		case "file is empty (0 bytes)":
			fmt.Printf("   Hint: the log file has no content\n")
		case "path is a directory, not a file":
			fmt.Printf("   Hint: give a file path, not a directory\n")
		}

	case errors.As(err, &configErr):
		fmt.Printf("❌ Config error: %s\n", configErr)
		fmt.Println()
		printUsage()

	case errors.As(err, &parseErr):
		fmt.Printf("❌ Parse error: %s\n", parseErr)

	default:
		fmt.Printf("❌ Error: %s\n", err)
	}
}
