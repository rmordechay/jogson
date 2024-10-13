package jsonmapper

import (
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"strconv"
	"time"
	"unicode"
)

var jsonIter = jsoniter.ConfigCompatibleWithStandardLibrary

// jc is JSON converter function type that convert type any to generic type T
type jc[T any] func(data *any, j jsonI) T

func getMapperFromField(data *any) JsonMapper {
	if data == nil {
		return JsonMapper{IsNull: true}
	}

	var mapper JsonMapper
	switch value := (*data).(type) {
	case bool:
		mapper.IsBool = true
		mapper.AsBool = value
	case int:
		mapper.IsInt = true
		mapper.AsInt = value
	case float64:
		if value == float64(int(value)) {
			mapper.IsInt = true
			mapper.AsInt = int(value)
		} else {
			mapper.IsFloat = true
		}
		mapper.AsFloat = value
	case string:
		mapper.IsString = true
		mapper.AsString = value
	case map[string]any:
		mapper.IsObject = true
		mapper.AsObject = convertAnyToObject(data, nil)
	case []float64:
		mapper.IsArray = true
		mapper.AsArray = convertSliceToJsonArray(value)
	case []int:
		mapper.IsArray = true
		mapper.AsArray = convertSliceToJsonArray(value)
	case []string:
		mapper.IsArray = true
		mapper.AsArray = convertSliceToJsonArray(value)
	case []bool:
		mapper.IsArray = true
		mapper.AsArray = convertSliceToJsonArray(value)
	case []*any:
		mapper.IsArray = true
		mapper.AsArray = *newArrayFromSlice(value)
	case []any:
		mapper.IsArray = true
		mapper.AsArray = *newArrayFromSlice(convertToSlicePtr(value))
	case nil:
		mapper.IsNull = true
	}
	return mapper
}

func convertToSlicePtr(data []any) []*any {
	array := make([]*any, len(data))
	for i, v := range data {
		v := v
		array[i] = &v
	}
	return array
}

func convertToMapValuesPtr(data map[string]any) map[string]*any {
	jsonObject := make(map[string]*any, len(data))
	for k, v := range data {
		v := v
		jsonObject[k] = &v
	}
	return jsonObject
}

func dataStartsWith(data []byte, brackOrParen byte) bool {
	if len(data) == 0 {
		return false
	}
	var firstChar byte
	for _, d := range data {
		firstChar = d
		if unicode.IsSpace(rune(firstChar)) {
			continue
		}
		return firstChar == brackOrParen
	}
	return false
}

func parseTime(t *any, j jsonI) time.Time {
	if t == nil {
		j.setLastError(createTypeConversionErr(nil, ""))
		return time.Time{}
	}
	timeAsString, ok := (*t).(string)
	if !ok {
		j.setLastError(fmt.Errorf("cannot convert type %T to type time.Time\n", t))
		return time.Time{}
	}
	for _, layout := range timeLayouts {
		parsedTime, err := time.Parse(layout, timeAsString)
		if err == nil {
			return parsedTime
		}
	}
	j.setLastError(fmt.Errorf("the value '%v' could not be converted to type time.Time", timeAsString))
	return time.Time{}
}

func marshal(v any) ([]byte, error) {
	jsonBytes, err := jsonIter.Marshal(v)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func marshalIndent(v any) ([]byte, error) {
	// jsoniter has a bug with indentation
	jsonBytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func unmarshal(data []byte, v any) error {
	return jsonIter.Unmarshal(data, &v)
}

func getGenericMap[T any](f jc[T], o JsonObject) map[string]T {
	genericMap := make(map[string]T)
	for k, v := range o.object {
		genericMap[k] = f(v, &o)
	}
	return genericMap
}

func getGenericArray[T any](f jc[T], o JsonArray) []T {
	arr := make([]T, 0, len(o.elements))
	for _, element := range o.elements {
		arr = append(arr, f(element, &o))
	}
	return arr
}

func getObjectScalar[T any](o *JsonObject, f jc[T], key string) T {
	var t T
	v, ok := o.object[key]
	if !ok {
		o.setLastError(createKeyNotFoundErr(key))
		return t
	}
	if v == nil {
		o.setLastError(createTypeConversionErr(nil, t))
		return t
	}
	return f(v, o)
}

func getArrayScalar[T any](a *JsonArray, f jc[T], i int) T {
	var t T
	if i >= a.Length() {
		a.setLastError(createIndexOutOfRangeErr(i, a.Length()))
		return t
	}
	data := a.elements[i]
	if data == nil {
		a.setLastError(createTypeConversionErr(nil, t))
		return t
	}
	return f(data, a)
}

func convertAnyToString(data *any, j jsonI) string {
	if data == nil {
		j.setLastError(createTypeConversionErr(nil, ""))
		return ""
	}
	switch v := (*data).(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		return strconv.Itoa(v)
	case bool:
		return strconv.FormatBool(v)
	default:
		j.setLastError(createTypeConversionErr(*data, ""))
		return ""
	}
}

func convertAnyToInt(data *any, j jsonI) int {
	if data == nil {
		j.setLastError(createTypeConversionErr(nil, 0))
		return 0
	}
	switch (*data).(type) {
	case float64:
		return int((*data).(float64))
	case int:
		return (*data).(int)
	default:
		j.setLastError(createTypeConversionErr(*data, 0))
		return 0
	}
}

func convertAnyToFloat(data *any, j jsonI) float64 {
	if data == nil {
		j.setLastError(createTypeConversionErr(nil, 0.0))
		return 0
	}
	v, ok := (*data).(float64)
	if !ok {
		j.setLastError(createTypeConversionErr(*data, 0.0))
		return 0
	}
	return v
}

func convertAnyToBool(data *any, j jsonI) bool {
	if data == nil {
		j.setLastError(createTypeConversionErr(nil, false))
		return false
	}
	v, ok := (*data).(bool)
	if !ok {
		j.setLastError(createTypeConversionErr(*data, false))
		return false
	}
	return v
}

func convertAnyToObject(data *any, j jsonI) JsonObject {
	if data == nil {
		j.setLastError(createTypeConversionErr(nil, JsonObject{}))
		return *EmptyObject()
	}
	v, ok := (*data).(map[string]any)
	if !ok {
		j.setLastError(createTypeConversionErr(data, JsonObject{}))
		return *EmptyObject()
	}
	obj := EmptyObject()
	var object = make(map[string]*any)
	for key, value := range v {
		object[key] = &value
	}
	obj.object = object
	return *obj
}

func convertAnyToArray(data *any, j jsonI) JsonArray {
	if data == nil {
		j.setLastError(createTypeConversionErr(nil, JsonArray{}))
		return *EmptyArray()
	}
	v, ok := (*data).([]any)
	if !ok {
		j.setLastError(createTypeConversionErr(data, JsonArray{}))
		return *EmptyArray()
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

//func transformKeys(m map[string]*any) map[string]*any {
//	newMap := make(map[string]*any)
//	for key, value := range m {
//		newKey := toSnakeCase(key)
//		if value == nil {
//			newMap[newKey] = nil
//			continue
//		}
//		nestedMap, ok := (*value).(map[string]any)
//		if ok {
//			nestedResult := transformKeys(convertToMapValuesPtr(nestedMap))
//			var nestedInterface any = nestedResult
//			newMap[newKey] = &nestedInterface
//		} else {
//			newMap[newKey] = value
//		}
//	}
//	return newMap
//}

//func toSnakeCase(str string) string {
//	var result []rune
//	for i, r := range str {
//		if unicode.IsUpper(r) {
//			if i > 0 {
//				result = append(result, '_')
//			}
//			result = append(result, unicode.ToLower(r))
//		} else {
//			result = append(result, r)
//		}
//	}
//	return string(result)
//}
