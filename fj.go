package fj

import (
	"strconv"
	"strings"
	"time"

	"github.com/sivaosorg/unify4g"
)

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

// Foreach iterates through the values of a JSON object or array, applying the provided iterator function.
//
// If the `Context` represents a non-existent value (Null or invalid JSON), no iteration occurs.
// For JSON objects, the iterator receives both the key and value of each item.
// For JSON arrays, the iterator receives only the value of each item.
// If the `Context` is not an array or object, the iterator is called once with the whole value.
//
// Example Usage:
//
//	ctx.Foreach(func(key, value Context) bool {
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
func (ctx Context) Foreach(iterator func(key, value Context) bool) {
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
func (ctx Context) Less(token Context, caseSensitive bool) bool {
	if ctx.kind < token.kind {
		return true
	}
	if ctx.kind > token.kind {
		return false
	}
	if ctx.kind == String {
		if caseSensitive {
			return ctx.strings < token.strings
		}
		return lessInsensitive(ctx.strings, token.strings)
	}
	if ctx.kind == Number {
		return ctx.numeric < token.numeric
	}
	return ctx.unprocessed < token.unprocessed
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

// Parse parses a JSON string and returns a Context representing the parsed value.
//
// This function processes the input JSON string and attempts to determine the type of the value it represents.
// It handles objects, arrays, numbers, strings, booleans, and null values. The function does not validate whether
// the JSON is well-formed, and instead returns a Context object that represents the first valid JSON element found
// in the string. Invalid JSON may result in unexpected behavior, so for input from unpredictable sources, consider
// using the `Valid` function first.
//
// Parameters:
//   - `json`: A string containing the JSON data to be parsed. This function expects well-formed JSON and does not
//     perform comprehensive validation.
//
// Returns:
//   - A `Context` that represents the parsed JSON element. The `Context` contains details about the type, value,
//     and position of the JSON element, including raw and unprocessed string data.
//
// Notes:
//   - The function attempts to determine the type of the JSON element by inspecting the first character in the
//     string. It supports the following types: Object (`{`), Array (`[`), Number, String (`"`), Boolean (`true` / `false`),
//     and Null (`null`).
//   - The function sets the `unprocessed` field of the `Context` to the raw JSON string for further processing, and
//     sets the `kind` field to represent the type of the value (e.g., `String`, `Number`, `True`, `False`, `JSON`, `Null`).
//
// Example Usage:
//
//	json := "{\"name\": \"John\", \"age\": 30}"
//	ctx := Parse(json)
//	fmt.Println(ctx.kind) // Output: JSON (if the input starts with '{')
//
//	json := "12345"
//	ctx := Parse(json)
//	fmt.Println(ctx.kind) // Output: Number (if the input is a numeric value)
//
//	json := "\"Hello, World!\""
//	ctx := Parse(json)
//	fmt.Println(ctx.kind) // Output: String (if the input is a string)
//
// Returns:
//   - `Context`: The parsed result, which may represent an object, array, string, number, boolean, or null.
func Parse(json string) Context {
	var value Context
	i := 0
	for ; i < len(json); i++ {
		if json[i] == '{' || json[i] == '[' {
			value.kind = JSON
			value.unprocessed = json[i:]
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

// ParseBytes parses a JSON byte slice and returns a Context representing the parsed value.
//
// This function is a wrapper around the `Parse` function, designed specifically for handling JSON data
// in the form of a byte slice. It converts the byte slice into a string and then calls `Parse` to process
// the JSON data. If you're working with raw JSON data as bytes, using this method is preferred over
// manually converting the bytes to a string and passing it to `Parse`.
//
// Parameters:
//   - `json`: A byte slice containing the JSON data to be parsed.
//
// Returns:
//   - A `Context` representing the parsed JSON element, similar to the behavior of `Parse`. The `Context`
//     contains information about the type, value, and position of the JSON element, including the raw and
//     unprocessed string data.
//
// Example Usage:
//
//	json := []byte("{\"name\": \"Alice\", \"age\": 25}")
//	ctx := ParseBytes(json)
//	fmt.Println(ctx.kind) // Output: JSON (if the input is an object)
//
// Returns:
//   - `Context`: The parsed result, representing the parsed JSON element, such as an object, array, string,
//     number, boolean, or null.
func ParseBytes(json []byte) Context {
	return Parse(string(json))
}

// Get searches for a specified path within the provided JSON string and returns the corresponding value as a Context.
// The path is provided in dot notation, where each segment represents a key or index. The function supports wildcards
// (`*` and `?`), array indexing, and special characters like '#' to access array lengths or child paths. The function
// will return the first matching result it finds along the specified path.
//
// Path Syntax:
// - Dot notation: "name.last" or "age" for direct key lookups.
// - Wildcards: "*" matches any key, "?" matches a single character.
// - Array indexing: "children.0" accesses the first item in the "children" array.
// - The '#' character returns the number of elements in an array (e.g., "children.#" returns the array length).
// - The dot (`.`) and wildcard characters (`*`, `?`) can be escaped with a backslash (`\`).
//
// Example Usage:
//
//	json := `{
//	  "user": {"firstName": "Alice", "lastName": "Johnson"},
//	  "age": 29,
//	  "siblings": ["Ben", "Clara", "David"],
//	  "friends": [
//	    {"firstName": "Tom", "lastName": "Smith"},
//	    {"firstName": "Sophia", "lastName": "Davis"}
//	  ],
//	  "address": {"city": "New York", "zipCode": "10001"}
//	}`
//
//	// Examples of Get function with paths:
//	Get(json, "user.lastName")        // Returns: "Johnson"
//	Get(json, "age")                  // Returns: 29
//	Get(json, "siblings.#")           // Returns: 3 (number of siblings)
//	Get(json, "siblings.1")           // Returns: "Clara" (second sibling)
//	Get(json, "friends.#.firstName")  // Returns: ["Tom", "Sophia"]
//	Get(json, "address.zipCode")      // Returns: "10001"
//
// Details:
//   - The function does not validate JSON format but expects well-formed input.
//     Invalid JSON may result in unexpected behavior.
//   - Modifiers (e.g., `@` for adjusting paths) and special sub-selectors (e.g., `[` and `{`) are supported and processed
//     in the path before extracting values.
//   - For complex structures, the function analyzes the provided path, handles nested arrays or objects, and returns
//     a Context containing the value found at the specified location.
//
// Parameters:
//   - `json`: A string containing the JSON data to search through.
//   - `path`: A string representing the path to the desired value, using dot notation or other special characters as described.
//
// Returns:
//   - `Context`: A Context object containing the value found at the specified path, including information such as the
//     type (`kind`), the raw JSON string (`unprocessed`), and the parsed value if available (e.g., `strings` for strings).
//
// Notes:
//   - If the path is not found, the returned Context will reflect this with an empty or null value.
func Get(json, path string) Context {
	if len(path) > 1 {
		if (path[0] == '@' && !DisableModifiers) || path[0] == '!' {
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
			kind := path[0] // using a sub-selector path
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
										b = appendJSON(b, sub.name)
									}
								} else {
									last := lastSegment(sub.path)
									if isValidName(last) {
										b = appendJSON(b, last)
									} else {
										b = appendJSON(b, "_")
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

// GetMul searches json for multiple paths.
// The return value is a slice of `Context` objects, where the number of items
// will be equal to the number of input paths. Each `Context` represents the value
// extracted for the corresponding path.
//
// Parameters:
//   - `json`: A string containing the JSON data to search through.
//   - `path`: A variadic list of paths to search for within the JSON data.
//
// Returns:
//   - A slice of `Context` objects, one for each path provided in the `path` parameter.
//
// Notes:
//   - The function will return a `Context` for each path, and the order of the `Context`
//     objects in the result will match the order of the paths provided.
//
// Example:
//
//	json := `{
//	  "user": {"firstName": "Alice", "lastName": "Johnson"},
//	  "age": 29,
//	  "siblings": ["Ben", "Clara", "David"],
//	  "friends": [
//	    {"firstName": "Tom", "lastName": "Smith"},
//	    {"firstName": "Sophia", "lastName": "Davis"}
//	  ]
//	}`
//	paths := []string{"user.lastName", "age", "siblings.#", "friends.#.firstName"}
//	results := GetMul(json, paths...)
//	// The result will contain Contexts for each path: ["Johnson", 29, 3, ["Tom", "Sophia"]]
func GetMul(json string, path ...string) []Context {
	ctx := make([]Context, len(path))
	for i, path := range path {
		ctx[i] = Get(json, path)
	}
	return ctx
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

// GetMulBytes searches json for multiple paths in the provided JSON byte slice.
// The return value is a slice of `Context` objects, where the number of items
// will be equal to the number of input paths. Each `Context` represents the value
// extracted for the corresponding path. This method operates directly on the byte slice,
// which is preferred when working with JSON data in byte format to minimize memory allocations.
//
// Parameters:
//   - `json`: A byte slice containing the JSON data to search through.
//   - `path`: A variadic list of paths to search for within the JSON data.
//
// Returns:
//   - A slice of `Context` objects, one for each path provided in the `path` parameter.
//
// Notes:
//   - The function will return a `Context` for each path, and the order of the `Context`
//     objects in the result will match the order of the paths provided.
//
// Example:
//
//	jsonBytes := []byte(`{"user": {"firstName": "Alice", "lastName": "Johnson"}, "age": 29}`)
//	paths := []string{"user.lastName", "age"}
//	results := GetMulBytes(jsonBytes, paths...)
//	// The result will contain Contexts for each path: ["Johnson", 29]
func GetMulBytes(json []byte, path ...string) []Context {
	ctx := make([]Context, len(path))
	for i, path := range path {
		ctx[i] = GetBytes(json, path)
	}
	return ctx
}

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

// ForeachLine iterates through each line of JSON data in the JSON Lines format (http://jsonlines.org/),
// and applies a provided iterator function to each line. This is useful for processing large JSON data
// sets where each line is a separate JSON object, allowing for efficient parsing and handling of each object.
//
// Parameters:
//   - `json`: A string containing JSON Lines formatted data, where each line is a separate JSON object.
//   - `iterator`: A callback function that is called for each line. It receives a `Context` representing
//     the parsed JSON object for the current line. The iterator function should return `true` to continue
//     processing the next line, or `false` to stop the iteration.
//
// Example Usage:
//
//	json := `{"name": "Alice"}\n{"name": "Bob"}`
//	iterator := func(line Context) bool {
//	    fmt.Println(line)
//	    return true
//	}
//	ForeachLine(json, iterator)
//	// Output:
//	// {"name": "Alice"}
//	// {"name": "Bob"}
//
// Notes:
//   - This function assumes the input `json` is formatted as JSON Lines, where each line is a valid JSON object.
//   - The function stops processing as soon as the `iterator` function returns `false` for a line.
//   - The function handles each line independently, meaning it processes one JSON object at a time and provides
//     it to the iterator, which can be used to process or filter lines.
//
// Returns:
//   - This function does not return a value. It processes the JSON data line-by-line and applies the iterator to each.
func ForeachLine(json string, iterator func(line Context) bool) {
	var ctx Context
	var i int
	for {
		i, ctx, _ = parseJSONAny(json, i, true)
		if !ctx.Exists() {
			break
		}
		if !iterator(ctx) {
			return
		}
	}
}

// Valid returns true if the input is valid json.
//
//	if !fj.Valid(json) {
//		return errors.New("invalid json")
//	}
//	value := fj.Get(json, "name.last")
func Valid(json string) bool {
	_, ok := verifyJson(fromStr2Bytes(json), 0)
	return ok
}

// ValidBytes returns true if the input is valid json.
//
//	if !fj.Valid(json) {
//		return errors.New("invalid json")
//	}
//	value := fj.Get(json, "name.last")
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

// AddModifier binds a custom modifier command to the fj syntax.
// This operation is not thread safe and should be executed prior to
// using all other fj function.
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
		Parse(arg).Foreach(func(key, value Context) bool {
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
		res.Foreach(func(_, value Context) bool {
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
		res.Foreach(func(key, value Context) bool {
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
		Parse(arg).Foreach(func(key, value Context) bool {
			if key.String() == "deep" {
				deep = value.Bool()
			}
			return true
		})
	}
	var out []byte
	out = append(out, '[')
	var idx int
	res.Foreach(func(_, value Context) bool {
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
	v.Foreach(func(key, _ Context) bool {
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
	v.Foreach(func(_, value Context) bool {
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
		Parse(arg).Foreach(func(key, value Context) bool {
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
		res.Foreach(func(_, value Context) bool {
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
		res.Foreach(func(_, value Context) bool {
			if !value.IsObject() {
				return true
			}
			value.Foreach(func(key, value Context) bool {
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
	return string(appendJSON(nil, str))
}

func modGroup(json, arg string) string {
	res := Parse(json)
	if !res.IsObject() {
		return ""
	}
	var all [][]byte
	res.Foreach(func(key, value Context) bool {
		if !value.IsArray() {
			return true
		}
		var idx int
		value.Foreach(func(_, value Context) bool {
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

// Paths returns the original fj paths for a Result where the Result came
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
func (ctx Context) Paths(json string) []string {
	if ctx.indexes == nil {
		return nil
	}
	paths := make([]string, 0, len(ctx.indexes))
	ctx.Foreach(func(_, value Context) bool {
		paths = append(paths, value.Path(json))
		return true
	})
	if len(paths) != len(ctx.indexes) {
		return nil
	}
	return paths
}

// Path returns the original fj path for a Result where the Result came
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
func (ctx Context) Path(json string) string {
	var path []byte
	var comps []string // raw components
	i := ctx.index - 1
	if ctx.index+len(ctx.unprocessed) > len(json) {
		// JSON cannot safely contain Result.
		goto fail
	}
	if !strings.HasPrefix(json[ctx.index:], ctx.unprocessed) {
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
		parent.Foreach(func(_, val Context) bool {
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
