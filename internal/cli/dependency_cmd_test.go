package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func writeDependencyEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunDependency_MissingArgs(t *testing.T) {
	err := RunDependency([]string{})
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRunDependency_InvalidFile(t *testing.T) {
	err := RunDependency([]string{"/nonexistent/.env", "DB_HOST:DB_PORT"})
	if err == nil {
		t.Fatal("expected error for invalid file")
	}
}

func TestRunDependency_InvalidRuleFormat(t *testing.T) {
	p := writeDependencyEnv(t, "DB_HOST=localhost\n")
	err := RunDependency([]string{p, "BADRULE"})
	if err == nil {
		t.Fatal("expected error for bad rule format")
	}
}

func TestRunDependency_NoIssues(t *testing.T) {
	p := writeDependencyEnv(t, "DB_HOST=localhost\nDB_PORT=5432\n")
	err := RunDependency([]string{p, "DB_HOST:DB_PORT"})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestRunDependency_MissingRef(t *testing.T) {
	p := writeDependencyEnv(t, "DB_HOST=localhost\n")
	err := RunDependency([]string{p, "DB_HOST:DB_PORT"})
	if err == nil {
		t.Fatal("expected error when dependency ref is missing")
	}
}

func TestRunDependency_EmptyRef(t *testing.T) {
	p := writeDependencyEnv(t, "DB_HOST=localhost\nDB_PORT=\n")
	err := RunDependency([]string{p, "DB_HOST:DB_PORT"})
	if err == nil {
		t.Fatal("expected error when dependency ref is empty")
	}
}

func TestRunDependency_MultipleRules(t *testing.T) {
	p := writeDependencyEnv(t, "APP_URL=http://example.com\nBASE_URL=http://base.com\nAPI_KEY=secret\n")
	err := RunDependency([]string{p, "APP_URL:BASE_URL", "APP_URL:API_KEY"})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}
