package patch

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBundle_MultiplePatches_PreservesOrder(t *testing.T) {
	dir := t.TempDir()
	dest := CompressedBundlePath(dir, "ordered")

	patches := []Patch{
		{ID: "a", Version: "1.0", Ops: []Op{{Op: "add", Path: "x", Value: 1}}},
		{ID: "b", Version: "1.0", Ops: []Op{{Op: "add", Path: "y", Value: 2}}},
		{ID: "c", Version: "1.0", Ops: []Op{{Op: "add", Path: "z", Value: 3}}},
	}

	if err := Bundle(patches, dest); err != nil {
		t.Fatalf("Bundle: %v", err)
	}

	got, err := Unbundle(dest)
	if err != nil {
		t.Fatalf("Unbundle: %v", err)
	}

	ids := []string{"a", "b", "c"}
	for i, id := range ids {
		if got[i].ID != id {
			t.Errorf("position %d: got ID %q, want %q", i, got[i].ID, id)
		}
	}
}

func TestBundle_CreatesIntermediateDirs(t *testing.T) {
	base := t.TempDir()
	dest := filepath.Join(base, "deep", "nested", "bundle.patch.gz")

	if err := Bundle(makeBundlePatches(), dest); err != nil {
		t.Fatalf("Bundle: %v", err)
	}

	if _, err := os.Stat(dest); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}

func TestBundle_EmptySlice(t *testing.T) {
	dir := t.TempDir()
	dest := CompressedBundlePath(dir, "empty")

	if err := Bundle([]Patch{}, dest); err != nil {
		t.Fatalf("Bundle empty: %v", err)
	}

	got, err := Unbundle(dest)
	if err != nil {
		t.Fatalf("Unbundle empty: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected 0 patches, got %d", len(got))
	}
}
