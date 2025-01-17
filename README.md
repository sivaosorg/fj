<h1>fj</h1>

**fj** (_Fast JSON_) is a Go package that provides a fast and simple way to retrieve values from a JSON document.

# Getting Started

## Requirements

Go version **1.19** or higher

## Installation

To start using `fj`, run `go get`:

- For a specific version:

  ```bash
  go get https://github.com/sivaosorg/fj@v0.0.1
  ```

- For the latest version (globally):
  ```bash
  go get -u https://github.com/sivaosorg/fj@latest
  ```

## Example

Given the [data JSON](./assets/data.json) string

## Syntax

A path is a string format used to define a pattern for efficiently retrieving values from a JSON structure.

### Path

A `fj` path is designed to be represented as a sequence of elements divided by a `.` symbol.

In addition to the `.` symbol, several other characters hold special significance, such as `|`, `#`, `@`, `\`, `*`, `!`, and `?`.

### Object

- **Basic**: in most situations, you'll simply need to access values using the object name or array index.

```shell
> id >> "http://subs/base-sample-schema.json"
> properties.alias.description >> "An unique identifier in a submission."
> properties.alias.minLength >> 1
> required >> ["alias", "taxonId", "releaseDate"]
> required.0 >> "alias"
> required.1 >> "taxonId"
> oneOf.0.required >> ["alias", "team"]
> oneOf.0.required.1 >> "team"
> properties.sampleRelationships >> { "$ref": "definitions-schema.json#/definitions/sampleRelationships" }
```

- **Wildcards**: A key can include special wildcard symbols like `*` and `?`. The `*` matches any sequence of characters (including none), while `?` matches exactly one character.

```shell
> anim*ls.1.name >> "Barky"
> *nimals.1.name >> "Barky"
```

- **Escape Character**: Characters with special meanings, like `.`, `*`, and `?`, can be escaped using the `\` symbol.

```shell
> properties.alias\.description >> "An unique identifier in a submission."
```

### Array

The `#` symbol enables navigation within JSON arrays. To retrieve the length of an array, simply use the `#` on its own.

```shell
> animals.# >> 3 (length of an array)
> animals.#.name >> ["Meowsy","Barky","Purrpaws"]
```

### Queries

- You can also search an array for the first match by using `#(...)`, or retrieve all matches with `#(...)#`.
  Queries support comparison operators such as `==`, `!=`, `<`, `<=`, `>`, `>=`, along with simple pattern matching operators `%` (like) and `!%` (not like).

```shell
> stock.#(price_2002==56.27).symbol >> "MMM"
> stock.#(company=="Amazon.com").symbol >> "AMZN"
> stock.#(initial_price>=10)#.symbol >> ["MMM","AMZN","CPB","DIS","DOW","XOM","F","GPS","GIS"]
> stock.#(company%"D*")#.symbol >> ["DIS","DOW"]
> stock.#(company!%"D*")#.symbol >> ["MMM","AMZN","CPB","XOM","F","GPS","GIS"]
> stock.#(company!%"F*")#.symbol >> ["MMM","AMZN","CPB","DIS","DOW","XOM","GPS","GIS"]
> stock.#(description%"*stores*")#.symbol >> ["CPB","GPS","GIS"]
> required.#(%"*as*")# >> ["alias","releaseDate"]
> required.#(%"*as*") >> "alias"
> required.#(!%"*s*") >> "taxonId"
> animals.#(foods.likes.#(%"*a*"))#.name >> ["Meowsy","Barky"]
```

- The `~` (tilde) operator evaluates a value as a boolean before performing a comparison.
  The most recent value that did not exist is considered `false`.
  The supported tilde comparison types are:

```shell
~true      Interprets truthy values as true
~false     Interprets falsy and undefined values as true
~null      Interprets null and undefined values as true
~*         Interprets any defined value as true
```

eg.

```shell
> bank.#(isActive==~true)#.name >> ["Davis Wade","Oneill Everett"]
> bank.#(isActive==~false)#.name >> ["Stark Jenkins","Odonnell Rollins","Rachelle Chang","Dalton Waters"]
> bank.#(eyeColor==~null)#.name >> ["Dalton Waters"]
> bank.#(company==~*)#.name >> ["Stark Jenkins","Odonnell Rollins","Rachelle Chang","Davis Wade","Oneill Everett","Dalton Waters"]
```

### Dot & Pipe

The `.` is the default separator, but you can also use a `|`.  
In most situations, both produce the same results.  
However, the `|` behaves differently from `.` when used after the `#` in the context of arrays and queries.

```shell
> bank.0.balance >> "$1,404.23"
> bank|0.balance >> "$1,404.23"
> bank.0|balance >> "$1,404.23"
> bank|0|balance >> "$1,404.23"
> bank.# >> 6
> bank|# >> 6
> bank.#(gender=="female")# >> [{"isActive":false,"balance":"$2,284.89","age":20,"eyeColor":"brown","name":"Rachelle Chang","gender":"female","company":"VERAQ","email":"rachellechang@veraq.com","phone":"+1 (955) 564-2002","address":"220 Drew Street, Ventress, Puerto Rico, 8432"},{"isActive":true,"balance":"$1,624.60","age":39,"eyeColor":"brown","name":"Davis Wade","gender":"female","company":"ASSISTIX","email":"daviswade@assistix.com","phone":"+1 (836) 432-2542","address":"532 Amity Street, Yukon, Palau, 3561"}]
> bank.#(gender=="female")#|# >> 2
> bank.#(gender=="female")#.name >> ["Rachelle Chang","Davis Wade"]
> bank.#(gender=="female")#|name >> not-present
> bank.#(gender=="female")#|0 >> {"isActive":false,"balance":"$2,284.89","age":20,"eyeColor":"brown","name":"Rachelle Chang","gender":"female","company":"VERAQ","email":"rachellechang@veraq.com","phone":"+1 (955) 564-2002","address":"220 Drew Street, Ventress, Puerto Rico, 8432"}
> bank.#(gender=="female")#.0 >> []
> bank.#(gender=="female")#.# >> []
```

### Multi-Selectors

The capability to combine multiple selectors into a single structure. Enclosing comma-separated selectors within `[...]` creates a new array, while using `{...}` generates a new object.

eg.

```shell
> {version,author,type,"stock_statics_symbol":stock.#(price_2007>=10)#.symbol} >> {"version":"1.0.0","author":"subs","type":"object","stock_statics_symbol":["MMM","AMZN","CPB","DIS","DOW","XOM","GPS","GIS"]}
> {version,author,type,"stock_statics_symbol":stock.#(company%"*m*")#.symbol} >> {"version":"1.0.0","author":"subs","type":"object","stock_statics_symbol":["AMZN","CPB","DOW"]}
```

### Literals

Support for JSON literals offers a straightforward method for creating static JSON blocks. This feature is especially helpful when building a new JSON document using [Multi-Selectors](#multi-selectors).
A JSON literal is introduced with the `!` declaration character.

eg. Add 2 fields, `marked` and `scope`

```shell
> {version,author,type,"stock_statics_symbol":stock.#(price_2007>=10)#.symbol,"marked":!true,"scope":!"static"} >> {"version":"1.0.0","author":"subs","type":"object","stock_statics_symbol":["MMM","AMZN","CPB","DIS","DOW","XOM","GPS","GIS"],"marked":true,"scope":"static"}
```

### Transformers

A transformer is a path component used to apply custom transformations to the JSON.
The following built-in transformers are currently available:

| Transformer   | Description                                                                                                                                                  | Arguments (optional)                                                         |
| ------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---------------------------------------------------------------------------- |
| `@trim`       | Remove leading and trailing whitespace from the input JSON string.                                                                                           |                                                                              |
| `@this`       | Use as a default case when no specific transformation is requested or supported.                                                                             |                                                                              |
| `@valid`      | Ensure the JSON is valid before processing it further. If the JSON string is not valid, then returns an empty string                                         |                                                                              |
| `@pretty`     | Format the JSON string into a human-readable, indented format.                                                                                               | `@pretty:{"sort_keys": true, "indent": "\t", "prefix": "tick", "width": 10}` |
| `@minify`     | Remove all whitespace characters from the JSON string, transforming it into a compact, single-line format order                                              |                                                                              |
| `@flip`       | Reverses the order of its characters                                                                                                                         |                                                                              |
| `@reverse`    | Reverse the order of elements in an array or the order of key-value pairs in an object function                                                              |                                                                              |
| `@flatten`    | Flattens a JSON array by removing any nested arrays within it                                                                                                | `@flatten:{"deep": true}`                                                    |
| `@join`       | Merges multiple JSON objects into a single object                                                                                                            | `@join:{"preserve": true}`                                                   |
| `@keys`       | Extracts the keys from a JSON object and returns them as a JSON array of strings                                                                             |                                                                              |
| `@values`     | Extracts the values from a JSON object and returns them as a JSON array of values                                                                            |                                                                              |
| `@string`     | Converts a regular string into a valid JSON string format                                                                                                    |                                                                              |
| `@json`       | Converts a string into a valid JSON representation                                                                                                           |                                                                              |
| `@group`      | Processes a JSON string containing objects and arrays, and groups the elements of arrays within objects by their keys                                        |                                                                              |
| `@search`     | Performs a value lookup on a JSON structure based on the specified path and returns a JSON-encoded string containing all matching values found at that path. |                                                                              |
| `@uppercase`  | Converts the JSON string to uppercase.                                                                                                                       |                                                                              |
| `@lowercase`  | Converts the JSON string to lowercase.                                                                                                                       |                                                                              |
| `@snakeCase`  | Converts the string to snake_case format                                                                                                                     |                                                                              |
| `@camelCase`  | Converts the string into camelCase                                                                                                                           |                                                                              |
| `@kebabCase`  | Converts the input string into kebab-case, often used for URL slugs                                                                                          |                                                                              |
| `@replace`    | Replaces a specific substring within the input string with another string                                                                                    | `@replace:{"target": "abc", "replacement": "111"}`                           |
| `@replaceAll` | Replaces all occurrences of a target substring with a replacement string                                                                                     | `@replaceAll:{"target": "abc", "replacement": "111"}`                        |
| `@hex`        | Converts the string to its hexadecimal representation                                                                                                        |                                                                              |
| `@bin`        | Converts the string to its binary representation                                                                                                             |                                                                              |
| `@insertAt`   | Inserts a specified string at a given index in the input string                                                                                              | `@insertAt:{"index": 0, "insert": "abc"}`                                    |
| `@wc`         | Counts the number of words in the input string                                                                                                               |                                                                              |
| `@padLeft`    | Pads the input string with a specified character on the left to a given length                                                                               | `@padLeft:{"padding": "*", "length": 30}`                                    |
| `@padRight`   | Pads the input string with a specified character on the right to a given length                                                                              | `@padRight:{"padding": "*", "length": 30}`                                   |

eg.

```shell
> required.1.@flip >> "dInoxat"
> required.@reverse >> ["releaseDate","taxonId","alias"]
> required.@reverse.0 >> "releaseDate"
> required.@reverse.1 >> "taxonId"
> animals.@join.@minify >> {"name":"Purrpaws","species":"cat","foods":{"likes":["mice"],"dislikes":["cookies"]}}
> animals.1.@keys >> ["name","species","foods"]
> animals.1.@values.@minify >> ["Barky","dog",{"likes":["bones","carrots"],"dislikes":["tuna"]}]
> {"id":bank.#.company,"details":bank.#(age>=10)#.eyeColor}|@group >> [{"id":"HINWAY","details":"blue"},{"id":"NEXGENE","details":"green"},{"id":"VERAQ","details":"brown"},{"id":"ASSISTIX","details":"brown"},{"id":"INCUBUS","details":"green"},{"id":"OVATION","details":null}]
> {"id":bank.#.company,"details":bank.#(age>=10)#.eyeColor}|@group|# >> 6
> stock.@search:#(price_2007>=50)|0.company >> "3M"
> stock.@search:#(price_2007>=50)|0.company.@lowercase >> "3m"
> stock.0.company.@hex >> "334d"
> stock.0.company.@bin >> "0011001101001101"
> stock.0.description.@wc >> 42
> author|@padLeft:{"padding": "*", "length": 15}|@string >> "***********subs"
> author|@padRight:{"padding": "*", "length": 15}|@string >> "subs***********"
> bank.0.@pretty:{"sort_keys": true}    >>
    {
        "address": "766 Cooke Court, Dunbar, Connecticut, 9512",
        "age": 26,
        "balance": "$1,404.23",
        "company": "HINWAY",
        "email": "starkjenkins@hinway.com",
        "eyeColor": "blue",
        "gender": "male",
        "isActive": false,
        "name": "Stark Jenkins",
        "phone": "+1 (943) 542-3591"
    }
```

## Usage

### Retrieve a value

Retrieves JSON data from the specified path, using dot notation (`.`) like "user.name" or "settings.theme". The value is returned as soon as it's located.

eg.

```go
package main

import (
	"fmt"

	"github.com/sivaosorg/fj"
)

const json = `
{"user":{"id":"12345","name":{"firstName":"John","lastName":"Doe"},"email":"john.doe@example.com","phone":"+1-555-555-5555","address":{"street":"123 Main St","city":"Anytown","state":"CA","postalCode":"12345","country":"USA"},"roles":[{"roleId":"1","roleName":"Admin","permissions":[{"permissionId":"101","permissionName":"View Reports","allowedActions":["view","download"]},{"permissionId":"102","permissionName":"Manage Users","allowedActions":["create","update","delete"]}]},{"roleId":"2","roleName":"Editor","permissions":[{"permissionId":"201","permissionName":"Edit Content","allowedActions":["create","edit","publish"]},{"permissionId":"202","permissionName":"View Analytics","allowedActions":["view"]}]}],"status":"active","createdAt":"2025-01-01T10:00:00Z","lastLogin":"2025-01-12T15:30:00Z"}}
`

func main() {
	value := fj.Get(json, "user.id")
	fmt.Println(value) // 12345
	value = fj.Get(json, "user.roles.#.roleName")
	fmt.Println(value) // ["Admin","Editor"]
	value = fj.Get(json, "user.roles.#.permissions.#.permissionName")
	fmt.Println(value) // [["View Reports","Manage Users"],["Edit Content","View Analytics"]]
	value = fj.Get(json, "user.address")
	fmt.Println(value) // {"street":"123 Main St","city":"Anytown","state":"CA","postalCode":"12345","country":"USA"}
}
```

### Working with Bytes

If your JSON is stored in a `[]byte` slice, you can use the `GetBytes([]byte(json), path)` function, which is recommended instead of using `Get(string(json), path)`.

eg.

```go
package main

import (
	"fmt"

	"github.com/sivaosorg/fj"
)

const json = `
{"user":{"id":"12345","name":{"firstName":"John","lastName":"Doe"},"email":"john.doe@example.com","phone":"+1-555-555-5555","address":{"street":"123 Main St","city":"Anytown","state":"CA","postalCode":"12345","country":"USA"},"roles":[{"roleId":"1","roleName":"Admin","permissions":[{"permissionId":"101","permissionName":"View Reports","allowedActions":["view","download"]},{"permissionId":"102","permissionName":"Manage Users","allowedActions":["create","update","delete"]}]},{"roleId":"2","roleName":"Editor","permissions":[{"permissionId":"201","permissionName":"Edit Content","allowedActions":["create","edit","publish"]},{"permissionId":"202","permissionName":"View Analytics","allowedActions":["view"]}]}],"status":"active","createdAt":"2025-01-01T10:00:00Z","lastLogin":"2025-01-12T15:30:00Z"}}
`

func main() {
	value := fj.GetBytes([]byte(json), "user.id")
	fmt.Println(value) // 12345
	value = fj.GetBytes([]byte(json), "user.roles.#.roleName")
	fmt.Println(value) // ["Admin","Editor"]
	value = fj.GetBytes([]byte(json), "user.roles.#.permissions.#.permissionName")
	fmt.Println(value) // [["View Reports","Manage Users"],["Edit Content","View Analytics"]]
	value = fj.GetBytes([]byte(json), "user.address")
	fmt.Println(value) // {"street":"123 Main St","city":"Anytown","state":"CA","postalCode":"12345","country":"USA"}
}
```

If you're using the `GetBytes([]byte(json), path)` function and want to avoid converting `ctx.Unprocessed()` into a `[]byte`. This is an attempt at a no-allocation sub slice of the original JSON. The method relies on the `ctx.Index()` field, which indicates the position of the raw data in the original JSON. It's possible that `ctx.Index()` could be zero, in which case `ctx.Unprocessed()` is converted into a `[]byte`. You can follow this approach:

eg.

```go
package main

import (
	"fmt"

	"github.com/sivaosorg/fj"
)

var json []byte = []byte(`{"user":{"id":"12345","name":{"firstName":"John","lastName":"Doe"},"email":"john.doe@example.com","phone":"+1-555-555-5555","address":{"street":"123 Main St","city":"Anytown","state":"CA","postalCode":"12345","country":"USA"},"roles":[{"roleId":"1","roleName":"Admin","permissions":[{"permissionId":"101","permissionName":"View Reports","allowedActions":["view","download"]},{"permissionId":"102","permissionName":"Manage Users","allowedActions":["create","update","delete"]}]},{"roleId":"2","roleName":"Editor","permissions":[{"permissionId":"201","permissionName":"Edit Content","allowedActions":["create","edit","publish"]},{"permissionId":"202","permissionName":"View Analytics","allowedActions":["view"]}]}],"status":"active","createdAt":"2025-01-01T10:00:00Z","lastLogin":"2025-01-12T15:30:00Z"}}`)

func main() {
	ctx := fj.GetBytes(json, "user.roles.#.roleName")
	var raw []byte
	if ctx.Index() > 0 {
		raw = json[ctx.Index() : ctx.Index()+len(ctx.Unprocessed())]
	} else {
		raw = []byte(ctx.Unprocessed())
	}
	fmt.Println(string(raw)) // ["Admin","Editor"]
}
```

### Read JSON by File

If you're working with the JSON dataset stored in a file (with the extension `.json` or `.txt`), you can use the `ParseBufio(in io.Reader)` function to read all the content of the JSON into the `Context` struct.

eg.

```go
package main

import (
	"fmt"
	"os"

	"github.com/sivaosorg/fj"
)

func main() {
	file, err := os.Open("./assets/data.json")
	if err != nil {
		return
	}
	defer file.Close()

	ctx := fj.ParseBufio(file)
	value := ctx.Get(`stock.0.symbol`)
	fmt.Println(value.String()) // MMM
	value = ctx.Get(`bank.0.@minify`)
	fmt.Println(value.String()) // {"isActive":false,"balance":"$1,404.23","age":26,"eyeColor":"blue","name":"Stark Jenkins","gender":"male","company":"HINWAY","email":"starkjenkins@hinway.com","phone":"+1 (943) 542-3591","address":"766 Cooke Court, Dunbar, Connecticut, 9512"}
}
```

> or

```go
package main

import (
	"fmt"

	"github.com/sivaosorg/fj"
)

func main() {
	ctx := fj.ParseFilepath("./assets/data.json")
	value := ctx.Get(`stock.0.symbol`)
	fmt.Println(value.String()) // MMM
	value = ctx.Get(`bank.0.@minify`)
	fmt.Println(value.String()) // {"isActive":false,"balance":"$1,404.23","age":26,"eyeColor":"blue","name":"Stark Jenkins","gender":"male","company":"HINWAY","email":"starkjenkins@hinway.com","phone":"+1 (943) 542-3591","address":"766 Cooke Court, Dunbar, Connecticut, 9512"}
}
```

### Validation

The `Get*` and `Parse*` functions assume that the JSON is properly structured. Invalid JSON won't cause a panic, but it could lead to unexpected outcomes.

eg.

```go
package main

import (
	"fmt"

	"github.com/sivaosorg/fj"
)

var json []byte = []byte(`{"user":{"id":"12345","name":{"firstName":"John","lastName":"Doe"},"email":"john.doe@example.com","phone":"+1-555-555-5555","address":{"street":"123 Main St","city":"Anytown","state":"CA","postalCode":"12345","country":"USA"},"roles":[{"roleId":"1","roleName":"Admin","permissions":[{"permissionId":"101","permissionName":"View Reports","allowedActions":["view","download"]},{"permissionId":"102","permissionName":"Manage Users","allowedActions":["create","update","delete"]}]},{"roleId":"2","roleName":"Editor","permissions":[{"permissionId":"201","permissionName":"Edit Content","allowedActions":["create","edit","publish"]},{"permissionId":"202","permissionName":"View Analytics","allowedActions":["view"]}]}],"status":"active","createdAt":"2025-01-01T10:00:00Z","lastLogin":"2025-01-12T15:30:00Z"}}`)

func main() {
	isValid := fj.IsValidJSONBytes(json) // or fj.IsValidJSON(string(json))
	fmt.Println(isValid)                 // true
	isValid = fj.IsValidJSON("{}")
	fmt.Println(isValid) // true
}
```

### Existence

Occasionally, you simply need to check if a value is present.

```go
package main

import (
	"fmt"

	"github.com/sivaosorg/fj"
)

var json []byte = []byte(`{"user":{"id":"12345","name":{"firstName":"John","lastName":"Doe"},"email":"john.doe@example.com","phone":"+1-555-555-5555","address":{"street":"123 Main St","city":"Anytown","state":"CA","postalCode":"12345","country":"USA"},"roles":[{"roleId":"1","roleName":"Admin","permissions":[{"permissionId":"101","permissionName":"View Reports","allowedActions":["view","download"]},{"permissionId":"102","permissionName":"Manage Users","allowedActions":["create","update","delete"]}]},{"roleId":"2","roleName":"Editor","permissions":[{"permissionId":"201","permissionName":"Edit Content","allowedActions":["create","edit","publish"]},{"permissionId":"202","permissionName":"View Analytics","allowedActions":["view"]}]}],"status":"active","createdAt":"2025-01-01T10:00:00Z","lastLogin":"2025-01-12T15:30:00Z"}}`)

func main() {
	ctx := fj.GetBytes(json, "user.name.firstName")
	if ctx.Exists() {
		fmt.Println(ctx.String()) // John
	} else {
		fmt.Println("Value of 'firstName' not found")
	}
}
```

### Loop

The `Foreach` method enables efficient iteration over objects or arrays. For objects, both the key and value are provided to the callback function, while for arrays, only the value is passed. The iteration can be halted by returning `false` from the callback.

eg.

```go
package main

import (
	"fmt"

	"github.com/sivaosorg/fj"
)

var json []byte = []byte(`{"user":{"id":"12345","name":{"firstName":"John","lastName":"Doe"},"email":"john.doe@example.com","phone":"+1-555-555-5555","address":{"street":"123 Main St","city":"Anytown","state":"CA","postalCode":"12345","country":"USA"},"roles":[{"roleId":"1","roleName":"Admin","permissions":[{"permissionId":"101","permissionName":"View Reports","allowedActions":["view","download"]},{"permissionId":"102","permissionName":"Manage Users","allowedActions":["create","update","delete"]}]},{"roleId":"2","roleName":"Editor","permissions":[{"permissionId":"201","permissionName":"Edit Content","allowedActions":["create","edit","publish"]},{"permissionId":"202","permissionName":"View Analytics","allowedActions":["view"]}]}],"status":"active","createdAt":"2025-01-01T10:00:00Z","lastLogin":"2025-01-12T15:30:00Z"}}`)

func main() {
	ctx := fj.GetBytes(json, "user.roles.#.permissions")
	ctx.Foreach(func(key, value fj.Context) bool {
		fmt.Println(value.String())
		return true
	})
}

// output:
// [{"permissionId":"101","permissionName":"View Reports","allowedActions":["view","download"]},{"permissionId":"102","permissionName":"Manage Users","allowedActions":["create","update","delete"]}]
// [{"permissionId":"201","permissionName":"Edit Content","allowedActions":["create","edit","publish"]},{"permissionId":"202","permissionName":"View Analytics","allowedActions":["view"]}]
```

### Unmarshal

eg.

```go
package main

import (
	"fmt"

	"github.com/sivaosorg/fj"
)

var json []byte = []byte(`{"user":{"id":"12345","name":{"firstName":"John","lastName":"Doe"},"email":"john.doe@example.com","phone":"+1-555-555-5555","address":{"street":"123 Main St","city":"Anytown","state":"CA","postalCode":"12345","country":"USA"},"roles":[{"roleId":"1","roleName":"Admin","permissions":[{"permissionId":"101","permissionName":"View Reports","allowedActions":["view","download"]},{"permissionId":"102","permissionName":"Manage Users","allowedActions":["create","update","delete"]}]},{"roleId":"2","roleName":"Editor","permissions":[{"permissionId":"201","permissionName":"Edit Content","allowedActions":["create","edit","publish"]},{"permissionId":"202","permissionName":"View Analytics","allowedActions":["view"]}]}],"status":"active","createdAt":"2025-01-01T10:00:00Z","lastLogin":"2025-01-12T15:30:00Z"}}`)

func main() {
	ctx := fj.GetBytes(json, "user.name")
	value, ok := ctx.Value().(map[string]interface{})
	if !ok {
		fmt.Println("value is not a map")
		return
	}
	fmt.Println(value) // map[firstName:John lastName:Doe]
}
```

```go
package main

import (
	"fmt"

	"github.com/sivaosorg/fj"
)

var json []byte = []byte(`{"user":{"id":12345,"name":{"firstName":"John","lastName":"Doe"},"email":"john.doe@example.com","phone":"+1-555-555-5555","address":{"street":"123 Main St","city":"Anytown","state":"CA","postalCode":"12345","country":"USA"},"roles":[{"roleId":"1","roleName":"Admin","permissions":[{"permissionId":"101","permissionName":"View Reports","allowedActions":["view","download"]},{"permissionId":"102","permissionName":"Manage Users","allowedActions":["create","update","delete"]}]},{"roleId":"2","roleName":"Editor","permissions":[{"permissionId":"201","permissionName":"Edit Content","allowedActions":["create","edit","publish"]},{"permissionId":"202","permissionName":"View Analytics","allowedActions":["view"]}]}],"status":"active","createdAt":"2025-01-01T10:00:00Z","lastLogin":"2025-01-12T15:30:00Z"}}`)

func main() {
	ctx := fj.GetBytes(json, "user.id")
	value, ok := ctx.Value().(float64)
	if !ok {
		fmt.Println("value is not a float64")
		return
	}
	fmt.Println(value) //12345
}
```

```go
package main

import (
	"fmt"

	"github.com/sivaosorg/fj"
)

var json []byte = []byte(`{"user":{"id":12345,"name":{"firstName":"John","lastName":"Doe"},"email":"john.doe@example.com","phone":"+1-555-555-5555","address":{"street":"123 Main St","city":"Anytown","state":"CA","postalCode":"12345","country":"USA"},"roles":[{"roleId":"1","roleName":"Admin","permissions":[{"permissionId":"101","permissionName":"View Reports","allowedActions":["view","download"]},{"permissionId":"102","permissionName":"Manage Users","allowedActions":["create","update","delete"]}]},{"roleId":"2","roleName":"Editor","permissions":[{"permissionId":"201","permissionName":"Edit Content","allowedActions":["create","edit","publish"]},{"permissionId":"202","permissionName":"View Analytics","allowedActions":["view"]}]}],"status":"active","createdAt":"2025-01-01T10:00:00Z","lastLogin":"2025-01-12T15:30:00Z"}}`)

func main() {
	ctx := fj.GetBytes(json, "user.id")
	fmt.Println(ctx.Int64()) // 12345
}
```

### Parse & Get

The `Parse*(json)` function to perform a straightforward parsing, and `ctx.Get*(path)` to retrieve a value from the parsed result.

eg.

```go
package main

import (
	"fmt"

	"github.com/sivaosorg/fj"
)

var json []byte = []byte(`{"user":{"id":12345,"name":{"firstName":"John","lastName":"Doe"},"email":"john.doe@example.com","phone":"+1-555-555-5555","address":{"street":"123 Main St","city":"Anytown","state":"CA","postalCode":"12345","country":"USA"},"roles":[{"roleId":"1","roleName":"Admin","permissions":[{"permissionId":"101","permissionName":"View Reports","allowedActions":["view","download"]},{"permissionId":"102","permissionName":"Manage Users","allowedActions":["create","update","delete"]}]},{"roleId":"2","roleName":"Editor","permissions":[{"permissionId":"201","permissionName":"Edit Content","allowedActions":["create","edit","publish"]},{"permissionId":"202","permissionName":"View Analytics","allowedActions":["view"]}]}],"status":"active","createdAt":"2025-01-01T10:00:00Z","lastLogin":"2025-01-12T15:30:00Z"}}`)

func main() {
	ctx := fj.ParseBytes(json).Get("user").Get("id")
	fmt.Println(ctx.Int64()) // 12345
	ctx = fj.GetBytes(json, "user.id")
	fmt.Println(ctx.Int64()) // 12345
	ctx = fj.GetBytes(json, "user").Get("id")
	fmt.Println(ctx.Int64()) // 12345
}
```

### JSON Lines

Support for [JSON Lines](https://jsonlines.org/) is available using the `..` prefix, enabling the treatment of a multi-line document as an array.

eg.

```go
package main

import (
	"fmt"

	"github.com/sivaosorg/fj"
)

var json []byte = []byte(`
		{"roleId":"1","roleName":"Admin","permissions":[{"permissionId":"101","permissionName":"View Reports","allowedActions":["view","download"]},{"permissionId":"102","permissionName":"Manage Users","allowedActions":["create","update","delete"]}]}
		{"roleId":"2","roleName":"Editor","permissions":[{"permissionId":"201","permissionName":"Edit Content","allowedActions":["create","edit","publish"]},{"permissionId":"202","permissionName":"View Analytics","allowedActions":["view"]}]}
`)

func main() {
	ctx := fj.ParseBytes(json).Get("..#")
	fmt.Println(ctx.String()) // 2
	ctx = fj.GetBytes(json, "..0")
	fmt.Println(ctx.String()) // {"roleId":"1","roleName":"Admin","permissions":[{"permissionId":"101","permissionName":"View Reports","allowedActions":["view","download"]},{"permissionId":"102","permissionName":"Manage Users","allowedActions":["create","update","delete"]}]}
	ctx = fj.GetBytes(json, "..#.roleName")
	fmt.Println(ctx.String()) // ["Admin","Editor"]
	ctx = fj.GetBytes(json, `..#.permissions.#(permissionId=="101").permissionName`)
	fmt.Println(ctx.String()) // ["View Reports"]
}
```

### Transformers

A transformer is a path component used to apply custom transformations to the JSON

eg.

```go
package main

import (
	"fmt"

	"github.com/sivaosorg/fj"
)

var json []byte = []byte(`
{"id":"http://subs/base-sample-schema.json","schema":"http://json-schema.org/draft-07/schema#","title":"Submissions Sample Schema","description":"A base validation sample schema","version":"1.0.0","author":"subs","type":"object","properties":{"alias":{"description":"An unique identifier in a submission.","type":"string","minLength":1},"title":{"description":"Title of the sample.","type":"string","minLength":1},"description":{"description":"More extensive free-form description.","type":"string","minLength":1},"taxonId":{"type":"integer","minimum":1},"taxon":{"type":"string","minLength":1},"releaseDate":{"type":"string","format":"date"},"attributes":{"description":"Attributes for describing a sample.","type":"object","properties":{},"patternProperties":{"^.*$":{"$ref":"definitions-schema.json#/definitions/attributes_structure"}}},"sampleRelationships":{"$ref":"definitions-schema.json#/definitions/sampleRelationships"}},"required":["alias","taxonId","releaseDate"],"oneOf":[{"required":["alias","team"]},{"required":["accession"]}],"animals":[{"name":"Meowsy","species":"cat","foods":{"likes":["tuna","catnip"],"dislikes":["ham","zucchini"]}},{"name":"Barky","species":"dog","foods":{"likes":["bones","carrots"],"dislikes":["tuna"]}},{"name":"Purrpaws","species":"cat","foods":{"likes":["mice"],"dislikes":["cookies"]}}],"stock":[{"company":"3M","description":"3M, based in Minnesota, may be best known for its Scotch tape and Post-It Notes, but it also produces sand paper, adhesives, medical products, computer screen filters, food safety items, stationery products and many products used in automotive, marine, and aircraft industries.","initial_price":44.28,"price_2002":56.27,"price_2007":95.85,"symbol":"MMM"},{"company":"Amazon.com","description":"Amazon.com, Inc. is an online retailer in North America and internationally. The company serves consumers through its retail Web sites and focuses on selection, price, and convenience. It also offers programs that enable sellers to sell their products on its Web sites, and their own branded Web sites. In addition, the company serves developer customers through Amazon Web Services, which provides access to technology infrastructure that developers can use to enable virtually various type of business. Further, it manufactures and sells the Kindle e-reader. Founded in 1994 and headquartered in Seattle, Washington.","initial_price":89.38,"price_2002":17.01,"price_2007":93.43,"symbol":"AMZN"},{"company":"Campbell Soup","description":"Campbell Soup is a worldwide food company, offering condensed and ready-to-serve soups; broth, stocks, and canned poultry; pasta sauces; Mexican sauces; canned pastas, gravies, and beans; juices and beverages; and tomato juices. Its customers include retail food chains, mass discounters, mass merchandisers, club stores, convenience stores, drug stores and other retail, and commercial and non-commercial establishments. Campbell Soup Company was founded in 1869 and is headquartered in Camden, New Jersey.","initial_price":37,"price_2002":22.27,"price_2007":36.4,"symbol":"CPB"},{"company":"Disney","description":"The Walt Disney Company, founded in 1923, is a worldwide entertainment company, with movies, cable networks, radio networks, movie production, musical recordings and live stage plays. Disney also operates Walt Disney World in Florida and Disneyland in California, Disney Cruise Line, and international Disney resorts. Disney owns countless licenses and literary properties and publishes books and magazines.","initial_price":40.68,"price_2002":15.24,"price_2007":35.47,"symbol":"DIS"},{"company":"Dow Chemical","description":"The Dow Chemical Company manufactures raw materials that go into consumer products and services. These materials include food and pharmaceutical ingredients, electronic displays, semiconductor packaging, water purification, insulation, adhesives, pest control, polyurethane, polystyrene (goes into plastics), and crude-oil based raw materials. Dow was founded in 1897 and is based in Midland, Michigan.","initial_price":38.83,"price_2002":27.65,"price_2007":44.67,"symbol":"DOW"},{"company":"Exxon Mobil","description":"Exxon Mobil engages in the exploration and production of crude oil and natural gas, and manufacture of petroleum products. The company manufactures commodity petrochemicals. The company has operations in the United States, Canada/South America, Europe, Africa, Asia, and Australia/Oceania. Exxon Mobil Corporation was founded in1870 and is based in Irving, Texas.","initial_price":39,"price_2002":32.82,"price_2007":91.36,"symbol":"XOM"},{"company":"Ford","description":"Ford Motor Co. develops, manufactures, sells and services vehicles and parts worldwide. Ford sells cars and trucks primarily under the Ford and Lincoln brands. It sells to consumers (through retail dealers) and to rental car companies, leasing companies, and governments. Ford also provides maintenance and repair services. Ford also offers financing to vehicle purchasers. Ford was founded in 1903 and is based in Dearborn, Michigan.","initial_price":27.34,"price_2002":9.63,"price_2007":8.37,"symbol":"F"},{"company":"The Gap","description":"The Gap, Inc. sells retail clothing, accessories and personal care products globally under the brand names Gap, Old Navy, Banana Republic, Piperlime, Athleta and Intermix. Products include sports apparel, casual clothing, sleepwear, footwear and infants’ and children’s clothing. The company has company-owned stores as well as franchise stores, online stores and catalogs. The Gap was founded in 1969 and is headquartered in San Francisco, California.","initial_price":46,"price_2002":11.56,"price_2007":18.9,"symbol":"GPS"},{"company":"General Mills","description":"General Mills manufactures and sells consumer foods worldwide. Products include cereals, frozen vegetables, dough, dessert and baking mixes, frozen pizzas, grains, fruits, ice creams and organic products. It sells to grocery stores as well as commercial food service distributors, restaurants and convenience stores. General Mills was founded in 1928 and is based in Minneapolis, Minnesota.","initial_price":15.59,"price_2002":22.1,"price_2007":28.76,"symbol":"GIS"}],"bank":[{"isActive":false,"balance":"$1,404.23","age":26,"eyeColor":"blue","name":"Stark Jenkins","gender":"male","company":"HINWAY","email":"starkjenkins@hinway.com","phone":"+1 (943) 542-3591","address":"766 Cooke Court, Dunbar, Connecticut, 9512"},{"isActive":false,"balance":"$1,247.08","age":36,"eyeColor":"green","name":"Odonnell Rollins","gender":"male","company":"NEXGENE","email":"odonnellrollins@nexgene.com","phone":"+1 (810) 521-2350","address":"210 Pleasant Place, Lloyd, Mississippi, 1636"},{"isActive":false,"balance":"$2,284.89","age":20,"eyeColor":"brown","name":"Rachelle Chang","gender":"female","company":"VERAQ","email":"rachellechang@veraq.com","phone":"+1 (955) 564-2002","address":"220 Drew Street, Ventress, Puerto Rico, 8432"},{"isActive":true,"balance":"$1,624.60","age":39,"eyeColor":"brown","name":"Davis Wade","gender":"female","company":"ASSISTIX","email":"daviswade@assistix.com","phone":"+1 (836) 432-2542","address":"532 Amity Street, Yukon, Palau, 3561"},{"isActive":true,"balance":"$3,818.97","age":23,"eyeColor":"green","name":"Oneill Everett","gender":"male","company":"INCUBUS","email":"oneilleverett@incubus.com","phone":"+1 (958) 522-2724","address":"273 Temple Court, Shelby, Georgia, 8682"},{"isActive":false,"balance":"$3,243.63","age":21,"eyeColor":null,"name":"Dalton Waters","gender":"male","company":"OVATION","email":"daltonwaters@ovation.com","phone":"+1 (899) 464-3878","address":"909 Wyona Street, Adelino, Hawaii, 6449"}]}
`)

func main() {
	ctx := fj.ParseBytes(json).Get("required.1.@flip")
	fmt.Println(ctx.String()) // "dInoxat"

	ctx = fj.GetBytes(json, "required.@reverse")
	fmt.Println(ctx.String()) // ["releaseDate","taxonId","alias"]

	ctx = fj.GetBytes(json, "required.@reverse.0")
	fmt.Println(ctx.String()) // "releaseDate"

	ctx = fj.GetBytes(json, `required.@reverse.1`)
	fmt.Println(ctx.String()) // "taxonId"

	ctx = fj.GetBytes(json, `animals.@join.@minify`)
	fmt.Println(ctx.String()) // {"name":"Purrpaws","species":"cat","foods":{"likes":["mice"],"dislikes":["cookies"]}}

	ctx = fj.GetBytes(json, `animals.1.@keys`)
	fmt.Println(ctx.String()) // ["name","species","foods"]

	ctx = fj.GetBytes(json, `animals.1.@values.@minify`)
	fmt.Println(ctx.String()) // ["Barky","dog",{"likes":["bones","carrots"],"dislikes":["tuna"]}]

	ctx = fj.GetBytes(json, `{"id":bank.#.company,"details":bank.#(age>=10)#.eyeColor}|@group`)
	fmt.Println(ctx.String()) // [{"id":"HINWAY","details":"blue"},{"id":"NEXGENE","details":"green"},{"id":"VERAQ","details":"brown"},{"id":"ASSISTIX","details":"brown"},{"id":"INCUBUS","details":"green"},{"id":"OVATION","details":null}]

	ctx = fj.GetBytes(json, `{"id":bank.#.company,"details":bank.#(age>=10)#.eyeColor}|@group|#`)
	fmt.Println(ctx.String()) // 6

	ctx = fj.GetBytes(json, `stock.@search:#(price_2007>=50)|0.company.@lowercase`)
	fmt.Println(ctx.String()) // "3m"

	ctx = fj.GetBytes(json, `stock.0.company.@hex`)
	fmt.Println(ctx.Unprocessed()) // "334d"

	ctx = fj.GetBytes(json, `stock.0.company.@bin`)
	fmt.Println(ctx.String()) // "0011001101001101"

	ctx = fj.GetBytes(json, `stock.0.description.@wc`)
	fmt.Println(ctx.String()) // 42

	ctx = fj.GetBytes(json, `author|@padLeft:{"padding": "*", "length": 15}|@string`)
	fmt.Println(ctx.Unprocessed()) // "***********subs"

	ctx = fj.GetBytes(json, `bank.0.@pretty:{"sort_keys": true}`)
	fmt.Println(ctx.String()) //
	/**
	{
	  "address": "766 Cooke Court, Dunbar, Connecticut, 9512",
	  "age": 26,
	  "balance": "$1,404.23",
	  "company": "HINWAY",
	  "email": "starkjenkins@hinway.com",
	  "eyeColor": "blue",
	  "gender": "male",
	  "isActive": false,
	  "name": "Stark Jenkins",
	  "phone": "+1 (943) 542-3591"
	}
	*/
}
```

### Custom Transformer

You can add custom transformer

eg.

```go
package main

import (
	"fmt"
	"strings"

	"github.com/sivaosorg/fj"
)

var json []byte = []byte(`{"user":{"id":12345,"name":{"firstName":"John","lastName":"Doe"},"email":"john.doe@example.com","phone":"+1-555-555-5555","address":{"street":"123 Main St","city":"Anytown","state":"CA","postalCode":"12345","country":"USA"},"roles":[{"roleId":"1","roleName":"Admin","permissions":[{"permissionId":"101","permissionName":"View Reports","allowedActions":["view","download"]},{"permissionId":"102","permissionName":"Manage Users","allowedActions":["create","update","delete"]}]},{"roleId":"2","roleName":"Editor","permissions":[{"permissionId":"201","permissionName":"Edit Content","allowedActions":["create","edit","publish"]},{"permissionId":"202","permissionName":"View Analytics","allowedActions":["view"]}]}],"status":"active","createdAt":"2025-01-01T10:00:00Z","lastLogin":"2025-01-12T15:30:00Z"}}`)

func main() {
	// customize the transformer
	wordTransformer := func(json, arg string) string {
		if arg == "upper" {
			return strings.ToUpper(json)
		}
		if arg == "lower" {
			return strings.ToLower(json)
		}
		return json
	}
	fj.AddTransformer("word", wordTransformer)

	ctx := fj.GetBytes(json, "user.name.firstName.@word:upper")
	fmt.Println(ctx.Value()) // "JOHN"
	ctx = fj.GetBytes(json, "user").Get("name.firstName.@word:lower")
	fmt.Println(ctx.Value()) // john
}
```

### JSON Color

You can view JSON data in a color-formatted format, which is most suitable for debugging in the console terminal log.

Currently, we already have some JSON color formats available as follows.

```go
fj.DarkStyle // DarkStyle uses darker tones for styling.
fj.NeonStyle // NeonStyle is a vibrant style using neon-like colors.
fj.PastelStyle // PastelStyle applies softer colors for a subdued look.
fj.HighContrastStyle // HighContrastStyle uses bold and contrasting colors for better visibility.
fj.VintageStyle // VintageStyle uses muted tones reminiscent of old terminal displays.
fj.CyberpunkStyle // CyberpunkStyle mimics a futuristic neon cyberpunk aesthetic.
fj.OceanStyle // OceanStyle is inspired by oceanic hues and soft contrasts.
fj.FieryStyle // FieryStyle uses intense warm colors like flames.
fj.GalaxyStyle // GalaxyStyle uses space-themed colors with a starry effect.
fj.SunsetStyle // SunsetStyle mimics the colors of a sunset, using warm hues and deep purples.
fj.JungleStyle // JungleStyle draws inspiration from a dense jungle with deep greens and browns.
fj.MonochromeStyle // MonochromeStyle uses different shades of black and white for a simple, high-contrast theme.
fj.ForestStyle // ForestStyle uses deep greens and browns to create a natural, earthy look.
fj.IceStyle // IceStyle brings a cool, frosty aesthetic with blues and whites.
fj.RetroStyle // RetroStyle brings back the vibrant colors from older computer systems and arcade games.
fj.AutumnStyle // AutumnStyle uses rich oranges, reds, and browns, evoking the colors of fall.
fj.GothicStyle // GothicStyle uses darker colors with a moody atmosphere, ideal for dark themes.
fj.VaporWaveStyle // VaporWaveStyle embraces the retro aesthetics of vapor-wave, with bright neon and pastel colors.
fj.VampireStyle // VampireStyle brings dark and sinister colors, with a touch of red for a spooky theme.
fj.CarnivalStyle // CarnivalStyle is inspired by a fun, bright carnival atmosphere, full of vivid, exciting colors.
fj.SteampunkStyle // SteampunkStyle has a vintage industrial look with brass and copper colors.
fj.WoodlandStyle // WoodlandStyle blends earthy tones with deep forest greens and browns.
fj.CandyStyle // CandyStyle is bright, with pastel hues that resemble candy colors.
fj.TwilightStyle // TwilightStyle brings in dusky, cool tones reminiscent of dusk.
fj.EarthStyle // EarthStyle reflects natural earthy colors with muted greens and browns.
fj.ElectricStyle // ElectricStyle uses electric, bright neon colors for a futuristic vibe.
fj.WitchingHourStyle // WitchingHourStyle combines deep purples with dark greens for a magical look.
fj.MidnightStyle // MidnightStyle gives a mysterious and dark aesthetic, like a quiet midnight scene.
fj.RetroFutureStyle // RetroFutureStyle combines retro tones with a futuristic neon palette for a vintage-tech feel.
fj.ForestMistStyle // ForestMistStyle invokes the serene and cool vibes of a misty forest.
fj.PrismStyle // PrismStyle offers a colorful, dazzling light prism effect for a modern, energetic look.
fj.SpringStyle // SpringStyle brings the fresh, light colors of spring to life.
fj.DesertStyle // DesertStyle evokes the warmth and serenity of a desert landscape.
fj.SolarFlareStyle // SolarFlareStyle uses vibrant oranges and fiery reds, inspired by the intense heat of the sun.
fj.IceQueenStyle // IceQueenStyle reflects a cool, frosty appearance with icy blues and silvers.
fj.ForestGroveStyle // ForestGroveStyle brings earthy tones with a dense forest theme.
fj.AutumnLeavesStyle // AutumnLeavesStyle uses warm, fall-inspired hues like browns, reds, and golden yellows.
fj.VaporStyle // VaporStyle uses pastel tones and calming shades of pink, purple, and blue.
fj.SunsetBoulevardStyle // SunsetBoulevardStyle mimics the stunning colors of a sunset, featuring warm oranges, pinks, and purples.
fj.NeonCityStyle // NeonCityStyle is bold and energetic, with electrifying neons of blue, pink, and green.
fj.MoonlitNightStyle // MoonlitNightStyle gives a serene and calm atmosphere with cool blues and soft silvers.
fj.CandyShopStyle // CandyShopStyle features bright, sugary tones of pinks, blues, and yellows for a fun and sweet theme.
fj.UnderwaterStyle // UnderwaterStyle is inspired by the deep ocean, featuring calming blues and aquatic greens.
fj.OceanBreezeStyle // OceanBreezeStyle reflects the calm and refreshing hues of the ocean.
fj.CandyPopStyle // CandyPopStyle brings a playful and sweet color palette, like a candy store.
fj.NoirStyle // NoirStyle gives a film-noir inspired look with dark, moody colors.
fj.GalacticStyle // GalacticStyle evokes the mysterious vastness of outer space with deep, cosmic hues.
fj.VintagePastelStyle // VintagePastelStyle offers a retro aesthetic with soft, pastel tones for a gentle, nostalgic atmosphere.
fj.VintageFilmStyle // VintageFilmStyle is inspired by the golden era of cinema, featuring muted golds, sepias, and classic black.
fj.FireworksStyle // FireworksStyle captures the excitement of a night sky lit up by colorful fireworks, featuring bold reds, yellows, and purples.
fj.ArcticSnowStyle // ArcticSnowStyle brings the cool, crisp whites and icy blues of the arctic tundra into the design.
fj.ElectricVibeStyle // ElectricVibeStyle takes on high-energy neon tones with a touch of electric brightness.
fj.DesertSunsetStyle // DesertSunsetStyle brings warm and deep hues inspired by the desert landscape at sunset.
fj.PastelDreamStyle // PastelDreamStyle evokes a dreamy, soft pastel palette perfect for relaxed and whimsical visuals.
fj.TropicalVibeStyle // TropicalVibeStyle draws inspiration from lush tropical jungles, with bright and vibrant greens and yellows.
```

eg.

```go
package main

import (
	"fmt"

	"github.com/sivaosorg/fj"
)

var json []byte = []byte(`{"user":{"id":12345,"name":{"firstName":"John","lastName":"Doe"},"email":"john.doe@example.com","phone":"+1-555-555-5555","address":{"street":"123 Main St","city":"Anytown","state":"CA","postalCode":"12345","country":"USA"},"roles":[{"roleId":"1","roleName":"Admin","permissions":[{"permissionId":"101","permissionName":"View Reports","allowedActions":["view","download"]},{"permissionId":"102","permissionName":"Manage Users","allowedActions":["create","update","delete"]}]},{"roleId":"2","roleName":"Editor","permissions":[{"permissionId":"201","permissionName":"Edit Content","allowedActions":["create","edit","publish"]},{"permissionId":"202","permissionName":"View Analytics","allowedActions":["view"]}]}],"status":"active","createdAt":"2025-01-01T10:00:00Z","lastLogin":"2025-01-12T15:30:00Z"}}`)

func main() {
	ctx := fj.ParseBytes(json)
	fmt.Println(ctx.WithStringColored(fj.ElectricVibeStyle))
}
```

> Output

![ElectricVibeStyle](./assets/JSON-ElectricVibeStyle.png)
![TropicalVibeStyle](./assets/JSON-TropicalVibeStyle.png)

## Types

The result type encapsulates one of the following JSON types: `string`, `number`, `boolean`, or `null`. Arrays and objects are represented as their raw JSON forms. The struct for accessing a JSON value:

```go
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
```

eg.

```go
package main

import (
	"fmt"

	"github.com/sivaosorg/fj"
)

var json []byte = []byte(`{"user":{"id":12345,"name":{"firstName":"John","lastName":"Doe"},"email":"john.doe@example.com","phone":"+1-555-555-5555","address":{"street":"123 Main St","city":"Anytown","state":"CA","postalCode":"12345","country":"USA"},"roles":[{"roleId":"1","roleName":"Admin","permissions":[{"permissionId":"101","permissionName":"View Reports","allowedActions":["view","download"]},{"permissionId":"102","permissionName":"Manage Users","allowedActions":["create","update","delete"]}]},{"roleId":"2","roleName":"Editor","permissions":[{"permissionId":"201","permissionName":"Edit Content","allowedActions":["create","edit","publish"]},{"permissionId":"202","permissionName":"View Analytics","allowedActions":["view"]}]}],"status":"active","createdAt":"2025-01-01T10:00:00Z","lastLogin":"2025-01-12T15:30:00Z"}}`)

func main() {
	ctx := fj.ParseBytes(json).Get("user.name") // hold the context: fj.Context
	fmt.Println(ctx.Kind().String())            // JSON
	fmt.Println(ctx.Unprocessed())              // {"firstName":"John","lastName":"Doe"}
	fmt.Println(ctx.String())                   // {"firstName":"John","lastName":"Doe"}
	fmt.Println(ctx.Numeric())                  // 0
	fmt.Println(ctx.Index())                    // 27
	fmt.Println(ctx.Indexes())                  // []
}
```

Type represents the different possible types for a JSON value.

```go
string
false (bool)
true (bool)
number (float64, float32, int64, uint64)
null (nil)
```

```go
const (
	// Null is a constant representing a JSON null value.
	// In JSON, null is used to represent the absence of a value.
	Null Type = iota
	// False is a constant representing a JSON false boolean value.
	// In JSON, false is a boolean value that represents a negative or off state.
	False
	// Number is a constant representing a JSON number value.
	// In JSON, numbers can be integers or floating-point values.
	Number
	// String is a constant representing a JSON string value.
	// In JSON, strings are sequences of characters enclosed in double quotes.
	String
	// True is a constant representing a JSON true boolean value.
	// In JSON, true is a boolean value that represents a positive or on state.
	True
	// JSON is a constant representing a raw JSON block.
	// This type can be used to represent any valid JSON object or array.
	JSON
)
```

## Functions

```go
func (ctx Context) Kind() Type // Kind returns the JSON type of the Context.

func (ctx Context) Unprocessed() string // Unprocessed returns the raw, unprocessed JSON string for the Context.

func (ctx Context) Numeric() float64 // Numeric returns the numeric value of the Context, if applicable.

func (ctx Context) Index() int // Index returns the index of the unprocessed JSON value in the original JSON string.

func (ctx Context) Indexes() []int // Indexes returns a slice of indices for elements matching a path containing the '#' character.

// String returns a string representation of the Context value.
// The output depends on the JSON type of the Context:
//   - For `False` type: Returns "false".
//   - For `True` type: Returns "true".
//   - For `Number` type: Returns the numeric value as a string.
//     If the numeric value was calculated, it formats the float value.
//     Otherwise, it preserves the original unprocessed string if valid.
//   - For `String` type: Returns the string value.
//   - For `JSON` type: Returns the raw unprocessed JSON string.
//   - For other types: Returns an empty string.
func (ctx Context) String() string

func (ctx Context) StringColored() string // StringColored returns a colored string representation of the Context value.

func (ctx Context) WithStringColored(style *unify4g.Style) string // WithStringColored applies a customizable colored styling to the string representation of the Context value.

func (ctx Context) Bool() bool // Bool converts the Context value into a boolean representation.

func (ctx Context) Int64() int64 // Int64 converts the Context value into an integer representation (int64).

func (ctx Context) Uint64() uint64 // Uint64 converts the Context value into an unsigned integer representation (uint64).

func (ctx Context) Float64() float64 // Float64 converts the Context value into a floating-point representation (float64).

func (ctx Context) Float32() float32 // Float32 converts the Context value into a floating-point representation (Float32).

func (ctx Context) Time() time.Time // Time converts the Context value into a time.Time representation.

func (ctx Context) WithTime(format string) (time.Time, error) // WithTime parses the Context value into a time.Time representation using a custom format.

func (ctx Context) Array() []Context // Array returns an array of `Context` values derived from the current `Context`.

func (ctx Context) IsObject() bool // IsObject checks if the current `Context` represents a JSON object.

func (ctx Context) IsArray() bool // IsArray checks if the current `Context` represents a JSON array.

func (ctx Context) IsBool() bool // IsBool checks if the current `Context` represents a JSON boolean value.

func (ctx Context) Exists() bool // Exists returns true if the value exists (i.e., it is not Null and contains data).

func (ctx Context) Value() interface{} // Value returns the corresponding Go type for the JSON value represented by the Context.

func (ctx Context) Map() map[string]Context // Map returns a map of values extracted from a JSON object.

func (ctx Context) Foreach(iterator func(key, value Context) bool) // Foreach iterates through the values of a JSON object or array, applying the provided iterator function.

func (ctx Context) Get(path string) Context // Get searches for a specified path within a JSON structure and returns the corresponding result.

func (ctx Context) GetMul(path ...string) []Context // GetMul searches for multiple paths within a JSON structure and returns a slice of results.

// Path returns the original fj path for a Result where the Result came
// from a simple query path that returns a single value.
func (ctx Context) Path(json string) string

// Paths returns the original fj paths for a Result where the Result came
// from a simple query path that returns an array.
func (ctx Context) Paths(json string) []string

// Less compares two Context values (tokens) and returns true if the first token is considered less than the second one.
// It performs comparisons based on the type of the tokens and their respective values.
// The comparison order follows: Null < False < Number < String < True < JSON.
// This function also supports case-insensitive comparisons for String type tokens based on the caseSensitive parameter.
func (ctx Context) Less(token Context, caseSensitive bool) bool
```

#### Array()

The `ctx.Array()` method produces an array of values. If the result corresponds to a non-existent value, it will return an empty array. If the result is not a JSON array, a single-element array containing the result will be returned.

#### Value()

The `ctx.Value()` function retrieves an `interface{}`, requiring a type assertion and is typically one of these Go types:

```go
boolean >> bool
number  >> float64, float32, int64, uint64
string  >> string
null    >> nil
array   >> []interface{}
object  >> map[string]interface{}
```

#### 64-bit un/integers

The `ctx.Int64()` and `ctx.Uint64()` functions can handle the full 64-bit range, enabling support for large JSON integer values (Refer this link: [Min Safe Integer](https://tc39.es/ecma262/#sec-number.min_safe_integer); [Max Safe Integer](https://tc39.es/ecma262/#sec-number.max_safe_integer)).

```go
ctx.Int64() int64    // -9223372036854775808 to 9223372036854775807
ctx.Uint64() uint64   // 0 to 18446744073709551615
```
