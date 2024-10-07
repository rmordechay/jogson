package jsonmapper

type Object interface {
	Keys() []string
	Values() []Mapper
	Has(key string) bool
	Get(key string) Mapper
	Find(key string) Mapper
	Elements() map[string]Mapper
	AddKeyValue(k string, value interface{}) JsonObject
	String() string
}

type JsonObject struct {
	object map[string]*interface{}
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

func ForEachObject(obj JsonObject, f func(key string, mapper Mapper)) {
	for k, element := range obj.object {
		f(k, getMapperFromField(element))
	}
}

func MapObject[T JsonType](obj JsonObject, f func(key string, mapper Mapper) T) []T {
	var jsonMappers []T
	for k, element := range obj.object {
		field := f(k, getMapperFromField(element))
		jsonMappers = append(jsonMappers, field)
	}
	return jsonMappers
}

func FilterObject(obj JsonObject, f func(key string, mapper Mapper) bool) []Mapper {
	var jsonMappers []Mapper
	for k, element := range obj.object {
		field := getMapperFromField(element)
		if f(k, field) {
			jsonMappers = append(jsonMappers, field)
		}
	}
	return jsonMappers
}

func (o JsonObject) String() string {
	jsonBytes, _ := marshal(o.object)
	return string(jsonBytes)
}

func createJsonObject(data interface{}) JsonObject {
	var obj JsonObject
	var object = make(map[string]*interface{})
	for k, v := range data.(map[string]interface{}) {
		object[k] = &v
	}
	obj.object = object
	return obj
}

func parseJsonObject(data []byte) (JsonObject, error) {
	var jo JsonObject
	err := unmarshal(data, &jo.object)
	if err != nil {
		return JsonObject{}, err
	}
	return jo, nil
}
