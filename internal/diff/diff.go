package diff

// Result holds the comparison result between two env files.
type Result struct {
	MissingInRight []string
	MissingInLeft  []string
	Mismatched     map[string][2]string // key -> [leftVal, rightVal]
}

// Compare compares two parsed env maps and returns a Result.
func Compare(left, right map[string]string) Result {
	r := Result{
		Mismatched: make(map[string][2]string),
	}

	for k, lv := range left {
		if rv, ok := right[k]; !ok {
			r.MissingInRight = append(r.MissingInRight, k)
		} else if lv != rv {
			r.Mismatched[k] = [2]string{lv, rv}
		}
	}

	for k := range right {
		if _, ok := left[k]; !ok {
			r.MissingInLeft = append(r.MissingInLeft, k)
		}
	}

	return r
}

// HasDiff returns true if there are any differences.
func (r Result) HasDiff() bool {
	return len(r.MissingInRight) > 0 || len(r.MissingInLeft) > 0 || len(r.Mismatched) > 0
}
