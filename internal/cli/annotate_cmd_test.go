package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeAnnotateEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunAnnotate_NoIssues(t *testing.T) {
	a := writeAnnotateEnv(t, "KEY=val\n")
	b := writeAnnotateEnv(t, "KEY=val\n")
	var buf bytes.Buffer
	if err := RunAnnotate([]string{a, b}, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no issues") {
		t.Errorf("expected 'no issues', got %q", buf.String())
	}
}

func TestRunAnnotate_MissingInB(t *testing.T) {
	a := writeAnnotateEnv(t, "ONLY_A=val\nSHARED=x\n")
	b := writeAnnotateEnv(t, "SHARED=x\n")
	var buf bytes.Buffer
	if err := RunAnnotate([]string{a, b}, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "ONLY_A") {
		t.Errorf("expected ONLY_A in output, got %q", buf.String())
	}
}

func TestRunAnnotate_Mismatched(t *testing.T) {
	a := writeAnnotateEnv(t, "KEY=foo\n")
	b := writeAnnotateEnv(t, "KEY=bar\n")
	var buf bytes.Buffer
	if err := RunAnnotate([]string{a, b}, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "mismatch") {
		t.Errorf("expected mismatch note, got %q", buf.String())
	}
}

func TestRunAnnotate_MissingArgs(t *testing.T) {
	var buf bytes.Buffer
	if err := RunAnnotate([]string{"only_one"}, &buf); err == nil {
		t.Error("expected error for missing args")
	}
}

func TestRunAnnotate_InvalidFile(t *testing.T) {
	var buf bytes.Buffer
	if err := RunAnnotate([]string{"/nonexistent/.env", "/nonexistent2/.env"}, &buf); err == nil {
		t.Error("expected error for invalid file")
	}
}
