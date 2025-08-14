package detector

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/terradrift-watcher/internal/config"
	"github.com/terradrift-watcher/internal/notifier"
	"github.com/terradrift-watcher/internal/terraform"
)

// Run executes the drift detection process for all configured projects
func Run(cfg *config.Config) error {
	_, err := RunWithResult(cfg)
	return err
}

// RunWithResult executes the drift detection process and returns whether any drift was found
func RunWithResult(cfg *config.Config) (bool, error) {
	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Create a done channel to signal when we're finished
	done := make(chan struct{})

	// Handle signals in a goroutine
	go func() {
		select {
		case sig := <-sigChan:
			log.Printf("INFO: Received signal %v, initiating graceful shutdown...", sig)
			// Clear any auth environment variables before exiting
			clearAuthEnvironment()
			log.Printf("INFO: Cleaned up authentication environment variables")
			os.Exit(130) // Exit code 130 is standard for SIGINT
		case <-done:
			// Normal completion
			return
		}
	}()

	// Ensure we signal completion when function returns
	defer close(done)

	// First, validate that Terraform is installed
	if err := terraform.ValidateTerraformInstallation(); err != nil {
		return false, fmt.Errorf("terraform validation failed: %w", err)
	}

	log.Println("INFO: Starting drift detection process...")

	// Track if any errors occurred and if any drift was detected
	var hasErrors bool
	var driftFound bool

	// Iterate through each project
	for _, project := range cfg.Projects {
		// Skip disabled projects (nil means default true)
		if project.Enabled != nil && (*project.Enabled) == false {
			log.Printf("INFO: Skipping disabled project '%s'", project.Name)
			continue
		}

		log.Printf("INFO: Checking for drift in '%s'...", project.Name)

		// Set authentication environment variables if auth profile is specified
		if project.AuthProfile != "" {
			if err := setAuthEnvironment(cfg, project.AuthProfile); err != nil {
				log.Printf("ERROR: Failed to set auth environment for project '%s': %v", project.Name, err)
				hasErrors = true
				continue
			}
			// Ensure cleanup happens even if we continue or an error occurs
			defer clearAuthEnvironment()
		}

		// Run Terraform drift check
		planOutput, exitCode, err := terraform.CheckDrift(project.Path)

		// Handle the results based on exit code
		switch exitCode {
		case 0:
			// No drift detected
			log.Printf("INFO: No drift detected in '%s'", project.Name)

		case 2:
			// Drift detected - send notifications
			driftFound = true
			log.Printf("ALERT: Drift detected in '%s'! Sending notifications...", project.Name)

			// Extract a summary from the plan output
			summary := terraform.ExtractPlanSummary(planOutput)

			// Always print the drift summary to console
			log.Printf("DRIFT SUMMARY for '%s':", project.Name)
			log.Printf("  %s", strings.ReplaceAll(summary, "\n", "\n  "))

			// Check if verbose mode is enabled
			isVerbose := os.Getenv("TERRADRIFT_VERBOSE") == "true"

			if isVerbose {
				// In verbose mode, show the full plan output
				log.Println("FULL TERRAFORM PLAN OUTPUT:")
				log.Println("=" + strings.Repeat("=", 79))
				for _, line := range strings.Split(planOutput, "\n") {
					log.Println(line)
				}
				log.Println("=" + strings.Repeat("=", 79))
			} else {
				// In normal mode, show a sample of the actual plan output
				planLines := strings.Split(planOutput, "\n")
				relevantLines := []string{}
				for _, line := range planLines {
					// Skip empty lines and certain terraform boilerplate
					trimmed := strings.TrimSpace(line)
					if trimmed != "" && !strings.HasPrefix(trimmed, "Refreshing") &&
						!strings.HasPrefix(trimmed, "Reading...") &&
						!strings.HasPrefix(trimmed, "Read complete") {
						relevantLines = append(relevantLines, line)
						if len(relevantLines) >= 10 {
							break
						}
					}
				}

				if len(relevantLines) > 0 {
					log.Println("DRIFT DETAILS (first 10 relevant lines):")
					for _, line := range relevantLines {
						log.Printf("  %s", line)
					}
					log.Println("  ... (use --verbose flag or run terraform plan manually for full details)")
				}
			}

			// Send notifications to all configured notifiers for this project
			notificationsSent := 0
			for _, notifierName := range project.Notifiers {
				if err := sendNotification(cfg, notifierName, project.Name, summary, planOutput); err != nil {
					log.Printf("ERROR: Failed to send notification via '%s' for project '%s': %v",
						notifierName, project.Name, err)
					hasErrors = true
				} else {
					log.Printf("INFO: Notification sent via '%s' for project '%s'", notifierName, project.Name)
					notificationsSent++
				}
			}

			// If no notifications were sent successfully, ensure the user knows about the drift
			if notificationsSent == 0 && len(project.Notifiers) > 0 {
				log.Printf("WARNING: Drift detected but no notifications were sent successfully!")
			}

		default:
			// Error occurred
			if err != nil {
				log.Printf("ERROR: Failed to check drift for project '%s': %v", project.Name, err)
				log.Printf("ERROR: Terraform output: %s", planOutput)
			} else {
				log.Printf("ERROR: Unexpected exit code %d for project '%s'", exitCode, project.Name)
			}
			hasErrors = true
		}
	}

	log.Println("INFO: Drift detection process completed")

	if hasErrors {
		return driftFound, fmt.Errorf("drift detection completed with errors")
	}

	return driftFound, nil
}

// setAuthEnvironment sets the environment variables for the specified auth profile
func setAuthEnvironment(cfg *config.Config, profileName string) error {
	profile, err := cfg.GetAuthProfile(profileName)
	if err != nil {
		return err
	}

	// Set environment variables based on provider type
	switch profile.Provider {
	case "aws":
		// Set AWS environment variables
		for key, value := range profile.Config {
			switch key {
			case "access_key_id":
				os.Setenv(config.AWSAccessKeyID, value)
			case "secret_access_key":
				os.Setenv(config.AWSSecretAccessKey, value)
			case "session_token":
				os.Setenv(config.AWSSessionToken, value)
			case "region":
				os.Setenv(config.AWSRegion, value)
			default:
				// Set any additional AWS environment variables
				os.Setenv(key, value)
			}
		}

	case "azure":
		// Set Azure environment variables
		for key, value := range profile.Config {
			switch key {
			case "client_id":
				os.Setenv(config.AzureClientID, value)
			case "client_secret":
				os.Setenv(config.AzureClientSecret, value)
			case "subscription_id":
				os.Setenv(config.AzureSubscriptionID, value)
			case "tenant_id":
				os.Setenv(config.AzureTenantID, value)
			default:
				// Set any additional Azure environment variables
				os.Setenv(key, value)
			}
		}

	case "gcp":
		// Set GCP environment variables
		for key, value := range profile.Config {
			// GCP typically uses GOOGLE_APPLICATION_CREDENTIALS pointing to a service account key file
			os.Setenv(key, value)
		}

	default:
		// For unknown providers, just set the config values as-is
		for key, value := range profile.Config {
			os.Setenv(key, value)
		}
	}

	return nil
}

// clearAuthEnvironment clears authentication-related environment variables
func clearAuthEnvironment() {
	// Clear AWS variables
	os.Unsetenv(config.AWSAccessKeyID)
	os.Unsetenv(config.AWSSecretAccessKey)
	os.Unsetenv(config.AWSSessionToken)
	os.Unsetenv(config.AWSRegion)

	// Clear Azure variables
	os.Unsetenv(config.AzureClientID)
	os.Unsetenv(config.AzureClientSecret)
	os.Unsetenv(config.AzureSubscriptionID)
	os.Unsetenv(config.AzureTenantID)

	// Clear GCP variables
	os.Unsetenv("GOOGLE_APPLICATION_CREDENTIALS")
	os.Unsetenv("GOOGLE_CLOUD_PROJECT")
}

// sendNotification sends a notification using the specified notifier
func sendNotification(cfg *config.Config, notifierName string, projectName string, summary string, planOutput string) error {
	notifierCfg, err := cfg.GetNotifier(notifierName)
	if err != nil {
		return err
	}

	// Skip disabled notifiers (nil means default true)
	if notifierCfg.Enabled != nil && (*notifierCfg.Enabled) == false {
		log.Printf("INFO: Skipping disabled notifier '%s'", notifierName)
		return nil
	}

	// Send notification based on type
	switch notifierCfg.Type {
	case "slack":
		webhookURL, ok := notifierCfg.Config[config.SlackWebhookURL]
		if !ok {
			return fmt.Errorf("slack webhook URL not configured for notifier '%s'", notifierName)
		}

		// Use the rich notification format for better visibility with retry logic (3 retries)
		return notifier.SendSlackRichNotificationWithRetry(webhookURL, projectName, summary, planOutput, 3)

	case "teams":
		// TODO: Implement Teams notification
		// For now, we'll just log that Teams is not yet implemented
		log.Printf("WARNING: Teams notifications not yet implemented for notifier '%s'", notifierName)
		return nil

	case "email":
		// TODO: Implement email notification
		// For now, we'll just log that email is not yet implemented
		log.Printf("WARNING: Email notifications not yet implemented for notifier '%s'", notifierName)
		return nil

	default:
		return fmt.Errorf("unknown notifier type '%s' for notifier '%s'", notifierCfg.Type, notifierName)
	}
}
