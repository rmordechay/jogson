package jsonmapper

type JsonObject struct {
	object map[string]*interface{}
}

func NewObject() JsonObject {
	var obj JsonObject
	obj.object = make(map[string]*interface{})
	return obj
}

func (o JsonObject) Keys() []string {
	keys := make([]string, 0, len(o.object))
	for key := range o.object {
		keys = append(keys, key)
	}
	return keys
}

func (o JsonObject) Values() []Mapper {
	values := make([]Mapper, 0, len(o.object))
	for _, v := range o.object {
		values = append(values, getMapperFromField(v))
	}
	return values
}

func (o JsonObject) Has(key string) bool {
	for k := range o.object {
		if k == key {
			return true
		}
	}
	return false
}

func (o JsonObject) Get(key string) Mapper {
	for k, v := range o.object {
		if k == key {
			return getMapperFromField(v)
		}
	}
	return Mapper{}
}

func (o JsonObject) Find(key string) Mapper {
	for k, v := range o.object {
		field := getMapperFromField(v)
		if k == key {
			return field
		}
		if field.IsObject {
			return field.Object.Find(key)
		}
	}
	return Mapper{}
}

func (o JsonObject) Elements() map[string]Mapper {
	jsons := make(map[string]Mapper)
	for k, v := range o.object {
		jsons[k] = getMapperFromField(v)
	}
	return jsons
}

func (o JsonObject) AddKeyValue(k string, value interface{}) {
	o.object[k] = &value
}

func (o JsonObject) ForEach(f func(key string, mapper Mapper)) {
	for k, element := range o.object {
		f(k, getMapperFromField(element))
	}
}

func (o JsonObject) Filter(f func(key string, mapper Mapper) bool) JsonObject {
	var obj = NewObject()
	for k, element := range o.object {
		field := getMapperFromField(element)
		if f(k, field) {
			obj.object[k] = element
		}
	}
	return obj
}

func (o JsonObject) PrettyString() string {
	jsonBytes, _ := marshalIndent(o.object)
	return string(jsonBytes)
}

func (o JsonObject) String() string {
	jsonBytes, _ := marshal(o.object)
	return string(jsonBytes)
}
