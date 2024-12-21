package fj

import (
	"strconv"
	"unicode/utf16"
	"unicode/utf8"
	"unsafe"
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

// getNumeric extracts the numeric portion of a JSON-encoded string and converts it to a float.
//
// This function processes a JSON-encoded string to extract a numeric value
// up to a specified delimiter or end of the string. It handles delimiters such as
// whitespace, commas, closing brackets (`]`), and closing braces (`}`). Additionally,
// it assumes that characters such as '+' and '-' may be part of the numeric value.
//
// Parameters:
//   - `json`: A string containing the JSON-encoded numeric data.
//
// Returns:
//   - `raw`: A substring of the input string containing the numeric portion.
//   - `num`: A float64 representation of the extracted numeric value. If the conversion fails, `num` will be 0.
//
// Notes:
//   - The function does not validate the structure of the input JSON beyond extracting the numeric portion.
//   - This function assumes that the input string contains valid JSON data.
//
// Example:
//
//	json := "123.45, other data"
//	raw, num := getNumeric(json) // raw: "123.45", num: 123.45
//
// Note:
//
//   - If the input does not contain a valid numeric value, `num` will be 0 and `raw` will contain
//     the unmodified input string. The function will attempt to extract and convert the first numeric
//     portion it encounters, stopping when a delimiter is found.
//
//   - If no delimiters are encountered, the entire input string will be processed as a numeric value.
func getNumeric(json string) (raw string, num float64) {
	for i := 1; i < len(json); i++ {
		// check for characters that signify the end of a numeric value.
		if json[i] <= '-' {
			// break if the character is a whitespace or comma.
			if json[i] <= ' ' || json[i] == ',' {
				raw = json[:i]
				num, _ = strconv.ParseFloat(raw, 64) // convert the numeric substring to a float.
				return
			}
			// if the character is '+' or '-', assume it could be part of the number.
		} else if json[i] == ']' || json[i] == '}' {
			// break on closing brackets or braces (']' or '}')
			raw = json[:i]
			num, _ = strconv.ParseFloat(raw, 64)
			return
		}
	}
	// if no delimiters are encountered, process the entire string.
	raw = json
	num, _ = strconv.ParseFloat(raw, 64) // convert the entire string to a float.
	return
}

// isSafeKeyChar checks if a given byte is a valid character for a safe path key.
//
// This function is used to determine if a character can be part of a safe key in a path,
// where safe characters are typically those that are printable and non-special. It checks if the
// byte is a valid alphanumeric character or one of a few other acceptable symbols, including
// underscore ('_'), hyphen ('-'), and colon (':').
//
// Parameters:
//   - `c`: The byte (character) to be checked.
//
// Returns:
//   - `true`: If the byte is considered a safe character for a path key.
//   - `false`: If the byte is not considered safe for a path key.
//
// Safe characters for a path key include:
//   - Whitespace characters (ASCII values <= 32)
//   - Printable characters from ASCII ('!' to '~')
//   - Letters (uppercase and lowercase), numbers, underscore ('_'), hyphen ('-'), and colon (':')
//
// Example:
//
//	isSafeKeyChar('a') // true
//	isSafeKeyChar('$') // false
func isSafeKeyChar(c byte) bool {
	return c <= ' ' || c > '~' || c == '_' || c == '-' || c == ':' ||
		(c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9')
}

// reverseSquash processes a JSON-like string in reverse, extracting the portion of the string
// starting from the last significant character (ignoring nested objects and arrays).
// It is designed to handle strings that end with a closing bracket (`]`), closing brace (`}`),
// closing parenthesis (`)`), or a double-quote (`"`), and it squashes the value by
// ignoring any nested structures within arrays or objects.
//
// Parameters:
//   - `json`: A string representing the JSON data to be processed.
//
// Returns:
//   - A substring starting from the last non-nested character, effectively squashing the value
//     by ignoring nested objects and arrays. If no nested structures are present, it returns the
//     portion of the string starting from the last significant character.
//
// Notes:
//   - The function assumes that the last character in the string is a valid closing character
//     (`]`, `}`, `)`, or `"`).
//   - It handles escaping of quotes within the string by checking for escape sequences (e.g., `\"`).
//   - It counts the depth of nested objects and arrays using a `depth` variable to ensure that nested
//     structures are ignored while processing.
//   - The function works in reverse from the end of the string, looking for the outermost value or
//     structure and returning the corresponding substring.
//
// Example:
//
//		json := "{\"key\": \"value\", \"nested\": {\"innerKey\": \"innerValue\"}}"
//		result := reverseSquash(json) // result: "{\"key\": \"value\", \"nested\": {\"innerKey\": \"innerValue\"}}"
//		json1 := `{"key": "value", "nested": {"innerKey": "innerValue"}}`
//
//	 result1 := reverseSquash(json1)
//	 fmt.Println("Result 1:", result1) // Expected output: {"key": "value", "nested": {"innerKey": "innerValue"}}
//
// // Example 2: A JSON array with nested objects
//
//	json2 := `["item1", {"key": "value"}, ["nestedItem1", "nestedItem2"]]`
//	result2 := reverseSquash(json2)
//	fmt.Println("Result 2:", result2) // Expected output: ["item1", {"key": "value"}, ["nestedItem1", "nestedItem2"]]
//
// // Example 3: A JSON string with nested arrays
//
//	json3 := `{"array": [1, 2, [3, 4], 5]}`
//	result3 := reverseSquash(json3)
//	fmt.Println("Result 3:", result3) // Expected output: {"array": [1, 2, [3, 4], 5]}
//
// // Example 4: A simple string value (no nested structures)
//
//	json4 := `"simpleString"`
//	result4 := reverseSquash(json4)
//	fmt.Println("Result 4:", result4) // Expected output: "simpleString"
//
// // Example 5: A deeply nested structure with complex objects and arrays
//
//	json5 := `{"outer": {"inner": {"key": "value"}}}`
//	result5 := reverseSquash(json5)
//	fmt.Println("Result 5:", result5) // Expected output: {"outer": {"inner": {"key": "value"}}}
func reverseSquash(json string) string {
	i := len(json) - 1
	var depth int
	// if the last character is not a quote, assume it's part of a value and increase depth
	if json[i] != '"' {
		depth++
	}
	// if the last character is a closing bracket, brace, or parenthesis, adjust the index to skip it
	if json[i] == '}' || json[i] == ']' || json[i] == ')' {
		i--
	}
	// loop backwards through the string
	for ; i >= 0; i-- {
		switch json[i] {
		case '"':
			// handle strings enclosed in double quotes
			i-- // skip the opening quote
			for ; i >= 0; i-- {
				if json[i] == '"' {
					// check for escape sequences (e.g., \")
					esc := 0
					for i > 0 && json[i-1] == '\\' {
						i-- // move back over escape characters
						esc++
					}
					// if the quote is escaped, continue
					if esc%2 == 1 {
						continue
					}
					// if the quote is not escaped, break out of the loop
					i += esc
					break
				}
			}
			// if the depth is 0, we've found the outermost value
			if depth == 0 {
				if i < 0 {
					i = 0
				}
				return json[i:] // return the substring starting from the outermost value
			}
		case '}', ']', ')': // increase depth when encountering closing brackets, braces, or parentheses
			depth++
		case '{', '[', '(': // decrease depth when encountering opening brackets, braces, or parentheses
			depth--
			// if depth reaches 0, we've found the outermost value
			if depth == 0 {
				return json[i:] // return the substring starting from the outermost value
			}
		}
	}
	return json
}

// getBytes efficiently processes a JSON byte slice and a path string to produce a `Context`.
// This function minimizes memory allocations and copies, leveraging unsafe operations
// to handle large JSON strings and slice conversions.
//
// Parameters:
//   - `json`: A byte slice containing the JSON data to process.
//   - `path`: A string representing the path to extract data from the JSON.
//
// Returns:
//   - A `Context` struct containing processed and unprocessed strings representing
//     the result of applying the path to the JSON data.
//
// Notes:
//   - The function uses unsafe pointer operations to avoid unnecessary allocations and copies.
//   - It extracts string and byte slice headers and ensures memory safety by copying headers
//     to strings when needed.
//   - The function checks whether the substring (`strings`) is part of the raw string (`unprocessed`)
//     and handles memory overlap efficiently.
//
// Example:
//
//	jsonBytes := []byte(`{"key": "value", "nested": {"innerKey": "innerValue"}}`)
//	path := "nested.innerKey"
//	context := getBytes(jsonBytes, path)
//	fmt.Println("Unprocessed:", context.unprocessed) // Output: `{"key": "value", "nested": {"innerKey": "innerValue"}}`
//	fmt.Println("Strings:", context.strings)         // Output: `{"innerKey": "innerValue"}`
func getBytes(json []byte, path string) Context {
	var result Context
	if json != nil {
		// unsafe cast json bytes to a string and process it using the Get function.
		result = Get(*(*string)(unsafe.Pointer(&json)), path)
		// extract the string headers for unprocessed and strings.
		rawSafe := *(*stringHeader)(unsafe.Pointer(&result.unprocessed))
		stringSafe := *(*stringHeader)(unsafe.Pointer(&result.strings))
		// create byte slice headers for the unprocessed and strings.
		rawHeader := sliceHeader{data: rawSafe.data, length: rawSafe.length, capacity: rawSafe.length}
		sliceHeader := sliceHeader{data: stringSafe.data, length: stringSafe.length, capacity: rawSafe.length}
		// check for nil data and safely copy headers to strings if necessary.
		if sliceHeader.data == nil {
			if rawHeader.data == nil {
				result.unprocessed = ""
			} else {
				// unprocessed has data, safely copy the slice header to a string
				result.unprocessed = string(*(*[]byte)(unsafe.Pointer(&rawHeader)))
			}
			result.strings = ""
		} else if rawHeader.data == nil {
			result.unprocessed = ""
			result.strings = string(*(*[]byte)(unsafe.Pointer(&sliceHeader)))
		} else if uintptr(sliceHeader.data) >= uintptr(rawHeader.data) &&
			uintptr(sliceHeader.data)+uintptr(sliceHeader.length) <=
				uintptr(rawHeader.data)+uintptr(rawHeader.length) {
			// strings is a substring of unprocessed.
			start := uintptr(sliceHeader.data) - uintptr(rawHeader.data)
			// safely copy the raw slice header
			result.unprocessed = string(*(*[]byte)(unsafe.Pointer(&rawHeader)))
			result.strings = result.unprocessed[start : start+uintptr(sliceHeader.length)]
		} else {
			// safely copy both headers to strings.
			result.unprocessed = string(*(*[]byte)(unsafe.Pointer(&rawHeader)))
			result.strings = string(*(*[]byte)(unsafe.Pointer(&sliceHeader)))
		}
	}
	return result
}

// calcSubstringIndex calculates and assigns the starting index of the `unprocessed` field in the `value`
// field of the `parser` struct relative to the `json` string.
//
// Parameters:
//   - `json`: The complete JSON string from which the index is derived.
//   - `c`: A pointer to a `parser` instance containing the `value` with the `unprocessed` field.
//
// Behavior:
//   - If the `unprocessed` field in `value` is non-empty and the `calc` flag in the parser is false:
//     1. It computes the relative index of `unprocessed` within the `json` string by comparing
//     the memory addresses of the respective string headers.
//     2. If the computed index is invalid (e.g., out of bounds of `json`), it sets the index to 0.
//
// Notes:
//   - The function uses unsafe operations to access the memory layout of strings and compute their
//     relative positions. This minimizes overhead but requires care to ensure memory safety.
//   - The `index` field is useful for tracking the position of a substring within the original JSON data.
//
// Example:
//
//	json := `{"key": "value"}`
//	value := Context{unprocessed: `"value"`}
//	c := &parser{json: json, value: value}
//	calcSubstringIndex(json, c)
//	fmt.Println(c.value.index) // Outputs the starting position of `"value"` in the JSON string.
func calcSubstringIndex(json string, c *parser) {
	if len(c.value.unprocessed) > 0 && !c.calc {
		jsonHeader := *(*stringHeader)(unsafe.Pointer(&json))
		unprocessedHeader := *(*stringHeader)(unsafe.Pointer(&(c.value.unprocessed)))
		c.value.index = int(uintptr(unprocessedHeader.data) - uintptr(jsonHeader.data))
		if c.value.index < 0 || c.value.index >= len(json) {
			c.value.index = 0
		}
	}
}

// fromStr2Bytes converts a string into a byte slice without allocating new memory for the data.
// This function uses unsafe operations to directly reinterpret the string's underlying data
// structure as a byte slice. This allows efficient access to the string's content as a mutable
// byte slice, but it also comes with risks.
//
// Parameters:
//   - `s`: The input string that needs to be converted to a byte slice.
//
// Returns:
//   - A byte slice (`[]byte`) that shares the same underlying data as the input string.
//
// Notes:
//   - This function leverages Go's `unsafe` package to bypass the usual safety mechanisms
//     of the Go runtime. It does this by manipulating memory layouts using `unsafe.Pointer`.
//   - The resulting byte slice must be treated with care. Modifying the byte slice can lead
//     to undefined behavior since strings in Go are immutable by design.
//   - Any operation that depends on the immutability of the original string should avoid using this function.
//
// Safety Considerations:
//   - Since this function operates on unsafe pointers, it is not portable across different
//     Go versions or architectures.
//   - Direct modifications to the returned byte slice will violate Go's immutability guarantees
//     for strings and may corrupt program state.
//
// Example Usage:
//
//	s := "immutable string"
//	b := fromStr2Bytes(s) // Efficiently converts the string to []byte
//	// WARNING: Modifying 'b' here can lead to undefined behavior.
func fromStr2Bytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&sliceHeader{
		data:     (*stringHeader)(unsafe.Pointer(&s)).data,
		length:   len(s),
		capacity: len(s),
	}))
}

// fromBytes2Str converts a byte slice into a string without allocating new memory for the data.
// This function uses unsafe operations to directly reinterpret the byte slice's underlying data
// structure as a string. This allows efficient conversion without copying the data.
//
// Parameters:
//   - `b`: The input byte slice that needs to be converted to a string.
//
// Returns:
//   - A string that shares the same underlying data as the input byte slice.
//
// Notes:
//   - This function leverages Go's `unsafe` package to bypass the usual safety mechanisms
//     of the Go runtime. It does this by directly converting the memory representation of
//     the byte slice to that of a string.
//   - The resulting string must be treated with care. Modifying the original byte slice
//     after conversion will affect the string and can lead to undefined behavior.
//
// Safety Considerations:
//   - Strings in Go are meant to be immutable, but this function creates a string that
//     shares the same underlying data as the mutable byte slice. Any modifications to the
//     byte slice will reflect in the string, violating immutability guarantees.
//   - Use this function only in performance-critical scenarios where avoiding memory allocations
//     is essential, and ensure the byte slice is not modified afterward.
//
// Example Usage:
//
//	b := []byte{'h', 'e', 'l', 'l', 'o'}
//	s := fromBytes2Str(b) // Efficiently converts the byte slice to a string
//	fmt.Println(s) // Output: "hello"
//	// WARNING: Modifying 'b' here will also modify 's', leading to unexpected behavior.
func fromBytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// lowerPrefix extracts the initial contiguous sequence of lowercase alphabetic characters
// ('a' to 'z') from the input string `json`. It stops when it encounters a character
// outside this range and returns the substring up to that point.
//
// Parameters:
//   - `json`: The input string from which the initial sequence of lowercase alphabetic
//     characters is extracted.
//
// Returns:
//   - `raw`: A substring of `json` containing only the leading lowercase alphabetic characters.
//     If the input string starts with no lowercase alphabetic characters, the return
//     value will be an empty string.
//
// Notes:
//   - The function starts iterating from the second character (index 1) since it assumes the
//     first character does not need validation or extraction.
//   - The comparison checks (`json[i] < 'a' || json[i] > 'z'`) ensure that the character
//     falls outside the range of lowercase alphabetic ASCII values.
//
// Example Usage:
//
//	s := "abc123xyz"
//	result := lowerPrefix(s) // result: "abc"
//
//	s = "123abc"
//	result := lowerPrefix(s) // result: ""
//
//	s = "only.lowercase"
//	result := lowerPrefix(s) // result: "only.lowercase"
func lowerPrefix(json string) (raw string) {
	for i := 1; i < len(json); i++ {
		if json[i] < 'a' || json[i] > 'z' {
			return json[:i]
		}
	}
	return json
}

// squash processes a JSON string and returns a version of it that "squashes" or reduces its
// structure by removing nested objects and arrays, stopping at the top-level object or array.
// It assumes that the input string starts with a valid JSON structure, i.e., a '[' (array),
// '{' (object), '(' (another variant of object), or '"' (string). The function also ignores
// escaped characters (like quotes within strings) to avoid premature termination of the string.
//
// Parameters:
//   - `json`: The input JSON string that needs to be squashed by ignoring all nested structures.
//
// Returns:
//   - A string that contains the squashed JSON, stopping at the first level of nested objects/arrays.
//     If no nested structures exist, it will return the original string.
//
// Notes:
//   - The function expects the input JSON string to begin with a valid starting character such as
//     an opening quote ('"'), opening brace ('{'), opening bracket ('['), or opening parenthesis ('(').
//     It processes the string to remove all nested objects or arrays, stopping when it reaches the
//     corresponding closing brace, bracket, or parenthesis.
//   - It tracks depth using the `depth` variable, which increments when encountering an opening
//     brace, bracket, or parenthesis, and decrements when encountering the corresponding closing
//     brace, bracket, or parenthesis.
//   - If the function encounters a string enclosed in double quotes, it processes the string carefully,
//     skipping over any escaped quotes to avoid breaking the structure prematurely.
//   - Once the top-level structure (array, object, or string) is fully processed, it returns the squashed
//     JSON as a string, effectively ignoring the contents of any nested arrays or objects.
//
// Example Usage:
//
//	json := `{"key": [1, 2, {"nestedKey": "value"}]}`
//	result := squash(json) // result: '{"key": [1, 2, {"nestedKey": "value"}]}'
//
//	json := `{"key": {"innerKey": "value"}}`
//	result := squash(json) // result: '{"key": {"innerKey": "value"}}'
func squash(json string) string {
	var i, depth int
	// If the first character is not a quote, initialize i and depth for the JSON object/array parsing.
	if json[0] != '"' {
		i, depth = 1, 1
	}
	// Iterate through the string starting from index 1 to process the content.
	for ; i < len(json); i++ {
		// Process characters that are within the range of valid JSON characters (from '"' to '}').
		if json[i] >= '"' && json[i] <= '}' {
			switch json[i] {
			// Handle string literals, ensuring to escape any escaped quotes inside.
			case '"':
				i++
				s2 := i
				for ; i < len(json); i++ {
					if json[i] > '\\' {
						continue
					}
					// If an unescaped quote is found, break out of the loop.
					if json[i] == '"' {
						// look for an escaped slash
						if json[i-1] == '\\' {
							n := 0
							// Count the number of preceding backslashes.
							for j := i - 2; j > s2-1; j-- {
								if json[j] != '\\' {
									break
								}
								n++
							}
							// If there is an even number of backslashes, continue, as this quote is escaped.
							if n%2 == 0 {
								continue
							}
						}
						// If quote is found and it's not escaped, break the loop.
						break
					}
				}
				// If depth is 0, we've finished processing the top-level string, return it.
				if depth == 0 {
					if i >= len(json) {
						return json
					}
					return json[:i+1]
				}
			// Process nested objects/arrays (opening braces or brackets).
			case '{', '[', '(':
				depth++
			// Process closing of nested objects/arrays (closing braces, brackets, or parentheses).
			case '}', ']', ')':
				depth--
				// If depth becomes 0, we've reached the end of the top-level object/array.
				if depth == 0 {
					return json[:i+1]
				}
			}
		}
	}
	return json
}

// unescape takes a JSON-encoded string as input and processes any escape sequences (e.g., \n, \t, \u) within it,
// returning a new string with the escape sequences replaced by their corresponding characters.
//
// Parameters:
//   - `json`: A string representing a JSON-encoded string which may contain escape sequences (e.g., `\n`, `\"`).
//
// Returns:
//   - A new string where all escape sequences in the input JSON string are replaced by their corresponding
//     character representations. If an escape sequence is invalid or incomplete, it returns the string up to
//     that point without applying the escape.
//
// Notes:
//   - The function processes escape sequences commonly found in JSON, such as `\\`, `\/`, `\b`, `\f`, `\n`, `\r`, `\t`, `\"`, and `\u` (Unicode).
//   - If an invalid or incomplete escape sequence is encountered (for example, an incomplete Unicode sequence), it returns the string up to that point.
//   - If a non-printable character (less than ASCII value 32) is encountered, the function terminates early and returns the string up to that point.
//   - The function handles Unicode escape sequences (e.g., `\uXXXX`) by decoding them into their respective Unicode characters and converting them into UTF-8.
//
// Example Usage:
//
//	input := "\"Hello\\nWorld\""
//	result := unescape(input)
//	// result: "Hello\nWorld"
//
//	input := "\"Unicode \\u0048\\u0065\\u006C\\u006C\\u006F\""
//	result := unescape(input)
//	// result: "Unicode Hello"
func unescape(json string) string {
	var str = make([]byte, 0, len(json))
	for i := 0; i < len(json); i++ {
		switch {
		default:
			str = append(str, json[i])
		case json[i] < ' ': // If the character is a non-printable character (ASCII value less than 32), terminate early.
			return string(str)
		case json[i] == '\\': // If the current character is a backslash, process the escape sequence.
			i++
			if i >= len(json) {
				return string(str)
			}
			switch json[i] {
			default:
				return string(str)
			case '\\':
				str = append(str, '\\')
			case '/':
				str = append(str, '/')
			case 'b':
				str = append(str, '\b')
			case 'f':
				str = append(str, '\f')
			case 'n':
				str = append(str, '\n')
			case 'r':
				str = append(str, '\r')
			case 't':
				str = append(str, '\t')
			case '"':
				str = append(str, '"')
			case 'u': // Handle Unicode escape sequences (\uXXXX).
				if i+5 > len(json) {
					return string(str)
				}
				r := hexToRune(json[i+1:]) // Decode the Unicode code point (assuming `goRune` is a helper function).
				i += 5
				if utf16.IsSurrogate(r) { // Check for surrogate pairs (used for characters outside the Basic Multilingual Plane).
					// If a second surrogate is found, decode it into the correct rune.
					if len(json[i:]) >= 6 && json[i] == '\\' &&
						json[i+1] == 'u' {
						// Decode the second part of the surrogate pair.
						r = utf16.DecodeRune(r, hexToRune(json[i+2:]))
						i += 6
					}
				}
				// Allocate enough space to encode the decoded rune as UTF-8.
				str = append(str, 0, 0, 0, 0, 0, 0, 0, 0)
				// Encode the rune as UTF-8 and append it to the result slice.
				n := utf8.EncodeRune(str[len(str)-8:], r)
				str = str[:len(str)-8+n]
				i-- // Backtrack the index to account for the additional character read.
			}
		}
	}
	return string(str)
}

// hexToRune converts a hexadecimal Unicode escape sequence (represented as a string)
// into the corresponding Unicode code point (rune).
//
// This function expects a string containing a 4-digit hexadecimal number that represents
// a Unicode code point (e.g., "0048" for the letter 'H') and converts it to a rune (i.e., a Unicode code point).
//
// Parameters:
//   - `json`: A string containing the first 4 characters of the Unicode escape sequence in hexadecimal format
//     (e.g., `json` would be `"0048"` for the Unicode code point 'H').
//
// Returns:
//   - A rune corresponding to the Unicode code point represented by the provided hexadecimal string.
//
// Notes:
//   - The function assumes that the input string (`json`) is at least 4 characters long and contains valid
//     hexadecimal digits (e.g., "0048"). If the input string is shorter or invalid, the function will panic or behave
//     unpredictably. In production code, input validation should be added to handle such cases safely.
//   - The function only parses the first 4 characters of the input string as a 16-bit hexadecimal number, suitable
//     for representing Basic Multilingual Plane (BMP) characters (Unicode code points U+0000 to U+FFFF). For surrogate pairs
//     (characters outside the BMP), additional handling is required.
//
// Example Usage:
//
//		input := "0048" // Hexadecimal for Unicode character 'H'
//		result := hexToRune(input)
//		// result: 'H' (rune corresponding to U+0048)
//
//	  Note: This function is specifically designed to handle only the first 4 characters of a Unicode escape sequence.
func hexToRune(json string) rune {
	n, _ := strconv.ParseUint(json[:4], 16, 64)
	return rune(n)
}

// lessInsensitive compares two strings a and b in a case-insensitive manner.
// It returns true if string a is lexicographically less than string b, ignoring case differences.
// If both strings are equal in a case-insensitive comparison, it returns false.
//
// Parameters:
//   - `a`: The first string to compare.
//   - `b`: The second string to compare.
//
// Returns:
//   - `true`: If string a is lexicographically smaller than string b in a case-insensitive comparison.
//   - `false`: Otherwise.
//
// Notes:
//   - The function compares the strings character by character. If both characters are uppercase, they are compared directly.
//   - If one character is uppercase and the other is lowercase, the uppercase character is treated as the corresponding lowercase character.
//   - If neither character is uppercase, they are compared directly without any transformation.
//   - The function handles cases where the strings have different lengths and returns true if the strings are equal up to the point where the shorter string ends.
//
// Example Usage:
//
//	result := lessInsensitive("apple", "Apple")
//	// result: false, because "apple" and "Apple" are considered equal when case is ignored
//
//	result := lessInsensitive("apple", "banana")
//	// result: true, because "apple" is lexicographically smaller than "banana"
func lessInsensitive(a, b string) bool {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] >= 'A' && a[i] <= 'Z' {
			if b[i] >= 'A' && b[i] <= 'Z' {
				// both are uppercase, do nothing
				if a[i] < b[i] {
					return true
				} else if a[i] > b[i] {
					return false
				}
			} else {
				// a is uppercase, convert a to lowercase
				if a[i]+32 < b[i] {
					return true
				} else if a[i]+32 > b[i] {
					return false
				}
			}
		} else if b[i] >= 'A' && b[i] <= 'Z' {
			// b is uppercase, convert b to lowercase
			if a[i] < b[i]+32 {
				return true
			} else if a[i] > b[i]+32 {
				return false
			}
		} else {
			// neither are uppercase
			if a[i] < b[i] {
				return true
			} else if a[i] > b[i] {
				return false
			}
		}
	}
	return len(a) < len(b)
}
