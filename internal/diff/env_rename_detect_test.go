package diff

import (
	"strings"
	"testing"
)

func TestDetectRenames_NoRenames(t *testing.T) {
	a := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := map[string]string{"FOO": "bar", "BAZ": "qux"}
	result := DetectRenames(a, b)
	if len(result) != 0 {
		t.Fatalf("expected no candidates, got %d", len(result))
	}
}

func TestDetectRenames_SingleCandidate(t *testing.T) {
	a := map[string]string{"OLD_KEY": "secret123"}
	b := map[string]string{"NEW_KEY": "secret123"}
	result := DetectRenames(a, b)
	if len(result) != 1 {
		t.Fatalf("expected 1 candidate, got %d", len(result))
	}
	if result[0].OldKey != "OLD_KEY" {
		t.Errorf("expected OldKey=OLD_KEY, got %s", result[0].OldKey)
	}
	if result[0].NewKey != "NEW_KEY" {
		t.Errorf("expected NewKey=NEW_KEY, got %s", result[0].NewKey)
	}
	if result[0].Value != "secret123" {
		t.Errorf("expected Value=secret123, got %s", result[0].Value)
	}
}

func TestDetectRenames_SkipsEmptyValues(t *testing.T) {
	a := map[string]string{"OLD_KEY": ""}
	b := map[string]string{"NEW_KEY": ""}
	result := DetectRenames(a, b)
	if len(result) != 0 {
		t.Fatalf("expected no candidates for empty values, got %d", len(result))
	}
}

func TestDetectRenames_AmbiguousValueSkipped(t *testing.T) {
	// Two keys in B share the same value — ambiguous, should not be reported
	a := map[string]string{"OLD_KEY": "shared"}
	b := map[string]string{"NEW_KEY1": "shared", "NEW_KEY2": "shared"}
	result := DetectRenames(a, b)
	if len(result) != 0 {
		t.Fatalf("expected no candidates for ambiguous value, got %d", len(result))
	}
}

func TestDetectRenames_SortedOutput(t *testing.T) {
	a := map[string]string{"Z_OLD": "val1", "A_OLD": "val2"}
	b := map[string]string{"Z_NEW": "val1", "A_NEW": "val2"}
	result := DetectRenames(a, b)
	if len(result) != 2 {
		t.Fatalf("expected 2 candidates, got %d", len(result))
	}
	if result[0].OldKey > result[1].OldKey {
		t.Errorf("expected sorted output, got %s before %s", result[0].OldKey, result[1].OldKey)
	}
}

func TestRenameCandidate_String(t *testing.T) {
	c := RenameCandidate{OldKey: "FOO", NewKey: "BAR", Value: "baz"}
	s := c.String()
	if !strings.Contains(s, "FOO") || !strings.Contains(s, "BAR") || !strings.Contains(s, "baz") {
		t.Errorf("unexpected String() output: %s", s)
	}
}

func TestFormatRenameReport_Empty(t *testing.T) {
	out := FormatRenameReport(nil)
	if !strings.Contains(out, "no rename") {
		t.Errorf("expected 'no rename' message, got: %s", out)
	}
}

func TestFormatRenameReport_WithCandidates(t *testing.T) {
	candidates := []RenameCandidate{
		{OldKey: "OLD", NewKey: "NEW", Value: "val"},
	}
	out := FormatRenameReport(candidates)
	if !strings.Contains(out, "OLD") || !strings.Contains(out, "NEW") {
		t.Errorf("expected candidate in report, got: %s", out)
	}
}
