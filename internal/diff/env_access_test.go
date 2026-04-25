package diff

import (
	"strings"
	"testing"
)

func TestCheckAccess_NoIssues(t *testing.T) {
	env := map[string]string{"SECRET_KEY": "abc", "APP_NAME": "myapp"}
	rules := []AccessRule{
		{Prefix: "SECRET_", Level: AccessSecret},
		{Prefix: "APP_", Level: AccessPublic},
	}
	actual := map[string]AccessLevel{
		"SECRET_KEY": AccessSecret,
		"APP_NAME":   AccessPublic,
	}
	issues := CheckAccess(env, rules, actual)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestCheckAccess_WrongLevel(t *testing.T) {
	env := map[string]string{"SECRET_TOKEN": "abc"}
	rules := []AccessRule{{Prefix: "SECRET_", Level: AccessSecret}}
	actual := map[string]AccessLevel{"SECRET_TOKEN": AccessPublic}
	issues := CheckAccess(env, rules, actual)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "SECRET_TOKEN" {
		t.Errorf("unexpected key %q", issues[0].Key)
	}
	if issues[0].Expected != AccessSecret {
		t.Errorf("expected AccessSecret, got %v", issues[0].Expected)
	}
}

func TestCheckAccess_DefaultsToPublic(t *testing.T) {
	env := map[string]string{"SECRET_PASS": "x"}
	rules := []AccessRule{{Prefix: "SECRET_", Level: AccessSecret}}
	// actual map omits the key — should default to AccessPublic
	issues := CheckAccess(env, rules, map[string]AccessLevel{})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Actual != AccessPublic {
		t.Errorf("expected actual=public, got %v", issues[0].Actual)
	}
}

func TestCheckAccess_NoMatchingRule(t *testing.T) {
	env := map[string]string{"UNRELATED": "val"}
	rules := []AccessRule{{Prefix: "SECRET_", Level: AccessSecret}}
	issues := CheckAccess(env, rules, map[string]AccessLevel{})
	if len(issues) != 0 {
		t.Fatalf("expected no issues for unmatched key, got %v", issues)
	}
}

func TestAccessLevel_String(t *testing.T) {
	if AccessPublic.String() != "public" {
		t.Errorf("unexpected %s", AccessPublic.String())
	}
	if AccessInternal.String() != "internal" {
		t.Errorf("unexpected %s", AccessInternal.String())
	}
	if AccessSecret.String() != "secret" {
		t.Errorf("unexpected %s", AccessSecret.String())
	}
}

func TestFormatAccessIssues_Empty(t *testing.T) {
	out := FormatAccessIssues(nil)
	if out != "no access issues found" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatAccessIssues_NonEmpty(t *testing.T) {
	issues := []AccessIssue{
		{Key: "SECRET_KEY", Actual: AccessPublic, Expected: AccessSecret},
	}
	out := FormatAccessIssues(issues)
	if !strings.Contains(out, "SECRET_KEY") {
		t.Errorf("expected key in output, got %q", out)
	}
	if !strings.Contains(out, "secret") {
		t.Errorf("expected 'secret' in output, got %q", out)
	}
}
