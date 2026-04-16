package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestRun_NoDiff(t *testing.T) {
	a := writeTempEnv(t, "KEY=value\nFOO=bar\n")
	b := writeTempEnv(t, "KEY=value\nFOO=bar\n")

	if err := Run([]string{a, b}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRun_StrictWithDiff(t *testing.T) {
	a := writeTempEnv(t, "KEY=value\n")
	b := writeTempEnv(t, "KEY=other\n")

	err := Run([]string{"--strict", a, b})
	if err == nil {
		t.Fatal("expected error in strict mode with diff")
	}
}

func TestRun_MissingArgs(t *testing.T) {
	if err := Run([]string{}); err == nil {
		t.Fatal("expected error for missing arguments")
	}
}

func TestRun_InvalidFile(t *testing.T) {
	err := Run([]string{"/nonexistent/a.env", "/nonexistent/b.env"})
	if err == nil {
		t.Fatal("expected error for invalid file path")
	}
}

func TestRun_OutputWritten(t *testing.T) {
	a := writeTempEnv(t, "KEY=value\nMISSING=yes\n")
	b := writeTempEnv(t, "KEY=changed\n")

	var buf bytes.Buffer
	cfg := &Config{
		FileA:  a,
		FileB:  b,
		Strict: false,
		Output: &buf,
	}
	if err := run(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected output to be written")
	}
}

func TestParseArgs_Strict(t *testing.T) {
	cfg, err := parseArgs([]string{"--strict", "a.env", "b.env"})
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.Strict {
		t.Error("expected strict=true")
	}
	if cfg.FileA != "a.env" || cfg.FileB != filepath.Clean("b.env") {
		t.Errorf("unexpected files: %s %s", cfg.FileA, cfg.FileB)
	}
}
