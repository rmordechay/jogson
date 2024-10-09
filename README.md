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
// Check if a key exists
var keyExists bool = mapper.Object.Has("children")

// Get object's keys
var keys []string = mapper.Object.Keys()

// Get object's values
var values []jsonmapper.Mapper = mapper.Object.Values()

// Iterating over an object with key, value pair
children := mapper.Object.GetObject("children")
for key, child := range children.Elements() {
    fmt.Println("Child name:", key)               // Rachel, Sara
    fmt.Println(child.Object.GetInt("age"))       // 15, 19
    fmt.Println(child.Object.GetBool("is_funny")) // false, true
}
```

### Primitive Types
Getting primitive types - `string`, `int`, `float64` or `bool` - is similar both for object and array and only differ
in the parameter type (objects take `string` as key and arrays take `int` as index)
#### From Object
```go
object := mapper.Object
// string 
var name string = object.GetString("name") // Jason
// int 
var age int = object.GetInt("age") // 15
// float64 
var height float64 = object.GetFloat("height") // 1.87
// bool 
var isFunny bool = object.GetBool("is_funny") // false
```
#### From Array
```go
array := mapper.Array
// string 
var s string = array.GetString(0)
// int 
var i int = array.GetInt(2)
// float64 
var f float64 = array.GetFloat(5)
// bool 
var b bool = array.GetBool(7)
```


### Arrays
```go
// Get array
array := mapper.Object.GetArray("features")

// Get array length
arrayLen := array.Length() // 2

// Get an element of array by index
secondElement := array.Get(1)

// Iterating over an array
array := mapper.Object.GetArray("features")
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
