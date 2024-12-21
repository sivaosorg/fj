package fj

import (
	"unsafe"
)

type Type int

// Context represents a json value that is returned from Get().
type Context struct {
	// kind is the json type
	kind Type
	// unprocessed is the unprocessed json
	unprocessed string
	// strings is the json string
	strings string
	// numeric is the json number
	numeric float64
	// index of unprocessed value in original json, zero means index unknown
	index int
	// indexes of all the elements that match on a path containing the '#'
	// query character.
	indexes []int
}

type tinyContext struct {
	ArrayResult []Context              `json:"-"`
	ArrayIns    []interface{}          `json:"-"`
	OpMap       map[string]Context     `json:"-"`
	OpIns       map[string]interface{} `json:"-"`
	valueN      byte                   `json:"-"`
}

type wildcard struct {
	Part  string `json:"-"`
	Path  string `json:"-"`
	Pipe  string `json:"-"`
	Piped bool   `json:"-"`
	Wild  bool   `json:"-"`
	More  bool   `json:"-"`
}

type deeper struct {
	Part    string `json:"-"`
	Path    string `json:"-"`
	Pipe    string `json:"-"`
	Piped   bool   `json:"-"`
	More    bool   `json:"-"`
	Arch    bool   `json:"-"`
	ALogOk  bool   `json:"-"`
	ALogKey string `json:"-"`
	query   struct {
		On        bool   `json:"-"`
		All       bool   `json:"-"`
		QueryPath string `json:"-"`
		Option    string `json:"-"`
		Value     string `json:"-"`
	} `json:"-"`
}

type parser struct {
	json  string
	value Context
	pipe  string
	piped bool
	calc  bool
	lines bool
}

// stringHeader instead of reflect.stringHeader
type stringHeader struct {
	data   unsafe.Pointer
	length int
}

// sliceHeader instead of reflect.sliceHeader
type sliceHeader struct {
	data     unsafe.Pointer
	length   int
	capacity int
}

type subSelector struct {
	name string
	path string
}
