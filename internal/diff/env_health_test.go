package diff

import (
	"testing"
)

func TestComputeHealth_Empty(t *testing.T) {
	h := ComputeHealth(map[string]string{})
	if h.Grade != "N/A" {
		t.Errorf("expected N/A, got %s", h.Grade)
	}
}

func TestComputeHealth_Perfect(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "8080",
	}
	h := ComputeHealth(env)
	if h.Score != 100 {
		t.Errorf("expected score 100, got %d", h.Score)
	}
	if h.Grade != "A" {
		t.Errorf("expected grade A, got %s", h.Grade)
	}
}

func TestComputeHealth_EmptyValues(t *testing.T) {
	env := map[string]string{
		"HOST": "",
		"PORT": "8080",
	}
	h := ComputeHealth(env)
	if h.EmptyValues != 1 {
		t.Errorf("expected 1 empty value, got %d", h.EmptyValues)
	}
	if h.Score >= 100 {
		t.Error("expected score below 100")
	}
}

func TestComputeHealth_Placeholders(t *testing.T) {
	env := map[string]string{
		"API_KEY": "CHANGEME",
		"HOST":    "localhost",
	}
	h := ComputeHealth(env)
	if h.Placeholders != 1 {
		t.Errorf("expected 1 placeholder, got %d", h.Placeholders)
	}
}

func TestComputeHealth_LowercaseKeys(t *testing.T) {
	env := map[string]string{
		"host": "localhost",
		"PORT": "8080",
	}
	h := ComputeHealth(env)
	if h.LowercaseKeys != 1 {
		t.Errorf("expected 1 lowercase key, got %d", h.LowercaseKeys)
	}
}

func TestComputeHealth_GradeF(t *testing.T) {
	env := map[string]string{
		"a": "",
		"b": "CHANGEME",
		"c": "",
		"d": "TODO",
	}
	h := ComputeHealth(env)
	if h.Grade != "F" {
		t.Errorf("expected grade F, got %s", h.Grade)
	}
}

func TestHealthStatus_String(t *testing.T) {
	h := HealthStatus{TotalKeys: 5, Score: 80, Grade: "B"}
	s := h.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}
