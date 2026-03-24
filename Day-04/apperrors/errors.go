package apperrors

import "fmt"

// FileError represents errors related to file operations
type FileError struct {
	Path    string
	Op      string // operation: open, read, write
	Message string
}

func (e *FileError) Error() string {
	return fmt.Sprintf("%s %s: %s", e.Op, e.Path, e.Message)
}

// ParseError represents errors when parsing log lines
type ParseError struct {
	Line    int
	Content string
	Message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error at line %d: %s", e.Line, e.Message)
}

// ConfigError represents errors with CLI arguments or config
type ConfigError struct {
	Flag    string
	Value   string
	Message string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("config error: flag --%s=%q: %s", e.Flag, e.Value, e.Message)
}
