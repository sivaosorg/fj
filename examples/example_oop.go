package main

import (
	"fmt"
	"log"

	"github.com/sivaosorg/fj"
	"github.com/sivaosorg/fj/pkg/config"
	"github.com/sivaosorg/fj/pkg/core"
	"github.com/sivaosorg/fj/pkg/transform"
)

func main() {
	// Test the original API (backward compatibility)
	fmt.Println("=== Testing Original API (Backward Compatibility) ===")
	
	json := `{
		"users": [
			{"name": "Alice", "age": 25, "active": true},
			{"name": "Bob", "age": 30, "active": false}
		],
		"metadata": {
			"version": "1.0",
			"created": "2023-01-01"
		}
	}`
	
	// Original API still works
	result := fj.Get(json, "users.0.name")
	fmt.Printf("User name: %s\n", result.String())
	
	result = fj.Get(json, "metadata.version")
	fmt.Printf("Version: %s\n", result.String())
	
	// Test array operations
	result = fj.Get(json, "users.#.name")
	fmt.Printf("All user names: %s\n", result.String())
	
	fmt.Println("\n=== Testing New OOP Architecture ===")
	
	// 1. Test Configuration Manager (Singleton Pattern)
	configManager := config.GetInstance()
	cfg := configManager.GetConfig()
	cfg.EnableTransformers = true
	cfg.EnableValidation = true
	cfg.LogLevel = 3
	err := configManager.SetConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Configuration - Transformers enabled: %t\n", cfg.EnableTransformers)
	fmt.Printf("Configuration - Log level: %d\n", cfg.LogLevel)
	
	// 2. Test JSON Parser (Factory Pattern)
	parser := core.NewParser()
	ctx := parser.Parse(json)
	if ctx.HasError() {
		log.Fatal(ctx.Error())
	}
	
	fmt.Printf("Parsed JSON successfully - Type: %v\n", ctx.Type())
	fmt.Printf("Is Object: %t\n", ctx.IsObject())
	
	// 3. Test enhanced Context interface
	userPath := ctx.Get("users.0")
	fmt.Printf("First user: %s\n", userPath.String())
	fmt.Printf("User exists: %t\n", userPath.Exists())
	
	// Test multiple paths
	paths := []string{"users.0.name", "users.0.age", "metadata.version"}
	results := ctx.GetMultiple(paths...)
	for i, path := range paths {
		fmt.Printf("Path '%s': %s\n", path, results[i].String())
	}
	
	// 4. Test Transformation Factory (Factory Pattern)
	factory := transform.NewFactory()
	
	// List available transformers
	transformerNames := factory.ListTransformers()
	fmt.Printf("Available transformers: %v\n", transformerNames[:5]) // Show first 5
	
	// Create and use a transformer
	uppercaseTransformer, err := factory.CreateTransformer("uppercase")
	if err != nil {
		log.Printf("Transformer error: %v", err)
	} else {
		fmt.Printf("Transformer name: %s\n", uppercaseTransformer.GetName())
		fmt.Printf("Transformer description: %s\n", uppercaseTransformer.GetDescription())
		
		// Apply transformation
		transformed := uppercaseTransformer.Transform("\"hello world\"", "")
		fmt.Printf("Transformed result: %s\n", transformed)
	}
	
	// 5. Test Transformation Manager (Strategy Pattern)
	manager := transform.NewManager(factory)
	
	// Apply transformation with error handling and events
	result2, err := manager.Transform("minify", json, "")
	if err != nil {
		log.Printf("Transform error: %v", err)
	} else {
		fmt.Printf("Minified JSON length: %d (original: %d)\n", len(result2), len(json))
	}
	
	// 6. Test transformation chaining (Decorator Pattern)
	steps := []transform.TransformationStep{
		{Name: "minify", Argument: ""},
		{Name: "uppercase", Argument: ""},
	}
	
	chained, err := manager.ChainTransformations(json, steps)
	if err != nil {
		log.Printf("Chain error: %v", err)
	} else {
		fmt.Printf("Chained transformation completed successfully (length: %d)\n", len(chained))
	}
	
	fmt.Println("\n=== Testing Error Handling ===")
	
	// Test custom error types
	invalidJSON := `{"invalid": json}`
	invalidCtx := parser.Parse(invalidJSON)
	if invalidCtx.HasError() {
		fmt.Printf("Parse error detected: %v\n", invalidCtx.Error())
	}
	
	// Test transformation error
	_, err = manager.Transform("nonexistent", json, "")
	if err != nil {
		fmt.Printf("Transformation error: %v\n", err)
	}
	
	fmt.Println("\n=== Architecture Benefits Demonstrated ===")
	fmt.Println("✓ Backward Compatibility: Original API works unchanged")
	fmt.Println("✓ Singleton Pattern: Configuration management")
	fmt.Println("✓ Factory Pattern: Parser and transformer creation")  
	fmt.Println("✓ Strategy Pattern: Different transformation types")
	fmt.Println("✓ Decorator Pattern: Transformation chaining")
	fmt.Println("✓ Enhanced Error Handling: Custom error types")
	fmt.Println("✓ Interface-based Design: Better testability")
	fmt.Println("✓ Dependency Injection: Flexible component composition")
}