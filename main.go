package main

import (
	"fmt"
)

func main() {
	const jsonStr = `
    {
  		"field1": [1, 2],
  		"field3": {
			"field1": 5.5,
			"field2": [
				{"name": "roi"},
				{"name": "adi"}
			]
		},
  		"field4": "null",
  		"field5": 6,
  		"field6": 7
	}`
	mapper, err := GetMapper(jsonStr)
	if err != nil {
		panic(err)
	}
	t := mapper.AsObject.Get("field3").AsObject.Get("field2").AsArray
	//get := t.Get("field3").AsObject
	fmt.Println(mapper)
	for _, v := range t.Elements() {
		object := v.AsObject
		fmt.Println(v)
		for k1, v1 := range object.Elements() {
			fmt.Println(k1)
			fmt.Println(v1)
		}
	}
}
