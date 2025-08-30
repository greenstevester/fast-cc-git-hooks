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

// ValidationError represents a validation failure.
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

// ValidationResult contains all validation errors.
type ValidationResult struct {
	Errors []error
	Valid  bool
}

// Error implements the error interface.
func (r *ValidationResult) Error() string {
	if r.Valid {
		return ""
	}

	messages := make([]string, 0, len(r.Errors))
	for _, err := range r.Errors {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// Validator validates commit messages according to configuration.
type Validator struct {
	config *config.Config
	parser *conventionalcommit.Parser
	// Compiled custom rules for performance.
	compiledRules map[string]*regexp.Regexp
	// Compiled ignore patterns for performance.
	compiledIgnorePatterns []*regexp.Regexp
}

// New creates a new validator with the given configuration.
func New(cfg *config.Config) (*Validator, error) {
	if cfg == nil {
		return nil, errors.New("config is required")
	}

	v := &Validator{
		config:        cfg,
		parser:        conventionalcommit.DefaultParser(),
		compiledRules: make(map[string]*regexp.Regexp),
	}

	// Compile custom rules.
	for _, rule := range cfg.CustomRules {
		re, err := regexp.Compile(rule.Pattern)
		if err != nil {
			return nil, fmt.Errorf("compiling custom rule %s: %w", rule.Name, err)
		}
		v.compiledRules[rule.Name] = re
	}

	// Compile JIRA ticket pattern if specified.
	if cfg.JIRATicketPattern != "" {
		re, err := regexp.Compile(cfg.JIRATicketPattern)
		if err != nil {
			return nil, fmt.Errorf("compiling JIRA ticket pattern: %w", err)
		}
		v.compiledRules["jira-pattern"] = re
	}

	// Compile ignore patterns.
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

// Validate validates a commit message.
func (v *Validator) Validate(ctx context.Context, message string) *ValidationResult {
	result := &ValidationResult{
		Errors: []error{},
		Valid:  true,
	}

	// Check for cancellation.
	if v.checkCancellation(ctx, result) {
		return result
	}

	// Check ignore patterns.
	if v.shouldIgnore(message) {
		return result
	}

	// Parse the commit message.
	commit, err := v.parser.Parse(message)
	if err != nil {
		v.addValidationError(result, "format", err.Error(), "")
		return result
	}

	// Run all validations.
	v.validateType(commit, result)
	v.validateScope(commit, result)
	v.validateSubjectLength(commit, result)
	v.validateBreakingChanges(commit, result)
	v.validateCustomRules(message, result)
	v.validateTicketRequirements(commit, result)

	return result
}

// ValidateFile validates commit messages from a file.
func (v *Validator) ValidateFile(ctx context.Context, path string) (*ValidationResult, error) {
	// Read commit message from file.
	content, err := readFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading commit file: %w", err)
	}

	// Remove comment lines (lines starting with #).
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

// shouldIgnore checks if a message matches any ignore pattern.
func (v *Validator) shouldIgnore(message string) bool {
	for _, re := range v.compiledIgnorePatterns {
		if re.MatchString(message) {
			return true
		}
	}
	return false
}

// validateTicketRequirements validates ticket reference requirements.
func (v *Validator) validateTicketRequirements(commit *conventionalcommit.Commit, result *ValidationResult) {
	v.validateJiraTicketRequired(commit, result)
	v.validateTicketRefRequired(commit, result)
	v.validateJiraTicketPattern(commit, result)
	v.validateJiraProjectPrefixes(commit, result)
}

// validateJiraTicketRequired checks if JIRA ticket is required.
func (v *Validator) validateJiraTicketRequired(commit *conventionalcommit.Commit, result *ValidationResult) {
	if v.config.RequireJIRATicket && !commit.HasJIRATicket() {
		v.addValidationError(result, "ticket", "JIRA ticket reference is required", "")
	}
}

// validateTicketRefRequired checks if any ticket reference is required.
func (v *Validator) validateTicketRefRequired(commit *conventionalcommit.Commit, result *ValidationResult) {
	if v.config.RequireTicketRef && !commit.HasTicketRefs() {
		v.addValidationError(result, "ticket", "ticket reference is required", "")
	}
}

// validateJiraTicketPattern validates JIRA ticket format if present.
func (v *Validator) validateJiraTicketPattern(commit *conventionalcommit.Commit, result *ValidationResult) {
	if v.config.JIRATicketPattern == "" || !commit.HasJIRATicket() {
		return
	}

	re, exists := v.compiledRules["jira-pattern"]
	if !exists {
		return
	}

	jiraTickets := commit.GetJIRATickets()
	for _, ticket := range jiraTickets {
		if !re.MatchString(ticket.ID) {
			message := fmt.Sprintf("JIRA ticket '%s' does not match required pattern", ticket.ID)
			v.addValidationError(result, "ticket", message, ticket.ID)
		}
	}
}

// validateJiraProjectPrefixes validates JIRA project prefixes if specified.
func (v *Validator) validateJiraProjectPrefixes(commit *conventionalcommit.Commit, result *ValidationResult) {
	if len(v.config.JIRAProjects) == 0 || !commit.HasJIRATicket() {
		return
	}

	jiraTickets := commit.GetJIRATickets()
	for _, ticket := range jiraTickets {
		v.validateJiraProjectPrefix(ticket, result)
	}
}

// validateJiraProjectPrefix validates a single JIRA project prefix.
func (v *Validator) validateJiraProjectPrefix(ticket conventionalcommit.TicketRef, result *ValidationResult) {
	parts := strings.Split(ticket.ID, "-")
	if len(parts) < 2 {
		return // Skip malformed tickets.
	}

	projectPrefix := parts[0]
	if v.isProjectAllowed(projectPrefix) {
		return
	}

	message := fmt.Sprintf("JIRA project '%s' is not allowed (allowed: %s)",
		projectPrefix, strings.Join(v.config.JIRAProjects, ", "))
	v.addValidationError(result, "ticket", message, ticket.ID)
}

// isProjectAllowed checks if a project prefix is in the allowed list.
func (v *Validator) isProjectAllowed(projectPrefix string) bool {
	for _, allowedProject := range v.config.JIRAProjects {
		if projectPrefix == allowedProject {
			return true
		}
	}
	return false
}

// addValidationError adds a validation error to the result.
func (*Validator) addValidationError(result *ValidationResult, field, message, value string) {
	result.Valid = false
	result.Errors = append(result.Errors, &ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	})
}

// readFile reads the contents of a file.
func readFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Quick validation helper for simple use cases.
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

// checkCancellation checks if the context is canceled and updates the result.
func (*Validator) checkCancellation(ctx context.Context, result *ValidationResult) bool {
	select {
	case <-ctx.Done():
		result.Valid = false
		result.Errors = append(result.Errors, ctx.Err())
		return true
	default:
		return false
	}
}

// validateType validates the commit type.
func (v *Validator) validateType(commit *conventionalcommit.Commit, result *ValidationResult) {
	if commit.Type != "" && !v.config.HasType(commit.Type) {
		v.addValidationError(result, "type",
			fmt.Sprintf("invalid type (allowed: %s)", strings.Join(v.config.Types, ", ")),
			commit.Type)
	}
}

// validateScope validates the commit scope.
func (v *Validator) validateScope(commit *conventionalcommit.Commit, result *ValidationResult) {
	if v.config.ScopeRequired && commit.Scope == "" {
		v.addValidationError(result, "scope", "scope is required", "")
	} else if commit.Scope != "" && !v.config.HasScope(commit.Scope) {
		v.addValidationError(result, "scope",
			fmt.Sprintf("invalid scope (allowed: %s)", strings.Join(v.config.Scopes, ", ")),
			commit.Scope)
	}
}

// validateSubjectLength validates the subject line length.
func (v *Validator) validateSubjectLength(commit *conventionalcommit.Commit, result *ValidationResult) {
	header := commit.Header()
	if len(header) > v.config.MaxSubjectLength {
		v.addValidationError(result, "subject",
			fmt.Sprintf("exceeds maximum length of %d characters", v.config.MaxSubjectLength),
			fmt.Sprintf("%d characters", len(header)))
	}
}

// validateBreakingChanges validates breaking change rules.
func (v *Validator) validateBreakingChanges(commit *conventionalcommit.Commit, result *ValidationResult) {
	if commit.Breaking && !v.config.AllowBreakingChanges {
		v.addValidationError(result, "breaking", "breaking changes are not allowed", "")
	}
}

// validateCustomRules applies custom validation rules.
func (v *Validator) validateCustomRules(message string, result *ValidationResult) {
	for _, rule := range v.config.CustomRules {
		re := v.compiledRules[rule.Name]
		if !re.MatchString(message) {
			msg := rule.Message
			if msg == "" {
				msg = fmt.Sprintf("failed custom rule: %s", rule.Name)
			}
			v.addValidationError(result, "custom", msg, "")
		}
	}
}
