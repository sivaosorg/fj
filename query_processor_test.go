package fj

import (
	"testing"
	"time"
)

// Test data for Phase 2 query processing
var testJSONData = `{
  "users": [
    {"name": "Alice", "age": 30, "email": "alice@example.com", "active": true},
    {"name": "Bob", "age": 25, "email": "bob@example.com", "active": false},
    {"name": "Charlie", "age": 35, "email": "charlie@example.com", "active": true},
    {"name": "Diana", "age": 28, "email": "diana@example.com", "active": true}
  ],
  "products": [
    {"title": "Laptop", "price": 999.99, "category": "Electronics", "inStock": true},
    {"title": "Book", "price": 19.99, "category": "Education", "inStock": false},
    {"title": "Phone", "price": 599.99, "category": "Electronics", "inStock": true}
  ],
  "config": {
    "database": {
      "host": "localhost",
      "port": 5432,
      "name": "testdb"
    },
    "api": {
      "version": "v1",
      "timeout": 30
    }
  }
}`

func TestPathHandlers(t *testing.T) {
	ctx := Parse(testJSONData)
	
	tests := []struct {
		name     string
		handler  PathHandler
		segment  string
		expected bool
	}{
		{
			name:     "ObjectPropertyHandler - simple property",
			handler:  &ObjectPropertyHandler{},
			segment:  "name",
			expected: true,
		},
		{
			name:     "ObjectPropertyHandler - nested property",
			handler:  &ObjectPropertyHandler{},
			segment:  "database.host",
			expected: true,
		},
		{
			name:     "ArrayElementHandler - array index",
			handler:  &ArrayElementHandler{},
			segment:  "users[0]",
			expected: true,
		},
		{
			name:     "WildcardHandler - wildcard",
			handler:  &WildcardHandler{},
			segment:  "users.*.name",
			expected: true,
		},
		{
			name:     "ConditionalHandler - conditional query",
			handler:  &ConditionalHandler{},
			segment:  "users[age>25]",
			expected: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			canHandle := tt.handler.CanHandle(tt.segment)
			if canHandle != tt.expected {
				t.Errorf("CanHandle() = %v, expected %v", canHandle, tt.expected)
			}
			
			if canHandle {
				parser := &parser{json: testJSONData, value: ctx}
				result, handled := tt.handler.Handle(ctx, tt.segment, parser)
				if !handled && tt.expected {
					t.Errorf("Handle() failed to handle segment: %s", tt.segment)
				}
				_ = result // Use result to avoid unused variable warning
			}
		})
	}
}

func TestHandlerChain(t *testing.T) {
	ctx := Parse(testJSONData)
	chain := NewHandlerChain()
	parser := &parser{json: testJSONData, value: ctx}
	
	tests := []struct {
		name     string
		segment  string
		expected bool
	}{
		{"Simple property", "users", true},
		{"Array access", "users[0]", true},
		{"Wildcard", "users.*", true},
		{"Conditional", "users[age>25]", true},
		{"Nested property", "config.database.host", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, handled := chain.Process(ctx, tt.segment, parser)
			if handled != tt.expected {
				t.Errorf("Process() handled = %v, expected %v for segment: %s", handled, tt.expected, tt.segment)
			}
			_ = result
		})
	}
}

func TestQueryStrategies(t *testing.T) {
	ctx := Parse(testJSONData)
	
	tests := []struct {
		name      string
		strategy  QueryStrategy
		query     *Query
		canExecute bool
	}{
		{
			name:     "SimplePathStrategy - basic path",
			strategy: &SimplePathStrategy{},
			query: &Query{
				Path: "config.database.host",
			},
			canExecute: true,
		},
		{
			name:     "WildcardStrategy - wildcard path", 
			strategy: &WildcardStrategy{},
			query: &Query{
				Path: "users.*.name",
			},
			canExecute: true,
		},
		{
			name:     "ArrayStrategy - array path",
			strategy: &ArrayStrategy{},
			query: &Query{
				Path: "users[0]",
			},
			canExecute: true,
		},
		{
			name:     "ConditionalStrategy - with filters",
			strategy: &ConditionalStrategy{},
			query: &Query{
				Path: "users",
				Filters: []Filter{
					{Field: "age", Operator: ">", Value: 25},
				},
			},
			canExecute: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			canExecute := tt.strategy.CanExecute(tt.query)
			if canExecute != tt.canExecute {
				t.Errorf("CanExecute() = %v, expected %v", canExecute, tt.canExecute)
			}
			
			if canExecute {
				result := tt.strategy.Execute(ctx, tt.query)
				_ = result // Use result to avoid unused variable warning
			}
		})
	}
}

func TestQueryBuilder(t *testing.T) {
	tests := []struct {
		name     string
		buildFn  func() *Query
		expected string
	}{
		{
			name: "Simple path query",
			buildFn: func() *Query {
				return NewQueryBuilder().
					Path("users").
					Build()
			},
			expected: "users",
		},
		{
			name: "Query with filters",
			buildFn: func() *Query {
				return NewQueryBuilder().
					Path("users").
					Filter("age", ">", 25).
					Filter("active", "=", true).
					Build()
			},
			expected: "users",
		},
		{
			name: "Query with selections",
			buildFn: func() *Query {
				return NewQueryBuilder().
					Path("users").
					Select("name", "email").
					Build()
			},
			expected: "users",
		},
		{
			name: "Complex query",
			buildFn: func() *Query {
				return NewQueryBuilder().
					Path("users").
					Filter("age", ">", 25).
					Select("name", "email").
					OrderBy("name").
					Limit(10).
					Build()
			},
			expected: "users",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := tt.buildFn()
			if query.Path != tt.expected {
				t.Errorf("Built query path = %v, expected %v", query.Path, tt.expected)
			}
		})
	}
}

func TestFluentQueryBuilder(t *testing.T) {
	ctx := Parse(testJSONData)
	
	// Test fluent API
	result := NewFluentQueryBuilder().
		Path("users").
		GreaterThan("age", 25).
		Select("name", "email").
		BuildAndExecute(ctx)
	
	if !result.Exists() {
		t.Error("Fluent query should return results")
	}
}

func TestQueryBuilderPool(t *testing.T) {
	pool := NewQueryBuilderPool()
	
	// Test getting builder from pool
	builder1 := pool.Get()
	if builder1 == nil {
		t.Error("Pool should return a valid builder")
	}
	
	// Test returning builder to pool
	pool.Put(builder1)
	
	// Test getting another builder
	builder2 := pool.Get()
	if builder2 == nil {
		t.Error("Pool should return a valid builder after put")
	}
}

func TestGlobalQueryBuilderPool(t *testing.T) {
	// Test global pool functions
	builder := GetQueryBuilderFromPool()
	if builder == nil {
		t.Error("Global pool should return a valid builder")
	}
	
	PutQueryBuilderToPool(builder)
}

func TestQueryExecution(t *testing.T) {
	ctx := Parse(testJSONData)
	
	tests := []struct {
		name     string
		query    *Query
		hasResult bool
	}{
		{
			name: "Simple path execution",
			query: NewQueryBuilder().
				Path("config.database.host").
				Build(),
			hasResult: true,
		},
		{
			name: "Array access execution",
			query: NewQueryBuilder().
				Path("users[0].name").
				Build(),
			hasResult: true,
		},
		{
			name: "Query with filters execution",
			query: NewQueryBuilder().
				Path("users").
				Filter("age", ">", 25).
				Build(),
			hasResult: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.query.Execute(ctx)
			hasResult := result.Exists()
			
			if hasResult != tt.hasResult {
				t.Errorf("Execute() hasResult = %v, expected %v", hasResult, tt.hasResult)
			}
		})
	}
}

func TestQueryProcessor(t *testing.T) {
	ctx := Parse(testJSONData)
	processor := GetQueryProcessor()
	
	query := &Query{
		Path: "users",
		Filters: []Filter{
			{Field: "age", Operator: ">", Value: 25},
		},
	}
	
	result := processor.ExecuteQuery(ctx, query)
	if !result.Exists() {
		t.Error("QueryProcessor should execute query successfully")
	}
}

func TestQueryCache(t *testing.T) {
	cache := NewQueryCache(10, time.Minute)
	testCtx := Parse(`{"test": "value"}`)
	
	// Test cache miss
	_, found := cache.Get("testkey")
	if found {
		t.Error("Cache should miss on empty cache")
	}
	
	// Test cache put and hit
	cache.Put("testkey", testCtx)
	result, found := cache.Get("testkey")
	if !found {
		t.Error("Cache should hit after putting value")
	}
	
	if result.String() != testCtx.String() {
		t.Error("Cached result should match original")
	}
	
	// Test cache stats
	stats := cache.Stats()
	if stats.Hits != 1 {
		t.Errorf("Cache hits = %v, expected 1", stats.Hits)
	}
	if stats.Misses != 1 {
		t.Errorf("Cache misses = %v, expected 1", stats.Misses)
	}
}

func TestPathCache(t *testing.T) {
	cache := NewPathCache(10)
	testMeta := &metadata{Part: "test", Path: "path"}
	
	// Test cache miss
	_, found := cache.Get("testpath")
	if found {
		t.Error("Path cache should miss on empty cache")
	}
	
	// Test cache put and hit
	cache.Put("testpath", testMeta)
	result, found := cache.Get("testpath")
	if !found {
		t.Error("Path cache should hit after putting value")
	}
	
	if result.Part != testMeta.Part {
		t.Error("Cached path metadata should match original")
	}
}

func TestQueryOptimizer(t *testing.T) {
	ctx := Parse(testJSONData)
	optimizer := NewQueryOptimizer()
	
	query := &Query{
		Path: "users[0].name",
	}
	
	// Test optimized execution
	result := optimizer.OptimizeQuery(ctx, query)
	if !result.Exists() {
		t.Error("Optimizer should execute query successfully")
	}
	
	// Test cache hit on second execution
	result2 := optimizer.OptimizeQuery(ctx, query)
	if result2.String() != result.String() {
		t.Error("Cached result should match original result")
	}
	
	// Test cache stats
	queryStats, pathStats := optimizer.GetCacheStats()
	_ = queryStats
	_ = pathStats
	
	// Test performance metrics
	metrics := optimizer.GetPerformanceMetrics()
	if metrics.CacheHitRatio < 0 || metrics.CacheHitRatio > 1 {
		t.Error("Cache hit ratio should be between 0 and 1")
	}
}

func TestGlobalOptimizer(t *testing.T) {
	ctx := Parse(testJSONData)
	query := &Query{Path: "config.api.version"}
	
	result := OptimizedExecute(ctx, query)
	if !result.Exists() {
		t.Error("Global optimizer should execute query successfully")
	}
}

func TestBackwardCompatibility(t *testing.T) {
	ctx := Parse(testJSONData)
	
	// Test that existing Get() function still works
	result1 := ctx.Get("users[0].name")
	if !result1.Exists() {
		t.Error("Existing Get() function should still work")
	}
	
	// Test that new query system produces same results
	query := NewQueryBuilder().Path("users[0].name").Build()
	result2 := query.Execute(ctx)
	
	if result1.String() != result2.String() {
		t.Errorf("New query system result (%s) should match existing Get() result (%s)", 
			result2.String(), result1.String())
	}
}

// Benchmark tests for performance comparison

func BenchmarkExistingGet(b *testing.B) {
	ctx := Parse(testJSONData)
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		result := ctx.Get("users[0].name")
		_ = result.String()
	}
}

func BenchmarkNewQuerySystem(b *testing.B) {
	ctx := Parse(testJSONData)
	query := NewQueryBuilder().Path("users[0].name").Build()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		result := query.Execute(ctx)
		_ = result.String()
	}
}

func BenchmarkOptimizedQuery(b *testing.B) {
	ctx := Parse(testJSONData)
	query := NewQueryBuilder().Path("users[0].name").Build()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		result := OptimizedExecute(ctx, query)
		_ = result.String()
	}
}

func BenchmarkQueryBuilderCreation(b *testing.B) {
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		builder := NewQueryBuilder().
			Path("users").
			Filter("age", ">", 25).
			Select("name", "email").
			Build()
		_ = builder
	}
}

func BenchmarkQueryBuilderPool(b *testing.B) {
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		builder := GetQueryBuilderFromPool()
		query := builder.Path("users").
			Filter("age", ">", 25).
			Select("name", "email").
			Build()
		PutQueryBuilderToPool(builder)
		_ = query
	}
}