// fast-cc-hooks is a fast conventional commits git hook system.
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/stevengreensill/fast-cc-git-hooks/internal/config"
	"github.com/stevengreensill/fast-cc-git-hooks/internal/hooks"
	"github.com/stevengreensill/fast-cc-git-hooks/internal/validator"
)

const version = "1.0.0"

// Command represents a CLI command
type Command struct {
	Name        string
	Description string
	Run         func(ctx context.Context, args []string) error
	Flags       *flag.FlagSet
}

var (
	// Global flags
	verbose    bool
	configFile string
	
	// Command-specific flags
	validateFile    string
	forceInstall    bool
	localInstall    bool
	
	logger *slog.Logger
)

func main() {
	// Setup base logger
	setupLogger(false)
	
	commands := map[string]*Command{
		"install":   installCommand(),
		"setup":     setupCommand(),
		"uninstall": uninstallCommand(),
		"remove":    removeCommand(),
		"validate":  validateCommand(),
		"init":      initCommand(),
		"version":   versionCommand(),
	}

	// Parse global flags
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.StringVar(&configFile, "config", "", "path to config file")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "🚀 fast-cc-hooks - Make your commit messages awesome!\n\n")
		fmt.Fprintf(os.Stderr, "📋 Super Easy Setup (just 2 steps!):\n")
		fmt.Fprintf(os.Stderr, "   1️⃣  %s setup     ← Start here! This sets everything up\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "   2️⃣  git commit -m \"feat: your message\"  ← Write better commits!\n\n")
		
		fmt.Fprintf(os.Stderr, "✨ All Commands:\n")
		fmt.Fprintf(os.Stderr, "  %-10s %s\n", "setup", "🚀 Easy setup - install git hooks everywhere!")
		fmt.Fprintf(os.Stderr, "  %-10s %s\n", "remove", "🗑️  Easy removal - uninstall git hooks")
		fmt.Fprintf(os.Stderr, "  %-10s %s\n", "validate", "🔍 Test a commit message")
		fmt.Fprintf(os.Stderr, "  %-10s %s\n", "init", "📝 Create a config file")
		fmt.Fprintf(os.Stderr, "  %-10s %s\n", "version", "ℹ️  Show version info")
		fmt.Fprintf(os.Stderr, "\n🤓 Advanced Commands:\n")
		fmt.Fprintf(os.Stderr, "  %-10s %s\n", "install", "Install git hooks globally for all repositories")
		fmt.Fprintf(os.Stderr, "  %-10s %s\n", "uninstall", "Remove git hooks from current repository")
		
		fmt.Fprintf(os.Stderr, "\n🏁 Quick Start:\n")
		fmt.Fprintf(os.Stderr, "   %s setup\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n🔧 Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n💡 Need help? Use '%s <command> -h' for more details\n", os.Args[0])
	}
	
	// Need at least command name
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	// Extract command
	cmdName := os.Args[1]
	
	// Handle help for commands
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

	// Parse command flags
	if err := cmd.Flags.Parse(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	// Update logger with verbose flag
	setupLogger(verbose)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Run command
	if err := cmd.Run(ctx, cmd.Flags.Args()); err != nil {
		logger.Error("command failed", "command", cmdName, "error", err)
		os.Exit(1)
	}
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

func installCommand() *Command {
	fs := flag.NewFlagSet("install", flag.ExitOnError)
	fs.BoolVar(&forceInstall, "force", false, "force installation, overwriting existing hooks")
	fs.BoolVar(&localInstall, "local", false, "install only for current repository (default: install globally)")
	
	return &Command{
		Name:        "install",
		Description: "Install git hooks globally for all repositories",
		Flags:       fs,
		Run: func(ctx context.Context, args []string) error {
			if localInstall {
				opts := hooks.Options{
					Logger:       logger,
					ForceInstall: forceInstall,
				}
				
				installer, err := hooks.New(opts)
				if err != nil {
					return fmt.Errorf("creating installer: %w", err)
				}
				
				return installer.Install(ctx)
			}
			
			return hooks.GlobalInstall(ctx, logger)
		},
	}
}

func uninstallCommand() *Command {
	fs := flag.NewFlagSet("uninstall", flag.ExitOnError)
	
	return &Command{
		Name:        "uninstall",
		Description: "Remove git hooks from current repository",
		Flags:       fs,
		Run: func(ctx context.Context, args []string) error {
			opts := hooks.Options{
				Logger: logger,
			}
			
			installer, err := hooks.New(opts)
			if err != nil {
				return fmt.Errorf("creating installer: %w", err)
			}
			
			return installer.Uninstall(ctx)
		},
	}
}

func validateCommand() *Command {
	fs := flag.NewFlagSet("validate", flag.ExitOnError)
	fs.StringVar(&validateFile, "file", "", "validate commit message from file")
	
	return &Command{
		Name:        "validate",
		Description: "🔍 Test a commit message",
		Flags:       fs,
		Run: func(ctx context.Context, args []string) error {
			// Load configuration
			cfg, err := config.Load(configFile)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}
			
			// Create validator
			v, err := validator.New(cfg)
			if err != nil {
				return fmt.Errorf("creating validator: %w", err)
			}
			
			var result *validator.ValidationResult
			
			if validateFile != "" {
				// Validate from file
				result, err = v.ValidateFile(ctx, validateFile)
				if err != nil {
					return fmt.Errorf("validating file: %w", err)
				}
			} else {
				// Validate from arguments or stdin
				var message string
				if len(args) > 0 {
					message = strings.Join(args, " ")
				} else {
					// Read from stdin
					buf := make([]byte, 0, 4096)
					for {
						n, err := os.Stdin.Read(buf[len(buf):cap(buf)])
						buf = buf[:len(buf)+n]
						if err != nil {
							break
						}
						if len(buf) == cap(buf) {
							// Grow buffer
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
		Run: func(ctx context.Context, args []string) error {
			path := configFile
			if path == "" {
				path = config.DefaultConfigFile
			}
			
			// Check if file exists
			if _, err := os.Stat(path); err == nil {
				return fmt.Errorf("config file already exists: %s", path)
			}
			
			// Create default config
			cfg := config.Default()
			
			// Save to file
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
		Run: func(ctx context.Context, args []string) error {
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
		Run: func(ctx context.Context, args []string) error {
			fmt.Println("🚀 Setting up fast-cc-hooks...")
			fmt.Println("   This will help you write better commit messages!")
			fmt.Println("")
			
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
			fmt.Println("💡 Try making a commit like: git commit -m \"feat: add awesome feature\"")
			return nil
		},
	}
}

func removeCommand() *Command {
	fs := flag.NewFlagSet("remove", flag.ExitOnError)
	
	return &Command{
		Name:        "remove",
		Description: "🗑️  Easy removal - uninstall git hooks",
		Flags:       fs,
		Run: func(ctx context.Context, args []string) error {
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