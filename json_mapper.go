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

type Json interface {
	AsBool() bool
	AsInt() int
	AsFloat() float64
	AsString() string
	Object() *JsonObject
	Array() *JsonArray
	IsBool() bool
	IsInt() bool
	IsFloat() bool
	IsString() bool
	IsObject() bool
	IsArray() bool
	IsNull() bool
}

type mapper struct {
	isBool   bool
	isInt    bool
	isFloat  bool
	isString bool
	isObject bool
	isArray  bool
	isNull   bool

	asBool   bool
	asInt    int
	asFloat  float64
	asString string
	object   JsonObject
	array    JsonArray
}

func (m *mapper) IsBool() bool {
	return m.isBool
}

func (m *mapper) IsInt() bool {
	return m.isInt
}

func (m *mapper) IsFloat() bool {
	return m.isFloat
}

func (m *mapper) IsString() bool {
	return m.isString
}

func (m *mapper) IsObject() bool {
	return m.isObject
}

func (m *mapper) IsArray() bool {
	return m.isArray
}

func (m *mapper) IsNull() bool {
	return m.isNull
}

func (m *mapper) AsBool() bool {
	return m.asBool
}

func (m *mapper) AsInt() int {
	return m.asInt
}

func (m *mapper) AsFloat() float64 {
	return m.asFloat
}

func (m *mapper) AsString() string {
	return m.asString
}

func (m *mapper) Object() *JsonObject {
	return &m.object
}

func (m *mapper) Array() *JsonArray {
	return &m.array
}

func FromBytes(data []byte) (Json, error) {
	var jsonMapper mapper
	if isObjectOrArray(data, '[') {
		jsonMapper.isArray = true
		array, err := parseJsonArray(data)
		if err != nil {
			return &mapper{}, err
		}
		jsonMapper.array = array
	} else if isObjectOrArray(data, '{') {
		jsonMapper.isObject = true
		object, err := parseJsonObject(data)
		if err != nil {
			return &mapper{}, err
		}
		jsonMapper.object = object
	} else {
		return &mapper{}, errors.New("could not parse JSON")
	}
	return &jsonMapper, nil
}

func FromStruct[T any](s T) (Json, error) {
	jsonBytes, err := marshal(s)
	if err != nil {
		return &mapper{}, err
	}
	return FromBytes(jsonBytes)
}

func FromFile(path string) (Json, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return &mapper{}, err
	}
	return FromBytes(file)
}

func FromString(data string) (Json, error) {
	return FromBytes([]byte(data))
}

func (m *mapper) AsTime() (time.Time, error) {
	if !m.isString {
		return time.Time{}, fmt.Errorf("cannot convert type %v to type time.Time\n", m.getType())
	}
	for _, layout := range timeLayouts {
		parsedTime, err := time.Parse(layout, m.asString)
		if err == nil {
			return parsedTime, nil
		}
	}
	return time.Time{}, fmt.Errorf("the value '%v' could not be converted to type time.Time", m.asString)
}

func (m *mapper) String() string {
	switch {
	case m.isBool:
		return fmt.Sprintf("%v", m.asBool)
	case m.isInt:
		return fmt.Sprintf("%v", m.asInt)
	case m.isFloat:
		return fmt.Sprintf("%v", m.asFloat)
	case m.isString:
		return fmt.Sprintf("%v", m.asString)
	case m.isObject:
		return fmt.Sprintf("%v", m.object)
	case m.isArray:
		return fmt.Sprintf("%v", m.array)
	}
	return ""
}

func (m *mapper) PrettyString() string {
	if m.isBool {
		return fmt.Sprintf("%v", m.asBool)
	} else if m.isInt {
		return fmt.Sprintf("%v", m.asInt)
	} else if m.isFloat {
		return fmt.Sprintf("%v", m.asFloat)
	} else if m.isString {
		return fmt.Sprintf("%v", m.asString)
	} else if m.isObject {
		return m.object.PrettyString()
	} else if m.isArray {
		return fmt.Sprintf("%v", m.array)
	}
	return ""
}

func (m *mapper) getType() JsonType {
	switch {
	case m.isBool:
		return Bool
	case m.isInt:
		return Int
	case m.isFloat:
		return Float
	case m.isString:
		return String
	case m.isObject:
		return Object
	case m.isNull:
		return Null
	case m.isArray:
		return Array
	}
	return Invalid
}

func getMapperFromField(data *interface{}) Json {
	var jsonMapper mapper
	if data == nil {
		return &mapper{isNull: true}
	}
	value := *data
	switch value.(type) {
	case bool:
		jsonMapper.isBool = true
		jsonMapper.asBool = value.(bool)
	case int:
		jsonMapper.isInt = true
		jsonMapper.asInt = value.(int)
	case float64:
		if value == float64(int(value.(float64))) {
			jsonMapper.isInt = true
			jsonMapper.asInt = int(value.(float64))
		} else {
			jsonMapper.isFloat = true
		}
		jsonMapper.asFloat = value.(float64)
	case string:
		jsonMapper.isString = true
		jsonMapper.asString = value.(string)
	case map[string]interface{}:
		jsonMapper.isObject = true
		jsonMapper.object = getAsJsonObject(value, nil)
	case []float64:
		jsonMapper.isArray = true
		jsonMapper.array = getAsJsonArray(value.([]float64))
	case []int:
		jsonMapper.isArray = true
		jsonMapper.array = getAsJsonArray(value.([]int))
	case []string:
		jsonMapper.isArray = true
		jsonMapper.array = getAsJsonArray(value.([]string))
	case []bool:
		jsonMapper.isArray = true
		jsonMapper.array = getAsJsonArray(value.([]bool))
	case []*interface{}:
		jsonMapper.isArray = true
		jsonMapper.array = JsonArray{elements: value.([]*interface{})}
	case []interface{}:
		jsonMapper.isArray = true
		jsonMapper.array = JsonArray{elements: convertToSlicePtr(value.([]interface{}))}
	case nil:
		jsonMapper.isNull = true
	default:
		panic(fmt.Errorf("JSON conversion for %v failed. %T not implemented", value, data))
	}
	return &jsonMapper
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
	var arr []*interface{}
	err := unmarshal(data, &arr)
	if err != nil {
		return JsonArray{}, err
	}
	ja.elements = arr
	return ja, nil
}

func convertToSlicePtr(data []interface{}) []*interface{} {
	array := make([]*interface{}, len(data))
	for i, v := range data {
		array[i] = &v
	}
	return array
}

func convertToMapValuesPtr(data map[string]interface{}) map[string]*interface{} {
	jsonObject := make(map[string]*interface{}, len(data))
	for k, v := range data {
		jsonObject[k] = &v
	}
	return jsonObject
}

func getAsJsonObject(data interface{}, j JsonError) JsonObject {
	v, ok := data.(map[string]interface{})
	if !ok {
		j.SetLastError(fmt.Errorf(TypeConversionErrStr, data, JsonObject{}))
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

func getAsJsonArray[T any](data []T) JsonArray {
	var arr JsonArray
	array := make([]*interface{}, len(data))
	for i, v := range data {
		var valAny interface{} = v
		array[i] = &valAny
	}
	arr.elements = array
	return arr
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

func marshal(v interface{}) ([]byte, error) {
	jsonBytes, err := jsonIter.Marshal(v)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func marshalIndent(v interface{}) ([]byte, error) {
	// jsoniter has a bug with indentation
	jsonBytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func unmarshal(data []byte, v interface{}) error {
	return jsonIter.Unmarshal(data, &v)
}
