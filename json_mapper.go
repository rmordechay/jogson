package jsonmapper

import (
	"encoding/json"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"log"
	"os"
	"reflect"
	"time"
	"unicode"
)

var jsonIter = jsoniter.ConfigCompatibleWithStandardLibrary

type Mapper struct {
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

	Err error
}

func FromBytes(data []byte) (Mapper, error) {
	var mapper Mapper
	if isArray(data) {
		mapper.IsArray = true
		array, err := parseJsonArray(data)
		if err != nil {
			return Mapper{}, err
		}
		mapper.Array = array
	} else if isObject(data) {
		mapper.IsObject = true
		object, err := parseJsonObject(data)
		if err != nil {
			return Mapper{}, err
		}
		mapper.Object = object
	} else {
		return Mapper{}, errors.New("could not parse JSON")
	}
	return mapper, nil
}

func FromStruct[T any](s T) (Mapper, error) {
	jsonBytes, err := marshal(s)
	if err != nil {
		return Mapper{}, err
	}
	return FromBytes(jsonBytes)
}

func FromFile(path string) (Mapper, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return Mapper{}, err
	}
	return FromBytes(file)
}

func FromString(data string) (Mapper, error) {
	return FromBytes([]byte(data))
}

func (m Mapper) AsTime() (time.Time, error) {
	if !m.IsString {
		return time.Time{}, fmt.Errorf("cannot convert type %v to type time.Time\n", m.GetType())
	}
	for _, layout := range timeLayouts {
		parsedTime, err := time.Parse(layout, m.AsString)
		if err == nil {
			return parsedTime, nil
		}
	}
	return time.Time{}, fmt.Errorf("the value '%v' could not be converted to type time.Time", m.AsString)
}

func (m Mapper) GetType() JsonType {
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

func (m Mapper) String() string {
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

func (m Mapper) PrettyString() string {
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

func createArray(data interface{}) JsonArray {
	var arr JsonArray
	switch data.(type) {
	case []*interface{}:
		arr.elements = data.([]*interface{})
	case []interface{}:
		array := convertToArrayPtr(data)
		arr.elements = array
	}
	return arr
}

func getMapperFromField(data *interface{}) Mapper {
	var mapper Mapper
	if data == nil {
		return Mapper{IsNull: true}
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
	case map[string]interface{}:
		mapper.IsObject = true
		mapper.Object = createJsonObject(value)
	case []float64:
		mapper.IsArray = true
		mapper.Array = convertArray(value.([]float64))
	case []int:
		mapper.IsArray = true
		mapper.Array = convertArray(value.([]int))
	case []string:
		mapper.IsArray = true
		mapper.Array = convertArray(value.([]string))
	case []bool:
		mapper.IsArray = true
		mapper.Array = convertArray(value.([]bool))
	case []interface{}:
		mapper.IsArray = true
		mapper.Array = createArray(value)
	case nil:
		mapper.IsNull = true
	default:
		log.Fatalf("JSON conversion for %v failed. %v not implemented.", value, reflect.TypeOf(data))
	}
	return mapper
}

func createJsonObject(data interface{}) JsonObject {
	var obj JsonObject
	var object = make(map[string]*interface{})
	for k, v := range data.(map[string]interface{}) {
		object[k] = &v
	}
	obj.object = object
	return obj
}

func parseJsonObject(data []byte) (JsonObject, error) {
	var jo JsonObject
	err := unmarshal(data, &jo.object)
	if err != nil {
		return JsonObject{}, err
	}
	return jo, nil
}

func convertToArrayPtr(data interface{}) []*interface{} {
	d := data.([]interface{})
	array := make([]*interface{}, len(d))
	for i, v := range d {
		array[i] = &v
	}
	return array
}

func isArray(data []byte) bool {
	return isObjectOrArray(data, '[')
}

func isObject(data []byte) bool {
	return isObjectOrArray(data, '{')
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
