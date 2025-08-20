package fj

import (
	"io"
)

// contextFactory implements the ContextFactory interface.
// It provides various methods for creating Context instances using the Factory pattern.
type contextFactory struct{}

var (
	// DefaultContextFactory is the default factory instance for creating Context objects.
	DefaultContextFactory ContextFactory = &contextFactory{}
)

// NewContextFactory creates a new ContextFactory instance.
// This function serves as a factory method for creating context factories.
//
// Returns:
//   - ContextFactory: A new context factory instance
//
// Example Usage:
//
//	factory := NewContextFactory()
//	ctx := factory.Create(`{"name": "John"}`, "name")
func NewContextFactory() ContextFactory {
	return &contextFactory{}
}

// Create creates a new Context from JSON string and optional path.
// This method parses the JSON and optionally applies a path query to extract specific values.
//
// Parameters:
//   - json: The JSON string to parse
//   - path: Optional path expressions to query specific values
//
// Returns:
//   - Context: The created Context instance
//
// Example Usage:
//
//	factory := NewContextFactory()
//	// Create context from entire JSON
//	ctx1 := factory.Create(`{"name": "John", "age": 30}`)
//	
//	// Create context with path query
//	ctx2 := factory.Create(`{"user": {"name": "John"}}`, "user.name")
func (cf *contextFactory) Create(json string, path ...string) Context {
	if len(path) == 0 {
		return Parse(json)
	}
	return Get(json, path[0])
}

// CreateFromBytes creates a new Context from byte slice and optional path.
// This method parses JSON from a byte slice and optionally applies a path query.
//
// Parameters:
//   - data: The JSON byte slice to parse
//   - path: Optional path expressions to query specific values
//
// Returns:
//   - Context: The created Context instance
//
// Example Usage:
//
//	factory := NewContextFactory()
//	jsonBytes := []byte(`{"name": "John", "age": 30}`)
//	
//	// Create context from entire JSON
//	ctx1 := factory.CreateFromBytes(jsonBytes)
//	
//	// Create context with path query  
//	ctx2 := factory.CreateFromBytes(jsonBytes, "name")
func (cf *contextFactory) CreateFromBytes(data []byte, path ...string) Context {
	if len(path) == 0 {
		return ParseBytes(data)
	}
	return GetBytes(data, path[0])
}

// CreateEmpty creates an empty Context.
// This method returns a Context that represents no value (similar to null).
//
// Returns:
//   - Context: An empty Context instance
//
// Example Usage:
//
//	factory := NewContextFactory()
//	emptyCtx := factory.CreateEmpty()
//	fmt.Println(emptyCtx.Exists()) // Output: false
func (cf *contextFactory) CreateEmpty() Context {
	return Context{}
}

// CreateWithError creates a Context with an error.
// This method creates a Context that contains error information.
//
// Parameters:
//   - err: The error to associate with the Context
//
// Returns:
//   - Context: A Context instance containing the error
//
// Example Usage:
//
//	factory := NewContextFactory()
//	ctx := factory.CreateWithError(errors.New("parsing failed"))
//	fmt.Println(ctx.Error()) // Output: "parsing failed"
func (cf *contextFactory) CreateWithError(err error) Context {
	return Context{err: err}
}

// CreateFromReader creates a Context from an io.Reader.
// This method reads JSON data from a Reader and parses it into a Context.
//
// Parameters:
//   - reader: The io.Reader containing JSON data
//   - path: Optional path expressions to query specific values
//
// Returns:
//   - Context: The created Context instance
//
// Example Usage:
//
//	factory := NewContextFactory()
//	file, _ := os.Open("data.json")
//	defer file.Close()
//	
//	ctx := CreateFromReader(file, "user.name")
func CreateFromReader(reader io.Reader, path ...string) Context {
	ctx := ParseBufio(reader)
	if ctx.err != nil {
		return ctx
	}
	if len(path) == 0 {
		return ctx
	}
	return ctx.Get(path[0])
}

// Convenience functions using the default factory

// CreateContext creates a Context using the default factory.
// This is a convenience function for common Context creation scenarios.
//
// Parameters:
//   - json: The JSON string to parse
//   - path: Optional path expressions to query specific values
//
// Returns:
//   - Context: The created Context instance
//
// Example Usage:
//
//	ctx := CreateContext(`{"name": "John"}`, "name")
//	fmt.Println(ctx.String()) // Output: "John"
func CreateContext(json string, path ...string) Context {
	return DefaultContextFactory.Create(json, path...)
}

// CreateContextFromBytes creates a Context from bytes using the default factory.
// This is a convenience function for creating Context from byte slices.
//
// Parameters:
//   - data: The JSON byte slice to parse
//   - path: Optional path expressions to query specific values
//
// Returns:
//   - Context: The created Context instance
//
// Example Usage:
//
//	data := []byte(`{"name": "John"}`)
//	ctx := CreateContextFromBytes(data, "name") 
//	fmt.Println(ctx.String()) // Output: "John"
func CreateContextFromBytes(data []byte, path ...string) Context {
	return DefaultContextFactory.CreateFromBytes(data, path...)
}

// CreateEmptyContext creates an empty Context using the default factory.
// This is a convenience function for creating empty Context instances.
//
// Returns:
//   - Context: An empty Context instance
//
// Example Usage:
//
//	ctx := CreateEmptyContext()
//	fmt.Println(ctx.Exists()) // Output: false
func CreateEmptyContext() Context {
	return DefaultContextFactory.CreateEmpty()
}

// CreateContextWithError creates a Context with error using the default factory.
// This is a convenience function for creating Context instances with errors.
//
// Parameters:
//   - err: The error to associate with the Context
//
// Returns:
//   - Context: A Context instance containing the error
//
// Example Usage:
//
//	ctx := CreateContextWithError(errors.New("invalid JSON"))
//	fmt.Println(ctx.Error()) // Output: "invalid JSON"
func CreateContextWithError(err error) Context {
	return DefaultContextFactory.CreateWithError(err)
}