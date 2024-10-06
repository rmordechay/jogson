package jsonmapper

import (
	"fmt"
	"testing"
)

func TestMain2(t *testing.T) {
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
	mapper, _ := GetMapper(jsonStr)
	find := mapper.AsObject.Find("field2")
	fmt.Println(find)
}
