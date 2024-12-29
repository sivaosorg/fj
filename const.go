package fj

import "regexp"

// DisableTransformers is a global flag that determines whether transformers should be applied
// when processing JSON values. If set to true, transformers will not be applied to the JSON values.
// If set to false, transformers will be applied as expected.
var DisableTransformers = false

// jsonTransformers is a map that associates a string key (the transformer type) with a function that
// takes two string arguments (`json` and `arg`), and returns a modified string. The map is used
// to apply various transformations to JSON data based on the specified jsonTransformers.
var jsonTransformers map[string]func(json, arg string) string

// hexDigits is an array of bytes representing the hexadecimal digits used in JSON encoding.
// It contains the characters '0' to '9' and 'a' to 'f', which are used for encoding hexadecimal numbers.
// This is commonly used for encoding special characters or byte sequences in JSON strings (e.g., for Unicode escape sequences).
var hexDigits = [...]byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'a', 'b', 'c', 'd', 'e', 'f',
}

const (
	// Null is a constant representing a JSON null value.
	// In JSON, null is used to represent the absence of a value.
	Null Type = iota
	// False is a constant representing a JSON false boolean value.
	// In JSON, false is a boolean value that represents a negative or off state.
	False
	// Number is a constant representing a JSON number value.
	// In JSON, numbers can be integers or floating-point values.
	Number
	// String is a constant representing a JSON string value.
	// In JSON, strings are sequences of characters enclosed in double quotes.
	String
	// True is a constant representing a JSON true boolean value.
	// In JSON, true is a boolean value that represents a positive or on state.
	True
	// JSON is a constant representing a raw JSON block.
	// This type can be used to represent any valid JSON object or array.
	JSON
)

var (
	// RegexpDupSpaces is a precompiled regular expression that matches one or more consecutive
	// whitespace characters (including spaces, tabs, and newlines). This can be used for tasks
	// such as normalizing whitespace in strings by replacing multiple whitespace characters
	// with a single space, or for validating string formats where excessive whitespace should
	// be trimmed or removed.
	RegexpDupSpaces = regexp.MustCompile(`\s+`)
)
