package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeChainEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunChain_MissingArgs(t *testing.T) {
	var buf bytes.Buffer
	err := RunChain([]string{"dev=/tmp/a.env"}, &buf)
	if err == nil {
		t.Fatal("expected error for too few args")
	}
}

func TestRunChain_InvalidArgFormat(t *testing.T) {
	dir := t.TempDir()
	f := writeChainEnv(t, dir, "dev.env", "PORT=3000\n")
	var buf bytes.Buffer
	err := RunChain([]string{"dev=" + f, "prodonly"}, &buf)
	if err == nil {
		t.Fatal("expected error for invalid arg format")
	}
}

func TestRunChain_InvalidFile(t *testing.T) {
	var buf bytes.Buffer
	err := RunChain([]string{"dev=/nonexistent.env", "prod=/also-missing.env"}, &buf)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRunChain_NoDiff(t *testing.T) {
	dir := t.TempDir()
	dev := writeChainEnv(t, dir, "dev.env", "HOST=localhost\nPORT=3000\n")
	prod := writeChainEnv(t, dir, "prod.env", "HOST=localhost\nPORT=3000\n")
	var buf bytes.Buffer
	err := RunChain([]string{"dev=" + dev, "prod=" + prod}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "override") {
		t.Errorf("expected no override, got: %s", out)
	}
}

func TestRunChain_WithOverride(t *testing.T) {
	dir := t.TempDir()
	dev := writeChainEnv(t, dir, "dev.env", "DB=dev-db\n")
	prod := writeChainEnv(t, dir, "prod.env", "DB=prod-db\n")
	var buf bytes.Buffer
	err := RunChain([]string{"dev=" + dev, "prod=" + prod}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "override") {
		t.Errorf("expected 'override' in output, got: %s", out)
	}
}

func TestRunChain_ThreeEnvs(t *testing.T) {
	dir := t.TempDir()
	base := writeChainEnv(t, dir, "base.env", "SECRET=base\nHOST=localhost\n")
	stage := writeChainEnv(t, dir, "staging.env", "SECRET=stage\nHOST=localhost\n")
	prod := writeChainEnv(t, dir, "prod.env", "SECRET=prod\nHOST=prod.example.com\n")
	var buf bytes.Buffer
	err := RunChain([]string{"base=" + base, "staging=" + stage, "prod=" + prod}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "SECRET") {
		t.Errorf("expected SECRET in output, got: %s", out)
	}
	if !strings.Contains(out, "HOST") {
		t.Errorf("expected HOST in output, got: %s", out)
	}
}
