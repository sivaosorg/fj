package fj

import (
	"strconv"
	"strings"
	"time"
)

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
//   - transformers (e.g., `@` for adjusting paths) and special sub-selectors (e.g., `[` and `{`) are supported and processed
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
		if (path[0] == '@' && !DisableTransformers) || path[0] == '!' {
			var ok bool
			var cPath string
			var cJson string
			if path[0] == '@' && !DisableTransformers {
				cPath, cJson, ok = adjustTransformer(json, path)
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
									if sub.name[0] == '"' && IsValidJSON(sub.name) {
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

// Foreach iterates through each line of JSON data in the JSON Lines format (http://jsonlines.org/),
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
//	Foreach(json, iterator)
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
func Foreach(json string, iterator func(line Context) bool) {
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

// IsValidJSON checks whether the provided string contains valid JSON data.
// It attempts to parse the JSON and returns a boolean indicating if the JSON is well-formed.
//
// Parameters:
//   - `json`: A string representing the JSON data that needs to be validated.
//
// Returns:
//   - A boolean value (`true` or `false`):
//   - `true`: The provided JSON string is valid and well-formed.
//   - `false`: The provided JSON string is invalid or malformed.
//
// Notes:
//   - This function utilizes the `fromStr2Bytes` function to efficiently convert the input string
//     into a byte slice without allocating new memory. It then passes the byte slice to the
//     `verifyJSON` function to check if the string conforms to valid JSON syntax.
//   - If the input JSON is invalid, the function will return `false`, indicating that the JSON
//     cannot be parsed or is improperly structured.
//   - The function does not perform deep validation of the content of the JSON, but merely
//     checks if the string is syntactically correct according to JSON rules.
//
// Example Usage:
//
//	json := `{"name": {"first": "Alice", "last": "Johnson"}, "age": 30}`
//	if !IsValidJSON(json) {
//	    fmt.Println("Invalid JSON")
//	} else {
//	    fmt.Println("IsValidJSON JSON")
//	}
//
//	// Output: "IsValidJSON JSON"
func IsValidJSON(json string) bool {
	_, ok := verifyJSON(fromStr2Bytes(json), 0)
	return ok
}

// IsValidJSONBytes checks whether the provided byte slice contains valid JSON data.
// It attempts to parse the JSON and returns a boolean indicating if the JSON is well-formed.
//
// Parameters:
//   - `json`: A byte slice (`[]byte`) representing the JSON data that needs to be validated.
//
// Returns:
//   - A boolean value (`true` or `false`):
//   - `true`: The provided JSON byte slice is valid and well-formed.
//   - `false`: The provided JSON byte slice is invalid or malformed.
//
// Notes:
//   - This function works directly with a byte slice (`[]byte`) rather than a string, making it more efficient
//     when dealing with raw byte data that represents JSON. It avoids the need to convert between strings and
//     byte slices, which can improve performance and memory usage when working with large or binary JSON data.
//   - The function utilizes the `verifyJSON` function to check if the byte slice conforms to valid JSON syntax.
//   - If the input byte slice represents invalid JSON, the function will return `false`, indicating that the JSON
//     cannot be parsed or is improperly structured.
//   - The function does not perform deep validation of the content of the JSON, but only checks whether the byte slice
//     adheres to the syntax rules defined for JSON data structures.
//
// Example Usage:
//
//	jsonBytes := []byte(`{"name": {"first": "Alice", "last": "Johnson"}, "age": 30}`)
//	if !IsValidJSONBytes(jsonBytes) {
//	    fmt.Println("Invalid JSON")
//	} else {
//	    fmt.Println("Valid JSON")
//	}
//
//	// Output: "Valid JSON"
func IsValidJSONBytes(json []byte) bool {
	_, ok := verifyJSON(json, 0)
	return ok
}

// AddTransformer binds a custom transformer function to the fj syntax.
//
// This function allows users to register custom transformer functions that can be applied
// to JSON data in the fj query language. A transformer is a transformation function that
// takes two string arguments — the JSON data and an argument (such as a key or value) —
// and returns a modified version of the JSON data. The registered transformer can then
// be used in queries to modify the JSON data dynamically.
//
// Parameters:
//   - `name`: A string representing the name of the transformer. This name will be used
//     in the fj query language to reference the transformer.
//   - `fn`: A function that takes two string arguments: `json` (the JSON data to be modified),
//     and `arg` (the argument that the transformer will use to transform the JSON). The function
//     should return a modified string (the transformed JSON data).
//
// Example Usage:
//
//	// Define a custom transformer to uppercase all values in a JSON array.
//	uppercaseTransformer := func(json, arg string) string {
//	  return strings.ToUpper(json)  // Modify the JSON data (this is just a simple example).
//	}
//
//	// Add the custom transformer to the fj system with the name "uppercase".
//	fj.AddTransformer("uppercase", uppercaseTransformer)
//
//	// Now you can use the "uppercase" transformer in a query:
//	json := `{
//	  "store": {
//	    "book": [
//	      { "category": "fiction", "author": "J.K. Rowling", "title": "Harry Potter" },
//	      { "category": "science", "author": "Stephen Hawking", "title": "A Brief History of Time" }
//	    ],
//	    "music": [
//	      { "artist": "The Beatles", "album": "Abbey Road" },
//	      { "artist": "Pink Floyd", "album": "The Wall" }
//	    ]
//	  }
//	}`
//	result := fj.Get(json, "store.music.1.album|@uppercase").String()  // Applies the uppercase transformer to each value in the array.
//	// result will contain: THE WALL
//
// Notes:
//   - This function is not thread-safe, so it should be called once, typically during
//     the initialization phase, before performing any queries that rely on custom transformers.
//   - Once registered, the transformer can be used in fj queries to transform the JSON data
//     according to the logic defined in the `fn` function.
func AddTransformer(name string, fn func(json, arg string) string) {
	jsonTransformers[name] = fn
}

// IsTransformerRegistered checks whether a specified transformer has been registered in the fj system.
//
// This function allows users to verify if a transformer with a given name has already
// been added to the `fj` query system. transformers are custom functions that transform
// JSON data in queries. This utility is useful to prevent duplicate registrations
// or to confirm the availability of a specific transformer before using it.
//
// Parameters:
//   - `name`: A string representing the name of the transformer to check for existence.
//
// Returns:
//   - `bool`: Returns `true` if a transformer with the given name exists, otherwise returns `false`.
//
// Example Usage:
//
//	// Check if a custom transformer named "uppercase" has already been registered.
//	if fj.IsTransformerRegistered("uppercase") {
//	  fmt.Println("The 'uppercase' transformer is available.")
//	} else {
//	  fmt.Println("The 'uppercase' transformer has not been registered.")
//	}
//
// Notes:
//   - This function does not modify the `transformers` map; it only queries it to check
//     for the existence of the specified transformer.
//   - It is thread-safe when used only to query the existence of a transformer.
func IsTransformerRegistered(name string) bool {
	if isEmpty(name) {
		return false
	}
	if len(jsonTransformers) == 0 {
		return false
	}
	_, ok := jsonTransformers[name]
	return ok
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
	case True:
		return "true"
	case False:
		return "false"
	case String:
		return ctx.strings
	case JSON:
		return ctx.unprocessed
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

// Path returns the original fj path for a Result where the Result came
// from a simple query path that returns a single value. For example, if the
// `Get` function was called with a query path like:
//
//	fj.Get(json, "employees.#(first=Admin)")
//
// This function will return the original path that corresponds to the single
// value in the result, formatted as a JSON path.
//
// The returned value will be in the form of a JSON string:
//
//	"employees.0"
//
// The param 'json' must be the original JSON used when calling Get.
//
// Returns:
//   - A string representing the original path for the single value in the result.
//   - If the paths cannot be determined (e.g., due to the result being from
//     a multi-path, transformer, or a nested query), an empty string will be returned.
//
// Notes:
//   - The `Path` function operates by tracing the position of the result within
//     the original JSON string and reconstructing the query path based on this position.
//   - The function checks the surrounding JSON context (such as whether the result
//     is within an array or object) and extracts the relevant path information.
//   - The path components are identified by traversing the string from the result's index
//     and extracting the array or object keys that lead to the specific value.
//
// Example Usage:
//
//	json := `{
//	  "employees": [
//	    {"id": 1, "name": {"first": "John", "last": "Doe"}, "department": "HR"},
//	    {"id": 2, "name": {"first": "Jane", "last": "Smith"}, "department": "Engineering"},
//	    {"id": 3, "name": {"first": "Admin", "last": "Land"}, "department": "Marketing"},
//	    {"id": 4, "name": {"first": "Emily", "last": "Jones"}, "department": "Engineering"}
//	  ],
//	  "companies": [
//	    {"name": "TechCorp", "employees": [1, 2]},
//	    {"name": "BizGroup", "employees": [3, 4]}
//	  ]
//	}`
//
//	// Get the employee's last name who works in the Engineering department
//	result := fj.Get(json, "employees.#(department=Engineering).name.last")
//	path := result.Path(json)
//
//	// Output: "employees.1.name.last"
//
//	// Explanation:
//	// The `Path` function returns the path to the "last" name of the second
//	// employee in the "employees" array who works in the "Engineering" department.
//	// The path "employees.1.name.last" corresponds to the "Jane Smith" employee,
//	// and the query specifically looks at the "last" name of that employee.
func (ctx Context) Path(json string) string {
	var path []byte
	var components []string
	i := ctx.index - 1
	// Ensure the index is within bounds of the original JSON
	if ctx.index+len(ctx.unprocessed) > len(json) {
		// JSON cannot safely contain the Result.
		goto fail
	}
	// Ensure that the unprocessed part matches the expected JSON structure
	if !strings.HasPrefix(json[ctx.index:], ctx.unprocessed) {
		// Result is not at the expected index in the JSON.
		goto fail
	}
	// Traverse the JSON from the result's index to extract the path
	for ; i >= 0; i-- {
		if json[i] <= ' ' {
			continue
		}
		if json[i] == ':' {
			for ; i >= 0; i-- {
				if json[i] != '"' {
					continue
				}
				break
			}
			raw := reverseSquash(json[:i+1])
			i = i - len(raw)
			components = append(components, raw)
			// Key obtained, now process the next component
			raw = reverseSquash(json[:i+1])
			i = i - len(raw)
			i++ // Move index for next loop step
		} else if json[i] == '{' {
			// Encountered an open object, this is likely not a valid result
			goto fail
		} else if json[i] == ',' || json[i] == '[' {
			// Inside an array, count the position of the element
			var arrayIdx int
			if json[i] == ',' {
				arrayIdx++
				i--
			}
			for ; i >= 0; i-- {
				if json[i] == ':' {
					// Unexpected colon indicates an object key
					goto fail
				} else if json[i] == ',' {
					arrayIdx++
				} else if json[i] == '[' {
					components = append(components, strconv.Itoa(arrayIdx))
					break
				} else if json[i] == ']' || json[i] == '}' || json[i] == '"' {
					raw := reverseSquash(json[:i+1])
					i = i - len(raw) + 1
				}
			}
		}
	}
	// If no components are found, return a default path for "this"
	if len(components) == 0 {
		if DisableTransformers {
			goto fail
		}
		return "@this"
	}
	// Build the final path by appending each component
	for i := len(components) - 1; i >= 0; i-- {
		rawComplexity := Parse(components[i])
		if !rawComplexity.Exists() {
			goto fail
		}
		comp := escapeUnsafeChars(rawComplexity.String())
		path = append(path, '.')
		path = append(path, comp...)
	}
	// Remove the leading dot and return the final path
	if len(path) > 0 {
		path = path[1:]
	}
	return string(path)
fail:
	// Return an empty string if the path could not be determined
	return ""
}

// Paths returns the original fj paths for a Result where the Result came
// from a simple query path that returns an array. For example, if the
// `Get` function was called with a query path like:
//
//	fj.Get(json, "friends.#.first")
//
// This function will return the paths for each element in the resulting array,
// formatted as a JSON array. The returned paths are the original query paths
// for each item in the array, reflecting the specific positions of the elements
// in the original JSON structure.
//
// The returned value will be in the form of a JSON array, such as:
//
//	["friends.0.first", "friends.1.first", "friends.2.first"]
//
// Parameters:
//   - `json`: A string representing the original JSON used in the query.
//     This is required for resolving the specific paths corresponding to
//     each element in the resulting array.
//
// Returns:
//   - A slice of strings (`[]string`), each containing the original path for
//     an element in the result array.
//   - If the result was a simple query that returns an array, each string
//     will be a path to an individual element in the array.
//   - If the paths cannot be determined (e.g., due to the result being
//     from a multi-path, transformer, or a nested query), an empty slice will
//     be returned.
//
// Notes:
//   - The `Paths` function relies on the `indexes` field in the `Context`
//     object. If the `indexes` field is `nil`, the function will return `nil`.
//   - The function iterates over each element in the result (which is expected
//     to be an array) and appends the corresponding path to the `paths` slice.
//   - If the paths cannot be determined (e.g., due to the result coming from
//     a multi-path or more complex query), an empty slice will be returned.
//   - This function is useful for extracting the specific query paths for
//     elements within a larger result array, providing a way to inspect or
//     manipulate the paths of individual items.
//
// Example Usage:
//
//	json := `{
//	  "friends": [
//	    {"first": "Tom", "last": "Smith"},
//	    {"first": "Sophia", "last": "Davis"},
//	    {"first": "James", "last": "Miller"}
//	  ]
//	}`
//
//	result := fj.Get(json, "friends.#.first")
//	paths := result.Paths(json)
//
//	// Output: ["friends.0.first", "friends.1.first", "friends.2.first"]
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

func init() {
	jsonTransformers = map[string]func(json, arg string) string{
		"trim":       transformTrim,
		"this":       transformDefault,
		"valid":      transformJSONValidity,
		"pretty":     transformPretty,
		"minify":     transformMinify,
		"reverse":    transformReverse,
		"flatten":    transformFlatten,
		"join":       transformJoin,
		"keys":       transformKeys,
		"values":     transformValues,
		"string":     transformToString,
		"json":       transformToJSON,
		"group":      transformGroup,
		"search":     transformSearch,
		"uppercase":  transformUppercase,
		"lowercase":  transformLowercase,
		"flip":       transformFlip,
		"snakeCase":  transformSnakeCase,
		"camelCase":  transformCamelCase,
		"kebabCase":  transformKebabCase,
		"replace":    transformReplace,
		"replaceAll": transformReplaceAll,
		"hex":        transformToHex,
		"bin":        transformToBinary,
		"insertAt":   transformInsertAt,
		"wc":         transformCountWords,
		"padLeft":    transformPadLeft,
		"padRight":   transformPadRight,
	}
}
