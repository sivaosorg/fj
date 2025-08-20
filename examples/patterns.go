package main

import (
	"fmt"
	"github.com/sivaosorg/fj"
)

func main() {
	fmt.Println("=== fj Design Patterns Examples ===\n")

	// Sample JSON data
	jsonData := `{
		"users": [
			{
				"id": 1,
				"name": "john doe",
				"email": " john@example.com ",
				"profile": {
					"settings": {
						"theme": "dark",
						"notifications": true
					}
				},
				"roles": ["admin", "user"]
			},
			{
				"id": 2,
				"name": "jane smith",
				"email": " jane@example.com ",
				"profile": {
					"settings": {
						"theme": "light",
						"notifications": false
					}
				},
				"roles": ["user"]
			}
		],
		"metadata": {
			"total": 2,
			"version": "1.0.0"
		}
	}`

	// ===== Builder Pattern Example =====
	fmt.Println("🏗️  Builder Pattern - Fluent Path Construction")
	
	// Build complex paths fluently
	userPath := fj.NewPathBuilder().
		Root("users").
		ArrayIndex(0).
		Field("profile").
		Field("settings").
		Field("theme").
		Build()
	fmt.Printf("Built path: %s\n", userPath)
	
	userTheme := fj.Extract(jsonData, userPath)
	fmt.Printf("User theme: %s\n\n", userTheme.String())

	// ===== Strategy Pattern Example =====
	fmt.Println("🎯 Strategy Pattern - Pluggable Transformations")
	
	// Register custom transformers
	fj.Register("emailClean", func(json, arg string) string {
		// Remove quotes and trim whitespace for email cleaning
		cleaned := json
		if len(cleaned) > 2 && cleaned[0] == '"' && cleaned[len(cleaned)-1] == '"' {
			cleaned = cleaned[1 : len(cleaned)-1]
		}
		// Simulate email cleaning
		return "\"" + "cleaned:" + cleaned + "\""
	})

	// Create transformation strategy
	emailStrategy := fj.NewTransformerStrategy("emailClean")
	if emailStrategy != nil {
		email := fj.Extract(jsonData, "users.0.email")
		cleanedEmail := emailStrategy.Execute(email.String(), "")
		fmt.Printf("Original email: %s\n", email.String())
		fmt.Printf("Cleaned email: %s\n\n", cleanedEmail)
	}

	// ===== Chain of Responsibility Example =====
	fmt.Println("⛓️  Chain of Responsibility - Transformation Chains")
	
	// Create transformation chain
	chain := fj.NewTransformerStrategyChain()
	
	// Register chain transformers
	fj.Register("addPrefix", func(json, arg string) string {
		return "\"PREFIX:" + json[1:len(json)-1] + "\""
	})
	
	fj.Register("addSuffix", func(json, arg string) string {
		return json[:len(json)-1] + ":SUFFIX\""
	})

	chain.Add("addPrefix").Add("addSuffix")
	
	name := fj.Extract(jsonData, "users.0.name")
	chainResult := chain.Execute(name.String())
	fmt.Printf("Original name: %s\n", name.String())
	fmt.Printf("Chain result: %s\n\n", chainResult)

	// ===== Factory Pattern Example =====
	fmt.Println("🏭 Factory Pattern - Flexible Context Creation")
	
	factory := fj.NewContextFactory()
	
	// Create contexts from different sources
	ctx1 := factory.Create(jsonData, "metadata.total")
	ctx2 := factory.CreateFromBytes([]byte(jsonData), "metadata.version")
	
	fmt.Printf("Total users: %d\n", ctx1.Int64())
	fmt.Printf("API version: %s\n", ctx2.String())
	
	// Use convenience From() function
	ctx3 := fj.From(jsonData, "users.1.name")
	fmt.Printf("Second user: %s\n\n", ctx3.String())

	// ===== Singleton Registry Example =====
	fmt.Println("🔄 Singleton Registry - Centralized Management")
	
	registry := fj.GetTransformerRegistry()
	
	// Show built-in transformers
	transformers := registry.List()
	fmt.Printf("Available transformers: %d\n", len(transformers))
	
	// Register and use custom transformer
	err := registry.RegisterFunc("shout", func(json, arg string) string {
		return "\"" + json[1:len(json)-1] + "!!!\""
	})
	if err == nil {
		if registry.Exists("shout") {
			transformer, _ := registry.Get("shout")
			result := transformer.Transform(`"hello"`, "")
			fmt.Printf("Shout transformer result: %s\n", result)
		}
	}

	fmt.Println()

	// ===== Enhanced API Examples =====
	fmt.Println("✨ Enhanced API - Clean Function Names")
	
	// Multiple extraction
	values := fj.ExtractMultiple(jsonData, "users.0.name", "users.1.name")
	fmt.Println("Multiple extraction:")
	for i, val := range values {
		fmt.Printf("  User %d: %s\n", i+1, val.String())
	}
	
	// Validation
	if fj.ValidateString(jsonData) {
		fmt.Println("✅ JSON is valid")
	}
	
	// Quick transformations
	result := fj.Transform(`"hello world"`, "shout")
	fmt.Printf("Quick transform: %s\n", result)
	
	// Path building convenience
	quickPath := fj.Build("users", "0", "roles", "0")
	role := fj.Extract(jsonData, quickPath)
	fmt.Printf("First user role: %s\n", role.String())

	fmt.Println("\n=== Examples Complete ===")
}