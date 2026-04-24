package diff

import (
	"strings"
	"testing"
)

func TestFindDuplicateConflicts_NoConflicts(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"HOST": "localhost", "PORT": "8080"},
		"prod": {"HOST": "localhost", "PORT": "8080"},
	}
	issues := FindDuplicateConflicts(envs)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestFindDuplicateConflicts_SingleConflict(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"HOST": "localhost"},
		"prod": {"HOST": "example.com"},
	}
	issues := FindDuplicateConflicts(envs)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "HOST" {
		t.Errorf("expected key HOST, got %s", issues[0].Key)
	}
}

func TestFindDuplicateConflicts_SkipsMissingKeys(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"HOST": "localhost"},
		"prod": {"PORT": "443"},
	}
	// HOST only in dev, PORT only in prod — no shared keys, no conflicts.
	issues := FindDuplicateConflicts(envs)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestFindDuplicateConflicts_ThreeEnvs(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":     {"DB": "dev_db"},
		"staging": {"DB": "staging_db"},
		"prod":    {"DB": "prod_db"},
	}
	issues := FindDuplicateConflicts(envs)
	// 3 pairs: dev/staging, dev/prod, staging/prod
	if len(issues) != 3 {
		t.Fatalf("expected 3 issues, got %d", len(issues))
	}
}

func TestFindDuplicateConflicts_SingleEnv(t *testing.T) {
	envs := map[string]map[string]string{
		"dev": {"HOST": "localhost"},
	}
	issues := FindDuplicateConflicts(envs)
	if len(issues) != 0 {
		t.Fatalf("expected no issues for single env, got %d", len(issues))
	}
}

func TestDuplicateIssue_String(t *testing.T) {
	issue := DuplicateIssue{Key: "HOST", EnvA: "dev", ValueA: "localhost", EnvB: "prod", ValueB: "example.com"}
	s := issue.String()
	if !strings.Contains(s, "HOST") {
		t.Errorf("expected key in string, got: %s", s)
	}
	if !strings.Contains(s, "dev") || !strings.Contains(s, "prod") {
		t.Errorf("expected env names in string, got: %s", s)
	}
}

func TestFormatDuplicateIssues_Empty(t *testing.T) {
	out := FormatDuplicateIssues(nil)
	if out != "no duplicate conflicts found" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatDuplicateIssues_WithIssues(t *testing.T) {
	issues := []DuplicateIssue{
		{Key: "PORT", EnvA: "dev", ValueA: "3000", EnvB: "prod", ValueB: "443"},
	}
	out := FormatDuplicateIssues(issues)
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected PORT in output, got: %s", out)
	}
}
