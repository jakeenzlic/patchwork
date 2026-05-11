package patch

import (
	"fmt"
	"strings"
)

// KnownOps is the set of supported patch operations.
var KnownOps = map[string]bool{
	"add":     true,
	"remove":  true,
	"replace": true,
}

// ValidationError holds a list of validation failures for a patch.
type ValidationError struct {
	Errors []string
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("patch validation failed:\n  %s", strings.Join(v.Errors, "\n  "))
}

// Validate checks that a Patch is well-formed before it is applied.
// It returns a *ValidationError if any issues are found, or nil on success.
func Validate(p *Patch) error {
	var errs []string

	if p.Version == "" {
		errs = append(errs, "missing required field: version")
	}

	for i, op := range p.Ops {
		if !KnownOps[op.Op] {
			errs = append(errs, fmt.Sprintf("op[%d]: unknown operation %q", i, op.Op))
		}
		if op.Path == "" {
			errs = append(errs, fmt.Sprintf("op[%d]: missing required field: path", i))
		}
		if op.Op != "remove" && op.Value == nil {
			errs = append(errs, fmt.Sprintf("op[%d]: operation %q requires a value", i, op.Op))
		}
	}

	if len(errs) > 0 {
		return &ValidationError{Errors: errs}
	}
	return nil
}
