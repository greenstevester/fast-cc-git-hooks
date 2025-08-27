package config

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	
	if cfg == nil {
		t.Fatal("Default() returned nil")
	}
	
	if len(cfg.Types) == 0 {
		t.Error("Default config should have types")
	}
	
	if cfg.MaxSubjectLength != DefaultMaxSubjectLength {
		t.Errorf("Default MaxSubjectLength = %d, want %d", cfg.MaxSubjectLength, DefaultMaxSubjectLength)
	}
	
	if !cfg.AllowBreakingChanges {
		t.Error("Default should allow breaking changes")
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "valid config",
			config:  Default(),
			wantErr: false,
		},
		{
			name: "no types",
			config: &Config{
				Types:            []string{},
				MaxSubjectLength: 72,
			},
			wantErr: true,
		},
		{
			name: "negative max length",
			config: &Config{
				Types:            DefaultTypes(),
				MaxSubjectLength: -1,
			},
			wantErr: true,
		},
		{
			name: "invalid custom rule - no name",
			config: &Config{
				Types:            DefaultTypes(),
				MaxSubjectLength: 72,
				CustomRules: []CustomRule{
					{Pattern: "test"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid custom rule - no pattern",
			config: &Config{
				Types:            DefaultTypes(),
				MaxSubjectLength: 72,
				CustomRules: []CustomRule{
					{Name: "test"},
				},
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_HasType(t *testing.T) {
	cfg := &Config{
		Types: []string{"feat", "fix", "docs"},
	}
	
	tests := []struct {
		typ  string
		want bool
	}{
		{"feat", true},
		{"fix", true},
		{"docs", true},
		{"chore", false},
		{"", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.typ, func(t *testing.T) {
			if got := cfg.HasType(tt.typ); got != tt.want {
				t.Errorf("Config.HasType(%q) = %v, want %v", tt.typ, got, tt.want)
			}
		})
	}
}

func TestConfig_HasScope(t *testing.T) {
	tests := []struct {
		name   string
		scopes []string
		scope  string
		want   bool
	}{
		{
			name:   "empty scopes allows any",
			scopes: []string{},
			scope:  "anything",
			want:   true,
		},
		{
			name:   "defined scope exists",
			scopes: []string{"api", "web", "cli"},
			scope:  "api",
			want:   true,
		},
		{
			name:   "undefined scope",
			scopes: []string{"api", "web", "cli"},
			scope:  "db",
			want:   false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{Scopes: tt.scopes}
			if got := cfg.HasScope(tt.scope); got != tt.want {
				t.Errorf("Config.HasScope(%q) = %v, want %v", tt.scope, got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		want    *Config
		wantErr bool
	}{
		{
			name: "basic config",
			yaml: `
types:
  - feat
  - fix
max_subject_length: 50
scope_required: true
`,
			want: &Config{
				Types:                []string{"feat", "fix"},
				MaxSubjectLength:     50,
				ScopeRequired:        true,
				AllowBreakingChanges: true, // default
			},
		},
		{
			name: "config with custom rules",
			yaml: `
types:
  - feat
max_subject_length: 72
custom_rules:
  - name: jira
    pattern: '\[JIRA-\d+\]'
    message: Must include JIRA ticket
`,
			want: &Config{
				Types:                []string{"feat"},
				MaxSubjectLength:     72,
				AllowBreakingChanges: true,
				CustomRules: []CustomRule{
					{
						Name:    "jira",
						Pattern: `\[JIRA-\d+\]`,
						Message: "Must include JIRA ticket",
					},
				},
			},
		},
		{
			name: "empty yaml uses defaults",
			yaml: "",
			want: Default(),
		},
		{
			name: "invalid yaml",
			yaml: `
types: [
  invalid yaml
`,
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.yaml)
			got, err := Parse(reader)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				// Compare relevant fields
				if !reflect.DeepEqual(got.Types, tt.want.Types) {
					t.Errorf("Parse() Types = %v, want %v", got.Types, tt.want.Types)
				}
				if got.MaxSubjectLength != tt.want.MaxSubjectLength {
					t.Errorf("Parse() MaxSubjectLength = %d, want %d", got.MaxSubjectLength, tt.want.MaxSubjectLength)
				}
				if got.ScopeRequired != tt.want.ScopeRequired {
					t.Errorf("Parse() ScopeRequired = %v, want %v", got.ScopeRequired, tt.want.ScopeRequired)
				}
			}
		})
	}
}

func TestConfig_SaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")
	
	// Create config with custom values
	original := &Config{
		Types:                []string{"feat", "fix", "custom"},
		Scopes:               []string{"api", "web"},
		ScopeRequired:        true,
		MaxSubjectLength:     100,
		AllowBreakingChanges: false,
		CustomRules: []CustomRule{
			{
				Name:    "test-rule",
				Pattern: "test.*",
				Message: "Test message",
			},
		},
		IgnorePatterns: []string{"^WIP"},
	}
	
	// Save config
	if err := original.Save(configPath); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}
	
	// Verify file exists
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("Config file not created: %v", err)
	}
	
	// Load config
	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	// Compare
	if !reflect.DeepEqual(loaded, original) {
		t.Errorf("Loaded config differs from original\nGot: %+v\nWant: %+v", loaded, original)
	}
}

func TestLoad_NonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "non-existent.yaml")
	
	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() with non-existent file should return default config, got error: %v", err)
	}
	
	// Should return default config
	if !reflect.DeepEqual(cfg, Default()) {
		t.Error("Load() with non-existent file should return default config")
	}
}

func BenchmarkConfig_HasType(b *testing.B) {
	cfg := Default()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cfg.HasType("feat")
	}
}

func BenchmarkConfig_Validate(b *testing.B) {
	cfg := Default()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cfg.Validate()
	}
}