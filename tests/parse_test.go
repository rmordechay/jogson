package tests

import (
	"github.com/rmordechay/json-mapper"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
	"time"
)

func TestParseJsonObjectFromString(t *testing.T) {
	mapper, err := jsonmapper.GetMapperFromString(jsonObjectTest)
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.Object.String())
	expected := removeWhiteSpaces(jsonObjectTest)

	assert.True(t, mapper.IsObject)
	assert.Equal(t, expected, actual)
}

func TestParseJsonArrayFromString(t *testing.T) {
	mapper, err := jsonmapper.GetMapperFromString(jsonArrayTest)
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.Array.String())
	expected := removeWhiteSpaces(jsonArrayTest)

	assert.True(t, mapper.IsArray)
	assert.Equal(t, expected, actual)
	assert.Equal(t, mapper.Array.Length(), 2)
}

func TestParseJsonArrayFromStringWithNulls(t *testing.T) {
	mapper, err := jsonmapper.GetMapperFromString(jsonArrayWithNullTest)
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.Array.String())
	expected := removeWhiteSpaces(jsonArrayWithNullTest)

	assert.True(t, mapper.IsArray)
	assert.Equal(t, expected, actual)
	assert.Equal(t, mapper.Array.Length(), 2)
}

func TestParseJsonObjectFromBytes(t *testing.T) {
	mapper, err := jsonmapper.GetMapperFromBytes([]byte(jsonObjectTest))
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.Object.String())
	expected := removeWhiteSpaces(jsonObjectTest)

	assert.True(t, mapper.IsObject)
	assert.Equal(t, expected, actual)
}

func TestParseJsonArrayFromBytes(t *testing.T) {
	mapper, err := jsonmapper.GetMapperFromBytes([]byte(jsonArrayTest))
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.Array.String())
	expected := removeWhiteSpaces(jsonArrayTest)

	assert.True(t, mapper.IsArray)
	assert.Equal(t, expected, actual)
	assert.Equal(t, mapper.Array.Length(), 2)
}

func TestParseJsonObjectFromFile(t *testing.T) {
	path := "files/test_object.json"
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
	path := "files/test_array.json"
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

func TestParseTime(t *testing.T) {
	mapper, err := jsonmapper.GetMapperFromString(jsonTimeTest)
	assert.NoError(t, err)
	actualTime1 := mapper.Object.Get("time1").AsTime()
	actualTime2 := mapper.Object.Get("time2").AsTime()
	actualTime3 := mapper.Object.Get("time3").AsTime()
	expectedTime1, _ := time.Parse(time.RFC3339, "2024-10-06T17:59:44Z")
	expectedTime2, _ := time.Parse(time.RFC3339, "2024-10-06T17:59:44+00:00")
	expectedTime3, _ := time.Parse(time.RFC850, "Sunday, 06-Oct-24 17:59:44 UTC")
	assert.Equal(t, expectedTime1, actualTime1)
	assert.Equal(t, expectedTime2, actualTime2)
	assert.Equal(t, expectedTime3, actualTime3)
}

func removeWhiteSpaces(data string) string {
	s := strings.ReplaceAll(data, " ", "")
	s = strings.ReplaceAll(s, "\n", "")
	return s
}
