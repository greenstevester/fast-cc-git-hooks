// fast-cc-hooks is a fast conventional commits git hook system.
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/greenstevester/fast-cc-git-hooks/internal/config"
	"github.com/greenstevester/fast-cc-git-hooks/internal/hooks"
	"github.com/greenstevester/fast-cc-git-hooks/internal/validator"
)

const version = "1.0.0"

// Command represents a CLI command.
type Command struct {
	Run         func(ctx context.Context, args []string) error
	Flags       *flag.FlagSet
	Name        string
	Description string
}

var (
	// Global flags.
	verbose    bool
	configFile string

	// Command-specific flags..
	validateFile string
	forceInstall bool
	localInstall bool

	logger *slog.Logger
)

func main() {
	// Setup base logger.
	setupLogger(false)

	commands := map[string]*Command{
		"setup":      setupCommand(),
		"setup-ent":  setupEnterpriseCommand(),
		"remove":     removeCommand(),
		"validate":   validateCommand(),
		"init":       initCommand(),
		"version":    versionCommand(),
	}

	// Parse global flags.
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.StringVar(&configFile, "config", "", "path to config file")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "🚀 fast-cc-hooks - Make your commit messages awesome!\n\n")
		fmt.Fprintf(os.Stderr, "📋 Super Easy Setup (just 2 steps!):\n")
		fmt.Fprintf(os.Stderr, "   1️⃣  %s setup     ← Start here! This sets everything up\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "   2️⃣  git commit -m \"feat: your message\"  ← Write better commits!\n\n")

		fmt.Fprintf(os.Stderr, "✨ All Commands:\n")
		fmt.Fprintf(os.Stderr, "  %-10s %s\n", "setup", "🚀 Easy setup - install git hooks (use --local for current repo only)")
		fmt.Fprintf(os.Stderr, "  %-10s %s\n", "setup-ent", "🏢 Enterprise setup - with JIRA validation (use --local for current repo only)")
		fmt.Fprintf(os.Stderr, "  %-10s %s\n", "remove", "🗑️  Easy removal - uninstall git hooks")
		fmt.Fprintf(os.Stderr, "  %-10s %s\n", "validate", "🔍 Test a commit message")
		fmt.Fprintf(os.Stderr, "  %-10s %s\n", "init", "📝 Create a config file")
		fmt.Fprintf(os.Stderr, "  %-10s %s\n", "version", "ℹ️  Show version info")
		fmt.Fprintf(os.Stderr, "\n🤓 Advanced Commands:\n")

		fmt.Fprintf(os.Stderr, "\n🏁 Quick Start:\n")
		fmt.Fprintf(os.Stderr, "   %s setup\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n🔧 Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n💡 Need help? Use '%s <command> -h' for more details\n", os.Args[0])
	}

	// Need at least command name...
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	// Extract command...
	cmdName := os.Args[1]

	// Handle help for commands...
	if cmdName == "-h" || cmdName == "--help" || cmdName == "help" {
		flag.Usage()
		os.Exit(0)
	}

	cmd, exists := commands[cmdName]
	if !exists {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmdName)
		flag.Usage()
		os.Exit(1)
	}

	// Parse command flags.
	if err := cmd.Flags.Parse(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	// Update logger with verbose flag..
	setupLogger(verbose)

	// Create context with timeout...
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	// Run command...
	if err := cmd.Run(ctx, cmd.Flags.Args()); err != nil {
		cancel()
		logger.Error("command failed", "command", cmdName, "error", err)
		os.Exit(1)
	}
	cancel()
}

func setupLogger(verbose bool) {
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewTextHandler(os.Stderr, opts)
	logger = slog.New(handler)
	slog.SetDefault(logger)
}


func validateCommand() *Command {
	fs := flag.NewFlagSet("validate", flag.ExitOnError)
	fs.StringVar(&validateFile, "file", "", "validate commit message from file")

	return &Command{
		Name:        "validate",
		Description: "🔍 Test a commit message",
		Flags:       fs,
		Run: func(ctx context.Context, args []string) error {
			// Load configuration.
			cfg, err := config.Load(configFile)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			// Create validator.
			v, err := validator.New(cfg)
			if err != nil {
				return fmt.Errorf("creating validator: %w", err)
			}

			var result *validator.ValidationResult

			if validateFile != "" {
				// Validate from file.
				result, err = v.ValidateFile(ctx, validateFile)
				if err != nil {
					return fmt.Errorf("validating file: %w", err)
				}
			} else {
				// Validate from arguments or stdin.
				var message string
				if len(args) > 0 {
					message = strings.Join(args, " ")
				} else {
					// Read from stdin.
					buf := make([]byte, 0, 4096)
					for {
						n, err := os.Stdin.Read(buf[len(buf):cap(buf)])
						buf = buf[:len(buf)+n]
						if err != nil {
							break
						}
						if len(buf) == cap(buf) {
							// Grow buffer.
							newBuf := make([]byte, len(buf), cap(buf)*2)
							copy(newBuf, buf)
							buf = newBuf
						}
					}
					message = string(buf)
				}

				if message == "" {
					return fmt.Errorf("no commit message provided")
				}

				result = v.Validate(ctx, message)
			}

			if !result.Valid {
				fmt.Fprintf(os.Stderr, "❌ Commit message validation failed:\n")
				for _, err := range result.Errors {
					fmt.Fprintf(os.Stderr, "  • %v\n", err)
				}
				return fmt.Errorf("validation failed")
			}

			fmt.Println("✅ Commit message is valid")
			return nil
		},
	}
}

func initCommand() *Command {
	fs := flag.NewFlagSet("init", flag.ExitOnError)

	return &Command{
		Name:        "init",
		Description: "📝 Create a config file",
		Flags:       fs,
		Run: func(_ context.Context, _ []string) error {
			path := configFile
			if path == "" {
				// Use default path in home directory
				if defaultPath, err := config.GetDefaultConfigPath(); err == nil {
					path = defaultPath
				} else {
					path = config.DefaultConfigFile
				}
			}

			// Check if file exists.
			if _, err := os.Stat(path); err == nil {
				return fmt.Errorf("config file already exists: %s", path)
			}

			// Create default config..
			cfg := config.Default()

			// Save to file..
			if err := cfg.Save(path); err != nil {
				return fmt.Errorf("saving config: %w", err)
			}

			logger.Info("created configuration file", "path", path)
			fmt.Printf("✅ Created configuration file: %s\n", path)
			fmt.Println("\nDefault configuration includes:")
			fmt.Printf("  • Commit types: %s\n", strings.Join(cfg.Types, ", "))
			fmt.Printf("  • Max subject length: %d\n", cfg.MaxSubjectLength)
			fmt.Printf("  • Scope required: %v\n", cfg.ScopeRequired)
			fmt.Printf("  • Breaking changes allowed: %v\n", cfg.AllowBreakingChanges)
			fmt.Println("\nEdit the file to customize your rules.")

			return nil
		},
	}
}

func versionCommand() *Command {
	fs := flag.NewFlagSet("version", flag.ExitOnError)

	return &Command{
		Name:        "version",
		Description: "ℹ️  Show version info",
		Flags:       fs,
		Run: func(_ context.Context, _ []string) error {
			fmt.Printf("fast-cc-hooks version %s\n", version)
			fmt.Printf("Go version: %s\n", runtime.Version())
			fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
			return nil
		},
	}
}

func setupCommand() *Command {
	fs := flag.NewFlagSet("setup", flag.ExitOnError)
	fs.BoolVar(&forceInstall, "force", false, "force installation, overwriting existing hooks")
	fs.BoolVar(&localInstall, "local", false, "install only for current repository (default: install globally)")

	return &Command{
		Name:        "setup",
		Description: "🚀 Easy setup - install git hooks everywhere!",
		Flags:       fs,
		Run: func(ctx context.Context, _ []string) error {
			fmt.Println("🚀 Setting up fast-cc-hooks...")
			fmt.Println("   This will help you write better commit messages!")
			fmt.Println("")

			// Step 1: Check/create configuration
			configPath, configCreated, configErr := ensureConfigExists()
			if configErr != nil {
				fmt.Printf("⚠️  Warning: Could not create config: %v\n", configErr)
				fmt.Println("   Hooks will use default settings.")
			} else if configCreated {
				fmt.Printf("📝 Created default configuration: %s\n", configPath)
			} else {
				fmt.Printf("📝 Using existing configuration: %s\n", configPath)
			}
			fmt.Println("")

			// Step 2: Install hooks
			var err error
			if localInstall {
				fmt.Println("📁 Installing hooks for this repository only...")
				opts := hooks.Options{
					Logger:       logger,
					ForceInstall: forceInstall,
				}

				installer, instErr := hooks.New(opts)
				if instErr != nil {
					return fmt.Errorf("creating installer: %w", instErr)
				}

				err = installer.Install(ctx)
			} else {
				fmt.Println("🌍 Installing hooks globally (for all your repositories)...")
				err = hooks.GlobalInstall(ctx, logger)
			}

			if err != nil {
				fmt.Println("❌ Setup failed:", err)
				return err
			}

			fmt.Println("")
			fmt.Println("✅ All done! Your commit messages will now be checked automatically!")
			if configPath != "" {
				fmt.Printf("⚙️  Configuration stored at: %s\n", configPath)
				fmt.Println("   Edit this file to customize commit rules.")
			}
			fmt.Println("💡 Try making a commit like: git commit -m \"feat: add awesome feature\"")
			return nil
		},
	}
}

func setupEnterpriseCommand() *Command {
	fs := flag.NewFlagSet("setup-ent", flag.ExitOnError)
	fs.BoolVar(&forceInstall, "force", false, "force installation, overwriting existing hooks")
	fs.BoolVar(&localInstall, "local", false, "install only for current repository (default: install globally)")

	return &Command{
		Name:        "setup-ent",
		Description: "🏢 Enterprise setup - with JIRA ticket validation!",
		Flags:       fs,
		Run: func(ctx context.Context, _ []string) error {
			fmt.Println("🏢 Setting up fast-cc-hooks for Enterprise...")
			fmt.Println("   This includes JIRA ticket validation and enterprise-ready rules!")
			fmt.Println("")

			// Step 1: Check/create enterprise configuration
			configPath, configCreated, configErr := ensureEnterpriseConfigExists()
			if configErr != nil {
				fmt.Printf("⚠️  Warning: Could not create enterprise config: %v\n", configErr)
				fmt.Println("   Hooks will use default settings.")
			} else if configCreated {
				fmt.Printf("📝 Created enterprise configuration: %s\n", configPath)
				fmt.Println("   ✅ JIRA ticket validation enabled")
				fmt.Println("   ✅ Enterprise scopes configured")
				fmt.Println("   ✅ Advanced validation rules ready")
			} else {
				fmt.Printf("📝 Using existing configuration: %s\n", configPath)
			}
			fmt.Println("")

			// Step 2: Install hooks
			var err error
			if localInstall {
				fmt.Println("📁 Installing hooks for this repository only...")
				opts := hooks.Options{
					Logger:       logger,
					ForceInstall: forceInstall,
				}

				installer, instErr := hooks.New(opts)
				if instErr != nil {
					return fmt.Errorf("creating installer: %w", instErr)
				}

				err = installer.Install(ctx)
			} else {
				fmt.Println("🌍 Installing hooks globally (for all your repositories)...")
				err = hooks.GlobalInstall(ctx, logger)
			}

			if err != nil {
				fmt.Println("❌ Setup failed:", err)
				return err
			}

			fmt.Println("")
			fmt.Println("✅ Enterprise setup complete! Your commit messages will be validated with:")
			fmt.Println("   🎫 JIRA ticket references (required)")
			fmt.Println("   📋 Enterprise scopes (api, web, cli, db, auth, core, etc.)")
			fmt.Println("   🔧 Advanced validation rules")
			if configPath != "" {
				fmt.Printf("⚙️  Configuration stored at: %s\n", configPath)
				fmt.Println("   Edit this file to customize enterprise rules.")
			}
			fmt.Println("💡 Try: git commit -m \"feat(api): PROJ-123 Add user authentication\"")
			return nil
		},
	}
}

// ensureEnterpriseConfigExists checks for existing config or creates enterprise config.
// Returns (configPath, wasCreated, error).
func ensureEnterpriseConfigExists() (string, bool, error) {
	// First check if there's already a config file specified
	if configFile != "" {
		if _, err := os.Stat(configFile); err == nil {
			return configFile, false, nil
		}
		return "", false, fmt.Errorf("specified config file not found: %s", configFile)
	}

	// Get the default config path in home directory
	defaultPath, err := config.GetDefaultConfigPath()
	if err != nil {
		return "", false, fmt.Errorf("cannot determine config path: %w", err)
	}

	// Check if any config already exists in home directory
	if _, err := os.Stat(defaultPath); err == nil {
		return defaultPath, false, nil
	}

	// Check for old filename in home directory for backward compatibility
	oldPath := filepath.Join(filepath.Dir(defaultPath), ".fast-cc-hooks.yaml")
	if _, err := os.Stat(oldPath); err == nil {
		return oldPath, false, nil
	}

	// Check if config exists in current directory (new filename first)
	if _, err := os.Stat(config.DefaultConfigFile); err == nil {
		return config.DefaultConfigFile, false, nil
	}

	// Check for old filename in current directory
	if _, err := os.Stat(".fast-cc-hooks.yaml"); err == nil {
		return ".fast-cc-hooks.yaml", false, nil
	}

	// Create enterprise config in home directory
	if err := copyEnterpriseConfig(defaultPath); err != nil {
		return "", false, fmt.Errorf("creating enterprise config: %w", err)
	}

	return defaultPath, true, nil
}

// copyEnterpriseConfig copies the enterprise config template to the specified path.
func copyEnterpriseConfig(destPath string) error {
	// Get the path to the enterprise config template
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("finding executable: %w", err)
	}
	
	// Look for enterprise config relative to executable
	exeDir := filepath.Dir(executable)
	templatePath := filepath.Join(exeDir, "example-configs", "fast-cc-hooks.enterprise.yaml")
	
	// If not found, try relative to current directory (development scenario)
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		templatePath = filepath.Join("example-configs", "fast-cc-hooks.enterprise.yaml")
	}

	// Read the enterprise config template
	templateData, err := os.ReadFile(templatePath)
	if err != nil {
		// If we can't find the template, create a basic enterprise config
		return createBasicEnterpriseConfig(destPath)
	}

	// Ensure the destination directory exists
	if err := os.MkdirAll(filepath.Dir(destPath), 0o750); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	// Write the enterprise config
	if err := os.WriteFile(destPath, templateData, 0o644); err != nil {
		return fmt.Errorf("writing enterprise config: %w", err)
	}

	return nil
}

// createBasicEnterpriseConfig creates a basic enterprise config if template is not found.
func createBasicEnterpriseConfig(destPath string) error {
	enterpriseConfig := `# fast-cc-hooks enterprise configuration

# Allowed commit types
types:
  - feat
  - fix
  - docs
  - style
  - refactor
  - test
  - chore
  - perf
  - ci
  - build
  - revert

# Enterprise scopes
scopes:
  - api
  - web
  - cli
  - db
  - auth
  - core
  - mw
  - net
  - sec
  - iam
  - app

# Scope is not required by default
scope_required: false

# Maximum length of the subject line
max_subject_length: 72

# Allow breaking changes
allow_breaking_changes: true

# Require JIRA ticket references in commits
require_jira_ticket: true

# No general ticket reference requirement
require_ticket_ref: false

# Custom rules (empty by default)
custom_rules: []
`

	// Ensure the destination directory exists
	if err := os.MkdirAll(filepath.Dir(destPath), 0o750); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	// Write the basic enterprise config
	if err := os.WriteFile(destPath, []byte(enterpriseConfig), 0o644); err != nil {
		return fmt.Errorf("writing basic enterprise config: %w", err)
	}

	return nil
}

// ensureConfigExists checks for existing config or creates a default one.
// Returns (configPath, wasCreated, error).
func ensureConfigExists() (string, bool, error) {
	// First check if there's already a config file specified
	if configFile != "" {
		if _, err := os.Stat(configFile); err == nil {
			return configFile, false, nil
		}
		return "", false, fmt.Errorf("specified config file not found: %s", configFile)
	}

	// Check for config in the default home directory location
	defaultPath, err := config.GetDefaultConfigPath()
	if err != nil {
		// Fallback to current directory (new filename first)
		if _, err := os.Stat(config.DefaultConfigFile); err == nil {
			return config.DefaultConfigFile, false, nil
		}
		// Check for old filename in current directory
		if _, err := os.Stat(".fast-cc-hooks.yaml"); err == nil {
			return ".fast-cc-hooks.yaml", false, nil
		}
		return "", false, fmt.Errorf("cannot determine config path: %w", err)
	}

	// Check if config already exists in home directory (new filename)
	if _, err := os.Stat(defaultPath); err == nil {
		return defaultPath, false, nil
	}

	// Check for old filename in home directory for backward compatibility
	oldPath := filepath.Join(filepath.Dir(defaultPath), ".fast-cc-hooks.yaml")
	if _, err := os.Stat(oldPath); err == nil {
		return oldPath, false, nil
	}

	// Check if config exists in current directory (new filename first)
	if _, err := os.Stat(config.DefaultConfigFile); err == nil {
		return config.DefaultConfigFile, false, nil
	}

	// Check for old filename in current directory
	if _, err := os.Stat(".fast-cc-hooks.yaml"); err == nil {
		return ".fast-cc-hooks.yaml", false, nil
	}

	// Create default config in home directory with new filename
	cfg := config.Default()
	if err := cfg.Save(defaultPath); err != nil {
		return "", false, fmt.Errorf("creating default config: %w", err)
	}

	return defaultPath, true, nil
}

func removeCommand() *Command {
	fs := flag.NewFlagSet("remove", flag.ExitOnError)

	return &Command{
		Name:        "remove",
		Description: "🗑️  Easy removal - uninstall git hooks",
		Flags:       fs,
		Run: func(ctx context.Context, _ []string) error {
			fmt.Println("🗑️  Removing fast-cc-hooks...")
			fmt.Println("   (Don't worry, your code stays safe!)")
			fmt.Println("")

			opts := hooks.Options{
				Logger: logger,
			}

			installer, err := hooks.New(opts)
			if err != nil {
				return fmt.Errorf("creating installer: %w", err)
			}

			err = installer.Uninstall(ctx)
			if err != nil {
				fmt.Println("❌ Removal failed:", err)
				return err
			}

			fmt.Println("")
			fmt.Println("✅ All removed! fast-cc-hooks is no longer checking your commits")
			fmt.Println("💭 Thanks for using fast-cc-hooks!")
			return nil
		},
	}
}
