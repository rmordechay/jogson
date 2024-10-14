package jogson

import (
	"time"
)

type jsonI interface {
	String() string
	PrettyString() string
	Length() int
	IsNull() bool
	IsEmpty() bool
	setLastError(err error)
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
