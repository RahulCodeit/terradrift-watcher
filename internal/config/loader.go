package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads and parses the configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	// Read the YAML file from disk
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	// Expand environment variables in the YAML content
	expandedData := os.ExpandEnv(string(data))

	// Parse the YAML content into the Config struct
	var config Config
	if err := yaml.Unmarshal([]byte(expandedData), &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	// Default enabled fields to true when omitted
	for i := range config.Projects {
		if config.Projects[i].Enabled == nil {
			def := true
			config.Projects[i].Enabled = &def
		}
	}
	for i := range config.Notifiers {
		if config.Notifiers[i].Enabled == nil {
			def := true
			config.Notifiers[i].Enabled = &def
		}
	}

	// Resolve relative project paths against the config file directory
	configDir := filepath.Dir(path)
	for i := range config.Projects {
		p := config.Projects[i].Path
		if p == "" {
			continue
		}
		if !filepath.IsAbs(p) {
			resolved := filepath.Clean(filepath.Join(configDir, p))
			config.Projects[i].Path = resolved
		}
	}

	// Validate the configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// validateConfig performs basic validation on the configuration
func validateConfig(config *Config) error {
	// Check if we have at least one project
	if len(config.Projects) == 0 {
		return fmt.Errorf("no projects defined in configuration")
	}

	// Create maps for quick lookup
	authProfiles := make(map[string]bool)
	for _, profile := range config.AuthProfiles {
		if profile.Name == "" {
			return fmt.Errorf("auth profile found with empty name")
		}
		if profile.Provider == "" {
			return fmt.Errorf("auth profile %s has no provider specified", profile.Name)
		}
		authProfiles[profile.Name] = true
	}

	notifiers := make(map[string]string)
	for _, notifier := range config.Notifiers {
		if notifier.Name == "" {
			return fmt.Errorf("notifier found with empty name")
		}
		if notifier.Type == "" {
			return fmt.Errorf("notifier %s has no type specified", notifier.Name)
		}
		notifiers[notifier.Name] = notifier.Type
	}

	// Validate each project
	for _, project := range config.Projects {
		if project.Name == "" {
			return fmt.Errorf("project found with empty name")
		}
		if project.Path == "" {
			return fmt.Errorf("project %s has no path specified", project.Name)
		}
		// Ensure the path exists
		if _, err := os.Stat(project.Path); err != nil {
			return fmt.Errorf("project %s path not found: %s", project.Name, project.Path)
		}

		// Check if auth profile exists
		if project.AuthProfile != "" && !authProfiles[project.AuthProfile] {
			return fmt.Errorf("project %s references unknown auth profile: %s", project.Name, project.AuthProfile)
		}

		// Check if all referenced notifiers exist
		for _, notifierName := range project.Notifiers {
			if _, ok := notifiers[notifierName]; !ok {
				return fmt.Errorf("project %s references unknown notifier: %s", project.Name, notifierName)
			}
		}
	}

	return nil
}

// GetAuthProfile returns the auth profile with the given name
func (c *Config) GetAuthProfile(name string) (*AuthProfile, error) {
	for _, profile := range c.AuthProfiles {
		if profile.Name == name {
			return &profile, nil
		}
	}
	return nil, fmt.Errorf("auth profile not found: %s", name)
}

// GetNotifier returns the notifier with the given name
func (c *Config) GetNotifier(name string) (*Notifier, error) {
	for _, notifier := range c.Notifiers {
		if notifier.Name == name {
			return &notifier, nil
		}
	}
	return nil, fmt.Errorf("notifier not found: %s", name)
}
