package diff

import (
	"testing"
)

func TestNewIgnoreList_Contains(t *testing.T) {
	il := NewIgnoreList([]string{"SECRET", "TOKEN"})
	if !il.Contains("SECRET") {
		t.Error("expected SECRET to be in ignore list")
	}
	if il.Contains("HOST") {
		t.Error("expected HOST not to be in ignore list")
	}
}

func TestIgnoreList_Apply_MissingInA(t *testing.T) {
	il := NewIgnoreList([]string{"SECRET"})
	r := Result{MissingInA: []string{"SECRET", "HOST"}, MissingInB: []string{}, Mismatched: map[string][2]string{}}
	out := il.Apply(r)
	if len(out.MissingInA) != 1 || out.MissingInA[0] != "HOST" {
		t.Errorf("unexpected MissingInA: %v", out.MissingInA)
	}
}

func TestIgnoreList_Apply_MissingInB(t *testing.T) {
	il := NewIgnoreList([]string{"TOKEN"})
	r := Result{MissingInA: []string{}, MissingInB: []string{"TOKEN", "DB"}, Mismatched: map[string][2]string{}}
	out := il.Apply(r)
	if len(out.MissingInB) != 1 || out.MissingInB[0] != "DB" {
		t.Errorf("unexpected MissingInB: %v", out.MissingInB)
	}
}

func TestIgnoreList_Apply_Mismatched(t *testing.T) {
	il := NewIgnoreList([]string{"PORT"})
	r := Result{
		MissingInA: []string{},
		MissingInB: []string{},
		Mismatched: map[string][2]string{
			"PORT": {"8080", "9090"},
			"HOST": {"a", "b"},
		},
	}
	out := il.Apply(r)
	if _, ok := out.Mismatched["PORT"]; ok {
		t.Error("expected PORT to be ignored")
	}
	if _, ok := out.Mismatched["HOST"]; !ok {
		t.Error("expected HOST to remain")
	}
}

func TestIgnoreList_Apply_Empty(t *testing.T) {
	il := NewIgnoreList([]string{})
	r := Result{MissingInA: []string{"A"}, MissingInB: []string{}, Mismatched: map[string][2]string{}}
	out := il.Apply(r)
	if len(out.MissingInA) != 1 {
		t.Errorf("expected 1 entry, got %d", len(out.MissingInA))
	}
}

func TestIgnoreList_Apply_AllIgnored(t *testing.T) {
	il := NewIgnoreList([]string{"SECRET", "TOKEN", "PORT"})
	r := Result{
		MissingInA: []string{"SECRET"},
		MissingInB: []string{"TOKEN"},
		Mismatched: map[string][2]string{"PORT": {"8080", "9090"}},
	}
	out := il.Apply(r)
	if len(out.MissingInA) != 0 || len(out.MissingInB) != 0 || len(out.Mismatched) != 0 {
		t.Errorf("expected all entries to be ignored, got %+v", out)
	}
}
