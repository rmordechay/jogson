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
	filteredArr := arr.Array.Filter(isNameChris)
	assert.Equal(t, 1, filteredArr.Length())
	assert.Equal(t, "Chris", filteredArr.Elements()[0].Object.Get("name").AsString)
}

func TestArrayForEach(t *testing.T) {
	arr, err := jsonmapper.FromString(jsonArrayTest)
	assert.NoError(t, err)
	wasVisited := false
	arr.Array.ForEach(func(mapper jsonmapper.Mapper) {
		wasVisited = true
		assert.NotNil(t, mapper)
	})
	assert.True(t, wasVisited)
}
