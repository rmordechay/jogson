# JSON Mapper Library for Go

[![GoDoc](https://pkg.go.dev/badge/badge)](https://pkg.go.dev/github.com/rmordechay/json-mapper)

A simple Go library to simplify working with JSON without the need to define structs.

## Installation
To install the library, use `go get`:

```bash
go get github.com/rmordechay/json-mapper
```

## Examples

```go
const jsonString = `
    {
      "name": "Jason",
      "age": 43
    }
`

func main() {
    mapper, err := jsonmapper.GetMapperFromString(jsonString)
    if err != nil {
        log.Fatal(err)
    }
    if mapper.IsObject {
        name := mapper.Object.Get("name").AsString
        age := mapper.Object.Get("age").AsInt
        fmt.Println("Name is ", name)  // outputs: Name is Jason
        fmt.Println("Age is ", age)    // outputs: Age is 43
    }
}

```