package gonfig

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFileProvider(t *testing.T) {
	fp := NewFileProvider("file.yml")
	require.NotNil(t, fp)
	assert.Equal(t, "file.yml", fp.FilePath)
	assert.Equal(t, ".yml", fp.FileExt)
	assert.False(t, fp.Required)
}

func TestFileProvider_Name(t *testing.T) {
	fp := FileProvider{
		FileExt: ".json",
	}

	assert.Equal(t, "File provider (json)", fp.Name())
}

func TestFileProvider_UnmarshalStruct(t *testing.T) {
	t.Run("file existence", func(t *testing.T) {
		fp := FileProvider{
			FilePath: "NotExistingFile.toml",
			FileExt:  ".toml",
			Required: false,
		}

		var i interface{}
		err := fp.UnmarshalStruct(i)
		assert.NoError(t, err)

		fp.Required = true
		err = fp.UnmarshalStruct(i)
		assert.Error(t, err)
	})

	t.Run("unsupported file extension", func(t *testing.T) {
		fp := FileProvider{
			FileExt: ".ini",
		}

		var i interface{}
		err := fp.UnmarshalStruct(i)
		require.Error(t, err)
		assert.Truef(
			t,
			errors.Is(err, ErrUnsupportedFileExt),
			"Error must wrap ErrUnsupportedFileExt error",
		)
	})

	t.Run("supported file extensions", func(t *testing.T) {
		for _, e := range []string{".json", ".yml", ".yaml", ".toml"} {
			s := struct{}{}
			fp := FileProvider{
				FilePath: "testdata/config" + e,
				FileExt:  e,
				Required: true,
			}

			err := fp.UnmarshalStruct(&s)
			require.NoError(t, err)
		}
	})
}

func TestFileProvider_Fill(t *testing.T) {
	t.Run("should be set", func(t *testing.T) {
		for _, e := range []string{".json", ".yml", ".yaml", ".toml"} {
			s := struct {
				Config struct {
					Host string
				}
			}{}
			in, err := NewInput(&s)
			require.NoError(t, err)
			require.NotNil(t, in)

			fp := FileProvider{
				FilePath: "testdata/config" + e,
				FileExt:  e,
				Required: true,
			}

			err = fp.Fill(in)
			require.NoError(t, err)
			for _, f := range in.Fields {
				assert.True(t, f.IsSet)
			}
		}
	})

	t.Run("config key", func(t *testing.T) {
		for _, e := range []string{".json", ".yml", ".yaml", ".toml"} {
			s := struct {
				Custom string `json:"custom_key" yaml:"custom_key" toml:"custom_key"`
			}{}
			in, err := NewInput(&s)
			require.NoError(t, err)
			require.NotNil(t, in)

			fp := FileProvider{
				FilePath: "testdata/config" + e,
				FileExt:  e,
				Required: true,
			}

			err = fp.Fill(in)
			require.NoError(t, err)
			for _, f := range in.Fields {
				assert.True(t, f.IsSet)
			}
		}
	})
}
