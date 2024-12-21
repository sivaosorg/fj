package fj

import (
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
	"unsafe"

	"github.com/sivaosorg/unify4g"
)

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

// calcSubstring calculates and assigns the starting index of the `unprocessed` field in the `value`
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
//	calcSubstring(json, c)
//	fmt.Println(c.value.index) // Outputs the starting position of `"value"` in the JSON string.
func calcSubstring(json string, c *parser) {
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
				r := hex2Rune(json[i+1:]) // Decode the Unicode code point (assuming `goRune` is a helper function).
				i += 5
				if utf16.IsSurrogate(r) { // Check for surrogate pairs (used for characters outside the Basic Multilingual Plane).
					// If a second surrogate is found, decode it into the correct rune.
					if len(json[i:]) >= 6 && json[i] == '\\' &&
						json[i+1] == 'u' {
						// Decode the second part of the surrogate pair.
						r = utf16.DecodeRune(r, hex2Rune(json[i+2:]))
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

// hex2Rune converts a hexadecimal Unicode escape sequence (represented as a string)
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
//		result := hex2Rune(input)
//		// result: 'H' (rune corresponding to U+0048)
//
//	  Note: This function is specifically designed to handle only the first 4 characters of a Unicode escape sequence.
func hex2Rune(json string) rune {
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

// verifyBoolTrue checks if the given byte slice starting at index i represents the string "true".
// It returns the next index after "true" and true if the sequence matches, otherwise it returns the current index and false.
//
// Parameters:
//   - data: A byte slice containing the data to validate.
//   - i: The index in the byte slice to start checking from.
//
// Returns:
//   - value: The index immediately after the "true" string if it matches, or the current index if it doesn't.
//   - ok: A boolean indicating whether the "true" string was found starting at index i.
//
// Notes:
//   - The function checks if the characters at positions i, i+1, and i+2 correspond to the letters 't', 'r', and 'u'.
//     If the substring matches the word "true", it returns the index after the "true" string and true.
//   - If the substring does not match "true", it returns the current index i and false.
//
// Example Usage:
//
//	data := []byte("this is true")
//	i := 10
//	value, ok := verifyBoolTrue(data, i)
//	// value: 13 (the index after the word "true")
//	// ok: true (because "true" was found starting at index 10)
func verifyBoolTrue(data []byte, i int) (val int, ok bool) {
	if i+3 <= len(data) && data[i] == 'r' && data[i+1] == 'u' &&
		data[i+2] == 'e' {
		return i + 3, true
	}
	return i, false
}

// verifyBoolFalse checks if the given byte slice starting at index i represents the string "false".
// It returns the next index after "false" and true if the sequence matches, otherwise it returns the current index and false.
//
// Parameters:
//   - data: A byte slice containing the data to validate.
//   - i: The index in the byte slice to start checking from.
//
// Returns:
//   - val: The index immediately after the "false" string if it matches, or the current index if it doesn't.
//   - ok: A boolean indicating whether the "false" string was found starting at index i.
//
// Notes:
//   - The function checks if the characters at positions i, i+1, i+2, and i+3 correspond to the letters 'f', 'a', 'l', and 's', respectively.
//     If the substring matches the word "false", it returns the index after the "false" string and true.
//   - If the substring does not match "false", it returns the current index i and false.
//
// Example Usage:
//
//	data := []byte("this is false")
//	i := 8
//	val, ok := verifyBoolFalse(data, i)
//	// val: 13 (the index after the word "false")
//	// ok: true (because "false" was found starting at index 8)
func verifyBoolFalse(data []byte, i int) (val int, ok bool) {
	if i+4 <= len(data) && data[i] == 'a' && data[i+1] == 'l' &&
		data[i+2] == 's' && data[i+3] == 'e' {
		return i + 4, true
	}
	return i, false
}

// verifyNullable checks if the given byte slice starting at index i represents the string "null".
// It returns the next index after "null" and true if the sequence matches, otherwise it returns the current index and false.
//
// Parameters:
//   - data: A byte slice containing the data to validate.
//   - i: The index in the byte slice to start checking from.
//
// Returns:
//   - val: The index immediately after the "null" string if it matches, or the current index if it doesn't.
//   - ok: A boolean indicating whether the "null" string was found starting at index i.
//
// Notes:
//   - The function checks if the characters at positions i, i+1, and i+2 correspond to the letters 'n', 'u', and 'l', respectively.
//     If the substring matches the word "null", it returns the index after the "null" string and true.
//   - If the substring does not match "null", it returns the current index i and false.
//
// Example Usage:
//
//	data := []byte("value is null")
//	i := 9
//	val, ok := verifyNullable(data, i)
//	// val: 13 (the index after the word "null")
//	// ok: true (because "null" was found starting at index 9)
func verifyNullable(data []byte, i int) (val int, ok bool) {
	if i+3 <= len(data) && data[i] == 'u' && data[i+1] == 'l' &&
		data[i+2] == 'l' {
		return i + 3, true
	}
	return i, false
}

// verifyNumeric validates whether the byte slice starting at index i represents a valid numeric value.
// It supports integer, floating-point, and exponential number formats as per JSON specifications.
//
// Parameters:
//   - data: A byte slice containing the input to validate.
//   - i: The starting index to check in the byte slice.
//
// Returns:
//   - val: The index immediately after the numeric value if it is valid, or the current index if it isn't.
//   - ok: A boolean indicating whether the input from index i represents a valid numeric value.
//
// Notes:
//   - This function validates numbers following the JSON number format, which includes:
//   - Optional sign ('-' for negative numbers).
//   - Integer component (digits, starting with '0' or other digits).
//   - Optional fractional part (a dot '.' followed by one or more digits).
//   - Optional exponent part ('e' or 'E', optionally signed, followed by one or more digits).
//   - The function iterates over the byte slice, checking each part of the number sequentially.
//   - If the numeric value is valid, the function returns the index after the number and true.
//     Otherwise, it returns the starting index and false.
//
// Example Usage:
//
//	data := []byte("-123.45e+6")
//	i := 1 // Start after the '-' sign
//	val, ok := verifyNumeric(data, i)
//	// val: 10 (the index after the number)
//	// ok: true (because "-123.45e+6" is a valid numeric value)
//
// Details:
//
//   - The function handles three major components of a number: sign, integer part, and optional components
//     (fractional and exponential parts).
//   - Each component is validated, and the function exits early with a false result if a part is invalid.
func verifyNumeric(data []byte, i int) (val int, ok bool) {
	// Check if i is within valid range
	if i <= 0 || i >= len(data) {
		return i, false
	}
	i--
	// Check for a sign ('-') at the start of the number.
	if data[i] == '-' {
		i++
		// A sign without any digits is invalid.
		// The character after the sign must be a digit.
		if i == len(data) || data[i] < '0' || data[i] > '9' {
			return i, false
		}
	}
	// Validate the integer part of the number.
	if i == len(data) {
		return i, false
	}
	if data[i] == '0' {
		// A leading '0' is valid but must not be followed by other digits.
		i++
	} else {
		// Consume digits in the integer part.
		for ; i < len(data); i++ {
			if data[i] >= '0' && data[i] <= '9' {
				continue
			}
			break
		}
	}
	// Validate the fractional part, if present.
	if i == len(data) {
		return i, true
	}
	if data[i] == '.' {
		i++
		// A dot without digits after it is invalid.
		// The character after the dot must be a digit.
		if i == len(data) || data[i] < '0' || data[i] > '9' {
			return i, false
		}
		i++
		// Consume digits in the fractional part.
		for ; i < len(data); i++ {
			if data[i] >= '0' && data[i] <= '9' {
				continue
			}
			break
		}
	}
	// Validate the exponential part, if present.
	if i == len(data) {
		return i, true
	}
	if data[i] == 'e' || data[i] == 'E' {
		i++
		if i == len(data) {
			// An 'e' or 'E' without any exponent value is invalid.
			return i, false
		}
		// Check for an optional sign in the exponent.
		if data[i] == '+' || data[i] == '-' {
			i++
		}
		// A sign without any digits in the exponent is invalid.
		// The character after the exponent must be a digit.
		if i == len(data) || data[i] < '0' || data[i] > '9' {
			return i, false
		}
		i++
		// Consume digits in the exponent part.
		for ; i < len(data); i++ {
			if data[i] >= '0' && data[i] <= '9' {
				continue
			}
			break
		}
	}
	return i, true
}

// verifyString validates whether the byte slice starting at index i represents a valid JSON string.
// The function ensures the string adheres to the JSON string format, including proper escaping of special characters.
//
// Parameters:
//   - data: A byte slice containing the input to validate.
//   - i: The starting index to check in the byte slice.
//
// Returns:
//   - val: The index immediately after the string if it is valid, or the current index if it isn't.
//   - ok: A boolean indicating whether the input from index i represents a valid JSON string.
//
// Notes:
//   - JSON strings must start and end with double quotes ('"').
//   - The function handles escaped characters such as '\\', '\"', and unicode escapes (e.g., '\\u1234').
//   - The function iterates over the byte slice, validating each character and ensuring proper escape sequences.
//   - If the string is valid, the function returns the index after the closing double quote and true.
//     Otherwise, it returns the current index and false.
//
// Example Usage:
//
//	data := []byte("\"Hello, \\\"world!\\\"\"")
//	i := 0 // Start at the first character
//	val, ok := verifyString(data, i)
//	// val: 20 (the index after the string)
//	// ok: true (because "\"Hello, \\\"world!\\\"\"" is a valid JSON string)
//
// Details:
//   - The function iterates over the characters, checking for valid JSON string content.
//   - It handles special escape sequences, ensuring their correctness.
//   - Early exit occurs if any invalid sequence or character is detected.
func verifyString(data []byte, i int) (val int, ok bool) {
	for ; i < len(data); i++ {
		if data[i] < ' ' {
			return i, false
		} else if data[i] == '\\' {
			i++
			if i == len(data) {
				return i, false
			}
			switch data[i] {
			default:
				return i, false
			case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
			case 'u':
				for j := 0; j < 4; j++ {
					i++
					if i >= len(data) {
						return i, false
					}
					if !((data[i] >= '0' && data[i] <= '9') ||
						(data[i] >= 'a' && data[i] <= 'f') ||
						(data[i] >= 'A' && data[i] <= 'F')) {
						return i, false
					}
				}
			}
		} else if data[i] == '"' {
			return i + 1, true
		}
	}
	return i, false
}

// verifyComma checks for the presence of a comma (',') or the specified end character in the given byte slice
// starting at index i. It skips over any whitespace characters and ensures valid structure.
//
// Parameters:
//   - data: A byte slice containing the input to validate.
//   - i: The starting index to check in the byte slice.
//   - end: The specific byte (character) to treat as a valid stopping point, in addition to a comma.
//
// Returns:
//   - val: The index of the comma or end character if found, or the current index if invalid.
//   - ok: A boolean indicating whether a valid comma or end character was found.
//
// Notes:
//   - Whitespace characters (' ', '\t', '\n', '\r') are skipped during validation.
//   - The function exits early with false if an invalid character is encountered before finding a comma or end character.
//   - If the comma or end character is found, the function returns its index and true.
//
// Example Usage:
//
//	data := []byte(" , next")
//	i := 0
//	end := byte('n')
//	val, ok := verifyComma(data, i, end)
//	// val: 1 (the index of the comma)
//	// ok: true (because a comma was found)
//
// Details:
//   - Iterates over characters, skipping valid whitespace.
//   - Checks for either a comma or the specified end character.
//   - Returns false if an invalid character is encountered.
func verifyComma(data []byte, i int, end byte) (val int, ok bool) {
	for ; i < len(data); i++ {
		switch data[i] {
		default:
			return i, false
		case ' ', '\t', '\n', '\r':
			continue
		case ',':
			return i, true
		case end:
			return i, true
		}
	}
	return i, false
}

// verifyColon checks for the presence of a colon (':') in the given byte slice starting at index i.
// It skips over any whitespace characters and ensures valid JSON structure.
//
// Parameters:
//   - data: A byte slice containing the input to validate.
//   - i: The starting index to check in the byte slice.
//
// Returns:
//   - val: The index immediately after the colon if found, or the current index if invalid.
//   - ok: A boolean indicating whether a valid colon was found.
//
// Notes:
//   - Whitespace characters (' ', '\t', '\n', '\r') are skipped during validation.
//   - The function exits early with false if an invalid character is encountered before finding the colon.
//   - If the colon is found, the function returns the index after it and true.
//
// Example Usage:
//
//	data := []byte(" : value")
//	i := 0
//	val, ok := verifyColon(data, i)
//	// val: 2 (the index after the colon)
//	// ok: true (because a colon was found)
//
// Details:
//   - Iterates over characters, skipping valid whitespace.
//   - Checks for a colon and returns the next index upon finding it.
//   - Returns false if an invalid character is encountered.
func verifyColon(data []byte, i int) (val int, ok bool) {
	for ; i < len(data); i++ {
		switch data[i] {
		default:
			return i, false
		case ' ', '\t', '\n', '\r':
			continue
		case ':':
			return i + 1, true
		}
	}
	return i, false
}

// verifyArray validates whether the byte slice starting at index i represents a valid JSON array.
// It ensures that the array starts and ends with square brackets ('[' and ']') and contains valid JSON values.
//
// Parameters:
//   - data: A byte slice containing the input to validate.
//   - i: The starting index to check in the byte slice.
//
// Returns:
//   - val: The index immediately after the valid array if found, or the current index if it isn't.
//   - ok: A boolean indicating whether the input from index i represents a valid JSON array.
//
// Notes:
//   - The function handles arrays that may contain:
//   - Whitespace (skipped).
//   - Comma-separated JSON values, validated using the `validateAny` function.
//   - Empty arrays ([]).
//   - The function ensures that the array ends with a closing square bracket (']').
//
// Example Usage:
//
//	data := []byte("[123, \"string\", false]")
//	i := 0
//	val, ok := verifyArray(data, i)
//	// val: 21 (the index after the array)
//	// ok: true (because the input is a valid JSON array)
//
// Details:
//   - Skips leading whitespace.
//   - Checks for an initial ']' to handle empty arrays.
//   - Iteratively validates JSON values and ensures proper use of commas.
//   - Returns false if an invalid character or structure is encountered.
func verifyArray(data []byte, i int) (val int, ok bool) {
	for ; i < len(data); i++ {
		switch data[i] {
		default:
			for ; i < len(data); i++ {
				if i, ok = verifyAny(data, i); !ok {
					return i, false
				}
				if i, ok = verifyComma(data, i, ']'); !ok {
					return i, false
				}
				if data[i] == ']' {
					return i + 1, true
				}
			}
		case ' ', '\t', '\n', '\r':
			continue
		case ']':
			return i + 1, true
		}
	}
	return i, false
}

// verifyObject validates whether the byte slice starting at index i represents a valid JSON object.
// It ensures that the object starts and ends with curly braces ('{' and '}') and contains valid key-value pairs.
//
// Parameters:
//   - data: A byte slice containing the input to validate.
//   - i: The starting index to check in the byte slice.
//
// Returns:
//   - val: The index immediately after the valid object if found, or the current index if it isn't.
//   - ok: A boolean indicating whether the input from index i represents a valid JSON object.
//
// Notes:
//   - The function handles objects that may contain:
//   - Whitespace (skipped).
//   - Key-value pairs, where keys are JSON strings and values are validated using `validateAny`.
//   - Empty objects ({}).
//   - Ensures proper use of colons (:) and commas (,) in separating keys and values.
//   - The function iteratively validates keys and values until the closing curly brace ('}') is found.
//
// Example Usage:
//
//	data := []byte(`{"key1": 123, "key2": "value"}`)
//	i := 0
//	val, ok := verifyObject(data, i)
//	// val: 28 (the index after the object)
//	// ok: true (because the input is a valid JSON object)
//
// Details:
//   - Skips leading whitespace.
//   - Validates the presence of keys (JSON strings) and their corresponding values.
//   - Ensures that the structure adheres to JSON object syntax, returning false for any invalid structure.
func verifyObject(data []byte, i int) (val int, ok bool) {
	for ; i < len(data); i++ {
		switch data[i] {
		default:
			return i, false
		case ' ', '\t', '\n', '\r':
			continue
		case '}':
			return i + 1, true
		case '"':
		key:
			if i, ok = verifyString(data, i+1); !ok {
				return i, false
			}
			if i, ok = verifyColon(data, i); !ok {
				return i, false
			}
			if i, ok = verifyAny(data, i); !ok {
				return i, false
			}
			if i, ok = verifyComma(data, i, '}'); !ok {
				return i, false
			}
			if data[i] == '}' {
				return i + 1, true
			}
			i++
			for ; i < len(data); i++ {
				switch data[i] {
				default:
					return i, false
				case ' ', '\t', '\n', '\r':
					continue
				case '"':
					goto key
				}
			}
			return i, false
		}
	}
	return i, false
}

// verifyAny attempts to validate the data starting at index i as one of the possible JSON value types.
// It recognizes and validates the following JSON types:
//   - Object: represented by curly braces `{}`
//   - Array: represented by square brackets `[]`
//   - String: represented by double quotes `""`
//   - Numeric values: including integers and floating-point numbers
//   - Boolean values: `true` or `false`
//   - Null: represented by `null`
//
// Parameters:
//   - data: A byte slice containing the JSON input to validate.
//   - i: The starting index in the byte slice where the validation should begin.
//
// Returns:
//   - val: The index immediately after the valid value if found, or the current index if it isn't.
//   - ok: A boolean indicating whether the input starting at index i is a valid JSON value of one of the recognized types.
//
// Notes:
//   - The function handles a variety of JSON data types, attempting to validate the input by matching
//     the character at the current index and calling the appropriate validation function for the recognized type.
//   - It will skip any whitespace characters (spaces, tabs, newlines, carriage returns) before checking the data.
//   - The function calls other helper functions to validate specific types of JSON values, such as `verifyObject`, `verifyArray`,
//     `verifyString`, `verifyNumeric`, `verifyBoolTrue`, `verifyBoolFalse`, and `verifyNullable`.
//
// Example Usage:
//
//	data := []byte(`{"key1": 123, "key2": "value"}`)
//	i := 0
//	val, ok := verifyAny(data, i)
//	// val: 28 (the index after the object)
//	// ok: true (because the input is a valid JSON object)
//
// Details:
//   - If the input data at index i is a valid JSON value (object, array, string, numeric, boolean, or null),
//     the function will return the index immediately after the valid value and true.
//   - If the data does not match a valid JSON value, it returns the current index and false.
func verifyAny(data []byte, i int) (val int, ok bool) {
	for ; i < len(data); i++ {
		switch data[i] {
		default:
			return i, false
		case ' ', '\t', '\n', '\r':
			continue
		case '{':
			return verifyObject(data, i+1)
		case '[':
			return verifyArray(data, i+1)
		case '"':
			return verifyString(data, i+1)
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return verifyNumeric(data, i+1)
		case 't':
			return verifyBoolTrue(data, i+1)
		case 'f':
			return verifyBoolFalse(data, i+1)
		case 'n':
			return verifyNullable(data, i+1)
		}
	}
	return i, false
}

// verifyJson attempts to validate the data starting at index i as a valid JSON payload. It checks for
// the presence of valid JSON values after skipping any whitespace characters (spaces, tabs, newlines, etc.).
// It calls the `verifyAny` function to validate the first JSON value and ensures that there are no unexpected
// characters after the valid value.
//
// Parameters:
//   - data: A byte slice containing the JSON input to validate.
//   - i: The starting index in the byte slice where the validation should begin.
//
// Returns:
//   - val: The index immediately after the validated JSON payload, or the current index if validation fails.
//   - ok: A boolean indicating whether the payload is valid. If true, the input is a valid JSON payload;
//     if false, the payload is not valid.
//
// Notes:
//   - The function starts by checking the first non-whitespace character in the data. It skips over any
//     whitespace (spaces, tabs, newlines, carriage returns) before trying to validate the payload.
//   - It then calls the `verifyAny` function to check for a valid JSON value (object, array, string, numeric,
//     boolean, or null) at the current index.
//   - After the first valid value is found, the function ensures that the rest of the input contains only
//     whitespace (ignoring spaces, tabs, newlines) before validating that the payload ends correctly.
//   - If the first value is valid and no unexpected characters are found, it returns true along with the
//     index after the valid value. If any issue is encountered, it returns false.
//
// Example Usage:
//
//	data := []byte(`{"key1": 123, "key2": "value"}`)
//	i := 0
//	val, ok := verifyJson(data, i)
//	// val: 28 (index after the valid payload)
//	// ok: true (because the input is a valid JSON payload)
//
// Details:
//   - The function ensures that the input data starts with a valid JSON value, as recognized by the `verifyAny`
//     function, and that there are no unexpected characters after that value.
//   - The function returns false if the JSON payload is incomplete or invalid.
func verifyJson(data []byte, i int) (val int, ok bool) {
	for ; i < len(data); i++ { // Iterate through the data starting from index i.
		// Handle unexpected characters in the payload.
		switch data[i] {
		default:
			// If the character is unexpected, call verifyAny to validate the first valid JSON value.
			i, ok = verifyAny(data, i)
			if !ok {
				return i, false // Return false if the value is not valid.
			}
			// After a valid JSON value, continue to check if the rest of the data only contains whitespace.
			for ; i < len(data); i++ {
				switch data[i] {
				default:
					return i, false // Return false if there is any invalid character after the value.
				case ' ', '\t', '\n', '\r': // Skip over whitespace characters.
					continue
				}
			}
			// If all subsequent characters are whitespace, return true along with the index after the valid value.
			return i, true
		case ' ', '\t', '\n', '\r': // Skip over whitespace characters (spaces, tabs, newlines, carriage returns).
			continue
		}
	}
	return i, false // Return false if the end of data is reached without a valid payload.
}

// lastSegment extracts the last part of a given path string, where the path segments are separated by
// either a pipe ('|') or a dot ('.'). The function returns the substring after the last separator,
// taking escape sequences (backslashes) into account. It ensures that any escaped separator is ignored.
//
// Parameters:
//   - path: A string representing the full path, which may contain segments separated by '|' or '.'.
//
// Returns:
//   - A string representing the last segment in the path after the last occurrence of either '|' or '.'.
//     If no separator is found, it returns the entire input string.
//
// Notes:
//   - The function handles escape sequences where separators are preceded by a backslash ('\').
//   - If there is no valid separator in the string, the entire path is returned as-is.
//   - The returned substring is the part after the last separator, which could be the last portion of the path.
//
// Example Usage:
//
//	path := "foo|bar.baz.qux"
//	segment := lastSegment(path)
//	// segment: "qux" (the last segment after the last dot or pipe)
//
// Details:
//   - The function iterates from the end of the string towards the beginning, looking for the last
//     occurrence of '|' or '.' that is not preceded by a backslash.
//   - It handles edge cases where the separator is escaped or there are no separators at all.
func lastSegment(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '|' || path[i] == '.' {
			if i > 0 {
				if path[i-1] == '\\' {
					continue
				}
			}
			return path[i+1:]
		}
	}
	return path
}

// isValidName checks if a given string component is a "simple name" according to specific rules.
// A "simple name" is a string that does not contain any control characters or any of the following special characters:
// '[' , ']' , '{' , '}' , '(' , ')' , '#' , '|' , '!'. The function returns true if the string meets these criteria.
//
// Parameters:
//   - component: A string to be checked for validity as a simple name.
//
// Returns:
//   - A boolean indicating whether the input string is a valid simple name.
//   - Returns true if the string contains only printable characters and does not include any of the restricted special characters.
//   - Returns false if the string contains any control characters or restricted special characters.
//
// Notes:
//   - The function checks each character of the string to ensure it is printable and does not contain any of the restricted characters.
//   - Control characters are defined as any character with a Unicode value less than a space (' ').
//   - The function assumes that the string is not empty and contains at least one character.
//
// Example Usage:
//
//	component := "validName"
//	isValid := isValidName(component)
//	// isValid: true (the string contains only valid characters)
//
//	component = "invalid|name"
//	isValid = isValidName(component)
//	// isValid: false (the string contains an invalid character '|')
//
// Details:
//   - The function iterates through each character of the string and checks whether it is a printable character and whether it
//     is not one of the restricted special characters. If any invalid character is found, the function returns false immediately.
func isValidName(component string) bool {
	if unify4g.IsEmpty(component) {
		return false
	}
	if unify4g.ContainsAny(component, " ") {
		return false
	}
	for i := 0; i < len(component); i++ {
		if component[i] < ' ' {
			return false
		}
		switch component[i] {
		case '[', ']', '{', '}', '(', ')', '#', '|', '!':
			return false
		}
	}
	return true
}

// appendHex16 appends the hexadecimal representation of a 16-bit unsigned integer (uint16)
// to a byte slice. The integer is converted to a 4-character hexadecimal string, and each character
// is appended to the input byte slice in sequence. The function uses a pre-defined set of hexadecimal
// digits ('0'–'9' and 'a'–'f') for the conversion.
//
// Parameters:
//   - bytes: A byte slice to which the hexadecimal characters will be appended.
//   - x: A 16-bit unsigned integer to be converted to hexadecimal and appended to the byte slice.
//
// Returns:
//   - A new byte slice containing the original bytes with the appended hexadecimal digits
//     representing the 16-bit integer.
//
// Example Usage:
//
//	var result []byte
//	x := uint16(3055) // Decimal 3055 is 0x0BEF in hexadecimal
//	result = appendHex16(result, x)
//	// result: []byte{'0', 'b', 'e', 'f'} (hexadecimal representation of 3055)
//
// Details:
//   - The function shifts and masks the 16-bit integer to extract each of the four hexadecimal digits.
//   - It uses the pre-defined `hexDigits` array to convert the integer's nibbles (4 bits) into their
//     corresponding hexadecimal characters.
func appendHex16(bytes []byte, x uint16) []byte {
	return append(bytes,
		hexDigits[x>>12&0xF], hexDigits[x>>8&0xF],
		hexDigits[x>>4&0xF], hexDigits[x>>0&0xF],
	)
}

// parseUint parses a string as an unsigned integer (uint64).
// It attempts to convert the given string to a numeric value, where each character in the string
// must be a digit between '0' and '9'. If any non-digit character is encountered, the function
// returns false, indicating the string does not represent a valid unsigned integer.
//
// Parameters:
//   - s: A string representing the unsigned integer to be parsed.
//
// Returns:
//   - n: The parsed unsigned integer value (of type uint64) if the string represents a valid number.
//   - ok: A boolean indicating whether the parsing was successful. If true, the string was successfully
//     parsed into an unsigned integer; if false, the string was invalid.
//
// Example Usage:
//
//	str := "12345"
//	n, ok := parseUint(str)
//	// n: 12345 (the parsed unsigned integer)
//	// ok: true (the string is a valid unsigned integer)
//
//	str = "12a45"
//	n, ok = parseUint(str)
//	// n: 0 (parsing failed)
//	// ok: false (the string contains invalid characters)
//
// Details:
//   - The function iterates through each character of the string. If it encounters a digit ('0'–'9'),
//     it accumulates the corresponding integer value into the result `n`. The result is multiplied by 10
//     with each new digit to shift the previous digits left.
//   - If any non-digit character is encountered, the function returns `0` and `false`.
//   - The function assumes that the input string is non-empty and only contains valid ASCII digits if valid.
func parseUint(s string) (n uint64, ok bool) {
	var i int
	if i == len(s) {
		return 0, false
	}
	for ; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			n = n*10 + uint64(s[i]-'0')
		} else {
			return 0, false
		}
	}
	return n, true
}

// parseInt parses a string as a signed integer (int64).
// It attempts to convert the given string to a numeric value, where each character in the string
// must be a digit between '0' and '9'. The function also supports negative numbers, indicated by a leading
// minus sign ('-'). If any non-digit character is encountered (except for the minus sign at the start),
// the function returns false, indicating the string does not represent a valid signed integer.
//
// Parameters:
//   - s: A string representing the signed integer to be parsed.
//
// Returns:
//   - n: The parsed signed integer value (of type int64) if the string represents a valid number.
//   - ok: A boolean indicating whether the parsing was successful. If true, the string was successfully
//     parsed into a signed integer; if false, the string was invalid.
//
// Example Usage:
//
//	str := "-12345"
//	n, ok := parseInt(str)
//	// n: -12345 (the parsed signed integer)
//	// ok: true (the string is a valid signed integer)
//
//	str = "12a45"
//	n, ok = parseInt(str)
//	// n: 0 (parsing failed)
//	// ok: false (the string contains invalid characters)
//
// Details:
//   - The function first checks for an optional leading minus sign ('-'). If found, it sets a `sign`
//     flag to indicate the number is negative.
//   - It then iterates through each character of the string. If it encounters a digit ('0'–'9'),
//     it accumulates the corresponding integer value into the result `n`, shifting the previous digits left.
//   - If any non-digit character is encountered (excluding the leading minus sign), the function returns `0` and `false`.
//   - If the `sign` flag is set, the result is negated before returning.
//   - The function assumes that the input string is non-empty and contains valid digits if valid, with an optional leading minus sign.
func parseInt(s string) (n int64, ok bool) {
	var i int
	var sign bool
	if len(s) > 0 && s[0] == '-' {
		sign = true
		i++
	}
	if i == len(s) {
		return 0, false
	}
	for ; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			n = n*10 + int64(s[i]-'0')
		} else {
			return 0, false
		}
	}
	if sign {
		return n * -1, true
	}
	return n, true
}

// ensureSafeInt validates a given floating-point number (float64) to ensure it lies within the safe range for integers
// in JavaScript (the Number type), as defined by the ECMAScript specification. The function checks if the number
// is within the range of valid safe integers that can be accurately represented in JavaScript without loss of precision.
// If the number is within the safe integer range, it returns the corresponding signed 64-bit integer value (int64).
// If the number is outside the safe range, it returns false, indicating that the number cannot be safely represented as an integer.
//
// Parameters:
//   - f: A floating-point number (float64) representing the number to be validated.
//
// Returns:
//   - n: The corresponding signed 64-bit integer (int64) value if the number is within the safe integer range.
//   - ok: A boolean indicating whether the number lies within the safe integer range. If true, the number is safe to
//     represent as an integer; if false, the number is outside the safe range.
//
// Example Usage:
//
//	f := 1234567890.0
//	n, ok := ensureSafeInt(f)
//	// n: 1234567890 (the number as a valid int64)
//	// ok: true (the number is within the safe integer range)
//
//	f = 9007199254740992.0
//	n, ok = ensureSafeInt(f)
//	// n: 0 (parsing failed)
//	// ok: false (the number is outside the safe integer range)
//
// Details:
//
//   - The function checks if the given number `f` falls within the range of -9007199254740991 to 9007199254740991,
//     which are the minimum and maximum safe integers in JavaScript as specified in the ECMAScript standard.
//
//   - If the number is within the safe range, it is converted to an int64 and returned, along with `true` to indicate success.
//
//   - If the number is outside this range, it returns `0` and `false`, as such numbers cannot be safely represented as
//     integers in JavaScript without losing precision.
//
//   - https://tc39.es/ecma262/#sec-number.min_safe_integer
//
//   - https://tc39.es/ecma262/#sec-number.max_safe_integer
func ensureSafeInt(f float64) (n int64, ok bool) {
	if f < -9007199254740991 || f > 9007199254740991 {
		return 0, false
	}
	return int64(f), true
}

// parseStaticValue parses a string path to find a static value, such as a boolean, null, or number.
// The function expects that the input path starts with a '!', indicating a static value. It identifies the static
// value by looking for valid characters and structures that represent literal values in a path. If a valid static value
// is found, it returns the remaining path, the static value, and a success flag. Otherwise, it returns false.
//
// Parameters:
//   - path: A string representing the path to parse. The path should start with a '!' to indicate a static value.
//     The function processes the string following the '!' to find the static value.
//
// Returns:
//   - pathOut: The remaining part of the path after the static value has been identified. This is the portion of the
//     string that follows the literal value, such as any further path segments or operators.
//   - result: The static value found in the path, which can be a boolean ("true" or "false"), null, NaN, or Inf, or
//     a numeric value, or an empty string if no valid static value is found.
//   - ok: A boolean indicating whether the function successfully identified a static value. Returns true if a valid
//     static value is found, and false otherwise.
//
// Example Usage:
//
//	path := "!true.some.other.path"
//	pathOut, result, ok := parseStaticValue(path)
//	// pathOut: ".some.other.path" (remaining path)
//	// result: "true" (the static value found)
//	// ok: true (successful identification of static value)
//
//	path = "!123.abc"
//	pathOut, result, ok = parseStaticValue(path)
//	// pathOut: ".abc" (remaining path)
//	// result: "123" (the static value found)
//	// ok: true (successful identification of static value)
//
// Details:
//
//   - The function looks for the first character after the '!' to determine if the value starts with a valid static
//     value, such as a number or a boolean literal ("true", "false"), null, NaN, or Inf.
//
//   - It processes the string to extract the static value, then identifies the rest of the path (if any) after the static
//     value, which is returned as the remaining portion of the path.
//
//   - If the function encounters a delimiter like '.' or '|', it stops further parsing of the static value and returns
//     the remaining path.
//
//   - If no static value is identified, the function returns false.
//
//     Notes:
//
//   - The function assumes that the input path is well-formed and follows the expected format (starting with '!').
//
//   - The value can be a boolean, null, NaN, Inf, or a number in the path.
func parseStaticValue(path string) (pathStatic, result string, ok bool) {
	name := path[1:]
	if len(name) > 0 {
		switch name[0] {
		case '{', '[', '"', '+', '-', '0', '1', '2', '3', '4', '5', '6', '7',
			'8', '9':
			_, result = parseJSONSquash(name, 0)
			pathStatic = name[len(result):]
			return pathStatic, result, true
		}
	}
	for i := 1; i < len(path); i++ {
		if path[i] == '|' {
			pathStatic = path[i:]
			name = path[1:i]
			break
		}
		if path[i] == '.' {
			pathStatic = path[i:]
			name = path[1:i]
			break
		}
	}
	switch strings.ToLower(name) {
	case "true", "false", "null", "nan", "inf":
		return pathStatic, name, true
	}
	return pathStatic, result, false
}

// trim removes leading and trailing whitespace characters from a string.
// The function iteratively checks and removes spaces (or any character less than or equal to a space)
// from both the left (beginning) and right (end) of the string.
//
// Parameters:
//   - s: A string that may contain leading and trailing whitespace characters that need to be removed.
//
// Returns:
//   - A new string with leading and trailing whitespace removed. The function does not modify the original string,
//     as strings in Go are immutable.
//
// Example Usage:
//
//	str := "  hello world  "
//	trimmed := trim(str)
//	// trimmed: "hello world" (leading and trailing spaces removed)
//
//	str = "\n\n   trim me   \t\n"
//	trimmed = trim(str)
//	// trimmed: "trim me" (leading and trailing spaces and newline characters removed)
//
// Details:
//
//   - The function works by iteratively removing any characters less than or equal to a space (ASCII 32) from the
//     left side of the string until no such characters remain. It then performs the same operation on the right side of
//     the string until no whitespace characters are left.
//
//   - The function uses a `goto` mechanism to handle the removal in a loop, which ensures all leading and trailing
//     spaces (or any whitespace characters) are removed without additional checks for length or condition evaluation
//     in every iteration.
//
//   - The trimmed result string will not contain leading or trailing whitespace characters after the function completes.
//
//   - The function returns an unchanged string if no whitespace is present.
func trim(s string) string {
	if unify4g.IsEmpty(s) {
		return s
	}
left:
	if len(s) > 0 && s[0] <= ' ' {
		s = s[1:]
		goto left
	}
right:
	if len(s) > 0 && s[len(s)-1] <= ' ' {
		s = s[:len(s)-1]
		goto right
	}
	return s
}

// removeOuterBraces removes the surrounding '[]' or '{}' characters from a JSON string.
// This function is useful when you want to extract the content inside a JSON array or object,
// effectively unwrapping the outermost brackets or braces.
//
// Parameters:
//   - json: A string representing a JSON object or array. The string may include square brackets ('[]') or
//     curly braces ('{}') at the beginning and end, which will be removed if they exist.
//
// Returns:
//   - A new string with the outermost '[]' or '{}' characters removed. If the string does not start
//     and end with matching brackets or braces, the string remains unchanged.
//
// Example Usage:
//
//	json := "[1, 2, 3]"
//	unwrapped := removeOuterBraces(json)
//	// unwrapped: "1, 2, 3" (the array removed)
//
//	json = "{ \"name\": \"John\" }"
//	unwrapped = removeOuterBraces(json)
//	// unwrapped: " \"name\": \"John\" " (the object removed)
//
//	str := "hello world"
//	unwrapped = removeOuterBraces(str)
//	// unwrapped: "hello world" (no change since no surrounding brackets or braces)
//
// Details:
//
//   - The function first trims any leading or trailing whitespace from the input string using the `trim` function.
//
//   - It then checks if the string has at least two characters and if the first character is either '[' or '{'.
//
//   - If the first character is an opening bracket or brace, and the last character matches its pair (']' or '}'),
//     the function removes both the first and last characters.
//
//   - If the string does not start and end with matching brackets or braces, the original string is returned unchanged.
//
//   - The function handles cases where the string may contain additional whitespace at the beginning or end by trimming it first.
func removeOuterBraces(json string) string {
	json = trim(json)
	if len(json) >= 2 && (json[0] == '[' || json[0] == '{') {
		json = json[1 : len(json)-1]
	}
	return json
}

// stripNonWhitespace removes all non-whitespace characters from the input string, leaving only whitespace characters.
// The function iterates over each character in the input string and appends only whitespace characters (' ', '\t', '\n', '\r')
// to a new string. All non-whitespace characters are ignored and not included in the result.
//
// Parameters:
//   - s: A string that may contain a mixture of whitespace and non-whitespace characters.
//
// Returns:
//   - A new string consisting only of whitespace characters from the original string. If there are no whitespace characters
//     in the input string, it returns an empty string.
//
// Example Usage:
//
//	str := "  \t\n   abc  "
//	result := stripNonWhitespace(str)
//	// result: "     " (all non-whitespace characters are removed)
//
//	str = "hello"
//	result = stripNonWhitespace(str)
//	// result: "" (no whitespace characters, returns an empty string)
//
// Details:
//
//   - The function iterates through each character in the input string `s` and skips any non-whitespace character.
//
//   - It appends each whitespace character to a new byte slice `s2`, which is later converted to a string and returned.
//
//   - If the input string contains no whitespace characters, the function returns an empty string.
//
//   - This function may not be very efficient for long strings, as it performs an inner loop on each non-whitespace character.
func stripNonWhitespace(s string) string {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case ' ', '\t', '\n', '\r':
			continue
		default:
			var s2 []byte
			for i := 0; i < len(s); i++ {
				switch s[i] {
				case ' ', '\t', '\n', '\r':
					s2 = append(s2, s[i])
				}
			}
			return string(s2)
		}
	}
	return s
}

// unescapeJsonEncoded extracts a JSON-encoded string and returns both the full JSON string (with quotes) and the unescaped string content.
// The function processes the input string to handle escaped characters and returns a clean, unescaped version of the string
// as well as the portion of the JSON string that includes the enclosing quotes.
//
// Parameters:
//   - `json`: A JSON-encoded string, which is expected to start and end with a double quote (").
//     This string may contain escape sequences (e.g., `\"`, `\\`, `\n`, etc.) within the string value.
//
// Returns:
//   - `raw`: The full JSON string including the enclosing quotes (e.g., `"[Hello]"`).
//   - `unescaped`: The unescaped content of the string inside the quotes (e.g., `"Hello"` becomes `Hello` after unescaping).
//
// Example Usage:
//
//	input := "\"Hello\\nWorld\""
//	raw, str := unescapeJsonEncoded(input)
//	// raw: "\"Hello\\nWorld\"" (the full JSON string with quotes)
//	// str: "Hello\nWorld" (the unescaped string content)
//
//	input := "\"This is a \\\"quoted\\\" word\""
//	raw, str := unescapeJsonEncoded(input)
//	// raw: "\"This is a \\\"quoted\\\" word\""
//	// str: "This is a \"quoted\" word" (the unescaped string content)
//
// Details:
//
//   - The function processes the input string starting from the second character (ignoring the initial quote).
//
//   - It handles escape sequences inside the string, skipping over escaped quotes (`\"`) and other escape sequences.
//
//   - When a closing quote (`"`) is encountered, the function checks for the escape sequences to ensure the string is correctly unescaped.
//
//   - The function also checks if there are any escaped slashes (`\\`) and validates if they are part of an even or odd sequence.
//     If an escaped slash is found, it is taken into account to avoid terminating the string early.
//
//   - If the string is well-formed, the function returns the entire JSON string with quotes (`raw`) and the unescaped string (`str`).
//
//   - If an unescaped string is not found or the JSON string doesn't match expected formats, the function returns the string as is.
func unescapeJsonEncoded(json string) (raw string, unescaped string) {
	for i := 1; i < len(json); i++ {
		if json[i] > '\\' {
			continue
		}
		if json[i] == '"' {
			return json[:i+1], json[1:i]
		}
		if json[i] == '\\' {
			i++
			for ; i < len(json); i++ {
				if json[i] > '\\' {
					continue
				}
				if json[i] == '"' {
					// look for an escaped slash
					if json[i-1] == '\\' {
						n := 0
						for j := i - 2; j > 0; j-- {
							if json[j] != '\\' {
								break
							}
							n++
						}
						if n%2 == 0 {
							continue
						}
					}
					return json[:i+1], unescape(json[1:i])
				}
			}
			var ret string
			if i+1 < len(json) {
				ret = json[:i+1]
			} else {
				ret = json[:i]
			}
			return ret, unescape(json[1:i])
		}
	}
	return json, json[1:]
}

// isModifierOrJsonStart checks whether the first character of the input string `s` is a special character
// (such as '@', '[', or '{') that might indicate a modifier or a JSON structure in the context of processing.
//
// The function performs the following checks:
//   - If the first character is '@', it further inspects if the following characters indicate a modifier.
//   - If the first character is '[' or '{', it returns `true`, indicating a potential JSON array or object.
//   - The function will return `false` for any other characters or if modifiers are disabled.
//
// Parameters:
//   - `s`: A string to be checked, which can be a part of a JSON structure or an identifier with a modifier.
//
// Returns:
//   - `bool`: `true` if the first character is '@' followed by a modifier, or if the first character is '[' or '{'.
//     `false` otherwise.
//
// Example Usage:
//
//	s1 := "@modifier|value"
//	isModifierOrJsonStart(s1)
//	// Returns: true (because it starts with '@' and is followed by a modifier)
//
//	s2 := "[1, 2, 3]"
//	isModifierOrJsonStart(s2)
//	// Returns: true (because it starts with '[')
//
//	s3 := "{ \"key\": \"value\" }"
//	isModifierOrJsonStart(s3)
//	// Returns: true (because it starts with '{')
//
//	s4 := "normalString"
//	isModifierOrJsonStart(s4)
//	// Returns: false (no '@', '[', or '{')
//
// Details:
//   - The function first checks if modifiers are disabled (by `DisableModifiers` flag). If they are, it returns `false` immediately.
//   - If the string starts with '@', it scans for a potential modifier by checking if there is a '.' or '|' after it,
//     and verifies whether the modifier exists in the `modifiers` map.
//   - If the string starts with '[' or '{', it immediately returns `true`, as those characters typically indicate the start of a JSON array or object.
func isModifierOrJsonStart(s string) bool {
	if DisableModifiers {
		return false
	}
	c := s[0]
	if c == '@' {
		i := 1
		for ; i < len(s); i++ {
			if s[i] == '.' || s[i] == '|' || s[i] == ':' {
				break
			}
		}
		_, ok := modifiers[s[1:i]]
		return ok
	}
	return c == '[' || c == '{'
}

// matchSafely checks if a string matches a pattern with a complexity limit to
// avoid excessive computational cost, such as those from ReDos (Regular Expression Denial of Service) attacks.
//
// This function utilizes the `MatchLimit` function from `unify4g` to perform the matching, enforcing a maximum
// complexity limit of 10,000. The function aims to prevent situations where matching could lead to long or
// excessive computation, particularly when dealing with user-controlled input.
//
// Parameters:
//   - `str`: The string to match against the pattern.
//   - `pattern`: The pattern string to match, which may include wildcards or other special characters.
//
// Returns:
//   - `bool`: `true` if the `str` matches the `pattern` within the set complexity limit; otherwise `false`.
//
// Example:
//
//	result := matchSafely("hello", "h*o") // Returns `true` if the pattern matches the string within the complexity limit.
func matchSafely(str, pattern string) bool {
	matched, _ := unify4g.MatchLimit(str, pattern, 10000)
	return matched
}

// splitPathPipe splits a given path into two parts around the first unescaped '|' character.
// It also handles nested structures and ensures correct parsing even when special characters
// (e.g., braces, brackets, or quotes) are present.
//
// Parameters:
//   - `path`: A string representing the input path that may contain nested objects, arrays, or a pipe character.
//
// Returns:
//   - `left`: The part of the string before the first unescaped '|' character.
//   - `right`: The part of the string after the first unescaped '|' character.
//   - `ok`: A boolean indicating whether a valid split was found.
//
// Details:
//   - If the path contains a '|' that is part of a nested structure or escaped, the function skips over it.
//   - If the path starts with '{', the function uses the `squash` function to handle nested structures,
//     ensuring correct splitting of the path while preserving JSON-like formats.
//
// Notes:
//   - The function supports nested structures, including JSON-like objects and arrays (`{}`, `[]`), as well as
//     selector expressions (e.g., `#[...]` or `#(...)`).
//   - The function carefully skips escaped characters (e.g., `\|` or `\"`) and ensures that string literals
//     enclosed in quotes are handled properly without premature termination.
//   - It stops and splits the path at the first valid unescaped '|' encountered, returning the left and right parts.
//
// Example Usage:
//
//	For Input: `path1|path2`
//	   Returns: `left="path1"`, `right="path2"`, `ok=true`
//
//	For Input: `{nested|structure}|path2`
//	   Returns: `left="{nested|structure}"`, `right="path2"`, `ok=true`
//
//	For Input: `path_without_pipe`
//	   Returns: `left=""`, `right=""`, `ok=false`
func splitPathPipe(path string) (left, right string, ok bool) {
	var possible bool
	for i := 0; i < len(path); i++ {
		if path[i] == '|' {
			possible = true
			break
		}
	}
	if !possible {
		return
	}
	if len(path) > 0 && path[0] == '{' {
		squashed := squash(path[1:])
		if len(squashed) < len(path)-1 {
			squashed = path[:len(squashed)+1]
			remain := path[len(squashed):]
			if remain[0] == '|' {
				return squashed, remain[1:], true
			}
		}
		return
	}
	for i := 0; i < len(path); i++ {
		if path[i] == '\\' {
			i++
		} else if path[i] == '.' {
			if i == len(path)-1 {
				return
			}
			if path[i+1] == '#' {
				i += 2
				if i == len(path) {
					return
				}
				if path[i] == '[' || path[i] == '(' {
					var start, end byte
					if path[i] == '[' {
						start, end = '[', ']'
					} else {
						start, end = '(', ')'
					}
					// inside selector, balance brackets
					i++
					depth := 1
					for ; i < len(path); i++ {
						if path[i] == '\\' {
							i++
						} else if path[i] == start {
							depth++
						} else if path[i] == end {
							depth--
							if depth == 0 {
								break
							}
						} else if path[i] == '"' {
							// inside selector string, balance quotes
							i++
							for ; i < len(path); i++ {
								if path[i] == '\\' {
									i++
								} else if path[i] == '"' {
									break
								}
							}
						}
					}
				}
			}
		} else if path[i] == '|' {
			return path[:i], path[i+1:], true
		}
	}
	return
}

// parseString parses a string enclosed in double quotes from a JSON-encoded input string, handling escape sequences.
// It starts from a given index `i` and extracts the next JSON string, taking care of any escape sequences like `\"`, `\\`, etc.
//
// Parameters:
//   - `json`: A JSON string that may contain one or more strings enclosed in double quotes, possibly with escape sequences.
//   - `i`: The index in the `json` string to begin parsing from. The function expects this to point to the starting quote of the string.
//
// Returns:
//   - `i`: The index immediately following the closing quote of the string, or the point where parsing ends if no valid string is found.
//   - `raw`: The substring of `json` that includes the entire quoted string, including the surrounding quotes.
//   - `escaped`: A boolean flag indicating whether escape sequences were found and processed in the string.
//   - `valid`: A boolean flag indicating whether the string was correctly parsed and closed with a quote.
//
// Example Usage:
//
//	json := "\"Hello\\nWorld\""
//	i, raw, escaped, valid := parseString(json, 1)
//	// raw: "\"Hello\\nWorld\"" (the full quoted string)
//	// escaped: true (escape sequences processed)
//	// valid: true (valid string enclosed with quotes)
//
//	json = "\"NoEscapeHere\""
//	i, raw, escaped, valid = parseString(json, 1)
//	// raw: "\"NoEscapeHere\"" (the full quoted string)
//	// escaped: false (no escape sequences)
//	// valid: true (valid string enclosed with quotes)
//
//	json = "\"Hello\\\"Quoted\\\"String\""
//	i, raw, escaped, valid = parseString(json, 1)
//	// raw: "\"Hello\\\"Quoted\\\"String\""
//	// escaped: true
//	// valid: true
//
// Details:
//   - The function starts at the given index `i` and looks for the next double quote (`"`) to indicate the end of the string.
//   - It processes escape sequences inside the string, such as escaped quotes (`\"`) and backslashes (`\\`).
//   - If a valid string is found, the function returns the index after the closing quote, the full quoted string, and flags indicating if escape sequences were processed and if the string was properly closed.
//   - If the string is not correctly closed or contains invalid escape sequences, the function will stop processing and return the current state.
func parseString(json string, i int) (int, string, bool, bool) {
	if unify4g.IsEmpty(json) || i < 0 {
		return i, json, false, false
	}
	var s = i
	for ; i < len(json); i++ {
		if json[i] > '\\' {
			continue
		}
		if json[i] == '"' {
			return i + 1, json[s-1 : i+1], false, true
		}
		if json[i] == '\\' {
			i++
			for ; i < len(json); i++ {
				if json[i] > '\\' {
					continue
				}
				if json[i] == '"' {
					// look for an escaped slash
					if json[i-1] == '\\' {
						n := 0
						for j := i - 2; j > 0; j-- {
							if json[j] != '\\' {
								break
							}
							n++
						}
						if n%2 == 0 {
							continue
						}
					}
					return i + 1, json[s-1 : i+1], true, true
				}
			}
			break
		}
	}
	return i, json[s-1:], false, false
}

// parseNumeric parses a numeric value (integer or floating-point) from a JSON-encoded input string,
// starting from a given index `i` and extracting the numeric value up to a non-numeric character or JSON delimiter.
//
// Parameters:
//   - `json`: A JSON string that may contain numeric values (such as integers or floats).
//   - `i`: The index in the `json` string to begin parsing from. The function expects this to point to the first digit or part of the number.
//
// Returns:
//   - `i`: The index immediately following the last character of the parsed numeric value, or the point where parsing ends if no valid number is found.
//   - `raw`: The substring of `json` that represents the numeric value (e.g., "123", "45.67", "-123").
//
// Example Usage:
//
//	json := "12345"
//	i, raw := parseNumeric(json, 0)
//	// raw: "12345" (the parsed numeric value)
//	// i: the index after the last character of the number (6)
//
//	json = "-3.14159"
//	i, raw = parseNumeric(json, 0)
//	// raw: "-3.14159" (the parsed numeric value)
//	// i: the index after the last character of the number (8)
//
// Details:
//   - The function begins parsing from the given index `i` and continues until it encounters a character that is not part of a number (e.g., a space, comma, or closing brace/bracket).
//   - It handles both integers and floating-point numbers, as well as negative numbers (e.g., "-123" or "-45.67").
//   - The function stops parsing as soon as it encounters a non-numeric character such as whitespace, a comma, or a closing JSON delimiter (`}` or `]`), which indicates the end of the numeric value in the JSON structure.
//   - The function returns the parsed numeric string along with the index that follows the number's last character.
func parseNumeric(json string, i int) (int, string) {
	if unify4g.IsEmpty(json) || i < 0 {
		return i, json
	}
	var s = i
	i++
	for ; i < len(json); i++ {
		if json[i] <= ' ' || json[i] == ',' || json[i] == ']' ||
			json[i] == '}' {
			return i, json[s:i]
		}
	}
	return i, json[s:]
}

// parsePathWithModifiers parses a given path string, extracting different components such as parts, pipes, paths, and wildcards.
// It identifies special characters ('.', '|', '*', '?', '\\') in the path and processes them accordingly. The function
// breaks the string into a part and further splits it into pipes or paths, marking certain flags when necessary.
// It also handles escaped characters by stripping escape sequences and processing them correctly.
//
// Parameters:
//   - `path`: A string representing the path to be parsed. It can contain various special characters like
//     dots ('.'), pipes ('|'), wildcards ('*', '?'), and escape sequences ('\\').
//
// Returns:
//   - `r`: A `wildcard` struct containing the parsed components of the path. It will include:
//   - `Part`: The part of the path before any special character or wildcard.
//   - `Path`: The portion of the path after a dot ('.') if present.
//   - `Pipe`: The portion of the path after a pipe ('|') if present.
//   - `Piped`: A boolean flag indicating if a pipe ('|') was encountered.
//   - `Wild`: A boolean flag indicating if a wildcard ('*' or '?') was encountered.
//   - `More`: A boolean flag indicating if the path is further segmented by a dot ('.').
//
// Example Usage:
//
//	path1 := "field.subfield|anotherField"
//	result := parsePathWithModifiers(path1)
//	// result.Part: "field"
//	// result.Path: "subfield"
//	// result.Pipe: "anotherField"
//	// result.Piped: true
//
//	path2 := "object.field"
//	result = parsePathWithModifiers(path2)
//	// result.Part: "object"
//	// result.Path: "field"
//	// result.More: true
//
//	path3 := "path\\.*.field"
//	result = parsePathWithModifiers(path3)
//	// result.Part: "path.*"
//	// result.Wild: true
//
// Details:
//   - The function scans through the path string character by character, processing the first encountered special
//     character (either '.', '|', '*', '?', '\\') and extracting the relevant components.
//   - If a '.' is encountered, the part before it is extracted as the `Part`, and the string following it is assigned
//     to `Path`. If there are modifiers or JSON structure indicators (like '[' or '{'), the path is marked accordingly.
//   - If a pipe ('|') is found, the `Part` is separated from the string after the pipe, and the `Piped` flag is set to true.
//   - Wildcard characters ('*' or '?') are detected, and the `Wild` flag is set.
//   - Escape sequences (indicated by '\\') are processed by appending the escaped character(s) and stripping the escape character.
//   - If no special characters are found, the entire path is assigned to `Part`, and the function returns the parsed result.
func parsePathWithModifiers(path string) (r wildcard) {
	for i := 0; i < len(path); i++ {
		if path[i] == '|' {
			r.Part = path[:i]
			r.Pipe = path[i+1:]
			r.Piped = true
			return
		}
		if path[i] == '.' {
			r.Part = path[:i]
			if i < len(path)-1 && isModifierOrJsonStart(path[i+1:]) {
				r.Pipe = path[i+1:]
				r.Piped = true
			} else {
				r.Path = path[i+1:]
				r.More = true
			}
			return
		}
		if path[i] == '*' || path[i] == '?' {
			r.Wild = true
			continue
		}
		if path[i] == '\\' {
			// go into escape mode
			// a slower path that strips off the escape character from the part.
			escapePart := []byte(path[:i])
			i++
			if i < len(path) {
				escapePart = append(escapePart, path[i])
				i++
				for ; i < len(path); i++ {
					if path[i] == '\\' {
						i++
						if i < len(path) {
							escapePart = append(escapePart, path[i])
						}
						continue
					} else if path[i] == '.' {
						r.Part = string(escapePart)
						if i < len(path)-1 && isModifierOrJsonStart(path[i+1:]) {
							r.Pipe = path[i+1:]
							r.Piped = true
						} else {
							r.Path = path[i+1:]
							r.More = true
						}
						return
					} else if path[i] == '|' {
						r.Part = string(escapePart)
						r.Pipe = path[i+1:]
						r.Piped = true
						return
					} else if path[i] == '*' || path[i] == '?' {
						r.Wild = true
					}
					escapePart = append(escapePart, path[i])
				}
			}
			r.Part = string(escapePart)
			return
		}
	}
	r.Part = path
	return
}

// parseJSONLiteral parses a literal value (e.g., "true", "false", or "null") from a JSON-encoded input string,
// starting from a given index `i` and extracting the literal value up to a non-alphabetic character or JSON delimiter.
//
// Parameters:
//   - `json`: A JSON string that may contain literal values (such as "true", "false", or "null").
//   - `i`: The index in the `json` string to begin parsing from. The function expects this to point to the first character of the literal.
//
// Returns:
//   - `i`: The index immediately following the last character of the parsed literal value, or the point where parsing ends if no valid literal is found.
//   - `raw`: The substring of `json` that represents the literal value (e.g., "true", "false", "null").
//
// Example Usage:
//
//	json := "true"
//	i, raw := parseJSONLiteral(json, 0)
//	// raw: "true" (the parsed literal value)
//	// i: the index after the last character of the literal (4)
//
//	json = "false"
//	i, raw = parseJSONLiteral(json, 0)
//	// raw: "false" (the parsed literal value)
//	// i: the index after the last character of the literal (5)
//
//	json = "null"
//	i, raw = parseJSONLiteral(json, 0)
//	// raw: "null" (the parsed literal value)
//	// i: the index after the last character of the literal (4)
//
// Details:
//   - The function begins parsing from the given index `i` and continues until it encounters a character that is not part of a valid literal (i.e., characters outside the range of 'a' to 'z').
//   - It handles JSON literal values like "true", "false", and "null" by checking for consecutive alphabetic characters.
//   - The function stops parsing as soon as it encounters a non-alphabetic character such as whitespace, a comma, or a closing JSON delimiter (`}` or `]`), which indicates the end of the literal value in the JSON structure.
//   - The function returns the parsed literal string along with the index that follows the literal's last character.
func parseJSONLiteral(json string, i int) (int, string) {
	if unify4g.IsEmpty(json) || i < 0 {
		return i, json
	}
	var s = i
	i++
	for ; i < len(json); i++ {
		if json[i] < 'a' || json[i] > 'z' {
			return i, json[s:i]
		}
	}
	return i, json[s:]
}

// parseJSONSquash processes a JSON string starting from a given index `i`, squashing (flattening) any nested JSON structures
// (such as arrays, objects, or even parentheses) into a single value. The function handles strings, nested objects,
// arrays, and parentheses while ignoring the nested structures themselves, only returning the top-level JSON structure
// from the starting point.
//
// Parameters:
//   - `json`: A string representing the JSON data to be parsed. This string can include various JSON constructs like
//     strings, objects, arrays, and nested structures within parentheses.
//   - `i`: The index in the `json` string from which parsing should begin. The function assumes that the character
//     at this index is the opening character of a JSON array ('['), object ('{'), or parentheses ('(').
//
// Returns:
//   - `int`: The new index after the parsing of the JSON structure, which is after the closing bracket/parenthesis/brace.
//   - `string`: A string containing the flattened JSON structure starting from the opening character and squashing all
//     nested structures until the corresponding closing character is reached.
//
// Example Usage:
//
//	json := "[{ \"key\": \"value\" }, { \"nested\": [1, 2, 3] }]"
//	i, result := parseJSONSquash(json, 0)
//	// result: "{ \"key\": \"value\" }, { \"nested\": [1, 2, 3] }" (flattened to top-level content)
//	// i: the index after the closing ']' of the outer array
//
// Details:
//   - The function expects that the character at index `i` is an opening character for an array, object, or parentheses,
//     and it will proceed to skip over any nested structures of the same type (i.e., arrays, objects, or parentheses).
//   - The depth of nesting is tracked, and whenever the function encounters a closing bracket (']'), brace ('}'), or parenthesis
//     (')'), it checks if the depth has returned to 0 (indicating the end of the top-level structure).
//   - If a string is encountered (enclosed in double quotes), it processes the string contents carefully, respecting escape sequences.
//   - The function ensures that nested structures (arrays, objects, or parentheses) are ignored, effectively "squashing" the
//     content into the outermost structure, while the depth ensures that only the highest-level structure is returned.
func parseJSONSquash(json string, i int) (int, string) {
	if unify4g.IsEmpty(json) || i < 0 {
		return i, json
	}
	s := i
	i++
	depth := 1
	for ; i < len(json); i++ {
		if json[i] >= '"' && json[i] <= '}' {
			switch json[i] {
			case '"':
				i++
				s2 := i
				for ; i < len(json); i++ {
					if json[i] > '\\' {
						continue
					}
					if json[i] == '"' {
						// look for an escaped slash
						if json[i-1] == '\\' {
							n := 0
							for j := i - 2; j > s2-1; j-- {
								if json[j] != '\\' {
									break
								}
								n++
							}
							if n%2 == 0 {
								continue
							}
						}
						break
					}
				}
			case '{', '[', '(':
				depth++
			case '}', ']', ')':
				depth--
				if depth == 0 {
					i++
					return i, json[s:i]
				}
			}
		}
	}
	return i, json[s:]
}

// parseJSONAny parses the next JSON value from a given JSON string starting at the specified index `i`.
// The function identifies and processes a variety of JSON value types including objects, arrays, strings, literals (true, false, null),
// and numeric values. The result of parsing is returned as a `Context` containing relevant information about the parsed value.
//
// Parameters:
//   - `json`: A string representing the JSON data to be parsed. This string can include objects, arrays, strings, literals, and numbers.
//   - `i`: The starting index in the `json` string where parsing should begin. The function will parse the value starting at this index.
//   - `hit`: A boolean flag indicating whether to capture the parsed result into the `Context` object. If true, the context will be populated with the parsed value.
//
// Returns:
//   - `i`: The updated index after parsing the JSON value. This is the index immediately after the parsed value.
//   - `ctx`: A `Context` object containing information about the parsed value, including the type (`kind`), the raw unprocessed string (`unprocessed`),
//     and for strings or numbers, the parsed value (e.g., `strings` for strings, `numeric` for numbers).
//   - `ok`: A boolean indicating whether the parsing was successful. If the function successfully identifies a JSON value, it returns true; otherwise, false.
//
// Example Usage:
//
//	json := `{"key": "value", "age": 25}`
//	i := 0
//	hit := true
//	i, ctx, ok := parseJSONAny(json, i, hit)
//	// i: the index after parsing the first JSON value (e.g., after the closing quote of "value").
//	// ctx: contains the parsed context information (e.g., for strings, the kind would be String, unprocessed would be the raw value, etc.)
//	// ok: true if the value was successfully parsed.
//
// Details:
//   - The function processes various JSON types, including objects, arrays, strings, literals (true, false, null), and numeric values.
//   - It recognizes objects (`{}`), arrays (`[]`), and string literals (`""`), and calls the appropriate helper functions for each type.
//   - When parsing a string, it handles escape sequences, and when parsing numeric values, it checks for valid numbers (including integers, floats, and special numeric literals like `NaN`).
//   - For literals like `true`, `false`, and `null`, the function parses the exact keywords and stores them in the `Context` object as `True`, `False`, or `Null` respectively.
//
// The function ensures flexibility by checking each character in the JSON string and delegating to specialized functions for handling different value types.
// If no valid JSON value is found at the given position, it returns false.
func parseJSONAny(json string, i int, hit bool) (int, Context, bool) {
	var ctx Context
	var val string
	for ; i < len(json); i++ {
		if json[i] == '{' || json[i] == '[' {
			i, val = parseJSONSquash(json, i)
			if hit {
				ctx.unprocessed = val
				ctx.kind = JSON
			}
			var tmp parser
			tmp.value = ctx
			calcSubstring(json, &tmp)
			return i, tmp.value, true
		}
		if json[i] <= ' ' {
			continue
		}
		var num bool
		switch json[i] {
		case '"':
			i++
			var escVal bool
			var ok bool
			i, val, escVal, ok = parseString(json, i)
			if !ok {
				return i, ctx, false
			}
			if hit {
				ctx.kind = String
				ctx.unprocessed = val
				if escVal {
					ctx.strings = unescape(val[1 : len(val)-1])
				} else {
					ctx.strings = val[1 : len(val)-1]
				}
			}
			return i, ctx, true
		case 'n':
			if i+1 < len(json) && json[i+1] != 'u' {
				num = true
				break
			}
			fallthrough
		case 't', 'f':
			vc := json[i]
			i, val = parseJSONLiteral(json, i)
			if hit {
				ctx.unprocessed = val
				switch vc {
				case 't':
					ctx.kind = True
				case 'f':
					ctx.kind = False
				}
				return i, ctx, true
			}
		case '+', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'i', 'I', 'N':
			num = true
		}
		if num {
			i, val = parseNumeric(json, i)
			if hit {
				ctx.unprocessed = val
				ctx.kind = Number
				ctx.numeric, _ = strconv.ParseFloat(val, 64)
			}
			return i, ctx, true
		}

	}
	return i, ctx, false
}

// parseJSONObject parses a JSON object structure from a given JSON string, extracting key-value pairs based on a specified path.
//
// The function processes a JSON object (denoted by curly braces '{' and '}') and looks for matching keys. It handles both
// simple key-value pairs and nested structures (objects or arrays) within the object. If the path to a key contains wildcards
// or modifiers, the function matches the keys accordingly. It also processes escape sequences for both keys and values,
// ensuring proper handling of special characters within JSON strings.
//
// Parameters:
//   - `c`: A pointer to a `parser` object that holds the JSON string (`json`), and context information (`value`).
//   - `i`: The current index in the JSON string from where the parsing should begin. This index should point to the
//     opening curly brace '{' of the JSON object.
//   - `path`: The string representing the path to be parsed. It may include modifiers or wildcards, guiding the matching
//     of specific keys in the object.
//
// Returns:
//   - `i` (int): The index in the JSON string after the parsing is completed. This index points to the character
//     immediately after the parsed object.
//   - `bool`: `true` if a match for the specified path was found, `false` if no match was found.
//
// Example Usage:
//
//	json := `{"name": "John", "age": 30, "address": {"city": "New York"}}`
//	i, found := parseJSONObject(c, 0, "name")
//	// found: true (if the "name" key was found in the JSON object)
//
// Details:
//   - The function first searches for a key enclosed in double quotes ('"'). It handles both normal keys and escaped keys.
//   - It then checks if the key matches the specified path, which may contain wildcards or exact matches.
//   - If the key matches and there are no more modifiers in the path, the corresponding value is extracted and stored in the `parser` object.
//   - If the key points to a nested object or array, the function recursively parses those structures to extract the required data.
//   - The function handles various types of JSON values including strings, numbers, booleans, objects, and arrays.
//   - The function also handles escape sequences within JSON strings and ensures that they are processed correctly.
//
// Notes:
//   - The function makes use of the `parsePathWithModifiers` function to parse and process the path for matching keys.
//   - If the path contains wildcards ('*' or '?'), the function uses `matchSafely` to ensure safe matching within a complexity limit.
//   - If the key is matched, the function will return the parsed value. If no match is found, the parsing continues.
//
// Key functions used:
//   - `parsePathWithModifiers`: Extracts and processes the path to identify the key and modifiers.
//   - `matchSafely`: Performs the safe matching of the key using a wildcard pattern, avoiding excessive complexity.
func parseJSONObject(c *parser, i int, path string) (int, bool) {
	var _match, keyEsc, escVal, ok, hit bool
	var key, val string
	pathModifiers := parsePathWithModifiers(path)
	if !pathModifiers.More && pathModifiers.Piped {
		c.pipe = pathModifiers.Pipe
		c.piped = true
	}
	for i < len(c.json) {
		for ; i < len(c.json); i++ {
			if c.json[i] == '"' {
				i++
				var s = i
				for ; i < len(c.json); i++ {
					if c.json[i] > '\\' {
						continue
					}
					if c.json[i] == '"' {
						i, key, keyEsc, ok = i+1, c.json[s:i], false, true
						goto parse_key_completed
					}
					if c.json[i] == '\\' {
						i++
						for ; i < len(c.json); i++ {
							if c.json[i] > '\\' {
								continue
							}
							if c.json[i] == '"' {
								// look for an escaped slash
								if c.json[i-1] == '\\' {
									n := 0
									for j := i - 2; j > 0; j-- {
										if c.json[j] != '\\' {
											break
										}
										n++
									}
									if n%2 == 0 {
										continue
									}
								}
								i, key, keyEsc, ok = i+1, c.json[s:i], true, true
								goto parse_key_completed
							}
						}
						break
					}
				}
				key, keyEsc, ok = c.json[s:], false, false
			parse_key_completed:
				break
			}
			if c.json[i] == '}' {
				return i + 1, false
			}
		}
		if !ok {
			return i, false
		}
		if pathModifiers.Wild {
			if keyEsc {
				_match = matchSafely(unescape(key), pathModifiers.Part)
			} else {
				_match = matchSafely(key, pathModifiers.Part)
			}
		} else {
			if keyEsc {
				_match = pathModifiers.Part == unescape(key)
			} else {
				_match = pathModifiers.Part == key
			}
		}
		hit = _match && !pathModifiers.More
		for ; i < len(c.json); i++ {
			var num bool
			switch c.json[i] {
			default:
				continue
			case '"':
				i++
				i, val, escVal, ok = parseString(c.json, i)
				if !ok {
					return i, false
				}
				if hit {
					if escVal {
						c.value.strings = unescape(val[1 : len(val)-1])
					} else {
						c.value.strings = val[1 : len(val)-1]
					}
					c.value.unprocessed = val
					c.value.kind = String
					return i, true
				}
			case '{':
				if _match && !hit {
					i, hit = parseJSONObject(c, i+1, pathModifiers.Path)
					if hit {
						return i, true
					}
				} else {
					i, val = parseJSONSquash(c.json, i)
					if hit {
						c.value.unprocessed = val
						c.value.kind = JSON
						return i, true
					}
				}
			case '[':
				if _match && !hit {
					i, hit = analyzeArray(c, i+1, pathModifiers.Path)
					if hit {
						return i, true
					}
				} else {
					i, val = parseJSONSquash(c.json, i)
					if hit {
						c.value.unprocessed = val
						c.value.kind = JSON
						return i, true
					}
				}
			case 'n':
				if i+1 < len(c.json) && c.json[i+1] != 'u' {
					num = true
					break
				}
				fallthrough
			case 't', 'f':
				vc := c.json[i]
				i, val = parseJSONLiteral(c.json, i)
				if hit {
					c.value.unprocessed = val
					switch vc {
					case 't':
						c.value.kind = True
					case 'f':
						c.value.kind = False
					}
					return i, true
				}
			case '+', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
				'i', 'I', 'N':
				num = true
			}
			if num {
				i, val = parseNumeric(c.json, i)
				if hit {
					c.value.unprocessed = val
					c.value.kind = Number
					c.value.numeric, _ = strconv.ParseFloat(val, 64)
					return i, true
				}
			}
			break
		}
	}
	return i, false
}

// analyzeQuery parses a query string into its constituent parts and identifies its structure.
// It is designed to handle queries that involve filtering or accessing nested structures,
// particularly in JSON-like or similar data representations.
//
// Parameters:
//   - query (string): The query string to parse. It must start with `#(` or `#[` to be valid.
//     The query string may include paths, operators, and values, and can contain nested structures.
//
// Returns:
//   - path (string): The portion of the query string representing the path to the field or property.
//   - op (string): The operator used in the query (e.g., `==`, `!=`, `<`, `>`, etc.).
//   - value (string): The value to compare against or use in the query.
//   - remain (string): The remaining portion of the query after processing.
//   - i (int): The index in the query string where parsing ended.
//   - _vEsc (bool): Indicates whether the value part contains any escaped characters.
//   - ok (bool): Indicates whether the parsing was successful or not.
//
// Example Usage:
//
//	For a query `#(first_name=="Aris").last`:
//	  - path: "first_name"
//	  - op: "=="
//	  - value: "Aris"
//	  - remain: ".last"
//
//	For a query `#(user_roles.#(=="admin")).privilege`:
//	  - path: "user_roles.#(=="admin")"
//	  - op: ""
//	  - value: ""
//	  - remain: ".privilege"
//
// Details:
//   - The function starts by verifying the query's validity, ensuring it begins with `#(` or `#[`.
//   - It processes the query character by character, accounting for nested structures, operators, and escaped characters.
//   - The `path` is extracted from the portion of the query before the operator or value.
//   - The `op` and `value` are identified and split if an operator is present.
//   - Remaining characters in the query, such as `.last` in the example, are captured in `remain`.
//
// Notes:
//   - The function supports a variety of operators (`==`, `!=`, `<`, `>`, etc.).
//   - It handles nested brackets or parentheses and ensures balanced nesting.
//   - Escaped characters (e.g., `\"`) within the query are processed correctly, with `_vEsc` indicating their presence.
//   - If the query is invalid or incomplete, the function will return `ok` as `false`.
//
// Edge Cases:
//   - Handles nested queries with multiple levels of depth.
//   - Ensures proper handling of invalid or malformed queries by returning appropriate values.
func analyzeQuery(query string) (
	path, op, value, remain string, i int, _vEsc, ok bool,
) {
	if len(query) < 2 || query[0] != '#' ||
		(query[1] != '(' && query[1] != '[') {
		return "", "", "", "", i, false, false
	}
	i = 2
	j := 0
	depth := 1
	for ; i < len(query); i++ {
		if depth == 1 && j == 0 {
			switch query[i] {
			case '!', '=', '<', '>', '%':
				j = i
				continue
			}
		}
		if query[i] == '\\' {
			i++
		} else if query[i] == '[' || query[i] == '(' {
			depth++
		} else if query[i] == ']' || query[i] == ')' {
			depth--
			if depth == 0 {
				break
			}
		} else if query[i] == '"' {
			i++
			for ; i < len(query); i++ {
				if query[i] == '\\' {
					_vEsc = true
					i++
				} else if query[i] == '"' {
					break
				}
			}
		}
	}
	if depth > 0 {
		return "", "", "", "", i, false, false
	}
	if j > 0 {
		path = trim(query[2:j])
		value = trim(query[j:i])
		remain = query[i+1:]
		var trail int
		switch {
		case len(value) == 1:
			trail = 1
		case value[0] == '!' && value[1] == '=':
			trail = 2
		case value[0] == '!' && value[1] == '%':
			trail = 2
		case value[0] == '<' && value[1] == '=':
			trail = 2
		case value[0] == '>' && value[1] == '=':
			trail = 2
		case value[0] == '=' && value[1] == '=':
			value = value[1:]
			trail = 1
		case value[0] == '<':
			trail = 1
		case value[0] == '>':
			trail = 1
		case value[0] == '=':
			trail = 1
		case value[0] == '%':
			trail = 1
		}
		op = value[:trail]
		value = trim(value[trail:])
	} else {
		path = trim(query[2:i])
		remain = query[i+1:]
	}
	return path, op, value, remain, i + 1, _vEsc, true
}

// analyzePath parses a string path into its structural components, breaking it into meaningful parts
// such as the main path, pipe (if present), query parameters, and nested paths. This function is particularly
// useful for processing JSON-like paths or other hierarchical data representations.
//
// Parameters:
//   - path (string): The input string path to be analyzed. It may contain various symbols such as '|', '.',
//     or '#' that represent different parts or behaviors.
//
// Returns:
//   - r (deeper): A struct containing the parsed components of the path, such as `Part`, `Pipe`, `Path`,
//     and additional metadata (e.g., `Piped`, `More`, `Arch`, and query-related fields).
//
// Fields in `deeper`:
//   - `Part`: The main part of the path before special characters like '.', '|', or '#'.
//   - `Path`: The remaining part of the path after '.' or other separators.
//   - `Pipe`: A piped portion of the path, if separated by '|', indicating a subsequent operation.
//   - `Piped`: A boolean indicating whether the path contains a pipe ('|').
//   - `More`: A boolean indicating whether there is more of the path to process after the first separator.
//   - `Arch`: A boolean indicating the presence of a '#' in the path, signifying an archive or query operation.
//   - `ALogOk`: A boolean indicating a valid archive log if the path starts with `#.`.
//   - `ALogKey`: The key following `#.` for archive logging, if applicable.
//   - `query`: A nested struct providing details about a query if the path contains query operations:
//   - `On`: Indicates whether the path contains a query (e.g., starting with `#(`).
//   - `All`: Indicates whether the query applies to all elements.
//   - `QueryPath`: The path portion of the query.
//   - `Option`: The operator used in the query (e.g., `==`, `!=`, etc.).
//   - `Value`: The value used in the query.
//   - `Option`: The operator for comparison.
//   - `Value`: The query value.
//
// Details:
//   - The function iterates through the `path` string character by character, identifying and processing special symbols
//     such as '|', '.', and '#'.
//   - If the path contains a '|', the portion before it is stored in `Part`, and the portion after it is stored in `Pipe`.
//     The `Piped` flag is set to `true`.
//   - If the path contains a '.', the portion before it is stored in `Part`, and the remaining part in `Path`.
//     If the path after the '.' starts with a modifier or JSON, it is stored in `Pipe` instead, with `Piped` set to `true`.
//   - If the path contains a '#', the `Arch` flag is set to `true`. It may also indicate an archive log (`#.key`) or a query (`#(...)`).
//     Queries are parsed using the `analyzeQuery` function, and relevant fields in the `query` struct are populated.
//   - For archive logs starting with `#.` (e.g., `#.key`), the `ALogOk` flag is set, and `ALogKey` contains the key.
//   - If the path contains a query, the function extracts and processes the query's path, operator, and value.
//     Queries are denoted by a '#' followed by '[' or '(' (e.g., `#[...]` or `#(...)`).
//
// Example Usage:
//
//	For Input: "data|filter.name"
//	   Part: "data"
//	   Pipe: "filter.name"
//	   Piped: true
//	   Path: ""
//	   More: false
//	   Arch: false
//
//	For Input: "items.#(value=='42').details"
//	   Part: "items"
//	   Path: "#(value=='42').details"
//	   Arch: true
//	   query.On: true
//	   query.QueryPath: "value"
//	   query.Option: "=="
//	   query.Value: "42"
//	   query.All: false
//
//	For Input: "#.log"
//	   Part: "#"
//	   Path: ""
//	   ALogOk: true
//	   ALogKey: "log"
//
// Notes:
//   - The function is robust against malformed paths but assumes valid inputs for proper operation.
//   - It ensures nested paths and queries are correctly identified and processed.
//
// Edge Cases:
//   - If no special characters are found, the entire input is stored in `Part`.
//   - If the path contains an incomplete or invalid query, the function skips the query parsing gracefully.
func analyzePath(path string) (r deeper) {
	for i := 0; i < len(path); i++ {
		if path[i] == '|' {
			r.Part = path[:i]
			r.Pipe = path[i+1:]
			r.Piped = true
			return
		}
		if path[i] == '.' {
			r.Part = path[:i]
			if !r.Arch && i < len(path)-1 && isModifierOrJsonStart(path[i+1:]) {
				r.Pipe = path[i+1:]
				r.Piped = true
			} else {
				r.Path = path[i+1:]
				r.More = true
			}
			return
		}
		if path[i] == '#' {
			r.Arch = true
			if i == 0 && len(path) > 1 {
				if path[1] == '.' {
					r.ALogOk = true
					r.ALogKey = path[2:]
					r.Path = path[:1]
				} else if path[1] == '[' || path[1] == '(' {
					// query
					r.query.On = true
					queryPath, op, value, _, fi, escVal, ok :=
						analyzeQuery(path[i:])
					if !ok {
						break
					}
					if len(value) >= 2 && value[0] == '"' &&
						value[len(value)-1] == '"' {
						value = value[1 : len(value)-1]
						if escVal {
							value = unescape(value)
						}
					}
					r.query.QueryPath = queryPath
					r.query.Option = op
					r.query.Value = value

					i = fi - 1
					if i+1 < len(path) && path[i+1] == '#' {
						r.query.All = true
					}
				}
			}
			continue
		}
	}
	r.Part = path
	r.Path = ""
	return
}

// analyzeArray processes and evaluates the path in the context of an array, checking
// for matches and executing queries on elements within the array. It is responsible
// for handling nested structures and queries, as well as determining if the analysis
// matches the current path and value in the context.
//
// Parameters:
//   - c (parser*): A pointer to the `parser` object that holds the current context,
//     including the JSON data and the parsed path.
//   - i (int): The current index in the JSON string being processed.
//   - path (string): The path to be analyzed for array processing.
//
// Returns:
//   - (int): The updated index after processing the array.
//   - (bool): A boolean indicating whether the analysis was successful or not.
//
// Details:
//   - The function analyzes a path related to arrays and performs various checks to
//     determine if the array elements match the specified path and conditions.
//   - It checks for array literals, objects, and nested structures, invoking appropriate
//     parsing functions for each type.
//   - If a query is present, it will evaluate the query on the current element and decide
//     whether to continue the search or return a match.
//   - The function supports queries on array elements (e.g., matching specific values),
//     and it can return results in JSON format or execute specific actions (like calculating a value).
//
// Flow:
//   - The function first processes the path and ensures that it is valid for array analysis.
//   - It checks if the path includes an archive log, and if so, handles logging operations.
//   - The core loop processes each element of the array, checking for string, numeric, object,
//     or array elements and evaluating whether they match the query conditions, if any.
//   - If the query is satisfied, the function performs further processing on the matching element,
//     such as storing the result or calculating a value. If no query is provided, it directly
//     sets the `c.value` with the matched result.
//   - It also handles special cases like archive logs and nested array structures.
//   - If no valid match is found, the function returns `false`, and the search continues.
//
// Example:
//
//	Input: `["apple", "banana", "cherry"]`
//	If the query was for "banana", the function would find a match and return the result.
//
// Edge Cases:
//   - Handles situations where no array is found or the query fails to match any element.
//   - Properly handles nested arrays or objects within the JSON data, maintaining structure.
//   - Takes into account escaped characters and special syntax (e.g., queries, JSON objects).
func analyzeArray(c *parser, i int, path string) (int, bool) {
	var _match, escVal, ok, hit bool
	var val string
	var h int
	var aLog []int
	var partIdx int
	var multics []byte
	var queryIndexes []int
	analysis := analyzePath(path)
	if !analysis.Arch {
		n, ok := parseUint(analysis.Part)
		if !ok {
			partIdx = -1
		} else {
			partIdx = int(n)
		}
	}
	if !analysis.More && analysis.Piped {
		c.pipe = analysis.Pipe
		c.piped = true
	}

	executeQuery := func(eVal Context) bool {
		if analysis.query.All {
			if len(multics) == 0 {
				multics = append(multics, '[')
			}
		}
		var tmp parser
		tmp.value = eVal
		calcSubstring(c.json, &tmp)
		parentIndex := tmp.value.index
		var res Context
		if eVal.kind == JSON {
			res = eVal.Get(analysis.query.QueryPath)
		} else {
			if analysis.query.QueryPath != "" {
				return false
			}
			res = eVal
		}
		if queryMatches(&analysis, res) {
			if analysis.More {
				left, right, ok := splitPathPipe(analysis.Path)
				if ok {
					analysis.Path = left
					c.pipe = right
					c.piped = true
				}
				res = eVal.Get(analysis.Path)
			} else {
				res = eVal
			}
			if analysis.query.All {
				raw := res.unprocessed
				if len(raw) == 0 {
					raw = res.String()
				}
				if raw != "" {
					if len(multics) > 1 {
						multics = append(multics, ',')
					}
					multics = append(multics, raw...)
					queryIndexes = append(queryIndexes, res.index+parentIndex)
				}
			} else {
				c.value = res
				return true
			}
		}
		return false
	}
	for i < len(c.json)+1 {
		if !analysis.Arch {
			_match = partIdx == h
			hit = _match && !analysis.More
		}
		h++
		if analysis.ALogOk {
			aLog = append(aLog, i)
		}
		for ; ; i++ {
			var ch byte
			if i > len(c.json) {
				break
			} else if i == len(c.json) {
				ch = ']'
			} else {
				ch = c.json[i]
			}
			var num bool
			switch ch {
			default:
				continue
			case '"':
				i++
				i, val, escVal, ok = parseString(c.json, i)
				if !ok {
					return i, false
				}
				if analysis.query.On {
					var cVal Context
					if escVal {
						cVal.strings = unescape(val[1 : len(val)-1])
					} else {
						cVal.strings = val[1 : len(val)-1]
					}
					cVal.unprocessed = val
					cVal.kind = String
					if executeQuery(cVal) {
						return i, true
					}
				} else if hit {
					if analysis.ALogOk {
						break
					}
					if escVal {
						c.value.strings = unescape(val[1 : len(val)-1])
					} else {
						c.value.strings = val[1 : len(val)-1]
					}
					c.value.unprocessed = val
					c.value.kind = String
					return i, true
				}
			case '{':
				if _match && !hit {
					i, hit = parseJSONObject(c, i+1, analysis.Path)
					if hit {
						if analysis.ALogOk {
							break
						}
						return i, true
					}
				} else {
					i, val = parseJSONSquash(c.json, i)
					if analysis.query.On {
						if executeQuery(Context{unprocessed: val, kind: JSON}) {
							return i, true
						}
					} else if hit {
						if analysis.ALogOk {
							break
						}
						c.value.unprocessed = val
						c.value.kind = JSON
						return i, true
					}
				}
			case '[':
				if _match && !hit {
					i, hit = analyzeArray(c, i+1, analysis.Path)
					if hit {
						if analysis.ALogOk {
							break
						}
						return i, true
					}
				} else {
					i, val = parseJSONSquash(c.json, i)
					if analysis.query.On {
						if executeQuery(Context{unprocessed: val, kind: JSON}) {
							return i, true
						}
					} else if hit {
						if analysis.ALogOk {
							break
						}
						c.value.unprocessed = val
						c.value.kind = JSON
						return i, true
					}
				}
			case 'n':
				if i+1 < len(c.json) && c.json[i+1] != 'u' {
					num = true
					break
				}
				fallthrough
			case 't', 'f':
				vc := c.json[i]
				i, val = parseJSONLiteral(c.json, i)
				if analysis.query.On {
					var cVal Context
					cVal.unprocessed = val
					switch vc {
					case 't':
						cVal.kind = True
					case 'f':
						cVal.kind = False
					}
					if executeQuery(cVal) {
						return i, true
					}
				} else if hit {
					if analysis.ALogOk {
						break
					}
					c.value.unprocessed = val
					switch vc {
					case 't':
						c.value.kind = True
					case 'f':
						c.value.kind = False
					}
					return i, true
				}
			case '+', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
				'i', 'I', 'N':
				num = true
			case ']':
				if analysis.Arch && analysis.Part == "#" {
					if analysis.ALogOk {
						left, right, ok := splitPathPipe(analysis.ALogKey)
						if ok {
							analysis.ALogKey = left
							c.pipe = right
							c.piped = true
						}
						var indexes = make([]int, 0, 64)
						var jsonVal = make([]byte, 0, 64)
						jsonVal = append(jsonVal, '[')
						for j, k := 0, 0; j < len(aLog); j++ {
							idx := aLog[j]
							for idx < len(c.json) {
								switch c.json[idx] {
								case ' ', '\t', '\r', '\n':
									idx++
									continue
								}
								break
							}
							if idx < len(c.json) && c.json[idx] != ']' {
								_, res, ok := parseJSONAny(c.json, idx, true)
								if ok {
									res := res.Get(analysis.ALogKey)
									if res.Exists() {
										if k > 0 {
											jsonVal = append(jsonVal, ',')
										}
										raw := res.unprocessed
										if len(raw) == 0 {
											raw = res.String()
										}
										jsonVal = append(jsonVal, []byte(raw)...)
										indexes = append(indexes, res.index)
										k++
									}
								}
							}
						}
						jsonVal = append(jsonVal, ']')
						c.value.kind = JSON
						c.value.unprocessed = string(jsonVal)
						c.value.indexes = indexes
						return i + 1, true
					}
					if analysis.ALogOk {
						break
					}

					c.value.kind = Number
					c.value.numeric = float64(h - 1)
					c.value.unprocessed = strconv.Itoa(h - 1)
					c.calc = true
					return i + 1, true
				}
				if !c.value.Exists() {
					if len(multics) > 0 {
						c.value = Context{
							unprocessed: string(append(multics, ']')),
							kind:        JSON,
							indexes:     queryIndexes,
						}
					} else if analysis.query.All {
						c.value = Context{
							unprocessed: "[]",
							kind:        JSON,
						}
					}
				}
				return i + 1, false
			}
			if num {
				i, val = parseNumeric(c.json, i)
				if analysis.query.On {
					var cVal Context
					cVal.unprocessed = val
					cVal.kind = Number
					cVal.numeric, _ = strconv.ParseFloat(val, 64)
					if executeQuery(cVal) {
						return i, true
					}
				} else if hit {
					if analysis.ALogOk {
						break
					}
					c.value.unprocessed = val
					c.value.kind = Number
					c.value.numeric, _ = strconv.ParseFloat(val, 64)
					return i, true
				}
			}
			break
		}
	}
	return i, false
}

// analyzeSubSelectors parses a sub-selection string, which can either be in the form of
// '[path1,path2]' or '{"field1":path1,"field2":path2}' type structure. It returns the parsed
// selectors from the given path, which includes the name and path of each selector within
// the structure. The function assumes that the first character in the path is either '[' or '{',
// and this check is expected to be performed before calling the function.
//
// Parameters:
//   - path: A string representing the sub-selection in either array or object format. The string
//     must begin with either '[' (array) or '{' (object), and the structure should contain
//     valid selectors or field-path pairs.
//
// Returns:
//   - selectors: A slice of `subSelector` structs containing the parsed selectors and their associated paths.
//   - out: The remaining part of the path after parsing the selectors. This will be the part following the
//     closing bracket (']') or brace ('}') if applicable.
//   - ok: A boolean indicating whether the parsing was successful. It returns true if the parsing was
//     successful and the structure was valid, or false if there was an error during parsing.
//
// Example Usage:
//
//	path := "[field1:subpath1,field2:subpath2]"
//	selectors, out, ok := analyzeSubSelectors(path)
//	// selectors: [{name: "field1", path: "subpath1"}, {name: "field2", path: "subpath2"}]
//	// out: "" (no remaining part of the path)
//	// ok: true (parsing was successful)
//
// Details:
//   - The function iterates through each character of the input path and identifies different types
//     of characters (e.g., commas, colons, brackets, braces, and quotes).
//   - It tracks the depth of nested structures (array or object) using the `depth` variable. This ensures
//     proper handling of nested elements within the sub-selection.
//   - The function supports escaping characters with backslashes ('\') and handles this case while parsing.
//   - If a colon (':') is encountered, it indicates a potential name-path pair. The function captures
//     the name and path accordingly, and if no colon is found, it assumes the value is just a path.
//   - The function handles both array-style sub-selections (e.g., [path1,path2]) and object-style
//     sub-selections (e.g., {"field1":path1,"field2":path2}).
//   - If an error is encountered during parsing (e.g., mismatched brackets or braces), the function
//     returns an empty slice and `false` to indicate a failure.
//
// Flow:
//   - The function first initializes tracking variables like `modifier`, `depth`, `colon`, and `start`.
//   - It iterates through the path, checking for different characters, such as backslashes (escape),
//     colons (for name-path pair separation), commas (for separating selectors), and brackets/braces (for
//     nested structures).
//   - If a valid selector is found, it is stored in the `selectors` slice.
//   - The function returns the parsed selectors, the remaining path, and a success flag.
func analyzeSubSelectors(path string) (selectors []subSelector, out string, ok bool) {
	modifier := 0
	depth := 1
	colon := 0
	start := 1
	i := 1
	pushSelectors := func() {
		var selector subSelector
		if colon == 0 {
			selector.path = path[start:i]
		} else {
			selector.name = path[start:colon]
			selector.path = path[colon+1 : i]
		}
		selectors = append(selectors, selector)
		colon = 0
		modifier = 0
		start = i + 1
	}
	for ; i < len(path); i++ {
		switch path[i] {
		case '\\':
			i++
		case '@':
			if modifier == 0 && i > 0 && (path[i-1] == '.' || path[i-1] == '|') {
				modifier = i
			}
		case ':':
			if modifier == 0 && colon == 0 && depth == 1 {
				colon = i
			}
		case ',':
			if depth == 1 {
				pushSelectors()
			}
		case '"':
			i++
		loop:
			for ; i < len(path); i++ {
				switch path[i] {
				case '\\':
					i++
				case '"':
					break loop
				}
			}
		case '[', '(', '{':
			depth++
		case ']', ')', '}':
			depth--
			if depth == 0 {
				pushSelectors()
				path = path[i+1:]
				return selectors, path, true
			}
		}
	}
	return
}

// adjustModifier parses a given path to identify a modifier function and its associated arguments,
// then applies the modifier to the provided JSON string based on the parsed path. This function expects
// that the path starts with a '@', indicating the presence of a modifier. It identifies the modifier's
// name, extracts any potential arguments, and returns the modified result along with the remaining path
// after processing the modifier.
//
// Parameters:
//   - json: A string containing the JSON data that the modifier will operate on.
//   - path: A string representing the path, which includes a modifier prefixed by '@'. The path may
//     contain an optional argument to be processed by the modifier.
//
// Returns:
//   - pYield: The remaining portion of the path after parsing the modifier and its arguments.
//   - result: The result obtained by applying the modifier to the JSON string, or an empty string
//     if no valid modifier is found.
//   - ok: A boolean indicating whether the modifier was successfully identified and applied. If true,
//     the modifier was found and applied; if false, the modifier was not found.
//
// Example Usage:
//
//	json := `{"key": "value"}`
//	path := "@modifierName:argument"
//	pYield, result, ok := adjustModifier(json, path)
//	// pYield: remaining path after the modifier
//	// result: the modified JSON result based on the modifier applied
//	// ok: true if the modifier was found and applied successfully
//
// Details:
//   - The function first removes the '@' character from the beginning of the path and processes the
//     remaining portion of the path to extract the modifier's name and its optional arguments.
//   - The function handles various formats of arguments, including JSON-like objects, arrays, strings,
//     and other specific cases based on the character delimiters such as '{', '[', '"', or '('.
//   - If a valid modifier function is found in the `modifiers` map, it applies the function to the JSON
//     string and returns the result along with the remaining path. If no valid modifier is found, it
//     returns the original path and an empty result.
func adjustModifier(json, path string) (pathYield, result string, ok bool) {
	name := path[1:] // remove the '@' character and initialize the name to the remaining path.
	var hasArgs bool
	// iterate over the path to find the modifier name and any arguments.
	for i := 1; i < len(path); i++ {
		// check for argument delimiter (':'), process if found.
		if path[i] == ':' {
			pathYield = path[i+1:]
			name = path[1:i]
			hasArgs = len(pathYield) > 0
			break
		}
		// check for pipe ('|'), dot ('.'), or other delimiters to separate the modifier name and arguments.
		if path[i] == '|' {
			pathYield = path[i:]
			name = path[1:i]
			break
		}
		if path[i] == '.' {
			pathYield = path[i:]
			name = path[1:i]
			break
		}
	}
	// check if the modifier exists in the modifiers map and apply it if found.
	if fn, ok := modifiers[name]; ok {
		var args string
		if hasArgs { // if arguments are found, parse and handle them.
			var parsedArgs bool
			// process the arguments based on their type (e.g., JSON, string, etc.).
			switch pathYield[0] {
			case '{', '[', '"': // handle JSON-like arguments.
				ctx := Parse(pathYield)
				if ctx.Exists() {
					args = squash(pathYield) // squash the JSON to remove nested structures.
					pathYield = pathYield[len(args):]
					parsedArgs = true
				}
			}
			if !parsedArgs { // process arguments if not already parsed as JSON.
				i := 0
				// iterate through the arguments and process any nested structures or strings.
				for ; i < len(pathYield); i++ {
					if pathYield[i] == '|' {
						break
					}
					switch pathYield[i] {
					case '{', '[', '"', '(': // handle nested structures like arrays or objects.
						s := squash(pathYield[i:])
						i += len(s) - 1
					}
				}
				args = pathYield[:i]      // extract the argument portion.
				pathYield = pathYield[i:] // update the remaining path.
			}
		}
		// apply the modifier function to the JSON data and return the result.
		return pathYield, fn(json, args), true
	}
	// if no modifier is found, return the path and an empty result.
	return pathYield, result, false
}
