package jsonmapper

import (
	"errors"
	"fmt"
)

var (
	nullConversionErrStr  = "null value cannot be converted to '%s'"
	typeConversionErrStr  = "'%T' could not be converted to %s"
	keyNotFoundErrStr     = "'%v'"
	indexOutOfRangeErrStr = "[%v] with length %v"
	invalidTime           = "'%v' could not be parsed as time"
)

var (
	NullConversionErr     = errors.New("ERROR: null value conversion error")
	TypeConversionErr     = errors.New("ERROR: type conversion error")
	KeyNotFoundErr        = errors.New("ERROR: key was not found")
	IndexOutOfRangeErr    = errors.New("ERROR: index out of range")
	TimeTypeConversionErr = errors.New("ERROR: time conversion error")
	InvalidTimeErr        = errors.New("ERROR: invalid time")
)

func createNullConversionErr(toType string) error {
	return errors.Join(NullConversionErr, fmt.Errorf(nullConversionErrStr, toType))
}

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
