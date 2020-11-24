package negotiate

import (
	"context"
	"net/http"
)

type handler struct {
	negotiate Negotiate
	header    string
	next      http.Handler
}

type ctxKey string

// Item returns the item negotiated for in a header processed by Middleware().
func Item(r *http.Request, header string) string {
	item, _ := r.Context().Value(ctxKey(http.CanonicalHeaderKey(header))).(string)
	return item
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Vary", h.header)

	switch value, err := h.negotiate.Process(r.Header.Get(h.header)); err {
	case nil:
		h.next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKey(h.header), value)))
	case ErrNotAcceptable:
		http.Error(w,
			"406: Not Acceptable\nSupported values for "+h.header+" header are: "+h.negotiate.String(),
			http.StatusNotAcceptable)
	default:
		http.Error(w,
			"400: Bad Request\nUnable to parse "+h.header+" header.",
			http.StatusBadRequest)
	}
}

// Middleware returns http middleware for negotiating on an http header.
//
// This function will panic if any of the passed items fail to parse.
//
// If a matching item is found, the Vary header will be added to the response, and the next handler will be invoked.
// The matching item can be retrieved using Item(r, header).
//
// A 406: Not Acceptable error will be generated if no items match.
//
// A 400: Bad Request error will be generated if any item in the given header fails to parse.
func Middleware(header string, parser ValueParser, items ...string) func(http.Handler) http.Handler {
	header = http.CanonicalHeaderKey(header)
	negotiate := Make(parser, items...)

	return func(next http.Handler) http.Handler {
		return handler{negotiate, header, next}
	}
}
