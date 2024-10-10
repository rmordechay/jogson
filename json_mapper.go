package jsonmapper

import (
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"os"
	"time"
)

var jsonIter = jsoniter.ConfigCompatibleWithStandardLibrary

// JsonMapper represents a generic JSON type. It contains fields for all supported JSON
// types like bool, int, float, string, object, and array, as well as Go supported types.
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

// FromBytes parses JSON data from a byte slice.
func FromBytes(data []byte) (JsonMapper, error) {
	if isObjectOrArray(data, '[') {
		return newJsonArray(data)
	} else if isObjectOrArray(data, '{') {
		return newJsonObject(data)
	} else {
		return JsonMapper{}, errors.New("could not parse JSON")
	}
}

// FromStruct serializes a Go struct into JSON and parses it into a JsonMapper object.
func FromStruct[T any](s T) (JsonMapper, error) {
	jsonBytes, err := marshal(s)
	if err != nil {
		return JsonMapper{}, err
	}
	return FromBytes(jsonBytes)
}

// FromFile reads a JSON file from the given path and parses it into a JsonMapper object.
func FromFile(path string) (JsonMapper, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return JsonMapper{}, err
	}
	return FromBytes(file)
}

// FromString parses JSON from a string into a JsonMapper object.
func FromString(data string) (JsonMapper, error) {
	return FromBytes([]byte(data))
}

// AsTime attempts to convert the JSON value to a time.Time object.
// Only works if the JSON value is a string and can be parsed as a valid time.
func (m *JsonMapper) AsTime() (time.Time, error) {
	if !m.IsString {
		return time.Time{}, NewTimeTypeConversionErr(m.getType())
	}
	for _, layout := range timeLayouts {
		parsedTime, err := time.Parse(layout, m.AsString)
		if err == nil {
			return parsedTime, nil
		}
	}
	return time.Time{}, NewInvalidTimeErr(m.AsString)
}

// PrettyString returns a formatted, human-readable string representation of the JsonMapper value.
func (m *JsonMapper) PrettyString() string {
	if m.IsBool {
		return fmt.Sprintf("%v", m.AsBool)
	} else if m.IsInt {
		return fmt.Sprintf("%v", m.AsInt)
	} else if m.IsFloat {
		return fmt.Sprintf("%v", m.AsFloat)
	} else if m.IsString {
		return fmt.Sprintf("%v", m.AsString)
	} else if m.IsObject {
		return m.AsObject.PrettyString()
	} else if m.IsArray {
		return fmt.Sprintf("%v", m.AsArray)
	}
	return ""
}

// String returns a string representation JsonMapper type in JSON format.
func (m *JsonMapper) String() string {
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
		return fmt.Sprintf("%v", m.AsObject)
	case m.IsArray:
		return fmt.Sprintf("%v", m.AsArray)
	}
	return ""
}
