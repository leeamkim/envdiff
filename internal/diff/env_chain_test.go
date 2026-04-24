package diff

import (
	"strings"
	"testing"
)

func TestBuildChain_NoDiff(t *testing.T) {
	chain := []string{"dev", "staging", "prod"}
	envs := map[string]map[string]string{
		"dev":     {"HOST": "localhost"},
		"staging": {"HOST": "localhost"},
		"prod":    {"HOST": "localhost"},
	}
	r := BuildChain(chain, envs)
	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.Entries))
	}
	for _, cv := range r.Entries[0].Values {
		if cv.Override {
			t.Errorf("expected no override for %s", cv.Env)
		}
	}
}

func TestBuildChain_Override(t *testing.T) {
	chain := []string{"dev", "prod"}
	envs := map[string]map[string]string{
		"dev":  {"DB": "dev-db"},
		"prod": {"DB": "prod-db"},
	}
	r := BuildChain(chain, envs)
	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.Entries))
	}
	prodVal := r.Entries[0].Values[1]
	if !prodVal.Override {
		t.Error("expected prod to be marked as override")
	}
}

func TestBuildChain_MissingInSome(t *testing.T) {
	chain := []string{"dev", "staging", "prod"}
	envs := map[string]map[string]string{
		"dev":     {"SECRET": "abc"},
		"staging": {},
		"prod":    {"SECRET": "abc"},
	}
	r := BuildChain(chain, envs)
	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.Entries))
	}
	staging := r.Entries[0].Values[1]
	if staging.Present {
		t.Error("expected staging to be missing")
	}
}

func TestBuildChain_KeysSorted(t *testing.T) {
	chain := []string{"dev"}
	envs := map[string]map[string]string{
		"dev": {"ZEBRA": "1", "ALPHA": "2", "MANGO": "3"},
	}
	r := BuildChain(chain, envs)
	keys := make([]string, len(r.Entries))
	for i, e := range r.Entries {
		keys[i] = e.Key
	}
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("keys not sorted: %v", keys)
		}
	}
}

func TestFormatChainReport_Empty(t *testing.T) {
	r := ChainResult{}
	out := FormatChainReport(r)
	if !strings.Contains(out, "no keys") {
		t.Errorf("expected 'no keys' message, got: %s", out)
	}
}

func TestFormatChainReport_ContainsOverride(t *testing.T) {
	chain := []string{"dev", "prod"}
	envs := map[string]map[string]string{
		"dev":  {"PORT": "3000"},
		"prod": {"PORT": "8080"},
	}
	r := BuildChain(chain, envs)
	out := FormatChainReport(r)
	if !strings.Contains(out, "override") {
		t.Errorf("expected 'override' in output, got: %s", out)
	}
}
