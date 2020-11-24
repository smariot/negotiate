package negotiate

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
)

func ExampleContentType() {
	middleware := ContentTypeMiddleware("image/png", "image/webp")

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Negotiated type is %s\n", ContentType(r))
	}))

	accepts := []string{"", "*/*, IMAGE/WEBP", "image/jpeg", "image/*, image/png; q=0", "i like waffles."}

	for _, accept := range accepts {
		r, _ := http.NewRequest("GET", "/foo", nil)
		r.Header.Set("Accept", accept)

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		fmt.Printf("Accept=%q\n", accept)
		io.Copy(os.Stdout, w.Result().Body)
		fmt.Println()
	}

	// Output:
	// Accept=""
	// Negotiated type is image/png
	//
	// Accept="*/*, IMAGE/WEBP"
	// Negotiated type is image/webp
	//
	// Accept="image/jpeg"
	// 406: Not Acceptable
	// Supported values for Accept header are: image/png, image/webp
	//
	// Accept="image/*, image/png; q=0"
	// Negotiated type is image/webp
	//
	// Accept="i like waffles."
	// 400: Bad Request
	// Unable to parse Accept header.
}
