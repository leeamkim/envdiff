package diff

import (
	"strings"
	"testing"
)

func TestCheckPrefixes_NoIssues(t *testing.T) {
	env := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
	}
	issues := CheckPrefixes(env, []string{"APP_"})
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestCheckPrefixes_ViolatingKey(t *testing.T) {
	env := map[string]string{
		"APP_HOST": "localhost",
		"DB_HOST":  "db",
	}
	issues := CheckPrefixes(env, []string{"APP_"})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", issues[0].Key)
	}
}

func TestCheckPrefixes_MultipleAllowed(t *testing.T) {
	env := map[string]string{
		"APP_HOST": "localhost",
		"SVC_NAME": "api",
		"SECRET":   "x",
	}
	issues := CheckPrefixes(env, []string{"APP_", "SVC_"})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "SECRET" {
		t.Errorf("expected SECRET, got %s", issues[0].Key)
	}
}

func TestCheckPrefixes_EmptyPrefixes(t *testing.T) {
	env := map[string]string{
		"ANYTHING": "val",
	}
	issues := CheckPrefixes(env, []string{})
	if len(issues) != 0 {
		t.Fatalf("expected no issues with empty prefixes, got %d", len(issues))
	}
}

func TestCheckPrefixes_SortedOutput(t *testing.T) {
	env := map[string]string{
		"Z_KEY": "z",
		"A_KEY": "a",
		"M_KEY": "m",
	}
	issues := CheckPrefixes(env, []string{"APP_"})
	if len(issues) != 3 {
		t.Fatalf("expected 3 issues, got %d", len(issues))
	}
	if issues[0].Key != "A_KEY" || issues[1].Key != "M_KEY" || issues[2].Key != "Z_KEY" {
		t.Errorf("issues not sorted: %v", issues)
	}
}

func TestPrefixIssue_String(t *testing.T) {
	issue := PrefixIssue{Key: "DB_HOST", Expected: "APP_"}
	s := issue.String()
	if !strings.Contains(s, "DB_HOST") || !strings.Contains(s, "APP_") {
		t.Errorf("unexpected string: %s", s)
	}
}

func TestFormatPrefixIssues_Empty(t *testing.T) {
	out := FormatPrefixIssues(nil)
	if out != "no prefix issues found" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatPrefixIssues_NonEmpty(t *testing.T) {
	issues := []PrefixIssue{{Key: "DB_HOST", Expected: "APP_"}}
	out := FormatPrefixIssues(issues)
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got: %s", out)
	}
}
