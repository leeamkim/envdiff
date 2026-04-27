package diff

import (
	"strings"
	"testing"
)

func baseSensitiveEnvs() map[string]map[string]string {
	return map[string]map[string]string{
		"prod": {
			"DB_PASSWORD": "secret123",
			"API_KEY":     "abc",
			"APP_NAME":    "myapp",
		},
		"staging": {
			"DB_PASSWORD": "staging-pass",
			"APP_NAME":    "myapp-staging",
		},
	}
}

func TestCountSensitiveKeys_Basic(t *testing.T) {
	result := CountSensitiveKeys(baseSensitiveEnvs(), nil)
	if len(result.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result.Entries))
	}
	// entries should be sorted by name: prod, staging
	if result.Entries[0].EnvName != "prod" {
		t.Errorf("expected first entry to be prod, got %s", result.Entries[0].EnvName)
	}
	if result.Entries[0].SensitiveKeys != 2 {
		t.Errorf("expected 2 sensitive keys for prod, got %d", result.Entries[0].SensitiveKeys)
	}
	if result.Entries[1].SensitiveKeys != 1 {
		t.Errorf("expected 1 sensitive key for staging, got %d", result.Entries[1].SensitiveKeys)
	}
}

func TestCountSensitiveKeys_TotalKeys(t *testing.T) {
	result := CountSensitiveKeys(baseSensitiveEnvs(), nil)
	prod := result.Entries[0]
	if prod.TotalKeys != 3 {
		t.Errorf("expected 3 total keys for prod, got %d", prod.TotalKeys)
	}
}

func TestCountSensitiveKeys_RedactedKeys(t *testing.T) {
	result := CountSensitiveKeys(baseSensitiveEnvs(), []string{"DB_PASSWORD"})
	prod := result.Entries[0]
	if prod.RedactedKeys != 1 {
		t.Errorf("expected 1 redacted key for prod, got %d", prod.RedactedKeys)
	}
	staging := result.Entries[1]
	if staging.RedactedKeys != 1 {
		t.Errorf("expected 1 redacted key for staging, got %d", staging.RedactedKeys)
	}
}

func TestCountSensitiveKeys_Empty(t *testing.T) {
	result := CountSensitiveKeys(map[string]map[string]string{}, nil)
	if len(result.Entries) != 0 {
		t.Errorf("expected 0 entries for empty input")
	}
}

func TestFormatSensitiveCountReport_ContainsHeaders(t *testing.T) {
	result := CountSensitiveKeys(baseSensitiveEnvs(), nil)
	out := FormatSensitiveCountReport(result)
	for _, header := range []string{"ENV", "TOTAL", "SENSITIVE", "REDACTED"} {
		if !strings.Contains(out, header) {
			t.Errorf("expected output to contain %q", header)
		}
	}
}

func TestFormatSensitiveCountReport_Empty(t *testing.T) {
	out := FormatSensitiveCountReport(SensitiveCountResult{})
	if !strings.Contains(out, "no environments") {
		t.Errorf("expected empty message, got: %s", out)
	}
}
