package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
	"unicode/utf8"
	
	"github.com/greenstevester/fast-cc-git-hooks/internal/banner"
)

const (
	maxSubjectLength  = 50
	maxBodyLineLength = 72
)

type ChangeType struct {
	Type        string
	Scope       string
	Description string
	Files       []string
	Priority    int
}

var (
	version   = "dev"
	buildTime = "unknown"
	commit    = "unknown"

	// Command line flags.
	noVerify = flag.Bool("no-verify", false, "Skip pre-commit hooks")
	execute  = flag.Bool("execute", false, "Execute the commit after generating message")
	verbose  = flag.Bool("verbose", false, "Show detailed analysis")
	help     = flag.Bool("help", false, "Show help")
)

func main() {
	// Print banner with terminal-appropriate formatting
	banner.Print()
	
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	if !isGitRepo() {
		log.Fatal("Not a git repository")
	}

	// Get git status and diffs.
	status, err := getGitStatus()
	if err != nil {
		log.Fatalf("Failed to get git status: %v", err)
	}

	if *verbose {
		fmt.Println("Git status:")
		fmt.Println(status)
		fmt.Println()
	}

	// Add all changes.
	if addErr := addAllChanges(); addErr != nil {
		log.Fatalf("Failed to add changes: %v", addErr)
	}

	// Get staged diff.
	diff, err := getStagedDiff()
	if err != nil {
		log.Fatalf("Failed to get diff: %v", err)
	}

	if strings.TrimSpace(diff) == "" {
		fmt.Println("No changes to commit")
		return
	}

	// Analyze changes.
	changes := analyzeDiff(diff)
	if *verbose {
		fmt.Println("Detected changes:")
		for _, change := range changes {
			fmt.Printf("- %s(%s): %s (files: %v)\n",
				change.Type, change.Scope, change.Description, change.Files)
		}
		fmt.Println()
	}

	// Generate commit message.
	message := generateCommitMessage(changes)
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("\n>>> based on your changes, cc created the following git commit message for you:")
	fmt.Println(message)

	if *execute {
		if err := executeCommit(message); err != nil {
			log.Fatalf("Failed to commit: %v", err)
		}
		fmt.Println("Commit created successfully!")
	}
}

func showHelp() {
	fmt.Printf("cc - Git Commit message generator v%s\n\n", version)
	fmt.Println("Analyzes staged changes and generates conventional commit messages.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  cc [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --execute      Execute the commit after generating message")
	fmt.Println("  --no-verify    Skip pre-commit hooks when committing")
	fmt.Println("  --verbose      Show detailed analysis of changes")
	fmt.Println("  --help         Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  cc                    # Generate commit message only")
	fmt.Println("  cc --execute          # Generate and commit")
	fmt.Println("  cc --verbose          # Show detailed analysis")
	fmt.Println("  cc --execute --no-verify  # Commit without hooks")
	fmt.Println()
	fmt.Printf("Build info: %s (%s)\n", buildTime, commit)
}

func isGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

func getGitStatus() (string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	return string(output), err
}

func addAllChanges() error {
	cmd := exec.Command("git", "add", ".")
	return cmd.Run()
}

func getStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--staged")
	output, err := cmd.Output()
	return string(output), err
}

func analyzeDiff(diff string) []ChangeType {
	fileChanges := make(map[string]*ChangeType)

	// Parse diff by files.
	files := strings.Split(diff, "diff --git")
	changes := make([]ChangeType, 0, len(files))
	for _, file := range files {
		if strings.TrimSpace(file) == "" {
			continue
		}

		change := analyzeFileChange(file)
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

func analyzeFileChange(fileDiff string) *ChangeType {
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
	changeType, scope := determineTypeAndScope(filename, fileDiff)
	description := generateDescription(filename, fileDiff, changeType)

	return &ChangeType{
		Type:        changeType,
		Scope:       scope,
		Description: description,
		Files:       []string{filename},
		Priority:    getTypePriority(changeType),
	}
}

func determineTypeAndScope(filename, diff string) (changeType, scope string) {
	scope = determineScope(filename)
	changeType = determineType(filename, diff)
	return changeType, scope
}

// determineScope determines the scope from the filename.
func determineScope(filename string) string {
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
func determineType(filename, diff string) string {
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
	case countAdditions(diff) > countDeletions(diff):
		return "feat"
	default:
		return "refactor"
	}
}

func generateDescription(filename, diff, changeType string) string {
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

func countAdditions(diff string) int {
	count := 0
	for _, line := range strings.Split(diff, "\n") {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			count++
		}
	}
	return count
}

func countDeletions(diff string) int {
	count := 0
	for _, line := range strings.Split(diff, "\n") {
		if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			count++
		}
	}
	return count
}

func getTypePriority(changeType string) int {
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

func generateCommitMessage(changes []ChangeType) string {
	if len(changes) == 0 {
		return "chore: update files"
	}

	// Use the highest priority change as primary.
	primary := changes[0]

	// Create subject line.
	subject := primary.Type
	if primary.Scope != "" {
		subject += fmt.Sprintf("(%s)", primary.Scope)
	}
	subject += fmt.Sprintf(": %s", primary.Description)

	// Truncate subject if too long.
	if utf8.RuneCountInString(subject) > maxSubjectLength {
		runes := []rune(subject)
		if len(runes) > maxSubjectLength-3 {
			subject = string(runes[:maxSubjectLength-3]) + "..."
		}
	}

	// Generate body for multiple changes or complex single change.
	var body []string

	if len(changes) > 1 {
		body = append(body, "", "Changes include:")
		for _, change := range changes {
			line := fmt.Sprintf("- %s", capitalizeFirst(change.Description))
			if len(change.Files) > 0 {
				line += fmt.Sprintf(" (%s)", change.Files[0])
				if len(change.Files) > 1 {
					line += fmt.Sprintf(" and %d more", len(change.Files)-1)
				}
			}
			body = append(body, wrapLine(line, maxBodyLineLength))
		}
	}


	if len(body) > 0 {
		return subject + strings.Join(body, "\n")
	}

	return subject
}

func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	firstRune, size := utf8.DecodeRuneInString(s)
	upperFirst, _ := utf8.DecodeRuneInString(strings.ToUpper(string(firstRune)))
	return string(upperFirst) + s[size:]
}

func wrapLine(line string, maxLength int) string {
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

func executeCommit(message string) error {
	args := []string{"commit", "-m", message}
	if *noVerify {
		args = append(args, "--no-verify")
	}

	cmd := exec.Command("git", args...) // #nosec G204 - args are validated git commands
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
