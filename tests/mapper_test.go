package tests

import (
	"github.com/rmordechay/jsonmapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
