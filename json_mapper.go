package jsonmapper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

const bufferSize = 4096

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

	buffer   []byte
	offset   int
	lastRead int
	reader   io.Reader
}

// FromBytes parses JSON data from a byte slice.
func FromBytes(data []byte) (JsonMapper, error) {
	if dataStartsWith(data, '[') {
		arrayBytes, err := NewArrayFromBytes(data)
		if err != nil {
			return JsonMapper{}, err
		}
		return JsonMapper{IsArray: true, AsArray: *arrayBytes}, nil
	}

	if dataStartsWith(data, '{') {
		objBytes, err := NewObjectFromBytes(data)
		if err != nil {
			return JsonMapper{}, err
		}
		return JsonMapper{IsObject: true, AsObject: *objBytes}, nil
	}

	asString := string(data)
	var mapper JsonMapper
	// check if value is int
	i, err := strconv.Atoi(asString)
	if err == nil {
		mapper.IsInt = true
		mapper.AsInt = i
		return JsonMapper{IsInt: true, AsInt: i}, nil
	}
	// check if value is float
	f, err := strconv.ParseFloat(asString, 64)
	if err == nil {
		mapper.IsFloat = true
		mapper.AsFloat = f
		return JsonMapper{IsFloat: true, AsFloat: f}, nil
	}
	// check if value is bool
	b, err := strconv.ParseBool(asString)
	if err == nil {
		mapper.IsBool = true
		mapper.AsBool = b
		return JsonMapper{IsBool: true, AsBool: b}, nil
	}
	// check if value is null
	if asString == "null" {
		mapper.IsNull = true
		return JsonMapper{IsNull: true}, nil
	}
	// fallback to string if no other type was found
	asString = strings.Trim(asString, `"`)
	mapper.IsString = true
	mapper.AsString = asString
	return JsonMapper{IsString: true, AsString: asString}, nil
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

func FromBuffer(reader io.Reader) (JsonMapper, error) {
	var m JsonMapper
	m.reader = reader
	m.buffer = make([]byte, bufferSize)
	return m, nil
}

// AsTime retrieves the value as uuid.UUID. Works only if the JSON value is a string.
func (m *JsonMapper) AsTime() (time.Time, error) {
	if !m.IsString {
		return time.Time{}, TimeTypeConversionErr
	}
	for _, layout := range timeLayouts {
		parsedTime, err := time.Parse(layout, m.AsString)
		if err == nil {
			return parsedTime, nil
		}
	}
	return time.Time{}, createNewInvalidTimeErr(m.AsString)
}

// AsUUID retrieves the value as uuid.UUID. Works only if the JSON value is a string.
func (m *JsonMapper) AsUUID() (uuid.UUID, error) {
	if !m.IsString {
		return uuid.Nil, nil
	}
	return uuid.Parse(m.AsString)
}

func (m *JsonMapper) ProcessObjectsWithArgs(numberOfWorkers int, f func(o JsonObject, args ...any), args ...any) error {
	if m.reader == nil {
		return errors.New("reader is not set")
	}
	if m.buffer == nil {
		m.buffer = make([]byte, bufferSize)
	}

	dec := json.NewDecoder(m.reader)
	_, err := dec.Token()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, numberOfWorkers)
	for dec.More() {
		var data map[string]*interface{}
		err = dec.Decode(&data)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		obj := newObjectFromMap(data)
		wg.Add(1)
		sem <- struct{}{}
		go func(o JsonObject) {
			defer wg.Done()
			defer func() { <-sem }()
			f(o, args...)
		}(*obj)
	}

	_, err = dec.Token()
	if err != nil {
		return err
	}

	wg.Wait()
	return nil
}

func (m *JsonMapper) ProcessObjects(numberOfWorkers int, f func(o JsonObject)) error {
	return m.ProcessObjectsWithArgs(numberOfWorkers, func(o JsonObject, args ...any) {
		f(o)
	})
}

// PrettyString returns a valid, pretty JSON string representation of the JsonMapper underlying value.
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
		return m.AsArray.PrettyString()
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
		return m.AsObject.String()
	case m.IsArray:
		return m.AsArray.String()
	}
	return ""
}
