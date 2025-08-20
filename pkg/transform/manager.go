package transform

import (
	"fmt"
	"sync"
	"time"

	"github.com/sivaosorg/fj/pkg/core"
	"github.com/sivaosorg/fj/pkg/errors"
)

// BaseTransformer provides a base implementation for transformers
type BaseTransformer struct {
	name        string
	description string
}

// GetName returns the transformer name
func (bt *BaseTransformer) GetName() string {
	return bt.name
}

// GetDescription returns the transformer description
func (bt *BaseTransformer) GetDescription() string {
	return bt.description
}

// Validate provides default validation (no validation)
func (bt *BaseTransformer) Validate(arg string) error {
	return nil // Default: accept any arguments
}

// Factory implements the TransformationFactory interface using the Factory pattern.
// It manages the creation and registration of JSON transformers.
type Factory struct {
	transformers map[string]core.JSONTransformer
	mutex        sync.RWMutex
}

// NewFactory creates a new transformation factory
func NewFactory() core.TransformationFactory {
	factory := &Factory{
		transformers: make(map[string]core.JSONTransformer),
	}
	
	// Register built-in transformers
	factory.registerBuiltInTransformers()
	
	return factory
}

// CreateTransformer creates a transformer by name
func (f *Factory) CreateTransformer(name string) (core.JSONTransformer, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	
	transformer, exists := f.transformers[name]
	if !exists {
		return nil, errors.NewTransformationError(name, 
			fmt.Sprintf("transformer '%s' not found", name))
	}
	
	return transformer, nil
}

// RegisterTransformer registers a new transformer
func (f *Factory) RegisterTransformer(name string, transformer core.JSONTransformer) error {
	if name == "" {
		return errors.NewTransformationError("", "transformer name cannot be empty")
	}
	
	if transformer == nil {
		return errors.NewTransformationError(name, "transformer cannot be nil")
	}
	
	f.mutex.Lock()
	defer f.mutex.Unlock()
	
	f.transformers[name] = transformer
	return nil
}

// UnregisterTransformer removes a transformer
func (f *Factory) UnregisterTransformer(name string) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	
	if _, exists := f.transformers[name]; !exists {
		return errors.NewTransformationError(name, 
			fmt.Sprintf("transformer '%s' not registered", name))
	}
	
	delete(f.transformers, name)
	return nil
}

// ListTransformers returns all registered transformer names
func (f *Factory) ListTransformers() []string {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	
	names := make([]string, 0, len(f.transformers))
	for name := range f.transformers {
		names = append(names, name)
	}
	return names
}

// IsRegistered checks if a transformer is registered
func (f *Factory) IsRegistered(name string) bool {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	
	_, exists := f.transformers[name]
	return exists
}

// registerBuiltInTransformers registers the built-in transformers
func (f *Factory) registerBuiltInTransformers() {
	// Register built-in transformers that wrap the existing transformer functions
	
	f.transformers["trim"] = &BuiltInTransformer{
		BaseTransformer: BaseTransformer{
			name:        "trim",
			description: "Trims whitespace from JSON strings",
		},
		transformFunc: f.wrapExistingTransformer("trim"),
	}
	
	f.transformers["pretty"] = &BuiltInTransformer{
		BaseTransformer: BaseTransformer{
			name:        "pretty",
			description: "Pretty prints JSON with indentation",
		},
		transformFunc: f.wrapExistingTransformer("pretty"),
	}
	
	f.transformers["minify"] = &BuiltInTransformer{
		BaseTransformer: BaseTransformer{
			name:        "minify",
			description: "Minifies JSON by removing whitespace",
		},
		transformFunc: f.wrapExistingTransformer("minify"),
	}
	
	f.transformers["uppercase"] = &BuiltInTransformer{
		BaseTransformer: BaseTransformer{
			name:        "uppercase",
			description: "Converts JSON string values to uppercase",
		},
		transformFunc: f.wrapExistingTransformer("uppercase"),
	}
	
	f.transformers["lowercase"] = &BuiltInTransformer{
		BaseTransformer: BaseTransformer{
			name:        "lowercase",
			description: "Converts JSON string values to lowercase",
		},
		transformFunc: f.wrapExistingTransformer("lowercase"),
	}
	
	f.transformers["reverse"] = &BuiltInTransformer{
		BaseTransformer: BaseTransformer{
			name:        "reverse",
			description: "Reverses the order of array elements or object properties",
		},
		transformFunc: f.wrapExistingTransformer("reverse"),
	}
	
	f.transformers["flatten"] = &BuiltInTransformer{
		BaseTransformer: BaseTransformer{
			name:        "flatten",
			description: "Flattens nested JSON structures",
		},
		transformFunc: f.wrapExistingTransformer("flatten"),
	}
	
	f.transformers["keys"] = &BuiltInTransformer{
		BaseTransformer: BaseTransformer{
			name:        "keys",
			description: "Extracts keys from JSON objects",
		},
		transformFunc: f.wrapExistingTransformer("keys"),
	}
	
	f.transformers["values"] = &BuiltInTransformer{
		BaseTransformer: BaseTransformer{
			name:        "values",
			description: "Extracts values from JSON objects",
		},
		transformFunc: f.wrapExistingTransformer("values"),
	}
	
	f.transformers["valid"] = &BuiltInTransformer{
		BaseTransformer: BaseTransformer{
			name:        "valid",
			description: "Validates JSON syntax and returns boolean result",
		},
		transformFunc: f.wrapExistingTransformer("valid"),
	}
}

// wrapExistingTransformer creates a wrapper function for existing transformers
func (f *Factory) wrapExistingTransformer(name string) func(string, string) string {
	return func(json, arg string) string {
		// This would call the existing transformer functions from the root package
		// For now, this is a placeholder that will be connected to the actual transformers
		return json // placeholder
	}
}

// BuiltInTransformer represents a built-in transformer that wraps existing functionality
type BuiltInTransformer struct {
	BaseTransformer
	transformFunc func(string, string) string
}

// Transform applies the transformation
func (bit *BuiltInTransformer) Transform(json, arg string) string {
	if bit.transformFunc == nil {
		return json
	}
	return bit.transformFunc(json, arg)
}

// Manager provides a high-level interface for managing transformations
// and implements the Strategy pattern for different transformation approaches
type Manager struct {
	factory       core.TransformationFactory
	eventPublisher core.EventPublisher
}

// NewManager creates a new transformation manager
func NewManager(factory core.TransformationFactory) *Manager {
	return &Manager{
		factory: factory,
	}
}

// SetEventPublisher sets the event publisher for transformation events
func (m *Manager) SetEventPublisher(publisher core.EventPublisher) {
	m.eventPublisher = publisher
}

// Transform applies a transformation with event publishing and error handling
func (m *Manager) Transform(transformerName, json, arg string) (string, error) {
	// Create transformer
	transformer, err := m.factory.CreateTransformer(transformerName)
	if err != nil {
		return "", err
	}
	
	// Validate arguments
	if err := transformer.Validate(arg); err != nil {
		return "", errors.WrapError(errors.TransformationError,
			fmt.Sprintf("invalid arguments for transformer '%s'", transformerName), err).
			WithContext("transformer", transformerName).
			WithContext("arguments", arg)
	}
	
	// Publish start event if available
	if m.eventPublisher != nil {
		m.eventPublisher.PublishTransformStart(transformerName, json, []string{arg})
	}
	
	startTime := time.Now()
	
	// Apply transformation
	result := transformer.Transform(json, arg)
	
	duration := time.Since(startTime)
	
	// Publish completion event if available
	if m.eventPublisher != nil {
		m.eventPublisher.PublishTransformComplete(transformerName, result, duration)
	}
	
	return result, nil
}

// ChainTransformations applies multiple transformations in sequence (Decorator pattern)
func (m *Manager) ChainTransformations(json string, transformations []TransformationStep) (string, error) {
	result := json
	
	for i, step := range transformations {
		transformed, err := m.Transform(step.Name, result, step.Argument)
		if err != nil {
			return "", errors.WrapError(errors.TransformationError,
				fmt.Sprintf("transformation chain failed at step %d", i+1), err).
				WithContext("step", i+1).
				WithContext("transformer", step.Name).
				WithContext("argument", step.Argument)
		}
		result = transformed
	}
	
	return result, nil
}

// TransformationStep represents a single step in a transformation chain
type TransformationStep struct {
	Name     string `json:"name"`
	Argument string `json:"argument"`
}

// GetFactory returns the transformation factory (for access to registration methods)
func (m *Manager) GetFactory() core.TransformationFactory {
	return m.factory
}