package tests

import (
	"encoding/json"
	"fmt"
	"github.com/rmordechay/jsonmapper"
	"github.com/rmordechay/jsonmapper/sandbox"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"
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

func TestProcessObjects(t *testing.T) {
	n := 1000
	array, _ := generateJSONArray(n)
	mapper, _ := jsonmapper.FromBuffer(strings.NewReader(array))
	c := 0
	var mu sync.Mutex
	err := mapper.ProcessObjects(10, func(o jsonmapper.JsonObject) {
		mu.Lock()
		c++
		mu.Unlock()
	})
	assert.NoError(t, err)
}

func TestExample(t *testing.T) {
	sandbox.RunExample()
}

// Function to generate a random JSON array with n elements
func generateJSONArray(n int) (string, error) {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	type Element struct {
		MyString string `json:"my_string"`
		MyNumber int    `json:"my_number"`
		MyBool   bool   `json:"my_bool"`
	}

	elements := make([]Element, n)
	for i := 0; i < n; i++ {
		elements[i] = Element{
			MyString: fmt.Sprintf("string_%d", i),
			MyNumber: r.Intn(10000),
			MyBool:   r.Intn(2) == 0,
		}
	}
	jsonData, err := json.Marshal(elements)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func humanReadableSize(bytes int) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
