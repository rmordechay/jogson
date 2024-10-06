package main

import (
	"encoding/json"
)

type JsonObject struct {
	object map[string]interface{}
}

func (o JsonObject) Get(key string) Json {
	for k, v := range o.object {
		if k == key {
			return getMapperFromField(v)
		}
	}
	return Json{}
}

func (o JsonObject) Elements() map[string]Json {
	jsons := make(map[string]Json)
	for k, v := range o.object {
		jsons[k] = getMapperFromField(v)
	}
	return jsons
}

func (o JsonObject) String() string {
	marshal, _ := json.Marshal(o.object)
	return string(marshal)
}
