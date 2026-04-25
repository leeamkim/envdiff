package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func writeAccessEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunAccess_MissingArgs(t *testing.T) {
	err := RunAccess([]string{})
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRunAccess_InvalidFile(t *testing.T) {
	err := RunAccess([]string{"/nonexistent/.env", "SECRET_:secret"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRunAccess_InvalidRuleFormat(t *testing.T) {
	p := writeAccessEnv(t, "SECRET_KEY=abc\n")
	err := RunAccess([]string{p, "BADRULE"})
	if err == nil {
		t.Fatal("expected error for bad rule format")
	}
}

func TestRunAccess_InvalidLevel(t *testing.T) {
	p := writeAccessEnv(t, "SECRET_KEY=abc\n")
	err := RunAccess([]string{p, "SECRET_:superadmin"})
	if err == nil {
		t.Fatal("expected error for unknown level")
	}
}

func TestRunAccess_NoIssues(t *testing.T) {
	// SECRET_KEY inferred as AccessSecret; rule expects secret — no violation.
	p := writeAccessEnv(t, "SECRET_KEY=abc\n")
	err := RunAccess([]string{p, "SECRET_:secret"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRunAccess_DetectsViolation(t *testing.T) {
	// APP_NAME inferred as AccessPublic; rule expects secret — violation.
	p := writeAccessEnv(t, "APP_NAME=myapp\n")
	err := RunAccess([]string{p, "APP_:secret"})
	if err == nil {
		t.Fatal("expected violation error")
	}
}

func TestRunAccess_MultipleRules(t *testing.T) {
	p := writeAccessEnv(t, "SECRET_TOKEN=x\nAPP_NAME=myapp\n")
	// SECRET_TOKEN -> secret (ok), APP_NAME -> public (ok)
	err := RunAccess([]string{p, "SECRET_:secret", "APP_:public"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
