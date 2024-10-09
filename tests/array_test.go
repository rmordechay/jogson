package tests

import (
	"github.com/rmordechay/jsonmapper"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestArrayFilter(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonObjectArrayTest)
	assert.NoError(t, err)
	filteredArr := mapper.Array.Filter(func(element jsonmapper.Json) bool {
		return element.Object.GetString("name") == "Chris"
	})
	assert.Equal(t, 1, filteredArr.Length())
	assert.Equal(t, "Chris", filteredArr.Elements()[0].Object.GetString("name"))
}

func TestArrayForEach(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonObjectArrayTest)
	assert.NoError(t, err)
	wasVisited := false
	mapper.Array.ForEach(func(mapper jsonmapper.Json) {
		wasVisited = true
		assert.NotNil(t, mapper)
	})
	assert.True(t, wasVisited)
}

func TestIndexOutOfBoundError(t *testing.T) {
	arr := jsonmapper.NewArray()
	arr.AddValue(1)
	assert.Equal(t, 1, arr.Length())
	assert.Equal(t, 0, arr.GetInt(3))
	assert.Error(t, arr.LastError)
	assert.Equal(t, arr.LastError.Error(), "index out of range [3] with length 1")
}

func TestArrayAsStringArray(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonStringArrayTest)
	array := mapper.Array
	assert.ElementsMatch(t, []string{"Jason", "Chris", "Rachel"}, array.AsStringArray())
}

func TestArrayAsIntArray(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonIntArrayTest)
	array := mapper.Array
	assert.ElementsMatch(t, []int{0, 15, -54, -346, 9223372036854775807}, array.AsIntArray())
}

func TestArrayAsFloatArray(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonFloatArrayTest)
	array := mapper.Array
	assert.ElementsMatch(t, []float64{15.13, 2, 45.3984, -1.81, 9.223372036854776}, array.AsFloatArray())
}

func TestArrayAs2DArray(t *testing.T) {
	mapper, _ := jsonmapper.FromString(json2DArrayTest)
	array := mapper.Array
	assert.ElementsMatch(t, []int{1, 2}, array.As2DArray()[0].AsIntArray())
	assert.ElementsMatch(t, []int{3, 4}, array.As2DArray()[1].AsIntArray())
}

func TestArrayAsObjectArray(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonObjectArrayTest)
	array := mapper.Array
	assert.Equal(t, "Jason", array.AsObjectArray()[0].GetString("name"))
	assert.Equal(t, "Chris", array.AsObjectArray()[1].GetString("name"))
}

func TestArrayGetString(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonAnyArrayTest)
	array := mapper.Array

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
	mapper, _ := jsonmapper.FromString(jsonAnyArrayTest)
	array := mapper.Array

	s := array.GetString(10)
	assert.Equal(t, "index out of range [10] with length 5", array.LastError.Error())
	assert.Equal(t, "", s)

	s = array.GetString(2)
	assert.Equal(t, "value is null and could not be converted to string", array.LastError.Error())
	assert.Equal(t, "", s)
}

func TestArrayGetInt(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonIntArrayTest)
	array := mapper.Array

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
	mapper, _ := jsonmapper.FromString(jsonAnyArrayTest)
	array := mapper.Array

	i := array.GetInt(10)
	assert.Equal(t, "index out of range [10] with length 5", array.LastError.Error())
	assert.Equal(t, 0, i)

	i = array.GetInt(0)
	assert.Equal(t, "the type 'string' could not be converted to int", array.LastError.Error())
	assert.Equal(t, 0, i)

	i = array.GetInt(2)
	assert.Equal(t, "value is null and could not be converted to int", array.LastError.Error())
	assert.Equal(t, 0, i)
}

func TestArrayGetFloat(t *testing.T) {
	mapper, _ := jsonmapper.FromString(jsonFloatArrayTest)
	array := mapper.Array

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
	mapper, _ := jsonmapper.FromString(jsonAnyArrayTest)
	array := mapper.Array

	f := array.GetFloat(10)
	assert.Equal(t, "index out of range [10] with length 5", array.LastError.Error())
	assert.Equal(t, float64(0), f)

	f = array.GetFloat(0)
	assert.Equal(t, "the type 'string' could not be converted to float64", array.LastError.Error())
	assert.Equal(t, float64(0), f)

	f = array.GetFloat(2)
	assert.Equal(t, "value is null and could not be converted to float64", array.LastError.Error())
	assert.Equal(t, float64(0), f)
}

func TestArrayGetArray(t *testing.T) {

}

func TestArrayGetArrayFails(t *testing.T) {

}

func TestArrayGetObject(t *testing.T) {

}

func TestArrayGetObjectFails(t *testing.T) {

}
