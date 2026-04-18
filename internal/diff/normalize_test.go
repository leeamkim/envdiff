package diff

import (
	"testing"
)

func TestNormalizeMap_TrimSpace(t *testing.T) {
	input := map[string]string{
		"KEY1": "  hello  ",
		"KEY2": "world",
	}
	out := NormalizeMap(input, NormalizeTrimSpace)
	if out["KEY1"] != "hello" {
		t.Errorf("expected 'hello', got %q", out["KEY1"])
	}
	if out["KEY2"] != "world" {
		t.Errorf("expected 'world', got %q", out["KEY2"])
	}
}

func TestNormalizeMap_Lowercase(t *testing.T) {
	input := map[string]string{"KEY": "Hello World"}
	out := NormalizeMap(input, NormalizeLowercase)
	if out["KEY"] != "hello world" {
		t.Errorf("expected 'hello world', got %q", out["KEY"])
	}
}

func TestNormalizeMap_MultipleOpts(t *testing.T) {
	input := map[string]string{"KEY": "  UPPER  "}
	out := NormalizeMap(input, NormalizeTrimSpace, NormalizeLowercase)
	if out["KEY"] != "upper" {
		t.Errorf("expected 'upper', got %q", out["KEY"])
	}
}

func TestNormalizeMap_DoesNotMutateOriginal(t *testing.T) {
	input := map[string]string{"KEY": "  value  "}
	NormalizeMap(input, NormalizeTrimSpace)
	if input["KEY"] != "  value  " {
		t.Error("original map was mutated")
	}
}

func TestNormalizeCompare_EqualAfterNormalize(t *testing.T) {
	a := map[string]string{"KEY": "  hello  "}
	b := map[string]string{"KEY": "hello"}
	result := NormalizeCompare(a, b, NormalizeTrimSpace)
	if len(result.Mismatched) != 0 {
		t.Errorf("expected no mismatches after normalize, got %v", result.Mismatched)
	}
}

func TestNormalizeCompare_StillMismatchedWithoutNormalize(t *testing.T) {
	a := map[string]string{"KEY": "  hello  "}
	b := map[string]string{"KEY": "hello"}
	result := Compare(a, b)
	if len(result.Mismatched) == 0 {
		t.Error("expected mismatch without normalization")
	}
}
