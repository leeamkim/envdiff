package diff

import (
	"testing"
)

func TestTagEnv_NoRules(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	entries := TagEnv(env)
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestTagEnv_SecretRule(t *testing.T) {
	env := map[string]string{"DB_PASSWORD": "secret123", "APP_NAME": "myapp"}
	entries := TagEnv(env, TagRuleSecrets)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Key != "DB_PASSWORD" {
		t.Errorf("expected DB_PASSWORD, got %s", entries[0].Key)
	}
	if entries[0].Tags[0] != "secret" {
		t.Errorf("expected tag 'secret', got %s", entries[0].Tags[0])
	}
}

func TestTagEnv_EmptyRule(t *testing.T) {
	env := map[string]string{"FOO": "", "BAR": "value"}
	entries := TagEnv(env, TagRuleEmpty)
	if len(entries) != 1 || entries[0].Key != "FOO" {
		t.Errorf("expected FOO tagged as empty")
	}
}

func TestTagEnv_PlaceholderRule(t *testing.T) {
	env := map[string]string{"API_KEY": "CHANGEME", "HOST": "localhost"}
	entries := TagEnv(env, TagRulePlaceholder)
	if len(entries) != 1 || entries[0].Key != "API_KEY" {
		t.Errorf("expected API_KEY tagged as placeholder")
	}
}

func TestTagEnv_MultipleRules(t *testing.T) {
	env := map[string]string{"DB_PASSWORD": "", "FOO": "bar"}
	entries := TagEnv(env, TagRuleSecrets, TagRuleEmpty)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if len(entries[0].Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(entries[0].Tags))
	}
}

func TestFormatTagEntries_Empty(t *testing.T) {
	out := FormatTagEntries(nil)
	if out != "no tags" {
		t.Errorf("expected 'no tags', got %s", out)
	}
}

func TestFormatTagEntries_WithEntries(t *testing.T) {
	entries := []TagEntry{{Key: "FOO", Tags: []string{"empty"}}}
	out := FormatTagEntries(entries)
	if out == "" || out == "no tags" {
		t.Errorf("expected formatted output, got: %s", out)
	}
}
