package negotiate

import (
	"fmt"
)

func ExampleNegotiate() {
	negotiate := Make(ParseSimple, "cake", "pie")

	for _, query := range []string{"", "*, CAKE;q=0.9", "pizza", "what is this?"} {
		if item, err := negotiate.Process(query); err != nil {
			fmt.Printf("%q -> error: %v\n", query, err)
		} else {
			fmt.Printf("%q -> %s\n", query, item)
		}
	}

	// Output:
	// "" -> cake
	// "*, CAKE;q=0.9" -> pie
	// "pizza" -> error: no item satisfies query
	// "what is this?" -> error: invalid simple item: "what is this?"
}
