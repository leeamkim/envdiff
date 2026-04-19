package diff

// MultiCompareResult holds pairwise comparison results across multiple envs.
type MultiCompareResult struct {
	EnvNames []string
	Pairs    map[string]Result
}

// PairKey returns a canonical key for a pair of env names.
func PairKey(a, b string) string {
	return a + ".." + b
}

// CompareMulti compares all pairs of envs in the provided map.
func CompareMulti(envs map[string]map[string]string) MultiCompareResult {
	names := make([]string, 0, len(envs))
	for n := range envs {
		names = append(names, n)
	}
	sortStrings(names)

	pairs := make(map[string]Result)
	for i := 0; i < len(names); i++ {
		for j := i + 1; j < len(names); j++ {
			a, b := names[i], names[j]
			pairs[PairKey(a, b)] = Compare(envs[a], envs[b])
		}
	}
	return MultiCompareResult{EnvNames: names, Pairs: pairs}
}

// HasAnyDiff returns true if any pair has a diff.
func (m MultiCompareResult) HasAnyDiff() bool {
	for _, r := range m.Pairs {
		if r.HasDiff() {
			return true
		}
	}
	return false
}
