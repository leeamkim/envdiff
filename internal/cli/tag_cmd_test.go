package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTagEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunTag_MissingArgs(t *testing.T) {
	err := RunTag([]string{})
	if err == nil {
		t.Error("expected error for missing args")
	}
}

func TestRunTag_InvalidFile(t *testing.T) {
	err := RunTag([]string{"/nonexistent/.env"})
	if err == nil {
		t.Error("expected error for invalid file")
	}
}

func TestRunTag_NoIssues(t *testing.T) {
	p := writeTagEnv(t, "APP_NAME=myapp\nHOST=localhost\n")
	err := RunTag([]string{p})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunTag_DefaultRules(t *testing.T) {
	p := writeTagEnv(t, "DB_PASSWORD=secret\nFOO=\nAPI_KEY=CHANGEME\n")
	err := RunTag([]string{p})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunTag_SecretsOnly(t *testing.T) {
	p := writeTagEnv(t, "DB_PASSWORD=secret\nFOO=\n")
	err := RunTag([]string{p, "--secrets"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunTag_EmptyOnly(t *testing.T) {
	p := writeTagEnv(t, "FOO=\nBAR=value\n")
	err := RunTag([]string{p, "--empty"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
