package diff

import (
	"testing"
)

func TestResolveAliases_NoAliasNeeded(t *testing.T) {
	env := map[string]string{"DATABASE_URL": "postgres://localhost"}
	aliases := AliasMap{"DATABASE_URL": {"DB_URL", "DB_CONNECTION"}}
	out, issues := ResolveAliases(env, aliases)
	if out["DATABASE_URL"] != "postgres://localhost" {
		t.Errorf("expected canonical key preserved")
	}
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d", len(issues))
	}
}

func TestResolveAliases_SubstitutesAlias(t *testing.T) {
	env := map[string]string{"DB_URL": "postgres://localhost"}
	aliases := AliasMap{"DATABASE_URL": {"DB_URL"}}
	out, issues := ResolveAliases(env, aliases)
	if out["DATABASE_URL"] != "postgres://localhost" {
		t.Errorf("expected canonical key set")
	}
	if _, ok := out["DB_URL"]; ok {
		t.Errorf("alias key should be removed")
	}
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Canonical != "DATABASE_URL" || issues[0].FoundAs != "DB_URL" {
		t.Errorf("unexpected issue: %v", issues[0])
	}
}

func TestResolveAliases_FirstAliasWins(t *testing.T) {
	env := map[string]string{"DB_URL": "first", "DB_CONNECTION": "second"}
	aliases := AliasMap{"DATABASE_URL": {"DB_URL", "DB_CONNECTION"}}
	out, issues := ResolveAliases(env, aliases)
	if out["DATABASE_URL"] != "first" {
		t.Errorf("expected first alias to win")
	}
	if len(issues) != 1 {
		t.Errorf("expected 1 issue")
	}
}

func TestAliasIssue_String(t *testing.T) {
	i := AliasIssue{Canonical: "FOO", FoundAs: "FOO_ALIAS", Value: "bar"}
	s := i.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}

func TestParseAliasMap_Valid(t *testing.T) {
	am, err := ParseAliasMap([]string{"DATABASE_URL=DB_URL,DB_CONNECTION"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(am["DATABASE_URL"]) != 2 {
		t.Errorf("expected 2 aliases, got %d", len(am["DATABASE_URL"]))
	}
}

func TestParseAliasMap_Invalid(t *testing.T) {
	_, err := ParseAliasMap([]string{"NODIVIDER"})
	if err == nil {
		t.Error("expected error for invalid entry")
	}
}
