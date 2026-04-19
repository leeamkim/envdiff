package diff

import (
	"testing"
	"time"
)

var t0 = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
var t1 = t0.Add(24 * time.Hour)
var t2 = t1.Add(24 * time.Hour)

func TestBuildTimeline_NoDiff(t *testing.T) {
	entries := []TimelineEntry{
		{Timestamp: t0, Label: "v1", Env: map[string]string{"A": "1"}},
		{Timestamp: t1, Label: "v2", Env: map[string]string{"A": "1"}},
	}
	events := BuildTimeline(entries)
	if len(events) != 0 {
		t.Fatalf("expected no events, got %d", len(events))
	}
}

func TestBuildTimeline_Added(t *testing.T) {
	entries := []TimelineEntry{
		{Timestamp: t0, Label: "v1", Env: map[string]string{}},
		{Timestamp: t1, Label: "v2", Env: map[string]string{"NEW": "val"}},
	}
	events := BuildTimeline(entries)
	if len(events) != 1 || events[0].Kind != "added" || events[0].Key != "NEW" {
		t.Fatalf("unexpected events: %+v", events)
	}
}

func TestBuildTimeline_Removed(t *testing.T) {
	entries := []TimelineEntry{
		{Timestamp: t0, Label: "v1", Env: map[string]string{"OLD": "x"}},
		{Timestamp: t1, Label: "v2", Env: map[string]string{}},
	}
	events := BuildTimeline(entries)
	if len(events) != 1 || events[0].Kind != "removed" || events[0].From != "x" {
		t.Fatalf("unexpected events: %+v", events)
	}
}

func TestBuildTimeline_Changed(t *testing.T) {
	entries := []TimelineEntry{
		{Timestamp: t0, Label: "v1", Env: map[string]string{"K": "old"}},
		{Timestamp: t1, Label: "v2", Env: map[string]string{"K": "new"}},
	}
	events := BuildTimeline(entries)
	if len(events) != 1 || events[0].Kind != "changed" || events[0].To != "new" {
		t.Fatalf("unexpected events: %+v", events)
	}
}

func TestBuildTimeline_MultipleSteps(t *testing.T) {
	entries := []TimelineEntry{
		{Timestamp: t0, Label: "v1", Env: map[string]string{"A": "1"}},
		{Timestamp: t1, Label: "v2", Env: map[string]string{"A": "2"}},
		{Timestamp: t2, Label: "v3", Env: map[string]string{"A": "2", "B": "3"}},
	}
	events := BuildTimeline(entries)
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
}

func TestTimelineEvent_String(t *testing.T) {
	e := TimelineEvent{Key: "X", To: "v", Kind: "added", Label: "v2"}
	s := e.String()
	if s == "" {
		t.Fatal("expected non-empty string")
	}
}
