package diff

import (
	"testing"
)

func baseAnnotateResult() Result {
	return Result{
		MissingInA:  map[string]string{"ONLY_B": "val"},
		MissingInB:  map[string]string{"ONLY_A": "val"},
		Mismatched:  map[string][2]string{"DIFF_KEY": {"v1", "v2"}},
		Matched:     map[string]string{"SAME": "val"},
	}
}

func TestAnnotate_MismatchedNote(t *testing.T) {
	r := baseAnnotateResult()
	ar := Annotate(r, []func(DiffEntry) string{AnnotatorMismatched})
	if note, ok := ar.Annotations["DIFF_KEY"]; !ok || note == "" {
		t.Errorf("expected annotation for DIFF_KEY, got %q", note)
	}
}

func TestAnnotate_MissingInB(t *testing.T) {
	r := baseAnnotateResult()
	ar := Annotate(r, []func(DiffEntry) string{AnnotatorMissingInB})
	if note, ok := ar.Annotations["ONLY_A"]; !ok || note == "" {
		t.Errorf("expected annotation for ONLY_A, got %q", note)
	}
}

func TestAnnotate_MissingInA(t *testing.T) {
	r := baseAnnotateResult()
	ar := Annotate(r, []func(DiffEntry) string{AnnotatorMissingInA})
	if note, ok := ar.Annotations["ONLY_B"]; !ok || note == "" {
		t.Errorf("expected annotation for ONLY_B, got %q", note)
	}
}

func TestAnnotate_NoAnnotationForMatched(t *testing.T) {
	r := baseAnnotateResult()
	ar := Annotate(r, []func(DiffEntry) string{
		AnnotatorMissingInA, AnnotatorMissingInB, AnnotatorMismatched,
	})
	if _, ok := ar.Annotations["SAME"]; ok {
		t.Error("expected no annotation for matched key")
	}
}

func TestAnnotate_MultipleAnnotators_FirstWins(t *testing.T) {
	r := baseAnnotateResult()
	first := func(e DiffEntry) string {
		if e.Key == "DIFF_KEY" {
			return "first"
		}
		return ""
	}
	second := func(e DiffEntry) string {
		if e.Key == "DIFF_KEY" {
			return "second"
		}
		return ""
	}
	ar := Annotate(r, []func(DiffEntry) string{first, second})
	if ar.Annotations["DIFF_KEY"] != "first" {
		t.Errorf("expected 'first', got %q", ar.Annotations["DIFF_KEY"])
	}
}
