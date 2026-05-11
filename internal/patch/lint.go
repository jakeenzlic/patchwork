package patch

import (
	"fmt"
	"strings"
)

// LintWarning represents a non-fatal advisory about a patch.
type LintWarning struct {
	Index   int
	Message string
}

func (w LintWarning) String() string {
	return fmt.Sprintf("op[%d]: %s", w.Index, w.Message)
}

// Lint inspects a Patch for common issues that are not hard errors but may
// indicate mistakes. It returns a (possibly empty) slice of LintWarnings.
func Lint(p *Patch) []LintWarning {
	var warnings []LintWarning

	seen := make(map[string]int)
	for i, op := range p.Ops {
		key := op.Op + ":" + op.Path
		if prev, ok := seen[key]; ok {
			warnings = append(warnings, LintWarning{
				Index:   i,
				Message: fmt.Sprintf("duplicate operation on path %q (first seen at op[%d])", op.Path, prev),
			})
		} else {
			seen[key] = i
		}

		if strings.HasPrefix(op.Path, "/") || strings.HasSuffix(op.Path, "/") {
			warnings = append(warnings, LintWarning{
				Index:   i,
				Message: fmt.Sprintf("path %q has leading or trailing slash", op.Path),
			})
		}

		if strings.Contains(op.Path, "//") {
			warnings = append(warnings, LintWarning{
				Index:   i,
				Message: fmt.Sprintf("path %q contains empty segment", op.Path),
			})
		}
	}

	return warnings
}
