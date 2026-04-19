package diff

import (
	"testing"
)

func TestEnvGraph_AddNodeAndNames(t *testing.T) {
	g := NewEnvGraph()
	g.AddNode("prod", map[string]string{"A": "1"})
	g.AddNode("staging", map[string]string{"A": "2"})
	names := g.NodeNames()
	if len(names) != 2 || names[0] != "prod" || names[1] != "staging" {
		t.Errorf("unexpected names: %v", names)
	}
}

func TestEnvGraph_AddNodeCopiesVars(t *testing.T) {
	g := NewEnvGraph()
	vars := map[string]string{"X": "1"}
	g.AddNode("dev", vars)
	vars["X"] = "mutated"
	if g.Nodes["dev"].Vars["X"] != "1" {
		t.Error("AddNode should copy vars")
	}
}

func TestEnvGraph_Neighbors(t *testing.T) {
	g := NewEnvGraph()
	g.AddNode("dev", map[string]string{})
	g.AddNode("staging", map[string]string{})
	g.AddNode("prod", map[string]string{})
	g.AddEdge("dev", "staging")
	g.AddEdge("dev", "prod")
	neighbors := g.Neighbors("dev")
	if len(neighbors) != 2 || neighbors[0] != "prod" || neighbors[1] != "staging" {
		t.Errorf("unexpected neighbors: %v", neighbors)
	}
}

func TestDiffGraph_NoDiff(t *testing.T) {
	g := NewEnvGraph()
	g.AddNode("a", map[string]string{"K": "v"})
	g.AddNode("b", map[string]string{"K": "v"})
	g.AddEdge("a", "b")
	entries := DiffGraph(g)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Result.HasDiff() {
		t.Error("expected no diff")
	}
}

func TestDiffGraph_WithDiff(t *testing.T) {
	g := NewEnvGraph()
	g.AddNode("a", map[string]string{"K": "1", "ONLY_A": "x"})
	g.AddNode("b", map[string]string{"K": "2"})
	g.AddEdge("a", "b")
	entries := DiffGraph(g)
	if !entries[0].Result.HasDiff() {
		t.Error("expected diff")
	}
}

func TestDiffGraph_SkipsMissingNodes(t *testing.T) {
	g := NewEnvGraph()
	g.AddNode("a", map[string]string{})
	g.Edges = append(g.Edges, EnvEdge{From: "a", To: "ghost"})
	entries := DiffGraph(g)
	if len(entries) != 0 {
		t.Error("expected no entries for missing node")
	}
}

func TestGraphDiffSummary(t *testing.T) {
	g := NewEnvGraph()
	g.AddNode("dev", map[string]string{"A": "1"})
	g.AddNode("prod", map[string]string{"A": "2"})
	g.AddEdge("dev", "prod")
	summary := GraphDiffSummary(g)
	if _, ok := summary["dev->prod"]; !ok {
		t.Error("expected key dev->prod in summary")
	}
}
