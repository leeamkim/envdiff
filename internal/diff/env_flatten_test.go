package diff

import (
	"strings"
	"testing"
)

func TestFlattenEnvs_ReturnsAllEntries(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"A": "1", "B": "2"},
		"prod": {"A": "10"},
	}
	entries := FlattenEnvs(envs)
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
}

func TestFlattenEnvs_Empty(t *testing.T) {
	entries := FlattenEnvs(map[string]map[string]string{})
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestFormatFlatEntries_NoEntries(t *testing.T) {
	out := FormatFlatEntries(nil)
	if out != "(no entries)" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatFlatEntries_ContainsEnvAndKey(t *testing.T) {
	entries := []FlatEntry{
		{Env: "dev", Key: "PORT", Value: "8080"},
	}
	out := FormatFlatEntries(entries)
	if !strings.Contains(out, "[dev]") || !strings.Contains(out, "PORT=8080") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestGroupFlatByKey_Groups(t *testing.T) {
	entries := []FlatEntry{
		{Env: "dev", Key: "A", Value: "1"},
		{Env: "prod", Key: "A", Value: "2"},
		{Env: "dev", Key: "B", Value: "3"},
	}
	grouped := GroupFlatByKey(entries)
	if len(grouped["A"]) != 2 {
		t.Errorf("expected 2 entries for key A, got %d", len(grouped["A"]))
	}
	if len(grouped["B"]) != 1 {
		t.Errorf("expected 1 entry for key B, got %d", len(grouped["B"]))
	}
}

func TestGroupFlatByKey_Empty(t *testing.T) {
	grouped := GroupFlatByKey(nil)
	if len(grouped) != 0 {
		t.Errorf("expected empty map")
	}
}
