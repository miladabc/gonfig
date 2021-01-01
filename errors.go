package gonfig

import (
	"errors"
	"reflect"
	"strings"
)

var (
	// Can not handle specified type
	ErrUnsupportedType = errors.New("unsupported type")

	// Only ".json", ".yml", ".yaml" and ".env" file types are supported
	ErrUnsupportedFileExt = errors.New("unsupported file extension")

	// Provider could not find value with specified key
	ErrKeyNotFound = errors.New("key not found")

	// Field is required but no value provided
	ErrRequiredField = errors.New("field is required")

	// Could not parse the string value
	ErrParsing = errors.New("failed parsing")

	// Value overflows type
	ErrValueOverflow = errors.New("value overflow")
)

const (
	unsupportedTypeErrFormat        = `%w: cannot handle type "%v" at "%v"`
	unsupportedFileExtErrFormat     = `%w: %v`
	decodeFailedErrFormat           = `failed to decode: %w`
	requiredFieldErrFormat          = `%w: no value found for "%v"`
	unsupportedElementTypeErrFormat = `%w: cannot handle slice/array of "%v" at "%v"`
	parseErrFormat                  = `%w at "%v": %v`
	overflowErrFormat               = `%w: "%v" overflows type "%v" at "%v"`
)

// An InvalidInputError describes an invalid argument passed to Into function
// The argument must be a non-nil struct pointer
type InvalidInputError struct {
	Type  reflect.Type
	Value reflect.Value
}

func (e *InvalidInputError) Error() string {
	msg := "gonfig: invalid input: "

	if e.Type == nil {
		msg += "nil"
	} else if e.Type.Kind() != reflect.Ptr {
		msg += "non-pointer type"
	} else if e.Value.IsNil() {
		msg += "nil pointer"
	} else {
		msg += "non-struct type"
	}

	return msg
}

// Collection of errors during populating the input struct
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
