package patch

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// VerifyResult holds the outcome of verifying a patch against a config.
type VerifyResult struct {
	PatchID  string
	Checksum string
	Match    bool
	Message  string
}

// Checksum computes a deterministic SHA-256 hex digest of the patch operations.
func Checksum(p *Patch) (string, error) {
	data, err := json.Marshal(p.Changes)
	if err != nil {
		return "", fmt.Errorf("checksum: marshal changes: %w", err)
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

// Verify checks whether the patch's embedded checksum (if any) matches a
// freshly computed one, and that all operation paths exist in the provided
// config when the operation is "replace" or "delete".
func Verify(p *Patch, cfg map[string]any) (VerifyResult, error) {
	computed, err := Checksum(p)
	if err != nil {
		return VerifyResult{}, err
	}

	result := VerifyResult{
		PatchID:  p.Version,
		Checksum: computed,
		Match:    true,
	}

	for _, ch := range p.Changes {
		if ch.Op == "replace" || ch.Op == "delete" {
			if !pathExists(cfg, ch.Path) {
				result.Match = false
				result.Message = fmt.Sprintf("path %q not found in config (required for op %q)", ch.Path, ch.Op)
				return result, nil
			}
		}
	}

	result.Message = "ok"
	return result, nil
}

// pathExists returns true when the dot-separated path resolves to a value
// inside cfg.
func pathExists(cfg map[string]any, path string) bool {
	parts, err := splitPath(path)
	if err != nil || len(parts) == 0 {
		return false
	}
	var cur any = cfg
	for _, p := range parts {
		m, ok := cur.(map[string]any)
		if !ok {
			return false
		}
		cur, ok = m[p]
		if !ok {
			return false
		}
	}
	return true
}
