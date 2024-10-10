package jsonmapper

import (
	"time"
)

type jsonI interface {
	String() string
	Length() int
	IsEmpty() bool
	setLastError(err error)
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

const (
	boolTypeStr   = "bool"
	intTypeStr    = "int"
	floatTypeStr  = "float64"
	stringTypeStr = "string"
	objectTypeStr = "JsonObject"
	arrayTypeStr  = "JsonArray"
	timeTypeStr   = "time.Time"
)

func (d JsonType) String() string {
	return [...]string{
		boolTypeStr,
		intTypeStr,
		floatTypeStr,
		stringTypeStr,
		objectTypeStr,
		arrayTypeStr,
		"null",
		"invalid",
	}[d]
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
