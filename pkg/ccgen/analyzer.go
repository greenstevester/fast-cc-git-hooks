// Package ccgen - analysis functions for commit message generation
package ccgen

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// getTypePriority returns priority for change type sorting
func (g *Generator) getTypePriority(changeType string) int {
	priorities := map[string]int{
		"feat":     1,
		"fix":      2,
		"perf":     3,
		"refactor": 4,
		"test":     5,
		"docs":     6,
		"ci":       7,
		"build":    8,
		"chore":    9,
	}

	if priority, ok := priorities[changeType]; ok {
		return priority
	}
	return 10
}

// GenerateCommitMessage generates the final commit message
func (g *Generator) GenerateCommitMessage(changes []ChangeType) string {
	if len(changes) == 0 {
		return "chore: update files"
	}

	// Use the highest priority change as primary.
	primary := changes[0]

	// Get JIRA ticket if available
	var jiraTicket string
	if g.options.JiraManager != nil {
		if ticket, err := g.options.JiraManager.GetCurrentJiraTicket(); err == nil && ticket != "" {
			jiraTicket = ticket
		}
	}

	// Create subject line.
	subject := primary.Type
	if primary.Scope != "" {
		subject += fmt.Sprintf("(%s)", primary.Scope)
	}
	subject += ": "

	// Add JIRA ticket if available
	if jiraTicket != "" {
		subject += fmt.Sprintf("%s ", jiraTicket)
	}

	subject += primary.Description

	// Truncate subject if too long.
	if utf8.RuneCountInString(subject) > MaxSubjectLength {
		runes := []rune(subject)
		if len(runes) > MaxSubjectLength-3 {
			subject = string(runes[:MaxSubjectLength-3]) + "..."
		}
	}

	// Generate body for multiple changes or complex single change.
	var body []string

	if len(changes) > 1 {
		body = append(body, "", "Changes include:")
		for _, change := range changes {
			line := fmt.Sprintf("- %s", g.capitalizeFirst(change.Description))
			if len(change.Files) > 0 {
				line += fmt.Sprintf(" (%s)", change.Files[0])
				if len(change.Files) > 1 {
					line += fmt.Sprintf(" and %d more", len(change.Files)-1)
				}
			}
			body = append(body, g.wrapLine(line, MaxBodyLineLength))
		}
	}

	if len(body) > 0 {
		return subject + strings.Join(body, "\n")
	}

	return subject
}

// capitalizeFirst capitalizes the first character of a string
func (g *Generator) capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	firstRune, size := utf8.DecodeRuneInString(s)
	upperFirst, _ := utf8.DecodeRuneInString(strings.ToUpper(string(firstRune)))
	return string(upperFirst) + s[size:]
}

// wrapLine wraps a line to the specified maximum length
func (g *Generator) wrapLine(line string, maxLength int) string {
	if utf8.RuneCountInString(line) <= maxLength {
		return line
	}

	words := strings.Fields(line)
	if len(words) == 0 {
		return line
	}

	var wrapped []string
	currentLine := words[0]

	for _, word := range words[1:] {
		testLine := currentLine + " " + word
		if utf8.RuneCountInString(testLine) <= maxLength {
			currentLine = testLine
		} else {
			wrapped = append(wrapped, currentLine)
			currentLine = word
		}
	}

	if currentLine != "" {
		wrapped = append(wrapped, currentLine)
	}

	return strings.Join(wrapped, "\n")
}
