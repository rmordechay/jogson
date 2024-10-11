# JSON Mapper Library for Go

[![GoDoc](https://pkg.go.dev/badge/badge)](https://pkg.go.dev/github.com/rmordechay/jsonmapper)

A simple Go library to simplify working with JSON without the need to define structs.

* [Installation](#Installation)
* [Create a Mapper](#Create-a-mapper)
* [Read from JSON](#Read-from-JSON)
    * [Objects](#Objects)
    * [Arrays](#Arrays)
    * [Scalars](#Scalars)
    * [Types](#types)
    * [Find Elements](#Find-Elements)
    * [Get as JSON String](#get-as-json-string)
    * [Error Handling](#Error-handling)
* [Write to JSON](#Write-to-JSON)
    * [Write Object](#Write-object)
    * [Write Array](#Write-array)
* [Design](#Design)
  * [JsonMapper](#JsonMapper)
  * [JsonObject](#JsonObject)
  * [JsonArray](#JsonArray)
  * [Good to know](#good-to-know)

## Installation

To install the library, use:
```bash
go get github.com/rmordechay/jsonmapper
```

## Create a Mapper

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

## Read from JSON

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

### Objects

```go
object := mapper.AsObject

// Check if a key exists
var keyExists bool = object.Has("children")

// Check if the object is empty
var objectEmpty bool = object.IsEmpty()

// Check if a key contains a null value
var valueNull bool = object.IsNull("children")

// Get the object's size
var keyExists bool = object.Length()

// Get object's keys
var keys []string = object.Keys()

// Get object's values
values := object.Values()

// Iterating over an object with key, value pair
children := object.GetObject("children")
for key, child := range children.Elements() {
    fmt.Println("Child name:", key)                 // Rachel, Sara
    fmt.Println(child.AsObject.GetInt("age"))         // 15, 19
    fmt.Println(child.AsObject.GetBool("is_funny"))   // false, true
}
```

### Arrays

```go
// Get array
array := mapper.AsObject.GetArray("features")

// Get array length
arrayLen := array.Length() // 2

// Get an element of array by index
secondElement := array.Get(1)

// Iterating over an array
for _, feature := range array.Elements() {
    fmt.Println(feature.AsString) // tall, ...
}

// Get as a slice of string
var stringArray []string = array.AsStringArray()

// Get as a slice of int
var intArray []int = array.AsIntArray()

// Get as a slice of float
var floatArray []float = array.AsFloatArray()

// Get as a slice of JsonArray
var nestedArray []JsonArray = array.As2DArray()

// Get as a slice of JsonObject
var objectArray []JsonObject = array.AsObjectArray()
```

### Scalars

Getting scalars - `string`, `int`, etc. - is similar both for object and array and only differ
in the parameter type (objects take `string` as key and arrays take `int` as index)

#### From Object

```go
object := mapper.AsObject
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
array := mapper.AsArray
// string 
var s string = array.GetString(0)
// int 
var i int = array.GetInt(2)
// float64 
var f float64 = array.GetFloat(5)
// bool 
var b bool = array.GetBool(7)
```

#### Time

To get a string as `time.Time`

```go
var birthday time.Time = array.GetTime(0)
var birthday string = object.GetTime("birthday".Format(time.RFC3339)) // 1981-10-08T00:00:00Z
```

The mapper will try to format the string against different time formats to increase the change of correct parsing. The following formats are supported:

`time.RFC3339` `time.RFC850` `time.RFC822` `time.RFC822Z` `time.RFC1123` `time.RFC1123Z` `time.RFC3339Nano` `time.ANSIC` `time.UnixDate` `time.RubyDate` `time.Layout` `time.Kitchen` `time.Stamp` `time.StampMilli` `time.StampMicro` `time.StampNano` `time.DateTime` `time.DateOnly` `time.TimeOnly`

### Types

```go
fmt.Println(mapper.IsObject)    // true
fmt.Println(mapper.IsBool)      // false
fmt.Println(mapper.IsInt)       // false
fmt.Println(mapper.IsFloat)     // false
fmt.Println(mapper.IsString)    // false
fmt.Println(mapper.IsArray)     // false
fmt.Println(mapper.IsNull)      // false
```

### Find Elements

You can search for a nested element.

```go
element := mapper.AsObject.Find("Rachel")
fmt.Println(element.IsObject) // true 
fmt.Println(element.Has("is_funny")) // true 
```

### Get as JSON String

You can get a string from every JSON element which is a valid JSON

```go
fmt.Println(mapper.AsObject.String())
// output: {"age":43,"children":{"Rachel":{"age":15,"is_funny":false},"Sara":{"age":19,"is_funny":true}},"features":["tall","blue eyes"],"is_funny":false,"name":"Jason"}

fmt.Println(mapper.AsObject.Get("children").String())
// output: {"Rachel":{"age":15,"is_funny":false},"Sara":{"age":19,"is_funny":true}}
```

or with pretty string

```go
fmt.Println(mapper.AsObject.Get("children").PrettyString())
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

### Error Handling

## Write to JSON

To write a JSON object or array is as simple as reading from it.

### Write Object

```go
// Create a new object
obj := jsonmapper.NewObject()
obj.AddKeyValue("name", "Chris")
fmt.Println(obj.String()) // {"name":"Chris"}
```

### Write Array

```go
// Create a new array
arr := jsonmapper.NewArray()
arr.AddElement(15)
arr.AddElement(19)
fmt.Println(arr.String()) // [15,19]
```

## Design

There are 3 structs that are important to know when working with the library
* `JsonMapper` struct for generic JSON.
* `JsonObject` represents JSON object.
* `JsonArray`  represents JSON array.

#### JsonMapper

Most likely, when you read JSON, you would start with a `JsonMapper` (see [Create a Mapper](#Create-a-mapper)).
`JsonMapper` is a struct that hold your JSON data. It has several `AsX` and `IsX` fields with which you can get the data and check
the type, respectively. For example, if your data is a JSON object, you can call `JsonMapper.AsObject`, or if it's a string, `JsonMapper.AsString`.
If you don't know the type and want to check it dynamically, you can use `JsonMapper.IsObject`, `JsonMapper.IsString`, etc. If the underlying data
is null, then `IsNull` will be set to true. `JsonMapper` is also returned in cases where the return type can be any JSON
type. For example, `JsonArray.Elements()` returns a slice `[]JsonMapper` which you can iterate over or query specific elements.

#### JsonObject

`JsonObject` holds JSON object data and has different methods to read from JSON object or write to it. The simplest way to get a `JsonObject`
is to call `JsonMapper.AsObject`. Once you have an instance, you get use the various methods to get data. 

#### JsonArray


### Good to know

#### `As`, `Is` and `Get`
The prefixes, `As`, `Is` and `Get` have similar semantics across the library.
* `IsX`: checks for the value's time. For example `JsonMapper.IsBool`
* `AsX`: converts the current data to other type representation. For example, `JsonArray.AsStringArray()` converts JsonArray to `[]string`.
* `GetX`: Fetches the data, usually with some sort of search in the underlying data.

#### `JsonObject` and `JsonArray` Similarity
`JsonObject` and `JsonArray` have very similar methods, both in naming and semantics, but with 2 differences
  * Input: 
    * `JsonObject`'s methods mostly asks for `string` as the key 
    * `JsonArray`'s methods mostly asks for `int` as the index
  * Output:
    * `JsonObject`'s methods mostly return a map
    * `JsonArray`'s methods mostly return a slice