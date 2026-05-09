package patch

import (
	"fmt"
	"time"
)

// OpType represents the type of patch operation.
type OpType string

const (
	OpAdd     OpType = "add"
	OpRemove  OpType = "remove"
	OpReplace OpType = "replace"
)

// Operation describes a single change to a config value.
type Operation struct {
	Op    OpType      `json:"op" yaml:"op"`
	Path  string      `json:"path" yaml:"path"`
	Value interface{} `json:"value,omitempty" yaml:"value,omitempty"`
}

// Patch represents a versioned set of operations to apply to a config.
type Patch struct {
	Version     string      `json:"version" yaml:"version"`
	Description string      `json:"description" yaml:"description"`
	CreatedAt   time.Time   `json:"created_at" yaml:"created_at"`
	Ops         []Operation `json:"ops" yaml:"ops"`
}

// Validate checks that all operations in the patch are well-formed.
func (p *Patch) Validate() error {
	if p.Version == "" {
		return fmt.Errorf("patch version must not be empty")
	}
	for i, op := range p.Ops {
		if op.Path == "" {
			return fmt.Errorf("op[%d]: path must not be empty", i)
		}
		switch op.Op {
		case OpAdd, OpRemove, OpReplace:
			// valid
		default:
			return fmt.Errorf("op[%d]: unknown op type %q", i, op.Op)
		}
		if op.Op != OpRemove && op.Value == nil {
			return fmt.Errorf("op[%d]: value required for op %q", i, op.Op)
		}
	}
	return nil
}

// Summary returns a human-readable summary of the patch, including its version,
// description, and the number of operations it contains.
func (p *Patch) Summary() string {
	desc := p.Description
	if desc == "" {
		desc = "(no description)"
	}
	return fmt.Sprintf("Patch %s: %s (%d op(s))", p.Version, desc, len(p.Ops))
}
