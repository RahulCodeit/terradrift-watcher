package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// configFile holds the path to the configuration file
	configFile string
	
	// version information (can be set during build)
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "terradrift-watcher",
	Short: "A CLI tool to detect configuration drift in Terraform projects",
	Long: `TerraDrift Watcher is a standalone CLI tool that detects configuration 
drift in Terraform projects by comparing the live infrastructure state 
against the code in Git.

It supports multiple cloud providers (AWS, Azure, GCP) and can send 
notifications via Slack, Microsoft Teams, or email when drift is detected.`,
	Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Define persistent flags that will be available to all subcommands
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.yml", 
		"Path to the configuration file")
	
	// Add version template
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "version %s" .Version}}
`)
} 