package diff

import (
	"fmt"
	"sort"
	"strings"
)

// FormatStyle controls the output format for env formatting.
type FormatStyle string

const (
	FormatStyleExport FormatStyle = "export"
	FormatStylePlain  FormatStyle = "plain"
	FormatStyleQuoted FormatStyle = "quoted"
)

// FormatOptions configures how an env map is formatted.
type FormatOptions struct {
	Style    FormatStyle
	SortKeys bool
	Prefix   string
}

// FormatEnv renders an env map as a string using the given options.
func FormatEnv(env map[string]string, opts FormatOptions) string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	if opts.SortKeys {
		sort.Strings(keys)
	}

	var sb strings.Builder
	for _, k := range keys {
		v := env[k]
		key := opts.Prefix + k
		switch opts.Style {
		case FormatStyleExport:
			fmt.Fprintf(&sb, "export %s=%s\n", key, shellQuote(v))
		case FormatStyleQuoted:
			fmt.Fprintf(&sb, "%s=%q\n", key, v)
		default:
			fmt.Fprintf(&sb, "%s=%s\n", key, v)
		}
	}
	return sb.String()
}

// shellQuote wraps a value in single quotes if it contains special characters.
func shellQuote(v string) string {
	if v == "" {
		return "\"\""
	}
	specials := " \t\n$\"'\\|&;<>(){}!`"
	for _, c := range specials {
		if strings.ContainsRune(v, c) {
			return "'" + strings.ReplaceAll(v, "'", "'\\'''") + "'"
		}
	}
	return v
}
