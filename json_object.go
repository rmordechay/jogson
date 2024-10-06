package jsonmapper

type Object interface {
	Has(key string) bool
	Get(key string) Mapper
	Find(key string) Mapper
	Elements() map[string]Mapper
	AddKeyValue(k string, value interface{}) JsonObject
	String() string
}

type JsonObject struct {
	object map[string]interface{}
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

func (o JsonObject) AddKeyValue(k string, value interface{}) JsonObject {
	o.object[k] = value
	return o
}

func CreateEmptyJsonObject() JsonObject {
	var obj JsonObject
	obj.object = make(map[string]interface{})
	return obj
}

func CreateJsonObject(data interface{}) JsonObject {
	var obj JsonObject
	obj.object = data.(map[string]interface{})
	return obj
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
	return string(marshal(o.object))
}

func parseJsonObject(data string) (JsonObject, error) {
	var jo JsonObject
	err := unmarshal([]byte(data), &jo.object)
	if err != nil {
		return JsonObject{}, err
	}
	return jo, nil
}
