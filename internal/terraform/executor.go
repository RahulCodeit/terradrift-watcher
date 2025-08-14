package terraform

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CheckDrift runs terraform plan to detect configuration drift
// Returns the plan output, exit code, and any error
// Exit codes:
//   - 0: No changes (no drift)
//   - 1: Error occurred
//   - 2: Changes detected (drift present)
func CheckDrift(projectPath string) (string, int, error) {
	// Validate that the project path exists
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return "", 1, fmt.Errorf("project path does not exist: %s", projectPath)
	}

	// Set up cleanup function for lock files on error
	cleanupLockFiles := func() {
		// Clean up Terraform lock files on failure
		tfLockFile := filepath.Join(projectPath, ".terraform.lock.hcl")
		if err := os.Remove(tfLockFile); err != nil && !os.IsNotExist(err) {
			fmt.Printf("WARNING: Failed to clean up .terraform.lock.hcl: %v\n", err)
		}

		// Also try to clean up any .terraform.tfstate.lock.info files
		tfStateLock := filepath.Join(projectPath, ".terraform.tfstate.lock.info")
		if err := os.Remove(tfStateLock); err != nil && !os.IsNotExist(err) {
			fmt.Printf("WARNING: Failed to clean up .terraform.tfstate.lock.info: %v\n", err)
		}
	}

	// Run terraform init
	initOutput, err := runTerraformInit(projectPath)
	if err != nil {
		cleanupLockFiles()
		return initOutput, 1, fmt.Errorf("terraform init failed: %w", err)
	}

	// Run terraform plan with detailed exit code
	planOutput, exitCode, err := runTerraformPlan(projectPath)
	if err != nil && exitCode != 2 {
		// Exit code 2 is expected when drift is detected, so we don't treat it as an error
		cleanupLockFiles()
		return planOutput, exitCode, fmt.Errorf("terraform plan failed: %w", err)
	}

	return planOutput, exitCode, nil
}

// buildEnv returns the environment to use for terraform commands
func buildEnv() []string {
	env := os.Environ()
	// Ensure automation-friendly output
	if os.Getenv("TF_IN_AUTOMATION") == "" {
		env = append(env, "TF_IN_AUTOMATION=true")
	}
	return env
}

// runTerraformInit executes terraform init command
func runTerraformInit(projectPath string) (string, error) {
	// Clean up any existing lock files first
	lockFile := filepath.Join(projectPath, ".terraform.lock.hcl")
	if _, err := os.Stat(lockFile); err == nil {
		// Lock file exists, try to remove it
		if err := os.Remove(lockFile); err != nil {
			// Log warning but continue
			fmt.Printf("WARNING: Could not remove existing lock file: %v\n", err)
		}
	}

	cmd := exec.Command("terraform", "init", "-input=false", "-no-color", "-upgrade=false")
	cmd.Dir = projectPath
	cmd.Env = buildEnv()

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()

	if err != nil {
		// Check for common backend initialization errors
		if strings.Contains(output, "Error loading backend config") ||
			strings.Contains(output, "Backend initialization required") ||
			strings.Contains(output, "Error configuring the backend") {
			return output, fmt.Errorf("backend initialization failed - may need manual intervention: %s", output)
		}
		if strings.Contains(output, "Could not load plugin") ||
			strings.Contains(output, "Provider produced inconsistent") {
			return output, fmt.Errorf("provider initialization failed - check provider versions: %s", output)
		}
		return output, fmt.Errorf("terraform init failed: %s", output)
	}

	return output, nil
}

// runTerraformPlan executes terraform plan command with detailed exit code
func runTerraformPlan(projectPath string) (string, int, error) {
	cmd := exec.Command("terraform", "plan", "-input=false", "-no-color", "-detailed-exitcode")
	cmd.Dir = projectPath
	cmd.Env = buildEnv()

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()

	// Get the exit code
	exitCode := 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	} else if err != nil {
		// If there's an error but it's not an ExitError, something went wrong
		return output, 1, fmt.Errorf("failed to execute terraform plan: %w", err)
	}

	// Exit code 2 means changes were detected (drift), which is not an error condition
	if exitCode == 2 {
		return output, exitCode, nil
	}

	// Any other non-zero exit code is an error
	if exitCode != 0 {
		return output, exitCode, fmt.Errorf("terraform plan failed with exit code %d: %s", exitCode, output)
	}

	return output, exitCode, nil
}

// ExtractPlanSummary extracts a summary from the terraform plan output
func ExtractPlanSummary(planOutput string) string {
	lines := strings.Split(planOutput, "\n")
	summary := []string{}
	resourceChanges := []string{}
	captureChanges := false

	// Look for the plan summary line and resource changes
	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Capture the main plan summary
		if strings.Contains(line, "Plan:") ||
			strings.Contains(line, "No changes") ||
			strings.Contains(line, "to add") ||
			strings.Contains(line, "to change") ||
			strings.Contains(line, "to destroy") {
			summary = append(summary, trimmedLine)
		}

		// Capture resource changes
		if strings.Contains(line, "Terraform will perform the following actions:") {
			captureChanges = true
			continue
		}

		// Look for resource change indicators
		if captureChanges && len(resourceChanges) < 10 {
			if strings.HasPrefix(trimmedLine, "#") ||
				strings.HasPrefix(trimmedLine, "~") ||
				strings.HasPrefix(trimmedLine, "+") ||
				strings.HasPrefix(trimmedLine, "-") {
				// This is likely a resource change line
				resourceChanges = append(resourceChanges, trimmedLine)
			}
		}

		// Also capture lines that show what resources will be modified
		if strings.Contains(line, "will be") && (strings.Contains(line, "created") ||
			strings.Contains(line, "destroyed") || strings.Contains(line, "updated") ||
			strings.Contains(line, "replaced")) {
			if len(resourceChanges) < 10 {
				resourceChanges = append(resourceChanges, trimmedLine)
			}
		}

		// Stop capturing after we hit the next section
		if captureChanges && (strings.Contains(line, "─────────────") || i > len(lines)-10) {
			captureChanges = false
		}
	}

	// Build the final summary
	var result strings.Builder

	if len(summary) > 0 {
		result.WriteString(strings.Join(summary, "\n"))
	} else {
		result.WriteString("Drift detected in Terraform configuration")
	}

	if len(resourceChanges) > 0 {
		result.WriteString("\n\nResource Changes Detected:")
		for _, change := range resourceChanges {
			result.WriteString("\n  " + change)
		}
		if len(resourceChanges) == 10 {
			result.WriteString("\n  ... (more changes, see full plan for details)")
		}
	}

	return result.String()
}

// ValidateTerraformInstallation checks if terraform is installed and accessible
func ValidateTerraformInstallation() error {
	cmd := exec.Command("terraform", "version")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("terraform is not installed or not in PATH: %w", err)
	}

	return nil
}
