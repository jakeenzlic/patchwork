package patch

import (
	"os"
	"path/filepath"
	"testing"
)

func makeBundlePatches() []Patch {
	return []Patch{
		{
			ID:      "p-001",
			Version: "1.0",
			Ops: []Op{
				{Op: "add", Path: "feature/enabled", Value: true},
			},
		},
		{
			ID:      "p-002",
			Version: "1.0",
			Ops: []Op{
				{Op: "replace", Path: "feature/limit", Value: float64(100)},
			},
		},
	}
}

func TestBundle_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	dest := CompressedBundlePath(dir, "release-1")

	if err := Bundle(makeBundlePatches(), dest); err != nil {
		t.Fatalf("Bundle: %v", err)
	}

	if _, err := os.Stat(dest); err != nil {
		t.Fatalf("expected bundle file to exist: %v", err)
	}
}

func TestUnbundle_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	dest := CompressedBundlePath(dir, "release-1")
	original := makeBundlePatches()

	if err := Bundle(original, dest); err != nil {
		t.Fatalf("Bundle: %v", err)
	}

	got, err := Unbundle(dest)
	if err != nil {
		t.Fatalf("Unbundle: %v", err)
	}

	if len(got) != len(original) {
		t.Fatalf("expected %d patches, got %d", len(original), len(got))
	}
	for i, p := range got {
		if p.ID != original[i].ID {
			t.Errorf("patch[%d] ID: got %q, want %q", i, p.ID, original[i].ID)
		}
	}
}

func TestUnbundle_NotExist(t *testing.T) {
	_, err := Unbundle(filepath.Join(t.TempDir(), "missing.patch.gz"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestIsBundle(t *testing.T) {
	cases := []struct {
		path string
		want bool
	}{
		{"release.patch.gz", true},
		{"release.json", false},
		{"release.yaml", false},
		{"/some/dir/v2.patch.gz", true},
	}
	for _, c := range cases {
		if got := IsBundle(c.path); got != c.want {
			t.Errorf("IsBundle(%q) = %v, want %v", c.path, got, c.want)
		}
	}
}

func TestCompressedBundlePath(t *testing.T) {
	got := CompressedBundlePath("/patches", "v1")
	want := "/patches/v1.patch.gz"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
