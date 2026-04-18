package diff

import (
	"testing"
)

func TestValidate_NoIssues(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "8080",
	}
	issues := Validate(env, []ValidationRule{RuleNoEmptyValues})
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestValidate_EmptyValue(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"SECRET": "",
	}
	issues := Validate(env, []ValidationRule{RuleNoEmptyValues})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "SECRET" {
		t.Errorf("expected key SECRET, got %s", issues[0].Key)
	}
}

func TestValidate_WhitespaceValue(t *testing.T) {
	env := map[string]string{
		"KEY": "   ",
	}
	issues := Validate(env, []ValidationRule{RuleNoWhitespaceValues})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "KEY" {
		t.Errorf("expected key KEY, got %s", issues[0].Key)
	}
}

func TestValidate_MultipleRules(t *testing.T) {
	env := map[string]string{
		"A": "",
		"B": "  ",
		"C": "ok",
	}
	issues := Validate(env, []ValidationRule{RuleNoEmptyValues, RuleNoWhitespaceValues})
	// A triggers RuleNoEmptyValues, B triggers RuleNoWhitespaceValues
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d: %v", len(issues), issues)
	}
}

func TestValidationIssue_String(t *testing.T) {
	v := ValidationIssue{Key: "FOO", Message: "value is empty"}
	expected := "FOO: value is empty"
	if v.String() != expected {
		t.Errorf("expected %q, got %q", expected, v.String())
	}
}
