package config

// Config represents the root configuration structure
type Config struct {
	Projects      []Project     `yaml:"projects"`
	AuthProfiles  []AuthProfile `yaml:"auth_profiles"`
	Notifiers     []Notifier    `yaml:"notifiers"`
	CheckInterval string        `yaml:"check_interval,omitempty"`
}

// Project represents a Terraform project to monitor
type Project struct {
	Name        string   `yaml:"name"`
	Path        string   `yaml:"path"`
	AuthProfile string   `yaml:"auth_profile"`
	Notifiers   []string `yaml:"notifiers"`
	Enabled     *bool    `yaml:"enabled,omitempty"`
}

// AuthProfile represents authentication credentials for cloud providers
type AuthProfile struct {
	Name     string            `yaml:"name"`
	Provider string            `yaml:"provider"` // aws, azure, gcp
	Config   map[string]string `yaml:"config"`   // Provider-specific config
}

// Notifier represents a notification channel configuration
type Notifier struct {
	Name    string            `yaml:"name"`
	Type    string            `yaml:"type"` // slack, teams, email
	Config  map[string]string `yaml:"config"`
	Enabled *bool             `yaml:"enabled,omitempty"`
}

// AWS-specific auth config keys
const (
	AWSAccessKeyID     = "AWS_ACCESS_KEY_ID"
	AWSSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	AWSSessionToken    = "AWS_SESSION_TOKEN"
	AWSRegion          = "AWS_DEFAULT_REGION"
)

// Azure-specific auth config keys
const (
	AzureClientID       = "ARM_CLIENT_ID"
	AzureClientSecret   = "ARM_CLIENT_SECRET"
	AzureSubscriptionID = "ARM_SUBSCRIPTION_ID"
	AzureTenantID       = "ARM_TENANT_ID"
)

// Notification config keys
const (
	SlackWebhookURL = "webhook_url"
	TeamsWebhookURL = "webhook_url"
	EmailSMTPHost   = "smtp_host"
	EmailSMTPPort   = "smtp_port"
	EmailFrom       = "from"
	EmailTo         = "to"
)
