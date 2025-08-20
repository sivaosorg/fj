package fj

import (
	"strings"
	"sync"
)

// QueryStrategy defines the interface for different query execution strategies
type QueryStrategy interface {
	// CanExecute determines if this strategy can handle the given query
	CanExecute(query *Query) bool
	
	// Execute runs the query using this strategy
	Execute(ctx Context, query *Query) Context
	
	// Name returns the name of this strategy for debugging/logging
	Name() string
}

// Query represents a parsed and optimized query ready for execution
type Query struct {
	Path        string            `json:"path"`        // The main path to query
	Filters     []Filter          `json:"filters"`     // Filter conditions
	Selections  []string          `json:"selections"`  // Fields to select
	OrderBy     string            `json:"orderBy"`     // Field to order by
	Limit       int               `json:"limit"`       // Maximum results to return
	Offset      int               `json:"offset"`      // Number of results to skip
	Options     map[string]interface{} `json:"options"` // Additional options
	
	// Internal fields for optimization
	parsedPath  *metadata
	handlerChain *HandlerChain
}

// Filter represents a filter condition in a query
type Filter struct {
	Field    string      `json:"field"`    // Field to filter on
	Operator string      `json:"operator"` // Comparison operator (=, !=, <, >, <=, >=, %)
	Value    interface{} `json:"value"`    // Value to compare against
}

// Execute runs the query against the given JSON context
func (q *Query) Execute(ctx Context) Context {
	processor := GetQueryProcessor()
	return processor.ExecuteQuery(ctx, q)
}

// QueryProcessor manages different query strategies and executes queries
type QueryProcessor struct {
	strategies []QueryStrategy
	mu         sync.RWMutex
}

var (
	queryProcessor *QueryProcessor
	queryProcOnce  sync.Once
)

// GetQueryProcessor returns the singleton query processor instance
func GetQueryProcessor() *QueryProcessor {
	queryProcOnce.Do(func() {
		queryProcessor = &QueryProcessor{
			strategies: []QueryStrategy{
				&ConditionalStrategy{},
				&ArrayStrategy{},
				&WildcardStrategy{}, 
				&SimplePathStrategy{},
			},
		}
	})
	return queryProcessor
}

// RegisterStrategy adds a new query strategy to the processor
func (p *QueryProcessor) RegisterStrategy(strategy QueryStrategy) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.strategies = append(p.strategies, strategy)
}

// ExecuteQuery executes a query using the appropriate strategy
func (p *QueryProcessor) ExecuteQuery(ctx Context, query *Query) Context {
	p.mu.RLock()
	strategies := make([]QueryStrategy, len(p.strategies))
	copy(strategies, p.strategies)
	p.mu.RUnlock()
	
	// Find the first strategy that can handle this query
	for _, strategy := range strategies {
		if strategy.CanExecute(query) {
			return strategy.Execute(ctx, query)
		}
	}
	
	// Fallback to simple path strategy if no other strategy matches
	fallback := &SimplePathStrategy{}
	return fallback.Execute(ctx, query)
}

// SimplePathStrategy handles basic property access queries
type SimplePathStrategy struct{}

func (s *SimplePathStrategy) Name() string {
	return "SimplePathStrategy"
}

func (s *SimplePathStrategy) CanExecute(query *Query) bool {
	// Handle simple paths without complex features
	path := query.Path
	return !strings.Contains(path, "*") &&
		   !strings.Contains(path, "?") &&
		   !strings.Contains(path, "[") &&
		   len(query.Filters) == 0
}

func (s *SimplePathStrategy) Execute(ctx Context, query *Query) Context {
	// Use existing Get functionality for simple paths
	result := ctx.Get(query.Path)
	
	// Apply selections if specified
	if len(query.Selections) > 0 {
		return s.applySelections(result, query.Selections)
	}
	
	return result
}

func (s *SimplePathStrategy) applySelections(ctx Context, selections []string) Context {
	if ctx.kind != JSON {
		return ctx
	}
	
	// If only one selection, return that field directly
	if len(selections) == 1 {
		return ctx.Get(selections[0])
	}
	
	// Multiple selections - create new object with selected fields
	results := make(map[string]Context)
	for _, field := range selections {
		value := ctx.Get(field)
		if value.Exists() {
			results[field] = value
		}
	}
	
	// Convert map to JSON string
	if len(results) == 0 {
		return Context{}
	}
	
	// Build JSON object string
	var builder strings.Builder
	builder.WriteString("{")
	first := true
	for key, value := range results {
		if !first {
			builder.WriteString(",")
		}
		first = false
		builder.WriteString("\"" + key + "\":")
		if value.kind == String {
			builder.WriteString("\"" + value.strings + "\"")
		} else {
			builder.WriteString(value.String())
		}
	}
	builder.WriteString("}")
	
	return Parse(builder.String())
}

// WildcardStrategy handles wildcard queries with * and ? operators  
type WildcardStrategy struct{}

func (s *WildcardStrategy) Name() string {
	return "WildcardStrategy"
}

func (s *WildcardStrategy) CanExecute(query *Query) bool {
	path := query.Path
	return strings.Contains(path, "*") || strings.Contains(path, "?")
}

func (s *WildcardStrategy) Execute(ctx Context, query *Query) Context {
	// Use existing Get functionality for wildcard processing
	result := ctx.Get(query.Path)
	
	// Apply filters if specified
	if len(query.Filters) > 0 {
		result = s.applyFilters(result, query.Filters)
	}
	
	// Apply selections if specified
	if len(query.Selections) > 0 {
		result = s.applySelections(result, query.Selections)
	}
	
	return result
}

func (s *WildcardStrategy) applyFilters(ctx Context, filters []Filter) Context {
	// For wildcard results, we need to filter array elements
	if !ctx.IsArray() {
		return ctx
	}
	
	// Get array elements and filter them
	array := ctx.Array()
	var filtered []Context
	
	for _, item := range array {
		if s.matchesFilters(item, filters) {
			filtered = append(filtered, item)
		}
	}
	
	// Convert back to JSON array
	return s.arrayToContext(filtered)
}

func (s *WildcardStrategy) applySelections(ctx Context, selections []string) Context {
	if !ctx.IsArray() {
		return ctx
	}
	
	array := ctx.Array()
	var results []Context
	
	for _, item := range array {
		if len(selections) == 1 {
			results = append(results, item.Get(selections[0]))
		} else {
			// Multiple selections - create object with selected fields
			selected := s.selectFields(item, selections)
			if selected.Exists() {
				results = append(results, selected)
			}
		}
	}
	
	return s.arrayToContext(results)
}

func (s *WildcardStrategy) selectFields(ctx Context, selections []string) Context {
	if ctx.kind != JSON {
		return Context{}
	}
	
	var builder strings.Builder
	builder.WriteString("{")
	first := true
	
	for _, field := range selections {
		value := ctx.Get(field)
		if value.Exists() {
			if !first {
				builder.WriteString(",")
			}
			first = false
			builder.WriteString("\"" + field + "\":")
			if value.kind == String {
				builder.WriteString("\"" + value.strings + "\"")
			} else {
				builder.WriteString(value.String())
			}
		}
	}
	
	builder.WriteString("}")
	if first { // No fields were added
		return Context{}
	}
	
	return Parse(builder.String())
}

func (s *WildcardStrategy) matchesFilters(ctx Context, filters []Filter) bool {
	for _, filter := range filters {
		if !s.matchesFilter(ctx, filter) {
			return false
		}
	}
	return true
}

func (s *WildcardStrategy) matchesFilter(ctx Context, filter Filter) bool {
	fieldValue := ctx.Get(filter.Field)
	if !fieldValue.Exists() {
		return false
	}
	
	// Convert filter value to string for comparison
	filterValueStr := ""
	switch v := filter.Value.(type) {
	case string:
		filterValueStr = v
	case int:
		filterValueStr = string(rune(v + '0'))
	case float64:
		filterValueStr = fieldValue.String() // Use field's string representation
	default:
		filterValueStr = fieldValue.String()
	}
	
	switch filter.Operator {
	case "=", "==":
		return fieldValue.String() == filterValueStr
	case "!=":
		return fieldValue.String() != filterValueStr
	case "<":
		return fieldValue.Numeric() < toFloat(filterValueStr)
	case ">":
		return fieldValue.Numeric() > toFloat(filterValueStr)
	case "<=":
		return fieldValue.Numeric() <= toFloat(filterValueStr)
	case ">=":
		return fieldValue.Numeric() >= toFloat(filterValueStr)
	case "%":
		return strings.Contains(fieldValue.String(), filterValueStr)
	default:
		return false
	}
}

func (s *WildcardStrategy) arrayToContext(array []Context) Context {
	if len(array) == 0 {
		return Context{}
	}
	
	var builder strings.Builder
	builder.WriteString("[")
	
	for i, item := range array {
		if i > 0 {
			builder.WriteString(",")
		}
		if item.kind == String {
			builder.WriteString("\"" + item.strings + "\"")
		} else {
			builder.WriteString(item.String())
		}
	}
	
	builder.WriteString("]")
	return Parse(builder.String())
}

// ArrayStrategy handles array-specific queries
type ArrayStrategy struct{}

func (s *ArrayStrategy) Name() string {
	return "ArrayStrategy"
}

func (s *ArrayStrategy) CanExecute(query *Query) bool {
	path := query.Path
	return strings.Contains(path, "[") && strings.Contains(path, "]")
}

func (s *ArrayStrategy) Execute(ctx Context, query *Query) Context {
	// Use existing Get functionality for array processing
	result := ctx.Get(query.Path)
	
	// Apply filters if specified
	if len(query.Filters) > 0 {
		wildcard := &WildcardStrategy{}
		result = wildcard.applyFilters(result, query.Filters)
	}
	
	// Apply selections if specified
	if len(query.Selections) > 0 {
		wildcard := &WildcardStrategy{}
		result = wildcard.applySelections(result, query.Selections)
	}
	
	return result
}

// ConditionalStrategy handles conditional queries with filters
type ConditionalStrategy struct{}

func (s *ConditionalStrategy) Name() string {
	return "ConditionalStrategy"
}

func (s *ConditionalStrategy) CanExecute(query *Query) bool {
	// Handle queries with explicit filters or conditional syntax in path
	if len(query.Filters) > 0 {
		return true
	}
	
	path := query.Path
	return strings.Contains(path, "[") && 
		   (strings.Contains(path, "=") || strings.Contains(path, ">") || 
		    strings.Contains(path, "<") || strings.Contains(path, "%"))
}

func (s *ConditionalStrategy) Execute(ctx Context, query *Query) Context {
	// Use existing Get functionality for conditional processing
	result := ctx.Get(query.Path)
	
	// Apply additional filters if specified
	if len(query.Filters) > 0 {
		wildcard := &WildcardStrategy{}
		result = wildcard.applyFilters(result, query.Filters)
	}
	
	// Apply selections if specified
	if len(query.Selections) > 0 {
		wildcard := &WildcardStrategy{}
		result = wildcard.applySelections(result, query.Selections)
	}
	
	return result
}

// Helper function to convert string to float
func toFloat(s string) float64 {
	return simpleFloatParse(s)
}