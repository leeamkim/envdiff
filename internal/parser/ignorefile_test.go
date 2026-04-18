package parser

import (
	"os"
	"testing"
)

func writeTempIgnore(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envignore-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestParseIgnoreFile_Basic(t *testing.T) {
	path := writeTempIgnore(t, "SECRET\nTOKEN\n")
	keys, err := ParseIgnoreFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 2 || keys[0] != "SECRET" || keys[1] != "TOKEN" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestParseIgnoreFile_SkipsCommentsAndBlanks(t *testing.T) {
	path := writeTempIgnore(t, "# this is a comment\n\nSECRET\n\n# another\nTOKEN\n")
	keys, err := ParseIgnoreFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}

func TestParseIgnoreFile_Empty(t *testing.T) {
	path := writeTempIgnore(t, "# only comments\n\n")
	keys, err := ParseIgnoreFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(keys))
	}
}

func TestParseIgnoreFile_NotFound(t *testing.T) {
	_, err := ParseIgnoreFile("/nonexistent/path/.envignore")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
