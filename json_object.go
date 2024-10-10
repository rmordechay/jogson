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

// NewObject creates and returns a new JsonObject.
func NewObject() JsonObject {
	var obj JsonObject
	obj.object = make(map[string]*any)
	return obj
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

// Values returns a slice of all values in the JsonObject as Json types.
func (o *JsonObject) Values() []Json {
	values := make([]Json, 0, len(o.object))
	for _, v := range o.object {
		values = append(values, getMapperFromField(v))
	}
	return values
}

// Elements returns a map of all elements in the JsonObject with their keys.
func (o *JsonObject) Elements() map[string]Json {
	jsons := make(map[string]Json)
	for k, v := range o.object {
		jsons[k] = getMapperFromField(v)
	}
	return jsons
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
		return getAsString(v, o)
	}
	o.LastError = fmt.Errorf(keyNotFoundErrStr, key)
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
			o.LastError = fmt.Errorf(nullConversionErrStr, "int")
			return 0
		}
		return getAsInt(v, o)
	}
	o.LastError = fmt.Errorf(keyNotFoundErrStr, key)
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
			o.LastError = fmt.Errorf(nullConversionErrStr, "float64")
			return 0
		}
		return getAsFloat(v, o)
	}
	o.LastError = fmt.Errorf(keyNotFoundErrStr, key)
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
			o.LastError = fmt.Errorf(nullConversionErrStr, "bool")
			return false
		}
		return getAsBool(v, o)
	}
	o.LastError = fmt.Errorf(keyNotFoundErrStr, key)
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
			return time.Time{}, fmt.Errorf(nullConversionErrStr, "time.Time")
		}
		return parseTime(v)
	}
	return time.Time{}, fmt.Errorf(keyNotFoundErrStr, key)
}

// GetObject retrieves a nested JsonObject associated with the specified key.
// If the key does not exist or the value is not a valid object, it sets LastError.
func (o *JsonObject) GetObject(key string) *JsonObject {
	for k, v := range o.object {
		if k != key {
			continue
		}
		if v == nil {
			o.LastError = fmt.Errorf(nullConversionErrStr, "JsonObject")
			return &JsonObject{}
		}
		switch (*v).(type) {
		case map[string]*any:
			object := (*v).(map[string]*any)
			return &JsonObject{object: object}
		case map[string]any:
			dataPtr := convertToMapValuesPtr((*v).(map[string]any))
			return &JsonObject{object: dataPtr}
		default:
			o.SetLastError(fmt.Errorf(typeConversionErrStr, *v, "JsonObject"))
			return &JsonObject{}
		}
	}
	o.LastError = fmt.Errorf(keyNotFoundErrStr, key)
	return &JsonObject{}
}

// GetArray retrieves an array of JsonArray associated with the specified key.
// If the key does not exist or the value is not a valid array, it sets LastError.
func (o *JsonObject) GetArray(key string) *JsonArray {
	for k, v := range o.object {
		if k != key {
			continue
		}
		if v == nil {
			o.LastError = fmt.Errorf(nullConversionErrStr, "JsonArray")
			return &JsonArray{}
		}
		switch (*v).(type) {
		case []any:
			return &JsonArray{elements: convertToSlicePtr((*v).([]any))}
		case []*any:
			return &JsonArray{elements: (*v).([]*any)}
		default:
			o.SetLastError(fmt.Errorf(typeConversionErrStr, *v, "[]*any"))
			return &JsonArray{}
		}
	}
	o.LastError = fmt.Errorf(keyNotFoundErrStr, key)
	return &JsonArray{}
}

// Find searches for a key in the JsonObject and its nested objects.
// Returns the Json associated with the key if found; otherwise, returns an empty Json.
func (o *JsonObject) Find(key string) Json {
	for k, v := range o.object {
		field := getMapperFromField(v)
		if k == key {
			return field
		}
		if field.IsObject {
			return field.Object.Find(key)
		}
	}
	return Json{}
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
		o.LastError = fmt.Errorf("could not add value of type %T", value)
	}
}

// ForEach applies the provided function to each key-value pair in the JsonObject.
func (o *JsonObject) ForEach(f func(key string, j Json)) {
	for k, element := range o.object {
		f(k, getMapperFromField(element))
	}
}

// Filter returns a new JsonObject containing only the key-value pairs for which the provided function returns true.
func (o *JsonObject) Filter(f func(key string, j Json) bool) JsonObject {
	var obj = NewObject()
	for k, element := range o.object {
		if f(k, getMapperFromField(element)) {
			obj.object[k] = element
		}
	}
	return obj
}

// SetLastError sets the LastError field of the JsonObject to the provided error.
func (o *JsonObject) SetLastError(err error) {
	o.LastError = err
}

// TransformObjectKeys returns a new JsonObject with transformed keys, where keys are converted to snake_case.
func (o *JsonObject) TransformObjectKeys() JsonObject {
	return JsonObject{object: transformKeys(o.object)}
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

func getAsJsonObject(data *any, j JsonError) JsonObject {
	if data == nil {
		j.SetLastError(fmt.Errorf(nullConversionErrStr, "string"))
		return JsonObject{}
	}
	v, ok := (*data).(map[string]any)
	if !ok {
		j.SetLastError(fmt.Errorf(typeConversionErrStr, data, "JsonObject"))
		return JsonObject{}
	}

	var obj JsonObject
	var object = make(map[string]*any)
	for key, value := range v {
		object[key] = &value
	}
	obj.object = object
	return obj
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
