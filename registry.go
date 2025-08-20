package fj

import (
	"errors"
	"sync"
)

// transformerRegistry implements the TransformerRegistry interface as a singleton.
// It provides thread-safe access to transformer registration and lookup operations.
type transformerRegistry struct {
	mu           sync.RWMutex
	transformers map[string]Transformer
	functions    map[string]func(json, arg string) string
}

var (
	// registryInstance holds the singleton instance of the transformer registry
	registryInstance *transformerRegistry
	// registryOnce ensures the singleton is created only once
	registryOnce sync.Once
)

// GetTransformerRegistry returns the singleton instance of the transformer registry.
// This function implements the Singleton pattern to ensure only one registry exists.
//
// Returns:
//   - TransformerRegistry: The singleton transformer registry instance
//
// Example Usage:
//
//	registry := GetTransformerRegistry()
//	err := registry.RegisterFunc("uppercase", func(json, arg string) string {
//		return strings.ToUpper(json)
//	})
func GetTransformerRegistry() TransformerRegistry {
	registryOnce.Do(func() {
		registryInstance = &transformerRegistry{
			transformers: make(map[string]Transformer),
			functions:    make(map[string]func(json, arg string) string),
		}
		// Initialize with existing transformers from global map
		registryInstance.initializeBuiltinTransformers()
	})
	return registryInstance
}

// Register adds a new transformer to the registry.
// This method provides thread-safe registration of transformer implementations.
//
// Parameters:
//   - name: Unique identifier for the transformer
//   - transformer: The transformer implementation to register
//
// Returns:
//   - error: Error if registration fails (e.g., duplicate name)
//
// Example Usage:
//
//	type upperTransformer struct{}
//	func (u upperTransformer) Transform(json, arg string) string { return strings.ToUpper(json) }
//	func (u upperTransformer) Name() string { return "upper" }
//	
//	registry := GetTransformerRegistry()
//	err := registry.Register("upper", upperTransformer{})
func (tr *transformerRegistry) Register(name string, transformer Transformer) error {
	if name == "" {
		return errors.New("transformer name cannot be empty")
	}
	if transformer == nil {
		return errors.New("transformer cannot be nil")
	}

	tr.mu.Lock()
	defer tr.mu.Unlock()

	// Check for duplicates
	if _, exists := tr.transformers[name]; exists {
		return errors.New("transformer with name '" + name + "' already exists")
	}
	if _, exists := tr.functions[name]; exists {
		return errors.New("function transformer with name '" + name + "' already exists")
	}

	tr.transformers[name] = transformer
	return nil
}

// RegisterFunc adds a function-based transformer to the registry.
// This method provides thread-safe registration of function transformers.
//
// Parameters:
//   - name: Unique identifier for the transformer
//   - fn: The transformation function to register
//
// Returns:
//   - error: Error if registration fails (e.g., duplicate name)
//
// Example Usage:
//
//	registry := GetTransformerRegistry()
//	err := registry.RegisterFunc("lowercase", func(json, arg string) string {
//		return strings.ToLower(json)
//	})
func (tr *transformerRegistry) RegisterFunc(name string, fn func(json, arg string) string) error {
	if name == "" {
		return errors.New("transformer name cannot be empty")
	}
	if fn == nil {
		return errors.New("transformer function cannot be nil")
	}

	tr.mu.Lock()
	defer tr.mu.Unlock()

	// Check for duplicates
	if _, exists := tr.transformers[name]; exists {
		return errors.New("object transformer with name '" + name + "' already exists")
	}
	if _, exists := tr.functions[name]; exists {
		return errors.New("function transformer with name '" + name + "' already exists")
	}

	tr.functions[name] = fn
	return nil
}

// Get retrieves a transformer by name.
// This method provides thread-safe lookup of registered transformers.
//
// Parameters:
//   - name: The name of the transformer to retrieve
//
// Returns:
//   - Transformer: The transformer implementation (nil if not found)
//   - bool: True if transformer was found, false otherwise
//
// Example Usage:
//
//	registry := GetTransformerRegistry()
//	transformer, exists := registry.Get("uppercase")
//	if exists {
//		result := transformer.Transform(`{"name": "john"}`, "")
//	}
func (tr *transformerRegistry) Get(name string) (Transformer, bool) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	// Check object transformers first
	if transformer, exists := tr.transformers[name]; exists {
		return transformer, true
	}

	// Check function transformers
	if fn, exists := tr.functions[name]; exists {
		return &functionTransformer{name: name, fn: fn}, true
	}

	return nil, false
}

// Exists checks if a transformer is registered.
// This method provides thread-safe existence checking.
//
// Parameters:
//   - name: The name of the transformer to check
//
// Returns:
//   - bool: True if transformer exists, false otherwise
//
// Example Usage:
//
//	registry := GetTransformerRegistry()
//	if registry.Exists("uppercase") {
//		// Use the transformer
//	}
func (tr *transformerRegistry) Exists(name string) bool {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	_, objExists := tr.transformers[name]
	_, fnExists := tr.functions[name]
	return objExists || fnExists
}

// List returns all registered transformer names.
// This method provides thread-safe listing of all transformer names.
//
// Returns:
//   - []string: Slice of all registered transformer names
//
// Example Usage:
//
//	registry := GetTransformerRegistry()
//	names := registry.List()
//	for _, name := range names {
//		fmt.Println("Available transformer:", name)
//	}
func (tr *transformerRegistry) List() []string {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	names := make([]string, 0, len(tr.transformers)+len(tr.functions))
	
	for name := range tr.transformers {
		names = append(names, name)
	}
	
	for name := range tr.functions {
		names = append(names, name)
	}
	
	return names
}

// Remove removes a transformer from the registry.
// This method provides thread-safe removal of transformers.
//
// Parameters:
//   - name: The name of the transformer to remove
//
// Returns:
//   - bool: True if transformer was removed, false if it didn't exist
//
// Example Usage:
//
//	registry := GetTransformerRegistry()
//	removed := registry.Remove("old-transformer")
func (tr *transformerRegistry) Remove(name string) bool {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	_, objExists := tr.transformers[name]
	_, fnExists := tr.functions[name]

	if objExists {
		delete(tr.transformers, name)
		return true
	}
	
	if fnExists {
		delete(tr.functions, name)
		return true
	}

	return false
}

// initializeBuiltinTransformers populates the registry with existing global transformers.
// This method migrates existing transformers from the global map to the registry.
func (tr *transformerRegistry) initializeBuiltinTransformers() {
	// Copy from global jsonTransformers map
	for name, fn := range jsonTransformers {
		tr.functions[name] = fn
	}
}

// functionTransformer wraps a function to implement the Transformer interface.
type functionTransformer struct {
	name string
	fn   func(json, arg string) string
}

// Transform applies the wrapped function transformation.
func (ft *functionTransformer) Transform(json, arg string) string {
	return ft.fn(json, arg)
}

// Name returns the transformer name.
func (ft *functionTransformer) Name() string {
	return ft.name
}

// RegisterTransformer is a convenience function that registers a transformer.
// This function provides a simple way to register transformers with the global registry.
//
// Parameters:
//   - name: Unique identifier for the transformer
//   - transformer: The transformer implementation to register
//
// Returns:
//   - error: Error if registration fails
//
// Example Usage:
//
//	err := RegisterTransformer("myTransformer", myTransformerInstance)
func RegisterTransformer(name string, transformer Transformer) error {
	return GetTransformerRegistry().Register(name, transformer)
}

// RegisterTransformerFunc is a convenience function that registers a function transformer.
// This function provides a simple way to register function transformers with the global registry.
//
// Parameters:
//   - name: Unique identifier for the transformer
//   - fn: The transformation function to register
//
// Returns:
//   - error: Error if registration fails
//
// Example Usage:
//
//	err := RegisterTransformerFunc("lowercase", func(json, arg string) string {
//		return strings.ToLower(json)
//	})
func RegisterTransformerFunc(name string, fn func(json, arg string) string) error {
	return GetTransformerRegistry().RegisterFunc(name, fn)
}