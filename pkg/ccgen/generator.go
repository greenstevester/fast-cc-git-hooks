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

// Options configures the commit generation behavior
type Options struct {
	NoVerify bool
	Execute  bool
	Copy     bool
	Verbose  bool
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
	// Check if we're in a git repo
	if !g.isGitRepo() {
		return nil, fmt.Errorf("not a git repository")
	}

	// Get git status
	status, err := g.getGitStatus()
	if err != nil {
		return nil, fmt.Errorf("failed to get git status: %w", err)
	}

	if g.options.Verbose {
		fmt.Println("Git status:")
		fmt.Println(status)
		fmt.Println()
	}

	// Add all changes
	if err := g.addAllChanges(); err != nil {
		return nil, fmt.Errorf("failed to add changes: %w", err)
	}

	// Get staged diff
	diff, err := g.getStagedDiff()
	if err != nil {
		return nil, fmt.Errorf("failed to get diff: %w", err)
	}

	if strings.TrimSpace(diff) == "" {
		return &Result{HasChanges: false}, nil
	}

	// Analyze changes
	changes := g.analyzeDiff(diff)

	if g.options.Verbose {
		fmt.Println("Detected changes:")
		for _, change := range changes {
			fmt.Printf("- %s(%s): %s (files: %v)\n",
				change.Type, change.Scope, change.Description, change.Files)
		}
		fmt.Println()
	}

	// Generate commit message
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
		fmt.Println("No changes to commit")
		return
	}

	fmt.Println("─────────────────────────────────────────")
	fmt.Println("\n>>> based on your changes, cc created the following git commit message for you:")
	fmt.Println(result.Message)

	if g.options.Copy {
		if err := g.CopyToClipboard(result.GitCommand); err != nil {
			fmt.Printf("Failed to copy to clipboard: %v\n", err)
		} else {
			if banner.UseASCII() {
				fmt.Println("\n[COPY] Git commit command copied to clipboard! Use Ctrl+V to paste.")
			} else {
				fmt.Println("\n📋 Git commit command copied to clipboard! Use Ctrl+V to paste.")
			}
			fmt.Printf("Command: %s\n", result.GitCommand)
		}
	}

	if g.options.Execute {
		if err := g.ExecuteCommit(result.Message); err != nil {
			fmt.Printf("Failed to commit: %v\n", err)
			return
		}
		fmt.Println("Commit created successfully!")
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
