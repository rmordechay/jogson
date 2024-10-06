package jsonmapper

import (
	"errors"
	"fmt"
)

type JsonArray struct {
	elements []interface{}
}

func (a JsonArray) Elements() []Mapper {
	jsons := make([]Mapper, 0, len(a.elements))
	for _, element := range a.elements {
		jsons = append(jsons, getMapperFromField(element))
	}
	return jsons
}

func (a JsonArray) Get(i int) Mapper {
	if i >= a.Length() {
		indexErr := fmt.Sprintf("index out of range [%v] with length %v", i, a.Length())
		var mapper Mapper
		mapper.Err = errors.New(indexErr)
		return mapper
	}
	return getMapperFromField((a.elements)[i])
}

func (a JsonArray) Length() int {
	return len(a.elements)
}

func (a JsonArray) AddValue(value interface{}) JsonArray {
	a.elements = append(a.elements, value)
	return a
}

func CreateEmptyJsonArray() JsonArray {
	var arr JsonArray
	elements := make([]interface{}, 0)
	arr.elements = elements
	return arr
}

func CreateJsonArray(data interface{}) JsonArray {
	var arr JsonArray
	arr.elements = data.([]interface{})
	return arr
}

func Map[T JsonType](arr JsonArray, f func(mapper Mapper) T) []T {
	var jsonMappers []T
	for _, element := range arr.elements {
		field := f(getMapperFromField(element))
		jsonMappers = append(jsonMappers, field)
	}
	return jsonMappers
}

func ForEach(arr JsonArray, f func(mapper Mapper)) {
	for _, element := range arr.elements {
		f(getMapperFromField(element))
	}
}

func Filter(arr JsonArray, f func(mapper Mapper) bool) []Mapper {
	var jsonMappers []Mapper
	for _, element := range arr.elements {
		field := getMapperFromField(element)
		if f(field) {
			jsonMappers = append(jsonMappers, field)
		}
	}
	return jsonMappers
}

func (a JsonArray) String() string {
	return string(marshal(a.elements))
}

func parseJsonArray(data string) (JsonArray, error) {
	var ja JsonArray
	var arr []interface{}
	err := unmarshal([]byte(data), &arr)
	if err != nil {
		return JsonArray{}, err
	}
	ja.elements = arr
	return ja, nil
}

func convertArray[T JsonType](data []T) JsonArray {
	var arr JsonArray
	result := make([]interface{}, len(data))
	for i, v := range data {
		result[i] = v
	}
	arr.elements = result
	return arr
}
