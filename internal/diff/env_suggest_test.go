package diff

import (
	"strings"
	"testing"
)

func TestSuggestRenames_NoIssues(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	ref := map[string]string{"DB_HOST": "prod", "DB_PORT": "5432"}
	issues := SuggestRenames(env, ref)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestSuggestRenames_CaseMismatch(t *testing.T) {
	env := map[string]string{"db_host": "localhost"}
	ref := map[string]string{"DB_HOST": "prod"}
	issues := SuggestRenames(env, ref)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "db_host" {
		t.Errorf("unexpected key: %s", issues[0].Key)
	}
	if issues[0].Suggestion != "DB_HOST" {
		t.Errorf("unexpected suggestion: %s", issues[0].Suggestion)
	}
	if !strings.Contains(issues[0].Reason, "case") {
		t.Errorf("expected reason to mention case, got: %s", issues[0].Reason)
	}
}

func TestSuggestRenames_DashUnderscoreMismatch(t *testing.T) {
	env := map[string]string{"DB-HOST": "localhost"}
	ref := map[string]string{"DB_HOST": "prod"}
	issues := SuggestRenames(env, ref)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Suggestion != "DB_HOST" {
		t.Errorf("unexpected suggestion: %s", issues[0].Suggestion)
	}
}

func TestSuggestRenames_NoMatchForUnknownKey(t *testing.T) {
	env := map[string]string{"TOTALLY_DIFFERENT": "val"}
	ref := map[string]string{"DB_HOST": "prod"}
	issues := SuggestRenames(env, ref)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestSuggestRenames_SortedOutput(t *testing.T) {
	env := map[string]string{"z_key": "1", "a_key": "2"}
	ref := map[string]string{"Z_KEY": "1", "A_KEY": "2"}
	issues := SuggestRenames(env, ref)
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}
	if issues[0].Key > issues[1].Key {
		t.Errorf("expected sorted output, got %s before %s", issues[0].Key, issues[1].Key)
	}
}

func TestFormatSuggestIssues_Empty(t *testing.T) {
	out := FormatSuggestIssues(nil)
	if out != "no rename suggestions" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatSuggestIssues_WithIssues(t *testing.T) {
	issues := []SuggestIssue{
		{Key: "db_host", Suggestion: "DB_HOST", Reason: "case/underscore mismatch"},
	}
	out := FormatSuggestIssues(issues)
	if !strings.Contains(out, "db_host") {
		t.Errorf("expected key in output, got: %s", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected suggestion in output, got: %s", out)
	}
}

func TestSuggestIssue_String(t *testing.T) {
	iss := SuggestIssue{Key: "db_host", Suggestion: "DB_HOST", Reason: "case/underscore mismatch"}
	s := iss.String()
	if !strings.Contains(s, "db_host") || !strings.Contains(s, "DB_HOST") {
		t.Errorf("unexpected string: %s", s)
	}
}
