package fj

import "testing"

func TestBasicFunctionality(t *testing.T) {
	// Use a simple JSON structure that should work  
	jsonData := `{"key": "value", "arr": [{"name": "Alice"}]}`
	
	// Test Parse
	ctx := Parse(jsonData)
	if !ctx.Exists() {
		t.Error("Parse should return existing context")
	}
	
	t.Logf("Parsed context: %s", ctx.String())
	
	// Test basic Get
	result := ctx.Get("key")
	if !result.Exists() {
		t.Error("Get should find key")
	}
	
	t.Logf("Key result: %s", result.String())
	expected := "value"
	if result.String() != expected {
		t.Errorf("Expected %s, got %s", expected, result.String())
	}
	
	// Test array element
	arrResult := ctx.Get("arr")
	t.Logf("Array result exists: %v", arrResult.Exists())
	t.Logf("Array result: '%s'", arrResult.String())
	
	// Test array access - try both arr.0 and arr[0] formats
	formats := []string{"arr.0.name", "arr[0].name"}
	for _, format := range formats {
		userResult := ctx.Get(format)
		t.Logf("Format '%s' - exists: %v, result: '%s'", format, userResult.Exists(), userResult.String())
		if userResult.Exists() && userResult.String() == "Alice" {
			t.Logf("Success with format: %s", format)
			return // Success!
		}
	}
	
	t.Error("No array access format worked")
}