package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// ================================
// TYPES
// ================================

type SystemInfo struct {
	Host    HostInfo   `json:"host"`
	CPU     CPUInfo    `json:"cpu"`
	Memory  MemInfo    `json:"memory"`
	Disk    []DiskInfo `json:"disk"`
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

type MemInfo struct {
	Total        string `json:"total"`
	Used         string `json:"used"`
	Available    string `json:"available"`
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

// ================================
// COBRA COMMAND
// ================================

var sysSection string

var sysinfoCmd = &cobra.Command{
	Use:   "sysinfo",
	Short: "Show system information",
	Long: `Collect and display system information like hostname, CPU, memory, and disk usage.

Examples:
  devctl sysinfo
  devctl sysinfo --json
  devctl sysinfo --section disk
  devctl sysinfo --section memory -o json`,
	Run: runSysinfo,
}

func init() {
	rootCmd.AddCommand(sysinfoCmd)
	sysinfoCmd.Flags().StringVar(&sysSection, "section", "all", "Show only: all, host, cpu, memory, disk")
}

func runSysinfo(cmd *cobra.Command, args []string) {
	info := collectSystemInfo()

	if OutputFormat == "json" {
		printSysinfoJSON(info, sysSection)
	} else {
		printSysinfoTable(info, sysSection)
	}
}

// ================================
// COLLECTORS
// ================================

func collectSystemInfo() SystemInfo {
	return SystemInfo{
		Host:   collectHost(),
		CPU:    collectCPU(),
		Memory: collectMem(),
		Disk:   collectDisk(),
	}
}

func collectHost() HostInfo {
	info := HostInfo{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	if hostname, err := os.Hostname(); err == nil {
		info.Hostname = hostname
	}

	if runtime.GOOS == "darwin" {
		if out, err := runCmd("uptime"); err == nil {
			info.Uptime = parseUptimeStr(out)
		}
	} else {
		if out, err := runCmd("uptime", "-p"); err == nil {
			info.Uptime = out
		}
	}

	return info
}

func collectCPU() CPUInfo {
	info := CPUInfo{Cores: runtime.NumCPU()}

	if runtime.GOOS == "darwin" {
		if out, err := runShellCmd("sysctl -n vm.loadavg"); err == nil {
			info.LoadAvg = strings.Trim(out, "{ }")
			info.LoadAvg = strings.Join(strings.Fields(info.LoadAvg), ", ")
		}
	} else {
		if out, err := runShellCmd("cat /proc/loadavg | awk '{print $1, $2, $3}'"); err == nil {
			info.LoadAvg = strings.Join(strings.Fields(out), ", ")
		}
	}

	return info
}

func collectMem() MemInfo {
	info := MemInfo{}

	if runtime.GOOS == "darwin" {
		// macOS - use sysctl for total mem and vm_stat for usage
		if out, err := runShellCmd("sysctl -n hw.memsize"); err == nil {
			bytes := strings.TrimSpace(out)
			info.Total = formatBytesStr(bytes)
		}
		if out, err := runShellCmd("vm_stat | head -5"); err == nil {
			info.Used, info.Available, info.UsagePercent = parseVmStat(out)
		}
	} else {
		// Linux - use free -h
		if out, err := runShellCmd("free -h | grep Mem"); err == nil {
			fields := strings.Fields(out)
			if len(fields) >= 7 {
				info.Total = fields[1]
				info.Used = fields[2]
				info.Available = fields[6]
			}
		}
		if out, err := runShellCmd("free | grep Mem | awk '{printf \"%.1f%%\", $3/$2*100}'"); err == nil {
			info.UsagePercent = out
		}
	}

	return info
}

func collectDisk() []DiskInfo {
	var disks []DiskInfo

	out, err := runShellCmd("df -h | grep -E '^/dev/' | head -5")
	if err != nil {
		return disks
	}

	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}
		disks = append(disks, DiskInfo{
			Filesystem: fields[0],
			Size:       fields[1],
			Used:       fields[2],
			Available:  fields[3],
			UsePercent: fields[4],
			MountedOn:  fields[len(fields)-1],
		})
	}

	return disks
}

// ================================
// DISPLAY
// ================================

func printSysinfoTable(info SystemInfo, section string) {
	fmt.Println("\n=== System Info ===")

	if section == "all" || section == "host" {
		fmt.Println("\nHost")
		fmt.Printf("  Hostname:   %s\n", info.Host.Hostname)
		fmt.Printf("  OS:         %s\n", info.Host.OS)
		fmt.Printf("  Arch:       %s\n", info.Host.Arch)
		fmt.Printf("  Uptime:     %s\n", info.Host.Uptime)
	}

	if section == "all" || section == "cpu" {
		fmt.Println("\nCPU")
		fmt.Printf("  Cores:      %d\n", info.CPU.Cores)
		fmt.Printf("  Load Avg:   %s\n", info.CPU.LoadAvg)
	}

	if section == "all" || section == "memory" {
		fmt.Println("\nMemory")
		fmt.Printf("  Total:      %s\n", info.Memory.Total)
		fmt.Printf("  Used:       %s\n", info.Memory.Used)
		fmt.Printf("  Available:  %s\n", info.Memory.Available)
		fmt.Printf("  Usage:      %s\n", info.Memory.UsagePercent)
	}

	if section == "all" || section == "disk" {
		fmt.Println("\nDisk")
		for _, d := range info.Disk {
			fmt.Printf("  %-12s %s total, %s used, %s free (%s)\n",
				d.MountedOn, d.Size, d.Used, d.Available, d.UsePercent)
		}
		if len(info.Disk) == 0 {
			fmt.Println("  (no disk info available)")
		}
	}

	fmt.Println()
}

func printSysinfoJSON(info SystemInfo, section string) {
	var output any

	switch section {
	case "host":
		output = info.Host
	case "cpu":
		output = info.CPU
	case "memory":
		output = info.Memory
	case "disk":
		output = info.Disk
	default:
		output = info
	}

	data, _ := json.MarshalIndent(output, "", "  ")
	fmt.Println(string(data))
}

// ================================
// HELPERS
// ================================

func runCmd(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func runShellCmd(command string) (string, error) {
	out, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func parseUptimeStr(raw string) string {
	// macOS uptime output: "17:30  up 5 days,  3:22, 4 users, load averages: ..."
	// Extract just the "up ..." part
	idx := strings.Index(raw, "up ")
	if idx == -1 {
		return raw
	}
	rest := raw[idx+3:]
	// Cut at "load" or end
	if li := strings.Index(rest, " load"); li != -1 {
		rest = rest[:li]
	}
	// Remove trailing ", N users"
	if ui := strings.Index(rest, " user"); ui != -1 {
		rest = rest[:ui]
		// Remove the number before "user"
		rest = strings.TrimSpace(rest)
		if li := strings.LastIndex(rest, ","); li != -1 {
			rest = strings.TrimSpace(rest[:li])
		}
	}
	return strings.TrimRight(rest, ", ")
}

func formatBytesStr(bytesStr string) string {
	var b float64
	fmt.Sscanf(strings.TrimSpace(bytesStr), "%f", &b)
	gb := b / (1024 * 1024 * 1024)
	if gb >= 1 {
		return fmt.Sprintf("%.0f GB", gb)
	}
	mb := b / (1024 * 1024)
	return fmt.Sprintf("%.0f MB", mb)
}

func parseVmStat(raw string) (used, available, percent string) {
	// Very rough parsing of vm_stat output for macOS
	lines := strings.Split(raw, "\n")
	var free, active, inactive, wired float64
	pageSize := 16384.0 // default for Apple Silicon

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		valStr := strings.TrimRight(fields[len(fields)-1], ".")
		var val float64
		fmt.Sscanf(valStr, "%f", &val)

		switch {
		case strings.Contains(line, "Pages free"):
			free = val
		case strings.Contains(line, "Pages active"):
			active = val
		case strings.Contains(line, "Pages inactive"):
			inactive = val
		case strings.Contains(line, "Pages wired"):
			wired = val
		}
	}

	usedBytes := (active + wired) * pageSize
	freeBytes := (free + inactive) * pageSize
	totalBytes := usedBytes + freeBytes

	usedGB := usedBytes / (1024 * 1024 * 1024)
	freeGB := freeBytes / (1024 * 1024 * 1024)

	pct := 0.0
	if totalBytes > 0 {
		pct = (usedBytes / totalBytes) * 100
	}

	return fmt.Sprintf("%.1f GB", usedGB), fmt.Sprintf("%.1f GB", freeGB), fmt.Sprintf("%.1f%%", pct)
}
