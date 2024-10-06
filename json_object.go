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

func (o JsonObject) Get(key string) JsonMapper {
	for k, v := range o.object {
		if k == key {
			return getMapperFromField(v)
		}
	}
	return JsonMapper{}
}

func (o JsonObject) Find(key string) JsonMapper {
	for k, v := range o.object {
		field := getMapperFromField(v)
		if k == key {
			return field
		}
		if field.IsObject {
			return field.Object.Find(key)
		}
	}
	return JsonMapper{}
}

func (o JsonObject) Elements() map[string]JsonMapper {
	jsons := make(map[string]JsonMapper)
	for k, v := range o.object {
		jsons[k] = getMapperFromField(v)
	}
	return jsons
}

func (o JsonObject) String() string {
	return string(marshal(o.object))
}
