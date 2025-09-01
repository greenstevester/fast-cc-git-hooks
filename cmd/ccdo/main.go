package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/greenstevester/fast-cc-git-hooks/internal/banner"
	"github.com/greenstevester/fast-cc-git-hooks/pkg/ccgen"
	"github.com/greenstevester/fast-cc-git-hooks/pkg/jira"
)

var (
	version   = "dev"
	buildTime = "unknown"
	commit    = "unknown"

	// Command line flags for ccdo.
	noVerify = flag.Bool("no-verify", false, "Skip pre-commit hooks")
	verbose  = flag.Bool("verbose", false, "Show detailed analysis")
	help     = flag.Bool("help", false, "Show help")
)

func main() {
	// Print banner with version, commit and build time information
	banner.PrintWithVersionAndBuildTime(version, commit, buildTime)

	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// Get current working directory for JIRA manager
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}

	// Create generator with execute option enabled
	generator := ccgen.New(ccgen.Options{
		NoVerify:    *noVerify,
		Execute:     true, // ccdo always executes
		Copy:        false,
		Verbose:     *verbose,
		JiraManager: jira.NewManager(cwd),
	})

	// Generate commit message and execute
	result, err := generator.Generate()
	if err != nil {
		log.Fatalf("Failed to generate commit: %v", err)
	}

	// Print result (includes execution)
	generator.PrintResult(result)

	// Exit with appropriate code
	if !result.HasChanges {
		os.Exit(0)
	}
}

func showHelp() {
	fmt.Printf(`ccdo - Conventional Commit Do... v %s

auto generates a conventional commit message and commits it for you (it doesn't get any easier than that).'.

USAGE:
    ccdo [options]

OPTIONS:
    --no-verify     Skip pre-commit hooks when committing
    --verbose       Show detailed analysis of changes
    --help          Show this help message

DESCRIPTION:
    ccdo is a convenience tool, that combines commit message generation with 
    immediate git commit. It analyzes your staged changes, generates an appropriate
    conventional commit message, and commits the changes automatically.

    This command is equivalent to running: ccg --execute

EXAMPLES:
    ccdo                    # Generate and commit with default settings
    ccdo --no-verify        # Generate and commit, skipping pre-commit hooks
    ccdo --verbose          # Generate and commit with detailed analysis

NOTES:
    - Works best with staged changes (git add your files first)
    - Follows conventional commit format (feat:, fix:, docs:, etc.)
    - Now uses shared commit generation logic (no external ccg dependency)

VERSION:
    Version: %s
    Built:   %s
    Commit:  %s

For more information about the underlying ccg command, run: ccg --help
`, version, version, buildTime, commit)
}
