// Package semantic provides semantic analysis integration for the cc command
package semantic

import (
	"context"
	"fmt"
	"strings"
)

// Enhanced cc command integration
type CCSemanticAnalyzer struct {
	analyzer *SemanticAnalyzer
	enabled  bool
}

// NewCCSemanticAnalyzer creates a new semantic analyzer for the cc command
func NewCCSemanticAnalyzer() *CCSemanticAnalyzer {
	registry := NewPluginRegistry()
	analyzer := NewSemanticAnalyzer(registry)

	return &CCSemanticAnalyzer{
		analyzer: analyzer,
		enabled:  true, // Can be configured
	}
}

// RegisterPlugins registers available semantic analysis plugins
func (c *CCSemanticAnalyzer) RegisterPlugins(plugins ...SemanticPlugin) error {
	for _, plugin := range plugins {
		if err := c.analyzer.registry.Register(plugin); err != nil {
			return fmt.Errorf("failed to register plugin %s: %w", plugin.Name(), err)
		}
	}
	return nil
}

// AnalyzeDiff analyzes a git diff for semantic changes
func (c *CCSemanticAnalyzer) AnalyzeDiff(diff string) (*SemanticChange, error) {
	if !c.enabled {
		return nil, nil
	}

	files := c.parseDiffToFileChanges(diff)
	if len(files) == 0 {
		return nil, nil
	}

	ctx := context.Background()
	changes, err := c.analyzer.AnalyzeChanges(ctx, files)
	if err != nil {
		return nil, fmt.Errorf("semantic analysis failed: %w", err)
	}

	if len(changes) == 0 {
		return nil, nil
	}

	// Return the highest confidence change
	return c.selectPrimaryChange(changes), nil
}

// parseDiffToFileChanges converts a git diff string to FileChange objects
func (c *CCSemanticAnalyzer) parseDiffToFileChanges(diff string) []FileChange {
	var files []FileChange

	// Split diff by files
	fileSections := strings.Split(diff, "diff --git")

	for _, section := range fileSections {
		if strings.TrimSpace(section) == "" {
			continue
		}

		file := c.parseFileSection(section)
		if file != nil {
			files = append(files, *file)
		}
	}

	return files
}

// parseFileSection parses a single file's diff section
func (c *CCSemanticAnalyzer) parseFileSection(section string) *FileChange {
	lines := strings.Split(section, "\n")
	if len(lines) < 2 {
		return nil
	}

	// Extract file path from the first line
	// Format: " a/path/to/file b/path/to/file"
	var filePath string
	if len(lines) > 0 {
		firstLine := lines[0]
		parts := strings.Fields(firstLine)
		if len(parts) >= 2 {
			// Remove 'a/' prefix
			filePath = strings.TrimPrefix(parts[0], "a/")
		}
	}

	if filePath == "" {
		return nil
	}

	// Determine change type
	changeType := "modified"
	if strings.Contains(section, "new file mode") {
		changeType = "added"
	} else if strings.Contains(section, "deleted file mode") {
		changeType = "deleted"
	}

	// Extract content changes
	var beforeContent, afterContent strings.Builder
	diffContent := section

	// Parse diff lines
	for _, line := range lines {
		if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			beforeContent.WriteString(strings.TrimPrefix(line, "-") + "\n")
		} else if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			afterContent.WriteString(strings.TrimPrefix(line, "+") + "\n")
		}
	}

	return &FileChange{
		Path:          filePath,
		Language:      c.detectLanguageFromPath(filePath),
		BeforeContent: beforeContent.String(),
		AfterContent:  afterContent.String(),
		DiffContent:   diffContent,
		ChangeType:    changeType,
	}
}

// detectLanguageFromPath detects language from file path
func (c *CCSemanticAnalyzer) detectLanguageFromPath(path string) string {
	switch {
	case strings.HasSuffix(path, ".tf") || strings.HasSuffix(path, ".tfvars"):
		return "terraform"
	case strings.HasSuffix(path, ".go"):
		return "go"
	case strings.HasSuffix(path, ".ts") || strings.HasSuffix(path, ".tsx"):
		return "typescript"
	case strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".jsx"):
		return "javascript"
	case strings.HasSuffix(path, ".py"):
		return "python"
	case strings.HasSuffix(path, ".java"):
		return "java"
	case strings.HasSuffix(path, ".rs"):
		return "rust"
	case strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml"):
		return "yaml"
	case strings.HasSuffix(path, ".json"):
		return "json"
	default:
		return "text"
	}
}

// selectPrimaryChange selects the most relevant change from multiple detected changes
func (c *CCSemanticAnalyzer) selectPrimaryChange(changes []*SemanticChange) *SemanticChange {
	if len(changes) == 0 {
		return nil
	}

	if len(changes) == 1 {
		return changes[0]
	}

	// Priority: breaking changes > feat > fix > perf > refactor > others
	priority := map[string]int{
		"feat":     1,
		"fix":      2,
		"perf":     3,
		"refactor": 4,
		"docs":     5,
		"test":     6,
		"build":    7,
		"ci":       8,
		"chore":    9,
	}

	var best *SemanticChange
	bestScore := 999

	for _, change := range changes {
		score := priority[change.Type]
		if score == 0 {
			score = 10 // unknown types get lowest priority
		}

		// Breaking changes get higher priority
		if change.BreakingChange {
			score -= 5
		}

		// Higher confidence gets priority
		confidenceBonus := int((1.0 - change.Confidence) * 3)
		score += confidenceBonus

		if score < bestScore {
			bestScore = score
			best = change
		}
	}

	return best
}

// Enable enables semantic analysis
func (c *CCSemanticAnalyzer) Enable() {
	c.enabled = true
}

// Disable disables semantic analysis
func (c *CCSemanticAnalyzer) Disable() {
	c.enabled = false
}

// IsEnabled returns whether semantic analysis is enabled
func (c *CCSemanticAnalyzer) IsEnabled() bool {
	return c.enabled
}

// GetAvailablePlugins returns list of registered plugins
func (c *CCSemanticAnalyzer) GetAvailablePlugins() []SemanticPlugin {
	return c.analyzer.registry.ListPlugins()
}
