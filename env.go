package gonfig

import (
	"errors"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// EnvProvider loads values from environment variables to provided struct
type EnvProvider struct {
	// Prefix is used when finding values from environment variables, defaults to ""
	EnvPrefix string

	// SnakeCase specifies whether to convert field names to snake_case or not, defaults to true
	SnakeCase bool

	// UpperCase specifies whether to convert field names to UPPERCASE or not, defaults to true
	UpperCase bool

	// FieldSeparator is used to separate field names, defaults to "_"
	FieldSeparator string

	// Source is used to retrieve environment variables
	// It can be either a path to a file or empty string, if empty OS will be used
	Source string

	// Whether to report error if env file is not found, defaults to false
	Required bool
}

// NewEnvProvider creates a new EnvProvider
func NewEnvProvider() *EnvProvider {
	return &EnvProvider{
		EnvPrefix:      "",
		SnakeCase:      true,
		UpperCase:      true,
		FieldSeparator: "_",
		Source:         "",
		Required:       false,
	}
}

// NewEnvFileProvider creates a new EnvProvider from .env file
func NewEnvFileProvider(path string) *EnvProvider {
	return &EnvProvider{
		EnvPrefix:      "",
		SnakeCase:      true,
		UpperCase:      true,
		FieldSeparator: "_",
		Source:         path,
		Required:       false,
	}
}

// Name of provider
func (ep *EnvProvider) Name() string {
	return "ENV provider"
}

// Fill takes struct fields and fills their values
func (ep *EnvProvider) Fill(in *Input) error {
	content, err := ep.envMap()
	if err != nil {
		return err
	}

	for _, f := range in.Fields {
		value, err := ep.provide(content, f.Tags.Config, f.Path)
		if err != nil {
			if errors.Is(err, ErrKeyNotFound) {
				continue
			}

			return err
		}

		err = in.setValue(f, value)
		if err != nil {
			return err
		}

		f.IsSet = true
	}

	return nil
}

// envMap returns environment variables map from either OS or file specified by source
// Defaults to operating system env variables
func (ep *EnvProvider) envMap() (map[string]string, error) {
	envs := envFromOS()
	var fileEnvs map[string]string
	var err error

	if ep.Source != "" {
		fileEnvs, err = envFromFile(ep.Source)
	}
	if err != nil {
		notExistsErr := errors.Is(err, os.ErrNotExist)
		if (notExistsErr && ep.Required) || !notExistsErr {
			return nil, err
		}
	}

	if len(envs) == 0 {
		if len(fileEnvs) == 0 {
			return nil, nil
		}

		envs = make(map[string]string)
	}

	for k, v := range fileEnvs {
		_, exists := envs[k]
		if !exists {
			envs[k] = v
		}
	}

	return envs, nil
}

// returns environment variables map retrieved from operating system
func envFromOS() map[string]string {
	envs := os.Environ()
	if len(envs) == 0 {
		return nil
	}

	envMap := make(map[string]string)

	for _, env := range envs {
		keyValue := strings.SplitN(env, "=", 2)
		if len(keyValue) < 2 {
			continue
		}

		envMap[keyValue[0]] = keyValue[1]
	}

	return envMap
}

// returns environment variables map retrieved from specified file
func envFromFile(path string) (map[string]string, error) {
	m, err := godotenv.Read(path)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// provide find a value from env variables based on specified key and path
func (ep *EnvProvider) provide(content map[string]string, key string, path []string) (string, error) {
	k := ep.buildKey(key, path)
	value, exists := content[k]
	if !exists {
		return "", ErrKeyNotFound
	}

	return value, nil
}

// buildKey prefix key with EnvPrefix, if not provided, path slice will be used
func (ep *EnvProvider) buildKey(key string, path []string) string {
	if key != "" {
		return ep.EnvPrefix + key
	}

	k := strings.Join(path, ep.FieldSeparator)
	if ep.SnakeCase {
		k = toSnakeCase(k)
	}
	if ep.UpperCase {
		k = strings.ToUpper(k)
	}

	k = ep.EnvPrefix + k

	return k
}
