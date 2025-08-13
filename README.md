# TerraDrift Watcher 🔍

[![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Terraform](https://img.shields.io/badge/Terraform-1.0%2B-purple.svg)](https://www.terraform.io/)

A robust, production-grade CLI tool for detecting configuration drift in Terraform projects by comparing live infrastructure state against code in Git. Prevent infrastructure drift before it becomes a problem!

## 🌟 Features

- 🔍 **Automated Drift Detection**: Continuously monitors Terraform projects for configuration drift
- ☁️ **Multi-Cloud Support**: Works with AWS, Azure, GCP, and any Terraform provider
- 📢 **Smart Notifications**: Slack integration with retry logic and rich formatting
- 🔐 **Secure Authentication**: Environment-based credential management with automatic cleanup
- 🔒 **Concurrent Run Protection**: File-based locking prevents conflicting executions
- 🛡️ **Graceful Shutdown**: Proper signal handling and resource cleanup
- 📦 **Standalone Binary**: Single executable, no runtime dependencies
- 🚀 **Production-Ready**: Battle-tested error handling and recovery mechanisms

## 📋 Table of Contents

- [Quick Start](#-quick-start)
- [Installation](#-installation)
- [Configuration](#-configuration)
- [Usage](#-usage)
- [Examples](#-examples)
- [Integration Examples](#-integration-examples)
- [Architecture](#-architecture)
- [Contributing](#-contributing)
- [License](#-license)

## 🚀 Quick Start

```bash
# 1. Download the latest release
wget https://github.com/yourusername/terradrift-watcher/releases/latest/download/terradrift-watcher-linux-amd64
chmod +x terradrift-watcher-linux-amd64
mv terradrift-watcher-linux-amd64 terradrift-watcher

# 2. Create a configuration file
cp config.example.yml config.yml
# Edit config.yml with your projects and settings

# 3. Set up authentication
export AWS_ACCESS_KEY_ID="your-key"
export AWS_SECRET_ACCESS_KEY="your-secret"
export SLACK_WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK"

# 4. Run drift detection
./terradrift-watcher run --config config.yml
```

## 📥 Installation

### Option 1: Download Pre-built Binary (Recommended)

Download the latest release for your platform from the [Releases](https://github.com/yourusername/terradrift-watcher/releases) page:

- **Linux**: `terradrift-watcher-linux-amd64`
- **macOS Intel**: `terradrift-watcher-darwin-amd64`
- **macOS Apple Silicon**: `terradrift-watcher-darwin-arm64`
- **Windows**: `terradrift-watcher-windows-amd64.exe`

### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/terradrift-watcher.git
cd terradrift-watcher

# Build the binary
go build -o terradrift-watcher .

# Or use the build script for cross-platform builds
./build.sh  # Unix/Linux/macOS
# or
build.bat   # Windows
```

### Option 3: Docker

```bash
# Using Docker Hub
docker pull yourusername/terradrift-watcher:latest

# Or build locally
docker build -t terradrift-watcher .

# Run with Docker
docker run --rm \
  -v $(pwd)/config.yml:/config.yml \
  -e AWS_ACCESS_KEY_ID \
  -e AWS_SECRET_ACCESS_KEY \
  -e SLACK_WEBHOOK_URL \
  yourusername/terradrift-watcher run --config /config.yml
```

### Prerequisites

- **Terraform**: Version 1.0.0 or higher must be installed and available in PATH
- **Go**: Version 1.21 or higher (only for building from source)

## ⚙️ Configuration

### Basic Configuration

Create a `config.yml` file:

```yaml
# Check interval for scheduled runs (optional)
check_interval: "1h"

# Projects to monitor
projects:
  - name: production-vpc
    path: ./terraform/production/vpc
    auth_profile: aws-prod
    notifiers:
      - slack-ops
    enabled: true

# Authentication profiles
auth_profiles:
  - name: aws-prod
    provider: aws
    config:
      access_key_id: ${AWS_ACCESS_KEY_ID}
      secret_access_key: ${AWS_SECRET_ACCESS_KEY}
      region: us-east-1

# Notification channels
notifiers:
  - name: slack-ops
    type: slack
    config:
      webhook_url: ${SLACK_WEBHOOK_URL}
    enabled: true
```

### Environment Variables

The tool supports environment variable substitution using `${VAR_NAME}` syntax:

```bash
export AWS_ACCESS_KEY_ID="your-access-key"
export AWS_SECRET_ACCESS_KEY="your-secret-key"
export SLACK_WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK"
```

For detailed configuration options, see [CONFIGURATION_GUIDE.md](CONFIGURATION_GUIDE.md).

## 🎮 Usage

### Basic Commands

```bash
# Run drift detection
terradrift-watcher run --config config.yml

# Run with verbose output
terradrift-watcher run --config config.yml --verbose

# Exit with code 2 if drift is detected (useful for CI/CD)
terradrift-watcher run --config config.yml --fail-on-drift

# Force run even if another instance is running
terradrift-watcher run --config config.yml --force

# Show version
terradrift-watcher --version

# Show help
terradrift-watcher --help
```

### Command Line Options

| Flag | Description | Default |
|------|-------------|---------|
| `-c, --config` | Path to configuration file | `config.yml` |
| `-v, --verbose` | Show full terraform plan output | `false` |
| `--fail-on-drift` | Exit with code 2 if drift detected | `false` |
| `--force` | Force release any existing lock | `false` |

## 📚 Examples

### Example 1: Multi-Environment Setup

```yaml
projects:
  - name: prod-infrastructure
    path: ./terraform/prod
    auth_profile: aws-prod
    notifiers: [slack-critical]
    
  - name: staging-infrastructure
    path: ./terraform/staging
    auth_profile: aws-staging
    notifiers: [slack-dev]

auth_profiles:
  - name: aws-prod
    provider: aws
    config:
      access_key_id: ${AWS_PROD_ACCESS_KEY}
      secret_access_key: ${AWS_PROD_SECRET_KEY}
      
  - name: aws-staging
    provider: aws
    config:
      access_key_id: ${AWS_STAGING_ACCESS_KEY}
      secret_access_key: ${AWS_STAGING_SECRET_KEY}
```

### Example 2: Multi-Cloud Setup

```yaml
projects:
  - name: aws-resources
    path: ./terraform/aws
    auth_profile: aws-main
    
  - name: azure-resources
    path: ./terraform/azure
    auth_profile: azure-main
    
  - name: gcp-resources
    path: ./terraform/gcp
    auth_profile: gcp-main

auth_profiles:
  - name: aws-main
    provider: aws
    config:
      access_key_id: ${AWS_ACCESS_KEY_ID}
      secret_access_key: ${AWS_SECRET_ACCESS_KEY}
      
  - name: azure-main
    provider: azure
    config:
      client_id: ${AZURE_CLIENT_ID}
      client_secret: ${AZURE_CLIENT_SECRET}
      subscription_id: ${AZURE_SUBSCRIPTION_ID}
      tenant_id: ${AZURE_TENANT_ID}
      
  - name: gcp-main
    provider: gcp
    config:
      GOOGLE_APPLICATION_CREDENTIALS: ${GOOGLE_APPLICATION_CREDENTIALS}
```

## 🔄 Integration Examples

### Cron Job (Linux/macOS)

```bash
# Add to crontab to run every 6 hours
crontab -e

# Add this line:
0 */6 * * * cd /path/to/terradrift-watcher && ./terradrift-watcher run --config config.yml
```

### Windows Task Scheduler

```powershell
# Create a scheduled task to run every 6 hours
$action = New-ScheduledTaskAction -Execute "C:\path\to\terradrift-watcher.exe" -Argument "run --config config.yml"
$trigger = New-ScheduledTaskTrigger -Daily -At 12am -DaysInterval 1 -RepetitionInterval (New-TimeSpan -Hours 6)
Register-ScheduledTask -TaskName "TerraDrift-Watcher" -Action $action -Trigger $trigger
```

### Docker Compose for Scheduled Runs

```yaml
version: '3'
services:
  terradrift-watcher:
    image: yourusername/terradrift-watcher:latest
    volumes:
      - ./config.yml:/config.yml:ro
      - ./terraform:/terraform:ro
    environment:
      - AWS_ACCESS_KEY_ID
      - AWS_SECRET_ACCESS_KEY
      - SLACK_WEBHOOK_URL
    command: run --config /config.yml
    restart: unless-stopped
```

## 🏗️ Architecture

### Project Structure

```
terradrift-watcher/
├── cmd/                    # CLI commands
│   ├── root.go            # Root command setup
│   └── run.go             # Run command implementation
├── internal/
│   ├── config/            # Configuration management
│   │   ├── loader.go      # YAML loading and validation
│   │   └── models.go      # Data structures
│   ├── detector/          # Drift detection engine
│   │   └── engine.go      # Orchestration logic
│   ├── lock/              # Concurrent run protection
│   │   └── filelock.go    # File-based locking
│   ├── notifier/          # Notification handlers
│   │   └── slack.go       # Slack integration with retry
│   └── terraform/         # Terraform wrapper
│       └── executor.go    # Command execution
├── testdata/              # Test fixtures
├── config.example.yml     # Example configuration
├── Dockerfile             # Container image
├── go.mod                 # Go dependencies
└── main.go               # Entry point
```

### Key Features Implementation

- **Concurrent Run Protection**: File-based locking with PID tracking
- **Retry Logic**: Exponential backoff for transient failures
- **Graceful Shutdown**: Signal handling (SIGINT/SIGTERM) with cleanup
- **Resource Cleanup**: Automatic cleanup of credentials and lock files
- **Error Recovery**: Enhanced error messages and recovery strategies

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/yourusername/terradrift-watcher.git
cd terradrift-watcher

# Install dependencies
go mod download

# Run tests
go test ./...

# Build locally
go build -o terradrift-watcher .
```

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI interface
- Uses [yaml.v3](https://gopkg.in/yaml.v3) for configuration parsing
- Inspired by infrastructure drift detection best practices

## 📞 Support

- 📧 Email: support@example.com
- 💬 Slack: [Join our community](https://slack.example.com)
- 🐛 Issues: [GitHub Issues](https://github.com/yourusername/terradrift-watcher/issues)
- 📖 Docs: [Documentation](https://docs.example.com)

---

Made with ❤️ by the TerraDrift team 