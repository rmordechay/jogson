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

func TestGetString(t *testing.T) {
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

func TestGetStringFails(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectTest)
	object := mapper.Object

	s := object.GetString("not found")
	assert.Equal(t, "the requested key 'not found' was not found", object.LastError.Error())
	assert.Equal(t, "", s)

	s = object.GetString("address")
	assert.Equal(t, "the value of 'address' is null and could not be converted to string", object.LastError.Error())
	assert.Equal(t, "", s)
}

func TestGetInt(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectTest)
	object := mapper.Object

	i := object.GetInt("age")
	assert.NoError(t, object.LastError)
	assert.Equal(t, 15, i)

	i = object.GetInt("height")
	assert.NoError(t, object.LastError)
	assert.Equal(t, 1, i)
}

func TestGetIntFails(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectTest)
	object := mapper.Object

	i := object.GetInt("not found")
	assert.Equal(t, "the requested key 'not found' was not found", object.LastError.Error())
	assert.Equal(t, 0, i)

	i = object.GetInt("name")
	assert.Equal(t, "the type of key 'name' (string) could not be converted to int", object.LastError.Error())
	assert.Equal(t, 0, i)

	i = object.GetInt("address")
	assert.Equal(t, "the value of 'address' is null and could not be converted to int", object.LastError.Error())
	assert.Equal(t, 0, i)
}

func TestGetFloat(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectTest)
	object := mapper.Object

	f := object.GetFloat("age")
	assert.NoError(t, object.LastError)
	assert.Equal(t, float64(15), f)

	f = object.GetFloat("height")
	assert.NoError(t, object.LastError)
	assert.Equal(t, 1.81, f)
}

func TestGetFloatFails(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectTest)
	object := mapper.Object

	f := object.GetFloat("not found")
	assert.Equal(t, "the requested key 'not found' was not found", object.LastError.Error())
	assert.Equal(t, float64(0), f)

	f = object.GetFloat("name")
	assert.Equal(t, "the type of key 'name' (string) could not be converted to float64", object.LastError.Error())
	assert.Equal(t, float64(0), f)

	f = object.GetFloat("address")
	assert.Equal(t, "the value of 'address' is null and could not be converted to float64", object.LastError.Error())
	assert.Equal(t, float64(0), f)
}

func TestGetBool(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectTest)
	object := mapper.Object

	b := object.GetBool("is_funny")
	assert.NoError(t, object.LastError)
	assert.Equal(t, true, b)
}

func TestGetBoolFails(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectTest)
	object := mapper.Object

	b := object.GetBool("not found")
	assert.Equal(t, "the requested key 'not found' was not found", object.LastError.Error())
	assert.Equal(t, false, b)

	b = object.GetBool("age")
	assert.Equal(t, "the type of key 'age' (float64) could not be converted to bool", object.LastError.Error())
	assert.Equal(t, false, b)

	b = object.GetBool("address")
	assert.Equal(t, "the value of 'address' is null and could not be converted to bool", object.LastError.Error())
	assert.Equal(t, false, b)
}

//func TestObjectFilter(t *testing.T) {
//	panic("not implemented")
//}
//
//func TestObjectForEach(t *testing.T) {
//	panic("not implemented")
//}
