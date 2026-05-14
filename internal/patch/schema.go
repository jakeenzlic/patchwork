package patch

import (
	"fmt"
	"strings"
)

// SchemaRule defines a constraint on a config path.
type SchemaRule struct {
	Path     string `json:"path" yaml:"path"`
	Type     string `json:"type" yaml:"type"`         // string, number, bool, object, array
	Required bool   `json:"required" yaml:"required"` // must exist after patch
}

// Schema holds a collection of rules for post-apply validation.
type Schema struct {
	Rules []SchemaRule `json:"rules" yaml:"rules"`
}

// ValidateSchema checks that the given config satisfies all schema rules.
func ValidateSchema(config map[string]any, schema Schema) []string {
	var violations []string
	for _, rule := range schema.Rules {
		val, exists := getValueAtPath(config, rule.Path)
		if !exists {
			if rule.Required {
				violations = append(violations, fmt.Sprintf("required path %q is missing", rule.Path))
			}
			continue
		}
		if rule.Type != "" {
			if err := checkType(val, rule.Type, rule.Path); err != nil {
				violations = append(violations, err.Error())
			}
		}
	}
	return violations
}

func getValueAtPath(config map[string]any, path string) (any, bool) {
	segments := strings.Split(strings.TrimPrefix(path, "."), ".")
	var cur any = config
	for _, seg := range segments {
		if seg == "" {
			continue
		}
		m, ok := cur.(map[string]any)
		if !ok {
			return nil, false
		}
		cur, ok = m[seg]
		if !ok {
			return nil, false
		}
	}
	return cur, true
}

func checkType(val any, expected string, path string) error {
	var actual string
	switch val.(type) {
	case string:
		actual = "string"
	case bool:
		actual = "bool"
	case float64, int, int64:
		actual = "number"
	case map[string]any:
		actual = "object"
	case []any:
		actual = "array"
	default:
		actual = "unknown"
	}
	if actual != expected {
		return fmt.Errorf("path %q: expected type %q but got %q", path, expected, actual)
	}
	return nil
}
