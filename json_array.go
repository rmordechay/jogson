package jsonmapper

import (
	"fmt"
	"time"
)

// JsonArray represents a JSON array
type JsonArray struct {
	elements  []*any
	LastError error
}

// NewArray initializes and returns a new instance of JsonArray with an empty list of elements.
func NewArray() *JsonArray {
	var arr JsonArray
	elements := make([]*any, 0)
	arr.elements = elements
	return &arr
}

// Length returns the number of elements in the JsonArray.
func (a *JsonArray) Length() int {
	return len(a.elements)
}

// Elements returns all elements in the JsonArray as a slice of Json objects.
func (a *JsonArray) Elements() []Json {
	jsons := make([]Json, 0, len(a.elements))
	for _, element := range a.elements {
		jsons = append(jsons, getMapperFromField(element))
	}
	return jsons
}

// As2DArray converts the elements of the JsonArray into a two-dimensional array, returning
// a slice of JsonArray objects.
func (a *JsonArray) As2DArray() []JsonArray {
	arr := make([]JsonArray, 0, len(a.elements))
	a.LastError = nil
	for _, element := range a.elements {
		asJsonObject := getAsJsonArray((*element).([]any))
		if a.LastError != nil {
			break
		}
		arr = append(arr, asJsonObject)
	}
	return arr
}

// AsObjectArray converts the elements of the JsonArray into a slice of JsonObject objects.
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

// AsStringArray converts the elements of the JsonArray into a slice of strings.
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

// AsIntArray converts the elements of the JsonArray into a slice of integers.
func (a *JsonArray) AsIntArray() []int {
	arr := make([]int, 0, len(a.elements))
	for _, element := range a.elements {
		asInt := getAsInt(element, a)
		arr = append(arr, asInt)
	}
	return arr
}

// AsFloatArray converts the elements of the JsonArray into a slice of floats.
func (a *JsonArray) AsFloatArray() []float64 {
	arr := make([]float64, 0, len(a.elements))
	for _, element := range a.elements {
		asFloat := getAsFloat(element, a)
		arr = append(arr, asFloat)
	}
	return arr
}

// GetString retrieves the string value from the element at the specified index.
func (a *JsonArray) GetString(i int) string {
	if i >= a.Length() {
		a.SetLastError(fmt.Errorf(indexOutOfRangeErrStr, i, a.Length()))
		return ""
	}
	return getAsString(a.elements[i], a)
}

// GetInt retrieves the integer value from the element at the specified index.
func (a *JsonArray) GetInt(i int) int {
	if i >= a.Length() {
		a.SetLastError(fmt.Errorf(indexOutOfRangeErrStr, i, a.Length()))
		return 0
	}
	return getAsInt(a.elements[i], a)
}

// GetFloat retrieves the float value from the element at the specified index.
func (a *JsonArray) GetFloat(i int) float64 {
	if i >= a.Length() {
		a.SetLastError(fmt.Errorf(indexOutOfRangeErrStr, i, a.Length()))
		return 0
	}
	return getAsFloat(a.elements[i], a)
}

// GetBool retrieves the boolean value from the element at the specified index.
func (a *JsonArray) GetBool(i int) bool {
	if i >= a.Length() {
		a.SetLastError(fmt.Errorf(indexOutOfRangeErrStr, i, a.Length()))
		return false
	}
	return getAsBool(a.elements[i], a)
}

// GetTime retrieves the time value from the element at the specified index, returning an error if conversion fails.
func (a *JsonArray) GetTime(key int) (time.Time, error) {
	for k, v := range a.elements {
		if k != key {
			continue
		}
		if v == nil {
			return time.Time{}, fmt.Errorf(nullConversionErrStr, "time.Time")
		}
		return parseTime(v)
	}
	return time.Time{}, fmt.Errorf(keyNotFoundErrStr, key)
}

// GetObject retrieves the JsonObject from the element at the specified index.
func (a *JsonArray) GetObject(i int) *JsonObject {
	if i >= a.Length() {
		a.SetLastError(fmt.Errorf(indexOutOfRangeErrStr, i, a.Length()))
		return &JsonObject{}
	}
	element := a.elements[i]
	if element == nil {
		a.SetLastError(fmt.Errorf(nullConversionErrStr, "JsonObject"))
		return &JsonObject{}
	}
	switch (*element).(type) {
	case map[string]*any:
		return &JsonObject{object: (*element).(map[string]*any)}
	case map[string]any:
		return &JsonObject{object: convertToMapValuesPtr((*element).(map[string]any))}
	default:
		a.SetLastError(fmt.Errorf(typeConversionErrStr, *element, "JsonObject"))
		return &JsonObject{}
	}
}

// GetArray retrieves the JsonArray from the element at the specified index.
func (a *JsonArray) GetArray(i int) JsonArray {
	if i >= a.Length() {
		a.SetLastError(fmt.Errorf(indexOutOfRangeErrStr, i, a.Length()))
		return JsonArray{}
	}
	element := a.elements[i]
	if element == nil {
		a.SetLastError(fmt.Errorf(nullConversionErrStr, "JsonArray"))
		return JsonArray{}
	}
	v, ok := (*element).([]any)
	if !ok {
		a.SetLastError(fmt.Errorf(typeConversionErrStr, *element, "JsonArray"))
		return JsonArray{}
	}
	return JsonArray{elements: convertToSlicePtr(v)}
}

// AddElement appends a new element to the JsonArray.
func (a *JsonArray) AddElement(value any) {
	a.elements = append(a.elements, &value)
}

// ForEach applies the given function to each element in the JsonArray.
func (a *JsonArray) ForEach(f func(j Json)) {
	for _, element := range a.elements {
		f(getMapperFromField(element))
	}
}

// Filter returns a new JsonArray containing only the elements that satisfy the given filter function.
func (a *JsonArray) Filter(f func(j Json) bool) JsonArray {
	var arr = NewArray()
	for _, element := range a.elements {
		if f(getMapperFromField(element)) {
			arr.elements = append(arr.elements, element)
		}
	}
	return *arr
}

// FilterNull returns a new JsonArray excluding any elements that are null.
func (a *JsonArray) FilterNull() JsonArray {
	var arr = NewArray()
	for _, element := range a.elements {
		field := getMapperFromField(element)
		if !field.IsNull {
			arr.elements = append(arr.elements, element)
		}
	}
	return *arr
}

// All returns true if all elements in the JsonArray are non-null.
func (a *JsonArray) All() bool {
	for _, element := range a.elements {
		field := getMapperFromField(element)
		if field.IsNull {
			return false
		}
	}
	return true
}

// Any returns true if any element in the JsonArray is non-null.
func (a *JsonArray) Any() bool {
	if len(a.elements) == 0 {
		return true
	}
	for _, element := range a.elements {
		field := getMapperFromField(element)
		if !field.IsNull {
			return true
		}
	}
	return false
}

// SetLastError sets the last error encountered in the JsonArray.
func (a *JsonArray) SetLastError(err error) {
	a.LastError = err
}

// String returns a string representation of the JsonArray in JSON format.
func (a *JsonArray) String() string {
	jsonBytes, _ := marshal(a.elements)
	return string(jsonBytes)
}

func getAsJsonArray[T any](data []T) JsonArray {
	var arr JsonArray
	array := make([]*any, len(data))
	for i, v := range data {
		var valAny any = v
		array[i] = &valAny
	}
	arr.elements = array
	return arr
}
