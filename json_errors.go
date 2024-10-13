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
	return fmt.Errorf("%w: %w", TypeConversionErr, fmt.Errorf(typeConversionErrStr, fromType, toType))
}

func createKeyNotFoundErr(key string) error {
	return fmt.Errorf("%w: %w", KeyNotFoundErr, fmt.Errorf(keyNotFoundErrStr, key))
}

func createIndexOutOfRangeErr(i int, length int) error {
	return fmt.Errorf("%w: %w", IndexOutOfRangeErr, fmt.Errorf(indexOutOfRangeErrStr, i, length))
}

func createNewInvalidTimeErr(v any) error {
	return fmt.Errorf("%w: %w", InvalidTimeErr, fmt.Errorf(invalidTime, v))
}
