package gonfig

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

// Supported file extensions
const (
	JSON = ".json"
	YML  = ".yml"
	YAML = ".yaml"
	ENV  = ".env"
	TOML = ".toml"
)

// FileProvider loads values from file to provided struct
type FileProvider struct {
	// Path to file
	FilePath string

	// File will be decoded based on extension
	// .json, .yml(.yaml), .env and .toml file extensions are supported
	FileExt string

	// Whether to report error if file is not found, defaults to false
	Required bool
}

// NewFileProvider creates a new FileProvider from specified path
func NewFileProvider(path string) *FileProvider {
	return &FileProvider{
		FilePath: path,
		FileExt:  filepath.Ext(path),
		Required: false,
	}
}

// Name of provider
func (fp *FileProvider) Name() string {
	return "File provider"
}

// UnmarshalStruct takes a struct pointer and loads values from provided file into it
func (fp *FileProvider) UnmarshalStruct(i interface{}) error {
	return fp.decode(i)
}

// Fill takes struct fields and and checks if their value is set
func (fp *FileProvider) Fill(in *Input) error {
	var content map[string]interface{}
	if err := fp.decode(&content); err != nil {
		return err
	}

	for _, f := range in.Fields {
		if f.IsSet {
			continue
		}

		var key string
		switch fp.FileExt {
		case JSON:
			key = f.Tags.Json
		case YML, YAML:
			key = f.Tags.Yaml
		case TOML:
			key = f.Tags.Toml
		}

		_, err := fp.provide(content, key, f.Path)
		if err == nil {
			f.IsSet = true
		}
	}

	return nil
}

// decode opens specified file and loads its content to input argument
func (fp *FileProvider) decode(i interface{}) (err error) {
	f, err := os.Open(fp.FilePath)
	if err != nil {
		if os.IsNotExist(err) && !fp.Required {
			return nil
		}

		return fmt.Errorf("file provider: %w", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	switch fp.FileExt {
	case JSON:
		err = json.NewDecoder(f).Decode(i)

	case YML, YAML:
		err = yaml.NewDecoder(f).Decode(i)

	case TOML:
		_, err = toml.DecodeReader(f, i)

	default:
		err = fmt.Errorf(unsupportedFileExtErrFormat, ErrUnsupportedFileExt, fp.FileExt)
	}

	if err != nil {
		return fmt.Errorf(decodeFailedErrFormat, err)
	}

	return nil
}

// provide find a value from file content based on specified key and path
func (fp *FileProvider) provide(content map[string]interface{}, key string, path []string) (string, error) {
	return traverseMap(content, fp.buildPath(key, path))
}

// buildPath makes a path from key and path slice
func (fp *FileProvider) buildPath(key string, path []string) []string {
	newPath := make([]string, len(path))
	copy(newPath, path)

	if key != "" {
		newPath[len(newPath)-1] = key
	}

	return newPath
}
