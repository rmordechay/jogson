package jsonmapper

import (
	"fmt"
)

type JsonArray struct {
	elements []*interface{}

	Err error
}

func NewArray() *JsonArray {
	var arr JsonArray
	elements := make([]*interface{}, 0)
	arr.elements = elements
	return &arr
}

func (a *JsonArray) Elements() []Mapper {
	jsons := make([]Mapper, 0, len(a.elements))
	for _, element := range a.elements {
		jsons = append(jsons, getMapperFromField(element))
	}
	return jsons
}

func (a *JsonArray) Get(i int) Mapper {
	if i >= a.Length() {
		a.Err = fmt.Errorf("index out of range [%v] with length %v", i, a.Length())
		return Mapper{}
	}
	return getMapperFromField((a.elements)[i])
}

func (a *JsonArray) Length() int {
	return len(a.elements)
}

func (a *JsonArray) AddValue(value interface{}) {
	a.elements = append(a.elements, &value)
}

func (a *JsonArray) ForEach(f func(mapper Mapper)) {
	for _, element := range a.elements {
		f(getMapperFromField(element))
	}
}

func (a *JsonArray) Filter(f func(mapper Mapper) bool) JsonArray {
	var arr = NewArray()
	for _, element := range a.elements {
		field := getMapperFromField(element)
		if f(field) {
			arr.elements = append(arr.elements, element)
		}
	}
	return *arr
}

func (a *JsonArray) String() string {
	jsonBytes, _ := marshal(a.elements)
	return string(jsonBytes)
}

func parseJsonArray(data []byte) (JsonArray, error) {
	var ja JsonArray
	var arr []*interface{}
	err := unmarshal(data, &arr)
	if err != nil {
		return JsonArray{}, err
	}
	ja.elements = arr
	return ja, nil
}

func convertArray[T any](data []T) JsonArray {
	var arr JsonArray
	array := make([]*interface{}, len(data))
	for i, v := range data {
		var valAny interface{} = v
		array[i] = &valAny
	}
	arr.elements = array
	return arr
}
