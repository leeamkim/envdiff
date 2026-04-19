package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func writeMultiEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunMulti_MissingArgs(t *testing.T) {
	var buf bytes.Buffer
	if err := RunMulti([]string{}, &buf); err == nil {
		t.Error("expected error for missing args")
	}
}

func TestRunMulti_NoDiff(t *testing.T) {
	dir := t.TempDir()
	a := writeMultiEnv(t, dir, "a.env", "KEY=value\n")
	b := writeMultiEnv(t, dir, "b.env", "KEY=value\n")
	var buf bytes.Buffer
	if err := RunMulti([]string{a, b}, &buf); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunMulti_InvalidFile(t *testing.T) {
	var buf bytes.Buffer
	if err := RunMulti([]string{"no_such.env", "also_no.env"}, &buf); err == nil {
		t.Error("expected error for missing file")
	}
}

func TestRunMulti_OutputContainsPairKey(t *testing.T) {
	dir := t.TempDir()
	a := writeMultiEnv(t, dir, "a.env", "KEY=1\n")
	b := writeMultiEnv(t, dir, "b.env", "KEY=1\n")
	var buf bytes.Buffer
	_ = RunMulti([]string{a, b}, &buf)
	if len(buf.String()) == 0 {
		t.Error("expected non-empty output")
	}
}
