package main

// ================================
// EXAMPLE 02 - Flags and Viper Config
// ================================
// This shows local flags, persistent flags, and Viper config file reading.
// Run: go run main.go
// Run: go run main.go --verbose deploy --env production
// Run: go run main.go deploy --env staging --replicas 5
// Run: go run main.go status --service web

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Global flag - available to ALL commands
var verbose bool

func main() {
	// ================================
	// 1. ROOT COMMAND WITH PERSISTENT FLAG
	// ================================

	rootCmd := &cobra.Command{
		Use:   "deployer",
		Short: "A deployment tool demo",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// This runs before ANY subcommand
			initConfig()
			if verbose {
				fmt.Println("[verbose] Config loaded")
				fmt.Println("[verbose] Config file used:", viper.ConfigFileUsed())
			}
		},
	}

	// Persistent flag - works for root AND all subcommands
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	// Bind the flag to Viper (so config file can set it too)
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	// ================================
	// 2. DEPLOY SUBCOMMAND (local flags)
	// ================================

	var env string
	var replicas int
	var dryRun bool

	deployCmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy an application",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("=== Deploying ===")
			fmt.Printf("Environment: %s\n", env)
			fmt.Printf("Replicas: %d\n", replicas)
			fmt.Printf("Dry run: %v\n", dryRun)

			if verbose {
				fmt.Println("[verbose] Using image:", viper.GetString("deploy.image"))
				fmt.Println("[verbose] Timeout:", viper.GetInt("deploy.timeout"), "seconds")
			}

			if dryRun {
				fmt.Println("\n(dry run - nothing was actually deployed)")
			} else {
				fmt.Println("\nDeployment started!")
			}
		},
	}

	// Local flags - only work with "deploy" command
	deployCmd.Flags().StringVarP(&env, "env", "e", "development", "Target environment")
	deployCmd.Flags().IntVarP(&replicas, "replicas", "r", 3, "Number of replicas")
	deployCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Simulate deployment without making changes")

	// ================================
	// 3. STATUS SUBCOMMAND
	// ================================

	var service string

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Check deployment status",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("=== Deployment Status ===")
			if service != "" {
				fmt.Printf("Service: %s\n", service)
				fmt.Printf("Status: running\n")
				fmt.Printf("Replicas: 3/3\n")
			} else {
				// Show all services from config
				services := viper.GetStringSlice("services")
				if len(services) > 0 {
					fmt.Println("Services from config:")
					for _, s := range services {
						fmt.Printf("  - %s: running\n", s)
					}
				} else {
					fmt.Println("No services configured. Use --service flag or add to config file.")
				}
			}
		},
	}

	statusCmd.Flags().StringVarP(&service, "service", "s", "", "Service name to check")

	// ================================
	// 4. CONFIG SUBCOMMAND
	// ================================

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Show current configuration",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("=== Current Config ===")
			fmt.Println()

			settings := viper.AllSettings()
			if len(settings) == 0 {
				fmt.Println("No config file found. Using defaults.")
				fmt.Println("Create a 'deployer.yaml' file to set config.")
				return
			}

			for key, value := range settings {
				fmt.Printf("  %s: %v\n", key, value)
			}
		},
	}

	// ================================
	// 5. WIRE IT UP
	// ================================

	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(configCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// initConfig reads the config file using Viper
func initConfig() {
	viper.SetConfigName("deployer") // file name without extension
	viper.SetConfigType("yaml")     // file format
	viper.AddConfigPath(".")        // look in current directory

	// Set defaults
	viper.SetDefault("deploy.image", "myapp:latest")
	viper.SetDefault("deploy.timeout", 60)
	viper.SetDefault("verbose", false)

	// Read environment variables with DEPLOYER_ prefix
	// DEPLOYER_VERBOSE=true would set "verbose"
	viper.SetEnvPrefix("DEPLOYER")
	viper.AutomaticEnv()

	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if verbose {
				fmt.Println("[verbose] No config file found, using defaults")
			}
		}
	}
}
