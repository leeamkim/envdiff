package diff

import "fmt"

// GraphDiffEntry records a diff between two connected environment nodes.
type GraphDiffEntry struct {
	From   string
	To     string
	Result Result
}

// String returns a human-readable summary of the graph diff entry.
func (e GraphDiffEntry) String() string {
	s := Summarize(e.Result)
	return fmt.Sprintf("%s -> %s: %s", e.From, e.To, s.String())
}

// DiffGraph computes diffs for all edges in the graph.
func DiffGraph(g *EnvGraph) []GraphDiffEntry {
	var entries []GraphDiffEntry
	for _, edge := range g.Edges {
		fromNode, ok1 := g.Nodes[edge.From]
		toNode, ok2 := g.Nodes[edge.To]
		if !ok1 || !ok2 {
			continue
		}
		result := Compare(fromNode.Vars, toNode.Vars)
		entries = append(entries, GraphDiffEntry{
			From:   edge.From,
			To:     edge.To,
			Result: result,
		})
	}
	return entries
}

// GraphDiffSummary returns a map of "from->to" to summary string for all edges.
func GraphDiffSummary(g *EnvGraph) map[string]string {
	entries := DiffGraph(g)
	out := make(map[string]string, len(entries))
	for _, e := range entries {
		key := fmt.Sprintf("%s->%s", e.From, e.To)
		out[key] = e.String()
	}
	return out
}
