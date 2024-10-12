package tests

import (
	"github.com/rmordechay/jsonmapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWriteObjectToObject(t *testing.T) {

}

func TestWriteArrayToObject(t *testing.T) {
	obj := jsonmapper.EmptyObject()
	arr := jsonmapper.EmptyArray()
	arr.AddInt(1)
	arr.AddInt(4)
	arr.AddInt(6)
	obj.AddJsonArray("children", arr)

	assert.NoError(t, obj.LastError)
	assert.True(t, obj.Has("children"))
	array := obj.GetArray("children")
	assert.NoError(t, obj.LastError)
	assert.Equal(t, 3, array.Length())
	assert.Equal(t, 4, array.GetInt(1))
}

func TestWriteStringToObject(t *testing.T) {
	obj := jsonmapper.EmptyObject()
	obj.AddString("name", "chris")
	assert.Equal(t, "chris", obj.GetString("name"))
}

func TestWriteIntToObject(t *testing.T) {
	obj := jsonmapper.EmptyObject()
	obj.AddInt("int", 2)
	objElements := obj.Elements()
	assert.Equal(t, 2, objElements["int"].AsInt)
}

func TestWriteArrayStringToObject(t *testing.T) {
	obj := jsonmapper.EmptyObject()
	obj.AddStringArray("strings", []string{"string1", "string2", "string4"})
	objElements := obj.Elements()

	stringArray := objElements["strings"].AsArray
	arrayElements := stringArray.Elements()

	assert.True(t, objElements["strings"].IsArray)
	assert.Equal(t, 3, stringArray.Length())
	assert.Equal(t, "string1", arrayElements[0].AsString)
	assert.Equal(t, "string4", arrayElements[2].AsString)
}

func TestWriteFloatToObject(t *testing.T) {
	obj := jsonmapper.EmptyObject()
	obj.AddFloat("float", 2.5)
	objElements := obj.Elements()
	assert.Equal(t, 2.5, objElements["float"].AsFloat)
}

func TestWriteBoolToObject(t *testing.T) {
	obj := jsonmapper.EmptyObject()
	obj.AddBool("bool_value", true)
	assert.Equal(t, true, obj.GetBool("bool_value"))
}

func TestWriteArrayStringsToObject(t *testing.T) {
	obj := jsonmapper.EmptyObject()
	obj.AddStringArray("strings", []string{"string1", "string2", "string4"})
	objElements := obj.Elements()

	stringArray := objElements["strings"].AsArray
	arrayElements := stringArray.Elements()

	assert.True(t, objElements["strings"].IsArray)
	assert.Equal(t, 3, stringArray.Length())
	assert.Equal(t, "string1", arrayElements[0].AsString)
	assert.Equal(t, "string4", arrayElements[2].AsString)
}

func TestWriteObjectArrayInt(t *testing.T) {
	obj := jsonmapper.EmptyObject()
	obj.AddIntArray("numbers", []int{1, 2, 4})
	objElements := obj.Elements()

	numberArray := objElements["numbers"].AsArray
	arrayElements := numberArray.Elements()

	assert.True(t, objElements["numbers"].IsArray)
	assert.Equal(t, 3, numberArray.Length())
	assert.Equal(t, 1, arrayElements[0].AsInt)
	assert.Equal(t, 4, arrayElements[2].AsInt)
}

func TestWriteArrayFloatToObject(t *testing.T) {
	obj := jsonmapper.EmptyObject()
	obj.AddFloatArray("numbers", []float64{1.5, 2.0, 4.2})
	objElements := obj.Elements()

	numberArray := objElements["numbers"].AsArray
	arrayElements := numberArray.Elements()

	assert.True(t, objElements["numbers"].IsArray)
	assert.Equal(t, 3, numberArray.Length())
	assert.Equal(t, 1.5, arrayElements[0].AsFloat)
	assert.Equal(t, 4.2, arrayElements[2].AsFloat)
}

func TestWriteFloatsArray(t *testing.T) {
	arr := jsonmapper.EmptyArray()
	arr.AddInt(1)
	arr.AddInt(4)
	arr.AddInt(6)
	assert.Equal(t, 3, arr.Length())
	assert.Equal(t, 1, arr.GetInt(0))
	assert.Equal(t, 4, arr.GetInt(1))
	assert.Equal(t, 6, arr.GetInt(2))
}
