package jsonmapper

type JsonWriter struct{}

func (w JsonWriter) CreateJsonObject() JsonMapper {
	return JsonMapper{}
}
