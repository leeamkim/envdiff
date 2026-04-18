package diff

import "fmt"

// DriftEntry represents a single key that has drifted from its baseline value.
type DriftEntry struct {
	Key      string
	OldValue string
	NewValue string
	Status   string // "added", "removed", "changed"
}

func (d DriftEntry) String() string {
	switch d.Status {
	case "added":
		return fmt.Sprintf("[added]   %s = %q", d.Key, d.NewValue)
	case "removed":
		return fmt.Sprintf("[removed] %s (was %q)", d.Key, d.OldValue)
	default:
		return fmt.Sprintf("[changed] %s: %q -> %q", d.Key, d.OldValue, d.NewValue)
	}
}

// DriftReport holds all drift entries and metadata.
type DriftReport struct {
	Entries []DriftEntry
}

func (r *DriftReport) HasDrift() bool {
	return len(r.Entries) > 0
}

func (r *DriftReport) Summary() string {
	if !r.HasDrift() {
		return "no drift detected"
	}
	return fmt.Sprintf("%d key(s) drifted", len(r.Entries))
}

// DetectDrift compares a current env map against a reference map and reports drift.
func DetectDrift(reference, current map[string]string) *DriftReport {
	report := &DriftReport{}

	for k, refVal := range reference {
		curVal, ok := current[k]
		if !ok {
			report.Entries = append(report.Entries, DriftEntry{Key: k, OldValue: refVal, Status: "removed"})
		} else if curVal != refVal {
			report.Entries = append(report.Entries, DriftEntry{Key: k, OldValue: refVal, NewValue: curVal, Status: "changed"})
		}
	}

	for k, curVal := range current {
		if _, ok := reference[k]; !ok {
			report.Entries = append(report.Entries, DriftEntry{Key: k, NewValue: curVal, Status: "added"})
		}
	}

	return report
}
