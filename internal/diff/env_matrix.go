package diff

import "sort"

// MatrixCell holds the comparison result between two environments for a key.
type MatrixCell struct {
	Status string // "match", "mismatch", "missing"
	ValueA string
	ValueB string
}

// EnvMatrix is a cross-environment comparison: keys x env-pairs.
type EnvMatrix struct {
	EnvNames []string
	Keys     []string
	// Cells[key][envName] = value
	Cells map[string]map[string]string
	// Status[key]["envA:envB"] = MatrixCell
	Status map[string]map[string]MatrixCell
}

// BuildMatrix compares multiple environments against the first (base) environment.
func BuildMatrix(envs map[string]map[string]string) EnvMatrix {
	keySet := map[string]struct{}{}
	for _, vars := range envs {
		for k := range vars {
			keySet[k] = struct{}{}
		}
	}
	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	envNames := make([]string, 0, len(envs))
	for name := range envs {
		envNames = append(envNames, name)
	}
	sort.Strings(envNames)

	cells := map[string]map[string]string{}
	for _, k := range keys {
		cells[k] = map[string]string{}
		for _, name := range envNames {
			cells[k][name] = envs[name][k]
		}
	}

	status := map[string]map[string]MatrixCell{}
	if len(envNames) < 2 {
		return EnvMatrix{EnvNames: envNames, Keys: keys, Cells: cells, Status: status}
	}
	base := envNames[0]
	for _, k := range keys {
		status[k] = map[string]MatrixCell{}
		for _, name := range envNames[1:] {
			pair := base + ":" + name
			aVal, aOk := envs[base][k]
			bVal, bOk := envs[name][k]
			var st string
			switch {
			case !aOk && !bOk:
				st = "missing"
			case !aOk || !bOk:
				st = "missing"
			case aVal == bVal:
				st = "match"
			default:
				st = "mismatch"
			}
			status[k][pair] = MatrixCell{Status: st, ValueA: aVal, ValueB: bVal}
		}
	}
	return EnvMatrix{EnvNames: envNames, Keys: keys, Cells: cells, Status: status}
}
