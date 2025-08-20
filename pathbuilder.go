package fj

import (
	"strconv"
	"strings"
)

// pathBuilder implements the PathBuilder interface using the Builder pattern.
// It provides a fluent API for constructing JSON paths in a readable manner.
type pathBuilder struct {
	segments []string
}

// NewPathBuilder creates a new PathBuilder instance.
// This function serves as a factory method for creating path builders.
//
// Returns:
//   - PathBuilder: A new path builder instance ready for use
//
// Example Usage:
//
//	builder := NewPathBuilder().
//		Root("user").
//		Field("profile").
//		ArrayIndex(0).
//		Field("name")
//	path := builder.Build() // Returns "user.profile.0.name"
func NewPathBuilder() PathBuilder {
	return &pathBuilder{
		segments: make([]string, 0),
	}
}

// Root starts a path from the root with the given key.
// This method initializes the path builder with a root key.
//
// Parameters:
//   - key: The root key to start the path with
//
// Returns:
//   - PathBuilder: The builder instance for method chaining
//
// Example Usage:
//
//	builder := NewPathBuilder().Root("users")
func (pb *pathBuilder) Root(key string) PathBuilder {
	pb.segments = []string{key}
	return pb
}

// Field adds a field to the path.
// This method appends a field key to the current path.
//
// Parameters:
//   - key: The field key to add to the path
//
// Returns:
//   - PathBuilder: The builder instance for method chaining
//
// Example Usage:
//
//	builder := NewPathBuilder().Root("user").Field("name")
func (pb *pathBuilder) Field(key string) PathBuilder {
	pb.segments = append(pb.segments, key)
	return pb
}

// ArrayIndex adds an array index to the path.
// This method appends an array index to the current path.
//
// Parameters:
//   - index: The array index to add to the path
//
// Returns:
//   - PathBuilder: The builder instance for method chaining
//
// Example Usage:
//
//	builder := NewPathBuilder().Root("users").ArrayIndex(0)
func (pb *pathBuilder) ArrayIndex(index int) PathBuilder {
	pb.segments = append(pb.segments, strconv.Itoa(index))
	return pb
}

// ArrayAll adds array wildcard to the path.
// This method appends the wildcard selector "*" to select all array elements.
//
// Returns:
//   - PathBuilder: The builder instance for method chaining
//
// Example Usage:
//
//	builder := NewPathBuilder().Root("users").ArrayAll()
func (pb *pathBuilder) ArrayAll() PathBuilder {
	pb.segments = append(pb.segments, "*")
	return pb
}

// Build constructs the final path string.
// This method joins all path segments with dots to create the final path.
//
// Returns:
//   - string: The constructed JSON path string
//
// Example Usage:
//
//	path := NewPathBuilder().
//		Root("user").
//		Field("roles").
//		ArrayIndex(0).
//		Field("permissions").
//		Build() // Returns "user.roles.0.permissions"
func (pb *pathBuilder) Build() string {
	return strings.Join(pb.segments, ".")
}

// Reset resets the builder to empty state.
// This method clears all segments and prepares the builder for reuse.
//
// Returns:
//   - PathBuilder: The builder instance for method chaining
//
// Example Usage:
//
//	builder := NewPathBuilder().Root("user").Field("name")
//	path1 := builder.Build()
//	path2 := builder.Reset().Root("orders").ArrayIndex(0).Build()
func (pb *pathBuilder) Reset() PathBuilder {
	pb.segments = pb.segments[:0]
	return pb
}

// PathBuilderFromString creates a PathBuilder from an existing path string.
// This function parses a dot-separated path string and creates a builder with those segments.
//
// Parameters:
//   - path: The existing path string to parse
//
// Returns:
//   - PathBuilder: A path builder initialized with the parsed segments
//
// Example Usage:
//
//	builder := PathBuilderFromString("user.profile.name")
//	newPath := builder.Field("first").Build() // Returns "user.profile.name.first"
func PathBuilderFromString(path string) PathBuilder {
	if path == "" {
		return NewPathBuilder()
	}
	
	segments := strings.Split(path, ".")
	return &pathBuilder{
		segments: segments,
	}
}