package fj

// Parser defines the interface for JSON parsing operations.
// This interface provides a contract for different parsing strategies
// and enables dependency injection and testing.
type Parser interface {
	// Parse parses a JSON string and returns a Context
	Parse(json string) Context
	
	// ParseBytes parses JSON from a byte slice
	ParseBytes(data []byte) Context
}

// PathBuilder defines the interface for building JSON paths fluently.
// This interface implements the Builder pattern for constructing
// JSON paths in a readable and chainable manner.
type PathBuilder interface {
	// Root starts a path from the root with the given key
	Root(key string) PathBuilder
	
	// Field adds a field to the path
	Field(key string) PathBuilder
	
	// ArrayIndex adds an array index to the path
	ArrayIndex(index int) PathBuilder
	
	// ArrayAll adds array wildcard to the path
	ArrayAll() PathBuilder
	
	// Build constructs the final path string
	Build() string
	
	// Reset resets the builder to empty state
	Reset() PathBuilder
}

// Transformer defines the interface for JSON transformation operations.
// This interface implements the Strategy pattern, allowing different
// transformation strategies to be plugged in at runtime.
type Transformer interface {
	// Transform applies the transformation to the JSON string
	Transform(json, arg string) string
	
	// Name returns the unique name identifier for this transformer
	Name() string
}

// TransformerRegistry defines the interface for managing transformers.
// This interface provides centralized transformer registration and lookup.
type TransformerRegistry interface {
	// Register adds a new transformer to the registry
	Register(name string, transformer Transformer) error
	
	// RegisterFunc adds a function-based transformer to the registry  
	RegisterFunc(name string, fn func(json, arg string) string) error
	
	// Get retrieves a transformer by name
	Get(name string) (Transformer, bool)
	
	// Exists checks if a transformer is registered
	Exists(name string) bool
	
	// List returns all registered transformer names
	List() []string
	
	// Remove removes a transformer from the registry
	Remove(name string) bool
}

// ContextFactory defines the interface for creating Context instances.
// This interface implements the Factory pattern for Context creation.
type ContextFactory interface {
	// Create creates a new Context from JSON string and optional path
	Create(json string, path ...string) Context
	
	// CreateFromBytes creates a new Context from byte slice and optional path  
	CreateFromBytes(data []byte, path ...string) Context
	
	// CreateEmpty creates an empty Context
	CreateEmpty() Context
	
	// CreateWithError creates a Context with an error
	CreateWithError(err error) Context
}

// PipeProcessor defines the interface for processing pipe operations.
// This interface implements the Chain of Responsibility pattern for
// handling sequential transformations.
type PipeProcessor interface {
	// Process applies pipe operations to the context
	Process(ctx Context, pipe string) Context
	
	// AddHandler adds a new pipe handler to the chain
	AddHandler(handler PipeHandler) PipeProcessor
}

// PipeHandler defines the interface for individual pipe operation handlers.
type PipeHandler interface {
	// Handle processes a single pipe operation
	Handle(ctx Context, operation string, arg string) (Context, bool)
	
	// CanHandle checks if this handler can process the given operation
	CanHandle(operation string) bool
}

// Validator defines the interface for JSON validation operations.
type Validator interface {
	// Validate checks if the given data is valid JSON
	Validate(data []byte) bool
	
	// ValidateString checks if the given string is valid JSON
	ValidateString(json string) bool
}

// QueryProcessor defines the interface for processing JSON queries.
type QueryProcessor interface {
	// Query processes a query against JSON data
	Query(json string, query string) []Context
	
	// QueryFirst processes a query and returns the first result
	QueryFirst(json string, query string) Context
}