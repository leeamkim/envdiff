package diff

// FilterOptions controls which diff results are included.
type FilterOptions struct {
	OnlyMissing    bool
	OnlyMismatched bool
}

// Filter returns a subset of Result based on FilterOptions.
func Filter(result Result, opts FilterOptions) Result {
	if !opts.OnlyMissing && !opts.OnlyMismatched {
		return result
	}

	filtered := Result{
		MissingInB:  make(map[string]string),
		MissingInA:  make(map[string]string),
		Mismatched:  make(map[string][2]string),
	}

	if opts.OnlyMissing || !opts.OnlyMismatched {
		for k, v := range result.MissingInB {
			filtered.MissingInB[k] = v
		}
		for k, v := range result.MissingInA {
			filtered.MissingInA[k] = v
		}
	}

	if opts.OnlyMismatched || !opts.OnlyMissing {
		for k, v := range result.Mismatched {
			filtered.Mismatched[k] = v
		}
	}

	return filtered
}
