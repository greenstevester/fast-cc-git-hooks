// Package config provides configuration management for fast-cc-hooks.
package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	// DefaultConfigFile is the default configuration filename
	DefaultConfigFile = ".fast-cc-hooks.yaml"
	// DefaultMaxSubjectLength is the default maximum subject line length
	DefaultMaxSubjectLength = 72
)

// Config represents the complete configuration for fast-cc-hooks
type Config struct {
	// Types defines allowed commit types
	Types []string `yaml:"types"`
	// Scopes defines allowed scopes (empty means any scope allowed)
	Scopes []string `yaml:"scopes,omitempty"`
	// ScopeRequired indicates if scope is mandatory
	ScopeRequired bool `yaml:"scope_required"`
	// MaxSubjectLength defines maximum subject line length
	MaxSubjectLength int `yaml:"max_subject_length"`
	// AllowBreakingChanges permits breaking change indicators (!)
	AllowBreakingChanges bool `yaml:"allow_breaking_changes"`
	// CustomRules defines additional validation rules
	CustomRules []CustomRule `yaml:"custom_rules,omitempty"`
	// IgnorePatterns defines patterns to skip validation
	IgnorePatterns []string `yaml:"ignore_patterns,omitempty"`
	
	// Ticket reference validation
	// RequireJIRATicket requires JIRA ticket references in commits
	RequireJIRATicket bool `yaml:"require_jira_ticket"`
	// RequireTicketRef requires any type of ticket reference in commits
	RequireTicketRef bool `yaml:"require_ticket_ref"`
	// JIRATicketPattern defines a regex pattern for valid JIRA tickets
	JIRATicketPattern string `yaml:"jira_ticket_pattern,omitempty"`
	// JIRAProjects defines allowed JIRA project prefixes
	JIRAProjects []string `yaml:"jira_projects,omitempty"`
}

// CustomRule defines a custom validation rule
type CustomRule struct {
	Name    string `yaml:"name"`
	Pattern string `yaml:"pattern"`
	Message string `yaml:"message"`
}

// DefaultTypes returns the standard conventional commit types
func DefaultTypes() []string {
	return []string{
		"feat",     // New feature
		"fix",      // Bug fix
		"docs",     // Documentation changes
		"style",    // Code style changes (formatting, semicolons, etc)
		"refactor", // Code refactoring
		"test",     // Adding or modifying tests
		"chore",    // Maintenance tasks
		"perf",     // Performance improvements
		"ci",       // CI/CD changes
		"build",    // Build system changes
		"revert",   // Reverts a previous commit
	}
}

// Default returns a default configuration
func Default() *Config {
	return &Config{
		Types:                DefaultTypes(),
		Scopes:               []string{},
		ScopeRequired:        false,
		MaxSubjectLength:     DefaultMaxSubjectLength,
		AllowBreakingChanges: true,
		CustomRules:          []CustomRule{},
		IgnorePatterns:       []string{},
	}
}

// Load reads configuration from a file
func Load(path string) (*Config, error) {
	if path == "" {
		path = DefaultConfigFile
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Return default config if file doesn't exist
		return Default(), nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening config file: %w", err)
	}
	defer file.Close()

	return Parse(file)
}

// Parse parses configuration from an io.Reader
func Parse(r io.Reader) (*Config, error) {
	cfg := Default()
	
	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(cfg); err != nil {
		if !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("parsing config: %w", err)
		}
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

// Save writes configuration to a file
func (c *Config) Save(path string) error {
	if path == "" {
		path = DefaultConfigFile
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating config file: %w", err)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	defer encoder.Close()

	if err := encoder.Encode(c); err != nil {
		return fmt.Errorf("encoding config: %w", err)
	}

	return nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if len(c.Types) == 0 {
		return errors.New("at least one commit type must be defined")
	}

	if c.MaxSubjectLength <= 0 {
		return errors.New("max_subject_length must be positive")
	}

	// Validate custom rules
	for i, rule := range c.CustomRules {
		if rule.Name == "" {
			return fmt.Errorf("custom rule %d: name is required", i)
		}
		if rule.Pattern == "" {
			return fmt.Errorf("custom rule %s: pattern is required", rule.Name)
		}
	}

	return nil
}

// HasType checks if a commit type is allowed
func (c *Config) HasType(t string) bool {
	for _, allowed := range c.Types {
		if allowed == t {
			return true
		}
	}
	return false
}

// HasScope checks if a scope is allowed (returns true if no scopes defined)
func (c *Config) HasScope(s string) bool {
	if len(c.Scopes) == 0 {
		return true // Any scope allowed if none defined
	}
	for _, allowed := range c.Scopes {
		if allowed == s {
			return true
		}
	}
	return false
}