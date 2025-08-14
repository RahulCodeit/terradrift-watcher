package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/terradrift-watcher/internal/config"
	"github.com/terradrift-watcher/internal/detector"
	"github.com/terradrift-watcher/internal/lock"
)

var verbose bool
var failOnDrift bool
var forceLock bool

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run drift detection for configured Terraform projects",
	Long: `Run executes the drift detection process for all configured Terraform projects.

It will:
1. Load the configuration from the specified config file
2. Iterate through each enabled project
3. Run 'terraform plan' to detect drift
4. Send notifications if drift is detected

Example:
  terradrift-watcher run --config config.yml
  terradrift-watcher run --config config.yml --verbose`,
	RunE: runDriftDetection,
}

func init() {
	// Add the run command to the root command
	rootCmd.AddCommand(runCmd)

	// Add verbose flag
	runCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show full terraform plan output")

	// Add fail-on-drift flag
	runCmd.Flags().BoolVar(&failOnDrift, "fail-on-drift", false, "Exit with code 2 if drift is detected")

	// Add force flag
	runCmd.Flags().BoolVar(&forceLock, "force", false, "Force release any existing lock and proceed")
}

// runDriftDetection is the main execution function for the run command
func runDriftDetection(cmd *cobra.Command, args []string) error {
	// Create and acquire lock
	fileLock := lock.NewFileLock("")

	if forceLock {
		// Force release any existing lock
		if err := fileLock.ForceRelease(); err != nil {
			log.Printf("WARNING: Failed to force release lock: %v", err)
		}
	}

	// Try to acquire the lock
	if err := fileLock.Acquire(); err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}
	defer func() {
		if err := fileLock.Release(); err != nil {
			log.Printf("WARNING: Failed to release lock: %v", err)
		}
	}()

	log.Printf("INFO: Loading configuration from %s", configFile)

	// Set verbose mode in environment for detector to use
	if verbose {
		os.Setenv("TERRADRIFT_VERBOSE", "true")
		log.Println("INFO: Verbose mode enabled - will show full plan output")
	}

	// Load the configuration
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	log.Printf("INFO: Configuration loaded successfully")
	log.Printf("INFO: Found %d projects, %d auth profiles, and %d notifiers",
		len(cfg.Projects), len(cfg.AuthProfiles), len(cfg.Notifiers))

	// Run the drift detection process
	driftFound, runErr := detector.RunWithResult(cfg)
	if runErr != nil {
		return fmt.Errorf("drift detection failed: %w", runErr)
	}

	if driftFound && failOnDrift {
		// Return an error that preserves exit code 2 via Cobra
		// Cobra will print this error; keep it concise
		return fmt.Errorf("drift detected (exiting with code 2)")
	}

	return nil
}
