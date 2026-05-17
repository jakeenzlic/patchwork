package patch

import "fmt"

// RenameResult holds the outcome of a rename operation.
type RenameResult struct {
	OldPath string
	NewPath string
	PatchID string
	Renamed bool
}

// Rename rewrites all patch operations whose path starts with oldPrefix
// to use newPrefix instead. It returns the updated patch and a result summary.
func Rename(p Patch, oldPrefix, newPrefix string) (Patch, RenameResult, error) {
	if oldPrefix == "" {
		return p, RenameResult{}, fmt.Errorf("rename: oldPrefix must not be empty")
	}
	if newPrefix == "" {
		return p, RenameResult{}, fmt.Errorf("rename: newPrefix must not be empty")
	}

	updated := deepCopyPatch(p)
	renamed := false

	for i, op := range updated.Ops {
		if pathHasPrefix(op.Path, oldPrefix) {
			updated.Ops[i].Path = newPrefix + op.Path[len(oldPrefix):]
			renamed = true
		}
	}

	return updated, RenameResult{
		OldPath: oldPrefix,
		NewPath: newPrefix,
		PatchID: p.ID,
		Renamed: renamed,
	}, nil
}

// FormatRenameResult returns a human-readable summary of a rename operation.
func FormatRenameResult(r RenameResult) string {
	if !r.Renamed {
		return fmt.Sprintf("rename: no ops matched prefix %q in patch %s", r.OldPath, r.PatchID)
	}
	return fmt.Sprintf("rename: %s — rewrote ops from %q to %q", r.PatchID, r.OldPath, r.NewPath)
}

// pathHasPrefix reports whether path starts with prefix, respecting segment boundaries.
func pathHasPrefix(path, prefix string) bool {
	if len(path) < len(prefix) {
		return false
	}
	if path == prefix {
		return true
	}
	if path[:len(prefix)] == prefix && (prefix[len(prefix)-1] == '/' || path[len(prefix)] == '/') {
		return true
	}
	return false
}

// deepCopyPatch returns a shallow-safe copy of a Patch.
func deepCopyPatch(p Patch) Patch {
	ops := make([]Op, len(p.Ops))
	copy(ops, p.Ops)
	p.Ops = ops
	return p
}
