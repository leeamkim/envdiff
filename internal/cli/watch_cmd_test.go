package cli

import (
	"bytes"
	"os"
	"testing"
	"time"
)

func writeWatchEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "watch*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestRunWatch_MissingArgs(t *testing.T) {
	err := RunWatch(WatchOptions{Files: []string{"only_one"}}, &bytes.Buffer{})
	if err == nil {
		t.Error("expected error for fewer than 2 files")
	}
}

func TestRunWatch_NoChange(t *testing.T) {
	a := writeWatchEnv(t, "KEY=val\n")
	b := writeWatchEnv(t, "KEY=val\n")
	var out bytes.Buffer
	err := RunWatch(WatchOptions{
		Files:    []string{a, b},
		Interval: 10 * time.Millisecond,
		MaxTicks: 1,
	}, &out)
	if err != nil {
		t.Fatal(err)
	}
	if out.Len() != 0 {
		t.Errorf("expected no output when files unchanged, got: %s", out.String())
	}
}

func TestRunWatch_DetectsChange(t *testing.T) {
	a := writeWatchEnv(t, "KEY=old\n")
	b := writeWatchEnv(t, "KEY=old\n")
	var out bytes.Buffer

	// Modify file a after a short delay in a goroutine.
	go func() {
		time.Sleep(15 * time.Millisecond)
		os.WriteFile(a, []byte("KEY=new\n"), 0644)
	}()

	err := RunWatch(WatchOptions{
		Files:    []string{a, b},
		Interval: 20 * time.Millisecond,
		MaxTicks: 2,
	}, &out)
	if err != nil {
		t.Fatal(err)
	}
	// At least one tick should have caught the change.
	if out.Len() == 0 {
		t.Error("expected output after file change")
	}
}
