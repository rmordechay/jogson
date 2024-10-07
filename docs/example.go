package docs

import (
	"fmt"
	"github.com/rmordechay/jsonmapper"
)

func RunExample() {
	jsonString := `
	{
		"name": "Jason",
		"age": 43,
		"is_funny": false,
		"features": ["tall", "blue eyes"],
		"birthday": "1981-10-08",
		"children": {
			"Rachel": {"age": 15, "is_funny": false}, 
			"Sara":   {"age": 19, "is_funny": true}
		}
	}
	`
	mapper, _ := jsonmapper.FromString(jsonString)
	//var birthday time.Time = mapper.Object.Get("birthday").AsTime()
	fmt.Println(mapper.Object.Get("children").PrettyString())

	//var name string = object.Get("name").AsString
	//var age int = object.Get("age").AsInt
	//var isFunny bool = object.Get("is_funny").AsBool
	//
	//fmt.Println(name)    // Jason
	//fmt.Println(age)     // 43
	//fmt.Println(isFunny) // false
	//
	//array := object.Get("features").Array
	//for _, feature := range array.Elements() {
	//	fmt.Println(feature.AsString) // tall, ...
	//}
	//
	//children := object.Get("children").Object
	//for key, child := range children.Elements() {
	//	fmt.Println("Child name:", key)
	//	fmt.Println(child.Object.Get("age").AsInt)       // 15, 19
	//	fmt.Println(child.Object.Get("is_funny").AsBool) // false, true
	//}
}
