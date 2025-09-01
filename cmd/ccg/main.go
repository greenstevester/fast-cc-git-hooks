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

	// Command line flags.
	noVerify = flag.Bool("no-verify", false, "Skip pre-commit hooks")
	execute  = flag.Bool("execute", false, "Execute the commit after generating message")
	noCopy   = flag.Bool("no-copy", false, "Disable copying git commit command to clipboard")
	verbose  = flag.Bool("verbose", false, "Show detailed analysis")
	help     = flag.Bool("help", false, "Show help")
)

func main() {
	// Print banner with version, commit and build time information
	banner.PrintWithVersionAndBuildTime(version, commit, buildTime)

	flag.Parse()

	// Handle subcommands
	args := flag.Args()
	if len(args) > 0 {
		if err := handleSubcommand(args); err != nil {
			log.Fatalf("Error: %v", err)
		}
		return
	}

	if *help {
		showHelp()
		return
	}

	// Get current working directory for JIRA manager
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}

	// Create generator with specified options
	generator := ccgen.New(ccgen.Options{
		NoVerify:    *noVerify,
		Execute:     *execute,
		Copy:        !*noCopy, // Copy by default unless --no-copy is specified
		Verbose:     *verbose,
		JiraManager: jira.NewManager(cwd),
	})

	// Generate commit message
	result, err := generator.Generate()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Print result
	generator.PrintResult(result)
}

func handleSubcommand(args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	jiraManager := jira.NewManager(cwd)

	switch args[0] {
	case "set-jira":
		if len(args) != 2 {
			return fmt.Errorf("usage: ccg set-jira <JIRA-TICKET>\nExample: ccg set-jira CGC-1234")
		}
		ticketID := args[1]
		if err := jiraManager.SetJiraTicket(ticketID); err != nil {
			return err
		}
		fmt.Printf("✅ **JIRA ticket set:** `%s`\n", ticketID)
		fmt.Println("\nThis ticket will now be automatically included in commit messages.")
		return nil

	case "clear-jira":
		if err := jiraManager.ClearJiraTicket(); err != nil {
			return err
		}
		fmt.Println("✅ **JIRA ticket cleared**")
		fmt.Println("\nNo JIRA ticket will be included in commit messages.")
		return nil

	case "jira-status":
		return jiraManager.ShowJiraStatus()

	case "jira-history":
		return jiraManager.ListJiraHistory()

	default:
		return fmt.Errorf("unknown subcommand: %s\n\nAvailable subcommands:\n  set-jira <TICKET>   Set current JIRA ticket\n  clear-jira          Clear current JIRA ticket\n  jira-status         Show current JIRA ticket status\n  jira-history        Show JIRA ticket history", args[0])
	}
}

func showHelp() {
	fmt.Printf("ccg - Git Commit message generator v%s\n\n", version)
	fmt.Println("Analyzes staged changes and generates conventional commit messages.")
	fmt.Println("Automatically copies git commit command to clipboard for easy pasting.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  ccg [flags]                    # Generate commit message")
	fmt.Println("  ccg <subcommand> [args]        # JIRA ticket management")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --execute      Execute the commit after generating message")
	fmt.Println("  --no-copy      Disable copying git commit command to clipboard")
	fmt.Println("  --no-verify    Skip pre-commit hooks when committing")
	fmt.Println("  --verbose      Show detailed analysis of changes")
	fmt.Println("  --help         Show this help message")
	fmt.Println()
	fmt.Println("JIRA Commands:")
	fmt.Println("  set-jira <TICKET>     Set current JIRA ticket (e.g., CGC-1234)")
	fmt.Println("  clear-jira            Clear current JIRA ticket")
	fmt.Println("  jira-status           Show current JIRA ticket status")
	fmt.Println("  jira-history          Show JIRA ticket history")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  ccg                    # Generate and copy git commit command")
	fmt.Println("  ccg --execute          # Generate and commit immediately")
	fmt.Println("  ccg set-jira CGC-1234  # Set JIRA ticket for future commits")
	fmt.Println("  ccg jira-status        # Check current JIRA ticket")
	fmt.Println("  ccg clear-jira         # Remove JIRA ticket from commits")
	fmt.Println()
	fmt.Printf("Build info: %s (%s)\n", buildTime, commit)
}
