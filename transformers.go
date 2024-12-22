package fj

import (
	"github.com/sivaosorg/unify4g"
)

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
