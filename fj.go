package fj

import (
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/sivaosorg/unify4g"
)

// EscapeUnsafeChars processes a string `component` to escape characters that are not considered safe
// according to the `isSafeKeyChar` function. It inserts a backslash (`\`) before each unsafe
// character, ensuring that the resulting string contains only safe characters.
//
// Parameters:
//   - `component`: A string that may contain unsafe characters that need to be escaped.
//
// Returns:
//   - A new string with unsafe characters escaped by prefixing them with a backslash (`\`).
//
// Notes:
//   - The function iterates through the input string and checks each character using the
//     `isSafeKeyChar` function. When it encounters an unsafe character, it escapes it with a backslash.
//   - Once an unsafe character is found, the function adds a backslash before each subsequent unsafe character
//     and continues until the end of the string.
//
// Example:
//
//	component := "key-with$pecial*chars"
//	escaped := EscapeUnsafeChars(component) // escaped: "key-with\$pecial\*chars"
func EscapeUnsafeChars(component string) string {
	for i := 0; i < len(component); i++ {
		if !isSafeKeyChar(component[i]) {
			noneComponent := []byte(component[:i])
			for ; i < len(component); i++ {
				if !isSafeKeyChar(component[i]) {
					noneComponent = append(noneComponent, '\\')
				}
				noneComponent = append(noneComponent, component[i])
			}
			return string(noneComponent)
		}
	}
	return component
}

// AppendJSON converts a given string into a valid JSON string format
// and appends it to the provided byte slice `dst`.
//
// This function escapes special characters in the input string `s` to ensure
// that it adheres to the JSON string encoding rules, such as escaping double
// quotes, backslashes, and control characters. Additionally, it handles UTF-8
// characters and appends them in their proper encoded format.
//
// Parameters:
//   - dst: A byte slice to which the encoded JSON string will be appended.
//   - s: The input string to be converted into JSON string format.
//
// Returns:
//   - []byte: The resulting byte slice containing the original content of `dst`
//     with the JSON-encoded string appended.
//
// Details:
//   - The function begins by appending space for the string `s` and wrapping
//     it in double quotes.
//   - It iterates through the input string `s` character by character and checks
//     for specific cases where escaping or additional encoding is required:
//   - Control characters (`\n`, `\r`, `\t`) are replaced with their escape
//     sequences (`\\n`, `\\r`, `\\t`).
//   - Characters like `<`, `>`, and `&` are escaped using Unicode notation
//     to ensure the resulting JSON string is safe for embedding in HTML or XML.
//   - Backslashes (`\`) and double quotes (`"`) are escaped with a preceding
//     backslash (`\\`).
//   - UTF-8 characters are properly encoded, and unsupported characters or
//     decoding errors are replaced with the Unicode replacement character
//     (`\ufffd`).
//
// Example Usage:
//
//	dst := []byte("Current JSON: ")
//	s := "Hello \"world\"\nLine break!"
//	result := AppendJSON(dst, s)
//	// result: []byte(`Current JSON: "Hello \"world\"\nLine break!"`)
//
// Notes:
//   - This function is useful for building JSON-encoded strings dynamically
//     without allocating new memory for each operation.
//   - It ensures that the resulting JSON string is safe and adheres to
//     encoding rules for use in various contexts such as web APIs or
//     configuration files.
func AppendJSON(target []byte, s string) []byte {
	target = append(target, make([]byte, len(s)+2)...)
	target = append(target[:len(target)-len(s)-2], '"')
	for i := 0; i < len(s); i++ {
		if s[i] < ' ' {
			target = append(target, '\\')
			switch s[i] {
			case '\n':
				target = append(target, 'n')
			case '\r':
				target = append(target, 'r')
			case '\t':
				target = append(target, 't')
			default:
				target = append(target, 'u')
				target = appendHex16(target, uint16(s[i]))
			}
		} else if s[i] == '>' || s[i] == '<' || s[i] == '&' {
			target = append(target, '\\', 'u')
			target = appendHex16(target, uint16(s[i]))
		} else if s[i] == '\\' {
			target = append(target, '\\', '\\')
		} else if s[i] == '"' {
			target = append(target, '\\', '"')
		} else if s[i] > 127 {
			r, n := utf8.DecodeRuneInString(s[i:]) // read utf8 character
			if n == 0 {
				break
			}
			if r == utf8.RuneError && n == 1 {
				target = append(target, `\ufffd`...)
			} else if r == '\u2028' || r == '\u2029' {
				target = append(target, `\u202`...)
				target = append(target, hexDigits[r&0xF])
			} else {
				target = append(target, s[i:i+n]...)
			}
			i = i + n - 1
		} else {
			target = append(target, s[i])
		}
	}
	return append(target, '"')
}

// String provides a string representation of the `Type` enumeration.
//
// This method converts the `Type` value into a human-readable string.
// It is particularly useful for debugging or logging purposes.
//
// Mapping of `Type` values to strings:
//   - Null: "Null"
//   - False: "False"
//   - Number: "Number"
//   - String: "String"
//   - True: "True"
//   - JSON: "JSON"
//   - Default (unknown type): An empty string is returned.
//
// Returns:
//   - string: A string representation of the `Type` value.
//
// Example Usage:
//
//	var t Type = True
//	fmt.Println(t.String())  // Output: "True"
func (t Type) String() string {
	switch t {
	default:
		return ""
	case Null:
		return "Null"
	case False:
		return "False"
	case Number:
		return "Number"
	case String:
		return "String"
	case True:
		return "True"
	case JSON:
		return "JSON"
	}
}

// Kind returns the JSON type of the Context.
// It provides the specific type of the JSON value, such as String, Number, Object, etc.
//
// Returns:
//   - Type: The type of the JSON value represented by the Context.
func (ctx Context) Kind() Type {
	return ctx.kind
}

// Unprocessed returns the raw, unprocessed JSON string for the Context.
// This can be useful for inspecting the original data without any parsing or transformations.
//
// Returns:
//   - string: The unprocessed JSON string.
func (ctx Context) Unprocessed() string {
	return ctx.unprocessed
}

// Numeric returns the numeric value of the Context, if applicable.
// This is relevant when the Context represents a JSON number.
//
// Returns:
//   - float64: The numeric value of the Context.
//     If the Context does not represent a number, the value may be undefined.
func (ctx Context) Numeric() float64 {
	return ctx.numeric
}

// Index returns the index of the unprocessed JSON value in the original JSON string.
// This can be used to track the position of the value in the source data.
// If the index is unknown, it defaults to 0.
//
// Returns:
//   - int: The position of the value in the original JSON string.
func (ctx Context) Index() int {
	return ctx.index
}

// Indexes returns a slice of indices for elements matching a path containing the '#' character.
// This is useful for handling path queries that involve multiple matches.
//
// Returns:
//   - []int: A slice of indices for matching elements.
func (ctx Context) Indexes() []int {
	return ctx.indexes
}

// String returns a string representation of the Context value.
// The output depends on the JSON type of the Context:
//   - For `False` type: Returns "false".
//   - For `True` type: Returns "true".
//   - For `Number` type: Returns the numeric value as a string.
//     If the numeric value was calculated, it formats the float value.
//     Otherwise, it preserves the original unprocessed string if valid.
//   - For `String` type: Returns the string value.
//   - For `JSON` type: Returns the raw unprocessed JSON string.
//   - For other types: Returns an empty string.
//
// Returns:
//   - string: A string representation of the Context value.
func (ctx Context) String() string {
	switch ctx.kind {
	default:
		return ""
	case False:
		return "false"
	case Number:
		if len(ctx.unprocessed) == 0 {
			return strconv.FormatFloat(ctx.numeric, 'f', -1, 64)
		}
		var i int
		if ctx.unprocessed[0] == '-' {
			i++
		}
		for ; i < len(ctx.unprocessed); i++ {
			if ctx.unprocessed[i] < '0' || ctx.unprocessed[i] > '9' {
				return strconv.FormatFloat(ctx.numeric, 'f', -1, 64)
			}
		}
		return ctx.unprocessed
	case String:
		return ctx.strings
	case JSON:
		return ctx.unprocessed
	case True:
		return "true"
	}
}

// Bool converts the Context value into a boolean representation.
// The conversion depends on the JSON type of the Context:
//   - For `True` type: Returns `true`.
//   - For `String` type: Attempts to parse the string as a boolean (case-insensitive).
//     If parsing fails, defaults to `false`.
//   - For `Number` type: Returns `true` if the numeric value is non-zero, otherwise `false`.
//   - For all other types: Returns `false`.
//
// Returns:
//   - bool: A boolean representation of the Context value.
func (ctx Context) Bool() bool {
	switch ctx.kind {
	default:
		return false
	case True:
		return true
	case String:
		b, _ := strconv.ParseBool(strings.ToLower(ctx.strings))
		return b
	case Number:
		return ctx.numeric != 0
	}
}

// Int64 converts the Context value into an integer representation (int64).
// The conversion depends on the JSON type of the Context:
//   - For `True` type: Returns 1.
//   - For `String` type: Attempts to parse the string into an integer. Defaults to 0 on failure.
//   - For `Number` type:
//   - Directly converts the numeric value to an integer if it's safe.
//   - Parses the unprocessed string for integer values as a fallback.
//   - Defaults to converting the float64 numeric value to an int64.
//
// Returns:
//   - int64: An integer representation of the Context value.
func (ctx Context) Int64() int64 {
	switch ctx.kind {
	default:
		return 0
	case True:
		return 1
	case String:
		n, _ := parseInt64(ctx.strings)
		return n
	case Number:
		i, ok := ensureSafeInt64(ctx.numeric)
		if ok {
			return i
		}
		i, ok = parseInt64(ctx.unprocessed)
		if ok {
			return i
		}
		return int64(ctx.numeric)
	}
}

// Uint64 converts the Context value into an unsigned integer representation (uint64).
// The conversion depends on the JSON type of the Context:
//   - For `True` type: Returns 1.
//   - For `String` type: Attempts to parse the string into an unsigned integer. Defaults to 0 on failure.
//   - For `Number` type:
//   - Directly converts the numeric value to a uint64 if it's safe and non-negative.
//   - Parses the unprocessed string for unsigned integer values as a fallback.
//   - Defaults to converting the float64 numeric value to a uint64.
//
// Returns:
//   - uint64: An unsigned integer representation of the Context value.
func (ctx Context) Uint64() uint64 {
	switch ctx.kind {
	default:
		return 0
	case True:
		return 1
	case String:
		n, _ := parseUint64(ctx.strings)
		return n
	case Number:
		i, ok := ensureSafeInt64(ctx.numeric)
		if ok && i >= 0 {
			return uint64(i)
		}
		u, ok := parseUint64(ctx.unprocessed)
		if ok {
			return u
		}
		return uint64(ctx.numeric)
	}
}

// Float64 converts the Context value into a floating-point representation (float64).
// The conversion depends on the JSON type of the Context:
//   - For `True` type: Returns 1.
//   - For `String` type: Attempts to parse the string as a floating-point number. Defaults to 0 on failure.
//   - For `Number` type: Returns the numeric value as a float64.
//
// Returns:
//   - float64: A floating-point representation of the Context value.
func (ctx Context) Float64() float64 {
	switch ctx.kind {
	default:
		return 0
	case True:
		return 1
	case String:
		n, _ := strconv.ParseFloat(ctx.strings, 64)
		return n
	case Number:
		return ctx.numeric
	}
}

// Time converts the Context value into a time.Time representation.
// The conversion interprets the Context value as a string in RFC3339 format.
// If parsing fails, the zero time (0001-01-01 00:00:00 UTC) is returned.
//
// Returns:
//   - time.Time: A time.Time representation of the Context value.
//     Defaults to the zero time if parsing fails.
func (ctx Context) Time() time.Time {
	v, _ := time.Parse(time.RFC3339, ctx.String())
	return v
}

// Array returns an array of `Context` values derived from the current `Context`.
//
// Behavior:
//   - If the current `Context` represents a `Null` value, it returns an empty array.
//   - If the current `Context` is not a JSON array, it returns an array containing itself as a single element.
//   - If the current `Context` is a JSON array, it parses and returns the array's elements.
//
// Returns:
//   - []Context: A slice of `Context` values representing the array elements.
//
// Example Usage:
//
//	ctx := Context{kind: Null}
//	arr := ctx.Array()
//	// arr: []
//
//	ctx = Context{kind: JSON, unprocessed: "[1, 2, 3]"}
//	arr = ctx.Array()
//	// arr: [Context, Context, Context]
//
// Notes:
//   - This function uses `parseJSONElements` internally to extract array elements.
//   - If the JSON is malformed or does not represent an array, the behavior may vary.
func (ctx Context) Array() []Context {
	if ctx.kind == Null {
		return []Context{}
	}
	if !ctx.IsArray() {
		return []Context{ctx}
	}
	r := ctx.parseJSONElements('[', false)
	return r.ArrayResult
}

// IsObject checks if the current `Context` represents a JSON object.
//
// A value is considered a JSON object if:
//   - The `kind` is `JSON`.
//   - The `unprocessed` string starts with the `{` character.
//
// Returns:
//   - bool: Returns `true` if the `Context` is a JSON object; otherwise, `false`.
//
// Example Usage:
//
//	ctx := Context{kind: JSON, unprocessed: "{"key": "value"}"}
//	isObj := ctx.IsObject()
//	// isObj: true
//
//	ctx = Context{kind: JSON, unprocessed: "[1, 2, 3]"}
//	isObj = ctx.IsObject()
//	// isObj: false
func (ctx Context) IsObject() bool {
	return ctx.kind == JSON && len(ctx.unprocessed) > 0 && ctx.unprocessed[0] == '{'
}

// IsArray checks if the current `Context` represents a JSON array.
//
// A value is considered a JSON array if:
//   - The `kind` is `JSON`.
//   - The `unprocessed` string starts with the `[` character.
//
// Returns:
//   - bool: Returns `true` if the `Context` is a JSON array; otherwise, `false`.
//
// Example Usage:
//
//	ctx := Context{kind: JSON, unprocessed: "[1, 2, 3]"}
//	isArr := ctx.IsArray()
//	// isArr: true
//
//	ctx = Context{kind: JSON, unprocessed: "{"key": "value"}"}
//	isArr = ctx.IsArray()
//	// isArr: false
func (ctx Context) IsArray() bool {
	return ctx.kind == JSON && len(ctx.unprocessed) > 0 && ctx.unprocessed[0] == '['
}

// IsBool checks if the current `Context` represents a JSON boolean value.
//
// A value is considered a JSON boolean if:
//   - The `kind` is `True` or `False`.
//
// Returns:
//   - bool: Returns `true` if the `Context` is a JSON boolean; otherwise, `false`.
//
// Example Usage:
//
//	ctx := Context{kind: True}
//	isBool := ctx.IsBool()
//	// isBool: true
//
//	ctx = Context{kind: String, strings: "true"}
//	isBool = ctx.IsBool()
//	// isBool: false
func (ctx Context) IsBool() bool {
	return ctx.kind == True || ctx.kind == False
}

// Exists returns true if the value exists (i.e., it is not Null and contains data).
//
// Example Usage:
//
//	if fj.Get(json, "user.name").Exists() {
//	  println("value exists")
//	}
//
// Returns:
//   - bool: Returns true if the value is not null and contains non-empty data, otherwise returns false.
func (ctx Context) Exists() bool {
	return ctx.kind != Null || len(ctx.unprocessed) != 0
}

// Value returns the corresponding Go type for the JSON value represented by the Context.
//
// The function returns one of the following types based on the JSON value:
//   - bool for JSON booleans (True or False)
//   - float64 for JSON numbers
//   - string for JSON string literals
//   - nil for JSON null
//   - map[string]interface{} for JSON objects
//   - []interface{} for JSON arrays
//
// Example Usage:
//
//	value := ctx.Value()
//	switch v := value.(type) {
//	  case bool:
//	    fmt.Println("Boolean:", v)
//	  case float64:
//	    fmt.Println("Number:", v)
//	  case string:
//	    fmt.Println("String:", v)
//	  case nil:
//	    fmt.Println("Null value")
//	  case map[string]interface{}:
//	    fmt.Println("Object:", v)
//	  case []interface{}:
//	    fmt.Println("Array:", v)
//	}
//
// Returns:
//
//   - interface{}: Returns the corresponding Go type for the JSON value, or nil if the type is not recognized.
func (ctx Context) Value() interface{} {
	if ctx.kind == String {
		return ctx.strings
	}
	switch ctx.kind {
	default:
		return nil
	case False:
		return false
	case Number:
		return ctx.numeric
	case JSON:
		r := ctx.parseJSONElements(0, true)
		if r.valueN == '{' {
			return r.OpIns
		} else if r.valueN == '[' {
			return r.ArrayIns
		}
		return nil
	case True:
		return true
	}
}

// ForEach iterates through the values of a JSON object or array, applying the provided iterator function.
//
// If the `Context` represents a non-existent value (Null or invalid JSON), no iteration occurs.
// For JSON objects, the iterator receives both the key and value of each item.
// For JSON arrays, the iterator receives only the value of each item.
// If the `Context` is not an array or object, the iterator is called once with the whole value.
//
// Example Usage:
//
//	ctx.ForEach(func(key, value Context) bool {
//	  if key.strings != "" {
//	    fmt.Printf("Key: %s, Value: %v\n", key.strings, value)
//	  } else {
//	    fmt.Printf("Value: %v\n", value)
//	  }
//	  return true // Continue iteration
//	})
//
// Parameters:
//   - iterator: A function that receives a `key` (for objects) and `value` (for both objects and arrays).
//     The function should return `true` to continue iteration or `false` to stop.
//
// Notes:
//   - If the result is a JSON object, the iterator receives key-value pairs.
//   - If the result is a JSON array, the iterator receives only the values.
//   - If the result is not an object or array, the iterator is invoked once with the value.
//
// Returns:
//   - None. The iteration continues until all items are processed or the iterator returns `false`.
func (ctx Context) ForEach(iterator func(key, value Context) bool) {
	if !ctx.Exists() {
		return
	}
	if ctx.kind != JSON {
		iterator(Context{}, ctx)
		return
	}
	json := ctx.unprocessed
	var obj bool
	var i int
	var key, value Context
	for ; i < len(json); i++ {
		if json[i] == '{' {
			i++
			key.kind = String
			obj = true
			break
		} else if json[i] == '[' {
			i++
			key.kind = Number
			key.numeric = -1
			break
		}
		if json[i] > ' ' {
			return
		}
	}
	var str string
	var _esc bool
	var ok bool
	var idx int
	for ; i < len(json); i++ {
		if obj {
			if json[i] != '"' {
				continue
			}
			s := i
			i, str, _esc, ok = parseString(json, i+1)
			if !ok {
				return
			}
			if _esc {
				key.strings = unescape(str[1 : len(str)-1])
			} else {
				key.strings = str[1 : len(str)-1]
			}
			key.unprocessed = str
			key.index = s + ctx.index
		} else {
			key.numeric += 1
		}
		for ; i < len(json); i++ {
			if json[i] <= ' ' || json[i] == ',' || json[i] == ':' {
				continue
			}
			break
		}
		s := i
		i, value, ok = parseJSONAny(json, i, true)
		if !ok {
			return
		}
		if ctx.indexes != nil {
			if idx < len(ctx.indexes) {
				value.index = ctx.indexes[idx]
			}
		} else {
			value.index = s + ctx.index
		}
		if !iterator(key, value) {
			return
		}
		idx++
	}
}

// Map returns a map of values extracted from a JSON object.
//
// The function assumes that the `Context` represents a JSON object. It parses the JSON object and returns a map
// where the keys are strings, and the values are `Context` elements representing the corresponding JSON values.
//
// If the `Context` does not represent a valid JSON object, the function will return an empty map.
//
// Parameters:
//   - ctx: The `Context` instance that holds the raw JSON string. The function checks if the context represents
//     a JSON object and processes it accordingly.
//
// Returns:
//   - map[string]Context: A map where the keys are strings (representing the keys in the JSON object), and
//     the values are `Context` instances representing the corresponding JSON values. If the context does not represent
//     a valid JSON object, a nil is returned.
//
// Example Usage:
//
//	ctx := Context{kind: JSON, unprocessed: "{\"key1\": \"value1\", \"key2\": 42}"}
//	result := ctx.Map()
//	// result.OpMap contains the parsed key-value pairs: {"key1": "value1", "key2": 42}
//
// Notes:
//   - The function calls `parseJSONElements` with the expected JSON object indicator ('{') to parse the JSON.
//   - If the `Context` is not a valid JSON object, it returns an empty map, which can be used to safely handle errors.
func (ctx Context) Map() map[string]Context {
	if ctx.kind != JSON {
		return nil
	}
	e := ctx.parseJSONElements('{', false)
	return e.OpMap
}

// Get searches for a specified path within a JSON structure and returns the corresponding result.
//
// This function allows you to search for a specific path in the JSON structure and retrieve the corresponding
// value as a `Context`. The path is represented as a string and can be used to navigate nested arrays or objects.
//
// The `path` parameter specifies the JSON path to search for, and the function will attempt to retrieve the value
// associated with that path. The result is returned as a `Context`, which contains information about the matched
// JSON value, including its type, string representation, numeric value, and index in the original JSON.
//
// Parameters:
//   - path: A string representing the path in the JSON structure to search for. The path may include array indices
//     and object keys separated by dots or brackets (e.g., "user.name", "items[0].price").
//
// Returns:
//   - Context: A `Context` instance containing the result of the search. The `Context` may represent various types of
//     JSON values (e.g., string, number, object, array). If no match is found, the `Context` will be empty.
//
// Example Usage:
//
//	ctx := Context{kind: JSON, unprocessed: "{\"user\": {\"name\": \"John\"}, \"items\": [1, 2, 3]}"}
//	result := ctx.Get("user.name")
//	// result.strings will contain "John", representing the value found at the "user.name" path.
//
// Notes:
//   - The function uses the `Get` function (presumably another function) to process the `unprocessed` JSON string
//     and search for the specified path.
//   - The function adjusts the indices of the results (if any) to account for the original position of the `Context`
//     in the JSON string.
func (ctx Context) Get(path string) Context {
	q := Get(ctx.unprocessed, path)
	if q.indexes != nil {
		for i := 0; i < len(q.indexes); i++ {
			q.indexes[i] += ctx.index
		}
	} else {
		q.index += ctx.index
	}
	return q
}

// parseJSONElements processes a JSON string (from the `Context`) and attempts to parse it as either a JSON array or a JSON object.
//
// The function examines the raw JSON string and determines whether it represents an array or an object by looking at
// the first character ('[' for arrays, '{' for objects). It then processes the content accordingly and returns the
// parsed results as a `queryContext`, which contains either an array or an object, depending on the type of the JSON structure.
//
// Parameters:
//   - vc: A byte representing the expected JSON structure type to parse ('[' for arrays, '{' for objects).
//   - valueSize: A boolean flag that indicates whether intermediary values should be stored as raw types (`true`)
//     or parsed into `Context` objects (`false`).
//
// Returns:
//   - queryContext: A `queryContext` struct containing the parsed elements. This can include:
//   - ArrayResult: A slice of `Context` elements for arrays.
//   - ArrayIns: A slice of `interface{}` elements for arrays when `valueSize` is true.
//   - OpMap: A map of string keys to `Context` values for objects when `valueSize` is false.
//   - OpIns: A map of string keys to `interface{}` values for objects when `valueSize` is true.
//   - valueN: The byte value indicating the start of the JSON array or object ('[' or '{').
//
// Function Process:
//
//  1. **Identifying JSON Structure**:
//     The function starts by checking the first non-whitespace character in the JSON string to determine if it's an object (`{`)
//     or an array (`[`). If the expected structure is detected, the function proceeds accordingly.
//
//  2. **Creating Appropriate Containers**:
//     Based on the type of JSON being parsed (array or object), the function initializes an empty slice or map
//     to store the parsed elements. The `OpMap` or `OpIns` is used for objects, while the `ArrayResult` or `ArrayIns`
//     is used for arrays. If `valueSize` is `true`, the values will be stored as raw types (`interface{}`), otherwise,
//     they will be stored as `Context` objects.
//
//  3. **Parsing JSON Elements**:
//     The function then loops through the JSON string, identifying and parsing individual elements. Each element could
//     be a string, number, boolean, `null`, array, or object. For each identified element, it is added to the appropriate
//     container (array or map) as determined by the type of JSON being processed.
//
//  4. **Handling Key-Value Pairs (for Objects)**:
//     If parsing an object (denoted by `{`), the function identifies key-value pairs and alternates between storing the
//     key (as a string) and its corresponding value (as a `Context` object or raw type) in the `OpMap` or `OpIns` container.
//
//  5. **Assigning Indices**:
//     After parsing the elements, the function assigns the correct index to each element in the `ArrayResult` based on
//     the `indexes` from the parent `Context`. If the number of elements in the array does not match the expected
//     number of indexes, the indices are reset to 0 for each element.
//
// Example Usage:
//
//	ctx := Context{kind: JSON, unprocessed: "[1, 2, 3]"}
//	result := ctx.parseJSONElements('[', false)
//	// result.ArrayResult contains the parsed `Context` elements for the array.
//
//	ctx = Context{kind: JSON, unprocessed: "{\"key\": \"value\"}"}
//	result = ctx.parseJSONElements('{', false)
//	// result.OpMap contains the parsed key-value pair for the object.
//
// Notes:
//   - The function handles various JSON value types, including numbers, strings, booleans, null, and nested arrays/objects.
//   - The function uses internal helper functions like `getNumeric`, `squash`, `lowerPrefix`, and `unescapeJSONEncoded`
//     to parse the raw JSON string into appropriate `Context` elements.
//   - The `valueSize` flag controls whether the elements are stored as raw types (`interface{}`) or as `Context` objects.
//   - If `valueSize` is `false`, the result will contain structured `Context` elements, which can be used for further processing or queries.
func (ctx Context) parseJSONElements(vc byte, valueSize bool) (result queryContext) {
	var json = ctx.unprocessed
	var i int
	var value Context
	var count int
	var key Context
	if vc == 0 {
		for ; i < len(json); i++ {
			if json[i] == '{' || json[i] == '[' {
				result.valueN = json[i]
				i++
				break
			}
			if json[i] > ' ' {
				goto end
			}
		}
	} else {
		for ; i < len(json); i++ {
			if json[i] == vc {
				i++
				break
			}
			if json[i] > ' ' {
				goto end
			}
		}
		result.valueN = vc
	}
	if result.valueN == '{' {
		if valueSize {
			result.OpIns = make(map[string]interface{})
		} else {
			result.OpMap = make(map[string]Context)
		}
	} else {
		if valueSize {
			result.ArrayIns = make([]interface{}, 0)
		} else {
			result.ArrayResult = make([]Context, 0)
		}
	}
	for ; i < len(json); i++ {
		if json[i] <= ' ' {
			continue
		}
		if json[i] == ']' || json[i] == '}' {
			break
		}
		switch json[i] {
		default:
			if (json[i] >= '0' && json[i] <= '9') || json[i] == '-' {
				value.kind = Number
				value.unprocessed, value.numeric = getNumeric(json[i:])
				value.strings = ""
			} else {
				continue
			}
		case '{', '[':
			value.kind = JSON
			value.unprocessed = squash(json[i:])
			value.strings, value.numeric = "", 0
		case 'n':
			value.kind = Null
			value.unprocessed = lowerPrefix(json[i:])
			value.strings, value.numeric = "", 0
		case 't':
			value.kind = True
			value.unprocessed = lowerPrefix(json[i:])
			value.strings, value.numeric = "", 0
		case 'f':
			value.kind = False
			value.unprocessed = lowerPrefix(json[i:])
			value.strings, value.numeric = "", 0
		case '"':
			value.kind = String
			value.unprocessed, value.strings = unescapeJSONEncoded(json[i:])
			value.numeric = 0
		}
		value.index = i + ctx.index

		i += len(value.unprocessed) - 1

		if result.valueN == '{' {
			if count%2 == 0 {
				key = value
			} else {
				if valueSize {
					if _, ok := result.OpIns[key.strings]; !ok {
						result.OpIns[key.strings] = value.Value()
					}
				} else {
					if _, ok := result.OpMap[key.strings]; !ok {
						result.OpMap[key.strings] = value
					}
				}
			}
			count++
		} else {
			if valueSize {
				result.ArrayIns = append(result.ArrayIns, value.Value())
			} else {
				result.ArrayResult = append(result.ArrayResult, value)
			}
		}
	}
end:
	if ctx.indexes != nil {
		if len(ctx.indexes) != len(result.ArrayResult) {
			for i := 0; i < len(result.ArrayResult); i++ {
				result.ArrayResult[i].index = 0
			}
		} else {
			for i := 0; i < len(result.ArrayResult); i++ {
				result.ArrayResult[i].index = ctx.indexes[i]
			}
		}
	}
	return
}

// Parse parses the json and returns a result.
//
// This function expects that the json is well-formed, and does not validate.
// Invalid json will not panic, but it may return back unexpected results.
// If you are consuming JSON from an unpredictable source then you may want to
// use the Valid function first.
func Parse(json string) Context {
	var value Context
	i := 0
	for ; i < len(json); i++ {
		if json[i] == '{' || json[i] == '[' {
			value.kind = JSON
			value.unprocessed = json[i:] // just take the entire raw
			break
		}
		if json[i] <= ' ' {
			continue
		}
		switch json[i] {
		case '+', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'i', 'I', 'N':
			value.kind = Number
			value.unprocessed, value.numeric = getNumeric(json[i:])
		case 'n':
			if i+1 < len(json) && json[i+1] != 'u' {
				// nan
				value.kind = Number
				value.unprocessed, value.numeric = getNumeric(json[i:])
			} else {
				// null
				value.kind = Null
				value.unprocessed = lowerPrefix(json[i:])
			}
		case 't':
			value.kind = True
			value.unprocessed = lowerPrefix(json[i:])
		case 'f':
			value.kind = False
			value.unprocessed = lowerPrefix(json[i:])
		case '"':
			value.kind = String
			value.unprocessed, value.strings = unescapeJSONEncoded(json[i:])
		default:
			return Context{}
		}
		break
	}
	if value.Exists() {
		value.index = i
	}
	return value
}

// ParseBytes parses the json and returns a result.
// If working with bytes, this method preferred over Parse(string(data))
func ParseBytes(json []byte) Context {
	return Parse(string(json))
}

// ForEachLine iterates through lines of JSON as specified by the JSON Lines
// format (http://jsonlines.org/).
// Each line is returned as a fj Result.
func ForEachLine(json string, iterator func(line Context) bool) {
	var res Context
	var i int
	for {
		i, res, _ = parseJSONAny(json, i, true)
		if !res.Exists() {
			break
		}
		if !iterator(res) {
			return
		}
	}
}

// Get searches json for the specified path.
// A path is in dot syntax, such as "name.last" or "age".
// When the value is found it's returned immediately.
//
// A path is a series of keys separated by a dot.
// A key may contain special wildcard characters '*' and '?'.
// To access an array value use the index as the key.
// To get the number of elements in an array or to access a child path, use
// the '#' character.
// The dot and wildcard character can be escaped with '\'.
//
//	{
//	  "name": {"first": "Tom", "last": "Anderson"},
//	  "age":37,
//	  "children": ["Sara","Alex","Jack"],
//	  "friends": [
//	    {"first": "James", "last": "Murphy"},
//	    {"first": "Roger", "last": "Craig"}
//	  ]
//	}
//	"name.last"          >> "Anderson"
//	"age"                >> 37
//	"children"           >> ["Sara","Alex","Jack"]
//	"children.#"         >> 3
//	"children.1"         >> "Alex"
//	"child*.2"           >> "Jack"
//	"children.0"         >> "Sara"
//	"friends.#.first"    >> ["James","Roger"]
//
// This function expects that the json is well-formed, and does not validate.
// Invalid json will not panic, but it may return back unexpected results.
// If you are consuming JSON from an unpredictable source then you may want to
// use the Valid function first.
func Get(json, path string) Context {
	if len(path) > 1 {
		if (path[0] == '@' && !DisableModifiers) || path[0] == '!' {
			// possible modifier
			var ok bool
			var cPath string
			var cJson string
			if path[0] == '@' && !DisableModifiers {
				cPath, cJson, ok = adjustModifier(json, path)
			} else if path[0] == '!' {
				cPath, cJson, ok = parseStaticValue(path)
			}
			if ok {
				path = cPath
				if len(path) > 0 && (path[0] == '|' || path[0] == '.') {
					res := Get(cJson, path[1:])
					res.index = 0
					res.indexes = nil
					return res
				}
				return Parse(cJson)
			}
		}
		if path[0] == '[' || path[0] == '{' {
			// using a sub-selector path
			kind := path[0]
			var ok bool
			var subs []subSelector
			subs, path, ok = analyzeSubSelectors(path)
			if ok {
				if len(path) == 0 || (path[0] == '|' || path[0] == '.') {
					var b []byte
					b = append(b, kind)
					var i int
					for _, sub := range subs {
						res := Get(json, sub.path)
						if res.Exists() {
							if i > 0 {
								b = append(b, ',')
							}
							if kind == '{' {
								if len(sub.name) > 0 {
									if sub.name[0] == '"' && Valid(sub.name) {
										b = append(b, sub.name...)
									} else {
										b = AppendJSON(b, sub.name)
									}
								} else {
									last := lastSegment(sub.path)
									if isValidName(last) {
										b = AppendJSON(b, last)
									} else {
										b = AppendJSON(b, "_")
									}
								}
								b = append(b, ':')
							}
							var raw string
							if len(res.unprocessed) == 0 {
								raw = res.String()
								if len(raw) == 0 {
									raw = "null"
								}
							} else {
								raw = res.unprocessed
							}
							b = append(b, raw...)
							i++
						}
					}
					b = append(b, kind+2)
					var res Context
					res.unprocessed = string(b)
					res.kind = JSON
					if len(path) > 0 {
						res = res.Get(path[1:])
					}
					res.index = 0
					return res
				}
			}
		}
	}
	var i int
	var c = &parser{json: json}
	if len(path) >= 2 && path[0] == '.' && path[1] == '.' {
		c.lines = true
		analyzeArray(c, 0, path[2:])
	} else {
		for ; i < len(c.json); i++ {
			if c.json[i] == '{' {
				i++
				parseJSONObject(c, i, path)
				break
			}
			if c.json[i] == '[' {
				i++
				analyzeArray(c, i, path)
				break
			}
		}
	}
	if c.piped {
		res := c.value.Get(c.pipe)
		res.index = 0
		return res
	}
	calcSubstring(json, c)
	return c.value
}

// GetBytes searches the provided JSON byte slice for the specified path and returns a `Context`
// representing the extracted data. This method is preferred over `Get(string(data), path)` when working
// with JSON data in byte slice format, as it directly operates on the byte slice, minimizing memory
// allocations and unnecessary copies.
//
// Parameters:
//   - `json`: A byte slice containing the JSON data to process.
//   - `path`: A string representing the path in the JSON data to extract.
//
// Returns:
//   - A `Context` struct containing the processed JSON data. The `Context` struct includes both
//     the raw unprocessed JSON string and the specific extracted string based on the given path.
//
// Notes:
//   - This function internally calls the `getBytes` function, which uses unsafe pointer operations
//     to minimize allocations and efficiently handle string slice headers.
//   - The function avoids unnecessary memory allocations by directly processing the byte slice and
//     utilizing memory safety features to manage substring extraction when the `strings` part is
//     a substring of the `unprocessed` part of the JSON data.
//
// Example:
//
//	jsonBytes := []byte(`{"key": "value", "nested": {"innerKey": "innerValue"}}`)
//	path := "nested.innerKey"
//	context := GetBytes(jsonBytes, path)
//	fmt.Println("Unprocessed:", context.unprocessed) // Output: `{"key": "value", "nested": {"innerKey": "innerValue"}}`
//	fmt.Println("Strings:", context.strings)         // Output: `"innerValue"`
func GetBytes(json []byte, path string) Context {
	return getBytes(json, path)
}

// Less compares two Context values (tokens) and returns true if the first token is considered less than the second one.
// It performs comparisons based on the type of the tokens and their respective values.
// The comparison order follows: Null < False < Number < String < True < JSON.
// This function also supports case-insensitive comparisons for String type tokens based on the caseSensitive parameter.
//
// Parameters:
//   - token: The other Context token to compare with the current one (t).
//   - caseSensitive: A boolean flag that indicates whether the comparison for String type tokens should be case-sensitive.
//   - If true, the comparison is case-sensitive (i.e., "a" < "b" but "A" < "b").
//   - If false, the comparison is case-insensitive (i.e., "a" == "A").
//
// Returns:
//   - true: If the current token (t) is considered less than the provided token.
//   - false: If the current token (t) is not considered less than the provided token.
//
// The function first compares the `kind` of both tokens, which represents their JSON types.
// If both tokens have the same kind, it proceeds to compare based on their specific types:
// - For String types, it compares the strings based on the case-sensitive flag.
// - For Number types, it compares the numeric values directly.
// - For other types, it compares the unprocessed JSON values as raw strings (this could be useful for types like Null, Boolean, etc.).
//
// Example usage:
//
//	context1 := Context{kind: String, strings: "apple"}
//	context2 := Context{kind: String, strings: "banana"}
//	result := context1.Less(context2, true) // This would return true because "apple" < "banana" and case-sensitive comparison is used.
func (t Context) Less(token Context, caseSensitive bool) bool {
	if t.kind < token.kind {
		return true
	}
	if t.kind > token.kind {
		return false
	}
	if t.kind == String {
		if caseSensitive {
			return t.strings < token.strings
		}
		return lessInsensitive(t.strings, token.strings)
	}
	if t.kind == Number {
		return t.numeric < token.numeric
	}
	return t.unprocessed < token.unprocessed
}

// GetMany searches json for the multiple paths.
// The return value is a Result array where the number of items
// will be equal to the number of input paths.
func GetMany(json string, path ...string) []Context {
	res := make([]Context, len(path))
	for i, path := range path {
		res[i] = Get(json, path)
	}
	return res
}

// GetManyBytes searches json for the multiple paths.
// The return value is a Result array where the number of items
// will be equal to the number of input paths.
func GetManyBytes(json []byte, path ...string) []Context {
	res := make([]Context, len(path))
	for i, path := range path {
		res[i] = GetBytes(json, path)
	}
	return res
}

// Valid returns true if the input is valid json.
//
//	if !bjson.Valid(json) {
//		return errors.New("invalid json")
//	}
//	value := bjson.Get(json, "name.last")
func Valid(json string) bool {
	_, ok := verifyJson(fromStr2Bytes(json), 0)
	return ok
}

// ValidBytes returns true if the input is valid json.
//
//	if !bjson.Valid(json) {
//		return errors.New("invalid json")
//	}
//	value := bjson.Get(json, "name.last")
//
// If working with bytes, this method preferred over ValidBytes(string(data))
func ValidBytes(json []byte) bool {
	_, ok := verifyJson(json, 0)
	return ok
}

func init() {
	modifiers = map[string]func(json, arg string) string{
		"pretty":  modPretty,
		"ugly":    modUgly,
		"reverse": modReverse,
		"this":    modThis,
		"flatten": modFlatten,
		"join":    modJoin,
		"valid":   modValid,
		"keys":    modKeys,
		"values":  modValues,
		"tostr":   modToStr,
		"fromstr": modFromStr,
		"group":   modGroup,
		"dig":     modDig,
	}
}

// AddModifier binds a custom modifier command to the bjson syntax.
// This operation is not thread safe and should be executed prior to
// using all other bjson function.
func AddModifier(name string, fn func(json, arg string) string) {
	modifiers[name] = fn
}

// ModifierExists returns true when the specified modifier exists.
func ModifierExists(name string, fn func(json, arg string) string) bool {
	_, ok := modifiers[name]
	return ok
}

// @pretty modifier makes the json look nice.
func modPretty(json, arg string) string {
	if len(arg) > 0 {
		opts := *unify4g.DefaultOptionsConfig
		Parse(arg).ForEach(func(key, value Context) bool {
			switch key.String() {
			case "sortKeys":
				opts.SortKeys = value.Bool()
			case "indent":
				opts.Indent = stripNonWhitespace(value.String())
			case "prefix":
				opts.Prefix = stripNonWhitespace(value.String())
			case "width":
				opts.Width = int(value.Int64())
			}
			return true
		})
		return fromBytes2Str(unify4g.PrettyOptions(fromStr2Bytes(json), &opts))
	}
	return fromBytes2Str(unify4g.Pretty(fromStr2Bytes(json)))
}

// @this returns the current element. Can be used to retrieve the root element.
func modThis(json, arg string) string {
	return json
}

// @ugly modifier removes all whitespace.
func modUgly(json, arg string) string {
	return fromBytes2Str(unify4g.Ugly(fromStr2Bytes(json)))
}

// @reverse reverses array elements or root object members.
func modReverse(json, arg string) string {
	res := Parse(json)
	if res.IsArray() {
		var values []Context
		res.ForEach(func(_, value Context) bool {
			values = append(values, value)
			return true
		})
		out := make([]byte, 0, len(json))
		out = append(out, '[')
		for i, j := len(values)-1, 0; i >= 0; i, j = i-1, j+1 {
			if j > 0 {
				out = append(out, ',')
			}
			out = append(out, values[i].unprocessed...)
		}
		out = append(out, ']')
		return fromBytes2Str(out)
	}
	if res.IsObject() {
		var keyValues []Context
		res.ForEach(func(key, value Context) bool {
			keyValues = append(keyValues, key, value)
			return true
		})
		out := make([]byte, 0, len(json))
		out = append(out, '{')
		for i, j := len(keyValues)-2, 0; i >= 0; i, j = i-2, j+1 {
			if j > 0 {
				out = append(out, ',')
			}
			out = append(out, keyValues[i+0].unprocessed...)
			out = append(out, ':')
			out = append(out, keyValues[i+1].unprocessed...)
		}
		out = append(out, '}')
		return fromBytes2Str(out)
	}
	return json
}

// @flatten an array with child arrays.
//
//	[1,[2],[3,4],[5,[6,7]]] -> [1,2,3,4,5,[6,7]]
//
// The {"deep":true} arg can be provide for deep flattening.
//
//	[1,[2],[3,4],[5,[6,7]]] -> [1,2,3,4,5,6,7]
//
// The original json is returned when the json is not an array.
func modFlatten(json, arg string) string {
	res := Parse(json)
	if !res.IsArray() {
		return json
	}
	var deep bool
	if arg != "" {
		Parse(arg).ForEach(func(key, value Context) bool {
			if key.String() == "deep" {
				deep = value.Bool()
			}
			return true
		})
	}
	var out []byte
	out = append(out, '[')
	var idx int
	res.ForEach(func(_, value Context) bool {
		var raw string
		if value.IsArray() {
			if deep {
				raw = removeOuterBraces(modFlatten(value.unprocessed, arg))
			} else {
				raw = removeOuterBraces(value.unprocessed)
			}
		} else {
			raw = value.unprocessed
		}
		raw = strings.TrimSpace(raw)
		if len(raw) > 0 {
			if idx > 0 {
				out = append(out, ',')
			}
			out = append(out, raw...)
			idx++
		}
		return true
	})
	out = append(out, ']')
	return fromBytes2Str(out)
}

// @keys extracts the keys from an object.
//
//	{"first":"Tom","last":"Smith"} -> ["first","last"]
func modKeys(json, arg string) string {
	v := Parse(json)
	if !v.Exists() {
		return "[]"
	}
	obj := v.IsObject()
	var out strings.Builder
	out.WriteByte('[')
	var i int
	v.ForEach(func(key, _ Context) bool {
		if i > 0 {
			out.WriteByte(',')
		}
		if obj {
			out.WriteString(key.unprocessed)
		} else {
			out.WriteString("null")
		}
		i++
		return true
	})
	out.WriteByte(']')
	return out.String()
}

// @values extracts the values from an object.
//
//	{"first":"Tom","last":"Smith"} -> ["Tom","Smith"]
func modValues(json, arg string) string {
	v := Parse(json)
	if !v.Exists() {
		return "[]"
	}
	if v.IsArray() {
		return json
	}
	var out strings.Builder
	out.WriteByte('[')
	var i int
	v.ForEach(func(_, value Context) bool {
		if i > 0 {
			out.WriteByte(',')
		}
		out.WriteString(value.unprocessed)
		i++
		return true
	})
	out.WriteByte(']')
	return out.String()
}

// @join multiple objects into a single object.
//
//	[{"first":"Tom"},{"last":"Smith"}] -> {"first","Tom","last":"Smith"}
//
// The arg can be "true" to specify that duplicate keys should be preserved.
//
//	[{"first":"Tom","age":37},{"age":41}] -> {"first","Tom","age":37,"age":41}
//
// Without preserved keys:
//
//	[{"first":"Tom","age":37},{"age":41}] -> {"first","Tom","age":41}
//
// The original json is returned when the json is not an object.
func modJoin(json, arg string) string {
	res := Parse(json)
	if !res.IsArray() {
		return json
	}
	var preserve bool
	if arg != "" {
		Parse(arg).ForEach(func(key, value Context) bool {
			if key.String() == "preserve" {
				preserve = value.Bool()
			}
			return true
		})
	}
	var out []byte
	out = append(out, '{')
	if preserve {
		// Preserve duplicate keys.
		var idx int
		res.ForEach(func(_, value Context) bool {
			if !value.IsObject() {
				return true
			}
			if idx > 0 {
				out = append(out, ',')
			}
			out = append(out, removeOuterBraces(value.unprocessed)...)
			idx++
			return true
		})
	} else {
		// Deduplicate keys and generate an object with stable ordering.
		var keys []Context
		keyVal := make(map[string]Context)
		res.ForEach(func(_, value Context) bool {
			if !value.IsObject() {
				return true
			}
			value.ForEach(func(key, value Context) bool {
				k := key.String()
				if _, ok := keyVal[k]; !ok {
					keys = append(keys, key)
				}
				keyVal[k] = value
				return true
			})
			return true
		})
		for i := 0; i < len(keys); i++ {
			if i > 0 {
				out = append(out, ',')
			}
			out = append(out, keys[i].unprocessed...)
			out = append(out, ':')
			out = append(out, keyVal[keys[i].String()].unprocessed...)
		}
	}
	out = append(out, '}')
	return fromBytes2Str(out)
}

// @valid ensures that the json is valid before moving on. An empty string is
// returned when the json is not valid, otherwise it returns the original json.
func modValid(json, arg string) string {
	if !Valid(json) {
		return ""
	}
	return json
}

// @fromstr converts a string to json
//
//	"{\"id\":1023,\"name\":\"alert\"}" -> {"id":1023,"name":"alert"}
func modFromStr(json, arg string) string {
	if !Valid(json) {
		return ""
	}
	return Parse(json).String()
}

// @tostr converts a string to json
//
//	{"id":1023,"name":"alert"} -> "{\"id\":1023,\"name\":\"alert\"}"
func modToStr(str, arg string) string {
	return string(AppendJSON(nil, str))
}

func modGroup(json, arg string) string {
	res := Parse(json)
	if !res.IsObject() {
		return ""
	}
	var all [][]byte
	res.ForEach(func(key, value Context) bool {
		if !value.IsArray() {
			return true
		}
		var idx int
		value.ForEach(func(_, value Context) bool {
			if idx == len(all) {
				all = append(all, []byte{})
			}
			all[idx] = append(all[idx], ("," + key.unprocessed + ":" + value.unprocessed)...)
			idx++
			return true
		})
		return true
	})
	var data []byte
	data = append(data, '[')
	for i, item := range all {
		if i > 0 {
			data = append(data, ',')
		}
		data = append(data, '{')
		data = append(data, item[1:]...)
		data = append(data, '}')
	}
	data = append(data, ']')
	return string(data)
}

// Paths returns the original bjson paths for a Result where the Result came
// from a simple query path that returns an array, like:
//
//	bjson.Get(json, "friends.#.first")
//
// The returned value will be in the form of a JSON array:
//
//	["friends.0.first","friends.1.first","friends.2.first"]
//
// The param 'json' must be the original JSON used when calling Get.
//
// Returns an empty string if the paths cannot be determined, which can happen
// when the Result came from a path that contained a multi-path, modifier,
// or a nested query.
func (t Context) Paths(json string) []string {
	if t.indexes == nil {
		return nil
	}
	paths := make([]string, 0, len(t.indexes))
	t.ForEach(func(_, value Context) bool {
		paths = append(paths, value.Path(json))
		return true
	})
	if len(paths) != len(t.indexes) {
		return nil
	}
	return paths
}

// Path returns the original bjson path for a Result where the Result came
// from a simple path that returns a single value, like:
//
//	bjson.Get(json, "friends.#(last=Murphy)")
//
// The returned value will be in the form of a JSON string:
//
//	"friends.0"
//
// The param 'json' must be the original JSON used when calling Get.
//
// Returns an empty string if the paths cannot be determined, which can happen
// when the Result came from a path that contained a multi-path, modifier,
// or a nested query.
func (t Context) Path(json string) string {
	var path []byte
	var comps []string // raw components
	i := t.index - 1
	if t.index+len(t.unprocessed) > len(json) {
		// JSON cannot safely contain Result.
		goto fail
	}
	if !strings.HasPrefix(json[t.index:], t.unprocessed) {
		// Result is not at the JSON index as expected.
		goto fail
	}
	for ; i >= 0; i-- {
		if json[i] <= ' ' {
			continue
		}
		if json[i] == ':' {
			// inside of object, get the key
			for ; i >= 0; i-- {
				if json[i] != '"' {
					continue
				}
				break
			}
			raw := reverseSquash(json[:i+1])
			i = i - len(raw)
			comps = append(comps, raw)
			// key gotten, now squash the rest
			raw = reverseSquash(json[:i+1])
			i = i - len(raw)
			i++ // increment the index for next loop step
		} else if json[i] == '{' {
			// Encountered an open object. The original result was probably an
			// object key.
			goto fail
		} else if json[i] == ',' || json[i] == '[' {
			// inside of an array, count the position
			var arrIdx int
			if json[i] == ',' {
				arrIdx++
				i--
			}
			for ; i >= 0; i-- {
				if json[i] == ':' {
					// Encountered an unexpected colon. The original result was
					// probably an object key.
					goto fail
				} else if json[i] == ',' {
					arrIdx++
				} else if json[i] == '[' {
					comps = append(comps, strconv.Itoa(arrIdx))
					break
				} else if json[i] == ']' || json[i] == '}' || json[i] == '"' {
					raw := reverseSquash(json[:i+1])
					i = i - len(raw) + 1
				}
			}
		}
	}
	if len(comps) == 0 {
		if DisableModifiers {
			goto fail
		}
		return "@this"
	}
	for i := len(comps) - 1; i >= 0; i-- {
		rawComplexity := Parse(comps[i])
		if !rawComplexity.Exists() {
			goto fail
		}
		comp := EscapeUnsafeChars(rawComplexity.String())
		path = append(path, '.')
		path = append(path, comp...)
	}
	if len(path) > 0 {
		path = path[1:]
	}
	return string(path)
fail:
	return ""
}

func parseRecursiveDescent(all []Context, parent Context, path string) []Context {
	if res := parent.Get(path); res.Exists() {
		all = append(all, res)
	}
	if parent.IsArray() || parent.IsObject() {
		parent.ForEach(func(_, val Context) bool {
			all = parseRecursiveDescent(all, val, path)
			return true
		})
	}
	return all
}

func modDig(json, arg string) string {
	all := parseRecursiveDescent(nil, Parse(json), arg)
	var out []byte
	out = append(out, '[')
	for i, res := range all {
		if i > 0 {
			out = append(out, ',')
		}
		out = append(out, res.unprocessed...)
	}
	out = append(out, ']')
	return string(out)
}
