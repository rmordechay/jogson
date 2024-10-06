package tests

import (
	jsonmapper "github.com/rmordechay/json-mapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestObjectGetKeys(t *testing.T) {
	obj, err := jsonmapper.CreateMapperFromString(jsonObjectTest)
	assert.NoError(t, err)
	keys := obj.Object.Keys()
	assert.Equal(t, 2, len(keys))
	assert.Contains(t, keys, "name")
	assert.Contains(t, keys, "age")
}

func TestObjectGetValues(t *testing.T) {
	obj, err := jsonmapper.CreateMapperFromString(jsonObjectTest)
	assert.NoError(t, err)
	values := obj.Object.Values()
	assert.Equal(t, 2, len(values))
	assert.Equal(t, "Jason", values[0].AsString)
	assert.Equal(t, 15, values[1].AsInt)
}
