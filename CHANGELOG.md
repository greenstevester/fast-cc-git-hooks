# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial implementation of fast conventional commits git hooks
- High-performance commit message parser with ~395ns validation time
- YAML-based configuration system with custom validation rules
- Cross-platform git hook installation and management
- CLI application with install, uninstall, validate, init, and version commands
- Docker container support with multi-stage builds
- Comprehensive test suite with benchmarks
- CI/CD pipelines with GitHub Actions
- GoReleaser configuration for automated releases
- Multi-platform builds (Linux, macOS, Windows) for AMD64 and ARM64
- Homebrew, Scoop, and Linux package manager support
- Documentation with usage examples and integration guides

### Features
- **Fast Validation**: Microsecond-level commit message validation
- **Flexible Configuration**: YAML config with custom types, scopes, and rules
- **Safe Installation**: Automatic backup of existing git hooks
- **Modern Go**: Built with Go 1.24+ using latest language features
- **Zero Dependencies**: Core functionality requires no external libraries
- **Context Support**: Cancellation and timeout support for validation
- **Parallel Processing**: Concurrent validation of multiple rules
- **Breaking Changes**: Full support for breaking change indicators
- **Custom Rules**: Regex-based custom validation patterns
- **Ignore Patterns**: Skip validation for specific commit patterns

### Performance
- Simple parse: ~395 ns/op with 4 allocations
- Complex parse: ~6.3 μs/op with 24 allocations
- Standard validation: ~873 ns/op with 8 allocations
- Parallel validation: ~721 ns/op with 8 allocations

### Supported Platforms
- Linux AMD64/ARM64
- macOS AMD64/ARM64 (Intel/Apple Silicon)
- Windows AMD64

### Standards Compliance
- [Conventional Commits 1.0.0](https://www.conventionalcommits.org/)
- [Semantic Versioning 2.0.0](https://semver.org/)
- [Keep a Changelog 1.0.0](https://keepachangelog.com/)

[Unreleased]: https://github.com/greenstevester/fast-cc-git-hooks/compare/v0.0.0...HEAD