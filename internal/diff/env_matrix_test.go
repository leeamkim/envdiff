package diff

import (
	"testing"
)

func TestBuildMatrix_Keys(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"A": "1", "B": "2"},
		"prod": {"A": "1", "C": "3"},
	}
	m := BuildMatrix(envs)
	if len(m.Keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(m.Keys))
	}
}

func TestBuildMatrix_MatchStatus(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"A": "1"},
		"prod": {"A": "1"},
	}
	m := BuildMatrix(envs)
	pair := m.EnvNames[0] + ":" + m.EnvNames[1]
	cell := m.Status["A"][pair]
	if cell.Status != "match" {
		t.Errorf("expected match, got %s", cell.Status)
	}
}

func TestBuildMatrix_MismatchStatus(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"A": "1"},
		"prod": {"A": "2"},
	}
	m := BuildMatrix(envs)
	pair := m.EnvNames[0] + ":" + m.EnvNames[1]
	cell := m.Status["A"][pair]
	if cell.Status != "mismatch" {
		t.Errorf("expected mismatch, got %s", cell.Status)
	}
}

func TestBuildMatrix_MissingStatus(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"A": "1"},
		"prod": {},
	}
	m := BuildMatrix(envs)
	pair := m.EnvNames[0] + ":" + m.EnvNames[1]
	cell := m.Status["A"][pair]
	if cell.Status != "missing" {
		t.Errorf("expected missing, got %s", cell.Status)
	}
}

func TestBuildMatrix_SingleEnv(t *testing.T) {
	envs := map[string]map[string]string{
		"dev": {"A": "1"},
	}
	m := BuildMatrix(envs)
	if len(m.Status) != 0 {
		t.Errorf("expected no status entries for single env")
	}
}

func TestBuildMatrix_CellValues(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"X": "hello"},
		"prod": {"X": "world"},
	}
	m := BuildMatrix(envs)
	if m.Cells["X"]["dev"] != "hello" {
		t.Errorf("expected hello")
	}
	if m.Cells["X"]["prod"] != "world" {
		t.Errorf("expected world")
	}
}
