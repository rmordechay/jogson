package jsonmapper

type JsonArray struct {
	length   int
	elements []interface{}
}

func Map[T JsonType](arr JsonArray, f func(mapper Mapper) T) []T {
	var jsonMappers []T
	for _, element := range arr.elements {
		field := f(getMapperFromField(element))
		jsonMappers = append(jsonMappers, field)
	}
	return jsonMappers
}

func ForEach(arr JsonArray, f func(mapper Mapper)) {
	for _, element := range arr.elements {
		f(getMapperFromField(element))
	}
}

func Filter(arr JsonArray, f func(mapper Mapper) bool) []Mapper {
	var jsonMappers []Mapper
	for _, element := range arr.elements {
		field := getMapperFromField(element)
		if f(field) {
			jsonMappers = append(jsonMappers, field)
		}
	}
	return jsonMappers
}

func (a JsonArray) Elements() []Mapper {
	jsons := make([]Mapper, 0, len(a.elements))
	for _, element := range a.elements {
		jsons = append(jsons, getMapperFromField(element))
	}
	return jsons
}

func (a JsonArray) Get(key int) Mapper {
	if key >= a.Length() {
		panic("index out of bound")
	}
	return getMapperFromField(a.elements[key])
}

func (a JsonArray) Length() int {
	return a.length
}

func (a JsonArray) String() string {
	return string(marshal(a.elements))
}

func CreateEmptyJsonArray() JsonArray {
	var arr JsonArray
	arr.elements = make([]interface{}, 0)
	return arr
}

func CreateJsonArray(data interface{}) JsonArray {
	var arr JsonArray
	arr.elements = data.([]interface{})
	arr.length = len(arr.elements)
	return arr
}

func parseJsonArray(data string) (JsonArray, error) {
	var ja JsonArray
	var arr []interface{}
	err := unmarshal([]byte(data), &arr)
	if err != nil {
		return JsonArray{}, err
	}
	ja.elements = arr
	ja.length = len(ja.elements)
	return ja, nil
}

func convertArray[T JsonType](data []T) JsonArray {
	var arr JsonArray
	result := make([]interface{}, len(data))
	for i, v := range data {
		result[i] = v
	}
	arr.elements = result
	arr.length = len(arr.elements)
	return arr
}
