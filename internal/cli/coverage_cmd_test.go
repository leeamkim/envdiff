package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeCoverageEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunCoverage_MissingArgs(t *testing.T) {
	err := RunCoverage([]string{})
	if err == nil || !strings.Contains(err.Error(), "usage") {
		t.Errorf("expected usage error, got %v", err)
	}
}

func TestRunCoverage_InvalidFile(t *testing.T) {
	err := RunCoverage([]string{"/nonexistent/ref.env", "/nonexistent/tgt.env"})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestRunCoverage_PerfectCoverage(t *testing.T) {
	dir := t.TempDir()
	ref := writeCoverageEnv(t, dir, "ref.env", "A=1\nB=2\nC=3\n")
	tgt := writeCoverageEnv(t, dir, "tgt.env", "A=1\nB=2\nC=3\n")

	err := RunCoverage([]string{ref, tgt})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestRunCoverage_LowCoverage(t *testing.T) {
	dir := t.TempDir()
	ref := writeCoverageEnv(t, dir, "ref.env", "A=1\nB=2\nC=3\nD=4\nE=5\n")
	tgt := writeCoverageEnv(t, dir, "tgt.env", "A=1\n")

	err := RunCoverage([]string{ref, tgt})
	if err == nil {
		t.Error("expected error for failing coverage")
	}
	if !strings.Contains(err.Error(), "failing coverage") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRunCoverage_MultipleTargets(t *testing.T) {
	dir := t.TempDir()
	ref := writeCoverageEnv(t, dir, "ref.env", "A=1\nB=2\n")
	tgt1 := writeCoverageEnv(t, dir, "tgt1.env", "A=1\nB=2\n")
	tgt2 := writeCoverageEnv(t, dir, "tgt2.env", "A=1\nB=2\n")

	err := RunCoverage([]string{ref, tgt1, tgt2})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
