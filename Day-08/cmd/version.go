package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// Set at build time with -ldflags
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		if OutputFormat == "json" {
			fmt.Printf(`{"version":"%s","build_time":"%s","git_commit":"%s","go_version":"%s","os_arch":"%s/%s"}`+"\n",
				Version, BuildTime, GitCommit, runtime.Version(), runtime.GOOS, runtime.GOARCH)
		} else {
			fmt.Println("devctl - DevOps Swiss Knife CLI")
			fmt.Println()
			fmt.Printf("Version:    %s\n", Version)
			fmt.Printf("Build Time: %s\n", BuildTime)
			fmt.Printf("Git Commit: %s\n", GitCommit)
			fmt.Printf("Go Version: %s\n", runtime.Version())
			fmt.Printf("OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
