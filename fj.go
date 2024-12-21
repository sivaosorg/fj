package fj

import (
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/sivaosorg/unify4g"
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

func (c Context) Kind() Type {
	return c.kind
}

func (c Context) Unprocessed() string {
	return c.unprocessed
}

func (c Context) Numeric() float64 {
	return c.numeric
}

func (c Context) Index() int {
	return c.index
}

func (c Context) Indexes() []int {
	return c.indexes
}

func (t Type) String() string {
	switch t {
	default:
		return ""
	case Null:
		return "Null"
	case False:
		return "False"
	case Number:
		return "Number"
	case String:
		return "String"
	case True:
		return "True"
	case JSON:
		return "JSON"
	}
}

// String returns a string representation of the value.
func (t Context) String() string {
	switch t.kind {
	default:
		return ""
	case False:
		return "false"
	case Number:
		if len(t.unprocessed) == 0 {
			// calculated result
			return strconv.FormatFloat(t.numeric, 'f', -1, 64)
		}
		var i int
		if t.unprocessed[0] == '-' {
			i++
		}
		for ; i < len(t.unprocessed); i++ {
			if t.unprocessed[i] < '0' || t.unprocessed[i] > '9' {
				return strconv.FormatFloat(t.numeric, 'f', -1, 64)
			}
		}
		return t.unprocessed
	case String:
		return t.strings
	case JSON:
		return t.unprocessed
	case True:
		return "true"
	}
}

// Bool returns an boolean representation.
func (t Context) Bool() bool {
	switch t.kind {
	default:
		return false
	case True:
		return true
	case String:
		b, _ := strconv.ParseBool(strings.ToLower(t.strings))
		return b
	case Number:
		return t.numeric != 0
	}
}

// Int returns an integer representation.
func (t Context) Int() int64 {
	switch t.kind {
	default:
		return 0
	case True:
		return 1
	case String:
		n, _ := parseInt(t.strings)
		return n
	case Number:
		// try to directly convert the float64 to int64
		i, ok := ensureSafeInt(t.numeric)
		if ok {
			return i
		}
		// now try to parse the raw string
		i, ok = parseInt(t.unprocessed)
		if ok {
			return i
		}
		// fallback to a standard conversion
		return int64(t.numeric)
	}
}

// Uint returns an unsigned integer representation.
func (t Context) Uint() uint64 {
	switch t.kind {
	default:
		return 0
	case True:
		return 1
	case String:
		n, _ := parseUint(t.strings)
		return n
	case Number:
		// try to directly convert the float64 to uint64
		i, ok := ensureSafeInt(t.numeric)
		if ok && i >= 0 {
			return uint64(i)
		}
		// now try to parse the raw string
		u, ok := parseUint(t.unprocessed)
		if ok {
			return u
		}
		// fallback to a standard conversion
		return uint64(t.numeric)
	}
}

// Float returns an float64 representation.
func (t Context) Float() float64 {
	switch t.kind {
	default:
		return 0
	case True:
		return 1
	case String:
		n, _ := strconv.ParseFloat(t.strings, 64)
		return n
	case Number:
		return t.numeric
	}
}

// Time returns a time.Time representation.
func (t Context) Time() time.Time {
	res, _ := time.Parse(time.RFC3339, t.String())
	return res
}

// Array returns back an array of values.
// If the result represents a null value or is non-existent, then an empty
// array will be returned.
// If the result is not a JSON array, the return value will be an
// array containing one result.
func (t Context) Array() []Context {
	if t.kind == Null {
		return []Context{}
	}
	if !t.IsArray() {
		return []Context{t}
	}
	r := t.arrayOrMap('[', false)
	return r.ArrayResult
}

// IsObject returns true if the result value is a JSON object.
func (t Context) IsObject() bool {
	return t.kind == JSON && len(t.unprocessed) > 0 && t.unprocessed[0] == '{'
}

// IsArray returns true if the result value is a JSON array.
func (t Context) IsArray() bool {
	return t.kind == JSON && len(t.unprocessed) > 0 && t.unprocessed[0] == '['
}

// IsBool returns true if the result value is a JSON boolean.
func (t Context) IsBool() bool {
	return t.kind == True || t.kind == False
}

// ForEach iterates through values.
// If the result represents a non-existent value, then no values will be
// iterated. If the result is an Object, the iterator will pass the key and
// value of each item. If the result is an Array, the iterator will only pass
// the value of each item. If the result is not a JSON array or object, the
// iterator will pass back one value equal to the result.
func (t Context) ForEach(iterator func(key, value Context) bool) {
	if !t.Exists() {
		return
	}
	if t.kind != JSON {
		iterator(Context{}, t)
		return
	}
	json := t.unprocessed
	var obj bool
	var i int
	var key, value Context
	for ; i < len(json); i++ {
		if json[i] == '{' {
			i++
			key.kind = String
			obj = true
			break
		} else if json[i] == '[' {
			i++
			key.kind = Number
			key.numeric = -1
			break
		}
		if json[i] > ' ' {
			return
		}
	}
	var str string
	var _esc bool
	var ok bool
	var idx int
	for ; i < len(json); i++ {
		if obj {
			if json[i] != '"' {
				continue
			}
			s := i
			i, str, _esc, ok = parseString(json, i+1)
			if !ok {
				return
			}
			if _esc {
				key.strings = unescape(str[1 : len(str)-1])
			} else {
				key.strings = str[1 : len(str)-1]
			}
			key.unprocessed = str
			key.index = s + t.index
		} else {
			key.numeric += 1
		}
		for ; i < len(json); i++ {
			if json[i] <= ' ' || json[i] == ',' || json[i] == ':' {
				continue
			}
			break
		}
		s := i
		i, value, ok = parseJSONAny(json, i, true)
		if !ok {
			return
		}
		if t.indexes != nil {
			if idx < len(t.indexes) {
				value.index = t.indexes[idx]
			}
		} else {
			value.index = s + t.index
		}
		if !iterator(key, value) {
			return
		}
		idx++
	}
}

// Map returns back a map of values. The result should be a JSON object.
// If the result is not a JSON object, the return value will be an empty map.
func (t Context) Map() map[string]Context {
	if t.kind != JSON {
		return map[string]Context{}
	}
	r := t.arrayOrMap('{', false)
	return r.OpMap
}

// Get searches result for the specified path.
// The result should be a JSON array or object.
func (t Context) Get(path string) Context {
	r := Get(t.unprocessed, path)
	if r.indexes != nil {
		for i := 0; i < len(r.indexes); i++ {
			r.indexes[i] += t.index
		}
	} else {
		r.index += t.index
	}
	return r
}

func (t Context) arrayOrMap(vc byte, valueSize bool) (result tinyContext) {
	var json = t.unprocessed
	var i int
	var value Context
	var count int
	var key Context
	if vc == 0 {
		for ; i < len(json); i++ {
			if json[i] == '{' || json[i] == '[' {
				result.valueN = json[i]
				i++
				break
			}
			if json[i] > ' ' {
				goto end
			}
		}
	} else {
		for ; i < len(json); i++ {
			if json[i] == vc {
				i++
				break
			}
			if json[i] > ' ' {
				goto end
			}
		}
		result.valueN = vc
	}
	if result.valueN == '{' {
		if valueSize {
			result.OpIns = make(map[string]interface{})
		} else {
			result.OpMap = make(map[string]Context)
		}
	} else {
		if valueSize {
			result.ArrayIns = make([]interface{}, 0)
		} else {
			result.ArrayResult = make([]Context, 0)
		}
	}
	for ; i < len(json); i++ {
		if json[i] <= ' ' {
			continue
		}
		// get next value
		if json[i] == ']' || json[i] == '}' {
			break
		}
		switch json[i] {
		default:
			if (json[i] >= '0' && json[i] <= '9') || json[i] == '-' {
				value.kind = Number
				value.unprocessed, value.numeric = getNumeric(json[i:])
				value.strings = ""
			} else {
				continue
			}
		case '{', '[':
			value.kind = JSON
			value.unprocessed = squash(json[i:])
			value.strings, value.numeric = "", 0
		case 'n':
			value.kind = Null
			value.unprocessed = lowerPrefix(json[i:])
			value.strings, value.numeric = "", 0
		case 't':
			value.kind = True
			value.unprocessed = lowerPrefix(json[i:])
			value.strings, value.numeric = "", 0
		case 'f':
			value.kind = False
			value.unprocessed = lowerPrefix(json[i:])
			value.strings, value.numeric = "", 0
		case '"':
			value.kind = String
			value.unprocessed, value.strings = unescapeJsonEncoded(json[i:])
			value.numeric = 0
		}
		value.index = i + t.index

		i += len(value.unprocessed) - 1

		if result.valueN == '{' {
			if count%2 == 0 {
				key = value
			} else {
				if valueSize {
					if _, ok := result.OpIns[key.strings]; !ok {
						result.OpIns[key.strings] = value.Value()
					}
				} else {
					if _, ok := result.OpMap[key.strings]; !ok {
						result.OpMap[key.strings] = value
					}
				}
			}
			count++
		} else {
			if valueSize {
				result.ArrayIns = append(result.ArrayIns, value.Value())
			} else {
				result.ArrayResult = append(result.ArrayResult, value)
			}
		}
	}
end:
	if t.indexes != nil {
		if len(t.indexes) != len(result.ArrayResult) {
			for i := 0; i < len(result.ArrayResult); i++ {
				result.ArrayResult[i].index = 0
			}
		} else {
			for i := 0; i < len(result.ArrayResult); i++ {
				result.ArrayResult[i].index = t.indexes[i]
			}
		}
	}
	return
}

// Parse parses the json and returns a result.
//
// This function expects that the json is well-formed, and does not validate.
// Invalid json will not panic, but it may return back unexpected results.
// If you are consuming JSON from an unpredictable source then you may want to
// use the Valid function first.
func Parse(json string) Context {
	var value Context
	i := 0
	for ; i < len(json); i++ {
		if json[i] == '{' || json[i] == '[' {
			value.kind = JSON
			value.unprocessed = json[i:] // just take the entire raw
			break
		}
		if json[i] <= ' ' {
			continue
		}
		switch json[i] {
		case '+', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'i', 'I', 'N':
			value.kind = Number
			value.unprocessed, value.numeric = getNumeric(json[i:])
		case 'n':
			if i+1 < len(json) && json[i+1] != 'u' {
				// nan
				value.kind = Number
				value.unprocessed, value.numeric = getNumeric(json[i:])
			} else {
				// null
				value.kind = Null
				value.unprocessed = lowerPrefix(json[i:])
			}
		case 't':
			value.kind = True
			value.unprocessed = lowerPrefix(json[i:])
		case 'f':
			value.kind = False
			value.unprocessed = lowerPrefix(json[i:])
		case '"':
			value.kind = String
			value.unprocessed, value.strings = unescapeJsonEncoded(json[i:])
		default:
			return Context{}
		}
		break
	}
	if value.Exists() {
		value.index = i
	}
	return value
}

// ParseBytes parses the json and returns a result.
// If working with bytes, this method preferred over Parse(string(data))
func ParseBytes(json []byte) Context {
	return Parse(string(json))
}

// Exists returns true if value exists.
//
//	 if bjson.Get(json, "name.last").Exists(){
//			println("value exists")
//	 }
func (t Context) Exists() bool {
	return t.kind != Null || len(t.unprocessed) != 0
}

// Value returns one of these types:
//
//	bool, for JSON booleans
//	float64, for JSON numbers
//	Number, for JSON numbers
//	string, for JSON string literals
//	nil, for JSON null
//	map[string]interface{}, for JSON objects
//	[]interface{}, for JSON arrays
func (t Context) Value() interface{} {
	if t.kind == String {
		return t.strings
	}
	switch t.kind {
	default:
		return nil
	case False:
		return false
	case Number:
		return t.numeric
	case JSON:
		r := t.arrayOrMap(0, true)
		if r.valueN == '{' {
			return r.OpIns
		} else if r.valueN == '[' {
			return r.ArrayIns
		}
		return nil
	case True:
		return true
	}
}

func ofFalse(t Context) bool {
	switch t.kind {
	case Null:
		return true
	case False:
		return true
	case String:
		b, err := strconv.ParseBool(strings.ToLower(t.strings))
		if err != nil {
			return false
		}
		return !b
	case Number:
		return t.numeric == 0
	default:
		return false
	}
}

func ofTrue(t Context) bool {
	switch t.kind {
	case True:
		return true
	case String:
		b, err := strconv.ParseBool(strings.ToLower(t.strings))
		if err != nil {
			return false
		}
		return b
	case Number:
		return t.numeric != 0
	default:
		return false
	}
}

func nullish(t Context) bool {
	return t.kind == Null
}

func queryMatches(dp *deeper, value Context) bool {
	rpv := dp.query.Value
	if len(rpv) > 0 {
		if rpv[0] == '~' {
			// convert to bool
			rpv = rpv[1:]
			var ish, ok bool
			switch rpv {
			case "*":
				ish, ok = value.Exists(), true
			case "null":
				ish, ok = nullish(value), true
			case "true":
				ish, ok = ofTrue(value), true
			case "false":
				ish, ok = ofFalse(value), true
			}
			if ok {
				rpv = "true"
				if ish {
					value = Context{kind: True}
				} else {
					value = Context{kind: False}
				}
			} else {
				rpv = ""
				value = Context{}
			}
		}
	}
	if !value.Exists() {
		return false
	}
	if dp.query.Option == "" {
		// the query is only looking for existence, such as:
		//   friends.#(name)
		// which makes sure that the array "friends" has an element of
		// "name" that exists
		return true
	}
	switch value.kind {
	case String:
		switch dp.query.Option {
		case "=":
			return value.strings == rpv
		case "!=":
			return value.strings != rpv
		case "<":
			return value.strings < rpv
		case "<=":
			return value.strings <= rpv
		case ">":
			return value.strings > rpv
		case ">=":
			return value.strings >= rpv
		case "%":
			return matchSafely(value.strings, rpv)
		case "!%":
			return !matchSafely(value.strings, rpv)
		}
	case Number:
		_rightVal, _ := strconv.ParseFloat(rpv, 64)
		switch dp.query.Option {
		case "=":
			return value.numeric == _rightVal
		case "!=":
			return value.numeric != _rightVal
		case "<":
			return value.numeric < _rightVal
		case "<=":
			return value.numeric <= _rightVal
		case ">":
			return value.numeric > _rightVal
		case ">=":
			return value.numeric >= _rightVal
		}
	case True:
		switch dp.query.Option {
		case "=":
			return rpv == "true"
		case "!=":
			return rpv != "true"
		case ">":
			return rpv == "false"
		case ">=":
			return true
		}
	case False:
		switch dp.query.Option {
		case "=":
			return rpv == "false"
		case "!=":
			return rpv != "false"
		case "<":
			return rpv == "true"
		case "<=":
			return true
		}
	}
	return false
}

// ForEachLine iterates through lines of JSON as specified by the JSON Lines
// format (http://jsonlines.org/).
// Each line is returned as a bjson Result.
func ForEachLine(json string, iterator func(line Context) bool) {
	var res Context
	var i int
	for {
		i, res, _ = parseJSONAny(json, i, true)
		if !res.Exists() {
			break
		}
		if !iterator(res) {
			return
		}
	}
}

// AppendJsonString is a convenience function that converts the provided string
// to a valid JSON string and appends it to dst.
func AppendJsonString(dst []byte, s string) []byte {
	dst = append(dst, make([]byte, len(s)+2)...)
	dst = append(dst[:len(dst)-len(s)-2], '"')
	for i := 0; i < len(s); i++ {
		if s[i] < ' ' {
			dst = append(dst, '\\')
			switch s[i] {
			case '\n':
				dst = append(dst, 'n')
			case '\r':
				dst = append(dst, 'r')
			case '\t':
				dst = append(dst, 't')
			default:
				dst = append(dst, 'u')
				dst = appendHex16(dst, uint16(s[i]))
			}
		} else if s[i] == '>' || s[i] == '<' || s[i] == '&' {
			dst = append(dst, '\\', 'u')
			dst = appendHex16(dst, uint16(s[i]))
		} else if s[i] == '\\' {
			dst = append(dst, '\\', '\\')
		} else if s[i] == '"' {
			dst = append(dst, '\\', '"')
		} else if s[i] > 127 {
			// read utf8 character
			r, n := utf8.DecodeRuneInString(s[i:])
			if n == 0 {
				break
			}
			if r == utf8.RuneError && n == 1 {
				dst = append(dst, `\ufffd`...)
			} else if r == '\u2028' || r == '\u2029' {
				dst = append(dst, `\u202`...)
				dst = append(dst, hexDigits[r&0xF])
			} else {
				dst = append(dst, s[i:i+n]...)
			}
			i = i + n - 1
		} else {
			dst = append(dst, s[i])
		}
	}
	return append(dst, '"')
}

// Get searches json for the specified path.
// A path is in dot syntax, such as "name.last" or "age".
// When the value is found it's returned immediately.
//
// A path is a series of keys separated by a dot.
// A key may contain special wildcard characters '*' and '?'.
// To access an array value use the index as the key.
// To get the number of elements in an array or to access a child path, use
// the '#' character.
// The dot and wildcard character can be escaped with '\'.
//
//	{
//	  "name": {"first": "Tom", "last": "Anderson"},
//	  "age":37,
//	  "children": ["Sara","Alex","Jack"],
//	  "friends": [
//	    {"first": "James", "last": "Murphy"},
//	    {"first": "Roger", "last": "Craig"}
//	  ]
//	}
//	"name.last"          >> "Anderson"
//	"age"                >> 37
//	"children"           >> ["Sara","Alex","Jack"]
//	"children.#"         >> 3
//	"children.1"         >> "Alex"
//	"child*.2"           >> "Jack"
//	"children.0"         >> "Sara"
//	"friends.#.first"    >> ["James","Roger"]
//
// This function expects that the json is well-formed, and does not validate.
// Invalid json will not panic, but it may return back unexpected results.
// If you are consuming JSON from an unpredictable source then you may want to
// use the Valid function first.
func Get(json, path string) Context {
	if len(path) > 1 {
		if (path[0] == '@' && !DisableModifiers) || path[0] == '!' {
			// possible modifier
			var ok bool
			var cPath string
			var cJson string
			if path[0] == '@' && !DisableModifiers {
				cPath, cJson, ok = adjustModifier(json, path)
			} else if path[0] == '!' {
				cPath, cJson, ok = parseStaticValue(path)
			}
			if ok {
				path = cPath
				if len(path) > 0 && (path[0] == '|' || path[0] == '.') {
					res := Get(cJson, path[1:])
					res.index = 0
					res.indexes = nil
					return res
				}
				return Parse(cJson)
			}
		}
		if path[0] == '[' || path[0] == '{' {
			// using a sub-selector path
			kind := path[0]
			var ok bool
			var subs []subSelector
			subs, path, ok = analyzeSubSelectors(path)
			if ok {
				if len(path) == 0 || (path[0] == '|' || path[0] == '.') {
					var b []byte
					b = append(b, kind)
					var i int
					for _, sub := range subs {
						res := Get(json, sub.path)
						if res.Exists() {
							if i > 0 {
								b = append(b, ',')
							}
							if kind == '{' {
								if len(sub.name) > 0 {
									if sub.name[0] == '"' && Valid(sub.name) {
										b = append(b, sub.name...)
									} else {
										b = AppendJsonString(b, sub.name)
									}
								} else {
									last := lastSegment(sub.path)
									if isValidName(last) {
										b = AppendJsonString(b, last)
									} else {
										b = AppendJsonString(b, "_")
									}
								}
								b = append(b, ':')
							}
							var raw string
							if len(res.unprocessed) == 0 {
								raw = res.String()
								if len(raw) == 0 {
									raw = "null"
								}
							} else {
								raw = res.unprocessed
							}
							b = append(b, raw...)
							i++
						}
					}
					b = append(b, kind+2)
					var res Context
					res.unprocessed = string(b)
					res.kind = JSON
					if len(path) > 0 {
						res = res.Get(path[1:])
					}
					res.index = 0
					return res
				}
			}
		}
	}
	var i int
	var c = &parser{json: json}
	if len(path) >= 2 && path[0] == '.' && path[1] == '.' {
		c.lines = true
		analyzeArray(c, 0, path[2:])
	} else {
		for ; i < len(c.json); i++ {
			if c.json[i] == '{' {
				i++
				parseJsonObject(c, i, path)
				break
			}
			if c.json[i] == '[' {
				i++
				analyzeArray(c, i, path)
				break
			}
		}
	}
	if c.piped {
		res := c.value.Get(c.pipe)
		res.index = 0
		return res
	}
	calcSubstring(json, c)
	return c.value
}

// GetBytes searches the provided JSON byte slice for the specified path and returns a `Context`
// representing the extracted data. This method is preferred over `Get(string(data), path)` when working
// with JSON data in byte slice format, as it directly operates on the byte slice, minimizing memory
// allocations and unnecessary copies.
//
// Parameters:
//   - `json`: A byte slice containing the JSON data to process.
//   - `path`: A string representing the path in the JSON data to extract.
//
// Returns:
//   - A `Context` struct containing the processed JSON data. The `Context` struct includes both
//     the raw unprocessed JSON string and the specific extracted string based on the given path.
//
// Notes:
//   - This function internally calls the `getBytes` function, which uses unsafe pointer operations
//     to minimize allocations and efficiently handle string slice headers.
//   - The function avoids unnecessary memory allocations by directly processing the byte slice and
//     utilizing memory safety features to manage substring extraction when the `strings` part is
//     a substring of the `unprocessed` part of the JSON data.
//
// Example:
//
//	jsonBytes := []byte(`{"key": "value", "nested": {"innerKey": "innerValue"}}`)
//	path := "nested.innerKey"
//	context := GetBytes(jsonBytes, path)
//	fmt.Println("Unprocessed:", context.unprocessed) // Output: `{"key": "value", "nested": {"innerKey": "innerValue"}}`
//	fmt.Println("Strings:", context.strings)         // Output: `"innerValue"`
func GetBytes(json []byte, path string) Context {
	return getBytes(json, path)
}

// Less compares two Context values (tokens) and returns true if the first token is considered less than the second one.
// It performs comparisons based on the type of the tokens and their respective values.
// The comparison order follows: Null < False < Number < String < True < JSON.
// This function also supports case-insensitive comparisons for String type tokens based on the caseSensitive parameter.
//
// Parameters:
//   - token: The other Context token to compare with the current one (t).
//   - caseSensitive: A boolean flag that indicates whether the comparison for String type tokens should be case-sensitive.
//   - If true, the comparison is case-sensitive (i.e., "a" < "b" but "A" < "b").
//   - If false, the comparison is case-insensitive (i.e., "a" == "A").
//
// Returns:
//   - true: If the current token (t) is considered less than the provided token.
//   - false: If the current token (t) is not considered less than the provided token.
//
// The function first compares the `kind` of both tokens, which represents their JSON types.
// If both tokens have the same kind, it proceeds to compare based on their specific types:
// - For String types, it compares the strings based on the case-sensitive flag.
// - For Number types, it compares the numeric values directly.
// - For other types, it compares the unprocessed JSON values as raw strings (this could be useful for types like Null, Boolean, etc.).
//
// Example usage:
//
//	context1 := Context{kind: String, strings: "apple"}
//	context2 := Context{kind: String, strings: "banana"}
//	result := context1.Less(context2, true) // This would return true because "apple" < "banana" and case-sensitive comparison is used.
func (t Context) Less(token Context, caseSensitive bool) bool {
	if t.kind < token.kind {
		return true
	}
	if t.kind > token.kind {
		return false
	}
	if t.kind == String {
		if caseSensitive {
			return t.strings < token.strings
		}
		return lessInsensitive(t.strings, token.strings)
	}
	if t.kind == Number {
		return t.numeric < token.numeric
	}
	return t.unprocessed < token.unprocessed
}

// GetMany searches json for the multiple paths.
// The return value is a Result array where the number of items
// will be equal to the number of input paths.
func GetMany(json string, path ...string) []Context {
	res := make([]Context, len(path))
	for i, path := range path {
		res[i] = Get(json, path)
	}
	return res
}

// GetManyBytes searches json for the multiple paths.
// The return value is a Result array where the number of items
// will be equal to the number of input paths.
func GetManyBytes(json []byte, path ...string) []Context {
	res := make([]Context, len(path))
	for i, path := range path {
		res[i] = GetBytes(json, path)
	}
	return res
}

// Valid returns true if the input is valid json.
//
//	if !bjson.Valid(json) {
//		return errors.New("invalid json")
//	}
//	value := bjson.Get(json, "name.last")
func Valid(json string) bool {
	_, ok := verifyJson(fromStr2Bytes(json), 0)
	return ok
}

// ValidBytes returns true if the input is valid json.
//
//	if !bjson.Valid(json) {
//		return errors.New("invalid json")
//	}
//	value := bjson.Get(json, "name.last")
//
// If working with bytes, this method preferred over ValidBytes(string(data))
func ValidBytes(json []byte) bool {
	_, ok := verifyJson(json, 0)
	return ok
}

func init() {
	modifiers = map[string]func(json, arg string) string{
		"pretty":  modPretty,
		"ugly":    modUgly,
		"reverse": modReverse,
		"this":    modThis,
		"flatten": modFlatten,
		"join":    modJoin,
		"valid":   modValid,
		"keys":    modKeys,
		"values":  modValues,
		"tostr":   modToStr,
		"fromstr": modFromStr,
		"group":   modGroup,
		"dig":     modDig,
	}
}

// AddModifier binds a custom modifier command to the bjson syntax.
// This operation is not thread safe and should be executed prior to
// using all other bjson function.
func AddModifier(name string, fn func(json, arg string) string) {
	modifiers[name] = fn
}

// ModifierExists returns true when the specified modifier exists.
func ModifierExists(name string, fn func(json, arg string) string) bool {
	_, ok := modifiers[name]
	return ok
}

// @pretty modifier makes the json look nice.
func modPretty(json, arg string) string {
	if len(arg) > 0 {
		opts := *unify4g.DefaultOptionsConfig
		Parse(arg).ForEach(func(key, value Context) bool {
			switch key.String() {
			case "sortKeys":
				opts.SortKeys = value.Bool()
			case "indent":
				opts.Indent = stripNonWhitespace(value.String())
			case "prefix":
				opts.Prefix = stripNonWhitespace(value.String())
			case "width":
				opts.Width = int(value.Int())
			}
			return true
		})
		return fromBytes2Str(unify4g.PrettyOptions(fromStr2Bytes(json), &opts))
	}
	return fromBytes2Str(unify4g.Pretty(fromStr2Bytes(json)))
}

// @this returns the current element. Can be used to retrieve the root element.
func modThis(json, arg string) string {
	return json
}

// @ugly modifier removes all whitespace.
func modUgly(json, arg string) string {
	return fromBytes2Str(unify4g.Ugly(fromStr2Bytes(json)))
}

// @reverse reverses array elements or root object members.
func modReverse(json, arg string) string {
	res := Parse(json)
	if res.IsArray() {
		var values []Context
		res.ForEach(func(_, value Context) bool {
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
	if res.IsObject() {
		var keyValues []Context
		res.ForEach(func(key, value Context) bool {
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

// @flatten an array with child arrays.
//
//	[1,[2],[3,4],[5,[6,7]]] -> [1,2,3,4,5,[6,7]]
//
// The {"deep":true} arg can be provide for deep flattening.
//
//	[1,[2],[3,4],[5,[6,7]]] -> [1,2,3,4,5,6,7]
//
// The original json is returned when the json is not an array.
func modFlatten(json, arg string) string {
	res := Parse(json)
	if !res.IsArray() {
		return json
	}
	var deep bool
	if arg != "" {
		Parse(arg).ForEach(func(key, value Context) bool {
			if key.String() == "deep" {
				deep = value.Bool()
			}
			return true
		})
	}
	var out []byte
	out = append(out, '[')
	var idx int
	res.ForEach(func(_, value Context) bool {
		var raw string
		if value.IsArray() {
			if deep {
				raw = removeOuterBraces(modFlatten(value.unprocessed, arg))
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

// @keys extracts the keys from an object.
//
//	{"first":"Tom","last":"Smith"} -> ["first","last"]
func modKeys(json, arg string) string {
	v := Parse(json)
	if !v.Exists() {
		return "[]"
	}
	obj := v.IsObject()
	var out strings.Builder
	out.WriteByte('[')
	var i int
	v.ForEach(func(key, _ Context) bool {
		if i > 0 {
			out.WriteByte(',')
		}
		if obj {
			out.WriteString(key.unprocessed)
		} else {
			out.WriteString("null")
		}
		i++
		return true
	})
	out.WriteByte(']')
	return out.String()
}

// @values extracts the values from an object.
//
//	{"first":"Tom","last":"Smith"} -> ["Tom","Smith"]
func modValues(json, arg string) string {
	v := Parse(json)
	if !v.Exists() {
		return "[]"
	}
	if v.IsArray() {
		return json
	}
	var out strings.Builder
	out.WriteByte('[')
	var i int
	v.ForEach(func(_, value Context) bool {
		if i > 0 {
			out.WriteByte(',')
		}
		out.WriteString(value.unprocessed)
		i++
		return true
	})
	out.WriteByte(']')
	return out.String()
}

// @join multiple objects into a single object.
//
//	[{"first":"Tom"},{"last":"Smith"}] -> {"first","Tom","last":"Smith"}
//
// The arg can be "true" to specify that duplicate keys should be preserved.
//
//	[{"first":"Tom","age":37},{"age":41}] -> {"first","Tom","age":37,"age":41}
//
// Without preserved keys:
//
//	[{"first":"Tom","age":37},{"age":41}] -> {"first","Tom","age":41}
//
// The original json is returned when the json is not an object.
func modJoin(json, arg string) string {
	res := Parse(json)
	if !res.IsArray() {
		return json
	}
	var preserve bool
	if arg != "" {
		Parse(arg).ForEach(func(key, value Context) bool {
			if key.String() == "preserve" {
				preserve = value.Bool()
			}
			return true
		})
	}
	var out []byte
	out = append(out, '{')
	if preserve {
		// Preserve duplicate keys.
		var idx int
		res.ForEach(func(_, value Context) bool {
			if !value.IsObject() {
				return true
			}
			if idx > 0 {
				out = append(out, ',')
			}
			out = append(out, removeOuterBraces(value.unprocessed)...)
			idx++
			return true
		})
	} else {
		// Deduplicate keys and generate an object with stable ordering.
		var keys []Context
		keyVal := make(map[string]Context)
		res.ForEach(func(_, value Context) bool {
			if !value.IsObject() {
				return true
			}
			value.ForEach(func(key, value Context) bool {
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
				out = append(out, ',')
			}
			out = append(out, keys[i].unprocessed...)
			out = append(out, ':')
			out = append(out, keyVal[keys[i].String()].unprocessed...)
		}
	}
	out = append(out, '}')
	return fromBytes2Str(out)
}

// @valid ensures that the json is valid before moving on. An empty string is
// returned when the json is not valid, otherwise it returns the original json.
func modValid(json, arg string) string {
	if !Valid(json) {
		return ""
	}
	return json
}

// @fromstr converts a string to json
//
//	"{\"id\":1023,\"name\":\"alert\"}" -> {"id":1023,"name":"alert"}
func modFromStr(json, arg string) string {
	if !Valid(json) {
		return ""
	}
	return Parse(json).String()
}

// @tostr converts a string to json
//
//	{"id":1023,"name":"alert"} -> "{\"id\":1023,\"name\":\"alert\"}"
func modToStr(str, arg string) string {
	return string(AppendJsonString(nil, str))
}

func modGroup(json, arg string) string {
	res := Parse(json)
	if !res.IsObject() {
		return ""
	}
	var all [][]byte
	res.ForEach(func(key, value Context) bool {
		if !value.IsArray() {
			return true
		}
		var idx int
		value.ForEach(func(_, value Context) bool {
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
	for i, item := range all {
		if i > 0 {
			data = append(data, ',')
		}
		data = append(data, '{')
		data = append(data, item[1:]...)
		data = append(data, '}')
	}
	data = append(data, ']')
	return string(data)
}

// Paths returns the original bjson paths for a Result where the Result came
// from a simple query path that returns an array, like:
//
//	bjson.Get(json, "friends.#.first")
//
// The returned value will be in the form of a JSON array:
//
//	["friends.0.first","friends.1.first","friends.2.first"]
//
// The param 'json' must be the original JSON used when calling Get.
//
// Returns an empty string if the paths cannot be determined, which can happen
// when the Result came from a path that contained a multi-path, modifier,
// or a nested query.
func (t Context) Paths(json string) []string {
	if t.indexes == nil {
		return nil
	}
	paths := make([]string, 0, len(t.indexes))
	t.ForEach(func(_, value Context) bool {
		paths = append(paths, value.Path(json))
		return true
	})
	if len(paths) != len(t.indexes) {
		return nil
	}
	return paths
}

// Path returns the original bjson path for a Result where the Result came
// from a simple path that returns a single value, like:
//
//	bjson.Get(json, "friends.#(last=Murphy)")
//
// The returned value will be in the form of a JSON string:
//
//	"friends.0"
//
// The param 'json' must be the original JSON used when calling Get.
//
// Returns an empty string if the paths cannot be determined, which can happen
// when the Result came from a path that contained a multi-path, modifier,
// or a nested query.
func (t Context) Path(json string) string {
	var path []byte
	var comps []string // raw components
	i := t.index - 1
	if t.index+len(t.unprocessed) > len(json) {
		// JSON cannot safely contain Result.
		goto fail
	}
	if !strings.HasPrefix(json[t.index:], t.unprocessed) {
		// Result is not at the JSON index as expected.
		goto fail
	}
	for ; i >= 0; i-- {
		if json[i] <= ' ' {
			continue
		}
		if json[i] == ':' {
			// inside of object, get the key
			for ; i >= 0; i-- {
				if json[i] != '"' {
					continue
				}
				break
			}
			raw := reverseSquash(json[:i+1])
			i = i - len(raw)
			comps = append(comps, raw)
			// key gotten, now squash the rest
			raw = reverseSquash(json[:i+1])
			i = i - len(raw)
			i++ // increment the index for next loop step
		} else if json[i] == '{' {
			// Encountered an open object. The original result was probably an
			// object key.
			goto fail
		} else if json[i] == ',' || json[i] == '[' {
			// inside of an array, count the position
			var arrIdx int
			if json[i] == ',' {
				arrIdx++
				i--
			}
			for ; i >= 0; i-- {
				if json[i] == ':' {
					// Encountered an unexpected colon. The original result was
					// probably an object key.
					goto fail
				} else if json[i] == ',' {
					arrIdx++
				} else if json[i] == '[' {
					comps = append(comps, strconv.Itoa(arrIdx))
					break
				} else if json[i] == ']' || json[i] == '}' || json[i] == '"' {
					raw := reverseSquash(json[:i+1])
					i = i - len(raw) + 1
				}
			}
		}
	}
	if len(comps) == 0 {
		if DisableModifiers {
			goto fail
		}
		return "@this"
	}
	for i := len(comps) - 1; i >= 0; i-- {
		rawComplexity := Parse(comps[i])
		if !rawComplexity.Exists() {
			goto fail
		}
		comp := EscapeUnsafeChars(rawComplexity.String())
		path = append(path, '.')
		path = append(path, comp...)
	}
	if len(path) > 0 {
		path = path[1:]
	}
	return string(path)
fail:
	return ""
}

func parseRecursiveDescent(all []Context, parent Context, path string) []Context {
	if res := parent.Get(path); res.Exists() {
		all = append(all, res)
	}
	if parent.IsArray() || parent.IsObject() {
		parent.ForEach(func(_, val Context) bool {
			all = parseRecursiveDescent(all, val, path)
			return true
		})
	}
	return all
}

func modDig(json, arg string) string {
	all := parseRecursiveDescent(nil, Parse(json), arg)
	var out []byte
	out = append(out, '[')
	for i, res := range all {
		if i > 0 {
			out = append(out, ',')
		}
		out = append(out, res.unprocessed...)
	}
	out = append(out, ']')
	return string(out)
}
