package patch

import "fmt"

// Op represents a single patch operation.
type Op struct {
	Op    string      `json:"op"    yaml:"op"`
	Path  string      `json:"path"  yaml:"path"`
	Value interface{} `json:"value,omitempty" yaml:"value,omitempty"`
}

// Patch is a versioned, ordered list of operations with optional metadata.
type Patch struct {
	ID          string   `json:"id"          yaml:"id"`
	Version     string   `json:"version"     yaml:"version"`
	Description string   `json:"description" yaml:"description"`
	Tags        []string `json:"tags,omitempty" yaml:"tags,omitempty"`
	Ops         []Op     `json:"ops"         yaml:"ops"`
}

// KnownOps is the set of operation types recognised by patchwork.
var KnownOps = map[string]struct{}{
	"add":     {},
	"delete":  {},
	"replace": {},
}

// String returns a short human-readable description of a patch.
func (p Patch) String() string {
	tags := ""
	if len(p.Tags) > 0 {
		tags = fmt.Sprintf(" [tags: %v]", p.Tags)
	}
	return fmt.Sprintf("%s (v%s)%s — %d op(s)", p.ID, p.Version, tags, len(p.Ops))
}
