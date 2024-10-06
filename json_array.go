package jsonmapper

import "encoding/json"

type JsonArray struct {
	Length   int
	elements []interface{}
}

func (a JsonArray) Elements() []JsonMapper {
	jsons := make([]JsonMapper, 0, len(a.elements))
	for _, element := range a.elements {
		jsons = append(jsons, getMapperFromField(element))
	}
	return jsons
}

func (a JsonArray) String() string {
	marshal, _ := json.Marshal(a.elements)
	return string(marshal)
}
