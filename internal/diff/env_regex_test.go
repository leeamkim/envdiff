package diff

import (
	"regexp"
	"strings"
	"testing"
)

func TestCheckRegex_NoIssues(t *testing.T) {
	env := map[string]string{
		"PORT": "8080",
		"HOST": "localhost",
	}
	rules := []RegexRule{
		{Key: "PORT", Pattern: regexp.MustCompile(`^\d+$`)},
	}
	issues := CheckRegex(env, rules)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestCheckRegex_PatternMismatch(t *testing.T) {
	env := map[string]string{
		"PORT": "not-a-number",
	}
	rules := []RegexRule{
		{Key: "PORT", Pattern: regexp.MustCompile(`^\d+$`)},
	}
	issues := CheckRegex(env, rules)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "PORT" {
		t.Errorf("expected key PORT, got %q", issues[0].Key)
	}
}

func TestCheckRegex_WildcardRule(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "abc",
		"APP_NAME": "envdiff",
	}
	rules := []RegexRule{
		{Key: "DB_*", Pattern: regexp.MustCompile(`^[a-z0-9]+$`)},
	}
	issues := CheckRegex(env, rules)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "DB_PORT" {
		t.Errorf("expected DB_PORT, got %q", issues[0].Key)
	}
}

func TestCheckRegex_SortedOutput(t *testing.T) {
	env := map[string]string{
		"Z_KEY": "bad",
		"A_KEY": "bad",
	}
	rules := []RegexRule{
		{Key: "Z_KEY", Pattern: regexp.MustCompile(`^\d+$`)},
		{Key: "A_KEY", Pattern: regexp.MustCompile(`^\d+$`)},
	}
	issues := CheckRegex(env, rules)
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}
	if issues[0].Key != "A_KEY" || issues[1].Key != "Z_KEY" {
		t.Errorf("issues not sorted: %v", issues)
	}
}

func TestRegexIssue_String(t *testing.T) {
	issue := RegexIssue{Key: "PORT", Value: "abc", Pattern: `^\d+$`}
	s := issue.String()
	if !strings.Contains(s, "PORT") || !strings.Contains(s, "abc") {
		t.Errorf("unexpected string: %q", s)
	}
}

func TestFormatRegexIssues_Empty(t *testing.T) {
	out := FormatRegexIssues(nil)
	if !strings.Contains(out, "no regex") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatRegexIssues_WithIssues(t *testing.T) {
	issues := []RegexIssue{
		{Key: "PORT", Value: "abc", Pattern: `^\d+$`},
	}
	out := FormatRegexIssues(issues)
	if !strings.Contains(out, "1 regex issue") {
		t.Errorf("unexpected output: %q", out)
	}
}
