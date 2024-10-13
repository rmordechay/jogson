package tests

import (
	"github.com/rmordechay/jsonmapper"
	"github.com/rmordechay/jsonmapper/sandbox"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseTimeInvalid(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonInvalidTimeTest)
	assert.NoError(t, err)
	for _, v := range mapper.AsObject.Elements() {
		_, err = v.AsTime()
		assert.Error(t, err)
	}
}

func TestMapperString(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonObjectTest)
	assert.NoError(t, err)
	expectedObj := `{"address":null,"age":15,"height":1.81,"is_funny":true,"name":"Jason"}`
	assert.Equal(t, expectedObj, mapper.String())

	mapper, err = jsonmapper.FromString(jsonObjectArrayTest)
	assert.NoError(t, err)
	expectedArray := `[{"name":"Jason"},{"name":"Chris"}]`
	assert.Equal(t, expectedArray, mapper.String())
}

func TestMapperPrettyString(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonObjectTest)
	assert.NoError(t, err)
	expectedObjStr := "{\n  \"address\": null,\n  \"age\": 15,\n  \"height\": 1.81,\n  \"is_funny\": true,\n  \"name\": \"Jason\"\n}"
	assert.Equal(t, expectedObjStr, mapper.PrettyString())

	mapper, err = jsonmapper.FromString(jsonObjectArrayTest)
	assert.NoError(t, err)
	expectedArrayStr := "[\n  {\n    \"name\": \"Jason\"\n  },\n  {\n    \"name\": \"Chris\"\n  }\n]"
	assert.Equal(t, expectedArrayStr, mapper.PrettyString())
}

func TestExample(t *testing.T) {
	sandbox.RunExample()
}
