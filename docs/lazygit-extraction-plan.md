# Lazygit Functionality Extraction Plan

## Executive Summary

This document outlines a comprehensive plan to extract relevant functionality from [lazygit](https://github.com/greenstevester/lazygit) into the fast-cc-git-hooks project. The goal is to improve code quality, reduce technical debt, and enhance maintainability by leveraging lazygit's battle-tested Git command abstractions.

**Key Recommendation:** Extract lazygit's git command builder and execution layer to replace current `exec.Command` approach, improving testability, error handling, and code maintainability.

---

## Current State Analysis

### Fast-cc-git-hooks Architecture

**Project Structure:**
```
├── cmd/
│   ├── ccg/          # Commit message generator
│   ├── ccdo/         # ???
│   └── fcgh/         # Main CLI
├── pkg/
│   ├── ccgen/        # Core commit generation logic
│   ├── conventionalcommit/  # Parser
│   ├── jira/         # JIRA integration
│   └── semantic/     # Semantic analysis
└── internal/
    ├── config/       # YAML configuration
    ├── hooks/        # Git hook management
    └── validator/    # Validation logic
```

**Current Git Interaction Approach:**

The project currently uses direct `exec.Command()` calls for all Git operations:

```go
// From pkg/ccgen/advanced_git_analyzer.go
cmd := exec.Command("git", "diff", "--stat", "HEAD~1", "HEAD")
output, err := cmd.Output()
```

**Pain Points:**
1. **No abstraction layer** - Direct shell execution throughout codebase
2. **Limited testability** - Difficult to mock git commands
3. **Error handling** - Basic error handling, no structured error types
4. **Code duplication** - Similar git command patterns repeated
5. **No command validation** - Arguments not validated before execution
6. **Limited debugging** - No built-in logging or command inspection

### Lazygit Architecture

**Relevant Packages:**

1. **`pkg/commands/git_commands/`** (55 files)
   - Comprehensive Git command wrappers
   - Structured command building
   - Type-safe API

2. **`pkg/commands/git_config/`**
   - Git configuration management
   - Cached config access

3. **`pkg/commands/oscommands/`**
   - OS command abstraction
   - Command execution layer

4. **`pkg/commands/models/`**
   - Domain models (Commit, Branch, File, etc.)

5. **`pkg/commands/patch/`**
   - Patch creation and manipulation

**Key Strengths:**
- ✅ Well-tested (extensive test coverage)
- ✅ Production-proven (thousands of users)
- ✅ Go 1.25+ compatible
- ✅ Minimal dependencies
- ✅ Excellent error handling
- ✅ Comprehensive git operation support

---

## Recommended Extractions

### Priority 1: Core Command Infrastructure (High Value, Low Risk)

#### 1.1 Git Command Builder

**Source:** `pkg/commands/git_commands/git_command_builder.go`

**Functionality:**
- Fluent builder pattern for Git commands
- Method chaining: `NewGitCmd("commit").Arg("-m", "msg").Dir("/repo")`
- Conditional arguments: `ArgIf()`, `ArgIfElse()`
- Smart argument ordering (config flags, directory flags, command args)
- Output formats: `ToArgv()`, `ToString()`

**Current Usage in Fast-cc-git-hooks:**
```go
// Before (current approach)
cmd := exec.Command("git", "diff", "--stat", "HEAD~1", "HEAD")

// After (with builder)
cmd := git.NewGitCmd("diff").
    Arg("--stat").
    Arg("HEAD~1", "HEAD").
    Build()
```

**Benefits:**
- 🎯 Type-safe command construction
- 🎯 Eliminates argument ordering bugs
- 🎯 Improves code readability
- 🎯 Easier to test and mock

**Extraction Effort:** Medium (2-3 days)

#### 1.2 Git Command Runner

**Source:** `pkg/commands/git_commands/git_cmd_obj_runner.go`

**Functionality:**
- Executes built commands with proper error handling
- Stream-based output (stdout/stderr separation)
- Context support for cancellation
- Credential handling for protected operations
- Environment variable management

**Benefits:**
- 🎯 Centralized error handling
- 🎯 Consistent logging
- 🎯 Better debugging capabilities
- 🎯 Timeout and cancellation support

**Extraction Effort:** Medium (2-3 days)

### Priority 2: Git Operation Wrappers (High Value, Medium Risk)

#### 2.1 Commit Operations

**Source:** `pkg/commands/git_commands/commit.go`

**Relevant Functions:**
- `GetCommitMessage(hash)` - Retrieve commit messages
- `GetCommitSubject(hash)` - Extract subject line
- `GetCommitAuthor(hash)` - Get author information
- `GetCommitsOneline()` - One-line commit history
- `AmendHead()` - Amend operations

**Current Usage:**
```go
// In advanced_git_analyzer.go:
cmd := exec.Command("git", "log", "--oneline", "-10")
output, err := cmd.Output()
// Manual parsing required
```

**Improved Approach:**
```go
commits, err := git.GetCommitsOneline(10)
for _, commit := range commits {
    // Structured commit objects with Hash, Message, etc.
}
```

**Benefits:**
- 🎯 Eliminates manual parsing
- 🎯 Structured commit data
- 🎯 Type-safe operations
- 🎯 Reduced code duplication

**Extraction Effort:** Medium-High (3-4 days)

#### 2.2 Diff Operations

**Source:** `pkg/commands/git_commands/diff.go`

**Relevant Functions:**
- Diff generation with various formats
- Word diff support
- Numstat parsing
- Directory statistics

**Benefits:**
- 🎯 Replaces 9 different git diff calls in `advanced_git_analyzer.go`
- 🎯 Structured diff output
- 🎯 Built-in parsing

**Extraction Effort:** Medium (2-3 days)

### Priority 3: Configuration & Infrastructure (Medium Value, Low Risk)

#### 3.1 Git Configuration

**Source:** `pkg/commands/git_commands/config.go`

**Functionality:**
- `GetCoreEditor()` - Get configured editor
- `GetRemoteURL()` - Remote URL lookup
- `GetCoreCommentChar()` - Comment character config
- Cached configuration access

**Use Cases:**
- Hook configuration
- User preference detection
- Editor integration (future)

**Extraction Effort:** Low-Medium (1-2 days)

#### 3.2 File Operations

**Source:** `pkg/commands/git_commands/file.go`, `file_loader.go`

**Functionality:**
- File status detection
- Staged/unstaged file lists
- File change type detection
- Working tree status

**Benefits:**
- 🎯 Better file analysis in ccgen
- 🎯 Improved change type detection
- 🎯 Structured file metadata

**Extraction Effort:** Medium (2-3 days)

### Priority 4: Advanced Features (Future Consideration)

#### 4.1 Patch Management
**Source:** `pkg/commands/patch/`
- Useful for future interactive features
- Low priority for current needs

#### 4.2 Worktree Support
**Source:** `pkg/commands/git_commands/worktree.go`
- Minimal value for hook use case
- Consider only if user requests

#### 4.3 Branch Operations
**Source:** `pkg/commands/git_commands/branch.go`
- Potential for scope detection improvements
- Future enhancement opportunity

---

## Implementation Plan

### Phase 1: Foundation (Week 1-2)

**Goals:** Establish core infrastructure without breaking existing functionality

**Tasks:**
1. Create `internal/gitcmd/` package structure
2. Extract and adapt `git_command_builder.go`
   - Simplify for fast-cc-git-hooks needs
   - Remove UI-specific dependencies
   - Add comprehensive tests
3. Extract and adapt command runner
   - Strip out interactive features
   - Keep core execution logic
   - Add context support
4. Create adapter layer for existing code
   - Allow gradual migration
   - No breaking changes to public APIs

**Deliverables:**
- ✅ `internal/gitcmd/builder.go`
- ✅ `internal/gitcmd/runner.go`
- ✅ `internal/gitcmd/builder_test.go`
- ✅ `internal/gitcmd/runner_test.go`
- ✅ Documentation

**Success Criteria:**
- All existing tests pass
- New package has >90% test coverage
- Benchmark shows no performance regression

### Phase 2: Migration (Week 3-4)

**Goals:** Migrate `pkg/ccgen/advanced_git_analyzer.go` to use new infrastructure

**Tasks:**
1. Replace git diff operations
   - `getFileStatistics()` → builder-based
   - `getChangeTypes()` → builder-based
   - `getNumStats()` → builder-based
   - `getWordDiff()` → builder-based
2. Replace git log operations
   - `analyzeRecentCommitPatterns()` → builder-based
3. Add structured parsing helpers
4. Improve error handling throughout
5. Add integration tests

**Deliverables:**
- ✅ Refactored `advanced_git_analyzer.go`
- ✅ New integration tests
- ✅ Updated documentation

**Success Criteria:**
- All tests pass (unit + integration)
- Code coverage maintained >90%
- No functional regressions
- Performance benchmarks show improvement or parity

### Phase 3: Enhancement (Week 5-6)

**Goals:** Extract higher-level git operations

**Tasks:**
1. Extract commit operations
   - Create `internal/gitcmd/commit.go`
   - Structured commit parsing
   - Message/author utilities
2. Extract diff operations
   - Create `internal/gitcmd/diff.go`
   - Structured diff output
   - Parsing helpers
3. Extract config operations
   - Create `internal/gitcmd/config.go`
   - Configuration caching
4. Create domain models
   - `internal/gitcmd/models/commit.go`
   - `internal/gitcmd/models/file.go`
   - `internal/gitcmd/models/diff.go`

**Deliverables:**
- ✅ Complete `internal/gitcmd/` package
- ✅ Domain models
- ✅ Comprehensive test suite
- ✅ API documentation
- ✅ Migration guide

**Success Criteria:**
- Package is production-ready
- Documentation is complete
- All legacy exec.Command calls replaced
- Performance benchmarks show improvement

### Phase 4: Polish & Documentation (Week 7)

**Goals:** Finalize implementation and document learnings

**Tasks:**
1. Code review and cleanup
2. Performance optimization
3. Add examples
4. Update CLAUDE.md
5. Create extraction retrospective

**Deliverables:**
- ✅ Updated project documentation
- ✅ Code examples
- ✅ Performance report
- ✅ Lessons learned document

---

## Detailed Extraction Examples

### Example 1: Command Builder Pattern

**File to Extract:** `pkg/commands/git_commands/git_command_builder.go`

**Adaptation Strategy:**

```go
// internal/gitcmd/builder.go
package gitcmd

import "os/exec"

// Builder constructs git commands using a fluent interface
type Builder struct {
    args      []string
    dir       string
    configArgs []string
}

// New creates a new git command builder
func New(subcommand string) *Builder {
    return &Builder{
        args: []string{subcommand},
        configArgs: make([]string, 0),
    }
}

// Arg adds arguments to the command
func (b *Builder) Arg(args ...string) *Builder {
    b.args = append(b.args, args...)
    return b
}

// ArgIf conditionally adds arguments
func (b *Builder) ArgIf(condition bool, args ...string) *Builder {
    if condition {
        b.args = append(b.args, args...)
    }
    return b
}

// Dir sets the working directory
func (b *Builder) Dir(dir string) *Builder {
    b.dir = dir
    return b
}

// Config adds a git config override
func (b *Builder) Config(key, value string) *Builder {
    b.configArgs = append(b.configArgs, "-c", key+"="+value)
    return b
}

// Build creates the exec.Cmd
func (b *Builder) Build() *exec.Cmd {
    // Combine: git + config flags + subcommand + args
    fullArgs := append([]string{"git"}, b.configArgs...)
    fullArgs = append(fullArgs, b.args...)

    cmd := exec.Command(fullArgs[0], fullArgs[1:]...)
    if b.dir != "" {
        cmd.Dir = b.dir
    }
    return cmd
}

// ToArgv returns the command as a string slice
func (b *Builder) ToArgv() []string {
    args := append([]string{"git"}, b.configArgs...)
    return append(args, b.args...)
}
```

**Usage in fast-cc-git-hooks:**

```go
// Before:
cmd := exec.Command("git", "diff", "--stat", "HEAD~1", "HEAD")
if g.hasPreviousCommits() {
    cmd = exec.Command("git", "diff", "--stat", "HEAD~1", "HEAD")
} else {
    cmd = exec.Command("git", "diff", "--stat", "--staged")
}

// After:
cmd := gitcmd.New("diff").
    Arg("--stat").
    ArgIf(g.hasPreviousCommits(), "HEAD~1", "HEAD").
    ArgIf(!g.hasPreviousCommits(), "--staged").
    Build()
```

### Example 2: Structured Diff Operations

**New Package:** `internal/gitcmd/diff.go`

```go
package gitcmd

import (
    "context"
    "fmt"
    "strconv"
    "strings"
)

// DiffOptions configures diff operation
type DiffOptions struct {
    Staged          bool
    Cached          bool
    FromCommit      string
    ToCommit        string
    WordDiff        bool
    Stat            bool
    NumStat         bool
    DirStat         bool
    FunctionContext bool
}

// NumStatResult represents a parsed numstat line
type NumStatResult struct {
    Additions int
    Deletions int
    Filename  string
}

// GetNumStat retrieves numerical statistics for changed files
func GetNumStat(ctx context.Context, opts DiffOptions) ([]NumStatResult, error) {
    builder := New("diff").Arg("--numstat")

    if opts.Staged {
        builder.Arg("--staged")
    } else if opts.FromCommit != "" && opts.ToCommit != "" {
        builder.Arg(opts.FromCommit, opts.ToCommit)
    }

    cmd := builder.Build()
    output, err := cmd.CombinedOutput()
    if err != nil {
        return nil, fmt.Errorf("git diff --numstat failed: %w", err)
    }

    return parseNumStat(string(output))
}

func parseNumStat(output string) ([]NumStatResult, error) {
    var results []NumStatResult
    lines := strings.Split(strings.TrimSpace(output), "\n")

    for _, line := range lines {
        if strings.TrimSpace(line) == "" {
            continue
        }

        parts := strings.Fields(line)
        if len(parts) < 3 {
            continue
        }

        additions, err1 := strconv.Atoi(parts[0])
        deletions, err2 := strconv.Atoi(parts[1])

        if err1 == nil && err2 == nil {
            results = append(results, NumStatResult{
                Additions: additions,
                Deletions: deletions,
                Filename:  parts[2],
            })
        }
    }

    return results, nil
}
```

**Usage:**

```go
// Before:
cmd := exec.Command("git", "diff", "--cached", "--numstat")
output, err := cmd.Output()
// Manual parsing...

// After:
results, err := gitcmd.GetNumStat(ctx, gitcmd.DiffOptions{Staged: true})
for _, stat := range results {
    fmt.Printf("%s: +%d -%d\n", stat.Filename, stat.Additions, stat.Deletions)
}
```

---

## Benefits & ROI Analysis

### Immediate Benefits (Phase 1-2)

1. **Code Quality**
   - ✅ Eliminate 50+ direct `exec.Command` calls
   - ✅ Centralized error handling
   - ✅ Consistent logging/debugging
   - ✅ Type-safe command construction

2. **Testability**
   - ✅ Mockable git operations
   - ✅ Easier integration testing
   - ✅ Command validation in tests
   - ✅ Better test coverage

3. **Maintainability**
   - ✅ Clearer code intent
   - ✅ Reduced duplication
   - ✅ Easier refactoring
   - ✅ Better documentation

### Long-term Benefits (Phase 3-4)

1. **Performance**
   - ⚡ Command caching opportunities
   - ⚡ Reduced subprocess overhead (batching)
   - ⚡ Configuration caching
   - ⚡ Optimized parsing

2. **Features**
   - 🚀 Foundation for interactive features
   - 🚀 Better error messages
   - 🚀 Enhanced commit analysis
   - 🚀 Advanced git operations support

3. **Developer Experience**
   - 💡 Easier to add git operations
   - 💡 Self-documenting API
   - 💡 Better IDE support (autocomplete)
   - 💡 Clearer abstractions

### Estimated ROI

**Investment:**
- Development time: 6-7 weeks (1 developer)
- Testing time: 1 week
- Documentation: 1 week
- **Total: 8-9 weeks**

**Return:**
- **Maintenance time reduction:** 30-40% (fewer bugs, clearer code)
- **Feature development acceleration:** 50% faster for git-related features
- **Bug reduction:** 60-70% fewer git command-related bugs
- **Onboarding time:** 40% faster for new contributors

**Break-even:** ~3-4 months after completion

---

## Risks & Mitigation

### Risk 1: Breaking Changes

**Probability:** Medium
**Impact:** High

**Mitigation:**
- Create adapter layer for gradual migration
- Maintain backward compatibility during migration
- Comprehensive integration tests
- Feature flags for new code paths
- Rollback plan for each phase

### Risk 2: Performance Regression

**Probability:** Low
**Impact:** Medium

**Mitigation:**
- Benchmark all changes
- Performance tests in CI
- Optimize hot paths
- Profile before/after
- Keep fallback to direct exec if needed

### Risk 3: Incomplete Extraction

**Probability:** Medium
**Impact:** Medium

**Mitigation:**
- Start with smallest useful subset
- Incremental extraction approach
- Clear success criteria per phase
- Regular checkpoints and reviews

### Risk 4: Dependency Management

**Probability:** Low
**Impact:** Low

**Mitigation:**
- Extract code, don't add dependency
- Minimize external dependencies
- Adapt code to fast-cc-git-hooks patterns
- Remove lazygit-specific features

---

## Success Metrics

### Code Quality Metrics

- [ ] Test coverage ≥ 90% for `internal/gitcmd/`
- [ ] Zero direct `exec.Command("git", ...)` calls in `pkg/`
- [ ] Cyclomatic complexity reduced by 30%
- [ ] Code duplication reduced by 50%

### Performance Metrics

- [ ] Benchmark performance within 5% of baseline
- [ ] Memory allocations reduced by 20%
- [ ] Average command execution overhead < 100μs

### Reliability Metrics

- [ ] Git operation error rate reduced by 60%
- [ ] Test flakiness reduced to zero
- [ ] CI success rate ≥ 99%

### Developer Experience Metrics

- [ ] Lines of code for new git operation: -50%
- [ ] Time to add new git feature: -40%
- [ ] Developer onboarding time: -30%

---

## Alternatives Considered

### Alternative 1: Use go-git Library

**Pros:**
- Pure Go implementation
- No git binary dependency
- Rich API

**Cons:**
- ❌ Large dependency (~100MB)
- ❌ Slower than native git
- ❌ Incomplete git feature coverage
- ❌ Different behavior than git CLI
- ❌ Not aligned with project philosophy (minimal dependencies)

**Decision:** Rejected - Too heavy, conflicts with project goals

### Alternative 2: Keep Current Approach

**Pros:**
- Zero effort
- No risk

**Cons:**
- ❌ Technical debt accumulates
- ❌ Maintainability issues grow
- ❌ Testing difficulties persist
- ❌ Feature development slows

**Decision:** Rejected - Unsustainable long-term

### Alternative 3: Minimal Wrapper Only

**Pros:**
- Low effort
- Some benefits

**Cons:**
- ❌ Doesn't solve parsing problems
- ❌ Limited testability improvement
- ❌ Still requires manual parsing

**Decision:** Considered but insufficient - Go with full extraction

### Alternative 4: Full Lazygit Dependency

**Pros:**
- All features available
- Well maintained

**Cons:**
- ❌ Massive dependency (~50+ packages)
- ❌ UI dependencies unnecessary
- ❌ Conflicts with minimal dependency goal
- ❌ Overkill for needs

**Decision:** Rejected - Extract only what's needed

---

## Recommended Next Steps

### Immediate Actions (This Week)

1. **Review & Approve Plan**
   - Team review of this document
   - Stakeholder approval
   - Timeline confirmation

2. **Setup Development Branch**
   - Create `feature/lazygit-extraction` branch
   - Setup tracking issues
   - Create project board

3. **Prototype Phase 1**
   - Quick POC of builder pattern
   - Performance baseline measurements
   - Validate extraction approach

### Short-term (Next 2 Weeks)

1. **Begin Phase 1 Implementation**
   - Create `internal/gitcmd/` package
   - Extract builder
   - Write comprehensive tests

2. **Documentation**
   - Update CLAUDE.md
   - Create API documentation
   - Write migration guide

### Medium-term (Weeks 3-8)

1. **Execute Phases 2-4**
   - Follow implementation plan
   - Regular checkpoints
   - Continuous integration

2. **Monitoring**
   - Track success metrics
   - Performance monitoring
   - User feedback collection

---

## Appendix: File Extraction Checklist

### From Lazygit - Priority 1

- [ ] `pkg/commands/git_commands/git_command_builder.go`
- [ ] `pkg/commands/git_commands/git_cmd_obj_builder.go`
- [ ] `pkg/commands/git_commands/git_cmd_obj_runner.go`
- [ ] `pkg/commands/git_commands/common.go` (utilities)

### From Lazygit - Priority 2

- [ ] `pkg/commands/git_commands/commit.go`
- [ ] `pkg/commands/git_commands/commit_loader.go`
- [ ] `pkg/commands/git_commands/diff.go`
- [ ] `pkg/commands/git_commands/status.go`
- [ ] `pkg/commands/git_commands/file.go`
- [ ] `pkg/commands/git_commands/file_loader.go`

### From Lazygit - Priority 3

- [ ] `pkg/commands/git_commands/config.go`
- [ ] `pkg/commands/models/commit.go`
- [ ] `pkg/commands/models/file.go`
- [ ] Supporting test files

### To Create (Adapters/Glue)

- [ ] `internal/gitcmd/adapter.go` - Backward compatibility layer
- [ ] `internal/gitcmd/context.go` - Context utilities
- [ ] `internal/gitcmd/errors.go` - Error types
- [ ] `internal/gitcmd/testing.go` - Test utilities

---

## Conclusion

Extracting functionality from lazygit represents a strategic investment in the fast-cc-git-hooks codebase. The recommended phased approach minimizes risk while delivering incremental value. The builder pattern extraction alone (Phase 1) will provide significant benefits, with subsequent phases building on this foundation.

**Primary Recommendation:** Proceed with Phase 1-2 extraction (builder + migration) as the minimum viable improvement, then evaluate Phase 3-4 based on results and team capacity.

**Timeline:** 8-9 weeks for full implementation
**Risk Level:** Low-Medium (with proper planning)
**Expected ROI:** High (3-4 month break-even)

---

**Document Version:** 1.0
**Author:** Claude (AI Assistant)
**Date:** 2025-11-08
**Status:** Draft - Pending Review
