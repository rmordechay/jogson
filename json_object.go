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
	return asGenericMap(convertAnyToString, *o)
}

// AsIntMap returns the object as map of (string, int)
func (o *JsonObject) AsIntMap() map[string]int {
	return asGenericMap(convertAnyToInt, *o)
}

// AsFloatMap returns the object as map of (string, float64)
func (o *JsonObject) AsFloatMap() map[string]float64 {
	return asGenericMap(convertAnyToFloat, *o)
}

// As2DMap returns the object as map of (string, JsonArray)
func (o *JsonObject) As2DMap() map[string]JsonArray {
	return asGenericMap(convertAnyToArray, *o)
}

// AsObjectMap returns the object as map of (string, JsonObject)
func (o *JsonObject) AsObjectMap() map[string]JsonObject {
	return asGenericMap(convertAnyToObject, *o)
}

// Has checks if the specified key exists in the JsonObject.
func (o *JsonObject) Has(key string) bool {
	for k := range o.object {
		if k == key {
			return true
		}
	}
	return false
}

// GetString retrieves the string value associated with the specified key.
// If the key does not exist or the value is not a string, it sets LastError.
func (o *JsonObject) GetString(key string) string {
	for k, v := range o.object {
		if k != key {
			continue
		}
		return convertAnyToString(v, o)
	}
	o.SetLastError(NewKeyNotFoundErr(key))
	return ""
}

// GetInt retrieves the int value associated with the specified key.
// If the key does not exist or the value is not an int, it sets LastError.
func (o *JsonObject) GetInt(key string) int {
	for k, v := range o.object {
		if k != key {
			continue
		}
		if v == nil {
			o.SetLastError(NewNullConversionErr("int"))
			return 0
		}
		return convertAnyToInt(v, o)
	}
	o.SetLastError(NewKeyNotFoundErr(key))
	return 0
}

// GetFloat retrieves the float64 value associated with the specified key.
// If the key does not exist or the value is not a float, it sets LastError.
func (o *JsonObject) GetFloat(key string) float64 {
	for k, v := range o.object {
		if k != key {
			continue
		}
		if v == nil {
			o.SetLastError(NewNullConversionErr("float64"))
			return 0
		}
		return convertAnyToFloat(v, o)
	}
	o.SetLastError(NewKeyNotFoundErr(key))
	return 0
}

// GetBool retrieves the bool value associated with the specified key.
// If the key does not exist or the value is not a bool, it sets LastError.
func (o *JsonObject) GetBool(key string) bool {
	for k, v := range o.object {
		if k != key {
			continue
		}
		if v == nil {
			o.SetLastError(NewNullConversionErr("bool"))
			return false
		}
		return convertAnyToBool(v, o)
	}
	o.SetLastError(NewKeyNotFoundErr(key))
	return false
}

// GetTime retrieves the time.Time value associated with the specified key.
// If the key does not exist or the value is not a valid time, it returns an error.
func (o *JsonObject) GetTime(key string) (time.Time, error) {
	for k, v := range o.object {
		if k != key {
			continue
		}
		if v == nil {
			return time.Time{}, NewNullConversionErr("time.Time")
		}
		return parseTime(v)
	}
	return time.Time{}, NewKeyNotFoundErr(key)
}

// GetObject retrieves a nested JsonObject associated with the specified key.
// If the key does not exist or the value is not a valid object, it sets LastError.
func (o *JsonObject) GetObject(key string) *JsonObject {
	for k, v := range o.object {
		if k != key {
			continue
		}
		if v == nil {
			o.SetLastError(NewNullConversionErr("JsonObject"))
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
			o.SetLastError(NewTypeConversionErr(*v, "JsonObject"))
			return EmptyObject()
		}
	}
	o.SetLastError(NewKeyNotFoundErr(key))
	return EmptyObject()
}

// GetArray retrieves an array of JsonArray associated with the specified key.
// If the key does not exist or the value is not a valid array, it sets LastError.
func (o *JsonObject) GetArray(key string) *JsonArray {
	for k, v := range o.object {
		if k != key {
			continue
		}
		if v == nil {
			o.SetLastError(NewNullConversionErr("JsonArray"))
			return EmptyArray()
		}
		switch (*v).(type) {
		case []any:
			return NewArray(convertToSlicePtr((*v).([]any)))
		case []*any:
			return NewArray((*v).([]*any))
		default:
			o.SetLastError(NewTypeConversionErr(*v, "[]*any"))
			return EmptyArray()
		}
	}
	o.SetLastError(NewKeyNotFoundErr(key))
	return EmptyArray()
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
	switch value.(type) {
	case JsonObject:
		var object any = value.(JsonObject).object
		o.object[k] = &object
	case *JsonObject:
		var object any = value.(JsonObject).object
		o.object[k] = &object
	case JsonArray:
		var elements any = value.(JsonArray).elements
		o.object[k] = &elements
	case *JsonArray:
		var elements any = value.(*JsonArray).elements
		o.object[k] = &elements
	case nil, string, int, float64, bool, []string, []int, []float64, []bool:
		o.object[k] = &value
	default:
		o.SetLastError(fmt.Errorf("could not add value of type %T", value))
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
func (o *JsonObject) SetLastError(err error) {
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
			newMap[newKey] = value
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
