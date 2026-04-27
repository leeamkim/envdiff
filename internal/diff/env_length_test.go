package diff

import (
	"strings"
	"testing"
)

func TestCheckLengths_NoIssues(t *testing.T) {
	env := map[string]string{
		"API_KEY": "abcdefgh",
		"DB_PASS": "strongpass",
	}
	rules := []LengthRule{
		{Pattern: "API_", Min: 4, Max: 32},
		{Pattern: "DB_", Min: 4, Max: 64},
	}
	issues := CheckLengths(env, rules)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestCheckLengths_TooShort(t *testing.T) {
	env := map[string]string{
		"API_KEY": "ab",
	}
	rules := []LengthRule{
		{Pattern: "API_", Min: 8, Max: 0},
	}
	issues := CheckLengths(env, rules)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Violation != "too short" {
		t.Errorf("expected 'too short', got %q", issues[0].Violation)
	}
	if issues[0].Actual != 2 {
		t.Errorf("expected actual=2, got %d", issues[0].Actual)
	}
}

func TestCheckLengths_TooLong(t *testing.T) {
	env := map[string]string{
		"TOKEN_X": strings.Repeat("x", 200),
	}
	rules := []LengthRule{
		{Pattern: "TOKEN_", Min: 0, Max: 100},
	}
	issues := CheckLengths(env, rules)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Violation != "too long" {
		t.Errorf("expected 'too long', got %q", issues[0].Violation)
	}
}

func TestCheckLengths_WildcardPattern(t *testing.T) {
	env := map[string]string{
		"ANYTHING": "x",
	}
	rules := []LengthRule{
		{Pattern: "*", Min: 4, Max: 0},
	}
	issues := CheckLengths(env, rules)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
}

func TestCheckLengths_SortedOutput(t *testing.T) {
	env := map[string]string{
		"Z_KEY": "a",
		"A_KEY": "b",
		"M_KEY": "c",
	}
	rules := []LengthRule{
		{Pattern: "*", Min: 5, Max: 0},
	}
	issues := CheckLengths(env, rules)
	if len(issues) != 3 {
		t.Fatalf("expected 3 issues, got %d", len(issues))
	}
	if issues[0].Key > issues[1].Key || issues[1].Key > issues[2].Key {
		t.Errorf("issues not sorted: %v %v %v", issues[0].Key, issues[1].Key, issues[2].Key)
	}
}

func TestLengthIssue_String(t *testing.T) {
	issue := LengthIssue{Key: "API_KEY", Actual: 2, Min: 8, Max: 32, Violation: "too short"}
	s := issue.String()
	if !strings.Contains(s, "API_KEY") {
		t.Errorf("expected key in string, got %q", s)
	}
	if !strings.Contains(s, "too short") {
		t.Errorf("expected violation in string, got %q", s)
	}
}

func TestFormatLengthIssues_Empty(t *testing.T) {
	out := FormatLengthIssues(nil)
	if out != "no length issues found" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatLengthIssues_NonEmpty(t *testing.T) {
	issues := []LengthIssue{
		{Key: "SECRET", Actual: 1, Min: 10, Max: 0, Violation: "too short"},
	}
	out := FormatLengthIssues(issues)
	if !strings.Contains(out, "SECRET") {
		t.Errorf("expected key in output, got %q", out)
	}
}
