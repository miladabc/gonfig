package gonfig

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidInputError_Error(t *testing.T) {
	var i interface{}
	var s struct{}
	var sp *struct{}
	num := 2

	tests := []struct {
		input    interface{}
		expected string
	}{
		{
			input:    nil,
			expected: "<nil>",
		},
		{
			input:    i,
			expected: "<nil>",
		},
		{
			input:    s,
			expected: "non-pointer type",
		},
		{
			input:    sp,
			expected: "nil pointer",
		},
		{
			input:    &num,
			expected: "non-struct type",
		},
	}

	for _, tc := range tests {
		ie := InvalidInputError{
			Value: reflect.ValueOf(tc.input),
		}

		assert.EqualError(t, &ie, "gonfig: invalid input: "+tc.expected)
	}
}

func TestConfigErrors_Error(t *testing.T) {
	ce := ConfigErrors{
		errors.New("first"),
		errors.New("second"),
		errors.New("third"),
	}

	assert.EqualError(t, ce, "gonfig:\n  * first\n  * second\n  * third")
}
