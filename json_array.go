package jogson

import (
	"os"
	"time"

	"github.com/google/uuid"
)

// JsonArray represents a JSON array
type JsonArray struct {
	elements  []*any
	LastError error
}

// NewArrayFromBytes parses JSON data from a byte slice.
func NewArrayFromBytes(data []byte) (*JsonArray, error) {
	jsonArray := EmptyArray()
	err := unmarshal(data, &jsonArray.elements)
	if err != nil {
		return &JsonArray{}, err
	}
	return jsonArray, nil
}

// NewArrayFromFile reads a JSON file from the given path and parses it into a JsonArray object.
func NewArrayFromFile(path string) (*JsonArray, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return &JsonArray{}, err
	}
	return NewArrayFromBytes(file)
}

// NewArrayFromString parses JSON from a string into a JsonArray object.
func NewArrayFromString(data string) (*JsonArray, error) {
	return NewArrayFromBytes([]byte(data))
}

// EmptyArray initializes and returns an empty new instance of JsonArray.
func EmptyArray() *JsonArray {
	var arr JsonArray
	arr.elements = make([]*any, 0)
	return &arr
}

// nullArray initializes and returns a null JsonArray.
func nullArray() *JsonArray {
	var arr JsonArray
	return &arr
}

// Length returns the number of elements in the JsonArray.
func (a *JsonArray) Length() int {
	return len(a.elements)
}

// IsEmpty checks if the JSON array has no elements
func (a *JsonArray) IsEmpty() bool {
	return len(a.elements) == 0
}

// IsNull checks if the JSON array is null
func (a *JsonArray) IsNull() bool {
	return a.elements == nil
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

// AsStringArrayN is the nullable version of AsStringArray, and returns a slice of string pointers instead values which
// imitates JSON's null as Go's nil.
func (a *JsonArray) AsStringArrayN() []*string {
	return getGenericArrayN(convertAnyToStringN, *a)
}

// AsIntArrayN is the nullable version of AsIntArray, and returns a slice of int pointers instead values which
// imitates JSON's null as Go's nil.
func (a *JsonArray) AsIntArrayN() []*int {
	return getGenericArrayN(convertAnyToIntN, *a)
}

// AsFloatArrayN is the nullable version of AsFloatArray, and returns a slice of float64 pointers instead values which
// imitates JSON's null as Go's nil.
func (a *JsonArray) AsFloatArrayN() []*float64 {
	return getGenericArrayN(convertAnyToFloatN, *a)
}

// As2DArray converts the elements of the JsonArray into a two-dimensional array, returning
// a slice of JsonArray objects.
func (a *JsonArray) As2DArray() []JsonArray {
	return getGenericArray(convertAnyToArray, *a)
}

// AsObjectArray converts the elements of the JsonArray into a slice of JsonObject objects.
func (a *JsonArray) AsObjectArray() []JsonObject {
	return getGenericArray(convertAnyToObject, *a)
}

// ContainsString checks if the JSON array contains the string s
func (a *JsonArray) ContainsString(s string) bool {
	for _, element := range a.elements {
		if element == nil {
			continue
		}
		if *element == s {
			return true
		}
	}
	return false
}

// ContainsInt checks if the JSON array contains the int i
func (a *JsonArray) ContainsInt(i int) bool {
	for _, element := range a.elements {
		if element == nil {
			continue
		}
		f, ok := (*element).(float64)
		if !ok {
			return false
		}
		if int(f) == i {
			return true
		}
	}
	return false
}

// ContainsFloat checks if the JSON array contains the float f
func (a *JsonArray) ContainsFloat(f float64) bool {
	for _, element := range a.elements {
		if element == nil {
			continue
		}
		if *element == f {
			return true
		}
	}
	return false
}

// Get retrieves the value at index i and returns it as a JsonMapper
// If the index is out of range, the value is invalid or is null, an error will be set to LastError.
func (a *JsonArray) Get(i int) JsonMapper {
	a.setLastError(nil)
	if i >= a.Length() {
		a.setLastError(createIndexOutOfRangeErr(i, a.Length()))
		return JsonMapper{}
	}
	return getMapperFromField(a.elements[i])
}

// GetString retrieves the string value from the element at the specified index. JSON values that are not
// proper string values, i.e. numbers or booleans, will still be converted to string. For example, the value 3
// will be converted to "3".
// If the index is out of range, the value is invalid or is null, an error will be set to LastError.
// In case of an error, the zero value will be returned. If you want to regard null values as well,
// use GetStringN()
func (a *JsonArray) GetString(i int) string {
	return getArrayScalar(a, convertAnyToString, i)
}

// GetInt retrieves the integer value from the element at the specified index.
// If the index is out of range, the value is invalid or is null, an error will be set to LastError.
// In case of an error, the zero value will be returned. If you want to regard null values as well,
// use GetIntN()
func (a *JsonArray) GetInt(i int) int {
	return getArrayScalar(a, convertAnyToInt, i)
}

// GetFloat retrieves the float value from the element at the specified index.
// If the index is out of range, the value is invalid or is null, an error will be set to LastError.
// In case of an error, the zero value will be returned. If you want to regard null values as well,
// use GetFloatN()
func (a *JsonArray) GetFloat(i int) float64 {
	return getArrayScalar(a, convertAnyToFloat, i)
}

// GetBool retrieves the boolean value from the element at the specified index.
// If the index is out of range, the value is invalid or is null, an error will be set to LastError.
// In case of an error, the zero value will be returned. If you want to regard null values as well,
// use GetBoolN()
func (a *JsonArray) GetBool(i int) bool {
	return getArrayScalar(a, convertAnyToBool, i)
}

// GetStringN is the nullable version of GetString, and returns a pointer instead of a zero value which
// imitates JSON's null as Go's nil. A nil pointer is returned if the index is out of bounds,
// it is null or could not be converted to string. The type of the error will be stored in LastError.
func (a *JsonArray) GetStringN(i int) *string {
	return getArrayScalarN(a, convertAnyToStringN, i)
}

// GetIntN is the nullable version of GetInt, and returns a pointer instead of a zero value which
// imitates JSON's null as Go's nil. A nil pointer is returned if the index is out of bounds found,
// it is null or could not be converted to int. The type of the error will be stored in LastError.
func (a *JsonArray) GetIntN(i int) *int {
	return getArrayScalarN(a, convertAnyToIntN, i)
}

// GetFloatN is the nullable version of GetFloat, and returns a pointer instead of a zero value which
// imitates JSON's null as Go's nil. A nil pointer is returned if the index is out of bounds
// it is null or could not be converted to float64. The type of the error will be stored in LastError.
func (a *JsonArray) GetFloatN(i int) *float64 {
	return getArrayScalarN(a, convertAnyToFloatN, i)
}

// GetBoolN is the nullable version of GetBool, and returns a pointer instead of a zero value which
// imitates JSON's null as Go's nil. A nil pointer is returned index is out of bounds found,
// it is null or could not be converted to bool. The type of the error will be stored in LastError.
func (a *JsonArray) GetBoolN(i int) *bool {
	return getArrayScalarN(a, convertAnyToBoolN, i)
}

// GetTime retrieves the value as time.Time at the specified index.
// If the index is out of range, the value is invalid, null or not a string, an error will be set to LastError.
// In case of an error, the zero value will be returned.
func (a *JsonArray) GetTime(i int) time.Time {
	return getArrayScalar(a, parseTime, i)
}

// GetUUID retrieves the value as uuid.UUID at the specified index.
// If the key does not exist, the value is invalid, null or not a string, an error will be set to LastError.
// In case of an error, the zero value will be returned.
func (a *JsonArray) GetUUID(i int) uuid.UUID {
	return getArrayScalar(a, parseUUID, i)
}

// GetObject retrieves the JsonObject from the element at the specified index.
// If the index is out of range, the value is invalid or is null, an error will be set to LastError.
func (a *JsonArray) GetObject(i int) *JsonObject {
	a.setLastError(nil)
	if i >= a.Length() {
		a.setLastError(createIndexOutOfRangeErr(i, a.Length()))
		return nullObject()
	}
	element := a.elements[i]
	if element == nil {
		a.setLastError(createTypeConversionErr(nil, ""))
		return nullObject()
	}
	switch (*element).(type) {
	case map[string]*any:
		data := (*element).(map[string]*any)
		return newObjectFromMap(data)
	case map[string]any:
		data := convertToMapValuesPtr((*element).(map[string]any))
		return newObjectFromMap(data)
	default:
		a.setLastError(createTypeConversionErr(*element, JsonObject{}))
		return nullObject()
	}
}

// GetArray retrieves the JsonArray from the element at the specified index.
// If the index is out of range, the value is invalid or is null, an error will be set to LastError.
func (a *JsonArray) GetArray(i int) *JsonArray {
	a.setLastError(nil)
	if i >= a.Length() {
		a.setLastError(createIndexOutOfRangeErr(i, a.Length()))
		return EmptyArray()
	}
	element := a.elements[i]
	if element == nil {
		a.setLastError(createTypeConversionErr(nil, JsonArray{}))
		return EmptyArray()
	}
	switch v := (*element).(type) {
	case []*any:
		return newArrayFromSlice(v)
	case []any:
		return newArrayFromSlice(convertToSlicePtr(v))
	default:
		a.setLastError(createTypeConversionErr(*element, JsonArray{}))
		return EmptyArray()
	}
}

// AddJsonObject appends a JsonObject to the JsonArray.
func (a *JsonArray) AddJsonObject(jsonObject *JsonObject) {
	var object any = jsonObject.object
	a.elements = append(a.elements, &object)
}

// AddJsonArray appends a nested JsonArray to the JsonArray.
func (a *JsonArray) AddJsonArray(jsonArray *JsonArray) {
	var elements any = jsonArray.elements
	a.elements = append(a.elements, &elements)
}

// AddString appends the string s to the JsonArray.
func (a *JsonArray) AddString(s string) {
	var value any = s
	a.elements = append(a.elements, &value)
}

// AddInt appends the int i to the JsonArray.
func (a *JsonArray) AddInt(i int) {
	var value any = i
	a.elements = append(a.elements, &value)
}

// AddFloat appends the float f to the JsonArray.
func (a *JsonArray) AddFloat(f float64) {
	var value any = f
	a.elements = append(a.elements, &value)
}

// AddBool appends the bool b to the JsonArray.
func (a *JsonArray) AddBool(b bool) {
	var value any = b
	a.elements = append(a.elements, &value)
}

// AddStringArray appends the []string s to the JsonArray.
func (a *JsonArray) AddStringArray(s []string) {
	var value any = s
	a.elements = append(a.elements, &value)
}

// AddIntArray appends the []int i to the JsonArray.
func (a *JsonArray) AddIntArray(i []int) {
	var value any = i
	a.elements = append(a.elements, &value)
}

// AddFloatArray appends the []float f to the JsonArray.
func (a *JsonArray) AddFloatArray(f []float64) {
	var value any = f
	a.elements = append(a.elements, &value)
}

// AddNull appends null to the JsonArray.
func (a *JsonArray) AddNull() {
	a.elements = append(a.elements, nil)
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

// Map returns a slice of the mapped the values of the array
func Map[T any](jsonArray *JsonArray, f func(j JsonMapper) T) []T {
	arr := make([]T, 0, len(jsonArray.elements))
	for _, element := range jsonArray.elements {
		mapper := f(getMapperFromField(element))
		arr = append(arr, mapper)
	}
	return arr
}

// MapNotNull returns a slice of the mapped the values of the array without null values
func MapNotNull[T any](jsonArray *JsonArray, f func(j JsonMapper) T) []T {
	arr := make([]T, 0, len(jsonArray.elements))
	for _, element := range jsonArray.elements {
		if element == nil {
			continue
		}
		mapper := f(getMapperFromField(element))
		arr = append(arr, mapper)
	}
	return arr
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

// All returns true if all elements in the JsonArray are not null. If the array is empty,
// it returns true.
func (a *JsonArray) All() bool {
	for _, element := range a.elements {
		field := getMapperFromField(element)
		if field.IsNull {
			return false
		}
	}
	return true
}

// Any returns true if any element in the JsonArray is non-null. If the array is empty,
// it returns true.
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

// PrettyString returns a pretty-printed string representation of the JsonArray.
func (a *JsonArray) PrettyString() string {
	jsonBytes, _ := marshalIndent(a.elements)
	return string(jsonBytes)
}

// String returns a string representation of the JsonArray in JSON format.
func (a *JsonArray) String() string {
	jsonBytes, _ := marshal(a.elements)
	return string(jsonBytes)
}

// SetLastError sets the last error encountered in the JsonArray.
func (a *JsonArray) setLastError(err error) {
	a.LastError = err
}

// newArrayFromSlice initializes and returns a new instance of JsonArray.
func newArrayFromSlice(data []*any) *JsonArray {
	var arr JsonArray
	arr.elements = data
	return &arr
}
