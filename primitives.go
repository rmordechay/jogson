package jsonmapper

import (
	"fmt"
	"strconv"
)

func getAsString(data *interface{}, j Json) string {
	if data == nil {
		j.SetLastError(fmt.Errorf(NullConversionErrStr, ""))
		return ""
	}
	switch (*data).(type) {
	case string:
		return (*data).(string)
	case float64:
		return strconv.FormatFloat((*data).(float64), 'f', -1, 64)
	case int:
		return strconv.Itoa((*data).(int))
	case bool:
		return strconv.FormatBool((*data).(bool))
	default:
		j.SetLastError(fmt.Errorf(TypeConversionErrStr, *data, ""))
		return ""
	}
}

func getAsInt(data *interface{}, j Json) int {
	if data == nil {
		j.SetLastError(fmt.Errorf(NullConversionErrStr, 0))
		return 0
	}
	switch (*data).(type) {
	case float64:
		return int((*data).(float64))
	case int:
		return (*data).(int)
	default:
		j.SetLastError(fmt.Errorf(TypeConversionErrStr, *data, 0))
		return 0
	}
}

func getAsFloat(data *interface{}, j Json) float64 {
	if data == nil {
		j.SetLastError(fmt.Errorf(NullConversionErrStr, 0.0))
		return 0
	}
	v, ok := (*data).(float64)
	if !ok {
		j.SetLastError(fmt.Errorf(TypeConversionErrStr, *data, 0.0))
		return 0
	}
	return v
}

func getAsBool(data *interface{}, j Json) bool {
	if data == nil {
		j.SetLastError(fmt.Errorf(NullConversionErrStr, false))
		return false
	}
	v, ok := (*data).(bool)
	if !ok {
		j.SetLastError(fmt.Errorf(TypeConversionErrStr, *data, false))
		return false
	}
	return v
}
