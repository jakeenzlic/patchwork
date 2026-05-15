package patch

import (
	"fmt"
	"strings"
)

// DependencyGraph represents patch ordering constraints derived from path overlap.
type DependencyGraph struct {
	Nodes []string
	Edges map[string][]string // node -> list of nodes it depends on
}

// BuildDependencyGraph analyses a slice of patches and constructs a dependency
// graph based on overlapping paths. A patch that writes to a path that another
// patch reads (replace/delete) is considered a dependency.
func BuildDependencyGraph(patches []Patch) (*DependencyGraph, error) {
	g := &DependencyGraph{
		Edges: make(map[string][]string),
	}

	for _, p := range patches {
		g.Nodes = append(g.Nodes, p.Version)
		if _, ok := g.Edges[p.Version]; !ok {
			g.Edges[p.Version] = []string{}
		}
	}

	for i, a := range patches {
		for j, b := range patches {
			if i >= j {
				continue
			}
			if patchesOverlap(a, b) {
				// b depends on a (a must run first)
				g.Edges[b.Version] = append(g.Edges[b.Version], a.Version)
			}
		}
	}

	if err := detectCycle(g); err != nil {
		return nil, err
	}

	return g, nil
}

// TopologicalSort returns patch versions in a valid application order.
func TopologicalSort(g *DependencyGraph) ([]string, error) {
	visited := make(map[string]bool)
	temp := make(map[string]bool)
	var order []string

	var visit func(n string) error
	visit = func(n string) error {
		if temp[n] {
			return fmt.Errorf("cycle detected at node %q", n)
		}
		if visited[n] {
			return nil
		}
		temp[n] = true
		for _, dep := range g.Edges[n] {
			if err := visit(dep); err != nil {
				return err
			}
		}
		temp[n] = false
		visited[n] = true
		order = append(order, n)
		return nil
	}

	for _, n := range g.Nodes {
		if err := visit(n); err != nil {
			return nil, err
		}
	}

	return order, nil
}

func patchesOverlap(a, b Patch) bool {
	for _, opA := range a.Ops {
		for _, opB := range b.Ops {
			if pathOverlaps(opA.Path, opB.Path) {
				return true
			}
		}
	}
	return false
}

func pathOverlaps(a, b string) bool {
	partsA := strings.Split(strings.Trim(a, "/"), "/")
	partsB := strings.Split(strings.Trim(b, "/"), "/")
	min := len(partsA)
	if len(partsB) < min {
		min = len(partsB)
	}
	for i := 0; i < min; i++ {
		if partsA[i] != partsB[i] {
			return false
		}
	}
	return true
}

func detectCycle(g *DependencyGraph) error {
	_, err := TopologicalSort(g)
	return err
}
