package diff

import (
	"strings"
	"testing"
)

func baseImmutableEnvs() map[string]map[string]string {
	return map[string]map[string]string{
		"dev":  {"APP_VERSION": "1.0", "DB_HOST": "localhost", "SECRET": "abc"},
		"prod": {"APP_VERSION": "1.0", "DB_HOST": "prod.db", "SECRET": "xyz"},
	}
}

func TestCheckImmutable_NoIssues(t *testing.T) {
	envs := baseImmutableEnvs()
	issues := CheckImmutable(envs, []string{"APP_VERSION"})
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestCheckImmutable_DetectsChange(t *testing.T) {
	envs := baseImmutableEnvs()
	issues := CheckImmutable(envs, []string{"DB_HOST"})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "DB_HOST" {
		t.Errorf("expected key DB_HOST, got %q", issues[0].Key)
	}
}

func TestCheckImmutable_MultipleKeys(t *testing.T) {
	envs := baseImmutableEnvs()
	issues := CheckImmutable(envs, []string{"DB_HOST", "SECRET"})
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}
}

func TestCheckImmutable_SkipsMissingKey(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"APP_VERSION": "1.0"},
		"prod": {"DB_HOST": "prod.db"},
	}
	issues := CheckImmutable(envs, []string{"APP_VERSION"})
	if len(issues) != 0 {
		t.Fatalf("expected no issues (key missing in one env), got %d", len(issues))
	}
}

func TestCheckImmutable_EmptyKeys(t *testing.T) {
	envs := baseImmutableEnvs()
	issues := CheckImmutable(envs, []string{})
	if len(issues) != 0 {
		t.Fatalf("expected no issues with empty key list, got %d", len(issues))
	}
}

func TestCheckImmutable_SortedOutput(t *testing.T) {
	envs := baseImmutableEnvs()
	issues := CheckImmutable(envs, []string{"SECRET", "DB_HOST"})
	if len(issues) < 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}
	if issues[0].Key > issues[1].Key {
		t.Errorf("issues not sorted by key: %q > %q", issues[0].Key, issues[1].Key)
	}
}

func TestImmutableIssue_String(t *testing.T) {
	issue := ImmutableIssue{Key: "DB_HOST", EnvA: "dev", EnvB: "prod", ValueA: "localhost", ValueB: "prod.db"}
	s := issue.String()
	if !strings.Contains(s, "DB_HOST") || !strings.Contains(s, "localhost") {
		t.Errorf("unexpected String() output: %q", s)
	}
}

func TestFormatImmutableIssues_Empty(t *testing.T) {
	out := FormatImmutableIssues(nil)
	if !strings.Contains(out, "no immutable") {
		t.Errorf("expected no-violation message, got %q", out)
	}
}

func TestFormatImmutableIssues_WithIssues(t *testing.T) {
	issues := []ImmutableIssue{
		{Key: "DB_HOST", EnvA: "dev", EnvB: "prod", ValueA: "localhost", ValueB: "prod.db"},
	}
	out := FormatImmutableIssues(issues)
	if !strings.Contains(out, "1 immutable violation") {
		t.Errorf("expected count in output, got %q", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected key in output, got %q", out)
	}
}
