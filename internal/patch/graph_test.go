package patch

import (
	"testing"
)

func makeGraphPatch(version string, paths []string) Patch {
	ops := make([]Op, len(paths))
	for i, p := range paths {
		ops[i] = Op{Op: "replace", Path: p, Value: "v"}
	}
	return Patch{Version: version, Ops: ops}
}

func TestBuildDependencyGraph_NoOverlap(t *testing.T) {
	patches := []Patch{
		makeGraphPatch("v1", []string{"database/host"}),
		makeGraphPatch("v2", []string{"server/port"}),
	}
	g, err := BuildDependencyGraph(patches)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(g.Edges["v2"]) != 0 {
		t.Errorf("expected no deps for v2, got %v", g.Edges["v2"])
	}
}

func TestBuildDependencyGraph_OverlapCreatesDep(t *testing.T) {
	patches := []Patch{
		makeGraphPatch("v1", []string{"database/host"}),
		makeGraphPatch("v2", []string{"database/host"}),
	}
	g, err := BuildDependencyGraph(patches)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(g.Edges["v2"]) == 0 {
		t.Error("expected v2 to depend on v1")
	}
}

func TestTopologicalSort_RespectsOrder(t *testing.T) {
	patches := []Patch{
		makeGraphPatch("v1", []string{"app/name"}),
		makeGraphPatch("v2", []string{"app/name"}),
		makeGraphPatch("v3", []string{"app/name"}),
	}
	g, err := BuildDependencyGraph(patches)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	order, err := TopologicalSort(g)
	if err != nil {
		t.Fatalf("sort error: %v", err)
	}
	if len(order) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(order))
	}
	index := func(v string) int {
		for i, n := range order {
			if n == v {
				return i
			}
		}
		return -1
	}
	if index("v1") > index("v2") {
		t.Error("v1 should come before v2")
	}
	if index("v2") > index("v3") {
		t.Error("v2 should come before v3")
	}
}

func TestPathOverlaps_NestedPath(t *testing.T) {
	if !pathOverlaps("database", "database/host") {
		t.Error("expected overlap between parent and child path")
	}
}

func TestPathOverlaps_DisjointPaths(t *testing.T) {
	if pathOverlaps("server/host", "database/host") {
		t.Error("expected no overlap between disjoint paths")
	}
}

func TestBuildDependencyGraph_AllNodes(t *testing.T) {
	patches := []Patch{
		makeGraphPatch("v1", []string{"x"}),
		makeGraphPatch("v2", []string{"y"}),
		makeGraphPatch("v3", []string{"z"}),
	}
	g, err := BuildDependencyGraph(patches)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(g.Nodes) != 3 {
		t.Errorf("expected 3 nodes, got %d", len(g.Nodes))
	}
}
