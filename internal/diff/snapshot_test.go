package diff

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewSnapshot_CopiesEntries(t *testing.T) {
	m := map[string]string{"A": "1", "B": "2"}
	s := NewSnapshot("test", m)
	m["A"] = "mutated"
	if s.Entries["A"] != "1" {
		t.Errorf("expected snapshot to be isolated from source map")
	}
	if s.Label != "test" {
		t.Errorf("expected label 'test'")
	}
}

func TestSaveLoadSnapshot(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	s := NewSnapshot("env-a", map[string]string{"KEY": "val"})
	if err := SaveSnapshot(path, s); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded.Label != "env-a" {
		t.Errorf("expected label env-a, got %s", loaded.Label)
	}
	if loaded.Entries["KEY"] != "val" {
		t.Errorf("expected KEY=val")
	}
}

func TestLoadSnapshot_NotFound(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestDiffSnapshots_Added(t *testing.T) {
	a := NewSnapshot("a", map[string]string{"X": "1"})
	b := NewSnapshot("b", map[string]string{"X": "1", "Y": "2"})
	d := DiffSnapshots(a, b)
	if _, ok := d.Added["Y"]; !ok {
		t.Error("expected Y to be added")
	}
	if d.HasChanges() == false {
		t.Error("expected HasChanges true")
	}
}

func TestDiffSnapshots_Removed(t *testing.T) {
	a := NewSnapshot("a", map[string]string{"X": "1", "Z": "3"})
	b := NewSnapshot("b", map[string]string{"X": "1"})
	d := DiffSnapshots(a, b)
	if _, ok := d.Removed["Z"]; !ok {
		t.Error("expected Z to be removed")
	}
}

func TestDiffSnapshots_Changed(t *testing.T) {
	a := NewSnapshot("a", map[string]string{"X": "old"})
	b := NewSnapshot("b", map[string]string{"X": "new"})
	d := DiffSnapshots(a, b)
	if pair, ok := d.Changed["X"]; !ok || pair[0] != "old" || pair[1] != "new" {
		t.Errorf("expected X changed old->new, got %v", d.Changed)
	}
}

func TestDiffSnapshots_NoChanges(t *testing.T) {
	a := NewSnapshot("a", map[string]string{"X": "1"})
	b := NewSnapshot("b", map[string]string{"X": "1"})
	d := DiffSnapshots(a, b)
	if d.HasChanges() {
		t.Error("expected no changes")
	}
}

var _ = os.WriteFile // suppress unused import
