package jsonmapper

import (
	"strconv"
)

func getAsString(data *any, j jsonEntity) string {
	if data == nil {
		j.SetLastError(NewNullConversionErr("string"))
		return ""
	}
	switch (*data).(type) {
	case string:
		return (*data).(string)
	case float64:
		return strconv.FormatFloat((*data).(float64), 'f', -1, 64)
	case int:
		return strconv.Itoa((*data).(int))
	case bool:
		return strconv.FormatBool((*data).(bool))
	default:
		j.SetLastError(NewTypeConversionErr(*data, "string"))
		return ""
	}
}

func getAsInt(data *any, j jsonEntity) int {
	if data == nil {
		j.SetLastError(NewNullConversionErr("int"))
		return 0
	}
	switch (*data).(type) {
	case float64:
		return int((*data).(float64))
	case int:
		return (*data).(int)
	default:
		j.SetLastError(NewTypeConversionErr(*data, "int"))
		return 0
	}
}

func getAsFloat(data *any, j jsonEntity) float64 {
	if data == nil {
		j.SetLastError(NewNullConversionErr("float64"))
		return 0
	}
	v, ok := (*data).(float64)
	if !ok {
		j.SetLastError(NewTypeConversionErr(*data, "float64"))
		return 0
	}
	return v
}

func getAsBool(data *any, j jsonEntity) bool {
	if data == nil {
		j.SetLastError(NewNullConversionErr("bool"))
		return false
	}
	v, ok := (*data).(bool)
	if !ok {
		j.SetLastError(NewTypeConversionErr(*data, "bool"))
		return false
	}
	return v
}

func getAsJsonObject(data *any, j jsonEntity) JsonObject {
	if data == nil {
		j.SetLastError(NewNullConversionErr("string"))
		return JsonObject{}
	}
	v, ok := (*data).(map[string]any)
	if !ok {
		j.SetLastError(NewTypeConversionErr(data, "JsonObject"))
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
