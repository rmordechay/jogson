package jsonmapper

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
)

type JsonType interface {
	int | string | float64 | bool | map[string]interface{} | []interface{}
}

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

func (j Mapper) String() string {
	if j.IsBool {
		return fmt.Sprintf("%v", j.AsBool)
	} else if j.IsInt {
		return fmt.Sprintf("%v", j.AsInt)
	} else if j.IsFloat {
		return fmt.Sprintf("%v", j.AsFloat)
	} else if j.IsString {
		return fmt.Sprintf("%v", j.AsString)
	} else if j.IsObject {
		return fmt.Sprintf("%v", j.Object)
	} else if j.IsArray {
		return fmt.Sprintf("%v", j.Array)
	}
	return ""
}

func getMapperFromField(data interface{}) Mapper {
	var mapper Mapper
	switch data.(type) {
	case bool:
		mapper.IsBool = true
		mapper.AsBool = data.(bool)
	case int:
		mapper.IsInt = true
		mapper.AsInt = data.(int)
	case float64:
		if data == float64(int(data.(float64))) {
			mapper.IsInt = true
			mapper.AsInt = int(data.(float64))
		} else {
			mapper.IsFloat = true
		}
		mapper.AsFloat = data.(float64)
	case string:
		mapper.IsString = true
		mapper.AsString = data.(string)
	case []float64:
		mapper.IsArray = true
		mapper.Array = convertArray(data.([]float64))
	case []int:
		mapper.IsArray = true
		mapper.Array = convertArray(data.([]int))
	case []string:
		mapper.IsArray = true
		mapper.Array = convertArray(data.([]string))
	case []bool:
		mapper.IsArray = true
		mapper.Array = convertArray(data.([]bool))
	case []interface{}:
		mapper.IsArray = true
		mapper.Array = CreateJsonArray(data)
	case map[string]interface{}:
		mapper.IsObject = true
		mapper.Object = CreateJsonObject(data)
	case nil:
		mapper.IsNull = true
	default:
		log.Fatalf("getMapperFromField failed. %v not implemented.", reflect.TypeOf(data))
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
