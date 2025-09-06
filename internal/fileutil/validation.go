// Package fileutil provides file utilities with security validations.
package fileutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/greenstevester/fast-cc-git-hooks/internal/errors"
)

const (
	// MaxFileSize is the maximum allowed file size (1MB)
	MaxFileSize = 1024 * 1024
	// MaxCommitFileSize is the maximum allowed commit message file size (64KB)
	MaxCommitFileSize = 64 * 1024
)

// ValidateFilePath validates a file path to prevent directory traversal attacks
func ValidateFilePath(path string) error {
	if path == "" {
		return errors.NewFileError("file path cannot be empty")
	}

	// Clean the path and check for directory traversal attempts
	cleaned := filepath.Clean(path)
	if strings.Contains(cleaned, "..") {
		return errors.NewFileError("file path contains directory traversal: %s", path)
	}

	// Check for null bytes which can be used to bypass security checks
	if strings.Contains(path, "\x00") {
		return errors.NewFileError("file path contains null bytes: %s", path)
	}

	return nil
}

// ValidateFileSize checks if a file exists and is within size limits
func ValidateFileSize(path string, maxSize int64) error {
	info, err := os.Stat(path)
	if err != nil {
		return errors.WrapFileError(fmt.Sprintf("cannot access file %s", path), err)
	}

	if info.Size() > maxSize {
		return errors.NewFileError("file %s is too large (%d bytes, max %d bytes)", path, info.Size(), maxSize)
	}

	return nil
}

// SafeReadFile reads a file with path validation and size checks
func SafeReadFile(path string, maxSize int64) ([]byte, error) {
	if err := ValidateFilePath(path); err != nil {
		return nil, err
	}

	if err := ValidateFileSize(path, maxSize); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path) // #nosec G304 - path is validated by ValidateFilePath before use
	if err != nil {
		return nil, errors.WrapFileError(fmt.Sprintf("reading file %s", path), err)
	}

	return data, nil
}

// SafeReadCommitFile is a specialized function for reading commit message files
func SafeReadCommitFile(path string) (string, error) {
	data, err := SafeReadFile(path, MaxCommitFileSize)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
