package fj

import (
	"strings"
)

// pipeProcessor implements the PipeProcessor interface using Chain of Responsibility pattern.
// It manages a chain of handlers for processing different pipe operations.
type pipeProcessor struct {
	handlers []PipeHandler
}

// NewPipeProcessor creates a new pipe processor with default handlers.
// This function initializes a processor with built-in pipe operation handlers.
//
// Returns:
//   - PipeProcessor: A new pipe processor instance
//
// Example Usage:
//
//	processor := NewPipeProcessor()
//	ctx := processor.Process(originalCtx, "uppercase|trim|minify")
func NewPipeProcessor() PipeProcessor {
	processor := &pipeProcessor{
		handlers: make([]PipeHandler, 0),
	}
	
	// Add default handlers
	processor.AddHandler(&TransformerPipeHandler{})
	processor.AddHandler(&PathPipeHandler{})
	processor.AddHandler(&QueryPipeHandler{})
	
	return processor
}

// Process applies pipe operations to the context using the chain of handlers.
// This method implements the Chain of Responsibility pattern.
//
// Parameters:
//   - ctx: The Context to process
//   - pipe: The pipe operation string
//
// Returns:
//   - Context: The processed Context
//
// Example Usage:
//
//	processor := NewPipeProcessor()
//	result := processor.Process(ctx, "uppercase|pretty")
func (pp *pipeProcessor) Process(ctx Context, pipe string) Context {
	if pipe == "" {
		return ctx
	}
	
	operations := strings.Split(pipe, "|")
	result := ctx
	
	for _, operation := range operations {
		operation = strings.TrimSpace(operation)
		if operation == "" {
			continue
		}
		
		// Extract operation name and arguments
		opName, args := parseOperation(operation)
		
		// Try each handler in the chain
		handled := false
		for _, handler := range pp.handlers {
			if handler.CanHandle(opName) {
				var ok bool
				result, ok = handler.Handle(result, opName, args)
				if ok {
					handled = true
					break
				}
			}
		}
		
		// If no handler processed the operation, skip it
		if !handled {
			continue
		}
	}
	
	return result
}

// AddHandler adds a new pipe handler to the chain.
// This method allows extending the processor with custom handlers.
//
// Parameters:
//   - handler: The pipe handler to add
//
// Returns:
//   - PipeProcessor: The processor instance for method chaining
//
// Example Usage:
//
//	processor := NewPipeProcessor().AddHandler(customHandler)
func (pp *pipeProcessor) AddHandler(handler PipeHandler) PipeProcessor {
	if handler != nil {
		pp.handlers = append(pp.handlers, handler)
	}
	return pp
}

// TransformerPipeHandler handles transformer-based pipe operations.
type TransformerPipeHandler struct{}

// Handle processes transformer operations.
func (tph *TransformerPipeHandler) Handle(ctx Context, operation string, arg string) (Context, bool) {
	// Remove @ prefix if present
	transformerName := operation
	if strings.HasPrefix(operation, "@") {
		transformerName = operation[1:]
	}
	
	registry := GetTransformerRegistry()
	transformer, exists := registry.Get(transformerName)
	if !exists {
		return ctx, false
	}
	
	// Apply transformation to the JSON content
	transformed := transformer.Transform(ctx.Unprocessed(), arg)
	
	// Create new context with transformed data
	return Parse(transformed), true
}

// CanHandle checks if this handler can process the given operation.
func (tph *TransformerPipeHandler) CanHandle(operation string) bool {
	transformerName := operation
	if strings.HasPrefix(operation, "@") {
		transformerName = operation[1:]
	}
	
	registry := GetTransformerRegistry()
	return registry.Exists(transformerName)
}

// PathPipeHandler handles path-based pipe operations.
type PathPipeHandler struct{}

// Handle processes path operations (e.g., extracting nested values).
func (pph *PathPipeHandler) Handle(ctx Context, operation string, arg string) (Context, bool) {
	// Handle path operations like "field.subfield"
	if !strings.Contains(operation, ".") && !strings.Contains(operation, "[") {
		// Simple field access
		return ctx.Get(operation), true
	}
	
	// Complex path access
	return ctx.Get(operation), true
}

// CanHandle checks if this handler can process path operations.
func (pph *PathPipeHandler) CanHandle(operation string) bool {
	// Handle operations that look like paths (contain dots or brackets)
	// but are not transformer operations (don't start with @)
	return !strings.HasPrefix(operation, "@") && 
		   (strings.Contains(operation, ".") || 
		    strings.Contains(operation, "[") || 
		    isValidFieldName(operation))
}

// QueryPipeHandler handles query-based pipe operations.
type QueryPipeHandler struct{}

// Handle processes query operations.
func (qph *QueryPipeHandler) Handle(ctx Context, operation string, arg string) (Context, bool) {
	// Handle query operations like "#(name==john)" 
	if strings.HasPrefix(operation, "#") {
		return ctx.Get(operation), true
	}
	
	return ctx, false
}

// CanHandle checks if this handler can process query operations.
func (qph *QueryPipeHandler) CanHandle(operation string) bool {
	return strings.HasPrefix(operation, "#")
}

// parseOperation extracts the operation name and arguments from a pipe operation.
// This function handles operations like "transformer:{args}" or "transformer:arg".
//
// Parameters:
//   - operation: The pipe operation string
//
// Returns:
//   - string: The operation name
//   - string: The operation arguments
//
// Example Usage:
//
//	name, args := parseOperation("replace:{\"target\":\"old\",\"replacement\":\"new\"}")
func parseOperation(operation string) (string, string) {
	if !strings.Contains(operation, ":") {
		return operation, ""
	}
	
	parts := strings.SplitN(operation, ":", 2)
	return parts[0], parts[1]
}

// isValidFieldName checks if a string is a valid field name for path operations.
func isValidFieldName(name string) bool {
	if name == "" {
		return false
	}
	
	// Check if it contains characters that would make it a transformer
	if strings.HasPrefix(name, "@") || strings.Contains(name, "{") {
		return false
	}
	
	// Simple validation - field names should be reasonable identifiers
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || 
			 (r >= 'A' && r <= 'Z') || 
			 (r >= '0' && r <= '9') || 
			 r == '_' || r == '-') {
			return false
		}
	}
	
	return true
}

// DefaultPipeProcessor is the default instance of pipe processor.
var DefaultPipeProcessor = NewPipeProcessor()

// ProcessPipe is a convenience function that uses the default pipe processor.
// This function provides a simple way to apply pipe operations to a Context.
//
// Parameters:
//   - ctx: The Context to process
//   - pipe: The pipe operation string
//
// Returns:
//   - Context: The processed Context
//
// Example Usage:
//
//	result := ProcessPipe(ctx, "uppercase|trim|minify")
func ProcessPipe(ctx Context, pipe string) Context {
	return DefaultPipeProcessor.Process(ctx, pipe)
}

// ChainedPipeProcessor allows building a custom pipe processor with specific handlers.
type ChainedPipeProcessor struct {
	*pipeProcessor
}

// NewChainedPipeProcessor creates a new chained pipe processor without default handlers.
// This allows building a completely custom processing chain.
//
// Returns:
//   - *ChainedPipeProcessor: A new chained processor instance
//
// Example Usage:
//
//	processor := NewChainedPipeProcessor().
//		WithTransformers().
//		WithPaths().
//		WithCustomHandler(myHandler)
func NewChainedPipeProcessor() *ChainedPipeProcessor {
	return &ChainedPipeProcessor{
		pipeProcessor: &pipeProcessor{
			handlers: make([]PipeHandler, 0),
		},
	}
}

// WithTransformers adds transformer handling to the processor.
func (cpp *ChainedPipeProcessor) WithTransformers() *ChainedPipeProcessor {
	cpp.AddHandler(&TransformerPipeHandler{})
	return cpp
}

// WithPaths adds path handling to the processor.
func (cpp *ChainedPipeProcessor) WithPaths() *ChainedPipeProcessor {
	cpp.AddHandler(&PathPipeHandler{})
	return cpp
}

// WithQueries adds query handling to the processor.
func (cpp *ChainedPipeProcessor) WithQueries() *ChainedPipeProcessor {
	cpp.AddHandler(&QueryPipeHandler{})
	return cpp
}

// WithCustomHandler adds a custom handler to the processor.
func (cpp *ChainedPipeProcessor) WithCustomHandler(handler PipeHandler) *ChainedPipeProcessor {
	cpp.AddHandler(handler)
	return cpp
}

// Build returns the final PipeProcessor.
func (cpp *ChainedPipeProcessor) Build() PipeProcessor {
	return cpp.pipeProcessor
}