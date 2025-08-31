// Package banner provides terminal-aware banner printing
package banner

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// Print displays the banner with appropriate formatting for the terminal
func Print() {
	if UseASCII() {
		// Use ASCII art heart for better compatibility
		const a = ">> Made with <3 for Boo"
		fmt.Println(a)
	} else {
		// Use emoji for terminals that support it
		fmt.Println(">> Made with ❤️ for Boo")
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
	if UseASCII() {
		return ">> Made with <3 for Boo"
	}
	return ">> Made with ❤️ for Boo"
}
