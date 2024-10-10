package jsonmapper

import (
	"strconv"
)

type F[T any] func(data *any, j jsonEntity) T

func asGenericMap[T any](f F[T], o JsonObject) map[string]T {
	genericMap := make(map[string]T)
	for k, v := range o.object {
		genericMap[k] = f(v, &o)
	}
	return genericMap
}

func asGenericArray[T any](f F[T], o JsonArray) []T {
	arr := make([]T, 0, len(o.elements))
	for _, element := range o.elements {
		arr = append(arr, f(element, &o))
	}
	return arr
}

func convertAnyToString(data *any, j jsonEntity) string {
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

func convertAnyToInt(data *any, j jsonEntity) int {
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

func convertAnyToFloat(data *any, j jsonEntity) float64 {
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

func convertAnyToBool(data *any, j jsonEntity) bool {
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

func convertAnyToObject(data *any, j jsonEntity) JsonObject {
	if data == nil {
		j.SetLastError(NewNullConversionErr("JsonObject"))
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

func convertAnyToArray(data *any, j jsonEntity) JsonArray {
	if data == nil {
		j.SetLastError(NewNullConversionErr("JsonArray"))
		return JsonArray{}
	}
	v, ok := (*data).([]any)
	if !ok {
		j.SetLastError(NewTypeConversionErr(data, "JsonArray"))
		return JsonArray{}
	}

	var array JsonArray
	var elements = make([]*any, 0, len(v))
	for _, value := range v {
		value := value
		elements = append(elements, &value)
	}
	array.elements = elements
	return array
}

func convertSliceToJsonArray[T any](data []T) JsonArray {
	var jsonArray JsonArray
	sliceAnyPtr := make([]*any, 0, len(data))
	for _, v := range data {
		var valAny any = v
		sliceAnyPtr = append(sliceAnyPtr, &valAny)
	}
	jsonArray.elements = sliceAnyPtr
	return jsonArray
}
