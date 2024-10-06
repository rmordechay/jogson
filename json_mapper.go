package jsonmapper

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
)

type JsonMapper struct {
	writer   JsonWriter
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

func GetMapperFromFile(path string) (JsonMapper, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return JsonMapper{}, err
	}
	return GetMapperFromString(string(file))
}

func GetMapperFromBytes(data []byte) (JsonMapper, error) {
	return GetMapperFromString(string(data))
}

func GetMapperFromString(data string) (JsonMapper, error) {
	var mapper JsonMapper
	if isJsonArray(data) {
		mapper.IsArray = true
		array, err := parseJsonArray(data)
		if err != nil {
			return JsonMapper{}, err
		}
		mapper.Array = array
	} else {
		mapper.IsObject = true
		object, err := parseJsonObject(data)
		if err != nil {
			return JsonMapper{}, err
		}
		mapper.Object = object
	}
	return mapper, nil
}

func getMapperFromField(data interface{}) JsonMapper {
	var mapper JsonMapper
	switch data.(type) {
	case bool:
		mapper.IsBool = true
		mapper.AsBool = data.(bool)
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
	case map[string]interface{}:
		mapper.IsObject = true
		var jo JsonObject
		jo.object = data.(map[string]interface{})
		mapper.Object = jo
	case []interface{}:
		mapper.IsArray = true
		var ja JsonArray
		ja.elements = data.([]interface{})
		ja.Length = len(ja.elements)
		mapper.Array = ja
	case nil:
		mapper.IsNull = true
	default:
		log.Fatalf("%v not implemented", reflect.TypeOf(data))
	}
	return mapper
}

func isJsonArray(data string) bool {
	return data[0] == '['
}

func parseJsonObject(data string) (JsonObject, error) {
	var jo JsonObject
	err := unmarshal([]byte(data), &jo.object)
	if err != nil {
		return JsonObject{}, err
	}
	return jo, nil
}

func parseJsonArray(data string) (JsonArray, error) {
	var ja JsonArray
	var arr []interface{}
	err := unmarshal([]byte(data), &arr)
	if err != nil {
		return JsonArray{}, err
	}
	ja.elements = arr
	return ja, nil
}

func (j JsonMapper) String() string {
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

func marshal(v any) []byte {
	jsonBytes, _ := json.Marshal(v)
	return jsonBytes
}

func unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, &v)
}
