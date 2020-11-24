package negotiate

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// QValue is a Value with an associated quality.
type QValue struct {
	Value
	Q float64
}

func (v QValue) String() string {
	// The weight is normalized to a real number in the range 0 through 1
	q := v.Q
	switch true {
	case q >= 1:
		// If no "q" parameter is present, the default weight is 1.
		return v.Value.String()
	case q > 0:
		// In range, do nothing.
	default:
		// Negative numbers, NaN.
		q = 0
	}

	// MUST NOT generate more than three digits after the decimal point.
	return fmt.Sprintf("%s; q=%.3g", v.Value.String(), q)
}

// Query is a slice of Value/Quality pairs.
type Query []QValue

func (q Query) Len() int {
	return len(q)
}

// Less enables you to sort the items in a query by precedence using sort.Stable(q).
//
// It returns returns true if the q[i] has a higher precedence than q[j].
func (q Query) Less(i, j int) bool {
	if a, b := q[i].Specificity(), q[j].Specificity(); a != b {
		return b < a
	}

	return q[i].Q > q[i].Q
}

func (q Query) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q Query) String() string {
	var b strings.Builder

	for i, v := range q {
		if i != 0 {
			b.WriteString(", ")
		}

		b.WriteString(v.String())
	}

	return b.String()
}

// Find finds the first value in the query that is satisfied by the given value.
//
// If that value has a positive quality, its index is returned,
// otherwise -1 is returned.
//
// If no value in the query is satisfied by the value, -1 is returned.
func (q Query) Find(v Value) int {
	for i, qv := range q {
		if v.Satisfies(qv.Value) {
			if qv.Q > 0 {
				return i
			}

			// a value of 0 means "not acceptable"
			break
		}
	}

	return -1
}

// Choose returns the index of the best value in the given list of choices,
// or -1 if none of the choices satisfy the query.
//
// "best" is the choice that yields the highest quality value.
// In the case of a tie, the query item with the higher precedence is used.
// If that query item can be satisfied by more than once choice, the one
// that appears first in the choices list is used.
func (q Query) Choose(choices []Value) int {
	var (
		bestChoiceIndex = -1
		bestQueryIndex  int
	)

	for choiceIndex, cv := range choices {
		if queryIndex := q.Find(cv); queryIndex != -1 {
			if bestChoiceIndex == -1 || q[queryIndex].Q > q[bestQueryIndex].Q || (q[queryIndex].Q == q[bestQueryIndex].Q && queryIndex < bestQueryIndex) {
				bestChoiceIndex, bestQueryIndex = choiceIndex, queryIndex
			}
		}
	}

	return bestChoiceIndex
}

//
var reItem = regexp.MustCompile(`^([^;]+?(?:;(?:[^;"]|"(?:[^"\\]|\\.)*")*)*?)(?:;\s*[qQ]=(1(?:\.0{0,3})?|0(?:\.[0-9]{0,3})?)\s*(?:;.*)?)?$`)

// ParseQuery parses a query, and returns the query with its values sorted by precedence.
//
// Extension parameters are not supported, and will be silently discarded if present.
//
// An empty query will be satisfied by anything.
func ParseQuery(parser ValueParser, query string) (q Query, err error) {
	if strings.TrimSpace(query) == "" {
		query = "*"
	}

	items := strings.Split(query, ",")
	q = make(Query, len(items))

	for i, item := range items {
		match := reItem.FindStringSubmatch(item)

		q[i].Q = 1.0
		q[i].Value, err = parser(strings.TrimSpace(match[1]))
		if err != nil {
			return nil, err
		}

		if value, err := strconv.ParseFloat(match[2], 10); err == nil {
			q[i].Q = value
		}
	}

	sort.Stable(q)

	return
}
