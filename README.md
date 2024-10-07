# JSON Mapper Library for Go

[![GoDoc](https://pkg.go.dev/badge/badge)](https://pkg.go.dev/github.com/rmordechay/jsonmapper)

A simple Go library to simplify working with JSON without the need to define structs.

## Installation
To install the library, use:

```bash
go get github.com/rmordechay/json-mapper
```

## Creating a mapper
There are multiple ways to create a mapper, depending on your source and if you want to read or write.
For reading JSON, you will need to pass the source. For writing JSON from scratch, you can create an empty mapper.
#### Bytes
```go
mapper, err := jsonmapper.FromBytes([]byte(jsonBytes))
```
#### String
```go
mapper, err := jsonmapper.FromString(jsonString)
```
#### Struct
```go
mapper, err := jsonmapper.FromStruct(jsonStruct)
```
#### File
```go
mapper, err := jsonmapper.FromFile("file.json")
```

## Reading from JSON
Once you have the `mapper`, you can read the data. For example:
```go
jsonString := `{
    "name": "Jason",
    "age": 43,
    "is_funny": false,
    "features": ["tall", "blue eyes"],
    "children": {"John": {"age": 15}, "Sara": {"age": 19}}
}`

fmt.Println(mapper.IsObject) // true    

var name string = mapper.Object.Get("name").AsString
var age int = mapper.Object.Get("age").AsInt
var isFunny bool = mapper.Object.Get("is_funny").AsBool

fmt.Println(name)       // Jason
fmt.Println(age)        // 43
fmt.Println(isFunny)    // false

`
```

If you need to iterate of over a list or an object, you can call the `Elements()` function:
```go
array := object.Get("features").Array
for _, feature := range array.Elements() {
    fmt.Println(feature.AsString) // tall, ...
}
```



## Writing to JSON