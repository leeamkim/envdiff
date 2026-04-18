package diff

import (
	"testing"
)

func TestApplyRenames_Basic(t *testing.T) {
	env := map[string]string{
		"OLD_KEY": "value1",
		"KEEP":    "value2",
	}
	renames := RenameMap{"OLD_KEY": "NEW_KEY"}
	updated, applied := ApplyRenames(env, renames)

	if updated["NEW_KEY"] != "value1" {
		t.Errorf("expected NEW_KEY=value1, got %q", updated["NEW_KEY"])
	}
	if updated["KEEP"] != "value2" {
		t.Errorf("expected KEEP=value2, got %q", updated["KEEP"])
	}
	if _, ok := updated["OLD_KEY"]; ok {
		t.Error("expected OLD_KEY to be removed")
	}
	if len(applied) != 1 || applied[0].OldKey != "OLD_KEY" || applied[0].NewKey != "NEW_KEY" {
		t.Errorf("unexpected applied renames: %+v", applied)
	}
}

func TestApplyRenames_MissingKey(t *testing.T) {
	env := map[string]string{"KEEP": "val"}
	renames := RenameMap{"GHOST": "NEW"}
	updated, applied := ApplyRenames(env, renames)

	if len(applied) != 0 {
		t.Errorf("expected no applied renames, got %+v", applied)
	}
	if updated["KEEP"] != "val" {
		t.Error("expected KEEP to remain")
	}
}

func TestApplyRenames_Empty(t *testing.T) {
	updated, applied := ApplyRenames(map[string]string{}, RenameMap{})
	if len(updated) != 0 || len(applied) != 0 {
		t.Error("expected empty results")
	}
}

func TestParseRenameMap_Valid(t *testing.T) {
	rm, err := ParseRenameMap([]string{"OLD=NEW", "FOO=BAR"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rm["OLD"] != "NEW" || rm["FOO"] != "BAR" {
		t.Errorf("unexpected rename map: %v", rm)
	}
}

func TestParseRenameMap_Invalid(t *testing.T) {
	_, err := ParseRenameMap([]string{"NODIVIDER"})
	if err == nil {
		t.Error("expected error for invalid pair")
	}
}

func TestParseRenameMap_EmptyParts(t *testing.T) {
	_, err := ParseRenameMap([]string{"=NEW"})
	if err == nil {
		t.Error("expected error for empty old key")
	}
}
