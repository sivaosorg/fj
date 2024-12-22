package fj

import (
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
	if unify4g.IsNoneEmpty(arg) {
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
