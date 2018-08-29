# Usage

# Command Line
```bash
Usage: go-mana-test [OPTIONS] COMMAND [arg...]

Making APIs Nice Again - Testing Framework
                  
Options:          
  -v, --verbose   Verbose mode
                  
Commands:         
  version         Shows version info
  test            Run tests
                  
Run 'go-mana-test COMMAND --help' for more information on a command.
```

#### Testing
Testing is done by pointing to your [project file](#project-file).
```bash
go-mana-test test ./exampleproj/project.yml
```

You can add a verbose flag to output more information.
```bash
go-mana-test -v test ./exampleproj/project.yml
```

You can also do a dry run, which will just collect and validate test files but not actually run them.
```bash
go-mana-test test ./exampleproj/project.yml -d
```

# Project File
The project file setups your project. Here is an example project file.
```yaml
name: My Sample Project
tests: ./tests
globals:
  - key: USER_AGENT
    value: go-mana-test
```
 - name: Defines the name of your project.
 - tests: Defines the path to your [test files](#test-file), relative to the path this project file is in.
 - globals: Defiles key/value global [variables](#variables) that can be later used in your tests.
 
# Test File
The test file defines a single test for your project. Here is an example of a test file.
```yaml
name: Update User Profile
index: 100
url: "https://api.sample.com/v1/profile"
request.method: POST
request.headers:
  - key: Content-Type
    value: application/json
  - key: User-Agent
    value: "{{globals.USER_AGENT}}"
request.body:
  username: john.doe
cache:
  - name: user_profile.username
    value: response.body.json.username
checks:
  - name: Proper response code
    check: response.code
    value: 200
  - name: Proper body for updated username
    check: response.body.json.username
    value: "{{cache.body.username}}"
```
 - name: Defines the name of your test.
 - url: Defines the URL to use in the test. This may make use of [variables](#variables).
 - method: Defines the HTTP method to use for the test.
 - index: Defines the index of the test. Lower value indexes will be run before higher value indexes.
 - request.headers: Defines the HTTP headers to send in the test in key/value format.
 - request.body: Defines the body to send in the test. If you are sending JSON [https://www.json2yaml.com](https://www.json2yaml.com/) is a nice tool to help you convert JSON to YAML.
 - cache: Defines [cache](#cache) to save from this test.
 - checks: Defines the checks to validate in this test.
 
 # Test Cache
Test cache is run before the [test checks][#test-checks]. Caching allows you to cache certain values that can cary onto the checks in your test file, or across all test files.

The test cache uses key/value methods to store data. 

 - name: Defines the name of your cache. Can be any string.
 - value: Defines the value to cache. Can be any of the following:
   - response.body.json: Refers the json body of the response. Anything following this prefix this will [query methods of the json](#json-query) body.

 # Test Checks
Test checks are used to validate results of the test. 

 - name: Defines the name of your test. Can be any string. 
 - check: Defines the check to use. Can be any of the following:
   - response.code: References response HTTP status code.
   - response.body.json: Refers the json body of the response. Anything following this prefix this will [query methods of the json](#json-query) body. 
 - value: Defines the value to use in the test. This may make use of [variables](#variables).

# Variables
Variable substitution is helpful for certain properties of the test and project files. Inside of YAML files, replacements must always be enclosed inside of `""`. The full syntax is:
```yaml
"{{variable.to.use}}
```

Below are some of the variables that can be used:

#### Globals
Global variables are defined in the [project file](#project-file). The full syntax is:
```yaml
"{{globals.MY_GLOBAL_VARIABLE}}"
```

So for the sample [project file](#project-file) above, referencing the ser Agent would be:
```yaml
"{{globals.USER_AGENT}}"
```

#### Cache
Cache variables are defined in the [test file](#test-cache). The full syntax is:
```yaml
"{{cache.MY_CACHE_VARIABLE}}" 
```

So for the sample [test file](#test-file) above, referencing the Username would be:
```yaml
"{{cache.user_profile.username}}"
```

#### Random String
Generates a random alpha-numeric string of a certain length. The syntax to generate a string of length 50 is:
```yaml
"{{rand.string.50}}"
```

#### Random String
Generates a random lowercase alpha-numeric string of a certain length. The syntax to generate a string of length 60 is:
```yaml
"{{rand.string.lower.60}}"
```

#### Random String
Generates a random uppercase alpha-numeric string of a certain length. The syntax to generate a string of length 70 is:
```yaml
"{{rand.string.upper.70}}"
```

#### Random Number
Generates a random number in between a certain range. The syntax to generate a number between 1 and 100 is:
```yaml
"{{rand.num.1.100}}"
```

# JSON Query

#### Path Syntax
A path is a series of keys separated by a dot.
A key may contain special wildcard characters '\*' and '?'.
To access an array value use the index as the key.
To get the number of elements in an array or to access a child path, use the '#' character.
The dot and wildcard characters can be escaped with '\\'.

```json
{
  "name": {"first": "Tom", "last": "Anderson"},
  "age":37,
  "children": ["Sara","Alex","Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age": 44},
    {"first": "Roger", "last": "Craig", "age": 68},
    {"first": "Jane", "last": "Murphy", "age": 47}
  ]
}
```
```
"name.last"          >> "Anderson"
"age"                >> 37
"children"           >> ["Sara","Alex","Jack"]
"children.#"         >> 3
"children.1"         >> "Alex"
"child*.2"           >> "Jack"
"c?ildren.0"         >> "Sara"
"fav\.movie"         >> "Deer Hunter"
"friends.#.first"    >> ["Dale","Roger","Jane"]
"friends.1.last"     >> "Craig"
```

You can also query an array for the first match by using `#[...]`, or find all matches with `#[...]#`. 
Queries support the `==`, `!=`, `<`, `<=`, `>`, `>=` comparison operators and the simple pattern matching `%` operator.

```json
friends.#[last=="Murphy"].first    >> "Dale"
friends.#[last=="Murphy"]#.first   >> ["Dale","Jane"]
friends.#[age>45]#.last            >> ["Craig","Murphy"]
friends.#[first%"D*"].last         >> "Murphy"
```

#### JSON Lines
There's support for [JSON Lines](http://jsonlines.org/) using the `..` prefix, which treats a multilined document as an array. 

For example:

```json
{"name": "Gilbert", "age": 61}
{"name": "Alexa", "age": 34}
{"name": "May", "age": 57}
{"name": "Deloise", "age": 44}
```

```
..#                   >> 4
..1                   >> {"name": "Alexa", "age": 34}
..3                   >> {"name": "Deloise", "age": 44}
..#.name              >> ["Gilbert","Alexa","May","Deloise"]
..#[name="May"].age   >> 57
```
