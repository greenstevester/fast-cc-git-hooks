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

	// Find the gc binary
	gcBinary, err := findGCBinary()
	if err != nil {
		log.Fatalf("Cannot find gc binary: %v", err)
	}

	// Build arguments for gc --execute
	args := []string{"--execute"}
	
	if *noVerify {
		args = append(args, "--no-verify")
	}
	
	if *verbose {
		args = append(args, "--verbose")
	}

	// Execute gc with --execute flag
	cmd := exec.Command(gcBinary, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		log.Fatalf("Error running gc: %v", err)
	}
}

// findGCBinary locates the gc binary relative to gce.
func findGCBinary() (string, error) {
	// Get the path of the current executable (gce)
	gceExe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("cannot find gce executable: %w", err)
	}

	// Look for gc in the same directory
	gceDir := filepath.Dir(gceExe)
	gcPath := filepath.Join(gceDir, "gc")
	
	// Check if gc exists in the same directory
	if _, err := os.Stat(gcPath); err == nil {
		return gcPath, nil
	}

	// Try gc with .exe extension (Windows)
	gcExePath := gcPath + ".exe"
	if _, err := os.Stat(gcExePath); err == nil {
		return gcExePath, nil
	}

	// Fallback: try to find gc in PATH
	gcInPath, err := exec.LookPath("gc")
	if err == nil {
		return gcInPath, nil
	}

	return "", fmt.Errorf("gc binary not found in %s or PATH", gceDir)
}

func showHelp() {
	fmt.Printf(`gce - Generate Commit & Execute v%s

A shortcut for 'gc --execute' that generates a conventional commit message and executes it.

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

    This command is equivalent to running: gc --execute

EXAMPLES:
    gce                    # Generate and commit with default settings
    gce --no-verify        # Generate and commit, skipping pre-commit hooks
    gce --verbose          # Generate and commit with detailed analysis

NOTES:
    - Requires the 'gc' command to be available
    - Works best with staged changes (git add your files first)
    - Follows conventional commit format (feat:, fix:, docs:, etc.)

VERSION:
    Version: %s
    Built:   %s
    Commit:  %s

For more information about the underlying gc command, run: gc --help
`, version, version, buildTime, commit)
}