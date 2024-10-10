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

// NewArray initializes and returns a new instance of JsonArray.
func NewArray(data []*any) *JsonArray {
	var arr JsonArray
	arr.elements = data
	return &arr
}

// EmptyArray initializes and returns an empty new instance of JsonArray.
func EmptyArray() *JsonArray {
	var arr JsonArray
	elements := make([]*any, 0)
	arr.elements = elements
	return &arr
}

// Length returns the number of elements in the JsonArray.
func (a *JsonArray) Length() int {
	return len(a.elements)
}

// IsEmpty checks if the JSON array has no elements in it
func (a *JsonArray) IsEmpty() bool {
	return len(a.elements) == 0
}

// Elements returns all elements in the JsonArray as a slice of JsonMapper objects.
func (a *JsonArray) Elements() []JsonMapper {
	jsons := make([]JsonMapper, 0, len(a.elements))
	for _, element := range a.elements {
		jsons = append(jsons, getMapperFromField(element))
	}
	return jsons
}

// AsStringArray converts the elements of the JsonArray into a slice of strings.
func (a *JsonArray) AsStringArray() []string {
	return getGenericArray(convertAnyToString, *a)
}

// AsIntArray converts the elements of the JsonArray into a slice of integers.
func (a *JsonArray) AsIntArray() []int {
	return getGenericArray(convertAnyToInt, *a)
}

// AsFloatArray converts the elements of the JsonArray into a slice of floats.
func (a *JsonArray) AsFloatArray() []float64 {
	return getGenericArray(convertAnyToFloat, *a)
}

// As2DArray converts the elements of the JsonArray into a two-dimensional array, returning
// a slice of JsonArray objects.
func (a *JsonArray) As2DArray() []JsonArray {
	array := getGenericArray(convertAnyToArray, *a)
	return array
}

// AsObjectArray converts the elements of the JsonArray into a slice of JsonObject objects.
func (a *JsonArray) AsObjectArray() []JsonObject {
	return getGenericArray(convertAnyToObject, *a)
}

// Get retrieves the value at index i and returns it as a JsonMapper
func (a *JsonArray) Get(i int) JsonMapper {
	return getMapperFromField(a.elements[i])
}

// GetString retrieves the string value from the element at the specified index.
// If the index is out of range, the value is invalid or is null, an error will be set to LastError.
func (a *JsonArray) GetString(i int) string {
	return getArrayScalar(a, convertAnyToString, i, stringTypeStr)
}

// GetInt retrieves the integer value from the element at the specified index.
// If the index is out of range, the value is invalid or is null, an error will be set to LastError.
func (a *JsonArray) GetInt(i int) int {
	return getArrayScalar(a, convertAnyToInt, i, intTypeStr)
}

// GetFloat retrieves the float value from the element at the specified index.
// If the index is out of range, the value is invalid or is null, an error will be set to LastError.
func (a *JsonArray) GetFloat(i int) float64 {
	return getArrayScalar(a, convertAnyToFloat, i, floatTypeStr)
}

// GetBool retrieves the boolean value from the element at the specified index.
// If the index is out of range, the value is invalid or is null, an error will be set to LastError.
func (a *JsonArray) GetBool(i int) bool {
	return getArrayScalar(a, convertAnyToBool, i, boolTypeStr)
}

// GetTime retrieves the time value from the element at the specified index, returning an error if conversion fails.
// If the index is out of range, the value is invalid or is null, an error will be set to LastError.
func (a *JsonArray) GetTime(i int) time.Time {
	return getArrayScalar(a, parseTime, i, timeTypeStr)
}

// GetObject retrieves the JsonObject from the element at the specified index.
// If the index is out of range, the value is invalid or is null, an error will be set to LastError.
func (a *JsonArray) GetObject(i int) *JsonObject {
	if i >= a.Length() {
		a.setLastError(createIndexOutOfRangeErr(i, a.Length()))
		return EmptyObject()
	}
	element := a.elements[i]
	if element == nil {
		a.setLastError(createNullConversionErr(objectTypeStr))
		return EmptyObject()
	}
	switch (*element).(type) {
	case map[string]*any:
		data := (*element).(map[string]*any)
		return NewObject(data)
	case map[string]any:
		data := convertToMapValuesPtr((*element).(map[string]any))
		return NewObject(data)
	default:
		a.setLastError(createTypeConversionErr(*element, objectTypeStr))
		return EmptyObject()
	}
}

// GetArray retrieves the JsonArray from the element at the specified index.
// If the index is out of range, the value is invalid or is null, an error will be set to LastError.
func (a *JsonArray) GetArray(i int) *JsonArray {
	if i >= a.Length() {
		a.setLastError(createIndexOutOfRangeErr(i, a.Length()))
		return EmptyArray()
	}
	element := a.elements[i]
	if element == nil {
		a.setLastError(createNullConversionErr(arrayTypeStr))
		return EmptyArray()
	}
	v, ok := (*element).([]any)
	if !ok {
		a.setLastError(createTypeConversionErr(*element, arrayTypeStr))
		return EmptyArray()
	}
	return NewArray(convertToSlicePtr(v))
}

// AddElement appends a new element to the JsonArray.
func (a *JsonArray) AddElement(value any) {
	switch value.(type) {
	case JsonObject:
		var object any = value.(JsonObject).object
		a.elements = append(a.elements, &object)
	case *JsonObject:
		var object any = value.(JsonObject).object
		a.elements = append(a.elements, &object)
	case JsonArray:
		var valueElements any = value.(JsonArray).elements
		a.elements = append(a.elements, &valueElements)
	case *JsonArray:
		var valueElements any = value.(*JsonArray).elements
		a.elements = append(a.elements, &valueElements)
	case nil, string, int, float64, bool, []string, []int, []float64, []bool:
		a.elements = append(a.elements, &value)
	default:
		a.setLastError(fmt.Errorf("could not add value of type %T", value))
	}
}

// ForEach applies the given function to each element in the JsonArray.
func (a *JsonArray) ForEach(f func(j JsonMapper)) {
	for _, element := range a.elements {
		f(getMapperFromField(element))
	}
}

// Filter returns a new JsonArray containing only the elements that satisfy the given filter function.
func (a *JsonArray) Filter(f func(j JsonMapper) bool) JsonArray {
	var arr = EmptyArray()
	for _, element := range a.elements {
		if f(getMapperFromField(element)) {
			arr.elements = append(arr.elements, element)
		}
	}
	return *arr
}

// FilterNull returns a new JsonArray excluding any elements that are null.
func (a *JsonArray) FilterNull() JsonArray {
	var arr = EmptyArray()
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
func (a *JsonArray) setLastError(err error) {
	a.LastError = err
}

// String returns a string representation of the JsonArray in JSON format.
func (a *JsonArray) String() string {
	jsonBytes, _ := marshal(a.elements)
	return string(jsonBytes)
}
