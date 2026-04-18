package diff

import (
	"testing"
)

func TestCompareToTemplate_NoIssues(t *testing.T) {
	tmpl := map[string]string{"A": "", "B": ""}
	env := map[string]string{"A": "1", "B": "2"}
	issues := CompareToTemplate(tmpl, env, false)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestCompareToTemplate_MissingKey(t *testing.T) {
	tmpl := map[string]string{"A": "", "B": "", "C": ""}
	env := map[string]string{"A": "1"}
	issues := CompareToTemplate(tmpl, env, false)
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}
	if issues[0].Key != "B" || issues[1].Key != "C" {
		t.Errorf("unexpected keys: %v", issues)
	}
}

func TestCompareToTemplate_StrictExtraKey(t *testing.T) {
	tmpl := map[string]string{"A": ""}
	env := map[string]string{"A": "1", "EXTRA": "x"}
	issues := CompareToTemplate(tmpl, env, true)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Reason != "extra key not in template" {
		t.Errorf("unexpected reason: %s", issues[0].Reason)
	}
}

func TestCompareToTemplate_StrictNoExtra(t *testing.T) {
	tmpl := map[string]string{"A": "", "B": ""}
	env := map[string]string{"A": "1", "B": "2"}
	issues := CompareToTemplate(tmpl, env, true)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestFormatTemplateIssues_Empty(t *testing.T) {
	out := FormatTemplateIssues(nil)
	if out != "no template issues found" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatTemplateIssues_WithIssues(t *testing.T) {
	issues := []TemplateIssue{{Key: "FOO", Reason: "missing from env"}}
	out := FormatTemplateIssues(issues)
	if out == "" {
		t.Error("expected non-empty output")
	}
}
