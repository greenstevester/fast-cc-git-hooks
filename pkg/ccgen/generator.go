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
	if banner.UseASCII() {
		fmt.Println("## Analyzing Git Repository")
	} else {
		fmt.Println("## 🔍 Analyzing Git Repository")
	}
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
	if err := g.addAllChanges(); err != nil {
		fmt.Println(" ❌")
		return nil, fmt.Errorf("failed to add changes: %w", err)
	}
	fmt.Println(" ✅")

	// Get staged diff
	fmt.Printf("Running `git diff --staged`")
	diff, err := g.getStagedDiff()
	if err != nil {
		fmt.Println(" ❌")
		return nil, fmt.Errorf("failed to get diff: %w", err)
	}
	fmt.Println(" ✅")

	if strings.TrimSpace(diff) == "" {
		fmt.Println("\n**No changes detected** - nothing to commit")
		return &Result{HasChanges: false}, nil
	}

	fmt.Println()
	// Analyze changes
	if banner.UseASCII() {
		fmt.Println("## Analyzing Changes")
	} else {
		fmt.Println("## 🧠 Analyzing Changes")
	}
	fmt.Println()

	changes := g.analyzeDiff(diff)

	fmt.Printf("**Found %d change type(s):**\n\n", len(changes))
	for i, change := range changes {
		fmt.Printf("%d. **%s", i+1, change.Type)
		if change.Scope != "" {
			fmt.Printf("(%s)", change.Scope)
		}
		fmt.Printf("**: %s", change.Description)
		if len(change.Files) > 0 {
			fmt.Printf("\n   - Files: `%s`", change.Files[0])
			if len(change.Files) > 1 {
				fmt.Printf(" (+%d more)", len(change.Files)-1)
			}
		}
		fmt.Printf("\n\n")
	}

	// Generate commit message
	if banner.UseASCII() {
		fmt.Println("## Generating Commit Message")
	} else {
		fmt.Println("## 📝 Generating Commit Message")
	}
	fmt.Println()

	// Check for JIRA ticket
	var jiraTicket string
	if g.options.JiraManager != nil {
		if ticket, err := g.options.JiraManager.GetCurrentJiraTicket(); err == nil && ticket != "" {
			jiraTicket = ticket
			fmt.Printf("**JIRA Ticket:** `%s` (will be included in commit)\n\n", jiraTicket)
		} else {
			fmt.Printf("**JIRA Ticket:** None set (use `cc set-jira CGC-1234` to set one)\n\n")
		}
	}

	message := g.GenerateCommitMessage(changes)

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

	if banner.UseASCII() {
		fmt.Println("## Generated Commit Message")
	} else {
		fmt.Println("## 🎯 Generated Commit Message")
	}
	fmt.Println()

	// Display the commit message in a code block
	fmt.Printf("```\n%s\n```\n\n", result.Message)

	if g.options.Copy {
		fmt.Printf("**Clipboard Action:** ")
		if err := g.CopyToClipboard(result.GitCommand); err != nil {
			fmt.Printf("❌ Failed to copy to clipboard: %v\n", err)
		} else {
			if banner.UseASCII() {
				fmt.Printf("✅ Git commit command copied to clipboard!\n\n")
				fmt.Printf("**Ready to paste:** `%s`\n\n", result.GitCommand)
				fmt.Println("Use **Ctrl+V** (or **Cmd+V** on Mac) to paste and execute.")
			} else {
				fmt.Printf("✅ Git commit command copied to clipboard!\n\n")
				fmt.Printf("**Ready to paste:** `%s`\n\n", result.GitCommand)
				fmt.Println("Use **Ctrl+V** (or **⌘+V** on Mac) to paste and execute.")
			}
		}
	}

	if g.options.Execute {
		fmt.Printf("**Executing commit:** ")
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

// getStagedDiff gets the staged diff
func (g *Generator) getStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--staged")
	output, err := cmd.Output()
	return string(output), err
}
