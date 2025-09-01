// Package banner provides terminal-aware banner printing
package banner

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

// Print displays the banner with appropriate formatting for the terminal
func Print() {
	PrintWithVersion("dev", "unknown")
}

// PrintWithVersion displays the banner with version and commit information
func PrintWithVersion(version, commit string) {
	PrintWithVersionAndBuildTime(version, commit, "")
}

// PrintWithVersionAndBuildTime displays the banner with version, commit and build time information
func PrintWithVersionAndBuildTime(version, commit, buildTime string) {
	var versionSuffix string
	
	// Format buildTime to dd.mm.yyyy if provided
	var formattedBuildTime string
	if buildTime != "" && buildTime != "unknown" {
		// Try to parse various date formats and convert to dd.mm.yyyy
		if parsedTime, err := time.Parse(time.RFC3339, buildTime); err == nil {
			formattedBuildTime = parsedTime.Format("02.01.2006")
		} else if parsedTime, err := time.Parse("2006-01-02T15:04:05Z", buildTime); err == nil {
			formattedBuildTime = parsedTime.Format("02.01.2006")
		} else if parsedTime, err := time.Parse("2006-01-02", buildTime); err == nil {
			formattedBuildTime = parsedTime.Format("02.01.2006")
		} else {
			// If parsing fails, use the buildTime as is
			formattedBuildTime = buildTime
		}
	}
	
	if version != "dev" && version != "unknown" && version != "" {
		if commit != "unknown" && commit != "" && len(commit) >= 7 {
			// Use short commit hash (first 7 characters)
			if formattedBuildTime != "" {
				versionSuffix = fmt.Sprintf(" / version %s (%s) built %s", version, commit[:7], formattedBuildTime)
			} else {
				versionSuffix = fmt.Sprintf(" / version %s (%s)", version, commit[:7])
			}
		} else {
			if formattedBuildTime != "" {
				versionSuffix = fmt.Sprintf(" / version %s built %s", version, formattedBuildTime)
			} else {
				versionSuffix = fmt.Sprintf(" / version %s", version)
			}
		}
	} else if commit != "unknown" && commit != "" && len(commit) >= 7 {
		// Just show commit if version is not available
		if formattedBuildTime != "" {
			versionSuffix = fmt.Sprintf(" / %s built %s", commit[:7], formattedBuildTime)
		} else {
			versionSuffix = fmt.Sprintf(" / %s", commit[:7])
		}
	} else if formattedBuildTime != "" {
		versionSuffix = fmt.Sprintf(" / built %s", formattedBuildTime)
	}

	if UseASCII() {
		// Use ASCII art heart for better compatibility
		fmt.Printf(">>> fast-cc gen / Made with <3 for Boo%s\n", versionSuffix)
	} else {
		// Use emoji for terminals that support it
		fmt.Printf(">>> fast-cc gen / Made with ❤️  for Boo%s\n", versionSuffix)
	}
}

// UseASCII determines if ASCII characters should be used instead of emojis
func UseASCII() bool {
	// Check various environment variables that indicate terminal type
	term := os.Getenv("TERM")
	msystem := os.Getenv("MSYSTEM") // MinGW/MSYS2
	termProgram := os.Getenv("TERM_PROGRAM")
	wtSession := os.Getenv("WT_SESSION") // Windows Terminal

	// Check if we're in Git Bash, MinGW, or similar
	if msystem != "" {
		return true // MinGW/MSYS2/Git Bash
	}

	// Check for specific terminal programs that don't handle emojis well
	if strings.Contains(strings.ToLower(term), "mingw") ||
		strings.Contains(strings.ToLower(term), "cygwin") ||
		strings.Contains(strings.ToLower(term), "msys") {
		return true
	}

	// Windows Command Prompt doesn't support emojis well
	if runtime.GOOS == "windows" {
		// Windows Terminal (newer) supports emojis
		if wtSession != "" {
			return false
		}
		// Check if we're in VSCode terminal which supports emojis
		if termProgram == "vscode" {
			return false
		}
		// Default to ASCII on Windows unless we know the terminal supports emojis
		return true
	}

	// Most modern Linux/Mac terminals support emojis
	return false
}

// GetBannerText returns the appropriate banner text based on terminal capabilities
func GetBannerText() string {
	return GetBannerTextWithVersion("dev", "unknown")
}

// GetBannerTextWithVersion returns banner text with version and commit information
func GetBannerTextWithVersion(version, commit string) string {
	var versionSuffix string
	if version != "dev" && version != "unknown" && version != "" {
		if commit != "unknown" && commit != "" && len(commit) >= 7 {
			// Use short commit hash (first 7 characters)
			versionSuffix = fmt.Sprintf(" / version %s (%s)", version, commit[:7])
		} else {
			versionSuffix = fmt.Sprintf(" / version %s", version)
		}
	} else if commit != "unknown" && commit != "" && len(commit) >= 7 {
		// Just show commit if version is not available
		versionSuffix = fmt.Sprintf(" / %s", commit[:7])
	}

	if UseASCII() {
		return fmt.Sprintf(">>> fast-cc gen / Made with <3 for Boo%s", versionSuffix)
	}
	return fmt.Sprintf(">>> fast-cc gen / Made with ❤️  for Boo%s", versionSuffix)
}
