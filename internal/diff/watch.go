package diff

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"time"
)

// WatchEvent represents a change detected during file watching.
type WatchEvent struct {
	File      string
	ChangedAt time.Time
	OldHash   string
	NewHash   string
}

func (e WatchEvent) String() string {
	return fmt.Sprintf("[%s] %s changed (was %s, now %s)",
		e.ChangedAt.Format(time.RFC3339), e.File, e.OldHash[:8], e.NewHash[:8])
}

// HashFile returns the MD5 hex digest of a file's contents.
func HashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// WatchOnce checks whether any of the given files have changed compared to
// the provided hash snapshot. Returns events for files that changed.
func WatchOnce(files []string, hashes map[string]string) ([]WatchEvent, map[string]string, error) {
	updated := make(map[string]string, len(files))
	var events []WatchEvent
	for _, f := range files {
		h, err := HashFile(f)
		if err != nil {
			return nil, nil, fmt.Errorf("watch: %w", err)
		}
		updated[f] = h
		prev, seen := hashes[f]
		if seen && prev != h {
			events = append(events, WatchEvent{
				File:      f,
				ChangedAt: time.Now(),
				OldHash:   prev,
				NewHash:   h,
			})
		}
	}
	return events, updated, nil
}
