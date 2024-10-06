package jsonmapper

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestParseJsonObjectFromString(t *testing.T) {
	assert.Empty(t, []int{})

}

func TestParseJsonArrayFromString(t *testing.T) {

}

func TestParseJsonObjectFromBytes(t *testing.T) {

}

func TestParseJsonArrayFromBytes(t *testing.T) {

}
func TestParseJsonObjectFromFile(t *testing.T) {

}

func TestParseJsonArrayFromFile(t *testing.T) {

}

func TestMapper(t *testing.T) {
	mapper, err := GetMapperFromFile("test.json")
	if err != nil {
		log.Fatal(err)
	}
	find := mapper.Object.Get("field1").Array
	get := find.Get(0)
	fmt.Println(get)
}
