package diff

import (
	"os"
	"testing"
)

func writeTempWatchFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "watch*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestHashFile_Stable(t *testing.T) {
	path := writeTempWatchFile(t, "KEY=value")
	h1, err := HashFile(path)
	if err != nil {
		t.Fatal(err)
	}
	h2, _ := HashFile(path)
	if h1 != h2 {
		t.Errorf("expected stable hash, got %s vs %s", h1, h2)
	}
}

func TestWatchOnce_NoChange(t *testing.T) {
	path := writeTempWatchFile(t, "KEY=value")
	h, _ := HashFile(path)
	events, _, err := WatchOnce([]string{path}, map[string]string{path: h})
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestWatchOnce_Detected(t *testing.T) {
	path := writeTempWatchFile(t, "KEY=old")
	oldHash, _ := HashFile(path)
	os.WriteFile(path, []byte("KEY=new"), 0644)
	events, updated, err := WatchOnce([]string{path}, map[string]string{path: oldHash})
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].OldHash != oldHash {
		t.Errorf("unexpected old hash")
	}
	if updated[path] == oldHash {
		t.Errorf("expected updated hash to differ")
	}
}

func TestWatchOnce_FirstSeen(t *testing.T) {
	path := writeTempWatchFile(t, "KEY=value")
	events, updated, err := WatchOnce([]string{path}, map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 0 {
		t.Errorf("first-seen file should not produce event")
	}
	if updated[path] == "" {
		t.Errorf("expected hash recorded for new file")
	}
}

func TestWatchEvent_String(t *testing.T) {
	e := WatchEvent{File: "a.env", OldHash: "aabbccdd1122", NewHash: "ffee99881122"}
	s := e.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}
