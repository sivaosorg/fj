package fj

// Enhanced function names with improved Go conventions and shorter, cleaner names.
// These functions provide more intuitive and professional naming while maintaining
// backward compatibility with existing APIs.

// Extract retrieves a value from JSON using a path expression.
// This is an enhanced alias for Get() with improved naming conventions.
//
// Parameters:
//   - json: The JSON string to extract from
//   - path: The path expression to the desired value
//
// Returns:
//   - Context: The extracted value context
//
// Example Usage:
//
//	ctx := Extract(`{"user":{"name":"John"}}`, "user.name")
//	fmt.Println(ctx.String()) // Output: "John"
func Extract(json, path string) Context {
	return Get(json, path)
}

// ExtractBytes retrieves a value from JSON bytes using a path expression.
// This is an enhanced alias for GetBytes() with improved naming.
//
// Parameters:
//   - data: The JSON byte slice to extract from
//   - path: The path expression to the desired value
//
// Returns:
//   - Context: The extracted value context
//
// Example Usage:
//
//	data := []byte(`{"user":{"name":"John"}}`)
//	ctx := ExtractBytes(data, "user.name")
//	fmt.Println(ctx.String()) // Output: "John"
func ExtractBytes(data []byte, path string) Context {
	return GetBytes(data, path)
}

// ExtractMultiple retrieves multiple values from JSON using path expressions.
// This is an enhanced alias for GetMul() with improved naming.
//
// Parameters:
//   - json: The JSON string to extract from
//   - paths: The path expressions to the desired values
//
// Returns:
//   - []Context: The extracted value contexts
//
// Example Usage:
//
//	contexts := ExtractMultiple(`{"name":"John","age":30}`, "name", "age")
//	fmt.Println(contexts[0].String()) // Output: "John"
//	fmt.Println(contexts[1].String()) // Output: "30"
func ExtractMultiple(json string, paths ...string) []Context {
	return GetMul(json, paths...)
}

// Validate checks if the provided data is valid JSON.
// This is an enhanced alias for IsValidJSONBytes() with improved naming.
//
// Parameters:
//   - data: The byte slice to validate
//
// Returns:
//   - bool: True if valid JSON, false otherwise
//
// Example Usage:
//
//	data := []byte(`{"name":"John"}`)
//	if Validate(data) {
//		fmt.Println("Valid JSON")
//	}
func Validate(data []byte) bool {
	return IsValidJSONBytes(data)
}

// ValidateString checks if the provided string is valid JSON.
// This provides a convenient string validation method.
//
// Parameters:
//   - json: The JSON string to validate
//
// Returns:
//   - bool: True if valid JSON, false otherwise
//
// Example Usage:
//
//	if ValidateString(`{"name":"John"}`) {
//		fmt.Println("Valid JSON")
//	}
func ValidateString(json string) bool {
	return IsValidJSONBytes([]byte(json))
}

// Transform applies a transformation to JSON using a transformer name.
// This provides a clean interface for applying transformations.
//
// Parameters:
//   - json: The JSON string to transform
//   - transformerName: The name of the transformer to apply
//   - args: Optional arguments for the transformation
//
// Returns:
//   - string: The transformed JSON string
//
// Example Usage:
//
//	result := Transform(`{"name":"john"}`, "uppercase")
//	fmt.Println(result) // Output: `{"NAME":"JOHN"}`
func Transform(json, transformerName string, args ...string) string {
	registry := GetTransformerRegistry()
	transformer, exists := registry.Get(transformerName)
	if !exists {
		return json
	}
	
	arg := ""
	if len(args) > 0 {
		arg = args[0]
	}
	
	return transformer.Transform(json, arg)
}

// Chain applies multiple transformations in sequence.
// This provides a convenient way to chain transformations.
//
// Parameters:
//   - json: The JSON string to transform
//   - transformerNames: The names of transformers to apply in order
//
// Returns:
//   - string: The final transformed JSON string
//
// Example Usage:
//
//	result := Chain(`{"name":" john "}`, "trim", "uppercase")
//	fmt.Println(result) // Output: transformed JSON
func Chain(json string, transformerNames ...string) string {
	chain := NewTransformerStrategyChain()
	for _, name := range transformerNames {
		chain.Add(name)
	}
	return chain.Execute(json)
}

// Query processes a query expression against JSON data.
// This provides a clean interface for JSON querying.
//
// Parameters:
//   - json: The JSON string to query
//   - query: The query expression
//
// Returns:
//   - Context: The query result context
//
// Example Usage:
//
//	ctx := Query(`[{"name":"john"},{"name":"jane"}]`, `#(name=="john")`)
//	fmt.Println(ctx.String()) // Output: matching result
func Query(json, query string) Context {
	return Get(json, query)
}

// Build creates a JSON path using the PathBuilder pattern.
// This is a convenience function for simple path building.
//
// Parameters:
//   - parts: The path parts to join
//
// Returns:
//   - string: The constructed path
//
// Example Usage:
//
//	path := Build("user", "profile", "name")
//	fmt.Println(path) // Output: "user.profile.name"
func Build(parts ...string) string {
	builder := NewPathBuilder()
	for i, part := range parts {
		if i == 0 {
			builder.Root(part)
		} else {
			builder.Field(part)
		}
	}
	return builder.Build()
}

// Register adds a transformer to the registry.
// This is a convenience function for transformer registration.
//
// Parameters:
//   - name: The unique name for the transformer
//   - fn: The transformation function
//
// Returns:
//   - error: Error if registration fails
//
// Example Usage:
//
//	err := Register("myTransform", func(json, arg string) string {
//		return strings.ToUpper(json)
//	})
func Register(name string, fn func(json, arg string) string) error {
	return GetTransformerRegistry().RegisterFunc(name, fn)
}

// Exists checks if a transformer is registered.
// This is a convenience function for checking transformer existence.
//
// Parameters:
//   - name: The transformer name to check
//
// Returns:
//   - bool: True if transformer exists, false otherwise
//
// Example Usage:
//
//	if Exists("uppercase") {
//		fmt.Println("Transformer available")
//	}
func Exists(name string) bool {
	return GetTransformerRegistry().Exists(name)
}

// Process applies pipe operations to a Context.
// This is a convenience function for pipe processing.
//
// Parameters:
//   - ctx: The Context to process
//   - operations: The pipe operations string
//
// Returns:
//   - Context: The processed Context
//
// Example Usage:
//
//	result := Process(ctx, "uppercase|trim")
func Process(ctx Context, operations string) Context {
	return ProcessPipe(ctx, operations)
}

// From creates a Context from various input types.
// This is a flexible factory function for Context creation.
//
// Parameters:
//   - input: The input data (string, []byte, or io.Reader)
//   - paths: Optional path expressions
//
// Returns:
//   - Context: The created Context
//
// Example Usage:
//
//	ctx := From(`{"name":"John"}`, "name")
//	ctx2 := From([]byte(`{"age":30}`))
func From(input interface{}, paths ...string) Context {
	factory := NewContextFactory()
	
	switch v := input.(type) {
	case string:
		return factory.Create(v, paths...)
	case []byte:
		return factory.CreateFromBytes(v, paths...)
	case Context:
		if len(paths) > 0 {
			return v.Get(paths[0])
		}
		return v
	default:
		return factory.CreateEmpty()
	}
}

// Deprecated function aliases with deprecation notices

// GetBytes extracts a value from JSON bytes using a path expression.
// 
// Deprecated: Use ExtractBytes instead. This function will be removed in v2.0.0.
// ExtractBytes provides a cleaner, more intuitive name following Go conventions.
//
// Example migration:
//   // Old way
//   ctx := GetBytes(data, "path")
//   
//   // New way
//   ctx := ExtractBytes(data, "path")
func GetBytesDeprecated(json []byte, path string) Context {
	return GetBytes(json, path)
}

// GetMul extracts multiple values from JSON using path expressions.
//
// Deprecated: Use ExtractMultiple instead. This function will be removed in v2.0.0.
// ExtractMultiple provides a cleaner, more descriptive name.
//
// Example migration:
//   // Old way
//   contexts := GetMul(json, "path1", "path2")
//   
//   // New way
//   contexts := ExtractMultiple(json, "path1", "path2")
func GetMulDeprecated(json string, path ...string) []Context {
	return GetMul(json, path...)
}

// IsValidJSONBytes checks if byte data is valid JSON.
//
// Deprecated: Use Validate instead. This function will be removed in v2.0.0.
// Validate provides a shorter, cleaner name following Go conventions.
//
// Example migration:
//   // Old way
//   valid := IsValidJSONBytes(data)
//   
//   // New way
//   valid := Validate(data)
func IsValidJSONBytesDeprecated(json []byte) bool {
	return IsValidJSONBytes(json)
}

// AddTransformer adds a transformer function to the registry.
//
// Deprecated: Use Register instead. This function will be removed in v2.0.0.
// Register provides a shorter, cleaner name following Go conventions.
//
// Example migration:
//   // Old way
//   AddTransformer("name", func(json, arg string) string { return json })
//   
//   // New way
//   Register("name", func(json, arg string) string { return json })
func AddTransformerDeprecated(name string, fn func(json, arg string) string) {
	AddTransformer(name, fn)
}

// IsTransformerRegistered checks if a transformer is registered.
//
// Deprecated: Use Exists instead. This function will be removed in v2.0.0.
// Exists provides a shorter, cleaner name following Go conventions.
//
// Example migration:
//   // Old way
//   registered := IsTransformerRegistered("name")
//   
//   // New way
//   registered := Exists("name")
func IsTransformerRegisteredDeprecated(name string) bool {
	return IsTransformerRegistered(name)
}