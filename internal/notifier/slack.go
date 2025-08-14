package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// SlackMessage represents a basic Slack webhook message
type SlackMessage struct {
	Text        string       `json:"text"`
	Username    string       `json:"username,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

// Attachment represents a Slack message attachment
type Attachment struct {
	Color      string  `json:"color,omitempty"`
	Title      string  `json:"title,omitempty"`
	Text       string  `json:"text,omitempty"`
	Footer     string  `json:"footer,omitempty"`
	FooterIcon string  `json:"footer_icon,omitempty"`
	Timestamp  int64   `json:"ts,omitempty"`
	Fields     []Field `json:"fields,omitempty"`
}

// Field represents a field in a Slack attachment
type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// SendSlackNotification sends a notification to a Slack webhook
func SendSlackNotification(webhookURL string, message string) error {
	if webhookURL == "" {
		return fmt.Errorf("webhook URL is empty")
	}

	if message == "" {
		return fmt.Errorf("message is empty")
	}

	// Create a simple Slack message
	slackMsg := SlackMessage{
		Text:      message,
		Username:  "TerraDrift Watcher",
		IconEmoji: ":warning:",
	}

	// Marshal the message to JSON
	jsonData, err := json.Marshal(slackMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal Slack message: %w", err)
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create the request
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send Slack notification: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Slack webhook returned status %d", resp.StatusCode)
	}

	return nil
}

// SendSlackRichNotification sends a rich formatted notification to Slack
func SendSlackRichNotification(webhookURL string, projectName string, driftSummary string, planOutput string) error {
	if webhookURL == "" {
		return fmt.Errorf("webhook URL is empty")
	}

	// Truncate plan output if it's too long
	const maxPlanLength = 2000
	if len(planOutput) > maxPlanLength {
		planOutput = planOutput[:maxPlanLength] + "\n... (truncated)"
	}

	// Create a rich Slack message with attachments
	slackMsg := SlackMessage{
		Text:      fmt.Sprintf(":rotating_light: *Drift Detected in Project: %s*", projectName),
		Username:  "TerraDrift Watcher",
		IconEmoji: ":warning:",
		Attachments: []Attachment{
			{
				Color: "danger",
				Title: "Configuration Drift Alert",
				Text:  driftSummary,
				Fields: []Field{
					{
						Title: "Project",
						Value: projectName,
						Short: true,
					},
					{
						Title: "Status",
						Value: "Drift Detected",
						Short: true,
					},
				},
				Footer:     "TerraDrift Watcher",
				FooterIcon: "https://www.terraform.io/favicon.ico",
				Timestamp:  time.Now().Unix(),
			},
			{
				Color: "warning",
				Title: "Plan Output",
				Text:  "```" + planOutput + "```",
			},
		},
	}

	// Marshal the message to JSON
	jsonData, err := json.Marshal(slackMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal Slack message: %w", err)
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create the request
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send Slack notification: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Slack webhook returned status %d", resp.StatusCode)
	}

	return nil
}

// SendSlackNotificationWithRetry sends a Slack notification with retry logic
func SendSlackNotificationWithRetry(webhookURL string, message string, maxRetries int) error {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 1s, 2s, 4s, etc.
			backoff := time.Duration(1<<uint(attempt-1)) * time.Second
			log.Printf("INFO: Retrying Slack notification (attempt %d/%d) after %v", attempt, maxRetries, backoff)
			time.Sleep(backoff)
		}

		err := SendSlackNotification(webhookURL, message)
		if err == nil {
			if attempt > 0 {
				log.Printf("INFO: Slack notification succeeded on attempt %d", attempt+1)
			}
			return nil
		}
		lastErr = err
	}

	return fmt.Errorf("failed after %d retries: %w", maxRetries+1, lastErr)
}

// SendSlackRichNotificationWithRetry sends a rich Slack notification with retry logic
func SendSlackRichNotificationWithRetry(webhookURL string, projectName string, driftSummary string, planOutput string, maxRetries int) error {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 1s, 2s, 4s, etc.
			backoff := time.Duration(1<<uint(attempt-1)) * time.Second
			log.Printf("INFO: Retrying Slack rich notification (attempt %d/%d) after %v", attempt, maxRetries, backoff)
			time.Sleep(backoff)
		}

		err := SendSlackRichNotification(webhookURL, projectName, driftSummary, planOutput)
		if err == nil {
			if attempt > 0 {
				log.Printf("INFO: Slack rich notification succeeded on attempt %d", attempt+1)
			}
			return nil
		}
		lastErr = err
	}

	return fmt.Errorf("failed after %d retries: %w", maxRetries+1, lastErr)
}
