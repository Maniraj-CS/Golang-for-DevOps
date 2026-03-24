package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Global flags
var (
	Verbose      bool
	OutputFormat string
)

// rootCmd is the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "devctl",
	Short: "DevOps swiss knife CLI",
	Long: `devctl - A DevOps swiss knife CLI tool

Combines health checks, system info, and GitHub operations
into one clean, professional tool.

Built as part of the "Golang for DevOps" course.`,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Persistent flags - available for all subcommands
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringVarP(&OutputFormat, "output", "o", "table", "Output format: table or json")

	// Bind flags to Viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
}

func initConfig() {
	// Look for config file
	viper.SetConfigName("devctl")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.devctl")

	// Read environment variables with DEVCTL_ prefix
	viper.SetEnvPrefix("DEVCTL")
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("health.timeout", 10)
	viper.SetDefault("health.concurrent", 5)
	viper.SetDefault("output", "table")

	// Try to read config file
	if err := viper.ReadInConfig(); err == nil {
		if Verbose {
			fmt.Println("[verbose] Using config file:", viper.ConfigFileUsed())
		}
	}
}
