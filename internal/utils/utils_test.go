package utils

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOSVar(t *testing.T) {
	assert := assert.New(t)
	os.Setenv("TESTVAR", "TESTVALUE")

	testvalue := "TESTVALUE"
	testvar := GetOSVar("TESTVAR")
	assert.Equal(testvalue, testvar, "should be equal")

	testvalue2 := ""
	testvar2 := GetOSVar("TESTVAR2")
	assert.Equal(testvalue2, testvar2, "should be equal")

}

func TestRespError(t *testing.T) {
	err := errors.New("test")
	assert.Equal(t, errors.New("there was an error during the call execution: test\n"), RespError(err))
}
