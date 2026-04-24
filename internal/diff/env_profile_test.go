package diff

import (
	"strings"
	"testing"
)

func TestBuildProfileReport_Keys(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"APP_NAME": "myapp", "DEBUG": "true"},
		"prod": {"APP_NAME": "myapp", "DEBUG": "false"},
	}
	r := BuildProfileReport(envs)
	if len(r.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(r.Entries))
	}
	if r.Entries[0].Key != "APP_NAME" {
		t.Errorf("expected APP_NAME first, got %s", r.Entries[0].Key)
	}
}

func TestBuildProfileReport_ProfilesSorted(t *testing.T) {
	envs := map[string]map[string]string{
		"staging": {"X": "1"},
		"dev":     {"X": "2"},
		"prod":    {"X": "3"},
	}
	r := BuildProfileReport(envs)
	if r.Profiles[0] != "dev" || r.Profiles[1] != "prod" || r.Profiles[2] != "staging" {
		t.Errorf("profiles not sorted: %v", r.Profiles)
	}
}

func TestBuildProfileReport_MissingKey(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"ONLY_DEV": "yes"},
		"prod": {},
	}
	r := BuildProfileReport(envs)
	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.Entries))
	}
	e := r.Entries[0]
	if e.Profiles["dev"] != "yes" {
		t.Errorf("expected dev=yes, got %s", e.Profiles["dev"])
	}
	if e.Profiles["prod"] != "" {
		t.Errorf("expected prod empty, got %s", e.Profiles["prod"])
	}
}

func TestBuildProfileReport_Empty(t *testing.T) {
	r := BuildProfileReport(map[string]map[string]string{})
	if len(r.Entries) != 0 {
		t.Errorf("expected no entries")
	}
}

func TestFormatProfileReport_ContainsHeaders(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"KEY": "val"},
		"prod": {"KEY": "other"},
	}
	r := BuildProfileReport(envs)
	out := FormatProfileReport(r)
	if !strings.Contains(out, "KEY") {
		t.Error("expected KEY in output")
	}
	if !strings.Contains(out, "dev") {
		t.Error("expected dev in output")
	}
	if !strings.Contains(out, "prod") {
		t.Error("expected prod in output")
	}
}

func TestFormatProfileReport_ShowsMissing(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"ONLY_DEV": "yes"},
		"prod": {},
	}
	r := BuildProfileReport(envs)
	out := FormatProfileReport(r)
	if !strings.Contains(out, "(missing)") {
		t.Error("expected (missing) in output")
	}
}

func TestFormatProfileReport_Empty(t *testing.T) {
	r := ProfileReport{}
	out := FormatProfileReport(r)
	if !strings.Contains(out, "no keys") {
		t.Error("expected 'no keys' message")
	}
}
