package mapper

import (
	"strconv"
)

type Json interface {
	GetByName(key string) Json
	Get(i int) Json
	Has(key string) bool
	AsString() string
	AsInt() int
}

type JsonObject struct {
	object map[string]interface{}
	value  []byte
}

func (o JsonObject) Get(i int) Json {
	panic("implement me")
}

func (o JsonObject) Has(key string) bool {
	for k := range o.object {
		if k == key {
			return true
		}
	}
	return false
}

func (o JsonObject) GetByName(key string) Json {
	for k, v := range o.object {
		if k == key {
			return createObjectMapper(v)
		}
	}
	return o
}

func (o JsonObject) AsString() string {
	return string(o.value)
}

func (o JsonObject) AsInt() int {
	toInt, _ := strconv.Atoi(string(o.value))
	return toInt
}
