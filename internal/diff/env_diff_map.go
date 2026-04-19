package diff

// EnvDiffMap holds a matrix-style diff across multiple named environments.
type EnvDiffMap struct {
	Keys  []string
	Envs  []string
	Cells map[string]map[string]string // key -> env -> value
}

// BuildEnvDiffMap constructs an EnvDiffMap from a map of env name -> key/value pairs.
func BuildEnvDiffMap(envs map[string]map[string]string) *EnvDiffMap {
	keySet := map[string]struct{}{}
	for _, vars := range envs {
		for k := range vars {
			keySet[k] = struct{}{}
		}
	}
	keys := sortStrings(keySetToSlice(keySet))
	envNames := sortStrings(mapKeys(envs))

	cells := make(map[string]map[string]string, len(keys))
	for _, k := range keys {
		cells[k] = make(map[string]string, len(envNames))
		for _, env := range envNames {
			if v, ok := envs[env][k]; ok {
				cells[k][env] = v
			} else {
				cells[k][env] = ""
			}
		}
	}
	return &EnvDiffMap{Keys: keys, Envs: envNames, Cells: cells}
}

// Consistent returns true if all envs have the same non-empty value for the key.
func (m *EnvDiffMap) Consistent(key string) bool {
	vals := m.Cells[key]
	var ref string
	set := false
	for _, v := range vals {
		if !set {
			ref = v
			set = true
			continue
		}
		if v != ref {
			return false
		}
	}
	return true
}

func keySetToSlice(s map[string]struct{}) []string {
	out := make([]string, 0, len(s))
	for k := range s {
		out = append(out, k)
	}
	return out
}

func mapKeys(m map[string]map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
