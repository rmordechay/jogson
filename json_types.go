package jsonmapper

import (
	"errors"
	"fmt"
	"time"
)

type jsonEntity interface {
	SetLastError(err error)
	String() string
	Length() int
}

type JsonType int

const (
	Bool JsonType = iota
	Int
	Float
	String
	Object
	Array
	Null
	Invalid
)

func (d JsonType) String() string {
	return [...]string{
		"bool",
		"int",
		"float64",
		"string",
		"Object",
		"Array",
		"null",
		"invalid",
	}[d]
}

var (
	nullConversionErrStr     = "null value cannot be converted to '%s'"
	typeConversionErrStr     = "'%T' could not be converted to %s"
	keyNotFoundErrStr        = "'%v'"
	indexOutOfRangeErrStr    = "[%v] with length %v"
	timeTypeConversionErrStr = "cannot convert type '%v' to time.Time"
	invalidTime              = "'%v' could not be parsed as time"
)

var (
	NullConversionErr     = errors.New("ERROR: null value conversion error")
	TypeConversionErr     = errors.New("ERROR: type conversion error")
	KeyNotFoundErr        = errors.New("ERROR: key was not found")
	IndexOutOfRangeErr    = errors.New("ERROR: index out of range")
	TimeTypeConversionErr = errors.New("ERROR: time conversion error")
	InvalidTimeErr        = errors.New("ERROR: invalid time")
)

func NewNullConversionErr(toType string) error {
	return errors.Join(NullConversionErr, fmt.Errorf(nullConversionErrStr, toType))
}

func NewTypeConversionErr(fromType any, toType any) error {
	return errors.Join(TypeConversionErr, fmt.Errorf(typeConversionErrStr, fromType, toType))
}

func NewKeyNotFoundErr(key string) error {
	return errors.Join(KeyNotFoundErr, fmt.Errorf(keyNotFoundErrStr, key))
}

func NewIndexOutOfRangeErr(i int, length int) error {
	return errors.Join(IndexOutOfRangeErr, fmt.Errorf(indexOutOfRangeErrStr, i, length))
}

func NewTimeTypeConversionErr(fromType JsonType) error {
	return errors.Join(TimeTypeConversionErr, fmt.Errorf(timeTypeConversionErrStr, fromType))
}

func NewInvalidTimeErr(v any) error {
	return errors.Join(InvalidTimeErr, fmt.Errorf(invalidTime, v))
}

var timeLayouts = []string{
	time.RFC3339,
	time.RFC850,
	time.RFC822,
	time.RFC822Z,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339Nano,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.Layout,
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
	time.DateTime,
	time.DateOnly,
	time.TimeOnly,
}
