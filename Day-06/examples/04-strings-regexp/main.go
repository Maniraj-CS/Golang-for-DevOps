package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	// ================================
	// 1. STRINGS PACKAGE - Everyday ops
	// ================================

	fmt.Println("=== strings Package ===")

	logLine := "2024-01-15 08:01:15 ERROR Connection refused: database db-02 (10.0.2.11:5432)"

	fmt.Println("Contains 'ERROR':", strings.Contains(logLine, "ERROR"))
	fmt.Println("HasPrefix '2024':", strings.HasPrefix(logLine, "2024"))
	fmt.Println("HasSuffix '5432)':", strings.HasSuffix(logLine, "5432)"))
	fmt.Println("ToUpper:", strings.ToUpper("error"))
	fmt.Println("ToLower:", strings.ToLower("ERROR"))

	// ================================
	// 2. SPLIT AND FIELDS
	// ================================

	fmt.Println("\n=== Split vs Fields ===")

	// Split by a specific separator
	parts := strings.Split("key=value", "=")
	fmt.Printf("Split by '=': %v (len: %d)\n", parts, len(parts))

	// SplitN - limit the number of parts
	kv := strings.SplitN("name=John Doe=Jr", "=", 2)
	fmt.Printf("SplitN by '=' (max 2): %v\n", kv)

	// Fields - splits by any whitespace (much better for command output!)
	output := "  /dev/disk1    228G   156G    72G    68%   /"
	fields := strings.Fields(output)
	fmt.Printf("Fields: %v\n", fields)
	fmt.Printf("  Device: %s, Size: %s, Used: %s, Avail: %s, Use%%: %s, Mount: %s\n",
		fields[0], fields[1], fields[2], fields[3], fields[4], fields[5])

	// ================================
	// 3. TRIM, REPLACE, JOIN
	// ================================

	fmt.Println("\n=== Trim, Replace, Join ===")

	fmt.Println("TrimSpace:", "'"+strings.TrimSpace("  hello  ")+"'")
	fmt.Println("TrimPrefix:", strings.TrimPrefix("/api/v1/users", "/api"))
	fmt.Println("TrimSuffix:", strings.TrimSuffix("server.log.gz", ".gz"))
	fmt.Println("ReplaceAll:", strings.ReplaceAll("foo-bar-baz", "-", "_"))

	tags := []string{"production", "us-east-1", "web-tier"}
	fmt.Println("Join:", strings.Join(tags, ", "))

	fmt.Println("Repeat:", strings.Repeat("=-", 20))

	// ================================
	// 4. STRINGS.BUILDER (efficient string building)
	// ================================

	fmt.Println("\n=== strings.Builder ===")

	var b strings.Builder
	servers := []string{"web-01", "web-02", "api-01", "db-01"}

	for i, s := range servers {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(s)
	}
	fmt.Println("Built string:", b.String())

	// ================================
	// 5. REGEXP - Basic patterns
	// ================================

	fmt.Println("\n=== Regexp Basics ===")

	// Find IP addresses
	ipRegex := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)

	text := "Server web-01 (10.0.1.10) connected to db-01 (10.0.2.10) via proxy 192.168.1.1"
	ips := ipRegex.FindAllString(text, -1)
	fmt.Println("Found IPs:", ips)

	// ================================
	// 6. REGEXP - Parsing command output
	// ================================

	fmt.Println("\n=== Regexp - Parsing Output ===")

	// Parse percentage values from text
	percentRegex := regexp.MustCompile(`(\d+\.?\d*)%`)

	sysOutput := "CPU: 78.5% used, Memory: 64.2% used, Disk: 91% used"
	matches := percentRegex.FindAllStringSubmatch(sysOutput, -1)

	for _, match := range matches {
		fmt.Printf("  Full match: %-8s  Number only: %s\n", match[0], match[1])
	}

	// ================================
	// 7. REGEXP - Named capture groups
	// ================================

	fmt.Println("\n=== Regexp - Named Groups ===")

	// Parse key=value pairs from log
	kvRegex := regexp.MustCompile(`(?P<key>\w+)=(?P<value>\S+)`)

	logEntry := "level=ERROR host=web-01 port=8080 msg=connection_refused"
	allMatches := kvRegex.FindAllStringSubmatch(logEntry, -1)

	for _, match := range allMatches {
		// match[0] = full match, match[1] = key, match[2] = value
		fmt.Printf("  %s -> %s\n", match[1], match[2])
	}

	// ================================
	// 8. REAL WORLD: Parsing disk usage output
	// ================================

	fmt.Println("\n=== Real World: Parse df output ===")

	// Simulated df -h output
	dfOutput := `Filesystem      Size  Used Avail Use% Mounted on
/dev/sda1       228G  156G   72G  68% /
/dev/sdb1       500G  350G  150G  70% /data
tmpfs           7.8G  1.2G  6.6G  16% /dev/shm`

	lines := strings.Split(dfOutput, "\n")
	fmt.Printf("  %-15s %-6s %-6s %-6s %-5s %s\n", "Device", "Size", "Used", "Avail", "Use%", "Mount")
	fmt.Printf("  %s\n", strings.Repeat("-", 60))

	for i, line := range lines {
		if i == 0 {
			continue // skip header
		}
		f := strings.Fields(line)
		if len(f) >= 6 {
			fmt.Printf("  %-15s %-6s %-6s %-6s %-5s %s\n", f[0], f[1], f[2], f[3], f[4], f[5])
		}
	}
}
