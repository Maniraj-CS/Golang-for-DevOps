package display

// ================================
// Terminal Color Helpers
// ================================
// Simple ANSI color codes for terminal output.
// We keep it simple - no external libraries needed!
// These work on macOS, Linux, and modern Windows terminals.

import "fmt"

// ANSI color codes
const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
	Gray   = "\033[90m"

	BoldRed    = "\033[1;31m"
	BoldGreen  = "\033[1;32m"
	BoldYellow = "\033[1;33m"
	BoldBlue   = "\033[1;34m"
	BoldCyan   = "\033[1;36m"
	BoldWhite  = "\033[1;37m"
)

// Colorize wraps text with an ANSI color code
func Colorize(color, text string) string {
	return color + text + Reset
}

// StatusIcon returns a colored emoji icon based on health status
func StatusIcon(health string) string {
	switch health {
	case "healthy":
		return "🟢"
	case "warning":
		return "🟡"
	case "critical":
		return "🔴"
	default:
		return "⚪"
	}
}

// DeploymentIcon returns an icon for deployment health
func DeploymentIcon(healthy bool) string {
	if healthy {
		return "✅"
	}
	return "⚠️ "
}

// PrintHeader prints a colored section header
func PrintHeader(title string) {
	fmt.Printf("\n%s--- %s ---%s\n\n", BoldCyan, title, Reset)
}

// PrintMainHeader prints the main title
func PrintMainHeader(context string) {
	fmt.Printf("\n%s=== Cluster Status ===%s\n", BoldWhite, Reset)
	fmt.Printf("Context: %s%s%s\n", BoldGreen, context, Reset)
}

// PrintSeparator prints a visual separator
func PrintSeparator() {
	fmt.Println()
}
