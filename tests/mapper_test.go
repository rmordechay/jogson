package tests

import (
	"testing"

	"github.com/rmordechay/jogson"
	"github.com/stretchr/testify/assert"
)

func TestParseTimeInvalid(t *testing.T) {
	mapper, err := jogson.NewMapperFromString(jsonInvalidTimeTest)
	assert.NoError(t, err)
	for _, v := range mapper.AsObject.Elements() {
		_, err = v.AsTime()
		assert.Error(t, err)
	}
}

func TestMapperString(t *testing.T) {
	mapper, err := jogson.NewMapperFromString(jsonObjectTest)
	assert.NoError(t, err)
	expectedObj := `{"address":null,"age":15,"height":1.81,"is_funny":true,"name":"Jason"}`
	assert.Equal(t, expectedObj, mapper.String())

	mapper, err = jogson.NewMapperFromString(jsonObjectArrayTest)
	assert.NoError(t, err)
	expectedArray := `[{"name":"Jason"},{"name":"Chris"}]`
	assert.Equal(t, expectedArray, mapper.String())

	mapper, err = jogson.NewMapperFromString(jsonOnlyStringTest)
	assert.NoError(t, err)
	assert.Equal(t, "test", mapper.String())

	mapper, err = jogson.NewMapperFromString(jsonOnlyIntTest)
	assert.NoError(t, err)
	assert.Equal(t, "56", mapper.String())

	mapper, err = jogson.NewMapperFromString(jsonOnlyFloatTest)
	assert.NoError(t, err)
	assert.Equal(t, "1.2", mapper.String())

	mapper, err = jogson.NewMapperFromString(jsonOnlyBoolTest)
	assert.NoError(t, err)
	assert.Equal(t, "true", mapper.String())

	mapper, err = jogson.NewMapperFromString(jsonOnlyNullTest)
	assert.NoError(t, err)
	assert.True(t, mapper.IsNull)
}

func TestJsonInvalid(t *testing.T) {
	mapper, err := jogson.NewMapperFromString(jsonInvalidObjectTest)
	assert.Zero(t, mapper)
	assert.Error(t, err)

	obj, err := jogson.NewObjectFromString(jsonInvalidObjectTest)
	assert.Error(t, err)
	assert.True(t, obj.IsNull())
	assert.True(t, obj.IsEmpty())

	arr, err := jogson.NewArrayFromString(jsonInvalidArrayTest)
	assert.Empty(t, arr)
	assert.Error(t, err)
}

//func TestProcessObjects(t *testing.T) {
//	n := 1000
//	array, _ := generateJSONArray(n)
//	mapper, _ := jogson.NewMapperFromBuffer(strings.NewReader(array))
//	c := 0
//	var mu sync.Mutex
//	err := mapper.processObjects(10, func(o jogson.JsonObject) {
//		mu.Lock()
//		c++
//		mu.Unlock()
//	})
//	assert.NoError(t, err)
//	assert.Equal(t, n, c)
//}
//
//func TestProcessObjectsWithArgs(t *testing.T) {
//	n := 1000
//	array, _ := generateJSONArray(n)
//	mapper, _ := jogson.NewMapperFromBuffer(strings.NewReader(array))
//	c := 0
//	var mu sync.Mutex
//	err := mapper.processObjectsWithArgs(10, worker, &c, &mu)
//	assert.NoError(t, err)
//	assert.Equal(t, n, c)
//}

func TestMapperPrettyString(t *testing.T) {
	mapper, err := jogson.NewMapperFromString(jsonObjectTest)
	assert.NoError(t, err)
	expectedObjStr := "{\n  \"address\": null,\n  \"age\": 15,\n  \"height\": 1.81,\n  \"is_funny\": true,\n  \"name\": \"Jason\"\n}"
	assert.Equal(t, expectedObjStr, mapper.PrettyString())

	mapper, err = jogson.NewMapperFromString(jsonObjectArrayTest)
	assert.NoError(t, err)
	expectedArrayStr := "[\n  {\n    \"name\": \"Jason\"\n  },\n  {\n    \"name\": \"Chris\"\n  }\n]"
	assert.Equal(t, expectedArrayStr, mapper.PrettyString())

	mapper, err = jogson.NewMapperFromString(jsonOnlyStringTest)
	assert.NoError(t, err)
	assert.Equal(t, "test", mapper.PrettyString())

	mapper, err = jogson.NewMapperFromString(jsonOnlyIntTest)
	assert.NoError(t, err)
	assert.Equal(t, "56", mapper.PrettyString())

	mapper, err = jogson.NewMapperFromString(jsonOnlyFloatTest)
	assert.NoError(t, err)
	assert.Equal(t, "1.2", mapper.PrettyString())

	mapper, err = jogson.NewMapperFromString(jsonOnlyBoolTest)
	assert.NoError(t, err)
	assert.Equal(t, "true", mapper.PrettyString())

	mapper, err = jogson.NewMapperFromString(jsonOnlyNullTest)
	assert.NoError(t, err)
	assert.True(t, mapper.IsNull)
}

func TestExample(t *testing.T) {
	//sandbox.RunExample()
}

//// Function to generate a random JSON array with n elements
//func generateJSONArray(n int) (string, error) {
//	src := rand.NewSource(time.Now().UnixNano())
//	r := rand.New(src)
//	type Element struct {
//		MyString string `json:"my_string"`
//		MyNumber int    `json:"my_number"`
//		MyBool   bool   `json:"my_bool"`
//	}
//
//	elements := make([]Element, n)
//	for i := 0; i < n; i++ {
//		elements[i] = Element{
//			MyString: fmt.Sprintf("string_%d", i),
//			MyNumber: r.Intn(10000),
//			MyBool:   r.Intn(2) == 0,
//		}
//	}
//	jsonData, err := json.Marshal(elements)
//	if err != nil {
//		return "", err
//	}
//
//	return string(jsonData), nil
//}
//
//func worker(o jogson.JsonObject, args ...any) {
//	c, ok := args[0].(*int)
//	if !ok {
//		return
//	}
//	mu, ok := args[1].(*sync.Mutex)
//	if !ok {
//		return
//	}
//	mu.Lock()
//	*c += 1
//	mu.Unlock()
//}
