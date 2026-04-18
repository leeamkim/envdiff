package diff

import (
	"testing"
)

func TestRedact_SensitiveKeys(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "secret123",
		"APP_NAME":    "myapp",
		"API_TOKEN":   "tok_abc",
	}
	opts := RedactOptions{Patterns: DefaultRedactPatterns}
	out := Redact(env, opts)

	if out["DB_PASSWORD"] != "***" {
		t.Errorf("expected DB_PASSWORD redacted, got %q", out["DB_PASSWORD"])
	}
	if out["API_TOKEN"] != "***" {
		t.Errorf("expected API_TOKEN redacted, got %q", out["API_TOKEN"])
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME unchanged, got %q", out["APP_NAME"])
	}
}

func TestRedact_NoPatterns(t *testing.T) {
	env := map[string]string{"SECRET_KEY": "val"}
	out := Redact(env, RedactOptions{})
	if out["SECRET_KEY"] != "val" {
		t.Errorf("expected value unchanged with no patterns")
	}
}

func TestRedactResult_Mismatched(t *testing.T) {
	r := Result{
		MissingInA: map[string]string{"DB_SECRET": "abc"},
		MissingInB: map[string]string{"HOST": "localhost"},
		Mismatched: map[string][2]string{
			"API_KEY": {"old", "new"},
			"PORT":    {"8080", "9090"},
		},
	}
	opts := RedactOptions{Patterns: DefaultRedactPatterns}
	out := RedactResult(r, opts)

	if out.MissingInA["DB_SECRET"] != "***" {
		t.Errorf("expected DB_SECRET redacted in MissingInA")
	}
	if out.MissingInB["HOST"] != "localhost" {
		t.Errorf("expected HOST unchanged in MissingInB")
	}
	if out.Mismatched["API_KEY"] != ([2]string{"***", "***"}) {
		t.Errorf("expected API_KEY redacted in Mismatched")
	}
	if out.Mismatched["PORT"] != ([2]string{"8080", "9090"}) {
		t.Errorf("expected PORT unchanged in Mismatched")
	}
}

func TestIsSensitive_CaseInsensitive(t *testing.T) {
	if !isSensitive("db_password", DefaultRedactPatterns) {
		t.Error("expected lowercase db_password to be sensitive")
	}
	if isSensitive("APP_ENV", DefaultRedactPatterns) {
		t.Error("expected APP_ENV to not be sensitive")
	}
}
