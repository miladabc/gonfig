package gonfig

import (
	"errors"
	"reflect"
	"strings"
)

var (
	// ErrUnsupportedType indicates unsupported struct field
	ErrUnsupportedType = errors.New("unsupported type")

	// ErrUnsupportedFileExt indicated unsupported file format
	// Only ".json", ".yml", ".yaml" and ".env" file types are supported
	ErrUnsupportedFileExt = errors.New("unsupported file extension")

	// ErrUnSettableField indicated unexported struct field
	ErrUnSettableField = errors.New("unSettable field")

	// ErrKeyNotFound is returned when no value found with specified key
	ErrKeyNotFound = errors.New("key not found")

	// ErrRequiredField indicates that Field is required but no value is provided
	ErrRequiredField = errors.New("field is required")

	// ErrParsing is returned in case of bad value
	ErrParsing = errors.New("failed parsing")

	// ErrValueOverflow indicates value overflow
	ErrValueOverflow = errors.New("value overflow")
)

const (
	unsupportedTypeErrFormat    = `%w: %v`
	badFieldErrFormat           = `bad field "%v": %w`
	unsupportedFileExtErrFormat = `%w: %v`
	unSettableFieldErrFormat    = `%w: %v`
	decodeFailedErrFormat       = `failed to decode: %w`
	requiredFieldErrFormat      = `%w: no value found for "%v"`
	parseErrFormat              = `%w at "%v": %v`
	overflowErrFormat           = `%w: "%v" overflows type "%v" at "%v"`
)

// An InvalidInputError describes an invalid argument passed to Into function
// The argument must be a non-nil struct pointer
type InvalidInputError struct {
	Value reflect.Value
}

func (e *InvalidInputError) Error() string {
	msg := "gonfig: invalid input: "
	var t reflect.Type

	if e.Value.IsValid() {
		t = e.Value.Type()
	}

	if t == nil {
		msg += "<nil>"
	} else if t.Kind() != reflect.Ptr {
		msg += "non-pointer type"
	} else if e.Value.IsNil() {
		msg += "nil pointer"
	} else {
		msg += "non-struct type"
	}

	return msg
}

// ConfigErrors is collection of errors during populating the input struct
type ConfigErrors []error

func (ce ConfigErrors) Error() string {
	if len(ce) == 0 {
		return ""
	}

	msg := "gonfig:\n"
	for i := range ce {
		msg += "  * " + ce[i].Error() + "\n"
	}

	return strings.TrimSpace(msg)
}
