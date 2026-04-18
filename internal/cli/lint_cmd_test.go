package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeLintEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunLint_NoIssues(t *testing.T) {
	f := writeLintEnv(t, "HOST=localhost\nPORT=8080\n")
	var buf bytes.Buffer
	err := RunLint(LintOptions{Files: []string{f}, NoEmpty: true, Out: &buf})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(buf.String(), "No lint issues") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestRunLint_EmptyValue(t *testing.T) {
	f := writeLintEnv(t, "API_KEY=\nHOST=localhost\n")
	var buf bytes.Buffer
	err := RunLint(LintOptions{Files: []string{f}, NoEmpty: true, Out: &buf})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(buf.String(), "API_KEY") {
		t.Errorf("expected API_KEY in output, got: %s", buf.String())
	}
}

func TestRunLint_NoFiles(t *testing.T) {
	err := RunLint(LintOptions{})
	if err == nil {
		t.Fatal("expected error for no files")
	}
}

func TestRunLint_DefaultRules(t *testing.T) {
	f := writeLintEnv(t, "foo=CHANGEME\n")
	var buf bytes.Buffer
	err := RunLint(LintOptions{Files: []string{f}, Out: &buf})
	if err == nil {
		t.Fatal("expected lint issues")
	}
	out := buf.String()
	if !strings.Contains(out, "foo") {
		t.Errorf("expected foo in output: %s", out)
	}
}

func TestRunLint_InvalidFile(t *testing.T) {
	err := RunLint(LintOptions{Files: []string{"/nonexistent/.env"}, Out: &bytes.Buffer{}})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
