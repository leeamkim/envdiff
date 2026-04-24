package diff

import (
	"strings"
	"testing"
)

func TestCheckTrim_NoIssues(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "8080",
	}
	issues := CheckTrim(env)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestCheckTrim_LeadingSpace(t *testing.T) {
	env := map[string]string{
		"HOST": "  localhost",
	}
	issues := CheckTrim(env)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "HOST" {
		t.Errorf("expected key HOST, got %s", issues[0].Key)
	}
	if issues[0].Trimmed != "localhost" {
		t.Errorf("expected trimmed 'localhost', got %q", issues[0].Trimmed)
	}
}

func TestCheckTrim_TrailingSpace(t *testing.T) {
	env := map[string]string{
		"DB_URL": "postgres://localhost ",
	}
	issues := CheckTrim(env)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Trimmed != "postgres://localhost" {
		t.Errorf("unexpected trimmed value: %q", issues[0].Trimmed)
	}
}

func TestCheckTrim_SortedOutput(t *testing.T) {
	env := map[string]string{
		"Z_KEY": " z ",
		"A_KEY": " a ",
		"M_KEY": " m ",
	}
	issues := CheckTrim(env)
	if len(issues) != 3 {
		t.Fatalf("expected 3 issues, got %d", len(issues))
	}
	if issues[0].Key != "A_KEY" || issues[1].Key != "M_KEY" || issues[2].Key != "Z_KEY" {
		t.Errorf("issues not sorted: %v", issues)
	}
}

func TestTrimIssue_String(t *testing.T) {
	issue := TrimIssue{Key: "FOO", Original: " bar ", Trimmed: "bar"}
	s := issue.String()
	if !strings.Contains(s, "FOO") {
		t.Errorf("expected key in string, got: %s", s)
	}
	if !strings.Contains(s, "bar") {
		t.Errorf("expected value in string, got: %s", s)
	}
}

func TestFormatTrimIssues_Empty(t *testing.T) {
	out := FormatTrimIssues(nil)
	if out != "no trim issues found" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatTrimIssues_WithIssues(t *testing.T) {
	issues := []TrimIssue{
		{Key: "API_KEY", Original: " secret ", Trimmed: "secret"},
	}
	out := FormatTrimIssues(issues)
	if !strings.Contains(out, "1 trim issue") {
		t.Errorf("expected count in output, got: %s", out)
	}
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected key in output, got: %s", out)
	}
}
