package diff

import (
	"testing"
)

func TestCompareMulti_NoDiff(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"A": "1", "B": "2"},
		"prod": {"A": "1", "B": "2"},
	}
	res := CompareMulti(envs)
	if res.HasAnyDiff() {
		t.Error("expected no diff")
	}
	if len(res.Pairs) != 1 {
		t.Errorf("expected 1 pair, got %d", len(res.Pairs))
	}
}

func TestCompareMulti_WithDiff(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"A": "1"},
		"prod": {"A": "2"},
	}
	res := CompareMulti(envs)
	if !res.HasAnyDiff() {
		t.Error("expected diff")
	}
}

func TestCompareMulti_ThreeEnvs(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":     {"A": "1"},
		"staging": {"A": "1"},
		"prod":    {"A": "1"},
	}
	res := CompareMulti(envs)
	if len(res.Pairs) != 3 {
		t.Errorf("expected 3 pairs, got %d", len(res.Pairs))
	}
	if res.HasAnyDiff() {
		t.Error("expected no diff")
	}
}

func TestCompareMulti_PairKey(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"X": "1"},
		"prod": {"X": "2"},
	}
	res := CompareMulti(envs)
	key := PairKey("dev", "prod")
	if _, ok := res.Pairs[key]; !ok {
		t.Errorf("expected pair key %q", key)
	}
}

func TestCompareMulti_EnvNames(t *testing.T) {
	envs := map[string]map[string]string{
		"b": {},
		"a": {},
	}
	res := CompareMulti(envs)
	if res.EnvNames[0] != "a" || res.EnvNames[1] != "b" {
		t.Errorf("expected sorted names, got %v", res.EnvNames)
	}
}
