package diff

import (
	"strings"
	"testing"
)

func TestCheckQuotas_NoIssues(t *testing.T) {
	env := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
	}
	rules := []QuotaRule{{Prefix: "APP_", MaxKeys: 5}}
	issues := CheckQuotas(env, rules)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestCheckQuotas_ExceedsLimit(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"DB_NAME": "mydb",
	}
	rules := []QuotaRule{{Prefix: "DB_", MaxKeys: 2}}
	issues := CheckQuotas(env, rules)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Actual != 3 {
		t.Errorf("expected actual=3, got %d", issues[0].Actual)
	}
	if issues[0].Allowed != 2 {
		t.Errorf("expected allowed=2, got %d", issues[0].Allowed)
	}
}

func TestCheckQuotas_MultipleRules(t *testing.T) {
	env := map[string]string{
		"APP_A": "1",
		"APP_B": "2",
		"SVC_X": "x",
		"SVC_Y": "y",
		"SVC_Z": "z",
	}
	rules := []QuotaRule{
		{Prefix: "APP_", MaxKeys: 3},
		{Prefix: "SVC_", MaxKeys: 2},
	}
	issues := CheckQuotas(env, rules)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Prefix != "SVC_" {
		t.Errorf("expected SVC_ violation, got %q", issues[0].Prefix)
	}
}

func TestCheckQuotas_SortedOutput(t *testing.T) {
	env := map[string]string{
		"Z_A": "1", "Z_B": "2", "Z_C": "3",
		"A_X": "x", "A_Y": "y", "A_Z": "z",
	}
	rules := []QuotaRule{
		{Prefix: "Z_", MaxKeys: 1},
		{Prefix: "A_", MaxKeys: 1},
	}
	issues := CheckQuotas(env, rules)
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}
	if issues[0].Prefix >= issues[1].Prefix {
		t.Errorf("expected sorted output, got %q then %q", issues[0].Prefix, issues[1].Prefix)
	}
}

func TestQuotaIssue_String(t *testing.T) {
	iss := QuotaIssue{Prefix: "DB_", Allowed: 2, Actual: 4, Keys: []string{"DB_A", "DB_B", "DB_C", "DB_D"}}
	s := iss.String()
	if !strings.Contains(s, "DB_") || !strings.Contains(s, "2") || !strings.Contains(s, "4") {
		t.Errorf("unexpected String output: %q", s)
	}
}

func TestFormatQuotaIssues_Empty(t *testing.T) {
	out := FormatQuotaIssues(nil)
	if !strings.Contains(out, "no quota") {
		t.Errorf("expected empty message, got %q", out)
	}
}

func TestFormatQuotaIssues_WithIssues(t *testing.T) {
	issues := []QuotaIssue{
		{Prefix: "APP_", Allowed: 1, Actual: 3, Keys: []string{"APP_A", "APP_B", "APP_C"}},
	}
	out := FormatQuotaIssues(issues)
	if !strings.Contains(out, "1 quota violation") {
		t.Errorf("expected violation count in output, got %q", out)
	}
	if !strings.Contains(out, "APP_") {
		t.Errorf("expected prefix in output, got %q", out)
	}
}
