package gonfig

import (
	"reflect"
	"strings"
	"time"
)

const (
	defaultSeparator = " "
	ignoreCharacter  = "-"
)

// Possible tags, all are optional
type ConfigTags struct {
	// Key to be used by providers to retrieve the needed value, defaults to field name.
	// Use "-" to ignore the field.
	Config string

	// json tag for json files
	Json string

	// yaml tag for yaml files
	Yaml string

	// toml tag for toml files
	Toml string

	// Default value for field.
	Default string

	// Specify if value should be present, defaults to false.
	Required bool

	// Specify if field should be ignored, defaults to false.
	Ignore bool

	// Specify if value should be expanded from env, defaults to false.
	Expand bool

	// Separator to be used for slice/array items, defaults to " ".
	Separator string

	// Format to be used for parsing time strings, defaults to time.RFC3339.
	Format string
}

// Returns default config tags.
func extractTags(st reflect.StructTag) *ConfigTags {
	tags := ConfigTags{
		Config:    st.Get("config"),
		Default:   st.Get("default"),
		Json:      extractKeyName(st.Get("json")),
		Yaml:      extractKeyName(st.Get("yaml")),
		Toml:      extractKeyName(st.Get("toml")),
		Required:  st.Get("required") == "true",
		Ignore:    st.Get("ignore") == "true",
		Expand:    st.Get("expand") == "true",
		Separator: st.Get("separator"),
		Format:    st.Get("format"),
	}

	if tags.Config == ignoreCharacter {
		tags.Ignore = true
	}
	if tags.Separator == "" {
		tags.Separator = defaultSeparator
	}
	if tags.Format == "" {
		tags.Format = time.RFC3339
	}

	return &tags
}

// It extracts name of the key from file tag, ignoring options
// e.g. calling with "field,omitempty" would return "field"
func extractKeyName(key string) string {
	slice := strings.Split(key, ",")
	if len(slice) == 0 {
		return ""
	}

	return slice[0]
}
