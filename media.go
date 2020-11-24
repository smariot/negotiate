package negotiate

import (
	"mime"
	"net/http"
	"strings"
)

type mediaValue struct {
	major, minor string
	params       map[string]string
}

func (m mediaValue) String() string {
	return mime.FormatMediaType(m.major+"/"+m.minor, m.params)
}

func (m mediaValue) Specificity() int {
	if m.major == "*" {
		return 0
	}

	if m.minor == "*" {
		return 1
	}

	if len(m.params) == 0 {
		return 2
	}

	return 3
}

func (m mediaValue) Satisfies(_ref Value) bool {
	ref := _ref.(mediaValue)

	if ref.major != "*" && m.major != ref.major {
		return false
	}

	if ref.minor != "*" && m.minor != ref.minor {
		return false
	}

	for key, refValue := range ref.params {
		if value, ok := m.params[key]; !ok || value != refValue {
			return false
		}
	}

	return true
}

// ParseMedia parses a media type and returns a Value.
func ParseMedia(mediaStr string) (Value, error) {
	media, params, err := mime.ParseMediaType(mediaStr)

	if err != nil {
		return nil, err
	}

	idx := strings.IndexByte(media, '/')

	if idx >= 0 {
		return mediaValue{media[:idx], media[idx+1:], params}, nil
	}

	return mediaValue{media, "*", params}, nil
}

// ContentTypeMiddleware is shorthand for Middleware("Accept", ParseMedia, items...)
func ContentTypeMiddleware(items ...string) func(http.Handler) http.Handler {
	return Middleware("Accept", ParseMedia, items...)
}

// ContentType is shorthand for Item(r, "Accept")
func ContentType(r *http.Request) string {
	return Item(r, "Accept")
}
