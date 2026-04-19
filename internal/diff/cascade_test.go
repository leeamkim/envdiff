package diff

import (
	"testing"
)

func TestCascade_NoDrift(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"PORT": "8080", "APP": "myapp"},
		"prod": {"PORT": "8080", "APP": "myapp"},
	}
	result := Cascade([]string{"dev", "prod"}, envs)
	for _, e := range result.Entries {
		if e.Drift {
			t.Errorf("expected no drift for key %s", e.Key)
		}
	}
}

func TestCascade_Drift(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"PORT": "8080"},
		"prod": {"PORT": "443"},
	}
	result := Cascade([]string{"dev", "prod"}, envs)
	if len(result.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result.Entries))
	}
	if !result.Entries[0].Drift {
		t.Error("expected drift for PORT")
	}
}

func TestCascade_MissingInSomeEnvs(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":     {"PORT": "8080", "DEBUG": "true"},
		"staging": {"PORT": "8080"},
		"prod":    {"PORT": "8080"},
	}
	result := Cascade([]string{"dev", "staging", "prod"}, envs)
	for _, e := range result.Entries {
		if e.Key == "DEBUG" {
			if len(e.Values) != 1 {
				t.Errorf("expected DEBUG in 1 env, got %d", len(e.Values))
			}
			return
		}
	}
	t.Error("DEBUG entry not found")
}

func TestCascadeEntry_String_Drift(t *testing.T) {
	e := CascadeEntry{Key: "PORT", Values: map[string]string{"dev": "8080", "prod": "443"}, Drift: true}
	s := e.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}

func TestCascadeEntry_String_Consistent(t *testing.T) {
	e := CascadeEntry{Key: "APP", Values: map[string]string{"dev": "myapp", "prod": "myapp"}, Drift: false}
	s := e.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}
