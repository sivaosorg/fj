package fj

import (
	"fmt"
	"strings"

	"github.com/sivaosorg/unify4g"
)

// transformDefault is a fallback transformation that simply returns the input JSON string
// without applying any modifications. This function is typically used as a default case
// when no specific transformation is requested or supported.
//
// Parameters:
//   - `json`: The JSON string to be returned as-is.
//   - `arg`: This parameter is unused for this transformation but is included for consistency
//     with other transform functions.
//
// Returns:
//   - The original input JSON string, unchanged.
//
// Example Usage:
//
//	// Input JSON
//	json := `{"name":"Alice","age":25}`
//
//	// No transformation applied, returns original JSON
//	unchangedJSON := transformDefault(json, "")
//	fmt.Println(unchangedJSON)
//	// Output: {"name":"Alice","age":25}
//
// Notes:
//   - This function is used when no transformation is specified or when the transformation
//     request is unsupported. It ensures that the input JSON is returned unmodified.
func transformDefault(json, arg string) string {
	return json
}

// transformPretty formats the input JSON string into a human-readable, indented format.
//
// This function applies "pretty printing" to the provided JSON data, making it easier to read
// and interpret. If additional formatting options are specified in the `arg` parameter, these
// options are parsed and applied to customize the output. Formatting options include sorting
// keys, setting indentation styles, specifying prefixes, and defining maximum line widths.
//
// Parameters:
//   - `json`: The JSON string to be formatted.
//   - `arg`: An optional string containing formatting configuration in JSON format. The configuration
//     can specify the following keys:
//   - `sortKeys`: A boolean value (`true` or `false`) that determines whether keys in JSON objects
//     should be sorted alphabetically.
//   - `indent`: A string containing whitespace characters (e.g., `"  "` or `"\t"`) used for indentation.
//   - `prefix`: A string prepended to each line of the formatted JSON.
//   - `width`: An integer specifying the maximum line width for the formatted output.
//
// Returns:
//   - A string representing the formatted JSON, transformed based on the specified or default options.
//
// Example Usage:
//
//	// Input JSON
//	json := `{"name":"Alice","age":25,"address":{"city":"New York","zip":"10001"}}`
//
//	// Format without additional options
//	prettyJSON := transformPretty(json, "")
//	fmt.Println(prettyJSON)
//	// Output:
//	// {
//	//   "name": "Alice",
//	//   "age": 25,
//	//   "address": {
//	//     "city": "New York",
//	//     "zip": "10001"
//	//   }
//	// }
//
//	// Format with additional options
//	arg := `{"indent": "\t", "sort_keys": true}`
//	prettyJSONWithOpts := transformPretty(json, arg)
//	fmt.Println(prettyJSONWithOpts)
//	// Output:
//	// {
//	// 	"address": {
//	// 		"city": "New York",
//	// 		"zip": "10001"
//	// 	},
//	// 	"age": 25,
//	// 	"name": "Alice"
//	// }
//
// Notes:
//   - If `arg` is empty, default formatting is applied with standard indentation.
//   - The function uses `unify4g.Pretty` or `unify4g.PrettyOptions` for the actual formatting.
//   - Invalid or unrecognized keys in the `arg` parameter are ignored.
//   - The function internally uses `fromStr2Bytes` and `fromBytes2Str` for efficient data conversion.
//
// Implementation Details:
//   - The `arg` string is parsed using the `Parse` function, and each key-value pair is applied
//     to configure the formatting options (`opts`).
//   - The `stripNonWhitespace` function ensures only whitespace characters are used for `indent`
//     and `prefix` settings to prevent formatting errors.
func transformPretty(json, arg string) string {
	if len(arg) > 0 {
		opts := *unify4g.DefaultOptionsConfig
		Parse(arg).Foreach(func(key, value Context) bool {
			switch key.String() {
			case "sort_keys":
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

// transformMinify removes all whitespace characters from the input JSON string,
// transforming it into a compact, single-line format.
//
// This function applies a "minified" transformation to the provided JSON data,
// removing all spaces, newlines, and other whitespace characters. The result is
// a more compact representation of the JSON, which is useful for minimizing
// data size, especially when transmitting JSON data over the network or storing
// it in a compact format.
//
// Parameters:
//   - `json`: The JSON string to be transformed into a compact form.
//   - `arg`: This parameter is unused for this transformation but is included
//     for consistency with other transform functions.
//
// Returns:
//   - A string representing the "ugly" JSON, with all whitespace removed.
//
// Example Usage:
//
//	// Input JSON
//	json := `{
//	  "name": "Alice",
//	  "age": 25,
//	  "address": {
//	    "city": "New York",
//	    "zip": "10001"
//	  }
//	}`
//
//	// Transform to minify (compact) JSON
//	uglyJSON := transformMinify(json, "")
//	fmt.Println(uglyJSON)
//	// Output: {"name":"Alice","age":25,"address":{"city":"New York","zip":"10001"}}
//
// Notes:
//   - The `arg` parameter is not used in this transformation, and its value is ignored.
//   - The function uses `unify4g.Ugly` for the actual transformation, which removes all
//     whitespace from the JSON data.
//   - This function is often used to reduce the size of JSON data for storage or transmission.
func transformMinify(json, arg string) string {
	return fromBytes2Str(unify4g.Ugly(fromStr2Bytes(json)))
}

// transformReverse reverses the order of elements in an array or the order of key-value
// pairs in an object. This function processes the JSON input and applies the reversal
// based on the type of JSON structure: array or object.
//
// If the JSON is an array, it reverses the array elements. If it's an object, it reverses
// the key-value pairs. If the input is neither an array nor an object, the original JSON
// string is returned unchanged.
//
// Parameters:
//   - `json`: The JSON string to be transformed, which may be an array or an object.
//   - `arg`: This parameter is unused for this transformation but is included for consistency
//     with other transform functions.
//
// Returns:
//   - A string representing the transformed JSON with reversed elements or key-value pairs.
//
// Example Usage:
//
//	// Input JSON (array)
//	jsonArray := `[1, 2, 3]`
//
//	// Reverse array elements
//	reversedJSON := transformReverse(jsonArray, "")
//	fmt.Println(reversedJSON)
//	// Output: [3,2,1]
//
//	// Input JSON (object)
//	jsonObject := `{"name":"Alice","age":25}`
//
//	// Reverse key-value pairs
//	reversedObject := transformReverse(jsonObject, "")
//	fmt.Println(reversedObject)
//	// Output: {"age":25,"name":"Alice"}
//
// Notes:
//   - If the input JSON is an array, the array elements are reversed.
//   - If the input JSON is an object, the key-value pairs are reversed.
//   - If the input JSON is neither an array nor an object, the original string is returned unchanged.
func transformReverse(json, arg string) string {
	ctx := Parse(json)
	if ctx.IsArray() {
		var values []Context
		ctx.Foreach(func(_, value Context) bool {
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
	if ctx.IsObject() {
		var keyValues []Context
		ctx.Foreach(func(key, value Context) bool {
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

// transformFlatten flattens a JSON array by removing any nested arrays within it.
//
// This function takes a JSON array (which may contain nested arrays) and flattens it
// into a single array by extracting the elements of any child arrays. The function
// supports both shallow and deep flattening based on the provided argument.
//
// Parameters:
//   - `json`: A string representing the JSON array to be flattened. The array may contain
//     nested arrays that will be flattened into the outer array.
//   - `arg`: An optional string containing configuration options in JSON format. The configuration
//     can specify the following key:
//   - `deep`: A boolean value (`true` or `false`) that determines whether nested arrays should
//     be recursively flattened (deep flattening). If `deep` is `true`, all nested arrays are
//     flattened into the main array, while if `false` (or absent), only the immediate nested arrays
//     are flattened.
//
// Returns:
//   - A string representing the flattened JSON array. The returned array may contain elements
//     from nested arrays, depending on whether deep flattening was requested.
//
// Example Usage:
//
//	// Input JSON (shallow flatten)
//	json := "[1,[2],[3,4],[5,[6,7]]]"
//	shallowFlattened := transformFlatten(json, "")
//	fmt.Println(shallowFlattened)
//	// Output: [1,2,3,4,5,[6,7]]
//
//	// Input JSON (deep flatten)
//	json := "[1,[2],[3,4],[5,[6,7]]]"
//	deepFlattened := transformFlatten(json, "{\"deep\": true}")
//	fmt.Println(deepFlattened)
//	// Output: [1,2,3,4,5,6,7]
//
// Notes:
//
//   - If the input JSON is not an array, the original JSON string is returned unchanged.
//
//   - The function first checks if the provided JSON is an array. If it is not an array, it returns
//     the original input string without any changes.
//
//   - The `deep` option controls whether nested arrays are flattened recursively. If the `deep`
//     option is set to `false` (or omitted), only the immediate nested arrays are flattened.
//
//   - Nested arrays can be flattened either shallowly or deeply depending on the configuration provided
//     in the `arg` parameter.
//
//   - The function uses `removeOuterBraces` to remove the surrounding brackets of nested arrays to
//     achieve the flattening effect.
//
//     [1,[2],[3,4],[5,[6,7]]] -> [1,2,3,4,5,[6,7]]
//
// The {"deep":true} arg can be provide for deep flattening.
//
//	[1,[2],[3,4],[5,[6,7]]] -> [1,2,3,4,5,6,7]
//
// The original json is returned when the json is not an array.
func transformFlatten(json, arg string) string {
	ctx := Parse(json)
	if !ctx.IsArray() {
		return json
	}
	var deep bool
	if unify4g.IsNotEmpty(arg) {
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
	ctx.Foreach(func(_, value Context) bool {
		var raw string
		if value.IsArray() {
			if deep {
				raw = removeOuterBraces(transformFlatten(value.unprocessed, arg))
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

// transformJoin merges multiple JSON objects into a single object.
// If the input is an array of JSON objects, it combines their key-value pairs
// into one object. Duplicate keys can be preserved or discarded based on the
// configuration provided in the `arg` parameter.
//
// Parameters:
//   - `json`: A string representing a JSON array, where each element is a JSON object.
//     The objects will be merged into a single object.
//   - `arg`: A string containing a JSON configuration that can specify whether
//     duplicate keys should be preserved. If `arg` is provided and contains
//     the key `preserve` set to `true`, duplicate keys will be kept in the output object.
//
// Returns:
//   - A string representing the merged JSON object. If the input is not an array
//     of JSON objects, the function returns the original `json` string unchanged.
//
// Example Usage:
//
//	// Input JSON (merge objects with duplicate keys discarded)
//	json := `[{"first":"Tom","age":37},{"age":41}]`
//	mergedJSON := transformJoin(json, "")
//	fmt.Println(mergedJSON)
//	// Output: {"first":"Tom","age":41}
//
//	// Input JSON (merge objects with duplicate keys preserved)
//	json := `[{"first":"Tom","age":37},{"age":41}]`
//	mergedJSONWithDupes := transformJoin(json, "{\"preserve\": true}")
//	fmt.Println(mergedJSONWithDupes)
//	// Output: {"first":"Tom","age":37,"age":41}
//
// Notes:
//   - If the input `json` is not a valid array of JSON objects, the function returns
//     the original input string unchanged.
//   - The `preserve` option controls whether duplicate keys should be kept in the merged object.
//     If `preserve` is `false` (or absent), only the last occurrence of each key is kept.
//   - The function uses `removeOuterBraces` to remove any extraneous brackets around JSON objects
//     before merging their contents.
//
// Implementation Details:
//   - If the `preserve` option is set to `true`, all key-value pairs from the objects are
//     appended to the resulting object, even if keys are repeated.
//   - If `preserve` is `false`, the function will deduplicate keys by selecting the last
//     value for each key across all objects in the array. The keys are also added in stable
//     order based on their appearance in the input objects.
func transformJoin(json, arg string) string {
	ctx := Parse(json)
	if !ctx.IsArray() {
		return json
	}
	var preserve bool
	if unify4g.IsNotEmpty(arg) {
		Parse(arg).Foreach(func(key, value Context) bool {
			if key.String() == "preserve" {
				preserve = value.Bool()
			}
			return true
		})
	}
	var target []byte
	target = append(target, '{')
	if preserve { // preserve duplicate keys.
		var idx int
		ctx.Foreach(func(_, value Context) bool {
			if !value.IsObject() {
				return true
			}
			if idx > 0 {
				target = append(target, ',')
			}
			target = append(target, removeOuterBraces(value.unprocessed)...)
			idx++
			return true
		})
	} else { // deduplicate keys and generate an object with stable ordering.
		var keys []Context
		keyVal := make(map[string]Context)
		ctx.Foreach(func(_, value Context) bool {
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
				target = append(target, ',')
			}
			target = append(target, keys[i].unprocessed...)
			target = append(target, ':')
			target = append(target, keyVal[keys[i].String()].unprocessed...)
		}
	}
	target = append(target, '}')
	return fromBytes2Str(target)
}

// transformJSONValidity ensures that the input JSON is valid before processing it further.
// If the provided JSON string is not valid, the function returns an empty string. If the
// JSON is valid, the original JSON string is returned unchanged.
//
// Parameters:
//   - `json`: A string representing the JSON data that needs to be validated.
//   - `arg`: This parameter is not used in this function but is included for consistency
//     with the signature of other transformation functions.
//
// Returns:
//   - A string: The original JSON string if it is valid, or an empty string if the JSON
//     is not valid.
//
// Example Usage:
//
//	// Input JSON (valid)
//	json := `{"name":"Alice","age":25}`
//	validJSON := transformJSONValidity(json, "")
//	fmt.Println(validJSON)
//	// Output: {"name":"Alice","age":25}
//
//	// Input JSON (invalid)
//	json := `{"name":"Alice","age":25`
//	validJSON := transformJSONValidity(json, "")
//	fmt.Println(validJSON)
//	// Output: (empty string)
//
// Notes:
//   - The function uses `IsValidJSON` (presumably another function in the codebase) to
//     check whether the input string is valid JSON.
//   - If the JSON is invalid, the function returns an empty string, which can be used
//     to stop further processing.
//   - If the JSON is valid, it returns the JSON string unmodified, allowing subsequent
//     transformations or processing to proceed.
func transformJSONValidity(json, arg string) string {
	if !IsValidJSON(json) {
		return strings.ToLower(False.String())
	}
	return strings.ToLower(True.String())
}

// transformKeys extracts the keys from a JSON object and returns them as a JSON array of strings.
// The function processes the input JSON, identifies whether it is an object, and then generates
// an array containing the keys of the object. If the input is not a valid JSON object, it returns
// an empty array.
//
// Parameters:
//   - `json`: A string representing the JSON data, which should be an object from which keys will be extracted.
//   - `arg`: This parameter is not used in this function but is included for consistency with other transformation functions.
//
// Returns:
//   - A string representing a JSON array of keys, or an empty array (`[]`) if the input is not a valid object.
//
// Example Usage:
//
//	// Input JSON (object)
//	json := `{"first":"Tom","last":"Smith"}`
//	keys := transformKeys(json, "")
//	fmt.Println(keys)
//	// Output: ["first","last"]
//
//	// Input JSON (non-object)
//	json := `"Tom"`
//	keys := transformKeys(json, "")
//	fmt.Println(keys)
//	// Output: []
//
// Notes:
//   - If the input JSON is an object, the function will iterate through the keys of the object and return them in
//     a JSON array format.
//   - If the input JSON is not an object (e.g., an array, string, or invalid), the function will return an empty array (`[]`).
//   - The function relies on the `Parse` function to parse the input JSON and the `Foreach` method to iterate over
//     the object keys.
//   - The `unprocessed` method is used to extract the raw key value as a string without further processing.
//
// Implementation Details:
//   - The function first checks if the parsed JSON object exists. If it does, it iterates through the object and extracts
//     the keys. Each key is added to a string builder, and the keys are wrapped in square brackets to form a valid JSON array.
//   - If the JSON is not an object, the function immediately returns an empty array (`[]`).
func transformKeys(json, arg string) string {
	ctx := Parse(json)
	if !ctx.Exists() {
		return "[]"
	}
	var i int
	var builder strings.Builder
	o := ctx.IsObject()
	builder.WriteByte('[')
	ctx.Foreach(func(key, _ Context) bool {
		if i > 0 {
			builder.WriteByte(',')
		}
		if o {
			builder.WriteString(key.unprocessed)
		} else {
			builder.WriteString("null")
		}
		i++
		return true
	})
	builder.WriteByte(']')
	return builder.String()
}

// transformValues extracts the values from a JSON object and returns them as a JSON array of values.
//
// This function parses the input JSON string, and if the JSON is an object, it extracts all the values
// from the key-value pairs and returns them as a JSON array of values. If the input JSON is already an array,
// it simply returns the original JSON string. If the input does not contain a valid JSON object or array,
// it returns an empty array ("[]").
//
// Parameters:
//   - `json`: The JSON string to extract values from. It can be a JSON object or array.
//   - `arg`: An optional argument that is not used in this function, but can be included for consistency
//     with other transformation functions.
//
// Returns:
//   - A string representing a JSON array containing the values extracted from the input JSON object.
//     If the input JSON is already an array, it is returned as-is. If the input is invalid or empty,
//     an empty array ("[]") is returned.
//
// Example Usage:
//
//	// Input JSON representing an object
//	json := `{"first":"Aris","last":"Nguyen"}`
//
//	// Extract the values from the object
//	values := transformValues(json, "")
//	fmt.Println(values) // Output: ["Aris","Nguyen"]
//
//	// Input JSON representing an array
//	jsonArray := `["apple", "banana", "cherry"]`
//
//	// Return the array as-is
//	values := transformValues(jsonArray, "")
//	fmt.Println(values) // Output: ["apple", "banana", "cherry"]
//
//	// Input JSON representing an invalid object
//	invalidJson := `{"key":}` // Invalid JSON
//
//	// Return empty array for invalid JSON
//	values := transformValues(invalidJson, "")
//	fmt.Println(values) // Output: []
//
// Details:
//   - The function first parses the input JSON string using `Parse`.
//   - If the input is an array, the function directly returns the original string as it is.
//   - If the input is an object, the function iterates over its key-value pairs, extracting only the values,
//     and then constructs a JSON array of these values.
//   - If the input JSON does not exist or is invalid, the function returns an empty JSON array ("[]").
func transformValues(json, arg string) string {
	ctx := Parse(json)
	if !ctx.Exists() {
		return "[]"
	}
	if ctx.IsArray() {
		return json
	}
	var i int
	var builder strings.Builder
	builder.WriteByte('[')
	ctx.Foreach(func(_, value Context) bool {
		if i > 0 {
			builder.WriteByte(',')
		}
		builder.WriteString(value.unprocessed)
		i++
		return true
	})
	builder.WriteByte(']')
	return builder.String()
}

// transformToJSON converts a string into a valid JSON representation.
//
// This function ensures that the input string is a valid JSON before attempting to
// parse and convert it into its corresponding JSON format. If the input string is
// a valid JSON, the function returns the formatted JSON as a string. Otherwise,
// it returns an empty string to indicate that the input was not valid JSON.
//
// Parameters:
//   - `json`: A string representing the data that needs to be converted to a valid JSON format.
//   - `arg`: An additional argument which is unused in this function. It may be a placeholder for future extensions.
//
// Returns:
//   - A string representing the input data in valid JSON format. If the input string is not valid JSON,
//     an empty string is returned.
//
// Example Usage:
//
//	// Input string
//	json := "{\"id\":1023,\"name\":\"alert\"}"
//
//	// Convert to valid JSON representation
//	result := transformToJSON(json, "")
//	fmt.Println(result)  // Output: {"id":1023,"name":"alert"}
//
//	// Invalid input string
//	invalidJson := "\"id\":1023,\"name\":\"alert\""
//	result = transformToJSON(invalidJson, "")
//	fmt.Println(result)  // Output: ""
//
// Notes:
//   - This function uses the `IsValidJSON` helper to check if the input string is a valid JSON format.
//   - If the input string is valid JSON, the `Parse` function is used to parse and format it, ensuring it is returned in the proper JSON format.
//   - If the input is invalid, an empty string is returned, indicating that the transformation failed.
func transformToJSON(json, arg string) string {
	if !IsValidJSON(json) {
		return ""
	}
	return Parse(json).String()
}

// transformToString converts a regular string into a valid JSON string format.
//
// This function takes an input string and converts it into a JSON-encoded string
// by wrapping it in double quotes and escaping any necessary characters (such as
// quotes or backslashes) to ensure the resulting string adheres to JSON encoding rules.
//
// Parameters:
//   - `str`: The input string to be converted into JSON string format.
//   - `arg`: This parameter is not used in this function but is present for consistency
//     with the function signature, which might be useful for future extensions or
//     compatibility with other transformations.
//
// Returns:
//   - A string that is a valid JSON representation of the input `str`, with special
//     characters properly escaped (e.g., double quotes, backslashes, control characters).
//
// Example Usage:
//
//	str := "Hello \"world\"\nLine break!"
//	result := transformToString(str, "")
//	fmt.Println(result)
//	// Output: "\"Hello \\\"world\\\"\nLine break!\""
//
// Notes:
//   - This function calls `appendJSON` to handle the conversion of the string into a
//     valid JSON string format, ensuring proper escaping of special characters (like
//     double quotes and newlines) to maintain valid JSON syntax.
//   - The `arg` parameter is included for consistency with other transformation functions
//     that may require it, though it does not affect the behavior of this specific function.
func transformToString(str, arg string) string {
	return string(appendJSON(nil, str))
}

// transformGroup processes a JSON string containing objects and arrays, and groups the
// elements of arrays within objects by their keys. It converts each array into a group of
// key-value pairs, resulting in a new JSON structure where each object is grouped by its array values.
//
// This function is primarily used to reformat JSON data by grouping values in arrays under
// the same key into new JSON objects, preserving the structure of the original JSON while
// transforming the array elements into key-value pairs associated with the respective key.
//
// Parameters:
//   - `json`: A string representing the JSON data that needs to be transformed. It is assumed
//     to be in the form of a JSON object that contains arrays.
//   - `arg`: A string argument that provides additional options for the transformation. In this
//     case, the argument is not used but is still passed for consistency with the other transformation
//     functions.
//
// Returns:
//   - A new string representing the transformed JSON data. The arrays in the input JSON are grouped
//     under their respective keys, and the resulting structure is a list of objects.
//
// Example Usage:
//
//	json := `{"a": [1, 2, 3], "b": [4, 5]}`
//	result := transformGroup(json, "")
//	fmt.Println(result)
//	// Output: `[
//	//  {"a":1},
//	//  {"a":2},
//	//  {"a":3},
//	//  {"b":4},
//	//  {"b":5}
//	// ]`
//
// Notes:
//   - The function works by iterating over each key-value pair in the input JSON object. If the value
//     associated with a key is an array, it processes the array's elements by creating new objects where
//     the key is used for each element, effectively turning each element into a separate object with its
//     corresponding key.
//   - If the input JSON does not contain arrays, it does not affect the final result.
//
// Implementation Details:
//   - The function first parses the input JSON string into a context object using `Parse(json)`.
//   - It checks if the context is a valid object (`ctx.IsObject()`). If not, it returns an empty string.
//   - It iterates over the object's key-value pairs using `ctx.Foreach()`. For each array value, it creates
//     new objects where each array element is associated with the respective key.
//   - The transformed groups of key-value pairs are collected into a new array format, where each array
//     element corresponds to a new object created from the array values.
func transformGroup(json, arg string) string {
	ctx := Parse(json)
	if !ctx.IsObject() {
		return ""
	}
	var all [][]byte
	ctx.Foreach(func(key, value Context) bool {
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
	for i, el := range all {
		if i > 0 {
			data = append(data, ',')
		}
		data = append(data, '{')
		data = append(data, el[1:]...)
		data = append(data, '}')
	}
	data = append(data, ']')
	return string(data)
}

// transformSearch performs a value lookup on a JSON structure based on the specified path
// and returns a JSON-encoded string containing all matching values found at that path.
//
// This function searches recursively through the JSON structure to find all occurrences
// of the specified path and then aggregates the results into a new JSON array. The results
// are returned as a string, which represents the matched values in their original, unprocessed form.
//
// Parameters:
//   - `json`: A string representing the input JSON data. The function will parse this JSON
//     and search for the specified path within the structure.
//   - `arg`: A string representing the JSON path to search for. The path is used to navigate
//     the nested JSON structure and retrieve values that match the specified key(s).
//
// Returns:
//   - A string representing a JSON array containing all matching values found in the input JSON
//     data. The values are presented in the same order they appear in the original JSON structure,
//     and are enclosed within square brackets ([]). If no matches are found, an empty array is returned.
//
// Example Usage:
//
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
//
//	arg := "book.author"
//	result := transformSearch(json, arg)
//
//	// Output: `["J.K. Rowling", "Stephen Hawking"]`
//	// The function will search for the "book.author" path and return all matching author names
//	// found in the nested `book` array.
//
// Notes:
//   - The `deepSearchRecursively` function is used to traverse the JSON structure and find all
//     matches for the specified path, ensuring that nested objects and arrays are searched as well.
//   - The results are accumulated in a slice, which is then converted into a JSON array and returned
//     as a string representation.
//   - If the path doesn't match any elements in the JSON structure, an empty array is returned.
//
// Implementation Details:
//   - The function utilizes `deepSearchRecursively` to perform a depth-first traversal of the JSON
//     structure and collect all matching values along the specified path.
//   - The results are then appended to a byte slice (`seg`), which is later converted into a string.
//   - The final output is a JSON array, even if no results are found.
func transformSearch(json, arg string) string {
	all := deepSearchRecursively(nil, Parse(json), arg)
	var seg []byte
	seg = append(seg, '[')
	for i, res := range all {
		if i > 0 {
			seg = append(seg, ',')
		}
		seg = append(seg, res.unprocessed...)
	}
	seg = append(seg, ']')
	return string(seg)
}

// transformUppercase converts the input JSON string to uppercase.
//
// This function takes a JSON string as input and converts all of its characters
// to uppercase. If the input string is empty, it returns the string unchanged.
//
// Parameters:
//   - `json`: The JSON string to be converted to uppercase.
//   - `arg`: An optional string parameter that is currently unused in this function,
//     but it could be extended for future use to modify the behavior of the transformation.
//
// Returns:
//   - A string with all characters converted to uppercase. If the input `json`
//     is empty, it returns the input string unchanged.
//
// Example Usage:
//
//	json := "{\"name\":\"Alice\",\"age\":25}"
//	result := transformUppercase(json, "")
//	fmt.Println(result) // Output: "{\"NAME\":\"ALICE\",\"AGE\":25}"
//
// Notes:
//   - This function uses the standard Go `strings.ToUpper` method to convert the string
//     to uppercase, which applies the transformation to every character in the string.
func transformUppercase(json, arg string) string {
	if unify4g.IsEmpty(json) {
		return json
	}
	return strings.ToUpper(json)
}

// transformLowercase converts the input JSON string to lowercase.
//
// This function takes a JSON string as input and converts all of its characters
// to lowercase. If the input string is empty, it returns the string unchanged.
//
// Parameters:
//   - `json`: The JSON string to be converted to lowercase.
//   - `arg`: An optional string parameter that is currently unused in this function,
//     but it could be extended for future use to modify the behavior of the transformation.
//
// Returns:
//   - A string with all characters converted to lowercase. If the input `json`
//     is empty, it returns the input string unchanged.
//
// Example Usage:
//
//	json := "{\"name\":\"Alice\",\"age\":25}"
//	result := transformLowercase(json, "")
//	fmt.Println(result) // Output: "{\"name\":\"alice\",\"age\":25}"
//
// Notes:
//   - This function uses the standard Go `strings.ToLower` method to convert the string
//     to lowercase, which applies the transformation to every character in the string.
func transformLowercase(json, arg string) string {
	if unify4g.IsEmpty(json) {
		return json
	}
	return strings.ToLower(json)
}

// transformFlip reverses the input JSON string.
//
// This function takes the input string and reverses the order of its characters.
// It returns the reversed string. If the input string is empty, it returns the
// string unchanged.
//
// Parameters:
//   - `json`: The JSON string to be reversed.
//   - `arg`: An optional argument that is currently unused, but could be extended
//     for future transformations.
//
// Returns:
//   - A string with its characters reversed. If the input is empty, the original
//     string is returned unchanged.
//
// Example Usage:
//
//	json := "{\"name\":\"Alice\",\"age\":25}"
//	result := transformReverse(json, "")
//	fmt.Println(result) // Output: "}52ega,\"ecilA\":\"emam{"
func transformFlip(json, arg string) string {
	if unify4g.IsEmpty(json) {
		return json
	}
	runes := []rune(json)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// transformTrim removes leading and trailing whitespace from the input JSON string.
//
// This function removes any whitespace characters at the beginning and end of the
// input string. If there is no whitespace, the string remains unchanged.
//
// Parameters:
//   - `json`: The JSON string to be trimmed of whitespace.
//   - `arg`: An optional string argument that is currently unused, but could be
//     extended for future use.
//
// Returns:
//   - A string with leading and trailing whitespace removed. If the input string
//     does not contain whitespace at the edges, it returns the original string.
//
// Example Usage:
//
//	json := "   {\"name\":\"Alice\"}   "
//	result := transformTrim(json, "")
//	fmt.Println(result) // Output: "{\"name\":\"Alice\"}"
//
// Notes:
//   - This function uses Go's `strings.TrimSpace` method to remove whitespace characters.
func transformTrim(json, arg string) string {
	return unify4g.TrimWhitespace(trim(json))
}

// transformSnakeCase converts the input string to snake_case format, which is typically used
// for variable names in many programming languages. The string is transformed to lowercase,
// and spaces or other delimiters are replaced with underscores ('_').
//
// Parameters:
//   - `json`: The input string that will be converted to snake_case.
//   - `arg`: An optional string argument that is currently unused.
//
// Returns:
//   - A string formatted in snake_case. If the input is empty, it returns unchanged.
//
// Example Usage:
//
//	json := "{\"First Name\":\"Alice\",\"Last Name\":\"Smith\"}"
//	result := transformSnakeCase(json, "")
//	fmt.Println(result) // Output: "{\"first_name\":\"alice\",\"last_name\":\"smith\"}"
func transformSnakeCase(json, arg string) string {
	if unify4g.IsEmpty(json) {
		return json
	}
	// Replace spaces with underscores, convert to lowercase.
	json = strings.ReplaceAll(json, " ", "_")
	return strings.ToLower(json)
}

// transformCamelCase converts the input string into camelCase, which is often used for
// variable and function names in JavaScript and other programming languages. The string
// is converted to lowercase with spaces removed, and the first letter of each word after
// the first is capitalized.
//
// Parameters:
//   - `json`: The input string that will be converted to camelCase.
//   - `arg`: An optional argument that is currently unused.
//
// Returns:
//   - A string formatted in camelCase. If the input is empty, it returns unchanged.
//
// Example Usage:
//
//	json := "{\"first name\":\"alice\",\"last name\":\"smith\"}"
//	result := transformCamelCase(json, "")
//	fmt.Println(result) // Output: "{\"firstName\":\"alice\",\"lastName\":\"smith\"}"
func transformCamelCase(json, arg string) string {
	if unify4g.IsEmpty(json) {
		return json
	}
	words := strings.Fields(json)
	for i := 1; i < len(words); i++ {
		words[i] = strings.Title(words[i])
	}
	return strings.Join(words, "")
}

// transformKebabCase converts the input string into kebab-case, often used for URL slugs.
// The string is transformed to lowercase, and spaces are replaced with hyphens ('-').
//
// Parameters:
//   - `json`: The input string that will be converted to kebab-case.
//   - `arg`: An optional string argument that is currently unused.
//
// Returns:
//   - A string formatted in kebab-case. If the input is empty, it returns unchanged.
//
// Example Usage:
//
//	json := "{\"First Name\":\"Alice\",\"Last Name\":\"Smith\"}"
//	result := transformKebabCase(json, "")
//	fmt.Println(result) // Output: "{\"first-name\":\"alice\",\"last-name\":\"smith\"}"
func transformKebabCase(json, arg string) string {
	if unify4g.IsEmpty(json) {
		return json
	}
	// Replace spaces with hyphens, convert to lowercase.
	json = strings.ReplaceAll(json, " ", "-")
	return strings.ToLower(json)
}

// transformReplaceSubstring replaces a specific substring within the input string with another string.
//
// This function can be used to replace any substring in the input string with a given replacement.
//
// Parameters:
//   - `json`: The input string in which the replacement will occur.
//   - `arg`: A JSON string containing the target and replacement substrings.
//     Example: `{"target": "apple", "replacement": "orange"}`.
//
// Returns:
//   - The string with the target substring replaced. If the input is empty, it returns unchanged.
//
// Example Usage:
//
//	json := "I love apple pie"
//	arg := "{\"target\":\"apple\",\"replacement\":\"orange\"}"
//	result := transformReplaceSubstring(json, arg)
//	fmt.Println(result) // Output: "I love orange pie"
func transformReplace(json, arg string) string {
	if unify4g.IsEmpty(json) {
		return json
	}
	var target, replacement string
	Parse(arg).Foreach(func(key, value Context) bool {
		if key.String() == "target" {
			target = value.String()
		}
		if key.String() == "replacement" {
			replacement = value.String()
		}
		return true
	})
	return strings.Replace(json, target, replacement, -1)
}

// transformReplaceAll replaces all occurrences of a target substring with a replacement string.
//
// This function performs a global replacement in the input string, replacing every occurrence
// of the target substring with the specified replacement.
//
// Parameters:
//   - `json`: The input string in which the replacement will be performed.
//   - `arg`: A JSON string containing the target and replacement substrings.
//     Example: `{"target": "foo", "replacement": "bar"}`.
//
// Returns:
//   - The string with all target substrings replaced by the replacement. If the input is empty,
//     it returns unchanged.
//
// Example Usage:
//
//	json := "foo bar foo"
//	arg := "{\"target\":\"foo\",\"replacement\":\"baz\"}"
//	result := transformReplaceAll(json, arg)
//	fmt.Println(result) // Output: "baz bar baz"
func transformReplaceAll(json, arg string) string {
	if unify4g.IsEmpty(json) {
		return json
	}
	var target, replacement string
	Parse(arg).Foreach(func(key, value Context) bool {
		if key.String() == "target" {
			target = value.String()
		}
		if key.String() == "replacement" {
			replacement = value.String()
		}
		return true
	})
	return strings.ReplaceAll(json, target, replacement)
}

// transformToHex converts the string to its hexadecimal representation.
//
// This function converts each character of the input string to its hexadecimal ASCII value.
//
// Parameters:
//   - `json`: The input string to be converted to hexadecimal.
//   - `arg`: An optional argument that is currently unused.
//
// Returns:
//   - A string containing the hexadecimal representation of the input string.
//
// Example Usage:
//
//	json := "hello"
//	result := transformToHex(json, "")
//	fmt.Println(result) // Output: "68656c6c6f"
func transformToHex(json, arg string) string {
	if unify4g.IsEmpty(json) {
		return json
	}
	ctx := Parse(json)
	if !ctx.IsArray() && !ctx.IsObject() && unify4g.IsPrimitive(ctx.String()) {
		return fmt.Sprintf("%x", removeDoubleQuotes(json))
	}
	return fmt.Sprintf("%x", json)
}

// transformToBinary converts the string to its binary representation.
//
// This function converts each character of the input string to its binary ASCII value.
//
// Parameters:
//   - `json`: The input string to be converted to binary.
//   - `arg`: An optional argument that is currently unused.
//
// Returns:
//   - A string containing the binary representation of the input string.
//
// Example Usage:
//
//	json := "hello"
//	result := transformToBinary(json, "")
//	fmt.Println(result) // Output: "11010001101101110110011011001101111"
func transformToBinary(json, arg string) string {
	ctx := Parse(json)
	if !ctx.IsArray() && !ctx.IsObject() && unify4g.IsPrimitive(ctx.String()) {
		var bin string
		for _, r := range ctx.String() {
			bin += fmt.Sprintf("%08b", r)
		}
		return bin
	}
	var bin string
	for _, r := range json {
		bin += fmt.Sprintf("%08b", r)
	}
	return bin
}

// transformInsertAt inserts a specified string at a given index in the input string.
//
// This function allows you to insert a string at a particular position in the original string.
//
// Parameters:
//   - `json`: The original string where the insertion will happen.
//   - `arg`: A JSON string containing the position (`index`) and the string to insert (`insert`).
//     Example: `{"index": 5, "insert": "XYZ"}`.
//
// Returns:
//   - A new string with the specified insertion. If the input is empty, it returns unchanged.
//
// Example Usage:
//
//	json := "HelloWorld"
//	arg := "{\"index\":5,\"insert\":\"XYZ\"}"
//	result := transformInsertAt(json, arg)
//	fmt.Println(result) // Output: "HelloXYZWorld"
func transformInsertAt(json, arg string) string {
	if unify4g.IsEmpty(json) {
		return json
	}
	var index int
	var insert string
	Parse(arg).Foreach(func(key, value Context) bool {
		if key.String() == "index" {
			index = int(value.Numeric())
		}
		if key.String() == "insert" {
			insert = value.String()
		}
		return true
	})
	if index < 0 || index > len(json) {
		return json // If index is out of bounds, return the original string.
	}
	return json[:index] + insert + json[index:]
}

// transformCountWords counts the number of words in the input string.
//
// This function splits the string into words (by spaces) and returns the count of the words.
//
// Parameters:
//   - `json`: The input string to count the words in.
//   - `arg`: An optional string argument that is currently unused.
//
// Returns:
//   - An integer representing the number of words in the input string. If the string is empty, it returns 0.
//
// Example Usage:
//
//	json := "Hello world"
//	result := transformCountWords(json, "")
//	fmt.Println(result) // Output: 2
func transformCountWords(json, arg string) string {
	ctx := Parse(json)
	json = unify4g.TrimWhitespace(ctx.String())
	if unify4g.IsEmpty(json) || unify4g.IsBlank(json) || unify4g.IsWhitespace(json) {
		return json
	}
	words := strings.Fields(json)
	return fmt.Sprintf("%v", len(words))
}

// transformPadLeft pads the input string with a specified character on the left to a given length.
//
// This function adds padding to the left of the string until the string reaches the specified length.
//
// Parameters:
//   - `json`: The input string to be padded.
//   - `arg`: A JSON string containing the padding character and the desired length. Example: `{"padding": "*", "length": 10}`.
//
// Returns:
//   - A new string with the specified padding added on the left. If the string is already the desired length, it returns unchanged.
//
// Example Usage:
//
//	json := "Hello"
//	arg := "{\"padding\": \"*\", \"length\": 10}"
//	result := transformPadLeft(json, arg)
//	fmt.Println(result) // Output: "*****Hello"
func transformPadLeft(json, arg string) string {
	ctx := Parse(json)
	value := ctx.String()
	var padding string
	var length int
	Parse(arg).Foreach(func(key, value Context) bool {
		if key.String() == "padding" {
			padding = value.String()
		}
		if key.String() == "length" {
			length = int(value.Int64())
		}
		return true
	})
	if length <= len(json) {
		return json
	}
	value = unify4g.TrimWhitespace(trim(value))
	t := unify4g.JsonN(value)
	v := strings.Repeat(padding, length-len(t)) + t
	return v
}

// transformPadRight pads the input string with a specified character on the right to a given length.
//
// This function adds padding to the right of the string until the string reaches the specified length.
//
// Parameters:
//   - `json`: The input string to be padded.
//   - `arg`: A JSON string containing the padding character and the desired length. Example: `{"padding": "*", "length": 10}`.
//
// Returns:
//   - A new string with the specified padding added on the right. If the string is already the desired length, it returns unchanged.
//
// Example Usage:
//
//	json := "Hello"
//	arg := "{\"padding\": \"*\", \"length\": 10}"
//	result := transformPadRight(json, arg)
//	fmt.Println(result) // Output: "Hello*****"
func transformPadRight(json, arg string) string {
	ctx := Parse(json)
	value := ctx.String()
	var padding string
	var length int
	Parse(arg).Foreach(func(key, value Context) bool {
		if key.String() == "padding" {
			padding = value.String()
		}
		if key.String() == "length" {
			length = int(value.Int64())
		}
		return true
	})
	if length <= len(json) {
		return json
	}
	value = unify4g.TrimWhitespace(trim(value))
	t := unify4g.JsonN(value)
	v := t + strings.Repeat(padding, length-len(t))
	return v
}
