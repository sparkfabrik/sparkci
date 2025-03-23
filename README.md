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

**Latest Development Build:**

```bash
# Linux (amd64)
curl -L https://github.com/sparkfabrik/sparkci/releases/download/latest/sparkci_linux_amd64.tar.gz | tar xz
sudo mv sparkci /usr/local/bin/

# Linux (arm64)
curl -L https://github.com/sparkfabrik/sparkci/releases/download/latest/sparkci_linux_arm64.tar.gz | tar xz
sudo mv sparkci /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/sparkfabrik/sparkci/releases/download/latest/sparkci_darwin_amd64.tar.gz | tar xz
sudo mv sparkci /usr/local/bin/

# macOS (Apple Silicon)
curl -L https://github.com/sparkfabrik/sparkci/releases/download/latest/sparkci_darwin_arm64.tar.gz | tar xz
sudo mv sparkci /usr/local/bin/
```

**Stable Release:**

```bash
# Download a specific version (replace v1.0.0 with the desired version)
curl -L https://github.com/sparkfabrik/sparkci/releases/download/v1.0.0/sparkci_v1.0.0_linux_amd64.tar.gz | tar xz
sudo mv sparkci /usr/local/bin/
```

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

## Release Process

SparkCI uses GitHub Actions and GoReleaser for automated releases. The release process creates both development and stable releases:

### Development Builds (Latest)

A "latest" development build is automatically created and updated whenever code is pushed to the `main` branch:

- This release is always tagged as `latest` on GitHub
- It is marked as a "pre-release" to indicate it's not a stable version
- The existing "latest" release is replaced with each new build
- This build contains the most recent changes from the `main` branch
- Download it from [github.com/sparkfabrik/sparkci/releases/tag/latest](https://github.com/sparkfabrik/sparkci/releases/tag/latest)

### Stable Releases

Stable releases are created by tagging the `main` branch:

1. Ensure you're on the latest `main` branch:

   ```bash
   git checkout main
   git pull origin main
   ```

2. Create and push a new tag with semantic versioning:

   ```bash
   git tag v1.0.0  # Use appropriate version
   git push origin v1.0.0
   ```

3. The GitHub Action will automatically:
   - Run tests
   - Build binaries for multiple platforms (Linux AMD64/ARM64, macOS AMD64/ARM64)
   - Create a GitHub release with the tagged version
   - Upload the binaries and checksums

Stable releases follow [Semantic Versioning](https://semver.org/) (MAJOR.MINOR.PATCH).

### Release Artifacts

Each release (development and stable) includes:

- Binary executables for supported platforms
- SHA256 checksums for verification
- Source code archives
- Release notes generated from git commits

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## Contributing

As this is an experimental project, contributions are welcome but may be subject to significant changes.

### Development Workflow

This project follows a simplified GitHub Flow workflow:

- `main` branch is the primary branch for both development and releases
- Feature branches should be created from and merged back into `main` via pull requests
- All work should be done in feature branches, not directly on `main`

Development builds are automatically generated from the `main` branch on each push.
Stable releases are created by tagging the `main` branch with a version tag (e.g., `v0.1.0`).

### Testing

Run the test suite before submitting a pull request:

```bash
make test
```

## Maintainers

This project is maintained by [SparkFabrik](https://www.sparkfabrik.com/).
