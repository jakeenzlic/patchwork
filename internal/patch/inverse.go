package patch

// invertPatch returns a new Patch whose operations are the logical inverse
// of the supplied patch, suitable for undoing its effects.
//
// Inversion rules:
//   - "add"     → "delete" (remove the key that was added)
//   - "delete"  → "add"    (restore the old value stored in OldValue)
//   - "replace" → "replace" with From/To swapped
func invertPatch(p Patch) Patch {
	inv := Patch{
		Version:     p.Version,
		Description: "inverse of: " + p.Description,
		Ops:         make([]Op, 0, len(p.Ops)),
	}
	// Reverse the op order so that undo is applied in reverse sequence.
	for i := len(p.Ops) - 1; i >= 0; i-- {
		op := p.Ops[i]
		switch op.Op {
		case "add":
			inv.Ops = append(inv.Ops, Op{
				Op:   "delete",
				Path: op.Path,
			})
		case "delete":
			inv.Ops = append(inv.Ops, Op{
				Op:    "add",
				Path:  op.Path,
				Value: op.OldValue,
			})
		case "replace":
			inv.Ops = append(inv.Ops, Op{
				Op:       "replace",
				Path:     op.Path,
				Value:    op.OldValue,
				OldValue: op.Value,
			})
		default:
			// Unknown ops are passed through unchanged.
			inv.Ops = append(inv.Ops, op)
		}
	}
	return inv
}
