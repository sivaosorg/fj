package fj

import (
	"errors"
	"strings"
	"testing"
)

// TestPathBuilder tests the PathBuilder implementation
func TestPathBuilder(t *testing.T) {
	tests := []struct {
		name     string
		builder  func() PathBuilder
		expected string
	}{
		{
			name: "Root only",
			builder: func() PathBuilder {
				return NewPathBuilder().Root("user")
			},
			expected: "user",
		},
		{
			name: "Root with field",
			builder: func() PathBuilder {
				return NewPathBuilder().Root("user").Field("name")
			},
			expected: "user.name",
		},
		{
			name: "Complex path with array index",
			builder: func() PathBuilder {
				return NewPathBuilder().
					Root("users").
					ArrayIndex(0).
					Field("profile").
					Field("name")
			},
			expected: "users.0.profile.name",
		},
		{
			name: "Path with array wildcard",
			builder: func() PathBuilder {
				return NewPathBuilder().
					Root("users").
					ArrayAll().
					Field("name")
			},
			expected: "users.*.name",
		},
		{
			name: "Reset and reuse builder",
			builder: func() PathBuilder {
				builder := NewPathBuilder().Root("old").Field("path")
				builder.Reset()
				return builder.Root("new").Field("path")
			},
			expected: "new.path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.builder().Build()
			if result != tt.expected {
				t.Errorf("PathBuilder.Build() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPathBuilderFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Parse existing path",
			input:    "user.profile.name",
			expected: "user.profile.name",
		},
		{
			name:     "Empty path",
			input:    "",
			expected: "",
		},
		{
			name:     "Single field",
			input:    "user",
			expected: "user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := PathBuilderFromString(tt.input)
			result := builder.Build()
			if result != tt.expected {
				t.Errorf("PathBuilderFromString(%v).Build() = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPathBuilderChaining(t *testing.T) {
	// Test that we can extend existing paths
	builder := PathBuilderFromString("user.profile")
	result := builder.Field("settings").Field("theme").Build()
	expected := "user.profile.settings.theme"

	if result != expected {
		t.Errorf("PathBuilder chaining = %v, want %v", result, expected)
	}
}

// TestTransformerRegistry tests the TransformerRegistry singleton
func TestTransformerRegistry(t *testing.T) {
	// Test singleton behavior
	registry1 := GetTransformerRegistry()
	registry2 := GetTransformerRegistry()

	if registry1 != registry2 {
		t.Error("GetTransformerRegistry() should return the same instance (singleton)")
	}
}

func TestTransformerRegistryFunctionRegistration(t *testing.T) {
	registry := GetTransformerRegistry()

	// Test function registration
	testFunc := func(json, arg string) string {
		return strings.ToUpper(json)
	}

	err := registry.RegisterFunc("testUpper", testFunc)
	if err != nil {
		t.Errorf("RegisterFunc() returned error: %v", err)
	}

	// Test existence check
	if !registry.Exists("testUpper") {
		t.Error("Registered transformer should exist")
	}

	// Test retrieval
	transformer, exists := registry.Get("testUpper")
	if !exists {
		t.Error("Get() should return true for existing transformer")
	}
	if transformer == nil {
		t.Error("Get() should return non-nil transformer")
	}

	// Test transformation
	result := transformer.Transform("hello", "")
	if result != "HELLO" {
		t.Errorf("Transform() = %v, want %v", result, "HELLO")
	}

	// Test removal
	removed := registry.Remove("testUpper")
	if !removed {
		t.Error("Remove() should return true for existing transformer")
	}

	// Verify removal
	if registry.Exists("testUpper") {
		t.Error("Transformer should not exist after removal")
	}
}

func TestTransformerRegistryObjectRegistration(t *testing.T) {
	registry := GetTransformerRegistry()

	// Test object registration
	testTransformer := &testTransformerImpl{name: "testLower"}

	err := registry.Register("testLower", testTransformer)
	if err != nil {
		t.Errorf("Register() returned error: %v", err)
	}

	// Test retrieval and transformation
	transformer, exists := registry.Get("testLower")
	if !exists {
		t.Error("Get() should return true for existing transformer")
	}

	result := transformer.Transform("HELLO", "")
	if result != "hello" {
		t.Errorf("Transform() = %v, want %v", result, "hello")
	}

	// Clean up
	registry.Remove("testLower")
}

func TestTransformerRegistryErrorHandling(t *testing.T) {
	registry := GetTransformerRegistry()

	// Test empty name error
	err := registry.RegisterFunc("", func(json, arg string) string { return json })
	if err == nil {
		t.Error("RegisterFunc() with empty name should return error")
	}

	// Test nil function error
	err = registry.RegisterFunc("nilFunc", nil)
	if err == nil {
		t.Error("RegisterFunc() with nil function should return error")
	}

	// Test nil transformer error
	err = registry.Register("nilTransformer", nil)
	if err == nil {
		t.Error("Register() with nil transformer should return error")
	}
}

func TestTransformerRegistryList(t *testing.T) {
	registry := GetTransformerRegistry()

	// Get initial count (built-in transformers)
	initialNames := registry.List()
	initialCount := len(initialNames)

	// Add a test transformer
	err := registry.RegisterFunc("listTest", func(json, arg string) string { return json })
	if err != nil {
		t.Errorf("RegisterFunc() returned error: %v", err)
	}

	// Check that list includes new transformer
	names := registry.List()
	if len(names) != initialCount+1 {
		t.Errorf("List() length = %v, want %v", len(names), initialCount+1)
	}

	// Check that our transformer is in the list
	found := false
	for _, name := range names {
		if name == "listTest" {
			found = true
			break
		}
	}
	if !found {
		t.Error("List() should include registered transformer")
	}

	// Clean up
	registry.Remove("listTest")
}

// TestContextFactory tests the ContextFactory implementation
func TestContextFactory(t *testing.T) {
	factory := NewContextFactory()

	// Test basic creation
	json := `{"name": "John", "age": 30}`
	ctx := factory.Create(json)

	if !ctx.Exists() {
		t.Error("Create() should return existing context for valid JSON")
	}

	// Test creation with path
	ctx = factory.Create(json, "name")
	if ctx.String() != "John" {
		t.Errorf("Create() with path = %v, want %v", ctx.String(), "John")
	}
}

func TestContextFactoryFromBytes(t *testing.T) {
	factory := NewContextFactory()

	jsonBytes := []byte(`{"name": "Alice", "age": 25}`)
	ctx := factory.CreateFromBytes(jsonBytes)

	if !ctx.Exists() {
		t.Error("CreateFromBytes() should return existing context for valid JSON")
	}

	// Test with path
	ctx = factory.CreateFromBytes(jsonBytes, "age")
	if ctx.Int64() != 25 {
		t.Errorf("CreateFromBytes() with path = %v, want %v", ctx.Int64(), 25)
	}
}

func TestContextFactoryEmpty(t *testing.T) {
	factory := NewContextFactory()

	ctx := factory.CreateEmpty()
	if ctx.Exists() {
		t.Error("CreateEmpty() should return non-existing context")
	}
}

func TestContextFactoryWithError(t *testing.T) {
	factory := NewContextFactory()

	testErr := errors.New("test error")
	ctx := factory.CreateWithError(testErr)

	if ctx.ErrMessage() != testErr.Error() {
		t.Errorf("CreateWithError() = %v, want %v", ctx.ErrMessage(), testErr.Error())
	}
}

func TestDefaultContextFactory(t *testing.T) {
	// Test convenience functions
	json := `{"user": {"name": "Bob"}}`

	ctx := CreateContext(json, "user.name")
	if ctx.String() != "Bob" {
		t.Errorf("CreateContext() = %v, want %v", ctx.String(), "Bob")
	}

	jsonBytes := []byte(json)
	ctx = CreateContextFromBytes(jsonBytes, "user.name")
	if ctx.String() != "Bob" {
		t.Errorf("CreateContextFromBytes() = %v, want %v", ctx.String(), "Bob")
	}

	emptyCtx := CreateEmptyContext()
	if emptyCtx.Exists() {
		t.Error("CreateEmptyContext() should return non-existing context")
	}

	testErr := errors.New("test error")
	errorCtx := CreateContextWithError(testErr)
	if errorCtx.ErrMessage() != testErr.Error() {
		t.Errorf("CreateContextWithError() = %v, want %v", errorCtx.ErrMessage(), testErr.Error())
	}
}

// Test helper implementations

type testTransformerImpl struct {
	name string
}

func (t *testTransformerImpl) Transform(json, arg string) string {
	return strings.ToLower(json)
}

func (t *testTransformerImpl) Name() string {
	return t.name
}

// TestConvenienceFunctions tests the convenience functions for registration
func TestConvenienceFunctions(t *testing.T) {
	// Test RegisterTransformerFunc
	err := RegisterTransformerFunc("convenienceTest", func(json, arg string) string {
		return "test-" + json
	})
	if err != nil {
		t.Errorf("RegisterTransformerFunc() returned error: %v", err)
	}

	// Test that it's available in registry
	registry := GetTransformerRegistry()
	if !registry.Exists("convenienceTest") {
		t.Error("RegisterTransformerFunc() should add transformer to registry")
	}

	// Test RegisterTransformer
	testTransformer := &testTransformerImpl{name: "convenienceObj"}
	err = RegisterTransformer("convenienceObj", testTransformer)
	if err != nil {
		t.Errorf("RegisterTransformer() returned error: %v", err)
	}

	if !registry.Exists("convenienceObj") {
		t.Error("RegisterTransformer() should add transformer to registry")
	}

	// Clean up
	registry.Remove("convenienceTest")
	registry.Remove("convenienceObj")
}