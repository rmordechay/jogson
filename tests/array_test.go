package tests

import (
	"github.com/rmordechay/jsonmapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArrayFilter(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonArrayTest)
	assert.NoError(t, err)
	isNameChris := func(element jsonmapper.Mapper) bool {
		return element.Object.Get("name").AsString == "Chris"
	}
	filteredArr := mapper.Array.Filter(isNameChris)
	assert.Equal(t, 1, filteredArr.Length())
	assert.Equal(t, "Chris", filteredArr.Elements()[0].Object.Get("name").AsString)
}

func TestArrayForEach(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonArrayTest)
	assert.NoError(t, err)
	wasVisited := false
	mapper.Array.ForEach(func(mapper jsonmapper.Mapper) {
		wasVisited = true
		assert.NotNil(t, mapper)
	})
	assert.True(t, wasVisited)
}

func TestIndexOutOfBoundError(t *testing.T) {
	arr := jsonmapper.NewArray()
	arr.AddValue(1)
	assert.Equal(t, 1, arr.Length())
	assert.Equal(t, 0, arr.Get(3).AsInt)
	assert.Error(t, arr.Err)
	assert.Equal(t, arr.Err.Error(), "index out of range [3] with length 1")
}
