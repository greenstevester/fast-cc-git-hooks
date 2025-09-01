// Package ccgen - Advanced Git analysis implementing comprehensive algorithm
package ccgen

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// GitAnalysisResult contains comprehensive git analysis data
type GitAnalysisResult struct {
	// File statistics and changes
	FileStats      map[string]*FileStatistics
	ChangeTypes    map[string]string // filename -> A/M/D
	TotalFiles     int
	TotalAdditions int
	TotalDeletions int

	// Directory-level statistics
	DirStats map[string]float64

	// Numerical statistics (precise counts)
	NumStats map[string]*NumStat

	// File operation summaries
	FileSummaries []string

	// Function context information
	ModifiedFunctions []string

	// Content analysis
	WordDiffContent string
	StagedDiff      string

	// Historical context
	RecentCommits  []CommitInfo
	CommitPatterns *CommitPatterns
}

// FileStatistics contains detailed stats for each file
type FileStatistics struct {
	Filename   string
	Additions  int
	Deletions  int
	ChangeType string // A/M/D
}

// NumStat contains precise numerical statistics from git diff --numstat
type NumStat struct {
	Additions int
	Deletions int
	Filename  string
}

// CommitInfo represents recent commit information
type CommitInfo struct {
	Hash    string
	Message string
	Files   int
}

// CommitPatterns analyzes patterns from recent commits
type CommitPatterns struct {
	CommonTypes    map[string]int
	CommonScopes   map[string]int
	AverageLength  int
	PreferredStyle string
}

// performAdvancedGitAnalysis implements the comprehensive algorithm
func (g *Generator) performAdvancedGitAnalysis() (*GitAnalysisResult, error) {
	result := &GitAnalysisResult{
		FileStats:         make(map[string]*FileStatistics),
		ChangeTypes:       make(map[string]string),
		DirStats:          make(map[string]float64),
		NumStats:          make(map[string]*NumStat),
		FileSummaries:     make([]string, 0),
		ModifiedFunctions: make([]string, 0),
	}

	// Step 1: Get change types (A/M/D) - fundamental file operations
	if err := g.getChangeTypes(result); err != nil {
		return nil, fmt.Errorf("getting change types: %w", err)
	}

	// Step 2: Get file operation summaries (create/delete/rename details)
	if err := g.getFileSummaries(result); err != nil {
		return nil, fmt.Errorf("getting file summaries: %w", err)
	}

	// Step 3: Get precise numerical statistics (exact line counts)
	if err := g.getNumStats(result); err != nil {
		return nil, fmt.Errorf("getting numerical statistics: %w", err)
	}

	// Step 4: Get file statistics (visual representation for compatibility)
	if err := g.getFileStatistics(result); err != nil {
		return nil, fmt.Errorf("getting file statistics: %w", err)
	}

	// Step 5: Get directory distribution statistics
	if err := g.getDirStats(result); err != nil {
		return nil, fmt.Errorf("getting directory statistics: %w", err)
	}

	// Step 6: Get staged diff (maintain compatibility)
	if err := g.getStagedDiffContent(result); err != nil {
		return nil, fmt.Errorf("getting staged diff: %w", err)
	}

	// Step 7: Get word-level diff for granular analysis
	if err := g.getWordDiff(result); err != nil {
		return nil, fmt.Errorf("getting word diff: %w", err)
	}

	// Step 8: Extract modified function contexts (specific change locations)
	if err := g.extractFunctionContexts(result); err != nil {
		return nil, fmt.Errorf("extracting function contexts: %w", err)
	}

	// Step 9: Analyze recent commit patterns
	g.analyzeRecentCommitPatterns(result)

	return result, nil
}

// getFileStatistics implements: git diff --stat HEAD~1 HEAD (or --staged if no HEAD~1)
func (g *Generator) getFileStatistics(result *GitAnalysisResult) error {
	fmt.Printf("Running `git diff --stat`")

	// Try staged first (for initial commits), fallback to HEAD~1 comparison
	var cmd *exec.Cmd
	if g.hasPreviousCommits() {
		cmd = exec.Command("git", "diff", "--stat", "HEAD~1", "HEAD")
	} else {
		cmd = exec.Command("git", "diff", "--stat", "--staged")
	}

	output, err := cmd.Output()
	if err != nil {
		// Fallback to staged if HEAD~1 fails
		cmd = exec.Command("git", "diff", "--stat", "--staged")
		output, err = cmd.Output()
		if err != nil {
			fmt.Println(" ❌")
			return fmt.Errorf("failed to get diff stat: %w", err)
		}
	}
	fmt.Println(" ✅")

	// Parse diff --stat output
	g.parseStatOutput(string(output), result)

	// Cross-reference with NumStats for more accurate data
	g.enhanceWithNumStats(result)

	return nil
}

// getChangeTypes implements: git diff --name-status HEAD~1 HEAD
func (g *Generator) getChangeTypes(result *GitAnalysisResult) error {
	fmt.Printf("Running `git diff --name-status`")

	var cmd *exec.Cmd
	if g.hasPreviousCommits() {
		cmd = exec.Command("git", "diff", "--name-status", "HEAD~1", "HEAD")
	} else {
		cmd = exec.Command("git", "diff", "--name-status", "--staged")
	}

	output, err := cmd.Output()
	if err != nil {
		// Fallback to staged
		cmd = exec.Command("git", "diff", "--name-status", "--staged")
		output, err = cmd.Output()
		if err != nil {
			fmt.Println(" ❌")
			return fmt.Errorf("failed to get name-status: %w", err)
		}
	}
	fmt.Println(" ✅")

	// Parse name-status output (format: "M\tfilename" or "A\tfilename")
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) >= 2 {
			changeType := parts[0]
			filename := parts[1]
			result.ChangeTypes[filename] = changeType

			// Update FileStatistics with change type
			if stat, exists := result.FileStats[filename]; exists {
				stat.ChangeType = changeType
			} else {
				result.FileStats[filename] = &FileStatistics{
					Filename:   filename,
					ChangeType: changeType,
				}
			}
		}
	}

	return nil
}

// getWordDiff implements: git diff HEAD~1 HEAD --word-diff
func (g *Generator) getWordDiff(result *GitAnalysisResult) error {
	fmt.Printf("Running `git diff --word-diff`")

	var cmd *exec.Cmd
	if g.hasPreviousCommits() {
		cmd = exec.Command("git", "diff", "HEAD~1", "HEAD", "--word-diff")
	} else {
		cmd = exec.Command("git", "diff", "--staged", "--word-diff")
	}

	output, err := cmd.Output()
	if err != nil {
		// Fallback to staged
		cmd = exec.Command("git", "diff", "--staged", "--word-diff")
		output, err = cmd.Output()
		if err != nil {
			fmt.Println(" ❌")
			return fmt.Errorf("failed to get word diff: %w", err)
		}
	}
	fmt.Println(" ✅")

	result.WordDiffContent = string(output)
	return nil
}

// getStagedDiffContent maintains compatibility with existing analyzer
func (g *Generator) getStagedDiffContent(result *GitAnalysisResult) error {
	fmt.Printf("Running `git diff --staged`")

	cmd := exec.Command("git", "diff", "--staged")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(" ❌")
		return fmt.Errorf("failed to get staged diff: %w", err)
	}
	fmt.Println(" ✅")

	result.StagedDiff = string(output)
	return nil
}

// analyzeRecentCommitPatterns implements: git log --oneline -10
func (g *Generator) analyzeRecentCommitPatterns(result *GitAnalysisResult) {
	fmt.Printf("Running `git log --oneline -10`")

	cmd := exec.Command("git", "log", "--oneline", "-10")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(" ❌")
		// Don't fail if no commits exist yet
		result.CommitPatterns = &CommitPatterns{
			CommonTypes:  make(map[string]int),
			CommonScopes: make(map[string]int),
		}
		return
	}
	fmt.Println(" ✅")

	// Parse recent commits
	result.RecentCommits = g.parseRecentCommits(string(output))
	result.CommitPatterns = g.analyzeCommitPatterns(result.RecentCommits)
}

// parseStatOutput parses git diff --stat output
func (g *Generator) parseStatOutput(output string, result *GitAnalysisResult) {
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		if strings.Contains(line, "|") {
			// Parse line: " filename.go | 23 +++++++++++----------"
			parts := strings.Split(line, "|")
			if len(parts) < 2 {
				continue
			}

			filename := strings.TrimSpace(parts[0])
			statsStr := strings.TrimSpace(parts[1])

			// Extract numbers and symbols
			additions, deletions := g.parseStatsLine(statsStr)

			result.FileStats[filename] = &FileStatistics{
				Filename:  filename,
				Additions: additions,
				Deletions: deletions,
			}

			result.TotalAdditions += additions
			result.TotalDeletions += deletions
			result.TotalFiles++
		}
	}
}

// parseStatsLine extracts addition/deletion counts from stats line
func (g *Generator) parseStatsLine(statsStr string) (additions, deletions int) {
	// Extract number at beginning (total changes)
	parts := strings.Fields(statsStr)
	if len(parts) == 0 {
		return 0, 0
	}

	totalChanges, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0 // Return zeros if parsing fails
	}

	// Count + and - symbols
	plusCount := strings.Count(statsStr, "+")
	minusCount := strings.Count(statsStr, "-")

	if plusCount+minusCount > 0 {
		// Proportional distribution based on symbols
		additions = (totalChanges * plusCount) / (plusCount + minusCount)
		deletions = totalChanges - additions
	}

	return additions, deletions
}

// parseRecentCommits parses git log --oneline output
func (g *Generator) parseRecentCommits(output string) []CommitInfo {
	var commits []CommitInfo
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) >= 2 {
			commits = append(commits, CommitInfo{
				Hash:    parts[0],
				Message: parts[1],
			})
		}
	}

	return commits
}

// analyzeCommitPatterns analyzes patterns from recent commits
func (g *Generator) analyzeCommitPatterns(commits []CommitInfo) *CommitPatterns {
	patterns := &CommitPatterns{
		CommonTypes:  make(map[string]int),
		CommonScopes: make(map[string]int),
	}

	totalLength := 0
	conventionalCount := 0

	// Regex for conventional commit format
	conventionalRegex := regexp.MustCompile(`^(\w+)(\([^)]+\))?: .+`)
	typeRegex := regexp.MustCompile(`^(\w+)`)
	scopeRegex := regexp.MustCompile(`^(\w+)\(([^)]+)\)`)

	for _, commit := range commits {
		message := commit.Message
		totalLength += len(message)

		// Check if it's conventional commit format
		if conventionalRegex.MatchString(message) {
			conventionalCount++

			// Extract type
			if typeMatches := typeRegex.FindStringSubmatch(message); len(typeMatches) > 1 {
				commitType := typeMatches[1]
				patterns.CommonTypes[commitType]++
			}

			// Extract scope if present
			if scopeMatches := scopeRegex.FindStringSubmatch(message); len(scopeMatches) > 2 {
				scope := scopeMatches[2]
				patterns.CommonScopes[scope]++
			}
		}
	}

	if len(commits) > 0 {
		patterns.AverageLength = totalLength / len(commits)

		// Determine preferred style
		if conventionalCount >= len(commits)/2 {
			patterns.PreferredStyle = "conventional"
		} else {
			patterns.PreferredStyle = "freeform"
		}
	}

	return patterns
}

// hasPreviousCommits checks if repository has any commits
func (g *Generator) hasPreviousCommits() bool {
	cmd := exec.Command("git", "rev-parse", "HEAD~1")
	return cmd.Run() == nil
}

// getAdvancedChangeAnalyses converts GitAnalysisResult to IntelligentChangeAnalysis
func (g *Generator) getAdvancedChangeAnalyses(analysis *GitAnalysisResult) []*IntelligentChangeAnalysis {
	var analyses []*IntelligentChangeAnalysis

	for filename, stats := range analysis.FileStats {
		changeAnalysis := g.createAdvancedChangeAnalysis(filename, stats, analysis)
		if changeAnalysis != nil {
			analyses = append(analyses, changeAnalysis)
		}
	}

	return analyses
}

// createAdvancedChangeAnalysis creates detailed analysis using comprehensive data
func (g *Generator) createAdvancedChangeAnalysis(filename string, stats *FileStatistics, gitAnalysis *GitAnalysisResult) *IntelligentChangeAnalysis {
	analysis := &IntelligentChangeAnalysis{
		FilePath: filename,
		Files:    []string{filename},
	}

	// Enhanced scope detection
	analysis.Scope = g.determineIntelligentScope(filename)

	// Advanced change type detection using change type + statistics
	analysis.ChangeType = g.determineAdvancedChangeType(stats, gitAnalysis)

	// Statistical impact assessment
	analysis.Impact = g.assessStatisticalImpact(stats, gitAnalysis)

	// Enhanced description using all available data
	analysis.Description = g.generateAdvancedDescription(filename, stats)

	// Context detection from word diff
	analysis.Context = g.detectContextFromWordDiff(gitAnalysis.WordDiffContent)

	// Priority based on change magnitude and type
	analysis.Priority = g.calculateAdvancedPriority(analysis.ChangeType, stats)

	return analysis
}

// determineAdvancedChangeType uses comprehensive data for better type detection
func (g *Generator) determineAdvancedChangeType(stats *FileStatistics, gitAnalysis *GitAnalysisResult) string {
	switch stats.ChangeType {
	case "A":
		return "feat"
	case "D":
		return "refactor"
	case "M":
		// For modifications, use ratio analysis
		total := stats.Additions + stats.Deletions
		if total == 0 {
			return "chore"
		}

		additionRatio := float64(stats.Additions) / float64(total)

		// Check patterns in word diff for more context
		if strings.Contains(gitAnalysis.WordDiffContent, "fix") || strings.Contains(gitAnalysis.WordDiffContent, "bug") {
			return "fix"
		}

		if strings.Contains(gitAnalysis.WordDiffContent, "test") || strings.HasSuffix(stats.Filename, "_test.go") {
			return "test"
		}

		if strings.HasSuffix(stats.Filename, ".md") {
			return "docs"
		}

		// Use addition ratio for feat vs refactor
		if additionRatio > 0.7 {
			return "feat"
		} else if additionRatio < 0.3 {
			return "refactor"
		} else {
			return "refactor"
		}
	default:
		return "chore"
	}
}

// assessStatisticalImpact uses statistical data for impact assessment
func (g *Generator) assessStatisticalImpact(stats *FileStatistics, gitAnalysis *GitAnalysisResult) string {
	total := stats.Additions + stats.Deletions
	avgChangesPerFile := 0
	if gitAnalysis.TotalFiles > 0 {
		avgChangesPerFile = (gitAnalysis.TotalAdditions + gitAnalysis.TotalDeletions) / gitAnalysis.TotalFiles
	}

	if total > avgChangesPerFile*2 {
		return "major changes"
	} else if total > avgChangesPerFile {
		return "moderate changes"
	} else {
		return "minor changes"
	}
}

// generateAdvancedDescription creates descriptions using comprehensive analysis
func (g *Generator) generateAdvancedDescription(filename string, stats *FileStatistics) string {
	baseName := g.extractFileName(filename)
	changeType := stats.ChangeType

	switch changeType {
	case "A":
		return fmt.Sprintf("add %s with %d lines", baseName, stats.Additions)
	case "D":
		return fmt.Sprintf("remove %s (%d lines deleted)", baseName, stats.Deletions)
	case "M":
		if stats.Additions > stats.Deletions*2 {
			return fmt.Sprintf("expand %s functionality (+%d lines)", baseName, stats.Additions)
		} else if stats.Deletions > stats.Additions*2 {
			return fmt.Sprintf("simplify %s implementation (-%d lines)", baseName, stats.Deletions)
		} else {
			return fmt.Sprintf("refactor %s (+%d/-%d lines)", baseName, stats.Additions, stats.Deletions)
		}
	default:
		return fmt.Sprintf("update %s", baseName)
	}
}

// detectContextFromWordDiff analyzes word-level changes for context
func (g *Generator) detectContextFromWordDiff(wordDiff string) string {
	contexts := []string{}

	// Look for specific patterns in word diff
	if strings.Contains(wordDiff, "{+error+}") || strings.Contains(wordDiff, "{+Error+}") {
		contexts = append(contexts, "improve error handling")
	}

	if strings.Contains(wordDiff, "{+performance+}") || strings.Contains(wordDiff, "{+optimize+}") {
		contexts = append(contexts, "enhance performance")
	}

	if strings.Contains(wordDiff, "{+test+}") || strings.Contains(wordDiff, "{+Test+}") {
		contexts = append(contexts, "improve test coverage")
	}

	if strings.Contains(wordDiff, "{+security+}") || strings.Contains(wordDiff, "{+validate+}") {
		contexts = append(contexts, "strengthen security")
	}

	if len(contexts) > 0 {
		return strings.Join(contexts, " and ")
	}

	return ""
}

// calculateAdvancedPriority uses comprehensive data for priority calculation
func (g *Generator) calculateAdvancedPriority(changeType string, stats *FileStatistics) int {
	basePriority := g.getTypePriority(changeType)

	// Adjust based on change magnitude
	total := stats.Additions + stats.Deletions
	if total > 100 {
		basePriority -= 1 // Higher priority for large changes
	} else if total < 10 {
		basePriority += 1 // Lower priority for small changes
	}

	return basePriority
}

// enhanceWithNumStats improves FileStatistics accuracy using NumStat data
func (g *Generator) enhanceWithNumStats(result *GitAnalysisResult) {
	for filename, numStat := range result.NumStats {
		if fileStat, exists := result.FileStats[filename]; exists {
			// Use precise NumStat data instead of approximated --stat parsing
			fileStat.Additions = numStat.Additions
			fileStat.Deletions = numStat.Deletions
		}
	}
}

// getDirStats implements: git diff --cached --dirstat=files,0
func (g *Generator) getDirStats(result *GitAnalysisResult) error {
	fmt.Printf("Running `git diff --cached --dirstat=files,0`")

	var cmd *exec.Cmd
	if g.hasPreviousCommits() {
		cmd = exec.Command("git", "diff", "HEAD~1", "HEAD", "--dirstat=files,0")
	} else {
		cmd = exec.Command("git", "diff", "--cached", "--dirstat=files,0")
	}

	output, err := cmd.Output()
	if err != nil {
		// Fallback to cached
		cmd = exec.Command("git", "diff", "--cached", "--dirstat=files,0")
		output, err = cmd.Output()
		if err != nil {
			fmt.Println(" ❌")
			return fmt.Errorf("failed to get dir stats: %w", err)
		}
	}
	fmt.Println(" ✅")

	// Parse dirstat output: " 28.5% pkg/semantic/plugins/"
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			percentStr := strings.TrimSuffix(parts[0], "%")
			if percent, err := strconv.ParseFloat(percentStr, 64); err == nil {
				directory := parts[1]
				result.DirStats[directory] = percent
			}
		}
	}

	return nil
}

// getNumStats implements: git diff --cached --numstat
func (g *Generator) getNumStats(result *GitAnalysisResult) error {
	fmt.Printf("Running `git diff --cached --numstat`")

	var cmd *exec.Cmd
	if g.hasPreviousCommits() {
		cmd = exec.Command("git", "diff", "HEAD~1", "HEAD", "--numstat")
	} else {
		cmd = exec.Command("git", "diff", "--cached", "--numstat")
	}

	output, err := cmd.Output()
	if err != nil {
		// Fallback to cached
		cmd = exec.Command("git", "diff", "--cached", "--numstat")
		output, err = cmd.Output()
		if err != nil {
			fmt.Println(" ❌")
			return fmt.Errorf("failed to get numstat: %w", err)
		}
	}
	fmt.Println(" ✅")

	// Parse numstat output: "78	78	pkg/ccgen/advanced_git_analyzer.go"
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 3 {
			additions, err1 := strconv.Atoi(parts[0])
			deletions, err2 := strconv.Atoi(parts[1])
			filename := parts[2]

			if err1 == nil && err2 == nil {
				result.NumStats[filename] = &NumStat{
					Additions: additions,
					Deletions: deletions,
					Filename:  filename,
				}
			}
		}
	}

	return nil
}

// getFileSummaries implements: git diff --cached --summary
func (g *Generator) getFileSummaries(result *GitAnalysisResult) error {
	fmt.Printf("Running `git diff --cached --summary`")

	var cmd *exec.Cmd
	if g.hasPreviousCommits() {
		cmd = exec.Command("git", "diff", "HEAD~1", "HEAD", "--summary")
	} else {
		cmd = exec.Command("git", "diff", "--cached", "--summary")
	}

	output, err := cmd.Output()
	if err != nil {
		// Fallback to cached
		cmd = exec.Command("git", "diff", "--cached", "--summary")
		output, err = cmd.Output()
		if err != nil {
			fmt.Println(" ❌")
			return fmt.Errorf("failed to get summary: %w", err)
		}
	}
	fmt.Println(" ✅")

	// Parse summary output: " create mode 100644 pkg/semantic/plugins/terraform_changeset_analyzer.go"
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			result.FileSummaries = append(result.FileSummaries, line)
		}
	}

	return nil
}

// extractFunctionContexts implements: git diff --cached --function-context --unified=0 | sed -n 's/^@@.* \(.*\) @@/\1/p' | sort -u | head -n 10
func (g *Generator) extractFunctionContexts(result *GitAnalysisResult) error {
	fmt.Printf("Running `git diff --cached --function-context --unified=0`")

	var cmd *exec.Cmd
	if g.hasPreviousCommits() {
		cmd = exec.Command("git", "diff", "HEAD~1", "HEAD", "--function-context", "--unified=0")
	} else {
		cmd = exec.Command("git", "diff", "--cached", "--function-context", "--unified=0")
	}

	output, err := cmd.Output()
	if err != nil {
		// Fallback to cached
		cmd = exec.Command("git", "diff", "--cached", "--function-context", "--unified=0")
		output, err = cmd.Output()
		if err != nil {
			fmt.Println(" ❌")
			return fmt.Errorf("failed to get function context: %w", err)
		}
	}
	fmt.Println(" ✅")

	// Extract function names from @@ lines using regex
	lines := strings.Split(string(output), "\n")
	functionMap := make(map[string]bool) // Use map to deduplicate

	for _, line := range lines {
		if strings.HasPrefix(line, "@@") && strings.HasSuffix(line, "@@") {
			// Extract function name from: "@@ -1,2 +1,3 @@ func methodName"
			parts := strings.Split(line, "@@")
			if len(parts) >= 3 {
				functionName := strings.TrimSpace(parts[2])
				if functionName != "" {
					functionMap[functionName] = true
				}
			}
		}
	}

	// Convert map to slice and limit to 10
	count := 0
	for funcName := range functionMap {
		if count >= 10 {
			break
		}
		result.ModifiedFunctions = append(result.ModifiedFunctions, funcName)
		count++
	}

	return nil
}
