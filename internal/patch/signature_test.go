package patch

import (
	"os"
	"path/filepath"
	"testing"
)

var baseSignaturePatch = &Patch{
	ID:      "sig-001",
	Version: "1.0",
	Ops: []Op{
		{Op: "add", Path: "feature/enabled", Value: true},
	},
}

func TestSignPatch_Deterministic(t *testing.T) {
	secret := []byte("supersecret")
	s1, err := SignPatch(baseSignaturePatch, secret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s2, err := SignPatch(baseSignaturePatch, secret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s1 != s2 {
		t.Errorf("expected deterministic signatures, got %q and %q", s1, s2)
	}
}

func TestSignPatch_EmptySecret(t *testing.T) {
	_, err := SignPatch(baseSignaturePatch, []byte{})
	if err == nil {
		t.Error("expected error for empty secret, got nil")
	}
}

func TestVerifySignature_Valid(t *testing.T) {
	secret := []byte("supersecret")
	p := &Patch{ID: "sig-002", Version: "1.0", Ops: baseSignaturePatch.Ops}
	sig, _ := SignPatch(p, secret)
	p.Signature = sig
	if err := VerifySignature(p, secret); err != nil {
		t.Errorf("expected valid signature, got error: %v", err)
	}
}

func TestVerifySignature_Mismatch(t *testing.T) {
	p := &Patch{ID: "sig-003", Version: "1.0", Ops: baseSignaturePatch.Ops, Signature: "deadbeef"}
	if err := VerifySignature(p, []byte("key")); err == nil {
		t.Error("expected mismatch error, got nil")
	}
}

func TestVerifySignature_Missing(t *testing.T) {
	p := &Patch{ID: "sig-004", Version: "1.0", Ops: baseSignaturePatch.Ops}
	if err := VerifySignature(p, []byte("key")); err == nil {
		t.Error("expected missing-signature error, got nil")
	}
}

func TestSaveAndLoadSignature(t *testing.T) {
	dir := t.TempDir()
	patchPath := filepath.Join(dir, "patch.yaml")
	const sig = "abc123"
	if err := SaveSignature(patchPath, sig); err != nil {
		t.Fatalf("SaveSignature: %v", err)
	}
	got, err := LoadSignature(patchPath)
	if err != nil {
		t.Fatalf("LoadSignature: %v", err)
	}
	if got != sig {
		t.Errorf("expected %q, got %q", sig, got)
	}
}

func TestLoadSignature_NotExist(t *testing.T) {
	_, err := LoadSignature(filepath.Join(t.TempDir(), "missing.yaml"))
	if !os.IsNotExist(err) && err == nil {
		t.Error("expected not-exist error")
	}
}
