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

// TODO: separate this into another struct
type Gonfig struct {
	Prefix     string
	structName string
	ce         ConfigErrors
}

func New(prefix string) *Gonfig {
	return &Gonfig{
		Prefix: prefix,
	}
}

// Input must be a non-nil struct pointer
func checkInput(i interface{}) error {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)

	if t == nil ||
		t.Kind() != reflect.Ptr ||
		v.IsNil() ||
		t.Elem().Kind() != reflect.Struct {
		return &InvalidInputError{
			Type:  t,
			Value: v,
		}
	}

	return nil
}

func (g *Gonfig) Into(i interface{}) error {
	if err := checkInput(i); err != nil {
		return err
	}

	v := reflect.ValueOf(i)
	g.structName = v.Type().String()
	v = v.Elem()

	g.populate(v, "", &ConfigTags{})

	if len(g.ce) != 0 {
		return g.ce
	}

	return nil
}

func (g *Gonfig) populate(v reflect.Value, value string, tags *ConfigTags, path ...string) {
	if tags.Ignore || !v.CanSet() {
		return
	}

	// TODO: it should not called here, if struct => bug!
	if v.Kind() != reflect.Struct && value == "" {
		var key string
		if tags.Config != "" {
			key = g.Prefix + tags.Config
		} else {
			key = g.Prefix + toScreamingSnakeCase(path)
		}

		var exists bool
		value, exists = os.LookupEnv(key)
		if !exists {
			if tags.Required {
				g.collectError(fmt.Errorf(missingValueErrFormat, ErrMissingValue, g.getPath(path)))
				return
			} else {
				value = tags.Default
			}
		}

		if tags.Expand {
			value = os.ExpandEnv(value)
		}
	}

	switch v.Kind() {
	case reflect.String:
		v.SetString(value)

	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			g.collectError(
				fmt.Errorf(
					parseErrFormat,
					ErrParsing, g.getPath(path), err,
				),
			)
			return
		}

		v.SetBool(b)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var d time.Duration
		var i int64
		var err error

		if isDuration(v) {
			d, err = time.ParseDuration(value)
			if err != nil {
				g.collectError(
					fmt.Errorf(
						parseErrFormat,
						ErrParsing, g.getPath(path), err,
					),
				)
				return
			}

			i = int64(d)
		} else {
			i, err = strconv.ParseInt(value, 0, 64)
			if err != nil {
				g.collectError(
					fmt.Errorf(
						parseErrFormat,
						ErrParsing, g.getPath(path), err,
					),
				)
				return
			}
		}

		if v.OverflowInt(i) {
			g.collectError(
				fmt.Errorf(
					overflowErrFormat,
					ErrValueOverflow, i, v.Kind(), g.getPath(path),
				),
			)
			return
		}

		v.SetInt(i)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(value, 0, 64)
		if err != nil {
			g.collectError(
				fmt.Errorf(
					parseErrFormat,
					ErrParsing, g.getPath(path), err,
				),
			)
			return
		}

		if v.OverflowUint(i) {
			g.collectError(
				fmt.Errorf(
					overflowErrFormat,
					ErrValueOverflow, i, v.Kind(), g.getPath(path),
				),
			)
			return
		}

		v.SetUint(i)

	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(value, v.Type().Bits())
		if err != nil {
			g.collectError(
				fmt.Errorf(
					parseErrFormat,
					ErrParsing, g.getPath(path), err,
				),
			)
			return
		}

		if v.OverflowFloat(f) {
			g.collectError(
				fmt.Errorf(
					overflowErrFormat,
					ErrValueOverflow, f, v.Kind(), g.getPath(path),
				),
			)
			return
		}

		v.SetFloat(f)

	case reflect.Complex64, reflect.Complex128:
		c, err := strconv.ParseComplex(value, v.Type().Bits())
		if err != nil {
			g.collectError(
				fmt.Errorf(
					parseErrFormat,
					ErrParsing, g.getPath(path), err,
				),
			)
			return
		}

		if v.OverflowComplex(c) {
			g.collectError(
				fmt.Errorf(
					overflowErrFormat,
					ErrValueOverflow, c, v.Kind(), g.getPath(path),
				),
			)
			return
		}

		v.SetComplex(c)

	case reflect.Slice, reflect.Array:
		switch v.Type().Elem().Kind() {
		case reflect.Slice,
			reflect.Array,
			reflect.Uintptr,
			reflect.Chan,
			reflect.Func,
			reflect.Interface,
			reflect.UnsafePointer:
			g.collectError(
				fmt.Errorf(
					unsupportedElementTypeErrFormat,
					ErrUnsupportedType, v.Type().Elem().Kind(), g.getPath(path),
				),
			)
			return
		}

		var items []string
		for _, v := range strings.Split(value, tags.Separator) {
			item := strings.TrimSpace(v)
			if len(item) > 0 {
				items = append(items, item)
			}
		}
		if len(items) == 0 {
			return
		}

		switch v.Kind() {
		// FIXME: in case of parse error slice should not get initialized
		case reflect.Slice:
			size := len(items)
			sv := reflect.MakeSlice(reflect.SliceOf(v.Type().Elem()), size, size)

			for i := range items {
				g.populate(sv.Index(i), items[i], tags, path...)
			}

			v.Set(sv)

		case reflect.Array:
			size := v.Len()
			if size == 0 {
				return
			}

			at := reflect.ArrayOf(size, v.Type().Elem())
			av := reflect.New(at).Elem()

			for i := 0; i < size; i++ {
				g.populate(av.Index(i), items[i], tags, path...)
			}

			v.Set(av)
		}

	case reflect.Map:
		// TODO

	case reflect.Ptr:
		pv := reflect.New(v.Type().Elem())
		g.populate(pv.Elem(), value, tags, path...)
		v.Set(pv)

	case reflect.Struct:
		if isTime(v) {
			format := tags.Format
			if format == "" {
				format = time.RFC3339
			}

			t, err := time.Parse(format, value)
			if err != nil {
				g.collectError(
					fmt.Errorf(
						parseErrFormat,
						ErrParsing, g.getPath(path), err,
					),
				)
				return
			}

			v.Set(reflect.ValueOf(t))
			return
		}

		if isURL(v) {
			u, err := url.Parse(value)
			if err != nil {
				g.collectError(
					fmt.Errorf(
						parseErrFormat,
						ErrParsing, g.getPath(path), err,
					),
				)
				return
			}

			v.Set(reflect.ValueOf(*u))
			return
		}

		for i := 0; i < v.NumField(); i++ {
			currentPath := append(path, v.Type().Field(i).Name)

			g.populate(
				v.Field(i),
				value,
				getTags(v.Type().Field(i).Tag),
				currentPath...,
			)
		}

	default:
		g.collectError(
			fmt.Errorf(
				unsupportedTypeErrFormat,
				ErrUnsupportedType, v.Kind(), g.getPath(path),
			),
		)
	}
}

func (g *Gonfig) collectError(e error) {
	g.ce = append(g.ce, e)
}

func (g *Gonfig) getPath(paths []string) string {
	return g.structName + "." + strings.Join(paths, ".")
}

func isDuration(v reflect.Value) bool {
	return v.Type().PkgPath() == "time" && v.Type().Name() == "Duration"
}

func isTime(v reflect.Value) bool {
	return v.Type().PkgPath() == "time" && v.Type().Name() == "Time"
}

func isURL(v reflect.Value) bool {
	return v.Type().PkgPath() == "net/url" && v.Type().Name() == "URL"
}
