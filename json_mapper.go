package jsonmapper

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"
)

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

func GetMapperFromFile(path string) (Mapper, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return Mapper{}, err
	}
	return GetMapperFromString(string(file))
}

func GetMapperFromBytes(data []byte) (Mapper, error) {
	return GetMapperFromString(string(data))
}

func GetMapperFromString(data string) (Mapper, error) {
	var mapper Mapper
	if isJsonArray(data) {
		mapper.IsArray = true
		array, err := parseJsonArray(data)
		if err != nil {
			return Mapper{}, err
		}
		mapper.Array = array
	} else {
		mapper.IsObject = true
		object, err := parseJsonObject(data)
		if err != nil {
			return Mapper{}, err
		}
		mapper.Object = object
	}
	return mapper, nil
}

func (m Mapper) AsTime() time.Time {
	if !m.IsString {
		m.Err = fmt.Errorf("could not convert current type %v to time", m.String())
		return time.Time{}
	}
	for _, layout := range timeLayouts {
		parsedTime, err := time.Parse(layout, m.AsString)
		if err == nil {
			return parsedTime
		}
	}
	return time.Time{}
}

func (m Mapper) String() string {
	if m.IsBool {
		return fmt.Sprintf("%v", m.AsBool)
	} else if m.IsInt {
		return fmt.Sprintf("%v", m.AsInt)
	} else if m.IsFloat {
		return fmt.Sprintf("%v", m.AsFloat)
	} else if m.IsString {
		return fmt.Sprintf("%v", m.AsString)
	} else if m.IsObject {
		return fmt.Sprintf("%v", m.Object)
	} else if m.IsArray {
		return fmt.Sprintf("%v", m.Array)
	}
	return ""
}

func getMapperFromField(data *interface{}) Mapper {
	var mapper Mapper
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
		mapper.Array = CreateJsonArray(value)
	case map[string]interface{}:
		mapper.IsObject = true
		mapper.Object = createJsonObject(value)
	case nil:
		mapper.IsNull = true
	default:
		log.Fatalf("JSON conversion for %v failed. %v not implemented.", value, reflect.TypeOf(data))
	}
	return mapper
}

func isJsonArray(data string) bool {
	return data[0] == '['
}

func marshal(v any) []byte {
	jsonBytes, _ := json.Marshal(v)
	return jsonBytes
}

func unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, &v)
}
