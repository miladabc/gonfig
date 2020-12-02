package gonfig

import "reflect"

const (
	defaultSeparator = " "
	ignoreCharacter  = "-"
)

// all possible useful tags
type ConfigTags struct {
	// config key name, use "-" to ignore, defaults to field name
	Config string

	// default value for field
	Default string

	// specify if value should be present, defaults to false
	Required bool

	// specify if field should be ignored, defaults to false
	Ignore bool

	// specify if value should be expanded from env, defaults to false
	Expand bool

	// separator to be used for slice/array items, defaults to " "
	Separator string

	// format to be used for parsing time strings, defaults to time.RFC3339
	Format string
}

func getTags(st reflect.StructTag) *ConfigTags {
	tags := ConfigTags{
		Config:    st.Get("config"),
		Default:   st.Get("default"),
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

	return &tags
}
