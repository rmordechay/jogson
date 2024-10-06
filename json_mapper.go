package jsonmapper

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

type JsonMapper struct {
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
	AsObject JsonObject
	AsArray  JsonArray
}

func GetMapper(data string) (JsonMapper, error) {
	var mapper JsonMapper
	if isJsonArray(data) {
		mapper.IsArray = true
		array, err := parseJsonArray(data)
		if err != nil {
			return JsonMapper{}, err
		}
		mapper.AsArray = array
	} else {
		mapper.IsObject = true
		object, err := parseJsonObject(data)
		if err != nil {
			return JsonMapper{}, err
		}
		mapper.AsObject = object
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
		mapper.AsObject = jo
	case []interface{}:
		mapper.IsArray = true
		var ja JsonArray
		ja.elements = data.([]interface{})
		ja.Length = len(ja.elements)
		mapper.AsArray = ja
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
	err := json.Unmarshal([]byte(data), &jo.object)
	if err != nil {
		return JsonObject{}, err
	}
	return jo, nil
}

func parseJsonArray(data string) (JsonArray, error) {
	var ja JsonArray
	var arr []interface{}
	err := json.Unmarshal([]byte(data), &arr)
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
		return fmt.Sprintf("%v", j.AsObject)
	} else if j.IsArray {
		return fmt.Sprintf("%v", j.AsArray)
	}
	return ""
}
