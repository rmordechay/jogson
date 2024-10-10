package jsonmapper

import (
	"fmt"
	"strconv"
)

func getAsString(data *any, j JsonError) string {
	if data == nil {
		j.SetLastError(fmt.Errorf(nullConversionErrStr, "string"))
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
		j.SetLastError(fmt.Errorf(typeConversionErrStr, *data, "string"))
		return ""
	}
}

func getAsInt(data *any, j JsonError) int {
	if data == nil {
		j.SetLastError(fmt.Errorf(nullConversionErrStr, "int"))
		return 0
	}
	switch (*data).(type) {
	case float64:
		return int((*data).(float64))
	case int:
		return (*data).(int)
	default:
		j.SetLastError(fmt.Errorf(typeConversionErrStr, *data, "int"))
		return 0
	}
}

func getAsFloat(data *any, j JsonError) float64 {
	if data == nil {
		j.SetLastError(fmt.Errorf(nullConversionErrStr, "float64"))
		return 0
	}
	v, ok := (*data).(float64)
	if !ok {
		j.SetLastError(fmt.Errorf(typeConversionErrStr, *data, "float64"))
		return 0
	}
	return v
}

func getAsBool(data *any, j JsonError) bool {
	if data == nil {
		j.SetLastError(fmt.Errorf(nullConversionErrStr, "bool"))
		return false
	}
	v, ok := (*data).(bool)
	if !ok {
		j.SetLastError(fmt.Errorf(typeConversionErrStr, *data, "bool"))
		return false
	}
	return v
}
