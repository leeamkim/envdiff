package diff

import (
	"strings"
	"testing"
)

func TestComputeCoverage_Perfect(t *testing.T) {
	ref := map[string]string{"A": "1", "B": "2", "C": "3"}
	tgt := map[string]string{"A": "1", "B": "2", "C": "3"}

	r := ComputeCoverage(ref, tgt)
	if r.Coverage != 1.0 {
		t.Errorf("expected 1.0, got %f", r.Coverage)
	}
	if r.Grade != "A" {
		t.Errorf("expected grade A, got %s", r.Grade)
	}
	if len(r.MissingKeys) != 0 {
		t.Errorf("expected no missing keys")
	}
}

func TestComputeCoverage_SomeMissing(t *testing.T) {
	ref := map[string]string{"A": "1", "B": "2", "C": "3", "D": "4"}
	tgt := map[string]string{"A": "1", "B": "2"}

	r := ComputeCoverage(ref, tgt)
	if r.PresentKeys != 2 {
		t.Errorf("expected 2 present, got %d", r.PresentKeys)
	}
	if r.TotalKeys != 4 {
		t.Errorf("expected 4 total, got %d", r.TotalKeys)
	}
	if len(r.MissingKeys) != 2 {
		t.Errorf("expected 2 missing, got %v", r.MissingKeys)
	}
}

func TestComputeCoverage_EmptyValue(t *testing.T) {
	ref := map[string]string{"A": "1", "B": "2"}
	tgt := map[string]string{"A": "1", "B": ""}

	r := ComputeCoverage(ref, tgt)
	if r.PresentKeys != 1 {
		t.Errorf("empty value should count as missing")
	}
	if !contains(r.MissingKeys, "B") {
		t.Errorf("B should be in missing keys")
	}
}

func TestComputeCoverage_EmptyReference(t *testing.T) {
	r := ComputeCoverage(map[string]string{}, map[string]string{"A": "1"})
	if r.Grade != "A" {
		t.Errorf("empty reference should yield grade A")
	}
}

func TestCoverageResult_String(t *testing.T) {
	r := CoverageResult{TotalKeys: 4, PresentKeys: 3, Coverage: 0.75, Grade: "C", MissingKeys: []string{"X"}}
	s := r.String()
	if !strings.Contains(s, "75.0%") {
		t.Errorf("expected percentage in string, got %s", s)
	}
	if !strings.Contains(s, "grade=C") {
		t.Errorf("expected grade in string, got %s", s)
	}
}

func TestCoverageGrade(t *testing.T) {
	cases := []struct {
		coverage float64
		grade    string
	}{
		{1.0, "A"}, {0.95, "A"}, {0.80, "B"}, {0.65, "C"}, {0.50, "D"}, {0.49, "F"},
	}
	for _, tc := range cases {
		g := coverageGrade(tc.coverage)
		if g != tc.grade {
			t.Errorf("coverage %.2f: expected %s got %s", tc.coverage, tc.grade, g)
		}
	}
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
