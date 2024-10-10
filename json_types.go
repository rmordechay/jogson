package jsonmapper

import "time"

type JsonError interface {
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
	nullConversionErrStr  = "value is null and could not be converted to %s"
	typeConversionErrStr  = "the type '%T' could not be converted to %s"
	keyNotFoundErrStr     = "the requested key '%v' was not found"
	indexOutOfRangeErrStr = "index out of range [%v] with length %v"
)

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
