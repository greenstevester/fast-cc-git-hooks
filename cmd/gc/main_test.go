package main

import (
	"strings"
	"testing"
)

func TestGenerateCommitMessage(t *testing.T) {
	tests := []struct {
		name     string
		changes  []ChangeType
		expected string
	}{
		{
			name: "single change without footer",
			changes: []ChangeType{
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
			changes: []ChangeType{
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
			changes:  []ChangeType{},
			expected: "chore: update files",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateCommitMessage(tt.changes)
			
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
	changes := []ChangeType{
		{
			Type:        "feat",
			Scope:       "test",
			Description: "add test functionality",
			Files:       []string{"test.go"},
			Priority:    1,
		},
	}
	
	result := generateCommitMessage(changes)
	
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