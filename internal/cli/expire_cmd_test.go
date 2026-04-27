package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func writeExpireEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunExpire_MissingArgs(t *testing.T) {
	err := RunExpire([]string{})
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRunExpire_InvalidFile(t *testing.T) {
	err := RunExpire([]string{"/nonexistent/.env", "API_KEY=2024-01-01"})
	if err == nil {
		t.Fatal("expected error for invalid file")
	}
}

func TestRunExpire_InvalidRuleFormat(t *testing.T) {
	p := writeExpireEnv(t, "API_KEY=secret\n")
	err := RunExpire([]string{p, "BADFORMAT"})
	if err == nil {
		t.Fatal("expected error for bad rule format")
	}
}

func TestRunExpire_InvalidDate(t *testing.T) {
	p := writeExpireEnv(t, "API_KEY=secret\n")
	err := RunExpire([]string{p, "API_KEY=not-a-date"})
	if err == nil {
		t.Fatal("expected error for invalid date")
	}
}

func TestRunExpire_NoIssues(t *testing.T) {
	p := writeExpireEnv(t, "API_KEY=secret\n")
	err := RunExpire([]string{p, "API_KEY=2099-01-01"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunExpire_DetectsExpired(t *testing.T) {
	p := writeExpireEnv(t, "API_KEY=secret\n")
	err := RunExpire([]string{p, "API_KEY=2020-01-01"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunExpire_WarnDaysFlag(t *testing.T) {
	p := writeExpireEnv(t, "API_KEY=secret\n")
	err := RunExpire([]string{p, "API_KEY=2099-12-31", "--warn-days=99999"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
