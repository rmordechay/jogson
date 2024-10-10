package tests

import (
	"github.com/rmordechay/jsonmapper"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
	"time"
)

func TestParseJsonObjectFromString(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonObjectTest)
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.AsObject.String())

	assert.True(t, mapper.IsObject)
	assert.Contains(t, actual, `"age":15`)
	assert.Contains(t, actual, `"name":"Jason"`)
}

func TestParseJsonArrayFromString(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonObjectArrayTest)
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.AsArray.String())
	expected := removeWhiteSpaces(jsonObjectArrayTest)

	assert.True(t, mapper.IsArray)
	assert.Equal(t, expected, actual)
	assert.Equal(t, mapper.AsArray.Length(), 2)
}

func TestParseJsonArrayFromStringWithNulls(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonArrayWithNullTest)
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.AsArray.String())
	expected := removeWhiteSpaces(jsonArrayWithNullTest)

	assert.True(t, mapper.IsArray)
	assert.Equal(t, expected, actual)
	assert.Equal(t, 4, mapper.AsArray.Length())
}

func TestParseJsonObjectFromBytes(t *testing.T) {
	mapper, err := jsonmapper.FromBytes([]byte(jsonObjectTest))
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.AsObject.String())

	assert.True(t, mapper.IsObject)
	assert.Contains(t, actual, `"age":15`)
	assert.Contains(t, actual, `"name":"Jason"`)
}

func TestParseJsonArrayFromBytes(t *testing.T) {
	mapper, err := jsonmapper.FromBytes([]byte(jsonObjectArrayTest))
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.AsArray.String())
	expected := removeWhiteSpaces(jsonObjectArrayTest)

	assert.True(t, mapper.IsArray)
	assert.Equal(t, expected, actual)
	assert.Equal(t, mapper.AsArray.Length(), 2)
}

func TestParseJsonArrayFromStruct(t *testing.T) {
	testStruct := struct {
		Name string `json:"name"`
		Age  int    `json:"Age"`
	}{"John", 15}
	mapper, err := jsonmapper.FromStruct(testStruct)
	assert.NoError(t, err)
	assert.True(t, mapper.IsObject)
	assert.Equal(t, "John", mapper.AsObject.GetString("name"))
	assert.Equal(t, 15, mapper.AsObject.GetInt("Age"))
}

func TestParseJsonArrayFromStruct2(t *testing.T) {
	type childTest struct {
		Age     int
		IsFunny bool
	}

	type personTest struct {
		Name     string
		Age      int
		Height   float64
		IsFunny  bool
		Birthday time.Time
		Features []string
		Children map[string]childTest
	}
	child1 := childTest{Age: 17, IsFunny: false}
	child2 := childTest{Age: 23, IsFunny: true}
	children := make(map[string]childTest)
	children["Rachel"] = child1
	children["Sara"] = child2
	birthday, _ := time.Parse(time.DateOnly, "1981-05-30")
	person := personTest{
		Name:     "Chris",
		Age:      45,
		Height:   1.85,
		IsFunny:  true,
		Birthday: birthday,
		Features: []string{"tall", "blue eyes"},
		Children: children,
	}
	mapper, err := jsonmapper.FromStruct(person)
	assert.NoError(t, err)
	assert.NotNil(t, mapper)
	assert.Equal(t, 45, mapper.AsObject.GetInt("Age"))
	assert.Equal(t, "1981-05-30T00:00:00Z", mapper.AsObject.GetString("Birthday"))
	getTime, err := mapper.AsObject.GetTime("Birthday")
	assert.NoError(t, err)
	assert.Equal(t, birthday, getTime)
	assert.Equal(t, 1.85, mapper.AsObject.GetFloat("Height"))
	assert.Equal(t, true, mapper.AsObject.GetBool("IsFunny"))
}

func TestParseJsonObjectFromFile(t *testing.T) {
	path := "files/test_object.json"
	mapper, err := jsonmapper.FromFile(path)
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.AsObject.String())
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

	actual := removeWhiteSpaces(mapper.AsArray.String())
	fileExpected, err := os.ReadFile(path)
	expected := removeWhiteSpaces(string(fileExpected))

	assert.NoError(t, err)
	assert.True(t, mapper.IsArray)
	assert.Equal(t, expected, actual)
	assert.Equal(t, mapper.AsArray.Length(), 2)
}

func TestParseTime(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonTimeTest)
	assert.NoError(t, err)
	actualTime1, err := mapper.AsObject.GetTime("time1")
	assert.NoError(t, err)
	actualTime2, err := mapper.AsObject.GetTime("time2")
	assert.NoError(t, err)
	actualTime3, err := mapper.AsObject.GetTime("time3")
	assert.NoError(t, err)
	expectedTime1, _ := time.Parse(time.RFC3339, "2024-10-06T17:59:44Z")
	expectedTime2, _ := time.Parse(time.RFC3339, "2024-10-06T17:59:44+00:00")
	expectedTime3, _ := time.Parse(time.RFC850, "Sunday, 06-Oct-24 17:59:44 UTC")
	assert.Equal(t, expectedTime1, actualTime1)
	assert.Equal(t, expectedTime2, actualTime2)
	assert.Equal(t, expectedTime3, actualTime3)
}

func TestParseTimeInvalid(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonInvalidTimeTest)
	assert.NoError(t, err)
	for _, v := range mapper.AsObject.Elements() {
		_, err = v.AsTime()
		assert.Error(t, err)
	}
}

func TestExample(t *testing.T) {
	//docs.RunExample()
}

func removeWhiteSpaces(data string) string {
	s := strings.ReplaceAll(data, " ", "")
	s = strings.ReplaceAll(s, "\n", "")
	return s
}
