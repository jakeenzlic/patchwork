package patch

import (
	"fmt"
	"strings"
)

// RedactRule describes a path pattern whose value should be masked.
type RedactRule struct {
	Path    string
	MaskWith string // defaults to "***"
}

// RedactResult holds the redacted config and a list of paths that were masked.
type RedactResult struct {
	Config  map[string]any
	Masked  []string
}

// Redact walks config and replaces values at paths matching any rule with a
// mask string, returning a deep copy so the original is not mutated.
func Redact(config map[string]any, rules []RedactRule) RedactResult {
	out := deepCopy(config)
	var masked []string

	for _, rule := range rules {
		mask := rule.MaskWith
		if mask == "" {
			mask = "***"
		}
		if redactPath(out, rule.Path, mask) {
			masked = append(masked, rule.Path)
		}
	}

	return RedactResult{Config: out, Masked: masked}
}

// redactPath sets the value at dotted/slash path to mask. Returns true if the
// path existed and was masked.
func redactPath(config map[string]any, path, mask string) bool {
	parts, err := splitPath(path)
	if err != nil || len(parts) == 0 {
		return false
	}

	cur := config
	for i, part := range parts {
		if i == len(parts)-1 {
			if _, ok := cur[part]; !ok {
				return false
			}
			cur[part] = mask
			return true
		}
		next, ok := cur[part]
		if !ok {
			return false
		}
		nextMap, ok := next.(map[string]any)
		if !ok {
			return false
		}
		cur = nextMap
	}
	return false
}

// FormatRedactResult returns a human-readable summary.
func FormatRedactResult(r RedactResult) string {
	if len(r.Masked) == 0 {
		return "redact: no paths matched\n"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "redact: masked %d path(s)\n", len(r.Masked))
	for _, p := range r.Masked {
		fmt.Fprintf(&sb, "  - %s\n", p)
	}
	return sb.String()
}
