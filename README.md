# fast-cc-git-hooks

A high-performance Git hook system for enforcing [Conventional Commits](https://www.conventionalcommits.org/) with extensive configuration options and blazing-fast validation.

## Features

- **Fast Validation**: Optimized for performance with minimal dependencies
- **Flexible Configuration**: YAML-based configuration with sensible defaults
- **Custom Rules**: Define project-specific validation rules using regex patterns
- **Cross-Platform**: Works on Linux, macOS, and Windows
- **Easy Installation**: Simple CLI commands for hook management
- **Go 1.21+**: Leverages modern Go features for better performance

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/stevengreensill/fast-cc-git-hooks.git
cd fast-cc-git-hooks

# Build and install
make build
make install  # Installs to /usr/local/bin
```

### Using Go

```bash
go install github.com/stevengreensill/fast-cc-git-hooks/cmd/fast-cc-hooks@latest
```

## Quick Start

1. **Initialize configuration** (optional):
```bash
fast-cc-hooks init
```

2. **Install hooks in your repository**:
```bash
fast-cc-hooks install
```

3. **Make commits** using conventional format:
```bash
git commit -m "feat: add new feature"
git commit -m "fix(api): resolve authentication issue"
git commit -m "docs: update README"
```

## Usage

### Commands

```bash
fast-cc-hooks [flags] <command> [args]

Commands:
  install    Install git hooks in current repository
  uninstall  Remove git hooks from current repository
  validate   Validate a commit message
  init       Create default configuration file
  version    Show version information

Global flags:
  -v         Verbose output
  -config    Path to config file (default: .fast-cc-hooks.yaml)
```

### Examples

#### Validate a message manually:
```bash
fast-cc-hooks validate "feat: add new feature"
fast-cc-hooks validate --file COMMIT_MSG
```

#### Install with force (overwrite existing hooks):
```bash
fast-cc-hooks install --force
```

#### Use custom configuration:
```bash
fast-cc-hooks --config=custom-config.yaml validate "fix: bug"
```

## Configuration

Configuration is done via `.fast-cc-hooks.yaml` file. If no configuration file exists, default settings are used.

### Example Configuration

```yaml
# Allowed commit types
types:
  - feat     # New feature
  - fix      # Bug fix
  - docs     # Documentation changes
  - style    # Code style changes
  - refactor # Code refactoring
  - test     # Testing
  - chore    # Maintenance
  - perf     # Performance improvements
  - ci       # CI/CD changes
  - build    # Build system
  - revert   # Revert commits

# Allowed scopes (empty = any scope allowed)
scopes:
  - api
  - web
  - cli

# Whether scope is required
scope_required: false

# Maximum subject line length
max_subject_length: 72

# Allow breaking changes
allow_breaking_changes: true

# Custom validation rules
custom_rules:
  - name: jira-ticket
    pattern: '\[JIRA-\d+\]'
    message: 'Must reference a JIRA ticket'

# Ignore patterns
ignore_patterns:
  - '^Merge'
  - '^WIP:'
```

### Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `types` | []string | Standard types | Allowed commit types |
| `scopes` | []string | [] (any) | Allowed scopes |
| `scope_required` | bool | false | Whether scope is mandatory |
| `max_subject_length` | int | 72 | Maximum subject line length |
| `allow_breaking_changes` | bool | true | Allow breaking change indicators |
| `require_jira_ticket` | bool | false | Require JIRA ticket in all commits |
| `require_ticket_ref` | bool | false | Require any ticket reference |
| `jira_ticket_pattern` | string | `^[A-Z]{3,4}-\d+$` | Custom JIRA ticket pattern |
| `jira_projects` | []string | [] (any) | Allowed JIRA project prefixes |
| `custom_rules` | []Rule | [] | Custom validation rules |
| `ignore_patterns` | []string | [] | Patterns to skip validation |

### Ticket Validation Configuration

```yaml
# Require JIRA tickets in all commits
require_jira_ticket: true

# Require any type of ticket reference  
require_ticket_ref: false

# Custom JIRA pattern (optional)
jira_ticket_pattern: "^[A-Z]{3}-\\d+$"  # Only 3-letter prefixes

# Allowed JIRA project codes (optional)
jira_projects:
  - CGC
  - PROJ
  - WORK

# Alternative: Use custom rules for ticket validation
custom_rules:
  - name: "require-jira"
    pattern: "\\b[A-Z]{3,4}-\\d+\\b"
    message: "Commit must include a JIRA ticket (e.g., CGC-1234)"
```

## Conventional Commits Format

```
<type>[optional scope]: [TICKET-ID] <description>

[optional body]

[optional footer(s)]
```

### Ticket Reference Support

The tool automatically detects and validates ticket references:

- **JIRA tickets**: `ABC-123`, `PROJ-456`, `WORK-789` (3-4 letter project prefixes)
- **GitHub issues**: `#123`, `GH-456` 
- **Linear tickets**: Can be configured for specific project prefixes
- **Generic format**: `[TICKET-123]`

Multiple ticket types can be referenced in the same commit.

### Examples:

```bash
# Standard conventional commits
feat: add user authentication
feat(auth): implement OAuth2 integration
fix!: correct critical security vulnerability
docs(api): update endpoint documentation
refactor(core): reorganize module structure

# With JIRA ticket references
feat: CGC-1234 add user authentication
feat(auth): PROJ-789 implement OAuth2 integration
fix: ABC-456 correct critical security vulnerability

# With GitHub issue references  
docs(api): update endpoint documentation #123
fix: resolve authentication bug GH-456

# Multi-line with tickets
feat: CGC-1234 add shopping cart functionality

This implements a shopping cart with:
- Add/remove items
- Calculate totals  
- Apply discounts

Closes: CGC-1234
Fixes: #456
```

## Development

### Building

```bash
make build        # Build for current platform
make build-all    # Build for all platforms
make test         # Run tests
make bench        # Run benchmarks
make coverage     # Generate coverage report
```

### Testing

The project includes comprehensive tests and benchmarks:

```bash
# Run all tests
go test ./...

# Run with race detection
go test -race ./...

# Run benchmarks
go test -bench=. ./...

# Generate coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Performance

The validator is optimized for speed with:
- Compiled regex patterns cached at startup
- Minimal allocations in hot paths
- Efficient string operations
- Zero dependencies for core functionality

Benchmark results (M1 Mac):
```
BenchmarkParser_Parse-8              500000      2341 ns/op
BenchmarkParser_ParseSimple-8       2000000       872 ns/op
BenchmarkValidator_Validate-8       1000000      1053 ns/op
```

## License

MIT License - see LICENSE file for details

## Contributing

Contributions are welcome! Please ensure:
1. Code follows Go best practices
2. All tests pass
3. New features include tests
4. Commits follow conventional format
