package jsonmapper

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

func (a JsonArray) Get(key int) JsonMapper {
	if key >= a.Length {
		panic("index out of bound")
	}
	return getMapperFromField(a.elements[key])
}

func (a JsonArray) String() string {
	return string(marshal(a.elements))
}
