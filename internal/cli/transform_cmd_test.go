package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTransformEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunTransform_MissingArgs(t *testing.T) {
	if err := RunTransform(nil); err == nil {
		t.Fatal("expected error")
	}
}

func TestRunTransform_Trim(t *testing.T) {
	p := writeTransformEnv(t, "KEY=  hello  \nOTHER=world\n")
	if err := RunTransform([]string{p, "--op=trim"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunTransform_Upper(t *testing.T) {
	p := writeTransformEnv(t, "KEY=hello\n")
	if err := RunTransform([]string{p, "--op=upper"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunTransform_Lower(t *testing.T) {
	p := writeTransformEnv(t, "KEY=HELLO\n")
	if err := RunTransform([]string{p, "--op=lower"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunTransform_SpecificKeys(t *testing.T) {
	p := writeTransformEnv(t, "A=hello\nB=world\n")
	if err := RunTransform([]string{p, "--op=upper", "--keys=A"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunTransform_UnknownOp(t *testing.T) {
	p := writeTransformEnv(t, "KEY=val\n")
	if err := RunTransform([]string{p, "--op=reverse"}); err == nil {
		t.Fatal("expected error for unknown op")
	}
}

func TestRunTransform_NoChange(t *testing.T) {
	p := writeTransformEnv(t, "KEY=HELLO\n")
	if err := RunTransform([]string{p, "--op=upper"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
