package jira

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestManager_SetAndGetJiraTicket(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	tests := []struct {
		name       string
		ticketID   string
		wantErr    bool
		expectText string
	}{
		{
			name:       "valid CGC ticket",
			ticketID:   "CGC-1234",
			wantErr:    false,
			expectText: "CGC-1234",
		},
		{
			name:       "valid PROJ ticket",
			ticketID:   "PROJ-999",
			wantErr:    false,
			expectText: "PROJ-999",
		},
		{
			name:       "lowercase ticket gets uppercased",
			ticketID:   "cgc-1234",
			wantErr:    false,
			expectText: "CGC-1234",
		},
		{
			name:     "invalid format - no hyphen",
			ticketID: "CGC1234",
			wantErr:  true,
		},
		{
			name:     "invalid format - too many digits",
			ticketID: "CGC-123456",
			wantErr:  true,
		},
		{
			name:     "invalid format - no letters",
			ticketID: "123-456",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.SetJiraTicket(tt.ticketID)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetJiraTicket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Get the ticket and verify it matches expected
				currentTicket, err := manager.GetCurrentJiraTicket()
				if err != nil {
					t.Errorf("GetCurrentJiraTicket() error = %v", err)
					return
				}
				if currentTicket != tt.expectText {
					t.Errorf("GetCurrentJiraTicket() = %v, want %v", currentTicket, tt.expectText)
				}
			}
		})
	}
}

func TestManager_MultipleTicketsCommentOut(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	// Set first ticket
	err := manager.SetJiraTicket("CGC-1234")
	if err != nil {
		t.Fatalf("Failed to set first ticket: %v", err)
	}

	// Set second ticket
	err = manager.SetJiraTicket("PROJ-5678")
	if err != nil {
		t.Fatalf("Failed to set second ticket: %v", err)
	}

	// Verify current ticket is the second one
	currentTicket, err := manager.GetCurrentJiraTicket()
	if err != nil {
		t.Fatalf("Failed to get current ticket: %v", err)
	}
	if currentTicket != "PROJ-5678" {
		t.Errorf("Expected current ticket to be PROJ-5678, got %v", currentTicket)
	}

	// Read file content and verify first ticket is commented out
	content, err := manager.readJiraRefFile()
	if err != nil {
		t.Fatalf("Failed to read jira ref file: %v", err)
	}

	if !strings.Contains(content, "PROJ-5678") {
		t.Error("File should contain current ticket PROJ-5678")
	}
	if !strings.Contains(content, "# CGC-1234") {
		t.Error("File should contain commented out previous ticket # CGC-1234")
	}
}

func TestManager_ClearJiraTicket(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	// Set a ticket first
	err := manager.SetJiraTicket("CGC-1234")
	if err != nil {
		t.Fatalf("Failed to set ticket: %v", err)
	}

	// Clear the ticket
	err = manager.ClearJiraTicket()
	if err != nil {
		t.Fatalf("Failed to clear ticket: %v", err)
	}

	// Verify no current ticket
	currentTicket, err := manager.GetCurrentJiraTicket()
	if err != nil {
		t.Fatalf("Failed to get current ticket: %v", err)
	}
	if currentTicket != "" {
		t.Errorf("Expected no current ticket, got %v", currentTicket)
	}

	// Verify file still exists and contains commented out ticket
	content, err := manager.readJiraRefFile()
	if err != nil {
		t.Fatalf("Failed to read jira ref file: %v", err)
	}

	if !strings.Contains(content, "# CGC-1234") {
		t.Error("File should contain commented out previous ticket # CGC-1234")
	}
}

func TestManager_GetJiraRefFileInfo(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	// Initially file shouldn't exist
	exists, path, err := manager.GetJiraRefFileInfo()
	if err != nil {
		t.Fatalf("GetJiraRefFileInfo() error = %v", err)
	}
	if exists {
		t.Error("File should not exist initially")
	}
	expectedPath := filepath.Join(tempDir, JiraRefFile)
	if path != expectedPath {
		t.Errorf("Expected path %v, got %v", expectedPath, path)
	}

	// Set a ticket (creates file)
	err = manager.SetJiraTicket("CGC-1234")
	if err != nil {
		t.Fatalf("Failed to set ticket: %v", err)
	}

	// Now file should exist
	exists, path, err = manager.GetJiraRefFileInfo()
	if err != nil {
		t.Fatalf("GetJiraRefFileInfo() error = %v", err)
	}
	if !exists {
		t.Error("File should exist after setting ticket")
	}
	if path != expectedPath {
		t.Errorf("Expected path %v, got %v", expectedPath, path)
	}
}

func TestManager_CreateEmptyJiraRefFile(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	// Try to get ticket when file doesn't exist (should create empty file)
	currentTicket, err := manager.GetCurrentJiraTicket()
	if err != nil {
		t.Fatalf("GetCurrentJiraTicket() error = %v", err)
	}
	if currentTicket != "" {
		t.Errorf("Expected empty ticket, got %v", currentTicket)
	}

	// File should now exist
	exists, _, err := manager.GetJiraRefFileInfo()
	if err != nil {
		t.Fatalf("GetJiraRefFileInfo() error = %v", err)
	}
	if !exists {
		t.Error("File should exist after GetCurrentJiraTicket() call")
	}
}

func TestManager_IsValidJiraFormat(t *testing.T) {
	manager := NewManager("")

	tests := []struct {
		ticketID string
		want     bool
	}{
		{"CGC-1234", true},
		{"PROJ-999", true},
		{"ABC-12345", true},
		{"A-1", false},           // too few letters
		{"ABCDEFGHIJK-1", false}, // too many letters
		{"CGC-123456", false},    // too many digits
		{"CGC-", false},          // no digits
		{"CGC1234", false},       // no hyphen
		{"cgc-1234", true},       // lowercase (will be converted)
		{"", false},              // empty
		{"123-456", false},       // no letters
	}

	for _, tt := range tests {
		t.Run(tt.ticketID, func(t *testing.T) {
			got := manager.isValidJiraFormat(tt.ticketID)
			if got != tt.want {
				t.Errorf("isValidJiraFormat(%v) = %v, want %v", tt.ticketID, got, tt.want)
			}
		})
	}
}

func TestManager_FileContent(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	// Set initial ticket
	err := manager.SetJiraTicket("CGC-1234")
	if err != nil {
		t.Fatalf("Failed to set ticket: %v", err)
	}

	// Read and verify file structure
	content, err := manager.readJiraRefFile()
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	expectedParts := []string{
		"# JIRA Commit Reference - Updated:",
		"# Current active ticket:",
		"CGC-1234",
	}

	for _, part := range expectedParts {
		if !strings.Contains(content, part) {
			t.Errorf("File content missing expected part: %v\nContent:\n%v", part, content)
		}
	}
}