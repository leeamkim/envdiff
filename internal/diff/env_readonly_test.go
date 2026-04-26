package diff

import (
	"strings"
	"testing"
)

func TestCheckReadonly_NoIssues(t *testing.T) {
	ref := map[string]string{"DB_HOST": "localhost", "APP_ENV": "production"}
	targets := map[string]map[string]string{
		"staging": {"DB_HOST": "localhost", "APP_ENV": "staging"},
	}
	issues := CheckReadonly(ref, targets, []string{"DB_HOST"})
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestCheckReadonly_DetectsChange(t *testing.T) {
	ref := map[string]string{"DB_HOST": "localhost"}
	targets := map[string]map[string]string{
		"staging": {"DB_HOST": "db.staging.internal"},
	}
	issues := CheckReadonly(ref, targets, []string{"DB_HOST"})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "DB_HOST" {
		t.Errorf("expected key DB_HOST, got %s", issues[0].Key)
	}
	if issues[0].Env != "staging" {
		t.Errorf("expected env staging, got %s", issues[0].Env)
	}
}

func TestCheckReadonly_EmptyReadonlyKeys(t *testing.T) {
	ref := map[string]string{"DB_HOST": "localhost"}
	targets := map[string]map[string]string{
		"staging": {"DB_HOST": "changed"},
	}
	issues := CheckReadonly(ref, targets, []string{})
	if len(issues) != 0 {
		t.Fatalf("expected no issues with empty readonly list, got %d", len(issues))
	}
}

func TestCheckReadonly_MissingKeyInTarget(t *testing.T) {
	ref := map[string]string{"DB_HOST": "localhost"}
	targets := map[string]map[string]string{
		"staging": {"OTHER_KEY": "value"},
	}
	issues := CheckReadonly(ref, targets, []string{"DB_HOST"})
	if len(issues) != 0 {
		t.Fatalf("expected no issues when key missing in target, got %d", len(issues))
	}
}

func TestCheckReadonly_MultipleEnvsSorted(t *testing.T) {
	ref := map[string]string{"SECRET_KEY": "abc123"}
	targets := map[string]map[string]string{
		"prod":    {"SECRET_KEY": "changed"},
		"staging": {"SECRET_KEY": "also-changed"},
	}
	issues := CheckReadonly(ref, targets, []string{"SECRET_KEY"})
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}
	if issues[0].Env != "prod" {
		t.Errorf("expected first issue env=prod, got %s", issues[0].Env)
	}
	if issues[1].Env != "staging" {
		t.Errorf("expected second issue env=staging, got %s", issues[1].Env)
	}
}

func TestReadonlyIssue_String(t *testing.T) {
	issue := ReadonlyIssue{Key: "DB_HOST", Env: "prod", OldValue: "localhost", NewValue: "remote"}
	s := issue.String()
	if !strings.Contains(s, "DB_HOST") || !strings.Contains(s, "prod") {
		t.Errorf("unexpected string output: %s", s)
	}
}

func TestFormatReadonlyIssues_Empty(t *testing.T) {
	out := FormatReadonlyIssues(nil)
	if out != "no readonly violations found" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatReadonlyIssues_WithIssues(t *testing.T) {
	issues := []ReadonlyIssue{
		{Key: "DB_HOST", Env: "prod", OldValue: "localhost", NewValue: "remote"},
	}
	out := FormatReadonlyIssues(issues)
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got: %s", out)
	}
}
