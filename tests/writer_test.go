package tests

import (
	"fmt"
	jsonmapper "github.com/rmordechay/json-mapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWrite(t *testing.T) {
	obj := jsonmapper.CreateEmptyJsonObject()
	obj.AddKeyValue("name", "chris")
	obj.AddKeyValue("numbers", []int{1, 2, 4})
	elements := obj.Elements()
	assert.Equal(t, "chris", elements["name"].String())
	assert.True(t, elements["numbers"].IsArray)
	assert.Equal(t, 3, elements["numbers"].Array.Length)
	fmt.Println(obj)
}
