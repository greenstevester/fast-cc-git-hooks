// Package conventionalcommit provides parsing and validation for conventional commit messages.
package conventionalcommit

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	// ErrInvalidFormat indicates the commit message doesn't match conventional format.
	ErrInvalidFormat = errors.New("invalid conventional commit format")
	// ErrEmptyMessage indicates the commit message is empty.
	ErrEmptyMessage = errors.New("empty commit message")
)

// Commit represents a parsed conventional commit message.
type Commit struct {
	TicketRefs  []TicketRef
	Type        string
	Scope       string
	Description string
	Body        string
	Footer      string
	Raw         string
	Breaking    bool
}

// TicketRef represents a ticket reference (e.g., JIRA ticket).
type TicketRef struct {
	Type string // e.g., "JIRA", "GITHUB", "LINEAR"
	ID   string // e.g., "PROJ-123", "#456", "ABC-789"
	Raw  string // Original reference as found in commit
}

// Parser provides conventional commit parsing with configurable options.
type Parser struct {
	// StrictMode enforces strict conventional commit format.
	StrictMode bool
	// AllowEmptyScope permits commits without scope.
	AllowEmptyScope bool
}

// DefaultParser returns a parser with default settings.
func DefaultParser() *Parser {
	return &Parser{
		StrictMode:      true,
		AllowEmptyScope: true,
	}
}

// conventionalCommitRegex matches: type(scope)!: description.
// Groups: 1=type, 2=scope with parens, 3=scope, 4=breaking indicator, 5=description.
var conventionalCommitRegex = regexp.MustCompile(`^(\w+)(\(([^)]*)\))?(!)?:\s*(.+)`)

// Ticket reference patterns.
var (
	// jiraTicketRegex matches JIRA tickets: PROJ-123, ABC-456 (3-4 letter prefixes).
	jiraTicketRegex = regexp.MustCompile(`\b([A-Z]{3,4}-\d+)\b`)

	// githubTicketRegex matches GitHub issues: #123, GH-456.
	githubTicketRegex = regexp.MustCompile(`(?:#(\d+)|GH-(\d+))\b`)

	// genericTicketRegex matches generic format: [TICKET-123] (3-4 letter prefixes).
	genericTicketRegex = regexp.MustCompile(`\[([A-Z]{3,4}-\d+)\]`)
)

// Parse parses a commit message into a Commit struct.
func (p *Parser) Parse(message string) (*Commit, error) {
	if message == "" {
		return nil, ErrEmptyMessage
	}

	lines := strings.Split(message, "\n")
	if len(lines) == 0 {
		return nil, ErrEmptyMessage
	}

	// Parse the header
	commit, err := p.parseHeader(lines[0], message)
	if err != nil {
		return nil, err
	}
	if commit.Type == "" {
		// Non-strict mode fallback already handled in parseHeader
		return commit, nil
	}

	// Parse body and footer if present
	if len(lines) > 1 {
		p.parseBodyAndFooter(commit, lines)
	}

	// Parse ticket references from entire commit message
	commit.TicketRefs = parseTicketRefs(message)

	return commit, nil
}

// parseHeader parses the commit header (first line) and returns a commit struct.
func (p *Parser) parseHeader(header, fullMessage string) (*Commit, error) {
	matches := conventionalCommitRegex.FindStringSubmatch(header)
	
	if matches == nil {
		if p.StrictMode {
			return nil, fmt.Errorf("%w: expected 'type(scope): description' format", ErrInvalidFormat)
		}
		// In non-strict mode, treat entire message as description.
		return &Commit{
			Description: header,
			Raw:         fullMessage,
		}, nil
	}

	return &Commit{
		Type:        matches[1],
		Scope:       matches[3],
		Breaking:    matches[4] == "!",
		Description: matches[5],
		Raw:         fullMessage,
	}, nil
}

// parseBodyAndFooter parses the body and footer sections of a commit message.
func (p *Parser) parseBodyAndFooter(commit *Commit, lines []string) {
	bodyStart := 1
	// Skip empty line after header if present.
	if bodyStart < len(lines) && lines[bodyStart] == "" {
		bodyStart++
	}

	// Find footer (starts with BREAKING CHANGE: or contains : ).
	footerStart := p.findFooterStart(lines, bodyStart)

	// Set body
	if footerStart > bodyStart {
		commit.Body = strings.TrimSpace(strings.Join(lines[bodyStart:footerStart], "\n"))
	} else if footerStart == -1 && bodyStart < len(lines) {
		commit.Body = strings.TrimSpace(strings.Join(lines[bodyStart:], "\n"))
	}

	// Set footer and check for breaking changes
	if footerStart != -1 {
		commit.Footer = strings.TrimSpace(strings.Join(lines[footerStart:], "\n"))
		if p.hasBreakingChangeInFooter(commit.Footer) {
			commit.Breaking = true
		}
	}
}

// findFooterStart finds the starting line index of the footer section.
func (p *Parser) findFooterStart(lines []string, bodyStart int) int {
	footerStart := -1
	for i := len(lines) - 1; i >= bodyStart; i-- {
		line := lines[i]
		if p.isBreakingChangeLine(line) || isFooterLine(line) {
			footerStart = i
		} else if line != "" && footerStart == -1 {
			// Non-footer line found, stop looking.
			break
		}
	}
	return footerStart
}

// isBreakingChangeLine checks if a line indicates a breaking change.
func (*Parser) isBreakingChangeLine(line string) bool {
	return strings.HasPrefix(line, "BREAKING CHANGE:") ||
		strings.HasPrefix(line, "BREAKING-CHANGE:")
}

// hasBreakingChangeInFooter checks if the footer contains breaking change indicators.
func (*Parser) hasBreakingChangeInFooter(footer string) bool {
	return strings.Contains(footer, "BREAKING CHANGE:") ||
		strings.Contains(footer, "BREAKING-CHANGE:")
}

// isFooterLine checks if a line looks like a footer token.
func isFooterLine(line string) bool {
	// Common footer tokens.
	footerTokens := []string{
		"Signed-off-by:",
		"Co-authored-by:",
		"Fixes:",
		"Closes:",
		"Refs:",
		"See-also:",
	}

	for _, token := range footerTokens {
		if strings.HasPrefix(line, token) {
			return true
		}
	}

	// Generic token format: Word-Word: or Word:.
	matched, err := regexp.MatchString(`^[A-Z][a-z]+(-[A-Z][a-z]+)*:\s+`, line)
	if err != nil {
		// If regex compilation fails, return false to be safe.
		return false
	}
	if matched {
		return true
	}

	return false
}

// parseTicketRefs extracts ticket references from a commit message.
func parseTicketRefs(message string) []TicketRef {
	var refs []TicketRef
	seen := make(map[string]bool)

	refs = parseGithubRefs(message, refs, seen)
	refs = parseGenericRefs(message, refs, seen)
	refs = parseJiraRefs(message, refs, seen)

	return refs
}

// parseGithubRefs extracts GitHub issue references.
func parseGithubRefs(message string, refs []TicketRef, seen map[string]bool) []TicketRef {
	matches := githubTicketRegex.FindAllStringSubmatch(message, -1)
	for _, match := range matches {
		if len(match) >= 3 {
			var id string
			if match[1] != "" { // #123 format.
				id = match[1]
			} else if match[2] != "" { // GH-456 format.
				id = match[2]
			}
			if id != "" {
				ref := TicketRef{
					Type: "GITHUB",
					ID:   id,
					Raw:  match[0],
				}
				refs = addUniqueRef(refs, ref, seen)
			}
		}
	}
	return refs
}

// parseGenericRefs extracts generic bracketed ticket references.
func parseGenericRefs(message string, refs []TicketRef, seen map[string]bool) []TicketRef {
	matches := genericTicketRegex.FindAllStringSubmatch(message, -1)
	for _, match := range matches {
		if len(match) > 1 {
			ref := TicketRef{
				Type: "GENERIC",
				ID:   match[1],
				Raw:  match[0],
			}
			refs = addUniqueRef(refs, ref, seen)
		}
	}
	return refs
}

// parseJiraRefs extracts JIRA ticket references.
func parseJiraRefs(message string, refs []TicketRef, seen map[string]bool) []TicketRef {
	matches := jiraTicketRegex.FindAllStringSubmatch(message, -1)
	for _, match := range matches {
		if len(match) > 1 {
			// Check if this was already classified as generic or github.
			if isAlreadyClassified(match[1], seen) {
				continue
			}

			// Skip GitHub-style references (GH-123 format).
			if strings.HasPrefix(match[1], "GH-") {
				continue
			}

			ref := TicketRef{
				Type: "JIRA",
				ID:   match[1],
				Raw:  match[0],
			}
			refs = addUniqueRef(refs, ref, seen)
		}
	}
	return refs
}

// addUniqueRef adds a ticket reference if it hasn't been seen before.
func addUniqueRef(refs []TicketRef, ref TicketRef, seen map[string]bool) []TicketRef {
	key := ref.Type + ":" + ref.ID
	if !seen[key] {
		refs = append(refs, ref)
		seen[key] = true
	}
	return refs
}

// isAlreadyClassified checks if a ticket ID was already classified.
func isAlreadyClassified(id string, seen map[string]bool) bool {
	genericKey := "GENERIC:" + id
	githubKey := "GITHUB:" + id
	return seen[genericKey] || seen[githubKey]
}

// HasTicketRefs returns true if the commit has any ticket references.
func (c *Commit) HasTicketRefs() bool {
	return len(c.TicketRefs) > 0
}

// HasJIRATicket returns true if the commit has JIRA ticket references.
func (c *Commit) HasJIRATicket() bool {
	for _, ref := range c.TicketRefs {
		if ref.Type == "JIRA" {
			return true
		}
	}
	return false
}

// GetJIRATickets returns all JIRA ticket references.
func (c *Commit) GetJIRATickets() []TicketRef {
	var jiraRefs []TicketRef
	for _, ref := range c.TicketRefs {
		if ref.Type == "JIRA" {
			jiraRefs = append(jiraRefs, ref)
		}
	}
	return jiraRefs
}

// Format formats a Commit back to conventional commit format.
func (c *Commit) Format() string {
	var sb strings.Builder

	// Write header.
	sb.WriteString(c.Type)
	if c.Scope != "" {
		sb.WriteString("(")
		sb.WriteString(c.Scope)
		sb.WriteString(")")
	}
	if c.Breaking {
		sb.WriteString("!")
	}
	sb.WriteString(": ")
	sb.WriteString(c.Description)

	// Write body if present.
	if c.Body != "" {
		sb.WriteString("\n\n")
		sb.WriteString(c.Body)
	}

	// Write footer if present.
	if c.Footer != "" {
		sb.WriteString("\n\n")
		sb.WriteString(c.Footer)
	}

	return sb.String()
}

// Header returns the first line of the commit message.
func (c *Commit) Header() string {
	var sb strings.Builder
	sb.WriteString(c.Type)
	if c.Scope != "" {
		sb.WriteString("(")
		sb.WriteString(c.Scope)
		sb.WriteString(")")
	}
	if c.Breaking {
		sb.WriteString("!")
	}
	sb.WriteString(": ")
	sb.WriteString(c.Description)
	return sb.String()
}
