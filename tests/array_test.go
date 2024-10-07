package tests

import (
	"github.com/rmordechay/jsonmapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArrayFilter(t *testing.T) {
	arr, err := jsonmapper.FromString(jsonArrayTest)
	assert.NoError(t, err)
	isNameChris := func(element jsonmapper.Mapper) bool {
		return element.Object.Get("name").AsString == "Chris"
	}
	filteredArr := jsonmapper.Filter(arr.Array, isNameChris)
	assert.Equal(t, 1, len(filteredArr))
	assert.Equal(t, "Chris", filteredArr[0].Object.Get("name").AsString)
}

func TestArrayMapString(t *testing.T) {
	arr, err := jsonmapper.FromString(jsonArrayTest)
	assert.NoError(t, err)
	getNames := func(element jsonmapper.Mapper) string {
		return element.Object.Get("name").AsString
	}
	mappedArr := jsonmapper.Map(arr.Array, getNames)
	assert.Equal(t, 2, len(mappedArr))
	assert.Contains(t, mappedArr, "Jason")
	assert.Contains(t, mappedArr, "Chris")
}

func TestArrayMapObject(t *testing.T) {
	arr, err := jsonmapper.FromString(jsonArrayTest)
	assert.NoError(t, err)
	getNames := func(element jsonmapper.Mapper) jsonmapper.JsonObject {
		return element.Object
	}
	var mappedArr = jsonmapper.Map(arr.Array, getNames)
	assert.Equal(t, 2, len(mappedArr))
	expected1 := mappedArr[0].Get("name").AsString
	expected2 := mappedArr[1].Get("name").AsString
	assert.Equal(t, expected1, "Jason")
	assert.Equal(t, expected2, "Chris")
}

func TestArrayForEach(t *testing.T) {
	arr, err := jsonmapper.FromString(jsonArrayTest)
	assert.NoError(t, err)
	wasVisited := false
	jsonmapper.ForEach(arr.Array, func(element jsonmapper.Mapper) {
		wasVisited = true
		assert.NotNil(t, element)
	})
	assert.True(t, wasVisited)
}
