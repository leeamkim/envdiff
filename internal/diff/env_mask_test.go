package diff

import (
	"strings"
	"testing"
)

func TestMaskEnv_NonSensitiveUnchanged(t *testing.T) {
	env := map[string]string{"APP_NAME": "myapp", "PORT": "8080"}
	entries := MaskEnv(env, MaskOptions{})
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	for _, e := range entries {
		if e.Original != e.Masked {
			t.Errorf("expected non-sensitive key %s to be unmasked", e.Key)
		}
	}
}

func TestMaskEnv_SensitiveKeyMasked(t *testing.T) {
	env := map[string]string{"DB_PASSWORD": "secret123"}
	entries := MaskEnv(env, MaskOptions{})
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry")
	}
	if entries[0].Masked == entries[0].Original {
		t.Error("expected sensitive key to be masked")
	}
	if !strings.Contains(entries[0].Masked, "****") {
		t.Errorf("expected masked value to contain ****, got %q", entries[0].Masked)
	}
}

func TestMaskEnv_ShowLength(t *testing.T) {
	env := map[string]string{"API_SECRET": "abcdefgh"}
	entries := MaskEnv(env, MaskOptions{ShowLength: true})
	if !strings.Contains(entries[0].Masked, "(8)") {
		t.Errorf("expected length hint in masked value, got %q", entries[0].Masked)
	}
}

func TestMaskEnv_VisibleChars(t *testing.T) {
	env := map[string]string{"API_KEY": "abcdefgh"}
	entries := MaskEnv(env, MaskOptions{VisibleChars: 3})
	if !strings.HasPrefix(entries[0].Masked, "abc") {
		t.Errorf("expected first 3 chars visible, got %q", entries[0].Masked)
	}
	if !strings.Contains(entries[0].Masked, "****") {
		t.Errorf("expected stars after visible chars, got %q", entries[0].Masked)
	}
}

func TestMaskEnv_EmptyValue(t *testing.T) {
	env := map[string]string{"DB_PASSWORD": ""}
	entries := MaskEnv(env, MaskOptions{})
	if entries[0].Masked != "" {
		t.Errorf("expected empty masked value for empty original, got %q", entries[0].Masked)
	}
}

func TestMaskEnv_SortedOutput(t *testing.T) {
	env := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	entries := MaskEnv(env, MaskOptions{})
	if entries[0].Key != "A_KEY" || entries[1].Key != "M_KEY" || entries[2].Key != "Z_KEY" {
		t.Error("expected entries sorted by key")
	}
}

func TestFormatMaskEntries_Empty(t *testing.T) {
	out := FormatMaskEntries(nil)
	if !strings.Contains(out, "no entries") {
		t.Errorf("expected 'no entries', got %q", out)
	}
}

func TestFormatMaskEntries_MaskedLabel(t *testing.T) {
	env := map[string]string{"API_TOKEN": "supersecret"}
	entries := MaskEnv(env, MaskOptions{})
	out := FormatMaskEntries(entries)
	if !strings.Contains(out, "[masked]") {
		t.Errorf("expected [masked] label in output, got %q", out)
	}
}
