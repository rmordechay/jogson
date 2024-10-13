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

// jcn is JSON converter function type that convert type any to generic type *T
type jcn[T any] func(data *any, j jsonI) *T

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

func getGenericMap[T any](f jc[T], o JsonObject) map[string]T {
	o.setLastError(nil)
	genericMap := make(map[string]T)
	for k, v := range o.object {
		if v == nil {
			continue
		}
		genericMap[k] = f(v, &o)
	}
	return genericMap
}

func getGenericMapN[T any](f jcn[T], o JsonObject) map[string]*T {
	o.setLastError(nil)
	genericMap := make(map[string]*T)
	for k, v := range o.object {
		if v == nil {
			genericMap[k] = nil
		} else {
			genericMap[k] = f(v, &o)
		}
	}
	return genericMap
}

func getGenericArray[T any](f jc[T], a JsonArray) []T {
	a.setLastError(nil)
	arr := make([]T, 0, len(a.elements))
	for _, v := range a.elements {
		if v == nil {
			continue
		}
		arr = append(arr, f(v, &a))
	}
	return arr
}

func getGenericArrayN[T any](f jcn[T], a JsonArray) []*T {
	a.setLastError(nil)
	arr := make([]*T, 0, len(a.elements))
	for _, v := range a.elements {
		if v == nil {
			arr = append(arr, nil)
		} else {
			arr = append(arr, f(v, &a))
		}
	}
	return arr
}

func getObjectScalar[T any](o *JsonObject, f jc[T], key string) T {
	o.setLastError(nil)
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

func getObjectScalarNullable[T any](o *JsonObject, f jcn[T], key string) *T {
	o.setLastError(nil)
	var t T
	v, ok := o.object[key]
	if !ok {
		o.setLastError(createKeyNotFoundErr(key))
		return nil
	}
	if v == nil {
		o.setLastError(createTypeConversionErr(nil, t))
		return nil
	}
	return f(v, o)
}

func getArrayScalar[T any](a *JsonArray, f jc[T], i int) T {
	a.setLastError(nil)
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

func getArrayScalarNullable[T any](a *JsonArray, f jcn[T], i int) *T {
	a.setLastError(nil)
	var t T
	if i >= a.Length() {
		a.setLastError(createIndexOutOfRangeErr(i, a.Length()))
		return nil
	}
	data := a.elements[i]
	if data == nil {
		a.setLastError(createTypeConversionErr(nil, t))
		return nil
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

func convertAnyToStringN(data *any, j jsonI) *string {
	if data == nil {
		j.setLastError(createTypeConversionErr(nil, ""))
		return nil
	}
	switch v := (*data).(type) {
	case string:
		return &v
	case float64:
		f := strconv.FormatFloat(v, 'f', -1, 64)
		return &f
	case int:
		i := strconv.Itoa(v)
		return &i
	case bool:
		f := strconv.FormatBool(v)
		return &f
	default:
		j.setLastError(createTypeConversionErr(*data, ""))
		return nil
	}
}

func convertAnyToInt(data *any, j jsonI) int {
	if data == nil {
		j.setLastError(createTypeConversionErr(nil, 0))
		return 0
	}
	switch v := (*data).(type) {
	case float64:
		return int(v)
	case int:
		return v
	default:
		j.setLastError(createTypeConversionErr(*data, 0))
		return 0
	}
}

func convertAnyToIntN(data *any, j jsonI) *int {
	if data == nil {
		j.setLastError(createTypeConversionErr(nil, 0))
		return nil
	}
	switch v := (*data).(type) {
	case float64:
		i := int(v)
		return &i
	case int:
		return &v
	default:
		j.setLastError(createTypeConversionErr(*data, 0))
		return nil
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

func convertAnyToFloatN(data *any, j jsonI) *float64 {
	if data == nil {
		j.setLastError(createTypeConversionErr(nil, 0.0))
		return nil
	}
	v, ok := (*data).(float64)
	if !ok {
		j.setLastError(createTypeConversionErr(*data, 0.0))
		return nil
	}
	return &v
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

func convertAnyToBoolN(data *any, j jsonI) *bool {
	if data == nil {
		j.setLastError(createTypeConversionErr(nil, false))
		return nil
	}
	v, ok := (*data).(bool)
	if !ok {
		j.setLastError(createTypeConversionErr(*data, false))
		return nil
	}
	return &v
}

func convertAnyToObject(data *any, j jsonI) JsonObject {
	if data == nil {
		j.setLastError(createTypeConversionErr(nil, JsonObject{}))
		return *nullObject()
	}
	v, ok := (*data).(map[string]any)
	if !ok {
		j.setLastError(createTypeConversionErr(data, JsonObject{}))
		return *nullObject()
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
		return *nullArray()
	}
	v, ok := (*data).([]any)
	if !ok {
		j.setLastError(createTypeConversionErr(data, JsonArray{}))
		return *nullArray()
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

func parseTime(t *any, j jsonI) time.Time {
	if t == nil {
		j.setLastError(createTypeConversionErr(nil, ""))
		return time.Time{}
	}
	timeAsString, ok := (*t).(string)
	if !ok {
		j.setLastError(createTypeConversionErr(t, time.Time{}))
		return time.Time{}
	}
	for _, layout := range timeLayouts {
		parsedTime, err := time.Parse(layout, timeAsString)
		if err == nil {
			return parsedTime
		}
	}
	j.setLastError(fmt.Errorf("'%v' could not be parsed as time.Time", timeAsString))
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
