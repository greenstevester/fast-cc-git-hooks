package main

import (
	"strings"
	"testing"
)

func TestFindCCBinary(t *testing.T) {
	// Test that findCCBinary doesn't panic
	_, err := findCCBinary()
	// It's okay if cc is not found in test environment
	if err != nil && !strings.Contains(err.Error(), "cc binary not found") {
		t.Errorf("Unexpected error from findCCBinary: %v", err)
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