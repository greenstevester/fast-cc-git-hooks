// Package ccgen - analysis functions for commit message generation
package ccgen

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"
)

// analyzeDiff analyzes a git diff and returns detected changes
func (g *Generator) analyzeDiff(diff string) []ChangeType {
	fileChanges := make(map[string]*ChangeType)

	// Parse diff by files.
	files := strings.Split(diff, "diff --git")
	changes := make([]ChangeType, 0, len(files))
	for _, file := range files {
		if strings.TrimSpace(file) == "" {
			continue
		}

		change := g.analyzeFileChange(file)
		if change != nil {
			// Merge similar changes.
			key := change.Type + ":" + change.Scope
			if existing, ok := fileChanges[key]; ok {
				existing.Files = append(existing.Files, change.Files...)
				if len(change.Description) > len(existing.Description) {
					existing.Description = change.Description
				}
			} else {
				fileChanges[key] = change
			}
		}
	}

	// Convert to slice and sort by priority.
	for _, change := range fileChanges {
		changes = append(changes, *change)
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Priority < changes[j].Priority
	})

	return changes
}

// analyzeDiffIntelligently performs Claude-style intelligent analysis
func (g *Generator) analyzeDiffIntelligently(diff string) []*IntelligentChangeAnalysis {
	// Parse diff by files
	files := strings.Split(diff, "diff --git")
	analyses := make([]*IntelligentChangeAnalysis, 0, len(files))
	
	for _, file := range files {
		if strings.TrimSpace(file) == "" {
			continue
		}

		analysis := g.analyzeChangeIntelligently(file)
		if analysis != nil {
			analyses = append(analyses, analysis)
		}
	}

	// Sort by priority
	sort.Slice(analyses, func(i, j int) bool {
		return analyses[i].Priority < analyses[j].Priority
	})

	return analyses
}

// analyzeFileChange analyzes a single file's changes
func (g *Generator) analyzeFileChange(fileDiff string) *ChangeType {
	lines := strings.Split(fileDiff, "\n")
	if len(lines) < 2 {
		return nil
	}

	// Extract filename.
	var filename string
	for _, line := range lines {
		if strings.HasPrefix(line, " a/") && strings.Contains(line, " b/") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				filename = strings.TrimPrefix(parts[0], "a/")
				break
			}
		}
	}

	if filename == "" {
		return nil
	}

	// Determine change type and scope.
	changeType, scope := g.determineTypeAndScope(filename, fileDiff)
	description := g.generateDescription(filename, fileDiff, changeType)

	return &ChangeType{
		Type:        changeType,
		Scope:       scope,
		Description: description,
		Files:       []string{filename},
		Priority:    g.getTypePriority(changeType),
	}
}

// determineTypeAndScope determines the commit type and scope
func (g *Generator) determineTypeAndScope(filename, diff string) (changeType, scope string) {
	scope = g.determineScope(filename)
	changeType = g.determineType(filename, diff)
	return changeType, scope
}

// determineScope determines the scope from the filename.
func (g *Generator) determineScope(filename string) string {
	switch {
	case strings.HasPrefix(filename, "cmd/"):
		return "cli"
	case strings.HasPrefix(filename, "internal/"):
		parts := strings.Split(filename, "/")
		if len(parts) > 1 {
			return parts[1]
		}
		return ""
	case strings.HasPrefix(filename, "pkg/"):
		return "api"
	case strings.HasPrefix(filename, ".github/"):
		return "ci"
	case strings.HasPrefix(filename, "docs/") || strings.HasSuffix(filename, ".md"):
		return "docs"
	case strings.Contains(filename, "test") || strings.HasSuffix(filename, "_test.go"):
		return "test"
	case filename == "Makefile" || filename == "go.mod" || filename == "go.sum":
		return "build"
	default:
		return ""
	}
}

// determineType determines the commit type from the filename and diff.
func (g *Generator) determineType(filename, diff string) string {
	switch {
	case strings.Contains(diff, "new file mode"):
		return "feat"
	case strings.Contains(diff, "deleted file mode"):
		return "refactor"
	case strings.HasSuffix(filename, "_test.go"):
		return "test"
	case strings.HasSuffix(filename, ".md"):
		return "docs"
	case strings.Contains(filename, "github/workflows"):
		return "ci"
	case filename == "Makefile" || filename == "go.mod" || filename == "go.sum":
		return "build"
	case strings.Contains(diff, "+func ") && !strings.Contains(diff, "-func "):
		return "feat"
	case strings.Contains(diff, "fix") || strings.Contains(diff, "Fix"):
		return "fix"
	case g.countAdditions(diff) > g.countDeletions(diff):
		return "feat"
	default:
		return "refactor"
	}
}

// generateDescription generates a human-readable description
func (g *Generator) generateDescription(filename, diff, changeType string) string {
	base := strings.TrimSuffix(filename, ".go")
	base = strings.TrimSuffix(base, ".md")

	// Extract meaningful part of filename.
	parts := strings.Split(base, "/")
	name := parts[len(parts)-1]

	switch changeType {
	case "feat":
		if strings.Contains(diff, "new file mode") {
			return fmt.Sprintf("add %s", name)
		}
		return fmt.Sprintf("enhance %s functionality", name)
	case "fix":
		return fmt.Sprintf("resolve %s issues", name)
	case "docs":
		return fmt.Sprintf("update %s documentation", name)
	case "test":
		return fmt.Sprintf("improve %s tests", name)
	case "ci":
		return fmt.Sprintf("update %s workflow", name)
	case "build":
		return fmt.Sprintf("update %s configuration", name)
	case "refactor":
		if strings.Contains(diff, "deleted file mode") {
			return fmt.Sprintf("remove %s", name)
		}
		return fmt.Sprintf("restructure %s", name)
	default:
		return fmt.Sprintf("update %s", name)
	}
}

// countAdditions counts added lines in diff
func (g *Generator) countAdditions(diff string) int {
	count := 0
	for _, line := range strings.Split(diff, "\n") {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			count++
		}
	}
	return count
}

// countDeletions counts deleted lines in diff
func (g *Generator) countDeletions(diff string) int {
	count := 0
	for _, line := range strings.Split(diff, "\n") {
		if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			count++
		}
	}
	return count
}

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
