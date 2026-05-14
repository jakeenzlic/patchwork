package patch

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// parseConfig deserialises raw bytes into a generic map.
// The format is inferred from the filename extension of hint.
// Supported extensions are .json, .yaml, and .yml.
func parseConfig(raw []byte, hint string) (map[string]any, error) {
	ext := strings.ToLower(filepath.Ext(hint))
	var out map[string]any
	switch ext {
	case ".json":
		if err := json.Unmarshal(raw, &out); err != nil {
			return nil, fmt.Errorf("parse JSON: %w", err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(raw, &out); err != nil {
			return nil, fmt.Errorf("parse YAML: %w", err)
		}
		out = normaliseYAML(out)
	default:
		return nil, fmt.Errorf("unsupported config extension %q", ext)
	}
	return out, nil
}

// normaliseYAML converts map[interface{}]interface{} trees produced by the
// YAML decoder into map[string]any so they are consistent with JSON output.
func normaliseYAML(in map[string]any) map[string]any {
	if in == nil {
		return nil
	}
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = normaliseValue(v)
	}
	return out
}

func normaliseValue(v any) any {
	switch t := v.(type) {
	case map[string]any:
		return normaliseYAML(t)
	case map[interface{}]interface{}:
		m := make(map[string]any, len(t))
		for k, val := range t {
			m[fmt.Sprintf("%v", k)] = normaliseValue(val)
		}
		return m
	case []any:
		for i, item := range t {
			t[i] = normaliseValue(item)
		}
		return t
	default:
		return v
	}
}
