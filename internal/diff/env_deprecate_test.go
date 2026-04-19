package diff

import (
	"strings"
	"testing"
)

var deprecateRules = []DeprecationRule{
	{Key: "OLD_API_KEY", Reason: "legacy auth", Replace: "API_KEY"},
	{Key: "LEGACY_HOST", Reason: "use new host config", Replace: ""},
}

func TestCheckDeprecations_NoIssues(t *testing.T) {
	env := map[string]string{"API_KEY": "abc", "HOST": "localhost"}
	issues := CheckDeprecations(env, deprecateRules)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestCheckDeprecations_SingleMatch(t *testing.T) {
	env := map[string]string{"OLD_API_KEY": "old", "HOST": "localhost"}
	issues := CheckDeprecations(env, deprecateRules)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "OLD_API_KEY" {
		t.Errorf("unexpected key: %s", issues[0].Key)
	}
	if issues[0].Replace != "API_KEY" {
		t.Errorf("unexpected replace: %s", issues[0].Replace)
	}
}

func TestCheckDeprecations_MultipleMatches(t *testing.T) {
	env := map[string]string{"OLD_API_KEY": "old", "LEGACY_HOST": "h"}
	issues := CheckDeprecations(env, deprecateRules)
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}
}

func TestDeprecationIssue_String_WithReplace(t *testing.T) {
	i := DeprecationIssue{Key: "OLD_API_KEY", Reason: "legacy auth", Replace: "API_KEY"}
	s := i.String()
	if !strings.Contains(s, "use API_KEY instead") {
		t.Errorf("expected replace hint, got: %s", s)
	}
}

func TestDeprecationIssue_String_NoReplace(t *testing.T) {
	i := DeprecationIssue{Key: "LEGACY_HOST", Reason: "use new host config"}
	s := i.String()
	if strings.Contains(s, "instead") {
		t.Errorf("did not expect replace hint, got: %s", s)
	}
}

func TestFormatDeprecationIssues_Empty(t *testing.T) {
	out := FormatDeprecationIssues(nil)
	if out != "no deprecated keys found" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatDeprecationIssues_NonEmpty(t *testing.T) {
	issues := []DeprecationIssue{{Key: "OLD_API_KEY", Reason: "legacy", Replace: "API_KEY"}}
	out := FormatDeprecationIssues(issues)
	if !strings.Contains(out, "OLD_API_KEY") {
		t.Errorf("expected key in output: %s", out)
	}
}
