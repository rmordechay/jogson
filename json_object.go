package jsonmapper

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

func (o JsonObject) AddKeyValue(k string, value interface{}) {
	o.object[k] = value
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
