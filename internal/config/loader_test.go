package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file for testing
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-config.yml")

	// Create a real project directory
	projectDir := filepath.Join(tempDir, "test", "path")
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("Failed to create project dir: %v", err)
	}

	configContent := fmt.Sprintf(`
auth_profiles:
  - name: test-aws
    provider: aws
    config:
      access_key_id: test-key
      secret_access_key: test-secret
      region: us-east-1

notifiers:
  - name: test-slack
    type: slack
    config:
      webhook_url: https://hooks.slack.com/test
    enabled: true

projects:
  - name: test-project
    path: '%s'
    auth_profile: test-aws
    notifiers:
      - test-slack
    enabled: true
`, projectDir)

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Test loading the configuration
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify the configuration was loaded correctly
	if len(config.AuthProfiles) != 1 {
		t.Errorf("Expected 1 auth profile, got %d", len(config.AuthProfiles))
	}

	if len(config.Notifiers) != 1 {
		t.Errorf("Expected 1 notifier, got %d", len(config.Notifiers))
	}

	if len(config.Projects) != 1 {
		t.Errorf("Expected 1 project, got %d", len(config.Projects))
	}

	// Test specific values
	if config.AuthProfiles[0].Name != "test-aws" {
		t.Errorf("Expected auth profile name 'test-aws', got '%s'", config.AuthProfiles[0].Name)
	}

	if config.Projects[0].Name != "test-project" {
		t.Errorf("Expected project name 'test-project', got '%s'", config.Projects[0].Name)
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("/non/existent/file.yml")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestGetAuthProfile(t *testing.T) {
	config := &Config{
		AuthProfiles: []AuthProfile{
			{
				Name:     "test-profile",
				Provider: "aws",
				Config:   map[string]string{"key": "value"},
			},
		},
	}

	// Test existing profile
	profile, err := config.GetAuthProfile("test-profile")
	if err != nil {
		t.Errorf("Failed to get existing profile: %v", err)
	}
	if profile.Name != "test-profile" {
		t.Errorf("Expected profile name 'test-profile', got '%s'", profile.Name)
	}

	// Test non-existing profile
	_, err = config.GetAuthProfile("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent profile, got nil")
	}
}

func TestGetNotifier(t *testing.T) {
	trueVal := true
	config := &Config{
		Notifiers: []Notifier{
			{
				Name:    "test-notifier",
				Type:    "slack",
				Config:  map[string]string{"webhook_url": "https://test.com"},
				Enabled: &trueVal,
			},
		},
	}

	// Test existing notifier
	notifier, err := config.GetNotifier("test-notifier")
	if err != nil {
		t.Errorf("Failed to get existing notifier: %v", err)
	}
	if notifier.Name != "test-notifier" {
		t.Errorf("Expected notifier name 'test-notifier', got '%s'", notifier.Name)
	}

	// Test non-existing notifier
	_, err = config.GetNotifier("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent notifier, got nil")
	}
}
