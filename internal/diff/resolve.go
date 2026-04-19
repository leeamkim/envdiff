package diff

import "fmt"

// ResolveStrategy defines how conflicts are resolved.
type ResolveStrategy int

const (
	StrategyPreferA ResolveStrategy = iota
	StrategyPreferB
	StrategyRequireMatch
)

// ResolveIssue represents a conflict that could not be resolved.
type ResolveIssue struct {
	Key    string
	ValueA string
	ValueB string
}

func (r ResolveIssue) String() string {
	return fmt.Sprintf("conflict on %q: %q vs %q", r.Key, r.ValueA, r.ValueB)
}

// ResolveResult holds the merged map and any unresolved conflicts.
type ResolveResult struct {
	Resolved map[string]string
	Conflicts []ResolveIssue
}

// Resolve merges two env maps using the given strategy.
func Resolve(a, b map[string]string, strategy ResolveStrategy) ResolveResult {
	resolved := make(map[string]string)
	var conflicts []ResolveIssue

	keys := make(map[string]struct{})
	for k := range a { keys[k] = struct{}{} }
	for k := range b { keys[k] = struct{}{} }

	for k := range keys {
		va, inA := a[k]
		vb, inB := b[k]

		switch {
		case inA && !inB:
			resolved[k] = va
		case !inA && inB:
			resolved[k] = vb
		case va == vb:
			resolved[k] = va
		default:
			switch strategy {
			case StrategyPreferA:
				resolved[k] = va
			case StrategyPreferB:
				resolved[k] = vb
			case StrategyRequireMatch:
				conflicts = append(conflicts, ResolveIssue{Key: k, ValueA: va, ValueB: vb})
			}
		}
	}

	return ResolveResult{Resolved: resolved, Conflicts: conflicts}
}
