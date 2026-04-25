package diff

import (
	"strings"
	"testing"
)

func TestFormatEnv_PlainStyle(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out := FormatEnv(env, FormatOptions{Style: FormatStylePlain, SortKeys: true})
	if !strings.Contains(out, "FOO=bar") {
		t.Errorf("expected FOO=bar in output, got: %s", out)
	}
	if !strings.Contains(out, "BAZ=qux") {
		t.Errorf("expected BAZ=qux in output, got: %s", out)
	}
}

func TestFormatEnv_ExportStyle(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	out := FormatEnv(env, FormatOptions{Style: FormatStyleExport, SortKeys: true})
	if !strings.Contains(out, "export KEY=value") {
		t.Errorf("expected export KEY=value, got: %s", out)
	}
}

func TestFormatEnv_QuotedStyle(t *testing.T) {
	env := map[string]string{"MSG": "hello world"}
	out := FormatEnv(env, FormatOptions{Style: FormatStyleQuoted, SortKeys: true})
	if !strings.Contains(out, `MSG="hello world"`) {
		t.Errorf("expected quoted value, got: %s", out)
	}
}

func TestFormatEnv_WithPrefix(t *testing.T) {
	env := map[string]string{"NAME": "alice"}
	out := FormatEnv(env, FormatOptions{Style: FormatStylePlain, SortKeys: true, Prefix: "APP_"})
	if !strings.Contains(out, "APP_NAME=alice") {
		t.Errorf("expected APP_NAME=alice, got: %s", out)
	}
}

func TestFormatEnv_SortedKeys(t *testing.T) {
	env := map[string]string{"Z": "last", "A": "first", "M": "mid"}
	out := FormatEnv(env, FormatOptions{Style: FormatStylePlain, SortKeys: true})
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "A=") {
		t.Errorf("expected first line to start with A=, got: %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "Z=") {
		t.Errorf("expected last line to start with Z=, got: %s", lines[2])
	}
}

func TestFormatEnv_EmptyValue_Export(t *testing.T) {
	env := map[string]string{"EMPTY": ""}
	out := FormatEnv(env, FormatOptions{Style: FormatStyleExport, SortKeys: true})
	if !strings.Contains(out, `export EMPTY=""`) {
		t.Errorf("expected export EMPTY=\"\", got: %s", out)
	}
}

func TestShellQuote_SpecialChars(t *testing.T) {
	tests := []struct {
		input    string
		contains string
	}{
		{"simple", "simple"},
		{"has space", "'"},
		{"has$dollar", "'"},
		{"", `""`},
	}
	for _, tt := range tests {
		out := shellQuote(tt.input)
		if !strings.Contains(out, tt.contains) {
			t.Errorf("shellQuote(%q) = %q, want to contain %q", tt.input, out, tt.contains)
		}
	}
}
