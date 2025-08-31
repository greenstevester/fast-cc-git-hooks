package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/greenstevester/fast-cc-git-hooks/internal/banner"
	"github.com/greenstevester/fast-cc-git-hooks/pkg/ccgen"
)

var (
	version   = "dev"
	buildTime = "unknown"
	commit    = "unknown"

	// Command line flags.
	noVerify = flag.Bool("no-verify", false, "Skip pre-commit hooks")
	execute  = flag.Bool("execute", false, "Execute the commit after generating message")
	copyCmd  = flag.Bool("copy", false, "Copy git commit command to clipboard")
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

	// Create generator with specified options
	generator := ccgen.New(ccgen.Options{
		NoVerify: *noVerify,
		Execute:  *execute,
		Copy:     *copyCmd,
		Verbose:  *verbose,
	})

	// Generate commit message
	result, err := generator.Generate()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Print result
	generator.PrintResult(result)
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
	fmt.Println("  --copy         Copy git commit command to clipboard")
	fmt.Println("  --no-verify    Skip pre-commit hooks when committing")
	fmt.Println("  --verbose      Show detailed analysis of changes")
	fmt.Println("  --help         Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  cc                    # Generate commit message only")
	fmt.Println("  cc --copy             # Copy git commit command to clipboard")
	fmt.Println("  cc --execute          # Generate and commit")
	fmt.Println("  cc --verbose          # Show detailed analysis")
	fmt.Println("  cc --execute --no-verify  # Commit without hooks")
	fmt.Println()
	fmt.Printf("Build info: %s (%s)\n", buildTime, commit)
}
