// Package jira provides JIRA ticket management for commit messages
package jira

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	JiraRefFile = "jira-commit-ref.txt"
)

// Manager handles JIRA ticket reference management
type Manager struct {
	repoPath string
}

// NewManager creates a new JIRA ticket manager
func NewManager(repoPath string) *Manager {
	return &Manager{
		repoPath: repoPath,
	}
}

// SetJiraTicket sets the current JIRA ticket, commenting out previous entries
func (m *Manager) SetJiraTicket(ticketID string) error {
	// Validate JIRA ticket format (e.g., CGC-1245)
	if !m.isValidJiraFormat(ticketID) {
		return fmt.Errorf("invalid JIRA ticket format: %s (expected format: XXX-####)", ticketID)
	}

	ticketID = strings.ToUpper(ticketID)

	// Read existing content (empty if file doesn't exist)
	existingContent, err := m.readJiraRefFile()
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read JIRA reference file: %w", err)
	}
	if os.IsNotExist(err) {
		existingContent = ""
	}

	// Comment out existing entries and add new one
	var newContent strings.Builder
	newContent.WriteString(fmt.Sprintf("# JIRA Commit Reference - Updated: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	newContent.WriteString("# Current active ticket:\n")
	newContent.WriteString(fmt.Sprintf("%s\n", ticketID))

	if existingContent != "" {
		newContent.WriteString("\n# Previous tickets (commented out):\n")
		lines := strings.Split(existingContent, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				// Comment out previous active tickets
				if m.isValidJiraFormat(line) {
					newContent.WriteString(fmt.Sprintf("# %s\n", line))
				}
			} else if strings.HasPrefix(line, "#") {
				// Keep existing comments
				newContent.WriteString(fmt.Sprintf("%s\n", line))
			}
		}
	}

	// Write to file
	return m.writeJiraRefFile(newContent.String())
}

// GetCurrentJiraTicket returns the current active JIRA ticket
func (m *Manager) GetCurrentJiraTicket() (string, error) {
	content, err := m.readJiraRefFile()
	if err != nil {
		// If file doesn't exist, create an empty one
		if os.IsNotExist(err) {
			if createErr := m.createEmptyJiraRefFile(); createErr != nil {
				return "", fmt.Errorf("failed to create JIRA reference file: %w", createErr)
			}
			return "", nil
		}
		return "", fmt.Errorf("failed to read JIRA reference file: %w", err)
	}

	// Find the first non-commented line that's a valid JIRA ticket
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			if m.isValidJiraFormat(line) {
				return strings.ToUpper(line), nil
			}
		}
	}

	return "", nil
}

// ShowJiraStatus displays the current JIRA ticket status
func (m *Manager) ShowJiraStatus() error {
	currentTicket, err := m.GetCurrentJiraTicket()
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("## 🎫 JIRA Ticket Status")
	fmt.Println()

	if currentTicket == "" {
		fmt.Println("**Current ticket:** None set")
		fmt.Println()
		fmt.Println("Use `cc set-jira CGC-1234` to set a JIRA ticket for commits.")
	} else {
		fmt.Printf("**Current ticket:** `%s`\n", currentTicket)
		fmt.Println()
		fmt.Printf("This ticket will be automatically included in commit messages.\n")
		fmt.Printf("Use `cc set-jira NEW-TICKET` to change or `cc clear-jira` to remove.\n")
	}

	return nil
}

// ClearJiraTicket removes the current JIRA ticket
func (m *Manager) ClearJiraTicket() error {
	// Read existing content to preserve history
	existingContent, err := m.readJiraRefFile()
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read JIRA reference file: %w", err)
	}

	// Create new content with no active ticket
	var newContent strings.Builder
	newContent.WriteString(fmt.Sprintf("# JIRA Commit Reference - Cleared: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	newContent.WriteString("# No active ticket set\n")

	if existingContent != "" {
		newContent.WriteString("\n# Previous tickets (commented out):\n")
		lines := strings.Split(existingContent, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				// Comment out any active tickets
				if m.isValidJiraFormat(line) {
					newContent.WriteString(fmt.Sprintf("# %s\n", line))
				}
			} else if strings.HasPrefix(line, "#") {
				// Keep existing comments
				newContent.WriteString(fmt.Sprintf("%s\n", line))
			}
		}
	}

	return m.writeJiraRefFile(newContent.String())
}

// getJiraRefFilePath returns the path to the JIRA reference file
func (m *Manager) getJiraRefFilePath() string {
	// Clean the path to prevent directory traversal attacks
	cleanRepoPath := filepath.Clean(m.repoPath)
	return filepath.Join(cleanRepoPath, JiraRefFile)
}

// readJiraRefFile reads the content of the JIRA reference file
func (m *Manager) readJiraRefFile() (string, error) {
	filePath := m.getJiraRefFilePath()
	
	// Validate that the file path is within the repository directory
	absRepoPath, err := filepath.Abs(m.repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve repository path: %w", err)
	}
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve file path: %w", err)
	}
	
	// Ensure the file is within the repository directory
	if !strings.HasPrefix(absFilePath, absRepoPath+string(filepath.Separator)) {
		return "", fmt.Errorf("file access outside repository directory not allowed")
	}
	
	// Additional validation: ensure we're only reading the specific JIRA reference file
	if filepath.Base(absFilePath) != JiraRefFile {
		return "", fmt.Errorf("unauthorized file access: only %s is allowed", JiraRefFile)
	}
	
	// Construct safe path directly from validated repository path
	safePath := filepath.Join(absRepoPath, JiraRefFile)
	content, err := os.ReadFile(safePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// writeJiraRefFile writes content to the JIRA reference file
func (m *Manager) writeJiraRefFile(content string) error {
	filePath := m.getJiraRefFilePath()
	return os.WriteFile(filePath, []byte(content), 0600)
}

// createEmptyJiraRefFile creates an empty JIRA reference file
func (m *Manager) createEmptyJiraRefFile() error {
	content := fmt.Sprintf("# JIRA Commit Reference - Created: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	content += "# No active ticket set\n"
	content += "# Use 'cc set-jira CGC-1234' to set a JIRA ticket for commits\n"
	return m.writeJiraRefFile(content)
}

// isValidJiraFormat validates JIRA ticket format (e.g., CGC-1245)
func (m *Manager) isValidJiraFormat(ticketID string) bool {
	// Pattern: 2-10 uppercase letters, hyphen, 1-5 digits
	pattern := `^[A-Z]{2,10}-\d{1,5}$`
	matched, err := regexp.MatchString(pattern, strings.ToUpper(ticketID))
	if err != nil {
		return false // Invalid regex pattern
	}
	return matched
}

// GetJiraRefFileInfo returns information about the JIRA reference file
func (m *Manager) GetJiraRefFileInfo() (exists bool, path string, err error) {
	filePath := m.getJiraRefFilePath()
	_, err = os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, filePath, nil
		}
		return false, filePath, err
	}
	return true, filePath, nil
}

// ListJiraHistory shows the history of JIRA tickets from the file
func (m *Manager) ListJiraHistory() error {
	content, err := m.readJiraRefFile()
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("**No JIRA reference file found.** Use `cc set-jira CGC-1234` to create one.")
			return nil
		}
		return fmt.Errorf("failed to read JIRA reference file: %w", err)
	}

	fmt.Println()
	fmt.Println("## 📋 JIRA Ticket History")
	fmt.Println()
	fmt.Printf("**File location:** `%s`\n", m.getJiraRefFilePath())
	fmt.Println()
	fmt.Println("```")
	fmt.Print(content)
	fmt.Println("```")

	return nil
}
