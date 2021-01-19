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

func isStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Struct && !isTime(t) && !isURL(t)
}

func isDuration(t reflect.Type) bool {
	return t.PkgPath() == "time" && t.Name() == "Duration"
}

func isTime(t reflect.Type) bool {
	return t.PkgPath() == "time" && t.Name() == "Time"
}

func isURL(t reflect.Type) bool {
	return t.PkgPath() == "net/url" && t.Name() == "URL"
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

// extractItems splits and trims input string based on separator
func extractItems(str string, sep string) []string {
	var items []string
	for _, v := range strings.Split(str, sep) {
		item := strings.TrimSpace(v)
		if len(item) > 0 {
			items = append(items, item)
		}
	}

	return items
}
