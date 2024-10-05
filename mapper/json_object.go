package mapper

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

type JsonType int

const (
	Bool JsonType = iota
	Int
	Float
	String
	Object
	Array
	Null
)

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
	AsObject JsonObject
	AsArray  JsonArray
}

func (j Json) getType() JsonType {
	if j.IsBool {
		return Bool
	} else if j.IsInt {
		return Int
	} else if j.IsFloat {
		return Float
	} else if j.IsString {
		return String
	} else if j.IsObject {
		return Object
	} else if j.IsArray {
		return Array
	}
	return Null
}

func (j Json) String() string {
	switch j.getType() {
	case Bool:
		return fmt.Sprintf("%v", j.AsBool)
	case Int:
		return fmt.Sprintf("%v", j.AsInt)
	case Float:
		return fmt.Sprintf("%v", j.AsFloat)
	case String:
		return fmt.Sprintf("%v", j.AsString)
	case Object:
		return fmt.Sprintf("%v", j.AsObject)
	case Array:
		return fmt.Sprintf("%v", j.AsArray)
	case Null:
		return fmt.Sprintf("%v", nil)
	}
	return ""
}

type JsonArray struct {
	elements []interface{}
}

type JsonObject struct {
	object map[string]interface{}
}

func (o JsonObject) Get(key string) Json {
	for k, v := range o.object {
		if k == key {
			return getMapperFromField(v)
		}
	}
	return Json{}
}

func GetMapper(data string) (Json, error) {
	var mapper Json
	if isJsonArray(data) {
		mapper.IsArray = true
		array, err := parseJsonArray(data)
		if err != nil {
			return Json{}, err
		}
		mapper.AsArray = array
	} else {
		mapper.IsObject = true
		object, err := parseJsonObject(data)
		if err != nil {
			return Json{}, err
		}
		mapper.AsObject = object
	}
	return mapper, nil
}

func getMapperFromField(data interface{}) Json {
	var mapper Json
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

func (a JsonArray) Elements() []Json {
	jsons := make([]Json, 0, len(a.elements))
	for _, element := range a.elements {
		jsons = append(jsons, getMapperFromField(element))
	}
	return jsons
}

func (o JsonObject) Elements() map[string]Json {
	jsons := make(map[string]Json)
	for k, v := range o.object {
		jsons[k] = getMapperFromField(v)
	}
	return jsons
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
