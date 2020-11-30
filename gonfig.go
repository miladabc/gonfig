package gonfig

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var ErrInvalidInput = errors.New("gonfig: input must be a struct pointer")

type Gonfig struct {
	Prefix string
}

func New(prefix string) *Gonfig {
	return &Gonfig{
		Prefix: prefix,
	}
}

func (g *Gonfig) Into(i interface{}) error {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return ErrInvalidInput
	}

	v = v.Elem()
	var prefix []string
	if g.Prefix != "" {
		prefix = append(prefix, g.Prefix)
	}

	populate(v, "", &ConfigTags{}, prefix...)

	return nil
}

func populate(v reflect.Value, value string, ct *ConfigTags, path ...string) {
	if ct.Ignore {
		return
	}
	if !v.CanSet() {
		fmt.Println("can not set")
		return
	}
	if value == "" {
		key := ct.Config
		if key == "" {
			key = toScreamingSnakeCase(path)
		}

		var exists bool
		value, exists = os.LookupEnv(key)
		if !exists {
			if ct.Required {
				return
			}

			value = ct.Default
		}

		if ct.Expand {
			value = os.ExpandEnv(value)
		}
	}

	switch v.Kind() {
	case reflect.String:
		v.SetString(value)

	case reflect.Bool:
		b, _ := strconv.ParseBool(value)
		v.SetBool(b)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var i int64

		if v.Type().PkgPath() == "time" && v.Type().Name() == "Duration" {
			d, _ := time.ParseDuration(value)
			i = int64(d)
		} else {
			i, _ = strconv.ParseInt(value, 0, 64)
		}

		if v.OverflowInt(i) {
			fmt.Printf("gonfig: value %v overflows type %v\n", i, v.Kind())
			return
		}

		v.SetInt(i)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, _ := strconv.ParseUint(value, 0, 64)
		if v.OverflowUint(i) {
			fmt.Printf("gonfig: value %v overflows type %v\n", i, v.Kind())
			return
		}

		v.SetUint(i)

	case reflect.Float32, reflect.Float64:
		f, _ := strconv.ParseFloat(value, v.Type().Bits())
		if v.OverflowFloat(f) {
			fmt.Printf("gonfig: value %v overflows type %v\n", f, v.Kind())
			return
		}

		v.SetFloat(f)

	case reflect.Complex64, reflect.Complex128:
		c, _ := strconv.ParseComplex(value, v.Type().Bits())
		if v.OverflowComplex(c) {
			fmt.Printf("gonfig: value %v overflows type %v\n", c, v.Kind())
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
			fmt.Printf("gonfig: cannot handle kind slice/array of %v\n", v.Type().Elem().Kind())
			return
		}

		var items []string
		for _, v := range strings.Split(value, ct.Separator) {
			item := strings.TrimSpace(v)
			if len(item) > 0 {
				items = append(items, item)
			}
		}
		if len(items) == 0 {
			return
		}

		switch v.Kind() {
		case reflect.Slice:
			size := len(items)
			sv := reflect.MakeSlice(reflect.SliceOf(v.Type().Elem()), size, size)

			for i := range items {
				populate(sv.Index(i), items[i], ct, path...)
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
				populate(av.Index(i), items[i], ct, path...)
			}

			v.Set(av)
		}

	case reflect.Map:
		// TODO

	case reflect.Ptr:
		pv := reflect.New(v.Type().Elem())
		populate(pv.Elem(), value, ct, path...)
		v.Set(pv)

	case reflect.Struct:
		if v.Type().Name() == "Time" {
			format := ct.Format
			if format == "" {
				format = time.RFC3339
			}

			t, _ := time.Parse(format, value)
			v.Set(reflect.ValueOf(t))
			return
		}

		if v.Type().Name() == "URL" {
			u, _ := url.Parse(value)
			v.Set(reflect.ValueOf(*u))
			return
		}

		for i := 0; i < v.NumField(); i++ {
			currentPath := append(path, v.Type().Field(i).Name)

			populate(
				v.Field(i),
				value,
				getTags(v.Type().Field(i).Tag),
				currentPath...,
			)
		}

	default:
		fmt.Printf("gonfig: cannot handle kind %v\n", v.Kind())
	}
}
