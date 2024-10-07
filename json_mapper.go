package jsonmapper

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
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

func FromBytes(data []byte) (Mapper, error) {
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

func EmptyObject() JsonObject {
	var obj JsonObject
	obj.object = make(map[string]*interface{})
	return obj
}

func EmptyArray() JsonArray {
	var arr JsonArray
	elements := make([]*interface{}, 0)
	arr.elements = elements
	return arr
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
		mapper.Array = createArray(value)
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

func convertToArrayPtr(data interface{}) []*interface{} {
	d := data.([]interface{})
	array := make([]*interface{}, len(d))
	for i, v := range d {
		array[i] = &v
	}
	return array
}

func isJsonArray(data []byte) bool {
	return data[0] == '['
}

func marshal(v interface{}) ([]byte, error) {
	var jsonIter = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonBytes, err := jsonIter.Marshal(v)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func unmarshal(data []byte, v interface{}) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Unmarshal(data, &v)
}

func iter(data []byte) {
	var jsonIter = jsoniter.ConfigCompatibleWithStandardLibrary
	iterator := jsonIter.BorrowIterator(data)
	defer jsonIter.ReturnIterator(iterator)

	fmt.Printf("%v\n", iterator.WhatIsNext())
	iterator.ReadArray()
	fmt.Printf("%v\n", iterator.WhatIsNext())
	fmt.Printf("%v\n", iterator.ReadObject())
}
