# fj - Fast JSON Processing Library

**fj** (_Fast JSON_) is a modern Go package that provides fast, flexible, and professional JSON processing capabilities with Object-Oriented design patterns and clean APIs.

[![Go Version](https://img.shields.io/badge/Go-1.19+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/tests-169%2B%20passing-brightgreen.svg)](#testing)

## 🚀 Features

- **Object-Oriented Architecture** - Built with modern OOP design patterns
- **Multiple Design Patterns** - Builder, Strategy, Factory, Singleton, Chain of Responsibility
- **Professional APIs** - Clean, intuitive function names following Go conventions
- **Full Backward Compatibility** - All existing APIs preserved
- **High Performance** - Optimized for speed and memory efficiency
- **Comprehensive Testing** - 169+ tests ensuring reliability
- **Rich Documentation** - Extensive examples and usage patterns

## 📦 Installation

```bash
# Latest version
go get -u github.com/sivaosorg/fj@latest

# Specific version
go get github.com/sivaosorg/fj@v0.0.1
```

**Requirements:** Go 1.19 or higher

## 🏗️ Architecture Overview

fj implements several design patterns for clean, maintainable code:

### Design Patterns Implemented

| Pattern | Implementation | Purpose |
|---------|----------------|---------|
| **Builder** | `PathBuilder` | Fluent JSON path construction |
| **Strategy** | `TransformerStrategy` | Pluggable transformation algorithms |
| **Factory** | `ContextFactory` | Flexible Context object creation |
| **Singleton** | `TransformerRegistry` | Centralized transformer management |
| **Chain of Responsibility** | `PipeProcessor` | Sequential operation processing |

## 🎯 Quick Start

### Basic Usage (Enhanced API)

```go
package main

import (
    "fmt"
    "github.com/sivaosorg/fj"
)

func main() {
    json := `{
        "user": {
            "name": "John Doe",
            "age": 30,
            "roles": ["admin", "user"]
        }
    }`

    // Extract single values
    name := fj.Extract(json, "user.name")
    fmt.Println(name.String()) // "John Doe"

    // Extract multiple values
    values := fj.ExtractMultiple(json, "user.name", "user.age")
    fmt.Println(values[0].String()) // "John Doe" 
    fmt.Println(values[1].Int64())  // 30

    // Validate JSON
    if fj.ValidateString(json) {
        fmt.Println("Valid JSON!")
    }
}
```

### Using Design Patterns

#### Builder Pattern - Fluent Path Construction

```go
// Build complex paths fluently
path := fj.NewPathBuilder().
    Root("users").
    ArrayIndex(0).
    Field("profile").
    Field("settings").
    Build()

fmt.Println(path) // "users.0.profile.settings"

// Use the path
ctx := fj.Extract(json, path)
```

#### Strategy Pattern - Pluggable Transformations

```go
// Create transformation strategies
strategy := fj.NewTransformerStrategy("uppercase")
result := strategy.Execute(`{"name": "john"}`, "")

// Chain multiple strategies
chain := fj.NewTransformerStrategyChain().
    Add("trim").
    Add("uppercase").
    Add("minify")

transformed := chain.Execute(`{"name": " john "}`, "")
```

#### Factory Pattern - Flexible Context Creation

```go
factory := fj.NewContextFactory()

// Create from different sources
ctx1 := factory.Create(`{"name": "John"}`, "name")
ctx2 := factory.CreateFromBytes([]byte(`{"age": 30}`), "age")
ctx3 := factory.CreateEmpty()

// Or use the convenience function
ctx4 := fj.From(`{"data": "value"}`, "data")
```

#### Singleton Registry - Centralized Management

```go
// Get the global registry
registry := fj.GetTransformerRegistry()

// Register custom transformers
err := registry.RegisterFunc("myTransform", func(json, arg string) string {
    return strings.ToUpper(json)
})

// Check if transformer exists
if registry.Exists("myTransform") {
    transformer, _ := registry.Get("myTransform")
    result := transformer.Transform("hello", "")
}
```

#### Chain of Responsibility - Pipe Processing

```go
processor := fj.NewPipeProcessor()

// Process operations in sequence
ctx := fj.Parse(`{"user": {"name": "John"}}`)
result := processor.Process(ctx, "user|name|uppercase")

// Or use the fluent builder
customProcessor := fj.NewChainedPipeProcessor().
    WithTransformers().
    WithPaths().
    WithQueries().
    Build()
```

## 🔧 Advanced Features

### Enhanced API Functions

| Enhanced Function | Original Function | Description |
|------------------|-------------------|-------------|
| `Extract()` | `Get()` | Extract JSON values with cleaner naming |
| `ExtractBytes()` | `GetBytes()` | Extract from byte slices |
| `ExtractMultiple()` | `GetMul()` | Extract multiple values at once |
| `Validate()` | `IsValidJSONBytes()` | Validate JSON data |
| `Transform()` | N/A | Apply single transformations |
| `Chain()` | N/A | Apply transformation chains |
| `Register()` | `AddTransformer()` | Register transformers |
| `From()` | N/A | Flexible Context creation |

### Convenience Functions

```go
// Simple path building
path := fj.Build("user", "profile", "name") // "user.profile.name"

// Quick transformer registration
fj.Register("reverse", func(json, arg string) string {
    return reverseString(json)
})

// Check transformer existence
if fj.Exists("reverse") {
    result := fj.Transform(json, "reverse")
}

// Process pipe operations
result := fj.Process(ctx, "field1|uppercase|trim")
```

### Transformation Chains

```go
// Apply multiple transformations
result := fj.Chain(json, "trim", "uppercase", "minify")

// Parse transformer chain from string
chain := fj.ParseTransformerChain("trim|uppercase|pretty")
result := chain.Execute(json)

// Apply transformation chain directly
result := fj.ApplyTransformerChain(json, "trim|uppercase|minify")
```

## 📊 Built-in Transformers

fj comes with 30+ built-in transformers:

| Transformer | Description | Example Usage |
|-------------|-------------|---------------|
| `@pretty` | Format JSON with indentation | `@pretty:{"indent":"  "}` |
| `@minify` | Minimize JSON size | `@minify` |
| `@uppercase` | Convert to uppercase | `@uppercase` |
| `@lowercase` | Convert to lowercase | `@lowercase` |
| `@trim` | Remove whitespace | `@trim` |
| `@reverse` | Reverse string/array | `@reverse` |
| `@keys` | Extract object keys | `@keys` |
| `@values` | Extract object values | `@values` |
| `@flatten` | Flatten nested arrays | `@flatten:{"deep":true}` |
| `@group` | Group array elements | `@group` |

[View complete transformer list →](docs/transformers.md)

## 🧪 Testing and Quality

fj maintains high quality standards with comprehensive testing:

```bash
# Run all tests
go test -v

# Run with coverage
go test -v -cover

# Run benchmarks
go test -v -bench=.
```

**Test Coverage:** 169+ tests covering all functionality including:
- Core JSON processing
- All design patterns
- Enhanced API functions
- Backward compatibility
- Error handling
- Performance benchmarks

## 📈 Performance

fj is optimized for high performance:

```go
// Benchmark results (example)
BenchmarkExtract-8       1000000    1234 ns/op    456 B/op    7 allocs/op
BenchmarkTransform-8     500000     2345 ns/op    789 B/op    12 allocs/op
BenchmarkChain-8         300000     3456 ns/op    1024 B/op   15 allocs/op
```

## 🔄 Migration Guide

### Upgrading from Previous Versions

The library maintains full backward compatibility. However, we recommend migrating to the enhanced API:

```go
// ❌ Old way (still works, but deprecated)
ctx := fj.GetBytes(data, "path")
valid := fj.IsValidJSONBytes(data)
fj.AddTransformer("name", transformerFunc)

// ✅ New way (recommended)
ctx := fj.ExtractBytes(data, "path")
valid := fj.Validate(data)
fj.Register("name", transformerFunc)
```

### Deprecation Timeline

- **v1.x**: All original functions available with deprecation notices
- **v2.0**: Original functions removed, enhanced API becomes standard

## 🛠️ Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/sivaosorg/fj.git
cd fj

# Install dependencies
go mod tidy

# Run tests
go test -v

# Run linter
golangci-lint run
```

## 📚 Documentation

- [API Reference](docs/api.md)
- [Design Patterns Guide](docs/patterns.md)
- [Transformer Reference](docs/transformers.md)
- [Performance Guide](docs/performance.md)
- [Migration Guide](docs/migration.md)

## 🤝 Support

- **Issues**: [GitHub Issues](https://github.com/sivaosorg/fj/issues)
- **Discussions**: [GitHub Discussions](https://github.com/sivaosorg/fj/discussions)
- **Documentation**: [Wiki](https://github.com/sivaosorg/fj/wiki)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Thanks to all contributors who helped improve this library
- Inspired by modern Go development practices and clean architecture principles
- Built with performance and developer experience in mind

---

<p align="center">Made with ❤️ by the fj team</p>