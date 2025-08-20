package core

import (
	"io"
	"time"
)

// Forward declarations to avoid circular imports
type Config interface {
	GetEnableTransformers() bool
	GetEnableEventPublishing() bool
	GetMaxParseDepth() int
	GetDefaultTimeout() time.Duration
	GetCacheSize() int
	GetEnableValidation() bool
	GetLogLevel() int
	GetCustomTransformers() map[string]string
}

// JSONParser defines the interface for JSON parsing operations.
// This interface provides a clean contract for different parsing implementations
// and enables dependency injection and testing with mocks.
type JSONParser interface {
	// Parse parses a JSON string and returns a Context representing the parsed data
	Parse(json string) Context

	// ParseBytes parses JSON from a byte slice
	ParseBytes(data []byte) Context

	// ParseReader parses JSON from an io.Reader
	ParseReader(reader io.Reader) Context

	// ParseFile parses JSON from a file path
	ParseFile(filepath string) Context

	// IsValid checks if the provided JSON string is valid
	IsValid(json string) bool

	// IsValidBytes checks if the provided JSON bytes are valid
	IsValidBytes(data []byte) bool
}

// JSONTransformer defines the interface for JSON transformation operations.
// This interface supports the Strategy pattern for different transformation types
// and enables pluggable transformation implementations.
type JSONTransformer interface {
	// Transform applies a transformation to the JSON data
	Transform(json, arg string) string

	// GetName returns the name/identifier of the transformer
	GetName() string

	// GetDescription returns a human-readable description of what this transformer does
	GetDescription() string

	// Validate checks if the provided arguments are valid for this transformer
	Validate(arg string) error
}

// QueryProcessor defines the interface for processing JSON path queries.
// This interface supports the Chain of Responsibility pattern for query processing
// and enables complex query operations with different strategies.
type QueryProcessor interface {
	// Get retrieves a value from JSON using a path query
	Get(json, path string) Context

	// GetMultiple retrieves multiple values using multiple path queries
	GetMultiple(json string, paths ...string) []Context

	// Exists checks if a path exists in the JSON
	Exists(json, path string) bool

	// Set sets a value at the specified path (if supported)
	Set(json, path, value string) (string, error)

	// Delete removes a value at the specified path (if supported)
	Delete(json, path string) (string, error)
}

// Context represents a JSON value with enhanced OOP methods.
// This interface defines the contract for context operations,
// enabling different implementations and better testability.
type Context interface {
	// Type and validation methods
	Type() Type
	Kind() Type
	Exists() bool
	IsObject() bool
	IsArray() bool
	IsString() bool
	IsNumber() bool
	IsBool() bool
	IsNull() bool

	// Value extraction methods
	String() string
	Int() int64
	Float() float64
	Bool() bool
	Time() time.Time
	TimeWithFormat(format string) (time.Time, error)
	Value() interface{}

	// Collection methods
	Array() []Context
	Map() map[string]Context
	Len() int

	// Query methods
	Get(path string) Context
	GetMultiple(paths ...string) []Context
	PathExists(path string) bool

	// Iteration methods
	Foreach(iterator func(key, value Context) bool)
	ForEachValue(iterator func(value Context) bool)

	// Transformation methods
	Transform(transformerName string, args ...string) Context
	Pretty() string
	Minify() string

	// Comparison and sorting
	Less(other Context) bool
	Equal(other Context) bool

	// Error handling
	Error() error
	HasError() bool

	// Raw data access
	Raw() string
	Bytes() []byte
	Index() int
}

// TransformationFactory defines the interface for creating transformers.
// This interface implements the Factory pattern for transformer creation
// and enables pluggable transformer registration.
type TransformationFactory interface {
	// CreateTransformer creates a transformer by name
	CreateTransformer(name string) (JSONTransformer, error)

	// RegisterTransformer registers a new transformer
	RegisterTransformer(name string, transformer JSONTransformer) error

	// UnregisterTransformer removes a transformer
	UnregisterTransformer(name string) error

	// ListTransformers returns all registered transformer names
	ListTransformers() []string

	// IsRegistered checks if a transformer is registered
	IsRegistered(name string) bool
}

// ConfigurationManager defines the interface for managing library configuration.
// This interface implements the Singleton pattern for global configuration
// and provides centralized configuration management.
type ConfigurationManager interface {
	// GetConfig returns the current configuration
	GetConfig() Config

	// SetConfig updates the configuration
	SetConfig(config Config) error

	// GetSetting gets a specific configuration setting
	GetSetting(key string) (interface{}, bool)

	// SetSetting sets a specific configuration setting
	SetSetting(key string, value interface{}) error

	// LoadFromFile loads configuration from a file
	LoadFromFile(filepath string) error

	// SaveToFile saves configuration to a file
	SaveToFile(filepath string) error

	// Reset resets configuration to defaults
	Reset()
}

// EventObserver defines the interface for observing parsing and transformation events.
// This interface implements the Observer pattern for event handling
// and enables extensible event-driven processing.
type EventObserver interface {
	// OnParseStart is called when parsing starts
	OnParseStart(data string)

	// OnParseComplete is called when parsing completes successfully
	OnParseComplete(context Context, duration time.Duration)

	// OnParseError is called when parsing encounters an error
	OnParseError(data string, err error, duration time.Duration)

	// OnTransformStart is called when transformation starts
	OnTransformStart(transformerName string, data string, args []string)

	// OnTransformComplete is called when transformation completes successfully
	OnTransformComplete(transformerName string, result string, duration time.Duration)

	// OnTransformError is called when transformation encounters an error
	OnTransformError(transformerName string, data string, args []string, err error, duration time.Duration)
}

// EventPublisher defines the interface for publishing events to observers.
// This interface supports the Observer pattern by managing event subscriptions
// and publishing events to registered observers.
type EventPublisher interface {
	// Subscribe adds an observer to receive events
	Subscribe(observer EventObserver) error

	// Unsubscribe removes an observer from receiving events
	Unsubscribe(observer EventObserver) error

	// PublishParseStart publishes a parse start event
	PublishParseStart(data string)

	// PublishParseComplete publishes a parse complete event
	PublishParseComplete(context Context, duration time.Duration)

	// PublishParseError publishes a parse error event
	PublishParseError(data string, err error, duration time.Duration)

	// PublishTransformStart publishes a transform start event
	PublishTransformStart(transformerName string, data string, args []string)

	// PublishTransformComplete publishes a transform complete event
	PublishTransformComplete(transformerName string, result string, duration time.Duration)

	// PublishTransformError publishes a transform error event
	PublishTransformError(transformerName string, data string, args []string, err error, duration time.Duration)
}

// ValidationEngine defines the interface for input validation.
// This interface provides comprehensive validation capabilities
// for JSON data and transformation arguments.
type ValidationEngine interface {
	// ValidateJSON validates JSON syntax and structure
	ValidateJSON(data string) error

	// ValidateJSONBytes validates JSON bytes
	ValidateJSONBytes(data []byte) error

	// ValidatePath validates a JSON path query
	ValidatePath(path string) error

	// ValidateTransformerArgs validates arguments for a specific transformer
	ValidateTransformerArgs(transformerName string, args []string) error

	// AddValidator adds a custom validator function
	AddValidator(name string, validator func(interface{}) error) error

	// RemoveValidator removes a custom validator
	RemoveValidator(name string) error
}

// Type represents the different possible types for a JSON value (from existing code)
type Type int