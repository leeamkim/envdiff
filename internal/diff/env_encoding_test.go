package diff

import (
	"strings"
	"testing"
)

func TestCheckEncoding_NoIssues(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"PORT":     "8080",
		"DEBUG":    "true",
	}
	issues := CheckEncoding(env)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestCheckEncoding_NonASCIIValue(t *testing.T) {
	env := map[string]string{
		"GREETING": "héllo",
	}
	issues := CheckEncoding(env)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if !strings.Contains(issues[0].Message, "non-ASCII") {
		t.Errorf("expected non-ASCII message, got %q", issues[0].Message)
	}
}

func TestCheckEncoding_NonASCIIKey(t *testing.T) {
	env := map[string]string{
		"HÉLLO": "world",
	}
	issues := CheckEncoding(env)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if !strings.Contains(issues[0].Message, "key contains non-ASCII") {
		t.Errorf("expected key non-ASCII message, got %q", issues[0].Message)
	}
}

func TestCheckEncoding_ControlCharInValue(t *testing.T) {
	env := map[string]string{
		"TOKEN": "abc\x01def",
	}
	issues := CheckEncoding(env)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if !strings.Contains(issues[0].Message, "control characters") {
		t.Errorf("expected control char message, got %q", issues[0].Message)
	}
}

func TestCheckEncoding_SortedOutput(t *testing.T) {
	env := map[string]string{
		"Z_KEY": "héllo",
		"A_KEY": "wörld",
	}
	issues := CheckEncoding(env)
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}
	if issues[0].Key != "A_KEY" {
		t.Errorf("expected A_KEY first, got %q", issues[0].Key)
	}
	if issues[1].Key != "Z_KEY" {
		t.Errorf("expected Z_KEY second, got %q", issues[1].Key)
	}
}

func TestEncodingIssue_String(t *testing.T) {
	issue := EncodingIssue{Key: "FOO", Value: "bàr", Message: "value contains non-ASCII characters"}
	got := issue.String()
	if !strings.Contains(got, "FOO") {
		t.Errorf("expected key in string, got %q", got)
	}
	if !strings.Contains(got, "non-ASCII") {
		t.Errorf("expected message in string, got %q", got)
	}
}

func TestFormatEncodingIssues_Empty(t *testing.T) {
	out := FormatEncodingIssues(nil)
	if out != "no encoding issues found" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatEncodingIssues_WithIssues(t *testing.T) {
	issues := []EncodingIssue{
		{Key: "FOO", Value: "bàr", Message: "value contains non-ASCII characters"},
	}
	out := FormatEncodingIssues(issues)
	if !strings.Contains(out, "[encoding]") {
		t.Errorf("expected [encoding] tag, got %q", out)
	}
	if !strings.Contains(out, "FOO") {
		t.Errorf("expected key in output, got %q", out)
	}
}
