package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func writeSecretsEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write %s: %v", p, err)
	}
	return p
}

func TestRunSecrets_MissingArgs(t *testing.T) {
	err := RunSecrets([]string{})
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRunSecrets_InvalidFile(t *testing.T) {
	err := RunSecrets([]string{"/nonexistent/.env"})
	if err == nil {
		t.Fatal("expected error for invalid file")
	}
}

func TestRunSecrets_NoIssues(t *testing.T) {
	dir := t.TempDir()
	f := writeSecretsEnv(t, dir, ".env", "APP_NAME=myapp\nPORT=8080\n")
	err := RunSecrets([]string{f})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestRunSecrets_WeakValue(t *testing.T) {
	dir := t.TempDir()
	f := writeSecretsEnv(t, dir, ".env", "DB_PASSWORD=password\nAPP_NAME=myapp\n")
	err := RunSecrets([]string{f})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunSecrets_EmptySensitiveKey(t *testing.T) {
	dir := t.TempDir()
	f := writeSecretsEnv(t, dir, ".env", "API_TOKEN=\nAPP_NAME=myapp\n")
	err := RunSecrets([]string{f})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunSecrets_OnlyWeakFlag(t *testing.T) {
	dir := t.TempDir()
	f := writeSecretsEnv(t, dir, ".env", "DB_PASSWORD=admin\n")
	err := RunSecrets([]string{f, "--weak"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
