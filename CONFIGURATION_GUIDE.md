# TerraDrift Watcher - Configuration Guide

## Overview
Configuration is done through a YAML file (`config.yml`) that defines:
- Cloud provider authentication
- Notification channels
- Terraform projects to monitor

## ZIP Package Configuration

### What You Need to Configure

When using the ZIP package, you'll configure:
1. **Local paths** to your Terraform projects
2. **Environment variables** for credentials
3. **Notification webhooks** for alerts

### Step-by-Step Configuration

#### 1. Create Your Config File
After extracting the ZIP package:
```powershell
cd C:\tools\terradrift-watcher
copy config.example.yml config.yml
```

#### 2. Edit config.yml
Open `config.yml` in your text editor and customize:

```yaml
# config.yml for ZIP package users
auth_profiles:
  # Define authentication profiles for different environments
  - name: production-aws
    provider: aws
    config:
      # Reference environment variables (set these in Step 3)
      access_key_id: ${AWS_ACCESS_KEY_ID}
      secret_access_key: ${AWS_SECRET_ACCESS_KEY}
      region: us-east-1
  
  - name: staging-aws
    provider: aws
    config:
      access_key_id: ${AWS_STAGING_ACCESS_KEY}
      secret_access_key: ${AWS_STAGING_SECRET_KEY}
      region: us-west-2

notifiers:
  # Slack notification
  - name: slack-ops
    type: slack
    config:
      webhook_url: ${SLACK_WEBHOOK_URL}
    enabled: true
  
  # Microsoft Teams notification (optional)
  - name: teams-alerts
    type: teams
    config:
      webhook_url: ${TEAMS_WEBHOOK_URL}
    enabled: false  # Set to true to enable

projects:
  # Your Terraform projects
  - name: web-application
    path: C:\repos\terraform\web-app        # Windows path to Terraform files
    auth_profile: production-aws            # Which auth profile to use
    notifiers:
      - slack-ops                           # Which notifiers to trigger
    enabled: true
  
  - name: database-layer
    path: C:\repos\terraform\database
    auth_profile: production-aws
    notifiers:
      - slack-ops
      - teams-alerts
    enabled: true
  
  - name: staging-environment
    path: C:\repos\terraform\staging
    auth_profile: staging-aws               # Different auth profile
    notifiers:
      - slack-ops
    enabled: false                          # Disabled for now
```

#### 3. Set Environment Variables

**Option A: Temporary (Current Session)**
```powershell
# AWS Credentials
$env:AWS_ACCESS_KEY_ID = "AKIAIOSFODNN7EXAMPLE"
$env:AWS_SECRET_ACCESS_KEY = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
$env:AWS_STAGING_ACCESS_KEY = "AKIASECONDEXAMPLE"
$env:AWS_STAGING_SECRET_KEY = "anotherSecretKeyExample"

# Notification Webhooks
$env:SLACK_WEBHOOK_URL = "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXX"
$env:TEAMS_WEBHOOK_URL = "https://outlook.office.com/webhook/YOUR-WEBHOOK-URL"
```

**Option B: Permanent (System-wide)**
1. Open System Properties → Advanced → Environment Variables
2. Add each variable under "User variables" or "System variables"
3. Restart PowerShell to load new variables

**Option C: Using a .env file (with a helper script)**
Create `.env` file:
```
AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/T00000000/B00000000/XXXX
```

Load it before running:
```powershell
# Load environment variables from .env file
Get-Content .env | ForEach-Object {
    if ($_ -match '^([^=]+)=(.*)$') {
        [System.Environment]::SetEnvironmentVariable($matches[1], $matches[2])
    }
}
```

#### 4. Test Your Configuration
```powershell
# Verify configuration is valid
.\terradrift-watcher.exe run --config config.yml --dry-run

# Run actual drift check
.\terradrift-watcher.exe run --config config.yml
```

---

## Docker Configuration

### What You Need to Configure

When using Docker, you'll configure:
1. **Container paths** for mounted volumes
2. **Environment variables** passed to container
3. **Volume mounts** for local directories

### Step-by-Step Configuration

#### 1. Create Directory Structure
```bash
mkdir terradrift-workspace
cd terradrift-workspace

# Create this structure:
# terradrift-workspace/
# ├── config.yml
# ├── .env
# └── terraform-projects/
#     ├── web-app/
#     └── database/
```

#### 2. Create config.yml for Docker
```yaml
# config.yml for Docker users
auth_profiles:
  - name: aws-main
    provider: aws
    config:
      # These will come from container environment
      access_key_id: ${AWS_ACCESS_KEY_ID}
      secret_access_key: ${AWS_SECRET_ACCESS_KEY}
      region: ${AWS_DEFAULT_REGION}

notifiers:
  - name: slack
    type: slack
    config:
      webhook_url: ${SLACK_WEBHOOK_URL}
    enabled: true

projects:
  # Use container paths (will be mounted from local)
  - name: web-application
    path: /terraform/web-app          # Path inside container
    auth_profile: aws-main
    notifiers:
      - slack
    enabled: true
  
  - name: database-layer
    path: /terraform/database         # Path inside container
    auth_profile: aws-main
    notifiers:
      - slack
    enabled: true
```

#### 3. Create .env File
```bash
# .env file for Docker
AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
AWS_DEFAULT_REGION=us-east-1
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/T00000000/B00000000/XXXX
```

#### 4. Run with Docker

**Simple Docker Run:**
```bash
docker run --rm \
  -v $(pwd)/config.yml:/config/config.yml:ro \
  -v $(pwd)/terraform-projects:/terraform:ro \
  --env-file .env \
  yourusername/terradrift-watcher:latest \
  run --config /config/config.yml
```

**Using Docker Compose:**
Create `docker-compose.yml`:
```yaml
version: '3.8'
services:
  terradrift-watcher:
    image: yourusername/terradrift-watcher:latest
    volumes:
      # Mount config file
      - ./config.yml:/config/config.yml:ro
      # Mount Terraform projects
      - ./terraform-projects:/terraform:ro
    env_file:
      - .env
    command: run --config /config/config.yml
```

Run with:
```bash
docker-compose run --rm terradrift-watcher
```

---

## Configuration Examples

### Minimal Configuration
```yaml
# Simplest possible configuration
auth_profiles:
  - name: default
    provider: aws
    config:
      access_key_id: ${AWS_ACCESS_KEY_ID}
      secret_access_key: ${AWS_SECRET_ACCESS_KEY}
      region: us-east-1

projects:
  - name: my-project
    path: ./terraform
    auth_profile: default
    enabled: true
```

### Multi-Cloud Configuration
```yaml
auth_profiles:
  # AWS Profile
  - name: aws-prod
    provider: aws
    config:
      access_key_id: ${AWS_ACCESS_KEY_ID}
      secret_access_key: ${AWS_SECRET_ACCESS_KEY}
      region: us-east-1
  
  # Azure Profile
  - name: azure-prod
    provider: azure
    config:
      client_id: ${AZURE_CLIENT_ID}
      client_secret: ${AZURE_CLIENT_SECRET}
      subscription_id: ${AZURE_SUBSCRIPTION_ID}
      tenant_id: ${AZURE_TENANT_ID}
  
  # GCP Profile
  - name: gcp-prod
    provider: gcp
    config:
      credentials_file: ${GOOGLE_APPLICATION_CREDENTIALS}
      project_id: ${GCP_PROJECT_ID}

projects:
  - name: aws-infrastructure
    path: ./terraform/aws
    auth_profile: aws-prod
    enabled: true
  
  - name: azure-resources
    path: ./terraform/azure
    auth_profile: azure-prod
    enabled: true
  
  - name: gcp-services
    path: ./terraform/gcp
    auth_profile: gcp-prod
    enabled: true
```

### Scheduled Monitoring Configuration
```yaml
# For automated/scheduled runs
check_interval: 6h  # Check every 6 hours

auth_profiles:
  - name: production
    provider: aws
    config:
      access_key_id: ${AWS_ACCESS_KEY_ID}
      secret_access_key: ${AWS_SECRET_ACCESS_KEY}
      region: us-east-1

notifiers:
  # Multiple notification channels
  - name: slack-critical
    type: slack
    config:
      webhook_url: ${SLACK_CRITICAL_WEBHOOK}
    enabled: true
  
  - name: slack-info
    type: slack
    config:
      webhook_url: ${SLACK_INFO_WEBHOOK}
    enabled: true
  
  - name: email-ops
    type: email
    config:
      smtp_host: smtp.gmail.com
      smtp_port: "587"
      from: alerts@company.com
      to: ops-team@company.com
      username: ${EMAIL_USERNAME}
      password: ${EMAIL_PASSWORD}
    enabled: true

projects:
  # Critical production infrastructure
  - name: production-core
    path: ./terraform/production/core
    auth_profile: production
    notifiers:
      - slack-critical
      - email-ops
    enabled: true
  
  # Less critical services
  - name: production-analytics
    path: ./terraform/production/analytics
    auth_profile: production
    notifiers:
      - slack-info
    enabled: true
```

---

## Troubleshooting Configuration

### Common Issues

1. **"Config file not found"**
   - Ensure config.yml exists in the specified location
   - Use full path: `--config C:\full\path\to\config.yml`

2. **"Invalid configuration"**
   - Check YAML syntax (proper indentation)
   - Validate all required fields are present
   - Run with `--validate` flag to check config

3. **"Environment variable not found"**
   - Verify variable is set: `echo $env:AWS_ACCESS_KEY_ID`
   - Check for typos in variable names
   - Restart shell after setting permanent variables

4. **"Path not found" (Docker)**
   - Ensure volumes are mounted correctly
   - Check container paths match config.yml paths
   - Verify local directories exist before mounting

5. **"Authentication failed"**
   - Verify credentials are correct
   - Check IAM/Azure/GCP permissions
   - Test credentials with cloud CLI tools first

### Configuration Best Practices

1. **Use Environment Variables for Secrets**
   - Never hardcode credentials in config.yml
   - Use `${VAR_NAME}` syntax for sensitive data

2. **Organize Projects Logically**
   - Group by environment (prod, staging, dev)
   - Or by team/service ownership
   - Use clear, descriptive names

3. **Start Simple**
   - Begin with one project
   - Test thoroughly
   - Add more projects gradually

4. **Version Control**
   - Commit config.example.yml to Git
   - Never commit actual config.yml with secrets
   - Use .gitignore to exclude sensitive files

5. **Regular Testing**
   - Run with `--dry-run` flag first
   - Test notifications work
   - Verify all projects are accessible 