package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot captures the state of an env file at a point in time.
type Snapshot struct {
	Label     string            `json:"label"`
	Timestamp time.Time         `json:"timestamp"`
	Entries   map[string]string `json:"entries"`
}

// SnapshotDiff describes changes between two snapshots.
type SnapshotDiff struct {
	From    string
	To      string
	Added   map[string]string
	Removed map[string]string
	Changed map[string][2]string // key -> [old, new]
}

// NewSnapshot creates a snapshot from an env map.
func NewSnapshot(label string, entries map[string]string) Snapshot {
	copy := make(map[string]string, len(entries))
	for k, v := range entries {
		copy[k] = v
	}
	return Snapshot{Label: label, Timestamp: time.Now(), Entries: copy}
}

// SaveSnapshot writes a snapshot to a JSON file.
func SaveSnapshot(path string, s Snapshot) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal snapshot: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadSnapshot reads a snapshot from a JSON file.
func LoadSnapshot(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("read snapshot: %w", err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return Snapshot{}, fmt.Errorf("unmarshal snapshot: %w", err)
	}
	return s, nil
}

// DiffSnapshots compares two snapshots and returns the differences.
func DiffSnapshots(a, b Snapshot) SnapshotDiff {
	d := SnapshotDiff{
		From:    a.Label,
		To:      b.Label,
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string][2]string),
	}
	for k, v := range b.Entries {
		if old, ok := a.Entries[k]; !ok {
			d.Added[k] = v
		} else if old != v {
			d.Changed[k] = [2]string{old, v}
		}
	}
	for k, v := range a.Entries {
		if _, ok := b.Entries[k]; !ok {
			d.Removed[k] = v
		}
	}
	return d
}

// HasChanges returns true if the diff contains any changes.
func (d SnapshotDiff) HasChanges() bool {
	return len(d.Added)+len(d.Removed)+len(d.Changed) > 0
}
