package docs

import (
	"fmt"
	jsonmapper "github.com/rmordechay/json-mapper"
	"log"
)

func main() {
	jsonString := `
	{
	  "name": "Jason",
	  "age": 43
	}
	`
	mapper, err := jsonmapper.GetMapperFromString(jsonString)
	if err != nil {
		log.Fatal(err)
	}
	if mapper.IsObject {
		name := mapper.Object.Get("name").AsString
		age := mapper.Object.Get("age").AsInt
		fmt.Println("Name is: ", name)
		fmt.Println("Name is: ", age)
	}
}
