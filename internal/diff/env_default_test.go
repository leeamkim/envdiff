package diff

import (
	"strings"
	"testing"
)

func TestCheckDefaults_NoIssues(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "prod.db.example.com",
		"API_KEY": "abc123secret",
	}
	issues := CheckDefaults(env, []DefaultRule{DefaultRuleCommonPlaceholders})
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestCheckDefaults_CommonPlaceholder(t *testing.T) {
	env := map[string]string{
		"DB_PASS": "changeme",
		"API_KEY": "real-key",
	}
	issues := CheckDefaults(env, []DefaultRule{DefaultRuleCommonPlaceholders})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "DB_PASS" {
		t.Errorf("expected DB_PASS, got %s", issues[0].Key)
	}
}

func TestCheckDefaults_CaseInsensitivePlaceholder(t *testing.T) {
	env := map[string]string{
		"SECRET": "CHANGEME",
	}
	issues := CheckDefaults(env, []DefaultRule{DefaultRuleCommonPlaceholders})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
}

func TestCheckDefaults_PinnedValues(t *testing.T) {
	pins := map[string]string{
		"REGION": "us-east-1",
		"TIMEOUT": "30",
	}
	env := map[string]string{
		"REGION":  "us-east-1",
		"TIMEOUT": "60",
		"HOST":    "myhost",
	}
	rule := DefaultRulePinnedValues(pins)
	issues := CheckDefaults(env, []DefaultRule{rule})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "REGION" {
		t.Errorf("expected REGION, got %s", issues[0].Key)
	}
}

func TestDefaultIssue_String(t *testing.T) {
	issue := DefaultIssue{Key: "TOKEN", Value: "todo", DefaultValue: "todo"}
	s := issue.String()
	if !strings.Contains(s, "TOKEN") || !strings.Contains(s, "todo") {
		t.Errorf("unexpected string: %s", s)
	}
}

func TestFormatDefaultIssues_Empty(t *testing.T) {
	out := FormatDefaultIssues(nil)
	if !strings.Contains(out, "no default") {
		t.Errorf("expected no-issues message, got: %s", out)
	}
}

func TestFormatDefaultIssues_WithIssues(t *testing.T) {
	issues := []DefaultIssue{
		{Key: "DB_PASS", Value: "changeme", DefaultValue: "changeme"},
	}
	out := FormatDefaultIssues(issues)
	if !strings.Contains(out, "[DEFAULT]") {
		t.Errorf("expected [DEFAULT] tag in output, got: %s", out)
	}
	if !strings.Contains(out, "DB_PASS") {
		t.Errorf("expected DB_PASS in output, got: %s", out)
	}
}
