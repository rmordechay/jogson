package jsonmapper

import (
	"fmt"
	"time"
)

var (
	NullConversionErrStr = "value is null and could not be converted to %T"
	TypeConversionErrStr = "the type '%T' could not be converted to %T"
	KeyNotFoundErrStr    = "the requested key '%v' was not found"
	IndexOutOfRange      = "index out of range [%v] with length %v"
)

type JsonObject struct {
	object map[string]*interface{}

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

func (o *JsonObject) Values() []*Mapper {
	values := make([]*Mapper, 0, len(o.object))
	for _, v := range o.object {
		values = append(values, getMapperFromField(v))
	}
	return values
}

func (o *JsonObject) Has(key string) bool {
	for k := range o.object {
		if k == key {
			return true
		}
	}
	return false
}

func (o *JsonObject) GetTime(key string) (time.Time, error) {
	for k, v := range o.object {
		if k != key {
			continue
		}
		if v == nil {
			return time.Time{}, fmt.Errorf(NullConversionErrStr, time.Time{})
		}
		return parseTime(k, v)
	}
	return time.Time{}, fmt.Errorf(KeyNotFoundErrStr, key)
}

func (o *JsonObject) GetString(key string) string {
	for k, v := range o.object {
		if k != key {
			continue
		}
		return getAsString(v, o)
	}
	o.LastError = fmt.Errorf(KeyNotFoundErrStr, key)
	return ""
}

func (o *JsonObject) SetLastError(err error) {
	o.LastError = err
}

func (o *JsonObject) GetInt(key string) int {
	for k, v := range o.object {
		if k != key {
			continue
		}
		if v == nil {
			o.LastError = fmt.Errorf(NullConversionErrStr, 0)
			return 0
		}
		return getAsInt(v, o)
	}
	o.LastError = fmt.Errorf(KeyNotFoundErrStr, key)
	return 0
}

func (o *JsonObject) GetFloat(key string) float64 {
	for k, v := range o.object {
		if k != key {
			continue
		}
		if v == nil {
			o.LastError = fmt.Errorf(NullConversionErrStr, 0.0)
			return 0
		}
		return getAsFloat(v, o)
	}
	o.LastError = fmt.Errorf(KeyNotFoundErrStr, key)
	return 0
}

func (o *JsonObject) GetBool(key string) bool {
	for k, v := range o.object {
		if k != key {
			continue
		}
		if v == nil {
			o.LastError = fmt.Errorf(NullConversionErrStr, false)
			return false
		}
		return getAsBool(v, o)
	}
	o.LastError = fmt.Errorf(KeyNotFoundErrStr, key)
	return false
}

func (o *JsonObject) GetObject(key string) JsonObject {
	for k, v := range o.object {
		if k != key {
			continue
		}
		if v == nil {
			o.LastError = fmt.Errorf(NullConversionErrStr, JsonObject{})
			return JsonObject{}
		}
		jsonObject, ok := (*v).(map[string]interface{})
		if !ok {
			o.LastError = fmt.Errorf(TypeConversionErrStr, *v, JsonObject{})
			return JsonObject{}
		}
		return JsonObject{object: convertToMapValuesPtr(jsonObject)}
	}
	o.LastError = fmt.Errorf(KeyNotFoundErrStr, key)
	return JsonObject{}
}

func (o *JsonObject) GetArray(key string) JsonArray {
	for k, v := range o.object {
		if k != key {
			continue
		}
		if v == nil {
			o.LastError = fmt.Errorf(NullConversionErrStr, false)
			return JsonArray{}
		}
		jsonArray, ok := (*v).([]interface{})
		if !ok {
			o.LastError = fmt.Errorf(TypeConversionErrStr, *v, false)
			return JsonArray{}
		}
		return JsonArray{elements: convertToSlicePtr(jsonArray)}
	}
	o.LastError = fmt.Errorf(KeyNotFoundErrStr, key)
	return JsonArray{}
}

func (o *JsonObject) Find(key string) *Mapper {
	for k, v := range o.object {
		field := getMapperFromField(v)
		if k == key {
			return field
		}
		if field.IsObject {
			return field.object.Find(key)
		}
	}
	return &Mapper{}
}

func (o *JsonObject) Elements() map[string]*Mapper {
	jsons := make(map[string]*Mapper)
	for k, v := range o.object {
		jsons[k] = getMapperFromField(v)
	}
	return jsons
}

func (o *JsonObject) AddKeyValue(k string, value interface{}) {
	o.object[k] = &value
}

func (o *JsonObject) ForEach(f func(key string, mapper Mapper)) {
	for k, element := range o.object {
		f(k, *getMapperFromField(element))
	}
}

func (o *JsonObject) Filter(f func(key string, mapper Mapper) bool) JsonObject {
	var obj = NewObject()
	for k, element := range o.object {
		field := getMapperFromField(element)
		if f(k, *field) {
			obj.object[k] = element
		}
	}
	return obj
}

func (o *JsonObject) PrettyString() string {
	jsonBytes, _ := marshalIndent(o.object)
	return string(jsonBytes)
}

func (o *JsonObject) String() string {
	jsonBytes, _ := marshal(o.object)
	return string(jsonBytes)
}

func parseTime(k string, t *interface{}) (time.Time, error) {
	if t == nil {
		return time.Time{}, fmt.Errorf(NullConversionErrStr, "")
	}
	timeAsString, ok := (*t).(string)
	if !ok {
		return time.Time{}, fmt.Errorf("cannot convert type %T to type time.Time\n", t)
	}
	for _, layout := range timeLayouts {
		parsedTime, err := time.Parse(layout, timeAsString)
		if err == nil {
			return parsedTime, nil
		}
	}
	return time.Time{}, fmt.Errorf("the value '%v' could not be converted to type time.Time", timeAsString)
}
