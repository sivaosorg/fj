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

Given the `JSON` string

```json
{
  "id": "http://subs/base-sample-schema.json",
  "schema": "http://json-schema.org/draft-07/schema#",
  "title": "Submissions Sample Schema",
  "description": "A base validation sample schema",
  "version": "1.0.0",
  "author": "subs",
  "type": "object",
  "properties": {
    "alias": {
      "description": "An unique identifier in a submission.",
      "type": "string",
      "minLength": 1
    },
    "title": {
      "description": "Title of the sample.",
      "type": "string",
      "minLength": 1
    },
    "description": {
      "description": "More extensive free-form description.",
      "type": "string",
      "minLength": 1
    },
    "taxonId": {
      "type": "integer",
      "minimum": 1
    },
    "taxon": {
      "type": "string",
      "minLength": 1
    },
    "releaseDate": {
      "type": "string",
      "format": "date"
    },
    "attributes": {
      "description": "Attributes for describing a sample.",
      "type": "object",
      "properties": {},
      "patternProperties": {
        "^.*$": {
          "$ref": "definitions-schema.json#/definitions/attributes_structure"
        }
      }
    },
    "sampleRelationships": {
      "$ref": "definitions-schema.json#/definitions/sampleRelationships"
    }
  },
  "required": ["alias", "taxonId", "releaseDate"],
  "oneOf": [
    {
      "required": ["alias", "team"]
    },
    {
      "required": ["accession"]
    }
  ],
  "animals": [
    {
      "name": "Meowsy",
      "species": "cat",
      "foods": {
        "likes": ["tuna", "catnip"],
        "dislikes": ["ham", "zucchini"]
      }
    },
    {
      "name": "Barky",
      "species": "dog",
      "foods": {
        "likes": ["bones", "carrots"],
        "dislikes": ["tuna"]
      }
    },
    {
      "name": "Purrpaws",
      "species": "cat",
      "foods": {
        "likes": ["mice"],
        "dislikes": ["cookies"]
      }
    }
  ],
  "stock": [
    {
      "company": "3M",
      "description": "3M, based in Minnesota, may be best known for its Scotch tape and Post-It Notes, but it also produces sand paper, adhesives, medical products, computer screen filters, food safety items, stationery products and many products used in automotive, marine, and aircraft industries.",
      "initial_price": 44.28,
      "price_2002": 56.27,
      "price_2007": 95.85,
      "symbol": "MMM"
    },
    {
      "company": "Amazon.com",
      "description": "Amazon.com, Inc. is an online retailer in North America and internationally. The company serves consumers through its retail Web sites and focuses on selection, price, and convenience. It also offers programs that enable sellers to sell their products on its Web sites, and their own branded Web sites. In addition, the company serves developer customers through Amazon Web Services, which provides access to technology infrastructure that developers can use to enable virtually various type of business. Further, it manufactures and sells the Kindle e-reader. Founded in 1994 and headquartered in Seattle, Washington.",
      "initial_price": 89.38,
      "price_2002": 17.01,
      "price_2007": 93.43,
      "symbol": "AMZN"
    },
    {
      "company": "Campbell Soup",
      "description": "Campbell Soup is a worldwide food company, offering condensed and ready-to-serve soups; broth, stocks, and canned poultry; pasta sauces; Mexican sauces; canned pastas, gravies, and beans; juices and beverages; and tomato juices. Its customers include retail food chains, mass discounters, mass merchandisers, club stores, convenience stores, drug stores and other retail, and commercial and non-commercial establishments. Campbell Soup Company was founded in 1869 and is headquartered in Camden, New Jersey.",
      "initial_price": 37.0,
      "price_2002": 22.27,
      "price_2007": 36.4,
      "symbol": "CPB"
    },
    {
      "company": "Disney",
      "description": "The Walt Disney Company, founded in 1923, is a worldwide entertainment company, with movies, cable networks, radio networks, movie production, musical recordings and live stage plays. Disney also operates Walt Disney World in Florida and Disneyland in California, Disney Cruise Line, and international Disney resorts. Disney owns countless licenses and literary properties and publishes books and magazines.",
      "initial_price": 40.68,
      "price_2002": 15.24,
      "price_2007": 35.47,
      "symbol": "DIS"
    },
    {
      "company": "Dow Chemical",
      "description": "The Dow Chemical Company manufactures raw materials that go into consumer products and services. These materials include food and pharmaceutical ingredients, electronic displays, semiconductor packaging, water purification, insulation, adhesives, pest control, polyurethane, polystyrene (goes into plastics), and crude-oil based raw materials. Dow was founded in 1897 and is based in Midland, Michigan.",
      "initial_price": 38.83,
      "price_2002": 27.65,
      "price_2007": 44.67,
      "symbol": "DOW"
    },
    {
      "company": "Exxon Mobil",
      "description": "Exxon Mobil engages in the exploration and production of crude oil and natural gas, and manufacture of petroleum products. The company manufactures commodity petrochemicals. The company has operations in the United States, Canada/South America, Europe, Africa, Asia, and Australia/Oceania. Exxon Mobil Corporation was founded in1870 and is based in Irving, Texas.",
      "initial_price": 39.0,
      "price_2002": 32.82,
      "price_2007": 91.36,
      "symbol": "XOM"
    },
    {
      "company": "Ford",
      "description": "Ford Motor Co. develops, manufactures, sells and services vehicles and parts worldwide. Ford sells cars and trucks primarily under the Ford and Lincoln brands. It sells to consumers (through retail dealers) and to rental car companies, leasing companies, and governments. Ford also provides maintenance and repair services. Ford also offers financing to vehicle purchasers. Ford was founded in 1903 and is based in Dearborn, Michigan.",
      "initial_price": 27.34,
      "price_2002": 9.63,
      "price_2007": 8.37,
      "symbol": "F"
    },
    {
      "company": "The Gap",
      "description": "The Gap, Inc. sells retail clothing, accessories and personal care products globally under the brand names Gap, Old Navy, Banana Republic, Piperlime, Athleta and Intermix. Products include sports apparel, casual clothing, sleepwear, footwear and infants\u2019 and children\u2019s clothing. The company has company-owned stores as well as franchise stores, online stores and catalogs. The Gap was founded in 1969 and is headquartered in San Francisco, California.",
      "initial_price": 46.0,
      "price_2002": 11.56,
      "price_2007": 18.9,
      "symbol": "GPS"
    },
    {
      "company": "General Mills",
      "description": "General Mills manufactures and sells consumer foods worldwide. Products include cereals, frozen vegetables, dough, dessert and baking mixes, frozen pizzas, grains, fruits, ice creams and organic products. It sells to grocery stores as well as commercial food service distributors, restaurants and convenience stores. General Mills was founded in 1928 and is based in Minneapolis, Minnesota.",
      "initial_price": 15.59,
      "price_2002": 22.1,
      "price_2007": 28.76,
      "symbol": "GIS"
    }
  ],
  "bank": [
    {
      "isActive": false,
      "balance": "$1,404.23",
      "age": 26,
      "eyeColor": "blue",
      "name": "Stark Jenkins",
      "gender": "male",
      "company": "HINWAY",
      "email": "starkjenkins@hinway.com",
      "phone": "+1 (943) 542-3591",
      "address": "766 Cooke Court, Dunbar, Connecticut, 9512"
    },
    {
      "isActive": false,
      "balance": "$1,247.08",
      "age": 36,
      "eyeColor": "green",
      "name": "Odonnell Rollins",
      "gender": "male",
      "company": "NEXGENE",
      "email": "odonnellrollins@nexgene.com",
      "phone": "+1 (810) 521-2350",
      "address": "210 Pleasant Place, Lloyd, Mississippi, 1636"
    },
    {
      "isActive": false,
      "balance": "$2,284.89",
      "age": 20,
      "eyeColor": "brown",
      "name": "Rachelle Chang",
      "gender": "female",
      "company": "VERAQ",
      "email": "rachellechang@veraq.com",
      "phone": "+1 (955) 564-2002",
      "address": "220 Drew Street, Ventress, Puerto Rico, 8432"
    },
    {
      "isActive": true,
      "balance": "$1,624.60",
      "age": 39,
      "eyeColor": "brown",
      "name": "Davis Wade",
      "gender": "female",
      "company": "ASSISTIX",
      "email": "daviswade@assistix.com",
      "phone": "+1 (836) 432-2542",
      "address": "532 Amity Street, Yukon, Palau, 3561"
    },
    {
      "isActive": true,
      "balance": "$3,818.97",
      "age": 23,
      "eyeColor": "green",
      "name": "Oneill Everett",
      "gender": "male",
      "company": "INCUBUS",
      "email": "oneilleverett@incubus.com",
      "phone": "+1 (958) 522-2724",
      "address": "273 Temple Court, Shelby, Georgia, 8682"
    },
    {
      "isActive": false,
      "balance": "$3,243.63",
      "age": 21,
      "eyeColor": null,
      "name": "Dalton Waters",
      "gender": "male",
      "company": "OVATION",
      "email": "daltonwaters@ovation.com",
      "phone": "+1 (899) 464-3878",
      "address": "909 Wyona Street, Adelino, Hawaii, 6449"
    }
  ]
}
```

## Syntax

A path is a string format used to define a pattern for efficiently retrieving values from a JSON structure.

### Path

A `fj` path is designed to be represented as a sequence of elements divided by a `.` symbol.

In addition to the `.` symbol, several other characters hold special significance, such as `|`, `#`, `@`, `\`, `*`, `!`, and `?`.

### Access Values - Object

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

### Access Values - Array

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
