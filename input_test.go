package gonfig

import (
	"errors"
	"fmt"
	"math"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type badTypes struct {
	Ui   uintptr
	Uip  *uintptr
	C    chan int
	Cp   *chan int
	F    func()
	Fp   *func()
	I    interface{}
	Ip   *interface{}
	Up   unsafe.Pointer
	Upp  *unsafe.Pointer
	Ss   [][]int
	Ssp  *[][]int
	Sa   [][0]int
	Sap  *[][0]int
	Aa   [0][0]int
	Aap  *[0][0]int
	As   [0][]int
	Asp  *[0][]int
	Sui  []uintptr
	PSui *[]uintptr
	SuiP []*uintptr
	Sc   []chan int
	PScp *[]chan int
	Scp  []*chan int
	Sf   []func()
	PSf  *[]func()
	Sfp  []*func()
	Si   []interface{}
	PSi  *[]interface{}
	Sip  []*interface{}
	Su   []unsafe.Pointer
	PSu  *[]unsafe.Pointer
	Sup  []*unsafe.Pointer
	Aui  [0]uintptr
	AuiP *[0]uintptr
	Ac   [0]chan int
	Acp  *[0]chan int
	Af   [0]func()
	Afp  *[0]func()
	Ai   [0]interface{}
	Aip  *[0]interface{}
	Au   [0]unsafe.Pointer
	Aup  *[0]unsafe.Pointer
}

type supportedTypes struct {
	Boolean         bool
	BooleanPtr      *bool
	BooleanArray    [2]bool
	BooleanArrayPtr *[2]bool
	BooleanPtrArray [2]*bool
	BooleanSlice    []bool
	BooleanPtrSlice []*bool

	Uint8         uint8
	Uint8Ptr      *uint8
	Uint8Array    [2]uint8
	Uint8ArrayPtr *[2]uint8
	Uint8PtrArray [2]*uint8
	Uint8Slice    []uint8
	Uint8PtrSlice []*uint8

	Uint16         uint16
	Uint16Ptr      *uint16
	Uint16Array    [2]uint16
	Uint16ArrayPtr *[2]uint16
	Uint16PtrArray [2]*uint16
	Uint16Slice    []uint16
	Uint16PtrSlice []*uint16

	Uint32         uint32
	Uint32Ptr      *uint32
	Uint32Array    [2]uint32
	Uint32ArrayPtr *[2]uint32
	Uint32PtrArray [2]*uint32
	Uint32Slice    []uint32
	Uint32PtrSlice []*uint32

	Uint64         uint64
	Uint64Ptr      *uint64
	Uint64Array    [2]uint64
	Uint64ArrayPtr *[2]uint64
	Uint64PtrArray [2]*uint64
	Uint64Slice    []uint64
	Uint64PtrSlice []*uint64

	Uint         uint
	UintPtr      *uint
	UintArray    [2]uint
	UintArrayPtr *[2]uint
	UintPtrArray [2]*uint
	UintSlice    []uint
	UintPtrSlice []*uint

	Int8         int8
	Int8Ptr      *int8
	Int8Array    [2]int8
	Int8ArrayPtr *[2]int8
	Int8PtrArray [2]*int8
	Int8Slice    []int8
	Int8PtrSlice []*int8

	Int16         int16
	Int16Ptr      *int16
	Int16Array    [2]int16
	Int16ArrayPtr *[2]int16
	Int16PtrArray [2]*int16
	Int16Slice    []int16
	Int16PtrSlice []*int16

	Int32         int32
	Int32Ptr      *int32
	Int32Array    [2]int32
	Int32ArrayPtr *[2]int32
	Int32PtrArray [2]*int32
	Int32Slice    []int32
	Int32PtrSlice []*int32

	Int64         int64
	Int64Ptr      *int64
	Int64Array    [2]int64
	Int64ArrayPtr *[2]int64
	Int64PtrArray [2]*int64
	Int64Slice    []int64
	Int64PtrSlice []*int64

	Int         int
	IntPtr      *int
	IntArray    [2]int
	IntArrayPtr *[2]int
	IntPtrArray [2]*int
	IntSlice    []int
	IntPtrSlice []*int

	Float32         float32
	Float32Ptr      *float32
	Float32Array    [2]float32
	Float32ArrayPtr *[2]float32
	Float32PtrArray [2]*float32
	Float32Slice    []float32
	Float32PtrSlice []*float32

	Float64         float64
	Float64Ptr      *float64
	Float64Array    [2]float64
	Float64ArrayPtr *[2]float64
	Float64PtrArray [2]*float64
	Float64Slice    []float64
	Float64PtrSlice []*float64

	Complex64         complex64
	Complex64Ptr      *complex64
	Complex64Array    [2]complex64
	Complex64ArrayPtr *[2]complex64
	Complex64PtrArray [2]*complex64
	Complex64Slice    []complex64
	Complex64PtrSlice []*complex64

	Complex128         complex128
	Complex128Ptr      *complex128
	Complex128Array    [2]complex128
	Complex128ArrayPtr *[2]complex128
	Complex128PtrArray [2]*complex128
	Complex128Slice    []complex128
	Complex128PtrSlice []*complex128

	// uint8
	Byte         byte
	BytePtr      *byte
	ByteArray    [2]byte
	ByteArrayPtr *[2]byte
	BytePtrArray [2]*byte
	ByteSlice    []byte
	BytePtrSlice []*byte

	// int32
	Rune         rune
	RunePtr      *rune
	RuneArray    [2]rune
	RuneArrayPtr *[2]rune
	RunePtrArray [2]*rune
	RuneSlice    []rune
	RunePtrSlice []*rune

	String         string
	StringPtr      *string
	StringArray    [2]string
	StringArrayPtr *[2]string
	StringPtrArray [2]*string
	StringSlice    []string
	StringPtrSlice []*string

	// int64
	Duration         time.Duration
	DurationPtr      *time.Duration
	DurationArray    [2]time.Duration
	DurationArrayPtr *[2]time.Duration
	DurationPtrArray [2]*time.Duration
	DurationSlice    []time.Duration
	DurationPtrSlice []*time.Duration

	// Map         map[string]string
	// MapPtr      *map[string]string
	// MapArray    [2]map[string]string
	// MapArrayPtr *[2]map[string]string
	// MapPtrArray [2]*map[string]string
	// MapSlice    []map[string]string
	// MapPtrSlice []*map[string]string

	Time         time.Time
	TimePtr      *time.Time
	TimeArray    [2]time.Time
	TimeArrayPtr *[2]time.Time
	TimePtrArray [2]*time.Time
	TimeSlice    []time.Time
	TimePtrSlice []*time.Time

	Url         url.URL
	UrlPtr      *url.URL
	UrlArray    [2]url.URL
	UrlArrayPtr *[2]url.URL
	UrlPtrArray [2]*url.URL
	UrlSlice    []url.URL
	UrlPtrSlice []*url.URL

	Struct struct {
		Int    int
		String string
		Embed
		*EmbedP
	}
	StructPtr *struct {
		Int    int
		String string
		Embed
		*EmbedP
	}
	Embed
	*EmbedP
}

type Embed struct {
	Int    int
	String string
}

type EmbedP struct {
	Int    int
	String string
}

func TestNewInput(t *testing.T) {
	t.Run("bad inputs", func(t *testing.T) {
		t.Parallel()

		var (
			sIn    string
			bIn    bool
			iIn    int
			i8In   int8
			i16In  int16
			i32In  int32
			i64In  int64
			uIn    uint
			u8In   uint8
			u16In  uint16
			u32In  uint32
			u64In  uint64
			uiIn   uintptr
			f32In  float32
			f64In  float64
			c64In  complex64
			c128In complex128
			aIn    [0]string
			slIn   []string
			chIn   = make(chan int)
			fuIn   = func() {}
			mIn    = make(map[string]string)
			st     struct{}
		)

		tests := []interface{}{
			nil,
			sIn, &sIn,
			bIn, &bIn,
			iIn, i8In, i16In, i32In, i64In,
			&iIn, &i8In, &i16In, &i32In, &i64In,
			uIn, u8In, u16In, u32In, u64In,
			&uIn, &u8In, &u16In, &u32In, &u64In,
			uiIn, &uiIn,
			f32In, &f32In,
			f64In, &f64In,
			c64In, &c64In,
			c128In, &c128In,
			aIn, &aIn,
			slIn, &slIn,
			chIn, &chIn,
			fuIn, &fuIn,
			mIn, &mIn,
			st,
		}

		for _, tc := range tests {
			tc := tc
			t.Run(fmt.Sprint(reflect.TypeOf(tc)), func(t *testing.T) {
				t.Parallel()
				_, err := NewInput(tc)
				assert.IsType(t, &InvalidInputError{}, err)
			})
		}
	})

	t.Run("bad types", func(t *testing.T) {
		t.Parallel()

		t.Run("unexported fields", func(t *testing.T) {
			unexportedFields := struct {
				ue int
			}{}

			in, err := NewInput(&unexportedFields)
			assert.NoError(t, err)
			assert.Empty(t, in.Fields)
		})

		v := reflect.ValueOf(badTypes{})

		for i := 0; i < v.NumField(); i++ {
			i := i
			t.Run(fmt.Sprint(v.Field(i).Type()), func(t *testing.T) {
				t.Parallel()

				structType := reflect.StructOf([]reflect.StructField{
					v.Type().Field(i),
				})
				structValue := reflect.New(structType)

				in, err := NewInput(structValue.Interface())
				assert.Nil(t, in)
				require.Error(t, err)
				assert.Truef(
					t,
					errors.Is(err, ErrUnsupportedType),
					"Error must wrap ErrUnsupportedType error",
				)
			})
		}
	})

	t.Run("tags", func(t *testing.T) {
		t.Parallel()

		req := require.New(t)
		ass := assert.New(t)

		tags := struct {
			Defaults int
			Keys     int `config:"TAGS" json:"tags,omitempty" yaml:"tags" toml:""`
			Others   int `default:"5" required:"true" expand:"true" separator:"," format:"good-format"`
			Ignored1 int `config:"-"`
			Ignored2 int `ignore:"true"`
		}{}

		in, err := NewInput(&tags)
		req.NoError(err)
		req.NotNil(in)
		req.Len(in.Fields, 3)

		defaults := in.Fields[0]
		ass.Equal(defaultSeparator, defaults.Tags.Separator)
		ass.Equal(defaultFormat, defaults.Tags.Format)
		ass.Equal(defaultFormat, defaults.Tags.Format)

		keys := in.Fields[1]
		ass.Equal("TAGS", keys.Tags.Config)
		ass.Equal("tags", keys.Tags.Json)
		ass.Equal("tags", keys.Tags.Yaml)
		ass.Equal("", keys.Tags.Toml)

		others := in.Fields[2]
		ass.Equal("5", others.Tags.Default)
		ass.True(others.Tags.Required)
		ass.True(others.Tags.Expand)
		ass.Equal(",", others.Tags.Separator)
		ass.Equal("good-format", others.Tags.Format)
	})

	t.Run("path", func(t *testing.T) {
		t.Parallel()

		req := require.New(t)
		ass := assert.New(t)

		paths := struct {
			First struct {
				Second struct {
					Third        int
					ThirdSibling int
				}

				SecondSibling int
			}
		}{}

		in, err := NewInput(&paths)
		req.NoError(err)
		req.NotNil(in)
		req.Len(in.Fields, 3)

		ass.Equal([]string{"First", "Second", "Third"}, in.Fields[0].Path)
		ass.Equal([]string{"First", "Second", "ThirdSibling"}, in.Fields[1].Path)
		ass.Equal([]string{"First", "SecondSibling"}, in.Fields[2].Path)
	})

	t.Run("supported types", func(t *testing.T) {
		t.Parallel()

		in, err := NewInput(new(supportedTypes))
		require.NoError(t, err)
		require.NotNil(t, in)
		require.Len(t, in.Fields, 163)
	})
}

func TestInput_SetValue(t *testing.T) {
	t.Run("bad types", func(t *testing.T) {
		t.Parallel()

		in := Input{}

		t.Run("unexported fields", func(t *testing.T) {
			t.Parallel()

			f := Field{
				Value: reflect.Zero(reflect.TypeOf(0)),
			}

			err := in.SetValue(&f, "")
			require.Error(t, err)
			assert.Truef(
				t,
				errors.Is(err, ErrUnSettableField),
				"Error must wrap ErrUnSettableField error",
			)
		})

		v := reflect.ValueOf(badTypes{})
		for i := 0; i < v.NumField(); i++ {
			i := i
			t.Run(fmt.Sprint(v.Field(i).Type()), func(t *testing.T) {
				t.Parallel()

				f := Field{
					Value: reflect.New(v.Field(i).Type()).Elem(),
					Tags:  new(ConfigTags),
				}

				err := in.SetValue(&f, "")
				require.Error(t, err)
				assert.Truef(
					t,
					errors.Is(err, ErrUnsupportedType),
					"Error must wrap ErrUnsupportedType error",
				)
			})
		}
	})

	t.Run("supported types", func(t *testing.T) {
		t.Parallel()

		var (
			b    bool          = true
			ui8  uint8         = 1
			ui16 uint16        = 1
			ui32 uint32        = 1
			ui64 uint64        = 1
			ui   uint          = 1
			i8   int8          = 1
			i16  int16         = 1
			i32  int32         = 1
			i64  int64         = 1
			i    int           = 1
			f32  float32       = 3.14
			f64  float64       = 3.14
			c64  complex64     = 2 + 3i
			c128 complex128    = 2 + 3i
			by   byte          = 1
			r    rune          = 1
			s    string        = "nice"
			d    time.Duration = 60000000000
		)
		ti, _ := time.Parse(time.RFC3339, "2020-11-26T18:26:14+03:30")
		ti2, _ := time.Parse(time.RFC3339, "2021-11-26T18:26:49+03:30")
		ur, _ := url.Parse("golang.org")
		ur2, _ := url.Parse("google.com")

		tests := []struct {
			input    string
			expected interface{}
		}{
			{"true", true},
			{"true", &b},
			{"true false", [2]bool{true, false}},
			{"true false", &[2]bool{true, false}},
			{"true true", [2]*bool{&b, &b}},
			{"true false", []bool{true, false}},
			{"true true", []*bool{&b, &b}},

			{"1", uint8(1)},
			{"1", &ui8},
			{"1 2", [2]uint8{1, 2}},
			{"1 2", &[2]uint8{1, 2}},
			{"1 1", [2]*uint8{&ui8, &ui8}},
			{"1 2", []uint8{1, 2}},
			{"1 1", []*uint8{&ui8, &ui8}},

			{"1", uint16(1)},
			{"1", &ui16},
			{"1 2", [2]uint16{1, 2}},
			{"1 2", &[2]uint16{1, 2}},
			{"1 1", [2]*uint16{&ui16, &ui16}},
			{"1 2", []uint16{1, 2}},
			{"1 1", []*uint16{&ui16, &ui16}},

			{"1", uint32(1)},
			{"1", &ui32},
			{"1 2", [2]uint32{1, 2}},
			{"1 2", &[2]uint32{1, 2}},
			{"1 1", [2]*uint32{&ui32, &ui32}},
			{"1 2", []uint32{1, 2}},
			{"1 1", []*uint32{&ui32, &ui32}},

			{"1", uint64(1)},
			{"1", &ui64},
			{"1 2", [2]uint64{1, 2}},
			{"1 2", &[2]uint64{1, 2}},
			{"1 1", [2]*uint64{&ui64, &ui64}},
			{"1 2", []uint64{1, 2}},
			{"1 1", []*uint64{&ui64, &ui64}},

			{"1", uint(1)},
			{"1", &ui},
			{"1 2", [2]uint{1, 2}},
			{"1 2", &[2]uint{1, 2}},
			{"1 1", [2]*uint{&ui, &ui}},
			{"1 2", []uint{1, 2}},
			{"1 1", []*uint{&ui, &ui}},

			{"1", int8(1)},
			{"1", &i8},
			{"1 2", [2]int8{1, 2}},
			{"1 2", &[2]int8{1, 2}},
			{"1 1", [2]*int8{&i8, &i8}},
			{"1 2", []int8{1, 2}},
			{"1 1", []*int8{&i8, &i8}},

			{"1", int16(1)},
			{"1", &i16},
			{"1 2", [2]int16{1, 2}},
			{"1 2", &[2]int16{1, 2}},
			{"1 1", [2]*int16{&i16, &i16}},
			{"1 2", []int16{1, 2}},
			{"1 1", []*int16{&i16, &i16}},

			{"1", int32(1)},
			{"1", &i32},
			{"1 2", [2]int32{1, 2}},
			{"1 2", &[2]int32{1, 2}},
			{"1 1", [2]*int32{&i32, &i32}},
			{"1 2", []int32{1, 2}},
			{"1 1", []*int32{&i32, &i32}},

			{"1", int64(1)},
			{"1", &i64},
			{"1 2", [2]int64{1, 2}},
			{"1 2", &[2]int64{1, 2}},
			{"1 1", [2]*int64{&i64, &i64}},
			{"1 2", []int64{1, 2}},
			{"1 1", []*int64{&i64, &i64}},

			{"1", int(1)},
			{"1", &i},
			{"1 2", [2]int{1, 2}},
			{"1 2", &[2]int{1, 2}},
			{"1 1", [2]*int{&i, &i}},
			{"1 2", []int{1, 2}},
			{"1 1", []*int{&i, &i}},

			{"3.14", float32(3.14)},
			{"3.14", &f32},
			{"3.14 0.1", [2]float32{3.14, 0.1}},
			{"3.14 0.1", &[2]float32{3.14, 0.1}},
			{"3.14 3.14", [2]*float32{&f32, &f32}},
			{"3.14 0.1", []float32{3.14, 0.1}},
			{"3.14 3.14", []*float32{&f32, &f32}},

			{"3.14", float64(3.14)},
			{"3.14", &f64},
			{"3.14 0.1", [2]float64{3.14, 0.1}},
			{"3.14 0.1", &[2]float64{3.14, 0.1}},
			{"3.14 3.14", [2]*float64{&f64, &f64}},
			{"3.14 0.1", []float64{3.14, 0.1}},
			{"3.14 3.14", []*float64{&f64, &f64}},

			{"2+3i", complex64(2 + 3i)},
			{"2+3i", &c64},
			{"2+3i 4-1i", [2]complex64{2 + 3i, 4 - 1i}},
			{"2+3i 4-1i", &[2]complex64{2 + 3i, 4 - 1i}},
			{"2+3i 2+3i", [2]*complex64{&c64, &c64}},
			{"2+3i 4-1i", []complex64{2 + 3i, 4 - 1i}},
			{"2+3i 2+3i", []*complex64{&c64, &c64}},

			{"2+3i", complex128(2 + 3i)},
			{"2+3i", &c128},
			{"2+3i 4-1i", [2]complex128{2 + 3i, 4 - 1i}},
			{"2+3i 4-1i", &[2]complex128{2 + 3i, 4 - 1i}},
			{"2+3i 2+3i", [2]*complex128{&c128, &c128}},
			{"2+3i 4-1i", []complex128{2 + 3i, 4 - 1i}},
			{"2+3i 2+3i", []*complex128{&c128, &c128}},

			{"1", byte(1)},
			{"1", &by},
			{"1 2", [2]byte{1, 2}},
			{"1 2", &[2]byte{1, 2}},
			{"1 1", [2]*byte{&by, &by}},
			{"1 2", []byte{1, 2}},
			{"1 1", []*byte{&by, &by}},

			{"1", rune(1)},
			{"1", &r},
			{"1 2", [2]rune{1, 2}},
			{"1 2", &[2]rune{1, 2}},
			{"1 1", [2]*rune{&r, &r}},
			{"1 2", []rune{1, 2}},
			{"1 1", []*rune{&r, &r}},

			{"config", "config"},
			{"nice", &s},
			{"nice config", [2]string{"nice", "config"}},
			{"nice config", &[2]string{"nice", "config"}},
			{"nice nice", [2]*string{&s, &s}},
			{"nice config", []string{"nice", "config"}},
			{"nice nice", []*string{&s, &s}},

			{"1m", time.Duration(60000000000)},
			{"1m", &d},
			{"1m 1s", [2]time.Duration{60000000000, 1000000000}},
			{"1m 1s", &[2]time.Duration{60000000000, 1000000000}},
			{"1m 1m", [2]*time.Duration{&d, &d}},
			{"1m 1s", []time.Duration{60000000000, 1000000000}},
			{"1m 1m", []*time.Duration{&d, &d}},

			{"2020-11-26T18:26:14+03:30", ti},
			{"2020-11-26T18:26:14+03:30", &ti},
			{"2020-11-26T18:26:14+03:30 2021-11-26T18:26:49+03:30", [2]time.Time{ti, ti2}},
			{"2020-11-26T18:26:14+03:30 2021-11-26T18:26:49+03:30", &[2]time.Time{ti, ti2}},
			{"2020-11-26T18:26:14+03:30 2021-11-26T18:26:49+03:30", [2]*time.Time{&ti, &ti2}},
			{"2020-11-26T18:26:14+03:30 2021-11-26T18:26:49+03:30", []time.Time{ti, ti2}},
			{"2020-11-26T18:26:14+03:30 2021-11-26T18:26:49+03:30", []*time.Time{&ti, &ti2}},

			{"golang.org", *ur},
			{"golang.org", ur},
			{"golang.org google.com", [2]url.URL{*ur, *ur2}},
			{"golang.org google.com", &[2]url.URL{*ur, *ur2}},
			{"golang.org google.com", [2]*url.URL{ur, ur2}},
			{"golang.org google.com", []url.URL{*ur, *ur2}},
			{"golang.org google.com", []*url.URL{ur, ur2}},

			{"1", int(1)},
			{"config", "config"},
			{"1", int(1)},
			{"config", "config"},
			{"1", int(1)},
			{"config", "config"},
			{"1", int(1)},
			{"config", "config"},
			{"1", int(1)},
			{"config", "config"},
			{"1", int(1)},
			{"config", "config"},
			{"1", int(1)},
			{"config", "config"},
			{"1", int(1)},
			{"config", "config"},
		}

		in, err := NewInput(new(supportedTypes))
		require.NoError(t, err)
		require.NotNil(t, in)
		require.Equal(t, len(in.Fields), len(tests))

		for i, tc := range tests {
			f := in.Fields[i]
			tc := tc

			t.Run(strings.Join(f.Path, "."), func(t *testing.T) {
				t.Parallel()
				err := in.SetValue(f, tc.input)
				require.NoError(t, err)
				require.Equal(t, tc.expected, f.Value.Interface())
			})
		}
	})

	t.Run("parse error", func(t *testing.T) {
		t.Parallel()

		input := struct {
			B  bool
			I  int
			U  uint
			F  float64
			C  complex128
			D  time.Duration
			T  time.Time
			Ur url.URL
		}{}
		tests := []string{
			"bool",
			"int",
			"uint",
			"float",
			"complex",
			"duration",
			"time",
			"!@#$%^&*()-=+",
		}

		in, err := NewInput(&input)
		require.NoError(t, err)
		require.NotNil(t, in)
		require.Equal(t, len(tests), len(in.Fields))

		for i, tc := range tests {
			f := in.Fields[i]
			tc := tc

			t.Run(tc, func(t *testing.T) {
				err := in.SetValue(f, tc)
				require.Error(t, err)
				assert.Truef(
					t,
					errors.Is(err, ErrParsing),
					"Error must wrap ErrParsing error",
				)
			})
		}
	})

	t.Run("value overflow", func(t *testing.T) {
		t.Parallel()

		input := struct {
			I int8
			U uint8
			F float32
			C complex64
		}{}
		tests := []string{
			"128",
			"256",
			fmt.Sprint(math.MaxFloat64),
			fmt.Sprint(complex(math.MaxFloat64, math.MaxFloat64)),
		}

		in, err := NewInput(&input)
		require.NoError(t, err)
		require.NotNil(t, in)
		require.Equal(t, len(tests), len(in.Fields))

		for i, tc := range tests {
			f := in.Fields[i]
			tc := tc

			t.Run(tc, func(t *testing.T) {
				err := in.SetValue(f, tc)
				require.Error(t, err)
				assert.Truef(
					t,
					errors.Is(err, ErrValueOverflow),
					"Error must wrap ErrValueOverflow error",
				)
			})
		}
	})
}
