package patch

import (
	"fmt"
	"strings"
)

// Annotation holds a user-defined note attached to a patch operation by path.
type Annotation struct {
	PatchID string `json:"patch_id" yaml:"patch_id"`
	Path    string `json:"path"     yaml:"path"`
	Note    string `json:"note"     yaml:"note"`
	Author  string `json:"author"   yaml:"author"`
}

// Annotate attaches notes to specific operation paths within a patch.
// It validates that each annotated path exists among the patch operations.
func Annotate(p Patch, annotations []Annotation) ([]Annotation, error) {
	pathSet := make(map[string]struct{}, len(p.Ops))
	for _, op := range p.Ops {
		pathSet[op.Path] = struct{}{}
	}

	var result []Annotation
	for _, ann := range annotations {
		if ann.PatchID != p.ID {
			return nil, fmt.Errorf("annotation patch_id %q does not match patch id %q", ann.PatchID, p.ID)
		}
		if strings.TrimSpace(ann.Note) == "" {
			return nil, fmt.Errorf("annotation for path %q has empty note", ann.Path)
		}
		if _, ok := pathSet[ann.Path]; !ok {
			return nil, fmt.Errorf("annotation path %q not found in patch %q", ann.Path, p.ID)
		}
		result = append(result, ann)
	}
	return result, nil
}

// FormatAnnotations returns a human-readable summary of annotations.
func FormatAnnotations(annotations []Annotation) string {
	if len(annotations) == 0 {
		return "no annotations"
	}
	var sb strings.Builder
	for _, ann := range annotations {
		author := ann.Author
		if author == "" {
			author = "unknown"
		}
		fmt.Fprintf(&sb, "[%s] %s — %s (by %s)\n", ann.PatchID, ann.Path, ann.Note, author)
	}
	return strings.TrimRight(sb.String(), "\n")
}
