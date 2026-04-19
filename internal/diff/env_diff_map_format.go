package diff

import (
	"encoding/json"
	"fmt"
	"strings"
)

// DiffMapRow is a single row in a serializable diff map report.
type DiffMapRow struct {
	Key        string            `json:"key"`
	Consistent bool              `json:"consistent"`
	Values     map[string]string `json:"values"`
}

// DiffMapReport is a serializable representation of an EnvDiffMap.
type DiffMapReport struct {
	Envs []string     `json:"envs"`
	Rows []DiffMapRow `json:"rows"`
}

// ToReport converts an EnvDiffMap to a DiffMapReport.
func (m *EnvDiffMap) ToReport() DiffMapReport {
	rows := make([]DiffMapRow, 0, len(m.Keys))
	for _, k := range m.Keys {
		vals := make(map[string]string, len(m.Envs))
		for _, env := range m.Envs {
			vals[env] = m.Cells[k][env]
		}
		rows = append(rows, DiffMapRow{
			Key:        k,
			Consistent: m.Consistent(k),
			Values:     vals,
		})
	}
	return DiffMapReport{Envs: m.Envs, Rows: rows}
}

// ToJSON serializes the report as JSON.
func (r DiffMapReport) ToJSON() (string, error) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Summary returns a short human-readable summary.
func (r DiffMapReport) Summary() string {
	inconsistent := 0
	for _, row := range r.Rows {
		if !row.Consistent {
			inconsistent++
		}
	}
	lines := []string{
		fmt.Sprintf("envs: %s", strings.Join(r.Envs, ", ")),
		fmt.Sprintf("total keys: %d", len(r.Rows)),
		fmt.Sprintf("inconsistent keys: %d", inconsistent),
	}
	return strings.Join(lines, "\n")
}
