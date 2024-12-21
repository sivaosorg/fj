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
			_, result = parseSquash(name, 0)
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
