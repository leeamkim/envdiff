package diff

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintReport_NoDiff(t *testing.T) {
	var buf bytes.Buffer
	PrintReport(&buf, Result{}, "a.env", "b.env")
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no-diff message, got: %s", buf.String())
	}
}

func TestPrintReport_MissingInB(t *testing.T) {
	var buf bytes.Buffer
	r := Result{
		MissingInB: map[string]string{"SECRET": "abc"},
	}
	PrintReport(&buf, r, "a.env", "b.env")
	out := buf.String()
	if !strings.Contains(out, "SECRET") {
		t.Errorf("expected SECRET in output, got: %s", out)
	}
	if !strings.Contains(out, "a.env") {
		t.Errorf("expected label a.env in output")
	}
}

func TestPrintReport_Mismatched(t *testing.T) {
	var buf bytes.Buffer
	r := Result{
		Mismatched: map[string][2]string{
			"DB_URL": {"postgres://local", "postgres://prod"},
		},
	}
	PrintReport(&buf, r, "local.env", "prod.env")
	out := buf.String()
	if !strings.Contains(out, "DB_URL") {
		t.Errorf("expected DB_URL in output")
	}
	if !strings.Contains(out, "local.env") || !strings.Contains(out, "prod.env") {
		t.Errorf("expected both labels in output")
	}
}

func TestPrintReport_AllDiffTypes(t *testing.T) {
	var buf bytes.Buffer
	r := Result{
		MissingInB: map[string]string{"ONLY_A": "1"},
		MissingInA: map[string]string{"ONLY_B": "2"},
		Mismatched: map[string][2]string{"SHARED": {"x", "y"}},
	}
	PrintReport(&buf, r, "a.env", "b.env")
	out := buf.String()
	for _, want := range []string{"ONLY_A", "ONLY_B", "SHARED"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %s in output", want)
		}
	}
}
