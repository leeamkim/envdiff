package cli

import (
	"os"
	"path/filepath"
	"testing"
)

// Integration: save snapshot from real env file, load and diff against modified version.
func TestRunSnapshot_Integration_ChangedKey(t *testing.T) {
	dir := t.TempDir()

	env1 := filepath.Join(dir, "v1.env")
	env2 := filepath.Join(dir, "v2.env")
	snap1 := filepath.Join(dir, "snap1.json")
	snap2 := filepath.Join(dir, "snap2.json")

	if err := os.WriteFile(env1, []byte("DB_HOST=localhost\nDB_PORT=5432\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(env2, []byte("DB_HOST=prod.db\nDB_PORT=5432\nDB_NAME=mydb\n"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := RunSnapshot([]string{"save", "v1", env1, snap1}); err != nil {
		t.Fatalf("save v1: %v", err)
	}
	if err := RunSnapshot([]string{"save", "v2", env2, snap2}); err != nil {
		t.Fatalf("save v2: %v", err)
	}
	if err := RunSnapshot([]string{"diff", snap1, snap2}); err != nil {
		t.Errorf("diff: %v", err)
	}
}

func TestRunSnapshot_InvalidEnvFile(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "snap.json")
	err := RunSnapshot([]string{"save", "label", "/nonexistent.env", out})
	if err == nil {
		t.Error("expected error for missing env file")
	}
}

func TestRunSnapshot_InvalidSnapFile(t *testing.T) {
	dir := t.TempDir()
	env := filepath.Join(dir, "a.env")
	if err := os.WriteFile(env, []byte("A=1\n"), 0644); err != nil {
		t.Fatal(err)
	}
	snap := filepath.Join(dir, "good.json")
	if err := RunSnapshot([]string{"save", "x", env, snap}); err != nil {
		t.Fatalf("save: %v", err)
	}
	err := RunSnapshot([]string{"diff", snap, "/nonexistent.json"})
	if err == nil {
		t.Error("expected error for missing snapshot file")
	}
}
