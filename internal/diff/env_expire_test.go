package diff

import (
	"testing"
	"time"
)

var (
	now      = time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	pastDate = time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	nearDate = time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	farDate  = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
)

func TestCheckExpiry_NoIssues(t *testing.T) {
	env := map[string]string{"API_KEY": "abc123"}
	rules := []ExpiryRule{{Pattern: "API_KEY", ExpiresAt: farDate}}
	issues := CheckExpiry(env, rules, 30, now)
	if len(issues) != 0 {
		t.Fatalf("expected 0 issues, got %d", len(issues))
	}
}

func TestCheckExpiry_Expired(t *testing.T) {
	env := map[string]string{"API_KEY": "abc123"}
	rules := []ExpiryRule{{Pattern: "API_KEY", ExpiresAt: pastDate}}
	issues := CheckExpiry(env, rules, 30, now)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if !issues[0].Expired {
		t.Error("expected issue to be marked expired")
	}
}

func TestCheckExpiry_WarningSoon(t *testing.T) {
	env := map[string]string{"API_KEY": "abc123"}
	rules := []ExpiryRule{{Pattern: "API_KEY", ExpiresAt: nearDate}}
	issues := CheckExpiry(env, rules, 30, now)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Expired {
		t.Error("expected issue to not be expired, just a warning")
	}
}

func TestCheckExpiry_WildcardPattern(t *testing.T) {
	env := map[string]string{"DB_SECRET": "x", "DB_TOKEN": "y", "OTHER": "z"}
	rules := []ExpiryRule{{Pattern: "DB_*", ExpiresAt: pastDate}}
	issues := CheckExpiry(env, rules, 30, now)
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}
}

func TestCheckExpiry_SortedOutput(t *testing.T) {
	env := map[string]string{"Z_KEY": "a", "A_KEY": "b"}
	rules := []ExpiryRule{{Pattern: "*", ExpiresAt: pastDate}}
	issues := CheckExpiry(env, rules, 30, now)
	if len(issues) < 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}
	if issues[0].Key > issues[1].Key {
		t.Error("expected issues to be sorted by key")
	}
}

func TestFormatExpireIssues_Empty(t *testing.T) {
	out := FormatExpireIssues(nil)
	if out != "no expiry issues found" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestExpireIssue_String(t *testing.T) {
	issue := ExpireIssue{Key: "API_KEY", ExpiresAt: pastDate, Expired: true}
	got := issue.String()
	if got != "API_KEY: expired on 2024-05-01" {
		t.Errorf("unexpected string: %q", got)
	}
}
