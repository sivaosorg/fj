package core

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/sivaosorg/fj/pkg/errors"
)

// contextImpl implements the Context interface while maintaining compatibility
// with the existing Context struct from the root package.
type contextImpl struct {
	// Embedded fields from the original Context struct (from types.go)
	kind        Type
	unprocessed string
	strings     string
	numeric     float64
	index       int
	indexes     []int
	err         error
	
	// Additional fields for enhanced functionality
	parser      JSONParser
	transformer JSONTransformer
}

// NewContext creates a new Context implementation
func NewContext(kind Type, unprocessed, strings string, numeric float64, parser JSONParser) Context {
	return &contextImpl{
		kind:        kind,
		unprocessed: unprocessed,
		strings:     strings,
		numeric:     numeric,
		parser:      parser,
	}
}

// NewContextFromExisting creates a new Context from existing data
func NewContextFromExisting(kind Type, unprocessed string, parser JSONParser) Context {
	return &contextImpl{
		kind:        kind,
		unprocessed: unprocessed,
		parser:      parser,
	}
}

// Type returns the JSON type of this context
func (c *contextImpl) Type() Type {
	return c.kind
}

// Kind returns the JSON type (alias for Type for compatibility)
func (c *contextImpl) Kind() Type {
	return c.kind
}

// Exists returns true if the value exists (i.e., it is not null and contains data)
func (c *contextImpl) Exists() bool {
	return c.kind != 0 && c.unprocessed != ""
}

// IsObject checks if the current context represents a JSON object
func (c *contextImpl) IsObject() bool {
	return c.kind == 6 // JSON type from const.go
}

// IsArray checks if the current context represents a JSON array
func (c *contextImpl) IsArray() bool {
	return c.kind == 6 && len(c.unprocessed) > 0 && c.unprocessed[0] == '['
}

// IsString checks if the current context represents a JSON string
func (c *contextImpl) IsString() bool {
	return c.kind == 4 // String type from const.go
}

// IsNumber checks if the current context represents a JSON number
func (c *contextImpl) IsNumber() bool {
	return c.kind == 2 // Number type from const.go
}

// IsBool checks if the current context represents a JSON boolean
func (c *contextImpl) IsBool() bool {
	return c.kind == 1 || c.kind == 5 // False or True type from const.go
}

// IsNull checks if the current context represents a JSON null value
func (c *contextImpl) IsNull() bool {
	return c.kind == 0 // Null type from const.go
}

// String returns the string representation of the context
func (c *contextImpl) String() string {
	switch c.kind {
	case 4: // String
		return c.strings
	case 2: // Number
		return strconv.FormatFloat(c.numeric, 'f', -1, 64)
	case 1: // False
		return "false"
	case 5: // True
		return "true"
	case 0: // Null
		return "null"
	case 6: // JSON
		return c.unprocessed
	default:
		return c.unprocessed
	}
}

// Int returns the integer representation of the context
func (c *contextImpl) Int() int64 {
	if c.kind == 2 { // Number
		return int64(c.numeric)
	}
	if c.kind == 4 { // String
		if val, err := strconv.ParseInt(c.strings, 10, 64); err == nil {
			return val
		}
	}
	return 0
}

// Float returns the float representation of the context
func (c *contextImpl) Float() float64 {
	if c.kind == 2 { // Number
		return c.numeric
	}
	if c.kind == 4 { // String
		if val, err := strconv.ParseFloat(c.strings, 64); err == nil {
			return val
		}
	}
	return 0
}

// Bool returns the boolean representation of the context
func (c *contextImpl) Bool() bool {
	switch c.kind {
	case 5: // True
		return true
	case 1: // False
		return false
	case 4: // String
		return c.strings != ""
	case 2: // Number
		return c.numeric != 0
	case 6: // JSON
		return c.unprocessed != ""
	default:
		return false
	}
}

// Time converts the context value into a time.Time representation
func (c *contextImpl) Time() time.Time {
	if c.kind == 4 { // String
		// Try common time formats
		formats := []string{
			time.RFC3339,
			time.RFC3339Nano,
			"2006-01-02T15:04:05Z",
			"2006-01-02 15:04:05",
			"2006-01-02",
		}
		
		for _, format := range formats {
			if t, err := time.Parse(format, c.strings); err == nil {
				return t
			}
		}
	} else if c.kind == 2 { // Number - assume Unix timestamp
		return time.Unix(int64(c.numeric), 0)
	}
	return time.Time{}
}

// TimeWithFormat parses the context value into a time.Time using the specified format
func (c *contextImpl) TimeWithFormat(format string) (time.Time, error) {
	if c.kind != 4 { // Must be string
		return time.Time{}, errors.NewValidationError("time parsing requires string value")
	}
	return time.Parse(format, c.strings)
}

// Value returns the corresponding Go type for the JSON value
func (c *contextImpl) Value() interface{} {
	switch c.kind {
	case 0: // Null
		return nil
	case 1: // False
		return false
	case 5: // True
		return true
	case 2: // Number
		// Try to return as int if it's a whole number
		if c.numeric == float64(int64(c.numeric)) {
			return int64(c.numeric)
		}
		return c.numeric
	case 4: // String
		return c.strings
	case 6: // JSON
		// Try to parse as Go types
		var result interface{}
		if err := json.Unmarshal([]byte(c.unprocessed), &result); err == nil {
			return result
		}
		return c.unprocessed
	default:
		return c.unprocessed
	}
}

// Array returns an array of Context values
func (c *contextImpl) Array() []Context {
	if !c.IsArray() {
		return []Context{c}
	}
	
	// For now, return empty slice - will be implemented with the existing array parsing logic
	return []Context{}
}

// Map returns a map of string to Context values
func (c *contextImpl) Map() map[string]Context {
	if !c.IsObject() {
		return nil
	}
	
	// For now, return empty map - will be implemented with the existing map parsing logic
	return make(map[string]Context)
}

// Len returns the length of arrays or objects, or 0 for other types
func (c *contextImpl) Len() int {
	if c.IsArray() {
		return len(c.Array())
	}
	if c.IsObject() {
		return len(c.Map())
	}
	return 0
}

// Get retrieves a value using a path query
func (c *contextImpl) Get(path string) Context {
	if c.parser != nil {
		// Use the parser's query functionality
		// This will be implemented when we integrate with the existing Get function
		return NewContextFromExisting(0, "", c.parser) // placeholder
	}
	return c
}

// GetMultiple retrieves multiple values using path queries
func (c *contextImpl) GetMultiple(paths ...string) []Context {
	results := make([]Context, len(paths))
	for i, path := range paths {
		results[i] = c.Get(path)
	}
	return results
}

// PathExists checks if a path exists in this context
func (c *contextImpl) PathExists(path string) bool {
	result := c.Get(path)
	return result.Exists()
}

// Foreach iterates over the context (for objects and arrays)
func (c *contextImpl) Foreach(iterator func(key, value Context) bool) {
	// For now, this is a placeholder - will be implemented with existing foreach logic
}

// ForEachValue iterates over values only (for arrays)
func (c *contextImpl) ForEachValue(iterator func(value Context) bool) {
	c.Foreach(func(key, value Context) bool {
		return iterator(value)
	})
}

// Transform applies a transformation to this context
func (c *contextImpl) Transform(transformerName string, args ...string) Context {
	// For now, this is a placeholder - will be implemented with transformer integration
	return c
}

// Pretty returns a pretty-printed JSON representation
func (c *contextImpl) Pretty() string {
	if c.IsObject() || c.IsArray() {
		var result interface{}
		if err := json.Unmarshal([]byte(c.unprocessed), &result); err == nil {
			if prettified, err := json.MarshalIndent(result, "", "  "); err == nil {
				return string(prettified)
			}
		}
	}
	return c.String()
}

// Minify returns a minified JSON representation
func (c *contextImpl) Minify() string {
	if c.IsObject() || c.IsArray() {
		var result interface{}
		if err := json.Unmarshal([]byte(c.unprocessed), &result); err == nil {
			if minified, err := json.Marshal(result); err == nil {
				return string(minified)
			}
		}
	}
	return c.String()
}

// Less compares this context with another for sorting
func (c *contextImpl) Less(other Context) bool {
	if c.Type() != other.Type() {
		return c.Type() < other.Type()
	}
	
	switch c.Type() {
	case 2: // Number
		return c.Float() < other.Float()
	case 4: // String
		return strings.ToLower(c.String()) < strings.ToLower(other.String())
	default:
		return c.String() < other.String()
	}
}

// Equal checks if this context equals another
func (c *contextImpl) Equal(other Context) bool {
	if c.Type() != other.Type() {
		return false
	}
	
	switch c.Type() {
	case 2: // Number
		return c.Float() == other.Float()
	case 4: // String
		return c.String() == other.String()
	default:
		return c.String() == other.String()
	}
}

// Error returns any error associated with this context
func (c *contextImpl) Error() error {
	return c.err
}

// HasError returns true if this context has an error
func (c *contextImpl) HasError() bool {
	return c.err != nil
}

// Raw returns the raw JSON string
func (c *contextImpl) Raw() string {
	return c.unprocessed
}

// Bytes returns the raw JSON as bytes
func (c *contextImpl) Bytes() []byte {
	return []byte(c.unprocessed)
}

// Index returns the character index position in the original JSON
func (c *contextImpl) Index() int {
	return c.index
}

// Helper method to create a context with error
func NewContextWithError(err error) Context {
	return &contextImpl{
		err: err,
	}
}