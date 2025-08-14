# Contributing to TerraDrift Watcher

First off, thank you for considering contributing to TerraDrift Watcher! It's people like you that make TerraDrift Watcher such a great tool.

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code. 

### Reporting Bugs

Before creating bug reports, please check existing issues as you might find out that you don't need to create one. When you are creating a bug report, please include as many details as possible:

* **Use a clear and descriptive title** for the issue to identify the problem
* **Describe the exact steps which reproduce the problem** in as many details as possible
* **Provide specific examples to demonstrate the steps**
* **Describe the behavior you observed after following the steps**
* **Explain which behavior you expected to see instead and why**
* **Include logs and error messages** (sanitize sensitive information)
* **Include your configuration file** (sanitize sensitive information)
* **Specify your environment** (OS, Terraform version, Go version if building from source)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, please include:

* **Use a clear and descriptive title** for the issue to identify the suggestion
* **Provide a step-by-step description of the suggested enhancement**
* **Provide specific examples to demonstrate the steps**
* **Describe the current behavior** and **explain which behavior you expected to see instead**
* **Explain why this enhancement would be useful**

### Pull Requests

1. Fork the repo and create your branch from `main`
2. If you've added code that should be tested, add tests
3. If you've changed APIs, update the documentation
4. Ensure the test suite passes
5. Make sure your code follows the existing code style
6. Issue that pull request!

### Prerequisites

- Go 1.21 or higher
- Terraform 1.0.0 or higher
- Git

### Setting Up Your Development Environment

1. Fork and clone the repository:
```bash
git clone https://github.com/yourusername/terradrift-watcher.git
cd terradrift-watcher
```

2. Install dependencies:
```bash
go mod download
```

3. Create a branch for your feature or fix:
```bash
git checkout -b feature/your-feature-name
```

### Building and Testing

```bash
# Build the project
go build -o terradrift-watcher .

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Run linter (install golangci-lint first)
golangci-lint run

# Format code
go fmt ./...
```

### Project Structure

```
terradrift-watcher/
├── cmd/                    # CLI commands
│   ├── root.go            # Root command setup
│   └── run.go             # Run command implementation
├── internal/              # Private application code
│   ├── config/            # Configuration management
│   ├── detector/          # Drift detection logic
│   ├── lock/              # Locking mechanism
│   ├── notifier/          # Notification handlers
│   └── terraform/         # Terraform integration
├── testdata/              # Test fixtures
└── main.go               # Application entry point
```

## Coding Guidelines

### Go Code Style

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` to format your code
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions small and focused
- Handle errors explicitly

### Example Code Style

```go
// CheckDrift runs terraform plan to detect configuration drift
// Returns the plan output, exit code, and any error
func CheckDrift(projectPath string) (string, int, error) {
    // Validate input
    if projectPath == "" {
        return "", 1, fmt.Errorf("project path cannot be empty")
    }
    
    // Implementation...
    return output, exitCode, nil
}
```

### Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line

Example:
```
Add retry logic for Slack notifications

- Implement exponential backoff for failed notifications
- Add configurable max retry attempts
- Update documentation with retry behavior

Fixes #123
```
### Writing Tests

- Write unit tests for new functionality
- Maintain or improve code coverage
- Use table-driven tests where appropriate
- Mock external dependencies

Example test:
```go
func TestLoadConfig(t *testing.T) {
    tests := []struct {
        name    string
        path    string
        wantErr bool
    }{
        {"valid config", "testdata/valid.yml", false},
        {"invalid yaml", "testdata/invalid.yml", true},
        {"missing file", "testdata/missing.yml", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := LoadConfig(tt.path)
            if (err != nil) != tt.wantErr {
                t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Documentation

- Update README.md if you change functionality
- Update CONFIGURATION_GUIDE.md for configuration changes
- Add inline documentation for complex logic
- Include examples for new features

## Release Process

1. Update version numbers in relevant files
2. Update CHANGELOG.md with release notes
3. Create a pull request with the changes
4. After merge, create a tagged release
5. Build binaries using `build.sh` or `build.bat`
6. Upload binaries to the GitHub release


## Recognition

Contributors will be recognized in the project README and release notes. Thank you for your contributions!

## License

By contributing, you agree that your contributions will be licensed under the MIT License. 
