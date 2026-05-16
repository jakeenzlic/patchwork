package patch

import (
	"fmt"
	"strings"
)

// Scope restricts patch operations to a subtree of the config.
type Scope struct {
	Prefix string
}

// NewScope creates a Scope that limits operations to paths under prefix.
func NewScope(prefix string) (*Scope, error) {
	prefix = strings.TrimRight(prefix, "/")
	if prefix == "" {
		return nil, fmt.Errorf("scope prefix must not be empty")
	}
	if strings.HasPrefix(prefix, "/") {
		return nil, fmt.Errorf("scope prefix must not start with '/'")
	}
	return &Scope{Prefix: prefix}, nil
}

// Filter returns only the operations whose paths fall under the scope prefix.
func (s *Scope) Filter(patches []Patch) []Patch {
	var out []Patch
	for _, p := range patches {
		filtered := filterOps(p, s.Prefix)
		if len(filtered.Operations) > 0 {
			out = append(out, filtered)
		}
	}
	return out
}

// InScope reports whether path falls under the given prefix.
func InScope(prefix, path string) bool {
	prefix = strings.TrimRight(prefix, "/")
	path = strings.TrimLeft(path, "/")
	pfx := strings.TrimLeft(prefix, "/")
	return path == pfx || strings.HasPrefix(path, pfx+"/")
}

// FormatScope returns a human-readable summary of a scoped filter result.
func FormatScope(prefix string, total, kept int) string {
	return fmt.Sprintf("scope=%q  total_ops=%d  in_scope=%d  skipped=%d",
		prefix, total, kept, total-kept)
}

func filterOps(p Patch, prefix string) Patch {
	out := p
	out.Operations = nil
	for _, op := range p.Operations {
		if InScope(prefix, op.Path) {
			out.Operations = append(out.Operations, op)
		}
	}
	return out
}
