package fj

import (
	"testing"
)

// TestEnhancedFunctions tests the new enhanced function names
func TestEnhancedExtract(t *testing.T) {
	json := `{"user":{"name":"John","age":30}}`
	
	// Test Extract function
	ctx := Extract(json, "user.name")
	if ctx.String() != "John" {
		t.Errorf("Extract() = %v, want %v", ctx.String(), "John")
	}
	
	// Test ExtractBytes function
	data := []byte(json)
	ctx2 := ExtractBytes(data, "user.age")
	if ctx2.Int64() != 30 {
		t.Errorf("ExtractBytes() = %v, want %v", ctx2.Int64(), 30)
	}
	
	// Test ExtractMultiple function
	contexts := ExtractMultiple(json, "user.name", "user.age")
	if len(contexts) != 2 {
		t.Errorf("ExtractMultiple() length = %v, want %v", len(contexts), 2)
	}
	if contexts[0].String() != "John" {
		t.Errorf("ExtractMultiple()[0] = %v, want %v", contexts[0].String(), "John")
	}
	if contexts[1].Int64() != 30 {
		t.Errorf("ExtractMultiple()[1] = %v, want %v", contexts[1].Int64(), 30)
	}
}

func TestEnhancedValidation(t *testing.T) {
	validJSON := []byte(`{"name":"John","age":30}`)
	invalidJSON := []byte(`{"name":"John","age":}`)
	
	// Test Validate function
	if !Validate(validJSON) {
		t.Error("Validate() should return true for valid JSON")
	}
	if Validate(invalidJSON) {
		t.Error("Validate() should return false for invalid JSON")
	}
	
	// Test ValidateString function
	if !ValidateString(`{"name":"John"}`) {
		t.Error("ValidateString() should return true for valid JSON string")
	}
	if ValidateString(`{"name":}`) {
		t.Error("ValidateString() should return false for invalid JSON string")
	}
}

func TestEnhancedTransform(t *testing.T) {
	// Register a test transformer
	registry := GetTransformerRegistry()
	err := registry.RegisterFunc("testEnhancedUpper", func(json, arg string) string {
		return "TRANSFORMED:" + json
	})
	if err != nil {
		t.Fatalf("Failed to register test transformer: %v", err)
	}
	defer registry.Remove("testEnhancedUpper")
	
	json := `{"name":"john"}`
	
	// Test Transform function
	result := Transform(json, "testEnhancedUpper")
	expected := "TRANSFORMED:" + json
	if result != expected {
		t.Errorf("Transform() = %v, want %v", result, expected)
	}
	
	// Test Transform with non-existent transformer
	result2 := Transform(json, "nonExistent")
	if result2 != json {
		t.Errorf("Transform() with non-existent transformer should return original JSON")
	}
}

func TestEnhancedChain(t *testing.T) {
	// Register test transformers
	registry := GetTransformerRegistry()
	
	err := registry.RegisterFunc("testChain1", func(json, arg string) string {
		return "STEP1:" + json
	})
	if err != nil {
		t.Fatalf("Failed to register testChain1: %v", err)
	}
	defer registry.Remove("testChain1")
	
	err = registry.RegisterFunc("testChain2", func(json, arg string) string {
		return "STEP2:" + json
	})
	if err != nil {
		t.Fatalf("Failed to register testChain2: %v", err)
	}
	defer registry.Remove("testChain2")
	
	json := `{"name":"john"}`
	
	// Test Chain function
	result := Chain(json, "testChain1", "testChain2")
	expected := "STEP2:STEP1:" + json
	if result != expected {
		t.Errorf("Chain() = %v, want %v", result, expected)
	}
}

func TestEnhancedQuery(t *testing.T) {
	json := `{"users":[{"name":"john"},{"name":"jane"}]}`
	
	// Test Query function (basic path query)
	ctx := Query(json, "users.0.name")
	if ctx.String() != "john" {
		t.Errorf("Query() = %v, want %v", ctx.String(), "john")
	}
}

func TestEnhancedBuild(t *testing.T) {
	// Test Build function
	path := Build("user", "profile", "name")
	expected := "user.profile.name"
	if path != expected {
		t.Errorf("Build() = %v, want %v", path, expected)
	}
	
	// Test Build with single part
	path2 := Build("user")
	if path2 != "user" {
		t.Errorf("Build() with single part = %v, want %v", path2, "user")
	}
	
	// Test Build with empty parts
	path3 := Build()
	if path3 != "" {
		t.Errorf("Build() with no parts = %v, want empty string", path3)
	}
}

func TestEnhancedRegister(t *testing.T) {
	// Test Register function
	err := Register("testEnhancedRegister", func(json, arg string) string {
		return "REGISTERED:" + json
	})
	if err != nil {
		t.Errorf("Register() returned error: %v", err)
	}
	
	// Test Exists function
	if !Exists("testEnhancedRegister") {
		t.Error("Exists() should return true for registered transformer")
	}
	
	// Test with non-existent transformer
	if Exists("nonExistentTransformer") {
		t.Error("Exists() should return false for non-existent transformer")
	}
	
	// Clean up
	registry := GetTransformerRegistry()
	registry.Remove("testEnhancedRegister")
}

func TestEnhancedProcess(t *testing.T) {
	json := `{"name":"John","age":30}`
	ctx := Parse(json)
	
	// Test Process function (simple field access)
	result := Process(ctx, "name")
	if result.String() != "John" {
		t.Errorf("Process() = %v, want %v", result.String(), "John")
	}
}

func TestEnhancedFrom(t *testing.T) {
	// Test From with string input
	json := `{"name":"John","age":30}`
	ctx := From(json, "name")
	if ctx.String() != "John" {
		t.Errorf("From() with string = %v, want %v", ctx.String(), "John")
	}
	
	// Test From with byte slice input
	data := []byte(json)
	ctx2 := From(data, "age")
	if ctx2.Int64() != 30 {
		t.Errorf("From() with bytes = %v, want %v", ctx2.Int64(), 30)
	}
	
	// Test From with Context input
	originalCtx := Parse(json)
	ctx3 := From(originalCtx, "name")
	if ctx3.String() != "John" {
		t.Errorf("From() with Context = %v, want %v", ctx3.String(), "John")
	}
	
	// Test From with invalid input type
	ctx4 := From(123)
	if ctx4.Exists() {
		t.Error("From() with invalid input should return empty context")
	}
	
	// Test From without path
	ctx5 := From(json)
	if !ctx5.Exists() {
		t.Error("From() without path should return parsed context")
	}
}

func TestDeprecatedFunctions(t *testing.T) {
	json := `{"name":"John","age":30}`
	data := []byte(json)
	
	// Test deprecated function still works
	ctx := GetBytesDeprecated(data, "name")
	if ctx.String() != "John" {
		t.Errorf("GetBytesDeprecated() = %v, want %v", ctx.String(), "John")
	}
	
	// Test deprecated GetMul
	contexts := GetMulDeprecated(json, "name", "age")
	if len(contexts) != 2 {
		t.Errorf("GetMulDeprecated() length = %v, want %v", len(contexts), 2)
	}
	
	// Test deprecated IsValidJSONBytes
	if !IsValidJSONBytesDeprecated(data) {
		t.Error("IsValidJSONBytesDeprecated() should return true for valid JSON")
	}
	
	// Test deprecated AddTransformer
	AddTransformerDeprecated("testDeprecated", func(json, arg string) string {
		return json
	})
	if !IsTransformerRegisteredDeprecated("testDeprecated") {
		t.Error("IsTransformerRegisteredDeprecated() should return true for registered transformer")
	}
	
	// Clean up
	registry := GetTransformerRegistry()
	registry.Remove("testDeprecated")
}

// TestBackwardCompatibility ensures all existing functions still work
func TestBackwardCompatibility(t *testing.T) {
	json := `{"user":{"name":"Alice","age":25}}`
	data := []byte(json)
	
	// Test original functions still work
	ctx1 := Get(json, "user.name")
	if ctx1.String() != "Alice" {
		t.Error("Original Get() function should still work")
	}
	
	ctx2 := GetBytes(data, "user.age")
	if ctx2.Int64() != 25 {
		t.Error("Original GetBytes() function should still work")
	}
	
	contexts := GetMul(json, "user.name", "user.age")
	if len(contexts) != 2 {
		t.Error("Original GetMul() function should still work")
	}
	
	if !IsValidJSONBytes(data) {
		t.Error("Original IsValidJSONBytes() function should still work")
	}
	
	// Test original transformer functions
	AddTransformer("backwardTest", func(json, arg string) string {
		return "BACKWARD:" + json
	})
	
	if !IsTransformerRegistered("backwardTest") {
		t.Error("Original IsTransformerRegistered() function should still work")
	}
	
	// Clean up
	registry := GetTransformerRegistry()
	registry.Remove("backwardTest")
}

func TestEnhancedFunctionConsistency(t *testing.T) {
	json := `{"data":{"value":"test"}}`
	data := []byte(json)
	path := "data.value"
	
	// Ensure enhanced functions return same results as originals
	original := Get(json, path)
	enhanced := Extract(json, path)
	
	if original.String() != enhanced.String() {
		t.Errorf("Enhanced Extract() inconsistent with original Get(): %v != %v", 
			enhanced.String(), original.String())
	}
	
	originalBytes := GetBytes(data, path)
	enhancedBytes := ExtractBytes(data, path)
	
	if originalBytes.String() != enhancedBytes.String() {
		t.Errorf("Enhanced ExtractBytes() inconsistent with original GetBytes(): %v != %v",
			enhancedBytes.String(), originalBytes.String())
	}
	
	originalValid := IsValidJSONBytes(data)
	enhancedValid := Validate(data)
	
	if originalValid != enhancedValid {
		t.Errorf("Enhanced Validate() inconsistent with original IsValidJSONBytes(): %v != %v",
			enhancedValid, originalValid)
	}
}