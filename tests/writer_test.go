package tests

import (
	jsonmapper "github.com/rmordechay/json-mapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWriteObjectString(t *testing.T) {
	obj := jsonmapper.CreateEmptyJsonObject()
	obj.AddKeyValue("name", "chris")
	objElements := obj.Elements()
	assert.Equal(t, "chris", objElements["name"].String())
}

func TestWriteObjectArrayInt(t *testing.T) {
	obj := jsonmapper.CreateEmptyJsonObject()
	obj.AddKeyValue("numbers", []int{1, 2, 4})
	objElements := obj.Elements()

	numberArray := objElements["numbers"].Array
	arrayElements := numberArray.Elements()

	assert.True(t, objElements["numbers"].IsArray)
	assert.Equal(t, 3, numberArray.Length())
	assert.Equal(t, 1, arrayElements[0].AsInt)
	assert.Equal(t, 4, arrayElements[2].AsInt)
}

func TestWriteObjectArrayFloat(t *testing.T) {
	obj := jsonmapper.CreateEmptyJsonObject()
	obj.AddKeyValue("numbers", []float64{1.5, 2.0, 4.2})
	objElements := obj.Elements()

	numberArray := objElements["numbers"].Array
	arrayElements := numberArray.Elements()

	assert.True(t, objElements["numbers"].IsArray)
	assert.Equal(t, 3, numberArray.Length())
	assert.Equal(t, 1.5, arrayElements[0].AsFloat)
	assert.Equal(t, 4.2, arrayElements[2].AsFloat)
}

func TestWriteObjectArrayStrings(t *testing.T) {
	obj := jsonmapper.CreateEmptyJsonObject()
	obj.AddKeyValue("strings", []string{"string1", "string2", "string4"})
	objElements := obj.Elements()

	stringArray := objElements["strings"].Array
	arrayElements := stringArray.Elements()

	assert.True(t, objElements["strings"].IsArray)
	assert.Equal(t, 3, stringArray.Length())
	assert.Equal(t, "string1", arrayElements[0].AsString)
	assert.Equal(t, "string4", arrayElements[2].AsString)
}
