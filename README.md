# JSON Mapper Library for Go

[![GoDoc](https://pkg.go.dev/badge/badge)](https://pkg.go.dev/github.com/rmordechay/jsonmapper)

A simple Go library to simplify working with JSON without the need to define structs.

* [Installation](#Installation)
* [Creating a mapper](#Creating-a-mapper)
* [Reading from JSON](#Reading-from-JSON)
  * [Creating a mapper](#Creating-a-mapper)
  * [Check Types](#Check-Types)
  * [Objects](#Objects)
  * [Get Values](#get-primitive-values)
  * [Arrays](#Arrays)
  * [Find Elements](#Find-Elements)
  * [Print JSON](#Get-JSON-as-string)
* [Writing to JSON](#Writing-to-JSON)
  * [Write elements](#Write-elements)


## Installation
To install the library, use:

```bash
go get github.com/rmordechay/jsonmapper
```

### Creating a mapper
There are multiple ways to create a mapper.
```go
// From bytes
mapper, err := jsonmapper.FromBytes([]byte(jsonBytes))
// From string
mapper, err := jsonmapper.FromString(jsonString)
// From struct
mapper, err := jsonmapper.FromStruct(jsonStruct)
// From file
mapper, err := jsonmapper.FromFile(jsonFilePath)
```

## Reading from JSON
Once you have the `mapper`, you can read the data easily. Consider the following JSON
```go
jsonString := `{
    "name": "Jason",
    "age": 43,
    "height": 1.87,
    "is_funny": false,
    "birthday": "1981-10-08",
    "features": ["tall", "blue eyes"],
    "children": {
        "Rachel": {"age": 15, "is_funny": false}, 
        "Sara":   {"age": 19, "is_funny": true}
    }
}`
```

### Check Types
```go
fmt.Println(mapper.IsObject)    // true
fmt.Println(mapper.IsBool)      // false
fmt.Println(mapper.IsInt)       // false
fmt.Println(mapper.IsFloat)     // false
fmt.Println(mapper.IsString)    // false
fmt.Println(mapper.IsArray)     // false
fmt.Println(mapper.IsNull)      // false
```

### Objects
```go
// Get a value by key
element := mapper.Object.Get("age")

// Check if a key exists
keyExists := mapper.Object.Has("children")

// Get object's keys
keys := mapper.Object.Keys()

// Get object's values
values := mapper.Object.Values()

// Iterating over an object with key, value pair
children := mapper.Object.Get("children").Object
for key, child := range children.Elements() {
    fmt.Println("Child name:", key)                  // Rachel, Sara
    fmt.Println(child.Object.Get("age").AsInt)       // 15, 19
    fmt.Println(child.Object.Get("is_funny").AsBool) // false, true
}
```

### Get Primitive Values
There are two ways to get a primitive value from an object. Either with the underlying Mapper or directly 
if you already know the type 
```go
object := mapper.Object

var name1 string = object.Get("name").AsString // Jason
var name2 string = object.GetString("name") // Jason

var age1 int = object.Get("age").AsInt // 15
var age2 int = object.GetInt("age") // 15

var height1 float64 = object.Get("height").AsFloat // 1.87
var height2 float64 = object.GetFloat("height") // 1.87

var isFunny1 bool = object.Get("is_funny").AsBool // false
var isFunny2 bool = object.GetBool("is_funny") // false
```

### Arrays
```go
// Get array length
arrayLen := mapper.Array.Length()

// Get an element of array by index
element := mapper.Array.Get(2)

// Iterating over an array
array := mapper.Object.Get("features").Array
for _, feature := range array.Elements() {
    fmt.Println(feature.AsString) // tall, ...
}
```

### Find Elements
You can search for a nested element. 
```go
element := mapper.Object.Find("Rachel")
fmt.Println(element.IsObject)         // true 
fmt.Println(element.Has("is_funny"))  // true 
```

### Get JSON as string
You can get a string from every JSON element which is a valid JSON
```go
fmt.Println(mapper.Object.String())
// output: {"age":43,"children":{"Rachel":{"age":15,"is_funny":false},"Sara":{"age":19,"is_funny":true}},"features":["tall","blue eyes"],"is_funny":false,"name":"Jason"}

fmt.Println(mapper.Object.Get("children").String())
// output: {"Rachel":{"age":15,"is_funny":false},"Sara":{"age":19,"is_funny":true}}
```
or with pretty string 
```go
fmt.Println(mapper.Object.Get("children").PrettyString())
// output:
// {
//    "Rachel": {
//        "age": 15,
//        "is_funny": false
//    },
//    "Sara": {
//        "age": 19,
//        "is_funny": true
//    }
// }
```

## Writing to JSON
To write a JSON object or array is as simple as reading from it.
To create an object or array from scratch, you can use
```go
// Object
obj := jsonmapper.CreateEmptyObject()
// Array
array := jsonmapper.CreateEmptyArray()
```

### Write elements
