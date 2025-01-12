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
> id #output: "http://subs/base-sample-schema.json"
> properties.alias.description #output: "An unique identifier in a submission."
> properties.alias.minLength #output: 1
> required #output: ["alias", "taxonId", "releaseDate"]
> required.0 #output: "alias"
> required.1 #output: "taxonId"
> oneOf.0.required #output: ["alias", "team"]
> oneOf.0.required.1 #output: "team"
> properties.sampleRelationships #output: { "$ref": "definitions-schema.json#/definitions/sampleRelationships" }
```

- **Wildcards**: A key can include special wildcard symbols like `*` and `?`. The `*` matches any sequence of characters (including none), while `?` matches exactly one character.

```shell
> anim*ls.1.name #output: "Barky"
> *nimals.1.name #output: "Barky"
```

- **Escape Character**: Characters with special meanings, like `.`, `*`, and `?`, can be escaped using the `\` symbol.

```shell
> properties.alias\.description #output: "An unique identifier in a submission."
```

### Access Values - Array

The `#` symbol enables navigation within JSON arrays. To retrieve the length of an array, simply use the `#` on its own.

```shell
> animals.# #output: 3 (length of an array)
> animals.#.name #output: ["Meowsy","Barky","Purrpaws"]
```

### Queries

You can also search an array for the first match by using `#(...)`, or retrieve all matches with `#(...)#`.
Queries support comparison operators such as `==`, `!=`, `<`, `<=`, `>`, `>=`, along with simple pattern matching operators `%` (like) and `!%` (not like).

```shell
> stock.#(price_2002==56.27).symbol #output: "MMM"
> stock.#(company=="Amazon.com").symbol #output: "AMZN"
> stock.#(initial_price>=10)#.symbol #output: ["MMM","AMZN","CPB","DIS","DOW","XOM","F","GPS","GIS"]
> stock.#(company%"D*")#.symbol #output: ["DIS","DOW"]
> stock.#(company!%"D*")#.symbol #output: ["MMM","AMZN","CPB","XOM","F","GPS","GIS"]
> stock.#(company!%"F*")#.symbol #output: ["MMM","AMZN","CPB","DIS","DOW","XOM","GPS","GIS"]
> stock.#(description%"*stores*")#.symbol #output: ["CPB","GPS","GIS"]
> required.#(%"*as*")# #output: ["alias","releaseDate"]
> required.#(%"*as*") #output: "alias"
> required.#(!%"*s*") #output: "taxonId"
> animals.#(foods.likes.#(%"*a*"))#.name #output: ["Meowsy","Barky"]
```
