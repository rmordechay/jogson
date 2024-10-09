package jsonmapper

import (
	"fmt"
)

type JsonArray struct {
	elements []*interface{}

	LastError error
}

func (a *JsonArray) SetLastError(err error) {
	a.LastError = err
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

func (a *JsonArray) AsNestedArray() []JsonArray {
	arr := make([]JsonArray, 0, len(a.elements))
	a.LastError = nil
	for _, element := range a.elements {
		asJsonObject := getAsJsonArray((*element).([]interface{}))
		if a.LastError != nil {
			break
		}
		arr = append(arr, asJsonObject)
	}
	return arr
}

func (a *JsonArray) AsObjectArray() []JsonObject {
	arr := make([]JsonObject, 0, len(a.elements))
	a.LastError = nil
	for _, element := range a.elements {
		asJsonObject := getAsJsonObject(element, a)
		if a.LastError != nil {
			break
		}
		arr = append(arr, asJsonObject)
	}
	return arr
}

func (a *JsonArray) AsStringArray() []string {
	arr := make([]string, 0, len(a.elements))
	a.LastError = nil
	for _, element := range a.elements {
		asString := getAsString(element, a)
		if a.LastError != nil {
			break
		}
		arr = append(arr, asString)
	}
	return arr
}

func (a *JsonArray) AsIntArray() []int {
	arr := make([]int, 0, len(a.elements))
	for _, element := range a.elements {
		asInt := getAsInt(element, a)
		arr = append(arr, asInt)
	}
	return arr
}

func (a *JsonArray) AsFloatArray() []float64 {
	arr := make([]float64, 0, len(a.elements))
	for _, element := range a.elements {
		asFloat := getAsFloat(element, a)
		arr = append(arr, asFloat)
	}
	return arr
}

func (a *JsonArray) GetString(i int) string {
	if i >= a.Length() {
		a.SetLastError(fmt.Errorf(IndexOutOfRange, i, a.Length()))
		return ""
	}
	return getAsString(a.elements[i], a)
}

func (a *JsonArray) GetInt(i int) int {
	if i >= a.Length() {
		a.SetLastError(fmt.Errorf(IndexOutOfRange, i, a.Length()))
		return 0
	}
	return getAsInt(a.elements[i], a)
}

func (a *JsonArray) GetFloat(i int) float64 {
	if i >= a.Length() {
		a.SetLastError(fmt.Errorf(IndexOutOfRange, i, a.Length()))
		return 0
	}
	return getAsFloat(a.elements[i], a)
}

func (a *JsonArray) GetBool(i int) bool {
	if i >= a.Length() {
		a.SetLastError(fmt.Errorf(IndexOutOfRange, i, a.Length()))
		return false
	}
	return getAsBool(a.elements[i], a)
}

func (a *JsonArray) GetObject(i int) JsonObject {
	if i >= a.Length() {
		a.SetLastError(fmt.Errorf(IndexOutOfRange, i, a.Length()))
		return JsonObject{}
	}
	element := *a.elements[i]
	v, ok := element.(map[string]interface{})
	if !ok {
		a.SetLastError(fmt.Errorf(TypeConversionErrStr, element, JsonObject{}))
		return JsonObject{}
	}
	return JsonObject{object: convertToMapValuesPtr(v)}
}

func (a *JsonArray) GetArray(i int) JsonArray {
	if i >= a.Length() {
		a.SetLastError(fmt.Errorf(IndexOutOfRange, i, a.Length()))
		return JsonArray{}
	}
	element := *a.elements[i]
	v, ok := element.([]interface{})
	if !ok {
		a.SetLastError(fmt.Errorf(TypeConversionErrStr, element, JsonArray{}))
		return JsonArray{}
	}
	return JsonArray{elements: convertToSlicePtr(v)}
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
