# JSON Mapper Library for Go

[![GoDoc](https://pkg.go.dev/badge/badge)](https://pkg.go.dev/github.com/rmordechay/json-mapper)

A simple Go library to simplify working with JSON without the need to define structs.

## Installation
To install the library, use:

```bash
go get github.com/rmordechay/json-mapper
```

## Creating a mapper
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
Once you have the `mapper`, you can read the data easily. Consider the following JSON:
```go
jsonString := `{
    "name": "Jason",
    "age": 43,
    "is_funny": false,
    "birthday": "1981-10-08",
    "features": ["tall", "blue eyes"],
    "children": {
        "Rachel": {"age": 15, "is_funny": false}, 
        "Sara":   {"age": 19, "is_funny": true}
    }
}`
mapper, err := jsonmapper.FromString(jsonString)
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
keyExists := mapper.Object.Has("children") // true
// Get a value by key
element := mapper.Object.Get("age")
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

### Get Values
```go
var name string = mapper.Object.Get("name").AsString
var age int = mapper.Object.Get("age").AsInt
var isFunny bool = mapper.Object.Get("is_funny").AsBool
var birthday time.Time = mapper.Object.Get("birthday").AsTime()

fmt.Println(name)        // Jason
fmt.Println(age)         // 43
fmt.Println(isFunny)     // false
fmt.Println(birthday)    // 1981-10-08 00:00:00 +0000 UTC
```

### Arrays
```go
// Iterating over an array
array := mapper.Object.Get("features").Array
for _, feature := range array.Elements() {
    fmt.Println(feature.AsString) // tall, ...
}

// Get array length
arrayLen := array.Length()
// Get an element of array by index
element := array.Get(2)
```

### Find Elements
You can search for a nested element. 
```go
element := mapper.Object.Find("Rachel")
fmt.Println(child.IsObject)         // true 
fmt.Println(child.Has("is_funny"))  // true 
```

### Print JSON
You can get a string from every JSON element which is a valid JSON
```go
fmt.Println(mapper.Object.String())
// output: {"age":43,"children":{"Rachel":{"age":15,"is_funny":false},"Sara":{"age":19,"is_funny":true}},"features":["tall","blue eyes"],"is_funny":false,"name":"Jason"}
fmt.Println(mapper.Object.Get("children").String())
// output: {"Rachel":{"age":15,"is_funny":false},"Sara":{"age":19,"is_funny":true}}
```
or a pretty string 
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
