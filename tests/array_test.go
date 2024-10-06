package tests

import (
	jsonmapper "github.com/rmordechay/json-mapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArrayFilter(t *testing.T) {
	arr, err := jsonmapper.GetMapperFromString(jsonStrArray)
	assert.NoError(t, err)
	arr.Array.Filter(func(element jsonmapper.Mapper) bool {
		return element.IsArray
	})
}
