package gonfig

import (
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Input stores information about given struct
type Input struct {
	// Struct name is used for error messages
	Name string

	// Fields information
	Fields []*Field
}

// Struct field information
type Field struct {
	// Field value
	Value reflect.Value

	// Field tags
	Tags *ConfigTags

	// Slice of field names from root of struct all the way down to the field
	Path []string

	// IsSet specifies whether field value is set by one of the providers
	IsSet bool
}

// NewInput validates and returns a new Input with all settable fields
// Input argument must be a non-nil struct pointer
func NewInput(i interface{}) (*Input, error) {
	v := reflect.ValueOf(i)

	if err := validateInput(v); err != nil {
		return nil, err
	}

	in := Input{
		Name: v.Type().String(),
	}

	f := Field{
		Value: v.Elem(),
		Tags:  new(ConfigTags),
	}

	if err := in.traverseField(&f); err != nil {
		return nil, err
	}

	return &in, nil
}

// validateInput checks for a non-nil struct pointer
func validateInput(v reflect.Value) error {
	if !v.IsValid() ||
		v.Type() == nil ||
		v.Type().Kind() != reflect.Ptr ||
		v.IsNil() ||
		v.Type().Elem().Kind() != reflect.Struct {
		return &InvalidInputError{
			Value: v,
		}
	}

	return nil
}

// traverseField recursively traverse all fields and collect their information
func (in *Input) traverseField(f *Field) error {
	if !f.Value.CanSet() || f.Tags.Ignore {
		return nil
	}

	if err := in.isSupportedType(f.Value.Type()); err != nil {
		return fmt.Errorf(badFieldErrFormat, in.getPath(f.Path), err)
	}

	if isStruct(f.Value.Type()) {
		for i := 0; i < f.Value.NumField(); i++ {
			nestedField := Field{
				Value: f.Value.Field(i),
				Tags:  extractTags(f.Value.Type().Field(i).Tag),
				Path:  append(f.Path, f.Value.Type().Field(i).Name),
			}

			if err := in.traverseField(&nestedField); err != nil {
				return err
			}
		}

		return nil
	}

	if f.Value.Kind() == reflect.Ptr && isStruct(f.Value.Type().Elem()) {
		if f.Value.IsNil() {
			initPtr(f.Value)
		}

		pointedField := Field{
			Value: f.Value.Elem(),
			Tags:  f.Tags,
			Path:  f.Path,
		}

		return in.traverseField(&pointedField)
	}

	in.collectField(f)
	return nil
}

func (in *Input) collectField(f *Field) {
	in.Fields = append(in.Fields, f)
}

func (in *Input) isSupportedType(t reflect.Type) error {
	switch t.Kind() {
	case reflect.Invalid,
		reflect.Uintptr,
		reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.UnsafePointer:
		return fmt.Errorf(unsupportedTypeErrFormat, ErrUnsupportedType, t.Kind())

	case reflect.Slice, reflect.Array:
		switch t.Elem().Kind() {
		case reflect.Slice, reflect.Array:
			return fmt.Errorf(unsupportedTypeErrFormat, ErrUnsupportedType, "multi-dimensional slice/array")

		default:
			return in.isSupportedType(t.Elem())
		}

	case reflect.Ptr:
		return in.isSupportedType(t.Elem())
	}

	return nil
}

// SetValue validates and sets the value of a struct field
// returns error in case of unSettable field or unsupported type
func (in *Input) SetValue(f *Field, value string) error {
	if !f.Value.CanSet() {
		return fmt.Errorf(
			unSettableFieldErrFormat,
			ErrUnSettableField, in.getPath(f.Path),
		)
	}
	if err := in.isSupportedType(f.Value.Type()); err != nil {
		return fmt.Errorf(badFieldErrFormat, in.getPath(f.Path), err)
	}

	if f.Tags.Expand {
		value = os.ExpandEnv(value)
	}

	switch f.Value.Kind() {
	case reflect.String:
		return in.setString(f, value)

	case reflect.Bool:
		return in.setBool(f, value)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if isDuration(f.Value.Type()) {
			return in.setDuration(f, value)
		}

		return in.setInt(f, value)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return in.setUint(f, value)

	case reflect.Float32, reflect.Float64:
		return in.setFloat(f, value)

	case reflect.Complex64, reflect.Complex128:
		return in.setComplex(f, value)

	case reflect.Slice:
		return in.setSlice(f, value)

	case reflect.Array:
		return in.setArray(f, value)

	case reflect.Map:
		return in.setMap(f, value)

	case reflect.Ptr:
		return in.setPointer(f, value)

	case reflect.Struct:
		if isTime(f.Value.Type()) {
			return in.setTime(f, value)
		}

		if isURL(f.Value.Type()) {
			return in.setUrl(f, value)
		}
	}

	return nil
}

func (in *Input) setString(f *Field, value string) error {
	f.Value.SetString(value)
	return nil
}

func (in *Input) setBool(f *Field, value string) error {
	b, err := strconv.ParseBool(value)
	if err != nil {
		return fmt.Errorf(
			parseErrFormat,
			ErrParsing, in.getPath(f.Path), err,
		)
	}

	f.Value.SetBool(b)
	return nil
}

func (in *Input) setInt(f *Field, value string) error {
	i, err := strconv.ParseInt(value, 0, 64)
	if err != nil {
		return fmt.Errorf(
			parseErrFormat,
			ErrParsing, in.getPath(f.Path), err,
		)
	}
	if f.Value.OverflowInt(i) {
		return fmt.Errorf(
			overflowErrFormat,
			ErrValueOverflow, i, f.Value.Kind(), in.getPath(f.Path),
		)
	}

	f.Value.SetInt(i)
	return nil
}

func (in *Input) setUint(f *Field, value string) error {
	i, err := strconv.ParseUint(value, 0, 64)
	if err != nil {
		return fmt.Errorf(
			parseErrFormat,
			ErrParsing, in.getPath(f.Path), err,
		)
	}
	if f.Value.OverflowUint(i) {
		return fmt.Errorf(
			overflowErrFormat,
			ErrValueOverflow, i, f.Value.Kind(), in.getPath(f.Path),
		)
	}

	f.Value.SetUint(i)
	return nil
}

func (in *Input) setFloat(f *Field, value string) error {
	fv, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf(
			parseErrFormat,
			ErrParsing, in.getPath(f.Path), err,
		)
	}
	if f.Value.OverflowFloat(fv) {
		return fmt.Errorf(
			overflowErrFormat,
			ErrValueOverflow, fv, f.Value.Kind(), in.getPath(f.Path),
		)
	}

	f.Value.SetFloat(fv)
	return nil
}

func (in *Input) setComplex(f *Field, value string) error {
	c, err := strconv.ParseComplex(value, 128)
	if err != nil {
		return fmt.Errorf(
			parseErrFormat,
			ErrParsing, in.getPath(f.Path), err,
		)
	}
	if f.Value.OverflowComplex(c) {
		return fmt.Errorf(
			overflowErrFormat,
			ErrValueOverflow, c, f.Value.Kind(), in.getPath(f.Path),
		)
	}

	f.Value.SetComplex(c)
	return nil
}

func (in *Input) setSlice(f *Field, value string) error {
	items := extractItems(value, f.Tags.Separator)
	size := len(items)
	if size == 0 {
		return nil
	}
	s := reflect.MakeSlice(reflect.SliceOf(f.Value.Type().Elem()), size, size)

	for i := range items {
		nestedField := Field{
			Value: s.Index(i),
			Tags:  f.Tags,
			Path:  f.Path,
		}

		if err := in.SetValue(&nestedField, items[i]); err != nil {
			return err
		}
	}

	f.Value.Set(s)
	return nil
}

func (in *Input) setArray(f *Field, value string) error {
	items := extractItems(value, f.Tags.Separator)
	size := f.Value.Len()
	if size == 0 || len(items) == 0 {
		return nil
	}

	a := reflect.New(reflect.ArrayOf(size, f.Value.Type().Elem())).Elem()

	for i := 0; i < size; i++ {
		nestedField := Field{
			Value: a.Index(i),
			Tags:  f.Tags,
			Path:  f.Path,
		}

		if err := in.SetValue(&nestedField, items[i]); err != nil {
			return err
		}
	}

	f.Value.Set(a)
	return nil
}

func (in *Input) setMap(f *Field, value string) error {
	// TODO
	return nil
}

func (in *Input) setPointer(f *Field, value string) error {
	if f.Value.IsNil() {
		initPtr(f.Value)
	}

	pointedField := Field{
		Value: f.Value.Elem(),
		Tags:  f.Tags,
		Path:  f.Path,
	}

	return in.SetValue(&pointedField, value)
}

func (in *Input) setDuration(f *Field, value string) error {
	d, err := time.ParseDuration(value)
	if err != nil {
		return fmt.Errorf(
			parseErrFormat,
			ErrParsing, in.getPath(f.Path), err,
		)
	}

	f.Value.SetInt(int64(d))
	return nil
}

func (in *Input) setTime(f *Field, value string) error {
	t, err := time.Parse(f.Tags.Format, value)
	if err != nil {
		return fmt.Errorf(
			parseErrFormat,
			ErrParsing, in.getPath(f.Path), err,
		)
	}

	f.Value.Set(reflect.ValueOf(t))
	return nil
}

func (in *Input) setUrl(f *Field, value string) error {
	u, err := url.Parse(value)
	if err != nil {
		return fmt.Errorf(
			parseErrFormat,
			ErrParsing, in.getPath(f.Path), err,
		)
	}

	f.Value.Set(reflect.ValueOf(*u))
	return nil
}

func initPtr(v reflect.Value) {
	v.Set(reflect.New(v.Type().Elem()))
}

// getPath returns a dot separated string prefixed with struct name
func (in *Input) getPath(paths []string) string {
	return in.Name + "." + strings.Join(paths, ".")
}
