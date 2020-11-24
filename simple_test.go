package negotiate

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
)

func ExampleCharset() {
	middleware := CharsetMiddleware("UTF-8", "US-ASCII")

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Negotiated charset is %s\n", Charset(r))
	}))

	charsets := []string{"", "*, us-ascii", "koi8-r", "i like waffles."}

	for _, charset := range charsets {
		r, _ := http.NewRequest("GET", "/", nil)
		r.Header.Set("Accept-Charset", charset)

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		fmt.Printf("Accept-Charset=%q\n", charset)
		io.Copy(os.Stdout, w.Result().Body)
		fmt.Println()
	}

	// Output:
	// Accept-Charset=""
	// Negotiated charset is UTF-8
	//
	// Accept-Charset="*, us-ascii"
	// Negotiated charset is US-ASCII
	//
	// Accept-Charset="koi8-r"
	// 406: Not Acceptable
	// Supported values for Accept-Charset header are: UTF-8, US-ASCII
	//
	// Accept-Charset="i like waffles."
	// 400: Bad Request
	// Unable to parse Accept-Charset header.
}

func ExampleEncoding() {
	middleware := EncodingMiddleware("identity", "gzip")

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Negotiated encoding is %s\n", Encoding(r))
	}))

	encodings := []string{"", "*", "GZIP", "br", "i like waffles."}

	for _, encoding := range encodings {
		r, _ := http.NewRequest("GET", "/foo", nil)
		r.Header.Set("Accept-Encoding", encoding)

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		fmt.Printf("Accept-Encoding=%q\n", encoding)
		io.Copy(os.Stdout, w.Result().Body)
		fmt.Println()
	}

	// Output:
	// Accept-Encoding=""
	// Negotiated encoding is identity
	//
	// Accept-Encoding="*"
	// Negotiated encoding is identity
	//
	// Accept-Encoding="GZIP"
	// Negotiated encoding is gzip
	//
	// Accept-Encoding="br"
	// 406: Not Acceptable
	// Supported values for Accept-Encoding header are: identity, gzip
	//
	// Accept-Encoding="i like waffles."
	// 400: Bad Request
	// Unable to parse Accept-Encoding header.
}
