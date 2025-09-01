package main

import (
	"strings"
	"testing"

	"github.com/greenstevester/fast-cc-git-hooks/pkg/ccgen"
)

func TestGenerateCommitMessage(t *testing.T) {
	tests := []struct {
		name     string
		changes  []ccgen.ChangeType
		expected string
	}{
		{
			name: "single change without footer",
			changes: []ccgen.ChangeType{
				{
					Type:        "feat",
					Scope:       "auth",
					Description: "add user login",
					Files:       []string{"auth.go"},
					Priority:    1,
				},
			},
			expected: "feat(auth): add user login",
		},
		{
			name: "multiple changes without footer",
			changes: []ccgen.ChangeType{
				{
					Type:        "feat",
					Scope:       "auth",
					Description: "add user login",
					Files:       []string{"auth.go"},
					Priority:    1,
				},
				{
					Type:        "fix",
					Scope:       "api",
					Description: "resolve validation error",
					Files:       []string{"validator.go"},
					Priority:    2,
				},
			},
			expected: "feat(auth): add user login\n\nChanges include:\n- Add user login (auth.go)\n- Resolve validation error (validator.go)",
		},
		{
			name:     "empty changes",
			changes:  []ccgen.ChangeType{},
			expected: "chore: update files",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ccgen.New(ccgen.Options{}).GenerateCommitMessage(tt.changes)

			// Verify the footer is not present
			if strings.Contains(result, "Generated with [Claude Code]") {
				t.Errorf("Expected no Claude Code footer in commit message, but found one")
			}
			if strings.Contains(result, "Co-Authored-By: Claude") {
				t.Errorf("Expected no Co-Authored-By footer in commit message, but found one")
			}

			// For the expected output, we only check the main part (not the full body for multi-change)
			if tt.name == "multiple changes without footer" {
				// For multiple changes, just verify it starts correctly and has no footer
				if !strings.HasPrefix(result, "feat(auth): add user login") {
					t.Errorf("Expected message to start with 'feat(auth): add user login', got: %s", result)
				}
				if !strings.Contains(result, "Changes include:") {
					t.Errorf("Expected message to contain 'Changes include:', got: %s", result)
				}
			} else {
				if result != tt.expected {
					t.Errorf("Expected: %s, got: %s", tt.expected, result)
				}
			}
		})
	}
}

func TestGenerateCommitMessageFooterRemoval(t *testing.T) {
	// Test specifically that no footer is added
	changes := []ccgen.ChangeType{
		{
			Type:        "feat",
			Scope:       "test",
			Description: "add test functionality",
			Files:       []string{"test.go"},
			Priority:    1,
		},
	}

	result := ccgen.New(ccgen.Options{}).GenerateCommitMessage(changes)

	// Verify no footer components are present
	footerElements := []string{
		"🤖 Generated with [Claude Code]",
		"Co-Authored-By: Claude",
		"claude.ai/code",
		"noreply@anthropic.com",
	}

	for _, element := range footerElements {
		if strings.Contains(result, element) {
			t.Errorf("Found unwanted footer element '%s' in commit message: %s", element, result)
		}
	}
}

func TestGenerateCommitMessageFormat(t *testing.T) {
	// Test that commit messages follow conventional commit format
	tests := []struct {
		name       string
		changes    []ccgen.ChangeType
		wantPrefix string
	}{
		{
			name: "feat with scope",
			changes: []ccgen.ChangeType{
				{
					Type:        "feat",
					Scope:       "auth",
					Description: "add user authentication",
					Files:       []string{"auth.go"},
					Priority:    1,
				},
			},
			wantPrefix: "feat(auth):",
		},
		{
			name: "fix without scope",
			changes: []ccgen.ChangeType{
				{
					Type:        "fix",
					Scope:       "",
					Description: "resolve memory leak",
					Files:       []string{"main.go"},
					Priority:    2,
				},
			},
			wantPrefix: "fix:",
		},
		{
			name: "docs with scope",
			changes: []ccgen.ChangeType{
				{
					Type:        "docs",
					Scope:       "api",
					Description: "update API documentation",
					Files:       []string{"README.md"},
					Priority:    6,
				},
			},
			wantPrefix: "docs(api):",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ccgen.New(ccgen.Options{}).GenerateCommitMessage(tt.changes)

			if !strings.HasPrefix(result, tt.wantPrefix) {
				t.Errorf("Expected commit message to start with '%s', got: %s", tt.wantPrefix, result)
			}

			// Verify conventional commit format (no footer)
			lines := strings.Split(result, "\n")
			firstLine := lines[0]

			// Check that first line doesn't exceed recommended length
			if len(firstLine) > 72 {
				t.Errorf("Subject line too long (%d chars): %s", len(firstLine), firstLine)
			}

			// Ensure no footer is present anywhere in the message
			if strings.Contains(result, "Generated with") || strings.Contains(result, "Co-Authored-By") {
				t.Errorf("Commit message should not contain footer, got: %s", result)
			}
		})
	}
}
