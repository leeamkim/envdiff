package diff

import (
	"strings"
	"testing"
)

func TestCheckCrossRefs_NoIssues(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_URL":  "${DB_HOST}:5432",
	}
	issues := CheckCrossRefs(env, nil, "prod")
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestCheckCrossRefs_UnresolvedInlineRef(t *testing.T) {
	env := map[string]string{
		"DB_URL": "${DB_HOST}:5432",
	}
	issues := CheckCrossRefs(env, nil, "prod")
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].RefKey != "DB_HOST" {
		t.Errorf("expected RefKey=DB_HOST, got %s", issues[0].RefKey)
	}
	if !strings.Contains(issues[0].Reason, "DB_HOST") {
		t.Errorf("expected reason to mention DB_HOST, got: %s", issues[0].Reason)
	}
}

func TestCheckCrossRefs_ExplicitRuleViolated(t *testing.T) {
	env := map[string]string{
		"USE_TLS": "true",
	}
	rules := []CrossRefRule{
		{SourceKey: "USE_TLS", TargetKey: "TLS_CERT"},
	}
	issues := CheckCrossRefs(env, rules, "staging")
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "USE_TLS" || issues[0].RefKey != "TLS_CERT" {
		t.Errorf("unexpected issue: %+v", issues[0])
	}
}

func TestCheckCrossRefs_ExplicitRuleSatisfied(t *testing.T) {
	env := map[string]string{
		"USE_TLS":  "true",
		"TLS_CERT": "/etc/certs/cert.pem",
	}
	rules := []CrossRefRule{
		{SourceKey: "USE_TLS", TargetKey: "TLS_CERT"},
	}
	issues := CheckCrossRefs(env, rules, "prod")
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestCheckCrossRefs_SortedOutput(t *testing.T) {
	env := map[string]string{
		"Z_URL": "${MISSING_Z}",
		"A_URL": "${MISSING_A}",
	}
	issues := CheckCrossRefs(env, nil, "dev")
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}
	if issues[0].Key > issues[1].Key {
		t.Errorf("expected sorted output, got %s before %s", issues[0].Key, issues[1].Key)
	}
}

func TestFormatCrossRefIssues_Empty(t *testing.T) {
	out := FormatCrossRefIssues(nil)
	if out != "no cross-reference issues found" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatCrossRefIssues_WithIssues(t *testing.T) {
	issues := []CrossRefIssue{
		{Key: "DB_URL", RefKey: "DB_HOST", EnvName: "prod", Reason: "missing"},
	}
	out := FormatCrossRefIssues(issues)
	if !strings.Contains(out, "DB_URL") || !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected issue in output, got: %s", out)
	}
}

func TestCrossRefIssue_String(t *testing.T) {
	i := CrossRefIssue{Key: "A", RefKey: "B", EnvName: "prod", Reason: "missing"}
	s := i.String()
	if !strings.Contains(s, "A") || !strings.Contains(s, "B") || !strings.Contains(s, "prod") {
		t.Errorf("unexpected String output: %s", s)
	}
}
