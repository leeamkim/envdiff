package diff

import "fmt"

// PromoteResult holds the outcome of promoting keys from one env to another.
type PromoteResult struct {
	Added    map[string]string
	Skipped  map[string]string
	Conflicts map[string]string
}

// PromoteOption controls promotion behaviour.
type PromoteOption func(*promoteConfig)

type promoteConfig struct {
	overwrite bool
	keys      []string
}

// PromoteOverwrite allows existing keys in dest to be overwritten.
func PromoteOverwrite() PromoteOption {
	return func(c *promoteConfig) { c.overwrite = true }
}

// PromoteKeys limits promotion to specific keys.
func PromoteKeys(keys ...string) PromoteOption {
	return func(c *promoteConfig) { c.keys = keys }
}

// Promote copies keys from src into dest, returning a PromoteResult.
func Promote(src, dest map[string]string, opts ...PromoteOption) PromoteResult {
	cfg := &promoteConfig{}
	for _, o := range opts {
		o(cfg)
	}

	result := PromoteResult{
		Added:     make(map[string]string),
		Skipped:   make(map[string]string),
		Conflicts: make(map[string]string),
	}

	keys := cfg.keys
	if len(keys) == 0 {
		for k := range src {
			keys = append(keys, k)
		}
	}

	for _, k := range keys {
		v, ok := src[k]
		if !ok {
			result.Skipped[k] = ""
			continue
		}
		if existing, exists := dest[k]; exists {
			if !cfg.overwrite {
				result.Conflicts[k] = existing
				continue
			}
		}
		dest[k] = v
		result.Added[k] = v
	}
	return result
}

// FormatPromoteResult returns a human-readable summary.
func FormatPromoteResult(r PromoteResult) string {
	return fmt.Sprintf("promoted=%d conflicts=%d skipped=%d",
		len(r.Added), len(r.Conflicts), len(r.Skipped))
}
