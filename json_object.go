package jsonmapper

type JsonObject interface {
	Has(key string) bool
	Get(key string) JsonMapper
	Find(key string) JsonMapper
	Elements() map[string]JsonMapper
	AddKeyValue(k string, value interface{})
	String() string
}

type jsonObject struct {
	object map[string]interface{}
}

func (o jsonObject) Has(key string) bool {
	for k := range o.object {
		if k == key {
			return true
		}
	}
	return false
}

func (o jsonObject) Get(key string) JsonMapper {
	for k, v := range o.object {
		if k == key {
			return getMapperFromField(v)
		}
	}
	return JsonMapper{}
}

func (o jsonObject) Find(key string) JsonMapper {
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

func (o jsonObject) Elements() map[string]JsonMapper {
	jsons := make(map[string]JsonMapper)
	for k, v := range o.object {
		jsons[k] = getMapperFromField(v)
	}
	return jsons
}

func (o jsonObject) AddKeyValue(k string, value interface{}) {
	o.object[k] = value
}

func (o jsonObject) String() string {
	return string(marshal(o.object))
}

func CreateEmptyJsonObject() JsonObject {
	var obj jsonObject
	obj.object = make(map[string]interface{})
	return obj
}

func CreateJsonObject(data interface{}) JsonObject {
	var obj jsonObject
	obj.object = data.(map[string]interface{})
	return obj
}

func parseJsonObject(data string) (jsonObject, error) {
	var jo jsonObject
	err := unmarshal([]byte(data), &jo.object)
	if err != nil {
		return jsonObject{}, err
	}
	return jo, nil
}
