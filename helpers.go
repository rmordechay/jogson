package jsonmapper

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"unicode"
)

func (m *JsonMapper) getType() JsonType {
	switch {
	case m.IsBool:
		return Bool
	case m.IsInt:
		return Int
	case m.IsFloat:
		return Float
	case m.IsString:
		return String
	case m.IsObject:
		return Object
	case m.IsNull:
		return Null
	case m.IsArray:
		return Array
	}
	return Invalid
}

func getMapperFromField(data *any) JsonMapper {
	var mapper JsonMapper
	if data == nil {
		return JsonMapper{IsNull: true}
	}
	value := *data
	switch value.(type) {
	case bool:
		mapper.IsBool = true
		mapper.AsBool = value.(bool)
	case int:
		mapper.IsInt = true
		mapper.AsInt = value.(int)
	case float64:
		if value == float64(int(value.(float64))) {
			mapper.IsInt = true
			mapper.AsInt = int(value.(float64))
		} else {
			mapper.IsFloat = true
		}
		mapper.AsFloat = value.(float64)
	case string:
		mapper.IsString = true
		mapper.AsString = value.(string)
	case map[string]any:
		mapper.IsObject = true
		mapper.AsObject = convertAnyToObject(&value, nil)
	case []float64:
		mapper.IsArray = true
		mapper.AsArray = convertSliceToJsonArray(value.([]float64))
	case []int:
		mapper.IsArray = true
		mapper.AsArray = convertSliceToJsonArray(value.([]int))
	case []string:
		mapper.IsArray = true
		mapper.AsArray = convertSliceToJsonArray(value.([]string))
	case []bool:
		mapper.IsArray = true
		mapper.AsArray = convertSliceToJsonArray(value.([]bool))
	case []*any:
		mapper.IsArray = true
		mapper.AsArray = *NewArray(value.([]*any))
	case []any:
		mapper.IsArray = true
		mapper.AsArray = *NewArray(convertToSlicePtr(value.([]any)))
	case nil:
		mapper.IsNull = true
	default:
		panic(fmt.Errorf("JSON conversion for %v failed. %T not implemented", value, data))
	}
	return mapper
}

func parseJsonObject(data []byte) (JsonObject, error) {
	var jo JsonObject
	err := unmarshal(data, &jo.object)
	if err != nil {
		return JsonObject{}, err
	}
	return jo, nil
}

func parseJsonArray(data []byte) (JsonArray, error) {
	var ja JsonArray
	var arr []*any
	err := unmarshal(data, &arr)
	if err != nil {
		return JsonArray{}, err
	}
	ja.elements = arr
	return ja, nil
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

func newJsonArray(data []byte) (JsonMapper, error) {
	var mapper JsonMapper
	mapper.IsArray = true
	array, err := parseJsonArray(data)
	if err != nil {
		return JsonMapper{}, err
	}
	mapper.AsArray = array
	return mapper, nil
}

func newJsonObject(data []byte) (JsonMapper, error) {
	var mapper JsonMapper
	mapper.IsObject = true
	object, err := parseJsonObject(data)
	if err != nil {
		return JsonMapper{}, err
	}
	mapper.AsObject = object
	return mapper, nil
}

func isObjectOrArray(data []byte, brackOrParen byte) bool {
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

func parseTime(t *any) (time.Time, error) {
	if t == nil {
		return time.Time{}, NewNullConversionErr("string")
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

func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

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
