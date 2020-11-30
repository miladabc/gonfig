package gonfig

import (
	"regexp"
	"strings"
)

var (
	firstCapRegex = regexp.MustCompile("([A-Z])([A-Z][a-z])")
	allCapRegex   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func toScreamingSnakeCase(in []string) string {
	s := strings.Join(in, "_")
	out := firstCapRegex.ReplaceAllString(s, "${1}_${2}")
	out = allCapRegex.ReplaceAllString(out, "${1}_${2}")
	out = strings.ReplaceAll(out, "-", "_")

	return strings.ToUpper(out)
}
