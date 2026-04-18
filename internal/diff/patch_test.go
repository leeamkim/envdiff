package diff

import (
	"strings"
	"testing"
)

func TestGeneratePatch_Add(t *testing.T) {
	a := map[string]string{}
	b := map[string]string{"FOO": "bar"}
	p := GeneratePatch(a, b)
	if len(p) != 1 || p[0].Action != "add" || p[0].Key != "FOO" {
		t.Errorf("unexpected patch: %+v", p)
	}
}

func TestGeneratePatch_Remove(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{}
	p := GeneratePatch(a, b)
	if len(p) != 1 || p[0].Action != "remove" || p[0].Key != "FOO" {
		t.Errorf("unexpected patch: %+v", p)
	}
}

func TestGeneratePatch_Update(t *testing.T) {
	a := map[string]string{"FOO": "old"}
	b := map[string]string{"FOO": "new"}
	p := GeneratePatch(a, b)
	if len(p) != 1 || p[0].Action != "update" || p[0].Value != "new" {
		t.Errorf("unexpected patch: %+v", p)
	}
}

func TestGeneratePatch_NoDiff(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "bar"}
	p := GeneratePatch(a, b)
	if len(p) != 0 {
		t.Errorf("expected empty patch, got %+v", p)
	}
}

func TestFormatPatch(t *testing.T) {
	a := map[string]string{"OLD": "val"}
	b := map[string]string{"NEW": "val", "OLD": "changed"}
	p := GeneratePatch(a, b)
	out := FormatPatch(p)
	if !strings.Contains(out, "+ NEW") {
		t.Errorf("expected add entry, got: %s", out)
	}
	if !strings.Contains(out, "~ OLD") {
		t.Errorf("expected update entry, got: %s", out)
	}
}
