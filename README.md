# CLI

## Usage

```
OpenLabs CLI is a command line interface for managing OpenLabs and its associated templates, ranges, and plugins.

Usage:
  openlabs [flags]
  openlabs [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  config      Manage CLI configuration
  help        Help about any command
  plugins     View and deploy plugins
  range       Deploy and manage ranges
  secrets     Upload and manage secrets
  templates   Upload and manage templates
  user        Manage user accounts
  version     Print version information

Flags:
      --api-url string   URL of the OpenLabs API server (default "http://localhost:8000")
      --debug            Enable debug mode to see detailed request/response information
  -h, --help             help for openlabs
      --token string     Authentication token for OpenLabs API

Use "openlabs [command] --help" for more information about a command.
```

## Development

### Build
```bash
make build
```

### Lint
The project uses [golangci-lint](https://golangci-lint.run/) for linting:

```bash
# Install golangci-lint if not already installed
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linting
make lint
```

### Build for all platforms
```bash
make build-all
```
