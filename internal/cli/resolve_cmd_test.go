package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func writeResolveEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunResolve_MissingArgs(t *testing.T) {
	if err := RunResolve([]string{}); err == nil {
		t.Error("expected error for missing args")
	}
}

func TestRunResolve_UnknownStrategy(t *testing.T) {
	dir := t.TempDir()
	a := writeResolveEnv(t, dir, "a.env", "FOO=a\n")
	b := writeResolveEnv(t, dir, "b.env", "FOO=b\n")
	if err := RunResolve([]string{"bad-strategy", a, b}); err == nil {
		t.Error("expected error for unknown strategy")
	}
}

func TestRunResolve_PreferA_NoConflict(t *testing.T) {
	dir := t.TempDir()
	a := writeResolveEnv(t, dir, "a.env", "FOO=from-a\nBAR=shared\n")
	b := writeResolveEnv(t, dir, "b.env", "FOO=from-b\nBAR=shared\n")
	if err := RunResolve([]string{"prefer-a", a, b}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunResolve_RequireMatch_Conflict(t *testing.T) {
	dir := t.TempDir()
	a := writeResolveEnv(t, dir, "a.env", "FOO=aaa\n")
	b := writeResolveEnv(t, dir, "b.env", "FOO=bbb\n")
	if err := RunResolve([]string{"require-match", a, b}); err == nil {
		t.Error("expected conflict error")
	}
}

func TestRunResolve_InvalidFile(t *testing.T) {
	if err := RunResolve([]string{"prefer-a", "/no/such/file.env", "/no/such/other.env"}); err == nil {
		t.Error("expected error for missing file")
	}
}
