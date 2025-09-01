// Package ccgen - Claude-inspired intelligent commit message generation
package ccgen

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"
)

// generateClaudeStyleCommitMessage generates a commit message using Claude's patterns
func (g *Generator) generateClaudeStyleCommitMessage(analyses []*IntelligentChangeAnalysis) string {
	return g.generateClaudeStyleCommitMessageWithPatterns(analyses, nil)
}

// generateClaudeStyleCommitMessageWithPatterns generates commit message using patterns from git analysis
func (g *Generator) generateClaudeStyleCommitMessageWithPatterns(analyses []*IntelligentChangeAnalysis, patterns *CommitPatterns) string {
	if len(analyses) == 0 {
		return "chore: update files"
	}

	// Sort by priority to get the primary change
	sort.Slice(analyses, func(i, j int) bool {
		return analyses[i].Priority < analyses[j].Priority
	})

	primary := analyses[0]

	// Adapt to existing repository style if patterns available
	if patterns != nil && patterns.PreferredStyle == "freeform" {
		// If repo uses freeform style, use simpler format
		return g.generateFreeformMessage(primary, analyses)
	}

	// Get JIRA ticket if available
	var jiraTicket string
	if g.options.JiraManager != nil {
		if ticket, err := g.options.JiraManager.GetCurrentJiraTicket(); err == nil && ticket != "" {
			jiraTicket = ticket
		}
	}

	// Create Claude-style subject line
	subject := g.buildClaudeSubject(primary, jiraTicket)

	// Adjust length based on repository patterns
	if patterns != nil && patterns.AverageLength > 0 {
		targetLength := patterns.AverageLength
		if len(subject) > targetLength && targetLength > 30 {
			subject = g.intelligentTruncate(subject, targetLength)
		}
	}

	// Create Claude-style body with detailed explanations
	body := g.buildClaudeBody(analyses, primary)

	if body != "" {
		return subject + "\n" + body
	}

	return subject
}

// generateFreeformMessage generates simpler messages for repos that don't use conventional commits
func (g *Generator) generateFreeformMessage(primary *IntelligentChangeAnalysis, analyses []*IntelligentChangeAnalysis) string {
	// Simple format without conventional commit structure
	if len(analyses) == 1 {
		return g.capitalizeFirst(primary.Description)
	}

	// Multiple changes - create simple list
	message := g.capitalizeFirst(primary.Description)
	if len(analyses) > 1 {
		message += fmt.Sprintf(" and %d other changes", len(analyses)-1)
	}

	return message
}

// buildClaudeSubject creates a Claude-style subject line
func (g *Generator) buildClaudeSubject(primary *IntelligentChangeAnalysis, jiraTicket string) string {
	subject := primary.ChangeType

	if primary.Scope != "" {
		subject += fmt.Sprintf("(%s)", primary.Scope)
	}
	subject += ": "

	// Add JIRA ticket if available
	if jiraTicket != "" {
		subject += fmt.Sprintf("%s ", jiraTicket)
	}

	// Use more descriptive, action-oriented language
	description := g.enhanceDescription(primary)
	subject += description

	// Truncate if too long, but prefer complete words
	if utf8.RuneCountInString(subject) > MaxSubjectLength {
		subject = g.intelligentTruncate(subject, MaxSubjectLength)
	}

	return subject
}

// enhanceDescription improves description with Claude-style language
func (g *Generator) enhanceDescription(analysis *IntelligentChangeAnalysis) string {
	desc := analysis.Description

	// Make descriptions more action-oriented and specific
	replacements := map[string]string{
		"update":     "improve",
		"enhance":    "improve",
		"change":     "update",
		"modify":     "refine",
		"fix issues": "resolve issues",
	}

	for old, new := range replacements {
		if strings.Contains(strings.ToLower(desc), old) {
			// Preserve case by using case-insensitive replacement
			desc = strings.ReplaceAll(desc, old, new)
			desc = strings.ReplaceAll(desc, g.capitalizeFirst(old), g.capitalizeFirst(new))
			break
		}
	}

	// Add context if available
	if analysis.Context != "" {
		// Integrate context naturally
		contextualWords := map[string]string{
			"improve error handling":  "with better error handling",
			"enhance performance":     "for better performance",
			"strengthen security":     "with enhanced security",
			"improve user experience": "for better usability",
		}

		if addition, exists := contextualWords[analysis.Context]; exists {
			desc += " " + addition
		}
	}

	return desc
}

// buildClaudeBody creates detailed body following Claude's patterns
func (g *Generator) buildClaudeBody(analyses []*IntelligentChangeAnalysis, primary *IntelligentChangeAnalysis) string {
	var bodyLines []string

	// Only create detailed body for complex changes
	shouldCreateBody := len(analyses) > 1 || len(primary.Details) > 2 || primary.Impact == "major additions"

	if !shouldCreateBody {
		return ""
	}

	// Add empty line before body
	bodyLines = append(bodyLines, "")

	if len(analyses) == 1 {
		// Single file with multiple changes
		if len(primary.Details) > 0 {
			for _, detail := range primary.Details {
				bodyLines = append(bodyLines, fmt.Sprintf("- %s", g.capitalizeFirst(detail)))
			}
		}

		// Add impact explanation for significant changes
		if primary.Impact != "" && primary.Impact != "targeted changes" {
			bodyLines = append(bodyLines, "")
			bodyLines = append(bodyLines, g.generateImpactExplanation(primary))
		}
	} else {
		// Multiple files - group by change type
		grouped := g.groupAnalysesByType(analyses)

		for changeType, group := range grouped {
			if len(group) == 1 {
				analysis := group[0]
				desc := g.enhanceDescription(analysis)
				bodyLines = append(bodyLines, fmt.Sprintf("- %s", g.capitalizeFirst(desc)))
			} else {
				// Multiple files of same type
				files := make([]string, 0, len(group))
				for _, analysis := range group {
					files = append(files, g.extractFileName(analysis.FilePath))
				}

				desc := g.generateGroupDescription(changeType, files)
				bodyLines = append(bodyLines, fmt.Sprintf("- %s", desc))
			}
		}

		// Add overall impact if significant
		if g.hasSignificantImpact(analyses) {
			bodyLines = append(bodyLines, "")
			bodyLines = append(bodyLines, g.generateOverallImpact(analyses))
		}
	}

	// Wrap lines to MaxBodyLineLength
	var wrappedLines []string
	for _, line := range bodyLines {
		if line == "" {
			wrappedLines = append(wrappedLines, line)
		} else {
			wrappedLines = append(wrappedLines, g.wrapLine(line, MaxBodyLineLength))
		}
	}

	return strings.Join(wrappedLines, "\n")
}

// groupAnalysesByType groups analyses by their change type
func (g *Generator) groupAnalysesByType(analyses []*IntelligentChangeAnalysis) map[string][]*IntelligentChangeAnalysis {
	grouped := make(map[string][]*IntelligentChangeAnalysis)

	for _, analysis := range analyses {
		key := analysis.ChangeType
		if analysis.Scope != "" {
			key += ":" + analysis.Scope
		}
		grouped[key] = append(grouped[key], analysis)
	}

	return grouped
}

// generateGroupDescription creates description for grouped changes
func (g *Generator) generateGroupDescription(changeType string, files []string) string {
	if len(files) == 1 {
		return fmt.Sprintf("%s %s functionality", g.getActionVerb(changeType), files[0])
	}

	if len(files) <= 3 {
		fileList := strings.Join(files[:len(files)-1], ", ")
		return fmt.Sprintf("%s %s and %s functionality", g.getActionVerb(changeType), fileList, files[len(files)-1])
	}

	return fmt.Sprintf("%s %s and %d other components", g.getActionVerb(changeType), files[0], len(files)-1)
}

// getActionVerb returns appropriate action verb for change type
func (g *Generator) getActionVerb(changeType string) string {
	verbs := map[string]string{
		"feat":     "Enhance",
		"fix":      "Resolve issues in",
		"refactor": "Refactor",
		"docs":     "Update documentation for",
		"test":     "Improve test coverage for",
		"ci":       "Update CI configuration for",
		"build":    "Update build configuration for",
		"chore":    "Update",
	}

	if verb, exists := verbs[changeType]; exists {
		return verb
	}
	return "Update"
}

// generateImpactExplanation creates explanation of change impact
func (g *Generator) generateImpactExplanation(analysis *IntelligentChangeAnalysis) string {
	switch analysis.Impact {
	case "major additions":
		return "These additions significantly expand the functionality and capabilities."
	case "significant cleanup":
		return "This cleanup improves code maintainability and reduces technical debt."
	case "multiple improvements":
		return "These changes collectively improve the overall system quality."
	default:
		return ""
	}
}

// hasSignificantImpact determines if changes have significant overall impact
func (g *Generator) hasSignificantImpact(analyses []*IntelligentChangeAnalysis) bool {
	totalFiles := len(analyses)
	majorChanges := 0

	for _, analysis := range analyses {
		if analysis.Impact == "major additions" || analysis.Impact == "significant cleanup" {
			majorChanges++
		}
	}

	return totalFiles >= 3 || majorChanges >= 2
}

// generateOverallImpact creates overall impact statement
func (g *Generator) generateOverallImpact(analyses []*IntelligentChangeAnalysis) string {
	hasFeatures := false
	hasFixes := false
	hasRefactoring := false

	for _, analysis := range analyses {
		switch analysis.ChangeType {
		case "feat":
			hasFeatures = true
		case "fix":
			hasFixes = true
		case "refactor":
			hasRefactoring = true
		}
	}

	var impacts []string
	if hasFeatures {
		impacts = append(impacts, "expand functionality")
	}
	if hasFixes {
		impacts = append(impacts, "improve reliability")
	}
	if hasRefactoring {
		impacts = append(impacts, "enhance maintainability")
	}

	if len(impacts) > 1 {
		last := impacts[len(impacts)-1]
		others := strings.Join(impacts[:len(impacts)-1], ", ")
		return fmt.Sprintf("These changes collectively %s and %s of the system.", others, last)
	} else if len(impacts) == 1 {
		return fmt.Sprintf("These changes %s of the system.", impacts[0])
	}

	return "These changes improve the overall quality and functionality of the system."
}

// intelligentTruncate truncates text at word boundaries when possible
func (g *Generator) intelligentTruncate(text string, maxLength int) string {
	if utf8.RuneCountInString(text) <= maxLength {
		return text
	}

	// Try to truncate at word boundary
	words := strings.Fields(text)
	result := ""

	for _, word := range words {
		test := result + word + " "
		if utf8.RuneCountInString(test) > maxLength-3 { // Leave room for "..."
			break
		}
		result = test
	}

	result = strings.TrimSpace(result) + "..."

	// If we couldn't fit even one word, do character truncation
	if result == "..." {
		runes := []rune(text)
		if len(runes) > maxLength-3 {
			result = string(runes[:maxLength-3]) + "..."
		}
	}

	return result
}
