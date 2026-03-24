package main

import (
	"errors"
	"fmt"
	"os"
)

// ================================
// CUSTOM ERROR TYPES
// ================================

// FileError - our custom error for file operations
type FileError struct {
	Path    string
	Op      string // what we were trying to do
	Message string
}

// Error() makes FileError satisfy the error interface
func (e *FileError) Error() string {
	return fmt.Sprintf("%s %s: %s", e.Op, e.Path, e.Message)
}

// ValidationError - for when data is invalid
type ValidationError struct {
	Field   string
	Value   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for %s=%q: %s", e.Field, e.Value, e.Message)
}

// ================================
// FUNCTIONS THAT RETURN ERRORS
// ================================

// validateLogLevel checks if the given level is valid
func validateLogLevel(level string) error {
	validLevels := map[string]bool{
		"INFO": true, "WARN": true, "ERROR": true, "DEBUG": true,
	}

	if !validLevels[level] {
		return &ValidationError{
			Field:   "level",
			Value:   level,
			Message: "must be one of: INFO, WARN, ERROR, DEBUG",
		}
	}
	return nil
}

// readLogFile tries to open and validate a log file
func readLogFile(path string) error {
	// Check if file exists
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &FileError{Path: path, Op: "open", Message: "file not found"}
		}
		return fmt.Errorf("readLogFile: %w", err) // wrap unknown errors
	}

	// Check if it's a directory
	if info.IsDir() {
		return &FileError{Path: path, Op: "open", Message: "path is a directory, not a file"}
	}

	// Check if file is empty
	if info.Size() == 0 {
		return &FileError{Path: path, Op: "read", Message: "file is empty"}
	}

	return nil
}

// processConfig wraps errors with context using %w
func processConfig(path string) error {
	err := readLogFile(path)
	if err != nil {
		// Wrap the error with extra context
		return fmt.Errorf("processConfig failed: %w", err)
	}
	return nil
}

func main() {
	// ================================
	// 1. BASIC ERROR HANDLING
	// ================================

	fmt.Println("=== Basic Error Handling ===")

	_, err := os.Open("this-file-does-not-exist.txt")
	if err != nil {
		fmt.Println("Error:", err)
	}

	// ================================
	// 2. errors.Is - CHECK FOR SPECIFIC ERROR
	// ================================

	fmt.Println("\n=== errors.Is ===")

	_, err = os.Open("missing.txt")
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("✅ errors.Is caught: file does not exist")
	}

	// Works through wrapped errors too!
	wrapped := fmt.Errorf("config load failed: %w", os.ErrNotExist)
	if errors.Is(wrapped, os.ErrNotExist) {
		fmt.Println("✅ errors.Is works through wrapping too!")
	}

	// ================================
	// 3. errors.As - EXTRACT ERROR TYPE
	// ================================

	fmt.Println("\n=== errors.As ===")

	_, err = os.Open("/root/secret.txt") // will likely fail with permission or not exist
	var pathErr *os.PathError
	if errors.As(err, &pathErr) {
		fmt.Printf("✅ errors.As extracted PathError:\n")
		fmt.Printf("   Op: %s\n", pathErr.Op)
		fmt.Printf("   Path: %s\n", pathErr.Path)
		fmt.Printf("   Err: %s\n", pathErr.Err)
	}

	// ================================
	// 4. CUSTOM ERRORS
	// ================================

	fmt.Println("\n=== Custom Errors ===")

	// Test validation error
	err = validateLogLevel("CRITICAL")
	if err != nil {
		fmt.Println("Error:", err)

		// Extract our custom type
		var valErr *ValidationError
		if errors.As(err, &valErr) {
			fmt.Printf("  Field: %s, Value: %s\n", valErr.Field, valErr.Value)
		}
	}

	err = validateLogLevel("ERROR")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("✅ ERROR is a valid log level")
	}

	// ================================
	// 5. FILE ERROR WITH WRAPPING
	// ================================

	fmt.Println("\n=== File Error with Wrapping ===")

	// Try a missing file
	err = processConfig("nonexistent.log")
	if err != nil {
		fmt.Println("Error:", err)

		// Even through wrapping, we can extract our custom error
		var fileErr *FileError
		if errors.As(err, &fileErr) {
			fmt.Printf("  Extracted FileError - Path: %s, Op: %s\n", fileErr.Path, fileErr.Op)
		}
	}

	// Try a directory (not a file)
	err = processConfig(".")
	if err != nil {
		fmt.Println("Error:", err)
	}

	// ================================
	// 6. MULTIPLE ERROR CHECKS IN ONE FUNCTION
	// ================================

	fmt.Println("\n=== Real World Pattern ===")

	files := []string{"sample.log", "missing.log", ".", "/root/nope.log"}

	for _, f := range files {
		err := processConfig(f)
		if err != nil {
			var fileErr *FileError
			if errors.As(err, &fileErr) {
				switch fileErr.Message {
				case "file not found":
					fmt.Printf("  ❌ %s - not found\n", f)
				case "file is empty":
					fmt.Printf("  ⚠️  %s - empty\n", f)
				case "path is a directory, not a file":
					fmt.Printf("  📁 %s - is a directory\n", f)
				default:
					fmt.Printf("  ❓ %s - %s\n", f, fileErr.Message)
				}
			} else {
				fmt.Printf("  💥 %s - unexpected: %s\n", f, err)
			}
		} else {
			fmt.Printf("  ✅ %s - OK\n", f)
		}
	}
}
