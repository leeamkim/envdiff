package diff

import "strings"

// TransformFunc is a function that transforms a value.
type TransformFunc func(string) string

// TransformOption configures a Transform call.
type TransformOption struct {
	Keys []string // if empty, applies to all keys
	Fn   TransformFunc
}

// TransformResult holds the transformed map and a log of changes.
type TransformResult struct {
	Out     map[string]string
	Changed []string
}

// Transform applies a set of TransformOptions to the given env map.
// A copy of the map is returned; the original is not mutated.
func Transform(env map[string]string, opts []TransformOption) TransformResult {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}
	changedSet := map[string]struct{}{}
	for _, opt := range opts {
		if len(opt.Keys) == 0 {
			for k, v := range out {
				nv := opt.Fn(v)
				if nv != v {
					out[k] = nv
					changedSet[k] = struct{}{}
				}
			}
		} else {
			for _, k := range opt.Keys {
				if v, ok := out[k]; ok {
					nv := opt.Fn(v)
					if nv != v {
						out[k] = nv
						changedSet[k] = struct{}{}
					}
				}
			}
		}
	}
	changed := make([]string, 0, len(changedSet))
	for k := range changedSet {
		changed = append(changed, k)
	}
	sortStrings(changed)
	return TransformResult{Out: out, Changed: changed}
}

// TransformTrimSpace returns a TransformFunc that trims whitespace.
func TransformTrimSpace() TransformFunc { return strings.TrimSpace }

// TransformUppercase returns a TransformFunc that uppercases a value.
func TransformUppercase() TransformFunc { return strings.ToUpper }

// TransformLowercase returns a TransformFunc that lowercases a value.
func TransformLowercase() TransformFunc { return strings.ToLower }

func sortStrings(ss []string) {
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && ss[j] < ss[j-1]; j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}
