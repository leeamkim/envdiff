package diff

import (
	"strings"
	"testing"
)

func TestCheckRequired_NoIssues(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "8080"}
	issues := CheckRequired(env, []string{"HOST", "PORT"})
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestCheckRequired_MissingKey(t *testing.T) {
	env := map[string]string{"HOST": "localhost"}
	issues := CheckRequired(env, []string{"HOST", "PORT"})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "PORT" {
		t.Errorf("expected PORT, got %s", issues[0].Key)
	}
}

func TestCheckRequired_EmptyValue(t *testing.T) {
	env := map[string]string{"HOST": "", "PORT": "8080"}
	issues := CheckRequired(env, []string{"HOST", "PORT"})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "HOST" {
		t.Errorf("expected HOST, got %s", issues[0].Key)
	}
}

func TestCheckRequired_MultipleIssues(t *testing.T) {
	env := map[string]string{}
	issues := CheckRequired(env, []string{"A", "B", "C"})
	if len(issues) != 3 {
		t.Fatalf("expected 3 issues, got %d", len(issues))
	}
}

func TestRequiredIssue_String(t *testing.T) {
	iss := RequiredIssue{Key: "SECRET"}
	if iss.String() != "missing required key: SECRET" {
		t.Errorf("unexpected string: %s", iss.String())
	}
}

func TestFormatRequiredIssues_Empty(t *testing.T) {
	out := FormatRequiredIssues(nil)
	if out != "all required keys present" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatRequiredIssues_WithIssues(t *testing.T) {
	issues := []RequiredIssue{{Key: "DB_URL"}, {Key: "API_KEY"}}
	out := FormatRequiredIssues(issues)
	if !strings.Contains(out, "2 required key(s)") {
		t.Errorf("expected count in output, got: %s", out)
	}
	if !strings.Contains(out, "DB_URL") || !strings.Contains(out, "API_KEY") {
		t.Errorf("expected keys in output, got: %s", out)
	}
}
