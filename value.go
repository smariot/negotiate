package negotiate

// Value represents one of the underlying values being negotiated for.
type Value interface {
	// String should function as the inverse of the ValueParser that was used to create the Value.
	String() string

	// Specificity returns an arbitrary value that only has meaning when compared to values of the same type.
	// Larger values are more specific.
	Specificity() int

	// Satisfies should return true if this value satisfies the passed value.
	//
	// Should panic if value is a different type of value.
	Satisfies(Value) bool
}

// ValueParser converts a string into a Value.
//
// The parser must interpret "*" to represent a wildcard value without returning an error.
type ValueParser func(string) (Value, error)

// Must takes a value and an error, and returns the value.
//
// If the error is not nil, it panics.
func Must(value Value, err error) Value {
	if err != nil {
		panic(err)
	}

	return value
}
