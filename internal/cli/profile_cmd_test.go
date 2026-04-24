package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeProfileEnv(t *testing.T, dir, name string, pairs map[string]string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	for k, v := range pairs {
		fmt.Fprintf(f, "%s=%s\n", k, v)
	}
	return path
}

func TestRunProfile_MissingArgs(t *testing.T) {
	var sb strings.Builder
	err := RunProfile([]string{"dev=file.env"}, &sb)
	if err == nil {
		t.Fatal("expected error for single arg")
	}
}

func TestRunProfile_InvalidArgFormat(t *testing.T) {
	var sb strings.Builder
	err := RunProfile([]string{"noequals", "dev=file.env"}, &sb)
	if err == nil {
		t.Fatal("expected error for invalid arg format")
	}
}

func TestRunProfile_InvalidFile(t *testing.T) {
	var sb strings.Builder
	err := RunProfile([]string{"dev=/no/such/file.env", "prod=/also/missing.env"}, &sb)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRunProfile_NoDiff(t *testing.T) {
	dir := t.TempDir()
	dev := writeProfileEnv(t, dir, "dev.env", map[string]string{"APP": "x", "PORT": "8080"})
	prod := writeProfileEnv(t, dir, "prod.env", map[string]string{"APP": "x", "PORT": "8080"})

	var sb strings.Builder
	err := RunProfile([]string{"dev=" + dev, "prod=" + prod}, &sb)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "APP") {
		t.Error("expected APP in output")
	}
	if !strings.Contains(out, "PORT") {
		t.Error("expected PORT in output")
	}
}

func TestRunProfile_ShowsMissingKey(t *testing.T) {
	dir := t.TempDir()
	dev := writeProfileEnv(t, dir, "dev.env", map[string]string{"ONLY_DEV": "yes"})
	prod := writeProfileEnv(t, dir, "prod.env", map[string]string{})

	var sb strings.Builder
	err := RunProfile([]string{"dev=" + dev, "prod=" + prod}, &sb)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "(missing)") {
		t.Errorf("expected (missing) in output, got:\n%s", out)
	}
}

func TestRunProfile_ThreeProfiles(t *testing.T) {
	dir := t.TempDir()
	dev := writeProfileEnv(t, dir, "dev.env", map[string]string{"X": "1"})
	staging := writeProfileEnv(t, dir, "staging.env", map[string]string{"X": "2"})
	prod := writeProfileEnv(t, dir, "prod.env", map[string]string{"X": "3"})

	var sb strings.Builder
	err := RunProfile([]string{"dev=" + dev, "staging=" + staging, "prod=" + prod}, &sb)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "staging") {
		t.Error("expected staging profile in output")
	}
}
