package diff

import (
	"testing"
)

func TestBuildEnvDiffMap_Keys(t *testing.T) {
	envs := map[string]map[string]string{
		"prod": {"A": "1", "B": "2"},
		"dev":  {"A": "1", "C": "3"},
	}
	m := BuildEnvDiffMap(envs)
	if len(m.Keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(m.Keys))
	}
}

func TestBuildEnvDiffMap_MissingValue(t *testing.T) {
	envs := map[string]map[string]string{
		"prod": {"A": "1"},
		"dev":  {"A": "1", "B": "2"},
	}
	m := BuildEnvDiffMap(envs)
	if m.Cells["B"]["prod"] != "" {
		t.Errorf("expected empty string for missing key in prod")
	}
}

func TestBuildEnvDiffMap_Consistent(t *testing.T) {
	envs := map[string]map[string]string{
		"prod": {"A": "same"},
		"dev":  {"A": "same"},
	}
	m := BuildEnvDiffMap(envs)
	if !m.Consistent("A") {
		t.Error("expected A to be consistent")
	}
}

func TestBuildEnvDiffMap_Inconsistent(t *testing.T) {
	envs := map[string]map[string]string{
		"prod": {"A": "x"},
		"dev":  {"A": "y"},
	}
	m := BuildEnvDiffMap(envs)
	if m.Consistent("A") {
		t.Error("expected A to be inconsistent")
	}
}

func TestBuildEnvDiffMap_EnvNames(t *testing.T) {
	envs := map[string]map[string]string{
		"staging": {"X": "1"},
		"prod":    {"X": "1"},
	}
	m := BuildEnvDiffMap(envs)
	if len(m.Envs) != 2 {
		t.Fatalf("expected 2 envs, got %d", len(m.Envs))
	}
	if m.Envs[0] != "prod" || m.Envs[1] != "staging" {
		t.Errorf("unexpected env order: %v", m.Envs)
	}
}
