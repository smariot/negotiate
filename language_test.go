package negotiate

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func ExampleLanguage() {
	middleware := LanguageMiddleware("en-CA", "en-US")

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Negotiated language is %s\n", Language(r))
	}))

	languages := []string{"", "en", "*, en, en-us", "fr-fr", "i like waffles."}

	for _, language := range languages {
		r, _ := http.NewRequest("GET", "/foo", nil)
		r.Header.Set("Accept-Language", language)

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		fmt.Printf("Accept-Language=%q\n", language)
		io.Copy(os.Stdout, w.Result().Body)
		fmt.Println()
	}

	// Output:
	// Accept-Language=""
	// Negotiated language is en-CA
	//
	// Accept-Language="en"
	// Negotiated language is en-CA
	//
	// Accept-Language="*, en, en-us"
	// Negotiated language is en-US
	//
	// Accept-Language="fr-fr"
	// 406: Not Acceptable
	// Supported values for Accept-Language header are: en-CA, en-US
	//
	// Accept-Language="i like waffles."
	// 400: Bad Request
	// Unable to parse Accept-Language header.
}

func TestParseLocale(t *testing.T) {
	tests := []struct {
		locale  string
		wantString string
		wantSpecificity int
		wantErr bool
	}{
		{"*", "*", 0, false},
		{"en", "en", 1, false},
		{"EN", "en", 1, false},
		{"en-ca", "en-CA", 2, false},
		{"EN-CA", "en-CA", 2, false},
		{"", "", 0, true},
		{"what is this", "", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.locale, func(t *testing.T) {
			got, err := ParseLocale(tt.locale)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLocale(%q) error = %v, wantErr %v", tt.locale, err, tt.wantErr)
				return
			}

			if err == nil {
				if s := got.String(); s != tt.wantString {
					t.Errorf("ParseLocale(%q).String() = %q, want %q", tt.locale, s, tt.wantString)
				}

				if s := got.Specificity(); s != tt.wantSpecificity {
					t.Errorf("ParseLocale(%q).Specificity() = %d, want %d", tt.locale, s, tt.wantSpecificity)
				}
			}
		})
	}
}
