package diff

import (
	"fmt"
	"regexp"
	"strings"
)

// TypeHint describes the inferred type of an env value.
type TypeHint string

const (
	TypeBoolean TypeHint = "boolean"
	TypeInteger TypeHint = "integer"
	TypeFloat   TypeHint = "float"
	TypeURL     TypeHint = "url"
	TypeEmpty   TypeHint = "empty"
	TypeString  TypeHint = "string"
)

// TypeEntry holds a key, its value, and the inferred type.
type TypeEntry struct {
	Key   string
	Value string
	Type  TypeHint
}

// TypeIssue represents a type inconsistency across environments.
type TypeIssue struct {
	Key   string
	Types map[string]TypeHint // env name -> inferred type
}

func (t TypeIssue) String() string {
	parts := make([]string, 0, len(t.Types))
	for env, hint := range t.Types {
		parts = append(parts, fmt.Sprintf("%s=%s", env, hint))
	}
	return fmt.Sprintf("key %q has inconsistent types: %s", t.Key, strings.Join(parts, ", "))
}

var (
	reURL     = regexp.MustCompile(`^https?://`)
	reInteger = regexp.MustCompile(`^-?[0-9]+$`)
	reFloat   = regexp.MustCompile(`^-?[0-9]+\.[0-9]+$`)
)

// InferType returns the TypeHint for a given string value.
func InferType(value string) TypeHint {
	switch {
	case value == "":
		return TypeEmpty
	case strings.EqualFold(value, "true") || strings.EqualFold(value, "false"):
		return TypeBoolean
	case reInteger.MatchString(value):
		return TypeInteger
	case reFloat.MatchString(value):
		return TypeFloat
	case reURL.MatchString(value):
		return TypeURL
	default:
		return TypeString
	}
}

// InferTypes annotates each key in an env map with its inferred type.
func InferTypes(env map[string]string) []TypeEntry {
	entries := make([]TypeEntry, 0, len(env))
	for k, v := range env {
		entries = append(entries, TypeEntry{Key: k, Value: v, Type: InferType(v)})
	}
	return entries
}

// CheckTypeConsistency compares inferred types for shared keys across multiple envs.
// It returns issues where the same key has different inferred types.
func CheckTypeConsistency(envs map[string]map[string]string) []TypeIssue {
	keyTypes := map[string]map[string]TypeHint{}
	for envName, vars := range envs {
		for k, v := range vars {
			if keyTypes[k] == nil {
				keyTypes[k] = map[string]TypeHint{}
			}
			keyTypes[k][envName] = InferType(v)
		}
	}

	var issues []TypeIssue
	for key, typeMap := range keyTypes {
		first := TypeHint("")
		consistent := true
		for _, t := range typeMap {
			if first == "" {
				first = t
			} else if t != first {
				consistent = false
				break
			}
		}
		if !consistent {
			issues = append(issues, TypeIssue{Key: key, Types: typeMap})
		}
	}
	return issues
}
