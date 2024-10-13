package jsonmapper

import (
	"errors"
	"fmt"
)

var (
	typeConversionErrStr  = "%T could not be converted to %T"
	keyNotFoundErrStr     = "'%v'"
	indexOutOfRangeErrStr = "[%v] with length %v"
	invalidTime           = "'%v' could not be parsed as time"
)

var (
	TypeConversionErr     = errors.New("type conversion error")
	KeyNotFoundErr        = errors.New("key was not found")
	IndexOutOfRangeErr    = errors.New("index out of range")
	TimeTypeConversionErr = errors.New("time conversion error")
	InvalidTimeErr        = errors.New("invalid time")
)

func createTypeConversionErr(fromType any, toType any) error {
	return errors.Join(TypeConversionErr, fmt.Errorf(typeConversionErrStr, fromType, toType))
}

func createKeyNotFoundErr(key string) error {
	return errors.Join(KeyNotFoundErr, fmt.Errorf(keyNotFoundErrStr, key))
}

func createIndexOutOfRangeErr(i int, length int) error {
	return errors.Join(IndexOutOfRangeErr, fmt.Errorf(indexOutOfRangeErrStr, i, length))
}

func createNewInvalidTimeErr(v any) error {
	return errors.Join(InvalidTimeErr, fmt.Errorf(invalidTime, v))
}
