// Package ccgen provides core commit message generation functionality
// shared between cc and ccc commands
package ccgen

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/greenstevester/fast-cc-git-hooks/internal/banner"
)

const (
	MaxSubjectLength  = 50
	MaxBodyLineLength = 72
)

// ChangeType represents a detected change in the repository
type ChangeType struct {
	Type        string
	Scope       string
	Description string
	Files       []string
	Priority    int
}

// JiraManager interface for JIRA ticket management
type JiraManager interface {
	GetCurrentJiraTicket() (string, error)
}

// Options configures the commit generation behavior
type Options struct {
	NoVerify    bool
	Execute     bool
	Copy        bool
	Verbose     bool
	JiraManager JiraManager
}

// Result contains the generated commit message and any additional information
type Result struct {
	Message    string
	Changes    []ChangeType
	GitCommand string
	HasChanges bool
}

// Generator handles commit message generation
type Generator struct {
	options Options
}

// New creates a new commit message generator with the given options
func New(opts Options) *Generator {
	return &Generator{
		options: opts,
	}
}

// Generate analyzes the repository and generates a commit message
func (g *Generator) Generate() (*Result, error) {
	fmt.Println()

	// Check if we're in a git repo
	fmt.Printf("Running `git rev-parse --git-dir`")
	if !g.isGitRepo() {
		fmt.Println(" ❌")
		return nil, fmt.Errorf("not a git repository")
	}
	fmt.Println(" ✅")

	// Get git status
	fmt.Printf("Running `git status --porcelain`")
	status, err := g.getGitStatus()
	if err != nil {
		fmt.Println(" ❌")
		return nil, fmt.Errorf("failed to get git status: %w", err)
	}
	fmt.Println(" ✅")

	if g.options.Verbose {
		fmt.Println("\n**Git status output:**")
		fmt.Printf("```\n%s```\n", status)
	}

	// Add all changes
	fmt.Printf("Running `git add .`")
	if addErr := g.addAllChanges(); addErr != nil {
		fmt.Println(" ❌")
		return nil, fmt.Errorf("failed to add changes: %w", addErr)
	}
	fmt.Println(" ✅")

	fmt.Println()
	// Perform advanced git analysis using comprehensive algorithm
	if banner.UseASCII() {
		fmt.Println("## Performing Advanced Git Analysis")
	} else {
		fmt.Println("## 🔬 Performing Advanced Git Analysis")
	}
	fmt.Println()

	// Use advanced git analysis algorithm
	gitAnalysis, err := g.performAdvancedGitAnalysis()
	if err != nil {
		return nil, fmt.Errorf("advanced git analysis failed: %w", err)
	}

	// Check if there are any changes
	if gitAnalysis.TotalFiles == 0 && strings.TrimSpace(gitAnalysis.StagedDiff) == "" {
		fmt.Println("\n**No changes detected** - nothing to commit")
		return &Result{HasChanges: false}, nil
	}

	// Convert advanced analysis to intelligent analyses
	intelligentAnalyses := g.getAdvancedChangeAnalyses(gitAnalysis)

	// Display advanced analysis results
	fmt.Printf("**Advanced Analysis Results:**\n")
	fmt.Printf("- Total files changed: %d\n", gitAnalysis.TotalFiles)
	fmt.Printf("- Total additions: +%d lines\n", gitAnalysis.TotalAdditions)
	fmt.Printf("- Total deletions: -%d lines\n", gitAnalysis.TotalDeletions)

	// Display directory statistics
	if len(gitAnalysis.DirStats) > 0 {
		fmt.Printf("- Directory distribution: ")
		var dirParts []string
		for dir, percent := range gitAnalysis.DirStats {
			dirParts = append(dirParts, fmt.Sprintf("%s (%.1f%%)", dir, percent))
		}
		fmt.Printf("%s\n", strings.Join(dirParts, ", "))
	}

	// Display file summaries
	if len(gitAnalysis.FileSummaries) > 0 {
		fmt.Printf("- File operations: %s\n", strings.Join(gitAnalysis.FileSummaries, ", "))
	}

	// Display modified functions
	if len(gitAnalysis.ModifiedFunctions) > 0 {
		fmt.Printf("- Modified functions: %s\n", strings.Join(gitAnalysis.ModifiedFunctions, ", "))
	}

	if gitAnalysis.CommitPatterns != nil && len(gitAnalysis.RecentCommits) > 0 {
		fmt.Printf("- Recent commit style: %s\n", gitAnalysis.CommitPatterns.PreferredStyle)
		fmt.Printf("- Average commit length: %d chars\n", gitAnalysis.CommitPatterns.AverageLength)
	}
	fmt.Printf("\n**Found %d change type(s):**\n\n", len(intelligentAnalyses))

	for i, analysis := range intelligentAnalyses {
		fmt.Printf("%d. **%s", i+1, analysis.ChangeType)
		if analysis.Scope != "" {
			fmt.Printf("(%s)", analysis.Scope)
		}
		fmt.Printf("**: %s", analysis.Description)
		if len(analysis.Files) > 0 {
			fmt.Printf("\n   - File: `%s`", analysis.Files[0])
		}
		if analysis.Impact != "" {
			fmt.Printf("\n   - Impact: %s", analysis.Impact)
		}
		if g.options.Verbose {
			// Show detailed statistics in verbose mode
			if stat, exists := gitAnalysis.FileStats[analysis.FilePath]; exists {
				fmt.Printf("\n   - Statistics: +%d/-%d lines, Type: %s",
					stat.Additions, stat.Deletions, stat.ChangeType)
			}
			if analysis.Context != "" {
				fmt.Printf("\n   - Context: %s", analysis.Context)
			}
			if len(analysis.Details) > 0 {
				fmt.Printf("\n   - Details:")
				for _, detail := range analysis.Details {
					fmt.Printf("\n     • %s", detail)
				}
			}
		}
		fmt.Printf("\n\n")
	}

	// Check for JIRA ticket
	if g.options.JiraManager != nil {
		if ticket, err := g.options.JiraManager.GetCurrentJiraTicket(); err == nil && ticket != "" {
			fmt.Printf("**JIRA Ticket:** `%s` (will be included in commit)\n\n", ticket)
		} else {
			fmt.Printf("**JIRA Ticket:** None set (use `cc set-jira CGC-1234` to set one)\n\n")
		}
	}

	// Generate Claude-style commit message using repository patterns
	message := g.generateClaudeStyleCommitMessageWithPatterns(intelligentAnalyses, gitAnalysis.CommitPatterns)

	// Also maintain backward compatibility by converting to old format for result
	changes := g.convertToLegacyFormat(intelligentAnalyses)

	// Build git command
	gitCommand := g.buildGitCommand(message)

	return &Result{
		Message:    message,
		Changes:    changes,
		GitCommand: gitCommand,
		HasChanges: true,
	}, nil
}

// ExecuteCommit commits the changes with the generated message
func (g *Generator) ExecuteCommit(message string) error {
	args := []string{"commit", "-m", message}
	if g.options.NoVerify {
		args = append(args, "--no-verify")
	}

	cmd := exec.Command("git", args...) // #nosec G204 - args are validated git commands
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// CopyToClipboard copies the git command to the system clipboard
func (g *Generator) CopyToClipboard(gitCommand string) error {
	return clipboard.WriteAll(gitCommand)
}

// PrintResult displays the result to the user
func (g *Generator) PrintResult(result *Result) {
	if !result.HasChanges {
		fmt.Println("**No changes to commit**")
		return
	}

	// Display the commit message in a code block
	fmt.Printf("```\n%s\n```\n\n", result.Message)

	if g.options.Copy {
		if err := g.CopyToClipboard(result.GitCommand); err != nil {
			fmt.Printf("❌ Failed to copy to clipboard: %v\n", err)
		} else {
			if banner.UseASCII() {
				fmt.Printf("✅ Git commit command copied to clipboard!\n\n")
			} else {
				fmt.Printf("✅ Git commit command copied to clipboard!\n\n")
			}
		}
	}

	if g.options.Execute {
		if err := g.ExecuteCommit(result.Message); err != nil {
			fmt.Printf("❌ Failed to commit: %v\n", err)
			return
		}
		fmt.Printf("✅ Commit created successfully!\n")
	}
}

// buildGitCommand builds the full git commit command string
func (g *Generator) buildGitCommand(message string) string {
	cmd := fmt.Sprintf("git commit -m %q", message)
	if g.options.NoVerify {
		cmd += " --no-verify"
	}
	return cmd
}

// isGitRepo checks if we're in a git repository
func (g *Generator) isGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

// getGitStatus gets git status output
func (g *Generator) getGitStatus() (string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	return string(output), err
}

// addAllChanges adds all changes to staging
func (g *Generator) addAllChanges() error {
	cmd := exec.Command("git", "add", ".")
	return cmd.Run()
}

// convertToLegacyFormat converts intelligent analyses to legacy ChangeType format for compatibility
func (g *Generator) convertToLegacyFormat(analyses []*IntelligentChangeAnalysis) []ChangeType {
	changes := make([]ChangeType, 0, len(analyses))

	for _, analysis := range analyses {
		change := ChangeType{
			Type:        analysis.ChangeType,
			Scope:       analysis.Scope,
			Description: analysis.Description,
			Files:       analysis.Files,
			Priority:    analysis.Priority,
		}
		changes = append(changes, change)
	}

	return changes
}
