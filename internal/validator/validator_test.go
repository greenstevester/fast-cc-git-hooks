package validator

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/greenstevester/fast-cc-git-hooks/internal/config"
)

func TestValidator_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *config.Config
		message string
		valid   bool
		errors  int
	}{
		{
			name:    "valid conventional commit",
			config:  config.Default(),
			message: "feat: add new feature",
			valid:   true,
			errors:  0,
		},
		{
			name:    "valid commit with scope",
			config:  config.Default(),
			message: "fix(api): resolve bug",
			valid:   true,
			errors:  0,
		},
		{
			name: "invalid type",
			config: &config.Config{
				Types:            []string{"feat", "fix"},
				MaxSubjectLength: 72,
			},
			message: "invalid: not a valid type",
			valid:   false,
			errors:  1,
		},
		{
			name: "scope required but missing",
			config: &config.Config{
				Types:            config.DefaultTypes(),
				ScopeRequired:    true,
				MaxSubjectLength: 72,
			},
			message: "feat: missing scope",
			valid:   false,
			errors:  1,
		},
		{
			name: "scope not in allowed list",
			config: &config.Config{
				Types:            config.DefaultTypes(),
				Scopes:           []string{"api", "web", "db"},
				MaxSubjectLength: 72,
			},
			message: "feat(cli): add CLI feature",
			valid:   false,
			errors:  1,
		},
		{
			name: "subject too long",
			config: &config.Config{
				Types:            config.DefaultTypes(),
				MaxSubjectLength: 50,
			},
			message: "feat: this is a very long commit message that exceeds the maximum allowed length",
			valid:   false,
			errors:  1,
		},
		{
			name: "breaking change not allowed",
			config: &config.Config{
				Types:                config.DefaultTypes(),
				MaxSubjectLength:     72,
				AllowBreakingChanges: false,
			},
			message: "feat!: breaking change",
			valid:   false,
			errors:  1,
		},
		{
			name: "breaking change allowed",
			config: &config.Config{
				Types:                config.DefaultTypes(),
				MaxSubjectLength:     72,
				AllowBreakingChanges: true,
			},
			message: "feat!: breaking change",
			valid:   true,
			errors:  0,
		},
		{
			name: "ignored pattern",
			config: &config.Config{
				Types:            config.DefaultTypes(),
				MaxSubjectLength: 72,
				IgnorePatterns:   []string{"^WIP:", "^Merge"},
			},
			message: "WIP: work in progress",
			valid:   true,
			errors:  0,
		},
		{
			name: "multiple errors",
			config: &config.Config{
				Types:                []string{"feat", "fix"},
				ScopeRequired:        true,
				MaxSubjectLength:     30,
				AllowBreakingChanges: false,
			},
			message: "invalid!: this message has multiple problems and is too long",
			valid:   false,
			errors:  4, // invalid type, missing scope, too long, breaking change
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := New(tt.config)
			if err != nil {
				t.Fatalf("Failed to create validator: %v", err)
			}

			result := v.Validate(context.Background(), tt.message)

			if result.Valid != tt.valid {
				t.Errorf("Validate() valid = %v, want %v", result.Valid, tt.valid)
			}

			if len(result.Errors) != tt.errors {
				t.Errorf("Validate() errors = %d, want %d", len(result.Errors), tt.errors)
				for _, err := range result.Errors {
					t.Logf("  Error: %v", err)
				}
			}
		})
	}
}

func TestValidator_ValidateFile(t *testing.T) {
	// Create temporary directory for test files.
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		content string
		valid   bool
	}{
		{
			name:    "valid commit from file",
			content: "feat: add new feature\n\nThis is the body",
			valid:   true,
		},
		{
			name:    "commit with comments",
			content: "fix: bug fix\n# This is a comment\n\nBody text\n# Another comment",
			valid:   true,
		},
		{
			name:    "empty file after comments",
			content: "# Only comments\n# Nothing else",
			valid:   false,
		},
	}

	cfg := config.Default()
	v, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file.
			file := filepath.Join(tmpDir, "COMMIT_MSG")
			if err := os.WriteFile(file, []byte(tt.content), 0o600); err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			result, err := v.ValidateFile(context.Background(), file)
			if err != nil {
				t.Fatalf("ValidateFile() error = %v", err)
			}

			if result.Valid != tt.valid {
				t.Errorf("ValidateFile() valid = %v, want %v", result.Valid, tt.valid)
			}
		})
	}
}

func TestValidator_CustomRules(t *testing.T) {
	cfg := &config.Config{
		Types:            config.DefaultTypes(),
		MaxSubjectLength: 72,
		CustomRules: []config.CustomRule{
			{
				Name:    "jira-ticket",
				Pattern: `\[JIRA-\d+\]`,
				Message: "commit must reference a JIRA ticket",
			},
		},
	}

	v, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	tests := []struct {
		name    string
		message string
		valid   bool
	}{
		{
			name:    "with JIRA ticket",
			message: "feat: [JIRA-123] add feature",
			valid:   true,
		},
		{
			name:    "without JIRA ticket",
			message: "feat: add feature",
			valid:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(context.Background(), tt.message)
			if result.Valid != tt.valid {
				t.Errorf("Validate() valid = %v, want %v", result.Valid, tt.valid)
				for _, err := range result.Errors {
					t.Logf("  Error: %v", err)
				}
			}
		})
	}
}

func TestValidator_ContextCancellation(t *testing.T) {
	cfg := config.Default()
	v, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	result := v.Validate(ctx, "feat: test message")
	if result.Valid {
		t.Error("Expected validation to fail with canceled context")
	}
}

func BenchmarkValidator_Validate(b *testing.B) {
	cfg := config.Default()
	v, err := New(cfg)
	if err != nil {
		b.Fatal(err)
	}
	ctx := context.Background()
	message := "feat(scope): add new feature with a reasonably long description"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := v.Validate(ctx, message)
		if !result.Valid {
			b.Fatal(result.Error())
		}
	}
}

func BenchmarkValidator_ValidateWithCustomRules(b *testing.B) {
	cfg := &config.Config{
		Types:            config.DefaultTypes(),
		MaxSubjectLength: 72,
		CustomRules: []config.CustomRule{
			{
				Name:    "rule1",
				Pattern: `\[JIRA-\d+\]`,
			},
			{
				Name:    "rule2",
				Pattern: `feat|fix|docs`,
			},
		},
	}
	v, err := New(cfg)
	if err != nil {
		b.Fatal(err)
	}
	ctx := context.Background()
	message := "feat: [JIRA-123] add new feature"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := v.Validate(ctx, message)
		if !result.Valid {
			b.Fatal(result.Error())
		}
	}
}
