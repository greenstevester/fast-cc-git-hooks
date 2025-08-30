package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	version   = "dev"
	buildTime = "unknown"
	commit    = "unknown"

	// Command line flags for gce.
	noVerify = flag.Bool("no-verify", false, "Skip pre-commit hooks")
	verbose  = flag.Bool("verbose", false, "Show detailed analysis")
	help     = flag.Bool("help", false, "Show help")
)

func main() {
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// Find the cc binary
	ccBinary, err := findCCBinary()
	if err != nil {
		log.Fatalf("Cannot find cc binary: %v", err)
	}

	// Build arguments for cc --execute
	args := []string{"--execute"}
	
	if *noVerify {
		args = append(args, "--no-verify")
	}
	
	if *verbose {
		args = append(args, "--verbose")
	}

	// Execute cc with --execute flag
	cmd := exec.Command(ccBinary, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		log.Fatalf("Error running cc: %v", err)
	}
}

// findCCBinary locates the cc binary relative to gce.
func findCCBinary() (string, error) {
	// Get the path of the current executable (gce)
	gceExe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("cannot find gce executable: %w", err)
	}

	// Look for cc in the same directory
	gceDir := filepath.Dir(gceExe)
	ccPath := filepath.Join(gceDir, "cc")
	
	// Check if cc exists in the same directory
	if _, err := os.Stat(ccPath); err == nil {
		return ccPath, nil
	}

	// Try cc with .exe extension (Windows)
	ccExePath := ccPath + ".exe"
	if _, err := os.Stat(ccExePath); err == nil {
		return ccExePath, nil
	}

	// Fallback: try to find cc in PATH
	ccInPath, err := exec.LookPath("cc")
	if err == nil {
		return ccInPath, nil
	}

	return "", fmt.Errorf("cc binary not found in %s or PATH", gceDir)
}

func showHelp() {
	fmt.Printf(`gce - Generate Commit & Execute v%s

A shortcut for 'cc --execute' that generates a conventional commit message and executes it.

USAGE:
    gce [options]

OPTIONS:
    --no-verify     Skip pre-commit hooks when committing
    --verbose       Show detailed analysis of changes
    --help          Show this help message

DESCRIPTION:
    gce is a convenient shortcut that combines commit message generation with 
    immediate execution. It analyzes your staged changes, generates an appropriate
    conventional commit message, and commits the changes automatically.

    This command is equivalent to running: cc --execute

EXAMPLES:
    gce                    # Generate and commit with default settings
    gce --no-verify        # Generate and commit, skipping pre-commit hooks
    gce --verbose          # Generate and commit with detailed analysis

NOTES:
    - Requires the 'cc' command to be available
    - Works best with staged changes (git add your files first)
    - Follows conventional commit format (feat:, fix:, docs:, etc.)

VERSION:
    Version: %s
    Built:   %s
    Commit:  %s

For more information about the underlying cc command, run: cc --help
`, version, version, buildTime, commit)
}