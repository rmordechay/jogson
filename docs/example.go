package docs

import (
	"fmt"
	"github.com/rmordechay/jsonmapper"
	"log"
)

func main() {
	jsonString := `
	{
	  "name": "Jason",
	  "age": 43
	}
	`
	mapper, err := jsonmapper.CreateMapperFromString(jsonString)
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
