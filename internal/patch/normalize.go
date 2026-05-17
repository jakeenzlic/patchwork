package patch

import (
	"fmt"
	"strings"
)

// NormalizeResult holds the outcome of a normalization pass over a patch.
type NormalizeResult struct {
	Patch    Patch
	Changes  []string
}

// Normalize applies a set of canonical transformations to a patch:
//   - trims whitespace from all path segments
//   - lower-cases op names
//   - removes duplicate consecutive ops targeting the same path+op
//   - strips leading/trailing slashes from paths
func Normalize(p Patch) (NormalizeResult, error) {
	if err := Validate(p); err != nil {
		return NormalizeResult{}, fmt.Errorf("normalize: invalid patch: %w", err)
	}

	result := NormalizeResult{Patch: deepCopyPatch(p)}
	seen := make(map[string]bool)
	var normalised []Op

	for i, op := range result.Patch.Ops {
		_ = i
		normOp := op

		// Lowercase op
		lower := strings.ToLower(string(normOp.Op))
		if lower != string(normOp.Op) {
			result.Changes = append(result.Changes, fmt.Sprintf("op[%d]: lowercased op %q -> %q", i, normOp.Op, lower))
			normOp.Op = OpType(lower)
		}

		// Strip leading/trailing slashes from path
		cleaned := strings.Trim(normOp.Path, "/")
		// Trim whitespace from each segment
		segs := strings.Split(cleaned, "/")
		for j, s := range segs {
			segs[j] = strings.TrimSpace(s)
		}
		cleaned = strings.Join(segs, "/")
		if cleaned != normOp.Path {
			result.Changes = append(result.Changes, fmt.Sprintf("op[%d]: normalised path %q -> %q", i, normOp.Path, cleaned))
			normOp.Path = cleaned
		}

		// Deduplicate consecutive identical op+path
		key := string(normOp.Op) + ":" + normOp.Path
		if seen[key] {
			result.Changes = append(result.Changes, fmt.Sprintf("op[%d]: removed duplicate %s on %q", i, normOp.Op, normOp.Path))
			continue
		}
		seen[key] = true
		normalised = append(normalised, normOp)
	}

	result.Patch.Ops = normalised
	return result, nil
}

// FormatNormalizeResult returns a human-readable summary of the normalization.
func FormatNormalizeResult(r NormalizeResult) string {
	if len(r.Changes) == 0 {
		return "normalize: no changes required\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("normalize: %d change(s):\n", len(r.Changes)))
	for _, c := range r.Changes {
		sb.WriteString("  - " + c + "\n")
	}
	return sb.String()
}
