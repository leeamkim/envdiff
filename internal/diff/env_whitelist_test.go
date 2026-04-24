package diff

import (
	"strings"
	"testing"
)

func TestCheckWhitelist_NoIssues(t *testing.T) {
	env := map[string]string{"APP_ENV": "production", "LOG_LEVEL": "info"}
	rules := []WhitelistRule{WhitelistRuleAllowedKeys([]string{"APP_ENV", "LOG_LEVEL"})}
	issues := CheckWhitelist(env, rules)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestCheckWhitelist_UnknownKey(t *testing.T) {
	env := map[string]string{"APP_ENV": "production", "UNKNOWN_KEY": "value"}
	rules := []WhitelistRule{WhitelistRuleAllowedKeys([]string{"APP_ENV"})}
	issues := CheckWhitelist(env, rules)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "UNKNOWN_KEY" {
		t.Errorf("expected UNKNOWN_KEY, got %s", issues[0].Key)
	}
}

func TestCheckWhitelist_AllowedValues_Valid(t *testing.T) {
	env := map[string]string{"LOG_LEVEL": "info"}
	rules := []WhitelistRule{WhitelistRuleAllowedValues("LOG_LEVEL", []string{"debug", "info", "warn", "error"})}
	issues := CheckWhitelist(env, rules)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestCheckWhitelist_AllowedValues_Invalid(t *testing.T) {
	env := map[string]string{"LOG_LEVEL": "verbose"}
	rules := []WhitelistRule{WhitelistRuleAllowedValues("LOG_LEVEL", []string{"debug", "info", "warn", "error"})}
	issues := CheckWhitelist(env, rules)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if !strings.Contains(issues[0].Reason, "verbose") {
		t.Errorf("expected reason to mention value, got: %s", issues[0].Reason)
	}
}

func TestCheckWhitelist_MultipleRules(t *testing.T) {
	env := map[string]string{"APP_ENV": "staging", "EXTRA": "x"}
	rules := []WhitelistRule{
		WhitelistRuleAllowedKeys([]string{"APP_ENV"}),
		WhitelistRuleAllowedValues("APP_ENV", []string{"production"}),
	}
	issues := CheckWhitelist(env, rules)
	// EXTRA triggers key rule, APP_ENV triggers value rule
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}
}

func TestWhitelistIssue_String(t *testing.T) {
	i := WhitelistIssue{Key: "FOO", Value: "bar", Reason: "key not in allowlist"}
	if i.String() != "FOO: key not in allowlist" {
		t.Errorf("unexpected string: %s", i.String())
	}
}

func TestFormatWhitelistIssues_Empty(t *testing.T) {
	out := FormatWhitelistIssues(nil)
	if out != "no whitelist violations found" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatWhitelistIssues_NonEmpty(t *testing.T) {
	issues := []WhitelistIssue{{Key: "BAD_KEY", Value: "x", Reason: "key not in allowlist"}}
	out := FormatWhitelistIssues(issues)
	if !strings.Contains(out, "VIOLATION") {
		t.Errorf("expected VIOLATION in output, got: %s", out)
	}
	if !strings.Contains(out, "BAD_KEY") {
		t.Errorf("expected BAD_KEY in output, got: %s", out)
	}
}
