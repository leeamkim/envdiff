package diff

import (
	"strings"
	"testing"
)

func TestCheckDependencies_NoIssues(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	rules := []DependencyRule{
		{Key: "DB_HOST", RefKey: "DB_PORT"},
	}
	issues := CheckDependencies(env, rules)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestCheckDependencies_MissingRef(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
	}
	rules := []DependencyRule{
		{Key: "DB_HOST", RefKey: "DB_PORT"},
	}
	issues := CheckDependencies(env, rules)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Reason != "dependency key is missing" {
		t.Errorf("unexpected reason: %s", issues[0].Reason)
	}
}

func TestCheckDependencies_EmptyRef(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "",
	}
	rules := []DependencyRule{
		{Key: "DB_HOST", RefKey: "DB_PORT"},
	}
	issues := CheckDependencies(env, rules)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Reason != "dependency key is empty" {
		t.Errorf("unexpected reason: %s", issues[0].Reason)
	}
}

func TestCheckDependencies_SkipsAbsentKey(t *testing.T) {
	env := map[string]string{}
	rules := []DependencyRule{
		{Key: "DB_HOST", RefKey: "DB_PORT"},
	}
	issues := CheckDependencies(env, rules)
	if len(issues) != 0 {
		t.Fatalf("expected no issues when key absent, got %d", len(issues))
	}
}

func TestCheckDependencies_SortedOutput(t *testing.T) {
	env := map[string]string{
		"Z_KEY": "val",
		"A_KEY": "val",
	}
	rules := []DependencyRule{
		{Key: "Z_KEY", RefKey: "MISSING_Z"},
		{Key: "A_KEY", RefKey: "MISSING_A"},
	}
	issues := CheckDependencies(env, rules)
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}
	if issues[0].Key > issues[1].Key {
		t.Errorf("issues not sorted: %s > %s", issues[0].Key, issues[1].Key)
	}
}

func TestDependencyIssue_String(t *testing.T) {
	issue := DependencyIssue{Key: "DB_HOST", RefKey: "DB_PORT", Reason: "dependency key is missing"}
	s := issue.String()
	if !strings.Contains(s, "DB_HOST") || !strings.Contains(s, "DB_PORT") {
		t.Errorf("unexpected string: %s", s)
	}
}

func TestFormatDependencyIssues_Empty(t *testing.T) {
	out := FormatDependencyIssues(nil)
	if out != "no dependency issues found" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatDependencyIssues_WithIssues(t *testing.T) {
	issues := []DependencyIssue{
		{Key: "DB_HOST", RefKey: "DB_PORT", Reason: "dependency key is missing"},
	}
	out := FormatDependencyIssues(issues)
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got: %s", out)
	}
}
