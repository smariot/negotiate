package negotiate

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type simpleValue string

func (v simpleValue) String() string {
	return string(v)
}

func (v simpleValue) Specificity() int {
	if v == "*" {
		return 0
	}

	return 1
}

func (v simpleValue) Satisfies(_ref Value) bool {
	ref := _ref.(simpleValue)

	if ref != "*" && v != ref {
		return false
	}

	return true
}

var reSimple = regexp.MustCompile(`^(?:\*|[a-zA-Z0-9_\-]+)$`)

// ParseSimple parses a simple item and returns a Value.
//
// The string is treated as case insensitive, and must consist of one or more alpha-numeric characters, or an underscore, or a dash.
//
// A single "*" is also acceptable, and treated as a wildcard that matches anything.
func ParseSimple(str string) (Value, error) {
	if !reSimple.MatchString(str) {
		return nil, fmt.Errorf("invalid simple item: %q", str)
	}

	return simpleValue(strings.ToLower(str)), nil
}

// CharsetMiddleware is shorthand for Middleware("Accept-Charset", ParseSimple, items...)
func CharsetMiddleware(items ...string) func(http.Handler) http.Handler {
	return Middleware("Accept-Charset", ParseSimple, items...)
}

// Charset is shorthand for Item(r, "Accept-Charset")
func Charset(r *http.Request) string {
	return Item(r, "Accept-Charset")
}

// EncodingMiddleware is shorthand for Middleware("Accept-Encoding", ParseSimple, items...)
//
// Note: This page has a list of chracter sets and their aliases.
func EncodingMiddleware(items ...string) func(http.Handler) http.Handler {
	return Middleware("Accept-Encoding", ParseSimple, items...)
}

// Encoding is shorthand for Item(r, "Accept-Encoding")
func Encoding(r *http.Request) string {
	return Item(r, "Accept-Encoding")
}
