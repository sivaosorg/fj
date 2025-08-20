package fj

import (
	"strings"
)

// TransformerStrategy implements the Strategy pattern for JSON transformations.
// This allows different transformation algorithms to be swapped at runtime.
type TransformerStrategy struct {
	name        string
	transformer Transformer
}

// NewTransformerStrategy creates a new transformer strategy.
// This function provides a convenient way to create strategy instances.
//
// Parameters:
//   - name: The name of the transformer to use
//
// Returns:
//   - *TransformerStrategy: A new strategy instance, nil if transformer not found
//
// Example Usage:
//
//	strategy := NewTransformerStrategy("uppercase")
//	if strategy != nil {
//		result := strategy.Execute(`{"name": "john"}`, "")
//	}
func NewTransformerStrategy(name string) *TransformerStrategy {
	registry := GetTransformerRegistry()
	transformer, exists := registry.Get(name)
	if !exists {
		return nil
	}
	
	return &TransformerStrategy{
		name:        name,
		transformer: transformer,
	}
}

// Execute applies the transformation strategy to the given JSON.
// This method encapsulates the transformation logic within the strategy.
//
// Parameters:
//   - json: The JSON string to transform
//   - arg: Additional arguments for the transformation
//
// Returns:
//   - string: The transformed JSON string
//
// Example Usage:
//
//	strategy := NewTransformerStrategy("uppercase")
//	result := strategy.Execute(`{"name": "john"}`, "")
func (ts *TransformerStrategy) Execute(json, arg string) string {
	if ts.transformer == nil {
		return json
	}
	return ts.transformer.Transform(json, arg)
}

// Name returns the name of the strategy.
func (ts *TransformerStrategy) Name() string {
	return ts.name
}

// IsValid checks if the strategy is valid and ready to use.
func (ts *TransformerStrategy) IsValid() bool {
	return ts.transformer != nil
}

// TransformerStrategyChain implements the Chain of Responsibility pattern
// for applying multiple transformer strategies in sequence.
type TransformerStrategyChain struct {
	strategies []*TransformerStrategy
}

// NewTransformerStrategyChain creates a new transformer strategy chain.
// This function initializes an empty chain ready for adding strategies.
//
// Returns:
//   - *TransformerStrategyChain: A new strategy chain instance
//
// Example Usage:
//
//	chain := NewTransformerStrategyChain().
//		Add("uppercase").
//		Add("trim")
//	result := chain.Execute(`{"name": " john "}`, "")
func NewTransformerStrategyChain() *TransformerStrategyChain {
	return &TransformerStrategyChain{
		strategies: make([]*TransformerStrategy, 0),
	}
}

// Add adds a transformer strategy to the chain.
// This method enables building a chain of transformations.
//
// Parameters:
//   - name: The name of the transformer to add to the chain
//
// Returns:
//   - *TransformerStrategyChain: The chain instance for method chaining
//
// Example Usage:
//
//	chain := NewTransformerStrategyChain().
//		Add("trim").
//		Add("uppercase").
//		Add("minify")
func (tsc *TransformerStrategyChain) Add(name string) *TransformerStrategyChain {
	strategy := NewTransformerStrategy(name)
	if strategy != nil {
		tsc.strategies = append(tsc.strategies, strategy)
	}
	return tsc
}

// AddStrategy adds a pre-created strategy to the chain.
// This method allows adding custom strategy instances.
//
// Parameters:
//   - strategy: The transformer strategy to add
//
// Returns:
//   - *TransformerStrategyChain: The chain instance for method chaining
//
// Example Usage:
//
//	customStrategy := NewTransformerStrategy("custom")
//	chain := NewTransformerStrategyChain().AddStrategy(customStrategy)
func (tsc *TransformerStrategyChain) AddStrategy(strategy *TransformerStrategy) *TransformerStrategyChain {
	if strategy != nil && strategy.IsValid() {
		tsc.strategies = append(tsc.strategies, strategy)
	}
	return tsc
}

// Execute applies all strategies in the chain sequentially.
// This method implements the Chain of Responsibility pattern.
//
// Parameters:
//   - json: The JSON string to transform
//   - args: Optional arguments for transformations (applied to all)
//
// Returns:
//   - string: The final transformed JSON string
//
// Example Usage:
//
//	chain := NewTransformerStrategyChain().Add("trim").Add("uppercase")
//	result := chain.Execute(`{"name": " john "}`, "")
func (tsc *TransformerStrategyChain) Execute(json string, args ...string) string {
	result := json
	arg := ""
	if len(args) > 0 {
		arg = args[0]
	}
	
	for _, strategy := range tsc.strategies {
		result = strategy.Execute(result, arg)
	}
	
	return result
}

// ExecuteWithArgs applies strategies with individual arguments.
// This method allows different arguments for each strategy in the chain.
//
// Parameters:
//   - json: The JSON string to transform
//   - args: Arguments for each strategy (matched by index)
//
// Returns:
//   - string: The final transformed JSON string
//
// Example Usage:
//
//	chain := NewTransformerStrategyChain().Add("replace").Add("uppercase")
//	result := chain.ExecuteWithArgs(`{"name": "john"}`, 
//		`{"target": "john", "replacement": "jane"}`, "")
func (tsc *TransformerStrategyChain) ExecuteWithArgs(json string, args ...string) string {
	result := json
	
	for i, strategy := range tsc.strategies {
		arg := ""
		if i < len(args) {
			arg = args[i]
		}
		result = strategy.Execute(result, arg)
	}
	
	return result
}

// Length returns the number of strategies in the chain.
func (tsc *TransformerStrategyChain) Length() int {
	return len(tsc.strategies)
}

// Clear removes all strategies from the chain.
func (tsc *TransformerStrategyChain) Clear() *TransformerStrategyChain {
	tsc.strategies = tsc.strategies[:0]
	return tsc
}

// Names returns the names of all strategies in the chain.
func (tsc *TransformerStrategyChain) Names() []string {
	names := make([]string, len(tsc.strategies))
	for i, strategy := range tsc.strategies {
		names[i] = strategy.Name()
	}
	return names
}

// ParseTransformerChain creates a transformer chain from a pipe string.
// This function parses pipe operations and creates appropriate strategies.
//
// Parameters:
//   - pipeStr: A pipe string like "trim|uppercase|minify"
//
// Returns:
//   - *TransformerStrategyChain: A configured strategy chain
//
// Example Usage:
//
//	chain := ParseTransformerChain("trim|uppercase|pretty")
//	result := chain.Execute(`{"name": " john "}`, "")
func ParseTransformerChain(pipeStr string) *TransformerStrategyChain {
	chain := NewTransformerStrategyChain()
	
	if pipeStr == "" {
		return chain
	}
	
	// Split by pipe character
	parts := strings.Split(pipeStr, "|")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		
		// Handle transformer with arguments: transformer:{args}
		transformerName := part
		if strings.Contains(part, ":") {
			transformerName = strings.Split(part, ":")[0]
		}
		
		// Remove @ prefix if present
		if strings.HasPrefix(transformerName, "@") {
			transformerName = transformerName[1:]
		}
		
		chain.Add(transformerName)
	}
	
	return chain
}

// ApplyTransformerChain is a convenience function that applies a transformer chain to JSON.
// This function provides a simple way to apply multiple transformations.
//
// Parameters:
//   - json: The JSON string to transform
//   - pipeStr: A pipe string defining the transformation chain
//
// Returns:
//   - string: The transformed JSON string
//
// Example Usage:
//
//	result := ApplyTransformerChain(`{"name": " john "}`, "trim|uppercase|minify")
func ApplyTransformerChain(json, pipeStr string) string {
	chain := ParseTransformerChain(pipeStr)
	return chain.Execute(json)
}