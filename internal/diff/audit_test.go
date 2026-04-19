package diff

import (
	"strings"
	"testing"
)

func TestGenerateAuditLog_NoChanges(t *testing.T) {
	m := map[string]string{"A": "1", "B": "2"}
	log := GenerateAuditLog(m, m)
	if log.HasChanges() {
		t.Fatalf("expected no changes, got %d", len(log))
	}
}

func TestGenerateAuditLog_Added(t *testing.T) {
	before := map[string]string{"A": "1"}
	after := map[string]string{"A": "1", "B": "2"}
	log := GenerateAuditLog(before, after)
	if len(log) != 1 || log[0].Action != "added" || log[0].Key != "B" {
		t.Fatalf("unexpected log: %+v", log)
	}
}

func TestGenerateAuditLog_Removed(t *testing.T) {
	before := map[string]string{"A": "1", "B": "2"}
	after := map[string]string{"A": "1"}
	log := GenerateAuditLog(before, after)
	if len(log) != 1 || log[0].Action != "removed" || log[0].Key != "B" {
		t.Fatalf("unexpected log: %+v", log)
	}
}

func TestGenerateAuditLog_Changed(t *testing.T) {
	before := map[string]string{"A": "old"}
	after := map[string]string{"A": "new"}
	log := GenerateAuditLog(before, after)
	if len(log) != 1 || log[0].Action != "changed" || log[0].OldVal != "old" || log[0].NewVal != "new" {
		t.Fatalf("unexpected log: %+v", log)
	}
}

func TestAuditEntry_String(t *testing.T) {
	before := map[string]string{"X": "foo"}
	after := map[string]string{"X": "bar"}
	log := GenerateAuditLog(before, after)
	s := log[0].String()
	if !strings.Contains(s, "CHANGED") || !strings.Contains(s, "X") {
		t.Fatalf("unexpected string: %s", s)
	}
}

func TestGenerateAuditLog_SortedByKey(t *testing.T) {
	before := map[string]string{}
	after := map[string]string{"Z": "1", "A": "2", "M": "3"}
	log := GenerateAuditLog(before, after)
	if log[0].Key != "A" || log[1].Key != "M" || log[2].Key != "Z" {
		t.Fatalf("not sorted: %+v", log)
	}
}
