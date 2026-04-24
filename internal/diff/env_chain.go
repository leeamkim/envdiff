package diff

import "fmt"

// ChainEntry represents a single key's value propagation across an ordered chain of envs.
type ChainEntry struct {
	Key    string
	Values []ChainValue // ordered by chain position
}

// ChainValue holds the value of a key in one environment of the chain.
type ChainValue struct {
	Env      string
	Value    string
	Present  bool
	Override bool // true if this env's value differs from the previous present value
}

// ChainResult is the output of BuildChain.
type ChainResult struct {
	Chain   []string
	Entries []ChainEntry
}

// BuildChain propagates values across an ordered list of env maps.
// Later envs override earlier ones. Each key's journey is tracked.
func BuildChain(chain []string, envs map[string]map[string]string) ChainResult {
	// collect all keys
	keySet := map[string]struct{}{}
	for _, name := range chain {
		for k := range envs[name] {
			keySet[k] = struct{}{}
		}
	}

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sortStrings(keys)

	entries := make([]ChainEntry, 0, len(keys))
	for _, key := range keys {
		var prevValue string
		hasPrev := false
		values := make([]ChainValue, 0, len(chain))
		for _, envName := range chain {
			v, ok := envs[envName][key]
			cv := ChainValue{Env: envName, Value: v, Present: ok}
			if ok {
				cv.Override = hasPrev && v != prevValue
				prevValue = v
				hasPrev = true
			}
			values = append(values, cv)
		}
		entries = append(entries, ChainEntry{Key: key, Values: values})
	}

	return ChainResult{Chain: chain, Entries: entries}
}

// FormatChainReport returns a human-readable report of the chain result.
func FormatChainReport(r ChainResult) string {
	if len(r.Entries) == 0 {
		return "no keys found in chain\n"
	}
	out := ""
	for _, e := range r.Entries {
		out += fmt.Sprintf("[%s]\n", e.Key)
		for _, cv := range e.Values {
			if !cv.Present {
				out += fmt.Sprintf("  %-16s (missing)\n", cv.Env)
			} else if cv.Override {
				out += fmt.Sprintf("  %-16s %s  (override)\n", cv.Env, cv.Value)
			} else {
				out += fmt.Sprintf("  %-16s %s\n", cv.Env, cv.Value)
			}
		}
	}
	return out
}
