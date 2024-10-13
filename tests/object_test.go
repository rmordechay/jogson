package tests

import (
	"github.com/rmordechay/jsonmapper"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestObjectGetKeys(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonObjectTest)
	assert.NoError(t, err)
	keys := mapper.AsObject.Keys()
	assert.Equal(t, 5, len(keys))
	assert.Contains(t, keys, "name")
	assert.Contains(t, keys, "age")
	assert.Contains(t, keys, "address")
}

func TestObjectGetValues(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonObjectTest)
	assert.NoError(t, err)
	values := mapper.AsObject.Values()
	assert.Equal(t, 5, len(values))
	for _, v := range values {
		assert.True(t, v.AsString == "Jason" || v.IsNull || v.AsInt == 15 || v.AsBool || v.AsFloat == 1.81)
	}
}

func TestElementNotFound(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonObjectTest)
	assert.NoError(t, err)
	_ = mapper.AsObject.GetFloat("not found")
	assert.Error(t, mapper.AsObject.LastError)
	assert.ErrorIs(t, mapper.AsObject.LastError, jsonmapper.KeyNotFoundErr)
}

func TestObjectGetString(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectTest)
	object := mapper.AsObject

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
	object := mapper.AsObject

	s := object.GetString("not found")
	assert.ErrorIs(t, object.LastError, jsonmapper.KeyNotFoundErr)
	assert.Equal(t, "", s)

	s = object.GetString("address")
	assert.ErrorIs(t, object.LastError, jsonmapper.TypeConversionErr)
	assert.Equal(t, "", s)
}

func TestObjectGetInt(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectTest)
	object := mapper.AsObject

	i := object.GetInt("age")
	assert.NoError(t, object.LastError)
	assert.Equal(t, 15, i)

	i = object.GetInt("height")
	assert.NoError(t, object.LastError)
	assert.Equal(t, 1, i)
}

func TestObjectGetIntFails(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectTest)
	object := mapper.AsObject

	i := object.GetInt("not found")
	assert.ErrorIs(t, object.LastError, jsonmapper.KeyNotFoundErr)
	assert.Equal(t, 0, i)

	i = object.GetInt("name")
	assert.ErrorIs(t, object.LastError, jsonmapper.TypeConversionErr)
	assert.Equal(t, 0, i)

	i = object.GetInt("address")
	assert.ErrorIs(t, object.LastError, jsonmapper.TypeConversionErr)
	assert.Equal(t, 0, i)
}

func TestObjectGetFloat(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectTest)
	object := mapper.AsObject

	f := object.GetFloat("age")
	assert.NoError(t, object.LastError)
	assert.Equal(t, float64(15), f)

	f = object.GetFloat("height")
	assert.NoError(t, object.LastError)
	assert.Equal(t, 1.81, f)
}

func TestObjectGetFloatFails(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectTest)
	object := mapper.AsObject

	f := object.GetFloat("not found")
	assert.ErrorIs(t, object.LastError, jsonmapper.KeyNotFoundErr)
	assert.Equal(t, float64(0), f)

	f = object.GetFloat("name")
	assert.ErrorIs(t, object.LastError, jsonmapper.TypeConversionErr)
	assert.Equal(t, float64(0), f)

	f = object.GetFloat("address")
	assert.ErrorIs(t, object.LastError, jsonmapper.TypeConversionErr)
	assert.Equal(t, float64(0), f)
}

func TestObjectGetBool(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectTest)
	object := mapper.AsObject

	b := object.GetBool("is_funny")
	assert.NoError(t, object.LastError)
	assert.Equal(t, true, b)
}

func TestObjectGetBoolFails(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectTest)
	object := mapper.AsObject

	b := object.GetBool("not found")
	assert.ErrorIs(t, object.LastError, jsonmapper.KeyNotFoundErr)
	assert.Equal(t, false, b)

	b = object.GetBool("age")
	assert.ErrorIs(t, object.LastError, jsonmapper.TypeConversionErr)
	assert.Equal(t, false, b)

	b = object.GetBool("address")
	assert.ErrorIs(t, object.LastError, jsonmapper.TypeConversionErr)
	assert.Equal(t, false, b)
}

func TestObjectGetArray(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectWithArrayTest)
	object := mapper.AsObject

	array := object.GetArray("names")
	assert.ElementsMatch(t, []string{"Jason", "Chris", "Rachel"}, array.AsStringArray())
}

func TestObjectGetArrayFails(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectWithArrayTest)
	object := mapper.AsObject

	arr := object.GetArray("not found")
	assert.ErrorIs(t, object.LastError, jsonmapper.KeyNotFoundErr)
	assert.Equal(t, jsonmapper.EmptyArray(), arr)

	arr = object.GetArray("name")
	assert.ErrorIs(t, object.LastError, jsonmapper.TypeConversionErr)
	assert.Equal(t, jsonmapper.EmptyArray(), arr)

	arr = object.GetArray("address")
	assert.ErrorIs(t, object.LastError, jsonmapper.TypeConversionErr)
	assert.Equal(t, jsonmapper.EmptyArray(), arr)
}

func TestObjectGetObject(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectNestedArrayTest)
	object := mapper.AsObject

	obj := object.GetObject("personTest")
	assert.Equal(t, "Jason", obj.GetString("name"))
}

func TestObjectGetObjectFails(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectTest)
	object := mapper.AsObject

	obj := object.GetObject("not found")
	assert.ErrorIs(t, object.LastError, jsonmapper.KeyNotFoundErr)
	assert.Equal(t, jsonmapper.EmptyObject(), obj)

	obj = object.GetObject("name")
	assert.ErrorIs(t, object.LastError, jsonmapper.TypeConversionErr)
	assert.Equal(t, jsonmapper.EmptyObject(), obj)

	obj = object.GetObject("address")
	assert.ErrorIs(t, object.LastError, jsonmapper.TypeConversionErr)
	assert.Equal(t, jsonmapper.EmptyObject(), obj)
}

func TestObjectGetTime(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonObjectTimeTest)
	assert.NoError(t, err)
	object := mapper.AsObject
	actualTime1 := object.GetTime("time1")
	assert.NoError(t, object.LastError)
	actualTime2 := object.GetTime("time2")
	assert.NoError(t, object.LastError)
	actualTime3 := object.GetTime("time3")
	assert.NoError(t, object.LastError)
	expectedTime1, _ := time.Parse(time.RFC3339, "2024-10-06T17:59:44Z")
	expectedTime2, _ := time.Parse(time.RFC3339, "2024-10-06T17:59:44+00:00")
	expectedTime3, _ := time.Parse(time.RFC850, "Sunday, 06-Oct-24 17:59:44 UTC")
	assert.Equal(t, expectedTime1, actualTime1)
	assert.Equal(t, expectedTime2, actualTime2)
	assert.Equal(t, expectedTime3, actualTime3)
}

func TestConvertKeysToSnakeCase(t *testing.T) {
	//mapper, err := jsonmapper.FromString(jsonObjectKeysPascalCaseTest)
	//assert.NoError(t, err)
	//object := mapper.AsObject
	//snakeCase := object.transformObjectKeys()
	//assert.NoError(t, object.LastError)
	//assert.ElementsMatch(t, []string{"children", "name", "age", "address", "second_address", "is_funny"}, snakeCase.Keys())
	//children := snakeCase.GetObject("children")
	//assert.NoError(t, snakeCase.LastError)
	//rachel := children.GetObject("rachel")
	//assert.NoError(t, children.LastError)
	//age := rachel.GetInt("age")
	//isFunny := rachel.GetBool("is_funny")
	//assert.Equal(t, 15, age)
	//assert.True(t, isFunny)
}

func TestObjectString(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonObjectTest)
	assert.NoError(t, err)
	s := mapper.AsObject.String()
	assert.Equal(t, `{"address":null,"age":15,"height":1.81,"is_funny":true,"name":"Jason"}`, s)
}

func TestObjectPrettyString(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonObjectTest)
	assert.NoError(t, err)
	expectedStr := "{\n  \"address\": null,\n  \"age\": 15,\n  \"height\": 1.81,\n  \"is_funny\": true,\n  \"name\": \"Jason\"\n}"
	assert.Equal(t, expectedStr, mapper.AsObject.PrettyString())
}
