package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// ParseFormat converts a string to an ExportFormat, returning an error if unknown.
func ParseFormat(s string) (diff.ExportFormat, error) {
	switch strings.ToLower(s) {
	case "json":
		return diff.FormatJSON, nil
	case "csv":
		return diff.FormatCSV, nil
	case "text", "":
		return diff.FormatText, nil
	default:
		return "", fmt.Errorf("unknown format %q: must be one of json, csv, text", s)
	}
}

// WriteOutput exports the diff result to w using the specified format.
func WriteOutput(w io.Writer, result diff.Result, format diff.ExportFormat) error {
	return diff.ExportResult(w, result, format)
}
