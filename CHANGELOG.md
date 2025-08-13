# Changelog

All notable changes to TerraDrift Watcher will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of TerraDrift Watcher
- Automated drift detection for Terraform projects
- Multi-cloud support (AWS, Azure, GCP)
- Slack notifications with retry logic
- Concurrent run protection with file locking
- Graceful shutdown handling
- Environment variable substitution in configuration
- Verbose mode for detailed output
- `--fail-on-drift` flag for CI/CD integration
- `--force` flag to override locks
- Comprehensive error handling and recovery
- Docker support
- Cross-platform binaries (Linux, macOS, Windows)

### Security
- Automatic cleanup of authentication environment variables
- Secure credential management

## [1.0.0] - 2024-01-15

### Added
- First stable release
- Core drift detection functionality
- Configuration validation
- Basic notification support

### Fixed
- Terraform lock file cleanup on errors
- Authentication environment variable leaks

### Changed
- Improved error messages for better debugging
- Enhanced Terraform backend initialization

## Version History

- **1.0.0** - Initial stable release with core functionality
- **0.9.0** - Beta release with testing feedback incorporated
- **0.8.0** - Alpha release for internal testing

---

For detailed release notes, see the [GitHub Releases](https://github.com/yourusername/terradrift-watcher/releases) page. 