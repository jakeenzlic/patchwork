package patch

import (
	"fmt"
	"strings"
)

// Apply applies a Patch to the given config map, returning a new modified map.
func Apply(config map[string]interface{}, p *Patch) (map[string]interface{}, error) {
	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("invalid patch: %w", err)
	}

	// Deep copy to avoid mutating the original.
	result := deepCopy(config)

	for _, op := range p.Ops {
		var err error
		switch op.Op {
		case OpAdd, OpReplace:
			err = setPath(result, op.Path, op.Value)
		case OpRemove:
			err = deletePath(result, op.Path)
		}
		if err != nil {
			return nil, fmt.Errorf("applying op %q on path %q: %w", op.Op, op.Path, err)
		}
	}
	return result, nil
}

func splitPath(path string) []string {
	return strings.Split(strings.TrimPrefix(path, "/"), "/")
}

func setPath(m map[string]interface{}, path string, value interface{}) error {
	keys := splitPath(path)
	current := m
	for i, key := range keys[:len(keys)-1] {
		next, ok := current[key]
		if !ok {
			nested := make(map[string]interface{})
			current[key] = nested
			current = nested
			continue
		}
		nested, ok := next.(map[string]interface{})
		if !ok {
			return fmt.Errorf("key %q (segment %d) is not an object", key, i)
		}
		current = nested
	}
	current[keys[len(keys)-1]] = value
	return nil
}

func deletePath(m map[string]interface{}, path string) error {
	keys := splitPath(path)
	current := m
	for _, key := range keys[:len(keys)-1] {
		next, ok := current[key]
		if !ok {
			return fmt.Errorf("path segment %q not found", key)
		}
		nested, ok := next.(map[string]interface{})
		if !ok {
			return fmt.Errorf("path segment %q is not an object", key)
		}
		current = nested
	}
	delete(current, keys[len(keys)-1])
	return nil
}

func deepCopy(m map[string]interface{}) map[string]interface{} {
	copy := make(map[string]interface{}, len(m))
	for k, v := range m {
		if nested, ok := v.(map[string]interface{}); ok {
			copy[k] = deepCopy(nested)
		} else {
			copy[k] = v
		}
	}
	return copy
}
