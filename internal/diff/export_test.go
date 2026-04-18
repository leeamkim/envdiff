package diff

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func baseExportResult() Result {
	return Result{
		MissingInA: map[string]string{"ONLY_B": "val_b"},
		MissingInB: map[string]string{"ONLY_A": "val_a"},
		Mismatched: map[string][2]string{"SHARED": {"foo", "bar"}},
	}
}

func TestExportResult_JSON(t *testing.T) {
	var buf bytes.Buffer
	err := ExportResult(&buf, baseExportResult(), FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var entries []map[string]string
	if err := json.Unmarshal(buf.Bytes(), &entries); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
}

func TestExportResult_JSON_Empty(t *testing.T) {
	var buf bytes.Buffer
	err := ExportResult(&buf, Result{}, FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "[]") {
		t.Errorf("expected empty JSON array, got: %s", buf.String())
	}
}

func TestExportResult_CSV(t *testing.T) {
	var buf bytes.Buffer
	err := ExportResult(&buf, baseExportResult(), FormatCSV)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if lines[0] != "key,status,value_a,value_b" {
		t.Errorf("unexpected CSV header: %s", lines[0])
	}
	if len(lines) != 4 {
		t.Errorf("expected 4 lines (header + 3 entries), got %d", len(lines))
	}
}

func TestExportResult_Text(t *testing.T) {
	var buf bytes.Buffer
	err := ExportResult(&buf, baseExportResult(), FormatText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty text output")
	}
}

func TestExportResult_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	err := ExportResult(&buf, baseExportResult(), ExportFormat("xml"))
	if err == nil {
		t.Error("expected error for unknown format")
	}
}
