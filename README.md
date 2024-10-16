# Jogson - JSON Mapper Library for Go

![Go Version](https://img.shields.io/badge/go%20version-%3E=1.18-61CFDD.svg)
[![GoDoc](https://pkg.go.dev/badge/badge)](https://pkg.go.dev/github.com/rmordechay/jogson)
[![Go Report Card](https://goreportcard.com/badge/github.com/rmordechay/jogson)](https://goreportcard.com/report/github.com/rmordechay/jogson)

A simple Go library to simplify working with JSON without the need to define structs.

* [Installation](#installation)
* [Create Object, Array or Mapper](#create-jsonobject-jsonarray-or-jsonmapper)
* [Read from JSON](#read-from-json)
    * [Scalars](#scalars)
    * [Objects](#objects)
    * [Arrays](#arrays)
    * [Time](#time)
    * [UUID](#uuid)
    * [Types](#types)
    * [Get JSON String](#get-json-string)
* [Write to JSON](#write-to-JSON)
    * [Write Object](#write-object)
    * [Write Array](#write-array)
* [Error Handling](#error-handling)
* [Design](#design)
  * [JsonMapper](#jsonmapper)
  * [JsonObject](#jsonobject)
  * [JsonArray](#jsonarray)

## Installation
Jogson requires Go version `1.18` or above. 

To install the library use:

```bash
go get github.com/rmordechay/jogson
```

## Create JsonObject, JsonArray or JsonMapper

For more information, see [design](#Design).

#### From String

```go
object, err := jogson.NewObjectFromString(jsonString)
array, err := jogson.NewArrayFromString(jsonString)
mapper, err := jogson.NewMapperFromString(jsonString)
```

#### From Bytes

```go
object, err := jogson.NewObjectFromBytes(jsonBytes)
array, err := jogson.NewArrayFromBytes(jsonBytes)
mapper, err := jogson.NewMapperFromBytes(jsonBytes)
```

#### From Struct

```go
object, err := jogson.NewObjectFromStruct(jsonStruct)
array, err := jogson.NewArrayFromStruct(jsonStruct)
mapper, err := jogson.NewMapperFromStruct(jsonStruct)
```

#### From File

```go
mapper, err := jogson.NewMapperFromFile(jsonFilePath)
object, err := jogson.NewObjectFromFile(jsonFilePath)
array, err := jogson.NewArrayFromFile(jsonFilePath)
```

## Read from JSON

Once you have an object, an array or a mapper, you can read the data easily. Consider the following JSON

```go
jsonString := `{
    "id": "748494b8-7d6e-4cad-8065-89e758797313",
    "name": "Jason",
    "age": 43,
    "height": 1.87,
    "is_funny": false,
    "address": null,
    "birthday": "1981-10-08",
    "features": ["tall", "blue eyes", null],
    "children": {
        "Rachel": {"age": 15, "is_funny": false}, 
        "Sara":   {"age": 19, "is_funny": true}
    }
}`
```

### Scalars

Getting scalars - `string`, `int`, etc. - is similar both for object and array and only differ
in the parameter type (objects take `string` as key and arrays take `int` as index). You can get scalars by value or by reference, where the latter allows JSON null values. Nullable methods have
the suffix 'N' in their names.

#### From Object

By value:

```go
// string 
var name string = object.GetString("name") // Jason

// int 
var age int = object.GetInt("age") // 15

// float64 
var height float64 = object.GetFloat("height") // 1.87

// bool 
var isFunny bool = object.GetBool("is_funny") // false
```

By reference, which allows JSON null values (note the suffix 'N' at the end of the method names):

```go
// string
var nameNullable *string = object.GetStringN("non-existent-key") // nil

// int
var ageNullable *int = object.GetIntN("non-existent-key") // nil

// float64
var heightNullable *float64 = object.GetFloatN("non-existent-key") // nil

// bool
var isFunnyNullable *bool = object.GetBoolN("non-existent-key") // nil
```

#### From Array

By value:

```go
// string 
var s string = array.GetString(0)

// int 
var i int = array.GetInt(2)

// float64 
var f float64 = array.GetFloat(5)

// bool 
var b bool = array.GetBool(7)
```

By reference, which allows JSON null values (Note the suffix 'N' at the end of the method names):

```go
// string
var nameNullable *string = array.GetStringN(100) // nil

// int
var ageNullable *int = array.GetIntN(100) // nil

// float64
var heightNullable *float64 = array.GetFloatN(100) // nil

// bool
var isFunnyNullable *bool = array.GetBoolN(100) // nil
```


### Objects

```go
// Check if a key exists
var keyExists bool = object.Contains("children")

// Check if the object is empty
var objectEmpty bool = object.IsEmpty()

// Get the object's size
var objectLen int = object.Length()

// Get object's keys
var keys []string = object.Keys()

// Get object's values
values := object.Values()

// Iterating over an object with key, value pair
children := object.GetObject("children")
for key, child := range children.Elements() {
    fmt.Println("Child name:", key)                   // Rachel, Sara
    fmt.Println(child.AsObject.GetInt("age"))         // 15, 19
    fmt.Println(child.AsObject.GetBool("is_funny"))   // false, true
}
```

#### As Maps

You can also get the object as a map of strings by scalars (`string`, `int`, etc.) Values that 
are JSON `null` will be returned as Go zero value. If you want to regard null values, call the 
function with the suffix `N` (see next section). If a null was found, `LastError` will be set
and report a null conversion error. For more information, see [error handling](#error-handling).

```go
// Get the JsonObject as map of string and string
var stringMap map[string]string = object.AsStringMap()

// Get the JsonObject as map of string and int
var intMap map[string]int = object.AsIntMap()

// Get the JsonObject as map of string and float
var floatMap map[string]float64 = object.AsFloatMap()

// Get as a slice of JsonArray
var nestedArray map[string]JsonArray = object.AsArrayMap()

// Get as a slice of JsonObject
var objectArray map[string]JsonObject = object.AsObjectMap()
```

#### As Maps with Nullable Values

If the object contains null values, and you want to represent that as nil value instead of zero value, you can 
use one of the following functions that returns a pointer rather than a value which allows nil. This set of methods 
do not set the `LastError` when a null value was found but instead returns it without reporting an error. 

Note! This applies only for the scalar types. `JsonObject` and `JsonArray` will return an instance anyway,
and you should use the `IsNull()` method to check if they are null.

```go
// Get the JsonObject as map of string and nullable strings
var nullableStringMap map[string]*string = object.AsStringMapN()

// Get the JsonObject as map of string and nullable ints
var nullableIntMap map[string]*int = object.AsIntMapN()

// Get the JsonObject as map of string and nullable floats
var nullableSloatMap map[string]*float64 = object.AsFloatMapN()
```

### Arrays

```go
// Get array from object
features := object.GetArray("features")

// Get features length
arrayLen := features.Length() // 2

// Get an element of features by index
secondElement := features.Get(1)

// Iterating over an array
for _, feature := range features.Elements() {
    fmt.Println(feature.AsString) // tall, ...
}
```

#### As Arrays

You can also get the array as a slice of scalars (`string`, `int`, etc.). Values that
are JSON `null` will be returned as Go zero value. If you want to regard null values, call the
function with the suffix `N` (see next section). If a null was found, `LastError` will be set
and report a null conversion error. For more information, see [error handling](#error-handling).

```go

// Get the JsonArray as a slice of string
var stringArray []string = array.AsStringArray()

// Get the JsonArray as a slice of int
var intArray []int = array.AsIntArray()

// Get the JsonArray as a slice of float
var floatArray []float64 = array.AsFloatArray()

// Get the JsonArray as a slice of JsonArray
var nestedArray []JsonArray = array.As2DArray()

// Get the JsonArray as a slice of JsonObject
var objectArray []JsonObject = array.AsObjectArray()
```

#### As Arrays of Nullable Values

If the array contains null values, and you want to represent that as nil value instead of zero value, you can
use one of the following functions that returns a pointer rather than a value which allows nil. This set of methods
do not set the `LastError` when a null value was found but instead returns it without reporting an error.

Note! This applies only for the scalar types. `JsonObject` and `JsonArray` will return an instance anyway, 
and you should use the `IsNull()` method to check if they are null.

```go
// Get the JsonArray as a slice of nullable strings
var nullableStringArray []*string = array.AsStringArrayN()

// Get the JsonArray as a slice of nullable ints
var nullableIntArray []*int = array.AsIntArrayN()

// Get the JsonArray as a slice of nullable floats
var nullableSloatArray []*float64 = array.AsFloatArrayN()
```

### Time

To get a string value as `time.Time`

```go
var birthday time.Time = object.GetTime("birthday") // 1981-10-08T00:00:00Z
var birthday time.Time = array.GetTime(0)
var birthday time.Time = mapper.AsTime()
```

The mapper will try to format the string against different time formats to increase the chance of correct parsing. The following 
formats are supported:

`time.RFC3339` `time.RFC850` `time.RFC822` `time.RFC822Z` `time.RFC1123` `time.RFC1123Z` `time.RFC3339Nano` `time.ANSIC` `time.UnixDate` `time.RubyDate` `time.Layout` `time.Kitchen` `time.Stamp` `time.StampMilli` `time.StampMicro` `time.StampNano` `time.DateTime` `time.DateOnly` `time.TimeOnly`


### UUID

You can also get a string value as `uuid.UUID` type

```go
var uuidValue uuid.UUID = object.GetUUID("id")
var uuidValue uuid.UUID = array.GetUUID(0)
var uuidValue uuid.UUID = mapper.AsUUID()
```

### Types

To check what type is your `JsonMapper` currently holding, use

```go
fmt.Println(mapper.IsObject)    // true
fmt.Println(mapper.IsBool)      // false
fmt.Println(mapper.IsInt)       // false
fmt.Println(mapper.IsFloat)     // false
fmt.Println(mapper.IsString)    // false
fmt.Println(mapper.IsArray)     // false
fmt.Println(mapper.IsNull)      // false
```

### Get JSON String

Every JSON element - `JsonObject`, `JsonArray` or `JsonMapper` - have the method `String()` and `PrettyString()`. 
The returned string from these methods will be a valid JSON. For example

```go
fmt.Println(mapper.AsObject.String())
// output: {"age":43,"children":{"Rachel":{"age":15,"is_funny":false},"Sara":{"age":19,"is_funny":true}},"features":["tall","blue eyes"],"is_funny":false,"name":"Jason"}

fmt.Println(mapper.AsObject.Get("children").String())
// output: {"Rachel":{"age":15,"is_funny":false},"Sara":{"age":19,"is_funny":true}}
```

or with `PrettyString()`

```go
fmt.Println(object.GetObject("children").PrettyString())
// output:
// {
//   "address": null,
//   "age": 43,
//   "birthday": "1981-10-08",
//   "children": {
//     "Rachel": {
//       "age": 15,
//       "is_funny": false
//     },
//     "Sara": {
//       "age": 19,
//       "is_funny": true
//     }
//   },
//   "features": [
//     "tall",
//     "blue eyes",
//     null
//   ],
//   "height": 1.87,
//   "is_funny": false,
//   "name": "Jason"
// }
```

## Write to JSON

To write a JSON object or array is as simple as reading from it.

### Write Object

```go
// Create a new object
obj := jogson.NewObject()
obj.AddKeyValue("name", "Chris")
fmt.Println(obj.String()) // {"name":"Chris"}
```

### Write Array

```go
// Create a new array
arr := jogson.NewArray()
arr.AddElement(15)
arr.AddElement(19)
fmt.Println(arr.String()) // [15,19]
```

## Error Handling
Error handling is designed in such a way to not break the flow and to allow a cleaner code. In most cases
errors are not returned, but instead, they are stored in an exported field, `LastError`, which can be found in 
both `JsonObject` and `JsonArray` and can be used to check for the type of error after an operation is made.

```go
_ = object.GetString("non-existent-key")
if object.LastError != nil {
    fmt.Println(object.LastError) // output: key was not found: 'non-existent-key'
}
```

Note, `LastError` is reset at the beginning of every operation, so if you need a reference to an old
error, you will have to store it in a variable. For example:

```go
_ = object.GetString("address")
fmt.Println(object.LastError) // output: type conversion error: <nil> could not be converted to string
_ = object.GetString("name")
fmt.Println(object.LastError) // output: <nil>
```

## Design

There are 3 structs that are important to know when working with the library
* `JsonMapper` struct for generic JSON.
* `JsonObject` represents JSON object.
* `JsonArray`  represents JSON array.

#### JsonMapper

`JsonMapper` is a struct that holds JSON data and serves as a generic type for all possible JSON types. It has `AsX` and `IsX` 
fields with which you can get the data and check the type, respectively. For example, if your data is a JSON object, 
you can call `JsonMapper.AsObject`, or if it's a string, `JsonMapper.AsString`. This struct is best used when you don't 
know the type at compile time and want to check it dynamically. In this case you can use `JsonMapper.IsArray`, 
`JsonMapper.IsString`, `JsonMapper.IsFloat`, etc. 

Note, in any case, `AsX` never returns nil, but always the zero value. If the underlying data is null, then `IsNull` will 
be set to true. 

`JsonMapper` is also returned in cases where the return type can be any JSON type. For example, `JsonArray.Elements()` 
returns a slice `[]JsonMapper` over which you can iterate or query specific elements. Other methods that return `JsonMapper` 
(instead of `JsonObject` or `JsonArray`) are `JsonObject.Values()`, `JsonObject.Get()`, `JsonArray.Get()`, `JsonArray.Find()` and more. 

#### JsonObject

`JsonObject` holds JSON object data and has different methods to read from JSON object and write to it. Once you have an instance, 
you can use the various methods to read or write data to you object. There are multiple ways to get a `JsonObject`. You can parse 
it directly, get it as an element of an array, or a value in a parent object.

There are 3 sets of methods that you can use when working with `JsonObject`.
* `GetX(key string)`: gets a value in the object as the requested type. For example, GetString("key") gets the value associated with `key` as type `string` 
* `AsXMap()`: gets the value as the requested type. For example, AsStringMap("key") gets the value associated with `key` as type `map[string]string` 
* `AddX(key string, value X)`: adds `value`, associated with `key`, to the object as the requested type. For example, `AddString("key", "string_value")` adds the string `string_value` associated with the `key`   

#### JsonArray
`JsonArray` is in many ways almost identical to `JsonObject`. It also contains the underlying array and methods to read and 
write data, and both also keep the same names for the methods. However, they have some minor differences. You can read more about them 
in the next section.

#### `JsonObject` and `JsonArray` Similarities and Differences
The structs `JsonObject` and `JsonArray` have very similar methods, both in naming and semantics. For example, both structs
have `GetString()`, `GetInt()`, `GetFloat()`, `IsNull()`, `IsEmpty()` and other methods that are identical in names. Additionally, 
they both have similar method names with the similar semantics. For example, `JsonObject.AsStringMap` will return `map[string]string`,
whereas `JsonArray.AsStringArray` will return `[]string`.


As a rule of thumb, they have 2 differences:
  * **Input**:
    * `JsonObject`'s methods mostly take a `string` as the key, e.g. `GetInt("age")`.
    * `JsonArray`'s methods mostly take an `int` as the index, e.g. `GetBool(1)`.
  * **Output**:
    * `JsonObject`'s methods mostly return a `map` or `JsonObject`, e.g. `GetFloatMap() -> map[string]string`. 
    * `JsonArray`'s methods mostly return a `slice` or `JsonArray`, e.g. `GetFloatArray(1)  -> []float`.

#### Methods and Variables Prefix
The prefixes, `As`, `Is`, `Get` and `Add` have similar semantics across the library and can be found in `JsonMapper`, `JsonObject`
and `JsonArray`.
* `IsX`: checks for the value's type. For example `JsonMapper.IsBool`
* `AsX`: converts the current data to other type representation. For example, `JsonArray.AsStringArray()` converts JsonArray to `[]string`.
* `GetX`: Fetches the data, usually with some sort of search in the underlying data.
* `AddX`: Adds the data to the JSON array or object
