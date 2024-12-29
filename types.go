package fj

import (
	"unsafe"
)

// Type represents the different possible types for a JSON value.
// It is used to indicate the specific type of a JSON value, such as a string, number, boolean, etc.
type Type int

// Context represents a JSON value returned from the Get() function.
// It stores information about a specific JSON element, including its type,
// unprocessed string data, string representation, numeric value, index in the original JSON,
// and the indexes of elements that match a path containing a '#'.
type Context struct {
	// kind is the JSON type (such as String, Number, Object, etc.).
	kind Type

	// unprocessed contains the raw JSON string that has not been processed or parsed.
	unprocessed string

	// strings contains the string value of the JSON element, if it is a string type.
	strings string

	// numeric contains the numeric value of the JSON element, if it is a number type.
	numeric float64

	// index holds the position of the unprocessed JSON value in the original JSON string.
	// A value of 0 means the index is unknown.
	index int

	// indexes holds the indices of all elements that match a path containing the '#' query character.
	indexes []int
}

// queryContext is a simplified version of the Context struct,
// primarily used to store intermediate results for JSON path processing.
type queryContext struct {
	// arrays stores a slice of Context elements representing an array result.
	arrays []Context

	// elements stores a slice of interface{} elements, used for intermediary results.
	elements []interface{}

	// operations maps string keys to Context values, used for handling operations on specific paths.
	operations map[string]Context

	// operationResults maps string keys to interface{} values, used for operation results on specific paths.
	operationResults map[string]interface{}

	// valueN stores a byte value for a specific operation, likely used for flagging or identifying states.
	valueN byte
}

// wildcard represents a wildcard element in a JSON path query.
// It is used to represent patterns that match various JSON values in path expressions.
type wildcard struct {
	// Part represents a specific segment of the wildcard pattern.
	Part string

	// Path represents the path expression that includes the wildcard.
	Path string

	// Pipe represents an operation or transformation to be applied to the wildcard.
	Pipe string

	// Piped indicates whether the wildcard has been piped for further processing.
	Piped bool

	// Wild indicates whether this wildcard is a true wildcard (e.g., "*").
	Wild bool

	// More indicates whether there are more operations or elements to process after this wildcard.
	More bool
}

// metadata represents a metadata, more complex structure for handling JSON path queries
// that may involve operations, conditions, and logging.
type metadata struct {
	// Part represents a segment of the deeper query structure.
	Part string

	// Path represents the full path for the deeper query structure.
	Path string

	// Pipe represents any operations or transformations applied to the deeper query structure.
	Pipe string

	// Piped indicates if the deeper query structure has been piped for further processing.
	Piped bool

	// More indicates whether there are additional elements or operations in the deeper query.
	More bool

	// Arch is a flag indicating some form of structure or architecture flag.
	Arch bool

	// ALogOk indicates whether logging is allowed for this deeper structure.
	ALogOk bool

	// ALogKey stores the key used for logging in the deeper structure.
	ALogKey string

	// query holds the query-related data for the deeper structure.
	query struct {
		// On is a flag that indicates whether the query is active.
		On bool

		// All indicates if the query applies to all elements, not just a single one.
		All bool

		// QueryPath is the path to query.
		QueryPath string

		// Option specifies an option related to the query.
		Option string

		// Value represents the value to search for or match in the query.
		Value string
	}
}

// parser holds the state and configuration for parsing JSON data.
type parser struct {
	// json is the raw JSON string that needs to be parsed.
	json string

	// value stores the parsed context value that represents a JSON element.
	value Context

	// pipe holds the pipe operation string to be applied during parsing.
	pipe string

	// piped indicates whether the value has been piped for additional processing.
	piped bool

	// calc indicates whether calculations or transformations should be performed during parsing.
	calc bool

	// lines indicates whether the JSON data should be processed line by line.
	lines bool
}

// stringHeader is a custom struct that mimics the reflect.stringHeader type
// and is used for low-level manipulation of string data in Go.
// It provides access to the internal representation of a string.
type stringHeader struct {
	data   unsafe.Pointer // data is a pointer to the underlying data of the string.
	length int            // length is the length of the string.
}

// sliceHeader is a custom struct that mimics the reflect.sliceHeader type
// and is used for low-level manipulation of slices in Go.
// It provides access to the internal representation of a slice.
type sliceHeader struct {
	data     unsafe.Pointer // data is a pointer to the underlying array of the slice.
	length   int            // length is the number of elements in the slice.
	capacity int            // capacity is the maximum number of elements the slice can hold before resizing.
}

// subSelector represents a selection in a JSON path query that specifies a name and path.
type subSelector struct {
	name string // name represents the name of the selector or key in the JSON path.
	path string //  path represents the full path expression for the selector.
}
