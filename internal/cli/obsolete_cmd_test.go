package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeObsoleteEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunObsolete_MissingArgs(t *testing.T) {
	err := RunObsolete([]string{})
	if err == nil || !strings.Contains(err.Error(), "usage") {
		t.Errorf("expected usage error, got %v", err)
	}
}

func TestRunObsolete_InvalidFile(t *testing.T) {
	err := RunObsolete([]string{"/nonexistent/env", "/nonexistent/ref"})
	if err == nil {
		t.Error("expected error for missing files")
	}
}

func TestRunObsolete_NoIssues(t *testing.T) {
	dir := t.TempDir()
	env := writeObsoleteEnv(t, dir, "env", "A=1\nB=2\n")
	ref := writeObsoleteEnv(t, dir, "ref", "A=x\nB=y\nC=z\n")
	err := RunObsolete([]string{env, ref})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunObsolete_DetectsObsoleteKey(t *testing.T) {
	dir := t.TempDir()
	env := writeObsoleteEnv(t, dir, "env", "A=1\nLEGACY=old\n")
	ref := writeObsoleteEnv(t, dir, "ref", "A=x\n")

	// Capture stdout via pipe
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := RunObsolete([]string{env, ref})
	w.Close()
	os.Stdout = old

	var buf strings.Builder
	buf.ReadFrom(r)
	out := buf.String()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "LEGACY") {
		t.Errorf("expected LEGACY in output, got: %q", out)
	}
}

func TestRunObsolete_InvalidRefFile(t *testing.T) {
	dir := t.TempDir()
	env := writeObsoleteEnv(t, dir, "env", "A=1\n")
	err := RunObsolete([]string{env, "/nonexistent/ref"})
	if err == nil {
		t.Error("expected error for missing reference file")
	}
}
