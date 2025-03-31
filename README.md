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
curl -L https://github.com/sparkfabrik/sparkci/releases/download/latest/sparkci_0.0.0-dev_linux_amd64.tar.gz | tar xz
sudo mv sparkci /usr/local/bin/

# Linux (arm64)
curl -L https://github.com/sparkfabrik/sparkci/releases/download/latest/sparkci_0.0.0-dev_linux_arm64.tar.gz | tar xz
sudo mv sparkci /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/sparkfabrik/sparkci/releases/download/latest/sparkci_0.0.0-dev_darwin_amd64.tar.gz | tar xz
sudo mv sparkci /usr/local/bin/

# macOS (Apple Silicon)
curl -L https://github.com/sparkfabrik/sparkci/releases/download/latest/sparkci_0.0.0-dev_darwin_arm64.tar.gz | tar xz
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

SparkCI provides several command groups for different operations:

```
sparkci [command]
```

Available Commands:

```
gitlab      GitLab CI utilities
gwif        Google Cloud Workload Identity Federation utilities
help        Help about any command
```

### Command Reference

#### GitLab Commands

```
sparkci gitlab [subcommand]
```

Subcommands:

- `print-env` - Print information about the current GitLab CI environment

  ```bash
  # Display GitLab CI environment in default text format
  sparkci gitlab print-env

  # Output in JSON format
  sparkci gitlab print-env --format json

  # Output in YAML format
  sparkci gitlab print-env --format yaml
  ```

- `format` - Format GitLab CI output with sections and banners

  - `section` - Create collapsible sections in GitLab CI output

    ```bash
    # Start a section
    sparkci gitlab format section --title "build-logs" --description "Build logs for the project"

    # End a section
    sparkci gitlab format section --title "build-logs" --end
    ```

  - `banner` - Print a highlighted banner in GitLab CI output
    ```bash
    # Print a banner
    sparkci gitlab format banner --text "Deployment Started"
    ```

#### Google Workload Identity Federation (GWIF) Commands

```
sparkci gwif [subcommand]
```

Subcommands:

- `configure` - Orchestrate the setup of Workload Identity Federation

  ```bash
  # Configure Workload Identity Federation by running all necessary setup steps
  sparkci gwif configure
  ```

- `gcloud-auth` - Authenticate gcloud CLI with WIF

  ```bash
  # Authenticate gcloud CLI with Workload Identity Federation
  sparkci gwif gcloud-auth
  ```

- `gcloud-exec` - Execute a gcloud command with WIF authentication

  ```bash
  # Run a gcloud command using WIF authentication
  sparkci gwif gcloud-exec -- secrets versions access latest --project="my-project"
  ```

- `print-gitlab-jwt` - Print the GitLab OIDC JWT token

  ```bash
  # Print JWT token in JSON format (default)
  sparkci gwif print-gitlab-jwt

  # Print JWT token in text format
  sparkci gwif print-gitlab-jwt --format text
  ```

- `print-sa-token` - Print Google Workload Identity token

  ```bash
  # Print the Google service account token
  sparkci gwif print-sa-token
  ```

- `print-vars` - Print the Workload Identity Federation standard env variables

  ```bash
  # Print WIF environment variables in formatted output
  sparkci gwif print-vars
  ```

- `status` - Check Workload Identity Federation status

  ```bash
  # Check if WIF is configured (exit code only)
  sparkci gwif status

  # Check WIF configuration and print errors if any
  sparkci gwif status --silent=false

  # Check and print WIF configuration details
  sparkci gwif status --print
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

## Coding conventions

SparkCI follows these Go coding patterns:

- Command structure follows [Cobra](https://github.com/spf13/cobra/blob/main/site/content/user_guide.md) best practices
- Main commands are organized in their own package under `cmd/` (e.g., `cmd/gitlab/`)
- Main command files export a PascalCase command variable (e.g., `GitlabCommand`) for imports
- Subcommands are defined in separate files with camelCase variables (e.g., `printEnv`)
- Use `RunE` instead of `Run` when the command can return errors
- Command flags are defined in each command's `init()` function
- Core functionality is implemented in packages under `pkg/`
- Use project logger utilities: `utils.Info()`, `utils.Debug()`, `utils.Warn()`, `utils.Error()`, `utils.Fatal()`
- Each command should be in its own file with a matching filename
- Main commands should provide a default `Run` function that calls `cmd.Help()`

### Development Workflow

This project follows a simplified GitHub Flow workflow:

- `main` branch is the primary branch for both development and releases
- Feature branches should be created from and merged back into `main` via pull requests
- All work should be done in feature branches, not directly on `main`
- Branch names should follow the pattern: `feature/description`, `bugfix/description`, or `docs/description`
- Commit messages should be clear and descriptive
- Pull requests should include adequate test coverage for new features

Development builds are automatically generated from the `main` branch on each push.
Stable releases are created by tagging the `main` branch with a version tag (e.g., `v0.1.0`).

### Testing

Run the test suite before submitting a pull request:

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage

# Run tests for a specific package
go test ./pkg/gitlab/...

# Run a specific test
go test ./pkg/gitlab/... -run TestFunctionName
```

For integration tests that require external services, use the appropriate flags:

```bash
# Run integration tests (requires credentials)
go test ./... -tags=integration
```

## Maintainers

This project is maintained by [SparkFabrik](https://www.sparkfabrik.com/).
