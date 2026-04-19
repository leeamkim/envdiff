package diff

import (
	"testing"
)

func TestResolve_NoConflicts(t *testing.T) {
	a := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := map[string]string{"FOO": "bar", "BAZ": "qux"}
	r := Resolve(a, b, StrategyRequireMatch)
	if len(r.Conflicts) != 0 {
		t.Fatalf("expected no conflicts, got %d", len(r.Conflicts))
	}
	if r.Resolved["FOO"] != "bar" {
		t.Errorf("expected FOO=bar")
	}
}

func TestResolve_PreferA(t *testing.T) {
	a := map[string]string{"KEY": "from-a"}
	b := map[string]string{"KEY": "from-b"}
	r := Resolve(a, b, StrategyPreferA)
	if r.Resolved["KEY"] != "from-a" {
		t.Errorf("expected from-a, got %s", r.Resolved["KEY"])
	}
	if len(r.Conflicts) != 0 {
		t.Errorf("expected no conflicts")
	}
}

func TestResolve_PreferB(t *testing.T) {
	a := map[string]string{"KEY": "from-a"}
	b := map[string]string{"KEY": "from-b"}
	r := Resolve(a, b, StrategyPreferB)
	if r.Resolved["KEY"] != "from-b" {
		t.Errorf("expected from-b, got %s", r.Resolved["KEY"])
	}
}

func TestResolve_RequireMatch_Conflict(t *testing.T) {
	a := map[string]string{"KEY": "aaa", "SHARED": "same"}
	b := map[string]string{"KEY": "bbb", "SHARED": "same"}
	r := Resolve(a, b, StrategyRequireMatch)
	if len(r.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(r.Conflicts))
	}
	if r.Conflicts[0].Key != "KEY" {
		t.Errorf("expected conflict on KEY")
	}
	if r.Resolved["SHARED"] != "same" {
		t.Errorf("expected SHARED resolved")
	}
}

func TestResolve_MissingInOne(t *testing.T) {
	a := map[string]string{"ONLY_A": "val"}
	b := map[string]string{"ONLY_B": "val2"}
	r := Resolve(a, b, StrategyRequireMatch)
	if r.Resolved["ONLY_A"] != "val" {
		t.Errorf("expected ONLY_A")
	}
	if r.Resolved["ONLY_B"] != "val2" {
		t.Errorf("expected ONLY_B")
	}
	if len(r.Conflicts) != 0 {
		t.Errorf("expected no conflicts for unique keys")
	}
}

func TestResolveIssue_String(t *testing.T) {
	i := ResolveIssue{Key: "X", ValueA: "a", ValueB: "b"}
	s := i.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}
