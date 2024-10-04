package main

import (
	"fmt"
	"json-mapper/mapper"
)

func main() {
	const jsonStr = `
  {
  	"field1": 0,
  	"field2": 0,
  	"field3": {
			"field1": 5
		}
	}
	`
	jsonMapper := mapper.CreateObjectMapperFromStr(jsonStr)
	int1 := jsonMapper.GetByName("field3").GetByName("field1").AsInt()
	int2 := jsonMapper.GetByName("field1").AsInt()
	fmt.Println(int1 + int2)
}
