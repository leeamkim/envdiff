package diff

import (
	"testing"
)

func TestInferType_Boolean(t *testing.T) {
	for _, v := range []string{"true", "false", "True", "FALSE"} {
		if got := InferType(v); got != TypeBoolean {
			t.Errorf("InferType(%q) = %q, want %q", v, got, TypeBoolean)
		}
	}
}

func TestInferType_Integer(t *testing.T) {
	for _, v := range []string{"0", "42", "-7", "1000"} {
		if got := InferType(v); got != TypeInteger {
			t.Errorf("InferType(%q) = %q, want %q", v, got, TypeInteger)
		}
	}
}

func TestInferType_Float(t *testing.T) {
	for _, v := range []string{"3.14", "-0.5", "100.0"} {
		if got := InferType(v); got != TypeFloat {
			t.Errorf("InferType(%q) = %q, want %q", v, got, TypeFloat)
		}
	}
}

func TestInferType_URL(t *testing.T) {
	for _, v := range []string{"http://example.com", "https://api.example.com/v1"} {
		if got := InferType(v); got != TypeURL {
			t.Errorf("InferType(%q) = %q, want %q", v, got, TypeURL)
		}
	}
}

func TestInferType_Empty(t *testing.T) {
	if got := InferType(""); got != TypeEmpty {
		t.Errorf("InferType(\"\") = %q, want %q", got, TypeEmpty)
	}
}

func TestInferType_String(t *testing.T) {
	for _, v := range []string{"hello", "some_value", "changeme"} {
		if got := InferType(v); got != TypeString {
			t.Errorf("InferType(%q) = %q, want %q", v, got, TypeString)
		}
	}
}

func TestInferTypes_ReturnsAllKeys(t *testing.T) {
	env := map[string]string{"PORT": "8080", "DEBUG": "true", "HOST": ""}
	entries := InferTypes(env)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestCheckTypeConsistency_NoIssues(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"PORT": "8080", "DEBUG": "true"},
		"prod": {"PORT": "443", "DEBUG": "false"},
	}
	issues := CheckTypeConsistency(envs)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d", len(issues))
	}
}

func TestCheckTypeConsistency_Inconsistent(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"PORT": "8080"},
		"prod": {"PORT": "https://proxy.example.com"},
	}
	issues := CheckTypeConsistency(envs)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "PORT" {
		t.Errorf("expected issue for PORT, got %q", issues[0].Key)
	}
}

func TestTypeIssue_String(t *testing.T) {
	issue := TypeIssue{
		Key:   "PORT",
		Types: map[string]TypeHint{"dev": TypeInteger, "prod": TypeURL},
	}
	s := issue.String()
	if s == "" {
		t.Error("expected non-empty string from TypeIssue.String()")
	}
	if !containsStr(s, "PORT") {
		t.Errorf("expected string to contain 'PORT', got: %s", s)
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && stringContains(s, sub))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
