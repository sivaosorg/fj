package main

import (
	"fmt"
	"strings"
	"github.com/sivaosorg/fj"
)

func main() {
	fmt.Println("=== Performance Benchmarks and Best Practices ===\n")

	// Sample data for benchmarking
	largeJSON := generateLargeJSON()
	smallJSON := `{"user": {"name": "John", "age": 30, "active": true}}`

	fmt.Println("🚀 Performance Comparison: Original vs Enhanced API")
	fmt.Println("---------------------------------------------------")

	// Benchmark extraction operations
	benchmarkExtraction(largeJSON)
	
	fmt.Println("\n📊 Memory Efficiency Examples")
	fmt.Println("------------------------------")
	
	// Memory efficient processing
	demonstrateMemoryEfficiency(smallJSON)
	
	fmt.Println("\n🔧 Best Practices")
	fmt.Println("------------------")
	
	// Best practices examples
	demonstrateBestPractices(largeJSON)

	fmt.Println("\n⚡ Advanced Performance Tips")
	fmt.Println("-----------------------------")
	
	// Advanced performance tips
	demonstrateAdvancedTips(largeJSON)
}

func generateLargeJSON() string {
	// Generate a larger JSON structure for testing
	users := make([]string, 100)
	for i := 0; i < 100; i++ {
		users[i] = fmt.Sprintf(`{
			"id": %d,
			"name": "User %d",
			"email": "user%d@example.com",
			"profile": {
				"settings": {
					"theme": "%s",
					"notifications": %t
				}
			}
		}`, i, i, i, []string{"dark", "light"}[i%2], i%2 == 0)
	}
	
	return fmt.Sprintf(`{
		"users": [%s],
		"metadata": {
			"total": %d,
			"version": "1.0.0"
		}
	}`, strings.Join(users, ","), len(users))
}

func benchmarkExtraction(jsonData string) {
	fmt.Println("Extraction Performance:")
	
	// Enhanced API
	fmt.Println("✅ Enhanced API:")
	ctx1 := fj.Extract(jsonData, "users.0.name")
	fmt.Printf("   - Extract(): %s\n", ctx1.String())
	
	contexts := fj.ExtractMultiple(jsonData, "users.0.name", "users.0.email", "metadata.total")
	fmt.Printf("   - ExtractMultiple(): %d values extracted\n", len(contexts))
	
	// Original API (still works)
	fmt.Println("✅ Original API (backward compatible):")
	ctx2 := fj.Get(jsonData, "users.0.name")  
	fmt.Printf("   - Get(): %s\n", ctx2.String())
	
	contexts2 := fj.GetMul(jsonData, "users.0.name", "users.0.email", "metadata.total")
	fmt.Printf("   - GetMul(): %d values extracted\n", len(contexts2))
}

func demonstrateMemoryEfficiency(jsonData string) {
	// Use PathBuilder for reusable paths
	fmt.Println("Memory-efficient path building:")
	
	pathBuilder := fj.NewPathBuilder()
	
	// Reuse the same builder
	userPath := pathBuilder.Root("user").Field("name").Build()
	fmt.Printf("   Built path: %s\n", userPath)
	
	agePath := pathBuilder.Reset().Root("user").Field("age").Build()
	fmt.Printf("   Built path: %s\n", agePath)
	
	// Factory pattern for efficient context creation
	fmt.Println("\nFactory-based context creation:")
	factory := fj.NewContextFactory()
	
	ctx1 := factory.Create(jsonData, userPath)
	ctx2 := factory.Create(jsonData, agePath)
	
	fmt.Printf("   Name: %s, Age: %d\n", ctx1.String(), ctx2.Int64())
}

func demonstrateBestPractices(jsonData string) {
	// 1. Validate before processing
	fmt.Println("1. Always validate JSON:")
	if fj.ValidateString(jsonData) {
		fmt.Println("   ✅ JSON is valid, safe to process")
	} else {
		fmt.Println("   ❌ Invalid JSON detected")
		return
	}
	
	// 2. Use appropriate extraction methods
	fmt.Println("\n2. Choose appropriate extraction methods:")
	
	// Single value
	name := fj.Extract(jsonData, "users.0.name")
	fmt.Printf("   Single value: %s\n", name.String())
	
	// Multiple values
	values := fj.ExtractMultiple(jsonData, "users.0.name", "users.0.email")
	fmt.Printf("   Multiple values: %d extracted\n", len(values))
	
	// 3. Efficient transformation chains
	fmt.Println("\n3. Efficient transformation chains:")
	
	// Register reusable transformers
	fj.Register("normalize", func(json, arg string) string {
		// Remove quotes, trim, lowercase
		normalized := strings.ToLower(strings.TrimSpace(json))
		if len(normalized) > 2 && normalized[0] == '"' && normalized[len(normalized)-1] == '"' {
			normalized = normalized[1 : len(normalized)-1]
		}
		return "\"" + normalized + "\""
	})
	
	// Use transformation chains
	email := fj.Extract(jsonData, "users.0.email")
	normalized := fj.Transform(email.String(), "normalize")
	fmt.Printf("   Original: %s\n", email.String())
	fmt.Printf("   Normalized: %s\n", normalized)
}

func demonstrateAdvancedTips(jsonData string) {
	// 1. Reuse processor instances
	fmt.Println("1. Reuse processor instances:")
	processor := fj.NewPipeProcessor()
	
	ctx := fj.Parse(jsonData)
	result := processor.Process(ctx, "users.0.name")
	fmt.Printf("   Processed result: %s\n", result.String())
	
	// 2. Use singleton registry
	fmt.Println("\n2. Use singleton registry efficiently:")
	registry := fj.GetTransformerRegistry()
	
	// Check existence before registration
	if !registry.Exists("customTransform") {
		registry.RegisterFunc("customTransform", func(json, arg string) string {
			return strings.ToUpper(json)
		})
		fmt.Println("   ✅ Custom transformer registered")
	} else {
		fmt.Println("   ✅ Custom transformer already exists")
	}
	
	// 3. Batch operations
	fmt.Println("\n3. Batch operations for efficiency:")
	
	// Extract multiple related values at once
	userPaths := []string{
		"users.0.name",
		"users.0.email", 
		"users.0.profile.settings.theme",
	}
	
	contexts := fj.ExtractMultiple(jsonData, userPaths...)
	fmt.Printf("   Batch extracted %d values:\n", len(contexts))
	for i, ctx := range contexts {
		fmt.Printf("     %s: %s\n", userPaths[i], ctx.String())
	}
	
	// 4. Strategy pattern for flexibility
	fmt.Println("\n4. Strategy pattern for flexible transformations:")
	
	// Create different strategies
	strategies := map[string]*fj.TransformerStrategyChain{
		"cleanup": fj.NewTransformerStrategyChain().Add("normalize"),
		"format":  fj.NewTransformerStrategyChain().Add("normalize"),
	}
	
	for name, strategy := range strategies {
		if strategy.Length() > 0 {
			fmt.Printf("   Strategy '%s': %d transformations\n", name, strategy.Length())
		}
	}
	
	fmt.Println("\n💡 Performance Tips Summary:")
	fmt.Println("   • Validate JSON before processing")
	fmt.Println("   • Reuse PathBuilder instances")
	fmt.Println("   • Use batch operations for multiple values")
	fmt.Println("   • Register transformers once, use multiple times")
	fmt.Println("   • Choose appropriate extraction methods")
	fmt.Println("   • Leverage singleton registry")
}