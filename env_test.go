package gonfig

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEnvProvider(t *testing.T) {
	ep := NewEnvProvider()
	require.NotNil(t, ep)
	assert.Equal(t, "", ep.EnvPrefix)
	assert.Equal(t, true, ep.SnakeCase)
	assert.Equal(t, true, ep.UpperCase)
	assert.Equal(t, "_", ep.FieldSeparator)
	assert.Equal(t, "", ep.Source)
	assert.Equal(t, false, ep.Required)
}

func TestEnvProvider_Name(t *testing.T) {
	ep := NewEnvProvider()
	assert.Equal(t, "ENV provider", ep.Name())
}

func TestEnvProvider_Fill(t *testing.T) {
	t.Run("env file existence", func(t *testing.T) {
		s := struct{}{}
		in, err := NewInput(&s)
		require.NoError(t, err)
		require.NotNil(t, in)
		ep := EnvProvider{
			Source:   "NotExistingFile",
			Required: false,
		}

		err = ep.Fill(in)
		assert.NoError(t, err)

		ep.Required = true
		err = ep.Fill(in)
		assert.Error(t, err)
	})

	t.Run("prioritize OS envs", func(t *testing.T) {
		os.Clearenv()
		err := os.Setenv("OS", "env from os")
		require.NoError(t, err)

		s := struct {
			Priority string
			Os       string
			File     string
		}{}
		in, err := NewInput(&s)
		require.NoError(t, err)
		require.NotNil(t, in)
		ep := NewEnvProvider()
		ep.Source = "testdata/.env"

		err = ep.Fill(in)
		require.NoError(t, err)
		assert.Equal(t, "FILE", s.Priority)
		assert.Equal(t, "env from os", s.Os)
		assert.Equal(t, "env from file", s.File)

		err = os.Setenv("PRIORITY", "OS")
		require.NoError(t, err)

		err = ep.Fill(in)
		require.NoError(t, err)
		assert.Equal(t, "OS", s.Priority)
	})

	t.Run("should be set", func(t *testing.T) {
		os.Clearenv()

		s := struct {
			Env1 string
			Env2 int
			Env3 bool
		}{}
		for i, v := range []string{"env1", "2", "false"} {
			err := os.Setenv("ENV"+fmt.Sprint(i+1), v)
			require.NoError(t, err)
		}

		in, err := NewInput(&s)
		require.NoError(t, err)
		require.NotNil(t, in)
		ep := NewEnvProvider()

		err = ep.Fill(in)
		require.NoError(t, err)

		for _, f := range in.Fields {
			assert.True(t, f.IsSet)
		}
	})

	t.Run("config key", func(t *testing.T) {
		os.Clearenv()
		err := os.Setenv("CUSTOM_KEY", "env")
		require.NoError(t, err)

		s := struct {
			Env string `config:"CUSTOM_KEY"`
		}{}
		in, err := NewInput(&s)
		require.NoError(t, err)
		require.NotNil(t, in)
		ep := EnvProvider{}

		err = ep.Fill(in)
		require.NoError(t, err)
		assert.Equal(t, "env", s.Env)
	})

	t.Run("env prefix", func(t *testing.T) {
		os.Clearenv()
		err := os.Setenv("APP_Env", "env")
		require.NoError(t, err)

		s := struct {
			Env string
		}{}

		in, err := NewInput(&s)
		require.NoError(t, err)
		require.NotNil(t, in)
		ep := EnvProvider{
			EnvPrefix: "APP_",
		}

		err = ep.Fill(in)
		require.NoError(t, err)
		assert.Equal(t, "env", s.Env)
	})

	t.Run("snake case", func(t *testing.T) {
		os.Clearenv()
		err := os.Setenv("Env_Var", "env")
		require.NoError(t, err)

		s := struct {
			EnvVar string
		}{}

		in, err := NewInput(&s)
		require.NoError(t, err)
		require.NotNil(t, in)
		ep := EnvProvider{
			SnakeCase: true,
		}

		err = ep.Fill(in)
		require.NoError(t, err)
		assert.Equal(t, "env", s.EnvVar)
	})

	t.Run("upper case", func(t *testing.T) {
		os.Clearenv()
		err := os.Setenv("ENV", "env")
		require.NoError(t, err)

		s := struct {
			Env string
		}{}

		in, err := NewInput(&s)
		require.NoError(t, err)
		require.NotNil(t, in)
		ep := EnvProvider{
			UpperCase: true,
		}

		err = ep.Fill(in)
		require.NoError(t, err)
		assert.Equal(t, "env", s.Env)
	})

	t.Run("field separator", func(t *testing.T) {
		os.Clearenv()
		err := os.Setenv("NestedEnvVar", "env")
		require.NoError(t, err)

		s := struct {
			Nested struct {
				EnvVar string
			}
		}{}

		in, err := NewInput(&s)
		require.NoError(t, err)
		require.NotNil(t, in)
		ep := EnvProvider{
			FieldSeparator: "",
		}

		err = ep.Fill(in)
		require.NoError(t, err)
		assert.Equal(t, "env", s.Nested.EnvVar)
	})
}
