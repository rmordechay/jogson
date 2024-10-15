package tests

import (
	"math"
	"testing"
	"time"

	"github.com/rmordechay/jogson"
	"github.com/stretchr/testify/assert"
)

func TestArrayAsStringArray(t *testing.T) {
	mapper, _ := jogson.NewMapperFromString(jsonStringArrayTest)
	array := mapper.AsArray
	assert.ElementsMatch(t, []string{"Jason", "Chris", "Rachel"}, array.AsStringArray())
}

func TestArrayAsIntArray(t *testing.T) {
	mapper, _ := jogson.NewMapperFromString(jsonIntArrayTest)
	array := mapper.AsArray
	assert.ElementsMatch(t, []int{0, 15, -54, -346, 9223372036854775807}, array.AsIntArray())
}

func TestArrayAsFloatArray(t *testing.T) {
	mapper, _ := jogson.NewMapperFromString(jsonFloatArrayTest)
	array := mapper.AsArray
	assert.ElementsMatch(t, []float64{15.13, 2, 45.3984, -1.81, 9.223372036854776}, array.AsFloatArray())
}

func TestArrayAsStringArrayNullable(t *testing.T) {
	array, err := jogson.NewArrayFromString(jsonStringArrayWithNullTest)
	assert.NoError(t, err)
	s := "Chris"
	s2 := "Rachel"
	assert.ElementsMatch(t, []*string{nil, &s, &s2}, array.AsStringArrayN())
}

func TestArrayAsIntArrayNullable(t *testing.T) {
	array, err := jogson.NewArrayFromString(jsonIntArrayWithNullTest)
	assert.NoError(t, err)
	i := 0
	assert.ElementsMatch(t, []*int{&i, nil}, array.AsIntArrayN())
}

func TestArrayAsFloatArrayNullable(t *testing.T) {
	array, err := jogson.NewArrayFromString(jsonFloatArrayWithNullTest)
	assert.NoError(t, err)
	var i float64 = 2
	f := 15.13
	assert.ElementsMatch(t, []*float64{&f, &i, nil}, array.AsFloatArrayN())
}

func TestArrayContainsString(t *testing.T) {
	array, err := jogson.NewArrayFromString(jsonStringArrayWithNullTest)
	assert.NoError(t, err)
	assert.True(t, array.ContainsString("Chris"))
	assert.True(t, array.ContainsString("Rachel"))
}

func TestArrayContainsInt(t *testing.T) {
	array, err := jogson.NewArrayFromString(jsonIntArrayTest)
	assert.NoError(t, err)
	assert.True(t, array.ContainsInt(9223372036854775807))
	assert.True(t, array.ContainsInt(-346))
}

func TestArrayContainsFloat(t *testing.T) {
	array, err := jogson.NewArrayFromString(jsonFloatArrayWithNullTest)
	assert.NoError(t, err)
	assert.True(t, array.ContainsFloat(2))
	assert.True(t, array.ContainsFloat(15.13))
}

func TestArrayAs2DArray(t *testing.T) {
	mapper, _ := jogson.NewMapperFromString(json2DIntArrayTest)
	array := mapper.AsArray
	jsonArray := array.As2DArray()
	assert.ElementsMatch(t, []int{1, 2}, jsonArray[0].AsIntArray())
	assert.ElementsMatch(t, []int{3, 4}, array.As2DArray()[1].AsIntArray())
}

func TestArrayAsObjectArray(t *testing.T) {
	mapper, _ := jogson.NewMapperFromString(jsonObjectArrayTest)
	array := mapper.AsArray
	assert.Equal(t, "Jason", array.AsObjectArray()[0].GetString("name"))
	assert.Equal(t, "Chris", array.AsObjectArray()[1].GetString("name"))
}

func TestArrayGetMapper(t *testing.T) {
	mapper, _ := jogson.NewMapperFromString(jsonAnyArrayTest)
	array := mapper.AsArray

	elementMapper := array.Get(0)
	assert.NoError(t, array.LastError)
	assert.Equal(t, "Jason", elementMapper.AsString)

	array.LastError = nil
	elementMapper = array.Get(1)
	assert.NoError(t, array.LastError)
	assert.Equal(t, 15, elementMapper.AsInt)
	assert.True(t, elementMapper.IsInt)

	array.LastError = nil
	elementMapper = array.Get(2)
	assert.NoError(t, array.LastError)
	assert.True(t, elementMapper.IsNull)

	array.LastError = nil
	elementMapper = array.Get(3)
	assert.NoError(t, array.LastError)
	assert.Equal(t, 1.81, elementMapper.AsFloat)
	assert.True(t, elementMapper.IsFloat)

	array.LastError = nil
	elementMapper = array.Get(4)
	assert.NoError(t, array.LastError)
	assert.Equal(t, true, elementMapper.AsBool)
	assert.True(t, elementMapper.IsBool)
}

func TestArrayGetStringN(t *testing.T) {
	array, err := jogson.NewArrayFromString(jsonAnyArrayTest)
	assert.NoError(t, err)
	assert.Nil(t, array.GetStringN(2))
	assert.Equal(t, "Jason", *array.GetStringN(0))
}

func TestArrayGetIntN(t *testing.T) {
	array, err := jogson.NewArrayFromString(jsonAnyArrayTest)
	assert.NoError(t, err)
	assert.Nil(t, array.GetIntN(2))
	assert.Equal(t, 15, *array.GetIntN(1))
}

func TestArrayGetFloatN(t *testing.T) {
	array, err := jogson.NewArrayFromString(jsonAnyArrayTest)
	assert.NoError(t, err)
	assert.Nil(t, array.GetFloatN(2))
	assert.Equal(t, 1.81, *array.GetFloatN(3))
}

func TestArrayGetBoolN(t *testing.T) {
	array, err := jogson.NewArrayFromString(jsonAnyArrayTest)
	assert.NoError(t, err)
	assert.Nil(t, array.GetBoolN(2))
	assert.Equal(t, true, *array.GetBoolN(4))
}

func TestArrayGetString(t *testing.T) {
	mapper, _ := jogson.NewMapperFromString(jsonAnyArrayTest)
	array := mapper.AsArray

	s := array.GetString(0)
	assert.NoError(t, array.LastError)
	assert.Equal(t, "Jason", s)

	array.LastError = nil
	s = array.GetString(1)
	assert.NoError(t, array.LastError)
	assert.Equal(t, "15", s)

	array.LastError = nil
	s = array.GetString(2)
	assert.Error(t, array.LastError)
	assert.Equal(t, "", s)

	array.LastError = nil
	s = array.GetString(3)
	assert.NoError(t, array.LastError)
	assert.Equal(t, "1.81", s)

	array.LastError = nil
	s = array.GetString(4)
	assert.NoError(t, array.LastError)
	assert.Equal(t, "true", s)
}

func TestArrayGetStringFails(t *testing.T) {
	mapper, _ := jogson.NewMapperFromString(jsonAnyArrayTest)
	array := mapper.AsArray

	s := array.GetString(10)
	assert.ErrorIs(t, array.LastError, jogson.IndexOutOfRangeErr)
	assert.Equal(t, "", s)

	s = array.GetString(2)
	assert.ErrorIs(t, array.LastError, jogson.TypeConversionErr)
	assert.Equal(t, "", s)
}

func TestArrayGetInt(t *testing.T) {
	mapper, _ := jogson.NewMapperFromString(jsonIntArrayTest)
	array := mapper.AsArray

	i := array.GetInt(0)
	assert.NoError(t, array.LastError)
	assert.Equal(t, 0, i)

	i = array.GetInt(1)
	assert.NoError(t, array.LastError)
	assert.Equal(t, 15, i)

	i = array.GetInt(2)
	assert.NoError(t, array.LastError)
	assert.Equal(t, -54, i)

	i = array.GetInt(4)
	assert.NoError(t, array.LastError)
	assert.Equal(t, math.MaxInt, i)
}

func TestArrayGetIntFails(t *testing.T) {
	mapper, _ := jogson.NewMapperFromString(jsonAnyArrayTest)
	array := mapper.AsArray

	i := array.GetInt(10)
	assert.ErrorIs(t, array.LastError, jogson.IndexOutOfRangeErr)
	assert.Equal(t, 0, i)

	i = array.GetInt(0)
	assert.ErrorIs(t, array.LastError, jogson.TypeConversionErr)
	assert.Equal(t, 0, i)

	i = array.GetInt(2)
	assert.ErrorIs(t, array.LastError, jogson.TypeConversionErr)
	assert.Equal(t, 0, i)
}

func TestArrayGetFloat(t *testing.T) {
	mapper, _ := jogson.NewMapperFromString(jsonFloatArrayTest)
	array := mapper.AsArray

	f := array.GetFloat(0)
	assert.NoError(t, array.LastError)
	assert.Equal(t, 15.13, f)

	f = array.GetFloat(1)
	assert.NoError(t, array.LastError)
	assert.Equal(t, float64(2), f)

	f = array.GetFloat(2)
	assert.NoError(t, array.LastError)
	assert.Equal(t, 45.3984, f)

	f = array.GetFloat(3)
	assert.NoError(t, array.LastError)
	assert.Equal(t, -1.81, f)

	f = array.GetFloat(4)
	assert.NoError(t, array.LastError)
	assert.Equal(t, 9.223372036854776, f)
}

func TestArrayGetFloatFails(t *testing.T) {
	mapper, _ := jogson.NewMapperFromString(jsonAnyArrayTest)
	array := mapper.AsArray

	f := array.GetFloat(10)
	assert.ErrorIs(t, array.LastError, jogson.IndexOutOfRangeErr)
	assert.Equal(t, float64(0), f)

	f = array.GetFloat(0)
	assert.ErrorIs(t, array.LastError, jogson.TypeConversionErr)
	assert.Equal(t, float64(0), f)

	f = array.GetFloat(2)
	assert.ErrorIs(t, array.LastError, jogson.TypeConversionErr)
	assert.Equal(t, float64(0), f)
}

func TestArrayGetArray(t *testing.T) {
	mapper, _ := jogson.NewMapperFromString(json2DArrayTest)
	array := mapper.AsArray

	nestedArray := array.GetArray(0)
	assert.ElementsMatch(t, []int{1, 2}, nestedArray.AsIntArray())
}

func TestArrayGetArrayFails(t *testing.T) {
	mapper, _ := jogson.NewMapperFromString(json2DArrayTest)
	array := mapper.AsArray

	innerArr := array.GetArray(5)
	assert.ErrorIs(t, array.LastError, jogson.IndexOutOfRangeErr)
	assert.Equal(t, jogson.EmptyArray(), innerArr)

	innerArr = array.GetArray(2)
	assert.ErrorIs(t, array.LastError, jogson.TypeConversionErr)
	assert.Equal(t, jogson.EmptyArray(), innerArr)

	innerArr = array.GetArray(3)
	assert.ErrorIs(t, array.LastError, jogson.TypeConversionErr)
	assert.Equal(t, jogson.EmptyArray(), innerArr)
}

func TestArrayGetObject(t *testing.T) {
	mapper, _ := jogson.NewMapperFromString(jsonObjectArrayTest)
	array := mapper.AsArray

	obj := array.GetObject(1)
	assert.Equal(t, "Chris", obj.GetString("name"))
}

func TestArrayGetObjectFails(t *testing.T) {
	mapper, _ := jogson.NewMapperFromString(jsonArrayWithNullTest)
	array := mapper.AsArray

	obj := array.GetObject(10)
	assert.ErrorIs(t, array.LastError, jogson.IndexOutOfRangeErr)
	assert.True(t, obj.IsNull())

	obj = array.GetObject(2)
	assert.ErrorIs(t, array.LastError, jogson.TypeConversionErr)
	assert.True(t, obj.IsNull())

	obj = array.GetObject(3)
	assert.ErrorIs(t, array.LastError, jogson.TypeConversionErr)
	assert.True(t, obj.IsNull())
}

func TestArrayGetTime(t *testing.T) {
	mapper, err := jogson.NewArrayFromString(jsonArrayTimeTest)
	assert.NoError(t, err)
	actualTime1 := mapper.GetTime(0)
	assert.NoError(t, mapper.LastError)
	actualTime2 := mapper.GetTime(1)
	assert.NoError(t, mapper.LastError)
	actualTime3 := mapper.GetTime(2)
	assert.NoError(t, mapper.LastError)
	expectedTime1, _ := time.Parse(time.RFC3339, "2024-10-06T17:59:44Z")
	expectedTime2, _ := time.Parse(time.RFC3339, "2024-10-06T17:59:44+00:00")
	expectedTime3, _ := time.Parse(time.RFC850, "Sunday, 06-Oct-24 17:59:44 UTC")
	assert.Equal(t, expectedTime1, actualTime1)
	assert.Equal(t, expectedTime2, actualTime2)
	assert.Equal(t, expectedTime3, actualTime3)
}

func TestArrayIsNull(t *testing.T) {
	mapper, _ := jogson.NewMapperFromString(jsonObjectTest)
	object := mapper.AsObject
	nullArray := object.GetArray("address")
	assert.True(t, nullArray.IsNull())
}

func TestArrayIsEmpty(t *testing.T) {
	array, _ := jogson.NewArrayFromFile(jsonEmptyArrayTest)
	assert.True(t, array.IsEmpty())
	obj, _ := jogson.NewObjectFromString(jsonEmptyArrayTest)
	assert.True(t, array.IsEmpty())
	assert.True(t, obj.IsEmpty())
	assert.True(t, obj.IsNull())
}

func TestArrayPrintString(t *testing.T) {
	mapper, err := jogson.NewMapperFromString(jsonObjectArrayTest)
	assert.NoError(t, err)
	s := mapper.AsArray.String()
	assert.Equal(t, `[{"name":"Jason"},{"name":"Chris"}]`, s)
}

func TestArrayPrettyString(t *testing.T) {
	mapper, err := jogson.NewMapperFromString(jsonObjectArrayTest)
	assert.NoError(t, err)
	expectedArrayStr := "[\n  {\n    \"name\": \"Jason\"\n  },\n  {\n    \"name\": \"Chris\"\n  }\n]"
	assert.Equal(t, expectedArrayStr, mapper.AsArray.PrettyString())
}

func TestArrayFilter(t *testing.T) {
	mapper, err := jogson.NewMapperFromString(jsonObjectArrayTest)
	assert.NoError(t, err)
	filteredArr := mapper.AsArray.Filter(func(element jogson.JsonMapper) bool {
		return element.AsObject.GetString("name") == "Chris"
	})
	assert.Equal(t, 1, filteredArr.Length())
	assert.Equal(t, "Chris", filteredArr.Elements()[0].AsObject.GetString("name"))
}

func TestArrayFilterNull(t *testing.T) {
	mapper, err := jogson.NewMapperFromString(jsonAnyArrayTest)
	assert.NoError(t, err)
	filteredArr := mapper.AsArray.FilterNull()
	assert.Equal(t, 5, mapper.AsArray.Length())
	assert.Equal(t, 4, filteredArr.Length())
	filteredArr.ForEach(func(j jogson.JsonMapper) {
		assert.True(t, !j.IsNull)
	})
}

func TestArrayForEach(t *testing.T) {
	mapper, err := jogson.NewMapperFromString(jsonObjectArrayTest)
	assert.NoError(t, err)
	wasVisited := false
	mapper.AsArray.ForEach(func(mapper jogson.JsonMapper) {
		wasVisited = true
		assert.NotNil(t, mapper)
	})
	assert.True(t, wasVisited)
}

func TestArrayAnyAndAll(t *testing.T) {
	arrayWithNull, err := jogson.NewArrayFromString(jsonArrayWithNullTest)
	assert.NoError(t, err)

	arrayWithOnlyNull, err := jogson.NewArrayFromString(jsonArrayWithOnlyNullTest)
	assert.NoError(t, err)

	arrayWithoutNull, err := jogson.NewArrayFromString(jsonStringArrayTest)
	assert.NoError(t, err)

	arrayEmpty, err := jogson.NewArrayFromString(jsonEmptyArrayTest)
	assert.NoError(t, err)

	assert.True(t, arrayEmpty.Any())
	assert.True(t, arrayEmpty.All())

	assert.False(t, arrayWithOnlyNull.Any())
	assert.False(t, arrayWithOnlyNull.All())

	assert.True(t, arrayWithNull.Any())
	assert.False(t, arrayWithNull.All())

	assert.True(t, arrayWithoutNull.Any())
	assert.True(t, arrayWithoutNull.All())
}

func TestIndexOutOfBoundError(t *testing.T) {
	array := jogson.EmptyArray()
	array.AddInt(1)
	assert.Equal(t, 1, array.Length())
	assert.Equal(t, 0, array.GetInt(3))
	assert.Error(t, array.LastError)
	assert.ErrorIs(t, array.LastError, jogson.IndexOutOfRangeErr)
}
