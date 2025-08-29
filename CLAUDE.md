# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Fast Conventional Commits Git Hooks is a high-performance Go application that enforces [Conventional Commits](https://www.conventionalcommits.org/) specification through git hooks. The project achieves microsecond-level validation performance while providing extensive configuration options.

## Architecture

```
├── cmd/fast-cc-hooks/          # Main CLI application
├── internal/                   # Private packages
│   ├── config/                # YAML configuration management
│   ├── hooks/                 # Git hook installation/management
│   └── validator/             # Core validation logic
├── pkg/conventionalcommit/     # Public reusable parser package
└── build/                     # Build artifacts
```

## Development Commands

### Building and Testing
```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests with race detection
make test

# Run benchmarks
make bench

# Generate coverage report
make coverage

# Format and lint code
make fmt lint

# Complete CI pipeline
make ci
```

### Git Hook Management
```bash
# Install pre-commit hook
make hook-install

# Uninstall pre-commit hook
make hook-uninstall

# Initialize configuration
make init-config
```

### Release and Packaging
```bash
# Create release snapshot
make release-snapshot

# Build Docker image
make docker-build
```

## Configuration

The application uses `.fast-cc-hooks.yaml` for configuration:
- **types**: Allowed commit types (feat, fix, docs, etc.)
- **scopes**: Allowed scopes (optional)
- **custom_rules**: Regex-based validation patterns
- **ignore_patterns**: Skip validation for specific commits
- **max_subject_length**: Subject line length limit

## Performance Requirements

- Parser should maintain ~400ns performance for simple commits
- Validator should process commits in under 1μs
- Memory allocations should be minimized in hot paths
- All regex patterns are compiled and cached at startup

## Testing Standards

- All packages must have comprehensive table-driven tests
- Benchmarks are required for performance-critical code
- Coverage should be >90% for core validation logic
- Context cancellation must be tested
- Cross-platform compatibility is verified in CI

## Code Style

- Use Go 1.24+ features (slices package, structured logging)
- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines  
- Error handling uses wrapped errors with context
- Public APIs are documented with examples
- Private functions focus on single responsibility

## Dependencies

The project minimizes external dependencies:
- Core functionality: zero dependencies
- Configuration: `gopkg.in/yaml.v3` only
- Testing: standard library only
- CLI: standard library `flag` package

## Release Process

Releases are automated via GoReleaser and support:
- Multi-platform binaries (Linux, macOS, Windows)
- Docker images with multi-arch support
- Homebrew formula updates
- Package manager distributions (deb, rpm, apk)
- Container signatures with cosign