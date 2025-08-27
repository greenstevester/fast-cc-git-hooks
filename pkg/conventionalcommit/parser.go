// Package conventionalcommit provides parsing and validation for conventional commit messages.
package conventionalcommit

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	// ErrInvalidFormat indicates the commit message doesn't match conventional format
	ErrInvalidFormat = errors.New("invalid conventional commit format")
	// ErrEmptyMessage indicates the commit message is empty
	ErrEmptyMessage = errors.New("empty commit message")
)

// Commit represents a parsed conventional commit message
type Commit struct {
	Type        string
	Scope       string
	Breaking    bool
	Description string
	Body        string
	Footer      string
	Raw         string
}

// Parser provides conventional commit parsing with configurable options
type Parser struct {
	// StrictMode enforces strict conventional commit format
	StrictMode bool
	// AllowEmptyScope permits commits without scope
	AllowEmptyScope bool
}

// DefaultParser returns a parser with default settings
func DefaultParser() *Parser {
	return &Parser{
		StrictMode:      true,
		AllowEmptyScope: true,
	}
}

// conventionalCommitRegex matches: type(scope)!: description
// Groups: 1=type, 2=scope with parens, 3=scope, 4=breaking indicator, 5=description
var conventionalCommitRegex = regexp.MustCompile(`^(\w+)(\(([^)]*)\))?(!)?:\s*(.+)`)

// Parse parses a commit message into a Commit struct
func (p *Parser) Parse(message string) (*Commit, error) {
	if message == "" {
		return nil, ErrEmptyMessage
	}

	lines := strings.Split(message, "\n")
	if len(lines) == 0 {
		return nil, ErrEmptyMessage
	}

	// Parse the first line (header)
	header := lines[0]
	matches := conventionalCommitRegex.FindStringSubmatch(header)
	
	if matches == nil {
		if p.StrictMode {
			return nil, fmt.Errorf("%w: expected 'type(scope): description' format", ErrInvalidFormat)
		}
		// In non-strict mode, treat entire message as description
		return &Commit{
			Description: header,
			Raw:         message,
		}, nil
	}

	commit := &Commit{
		Type:        matches[1],
		Scope:       matches[3],
		Breaking:    matches[4] == "!",
		Description: matches[5],
		Raw:         message,
	}

	// Parse body and footer if present
	if len(lines) > 1 {
		bodyStart := 1
		// Skip empty line after header if present
		if bodyStart < len(lines) && lines[bodyStart] == "" {
			bodyStart++
		}

		// Find footer (starts with BREAKING CHANGE: or contains : )
		footerStart := -1
		for i := len(lines) - 1; i >= bodyStart; i-- {
			line := lines[i]
			if strings.HasPrefix(line, "BREAKING CHANGE:") || 
			   strings.HasPrefix(line, "BREAKING-CHANGE:") ||
			   isFooterLine(line) {
				footerStart = i
			} else if line != "" && footerStart == -1 {
				// Non-footer line found, stop looking
				break
			}
		}

		if footerStart > bodyStart {
			commit.Body = strings.TrimSpace(strings.Join(lines[bodyStart:footerStart], "\n"))
		} else if footerStart == -1 && bodyStart < len(lines) {
			commit.Body = strings.TrimSpace(strings.Join(lines[bodyStart:], "\n"))
		}

		if footerStart != -1 {
			commit.Footer = strings.TrimSpace(strings.Join(lines[footerStart:], "\n"))
			// Check for breaking change in footer
			if strings.Contains(commit.Footer, "BREAKING CHANGE:") ||
			   strings.Contains(commit.Footer, "BREAKING-CHANGE:") {
				commit.Breaking = true
			}
		}
	}

	return commit, nil
}

// isFooterLine checks if a line looks like a footer token
func isFooterLine(line string) bool {
	// Common footer tokens
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
	
	// Generic token format: Word-Word: or Word:
	if matched, _ := regexp.MatchString(`^[A-Z][a-z]+(-[A-Z][a-z]+)*:\s+`, line); matched {
		return true
	}
	
	return false
}

// Format formats a Commit back to conventional commit format
func (c *Commit) Format() string {
	var sb strings.Builder
	
	// Write header
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
	
	// Write body if present
	if c.Body != "" {
		sb.WriteString("\n\n")
		sb.WriteString(c.Body)
	}
	
	// Write footer if present
	if c.Footer != "" {
		sb.WriteString("\n\n")
		sb.WriteString(c.Footer)
	}
	
	return sb.String()
}

// Header returns the first line of the commit message
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