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
  ]
}
```

## Syntax

A path is a string format used to define a pattern for efficiently retrieving values from a JSON structure.

## Path

A `fj` path is designed to be represented as a sequence of elements divided by a `.` symbol.

In addition to the `.` symbol, several other characters hold special significance, such as `|`, `#`, `@`, `\`, `*`, `!`, and `?`.

## Access values

- **Basic**: in most situations, you'll simply need to access values using the object name or array index.

```shell
> id # "http://subs/base-sample-schema.json"
> properties.alias.description # "An unique identifier in a submission."
> properties.alias.minLength # 1
> required # ["alias", "taxonId", "releaseDate"]
> required.0 # "alias"
> required.1 # "taxonId"
> oneOf.0.required # ["alias", "team"]
> oneOf.0.required.1 # "team"
```

- **Wildcards**: A key can include special wildcard symbols like `*` and `?`. The `*` matches any sequence of characters (including none), while `?` matches exactly one character.

```shell
> anim*ls.1.name # "Barky"
```
