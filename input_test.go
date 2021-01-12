package gonfig

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

	t.Run("bad fields", func(t *testing.T) {
		t.Parallel()

		t.Run("unexported fields", func(t *testing.T) {
			unexportedFields := struct {
				ue int
			}{}

			in, err := NewInput(&unexportedFields)
			assert.NoError(t, err)
			assert.Empty(t, in.Fields)
		})

		badFields := struct {
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
			SuiP *[]uintptr
			Sc   []chan int
			Scp  *[]chan int
			Sf   []func()
			Sfp  *[]func()
			Si   []interface{}
			Sip  *[]interface{}
			Su   []unsafe.Pointer
			Sup  *[]unsafe.Pointer
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
		}{}
		v := reflect.ValueOf(badFields)

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

	t.Run("pointer field", func(t *testing.T) {
		t.Parallel()

		pf := struct {
			Ip *int
		}{}

		in, err := NewInput(&pf)
		require.NoError(t, err)
		require.NotNil(t, in)
		require.Len(t, in.Fields, 1)

		require.NotNil(t, pf.Ip)
		assert.Zero(t, *pf.Ip)
	})

	t.Run("all possible types", func(t *testing.T) {
		t.Parallel()

		type Embed struct {
			Int    int
			String string
		}

		type possibleTypes struct {
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

			Map         map[string]string
			MapPtr      *map[string]string
			MapArray    [2]map[string]string
			MapArrayPtr *[2]map[string]string
			MapPtrArray [2]*map[string]string
			MapSlice    []map[string]string
			MapPtrSlice []*map[string]string

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
				Int int
			}
			StructPtr *struct {
				Int int
			}
			Embed
		}

		in, err := NewInput(new(possibleTypes))
		require.NoError(t, err)
		require.NotNil(t, in)
		require.Len(t, in.Fields, 158)
	})
}
