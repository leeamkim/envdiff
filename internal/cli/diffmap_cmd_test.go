package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeDiffMapEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunDiffMap_MissingArgs(t *testing.T) {
	err := RunDiffMap([]string{"prod=/tmp/x"}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "at least 2") {
		t.Errorf("expected at-least-2 error, got %v", err)
	}
}

func TestRunDiffMap_InvalidArg(t *testing.T) {
	err := RunDiffMap([]string{"prod=/tmp/x", "badarg"}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "invalid argument") {
		t.Errorf("expected invalid argument error, got %v", err)
	}
}

func TestRunDiffMap_NoDiff(t *testing.T) {
	dir := t.TempDir()
	p1 := writeDiffMapEnv(t, dir, "prod.env", "A=1\nB=2\n")
	p2 := writeDiffMapEnv(t, dir, "dev.env", "A=1\nB=2\n")
	var buf bytes.Buffer
	err := RunDiffMap([]string{"prod=" + p1, "dev=" + p2}, &buf)
	if err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if strings.Contains(out, "!") {
		t.Errorf("expected no inconsistency markers, got:\n%s", out)
	}
}

func TestRunDiffMap_Inconsistent(t *testing.T) {
	dir := t.TempDir()
	p1 := writeDiffMapEnv(t, dir, "prod.env", "A=hello\n")
	p2 := writeDiffMapEnv(t, dir, "dev.env", "A=world\n")
	var buf bytes.Buffer
	err := RunDiffMap([]string{"prod=" + p1, "dev=" + p2}, &buf)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "!") {
		t.Errorf("expected inconsistency marker")
	}
}

func TestRunDiffMap_InvalidFile(t *testing.T) {
	err := RunDiffMap([]string{"prod=/nonexistent", "dev=/nonexistent2"}, &bytes.Buffer{})
	if err == nil {
		t.Error("expected error for missing file")
	}
}
