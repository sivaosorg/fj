package fj

import (
	"strings"
	"testing"
)

// TestTransformerStrategy tests the Strategy pattern implementation
func TestTransformerStrategy(t *testing.T) {
	// First register a test transformer
	registry := GetTransformerRegistry()
	err := registry.RegisterFunc("testUpper", func(json, arg string) string {
		return strings.ToUpper(json)
	})
	if err != nil {
		t.Fatalf("Failed to register test transformer: %v", err)
	}
	defer registry.Remove("testUpper")

	// Test creating a valid strategy
	strategy := NewTransformerStrategy("testUpper")
	if strategy == nil {
		t.Error("NewTransformerStrategy should return non-nil for existing transformer")
	}
	if !strategy.IsValid() {
		t.Error("Strategy should be valid")
	}
	if strategy.Name() != "testUpper" {
		t.Errorf("Strategy name = %v, want %v", strategy.Name(), "testUpper")
	}

	// Test executing the strategy
	result := strategy.Execute("hello", "")
	if result != "HELLO" {
		t.Errorf("Strategy Execute() = %v, want %v", result, "HELLO")
	}

	// Test creating strategy for non-existent transformer
	invalidStrategy := NewTransformerStrategy("nonExistent")
	if invalidStrategy != nil {
		t.Error("NewTransformerStrategy should return nil for non-existent transformer")
	}
}

func TestTransformerStrategyChain(t *testing.T) {
	// Register test transformers
	registry := GetTransformerRegistry()
	
	err := registry.RegisterFunc("testUpper", func(json, arg string) string {
		return strings.ToUpper(json)
	})
	if err != nil {
		t.Fatalf("Failed to register testUpper: %v", err)
	}
	defer registry.Remove("testUpper")
	
	err = registry.RegisterFunc("testReplace", func(json, arg string) string {
		return strings.ReplaceAll(json, "HELLO", "HI")
	})
	if err != nil {
		t.Fatalf("Failed to register testReplace: %v", err)
	}
	defer registry.Remove("testReplace")

	// Test creating and using a chain
	chain := NewTransformerStrategyChain().
		Add("testUpper").
		Add("testReplace")

	if chain.Length() != 2 {
		t.Errorf("Chain length = %v, want %v", chain.Length(), 2)
	}

	result := chain.Execute("hello world", "")
	expected := "HI WORLD"
	if result != expected {
		t.Errorf("Chain Execute() = %v, want %v", result, expected)
	}

	// Test chain names
	names := chain.Names()
	expectedNames := []string{"testUpper", "testReplace"}
	if len(names) != len(expectedNames) {
		t.Errorf("Chain names length = %v, want %v", len(names), len(expectedNames))
	}
	for i, name := range names {
		if name != expectedNames[i] {
			t.Errorf("Chain name[%d] = %v, want %v", i, name, expectedNames[i])
		}
	}

	// Test clearing the chain
	chain.Clear()
	if chain.Length() != 0 {
		t.Errorf("Chain length after clear = %v, want %v", chain.Length(), 0)
	}
}

func TestTransformerStrategyChainWithArgs(t *testing.T) {
	// Register a test transformer that uses arguments
	registry := GetTransformerRegistry()
	
	err := registry.RegisterFunc("testPrefix", func(json, arg string) string {
		if arg == "" {
			return json
		}
		return arg + ":" + json
	})
	if err != nil {
		t.Fatalf("Failed to register testPrefix: %v", err)
	}
	defer registry.Remove("testPrefix")
	
	err = registry.RegisterFunc("testSuffix", func(json, arg string) string {
		if arg == "" {
			return json
		}
		return json + ":" + arg
	})
	if err != nil {
		t.Fatalf("Failed to register testSuffix: %v", err)
	}
	defer registry.Remove("testSuffix")

	chain := NewTransformerStrategyChain().
		Add("testPrefix").
		Add("testSuffix")

	result := chain.ExecuteWithArgs("hello", "start", "end")
	expected := "start:hello:end"
	if result != expected {
		t.Errorf("Chain ExecuteWithArgs() = %v, want %v", result, expected)
	}
}

func TestParseTransformerChain(t *testing.T) {
	// Test parsing simple chain
	chain := ParseTransformerChain("trim|uppercase")
	names := chain.Names()
	
	// The exact transformers depend on what's registered, but we should get some
	if len(names) == 0 {
		t.Error("ParseTransformerChain should parse at least some transformers")
	}

	// Test parsing with @ prefix
	chain = ParseTransformerChain("@trim|@uppercase")
	if chain.Length() == 0 {
		t.Error("ParseTransformerChain should handle @ prefixes")
	}

	// Test empty string
	chain = ParseTransformerChain("")
	if chain.Length() != 0 {
		t.Error("ParseTransformerChain with empty string should return empty chain")
	}

	// Test with transformer arguments
	chain = ParseTransformerChain("replace:{target:old}|uppercase")
	if chain.Length() == 0 {
		t.Error("ParseTransformerChain should handle transformers with arguments")
	}
}

func TestApplyTransformerChain(t *testing.T) {
	// This will test with whatever built-in transformers exist
	json := `{"name": "john"}`
	result := ApplyTransformerChain(json, "")
	
	// Empty chain should return original
	if result != json {
		t.Errorf("ApplyTransformerChain with empty chain = %v, want %v", result, json)
	}
}

// TestPipeProcessor tests the Chain of Responsibility pattern for pipe processing
func TestPipeProcessor(t *testing.T) {
	processor := NewPipeProcessor()

	// Test with empty pipe
	ctx := Parse(`{"name": "John"}`)
	result := processor.Process(ctx, "")
	if result.String() != ctx.String() {
		t.Error("Process with empty pipe should return original context")
	}

	// Test with simple field access
	result = processor.Process(ctx, "name")
	if result.String() != "John" {
		t.Errorf("Process with field access = %v, want %v", result.String(), "John")
	}
}

func TestTransformerPipeHandler(t *testing.T) {
	handler := &TransformerPipeHandler{}
	
	// Register a test transformer
	registry := GetTransformerRegistry()
	err := registry.RegisterFunc("testPipeUpper", func(json, arg string) string {
		return strings.ToUpper(json)
	})
	if err != nil {
		t.Fatalf("Failed to register test transformer: %v", err)
	}
	defer registry.Remove("testPipeUpper")

	// Test CanHandle
	if !handler.CanHandle("testPipeUpper") {
		t.Error("TransformerPipeHandler should handle registered transformer")
	}
	if !handler.CanHandle("@testPipeUpper") {
		t.Error("TransformerPipeHandler should handle transformer with @ prefix")
	}
	if handler.CanHandle("nonExistent") {
		t.Error("TransformerPipeHandler should not handle non-existent transformer")
	}

	// Test Handle
	ctx := Parse(`{"name": "john"}`)
	result, ok := handler.Handle(ctx, "testPipeUpper", "")
	if !ok {
		t.Error("Handle should return true for successful transformation")
	}
	if !strings.Contains(result.String(), "JOHN") {
		t.Errorf("Handle result should contain uppercased content: %v", result.String())
	}

	// Test Handle with @ prefix
	result, ok = handler.Handle(ctx, "@testPipeUpper", "")
	if !ok {
		t.Error("Handle should work with @ prefix")
	}
}

func TestPathPipeHandler(t *testing.T) {
	handler := &PathPipeHandler{}

	// Test CanHandle
	if !handler.CanHandle("name") {
		t.Error("PathPipeHandler should handle simple field names")
	}
	if !handler.CanHandle("user.name") {
		t.Error("PathPipeHandler should handle dotted paths")
	}
	if handler.CanHandle("@transformer") {
		t.Error("PathPipeHandler should not handle transformer operations")
	}

	// Test Handle
	ctx := Parse(`{"user": {"name": "Alice"}}`)
	result, ok := handler.Handle(ctx, "user", "")
	if !ok {
		t.Error("Handle should return true for valid field access")
	}
	if !result.Exists() {
		t.Error("Handle should return existing context for valid path")
	}
}

func TestQueryPipeHandler(t *testing.T) {
	handler := &QueryPipeHandler{}

	// Test CanHandle
	if !handler.CanHandle("#(name==john)") {
		t.Error("QueryPipeHandler should handle query operations")
	}
	if handler.CanHandle("name") {
		t.Error("QueryPipeHandler should not handle regular field names")
	}

	ctx := Parse(`[{"name": "john"}, {"name": "jane"}]`)
	_, ok := handler.Handle(ctx, "#(name==john)", "")
	if !ok {
		t.Error("Handle should return true for query operations")
	}
}

func TestParseOperation(t *testing.T) {
	tests := []struct {
		input    string
		wantName string
		wantArg  string
	}{
		{
			input:    "transformer",
			wantName: "transformer",
			wantArg:  "",
		},
		{
			input:    "replace:{\"target\":\"old\",\"replacement\":\"new\"}",
			wantName: "replace",
			wantArg:  "{\"target\":\"old\",\"replacement\":\"new\"}",
		},
		{
			input:    "uppercase:arg",
			wantName: "uppercase",
			wantArg:  "arg",
		},
		{
			input:    "complex:arg:with:colons",
			wantName: "complex",
			wantArg:  "arg:with:colons",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			name, arg := parseOperation(tt.input)
			if name != tt.wantName {
				t.Errorf("parseOperation() name = %v, want %v", name, tt.wantName)
			}
			if arg != tt.wantArg {
				t.Errorf("parseOperation() arg = %v, want %v", arg, tt.wantArg)
			}
		})
	}
}

func TestIsValidFieldName(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"name", true},
		{"user_name", true},
		{"userName", true},
		{"user-name", true},
		{"user123", true},
		{"", false},
		{"@transformer", false},
		{"name{with:args}", false},
		{"name.with.dots", false}, // dots make it a path, not a simple field
		{"name with space", false},
		{"name/with/slash", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := isValidFieldName(tt.input)
			if got != tt.want {
				t.Errorf("isValidFieldName(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestProcessPipeConvenience(t *testing.T) {
	ctx := Parse(`{"name": "John", "age": 30}`)
	
	// Test convenience function
	result := ProcessPipe(ctx, "name")
	if result.String() != "John" {
		t.Errorf("ProcessPipe convenience function = %v, want %v", result.String(), "John")
	}
}

func TestChainedPipeProcessor(t *testing.T) {
	processor := NewChainedPipeProcessor().
		WithTransformers().
		WithPaths().
		WithQueries().
		Build()

	ctx := Parse(`{"name": "John"}`)
	result := processor.Process(ctx, "name")
	if result.String() != "John" {
		t.Errorf("ChainedPipeProcessor = %v, want %v", result.String(), "John")
	}

	// Test empty processor
	emptyProcessor := NewChainedPipeProcessor().Build()
	result = emptyProcessor.Process(ctx, "name")
	// Should return original context since no handlers are present
	if !result.Exists() {
		t.Error("Empty processor should still return a context")
	}
}