// Package errors provides standardized error formatting and types.
package errors

import (
	"fmt"
)

// ErrorType represents different categories of errors
type ErrorType string

const (
	// ValidationError represents validation failures
	ValidationError ErrorType = "validation"
	// ConfigError represents configuration issues
	ConfigError ErrorType = "config"
	// FileError represents file operation issues
	FileError ErrorType = "file"
	// GitError represents git-related issues
	GitError ErrorType = "git"
	// NetworkError represents network-related issues
	NetworkError ErrorType = "network"
	// InternalError represents internal application issues
	InternalError ErrorType = "internal"
)

// StandardError represents a formatted error with context
type StandardError struct {
	Type    ErrorType
	Message string
	Cause   error
	Context map[string]interface{}
}

func (e *StandardError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *StandardError) Unwrap() error {
	return e.Cause
}

// New creates a new standardized error
func New(errorType ErrorType, message string) *StandardError {
	return &StandardError{
		Type:    errorType,
		Message: message,
		Context: make(map[string]interface{}),
	}
}

// Wrap wraps an existing error with context
func Wrap(errorType ErrorType, message string, cause error) *StandardError {
	return &StandardError{
		Type:    errorType,
		Message: message,
		Cause:   cause,
		Context: make(map[string]interface{}),
	}
}

// WithContext adds context information to an error
func (e *StandardError) WithContext(key string, value interface{}) *StandardError {
	e.Context[key] = value
	return e
}

// Convenience functions for common error types

// NewValidationError creates a validation error
func NewValidationError(message string, args ...interface{}) *StandardError {
	return New(ValidationError, fmt.Sprintf(message, args...))
}

// WrapValidationError wraps a validation error
func WrapValidationError(message string, cause error) *StandardError {
	return Wrap(ValidationError, message, cause)
}

// NewConfigError creates a config error
func NewConfigError(message string, args ...interface{}) *StandardError {
	return New(ConfigError, fmt.Sprintf(message, args...))
}

// WrapConfigError wraps a config error
func WrapConfigError(message string, cause error) *StandardError {
	return Wrap(ConfigError, message, cause)
}

// NewFileError creates a file error
func NewFileError(message string, args ...interface{}) *StandardError {
	return New(FileError, fmt.Sprintf(message, args...))
}

// WrapFileError wraps a file error
func WrapFileError(message string, cause error) *StandardError {
	return Wrap(FileError, message, cause)
}

// NewGitError creates a git error
func NewGitError(message string, args ...interface{}) *StandardError {
	return New(GitError, fmt.Sprintf(message, args...))
}

// WrapGitError wraps a git error
func WrapGitError(message string, cause error) *StandardError {
	return Wrap(GitError, message, cause)
}
