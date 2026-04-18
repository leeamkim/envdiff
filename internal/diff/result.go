package diff

// Result holds the outcome of comparing two env maps.
type Result struct {
	// MissingInB contains keys present in A but absent in B, with A's value.
	MissingInB map[string]string
	// MissingInA contains keys present in B but absent in A, with B's value.
	MissingInA map[string]string
	// Mismatched contains keys present in both but with differing values.
	// The array is [valueA, valueB].
	Mismatched map[string][2]string
}

// NewResult initialises a Result with empty maps.
func NewResult() Result {
	return Result{
		MissingInB: make(map[string]string),
		MissingInA: make(map[string]string),
		Mismatched: make(map[string][2]string),
	}
}

// HasDiff returns true when any differences exist.
func (r Result) HasDiff() bool {
	return len(r.MissingInA) > 0 || len(r.MissingInB) > 0 || len(r.Mismatched) > 0
}
