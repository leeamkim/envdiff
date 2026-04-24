package diff

import (
	"strings"
	"testing"
)

func TestCheckClamps_NoIssues(t *testing.T) {
	env := map[string]string{
		"TOKEN": "abcdef",
		"NAME":  "alice",
	}
	rules := []ClampRule{
		{Key: "TOKEN", Min: 4, Max: 10},
		{Key: "NAME", Min: 1, Max: 10},
	}
	issues := CheckClamps(env, rules)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestCheckClamps_TooShort(t *testing.T) {
	env := map[string]string{"TOKEN": "ab"}
	rules := []ClampRule{{Key: "TOKEN", Min: 5, Max: 20}}
	issues := CheckClamps(env, rules)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Length != 2 {
		t.Errorf("expected length 2, got %d", issues[0].Length)
	}
}

func TestCheckClamps_TooLong(t *testing.T) {
	env := map[string]string{"SECRET": "verylongsecretvalue"}
	rules := []ClampRule{{Key: "SECRET", Min: 1, Max: 8}}
	issues := CheckClamps(env, rules)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "SECRET" {
		t.Errorf("unexpected key %q", issues[0].Key)
	}
}

func TestCheckClamps_SkipsMissingKeys(t *testing.T) {
	env := map[string]string{"OTHER": "value"}
	rules := []ClampRule{{Key: "TOKEN", Min: 1, Max: 10}}
	issues := CheckClamps(env, rules)
	if len(issues) != 0 {
		t.Fatalf("expected no issues for missing key, got %d", len(issues))
	}
}

func TestClampIssue_String(t *testing.T) {
	i := ClampIssue{Key: "TOKEN", Value: "ab", Min: 5, Max: 20, Length: 2}
	s := i.String()
	if !strings.Contains(s, "TOKEN") || !strings.Contains(s, "2") {
		t.Errorf("unexpected string output: %q", s)
	}
}

func TestFormatClampIssues_Empty(t *testing.T) {
	out := FormatClampIssues(nil)
	if !strings.Contains(out, "no clamp issues") {
		t.Errorf("expected no-issues message, got %q", out)
	}
}

func TestFormatClampIssues_WithIssues(t *testing.T) {
	issues := []ClampIssue{
		{Key: "A", Value: "x", Min: 5, Max: 10, Length: 1},
		{Key: "B", Value: "toolongvalue", Min: 1, Max: 5, Length: 12},
	}
	out := FormatClampIssues(issues)
	if !strings.Contains(out, "2 clamp issue") {
		t.Errorf("expected count in output, got %q", out)
	}
	if !strings.Contains(out, "\"A\"") || !strings.Contains(out, "\"B\"") {
		t.Errorf("expected both keys in output, got %q", out)
	}
}
