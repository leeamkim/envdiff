package cli

import (
	"bytes"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct {
		input    string
		want     diff.ExportFormat
	}{
		{"json", diff.FormatJSON},
		{"JSON", diff.FormatJSON},
		{"csv", diff.FormatCSV},
		{"text", diff.FormatText},
		{"", diff.FormatText},
	}
	for _, tc := range cases {
		got, err := ParseFormat(tc.input)
		if err != nil {
			t.Errorf("ParseFormat(%q) unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ParseFormat(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := ParseFormat("xml")
	if err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestWriteOutput_JSON(t *testing.T) {
	result := diff.Result{
		MissingInA: map[string]string{"KEY": "val"},
		MissingInB: map[string]string{},
		Mismatched: map[string][2]string{},
	}
	var buf bytes.Buffer
	err := WriteOutput(&buf, result, diff.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

func TestWriteOutput_CSV(t *testing.T) {
	result := diff.Result{
		MissingInA: map[string]string{},
		MissingInB: map[string]string{"ONLY_A": "v"},
		Mismatched: map[string][2]string{},
	}
	var buf bytes.Buffer
	err := WriteOutput(&buf, result, diff.FormatCSV)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}
