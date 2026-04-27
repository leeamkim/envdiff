package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func writeImmutableEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeImmutableEnv: %v", err)
	}
	return p
}

func TestRunImmutable_MissingArgs(t *testing.T) {
	err := RunImmutable([]string{})
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRunImmutable_MissingKeys(t *testing.T) {
	dir := t.TempDir()
	f1 := writeImmutableEnv(t, dir, "dev.env", "APP=1\n")
	f2 := writeImmutableEnv(t, dir, "prod.env", "APP=2\n")
	err := RunImmutable([]string{"dev=" + f1, "prod=" + f2})
	if err == nil {
		t.Fatal("expected error when --keys not provided")
	}
}

func TestRunImmutable_NoIssues(t *testing.T) {
	dir := t.TempDir()
	f1 := writeImmutableEnv(t, dir, "dev.env", "APP_VERSION=1.0\nDB_HOST=localhost\n")
	f2 := writeImmutableEnv(t, dir, "prod.env", "APP_VERSION=1.0\nDB_HOST=prod.db\n")
	err := RunImmutable([]string{"--keys", "APP_VERSION", "dev=" + f1, "prod=" + f2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunImmutable_DetectsChange(t *testing.T) {
	dir := t.TempDir()
	f1 := writeImmutableEnv(t, dir, "dev.env", "APP_VERSION=1.0\nDB_HOST=localhost\n")
	f2 := writeImmutableEnv(t, dir, "prod.env", "APP_VERSION=1.0\nDB_HOST=prod.db\n")
	err := RunImmutable([]string{"--keys", "DB_HOST", "dev=" + f1, "prod=" + f2})
	if err == nil {
		t.Fatal("expected error due to immutable violation")
	}
}

func TestRunImmutable_InvalidFile(t *testing.T) {
	err := RunImmutable([]string{"--keys", "APP", "dev=/nonexistent/dev.env", "prod=/nonexistent/prod.env"})
	if err == nil {
		t.Fatal("expected error for invalid file")
	}
}

func TestRunImmutable_InvalidArgFormat(t *testing.T) {
	dir := t.TempDir()
	f1 := writeImmutableEnv(t, dir, "dev.env", "APP=1\n")
	err := RunImmutable([]string{"--keys", "APP", "dev=" + f1, "badarg"})
	if err == nil {
		t.Fatal("expected error for malformed name=file argument")
	}
}
