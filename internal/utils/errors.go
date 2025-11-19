package utils

import (
	"errors"
	"fmt"
)

// Common error types
var (
	ErrClusterNotFound      = errors.New("cluster not found")
	ErrClusterAlreadyExists = errors.New("cluster already exists")
	ErrInvalidConfig        = errors.New("invalid configuration")
	ErrConnectionFailed     = errors.New("connection failed")
	ErrAdapterNotFound      = errors.New("adapter not found")
	ErrOperationTimeout     = errors.New("operation timeout")
	ErrServiceUnavailable   = errors.New("service unavailable")
	ErrInvalidClusterID     = errors.New("invalid cluster ID")
	ErrInvalidOperation     = errors.New("invalid operation")
	ErrUnauthorized         = errors.New("unauthorized")
)

// ThroomError represents a custom error with additional context
type ThroomError struct {
	Code    string
	Message string
	Err     error
	Context map[string]interface{}
}

// Error implements the error interface
func (e *ThroomError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the wrapped error
func (e *ThroomError) Unwrap() error {
	return e.Err
}

// NewError creates a new ThroomError
func NewError(code, message string, err error) *ThroomError {
	return &ThroomError{
		Code:    code,
		Message: message,
		Err:     err,
		Context: make(map[string]interface{}),
	}
}

// WithContext adds context to the error
func (e *ThroomError) WithContext(key string, value interface{}) *ThroomError {
	e.Context[key] = value
	return e
}

// IsRetryable checks if an error is retryable
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// Check for specific error types
	if errors.Is(err, ErrConnectionFailed) ||
		errors.Is(err, ErrOperationTimeout) ||
		errors.Is(err, ErrServiceUnavailable) {
		return true
	}

	return false
}
