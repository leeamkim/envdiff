package diff

// Result holds the categorised diff between two env maps.
type Result struct {
	MissingInA map[string]string
	MissingInB map[string]string
	Mismatched map[string][2]string
	Matching   map[string]string
}

// NewResult allocates an empty Result.
func NewResult() Result {
	return Result{
		MissingInA: make(map[string]string),
		MissingInB: make(map[string]string),
		Mismatched: make(map[string][2]string),
		Matching:   make(map[string]string),
	}
}
