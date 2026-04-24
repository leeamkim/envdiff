package diff

import (
	"fmt"
	"strings"
)

// InterpolateIssue describes a variable reference that could not be resolved.
type InterpolateIssue struct {
	Key string
	Ref string // the unresolved ${REF} token
}

func (i InterpolateIssue) String() string {
	return fmt.Sprintf("key %q references unresolved variable %q", i.Key, i.Ref)
}

// Interpolate expands ${VAR} references in env values using the same map.
// It returns a new map with values expanded and a list of unresolved references.
func Interpolate(env map[string]string) (map[string]string, []InterpolateIssue) {
	result := make(map[string]string, len(env))
	var issues []InterpolateIssue

	for k, v := range env {
		expanded, unresolved := expandValue(v, env)
		result[k] = expanded
		for _, ref := range unresolved {
			issues = append(issues, InterpolateIssue{Key: k, Ref: ref})
		}
	}
	return result, issues
}

// FormatInterpolateIssues returns a human-readable summary of interpolation issues.
func FormatInterpolateIssues(issues []InterpolateIssue) string {
	if len(issues) == 0 {
		return "no interpolation issues"
	}
	var sb strings.Builder
	for _, iss := range issues {
		sb.WriteString("  " + iss.String() + "\n")
	}
	return sb.String()
}

// expandValue replaces ${VAR} tokens in s using env. Returns expanded string
// and a slice of token names that were not found in env.
func expandValue(s string, env map[string]string) (string, []string) {
	var unresolved []string
	result := s
	for {
		start := strings.Index(result, "${") 
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], "}")
		if end == -1 {
			break
		}
		end += start
		token := result[start+2 : end]
		if val, ok := env[token]; ok {
			result = result[:start] + val + result[end+1:]
		} else {
			unresolved = append(unresolved, token)
			// advance past this token to avoid infinite loop
			result = result[:start] + result[end+1:]
		}
	}
	return result, unresolved
}
