package diff

import (
	"fmt"
	"sort"
	"time"
)

// TimelineEntry records a snapshot of an env map at a point in time.
type TimelineEntry struct {
	Timestamp time.Time
	Label     string
	Env       map[string]string
}

// TimelineEvent describes a change between two timeline entries.
type TimelineEvent struct {
	Key      string
	From     string
	To       string
	Kind     string // "added", "removed", "changed"
	At       time.Time
	Label    string
}

func (e TimelineEvent) String() string {
	switch e.Kind {
	case "added":
		return fmt.Sprintf("[%s] +%s=%q", e.Label, e.Key, e.To)
	case "removed":
		return fmt.Sprintf("[%s] -%s (was %q)", e.Label, e.Key, e.From)
	default:
		return fmt.Sprintf("[%s] ~%s: %q -> %q", e.Label, e.Key, e.From, e.To)
	}
}

// BuildTimeline compares consecutive entries and returns all change events.
func BuildTimeline(entries []TimelineEntry) []TimelineEvent {
	var events []TimelineEvent
	for i := 1; i < len(entries); i++ {
		prev := entries[i-1]
		curr := entries[i]
		events = append(events, diffEnvs(prev.Env, curr.Env, curr.Timestamp, curr.Label)...)
	}
	return events
}

func diffEnvs(a, b map[string]string, at time.Time, label string) []TimelineEvent {
	var events []TimelineEvent
	seen := map[string]bool{}
	for k, va := range a {
		seen[k] = true
		vb, ok := b[k]
		if !ok {
			events = append(events, TimelineEvent{Key: k, From: va, Kind: "removed", At: at, Label: label})
		} else if va != vb {
			events = append(events, TimelineEvent{Key: k, From: va, To: vb, Kind: "changed", At: at, Label: label})
		}
	}
	for k, vb := range b {
		if !seen[k] {
			events = append(events, TimelineEvent{Key: k, To: vb, Kind: "added", At: at, Label: label})
		}
	}
	sort.Slice(events, func(i, j int) bool { return events[i].Key < events[j].Key })
	return events
}
