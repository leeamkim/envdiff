package diff

import (
	"strings"
	"testing"
)

func TestCheckObsolete_NoIssues(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	ref := map[string]string{"A": "x", "B": "y", "C": "z"}
	issues := CheckObsolete(env, ref)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestCheckObsolete_SingleObsolete(t *testing.T) {
	env := map[string]string{"A": "1", "LEGACY": "old"}
	ref := map[string]string{"A": "x"}
	issues := CheckObsolete(env, ref)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "LEGACY" {
		t.Errorf("expected key LEGACY, got %q", issues[0].Key)
	}
}

func TestCheckObsolete_MultipleObsolete(t *testing.T) {
	env := map[string]string{"A": "1", "OLD1": "x", "OLD2": "y"}
	ref := map[string]string{"A": "z"}
	issues := CheckObsolete(env, ref)
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}
	if issues[0].Key != "OLD1" || issues[1].Key != "OLD2" {
		t.Errorf("unexpected order: %v, %v", issues[0].Key, issues[1].Key)
	}
}

func TestCheckObsolete_SortedOutput(t *testing.T) {
	env := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	ref := map[string]string{}
	issues := CheckObsolete(env, ref)
	for i := 1; i < len(issues); i++ {
		if issues[i-1].Key > issues[i].Key {
			t.Errorf("not sorted: %q before %q", issues[i-1].Key, issues[i].Key)
		}
	}
}

func TestObsoleteIssue_String(t *testing.T) {
	issue := ObsoleteIssue{Key: "OLD_KEY", Value: "val"}
	s := issue.String()
	if !strings.Contains(s, "OLD_KEY") {
		t.Errorf("expected key in string, got %q", s)
	}
}

func TestFormatObsoleteIssues_Empty(t *testing.T) {
	out := FormatObsoleteIssues(nil)
	if !strings.Contains(out, "no obsolete") {
		t.Errorf("expected no-issue message, got %q", out)
	}
}

func TestFormatObsoleteIssues_WithIssues(t *testing.T) {
	issues := []ObsoleteIssue{{Key: "DEAD_KEY", Value: "old"}}
	out := FormatObsoleteIssues(issues)
	if !strings.Contains(out, "DEAD_KEY") {
		t.Errorf("expected DEAD_KEY in output, got %q", out)
	}
}
