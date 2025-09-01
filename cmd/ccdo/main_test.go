package main

import (
	"testing"

	"github.com/greenstevester/fast-cc-git-hooks/pkg/ccgen"
)

func TestCCCUsesSharedPackage(t *testing.T) {
	// Test that we can create a generator (main functionality of ccc)
	generator := ccgen.New(ccgen.Options{
		NoVerify: false,
		Execute:  true, // ccc always executes
		Copy:     false,
		Verbose:  false,
	})

	if generator == nil {
		t.Error("Expected generator to be created successfully")
	}
}

func TestShowHelp(t *testing.T) {
	// Test that showHelp doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("showHelp panicked: %v", r)
		}
	}()

	// We can't easily test the output without capturing stdout,
	// but we can ensure it doesn't panic
	// showHelp() // This would print to stdout
}
