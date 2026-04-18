package diff

import (
	"strings"
	"testing"
)

func TestLint_NoIssues(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "8080"}
	result := Lint("a.env", env, []LintRule{LintRuleNoEmptyValues, LintRuleKeyUppercase})
	if result.HasIssues() {
		t.Fatalf("expected no issues, got %v", result.Issues)
	}
}

func TestLint_EmptyValue(t *testing.T) {
	env := map[string]string{"API_KEY": ""}
	result := Lint("a.env", env, []LintRule{LintRuleNoEmptyValues})
	if !result.HasIssues() {
		t.Fatal("expected issues")
	}
	if result.Issues[0].Key != "API_KEY" {
		t.Errorf("unexpected key: %s", result.Issues[0].Key)
	}
}

func TestLint_PlaceholderValue(t *testing.T) {
	env := map[string]string{"SECRET": "CHANGEME"}
	result := Lint("b.env", env, []LintRule{LintRuleNoDuplicatePlaceholder})
	if !result.HasIssues() {
		t.Fatal("expected issues")
	}
	if !strings.Contains(result.Issues[0].Msg, "placeholder") {
		t.Errorf("unexpected msg: %s", result.Issues[0].Msg)
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	env := map[string]string{"db_host": "localhost"}
	result := Lint("c.env", env, []LintRule{LintRuleKeyUppercase})
	if !result.HasIssues() {
		t.Fatal("expected issues for lowercase key")
	}
}

func TestLintIssue_String(t *testing.T) {
	issue := LintIssue{File: "a.env", Key: "FOO", Msg: "value is empty"}
	s := issue.String()
	if !strings.Contains(s, "a.env") || !strings.Contains(s, "FOO") {
		t.Errorf("unexpected string: %s", s)
	}
}

func TestLint_MultipleRules(t *testing.T) {
	env := map[string]string{"foo": "", "BAR": "CHANGEME"}
	rules := []LintRule{LintRuleNoEmptyValues, LintRuleKeyUppercase, LintRuleNoDuplicatePlaceholder}
	result := Lint("x.env", env, rules)
	if len(result.Issues) < 2 {
		t.Errorf("expected at least 2 issues, got %d", len(result.Issues))
	}
}
