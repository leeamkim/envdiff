package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func writeSchemaEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunSchema_NoIssues(t *testing.T) {
	p := writeSchemaEnv(t, "HOST=localhost\nPORT=8080\n")
	if err := RunSchema([]string{p}); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestRunSchema_EmptyValue(t *testing.T) {
	p := writeSchemaEnv(t, "HOST=\nPORT=8080\n")
	err := RunSchema([]string{p})
	if err == nil {
		t.Error("expected error for empty value, got nil")
	}
}

func TestRunSchema_RequiredKeysMissing(t *testing.T) {
	p := writeSchemaEnv(t, "HOST=localhost\n")
	err := RunSchema([]string{p, "--require=HOST,PORT,DB_URL"})
	if err == nil {
		t.Error("expected error for missing required keys")
	}
}

func TestRunSchema_RequiredKeysAllPresent(t *testing.T) {
	p := writeSchemaEnv(t, "HOST=localhost\nPORT=8080\n")
	if err := RunSchema([]string{p, "--require=HOST,PORT"}); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestRunSchema_NoArgs(t *testing.T) {
	err := RunSchema([]string{})
	if err == nil {
		t.Error("expected error for missing args")
	}
}

func TestRunSchema_InvalidFile(t *testing.T) {
	err := RunSchema([]string{"/nonexistent/.env"})
	if err == nil {
		t.Error("expected error for invalid file")
	}
}
