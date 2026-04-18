package diff

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// ExportFormat represents the output format for exporting diff results.
type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatCSV  ExportFormat = "csv"
	FormatText ExportFormat = "text"
)

// ExportResult writes the diff result to w in the given format.
func ExportResult(w io.Writer, result Result, format ExportFormat) error {
	switch format {
	case FormatJSON:
		return exportJSON(w, result)
	case FormatCSV:
		return exportCSV(w, result)
	case FormatText:
		return exportText(w, result)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

type jsonEntry struct {
	Key    string `json:"key"`
	Status string `json:"status"`
	ValueA string `json:"value_a,omitempty"`
	ValueB string `json:"value_b,omitempty"`
}

func exportJSON(w io.Writer, result Result) error {
	var entries []jsonEntry
	for _, k := range sortedKeys(result.MissingInB) {
		entries = append(entries, jsonEntry{Key: k, Status: "missing_in_b", ValueA: result.MissingInB[k]})
	}
	for _, k := range sortedKeys(result.MissingInA) {
		entries = append(entries, jsonEntry{Key: k, Status: "missing_in_a", ValueB: result.MissingInA[k]})
	}
	for _, k := range sortedKeys(result.Mismatched) {
		p := result.Mismatched[k]
		entries = append(entries, jsonEntry{Key: k, Status: "mismatched", ValueA: p[0], ValueB: p[1]})
	}
	if entries == nil {
		entries = []jsonEntry{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}

func exportCSV(w io.Writer, result Result) error {
	fmt.Fprintln(w, "key,status,value_a,value_b")
	var keys []string
	for k := range result.MissingInB {
		keys = append(keys, "mb:"+k)
	}
	for k := range result.MissingInA {
		keys = append(keys, "ma:"+k)
	}
	for k := range result.Mismatched {
		keys = append(keys, "mm:"+k)
	}
	sort.Strings(keys)
	for _, raw := range keys {
		parts := strings.SplitN(raw, ":", 2)
		tag, k := parts[0], parts[1]
		switch tag {
		case "mb":
			fmt.Fprintf(w, "%s,missing_in_b,%s,\n", k, result.MissingInB[k])
		case "ma":
			fmt.Fprintf(w, "%s,missing_in_a,,%s\n", k, result.MissingInA[k])
		case "mm":
			p := result.Mismatched[k]
			fmt.Fprintf(w, "%s,mismatched,%s,%s\n", k, p[0], p[1])
		}
	}
	return nil
}

func exportText(w io.Writer, result Result) error {
	PrintReport(w, result)
	return nil
}
