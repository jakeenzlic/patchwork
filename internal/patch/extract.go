package patch

import "fmt"

// ExtractResult holds the subset of a config matching the given path prefix.
type ExtractResult struct {
	Prefix string
	Data   map[string]any
}

// Extract returns the nested value rooted at prefix from config, presented as
// a flat map relative to that prefix. Returns an error when the prefix does
// not exist or resolves to a non-map value.
func Extract(config map[string]any, prefix string) (ExtractResult, error) {
	if prefix == "" {
		return ExtractResult{}, fmt.Errorf("extract: prefix must not be empty")
	}
	if prefix[0] == '/' {
		prefix = prefix[1:]
	}

	segments := splitPath(prefix)
	var cur any = config
	for _, seg := range segments {
		m, ok := cur.(map[string]any)
		if !ok {
			return ExtractResult{}, fmt.Errorf("extract: path segment %q is not a map", seg)
		}
		cur, ok = m[seg]
		if !ok {
			return ExtractResult{}, fmt.Errorf("extract: path segment %q not found", seg)
		}
	}

	m, ok := cur.(map[string]any)
	if !ok {
		return ExtractResult{}, fmt.Errorf("extract: value at %q is not a map", prefix)
	}

	return ExtractResult{
		Prefix: prefix,
		Data:   deepCopyMap(m),
	}, nil
}

// FormatExtractResult returns a human-readable summary of the extraction.
func FormatExtractResult(r ExtractResult) string {
	out := fmt.Sprintf("Extracted %d key(s) from prefix %q:\n", len(r.Data), r.Prefix)
	for k := range r.Data {
		out += fmt.Sprintf("  %s\n", k)
	}
	return out
}

func deepCopyMap(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		switch val := v.(type) {
		case map[string]any:
			out[k] = deepCopyMap(val)
		default:
			out[k] = val
		}
	}
	return out
}
