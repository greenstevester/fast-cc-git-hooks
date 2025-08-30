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

func TestValidator_CGCCommitFormat(t *testing.T) {
	// Test for CGC-style commit messages with format: "feat(db): CGC-1425 Added new database"
	tests := []struct {
		name    string
		config  *config.Config
		message string
		valid   bool
		errors  int
	}{
		{
			name: "valid CGC format with ticket number",
			config: &config.Config{
				Types:            config.DefaultTypes(),
				Scopes:           []string{"db", "api", "web", "cli", "core"},
				MaxSubjectLength: 100,
				CustomRules: []config.CustomRule{
					{
						Name:    "cgc-ticket",
						Pattern: `CGC-\d{1,5}`,
						Message: "commit must reference a CGC ticket number",
					},
				},
			},
			message: "feat(db): CGC-1425 Added new database",
			valid:   true,
			errors:  0,
		},
		{
			name: "valid CGC format with different ticket number",
			config: &config.Config{
				Types:            config.DefaultTypes(),
				Scopes:           []string{"db", "api", "web", "cli", "core"},
				MaxSubjectLength: 100,
				CustomRules: []config.CustomRule{
					{
						Name:    "cgc-ticket",
						Pattern: `CGC-\d{1,5}`,
						Message: "commit must reference a CGC ticket number",
					},
				},
			},
			message: "fix(api): CGC-99 Fixed authentication issue",
			valid:   true,
			errors:  0,
		},
		{
			name: "missing CGC ticket number",
			config: &config.Config{
				Types:            config.DefaultTypes(),
				Scopes:           []string{"db", "api", "web", "cli", "core"},
				MaxSubjectLength: 100,
				CustomRules: []config.CustomRule{
					{
						Name:    "cgc-ticket",
						Pattern: `CGC-\d{1,5}`,
						Message: "commit must reference a CGC ticket number",
					},
				},
			},
			message: "feat(db): Added new database",
			valid:   false,
			errors:  1,
		},
		{
			name: "invalid CGC ticket format - wrong prefix",
			config: &config.Config{
				Types:            config.DefaultTypes(),
				Scopes:           []string{"db", "api", "web", "cli", "core"},
				MaxSubjectLength: 100,
				CustomRules: []config.CustomRule{
					{
						Name:    "cgc-ticket",
						Pattern: `CGC-\d{1,5}`,
						Message: "commit must reference a CGC ticket number",
					},
				},
			},
			message: "feat(db): ABC-1425 Added new database",
			valid:   false,
			errors:  1,
		},
		{
			name: "invalid CGC ticket format - no hyphen",
			config: &config.Config{
				Types:            config.DefaultTypes(),
				Scopes:           []string{"db", "api", "web", "cli", "core"},
				MaxSubjectLength: 100,
				CustomRules: []config.CustomRule{
					{
						Name:    "cgc-ticket",
						Pattern: `CGC-\d{1,5}`,
						Message: "commit must reference a CGC ticket number",
					},
				},
			},
			message: "feat(db): CGC1425 Added new database",
			valid:   false,
			errors:  1,
		},
		{
			name: "invalid CGC ticket format - too many digits",
			config: &config.Config{
				Types:            config.DefaultTypes(),
				Scopes:           []string{"db", "api", "web", "cli", "core"},
				MaxSubjectLength: 100,
				CustomRules: []config.CustomRule{
					{
						Name:    "cgc-ticket",
						Pattern: `^[a-z]+\([a-z]+\): CGC-\d{1,5} `,
						Message: "commit must start with type(scope): CGC-XXXXX followed by space",
					},
				},
			},
			message: "feat(db): CGC-142599 Added new database",
			valid:   false,
			errors:  1,
		},
		{
			name: "CGC ticket in wrong position",
			config: &config.Config{
				Types:            config.DefaultTypes(),
				Scopes:           []string{"db", "api", "web", "cli", "core"},
				MaxSubjectLength: 100,
				CustomRules: []config.CustomRule{
					{
						Name:    "cgc-ticket-position",
						Pattern: `^[a-z]+\([a-z]+\): CGC-\d{1,5} `,
						Message: "CGC ticket must appear immediately after type(scope): ",
					},
				},
			},
			message: "feat(db): Added new database CGC-1425",
			valid:   false,
			errors:  1,
		},
		{
			name: "invalid scope with CGC format",
			config: &config.Config{
				Types:            config.DefaultTypes(),
				Scopes:           []string{"db", "api", "web"},
				MaxSubjectLength: 100,
				CustomRules: []config.CustomRule{
					{
						Name:    "cgc-ticket",
						Pattern: `CGC-\d{1,5}`,
						Message: "commit must reference a CGC ticket number",
					},
				},
			},
			message: "feat(invalid): CGC-1425 Added new database",
			valid:   false,
			errors:  1,
		},
		{
			name: "invalid type with CGC format",
			config: &config.Config{
				Types:            []string{"feat", "fix", "docs"},
				Scopes:           []string{"db", "api", "web"},
				MaxSubjectLength: 100,
				CustomRules: []config.CustomRule{
					{
						Name:    "cgc-ticket",
						Pattern: `CGC-\d{1,5}`,
						Message: "commit must reference a CGC ticket number",
					},
				},
			},
			message: "invalid(db): CGC-1425 Added new database",
			valid:   false,
			errors:  1,
		},
		{
			name: "subject too long with CGC format",
			config: &config.Config{
				Types:            config.DefaultTypes(),
				Scopes:           []string{"db", "api", "web"},
				MaxSubjectLength: 50,
				CustomRules: []config.CustomRule{
					{
						Name:    "cgc-ticket",
						Pattern: `CGC-\d{1,5}`,
						Message: "commit must reference a CGC ticket number",
					},
				},
			},
			message: "feat(db): CGC-1425 Added new database with extremely long description that exceeds the maximum allowed length",
			valid:   false,
			errors:  1,
		},
		{
			name: "multiple CGC tickets in message",
			config: &config.Config{
				Types:            config.DefaultTypes(),
				Scopes:           []string{"db", "api", "web"},
				MaxSubjectLength: 100,
				CustomRules: []config.CustomRule{
					{
						Name:    "cgc-ticket",
						Pattern: `CGC-\d{1,5}`,
						Message: "commit must reference a CGC ticket number",
					},
				},
			},
			message: "feat(db): CGC-1425 CGC-1426 Added new database",
			valid:   true,
			errors:  0,
		},
		{
			name: "strict CGC format - exactly one ticket at start",
			config: &config.Config{
				Types:            config.DefaultTypes(),
				Scopes:           []string{"db", "api", "web"},
				MaxSubjectLength: 100,
				CustomRules: []config.CustomRule{
					{
						Name:    "cgc-ticket-strict",
						Pattern: `^[a-z]+\([a-z]+\): CGC-\d{1,5} [A-Z]`,
						Message: "commit must follow format: type(scope): CGC-XXXXX Description",
					},
				},
			},
			message: "feat(db): CGC-1425 Added new database",
			valid:   true,
			errors:  0,
		},
		{
			name: "CGC format with breaking change",
			config: &config.Config{
				Types:                config.DefaultTypes(),
				Scopes:               []string{"db", "api", "web"},
				MaxSubjectLength:     100,
				AllowBreakingChanges: true,
				CustomRules: []config.CustomRule{
					{
						Name:    "cgc-ticket",
						Pattern: `CGC-\d{1,5}`,
						Message: "commit must reference a CGC ticket number",
					},
				},
			},
			message: "feat(db)!: CGC-1425 Breaking database change",
			valid:   true,
			errors:  0,
		},
		{
			name: "CGC format missing scope",
			config: &config.Config{
				Types:            config.DefaultTypes(),
				ScopeRequired:    true,
				MaxSubjectLength: 100,
				CustomRules: []config.CustomRule{
					{
						Name:    "cgc-ticket",
						Pattern: `CGC-\d{1,5}`,
						Message: "commit must reference a CGC ticket number",
					},
				},
			},
			message: "feat: CGC-1425 Added new database",
			valid:   false,
			errors:  1,
		},
		{
			name: "CGC format with lowercase description start",
			config: &config.Config{
				Types:            config.DefaultTypes(),
				Scopes:           []string{"db", "api", "web"},
				MaxSubjectLength: 100,
				CustomRules: []config.CustomRule{
					{
						Name:    "cgc-ticket",
						Pattern: `CGC-\d{1,5}`,
						Message: "commit must reference a CGC ticket number",
					},
					{
						Name:    "description-capitalized",
						Pattern: `CGC-\d{1,5} [A-Z]`,
						Message: "description after ticket must start with capital letter",
					},
				},
			},
			message: "feat(db): CGC-1425 added new database",
			valid:   false,
			errors:  1,
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
