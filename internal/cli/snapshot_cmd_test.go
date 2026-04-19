package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func writeSnapshotEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	return p
}

func TestRunSnapshot_MissingArgs(t *testing.T) {
	if err := RunSnapshot([]string{}); err == nil {
		t.Error("expected error for missing args")
	}
}

func TestRunSnapshot_UnknownCommand(t *testing.T) {
	if err := RunSnapshot([]string{"export"}); err == nil {
		t.Error("expected error for unknown command")
	}
}

func TestRunSnapshot_SaveAndDiff_NoChanges(t *testing.T) {
	dir := t.TempDir()
	env := writeSnapshotEnv(t, dir, "a.env", "KEY=value\nFOO=bar\n")
	snap1 := filepath.Join(dir, "snap1.json")
	snap2 := filepath.Join(dir, "snap2.json")

	if err := RunSnapshot([]string{"save", "v1", env, snap1}); err != nil {
		t.Fatalf("save snap1: %v", err)
	}
	if err := RunSnapshot([]string{"save", "v2", env, snap2}); err != nil {
		t.Fatalf("save snap2: %v", err)
	}
	if err := RunSnapshot([]string{"diff", snap1, snap2}); err != nil {
		t.Errorf("diff: %v", err)
	}
}

func TestRunSnapshot_DiffDetectsChanges(t *testing.T) {
	dir := t.TempDir()
	snap1 := filepath.Join(dir, "s1.json")
	snap2 := filepath.Join(dir, "s2.json")

	s1 := diff.NewSnapshot("old", map[string]string{"A": "1", "B": "2"})
	s2 := diff.NewSnapshot("new", map[string]string{"A": "changed", "C": "3"})

	if err := diff.SaveSnapshot(snap1, s1); err != nil {
		t.Fatalf("save s1: %v", err)
	}
	if err := diff.SaveSnapshot(snap2, s2); err != nil {
		t.Fatalf("save s2: %v", err)
	}
	if err := RunSnapshot([]string{"diff", snap1, snap2}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunSnapshot_SaveMissingArgs(t *testing.T) {
	if err := RunSnapshot([]string{"save", "label"}); err == nil {
		t.Error("expected error")
	}
}

func TestRunSnapshot_DiffMissingArgs(t *testing.T) {
	if err := RunSnapshot([]string{"diff", "only-one"}); err == nil {
		t.Error("expected error")
	}
}
