package patch

import "fmt"

// Conflict describes a detected conflict between two patches.
type Conflict struct {
	PatchA string
	PatchB string
	Path   string
	Reason string
}

func (c Conflict) String() string {
	return fmt.Sprintf("conflict between %q and %q at path %q: %s",
		c.PatchA, c.PatchB, c.Path, c.Reason)
}

// DetectConflicts returns all conflicts found among the given patches.
// Two patches conflict when they both write to the same path with
// incompatible operations, or when one deletes a path the other depends on.
func DetectConflicts(patches []Patch) []Conflict {
	var conflicts []Conflict

	for i := 0; i < len(patches); i++ {
		for j := i + 1; j < len(patches); j++ {
			a, b := patches[i], patches[j]
			for _, opA := range a.Ops {
				for _, opB := range b.Ops {
					if !pathOverlaps(opA.Path, opB.Path) {
						continue
					}
					reason := conflictReason(opA.Op, opB.Op)
					if reason == "" {
						continue
					}
					conflicts = append(conflicts, Conflict{
						PatchA: a.ID,
						PatchB: b.ID,
						Path:   opA.Path,
						Reason: reason,
					})
				}
			}
		}
	}
	return conflicts
}

// conflictReason returns a human-readable reason when opA and opB conflict,
// or an empty string when they are compatible.
func conflictReason(opA, opB string) string {
	switch {
	case opA == "replace" && opB == "replace":
		return "both patches replace the same path"
	case opA == "delete" && (opB == "replace" || opB == "add"):
		return "patch deletes a path that another patch writes"
	case opB == "delete" && (opA == "replace" || opA == "add"):
		return "patch deletes a path that another patch writes"
	case opA == "delete" && opB == "delete":
		return "both patches delete the same path"
	default:
		return ""
	}
}

// FormatConflicts returns a human-readable summary of the given conflicts.
func FormatConflicts(conflicts []Conflict) string {
	if len(conflicts) == 0 {
		return "no conflicts detected\n"
	}
	out := fmt.Sprintf("%d conflict(s) detected:\n", len(conflicts))
	for _, c := range conflicts {
		out += "  " + c.String() + "\n"
	}
	return out
}
