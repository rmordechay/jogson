package jsonmapper

import (
	"fmt"
	"time"
)

type JsonObject struct {
	object    map[string]*interface{}
	LastError error
}

func NewObject() JsonObject {
	var obj JsonObject
	obj.object = make(map[string]*interface{})
	return obj
}

func (o *JsonObject) Keys() []string {
	keys := make([]string, 0, len(o.object))
	for key := range o.object {
		keys = append(keys, key)
	}
	return keys
}

func (o *JsonObject) Values() []Json {
	values := make([]Json, 0, len(o.object))
	for _, v := range o.object {
		values = append(values, getMapperFromField(v))
	}
	return values
}

func (o *JsonObject) Elements() map[string]Json {
	jsons := make(map[string]Json)
	for k, v := range o.object {
		jsons[k] = getMapperFromField(v)
	}
	return jsons
}

func (o *JsonObject) Has(key string) bool {
	for k := range o.object {
		if k == key {
			return true
		}
	}
	return false
}

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

func (o *JsonObject) GetInt(key string) int {
	for k, v := range o.object {
		if k != key {
			continue
		}
		if v == nil {
			o.LastError = fmt.Errorf(nullConversionErrStr, 0)
			return 0
		}
		return getAsInt(v, o)
	}
	o.LastError = fmt.Errorf(keyNotFoundErrStr, key)
	return 0
}

func (o *JsonObject) GetFloat(key string) float64 {
	for k, v := range o.object {
		if k != key {
			continue
		}
		if v == nil {
			o.LastError = fmt.Errorf(nullConversionErrStr, 0.0)
			return 0
		}
		return getAsFloat(v, o)
	}
	o.LastError = fmt.Errorf(keyNotFoundErrStr, key)
	return 0
}

func (o *JsonObject) GetBool(key string) bool {
	for k, v := range o.object {
		if k != key {
			continue
		}
		if v == nil {
			o.LastError = fmt.Errorf(nullConversionErrStr, false)
			return false
		}
		return getAsBool(v, o)
	}
	o.LastError = fmt.Errorf(keyNotFoundErrStr, key)
	return false
}

func (o *JsonObject) GetTime(key string) (time.Time, error) {
	for k, v := range o.object {
		if k != key {
			continue
		}
		if v == nil {
			return time.Time{}, fmt.Errorf(nullConversionErrStr, time.Time{})
		}
		return parseTime(v)
	}
	return time.Time{}, fmt.Errorf(keyNotFoundErrStr, key)
}

func (o *JsonObject) GetObject(key string) *JsonObject {
	for k, v := range o.object {
		if k != key {
			continue
		}
		if v == nil {
			o.LastError = fmt.Errorf(nullConversionErrStr, JsonObject{})
			return &JsonObject{}
		}
		switch (*v).(type) {
		case map[string]*interface{}:
			return &JsonObject{object: (*v).(map[string]*interface{})}
		case map[string]interface{}:
			return &JsonObject{object: convertToMapValuesPtr((*v).(map[string]interface{}))}
		default:
			o.LastError = fmt.Errorf(typeConversionErrStr, *v, JsonObject{})
			return &JsonObject{}
		}
	}
	o.LastError = fmt.Errorf(keyNotFoundErrStr, key)
	return &JsonObject{}
}

func (o *JsonObject) GetArray(key string) *JsonArray {
	for k, v := range o.object {
		if k != key {
			continue
		}
		if v == nil {
			o.LastError = fmt.Errorf(nullConversionErrStr, JsonArray{})
			return &JsonArray{}
		}
		jsonArray, ok := (*v).([]interface{})
		if !ok {
			o.LastError = fmt.Errorf(typeConversionErrStr, *v, JsonArray{})
			return &JsonArray{}
		}
		return &JsonArray{elements: convertToSlicePtr(jsonArray)}
	}
	o.LastError = fmt.Errorf(keyNotFoundErrStr, key)
	return &JsonArray{}
}

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

func (o *JsonObject) AddKeyValue(k string, value interface{}) {
	o.object[k] = &value
}

func (o *JsonObject) ForEach(f func(key string, j Json)) {
	for k, element := range o.object {
		f(k, getMapperFromField(element))
	}
}

func (o *JsonObject) Filter(f func(key string, j Json) bool) JsonObject {
	var obj = NewObject()
	for k, element := range o.object {
		if f(k, getMapperFromField(element)) {
			obj.object[k] = element
		}
	}
	return obj
}

func (o *JsonObject) SetLastError(err error) {
	o.LastError = err
}

func (o *JsonObject) TransformObjectKeys() JsonObject {
	return JsonObject{object: transformKeys(o.object)}
}

func (o *JsonObject) PrettyString() string {
	jsonBytes, _ := marshalIndent(o.object)
	return string(jsonBytes)
}

func (o *JsonObject) String() string {
	jsonBytes, _ := marshal(o.object)
	return string(jsonBytes)
}

func getAsJsonObject(data *interface{}, j JsonError) JsonObject {
	if data == nil {
		j.SetLastError(fmt.Errorf(nullConversionErrStr, ""))
		return JsonObject{}
	}
	v, ok := (*data).(map[string]interface{})
	if !ok {
		j.SetLastError(fmt.Errorf(typeConversionErrStr, data, JsonObject{}))
		return JsonObject{}
	}

	var obj JsonObject
	var object = make(map[string]*interface{})
	for key, value := range v {
		object[key] = &value
	}
	obj.object = object
	return obj
}

func transformKeys(m map[string]*interface{}) map[string]*interface{} {
	newMap := make(map[string]*interface{})
	for key, value := range m {
		newKey := toSnakeCase(key)
		if value == nil {
			newMap[newKey] = value
			continue
		}
		nestedMap, ok := (*value).(map[string]interface{})
		if ok {
			nestedResult := transformKeys(convertToMapValuesPtr(nestedMap))
			var nestedInterface interface{} = nestedResult
			newMap[newKey] = &nestedInterface
		} else {
			newMap[newKey] = value
		}
	}
	return newMap
}
