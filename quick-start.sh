#!/bin/bash

# Quick Start Script for TerraDrift Watcher
# This script helps you get started with TerraDrift Watcher quickly

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}╔══════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     TerraDrift Watcher Quick Start       ║${NC}"
echo -e "${BLUE}╚══════════════════════════════════════════╝${NC}"
echo ""

# Check if Terraform is installed
echo -e "${YELLOW}Checking prerequisites...${NC}"
if ! command -v terraform &> /dev/null; then
    echo -e "${RED}✗ Terraform is not installed${NC}"
    echo "Please install Terraform first: https://www.terraform.io/downloads"
    exit 1
fi
echo -e "${GREEN}✓ Terraform is installed${NC}"

# Check if Go is installed (optional, for building from source)
if command -v go &> /dev/null; then
    echo -e "${GREEN}✓ Go is installed (optional)${NC}"
else
    echo -e "${YELLOW}ℹ Go is not installed (only needed for building from source)${NC}"
fi

# Download or build TerraDrift Watcher
echo ""
echo -e "${YELLOW}Setting up TerraDrift Watcher...${NC}"

if [ -f "terradrift-watcher" ] || [ -f "terradrift-watcher.exe" ]; then
    echo -e "${GREEN}✓ TerraDrift Watcher binary found${NC}"
elif command -v go &> /dev/null && [ -f "go.mod" ]; then
    echo "Building TerraDrift Watcher from source..."
    go build -o terradrift-watcher .
    echo -e "${GREEN}✓ Built TerraDrift Watcher${NC}"
else
    echo "Downloading TerraDrift Watcher..."
    # Detect OS and architecture
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case "$ARCH" in
        x86_64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
    esac
    
    case "$OS" in
        darwin) OS="darwin" ;;
        linux) OS="linux" ;;
        mingw*|msys*|cygwin*) OS="windows"; EXT=".exe" ;;
        *) echo "Unsupported OS: $OS"; exit 1 ;;
    esac
    
    BINARY="terradrift-watcher-${OS}-${ARCH}${EXT}"
    URL="https://github.com/yourusername/terradrift-watcher/releases/latest/download/${BINARY}"
    
    echo "Downloading from: $URL"
    curl -L -o terradrift-watcher${EXT} "$URL" || wget -O terradrift-watcher${EXT} "$URL"
    chmod +x terradrift-watcher${EXT}
    echo -e "${GREEN}✓ Downloaded TerraDrift Watcher${NC}"
fi

# Create example configuration if it doesn't exist
if [ ! -f "config.yml" ]; then
    echo ""
    echo -e "${YELLOW}Creating example configuration...${NC}"
    
    if [ -f "config.example.yml" ]; then
        cp config.example.yml config.yml
    else
        cat > config.yml << 'EOF'
# TerraDrift Watcher Configuration
projects:
  - name: my-terraform-project
    path: ./terraform  # Update this path
    auth_profile: aws  # Remove if not using AWS
    notifiers:
      - slack
    enabled: true

auth_profiles:
  - name: aws
    provider: aws
    config:
      access_key_id: ${AWS_ACCESS_KEY_ID}
      secret_access_key: ${AWS_SECRET_ACCESS_KEY}
      region: us-east-1

notifiers:
  - name: slack
    type: slack
    config:
      webhook_url: ${SLACK_WEBHOOK_URL}
    enabled: true
EOF
    fi
    echo -e "${GREEN}✓ Created config.yml${NC}"
    echo -e "${YELLOW}ℹ Please edit config.yml to match your setup${NC}"
fi

# Create .env.example file
echo ""
echo -e "${YELLOW}Creating environment variables template...${NC}"
cat > .env.example << 'EOF'
# AWS Credentials (if using AWS)
AWS_ACCESS_KEY_ID=your-access-key-here
AWS_SECRET_ACCESS_KEY=your-secret-key-here

# Slack Webhook (if using Slack notifications)
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/WEBHOOK/URL

# Azure Credentials (if using Azure)
# AZURE_CLIENT_ID=your-client-id
# AZURE_CLIENT_SECRET=your-client-secret
# AZURE_SUBSCRIPTION_ID=your-subscription-id
# AZURE_TENANT_ID=your-tenant-id

# GCP Credentials (if using GCP)
# GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account.json
EOF
echo -e "${GREEN}✓ Created .env.example${NC}"

# Instructions
echo ""
echo -e "${BLUE}═══════════════════════════════════════════${NC}"
echo -e "${GREEN}Setup Complete!${NC}"
echo -e "${BLUE}═══════════════════════════════════════════${NC}"
echo ""
echo "Next steps:"
echo "1. Edit config.yml to configure your Terraform projects"
echo "2. Copy .env.example to .env and add your credentials"
echo "3. Run: source .env (or use your preferred method to set environment variables)"
echo "4. Test: ./terradrift-watcher run --config config.yml"
echo ""
echo "For more information:"
echo "- Configuration Guide: CONFIGURATION_GUIDE.md"
echo "- Full Documentation: README.md"
echo "- Examples: examples/"
echo ""
echo -e "${GREEN}Happy drift detecting! 🚀${NC}" 