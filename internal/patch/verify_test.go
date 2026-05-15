package patch

import (
	"testing"
)

var baseVerifyPatch = &Patch{
	Version: "v1.2.0",
	Changes: []Change{
		{Op: "replace", Path: "server.port", Value: 9090},
		{Op: "add", Path: "feature.dark_mode", Value: true},
	},
}

func TestChecksum_Deterministic(t *testing.T) {
	sum1, err := Checksum(baseVerifyPatch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sum2, err := Checksum(baseVerifyPatch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum1 != sum2 {
		t.Errorf("checksums differ: %q vs %q", sum1, sum2)
	}
}

func TestChecksum_DiffersOnChange(t *testing.T) {
	other := &Patch{
		Version: "v1.2.0",
		Changes: []Change{
			{Op: "replace", Path: "server.port", Value: 8080},
		},
	}
	sum1, _ := Checksum(baseVerifyPatch)
	sum2, _ := Checksum(other)
	if sum1 == sum2 {
		t.Error("expected different checksums for different patches")
	}
}

func TestChecksum_DiffersOnVersion(t *testing.T) {
	other := &Patch{
		Version: "v9.9.9",
		Changes: []Change{
			{Op: "replace", Path: "server.port", Value: 9090},
			{Op: "add", Path: "feature.dark_mode", Value: true},
		},
	}
	sum1, _ := Checksum(baseVerifyPatch)
	sum2, _ := Checksum(other)
	if sum1 == sum2 {
		t.Error("expected different checksums for patches with different versions")
	}
}

func TestVerify_AllPathsPresent(t *testing.T) {
	cfg := map[string]any{
		"server": map[string]any{"port": 8080},
	}
	res, err := Verify(baseVerifyPatch, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Match {
		t.Errorf("expected match=true, got message: %s", res.Message)
	}
}

func TestVerify_MissingReplacePath(t *testing.T) {
	p := &Patch{
		Version: "v1.0.0",
		Changes: []Change{
			{Op: "replace", Path: "database.host", Value: "newhost"},
		},
	}
	cfg := map[string]any{
		"server": map[string]any{"port": 8080},
	}
	res, err := Verify(p, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Match {
		t.Error("expected match=false for missing replace path")
	}
}

func TestVerify_AddDoesNotRequirePath(t *testing.T) {
	p := &Patch{
		Version: "v1.0.0",
		Changes: []Change{
			{Op: "add", Path: "brand.new.key", Value: 42},
		},
	}
	res, err := Verify(p, map[string]any{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Match {
		t.Errorf("expected match=true for add op, got: %s", res.Message)
	}
}

func TestPathExists_NestedTrue(t *testing.T) {
	cfg := map[string]any{
		"a": map[string]any{"b": map[string]any{"c": 1}},
	}
	if !pathExists(cfg, "a.b.c") {
		t.Error("expected path a.b.c to exist")
	}
}

func TestPathExists_Missing(t *testing.T) {
	cfg := map[string]any{"x": 1}
	if pathExists(cfg, "x.y") {
		t.Error("expected path x.y to not exist")
	}
}
