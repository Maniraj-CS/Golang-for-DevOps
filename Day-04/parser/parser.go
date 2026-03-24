package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"log-parser/apperrors"
)

// LogEntry represents a single parsed log line
type LogEntry struct {
	Timestamp string
	Level     string
	Message   string
	Raw       string // the original line
}

// ParseResult holds the results of parsing a log file
type ParseResult struct {
	Entries    []LogEntry
	LevelCount map[string]int
	TotalLines int
}

// ValidLevels are the log levels we understand
var ValidLevels = map[string]bool{
	"INFO": true, "WARN": true, "ERROR": true, "DEBUG": true,
}

// ValidateFile checks if the log file is valid before parsing
func ValidateFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &apperrors.FileError{
				Path:    path,
				Op:      "open",
				Message: "file not found",
			}
		}
		return fmt.Errorf("validate file: %w", err)
	}

	if info.IsDir() {
		return &apperrors.FileError{
			Path:    path,
			Op:      "open",
			Message: "path is a directory, not a file",
		}
	}

	if info.Size() == 0 {
		return &apperrors.FileError{
			Path:    path,
			Op:      "read",
			Message: "file is empty (0 bytes)",
		}
	}

	return nil
}

// ValidateLevel checks if the filter level is valid
func ValidateLevel(level string) error {
	if level == "" || level == "ALL" {
		return nil // no filter = valid
	}

	upper := strings.ToUpper(level)
	if !ValidLevels[upper] {
		return &apperrors.ConfigError{
			Flag:    "level",
			Value:   level,
			Message: "must be one of: INFO, WARN, ERROR, DEBUG, ALL",
		}
	}
	return nil
}

// ParseFile reads a log file and returns parsed results
func ParseFile(path string) (*ParseResult, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}
	defer file.Close()

	result := &ParseResult{
		Entries:    []LogEntry{},
		LevelCount: map[string]int{},
	}

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if strings.TrimSpace(line) == "" {
			continue // skip empty lines
		}

		entry := parseLine(line, lineNum)
		result.Entries = append(result.Entries, entry)
		result.LevelCount[entry.Level]++
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read log file at line %d: %w", lineNum, err)
	}

	result.TotalLines = lineNum
	return result, nil
}

// parseLine extracts timestamp, level, and message from a log line
// Format expected: "2024-01-15 08:00:01 INFO some message here"
func parseLine(line string, lineNum int) LogEntry {
	entry := LogEntry{Raw: line}

	// Try to extract timestamp (first 19 chars: "2024-01-15 08:00:01")
	if len(line) >= 19 {
		entry.Timestamp = line[:19]
		rest := strings.TrimSpace(line[19:])

		// Try to extract level (first word after timestamp)
		parts := strings.SplitN(rest, " ", 2)
		if len(parts) >= 1 {
			level := strings.ToUpper(parts[0])
			if ValidLevels[level] {
				entry.Level = level
				if len(parts) >= 2 {
					entry.Message = parts[1]
				}
			} else {
				entry.Level = "UNKNOWN"
				entry.Message = rest
			}
		}
	} else {
		entry.Level = "UNKNOWN"
		entry.Message = line
	}

	return entry
}

// FilterByLevel returns only entries matching the given level
func FilterByLevel(entries []LogEntry, level string) []LogEntry {
	if level == "" || level == "ALL" {
		return entries
	}

	upper := strings.ToUpper(level)
	filtered := []LogEntry{}

	for _, e := range entries {
		if e.Level == upper {
			filtered = append(filtered, e)
		}
	}

	return filtered
}

// WriteToFile writes log entries to an output file
func WriteToFile(entries []LogEntry, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return &apperrors.FileError{
			Path:    path,
			Op:      "write",
			Message: fmt.Sprintf("cannot create file: %s", err),
		}
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, e := range entries {
		_, err := writer.WriteString(e.Raw + "\n")
		if err != nil {
			return &apperrors.FileError{
				Path:    path,
				Op:      "write",
				Message: fmt.Sprintf("write failed: %s", err),
			}
		}
	}

	err = writer.Flush()
	if err != nil {
		return &apperrors.FileError{
			Path:    path,
			Op:      "write",
			Message: fmt.Sprintf("flush failed: %s", err),
		}
	}

	return nil
}
