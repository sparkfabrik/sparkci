# SparkCI

A CLI tool designed for GitLab CI automation and integration.

## ⚠️ Experimental Project Notice

**This is an experimental project in early development.**

This project is under active development and is not yet ready for production use. APIs, commands, and behavior may change without notice.

## Overview

SparkCI provides command-line utilities for GitLab CI operations. It's designed to be used both within GitLab CI pipelines and locally for development purposes.

Key features:

- GitLab CI environment management
- Google Cloud Workload Identity Federation (WIF) integration
- Nicely formatted output for CI environment variables

## Installation

### Binary Installation

Download the latest release from the [releases page](https://github.com/sparkfabrik/sparkci/releases).

### Build from Source

```bash
# Clone the repository
git clone https://github.com/sparkfabrik/sparkci.git
cd sparkci

# Build
make build
```

## Usage

### GitLab CI Environment

Display information about the current GitLab CI environment:

```bash
sparkci gitlab env
```

With different output formats:

```bash
sparkci gitlab env --format json
sparkci gitlab env --format yaml
sparkci gitlab env --format table
```

### Google Cloud Workload Identity Federation

Commands for working with Google Cloud Workload Identity Federation (WIF) in GitLab CI:

```bash
sparkci gwif auth    # Authenticate with Google Cloud using WIF
sparkci gwif config  # Configure WIF settings
sparkci gwif exec    # Execute a command with WIF authentication
```

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## Contributing

As this is an experimental project, contributions are welcome but may be subject to significant changes.

### Development Workflow

This project follows a simplified GitHub Flow workflow:

- `main` branch is the primary branch for both development and releases
- Feature branches should be created from and merged back into `main` via pull requests
- All work should be done in feature branches, not directly on `main`

Development releases are automatically generated from the `main` branch on each push.
Stable releases are created by tagging the `main` branch with a version tag (e.g., `v0.1.0`).

## Maintainers

This project is maintained by [SparkFabrik](https://www.sparkfabrik.com/).