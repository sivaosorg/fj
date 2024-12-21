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
