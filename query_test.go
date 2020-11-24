package negotiate

import (
	"fmt"
	"math"
	"testing"
)

func ExampleQuery() {
	// Example from rfc7231 page 40

	q, _ := ParseQuery(ParseMedia, "text/*;q=0.3, text/html;q=0.7, text/html;level=1, text/html;level=2;q=0.4, */*;q=0.5")

	values := []Value{
		Must(ParseMedia("text/html;level=1")),
		Must(ParseMedia("text/html")),
		Must(ParseMedia("text/plain")),
		Must(ParseMedia("image/jpeg")),
		Must(ParseMedia("text/html;level=2")),
		Must(ParseMedia("text/html;level=3")),
	}

	fmt.Println("Query:", q)
	fmt.Println()
	fmt.Println("Media Type           Quality Value")
	for _, value := range values {
		fmt.Printf("%-20s %g\n", value, q[q.Find(value)].Q)
	}

	// Output:
	// Query: text/html; level=1, text/html; level=2; q=0.4, text/html; q=0.7, text/*; q=0.3, */*; q=0.5
	//
	// Media Type           Quality Value
	// text/html; level=1   1
	// text/html            0.7
	// text/plain           0.3
	// image/jpeg           0.5
	// text/html; level=2   0.4
	// text/html; level=3   0.7
}

func TestQValue_String(t *testing.T) {
	tests := []struct {
		name   string
		value QValue
		want   string
	}{
		// These have legal q values.
		{"q = 1", QValue{simpleValue("test"), 1}, "test"},
		{"q = 0", QValue{simpleValue("test"), 0}, "test; q=0"},
		{"q = 0.5", QValue{simpleValue("test"), 0.5}, "test; q=0.5"},

		// These have illegal q values.
		{"q = 0.12345", QValue{simpleValue("test"), 0.12345}, "test; q=0.123"},
		{"q > 1", QValue{simpleValue("test"), 2}, "test"},
		{"q < 0", QValue{simpleValue("test"), -1}, "test; q=0"},
		{"q = +Inf", QValue{simpleValue("test"), math.Inf(1)}, "test"},
		{"q = -Inf", QValue{simpleValue("test"), math.Inf(-1)}, "test; q=0"},
		{"q = NaN", QValue{simpleValue("test"), math.NaN()}, "test; q=0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.String(); got != tt.want {
				t.Errorf("QValue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
