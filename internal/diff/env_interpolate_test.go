package diff

import (
	"strings"
	"testing"
)

func TestInterpolate_NoRefs(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
	}
	out, issues := Interpolate(env)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
	if out["HOST"] != "localhost" || out["PORT"] != "5432" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestInterpolate_ResolvesRef(t *testing.T) {
	env := map[string]string{
		"BASE_URL": "http://localhost",
		"API_URL":  "${BASE_URL}/api",
	}
	out, issues := Interpolate(env)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
	if out["API_URL"] != "http://localhost/api" {
		t.Errorf("expected expanded URL, got %q", out["API_URL"])
	}
}

func TestInterpolate_UnresolvedRef(t *testing.T) {
	env := map[string]string{
		"API_URL": "${MISSING_VAR}/api",
	}
	out, issues := Interpolate(env)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "API_URL" || issues[0].Ref != "MISSING_VAR" {
		t.Errorf("unexpected issue: %v", issues[0])
	}
	if out["API_URL"] != "/api" {
		t.Errorf("expected token removed, got %q", out["API_URL"])
	}
}

func TestInterpolate_MultipleRefs(t *testing.T) {
	env := map[string]string{
		"SCHEME": "https",
		"HOST":   "example.com",
		"URL":    "${SCHEME}://${HOST}",
	}
	out, issues := Interpolate(env)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
	if out["URL"] != "https://example.com" {
		t.Errorf("expected full URL, got %q", out["URL"])
	}
}

func TestInterpolate_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{
		"A": "${B}",
		"B": "hello",
	}
	original := env["A"]
	Interpolate(env)
	if env["A"] != original {
		t.Error("input map was mutated")
	}
}

func TestInterpolateIssue_String(t *testing.T) {
	iss := InterpolateIssue{Key: "URL", Ref: "MISSING"}
	got := iss.String()
	if !strings.Contains(got, "URL") || !strings.Contains(got, "MISSING") {
		t.Errorf("unexpected string: %q", got)
	}
}

func TestFormatInterpolateIssues_Empty(t *testing.T) {
	out := FormatInterpolateIssues(nil)
	if !strings.Contains(out, "no interpolation issues") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatInterpolateIssues_WithIssues(t *testing.T) {
	issues := []InterpolateIssue{
		{Key: "API_URL", Ref: "BASE"},
	}
	out := FormatInterpolateIssues(issues)
	if !strings.Contains(out, "API_URL") || !strings.Contains(out, "BASE") {
		t.Errorf("expected issue in output, got: %q", out)
	}
}
