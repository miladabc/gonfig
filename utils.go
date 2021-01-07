package gonfig

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var (
	firstCapRegex = regexp.MustCompile("([A-Z])([A-Z][a-z])")
	allCapRegex   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

// toSnakeCase converts input string into snake_case form
func toSnakeCase(s string) string {
	out := firstCapRegex.ReplaceAllString(s, "${1}_${2}")
	out = allCapRegex.ReplaceAllString(out, "${1}_${2}")
	out = strings.ReplaceAll(out, "-", "_")

	return out
}

func isDuration(v reflect.Value) bool {
	return v.Type().PkgPath() == "time" && v.Type().Name() == "Duration"
}

func isTime(v reflect.Value) bool {
	return v.Type().PkgPath() == "time" && v.Type().Name() == "Time"
}

func isURL(v reflect.Value) bool {
	return v.Type().PkgPath() == "net/url" && v.Type().Name() == "URL"
}

// traverseMap finds a value in a map based on provided path
func traverseMap(m map[string]interface{}, path []string) (string, bool) {
	if len(path) == 0 {
		return "", false
	}
	first, path := path[0], path[1:]

	value, exists := m[first]
	if !exists {
		value, exists = m[strings.ToLower(first)]
		if !exists {
			return "", false
		}
	}

	if len(path) == 0 {
		return fmt.Sprint(value), true
	}

	nestedMap, ok := value.(map[string]interface{})
	if !ok {
		return "", false
	}

	return traverseMap(nestedMap, path)
}
