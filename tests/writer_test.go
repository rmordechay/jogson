package tests

import (
	"github.com/rmordechay/jsonmapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWriteStringToObject(t *testing.T) {
	obj := jsonmapper.NewObject()
	obj.AddKeyValue("name", "chris")
	objElements := obj.Elements()
	assert.Equal(t, "chris", objElements["name"].String())
}

func TestWriteIntToObject(t *testing.T) {
	obj := jsonmapper.NewObject()
	obj.AddKeyValue("int", 2)
	objElements := obj.Elements()
	assert.Equal(t, 2, objElements["int"].AsInt)
}

func TestWriteArrayStringToObject(t *testing.T) {
	obj := jsonmapper.NewObject()
	obj.AddKeyValue("strings", []string{"string1", "string2", "string4"})
	objElements := obj.Elements()

	stringArray := objElements["strings"].Array
	arrayElements := stringArray.Elements()

	assert.True(t, objElements["strings"].IsArray)
	assert.Equal(t, 3, stringArray.Length())
	assert.Equal(t, "string1", arrayElements[0].AsString)
	assert.Equal(t, "string4", arrayElements[2].AsString)
}

func TestWriteFloatToObject(t *testing.T) {
	obj := jsonmapper.NewObject()
	obj.AddKeyValue("float", 2.5)
	objElements := obj.Elements()
	assert.Equal(t, 2.5, objElements["float"].AsFloat)
}

func TestWriteBoolToObject(t *testing.T) {
	obj := jsonmapper.NewObject()
	obj.AddKeyValue("bool", true)
	objElements := obj.Elements()
	assert.Equal(t, true, objElements["bool"].AsBool)
}

func TestWriteArrayStringsToObject(t *testing.T) {
	obj := jsonmapper.NewObject()
	obj.AddKeyValue("strings", []string{"string1", "string2", "string4"})
	objElements := obj.Elements()

	stringArray := objElements["strings"].Array
	arrayElements := stringArray.Elements()

	assert.True(t, objElements["strings"].IsArray)
	assert.Equal(t, 3, stringArray.Length())
	assert.Equal(t, "string1", arrayElements[0].AsString)
	assert.Equal(t, "string4", arrayElements[2].AsString)
}

func TestWriteObjectArrayInt(t *testing.T) {
	obj := jsonmapper.NewObject()
	obj.AddKeyValue("numbers", []int{1, 2, 4})
	objElements := obj.Elements()

	numberArray := objElements["numbers"].Array
	arrayElements := numberArray.Elements()

	assert.True(t, objElements["numbers"].IsArray)
	assert.Equal(t, 3, numberArray.Length())
	assert.Equal(t, 1, arrayElements[0].AsInt)
	assert.Equal(t, 4, arrayElements[2].AsInt)
}

func TestWriteArrayFloatToObject(t *testing.T) {
	obj := jsonmapper.NewObject()
	obj.AddKeyValue("numbers", []float64{1.5, 2.0, 4.2})
	objElements := obj.Elements()

	numberArray := objElements["numbers"].Array
	arrayElements := numberArray.Elements()

	assert.True(t, objElements["numbers"].IsArray)
	assert.Equal(t, 3, numberArray.Length())
	assert.Equal(t, 1.5, arrayElements[0].AsFloat)
	assert.Equal(t, 4.2, arrayElements[2].AsFloat)
}

func TestWriteFloatsArray(t *testing.T) {
	arr := jsonmapper.NewArray()
	arr = arr.AddValue(1)
	arr = arr.AddValue(4)
	arr = arr.AddValue(6)
	assert.Equal(t, 3, arr.Length())
	assert.Equal(t, 1, arr.Get(0).AsInt)
	assert.Equal(t, 4, arr.Get(1).AsInt)
	assert.Equal(t, 6, arr.Get(2).AsInt)
	assert.Equal(t, 0, arr.Get(3).AsInt)
}
