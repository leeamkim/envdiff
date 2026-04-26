package diff

import (
	"strings"
	"testing"
)

func TestCheckLifecycle_NoIssues(t *testing.T) {
	env := map[string]string{"APP_NAME": "myapp", "PORT": "8080"}
	rules := []LifecycleRule{
		{Pattern: "OLD_*", Stage: StageDeprecated},
	}
	issues := CheckLifecycle(env, rules)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestCheckLifecycle_DeprecatedKey(t *testing.T) {
	env := map[string]string{"OLD_API_KEY": "abc", "APP_NAME": "myapp"}
	rules := []LifecycleRule{
		{Pattern: "OLD_*", Stage: StageDeprecated},
	}
	issues := CheckLifecycle(env, rules)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "OLD_API_KEY" {
		t.Errorf("expected key OLD_API_KEY, got %s", issues[0].Key)
	}
	if issues[0].Stage != StageDeprecated {
		t.Errorf("expected stage deprecated, got %s", issues[0].Stage)
	}
}

func TestCheckLifecycle_RetiredKey(t *testing.T) {
	env := map[string]string{"LEGACY_TOKEN": "xyz"}
	rules := []LifecycleRule{
		{Pattern: "LEGACY_TOKEN", Stage: StageRetired},
	}
	issues := CheckLifecycle(env, rules)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if !strings.Contains(issues[0].Note, "retired") {
		t.Errorf("expected retired note, got %s", issues[0].Note)
	}
}

func TestCheckLifecycle_ActiveAndExperimentalIgnored(t *testing.T) {
	env := map[string]string{"BETA_FEATURE": "on", "NEW_KEY": "val"}
	rules := []LifecycleRule{
		{Pattern: "BETA_*", Stage: StageExperimental},
		{Pattern: "NEW_KEY", Stage: StageActive},
	}
	issues := CheckLifecycle(env, rules)
	if len(issues) != 0 {
		t.Fatalf("expected no issues for active/experimental, got %d", len(issues))
	}
}

func TestCheckLifecycle_SortedOutput(t *testing.T) {
	env := map[string]string{"OLD_Z": "1", "OLD_A": "2", "OLD_M": "3"}
	rules := []LifecycleRule{
		{Pattern: "OLD_*", Stage: StageDeprecated},
	}
	issues := CheckLifecycle(env, rules)
	if len(issues) != 3 {
		t.Fatalf("expected 3 issues, got %d", len(issues))
	}
	if issues[0].Key != "OLD_A" || issues[1].Key != "OLD_M" || issues[2].Key != "OLD_Z" {
		t.Errorf("expected sorted keys, got %v %v %v", issues[0].Key, issues[1].Key, issues[2].Key)
	}
}

func TestLifecycleIssue_String(t *testing.T) {
	issue := LifecycleIssue{Key: "OLD_KEY", Stage: StageDeprecated, Note: "migrate soon"}
	s := issue.String()
	if !strings.Contains(s, "OLD_KEY") || !strings.Contains(s, "deprecated") {
		t.Errorf("unexpected string: %s", s)
	}
}

func TestFormatLifecycleIssues_Empty(t *testing.T) {
	out := FormatLifecycleIssues(nil)
	if out != "no lifecycle issues found" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatLifecycleIssues_NonEmpty(t *testing.T) {
	issues := []LifecycleIssue{
		{Key: "OLD_KEY", Stage: StageRetired, Note: "remove it"},
	}
	out := FormatLifecycleIssues(issues)
	if !strings.Contains(out, "OLD_KEY") {
		t.Errorf("expected key in output, got: %s", out)
	}
}
