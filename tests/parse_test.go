package tests

import (
	"github.com/rmordechay/json-mapper"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestParseJsonObjectFromString(t *testing.T) {
	mapper, err := jsonmapper.GetMapperFromString(jsonStrObject)
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.Object.String())
	expected := removeWhiteSpaces(jsonStrObject)

	assert.True(t, mapper.IsObject)
	assert.Equal(t, expected, actual)
}

func TestParseJsonArrayFromString(t *testing.T) {
	mapper, err := jsonmapper.GetMapperFromString(jsonStrArray)
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.Array.String())
	expected := removeWhiteSpaces(jsonStrArray)

	assert.True(t, mapper.IsArray)
	assert.Equal(t, expected, actual)
	assert.Equal(t, mapper.Array.Length(), 2)
}

func TestParseJsonObjectFromBytes(t *testing.T) {
	mapper, err := jsonmapper.GetMapperFromBytes([]byte(jsonStrObject))
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.Object.String())
	expected := removeWhiteSpaces(jsonStrObject)

	assert.True(t, mapper.IsObject)
	assert.Equal(t, expected, actual)
}

func TestParseJsonArrayFromBytes(t *testing.T) {
	mapper, err := jsonmapper.GetMapperFromBytes([]byte(jsonStrArray))
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.Array.String())
	expected := removeWhiteSpaces(jsonStrArray)

	assert.True(t, mapper.IsArray)
	assert.Equal(t, expected, actual)
	assert.Equal(t, mapper.Array.Length(), 2)
}

func TestParseJsonObjectFromFile(t *testing.T) {
	path := "test_object.json"
	mapper, err := jsonmapper.GetMapperFromFile(path)
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.Object.String())
	fileExpected, err := os.ReadFile(path)
	expected := removeWhiteSpaces(string(fileExpected))

	assert.NoError(t, err)
	assert.True(t, mapper.IsObject)
	assert.Equal(t, expected, actual)
}

func TestParseJsonArrayFromFile(t *testing.T) {
	path := "test_array.json"
	mapper, err := jsonmapper.GetMapperFromFile(path)
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.Array.String())
	fileExpected, err := os.ReadFile(path)
	expected := removeWhiteSpaces(string(fileExpected))

	assert.NoError(t, err)
	assert.True(t, mapper.IsArray)
	assert.Equal(t, expected, actual)
	assert.Equal(t, mapper.Array.Length(), 2)
}

func removeWhiteSpaces(data string) string {
	s := strings.ReplaceAll(data, " ", "")
	s = strings.ReplaceAll(s, "\n", "")
	return s
}
