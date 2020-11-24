// This package provides http middlewares for negotiating content-type, encoding, language, and charset.

package negotiate

import (
	"errors"
	"strings"
)

// Negotiate holds a list of items that can be negotiated for,
// and provides a Process() method for selecting one of those items for a given query.
type Negotiate struct {
	parser ValueParser
	items  []string
	values []Value
}

// Make returns a Negotiate object for the given items.
//
// The items should be ordered with the more compatible ones first, so that they
// will be preferred for wildcard matches.
func Make(parser ValueParser, items ...string) Negotiate {
	values := make([]Value, len(items))

	for i, item := range items {
		values[i] = Must(parser(item))
	}

	return Negotiate{parser, items, values}
}

// String returns all of the items as a comma seperated list.
func (n Negotiate) String() string {
	return strings.Join(n.items, ", ")
}

// ErrNotAcceptable is returned when no item satisfies a query.
var ErrNotAcceptable = errors.New("no item satisfies query")

// Process returns the item with the highest quality satisfying query,
// prefering the earlier items given to New in the event of a tie.
//
// Returns ErrNotAcceptable if no item satisfies the query.
//
// If the value parser returns an error, that error will be returned.
func (n Negotiate) Process(query string) (item string, err error) {
	q, err := ParseQuery(n.parser, query)

	if err != nil {
		return "", err
	}

	if i := q.Choose(n.values); i != -1 {
		return n.items[i], nil
	}

	return "", ErrNotAcceptable
}
