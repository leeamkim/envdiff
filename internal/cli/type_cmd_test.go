package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTypeEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeTypeEnv: %v", err)
	}
	return p
}

func TestRunType_MissingArgs(t *testing.T) {
	var buf bytes.Buffer
	err := RunType([]string{}, &buf)
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRunType_InvalidArgFormat(t *testing.T) {
	var buf bytes.Buffer
	err := RunType([]string{"noequalsign"}, &buf)
	if err == nil {
		t.Fatal("expected error for invalid arg format")
	}
}

func TestRunType_InvalidFile(t *testing.T) {
	var buf bytes.Buffer
	err := RunType([]string{"dev=/nonexistent/path.env"}, &buf)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRunType_SingleEnv(t *testing.T) {
	dir := t.TempDir()
	p := writeTypeEnv(t, dir, "dev.env", "PORT=8080\nDEBUG=true\nHOST=\n")
	var buf bytes.Buffer
	err := RunType([]string{"dev=" + p}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[dev]") {
		t.Errorf("expected [dev] header in output, got:\n%s", out)
	}
	if !strings.Contains(out, "integer") {
		t.Errorf("expected 'integer' type in output, got:\n%s", out)
	}
	if !strings.Contains(out, "boolean") {
		t.Errorf("expected 'boolean' type in output, got:\n%s", out)
	}
}

func TestRunType_ConsistentTypes(t *testing.T) {
	dir := t.TempDir()
	dev := writeTypeEnv(t, dir, "dev.env", "PORT=8080\nDEBUG=true\n")
	prod := writeTypeEnv(t, dir, "prod.env", "PORT=443\nDEBUG=false\n")
	var buf bytes.Buffer
	err := RunType([]string{"dev=" + dev, "prod=" + prod}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "consistent") {
		t.Errorf("expected 'consistent' message, got:\n%s", buf.String())
	}
}

func TestRunType_InconsistentTypes(t *testing.T) {
	dir := t.TempDir()
	dev := writeTypeEnv(t, dir, "dev.env", "PORT=8080\n")
	prod := writeTypeEnv(t, dir, "prod.env", "PORT=https://proxy.example.com\n")
	var buf bytes.Buffer
	// RunType calls os.Exit(1) on inconsistency; we only check that the output mentions PORT
	// In a real test harness we'd use a subprocess; here we skip the exit and just verify output
	_ = RunType([]string{"dev=" + dev, "prod=" + prod}, &buf)
	if !strings.Contains(buf.String(), "PORT") {
		t.Errorf("expected PORT in inconsistency output, got:\n%s", buf.String())
	}
}
