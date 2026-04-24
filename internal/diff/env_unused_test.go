package diff

import (
	"strings"
	"testing"
)

func TestCheckUnused_NoIssues(t *testing.T) {
	ref := map[string]string{"A": "1", "B": "2"}
	targets := []map[string]string{
		{"A": "1", "B": "2"},
	}
	issues := CheckUnused(ref, targets)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestCheckUnused_MissingInAll(t *testing.T) {
	ref := map[string]string{"A": "1", "B": "2", "C": "3"}
	targets := []map[string]string{
		{"A": "1"},
		{"A": "1", "B": "2"},
	}
	issues := CheckUnused(ref, targets)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d: %v", len(issues), issues)
	}
	if issues[0].Key != "C" {
		t.Errorf("expected key C, got %s", issues[0].Key)
	}
}

func TestCheckUnused_PresentInAtLeastOne(t *testing.T) {
	ref := map[string]string{"X": "val"}
	targets := []map[string]string{
		{},
		{"X": "val"},
	}
	issues := CheckUnused(ref, targets)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestCheckUnused_NoTargets(t *testing.T) {
	ref := map[string]string{"A": "1"}
	issues := CheckUnused(ref, nil)
	if len(issues) != 0 {
		t.Fatalf("expected no issues with empty targets, got %v", issues)
	}
}

func TestCheckUnused_SortedOutput(t *testing.T) {
	ref := map[string]string{"Z": "1", "A": "2", "M": "3"}
	targets := []map[string]string{{}}
	issues := CheckUnused(ref, targets)
	if len(issues) != 3 {
		t.Fatalf("expected 3 issues, got %d", len(issues))
	}
	if issues[0].Key != "A" || issues[1].Key != "M" || issues[2].Key != "Z" {
		t.Errorf("unexpected order: %v", issues)
	}
}

func TestUnusedIssue_String(t *testing.T) {
	iss := UnusedIssue{Key: "MY_KEY"}
	if iss.String() != "unused key: MY_KEY" {
		t.Errorf("unexpected string: %s", iss.String())
	}
}

func TestFormatUnusedIssues_Empty(t *testing.T) {
	out := FormatUnusedIssues(nil)
	if !strings.Contains(out, "no unused") {
		t.Errorf("expected no-unused message, got: %s", out)
	}
}

func TestFormatUnusedIssues_WithIssues(t *testing.T) {
	issues := []UnusedIssue{{Key: "ALPHA"}, {Key: "BETA"}}
	out := FormatUnusedIssues(issues)
	if !strings.Contains(out, "2 unused") {
		t.Errorf("expected count in output, got: %s", out)
	}
	if !strings.Contains(out, "ALPHA") || !strings.Contains(out, "BETA") {
		t.Errorf("expected keys in output, got: %s", out)
	}
}
