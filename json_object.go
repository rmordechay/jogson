package jogson

import (
	"os"
	"time"

	"github.com/google/uuid"
)

// JsonObject represents a JSON object
type JsonObject struct {
	object    map[string]*any
	LastError error
}

// NewObjectFromBytes parses JSON data from a byte slice.
func NewObjectFromBytes(data []byte) (*JsonObject, error) {
	jsonObject := EmptyObject()
	err := unmarshal(data, &jsonObject.object)
	if err != nil {
		return &JsonObject{}, err
	}
	return jsonObject, nil
}

// NewObjectFromFile reads a JSON file from the given path and parses it into a JsonObject object.
func NewObjectFromFile(path string) (*JsonObject, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return &JsonObject{}, err
	}
	return NewObjectFromBytes(file)
}

// NewObjectFromStruct serializes a Go struct into JsonObject.
func NewObjectFromStruct[T any](s T) (*JsonObject, error) {
	jsonBytes, err := marshal(s)
	if err != nil {
		return &JsonObject{}, err
	}
	return NewObjectFromBytes(jsonBytes)
}

// NewObjectFromString parses JSON from a string into a JsonObject object.
func NewObjectFromString(data string) (*JsonObject, error) {
	return NewObjectFromBytes([]byte(data))
}

// EmptyObject initializes and returns an empty new instance of JsonObject.
func EmptyObject() *JsonObject {
	var obj JsonObject
	obj.object = make(map[string]*any)
	return &obj
}

// nullObject initializes and returns a null JsonObject.
func nullObject() *JsonObject {
	var obj JsonObject
	return &obj
}

// Length returns the number of elements in the JsonObject.
func (o *JsonObject) Length() int {
	return len(o.object)
}

// Contains checks if the specified key exists in the JsonObject.
func (o *JsonObject) Contains(key string) bool {
	_, ok := o.object[key]
	return ok
}

// IsEmpty checks if the JSON object has no fields in it
func (o *JsonObject) IsEmpty() bool {
	return len(o.object) == 0
}

// IsNull checks if the JSON object is null
func (o *JsonObject) IsNull() bool {
	return o.object == nil
}

// Keys returns a slice of all keys in the JsonObject.
func (o *JsonObject) Keys() []string {
	keys := make([]string, 0, len(o.object))
	for key := range o.object {
		keys = append(keys, key)
	}
	return keys
}

// Values returns a slice of all values in the JsonObject as JsonMapper types.
func (o *JsonObject) Values() []JsonMapper {
	values := make([]JsonMapper, 0, len(o.object))
	for _, v := range o.object {
		values = append(values, getMapperFromField(v))
	}
	return values
}

// Elements returns a map of all elements in the JsonObject with their keys.
func (o *JsonObject) Elements() map[string]JsonMapper {
	jsons := make(map[string]JsonMapper)
	for k, v := range o.object {
		jsons[k] = getMapperFromField(v)
	}
	return jsons
}

// AsStringMap returns the object as map[string]string
func (o *JsonObject) AsStringMap() map[string]string {
	return getGenericMap(convertAnyToString, *o)
}

// AsIntMap returns the object as map[string]int
func (o *JsonObject) AsIntMap() map[string]int {
	return getGenericMap(convertAnyToInt, *o)
}

// AsFloatMap returns the object as map[string]float64
func (o *JsonObject) AsFloatMap() map[string]float64 {
	return getGenericMap(convertAnyToFloat, *o)
}

// AsStringMapN is the nullable version of AsStringMap, and returns a map of string and string pointers instead values which
// imitates JSON's null as Go's nil.
func (o *JsonObject) AsStringMapN() map[string]*string {
	return getGenericMapN(convertAnyToStringN, *o)
}

// AsIntMapN is the nullable version of AsIntMap, and returns a map of string and pointers int instead values which
// imitates JSON's null as Go's nil.
func (o *JsonObject) AsIntMapN() map[string]*int {
	return getGenericMapN(convertAnyToIntN, *o)
}

// AsFloatMapN is the nullable version of AsFloatMap, and returns a map of string and float64 pointers instead values which
// imitates JSON's null as Go's nil.
func (o *JsonObject) AsFloatMapN() map[string]*float64 {
	return getGenericMapN(convertAnyToFloatN, *o)
}

// AsArrayMap returns the object as map of (string, JsonArray)
func (o *JsonObject) AsArrayMap() map[string]*JsonArray {
	return getGenericMap(convertAnyToArray, *o)
}

// AsObjectMap returns the object as map of (string, JsonObject)
func (o *JsonObject) AsObjectMap() map[string]JsonObject {
	return getGenericMap(convertAnyToObject, *o)
}

// Get retrieves the value associated with the key and returns it as a JsonMapper
// If the key does not exist, the value is invalid or is null, an error will be set to LastError.
// In case of an error, the zero value will be returned.
func (o *JsonObject) Get(key string) *JsonMapper {
	o.setLastError(nil)
	v, ok := o.object[key]
	if !ok {
		o.setLastError(createKeyNotFoundErr(key))
		return &JsonMapper{}
	}
	mapper := getMapperFromField(v)
	return &mapper
}

// ToStruct converts the JsonObject to the type v.
func (o *JsonObject) ToStruct(v any) {
	o.setLastError(nil)
	bytes, err := marshal(o.object)
	if err != nil {
		o.setLastError(err)
		return
	}
	err = unmarshal(bytes, &v)
	if err != nil {
		o.setLastError(err)
		return
	}
}

// GetString retrieves the value associated with the specified key as string. JSON values that are not
// proper string values, i.e. numbers or booleans, will still be converted to string. For example, the value 3
// will be converted to "3".
// If the key does not exist, the value is invalid or is null, an error will be set to LastError.
// In case of an error, the zero value will be returned. If you want to regard null values as well,
// use GetStringN()
func (o *JsonObject) GetString(key string) string {
	return getObjectScalar(o, convertAnyToString, key)
}

// GetInt retrieves the int value associated with the specified key.
// If the key does not exist, the value is invalid or is null, an error will be set to LastError.
// In case of an error, the zero value will be returned. If you want to regard null values as well,
// use GetIntN()
func (o *JsonObject) GetInt(key string) int {
	return getObjectScalar(o, convertAnyToInt, key)
}

// GetFloat retrieves the float64 value associated with the specified key.
// If the key does not exist, the value is invalid or is null, an error will be set to LastError.
// In case of an error, the zero value will be returned. If you want to regard null values as well,
// use GetFloatN()
func (o *JsonObject) GetFloat(key string) float64 {
	return getObjectScalar(o, convertAnyToFloat, key)
}

// GetBool retrieves the bool value associated with the specified key.
// If the key does not exist, the value is invalid or is null, an error will be set to LastError.
// In case of an error, the zero value will be returned. If you want to regard null values as well,
// use GetBoolN()
func (o *JsonObject) GetBool(key string) bool {
	return getObjectScalar(o, convertAnyToBool, key)
}

// GetStringN is the nullable version of GetString, and returns a pointer instead of a zero value which
// imitates JSON's null as Go's nil. A nil pointer is returned if the key was not found, it is null or
// could not be converted to string. The type of the error will be stored in LastError.
func (o *JsonObject) GetStringN(key string) *string {
	return getObjectScalarN(o, convertAnyToStringN, key)
}

// GetIntN is the nullable version of GetInt, and returns a pointer instead of a zero value which
// imitates JSON's null as Go's nil. A nil pointer is returned if the key was not found, it is null or
// could not be converted to int. The type of the error will be stored in LastError.
func (o *JsonObject) GetIntN(key string) *int {
	return getObjectScalarN(o, convertAnyToIntN, key)
}

// GetFloatN is the nullable version of GetFloat, and returns a pointer instead of a zero value which
// imitates JSON's null as Go's nil. A nil pointer is returned if the key was not found, it is null or
// could not be converted to float64. The type of the error will be stored in LastError.
func (o *JsonObject) GetFloatN(key string) *float64 {
	return getObjectScalarN(o, convertAnyToFloatN, key)
}

// GetBoolN is the nullable version of GetBool, and returns a pointer instead of a zero value which
// imitates JSON's null as Go's nil. A nil pointer is returned if the key was not found, it is null or
// could not be converted to bool. The type of the error will be stored in LastError.
func (o *JsonObject) GetBoolN(key string) *bool {
	return getObjectScalarN(o, convertAnyToBoolN, key)
}

// GetTime retrieves the value as time.Time associated with the specified key.
// If the key does not exist, the value is invalid, null or not a string, an error will be set to LastError.
// In case of an error, the zero value will be returned.
func (o *JsonObject) GetTime(key string) time.Time {
	return getObjectScalar(o, parseTime, key)
}

// GetUUID retrieves the value as uuid.UUID associated with the specified key.
// If the key does not exist, the value is invalid, null or not a string, an error will be set to LastError.
// In case of an error, the zero value will be returned.
func (o *JsonObject) GetUUID(key string) uuid.UUID {
	return getObjectScalar(o, parseUUID, key)
}

// GetObject retrieves a nested JsonObject associated with the specified key.
// If the key does not exist, the value is invalid or is null, an error will be set to LastError.
func (o *JsonObject) GetObject(key string) *JsonObject {
	o.setLastError(nil)
	v, ok := o.object[key]
	if !ok {
		o.setLastError(createKeyNotFoundErr(key))
		return nullObject()
	}
	if v == nil {
		o.setLastError(createTypeConversionErr(nil, JsonObject{}))
		return nullObject()
	}
	switch value := (*v).(type) {
	case map[string]*any:
		return newObjectFromMap(value)
	case map[string]any:
		return newObjectFromMap(convertToMapValuesPtr(value))
	default:
		o.setLastError(createTypeConversionErr(*v, JsonObject{}))
		return nullObject()
	}
}

// GetArray retrieves an array of JsonArray associated with the specified key.
// If the key does not exist, the value is invalid or is null, an error will be set to LastError.
func (o *JsonObject) GetArray(key string) *JsonArray {
	o.setLastError(nil)
	v, ok := o.object[key]
	if !ok {
		o.setLastError(createKeyNotFoundErr(key))
		return nullArray()
	}
	if v == nil {
		o.setLastError(createTypeConversionErr(nil, JsonArray{}))
		return nullArray()
	}
	switch value := (*v).(type) {
	case []any:
		return newArrayFromSlice(convertToSlicePtr(value))
	case []*any:
		return newArrayFromSlice(value)
	default:
		o.setLastError(createTypeConversionErr(*v, JsonArray{}))
		return EmptyArray()
	}
}

// AddJsonObject adds a nested JsonObject to the JsonObject associated with the key.
func (o *JsonObject) AddJsonObject(key string, jsonObject *JsonObject) {
	var object any = jsonObject.object
	o.object[key] = &object
}

// AddJsonArray adds a JsonArray to the JsonObject associated with the key.
func (o *JsonObject) AddJsonArray(key string, jsonArray *JsonArray) {
	var elements any = jsonArray.elements
	o.object[key] = &elements
}

// AddString adds a string to the JsonObject associated with the key.
func (o *JsonObject) AddString(key string, s string) {
	var value any = s
	o.object[key] = &value
}

// AddInt adds an int to the JsonObject associated with the key.
func (o *JsonObject) AddInt(key string, i int) {
	var value any = i
	o.object[key] = &value
}

// AddFloat adds a float to the JsonObject associated with the key.
func (o *JsonObject) AddFloat(key string, f float64) {
	var value any = f
	o.object[key] = &value
}

// AddBool adds a bool to the JsonObject associated with the key.
func (o *JsonObject) AddBool(key string, b bool) {
	var value any = b
	o.object[key] = &value
}

// AddStringArray adds a string array to the JsonObject associated with the key.
func (o *JsonObject) AddStringArray(key string, s []string) {
	var value any = s
	o.object[key] = &value
}

// AddIntArray adds an int array to the JsonObject associated with the key.
func (o *JsonObject) AddIntArray(key string, i []int) {
	var value any = i
	o.object[key] = &value
}

// AddFloatArray adds a float array to the JsonObject associated with the key.
func (o *JsonObject) AddFloatArray(key string, f []float64) {
	var value any = f
	o.object[key] = &value
}

// AddNull adds nil to the JsonObject associated with the key.
func (o *JsonObject) AddNull(key string) {
	o.object[key] = nil
}

// TransformKeys returns a new JsonObject with transformed keys. It takes a
// transformation function as parameter f that takes a string, the original key,
// and returns a new string, the new key.
func (o *JsonObject) TransformKeys(f func(string) string) *JsonObject {
	return newObjectFromMap(transformKeys(o.object, f))
}

// ForEach applies the provided function to each key-value pair in the JsonObject.
func (o *JsonObject) ForEach(f func(key string, j JsonMapper)) {
	for k, element := range o.object {
		f(k, getMapperFromField(element))
	}
}

// Filter returns a new JsonObject containing only the key-value pairs for which the provided function returns true.
func (o *JsonObject) Filter(f func(key string, j JsonMapper) bool) *JsonObject {
	var obj = EmptyObject()
	for k, element := range o.object {
		if f(k, getMapperFromField(element)) {
			obj.object[k] = element
		}
	}
	return obj
}

// PrettyString returns a pretty-printed string representation of the JsonObject.
func (o *JsonObject) PrettyString() string {
	jsonBytes, _ := marshalIndent(o.object)
	return string(jsonBytes)
}

// String returns a string representation of the JsonObject in JSON format.
// Returns: A string containing the JSON representation.
func (o *JsonObject) String() string {
	jsonBytes, _ := marshal(o.object)
	return string(jsonBytes)
}

// SetLastError sets the LastError field of the JsonObject to the provided error.
func (o *JsonObject) setLastError(err error) {
	o.LastError = err
}

// newObjectFromMap initializes and returns a new instance of JsonObject.
func newObjectFromMap(data map[string]*any) *JsonObject {
	var obj JsonObject
	obj.object = data
	return &obj
}
