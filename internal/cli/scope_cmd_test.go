package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func writeScopeEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeScopeEnv: %v", err)
	}
	return p
}

func TestRunScope_MissingArgs(t *testing.T) {
	err := RunScope(nil)
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRunScope_InvalidArgFormat(t *testing.T) {
	err := RunScope([]string{"nodeclaration"})
	if err == nil {
		t.Fatal("expected error for invalid arg format")
	}
}

func TestRunScope_InvalidScopeDecl(t *testing.T) {
	dir := t.TempDir()
	f := writeScopeEnv(t, dir, ".env.prod", "APP_HOST=example.com\n")
	err := RunScope([]string{"nocolon=" + f})
	if err == nil {
		t.Fatal("expected error for missing colon in scope declaration")
	}
}

func TestRunScope_InvalidFile(t *testing.T) {
	err := RunScope([]string{"prod:APP_=/nonexistent/.env"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRunScope_NoIssues(t *testing.T) {
	dir := t.TempDir()
	f := writeScopeEnv(t, dir, ".env.prod", "APP_HOST=example.com\nAPP_PORT=8080\n")
	err := RunScope([]string{"prod:APP_=" + f})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestRunScope_WithViolation(t *testing.T) {
	dir := t.TempDir()
	f := writeScopeEnv(t, dir, ".env.prod", "APP_HOST=example.com\nDB_PASS=secret\n")
	err := RunScope([]string{"prod:APP_=" + f})
	if err == nil {
		t.Fatal("expected error for scope violation")
	}
}

func TestRunScope_MultipleEnvs(t *testing.T) {
	dir := t.TempDir()
	prod := writeScopeEnv(t, dir, ".env.prod", "APP_HOST=example.com\n")
	stg := writeScopeEnv(t, dir, ".env.stg", "APP_HOST=stg.example.com\nSVC_URL=http://svc\n")
	err := RunScope([]string{
		"prod:APP_=" + prod,
		"staging:APP_,SVC_=" + stg,
	})
	if err != nil {
		t.Fatalf("expected no error for valid multi-env scopes, got: %v", err)
	}
}
