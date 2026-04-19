package diff

import "sort"

// EnvNode represents a single environment in the graph.
type EnvNode struct {
	Name string
	Vars map[string]string
}

// EnvEdge represents a directed relationship between two environments.
type EnvEdge struct {
	From string
	To   string
}

// EnvGraph holds multiple environments and their relationships.
type EnvGraph struct {
	Nodes map[string]*EnvNode
	Edges []EnvEdge
}

// NewEnvGraph creates an empty EnvGraph.
func NewEnvGraph() *EnvGraph {
	return &EnvGraph{Nodes: make(map[string]*EnvNode)}
}

// AddNode adds or replaces an environment node.
func (g *EnvGraph) AddNode(name string, vars map[string]string) {
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	g.Nodes[name] = &EnvNode{Name: name, Vars: copy}
}

// AddEdge adds a directed edge from one environment to another.
func (g *EnvGraph) AddEdge(from, to string) {
	g.Edges = append(g.Edges, EnvEdge{From: from, To: to})
}

// Neighbors returns all nodes reachable from the given node name.
func (g *EnvGraph) Neighbors(name string) []string {
	var result []string
	for _, e := range g.Edges {
		if e.From == name {
			result = append(result, e.To)
		}
	}
	sort.Strings(result)
	return result
}

// NodeNames returns sorted list of all node names.
func (g *EnvGraph) NodeNames() []string {
	names := make([]string, 0, len(g.Nodes))
	for n := range g.Nodes {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
