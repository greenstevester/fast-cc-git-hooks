// Package ccgen - Intelligent analysis inspired by Claude's commit message generation
package ccgen

import (
	"fmt"
	"regexp"
	"strings"
)

// IntelligentChangeAnalysis provides advanced analysis of code changes
type IntelligentChangeAnalysis struct {
	FilePath    string
	ChangeType  string
	Scope       string
	Description string
	Details     []string
	Files       []string
	Priority    int
	Impact      string
	Context     string
}

// analyzeChangeIntelligently performs Claude-inspired intelligent analysis
func (g *Generator) analyzeChangeIntelligently(fileDiff string) *IntelligentChangeAnalysis {
	lines := strings.Split(fileDiff, "\n")
	if len(lines) < 2 {
		return nil
	}

	// Extract filename
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

	analysis := &IntelligentChangeAnalysis{
		FilePath: filename,
		Files:    []string{filename},
	}

	// Perform deep analysis of the change
	g.analyzeChangeSemantics(analysis, fileDiff)
	g.analyzeChangeImpact(analysis, fileDiff)
	g.analyzeChangeContext(analysis, fileDiff)

	return analysis
}

// analyzeChangeSemantics performs semantic analysis of code changes
func (g *Generator) analyzeChangeSemantics(analysis *IntelligentChangeAnalysis, diff string) {
	// Determine scope
	analysis.Scope = g.determineIntelligentScope(analysis.FilePath)

	// Analyze change patterns
	isNewFile := strings.Contains(diff, "new file mode")
	isDeletedFile := strings.Contains(diff, "deleted file mode")

	addedLines := g.extractAddedLines(diff)
	removedLines := g.extractRemovedLines(diff)

	// Intelligent type determination
	if isNewFile {
		analysis.ChangeType = "feat"
		analysis.Description = g.generateNewFileDescription(analysis.FilePath, addedLines)
	} else if isDeletedFile {
		analysis.ChangeType = "refactor"
		analysis.Description = g.generateDeletedFileDescription(analysis.FilePath)
	} else {
		analysis.ChangeType, analysis.Description = g.analyzeModificationSemantics(analysis.FilePath, addedLines, removedLines)
	}

	analysis.Priority = g.getTypePriority(analysis.ChangeType)
}

// analyzeChangeImpact determines the impact and details of changes
func (g *Generator) analyzeChangeImpact(analysis *IntelligentChangeAnalysis, diff string) {
	var details []string

	addedLines := g.extractAddedLines(diff)
	removedLines := g.extractRemovedLines(diff)

	// Analyze function changes
	if functionDetails := g.analyzeFunctionChanges(addedLines, removedLines); len(functionDetails) > 0 {
		details = append(details, functionDetails...)
	}

	// Analyze structural changes
	if structuralDetails := g.analyzeStructuralChanges(addedLines, removedLines); len(structuralDetails) > 0 {
		details = append(details, structuralDetails...)
	}

	// Analyze configuration changes
	if configDetails := g.analyzeConfigurationChanges(analysis.FilePath, addedLines, removedLines); len(configDetails) > 0 {
		details = append(details, configDetails...)
	}

	// Analyze documentation changes
	if docDetails := g.analyzeDocumentationChanges(analysis.FilePath, addedLines, removedLines); len(docDetails) > 0 {
		details = append(details, docDetails...)
	}

	analysis.Details = details

	// Determine overall impact
	if len(addedLines) > len(removedLines)*2 {
		analysis.Impact = "major additions"
	} else if len(removedLines) > len(addedLines)*2 {
		analysis.Impact = "significant cleanup"
	} else if len(details) > 3 {
		analysis.Impact = "multiple improvements"
	} else {
		analysis.Impact = "targeted changes"
	}
}

// analyzeChangeContext provides context about why changes were made
func (g *Generator) analyzeChangeContext(analysis *IntelligentChangeAnalysis, diff string) {
	contexts := []string{}

	// Look for error handling improvements
	if g.containsPattern(strings.Split(diff, "\n"), `\+.*(?i)(error|Error)`) {
		contexts = append(contexts, "improve error handling")
	}

	// Look for performance improvements
	if g.containsPattern(strings.Split(diff, "\n"), `\+.*(?i)(cache|Cache|optimize|Optimize)`) {
		contexts = append(contexts, "enhance performance")
	}

	// Look for security improvements
	if g.containsPattern(strings.Split(diff, "\n"), `\+.*(?i)(security|Security|validate|Validate)`) {
		contexts = append(contexts, "strengthen security")
	}

	// Look for user experience improvements
	if g.containsPattern(strings.Split(diff, "\n"), `\+.*(?i)(user|User|help|Help)`) {
		contexts = append(contexts, "improve user experience")
	}

	if len(contexts) > 0 {
		analysis.Context = strings.Join(contexts, " and ")
	}
}

// determineIntelligentScope provides more granular scope detection
func (g *Generator) determineIntelligentScope(filename string) string {
	switch {
	case strings.HasPrefix(filename, "cmd/"):
		// Extract specific command name
		parts := strings.Split(filename, "/")
		if len(parts) > 1 {
			return parts[1]
		}
		return "cli"
	case strings.HasPrefix(filename, "internal/"):
		parts := strings.Split(filename, "/")
		if len(parts) > 1 {
			return parts[1]
		}
		return "internal"
	case strings.HasPrefix(filename, "pkg/"):
		parts := strings.Split(filename, "/")
		if len(parts) > 1 {
			return parts[1]
		}
		return "pkg"
	case strings.HasPrefix(filename, ".github/workflows/"):
		return "ci"
	case strings.HasPrefix(filename, "docs/") || strings.HasSuffix(filename, ".md"):
		if strings.Contains(filename, "README") {
			return "docs"
		}
		return "docs"
	case strings.Contains(filename, "test") || strings.HasSuffix(filename, "_test.go"):
		return "test"
	case filename == "Makefile" || filename == "go.mod" || filename == "go.sum":
		return "build"
	default:
		return ""
	}
}

// extractAddedLines extracts lines that were added
func (g *Generator) extractAddedLines(diff string) []string {
	var addedLines []string
	for _, line := range strings.Split(diff, "\n") {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			addedLines = append(addedLines, strings.TrimPrefix(line, "+"))
		}
	}
	return addedLines
}

// extractRemovedLines extracts lines that were removed
func (g *Generator) extractRemovedLines(diff string) []string {
	var removedLines []string
	for _, line := range strings.Split(diff, "\n") {
		if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			removedLines = append(removedLines, strings.TrimPrefix(line, "-"))
		}
	}
	return removedLines
}

// generateNewFileDescription creates description for new files
func (g *Generator) generateNewFileDescription(filePath string, addedLines []string) string {
	filename := g.extractFileName(filePath)

	// Analyze what kind of file was added
	if strings.HasSuffix(filePath, "_test.go") {
		return fmt.Sprintf("add comprehensive tests for %s", strings.TrimSuffix(filename, "_test"))
	}

	if strings.HasSuffix(filePath, ".md") {
		return fmt.Sprintf("add %s documentation", strings.TrimSuffix(filename, ".md"))
	}

	// Look for main functions, interfaces, structs
	hasMainFunc := g.containsPattern(addedLines, `func main\(`)
	hasInterface := g.containsPattern(addedLines, `type \w+ interface`)
	hasStruct := g.containsPattern(addedLines, `type \w+ struct`)

	if hasMainFunc {
		return fmt.Sprintf("add %s command implementation", filename)
	} else if hasInterface {
		return fmt.Sprintf("add %s interface definition", filename)
	} else if hasStruct {
		return fmt.Sprintf("add %s data structures", filename)
	}

	return fmt.Sprintf("add %s implementation", filename)
}

// generateDeletedFileDescription creates description for deleted files
func (g *Generator) generateDeletedFileDescription(filePath string) string {
	filename := g.extractFileName(filePath)
	return fmt.Sprintf("remove %s implementation", filename)
}

// analyzeModificationSemantics analyzes modifications to existing files
func (g *Generator) analyzeModificationSemantics(filePath string, addedLines, removedLines []string) (string, string) {
	filename := g.extractFileName(filePath)

	// Test files
	if strings.HasSuffix(filePath, "_test.go") {
		if len(addedLines) > len(removedLines) {
			return "test", fmt.Sprintf("enhance %s test coverage", strings.TrimSuffix(filename, "_test"))
		}
		return "test", fmt.Sprintf("update %s test cases", strings.TrimSuffix(filename, "_test"))
	}

	// Documentation
	if strings.HasSuffix(filePath, ".md") {
		return "docs", fmt.Sprintf("update %s documentation", strings.TrimSuffix(filename, ".md"))
	}

	// Look for fix patterns
	if g.containsPattern(addedLines, `(?i)(fix|resolve|correct|repair)`) ||
		g.containsPattern(removedLines, `(?i)(bug|error|issue|problem)`) {
		return "fix", fmt.Sprintf("resolve %s issues", filename)
	}

	// Look for refactoring patterns
	if len(removedLines) > len(addedLines) ||
		g.containsPattern(addedLines, `(?i)(refactor|restructure|reorganize)`) {
		return "refactor", fmt.Sprintf("refactor %s for better maintainability", filename)
	}

	// Look for feature patterns
	if g.containsPattern(addedLines, `func \w+`) && !g.containsPattern(removedLines, `func \w+`) {
		return "feat", fmt.Sprintf("add new functionality to %s", filename)
	}

	// Default to enhancement
	return "feat", fmt.Sprintf("enhance %s functionality", filename)
}

// analyzeFunctionChanges analyzes function-level changes
func (g *Generator) analyzeFunctionChanges(addedLines, removedLines []string) []string {
	var details []string

	// Find new functions
	funcRegex := regexp.MustCompile(`func\s+(\w+)`)
	for _, line := range addedLines {
		if matches := funcRegex.FindStringSubmatch(line); len(matches) > 1 {
			details = append(details, fmt.Sprintf("Add %s function implementation", matches[1]))
		}
	}

	// Find removed functions
	for _, line := range removedLines {
		if matches := funcRegex.FindStringSubmatch(line); len(matches) > 1 {
			details = append(details, fmt.Sprintf("Remove %s function", matches[1]))
		}
	}

	return details
}

// analyzeStructuralChanges analyzes structural code changes
func (g *Generator) analyzeStructuralChanges(addedLines, removedLines []string) []string {
	var details []string

	// Look for new imports
	if g.containsPattern(addedLines, `import\s+`) {
		details = append(details, "Add new dependencies")
	}

	// Look for new types
	if g.containsPattern(addedLines, `type\s+\w+\s+struct`) {
		details = append(details, "Define new data structures")
	}

	// Look for interface changes
	if g.containsPattern(addedLines, `type\s+\w+\s+interface`) {
		details = append(details, "Define new interfaces")
	}

	return details
}

// analyzeConfigurationChanges analyzes configuration-related changes
func (g *Generator) analyzeConfigurationChanges(filePath string, addedLines, removedLines []string) []string {
	var details []string

	if strings.Contains(filePath, "config") || strings.HasSuffix(filePath, ".yaml") || strings.HasSuffix(filePath, ".yml") {
		if len(addedLines) > 0 {
			details = append(details, "Update configuration settings")
		}
	}

	if filePath == "Makefile" {
		details = append(details, "Update build configuration")
	}

	if filePath == "go.mod" || filePath == "go.sum" {
		details = append(details, "Update project dependencies")
	}

	return details
}

// analyzeDocumentationChanges analyzes documentation changes
func (g *Generator) analyzeDocumentationChanges(filePath string, addedLines, removedLines []string) []string {
	var details []string

	if strings.HasSuffix(filePath, ".md") {
		if len(addedLines) > len(removedLines) {
			details = append(details, "Expand documentation with new information")
		} else {
			details = append(details, "Update documentation for clarity")
		}
	}

	// Look for comment changes
	commentCount := 0
	for _, line := range addedLines {
		if strings.TrimSpace(line) != "" && strings.HasPrefix(strings.TrimSpace(line), "//") {
			commentCount++
		}
	}

	if commentCount > 2 {
		details = append(details, "Improve code documentation and comments")
	}

	return details
}

// containsPattern checks if any line matches a regex pattern
func (g *Generator) containsPattern(lines []string, pattern string) bool {
	regex := regexp.MustCompile(pattern)
	for _, line := range lines {
		if regex.MatchString(line) {
			return true
		}
	}
	return false
}

// extractFileName extracts the base filename from a path
func (g *Generator) extractFileName(filePath string) string {
	parts := strings.Split(filePath, "/")
	filename := parts[len(parts)-1]

	// Remove common extensions
	filename = strings.TrimSuffix(filename, ".go")
	filename = strings.TrimSuffix(filename, ".md")
	filename = strings.TrimSuffix(filename, ".yaml")
	filename = strings.TrimSuffix(filename, ".yml")

	return filename
}
