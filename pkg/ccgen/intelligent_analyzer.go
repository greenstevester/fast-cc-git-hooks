// Package ccgen - Intelligent analysis inspired by Claude's commit message generation
package ccgen

import (
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
