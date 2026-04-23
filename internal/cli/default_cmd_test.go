package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func writeDefaultEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunDefault_MissingArgs(t *testing.T) {
	err := RunDefault([]string{})
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRunDefault_NoIssues(t *testing.T) {
	f := writeDefaultEnv(t, "DB_HOST=prod.db.example.com\nAPI_KEY=abc123\n")
	err := RunDefault([]string{f})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunDefault_DetectsChangeme(t *testing.T) {
	f := writeDefaultEnv(t, "DB_PASS=changeme\nAPI_KEY=real-key\n")
	err := RunDefault([]string{f})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunDefault_DetectsTodo(t *testing.T) {
	f := writeDefaultEnv(t, "SECRET=todo\n")
	err := RunDefault([]string{f})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunDefault_WithPins(t *testing.T) {
	f := writeDefaultEnv(t, "REGION=us-east-1\nTIMEOUT=60\n")
	err := RunDefault([]string{"--pins", "REGION=us-east-1", f})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunDefault_InvalidPins(t *testing.T) {
	f := writeDefaultEnv(t, "REGION=us-east-1\n")
	err := RunDefault([]string{"--pins", "BADPIN", f})
	if err == nil {
		t.Fatal("expected error for invalid pin format")
	}
}

func TestRunDefault_InvalidFile(t *testing.T) {
	err := RunDefault([]string{"/nonexistent/.env"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRunDefault_MultipleFiles(t *testing.T) {
	f1 := writeDefaultEnv(t, "KEY=real\n")
	f2 := writeDefaultEnv(t, "KEY=example\n")
	err := RunDefault([]string{f1, f2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
