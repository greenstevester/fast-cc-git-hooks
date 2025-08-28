// Package validator provides commit message validation against conventional commit rules.
package validator

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/stevengreensill/fast-cc-git-hooks/internal/config"
	"github.com/stevengreensill/fast-cc-git-hooks/pkg/conventionalcommit"
)

// ValidationError represents a validation failure
type ValidationError struct {
	Field   string
	Message string
	Value   string
}

func (e *ValidationError) Error() string {
	if e.Value != "" {
		return fmt.Sprintf("%s: %s (got: %q)", e.Field, e.Message, e.Value)
	}
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationResult contains all validation errors
type ValidationResult struct {
	Errors []error
	Valid  bool
}

// Error implements the error interface
func (r *ValidationResult) Error() string {
	if r.Valid {
		return ""
	}

	var messages []string
	for _, err := range r.Errors {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// Validator validates commit messages according to configuration
type Validator struct {
	config *config.Config
	parser *conventionalcommit.Parser
	// Compiled custom rules for performance
	compiledRules map[string]*regexp.Regexp
	// Compiled ignore patterns for performance
	compiledIgnorePatterns []*regexp.Regexp
}

// New creates a new validator with the given configuration
func New(cfg *config.Config) (*Validator, error) {
	if cfg == nil {
		return nil, errors.New("config is required")
	}

	v := &Validator{
		config:        cfg,
		parser:        conventionalcommit.DefaultParser(),
		compiledRules: make(map[string]*regexp.Regexp),
	}

	// Compile custom rules
	for _, rule := range cfg.CustomRules {
		re, err := regexp.Compile(rule.Pattern)
		if err != nil {
			return nil, fmt.Errorf("compiling custom rule %s: %w", rule.Name, err)
		}
		v.compiledRules[rule.Name] = re
	}

	// Compile JIRA ticket pattern if specified
	if cfg.JIRATicketPattern != "" {
		re, err := regexp.Compile(cfg.JIRATicketPattern)
		if err != nil {
			return nil, fmt.Errorf("compiling JIRA ticket pattern: %w", err)
		}
		v.compiledRules["jira-pattern"] = re
	}

	// Compile ignore patterns
	v.compiledIgnorePatterns = make([]*regexp.Regexp, 0, len(cfg.IgnorePatterns))
	for _, pattern := range cfg.IgnorePatterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("compiling ignore pattern %q: %w", pattern, err)
		}
		v.compiledIgnorePatterns = append(v.compiledIgnorePatterns, re)
	}

	return v, nil
}

// Validate validates a commit message
func (v *Validator) Validate(ctx context.Context, message string) *ValidationResult {
	result := &ValidationResult{
		Errors: []error{},
		Valid:  true,
	}

	// Check for cancellation
	select {
	case <-ctx.Done():
		result.Valid = false
		result.Errors = append(result.Errors, ctx.Err())
		return result
	default:
	}

	// Check ignore patterns
	if v.shouldIgnore(message) {
		return result
	}

	// Parse the commit message
	commit, err := v.parser.Parse(message)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, &ValidationError{
			Field:   "format",
			Message: err.Error(),
		})
		return result
	}

	// Validate type
	if commit.Type != "" && !v.config.HasType(commit.Type) {
		result.Valid = false
		result.Errors = append(result.Errors, &ValidationError{
			Field:   "type",
			Message: fmt.Sprintf("invalid type (allowed: %s)", strings.Join(v.config.Types, ", ")),
			Value:   commit.Type,
		})
	}

	// Validate scope
	if v.config.ScopeRequired && commit.Scope == "" {
		result.Valid = false
		result.Errors = append(result.Errors, &ValidationError{
			Field:   "scope",
			Message: "scope is required",
		})
	} else if commit.Scope != "" && !v.config.HasScope(commit.Scope) {
		result.Valid = false
		result.Errors = append(result.Errors, &ValidationError{
			Field:   "scope",
			Message: fmt.Sprintf("invalid scope (allowed: %s)", strings.Join(v.config.Scopes, ", ")),
			Value:   commit.Scope,
		})
	}

	// Validate subject length
	header := commit.Header()
	if len(header) > v.config.MaxSubjectLength {
		result.Valid = false
		result.Errors = append(result.Errors, &ValidationError{
			Field:   "subject",
			Message: fmt.Sprintf("exceeds maximum length of %d characters", v.config.MaxSubjectLength),
			Value:   fmt.Sprintf("%d characters", len(header)),
		})
	}

	// Validate breaking changes
	if commit.Breaking && !v.config.AllowBreakingChanges {
		result.Valid = false
		result.Errors = append(result.Errors, &ValidationError{
			Field:   "breaking",
			Message: "breaking changes are not allowed",
		})
	}

	// Apply custom rules
	for _, rule := range v.config.CustomRules {
		re := v.compiledRules[rule.Name]
		if !re.MatchString(message) {
			result.Valid = false
			msg := rule.Message
			if msg == "" {
				msg = fmt.Sprintf("failed custom rule: %s", rule.Name)
			}
			result.Errors = append(result.Errors, &ValidationError{
				Field:   "custom",
				Message: msg,
			})
		}
	}

	// Validate ticket requirements
	v.validateTicketRequirements(commit, result)

	return result
}

// ValidateFile validates commit messages from a file
func (v *Validator) ValidateFile(ctx context.Context, path string) (*ValidationResult, error) {
	// Read commit message from file
	content, err := readFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading commit file: %w", err)
	}

	// Remove comment lines (lines starting with #)
	lines := strings.Split(content, "\n")
	var messageLines []string
	for _, line := range lines {
		if !strings.HasPrefix(strings.TrimSpace(line), "#") {
			messageLines = append(messageLines, line)
		}
	}

	message := strings.TrimSpace(strings.Join(messageLines, "\n"))
	if message == "" {
		return &ValidationResult{
			Valid: false,
			Errors: []error{&ValidationError{
				Field:   "message",
				Message: "commit message is empty",
			}},
		}, nil
	}

	return v.Validate(ctx, message), nil
}

// shouldIgnore checks if a message matches any ignore pattern
func (v *Validator) shouldIgnore(message string) bool {
	for _, re := range v.compiledIgnorePatterns {
		if re.MatchString(message) {
			return true
		}
	}
	return false
}

// validateTicketRequirements validates ticket reference requirements
func (v *Validator) validateTicketRequirements(commit *conventionalcommit.Commit, result *ValidationResult) {
	// Check if JIRA ticket is required
	if v.config.RequireJIRATicket && !commit.HasJIRATicket() {
		result.Valid = false
		result.Errors = append(result.Errors, &ValidationError{
			Field:   "ticket",
			Message: "JIRA ticket reference is required",
		})
	}

	// Check if any ticket reference is required
	if v.config.RequireTicketRef && !commit.HasTicketRefs() {
		result.Valid = false
		result.Errors = append(result.Errors, &ValidationError{
			Field:   "ticket",
			Message: "ticket reference is required",
		})
	}

	// Validate JIRA ticket format if present
	if v.config.JIRATicketPattern != "" && commit.HasJIRATicket() {
		re, exists := v.compiledRules["jira-pattern"]
		if exists {
			jiraTickets := commit.GetJIRATickets()
			for _, ticket := range jiraTickets {
				if !re.MatchString(ticket.ID) {
					result.Valid = false
					result.Errors = append(result.Errors, &ValidationError{
						Field:   "ticket",
						Message: fmt.Sprintf("JIRA ticket '%s' does not match required pattern", ticket.ID),
						Value:   ticket.ID,
					})
				}
			}
		}
	}

	// Validate JIRA project prefixes if specified
	if len(v.config.JIRAProjects) > 0 && commit.HasJIRATicket() {
		jiraTickets := commit.GetJIRATickets()
		for _, ticket := range jiraTickets {
			// Extract project prefix (part before the dash)
			parts := strings.Split(ticket.ID, "-")
			if len(parts) < 2 {
				continue // Skip malformed tickets
			}

			projectPrefix := parts[0]
			allowed := false
			for _, allowedProject := range v.config.JIRAProjects {
				if projectPrefix == allowedProject {
					allowed = true
					break
				}
			}

			if !allowed {
				result.Valid = false
				result.Errors = append(result.Errors, &ValidationError{
					Field:   "ticket",
					Message: fmt.Sprintf("JIRA project '%s' is not allowed (allowed: %s)", projectPrefix, strings.Join(v.config.JIRAProjects, ", ")),
					Value:   ticket.ID,
				})
			}
		}
	}
}

// readFile reads the contents of a file
func readFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Quick validation helper for simple use cases
func Quick(message string) error {
	cfg := config.Default()
	v, err := New(cfg)
	if err != nil {
		return err
	}

	result := v.Validate(context.Background(), message)
	if !result.Valid {
		return result
	}
	return nil
}
