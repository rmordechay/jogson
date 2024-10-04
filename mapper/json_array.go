package mapper

type JsonArray struct {
	object []interface{}
	value  []byte
}

func (j JsonArray) Has(key string) bool {
	panic("implement me")
}

func (j JsonArray) Get(i int) Json {
	panic("implement me")
}

func (j JsonArray) GetByName(key string) Json {
	return j
}

func (j JsonArray) AsString() string {
	panic("implement me")
}

func (j JsonArray) AsInt() int {
	panic("implement me")
}
