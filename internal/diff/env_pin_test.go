package diff

import (
	"strings"
	"testing"
)

func TestCheckPins_NoIssues(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "8080"}
	pins := map[string]string{"HOST": "localhost", "PORT": "8080"}
	issues := CheckPins(env, pins)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestCheckPins_MismatchedValue(t *testing.T) {
	env := map[string]string{"HOST": "prod.example.com"}
	pins := map[string]string{"HOST": "localhost"}
	issues := CheckPins(env, pins)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "HOST" || issues[0].Pinned != "localhost" || issues[0].Actual != "prod.example.com" {
		t.Errorf("unexpected entry: %+v", issues[0])
	}
}

func TestCheckPins_MissingKey(t *testing.T) {
	env := map[string]string{}
	pins := map[string]string{"SECRET": "abc123"}
	issues := CheckPins(env, pins)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if !issues[0].Missing {
		t.Error("expected Missing=true")
	}
}

func TestCheckPins_SortedOutput(t *testing.T) {
	env := map[string]string{"Z": "wrong", "A": "wrong"}
	pins := map[string]string{"Z": "right", "A": "right"}
	issues := CheckPins(env, pins)
	if issues[0].Key != "A" || issues[1].Key != "Z" {
		t.Error("expected sorted by key")
	}
}

func TestFormatPinIssues_Empty(t *testing.T) {
	out := FormatPinIssues(nil)
	if out != "all pinned values match" {
		t.Errorf("unexpected: %s", out)
	}
}

func TestFormatPinIssues_WithIssues(t *testing.T) {
	issues := []PinEntry{
		{Key: "HOST", Pinned: "localhost", Actual: "remote"},
		{Key: "PORT", Pinned: "8080", Missing: true},
	}
	out := FormatPinIssues(issues)
	if !strings.Contains(out, "2 pin violation") {
		t.Errorf("expected count in output, got: %s", out)
	}
	if !strings.Contains(out, "HOST") || !strings.Contains(out, "PORT") {
		t.Errorf("expected keys in output, got: %s", out)
	}
}

func TestPinEntry_String_Missing(t *testing.T) {
	e := PinEntry{Key: "X", Pinned: "val", Missing: true}
	if !strings.Contains(e.String(), "missing in env") {
		t.Errorf("unexpected: %s", e.String())
	}
}
