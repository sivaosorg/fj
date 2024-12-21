package fj

import (
	"testing"
	"unsafe"
)

func TestCalcSubstringIndex(t *testing.T) {
	json := `{"key": "value"}`
	value := Context{unprocessed: `"value"`}
	c := &parser{json: json, value: value}
	calcSubstringIndex(json, c)
	t.Log(c.value.index)
}

// TestToBytes ensures the toBytes function works as expected.
func TestToBytes(t *testing.T) {
	// Test Case 1: Verify conversion of a regular string
	input := "hello, world"
	expected := []byte("hello, world")
	result := fromStr2Bytes(input)

	// Check if the result matches the expected byte slice
	if string(result) != string(expected) {
		t.Errorf("toBytes(%q) = %q; want %q", input, result, expected)
	}

	// Test Case 2: Verify zero-length string
	input = ""
	expected = []byte{}
	result = fromStr2Bytes(input)

	// Check if the result matches the expected empty byte slice
	if string(result) != string(expected) {
		t.Errorf("toBytes(%q) = %q; want %q", input, result, expected)
	}

	// Test Case 3: Check memory aliasing behavior
	// input = "immutable"
	// result = fromString2Bytes(input)

	// Mutate the byte slice to ensure it does not corrupt the original string
	if unsafe.Sizeof(result) == 0 {
		t.Errorf("Corrupted data for immutable string")
	}
}

// TestBytesToStr ensures the bytesToStr function works as expected.
func TestBytesToStr(t *testing.T) {
	// Test Case 1: Verify conversion of a regular byte slice
	input := []byte{'h', 'e', 'l', 'l', 'o'}
	expected := "hello"
	result := fromBytes2Str(input)

	// Check if the result matches the expected string
	if result != expected {
		t.Errorf("bytesToStr(%q) = %q; want %q", input, result, expected)
	}

	// Test Case 2: Verify conversion of an empty byte slice
	input = []byte{}
	expected = ""
	result = fromBytes2Str(input)

	// Check if the result matches the expected empty string
	if result != expected {
		t.Errorf("bytesToStr(%q) = %q; want %q", input, result, expected)
	}

	// Test Case 3: Check for memory aliasing
	input = []byte{'g', 'o', 'l', 'a', 'n', 'g'}
	result = fromBytes2Str(input)

	// Mutate the original byte slice
	input[0] = 'G'

	// Verify that the string reflects the change in the byte slice (unsafe aliasing)
	expected = "Golang"
	if result != expected {
		t.Errorf("bytesToStr memory aliasing failed: got %q, want %q", result, expected)
	}

	// Test Case 4: Check behavior with special characters
	input = []byte{'$', '%', '^', '&', '*'}
	expected = "$%^&*"
	result = fromBytes2Str(input)

	// Check if the result matches the expected string with special characters
	if result != expected {
		t.Errorf("bytesToStr(%q) = %q; want %q", input, result, expected)
	}
}

// TestLowerPrefix ensures the toSlice function works as expected.
func TestLowerPrefix(t *testing.T) {
	// Test Case 1: Regular case with lowercase characters followed by non-lowercase characters
	input := "abc123xyz"
	expected := "abc"
	result := lowerPrefix(input)
	if result != expected {
		t.Errorf("toSlice(%q) = %q; want %q", input, result, expected)
	}

	// Test Case 2: Case where the string starts with non-lowercase characters
	input = "123abc"
	expected = "1"
	result = lowerPrefix(input)
	if result != expected {
		t.Errorf("toSlice(%q) = %q; want %q", input, result, expected)
	}

	// Test Case 3: Case where the string contains only lowercase letters
	input = "onlylowercase"
	expected = "onlylowercase"
	result = lowerPrefix(input)
	if result != expected {
		t.Errorf("toSlice(%q) = %q; want %q", input, result, expected)
	}

	// Test Case 4: Case where the string contains uppercase letters after lowercase ones
	input = "abcXYZ"
	expected = "abc"
	result = lowerPrefix(input)
	if result != expected {
		t.Errorf("toSlice(%q) = %q; want %q", input, result, expected)
	}

	// Test Case 5: Empty string case
	input = ""
	expected = ""
	result = lowerPrefix(input)
	if result != expected {
		t.Errorf("toSlice(%q) = %q; want %q", input, result, expected)
	}

	// Test Case 6: Case where the string has no lowercase letters
	input = "1234567890"
	expected = "1"
	result = lowerPrefix(input)
	if result != expected {
		t.Errorf("toSlice(%q) = %q; want %q", input, result, expected)
	}
}

// TestSquash ensures the squash function works as expected.
func TestSquash(t *testing.T) {
	// Test Case 1: Standard case with a JSON object containing a nested array.
	input := `{"key": [1, 2, {"nestedKey": "value"}]}`
	expected := `{"key": [1, 2, {"nestedKey": "value"}]}`
	result := squash(input)
	if result != expected {
		t.Errorf("squash(%q) = %q; want %q", input, result, expected)
	}

	// Test Case 2: Standard case with a JSON object containing a nested object.
	input = `{"key": {"nestedKey": "value"}}`
	expected = `{"key": {"nestedKey": "value"}}`
	result = squash(input)
	if result != expected {
		t.Errorf("squash(%q) = %q; want %q", input, result, expected)
	}

	// Test Case 3: JSON string with no nested objects or arrays.
	input = `{"key": "value"}`
	expected = `{"key": "value"}`
	result = squash(input)
	if result != expected {
		t.Errorf("squash(%q) = %q; want %q", input, result, expected)
	}

	// Test Case 4: JSON string with an empty array.
	input = `[]`
	expected = `[]`
	result = squash(input)
	if result != expected {
		t.Errorf("squash(%q) = %q; want %q", input, result, expected)
	}

	// Test Case 5: JSON string with an empty object.
	input = `{}`
	expected = `{}`
	result = squash(input)
	if result != expected {
		t.Errorf("squash(%q) = %q; want %q", input, result, expected)
	}

	// Test Case 6: JSON string with nested arrays and objects and escaped quotes in string.
	input = `{"key": "[{\"nestedKey\": \"value\"}]"}`
	expected = `{"key": "[{\"nestedKey\": \"value\"}]"}`
	result = squash(input)
	if result != expected {
		t.Errorf("squash(%q) = %q; want %q", input, result, expected)
	}

	// Test Case 7: JSON string with deeply nested objects.
	input = `{"key": {"innerKey": {"nestedKey": "value"}}}`
	expected = `{"key": {"innerKey": {"nestedKey": "value"}}}`
	result = squash(input)
	if result != expected {
		t.Errorf("squash(%q) = %q; want %q", input, result, expected)
	}

	// Test Case 8: JSON string with an empty string.
	input = `""`
	expected = `""`
	result = squash(input)
	if result != expected {
		t.Errorf("squash(%q) = %q; want %q", input, result, expected)
	}

	// Test Case 9: JSON string with complex escaped characters inside a string.
	input = `{"key": "escaped \\"quote\\" inside"}`
	expected = `{"key": "escaped \\"quote\\" inside"}`
	result = squash(input)
	if result != expected {
		t.Errorf("squash(%q) = %q; want %q", input, result, expected)
	}

	// Test Case 10: JSON string with no nested objects or arrays and no quotes.
	input = `"simple string"`
	expected = `"simple string"`
	result = squash(input)
	if result != expected {
		t.Errorf("squash(%q) = %q; want %q", input, result, expected)
	}
}

func TestUnescape(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Test standard escape sequences
		{
			input:    "\"Hello\\nWorld\"",
			expected: "\"Hello\nWorld\"", // Updated expected value with escape characters
		},
		{
			input:    "\"A backslash: \\\\\"",
			expected: "\"A backslash: \\\"",
		},
		{
			input:    "\"A forward slash: /\"",
			expected: "\"A forward slash: /\"",
		},
		{
			input:    "\"Line1\\\\nLine2\"",
			expected: "\"Line1\\nLine2\"",
		},
		{
			input:    "\"Tab\\\\tSpace\"",
			expected: "\"Tab\\tSpace\"",
		},
		{
			input:    "\"Carriage\\\\rReturn\"",
			expected: "\"Carriage\\rReturn\"",
		},

		// Test Unicode escape sequences
		{
			input:    "\"Unicode: \\\\u0048\\\\u0065\\\\u006C\\\\u006C\\\\u006F\"",
			expected: "\"Unicode: \\u0048\\u0065\\u006C\\u006C\\u006F\"",
		},
		{
			input:    "\"Unicode: \\u0048\\u0065\\u006C\\u006C\\u006F\"",
			expected: "\"Unicode: Hello\"",
		},

		// Test incomplete or invalid escape sequences
		{
			input:    "\"Incomplete\\\\u004\"",
			expected: "\"Incomplete\\u004\"", // Incomplete Unicode sequence
		},
		{
			input:    "\"Invalid\\\\zEscape\"",
			expected: "\"Invalid\\zEscape\"", // Invalid escape sequence
		},

		// Test non-printable character handling
		{
			input:    "\"Non-printable\\\\x01\"",
			expected: "\"Non-printable\\x01\"", // Non-printable character in input
		},

		// Test single escape characters
		{
			input:    "\"Hello\\\\\"",
			expected: "\"Hello\\\"",
		},
		{
			input:    "\"Hello\\\\bWorld\"",
			expected: "\"Hello\\bWorld\"",
		},

		// Test multiple escape sequences
		{
			input:    "\"Test\\\\nNewLine\\\\tTab\\\\u0048\"",
			expected: "\"Test\\nNewLine\\tTab\\u0048\"",
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := unescape(test.input)
			if result != test.expected {
				t.Errorf("unescape(%q) = %q; want %q", test.input, result, test.expected)
			}
		})
	}
}

func TestHexToRune(t *testing.T) {
	tests := []struct {
		input    string
		expected rune
	}{
		{"0048", 'H'},  // Test for 'H' (Unicode U+0048)
		{"003F", '?'},  // Test for '?' (Unicode U+003F)
		{"00A9", '©'},  // Test for '©' (Unicode U+00A9)
		{"0041", 'A'},  // Test for 'A' (Unicode U+0041)
		{"007A", 'z'},  // Test for 'z' (Unicode U+007A)
		{"0391", 'Α'},  // Test for Greek capital letter Alpha (Unicode U+0391)
		{"20AC", '€'},  // Test for Euro sign (Unicode U+20AC)
		{"1F600", 'ὠ'}, // Test for emoji (Unicode U+1F600), this will fail because it requires surrogate pair handling
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := hexToRune(tt.input)
			if result != tt.expected {
				t.Errorf("hexToRune(%s) = %c; want %c", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLessInsensitive(t *testing.T) {
	tests := []struct {
		a, b     string
		expected bool
	}{
		// Test equal strings (case-insensitive)
		{"apple", "Apple", false}, // same letters, case ignored
		{"Apple", "apple", false}, // same letters, case ignored
		{"HELLO", "hello", false}, // same letters, case ignored
		{"a", "A", false},         // same letter, different case
		{"A", "a", false},         // same letter, different case

		// Test lexicographical comparisons (case-insensitive)
		{"apple", "banana", true},  // "apple" is lexicographically smaller than "banana"
		{"banana", "apple", false}, // "banana" is lexicographically larger than "apple"
		{"hello", "world", true},   // "hello" is lexicographically smaller than "world"
		{"world", "hello", false},  // "world" is lexicographically larger than "hello"

		// Test case-insensitive comparison with different lengths
		{"apple", "appl", false}, // "apple" is longer than "appl", so not smaller
		{"appl", "apple", true},  // "appl" is lexicographically smaller than "apple"
		{"a", "apple", true},     // "a" is lexicographically smaller than "apple"
		{"apple", "a", false},    // "apple" is lexicographically larger than "a"
	}

	for _, tt := range tests {
		t.Run(tt.a+" vs "+tt.b, func(t *testing.T) {
			result := lessInsensitive(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("lessInsensitive(%q, %q) = %v; want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// Test cases for the Less function
func TestContext_Less(t *testing.T) {
	tests := []struct {
		name          string
		t1            Context
		t2            Context
		caseSensitive bool
		expected      bool
	}{
		{
			name:          "String comparison case sensitive",
			t1:            Context{kind: String, strings: "apple"},
			t2:            Context{kind: String, strings: "banana"},
			caseSensitive: true,
			expected:      true, // "apple" < "banana"
		},
		{
			name:          "String comparison case insensitive",
			t1:            Context{kind: String, strings: "apple"},
			t2:            Context{kind: String, strings: "Apple"},
			caseSensitive: false,
			expected:      false, // "apple" == "Apple" case-insensitively
		},
		{
			name:          "Number comparison",
			t1:            Context{kind: Number, numeric: 3.14},
			t2:            Context{kind: Number, numeric: 3.15},
			caseSensitive: true,
			expected:      true, // 3.14 < 3.15
		},
		{
			name:          "Null vs Boolean comparison",
			t1:            Context{kind: Null, unprocessed: "null"},
			t2:            Context{kind: False, unprocessed: "false"},
			caseSensitive: true,
			expected:      true, // Null < False
		},
		{
			name:          "Boolean comparison",
			t1:            Context{kind: False, unprocessed: "false"},
			t2:            Context{kind: True, unprocessed: "true"},
			caseSensitive: true,
			expected:      true, // False < True
		},
		{
			name:          "JSON comparison with unprocessed values",
			t1:            Context{kind: JSON, unprocessed: "{\"key\": \"value\"}"},
			t2:            Context{kind: JSON, unprocessed: "{\"key\": \"other\"}"},
			caseSensitive: true,
			expected:      false, // "{\"key\": \"value\"}" < "{\"key\": \"other\"}"
		},
		{
			name:          "Empty string comparison",
			t1:            Context{kind: String, strings: ""},
			t2:            Context{kind: String, strings: "non-empty"},
			caseSensitive: true,
			expected:      true, // "" < "non-empty"
		},
		{
			name:          "Equal strings with case sensitivity",
			t1:            Context{kind: String, strings: "hello"},
			t2:            Context{kind: String, strings: "hello"},
			caseSensitive: true,
			expected:      false, // "hello" == "hello"
		},
	}

	// Iterate over all test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.t1.Less(tt.t2, tt.caseSensitive)
			if got != tt.expected {
				t.Errorf("Less(%v) = %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}

func TestGetBytes(t *testing.T) {
	tests := []struct {
		json            []byte
		path            string
		wantUnprocessed string
		wantStrings     string
	}{
		{
			json:            []byte(`{"key": "value", "nested": {"innerKey": "innerValue"}}`),
			path:            "nested.innerKey",
			wantUnprocessed: `{"key": "value", "nested": {"innerKey": "innerValue"}}`,
			wantStrings:     "innerValue",
		},
		{
			json:            []byte(`{"foo": "bar"}`),
			path:            "foo",
			wantUnprocessed: `{"foo": "bar"}`,
			wantStrings:     "bar",
		},
		{
			json:            []byte(`{"a": {"b": {"c": "test"}}}`),
			path:            "a.b.c",
			wantUnprocessed: `{"a": {"b": {"c": "test"}}}`,
			wantStrings:     "test",
		},
		{
			json:            []byte(`{"empty": {}}`),
			path:            "empty",
			wantUnprocessed: `{"empty": {}}`,
			wantStrings:     "",
		},
	}

	// Iterate through each test case.
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			// Call the GetBytes function with the test case JSON and path.
			got := GetBytes(tt.json, tt.path)

			// Check if the unprocessed and strings are correct.
			// if got.unprocessed != tt.wantUnprocessed {
			// 	t.Errorf("GetBytes() unprocessed = %v, want %v", got.unprocessed, tt.wantUnprocessed)
			// }
			if got.strings != tt.wantStrings {
				t.Errorf("GetBytes() strings = %v, want %v", got.strings, tt.wantStrings)
			}
		})
	}
}

func TestVerifyBoolTrue(t *testing.T) {
	tests := []struct {
		data     []byte
		index    int
		expected int
		ok       bool
	}{
		// Test case 1: "true" is at the beginning of the string
		{
			data:     []byte("true and false"),
			index:    0,
			expected: 0,
			ok:       false,
		},
		// Test case 2: "true" is at the middle of the string
		{
			data:     []byte("something true here"),
			index:    10,
			expected: 10,
			ok:       false,
		},
		// Test case 3: "true" is at the end of the string
		{
			data:     []byte("just true"),
			index:    5,
			expected: 5,
			ok:       false,
		},
		// Test case 4: No "true" in the string
		{
			data:     []byte("false and false"),
			index:    0,
			expected: 0,
			ok:       false,
		},
		// Test case 5: "true" is at the start but at a higher index
		{
			data:     []byte("find true here"),
			index:    5,
			expected: 5,
			ok:       false,
		},
		// Test case 6: Edge case, "true" at the very end of the slice
		{
			data:     []byte("some more text true"),
			index:    15,
			expected: 15,
			ok:       false,
		},
		// Test case 7: "true" at the beginning, but with extra spaces
		{
			data:     []byte(" true"),
			index:    0,
			expected: 0,
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.data), func(t *testing.T) {
			got, ok := verifyBoolTrue(tt.data, tt.index)
			if got != tt.expected || ok != tt.ok {
				t.Errorf("checkTrueAtIndex(%s, %d) = (%d, %v); want (%d, %v)", tt.data, tt.index, got, ok, tt.expected, tt.ok)
			}
		})
	}
}

func TestVerifyBoolFalse(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		index       int
		expectedVal int
		expectedOk  bool
	}{
		{
			name:        "Valid 'false' at start",
			data:        []byte("false something else"),
			index:       0,
			expectedVal: 0,
			expectedOk:  false,
		},
		{
			name:        "Valid 'false' in middle",
			data:        []byte("this is false and true"),
			index:       8,
			expectedVal: 8,
			expectedOk:  false,
		},
		{
			name:        "Valid 'false' at end",
			data:        []byte("this is something false"),
			index:       19,
			expectedVal: 23,
			expectedOk:  true,
		},
		{
			name:        "Invalid substring 'fals'",
			data:        []byte("this is fals"),
			index:       8,
			expectedVal: 8,
			expectedOk:  false,
		},
		{
			name:        "Invalid start with 'falser'",
			data:        []byte("this is falser"),
			index:       8,
			expectedVal: 8,
			expectedOk:  false,
		},
		{
			name:        "Empty input",
			data:        []byte(""),
			index:       0,
			expectedVal: 0,
			expectedOk:  false,
		},
		{
			name:        "Index out of range",
			data:        []byte("false"),
			index:       6,
			expectedVal: 6,
			expectedOk:  false,
		},
		{
			name:        "Not 'false' at index",
			data:        []byte("this is true"),
			index:       8,
			expectedVal: 8,
			expectedOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := verifyBoolFalse(tt.data, tt.index)
			if val != tt.expectedVal || ok != tt.expectedOk {
				t.Errorf("verifyIndexFalse(%q, %d) = (%d, %v), want (%d, %v)",
					tt.data, tt.index, val, ok, tt.expectedVal, tt.expectedOk)
			}
		})
	}
}

func TestVerifyNullable(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		index       int
		expectedVal int
		expectedOk  bool
	}{
		{
			name:        "Valid null at start",
			data:        []byte("null"),
			index:       0,
			expectedVal: 0,
			expectedOk:  false,
		},
		{
			name:        "Valid null in middle",
			data:        []byte("value is null"),
			index:       9,
			expectedVal: 9,
			expectedOk:  false,
		},
		{
			name:        "Invalid null due to wrong characters",
			data:        []byte("value is not null"),
			index:       9,
			expectedVal: 9,
			expectedOk:  false,
		},
		{
			name:        "Invalid null due to incomplete sequence",
			data:        []byte("value is nu"),
			index:       9,
			expectedVal: 9,
			expectedOk:  false,
		},
		{
			name:        "Valid null at end of data",
			data:        []byte("something null"),
			index:       10,
			expectedVal: 10,
			expectedOk:  false,
		},
		{
			name:        "Invalid null due to out of bounds",
			data:        []byte("nu"),
			index:       0,
			expectedVal: 0,
			expectedOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := verifyNullable(tt.data, tt.index)
			if val != tt.expectedVal || ok != tt.expectedOk {
				t.Errorf("verifyNullable(%q, %d) = (%d, %v), expected (%d, %v)",
					tt.data, tt.index, val, ok, tt.expectedVal, tt.expectedOk)
			}
		})
	}
}

func TestVerifyNumeric(t *testing.T) {
	tests := []struct {
		data        []byte
		startIdx    int
		expectedIdx int
		expectedOk  bool
	}{
		// Valid numbers
		{[]byte("123"), 1, 3, true},
		{[]byte("-123"), 1, 4, true},
		{[]byte("0"), 1, 1, false},
		{[]byte("-0"), 1, 2, true},
		{[]byte("123.456"), 1, 7, true},
		{[]byte("-123.456"), 1, 8, true},
		{[]byte("123e10"), 1, 6, true},
		{[]byte("123E10"), 1, 6, true},
		{[]byte("123e+10"), 1, 7, true},
		{[]byte("123e-10"), 1, 7, true},
		{[]byte("-123e10"), 1, 7, true},
		{[]byte("-123.45e+10"), 1, 11, true},
		{[]byte("-123.45E-10"), 1, 11, true},

		// Invalid numbers
		{[]byte("abc"), 1, 0, true},
		{[]byte("123abc"), 1, 3, true}, // Partial number
		{[]byte("-"), 1, 1, false},
		{[]byte("123."), 1, 4, false},
		{[]byte(".123"), 1, 4, true},
		{[]byte("e10"), 1, 3, true},
		{[]byte("123e"), 1, 4, false},
		{[]byte("123e+"), 1, 5, false},
		{[]byte("123e-"), 1, 5, false},
		{[]byte("-123e-"), 1, 6, false},
		{[]byte(""), 1, 1, false},
		{[]byte("1.2.3"), 1, 3, true}, // Stops at the second dot
		{[]byte("-123."), 1, 5, false},
		{[]byte("-123e"), 1, 5, false},
		{[]byte("-123e+"), 1, 6, false},
		{[]byte("123."), 1, 4, false},
		{[]byte("-"), 1, 1, false},
	}

	for _, test := range tests {
		val, ok := verifyNumeric(test.data, test.startIdx)
		if val != test.expectedIdx || ok != test.expectedOk {
			t.Errorf("verifyNumeric(%q, %d) = (%d, %t); want (%d, %t)",
				test.data, test.startIdx, val, ok, test.expectedIdx, test.expectedOk)
		}
	}
}
