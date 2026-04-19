package diff

import (
	"fmt"
	"sort"
	"time"
)

// AuditEntry records a single change event between two env states.
type AuditEntry struct {
	Key    string
	OldVal string
	NewVal string
	Action string // "added", "removed", "changed"
	At     time.Time
}

func (a AuditEntry) String() string {
	switch a.Action {
	case "added":
		return fmt.Sprintf("[%s] ADDED %s=%q", a.At.Format(time.RFC3339), a.Key, a.NewVal)
	case "removed":
		return fmt.Sprintf("[%s] REMOVED %s (was %q)", a.At.Format(time.RFC3339), a.Key, a.OldVal)
	default:
		return fmt.Sprintf("[%s] CHANGED %s: %q -> %q", a.At.Format(time.RFC3339), a.Key, a.OldVal, a.NewVal)
	}
}

// AuditLog is an ordered list of AuditEntry values.
type AuditLog []AuditEntry

// GenerateAuditLog compares two env maps and returns an audit log stamped with now.
func GenerateAuditLog(before, after map[string]string) AuditLog {
	now := time.Now().UTC()
	var log AuditLog

	for k, newVal := range after {
		if oldVal, ok := before[k]; !ok {
			log = append(log, AuditEntry{Key: k, NewVal: newVal, Action: "added", At: now})
		} else if oldVal != newVal {
			log = append(log, AuditEntry{Key: k, OldVal: oldVal, NewVal: newVal, Action: "changed", At: now})
		}
	}

	for k, oldVal := range before {
		if _, ok := after[k]; !ok {
			log = append(log, AuditEntry{Key: k, OldVal: oldVal, Action: "removed", At: now})
		}
	}

	sort.Slice(log, func(i, j int) bool { return log[i].Key < log[j].Key })
	return log
}

// HasChanges returns true if the audit log contains any entries.
func (l AuditLog) HasChanges() bool { return len(l) > 0 }
