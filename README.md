# JSON Mapper Library for Go

[![GoDoc](https://pkg.go.dev/badge/badge)](https://pkg.go.dev/github.com/rmordechay/jsonmapper)

A simple Go library to simplify working with JSON without the need to define structs.

* [Installation](#Installation)
* [Create a Mapper](#Create-a-mapper)
* [Read from JSON](#Read-from-JSON)
  * [Check Types](#Check-Types)
  * [Objects](#Objects)
  * [Get Values](#Primitive-Types)
  * [Arrays](#Arrays)
  * [Find Elements](#Find-Elements)
  * [Print JSON](#Get-JSON-as-string)
* [Write to JSON](#Write-to-JSON)
  * [Write Object](#Write-object)
  * [Write Array](#Write-array)


## Installation
To install the library, use:

```bash
go get github.com/rmordechay/jsonmapper
```

### Create a Mapper
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
object := mapper.Object

// Check if a key exists
var keyExists bool = object.Has("children")

// Get object's keys
var keys []string = object.Keys()

// Get object's values
values := object.Values()

// Iterating over an object with key, value pair
children := object.GetObject("children")
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

#### Time
```go
var birthday time.Time = array.GetTime(0)
var birthday string = object.GetTime("birthday".Format(time.RFC3339)) // 1981-10-08T00:00:00Z
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