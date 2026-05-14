package patch

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// TagIndex maps tag names to lists of patch IDs.
type TagIndex map[string][]string

var validTagRe = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)

// ValidateTag returns an error if the tag name contains invalid characters.
func ValidateTag(tag string) error {
	if tag == "" {
		return fmt.Errorf("tag name must not be empty")
	}
	if !validTagRe.MatchString(tag) {
		return fmt.Errorf("tag %q contains invalid characters (allowed: a-z, A-Z, 0-9, _, -)" , tag)
	}
	return nil
}

// BuildTagIndex scans a slice of patches and builds an index from tag → patch IDs.
func BuildTagIndex(patches []Patch) (TagIndex, error) {
	idx := make(TagIndex)
	for _, p := range patches {
		for _, tag := range p.Tags {
			if err := ValidateTag(tag); err != nil {
				return nil, fmt.Errorf("patch %q: %w", p.ID, err)
			}
			idx[tag] = append(idx[tag], p.ID)
		}
	}
	return idx, nil
}

// FilterByTag returns only the patches that carry the given tag.
func FilterByTag(patches []Patch, tag string) []Patch {
	var out []Patch
	for _, p := range patches {
		for _, t := range p.Tags {
			if t == tag {
				out = append(out, p)
				break
			}
		}
	}
	return out
}

// FormatTagIndex returns a human-readable summary of the tag index.
func FormatTagIndex(idx TagIndex) string {
	if len(idx) == 0 {
		return "(no tags)"
	}
	keys := make([]string, 0, len(idx))
	for k := range idx {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%-20s %s\n", k, strings.Join(idx[k], ", "))
	}
	return sb.String()
}
