package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// ================================
// TYPES
// ================================

type SystemInfo struct {
	CollectedAt string     `json:"collected_at"`
	Host        HostInfo   `json:"host"`
	CPU         CPUInfo    `json:"cpu"`
	Memory      MemoryInfo `json:"memory"`
	Disk        []DiskInfo `json:"disk"`
	Network     []NetInfo  `json:"network"`
}

type HostInfo struct {
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Uptime   string `json:"uptime"`
}

type CPUInfo struct {
	Cores   int    `json:"cores"`
	LoadAvg string `json:"load_avg"`
}

type MemoryInfo struct {
	Total     string `json:"total"`
	Used      string `json:"used"`
	Available string `json:"available"`
	UsagePercent string `json:"usage_percent"`
}

type DiskInfo struct {
	Filesystem string `json:"filesystem"`
	Size       string `json:"size"`
	Used       string `json:"used"`
	Available  string `json:"available"`
	UsePercent string `json:"use_percent"`
	MountedOn  string `json:"mounted_on"`
}

type NetInfo struct {
	Interface string `json:"interface"`
	IPAddress string `json:"ip_address"`
}

// Config for CLI
type Config struct {
	OutputJSON bool
	Section    string // "all", "host", "cpu", "memory", "disk", "network"
}

// ================================
// CLI PARSING
// ================================

func parseArgs() *Config {
	config := &Config{
		Section: "all",
	}

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--json":
			config.OutputJSON = true
		case "--section":
			if i+1 < len(args) {
				i++
				config.Section = strings.ToLower(args[i])
			}
		case "--help", "-h":
			printUsage()
			os.Exit(0)
		}
	}

	return config
}

func printUsage() {
	fmt.Println("Usage: sysinfo [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --json              Output as JSON")
	fmt.Println("  --section SECTION   Show only: all, host, cpu, memory, disk, network")
	fmt.Println("  --help              Show this help")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run main.go")
	fmt.Println("  go run main.go --json")
	fmt.Println("  go run main.go --section disk")
}

// ================================
// COMMAND RUNNERS (safe execution)
// ================================

// runCommand runs a command and returns trimmed output
func runCommand(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// runShell runs a shell command (supports pipes, etc.)
func runShell(command string) (string, error) {
	out, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// ================================
// COLLECTORS
// ================================

func collectHostInfo() HostInfo {
	info := HostInfo{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	// Hostname
	if hostname, err := os.Hostname(); err == nil {
		info.Hostname = hostname
	}

	// Uptime
	if runtime.GOOS == "darwin" {
		// macOS: use uptime command
		if out, err := runCommand("uptime"); err == nil {
			// Extract just the uptime part
			info.Uptime = parseUptime(out)
		}
	} else {
		// Linux: use uptime -p for pretty output
		if out, err := runCommand("uptime", "-p"); err == nil {
			info.Uptime = out
		} else if out, err := runCommand("uptime"); err == nil {
			info.Uptime = parseUptime(out)
		}
	}

	return info
}

func collectCPUInfo() CPUInfo {
	info := CPUInfo{
		Cores: runtime.NumCPU(),
	}

	// Load average
	if runtime.GOOS == "darwin" {
		if out, err := runShell("sysctl -n vm.loadavg"); err == nil {
			// macOS returns "{ 1.23 0.98 0.87 }"
			info.LoadAvg = strings.Trim(out, "{ }")
			info.LoadAvg = strings.Join(strings.Fields(info.LoadAvg), ", ")
		}
	} else {
		if out, err := runShell("cat /proc/loadavg | awk '{print $1, $2, $3}'"); err == nil {
			fields := strings.Fields(out)
			info.LoadAvg = strings.Join(fields, ", ")
		}
	}

	return info
}

func collectMemoryInfo() MemoryInfo {
	info := MemoryInfo{}

	if runtime.GOOS == "darwin" {
		// macOS: use vm_stat and sysctl
		collectMemoryMacOS(&info)
	} else {
		// Linux: use free -h
		collectMemoryLinux(&info)
	}

	return info
}

func collectMemoryMacOS(info *MemoryInfo) {
	// Get total memory from sysctl
	if out, err := runShell("sysctl -n hw.memsize"); err == nil {
		var totalBytes int64
		fmt.Sscanf(out, "%d", &totalBytes)
		totalGB := float64(totalBytes) / (1024 * 1024 * 1024)
		info.Total = fmt.Sprintf("%.2f GB", totalGB)

		// Get page size and vm_stat for used/free
		if vmOut, err := runCommand("vm_stat"); err == nil {
			pageSize := int64(16384) // default on Apple Silicon
			if ps, err := runShell("sysctl -n hw.pagesize"); err == nil {
				fmt.Sscanf(ps, "%d", &pageSize)
			}

			freePages := extractVMStatValue(vmOut, "Pages free")
			inactivePages := extractVMStatValue(vmOut, "Pages inactive")
			activePages := extractVMStatValue(vmOut, "Pages active")
			wiredPages := extractVMStatValue(vmOut, "Pages wired down")
			compressedPages := extractVMStatValue(vmOut, "Pages occupied by compressor")

			usedBytes := (activePages + wiredPages + compressedPages) * pageSize
			availBytes := (freePages + inactivePages) * pageSize

			usedGB := float64(usedBytes) / (1024 * 1024 * 1024)
			availGB := float64(availBytes) / (1024 * 1024 * 1024)
			usagePercent := (usedGB / totalGB) * 100

			info.Used = fmt.Sprintf("%.2f GB", usedGB)
			info.Available = fmt.Sprintf("%.2f GB", availGB)
			info.UsagePercent = fmt.Sprintf("%.1f%%", usagePercent)
		}
	}
}

func collectMemoryLinux(info *MemoryInfo) {
	if out, err := runShell("free -h | grep Mem"); err == nil {
		fields := strings.Fields(out)
		if len(fields) >= 7 {
			info.Total = fields[1]
			info.Used = fields[2]
			info.Available = fields[6]

			// Calculate percentage from free -b (bytes) for accuracy
			if bOut, err := runShell("free -b | grep Mem"); err == nil {
				bFields := strings.Fields(bOut)
				if len(bFields) >= 7 {
					var total, avail int64
					fmt.Sscanf(bFields[1], "%d", &total)
					fmt.Sscanf(bFields[6], "%d", &avail)
					if total > 0 {
						usage := float64(total-avail) / float64(total) * 100
						info.UsagePercent = fmt.Sprintf("%.1f%%", usage)
					}
				}
			}
		}
	}
}

func collectDiskInfo() []DiskInfo {
	var disks []DiskInfo

	var out string
	var err error

	if runtime.GOOS == "darwin" {
		out, err = runShell("df -h | grep -E '^/dev/'")
	} else {
		out, err = runShell("df -h --output=source,size,used,avail,pcent,target | grep -E '^/dev/'")
	}

	if err != nil || out == "" {
		// Fallback: try without grep filter
		out, _ = runShell("df -h")
	}

	lines := strings.Split(out, "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		// Skip header-like lines
		if fields[0] == "Filesystem" {
			continue
		}

		disk := DiskInfo{
			Filesystem: fields[0],
			Size:       fields[1],
			Used:       fields[2],
			Available:  fields[3],
			UsePercent: fields[4],
			MountedOn:  fields[len(fields)-1], // last field (mount point might have spaces)
		}
		disks = append(disks, disk)
	}

	return disks
}

func collectNetworkInfo() []NetInfo {
	var interfaces []NetInfo

	if runtime.GOOS == "darwin" {
		// macOS: use ifconfig and parse
		out, err := runShell("ifconfig | grep -E '(^\\S|inet )' | grep -v inet6")
		if err == nil {
			interfaces = parseIfconfig(out)
		}
	} else {
		// Linux: try ip addr first, fallback to ifconfig
		out, err := runShell("ip -4 addr show | grep -E '(^\\d|inet )'")
		if err == nil {
			interfaces = parseIPAddr(out)
		} else {
			out, err = runShell("ifconfig | grep -E '(^\\S|inet )' | grep -v inet6")
			if err == nil {
				interfaces = parseIfconfig(out)
			}
		}
	}

	return interfaces
}

// ================================
// PARSERS
// ================================

func parseUptime(raw string) string {
	// Try to extract the "up X days, HH:MM" part
	re := regexp.MustCompile(`up\s+(.+?),\s+\d+\s+user`)
	matches := re.FindStringSubmatch(raw)
	if len(matches) >= 2 {
		return "up " + strings.TrimSpace(matches[1])
	}
	// Fallback: return as-is but trimmed
	if idx := strings.Index(raw, "up"); idx >= 0 {
		result := raw[idx:]
		if cIdx := strings.Index(result, ",  "); cIdx > 0 {
			return strings.TrimSpace(result[:cIdx])
		}
		return strings.TrimSpace(result)
	}
	return raw
}

func extractVMStatValue(vmstat string, key string) int64 {
	for _, line := range strings.Split(vmstat, "\n") {
		if strings.Contains(line, key) {
			re := regexp.MustCompile(`(\d+)`)
			match := re.FindString(line)
			var val int64
			fmt.Sscanf(match, "%d", &val)
			return val
		}
	}
	return 0
}

func parseIfconfig(output string) []NetInfo {
	var interfaces []NetInfo
	var currentIface string

	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Interface line (starts without whitespace in raw, but we trimmed)
		if !strings.HasPrefix(line, "inet") {
			// It's an interface name like "en0: flags=..."
			parts := strings.SplitN(line, ":", 2)
			if len(parts) >= 1 {
				currentIface = strings.TrimSpace(parts[0])
			}
			continue
		}

		// inet line
		if strings.HasPrefix(line, "inet ") && currentIface != "" {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				ip := fields[1]
				interfaces = append(interfaces, NetInfo{
					Interface: currentIface,
					IPAddress: ip,
				})
			}
		}
	}

	return interfaces
}

func parseIPAddr(output string) []NetInfo {
	var interfaces []NetInfo
	var currentIface string

	ifaceRegex := regexp.MustCompile(`^\d+:\s+(\S+):`)
	inetRegex := regexp.MustCompile(`inet\s+(\S+)`)

	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)

		if matches := ifaceRegex.FindStringSubmatch(line); len(matches) >= 2 {
			currentIface = matches[1]
			continue
		}

		if matches := inetRegex.FindStringSubmatch(line); len(matches) >= 2 && currentIface != "" {
			ip := strings.Split(matches[1], "/")[0] // remove CIDR prefix
			interfaces = append(interfaces, NetInfo{
				Interface: currentIface,
				IPAddress: ip,
			})
		}
	}

	return interfaces
}

// ================================
// DISPLAY
// ================================

func printTable(info *SystemInfo, section string) {
	fmt.Println("=== System Info Collector ===")
	fmt.Printf("Collected at: %s\n", info.CollectedAt)

	if section == "all" || section == "host" {
		fmt.Println("\n--- Host ---")
		fmt.Printf("  %-14s: %s\n", "Hostname", info.Host.Hostname)
		fmt.Printf("  %-14s: %s (%s)\n", "OS", info.Host.OS, info.Host.Arch)
		fmt.Printf("  %-14s: %s\n", "Uptime", info.Host.Uptime)
	}

	if section == "all" || section == "cpu" {
		fmt.Println("\n--- CPU ---")
		fmt.Printf("  %-14s: %d\n", "Cores", info.CPU.Cores)
		fmt.Printf("  %-14s: %s\n", "Load Avg", info.CPU.LoadAvg)
	}

	if section == "all" || section == "memory" {
		fmt.Println("\n--- Memory ---")
		fmt.Printf("  %-14s: %s\n", "Total", info.Memory.Total)
		fmt.Printf("  %-14s: %s\n", "Used", info.Memory.Used)
		fmt.Printf("  %-14s: %s\n", "Available", info.Memory.Available)
		fmt.Printf("  %-14s: %s\n", "Usage", info.Memory.UsagePercent)
	}

	if section == "all" || section == "disk" {
		fmt.Println("\n--- Disk ---")
		if len(info.Disk) > 0 {
			fmt.Printf("  %-20s %-6s %-6s %-6s %-5s %s\n",
				"Filesystem", "Size", "Used", "Avail", "Use%", "Mounted")
			for _, d := range info.Disk {
				fmt.Printf("  %-20s %-6s %-6s %-6s %-5s %s\n",
					d.Filesystem, d.Size, d.Used, d.Available, d.UsePercent, d.MountedOn)
			}
		} else {
			fmt.Println("  (no disk info available)")
		}
	}

	if section == "all" || section == "network" {
		fmt.Println("\n--- Network ---")
		if len(info.Network) > 0 {
			fmt.Printf("  %-14s %s\n", "Interface", "IP Address")
			for _, n := range info.Network {
				fmt.Printf("  %-14s %s\n", n.Interface, n.IPAddress)
			}
		} else {
			fmt.Println("  (no network info available)")
		}
	}

	fmt.Println()
}

func printJSON(info *SystemInfo) {
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		fmt.Println("Error generating JSON:", err)
		return
	}
	fmt.Println(string(data))
}

// ================================
// MAIN
// ================================

func main() {
	config := parseArgs()

	// Collect all system info
	info := &SystemInfo{
		CollectedAt: time.Now().Format("2006-01-02 15:04:05"),
		Host:        collectHostInfo(),
		CPU:         collectCPUInfo(),
		Memory:      collectMemoryInfo(),
		Disk:        collectDiskInfo(),
		Network:     collectNetworkInfo(),
	}

	// Display
	if config.OutputJSON {
		printJSON(info)
	} else {
		printTable(info, config.Section)
	}
}
