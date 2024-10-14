package tests

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rmordechay/jogson"
	"github.com/stretchr/testify/assert"
)

func TestParseJsonObjectFromString(t *testing.T) {
	mapper, err := jogson.FromString(jsonObjectTest)
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.AsObject.String())

	assert.True(t, mapper.IsObject)
	assert.Contains(t, actual, `"age":15`)
	assert.Contains(t, actual, `"name":"Jason"`)

	obj, err := jogson.NewObjectFromString(jsonObjectTest)
	assert.NoError(t, err)
	assert.Equal(t, 15, obj.GetInt("age"))
	assert.Equal(t, "Jason", obj.GetString("name"))
}

func TestParseJsonArrayFromString(t *testing.T) {
	mapper, err := jogson.FromString(jsonObjectArrayTest)
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.AsArray.String())
	expected := removeWhiteSpaces(jsonObjectArrayTest)

	assert.True(t, mapper.IsArray)
	assert.Equal(t, expected, actual)
	assert.Equal(t, mapper.AsArray.Length(), 2)

	array, err := jogson.NewArrayFromString(jsonObjectArrayTest)
	assert.NoError(t, err)
	assert.Equal(t, 2, array.Length())
	assert.Equal(t, "Jason", array.GetObject(0).GetString("name"))
}

func TestParseJsonArrayFromStringWithNulls(t *testing.T) {
	mapper, err := jogson.FromString(jsonArrayWithNullTest)
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.AsArray.String())
	expected := removeWhiteSpaces(jsonArrayWithNullTest)

	assert.True(t, mapper.IsArray)
	assert.Equal(t, expected, actual)
	assert.Equal(t, 4, mapper.AsArray.Length())

	array, err := jogson.NewArrayFromString(jsonArrayWithNullTest)
	assert.NoError(t, err)
	assert.Equal(t, 4, array.Length())
	assert.Equal(t, "string", array.GetString(2))
}

func TestParseJsonObjectFromBytes(t *testing.T) {
	mapper, err := jogson.FromBytes([]byte(jsonObjectTest))
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.AsObject.String())

	assert.True(t, mapper.IsObject)
	assert.Contains(t, actual, `"age":15`)
	assert.Contains(t, actual, `"name":"Jason"`)

	obj, err := jogson.NewObjectFromBytes([]byte(jsonObjectTest))
	assert.NoError(t, err)
	assert.Equal(t, 15, obj.GetInt("age"))
	assert.Equal(t, "Jason", obj.GetString("name"))
}

func TestParseJsonArrayFromBytes(t *testing.T) {
	mapper, err := jogson.FromBytes([]byte(jsonObjectArrayTest))
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.AsArray.String())
	expected := removeWhiteSpaces(jsonObjectArrayTest)

	assert.True(t, mapper.IsArray)
	assert.Equal(t, expected, actual)
	assert.Equal(t, mapper.AsArray.Length(), 2)

	array, err := jogson.NewArrayFromBytes([]byte(jsonObjectArrayTest))
	assert.NoError(t, err)
	assert.Equal(t, 2, array.Length())
	assert.Equal(t, "Jason", array.GetObject(0).GetString("name"))
}

func TestParseJsonArrayFromStruct(t *testing.T) {
	testStruct := struct {
		Name string `json:"name"`
		Age  int    `json:"Age"`
	}{"John", 15}
	mapper, err := jogson.FromStruct(testStruct)
	assert.NoError(t, err)
	assert.True(t, mapper.IsObject)
	assert.Equal(t, "John", mapper.AsObject.GetString("name"))
	assert.Equal(t, 15, mapper.AsObject.GetInt("Age"))

	obj, err := jogson.NewObjectFromStruct(testStruct)
	assert.NoError(t, err)
	assert.Equal(t, 15, obj.GetInt("Age"))
	assert.Equal(t, "John", obj.GetString("name"))
}

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

func TestParseJsonArrayFromStruct2(t *testing.T) {
	person := getTestPerson()
	mapper, err := jogson.FromStruct(person)
	assert.NoError(t, err)
	getTime := mapper.AsObject.GetTime("Birthday")
	expectedBirthday, _ := time.Parse(time.DateOnly, "1981-05-30")

	assert.Equal(t, 45, mapper.AsObject.GetInt("Age"))
	assert.Equal(t, "1981-05-30T00:00:00Z", mapper.AsObject.GetString("Birthday"))
	assert.NoError(t, mapper.AsObject.LastError)
	assert.Equal(t, expectedBirthday, getTime)
	assert.Equal(t, 1.85, mapper.AsObject.GetFloat("Height"))
	assert.Equal(t, true, mapper.AsObject.GetBool("IsFunny"))

	obj, err := jogson.NewObjectFromStruct(person)
	assert.NoError(t, err)

	assert.Equal(t, 45, obj.GetInt("Age"))
	assert.Equal(t, "1981-05-30T00:00:00Z", obj.GetString("Birthday"))
	assert.NoError(t, obj.LastError)
	assert.Equal(t, expectedBirthday, getTime)
	assert.Equal(t, 1.85, obj.GetFloat("Height"))
	assert.Equal(t, true, obj.GetBool("IsFunny"))
}

func TestParseJsonObjectFromFile(t *testing.T) {
	path := "files/test_object.json"
	mapper, err := jogson.FromFile(path)
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.AsObject.String())
	fileExpected, err := os.ReadFile(path)
	expected := removeWhiteSpaces(string(fileExpected))

	assert.NoError(t, err)
	assert.True(t, mapper.IsObject)
	assert.Equal(t, expected, actual)
	assert.Equal(t, 43, mapper.AsObject.GetInt("age"))

	obj, err := jogson.NewObjectFromFile(path)
	assert.NoError(t, err)
	assert.True(t, mapper.IsObject)
	assert.Equal(t, expected, removeWhiteSpaces(obj.String()))
	assert.Equal(t, 43, obj.GetInt("age"))
}

func TestParseJsonArrayFromFile(t *testing.T) {
	path := "files/test_array.json"
	mapper, err := jogson.FromFile(path)
	assert.NoError(t, err)

	actual := removeWhiteSpaces(mapper.AsArray.String())
	fileExpected, err := os.ReadFile(path)
	expected := removeWhiteSpaces(string(fileExpected))

	assert.NoError(t, err)
	assert.True(t, mapper.IsArray)
	assert.Equal(t, expected, actual)
	assert.Equal(t, mapper.AsArray.Length(), 2)

	array, err := jogson.NewArrayFromFile(path)
	assert.NoError(t, err)
	assert.Equal(t, expected, removeWhiteSpaces(array.String()))
	assert.Equal(t, array.Length(), 2)
}

func TestParseOnlyString(t *testing.T) {
	mapper, err := jogson.FromString(jsonOnlyStringTest)
	assert.NoError(t, err)
	assert.True(t, mapper.IsString)
	assert.Equal(t, "test", mapper.AsString)
}

func TestParseOnlyInt(t *testing.T) {
	mapper, err := jogson.FromString(jsonOnlyIntTest)
	assert.NoError(t, err)
	assert.True(t, mapper.IsInt)
	assert.Equal(t, 56, mapper.AsInt)
}

func TestParseOnlyFloat(t *testing.T) {
	mapper, err := jogson.FromString(jsonOnlyFloatTest)
	assert.NoError(t, err)
	assert.True(t, mapper.IsFloat)
	assert.Equal(t, 1.2, mapper.AsFloat)
}

func TestParseOnlyBool(t *testing.T) {
	mapper, err := jogson.FromString(jsonOnlyBoolTest)
	assert.NoError(t, err)
	assert.True(t, mapper.IsBool)
	assert.Equal(t, true, mapper.AsBool)
}

func TestParseOnlyNull(t *testing.T) {
	mapper, err := jogson.FromString(jsonOnlyNullTest)
	assert.NoError(t, err)
	assert.True(t, mapper.IsNull)
}

func removeWhiteSpaces(data string) string {
	s := strings.ReplaceAll(data, " ", "")
	s = strings.ReplaceAll(s, "\n", "")
	return s
}

func getTestPerson() personTest {
	child1 := childTest{Age: 17, IsFunny: false}
	child2 := childTest{Age: 23, IsFunny: true}
	children := make(map[string]childTest)
	children["Rachel"] = child1
	children["Sara"] = child2
	birthday, _ := time.Parse(time.DateOnly, "1981-05-30")
	return personTest{
		Name:     "Chris",
		Age:      45,
		Height:   1.85,
		IsFunny:  true,
		Birthday: birthday,
		Features: []string{"tall", "blue eyes"},
		Children: children,
	}
}
