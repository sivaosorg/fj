package fj

import (
	"sync"
)

// QueryBuilder provides a fluent interface for building queries
type QueryBuilder struct {
	query *Query
	mu    sync.Mutex
}

// NewQueryBuilder creates a new query builder
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		query: &Query{
			Filters:   make([]Filter, 0),
			Selections: make([]string, 0),
			Options:   make(map[string]interface{}),
			Limit:     -1, // No limit by default
			Offset:    0,
		},
	}
}

// Path sets the base path for the query
func (b *QueryBuilder) Path(path string) *QueryBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	b.query.Path = path
	return b
}

// Filter adds a filter condition to the query
func (b *QueryBuilder) Filter(field, operator string, value interface{}) *QueryBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	filter := Filter{
		Field:    field,
		Operator: operator,
		Value:    value,
	}
	
	b.query.Filters = append(b.query.Filters, filter)
	return b
}

// Where is an alias for Filter for more natural language
func (b *QueryBuilder) Where(field, operator string, value interface{}) *QueryBuilder {
	return b.Filter(field, operator, value)
}

// Equals adds an equality filter (shorthand for Filter with "=" operator)
func (b *QueryBuilder) Equals(field string, value interface{}) *QueryBuilder {
	return b.Filter(field, "=", value)
}

// NotEquals adds a not-equals filter (shorthand for Filter with "!=" operator)  
func (b *QueryBuilder) NotEquals(field string, value interface{}) *QueryBuilder {
	return b.Filter(field, "!=", value)
}

// GreaterThan adds a greater-than filter (shorthand for Filter with ">" operator)
func (b *QueryBuilder) GreaterThan(field string, value interface{}) *QueryBuilder {
	return b.Filter(field, ">", value)
}

// LessThan adds a less-than filter (shorthand for Filter with "<" operator)
func (b *QueryBuilder) LessThan(field string, value interface{}) *QueryBuilder {
	return b.Filter(field, "<", value)
}

// GreaterThanOrEqual adds a >= filter (shorthand for Filter with ">=" operator)
func (b *QueryBuilder) GreaterThanOrEqual(field string, value interface{}) *QueryBuilder {
	return b.Filter(field, ">=", value)
}

// LessThanOrEqual adds a <= filter (shorthand for Filter with "<=" operator)
func (b *QueryBuilder) LessThanOrEqual(field string, value interface{}) *QueryBuilder {
	return b.Filter(field, "<=", value)
}

// Contains adds a contains filter (shorthand for Filter with "%" operator)
func (b *QueryBuilder) Contains(field string, value interface{}) *QueryBuilder {
	return b.Filter(field, "%", value)
}

// Select adds fields to select from the query results
func (b *QueryBuilder) Select(fields ...string) *QueryBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	b.query.Selections = append(b.query.Selections, fields...)
	return b
}

// OrderBy sets the field to order results by
func (b *QueryBuilder) OrderBy(field string) *QueryBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	b.query.OrderBy = field
	return b
}

// Limit sets the maximum number of results to return
func (b *QueryBuilder) Limit(limit int) *QueryBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	b.query.Limit = limit
	return b
}

// Offset sets the number of results to skip
func (b *QueryBuilder) Offset(offset int) *QueryBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	b.query.Offset = offset
	return b
}

// Option adds a custom option to the query
func (b *QueryBuilder) Option(key string, value interface{}) *QueryBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	b.query.Options[key] = value
	return b
}

// Build finalizes the query and returns an executable Query object
func (b *QueryBuilder) Build() *Query {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	// Create a copy of the query to prevent further modification
	builtQuery := &Query{
		Path:        b.query.Path,
		Filters:     make([]Filter, len(b.query.Filters)),
		Selections:  make([]string, len(b.query.Selections)),
		OrderBy:     b.query.OrderBy,
		Limit:       b.query.Limit,
		Offset:      b.query.Offset,
		Options:     make(map[string]interface{}),
	}
	
	// Deep copy filters
	copy(builtQuery.Filters, b.query.Filters)
	
	// Deep copy selections  
	copy(builtQuery.Selections, b.query.Selections)
	
	// Deep copy options
	for k, v := range b.query.Options {
		builtQuery.Options[k] = v
	}
	
	// Pre-parse and optimize the query
	builtQuery.optimize()
	
	return builtQuery
}

// BuildAndExecute builds the query and immediately executes it against the given context
func (b *QueryBuilder) BuildAndExecute(ctx Context) Context {
	query := b.Build()
	return query.Execute(ctx)
}

// optimize prepares the query for efficient execution
func (q *Query) optimize() {
	// Parse the path using existing analyzePath functionality
	if !isEmptyString(q.Path) {
		q.parsedPath = &metadata{}
		*q.parsedPath = analyzePath(q.Path)
	}
	
	// Initialize handler chain for path processing
	q.handlerChain = NewHandlerChain()
	
	// Additional optimizations could be added here:
	// - Index usage hints
	// - Query plan generation  
	// - Caching strategy selection
}

// QueryBuilderPool provides a pool of reusable QueryBuilder instances for performance
type QueryBuilderPool struct {
	pool sync.Pool
}

// NewQueryBuilderPool creates a new query builder pool
func NewQueryBuilderPool() *QueryBuilderPool {
	return &QueryBuilderPool{
		pool: sync.Pool{
			New: func() interface{} {
				return NewQueryBuilder()
			},
		},
	}
}

// Get retrieves a QueryBuilder from the pool
func (p *QueryBuilderPool) Get() *QueryBuilder {
	builder := p.pool.Get().(*QueryBuilder)
	
	// Reset the builder state
	builder.mu.Lock()
	builder.query = &Query{
		Filters:   make([]Filter, 0),
		Selections: make([]string, 0), 
		Options:   make(map[string]interface{}),
		Limit:     -1,
		Offset:    0,
	}
	builder.mu.Unlock()
	
	return builder
}

// Put returns a QueryBuilder to the pool
func (p *QueryBuilderPool) Put(builder *QueryBuilder) {
	if builder != nil {
		p.pool.Put(builder)
	}
}

// Global query builder pool for performance
var (
	globalBuilderPool *QueryBuilderPool
	builderPoolOnce   sync.Once
)

// GetQueryBuilderFromPool returns a QueryBuilder from the global pool
func GetQueryBuilderFromPool() *QueryBuilder {
	builderPoolOnce.Do(func() {
		globalBuilderPool = NewQueryBuilderPool()
	})
	return globalBuilderPool.Get()
}

// PutQueryBuilderToPool returns a QueryBuilder to the global pool
func PutQueryBuilderToPool(builder *QueryBuilder) {
	if globalBuilderPool != nil {
		globalBuilderPool.Put(builder)
	}
}

// FluentQueryBuilder provides additional fluent methods for common patterns
type FluentQueryBuilder struct {
	*QueryBuilder
}

// NewFluentQueryBuilder creates a new fluent query builder with enhanced methods
func NewFluentQueryBuilder() *FluentQueryBuilder {
	return &FluentQueryBuilder{
		QueryBuilder: NewQueryBuilder(),
	}
}

// FromPath creates a query builder starting with the given path (alias for Path)
func FromPath(path string) *QueryBuilder {
	return NewQueryBuilder().Path(path)
}

// FromPathFluent creates a fluent query builder starting with the given path
func FromPathFluent(path string) *FluentQueryBuilder {
	return NewFluentQueryBuilder().Path(path)
}

// Path overrides the base Path method to return FluentQueryBuilder
func (b *FluentQueryBuilder) Path(path string) *FluentQueryBuilder {
	b.QueryBuilder.Path(path)
	return b
}

// Filter overrides the base Filter method to return FluentQueryBuilder  
func (b *FluentQueryBuilder) Filter(field, operator string, value interface{}) *FluentQueryBuilder {
	b.QueryBuilder.Filter(field, operator, value)
	return b
}

// Where overrides the base Where method to return FluentQueryBuilder
func (b *FluentQueryBuilder) Where(field, operator string, value interface{}) *FluentQueryBuilder {
	b.QueryBuilder.Where(field, operator, value)
	return b
}

// Select overrides the base Select method to return FluentQueryBuilder
func (b *FluentQueryBuilder) Select(fields ...string) *FluentQueryBuilder {
	b.QueryBuilder.Select(fields...)
	return b
}

// OrderBy overrides the base OrderBy method to return FluentQueryBuilder
func (b *FluentQueryBuilder) OrderBy(field string) *FluentQueryBuilder {
	b.QueryBuilder.OrderBy(field)
	return b
}

// Limit overrides the base Limit method to return FluentQueryBuilder
func (b *FluentQueryBuilder) Limit(limit int) *FluentQueryBuilder {
	b.QueryBuilder.Limit(limit)
	return b
}

// GreaterThanOrEqual overrides the base GreaterThanOrEqual method to return FluentQueryBuilder
func (b *FluentQueryBuilder) GreaterThanOrEqual(field string, value interface{}) *FluentQueryBuilder {
	b.QueryBuilder.GreaterThanOrEqual(field, value)
	return b
}

// LessThanOrEqual overrides the base LessThanOrEqual method to return FluentQueryBuilder  
func (b *FluentQueryBuilder) LessThanOrEqual(field string, value interface{}) *FluentQueryBuilder {
	b.QueryBuilder.LessThanOrEqual(field, value)
	return b
}

// Contains overrides the base Contains method to return FluentQueryBuilder
func (b *FluentQueryBuilder) Contains(field string, value interface{}) *FluentQueryBuilder {
	b.QueryBuilder.Contains(field, value)
	return b
}

// Equals overrides the base Equals method to return FluentQueryBuilder
func (b *FluentQueryBuilder) Equals(field string, value interface{}) *FluentQueryBuilder {
	b.QueryBuilder.Equals(field, value)
	return b
}

// NotEquals overrides the base NotEquals method to return FluentQueryBuilder
func (b *FluentQueryBuilder) NotEquals(field string, value interface{}) *FluentQueryBuilder {
	b.QueryBuilder.NotEquals(field, value)
	return b
}

// GreaterThan overrides the base GreaterThan method to return FluentQueryBuilder
func (b *FluentQueryBuilder) GreaterThan(field string, value interface{}) *FluentQueryBuilder {
	b.QueryBuilder.GreaterThan(field, value)
	return b
}

// LessThan overrides the base LessThan method to return FluentQueryBuilder
func (b *FluentQueryBuilder) LessThan(field string, value interface{}) *FluentQueryBuilder {
	b.QueryBuilder.LessThan(field, value)
	return b
}

// IfCondition applies the given function only if the condition is true
func (b *FluentQueryBuilder) IfCondition(condition bool, fn func(*FluentQueryBuilder) *FluentQueryBuilder) *FluentQueryBuilder {
	if condition {
		return fn(b)
	}
	return b
}

// WhenEmpty applies the given function only if the path is empty
func (b *FluentQueryBuilder) WhenEmpty(fn func(*FluentQueryBuilder) *FluentQueryBuilder) *FluentQueryBuilder {
	b.QueryBuilder.mu.Lock()
	isEmpty := b.QueryBuilder.query.Path == ""
	b.QueryBuilder.mu.Unlock()
	
	if isEmpty {
		return fn(b)
	}
	return b
}

// WhenNotEmpty applies the given function only if the path is not empty
func (b *FluentQueryBuilder) WhenNotEmpty(fn func(*FluentQueryBuilder) *FluentQueryBuilder) *FluentQueryBuilder {
	b.QueryBuilder.mu.Lock()
	isNotEmpty := b.QueryBuilder.query.Path != ""
	b.QueryBuilder.mu.Unlock()
	
	if isNotEmpty {
		return fn(b)
	}
	return b
}

// Pipe applies a custom transformation function to the builder
func (b *FluentQueryBuilder) Pipe(fn func(*FluentQueryBuilder) *FluentQueryBuilder) *FluentQueryBuilder {
	return fn(b)
}

// Chain applies multiple transformation functions in sequence
func (b *FluentQueryBuilder) Chain(fns ...func(*FluentQueryBuilder) *FluentQueryBuilder) *FluentQueryBuilder {
	result := b
	for _, fn := range fns {
		result = fn(result)
	}
	return result
}