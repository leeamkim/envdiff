package diff

import "fmt"

// AliasMap maps canonical key names to one or more aliases.
type AliasMap map[string][]string

// AliasIssue represents a key found under an alias instead of its canonical name.
type AliasIssue struct {
	Canonical string
	FoundAs   string
	Value     string
}

func (a AliasIssue) String() string {
	return fmt.Sprintf("key %q found as alias %q (value: %q)", a.Canonical, a.FoundAs, a.Value)
}

// ResolveAliases checks env for keys that are stored under an alias rather than
// their canonical name. It returns a deduplicated map with canonical keys and a
// list of issues describing each substitution made.
func ResolveAliases(env map[string]string, aliases AliasMap) (map[string]string, []AliasIssue) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	var issues []AliasIssue
	for canonical, aliasList := range aliases {
		if _, ok := out[canonical]; ok {
			continue
		}
		for _, alias := range aliasList {
			if val, ok := out[alias]; ok {
				issue := AliasIssue{Canonical: canonical, FoundAs: alias, Value: val}
				issues = append(issues, issue)
				out[canonical] = val
				delete(out, alias)
				break
			}
		}
	}
	return out, issues
}

// ParseAliasMap parses a slice of "canonical=alias1,alias2" strings into an AliasMap.
func ParseAliasMap(entries []string) (AliasMap, error) {
	am := make(AliasMap)
	for _, entry := range entries {
		canonical, rest, ok := splitOnce(entry, '=')
		if !ok || canonical == "" || rest == "" {
			return nil, fmt.Errorf("invalid alias entry: %q", entry)
		}
		aliases := splitCSVAlias(rest)
		am[canonical] = append(am[canonical], aliases...)
	}
	return am, nil
}

func splitOnce(s string, sep byte) (string, string, bool) {
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			return s[:i], s[i+1:], true
		}
	}
	return "", "", false
}

func splitCSVAlias(s string) []string {
	var out []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			if tok := s[start:i]; tok != "" {
				out = append(out, tok)
			}
			start = i + 1
		}
	}
	return out
}
