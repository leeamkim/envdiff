package diff

import "strings"

// IgnoreList holds a set of keys to exclude from diff results.
type IgnoreList struct {
	keys map[string]struct{}
}

// NewIgnoreList creates an IgnoreList from a slice of key names.
func NewIgnoreList(keys []string) *IgnoreList {
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[strings.TrimSpace(k)] = struct{}{}
	}
	return &IgnoreList{keys: m}
}

// Contains reports whether the key is in the ignore list.
func (il *IgnoreList) Contains(key string) bool {
	_, ok := il.keys[key]
	return ok
}

// Apply removes ignored keys from a Result.
func (il *IgnoreList) Apply(r Result) Result {
	out := Result{
		MissingInA:  []string{},
		MissingInB:  []string{},
		Mismatched:  map[string][2]string{},
	}
	for _, k := range r.MissingInA {
		if !il.Contains(k) {
			out.MissingInA = append(out.MissingInA, k)
		}
	}
	for _, k := range r.MissingInB {
		if !il.Contains(k) {
			out.MissingInB = append(out.MissingInB, k)
		}
	}
	for k, v := range r.Mismatched {
		if !il.Contains(k) {
			out.Mismatched[k] = v
		}
	}
	return out
}
