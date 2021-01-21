package gonfig

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	c := Load()
	assert.NotNil(t, c)
}

func TestConfig_FromEnv(t *testing.T) {
	c := Config{}

	c.FromEnv()
	require.Len(t, c.Providers, 1)
	assert.IsType(t, new(EnvProvider), c.Providers[0])
}

func TestConfig_FromFile(t *testing.T) {
	t.Run("json-yaml-toml files", func(t *testing.T) {
		c := Config{}

		c.FromFile("file.json")
		require.Len(t, c.Providers, 1)
		assert.IsType(t, new(FileProvider), c.Providers[0])
	})

	t.Run("env files", func(t *testing.T) {
		c := Config{}

		c.FromFile("file.env")
		require.Len(t, c.Providers, 1)
		assert.IsType(t, new(EnvProvider), c.Providers[0])
	})
}

func TestConfig_AddProvider(t *testing.T) {
	c := Config{}

	var p Provider
	c.AddProvider(p)
	require.Len(t, c.Providers, 1)
}

func TestConfig_Into(t *testing.T) {
	s := struct {
		Required string `required:"true"`
		Default  string `default:"default_value"`
		Expand   string `expand:"true"`
	}{
		Expand: "${EXPAND}",
	}
	err := os.Setenv("EXPAND", "expand_value")
	require.NoError(t, err)

	err = Load().FromEnv().FromFile("testdata/config.json").Into(&s)
	require.Error(t, err)
	ce := err.(ConfigErrors)
	assert.Truef(
		t,
		errors.Is(ce[0], ErrRequiredField),
		"must wrap ErrRequiredField error",
	)
	assert.Equal(t, "default_value", s.Default)
	assert.Equal(t, "expand_value", s.Expand)
}
