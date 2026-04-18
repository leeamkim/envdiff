package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemplateEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunTemplate_NoIssues(t *testing.T) {
	dir := t.TempDir()
	tmpl := writeTemplateEnv(t, dir, "tmpl.env", "A=\nB=\n")
	env := writeTemplateEnv(t, dir, "env.env", "A=1\nB=2\n")
	code := RunTemplate([]string{tmpl, env})
	if code != 0 {
		t.Errorf("expected exit 0, got %d", code)
	}
}

func TestRunTemplate_MissingKey(t *testing.T) {
	dir := t.TempDir()
	tmpl := writeTemplateEnv(t, dir, "tmpl.env", "A=\nB=\nC=\n")
	env := writeTemplateEnv(t, dir, "env.env", "A=1\n")
	code := RunTemplate([]string{tmpl, env})
	if code != 1 {
		t.Errorf("expected exit 1, got %d", code)
	}
}

func TestRunTemplate_StrictExtraKey(t *testing.T) {
	dir := t.TempDir()
	tmpl := writeTemplateEnv(t, dir, "tmpl.env", "A=\n")
	env := writeTemplateEnv(t, dir, "env.env", "A=1\nEXTRA=x\n")
	code := RunTemplate([]string{"--strict", tmpl, env})
	if code != 1 {
		t.Errorf("expected exit 1, got %d", code)
	}
}

func TestRunTemplate_MissingArgs(t *testing.T) {
	code := RunTemplate([]string{})
	if code != 1 {
		t.Errorf("expected exit 1, got %d", code)
	}
}

func TestRunTemplate_InvalidFile(t *testing.T) {
	dir := t.TempDir()
	tmpl := writeTemplateEnv(t, dir, "tmpl.env", "A=\n")
	code := RunTemplate([]string{tmpl, "/nonexistent/env.env"})
	if code != 1 {
		t.Errorf("expected exit 1, got %d", code)
	}
}
