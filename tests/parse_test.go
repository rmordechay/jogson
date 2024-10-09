package tests

import (
	"github.com/rmordechay/jsonmapper"
	"github.com/rmordechay/jsonmapper/docs"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
	"time"
)

func TestParseJsonObjectFromString(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonObjectTest)
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.Object.String())

	assert.True(t, mapper.IsObject)
	assert.Contains(t, actual, `"age":15`)
	assert.Contains(t, actual, `"name":"Jason"`)
}

func TestParseJsonArrayFromString(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonArrayTest)
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.Array.String())
	expected := removeWhiteSpaces(jsonArrayTest)

	assert.True(t, mapper.IsArray)
	assert.Equal(t, expected, actual)
	assert.Equal(t, mapper.Array.Length(), 2)
}

func TestParseJsonArrayFromStringWithNulls(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonArrayWithNullTest)
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.Array.String())
	expected := removeWhiteSpaces(jsonArrayWithNullTest)

	assert.True(t, mapper.IsArray)
	assert.Equal(t, expected, actual)
	assert.Equal(t, 3, mapper.Array.Length())
}

func TestParseJsonObjectFromBytes(t *testing.T) {
	mapper, err := jsonmapper.FromBytes([]byte(jsonObjectTest))
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.Object.String())

	assert.True(t, mapper.IsObject)
	assert.Contains(t, actual, `"age":15`)
	assert.Contains(t, actual, `"name":"Jason"`)
}

func TestParseJsonArrayFromBytes(t *testing.T) {
	mapper, err := jsonmapper.FromBytes([]byte(jsonArrayTest))
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.Array.String())
	expected := removeWhiteSpaces(jsonArrayTest)

	assert.True(t, mapper.IsArray)
	assert.Equal(t, expected, actual)
	assert.Equal(t, mapper.Array.Length(), 2)
}

func TestParseJsonArrayFromStruct(t *testing.T) {
	testStruct := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{"John", 15}
	mapper, err := jsonmapper.FromStruct(testStruct)
	assert.NoError(t, err)
	assert.True(t, mapper.IsObject)
	assert.Equal(t, "John", mapper.Object.GetString("name"))
	assert.Equal(t, 15, mapper.Object.GetInt("age"))
}

func TestParseJsonObjectFromFile(t *testing.T) {
	path := "files/test_object.json"
	mapper, err := jsonmapper.FromFile(path)
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
	mapper, err := jsonmapper.FromFile(path)
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
	mapper, err := jsonmapper.FromString(jsonTimeTest)
	assert.NoError(t, err)
	actualTime1, err := mapper.Object.GetTime("time1")
	assert.NoError(t, err)
	actualTime2, err := mapper.Object.GetTime("time2")
	assert.NoError(t, err)
	actualTime3, err := mapper.Object.GetTime("time3")
	assert.NoError(t, err)
	expectedTime1, _ := time.Parse(time.RFC3339, "2024-10-06T17:59:44Z")
	expectedTime2, _ := time.Parse(time.RFC3339, "2024-10-06T17:59:44+00:00")
	expectedTime3, _ := time.Parse(time.RFC850, "Sunday, 06-Oct-24 17:59:44 UTC")
	assert.Equal(t, expectedTime1, actualTime1)
	assert.Equal(t, expectedTime2, actualTime2)
	assert.Equal(t, expectedTime3, actualTime3)
}

func TestParseTimeInvalid(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonTimeTestInvalid)
	assert.NoError(t, err)
	for _, v := range mapper.Object.Elements() {
		_, err = v.AsTime()
		assert.Error(t, err)
	}
}

func TestExample(t *testing.T) {
	docs.RunExample()
}

func removeWhiteSpaces(data string) string {
	s := strings.ReplaceAll(data, " ", "")
	s = strings.ReplaceAll(s, "\n", "")
	return s
}
