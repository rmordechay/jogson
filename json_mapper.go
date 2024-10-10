package jsonmapper

import (
	"encoding/json"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"os"
	"time"
	"unicode"
)

var jsonIter = jsoniter.ConfigCompatibleWithStandardLibrary

// Json represents a generic JSON type.
// It contains fields for all supported JSON types like bool, int, float, string, object, and array,
// as well as Go supported types.
type Json struct {
	IsBool   bool
	IsInt    bool
	IsFloat  bool
	IsString bool
	IsObject bool
	IsArray  bool
	IsNull   bool

	AsBool   bool
	AsInt    int
	AsFloat  float64
	AsString string
	Object   JsonObject
	Array    JsonArray
}

// FromBytes parses JSON data from a byte slice.
// It automatically determines whether the input is a JSON object or array.
func FromBytes(data []byte) (Json, error) {
	if isObjectOrArray(data, '[') {
		return newJsonArray(data)
	} else if isObjectOrArray(data, '{') {
		return newJsonObject(data)
	} else {
		return Json{}, errors.New("could not parse JSON")
	}
}

// FromStruct serializes a Go struct into JSON and parses it into a Json object.
func FromStruct[T any](s T) (Json, error) {
	jsonBytes, err := marshal(s)
	if err != nil {
		return Json{}, err
	}
	return FromBytes(jsonBytes)
}

// FromFile reads a JSON file from the given path and parses it into a Json object.
func FromFile(path string) (Json, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return Json{}, err
	}
	return FromBytes(file)
}

// FromString parses JSON from a string into a Json object.
func FromString(data string) (Json, error) {
	return FromBytes([]byte(data))
}

// AsTime attempts to convert the JSON value to a time.Time object.
// Only works if the JSON value is a string and can be parsed as a valid time.
func (m *Json) AsTime() (time.Time, error) {
	if !m.IsString {
		return time.Time{}, fmt.Errorf("cannot convert type %v to type time.Time\n", m.getType())
	}
	for _, layout := range timeLayouts {
		parsedTime, err := time.Parse(layout, m.AsString)
		if err == nil {
			return parsedTime, nil
		}
	}
	return time.Time{}, fmt.Errorf("the value '%v' could not be converted to type time.Time", m.AsString)
}

// PrettyString returns a formatted, human-readable string representation of the Json value.
func (m *Json) PrettyString() string {
	if m.IsBool {
		return fmt.Sprintf("%v", m.AsBool)
	} else if m.IsInt {
		return fmt.Sprintf("%v", m.AsInt)
	} else if m.IsFloat {
		return fmt.Sprintf("%v", m.AsFloat)
	} else if m.IsString {
		return fmt.Sprintf("%v", m.AsString)
	} else if m.IsObject {
		return m.Object.PrettyString()
	} else if m.IsArray {
		return fmt.Sprintf("%v", m.Array)
	}
	return ""
}

// String returns a string representation Json type in JSON format.
func (m *Json) String() string {
	switch {
	case m.IsBool:
		return fmt.Sprintf("%v", m.AsBool)
	case m.IsInt:
		return fmt.Sprintf("%v", m.AsInt)
	case m.IsFloat:
		return fmt.Sprintf("%v", m.AsFloat)
	case m.IsString:
		return fmt.Sprintf("%v", m.AsString)
	case m.IsObject:
		return fmt.Sprintf("%v", m.Object)
	case m.IsArray:
		return fmt.Sprintf("%v", m.Array)
	}
	return ""
}

func (m *Json) getType() JsonType {
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

func getMapperFromField(data *any) Json {
	var mapper Json
	if data == nil {
		return Json{IsNull: true}
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
		mapper.Object = getAsJsonObject(&value, nil)
	case []float64:
		mapper.IsArray = true
		mapper.Array = getAsJsonArray(value.([]float64))
	case []int:
		mapper.IsArray = true
		mapper.Array = getAsJsonArray(value.([]int))
	case []string:
		mapper.IsArray = true
		mapper.Array = getAsJsonArray(value.([]string))
	case []bool:
		mapper.IsArray = true
		mapper.Array = getAsJsonArray(value.([]bool))
	case []*any:
		mapper.IsArray = true
		mapper.Array = JsonArray{elements: value.([]*any)}
	case []any:
		mapper.IsArray = true
		mapper.Array = JsonArray{elements: convertToSlicePtr(value.([]any))}
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
		array[i] = &v
	}
	return array
}

func convertToMapValuesPtr(data map[string]any) map[string]*any {
	jsonObject := make(map[string]*any, len(data))
	for k, v := range data {
		jsonObject[k] = &v
	}
	return jsonObject
}

func newJsonArray(data []byte) (Json, error) {
	var mapper Json
	mapper.IsArray = true
	array, err := parseJsonArray(data)
	if err != nil {
		return Json{}, err
	}
	mapper.Array = array
	return mapper, nil
}

func newJsonObject(data []byte) (Json, error) {
	var mapper Json
	mapper.IsObject = true
	object, err := parseJsonObject(data)
	if err != nil {
		return Json{}, err
	}
	mapper.Object = object
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
		return time.Time{}, fmt.Errorf(nullConversionErrStr, "string")
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
