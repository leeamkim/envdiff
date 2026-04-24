package diff

import (
	"strings"
	"testing"
)

func TestCheckSecrets_NoIssues(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"PORT":     "8080",
	}
	issues := CheckSecrets(env, []SecretRule{SecretRuleWeakValue, SecretRuleKeyWithoutValue})
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d", len(issues))
	}
}

func TestCheckSecrets_WeakValue(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "password",
	}
	issues := CheckSecrets(env, []SecretRule{SecretRuleWeakValue})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "DB_PASSWORD" {
		t.Errorf("expected key DB_PASSWORD, got %s", issues[0].Key)
	}
	if !strings.Contains(issues[0].Reason, "weak secret") {
		t.Errorf("expected reason to mention weak secret, got %s", issues[0].Reason)
	}
}

func TestCheckSecrets_EmptySensitiveKey(t *testing.T) {
	env := map[string]string{
		"API_TOKEN": "",
		"APP_NAME": "myapp",
	}
	issues := CheckSecrets(env, []SecretRule{SecretRuleKeyWithoutValue})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "API_TOKEN" {
		t.Errorf("expected key API_TOKEN, got %s", issues[0].Key)
	}
}

func TestCheckSecrets_SortedOutput(t *testing.T) {
	env := map[string]string{
		"Z_SECRET": "changeme",
		"A_SECRET": "admin",
	}
	issues := CheckSecrets(env, []SecretRule{SecretRuleWeakValue})
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}
	if issues[0].Key != "A_SECRET" {
		t.Errorf("expected first issue key to be A_SECRET, got %s", issues[0].Key)
	}
}

func TestSecretIssue_String(t *testing.T) {
	issue := SecretIssue{Key: "API_KEY", Value: "secret", Reason: "matches known weak secret"}
	s := issue.String()
	if !strings.Contains(s, "[SECRET]") {
		t.Errorf("expected [SECRET] prefix, got %s", s)
	}
	if !strings.Contains(s, "API_KEY") {
		t.Errorf("expected key in output, got %s", s)
	}
}

func TestFormatSecretIssues_Empty(t *testing.T) {
	out := FormatSecretIssues(nil)
	if out != "no secret issues found" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatSecretIssues_Multiple(t *testing.T) {
	issues := []SecretIssue{
		{Key: "DB_PASS", Value: "admin", Reason: "weak"},
		{Key: "API_KEY", Value: "", Reason: "empty"},
	}
	out := FormatSecretIssues(issues)
	if !strings.Contains(out, "DB_PASS") {
		t.Errorf("expected DB_PASS in output")
	}
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in output")
	}
}
