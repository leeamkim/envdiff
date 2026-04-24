package diff

import (
	"testing"
)

func baseOverlapEnvs() map[string]map[string]string {
	return map[string]map[string]string{
		"dev": {
			"APP_NAME": "myapp",
			"DB_HOST":  "localhost",
			"DEBUG":    "true",
		},
		"staging": {
			"APP_NAME": "myapp",
			"DB_HOST":  "staging.db",
			"LOG_LEVEL": "info",
		},
		"prod": {
			"APP_NAME": "myapp",
			"DB_HOST":  "prod.db",
			"LOG_LEVEL": "warn",
		},
	}
}

func TestComputeOverlap_AllEnvs(t *testing.T) {
	result := ComputeOverlap(baseOverlapEnvs(), 3)
	keys := map[string]bool{}
	for _, e := range result.Entries {
		keys[e.Key] = true
	}
	if !keys["APP_NAME"] {
		t.Error("expected APP_NAME in overlap")
	}
	if !keys["DB_HOST"] {
		t.Error("expected DB_HOST in overlap")
	}
	if keys["DEBUG"] {
		t.Error("DEBUG should not appear in 3-env overlap")
	}
}

func TestComputeOverlap_Consistent(t *testing.T) {
	result := ComputeOverlap(baseOverlapEnvs(), 3)
	for _, e := range result.Entries {
		if e.Key == "APP_NAME" && !e.Consistent {
			t.Error("APP_NAME should be consistent across all envs")
		}
		if e.Key == "DB_HOST" && e.Consistent {
			t.Error("DB_HOST should be inconsistent across envs")
		}
	}
}

func TestComputeOverlap_MinTwo(t *testing.T) {
	result := ComputeOverlap(baseOverlapEnvs(), 2)
	keys := map[string]bool{}
	for _, e := range result.Entries {
		keys[e.Key] = true
	}
	if !keys["LOG_LEVEL"] {
		t.Error("expected LOG_LEVEL in 2-env overlap")
	}
	if keys["DEBUG"] {
		t.Error("DEBUG only in one env, should not appear")
	}
}

func TestComputeOverlap_SortedEntries(t *testing.T) {
	result := ComputeOverlap(baseOverlapEnvs(), 1)
	for i := 1; i < len(result.Entries); i++ {
		if result.Entries[i].Key < result.Entries[i-1].Key {
			t.Errorf("entries not sorted: %s before %s", result.Entries[i-1].Key, result.Entries[i].Key)
		}
	}
}

func TestComputeOverlap_Empty(t *testing.T) {
	result := ComputeOverlap(map[string]map[string]string{}, 1)
	if len(result.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(result.Entries))
	}
}

func TestComputeOverlap_ValuesPopulated(t *testing.T) {
	result := ComputeOverlap(baseOverlapEnvs(), 3)
	for _, e := range result.Entries {
		if e.Key == "DB_HOST" {
			if e.Values["dev"] != "localhost" {
				t.Errorf("expected dev DB_HOST=localhost, got %s", e.Values["dev"])
			}
			if e.Values["prod"] != "prod.db" {
				t.Errorf("expected prod DB_HOST=prod.db, got %s", e.Values["prod"])
			}
		}
	}
}
