package jsonmapper

type JsonArray interface {
	Elements() []JsonMapper
	Get(key int) JsonMapper
	String() string
	Length() int
}

type jsonArray struct {
	length   int
	elements []interface{}
}

func (a jsonArray) Elements() []JsonMapper {
	jsons := make([]JsonMapper, 0, len(a.elements))
	for _, element := range a.elements {
		jsons = append(jsons, getMapperFromField(element))
	}
	return jsons
}

func (a jsonArray) Get(key int) JsonMapper {
	if key >= a.Length() {
		panic("index out of bound")
	}
	return getMapperFromField(a.elements[key])
}

func (a jsonArray) Length() int {
	return a.length
}

func (a jsonArray) String() string {
	return string(marshal(a.elements))
}

func CreateEmptyJsonArray() JsonArray {
	var arr jsonArray
	arr.elements = make([]interface{}, 0)
	return arr
}

func CreateJsonArray(data interface{}) JsonArray {
	var arr jsonArray
	arr.elements = data.([]interface{})
	arr.length = len(arr.elements)
	return arr
}

func parseJsonArray(data string) (jsonArray, error) {
	var ja jsonArray
	var arr []interface{}
	err := unmarshal([]byte(data), &arr)
	if err != nil {
		return jsonArray{}, err
	}
	ja.elements = arr
	ja.length = len(ja.elements)
	return ja, nil
}

func convertArray[T JsonType](data []T) jsonArray {
	var arr jsonArray
	result := make([]interface{}, len(data))
	for i, v := range data {
		result[i] = v
	}
	arr.elements = result
	arr.length = len(arr.elements)
	return arr
}
