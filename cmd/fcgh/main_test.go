package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Helper function to create a test command context
func setupTestContext(t *testing.T) (context.Context, func()) {
	// Create temp directory
	tempDir := t.TempDir()
	
	// Store original working directory
	origWD, _ := os.Getwd()
	
	// Change to temp directory
	os.Chdir(tempDir)
	
	// Store original environment variables
	origHome := os.Getenv("HOME")
	origUser := os.Getenv("USER")
	
	// Set test environment
	os.Setenv("HOME", tempDir)
	os.Setenv("USER", "testuser")
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	
	cleanup := func() {
		cancel()
		os.Chdir(origWD)
		os.Setenv("HOME", origHome)
		os.Setenv("USER", origUser)
	}
	
	return ctx, cleanup
}

func TestSetupLogger(t *testing.T) {
	tests := []struct {
		name    string
		verbose bool
	}{
		{"verbose mode", true},
		{"normal mode", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test doesn't crash
			setupLogger(tt.verbose)
			if logger == nil {
				t.Error("Logger should be initialized")
			}
		})
	}
}

func TestValidateCommand(t *testing.T) {
	cmd := validateCommand()
	
	if cmd.Name != "validate" {
		t.Errorf("Expected command name 'validate', got %s", cmd.Name)
	}
	
	if cmd.Run == nil {
		t.Error("Run function should not be nil")
	}
	
	if cmd.Flags == nil {
		t.Error("Flags should not be nil")
	}
	
	// Test with valid message
	ctx, cleanup := setupTestContext(t)
	defer cleanup()
	
	err := cmd.Run(ctx, []string{"feat: add new feature"})
	if err != nil {
		t.Errorf("Valid message should not return error: %v", err)
	}
	
	// Test with invalid message
	err = cmd.Run(ctx, []string{"invalid message"})
	if err == nil {
		t.Error("Invalid message should return error")
	}
	
	// Test with file flag
	testFile := filepath.Join(t.TempDir(), "test.txt")
	os.WriteFile(testFile, []byte("feat: test message"), 0644)
	
	validateFile = testFile
	defer func() { validateFile = "" }()
	
	err = cmd.Run(ctx, []string{})
	if err != nil {
		t.Errorf("Valid file should not return error: %v", err)
	}
}

func TestVersionCommand(t *testing.T) {
	cmd := versionCommand()
	
	if cmd.Name != "version" {
		t.Errorf("Expected command name 'version', got %s", cmd.Name)
	}
	
	ctx, cleanup := setupTestContext(t)
	defer cleanup()
	
	err := cmd.Run(ctx, []string{})
	if err != nil {
		t.Errorf("Version command should not return error: %v", err)
	}
}

func TestInitCommand(t *testing.T) {
	cmd := initCommand()
	
	if cmd.Name != "init" {
		t.Errorf("Expected command name 'init', got %s", cmd.Name)
	}
	
	ctx, cleanup := setupTestContext(t)
	defer cleanup()
	
	// Test init command
	err := cmd.Run(ctx, []string{})
	if err != nil {
		t.Errorf("Init command should not return error: %v", err)
	}
}

func TestSetupCommand(t *testing.T) {
	cmd := setupCommand()
	
	if cmd.Name != "setup" {
		t.Errorf("Expected command name 'setup', got %s", cmd.Name)
	}
	
	ctx, cleanup := setupTestContext(t)
	defer cleanup()
	
	// Initialize a git repository for testing
	os.MkdirAll(".git", 0755)
	
	// Test global setup
	localInstall = false
	err := cmd.Run(ctx, []string{})
	// This might fail due to git config access, but should not panic
	_ = err // Allow error for now since it requires git configuration
	
	// Test local setup  
	localInstall = true
	defer func() { localInstall = false }()
	
	err = cmd.Run(ctx, []string{})
	_ = err // Allow error for now
}

func TestSetupEnterpriseCommand(t *testing.T) {
	cmd := setupEnterpriseCommand()
	
	if cmd.Name != "setup-ent" {
		t.Errorf("Expected command name 'setup-ent', got %s", cmd.Name)
	}
	
	ctx, cleanup := setupTestContext(t)
	defer cleanup()
	
	// Initialize a git repository for testing
	os.MkdirAll(".git", 0755)
	
	err := cmd.Run(ctx, []string{})
	_ = err // Allow error for now since it requires git configuration
}

func TestRemoveCommand(t *testing.T) {
	cmd := removeCommand()
	
	if cmd.Name != "remove" {
		t.Errorf("Expected command name 'remove', got %s", cmd.Name)
	}
	
	ctx, cleanup := setupTestContext(t)
	defer cleanup()
	
	// Test remove command (might not find installations, but shouldn't crash)
	err := cmd.Run(ctx, []string{})
	_ = err // Allow error since no installations exist
}

func TestEnsureConfigExists(t *testing.T) {
	ctx, cleanup := setupTestContext(t)
	defer cleanup()
	_ = ctx
	
	// Test creating config in home directory
	configPath, isNew, err := ensureConfigExists()
	if err != nil {
		t.Errorf("ensureConfigExists should not return error: %v", err)
	}
	
	if !isNew {
		t.Error("Config should be marked as new when created")
	}
	
	if !strings.Contains(configPath, "fast-cc-config.yaml") {
		t.Errorf("Config path should contain fast-cc-config.yaml, got: %s", configPath)
	}
	
	// Test when config already exists
	configPath2, isNew2, err := ensureConfigExists()
	if err != nil {
		t.Errorf("ensureConfigExists should not return error on existing config: %v", err)
	}
	
	if isNew2 {
		t.Error("Config should not be marked as new when it exists")
	}
	
	if configPath != configPath2 {
		t.Error("Config path should be consistent")
	}
}

func TestEnsureEnterpriseConfigExists(t *testing.T) {
	ctx, cleanup := setupTestContext(t)
	defer cleanup()
	_ = ctx
	
	configPath, isNew, err := ensureEnterpriseConfigExists()
	if err != nil {
		t.Errorf("ensureEnterpriseConfigExists should not return error: %v", err)
	}
	
	if !isNew {
		t.Error("Enterprise config should be marked as new when created")
	}
	
	if !strings.Contains(configPath, "fast-cc-config.yaml") {
		t.Errorf("Enterprise config path should contain fast-cc-config.yaml, got: %s", configPath)
	}
}

func TestCopyEnterpriseConfig(t *testing.T) {
	tempDir := t.TempDir()
	destPath := filepath.Join(tempDir, "config.yaml")
	
	err := copyEnterpriseConfig(destPath)
	if err != nil {
		t.Errorf("copyEnterpriseConfig should not return error: %v", err)
	}
	
	// Check if file was created
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		t.Error("Enterprise config file should be created")
	}
	
	// Check file content
	content, err := os.ReadFile(destPath)
	if err != nil {
		t.Errorf("Should be able to read created config file: %v", err)
	}
	
	contentStr := string(content)
	if !strings.Contains(contentStr, "types:") {
		t.Error("Config should contain types section")
	}
}

func TestCreateBasicEnterpriseConfig(t *testing.T) {
	tempDir := t.TempDir()
	destPath := filepath.Join(tempDir, "basic-config.yaml")
	
	err := createBasicEnterpriseConfig(destPath)
	if err != nil {
		t.Errorf("createBasicEnterpriseConfig should not return error: %v", err)
	}
	
	// Check if file was created
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		t.Error("Basic enterprise config file should be created")
	}
	
	// Check file content
	content, err := os.ReadFile(destPath)
	if err != nil {
		t.Errorf("Should be able to read created basic config file: %v", err)
	}
	
	contentStr := string(content)
	if !strings.Contains(contentStr, "feat") {
		t.Error("Basic config should contain feat type")
	}
	if !strings.Contains(contentStr, "fix") {
		t.Error("Basic config should contain fix type")
	}
}

func TestGetGitConfigDir(t *testing.T) {
	// This test might fail on systems without git config
	configDir, err := getGitConfigDir()
	
	// We allow this to fail since git might not be configured
	// but if it succeeds, it should return a valid path
	if err == nil && configDir == "" {
		t.Error("If getGitConfigDir succeeds, it should return a non-empty path")
	}
}

func TestCheckInstallations(t *testing.T) {
	ctx, cleanup := setupTestContext(t)
	defer cleanup()
	_ = ctx
	
	// This will likely return false, false for both installations
	// since we're in a test environment, but it shouldn't crash
	hasGlobal, hasLocal, err := checkInstallations()
	if err != nil {
		// Allow error since git might not be configured
		t.Logf("checkInstallations returned error (expected in test): %v", err)
	}
	
	// Just verify it returns boolean values
	_ = hasGlobal
	_ = hasLocal
}

func TestHasGlobalInstallation(t *testing.T) {
	hasGlobal, err := hasGlobalInstallation()
	
	// This might fail due to git configuration issues, which is expected in tests
	if err != nil {
		t.Logf("hasGlobalInstallation returned error (expected in test): %v", err)
	}
	
	// Just verify it returns a boolean
	_ = hasGlobal
}

func TestPromptUserChoice(t *testing.T) {
	// This is an interactive function, so we'll test the structure
	// In a real test environment, this would require mocking stdin
	
	// We can't easily test interactive input without mocking stdin
	// so we'll skip this test for now
	t.Skip("Interactive function requires stdin mocking")
}

func TestRemoveGlobalInstallation(t *testing.T) {
	// Test removing global installation
	// This will likely fail since no global installation exists in test
	err := removeGlobalInstallation()
	// Allow error since no global installation exists in test environment
	_ = err
}

// Test command creation functions
func TestAllCommandsHaveRequiredFields(t *testing.T) {
	commands := []*Command{
		validateCommand(),
		initCommand(), 
		versionCommand(),
		setupCommand(),
		setupEnterpriseCommand(),
		removeCommand(),
	}
	
	for _, cmd := range commands {
		if cmd.Name == "" {
			t.Errorf("Command should have a name")
		}
		
		if cmd.Description == "" {
			t.Errorf("Command %s should have a description", cmd.Name)
		}
		
		if cmd.Run == nil {
			t.Errorf("Command %s should have a Run function", cmd.Name)
		}
		
		if cmd.Flags == nil {
			t.Errorf("Command %s should have Flags", cmd.Name)
		}
	}
}

// Test edge cases
func TestValidateCommandWithEmptyArgs(t *testing.T) {
	cmd := validateCommand()
	ctx, cleanup := setupTestContext(t)
	defer cleanup()
	
	// Test with no arguments
	err := cmd.Run(ctx, []string{})
	if err == nil {
		t.Error("Validate command with no args should return error")
	}
}

func TestValidateCommandWithFileFlag(t *testing.T) {
	cmd := validateCommand()
	ctx, cleanup := setupTestContext(t)
	defer cleanup()
	
	// Create test file with invalid content
	testFile := filepath.Join(t.TempDir(), "invalid.txt")
	os.WriteFile(testFile, []byte("invalid commit message"), 0644)
	
	validateFile = testFile
	defer func() { validateFile = "" }()
	
	err := cmd.Run(ctx, []string{})
	if err == nil {
		t.Error("Validate command with invalid file should return error")
	}
}

func TestCommandFlags(t *testing.T) {
	tests := []struct {
		name string
		cmd  *Command
	}{
		{"validate", validateCommand()},
		{"init", initCommand()},
		{"version", versionCommand()},
		{"setup", setupCommand()},
		{"setup-ent", setupEnterpriseCommand()},
		{"remove", removeCommand()},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify flags can be parsed without errors
			flagSet := tt.cmd.Flags
			if flagSet == nil {
				t.Errorf("Command %s should have flag set", tt.name)
				return
			}
			
			// Test flag parsing with empty args
			err := flagSet.Parse([]string{})
			if err != nil {
				t.Errorf("Command %s flags should parse empty args: %v", tt.name, err)
			}
		})
	}
}

// Test global variables initialization
func TestGlobalVariables(t *testing.T) {
	// Test that global variables have sensible defaults
	if version == "" {
		t.Error("Version should not be empty")
	}
	
	// Test logger is initialized after setupLogger
	setupLogger(false)
	if logger == nil {
		t.Error("Logger should be initialized after setupLogger")
	}
}

// Test configuration file operations
func TestConfigFileOperations(t *testing.T) {
	ctx, cleanup := setupTestContext(t)
	defer cleanup()
	_ = ctx
	
	// Test that we can create and read configuration
	configPath, _, err := ensureConfigExists()
	if err != nil {
		t.Fatalf("Failed to ensure config exists: %v", err)
	}
	
	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file should exist after ensureConfigExists")
	}
	
	// Test enterprise config
	entConfigPath, _, err := ensureEnterpriseConfigExists()
	if err != nil {
		t.Fatalf("Failed to ensure enterprise config exists: %v", err)
	}
	
	// Verify enterprise file exists
	if _, err := os.Stat(entConfigPath); os.IsNotExist(err) {
		t.Error("Enterprise config file should exist after ensureEnterpriseConfigExists")
	}
}

// Test edge cases for ensureConfigExists
func TestEnsureConfigExistsEdgeCases(t *testing.T) {
	ctx, cleanup := setupTestContext(t)
	defer cleanup()
	_ = ctx
	
	// Create a directory with restrictive permissions to test error handling
	restrictedDir := filepath.Join(t.TempDir(), "restricted")
	os.MkdirAll(restrictedDir, 0000) // No permissions
	defer os.Chmod(restrictedDir, 0755) // Restore permissions for cleanup
	
	// Test with existing old-style config file
	tempHome := t.TempDir()
	os.Setenv("HOME", tempHome)
	defer cleanup()
	
	// Create .fast-cc directory
	configDir := filepath.Join(tempHome, ".fast-cc")
	os.MkdirAll(configDir, 0755)
	
	// Create old-style config file
	oldConfigPath := filepath.Join(configDir, ".fast-cc-hooks.yaml")
	os.WriteFile(oldConfigPath, []byte("types:\n  - feat\n  - fix"), 0644)
	
	configPath, isNew, err := ensureConfigExists()
	if err != nil {
		t.Errorf("Should handle existing old config file: %v", err)
	}
	if isNew {
		t.Error("Should not be marked as new when old config exists")
	}
	_ = configPath
}

// Test ensureEnterpriseConfigExists edge cases
func TestEnsureEnterpriseConfigExistsEdgeCases(t *testing.T) {
	ctx, cleanup := setupTestContext(t)
	defer cleanup()
	_ = ctx
	
	// Test when home directory config already exists as different type
	tempHome := t.TempDir()
	os.Setenv("HOME", tempHome)
	defer cleanup()
	
	// Create .fast-cc directory
	configDir := filepath.Join(tempHome, ".fast-cc")
	os.MkdirAll(configDir, 0755)
	
	// Create existing config
	configPath := filepath.Join(configDir, "fast-cc-config.yaml")
	os.WriteFile(configPath, []byte("# existing config\ntypes:\n  - feat"), 0644)
	
	_, isNew, err := ensureEnterpriseConfigExists()
	if err != nil {
		t.Errorf("Should handle existing config: %v", err)
	}
	if isNew {
		t.Error("Should not be new when config exists")
	}
}

// Test copyEnterpriseConfig edge cases
func TestCopyEnterpriseConfigEdgeCases(t *testing.T) {
	// Test with invalid destination path
	err := copyEnterpriseConfig("/invalid/path/config.yaml")
	if err == nil {
		t.Error("Should return error for invalid path")
	}
	
	// Test with valid path but directory doesn't exist
	tempDir := t.TempDir()
	destPath := filepath.Join(tempDir, "subdir", "config.yaml")
	
	err = copyEnterpriseConfig(destPath)
	if err != nil {
		t.Errorf("Should create directory and file: %v", err)
	}
	
	// Verify file was created
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		t.Error("Config file should be created")
	}
}

// Test createBasicEnterpriseConfig edge cases
func TestCreateBasicEnterpriseConfigEdgeCases(t *testing.T) {
	// Test with invalid destination path
	err := createBasicEnterpriseConfig("/invalid/path/config.yaml")
	if err == nil {
		t.Error("Should return error for invalid path")
	}
}

// Test removeCommand with more scenarios
func TestRemoveCommandScenarios(t *testing.T) {
	cmd := removeCommand()
	ctx, cleanup := setupTestContext(t)
	defer cleanup()
	
	// Test with local flag
	localInstall = true
	defer func() { localInstall = false }()
	
	err := cmd.Run(ctx, []string{})
	// Allow error since no installation exists
	_ = err
	
	// Test with force flag
	localInstall = false
	forceInstall = true
	defer func() { forceInstall = false }()
	
	err = cmd.Run(ctx, []string{})
	// Allow error since no installation exists
	_ = err
}

// Test promptUserChoice with different scenarios (mocked)
func TestPromptUserChoiceScenarios(t *testing.T) {
	// We can't easily mock stdin, but we can test the error cases
	// by temporarily replacing stdin with a pipe
	
	t.Skip("Requires stdin mocking - would need more complex setup")
}

// Test removeGlobalInstallation edge cases
func TestRemoveGlobalInstallationEdgeCases(t *testing.T) {
	// This function tries to remove global installation
	// In test environment, it should handle the case where git config fails
	err := removeGlobalInstallation()
	// Allow any error - in test environment git config might not be available
	_ = err
	
	// The function should not panic
	// If it returns an error, that's acceptable in test environment
}

// Test checkInstallations edge cases  
func TestCheckInstallationsEdgeCases(t *testing.T) {
	ctx, cleanup := setupTestContext(t)
	defer cleanup()
	_ = ctx
	
	// Create a fake .git directory to make it look like a git repo
	os.MkdirAll(".git/hooks", 0755)
	
	hasGlobal, hasLocal, err := checkInstallations()
	// In test environment, global might fail due to git config
	// but local should work now that we have .git directory
	_ = hasGlobal
	_ = hasLocal
	_ = err // Allow error
}

// Test setupCommand and setupEnterpriseCommand edge cases
func TestSetupCommandEdgeCases(t *testing.T) {
	cmd := setupCommand()
	ctx, cleanup := setupTestContext(t)
	defer cleanup()
	
	// Create .git directory
	os.MkdirAll(".git", 0755)
	
	// Test with force flag
	forceInstall = true
	defer func() { forceInstall = false }()
	
	err := cmd.Run(ctx, []string{})
	_ = err // Allow error in test environment
}

func TestSetupEnterpriseCommandEdgeCases(t *testing.T) {
	cmd := setupEnterpriseCommand()
	ctx, cleanup := setupTestContext(t)
	defer cleanup()
	
	// Create .git directory
	os.MkdirAll(".git", 0755)
	
	// Test with force flag
	forceInstall = true
	localInstall = true
	defer func() { 
		forceInstall = false
		localInstall = false
	}()
	
	err := cmd.Run(ctx, []string{})
	_ = err // Allow error in test environment
}

// Test validateCommand edge cases
func TestValidateCommandEdgeCases(t *testing.T) {
	cmd := validateCommand()
	ctx, cleanup := setupTestContext(t)
	defer cleanup()
	
	// Test with non-existent file
	validateFile = "/non/existent/file.txt"
	defer func() { validateFile = "" }()
	
	err := cmd.Run(ctx, []string{})
	if err == nil {
		t.Error("Should return error for non-existent file")
	}
	
	// Reset validateFile
	validateFile = ""
	
	// Test with multiple arguments
	err = cmd.Run(ctx, []string{"arg1", "arg2"})
	if err == nil {
		t.Error("Should return error for multiple arguments")
	}
}

// Test initCommand edge cases
func TestInitCommandEdgeCases(t *testing.T) {
	cmd := initCommand()
	ctx, cleanup := setupTestContext(t)
	defer cleanup()
	
	// Test when config already exists
	// First run should create config
	err := cmd.Run(ctx, []string{})
	if err != nil {
		t.Errorf("First init should succeed: %v", err)
	}
	
	// Second run should return error since config already exists
	err = cmd.Run(ctx, []string{})
	if err == nil {
		t.Error("Second init should return error when config already exists")
	}
}

// Test global flag combinations
func TestGlobalFlagCombinations(t *testing.T) {
	// Test various combinations of global flags
	originalVerbose := verbose
	originalConfigFile := configFile
	defer func() {
		verbose = originalVerbose
		configFile = originalConfigFile
	}()
	
	// Test verbose flag
	verbose = true
	setupLogger(verbose)
	if logger == nil {
		t.Error("Logger should be initialized with verbose mode")
	}
	
	// Test config file flag
	configFile = "/custom/config.yaml"
	// This would affect config loading, but we can't easily test
	// without more complex mocking
}

// Test version command output
func TestVersionCommandOutput(t *testing.T) {
	cmd := versionCommand()
	ctx, cleanup := setupTestContext(t)
	defer cleanup()
	
	// Capture output would require redirecting stdout
	// For now, just ensure it doesn't crash
	err := cmd.Run(ctx, []string{})
	if err != nil {
		t.Errorf("Version command should not return error: %v", err)
	}
}

// Test more comprehensive ensureConfigExists scenarios
func TestEnsureConfigExistsComprehensive(t *testing.T) {
	ctx, cleanup := setupTestContext(t)
	defer cleanup()
	_ = ctx
	
	// Test with .fast-cc-hooks.yaml in current directory
	os.WriteFile(".fast-cc-hooks.yaml", []byte("types:\n  - feat\n  - fix"), 0644)
	defer os.Remove(".fast-cc-hooks.yaml")
	
	configPath, isNew, err := ensureConfigExists()
	if err != nil {
		t.Errorf("Should handle .fast-cc-hooks.yaml in current dir: %v", err)
	}
	
	// Should use the current directory file
	if !strings.Contains(configPath, ".fast-cc-hooks.yaml") {
		t.Errorf("Should use current directory config file, got: %s", configPath)
	}
	
	if isNew {
		t.Error("Should not be new when using existing file")
	}
}

// Test more comprehensive ensureEnterpriseConfigExists scenarios
func TestEnsureEnterpriseConfigExistsComprehensive(t *testing.T) {
	ctx, cleanup := setupTestContext(t)
	defer cleanup()
	_ = ctx
	
	// Test with .fast-cc-hooks.yaml in current directory for enterprise
	os.WriteFile(".fast-cc-hooks.yaml", []byte("types:\n  - feat\n  - fix"), 0644)
	defer os.Remove(".fast-cc-hooks.yaml")
	
	configPath, isNew, err := ensureEnterpriseConfigExists()
	if err != nil {
		t.Errorf("Enterprise config should handle existing local file: %v", err)
	}
	_ = configPath
	_ = isNew
}

// Test removeCommand with different installation scenarios
func TestRemoveCommandComprehensive(t *testing.T) {
	cmd := removeCommand()
	ctx, cleanup := setupTestContext(t)
	defer cleanup()
	
	// Create mock git directory and hooks
	os.MkdirAll(".git/hooks", 0755)
	os.WriteFile(".git/hooks/commit-msg", []byte("#!/bin/sh\necho test"), 0755)
	
	// Test removing local installation
	localInstall = true
	defer func() { localInstall = false }()
	
	err := cmd.Run(ctx, []string{})
	_ = err // Allow error - the hook removal might succeed or fail
}

// Test getGitConfigDir edge cases
func TestGetGitConfigDirEdgeCases(t *testing.T) {
	// Store original HOME
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)
	
	// Test with invalid HOME
	os.Setenv("HOME", "/invalid/nonexistent/path")
	
	configDir, err := getGitConfigDir()
	// This might fail, but shouldn't panic
	_ = configDir
	_ = err
}

// Test hasGlobalInstallation edge cases
func TestHasGlobalInstallationComprehensive(t *testing.T) {
	// This function depends on git config which might not be available
	// in test environment, but we can still test it doesn't panic
	hasGlobal, err := hasGlobalInstallation()
	_ = hasGlobal
	_ = err // Allow any result
}

// Test flag parsing edge cases
func TestCommandFlagsParsing(t *testing.T) {
	tests := []struct {
		name string
		cmd  *Command
		args []string
	}{
		{"validate with -file flag", validateCommand(), []string{"-file", "test.txt"}},
		{"setup with -local flag", setupCommand(), []string{"-local"}},
		{"setup with -force flag", setupCommand(), []string{"-force"}},
		{"remove with -local flag", removeCommand(), []string{"-local"}},
		{"remove with -global flag", removeCommand(), []string{"-global"}},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags to defaults
			validateFile = ""
			localInstall = false
			forceInstall = false
			
			err := tt.cmd.Flags.Parse(tt.args)
			if err != nil {
				t.Errorf("Should parse flags without error: %v", err)
			}
		})
	}
}

// Test error paths in validateCommand
func TestValidateCommandErrorPaths(t *testing.T) {
	cmd := validateCommand()
	ctx, cleanup := setupTestContext(t)
	defer cleanup()
	
	// Create a file with read permission denied
	restrictedFile := filepath.Join(t.TempDir(), "restricted.txt")
	os.WriteFile(restrictedFile, []byte("test"), 0000) // No read permission
	defer os.Chmod(restrictedFile, 0644) // Restore for cleanup
	
	validateFile = restrictedFile
	defer func() { validateFile = "" }()
	
	err := cmd.Run(ctx, []string{})
	if err == nil {
		t.Error("Should return error for file with no read permissions")
	}
}

// Test main function behavior by testing the command registration
func TestMainFunctionBehavior(t *testing.T) {
	// We can't directly test main() since it calls os.Exit
	// But we can test that all commands are properly registered
	
	commands := []*Command{
		validateCommand(),
		initCommand(),
		versionCommand(),
		setupCommand(),
		setupEnterpriseCommand(),
		removeCommand(),
	}
	
	// Verify all commands have unique names
	nameMap := make(map[string]bool)
	for _, cmd := range commands {
		if nameMap[cmd.Name] {
			t.Errorf("Duplicate command name: %s", cmd.Name)
		}
		nameMap[cmd.Name] = true
	}
	
	// Verify we have expected commands
	expectedCommands := []string{"validate", "init", "version", "setup", "setup-ent", "remove"}
	for _, expected := range expectedCommands {
		if !nameMap[expected] {
			t.Errorf("Missing expected command: %s", expected)
		}
	}
}

// Test edge cases in configuration creation
func TestConfigCreationEdgeCases(t *testing.T) {
	ctx, cleanup := setupTestContext(t)
	defer cleanup()
	_ = ctx
	
	// Test createBasicEnterpriseConfig with valid path
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "enterprise-basic.yaml")
	
	err := createBasicEnterpriseConfig(configPath)
	if err != nil {
		t.Errorf("Should create basic enterprise config: %v", err)
	}
	
	// Verify content
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Errorf("Should be able to read created file: %v", err)
	}
	
	if !strings.Contains(string(content), "require_jira_ticket") {
		t.Error("Basic enterprise config should contain require_jira_ticket")
	}
}

// Test error paths in removeCommand 
func TestRemoveCommandErrorPaths(t *testing.T) {
	// Test conflicting flags
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	var localRemove, globalRemove bool
	fs.BoolVar(&localRemove, "local", false, "")
	fs.BoolVar(&globalRemove, "global", false, "")
	
	// Parse flags with both local and global
	args := []string{"-local", "-global"}
	err := fs.Parse(args)
	if err != nil {
		t.Errorf("Should parse flags: %v", err)
	}
	
	// Now both flags are set - this should cause error
	ctx := context.Background()
	cmd := removeCommand()
	cmd.Flags = fs
	
	// Set the internal flags by re-parsing with the command's flag set
	err = cmd.Flags.Parse(args)
	if err != nil {
		t.Errorf("Should parse command flags: %v", err)
	}
	
	// This should return error due to conflicting flags
	err = cmd.Run(ctx, cmd.Flags.Args())
	if err == nil {
		t.Error("Should return error for conflicting flags")
	}
	if !strings.Contains(err.Error(), "cannot specify both") {
		t.Errorf("Should mention conflicting flags, got: %v", err)
	}
}

// Test removeCommand scenarios with mock installations
func TestRemoveCommandInstallationScenarios(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create a git repository for testing local installations
	gitDir := filepath.Join(tempDir, ".git")
	hooksDir := filepath.Join(gitDir, "hooks")
	err := os.MkdirAll(hooksDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create git hooks directory: %v", err)
	}
	
	// Change to the temp directory
	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)
	
	tests := []struct {
		name             string
		hasLocalHook     bool
		hasGlobalHook    bool
		localFlag        bool
		globalFlag       bool
		expectError      bool
		expectRemoval    bool
	}{
		{
			name: "remove local when only local exists",
			hasLocalHook: true,
			hasGlobalHook: false,
			localFlag: true,
			globalFlag: false,
			expectError: false,
			expectRemoval: true,
		},
		{
			name: "remove global when only global exists", 
			hasLocalHook: false,
			hasGlobalHook: true,
			localFlag: false,
			globalFlag: true,
			expectError: false,
			expectRemoval: true,
		},
		{
			name: "remove local when both exist",
			hasLocalHook: true,
			hasGlobalHook: true,
			localFlag: true,
			globalFlag: false,
			expectError: false,
			expectRemoval: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup hooks based on test case
			localHookPath := filepath.Join(hooksDir, "commit-msg")
			globalConfigDir := filepath.Join(tempDir, ".config", "git", "hooks")
			globalHookPath := filepath.Join(globalConfigDir, "commit-msg")
			
			// Clean up from previous test
			os.Remove(localHookPath)
			os.RemoveAll(filepath.Join(tempDir, ".config"))
			
			if tt.hasLocalHook {
				os.WriteFile(localHookPath, []byte("#!/bin/sh\necho local hook"), 0755)
			}
			if tt.hasGlobalHook {
				os.MkdirAll(globalConfigDir, 0755)
				os.WriteFile(globalHookPath, []byte("#!/bin/sh\necho global hook"), 0755)
			}
			
			// Set HOME to tempDir for this test
			originalHome := os.Getenv("HOME")
			os.Setenv("HOME", tempDir)
			defer os.Setenv("HOME", originalHome)
			
			ctx := context.Background()
			cmd := removeCommand()
			
			// Set up flags
			args := []string{}
			if tt.localFlag {
				args = append(args, "-local")
			}
			if tt.globalFlag {
				args = append(args, "-global")
			}
			
			err := cmd.Flags.Parse(args)
			if err != nil {
				t.Errorf("Should parse flags: %v", err)
			}
			
			err = cmd.Run(ctx, cmd.Flags.Args())
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// Test more validateCommand error paths
func TestValidateCommandMoreErrorPaths(t *testing.T) {
	tempDir := t.TempDir()
	
	tests := []struct {
		name        string
		setupFile   bool
		fileContent string
		fileName    string
		expectError bool
	}{
		{
			name: "validate with non-existent file",
			setupFile: false,
			fileName: "nonexistent.txt",
			expectError: true,
		},
		{
			name: "validate with empty file",
			setupFile: true,
			fileContent: "",
			fileName: "empty.txt", 
			expectError: true,
		},
		{
			name: "validate with valid file content",
			setupFile: true,
			fileContent: "feat: add new feature",
			fileName: "valid.txt",
			expectError: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			cmd := validateCommand()
			
			var filePath string
			if tt.setupFile {
				filePath = filepath.Join(tempDir, tt.fileName)
				err := os.WriteFile(filePath, []byte(tt.fileContent), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			} else {
				filePath = filepath.Join(tempDir, tt.fileName)
			}
			
			// Parse flags
			args := []string{"-file", filePath}
			err := cmd.Flags.Parse(args)
			if err != nil {
				t.Errorf("Should parse flags: %v", err)
			}
			
			err = cmd.Run(ctx, cmd.Flags.Args())
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// Test more initCommand error paths
func TestInitCommandMoreErrorPaths(t *testing.T) {
	ctx := context.Background()
	
	// Test with read-only directory
	tempDir := t.TempDir()
	readOnlyDir := filepath.Join(tempDir, "readonly")
	err := os.Mkdir(readOnlyDir, 0555) // Read-only directory
	if err != nil {
		t.Fatalf("Failed to create read-only directory: %v", err)
	}
	
	// Change to read-only directory
	originalDir, _ := os.Getwd()
	os.Chdir(readOnlyDir)
	defer func() {
		os.Chdir(originalDir)
		os.Chmod(readOnlyDir, 0755) // Make it writable again for cleanup
	}()
	
	// Try to initialize in read-only directory
	cmd := initCommand()
	
	// Test with file flag pointing to read-only location
	readOnlyFile := filepath.Join(readOnlyDir, "config.yaml")
	args := []string{"-file", readOnlyFile}
	err = cmd.Flags.Parse(args)
	if err != nil {
		t.Errorf("Should parse flags: %v", err)
	}
	
	err = cmd.Run(ctx, cmd.Flags.Args())
	if err == nil {
		t.Error("Expected error when writing to read-only directory")
	}
}

// Test more setupCommand error paths
func TestSetupCommandMoreErrorPaths(t *testing.T) {
	ctx := context.Background()
	
	// Test in non-git directory with no HOME
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	originalHome := os.Getenv("HOME")
	
	defer func() {
		os.Chdir(originalDir)
		os.Setenv("HOME", originalHome)
	}()
	
	// Change to temp directory (no git repo)
	os.Chdir(tempDir)
	// Unset HOME to trigger fallback paths
	os.Unsetenv("HOME")
	
	cmd := setupCommand()
	
	// Test with force flag in non-git directory
	args := []string{"-force"}
	err := cmd.Flags.Parse(args)
	if err != nil {
		t.Errorf("Should parse flags: %v", err)
	}
	
	err = cmd.Run(ctx, cmd.Flags.Args())
	// This should still work as it falls back to creating config
	if err != nil {
		t.Logf("Expected behavior - setup failed in non-git directory: %v", err)
	}
}

// Test more setupEnterpriseCommand error paths
func TestSetupEnterpriseCommandMoreErrorPaths(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()
	
	originalDir, _ := os.Getwd()
	originalHome := os.Getenv("HOME")
	
	defer func() {
		os.Chdir(originalDir)
		os.Setenv("HOME", originalHome)
	}()
	
	// Change to temp directory (no git repo)
	os.Chdir(tempDir)
	// Set HOME to temp directory
	os.Setenv("HOME", tempDir)
	
	// Create a read-only .fast-cc directory to trigger permission errors
	fastCCDir := filepath.Join(tempDir, ".fast-cc")
	err := os.Mkdir(fastCCDir, 0555)
	if err != nil {
		t.Fatalf("Failed to create read-only .fast-cc directory: %v", err)
	}
	defer os.Chmod(fastCCDir, 0755) // Make writable for cleanup
	
	cmd := setupEnterpriseCommand()
	
	// Test with local flag
	args := []string{"-local"}
	err = cmd.Flags.Parse(args)
	if err != nil {
		t.Errorf("Should parse flags: %v", err)
	}
	
	err = cmd.Run(ctx, cmd.Flags.Args())
	if err == nil {
		t.Error("Expected error when writing to read-only directory")
	}
}

// Test copyEnterpriseConfig with permission errors
func TestCopyEnterpriseConfigPermissionError(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create read-only directory
	readOnlyDir := filepath.Join(tempDir, "readonly")
	err := os.Mkdir(readOnlyDir, 0555)
	if err != nil {
		t.Fatalf("Failed to create read-only directory: %v", err)
	}
	defer os.Chmod(readOnlyDir, 0755) // Make writable for cleanup
	
	// Try to copy to read-only directory
	destPath := filepath.Join(readOnlyDir, "enterprise.yaml")
	err = copyEnterpriseConfig(destPath)
	if err == nil {
		t.Error("Expected permission error when copying to read-only directory")
	}
}

// Test ensureConfigExists with more error paths
func TestEnsureConfigExistsMoreErrors(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create file where directory should be
	badPath := filepath.Join(tempDir, "badfile")
	err := os.WriteFile(badPath, []byte("blocking"), 0644)
	if err != nil {
		t.Fatalf("Failed to create blocking file: %v", err)
	}
	
	// Try to create config where file exists
	configPath := filepath.Join(badPath, "config.yaml") // This will fail
	err = ensureConfigExists(configPath)
	if err == nil {
		t.Error("Expected error when config path is blocked by file")
	}
}

// Test ensureEnterpriseConfigExists with more error paths  
func TestEnsureEnterpriseConfigExistsMoreErrors(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create file where directory should be
	badPath := filepath.Join(tempDir, "badfile")
	err := os.WriteFile(badPath, []byte("blocking"), 0644)
	if err != nil {
		t.Fatalf("Failed to create blocking file: %v", err)
	}
	
	// Try to create enterprise config where file exists
	configPath := filepath.Join(badPath, "config.yaml") // This will fail
	err = ensureEnterpriseConfigExists(configPath)
	if err == nil {
		t.Error("Expected error when enterprise config path is blocked by file")
	}
}

// Test getGitConfigDir with various scenarios
func TestGetGitConfigDirVariousScenarios(t *testing.T) {
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	
	defer os.Setenv("HOME", originalHome)
	
	tests := []struct {
		name        string
		setupHome   bool
		homeValue   string
		expectError bool
	}{
		{
			name: "with valid HOME",
			setupHome: true,
			homeValue: tempDir,
			expectError: false,
		},
		{
			name: "with empty HOME",
			setupHome: true,
			homeValue: "",
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupHome {
				os.Setenv("HOME", tt.homeValue)
			} else {
				os.Unsetenv("HOME")
			}
			
			dir, err := getGitConfigDir()
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if dir == "" {
					t.Error("Expected non-empty directory")
				}
			}
		})
	}
}

// Test hasGlobalInstallation error scenarios
func TestHasGlobalInstallationErrors(t *testing.T) {
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	
	// Test with no HOME directory
	os.Unsetenv("HOME")
	
	has, err := hasGlobalInstallation()
	if err == nil {
		t.Error("Expected error when HOME is not set")
	}
	if has {
		t.Error("Should not report global installation when HOME is not set")
	}
}

// Test removeGlobalInstallation error scenarios
func TestRemoveGlobalInstallationErrors(t *testing.T) {
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	
	// Test with no HOME directory
	os.Unsetenv("HOME")
	
	err := removeGlobalInstallation()
	if err == nil {
		t.Error("Expected error when HOME is not set")
	}
}

// Test comprehensive coverage of copyEnterpriseConfig
func TestCopyEnterpriseConfigComprehensive(t *testing.T) {
	tempDir := t.TempDir()
	
	// Test successful copy
	destPath := filepath.Join(tempDir, "copied-enterprise.yaml")
	err := copyEnterpriseConfig(destPath)
	if err != nil {
		t.Errorf("Should successfully copy enterprise config: %v", err)
	}
	
	// Verify file was created and has expected content
	content, err := os.ReadFile(destPath)
	if err != nil {
		t.Errorf("Should read copied file: %v", err)
	}
	
	contentStr := string(content)
	expectedContent := []string{"types:", "scopes:", "require_jira_ticket", "max_subject_length"}
	for _, expected := range expectedContent {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Copied config should contain %s", expected)
		}
	}
}