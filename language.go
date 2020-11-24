package negotiate

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type localeValue struct {
	language, territory string
}

func (l localeValue) String() string {
	if l.territory == "" {
		return l.language
	}

	return l.language+"-"+l.territory
}

func (l localeValue) Specificity() int {
	if l.language == "*" {
		return 0
	}

	if l.territory == "" {
		return 1
	}

	return 2
}

func (l localeValue) Satisfies(_ref Value) bool {
	ref := _ref.(localeValue)

	if ref.language != "*" && l.language != ref.language {
		return false
	}

	if ref.territory != "" && l.territory != ref.territory {
		return false
	}

	return true
}

var reLocale = regexp.MustCompile(`^(\*|[a-zA-Z0-9_]+)(?:-([a-zA-Z0-9_]+))?$`)

// ParseLocale parses a locale and returns a Value.
func ParseLocale(locale string) (Value, error) {
	match := reLocale.FindStringSubmatch(locale)
	if match == nil {
		return nil, fmt.Errorf("bad locale: %q", locale)
	}

	return localeValue{strings.ToLower(match[1]), strings.ToUpper(match[2])}, nil
}

// LanguageMiddleware is shorthand for Middleware("Accept-Language", ParseLocale, items...)
func LanguageMiddleware(items ...string) func(http.Handler) http.Handler {
	return Middleware("Accept-Language", ParseLocale, items...)
}

// Language is shorthand for Item(r, "Accept-Charset")
func Language(r *http.Request) string {
	return Item(r, "Accept-Language")
}
