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

func TestObjectAsStringMap(t *testing.T) {
	object, _ := jsonmapper.NewObjectFromString(jsonObjectOnlyStringTest)
	expectedMap := map[string]string{"first": "string1", "second": "string2", "third": "string3"}
	assert.Equal(t, expectedMap, object.AsStringMap())
	assert.Equal(t, 3, object.Length())
}

func TestObjectAsIntMap(t *testing.T) {
	object, _ := jsonmapper.NewObjectFromString(jsonObjectOnlyIntTest)
	expectedMap := map[string]int{"first": 1, "second": 3, "third": 54}
	assert.Equal(t, expectedMap, object.AsIntMap())
	assert.Equal(t, 3, object.Length())
}

func TestObjectAsFloatMap(t *testing.T) {
	object, _ := jsonmapper.NewObjectFromString(jsonObjectOnlyFloatTest)
	expectedMap := map[string]float64{"first": 5.3, "second": 1.4, "third": -0.3}
	assert.Equal(t, expectedMap, object.AsFloatMap())
	assert.Equal(t, 3, object.Length())
}

func TestObjectAsStringMapNullable(t *testing.T) {
	object, err := jsonmapper.NewObjectFromString(jsonObjectOnlyStringWithNullTest)
	assert.NoError(t, err)
	string1 := "string1"
	string2 := "string2"
	string3 := "string3"
	expectedMap := map[string]*string{
		"first": &string1, "second": &string2, "third": &string3, "fourth": nil,
	}
	stringMap := object.AsStringMapN()
	assert.Equal(t, expectedMap, stringMap)
	assert.Equal(t, 4, len(stringMap))
	assert.Equal(t, 4, object.Length())
	assert.Equal(t, 3, len(object.AsStringMap()))
}

func TestObjectAsIntMapNullable(t *testing.T) {
	object, err := jsonmapper.NewObjectFromString(jsonObjectOnlyIntWithNullTest)
	assert.NoError(t, err)
	i1 := 1
	i2 := 3
	i3 := 54
	expectedMap := map[string]*int{
		"first": &i1, "second": &i2, "third": &i3, "fourth": nil,
	}
	intMap := object.AsIntMapN()
	assert.Equal(t, expectedMap, intMap)
	assert.Equal(t, 4, len(intMap))
	assert.Equal(t, 4, object.Length())
	assert.Equal(t, 3, len(object.AsIntMap()))
}

func TestObjectAsFloatMapNullable(t *testing.T) {
	object, err := jsonmapper.NewObjectFromString(jsonObjectOnlyFloatWithNullTest)
	assert.NoError(t, err)
	f1 := 5.3
	f2 := 1.4
	f3 := -0.3
	expectedMap := map[string]*float64{
		"first": &f1, "second": &f2, "third": &f3, "fourth": nil,
	}
	floatMap := object.AsFloatMapN()
	assert.Equal(t, expectedMap, floatMap)
	assert.Equal(t, 4, len(floatMap))
	assert.Equal(t, 4, object.Length())
	assert.Equal(t, 3, len(object.AsFloatMap()))
}

func TestObjectGetMapper(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectTest)
	array := mapper.AsObject

	elementMapper := array.Get("name")
	assert.NoError(t, array.LastError)
	assert.Equal(t, "Jason", elementMapper.AsString)

	array.LastError = nil
	elementMapper = array.Get("age")
	assert.NoError(t, array.LastError)
	assert.Equal(t, 15, elementMapper.AsInt)
	assert.True(t, elementMapper.IsInt)

	array.LastError = nil
	elementMapper = array.Get("address")
	assert.NoError(t, array.LastError)
	assert.True(t, elementMapper.IsNull)

	array.LastError = nil
	elementMapper = array.Get("is_funny")
	assert.NoError(t, array.LastError)
	assert.Equal(t, true, elementMapper.AsBool)
	assert.True(t, elementMapper.IsBool)

	array.LastError = nil
	elementMapper = array.Get("height")
	assert.NoError(t, array.LastError)
	assert.Equal(t, 1.81, elementMapper.AsFloat)
	assert.True(t, elementMapper.IsFloat)

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
	assert.True(t, obj.IsNull())

	obj = object.GetObject("name")
	assert.ErrorIs(t, object.LastError, jsonmapper.TypeConversionErr)
	assert.True(t, obj.IsNull())

	obj = object.GetObject("address")
	assert.ErrorIs(t, object.LastError, jsonmapper.TypeConversionErr)
	assert.True(t, obj.IsNull())
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

func TestObjectGetUUID(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonUUIDTest)
	assert.NoError(t, err)
	uuid, err := mapper.AsObject.Get("uuid").AsUUID()
	assert.NoError(t, err)
	assert.Equal(t, "870fb3fd-d177-4ac4-a648-a33afd5ab288", uuid.String())

}

func TestObjectIsNull(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectTest)
	object := mapper.AsObject

	obj := object.GetObject("address")
	assert.True(t, obj.IsNull())
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

func TestElementNotFound(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonObjectTest)
	assert.NoError(t, err)
	_ = mapper.AsObject.GetFloat("not found")
	assert.Error(t, mapper.AsObject.LastError)
	assert.ErrorIs(t, mapper.AsObject.LastError, jsonmapper.KeyNotFoundErr)
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
