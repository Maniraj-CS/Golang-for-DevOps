package utils

import (
	"fmt"
	"strings"
)

// FormatLevel normalizes a log level string to uppercase
// Exported function - starts with uppercase
func FormatLevel(level string) string {
	return strings.ToUpper(strings.TrimSpace(level))
}

// IsValidLevel checks if the given level is a known log level
// Exported function
func IsValidLevel(level string) bool {
	valid := map[string]bool{
		"INFO": true, "WARN": true, "ERROR": true, "DEBUG": true,
	}
	return valid[FormatLevel(level)]
}

// FormatBytes converts bytes to human readable format
// Exported function
func FormatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}

// repeatChar is unexported - only visible inside this package
// starts with lowercase
func repeatChar(char string, count int) string {
	return strings.Repeat(char, count)
}

// PrintSeparator uses the unexported helper
// Exported function can use unexported helpers
func PrintSeparator(width int) {
	fmt.Println(repeatChar("-", width))
}
