package errors

import (
	"fmt"
	"time"
)

// ErrorType represents different categories of errors that can occur
type ErrorType int

const (
	// JSONParseError indicates an error during JSON parsing
	JSONParseError ErrorType = iota

	// PathQueryError indicates an error in path query processing
	PathQueryError

	// TransformationError indicates an error during transformation
	TransformationError

	// ValidationError indicates a validation failure
	ValidationError

	// ConfigurationError indicates a configuration-related error
	ConfigurationError

	// IOError indicates an input/output error
	IOError

	// TimeoutError indicates an operation timeout
	TimeoutError

	// MemoryError indicates a memory-related error
	MemoryError

	// PluginError indicates a plugin-related error
	PluginError
)

// String returns a string representation of the ErrorType
func (et ErrorType) String() string {
	switch et {
	case JSONParseError:
		return "JSON Parse Error"
	case PathQueryError:
		return "Path Query Error"
	case TransformationError:
		return "Transformation Error"
	case ValidationError:
		return "Validation Error"
	case ConfigurationError:
		return "Configuration Error"
	case IOError:
		return "I/O Error"
	case TimeoutError:
		return "Timeout Error"
	case MemoryError:
		return "Memory Error"
	case PluginError:
		return "Plugin Error"
	default:
		return "Unknown Error"
	}
}

// FJError represents a comprehensive error type for the FJ library
// with detailed context, error chaining, and debugging information.
type FJError struct {
	// Type categorizes the error
	Type ErrorType `json:"type"`

	// Message provides a human-readable error description
	Message string `json:"message"`

	// Code provides a machine-readable error code
	Code string `json:"code"`

	// Context provides additional contextual information
	Context map[string]interface{} `json:"context,omitempty"`

	// Timestamp records when the error occurred
	Timestamp time.Time `json:"timestamp"`

	// Source indicates where the error originated (function/method name)
	Source string `json:"source,omitempty"`

	// Input contains the input data that caused the error (truncated if too long)
	Input string `json:"input,omitempty"`

	// Position indicates the character position where the error occurred (for parsing errors)
	Position int `json:"position,omitempty"`

	// Path indicates the JSON path where the error occurred
	Path string `json:"path,omitempty"`

	// Cause contains the underlying error that caused this error
	Cause error `json:"cause,omitempty"`

	// Stack contains the call stack information
	Stack []string `json:"stack,omitempty"`

	// Recoverable indicates whether the error is recoverable
	Recoverable bool `json:"recoverable"`
}

// Error implements the error interface
func (e *FJError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("%s at path '%s': %s", e.Type.String(), e.Path, e.Message)
	}
	return fmt.Sprintf("%s: %s", e.Type.String(), e.Message)
}

// Unwrap returns the underlying cause error for error unwrapping
func (e *FJError) Unwrap() error {
	return e.Cause
}

// IsType checks if the error is of a specific type
func (e *FJError) IsType(errorType ErrorType) bool {
	return e.Type == errorType
}

// IsRecoverable returns whether the error is recoverable
func (e *FJError) IsRecoverable() bool {
	return e.Recoverable
}

// WithContext adds contextual information to the error
func (e *FJError) WithContext(key string, value interface{}) *FJError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithPath sets the JSON path where the error occurred
func (e *FJError) WithPath(path string) *FJError {
	e.Path = path
	return e
}

// WithPosition sets the character position where the error occurred
func (e *FJError) WithPosition(position int) *FJError {
	e.Position = position
	return e
}

// WithInput sets the input data that caused the error (truncated if necessary)
func (e *FJError) WithInput(input string) *FJError {
	maxInputLength := 200
	if len(input) > maxInputLength {
		e.Input = input[:maxInputLength] + "..."
	} else {
		e.Input = input
	}
	return e
}

// WithCause sets the underlying cause error
func (e *FJError) WithCause(cause error) *FJError {
	e.Cause = cause
	return e
}

// NewFJError creates a new FJ error with the specified type and message
func NewFJError(errorType ErrorType, message string) *FJError {
	return &FJError{
		Type:        errorType,
		Message:     message,
		Timestamp:   time.Now(),
		Recoverable: isRecoverableByDefault(errorType),
	}
}

// NewFJErrorWithCode creates a new FJ error with type, code, and message
func NewFJErrorWithCode(errorType ErrorType, code, message string) *FJError {
	return &FJError{
		Type:        errorType,
		Code:        code,
		Message:     message,
		Timestamp:   time.Now(),
		Recoverable: isRecoverableByDefault(errorType),
	}
}

// WrapError wraps an existing error with FJ error context
func WrapError(errorType ErrorType, message string, cause error) *FJError {
	return &FJError{
		Type:        errorType,
		Message:     message,
		Cause:       cause,
		Timestamp:   time.Now(),
		Recoverable: isRecoverableByDefault(errorType),
	}
}

// isRecoverableByDefault determines if an error type is recoverable by default
func isRecoverableByDefault(errorType ErrorType) bool {
	switch errorType {
	case JSONParseError, PathQueryError, ValidationError:
		return true
	case TransformationError, ConfigurationError, IOError:
		return true
	case TimeoutError:
		return true
	case MemoryError, PluginError:
		return false
	default:
		return false
	}
}

// Common error codes for consistent error reporting
const (
	// JSON Parse Error Codes
	CodeInvalidJSON          = "E001"
	CodeUnexpectedToken      = "E002"
	CodeUnterminatedString   = "E003"
	CodeInvalidNumber        = "E004"
	CodeInvalidEscape        = "E005"
	CodeMaxDepthExceeded     = "E006"

	// Path Query Error Codes
	CodeInvalidPath          = "E101"
	CodePathNotFound         = "E102"
	CodeInvalidPathSyntax    = "E103"
	CodeWildcardError        = "E104"

	// Transformation Error Codes
	CodeTransformerNotFound  = "E201"
	CodeInvalidArguments     = "E202"
	CodeTransformationFailed = "E203"

	// Validation Error Codes
	CodeValidationFailed     = "E301"
	CodeSchemaViolation      = "E302"

	// Configuration Error Codes
	CodeInvalidConfig        = "E401"
	CodeConfigNotFound       = "E402"

	// I/O Error Codes
	CodeFileNotFound         = "E501"
	CodePermissionDenied     = "E502"
	CodeReadError            = "E503"
	CodeWriteError           = "E504"

	// Timeout Error Codes
	CodeOperationTimeout     = "E601"

	// Memory Error Codes
	CodeOutOfMemory          = "E701"
	CodeMemoryLimitExceeded  = "E702"

	// Plugin Error Codes
	CodePluginNotFound       = "E801"
	CodePluginLoadError      = "E802"
	CodePluginExecutionError = "E803"
)

// Predefined error constructors for common scenarios

// NewJSONParseError creates a JSON parsing error
func NewJSONParseError(message string, position int) *FJError {
	return NewFJErrorWithCode(JSONParseError, CodeInvalidJSON, message).
		WithPosition(position)
}

// NewPathQueryError creates a path query error
func NewPathQueryError(message, path string) *FJError {
	return NewFJErrorWithCode(PathQueryError, CodeInvalidPath, message).
		WithPath(path)
}

// NewTransformationError creates a transformation error
func NewTransformationError(transformerName, message string) *FJError {
	return NewFJErrorWithCode(TransformationError, CodeTransformationFailed, message).
		WithContext("transformer", transformerName)
}

// NewValidationError creates a validation error
func NewValidationError(message string) *FJError {
	return NewFJErrorWithCode(ValidationError, CodeValidationFailed, message)
}

// NewConfigurationError creates a configuration error
func NewConfigurationError(message string) *FJError {
	return NewFJErrorWithCode(ConfigurationError, CodeInvalidConfig, message)
}

// NewIOError creates an I/O error
func NewIOError(message string, cause error) *FJError {
	return NewFJErrorWithCode(IOError, CodeReadError, message).
		WithCause(cause)
}

// NewTimeoutError creates a timeout error
func NewTimeoutError(operation string, timeout time.Duration) *FJError {
	return NewFJErrorWithCode(TimeoutError, CodeOperationTimeout,
		fmt.Sprintf("operation '%s' timed out after %v", operation, timeout)).
		WithContext("operation", operation).
		WithContext("timeout", timeout)
}