package jsonmapper

import (
	"fmt"
	"time"
)

// JsonObject represents a JSON object
type JsonObject struct {
	object    map[string]*any
	LastError error
}

// NewObject initializes and returns a new instance of JsonObject.
func NewObject(data map[string]*any) *JsonObject {
	var obj JsonObject
	obj.object = data
	return &obj
}

// EmptyObject initializes and returns an empty new instance of JsonObject.
func EmptyObject() *JsonObject {
	var obj JsonObject
	obj.object = make(map[string]*any)
	return &obj
}

// Length returns the number of elements in the JsonObject.
func (o *JsonObject) Length() int {
	return len(o.object)
}

// Has checks if the specified key exists in the JsonObject.
func (o *JsonObject) Has(key string) bool {
	_, ok := o.object[key]
	return ok
}

// IsEmpty checks if the JSON object has no fields in it
func (o *JsonObject) IsEmpty() bool {
	return len(o.object) == 0
}

// IsNull checks if the value associated with the key is null
func (o *JsonObject) IsNull(key string) bool {
	v := o.object[key]
	return v == nil
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

// AsStringMap returns the object as map of string, string
func (o *JsonObject) AsStringMap() map[string]string {
	return getGenericMap(convertAnyToString, *o)
}

// AsIntMap returns the object as map of (string, int)
func (o *JsonObject) AsIntMap() map[string]int {
	return getGenericMap(convertAnyToInt, *o)
}

// AsFloatMap returns the object as map of (string, float64)
func (o *JsonObject) AsFloatMap() map[string]float64 {
	return getGenericMap(convertAnyToFloat, *o)
}

// AsArrayMap returns the object as map of (string, JsonArray)
func (o *JsonObject) AsArrayMap() map[string]JsonArray {
	return getGenericMap(convertAnyToArray, *o)
}

// AsObjectMap returns the object as map of (string, JsonObject)
func (o *JsonObject) AsObjectMap() map[string]JsonObject {
	return getGenericMap(convertAnyToObject, *o)
}

// Get retrieves the value associated with the key and returns it as a JsonMapper
func (o *JsonObject) Get(key string) JsonMapper {
	return getMapperFromField(o.object[key])
}

// GetString retrieves the string value associated with the specified key.
// If the key does not exist, the value is invalid or is null, an error will be set to LastError.
func (o *JsonObject) GetString(key string) string {
	return getObjectScalar(o, convertAnyToString, key, stringTypeStr)
}

// GetInt retrieves the int value associated with the specified key.
// If the key does not exist, the value is invalid or is null, an error will be set to LastError.
func (o *JsonObject) GetInt(key string) int {
	return getObjectScalar(o, convertAnyToInt, key, intTypeStr)
}

// GetFloat retrieves the float64 value associated with the specified key.
// If the key does not exist, the value is invalid or is null, an error will be set to LastError.
func (o *JsonObject) GetFloat(key string) float64 {
	return getObjectScalar(o, convertAnyToFloat, key, floatTypeStr)
}

// GetBool retrieves the bool value associated with the specified key.
// If the key does not exist, the value is invalid or is null, an error will be set to LastError.
func (o *JsonObject) GetBool(key string) bool {
	return getObjectScalar(o, convertAnyToBool, key, boolTypeStr)
}

// GetTime retrieves the time.Time value associated with the specified key.
// If the key does not exist, the value is invalid or is null, an error will be set to LastError.
func (o *JsonObject) GetTime(key string) time.Time {
	return getObjectScalar(o, parseTime, key, timeTypeStr)
}

// GetObject retrieves a nested JsonObject associated with the specified key.
// If the key does not exist, the value is invalid or is null, an error will be set to LastError.
func (o *JsonObject) GetObject(key string) *JsonObject {
	v, ok := o.object[key]
	if !ok {
		o.setLastError(createKeyNotFoundErr(key))
		return EmptyObject()
	}
	if v == nil {
		o.setLastError(createNullConversionErr(objectTypeStr))
		return EmptyObject()
	}
	switch (*v).(type) {
	case map[string]*any:
		data := (*v).(map[string]*any)
		return NewObject(data)
	case map[string]any:
		dataPtr := convertToMapValuesPtr((*v).(map[string]any))
		return NewObject(dataPtr)
	default:
		o.setLastError(createTypeConversionErr(*v, objectTypeStr))
		return EmptyObject()
	}
}

// GetArray retrieves an array of JsonArray associated with the specified key.
// If the key does not exist, the value is invalid or is null, an error will be set to LastError.
func (o *JsonObject) GetArray(key string) *JsonArray {
	v, ok := o.object[key]
	if !ok {
		o.setLastError(createKeyNotFoundErr(key))
		return EmptyArray()
	}
	if v == nil {
		o.setLastError(createNullConversionErr(arrayTypeStr))
		return EmptyArray()
	}
	switch (*v).(type) {
	case []any:
		return NewArray(convertToSlicePtr((*v).([]any)))
	case []*any:
		return NewArray((*v).([]*any))
	default:
		o.setLastError(createTypeConversionErr(*v, arrayTypeStr))
		return EmptyArray()
	}
}

// Find searches for a key in the JsonObject and its nested objects.
// Returns the JsonMapper associated with the key if found; otherwise, returns an empty JsonMapper.
func (o *JsonObject) Find(key string) JsonMapper {
	for k, v := range o.object {
		field := getMapperFromField(v)
		if k == key {
			return field
		}
		if field.IsObject {
			return field.AsObject.Find(key)
		}
	}
	return JsonMapper{}
}

// AddKeyValue adds a key-value pair to the JsonObject.
func (o *JsonObject) AddKeyValue(k string, value any) {
	switch value := value.(type) {
	case JsonObject:
		var object any = value.object
		o.object[k] = &object
	case *JsonObject:
		var object any = value.object
		o.object[k] = &object
	case JsonArray:
		var elements any = value.elements
		o.object[k] = &elements
	case *JsonArray:
		var elements any = value.elements
		o.object[k] = &elements
	case nil, string, int, float64, bool, []string, []int, []float64, []bool:
		o.object[k] = &value
	default:
		o.setLastError(fmt.Errorf("could not add value of type %T", value))
	}
}

// ForEach applies the provided function to each key-value pair in the JsonObject.
func (o *JsonObject) ForEach(f func(key string, j JsonMapper)) {
	for k, element := range o.object {
		f(k, getMapperFromField(element))
	}
}

// Filter returns a new JsonObject containing only the key-value pairs for which the provided function returns true.
func (o *JsonObject) Filter(f func(key string, j JsonMapper) bool) JsonObject {
	var obj = EmptyObject()
	for k, element := range o.object {
		if f(k, getMapperFromField(element)) {
			obj.object[k] = element
		}
	}
	return *obj
}

// TransformObjectKeys returns a new JsonObject with transformed keys, where keys are converted to snake_case.
func (o *JsonObject) TransformObjectKeys() JsonObject {
	return *NewObject(transformKeys(o.object))
}

// SetLastError sets the LastError field of the JsonObject to the provided error.
func (o *JsonObject) setLastError(err error) {
	o.LastError = err
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

func transformKeys(m map[string]*any) map[string]*any {
	newMap := make(map[string]*any)
	for key, value := range m {
		newKey := toSnakeCase(key)
		if value == nil {
			newMap[newKey] = nil
			continue
		}
		nestedMap, ok := (*value).(map[string]any)
		if ok {
			nestedResult := transformKeys(convertToMapValuesPtr(nestedMap))
			var nestedInterface any = nestedResult
			newMap[newKey] = &nestedInterface
		} else {
			newMap[newKey] = value
		}
	}
	return newMap
}
