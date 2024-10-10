package tests

import (
	"github.com/rmordechay/jsonmapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestObjectGetKeys(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonObjectTest)
	assert.NoError(t, err)
	keys := mapper.Object.Keys()
	assert.Equal(t, 5, len(keys))
	assert.Contains(t, keys, "name")
	assert.Contains(t, keys, "age")
	assert.Contains(t, keys, "address")
}

func TestObjectGetValues(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonObjectTest)
	assert.NoError(t, err)
	values := mapper.Object.Values()
	assert.Equal(t, 5, len(values))
	for _, v := range values {
		assert.True(t, v.AsString == "Jason" || v.IsNull || v.AsInt == 15 || v.AsBool || v.AsFloat == 1.81)
	}
}

func TestElementNotFound(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonObjectTest)
	assert.NoError(t, err)
	_ = mapper.Object.GetFloat("not found")
	assert.Error(t, mapper.Object.LastError)
	assert.Equal(t, mapper.Object.LastError.Error(), "the requested key 'not found' was not found")
}

func TestObjectGetString(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectTest)
	object := mapper.Object

	s := object.GetString("name")
	assert.NoError(t, object.LastError)
	assert.Equal(t, "Jason", s)

	s = object.GetString("age")
	assert.NoError(t, object.LastError)
	assert.Equal(t, "15", s)

	s = object.GetString("height")
	assert.NoError(t, object.LastError)
	assert.Equal(t, "1.81", s)

	s = object.GetString("is_funny")
	assert.NoError(t, object.LastError)
	assert.Equal(t, "true", s)
}

func TestObjectGetStringFails(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectTest)
	object := mapper.Object

	s := object.GetString("not found")
	assert.Equal(t, "the requested key 'not found' was not found", object.LastError.Error())
	assert.Equal(t, "", s)

	s = object.GetString("address")
	assert.Equal(t, "value is null and could not be converted to string", object.LastError.Error())
	assert.Equal(t, "", s)
}

func TestObjectGetInt(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectTest)
	object := mapper.Object

	i := object.GetInt("age")
	assert.NoError(t, object.LastError)
	assert.Equal(t, 15, i)

	i = object.GetInt("height")
	assert.NoError(t, object.LastError)
	assert.Equal(t, 1, i)
}

func TestObjectGetIntFails(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectTest)
	object := mapper.Object

	i := object.GetInt("not found")
	assert.Equal(t, "the requested key 'not found' was not found", object.LastError.Error())
	assert.Equal(t, 0, i)

	i = object.GetInt("name")
	assert.Equal(t, "the type 'string' could not be converted to int", object.LastError.Error())
	assert.Equal(t, 0, i)

	i = object.GetInt("address")
	assert.Equal(t, "value is null and could not be converted to int", object.LastError.Error())
	assert.Equal(t, 0, i)
}

func TestObjectGetFloat(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectTest)
	object := mapper.Object

	f := object.GetFloat("age")
	assert.NoError(t, object.LastError)
	assert.Equal(t, float64(15), f)

	f = object.GetFloat("height")
	assert.NoError(t, object.LastError)
	assert.Equal(t, 1.81, f)
}

func TestObjectGetFloatFails(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectTest)
	object := mapper.Object

	f := object.GetFloat("not found")
	assert.Equal(t, "the requested key 'not found' was not found", object.LastError.Error())
	assert.Equal(t, float64(0), f)

	f = object.GetFloat("name")
	assert.Equal(t, "the type 'string' could not be converted to float64", object.LastError.Error())
	assert.Equal(t, float64(0), f)

	f = object.GetFloat("address")
	assert.Equal(t, "value is null and could not be converted to float64", object.LastError.Error())
	assert.Equal(t, float64(0), f)
}

func TestObjectGetBool(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectTest)
	object := mapper.Object

	b := object.GetBool("is_funny")
	assert.NoError(t, object.LastError)
	assert.Equal(t, true, b)
}

func TestObjectGetBoolFails(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectTest)
	object := mapper.Object

	b := object.GetBool("not found")
	assert.Equal(t, "the requested key 'not found' was not found", object.LastError.Error())
	assert.Equal(t, false, b)

	b = object.GetBool("age")
	assert.Equal(t, "the type 'float64' could not be converted to bool", object.LastError.Error())
	assert.Equal(t, false, b)

	b = object.GetBool("address")
	assert.Equal(t, "value is null and could not be converted to bool", object.LastError.Error())
	assert.Equal(t, false, b)
}

func TestObjectGetArray(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectWithArrayTest)
	object := mapper.Object

	array := object.GetArray("names")
	assert.ElementsMatch(t, []string{"Jason", "Chris", "Rachel"}, array.AsStringArray())
}

func TestObjectGetArrayFails(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectWithArrayTest)
	object := mapper.Object

	arr := object.GetArray("not found")
	assert.Equal(t, "the requested key 'not found' was not found", object.LastError.Error())
	assert.Equal(t, &jsonmapper.JsonArray{}, arr)

	arr = object.GetArray("name")
	assert.Equal(t, "the type 'string' could not be converted to []*any", object.LastError.Error())
	assert.Equal(t, &jsonmapper.JsonArray{}, arr)

	arr = object.GetArray("address")
	assert.Equal(t, "value is null and could not be converted to JsonArray", object.LastError.Error())
	assert.Equal(t, &jsonmapper.JsonArray{}, arr)
}

func TestObjectGetObject(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectNestedArrayTest)
	object := mapper.Object

	obj := object.GetObject("personTest")
	assert.Equal(t, "Jason", obj.GetString("name"))
}

func TestObjectGetObjectFails(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectTest)
	object := mapper.Object

	obj := object.GetObject("not found")
	assert.Equal(t, "the requested key 'not found' was not found", object.LastError.Error())
	assert.Equal(t, &jsonmapper.JsonObject{}, obj)

	obj = object.GetObject("name")
	assert.Equal(t, "the type 'string' could not be converted to JsonObject", object.LastError.Error())
	assert.Equal(t, &jsonmapper.JsonObject{}, obj)

	obj = object.GetObject("address")
	assert.Equal(t, "value is null and could not be converted to JsonObject", object.LastError.Error())
	assert.Equal(t, &jsonmapper.JsonObject{}, obj)

}

func TestConvertKeysToSnakeCase(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonObjectKeysPascalCaseTest)
	assert.NoError(t, err)
	object := mapper.Object
	snakeCase := object.TransformObjectKeys()
	assert.NoError(t, object.LastError)
	assert.ElementsMatch(t, []string{"children", "name", "age", "address", "second_address", "is_funny"}, snakeCase.Keys())
	children := snakeCase.GetObject("children")
	assert.NoError(t, snakeCase.LastError)
	rachel := children.GetObject("rachel")
	assert.NoError(t, children.LastError)
	age := rachel.GetInt("age")
	isFunny := rachel.GetBool("is_funny")
	assert.Equal(t, 15, age)
	assert.True(t, isFunny)
}
