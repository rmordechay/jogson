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

	obj, err := jsonmapper.NewObjectFromString(jsonObjectTest)
	assert.NoError(t, err)
	assert.Equal(t, 15, obj.GetInt("age"))
	assert.Equal(t, "Jason", obj.GetString("name"))
}

func TestParseJsonArrayFromString(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonObjectArrayTest)
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.AsArray.String())
	expected := removeWhiteSpaces(jsonObjectArrayTest)

	assert.True(t, mapper.IsArray)
	assert.Equal(t, expected, actual)
	assert.Equal(t, mapper.AsArray.Length(), 2)

	array, err := jsonmapper.NewArrayFromString(jsonObjectArrayTest)
	assert.NoError(t, err)
	assert.Equal(t, 2, array.Length())
	assert.Equal(t, "Jason", array.GetObject(0).GetString("name"))
}

func TestParseJsonArrayFromStringWithNulls(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonArrayWithNullTest)
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.AsArray.String())
	expected := removeWhiteSpaces(jsonArrayWithNullTest)

	assert.True(t, mapper.IsArray)
	assert.Equal(t, expected, actual)
	assert.Equal(t, 4, mapper.AsArray.Length())

	array, err := jsonmapper.NewArrayFromString(jsonArrayWithNullTest)
	assert.NoError(t, err)
	assert.Equal(t, 4, array.Length())
	assert.Equal(t, "string", array.GetString(2))
}

func TestParseJsonObjectFromBytes(t *testing.T) {
	mapper, err := jsonmapper.FromBytes([]byte(jsonObjectTest))
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.AsObject.String())

	assert.True(t, mapper.IsObject)
	assert.Contains(t, actual, `"age":15`)
	assert.Contains(t, actual, `"name":"Jason"`)

	obj, err := jsonmapper.NewObjectFromBytes([]byte(jsonObjectTest))
	assert.NoError(t, err)
	assert.Equal(t, 15, obj.GetInt("age"))
	assert.Equal(t, "Jason", obj.GetString("name"))
}

func TestParseJsonArrayFromBytes(t *testing.T) {
	mapper, err := jsonmapper.FromBytes([]byte(jsonObjectArrayTest))
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.AsArray.String())
	expected := removeWhiteSpaces(jsonObjectArrayTest)

	assert.True(t, mapper.IsArray)
	assert.Equal(t, expected, actual)
	assert.Equal(t, mapper.AsArray.Length(), 2)

	array, err := jsonmapper.NewArrayFromBytes([]byte(jsonObjectArrayTest))
	assert.NoError(t, err)
	assert.Equal(t, 2, array.Length())
	assert.Equal(t, "Jason", array.GetObject(0).GetString("name"))
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

	obj, err := jsonmapper.NewObjectFromStruct(testStruct)
	assert.NoError(t, err)
	assert.Equal(t, 15, obj.GetInt("Age"))
	assert.Equal(t, "John", obj.GetString("name"))
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
	getTime := mapper.AsObject.GetTime("Birthday")

	assert.Equal(t, 45, mapper.AsObject.GetInt("Age"))
	assert.Equal(t, "1981-05-30T00:00:00Z", mapper.AsObject.GetString("Birthday"))
	assert.NoError(t, mapper.AsObject.LastError)
	assert.Equal(t, birthday, getTime)
	assert.Equal(t, 1.85, mapper.AsObject.GetFloat("Height"))
	assert.Equal(t, true, mapper.AsObject.GetBool("IsFunny"))

	obj, err := jsonmapper.NewObjectFromStruct(person)
	assert.NoError(t, err)

	assert.Equal(t, 45, obj.GetInt("Age"))
	assert.Equal(t, "1981-05-30T00:00:00Z", obj.GetString("Birthday"))
	assert.NoError(t, obj.LastError)
	assert.Equal(t, birthday, getTime)
	assert.Equal(t, 1.85, obj.GetFloat("Height"))
	assert.Equal(t, true, obj.GetBool("IsFunny"))
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
	assert.Equal(t, 43, mapper.AsObject.GetInt("age"))

	obj, err := jsonmapper.NewObjectFromFile(path)
	assert.NoError(t, err)
	assert.True(t, mapper.IsObject)
	assert.Equal(t, expected, removeWhiteSpaces(obj.String()))
	assert.Equal(t, 43, obj.GetInt("age"))
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

	array, err := jsonmapper.NewArrayFromFile(path)
	assert.NoError(t, err)
	assert.Equal(t, expected, removeWhiteSpaces(array.String()))
	assert.Equal(t, array.Length(), 2)
}

func TestParseOnlyString(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonOnlyStringTest)
	assert.NoError(t, err)
	assert.True(t, mapper.IsString)
	assert.Equal(t, "test", mapper.AsString)
}

func TestParseOnlyInt(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonOnlyIntTest)
	assert.NoError(t, err)
	assert.True(t, mapper.IsInt)
	assert.Equal(t, 56, mapper.AsInt)
}

func TestParseOnlyFloat(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonOnlyFloatTest)
	assert.NoError(t, err)
	assert.True(t, mapper.IsFloat)
	assert.Equal(t, 1.2, mapper.AsFloat)
}

func TestParseOnlyBool(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonOnlyBoolTest)
	assert.NoError(t, err)
	assert.True(t, mapper.IsBool)
	assert.Equal(t, true, mapper.AsBool)
}

func TestParseOnlyNull(t *testing.T) {
	mapper, err := jsonmapper.FromString(jsonOnlyNullTest)
	assert.NoError(t, err)
	assert.True(t, mapper.IsNull)
}

func removeWhiteSpaces(data string) string {
	s := strings.ReplaceAll(data, " ", "")
	s = strings.ReplaceAll(s, "\n", "")
	return s
}
