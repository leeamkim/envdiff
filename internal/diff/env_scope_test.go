package diff

import (
	"strings"
	"testing"
)

func TestCheckScopes_NoIssues(t *testing.T) {
	envs := map[string]map[string]string{
		"prod": {"APP_HOST": "example.com", "APP_PORT": "8080"},
	}
	rules := []ScopeRule{
		{Scope: "prod", Prefixes: []string{"APP_"}},
	}
	issues := CheckScopes(envs, rules)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestCheckScopes_ViolatingKey(t *testing.T) {
	envs := map[string]map[string]string{
		"prod": {"APP_HOST": "example.com", "DB_PASS": "secret"},
	}
	rules := []ScopeRule{
		{Scope: "prod", Prefixes: []string{"APP_"}},
	}
	issues := CheckScopes(envs, rules)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "DB_PASS" {
		t.Errorf("expected key DB_PASS, got %s", issues[0].Key)
	}
}

func TestCheckScopes_MultipleAllowedPrefixes(t *testing.T) {
	envs := map[string]map[string]string{
		"staging": {"APP_HOST": "stg.example.com", "SVC_URL": "http://svc", "UNKNOWN_KEY": "x"},
	}
	rules := []ScopeRule{
		{Scope: "staging", Prefixes: []string{"APP_", "SVC_"}},
	}
	issues := CheckScopes(envs, rules)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "UNKNOWN_KEY" {
		t.Errorf("expected UNKNOWN_KEY, got %s", issues[0].Key)
	}
}

func TestCheckScopes_EmptyRules(t *testing.T) {
	envs := map[string]map[string]string{
		"prod": {"ANYTHING": "val"},
	}
	issues := CheckScopes(envs, nil)
	if len(issues) != 0 {
		t.Fatalf("expected no issues with empty rules, got %d", len(issues))
	}
}

func TestCheckScopes_ScopeNotInEnvs(t *testing.T) {
	envs := map[string]map[string]string{
		"dev": {"APP_KEY": "val"},
	}
	rules := []ScopeRule{
		{Scope: "prod", Prefixes: []string{"APP_"}},
	}
	issues := CheckScopes(envs, rules)
	if len(issues) != 0 {
		t.Fatalf("expected no issues when scope not in envs, got %d", len(issues))
	}
}

func TestScopeIssue_String(t *testing.T) {
	issue := ScopeIssue{Key: "DB_PASS", Env: "prod", Scope: "prod", Message: "does not match allowed prefixes: APP_"}
	s := issue.String()
	if !strings.Contains(s, "DB_PASS") || !strings.Contains(s, "prod") {
		t.Errorf("unexpected string: %s", s)
	}
}

func TestFormatScopeIssues_Empty(t *testing.T) {
	out := FormatScopeIssues(nil)
	if out != "no scope violations found" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatScopeIssues_WithIssues(t *testing.T) {
	issues := []ScopeIssue{
		{Key: "DB_PASS", Env: "prod", Scope: "prod", Message: "does not match allowed prefixes: APP_"},
	}
	out := FormatScopeIssues(issues)
	if !strings.Contains(out, "DB_PASS") {
		t.Errorf("expected DB_PASS in output, got: %s", out)
	}
}
