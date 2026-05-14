package patch

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadSchema reads a schema file (JSON or YAML) from disk.
func LoadSchema(path string) (Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Schema{}, fmt.Errorf("read schema: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(path))
	var schema Schema
	switch ext {
	case ".json":
		if err := json.Unmarshal(data, &schema); err != nil {
			return Schema{}, fmt.Errorf("parse schema JSON: %w", err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &schema); err != nil {
			return Schema{}, fmt.Errorf("parse schema YAML: %w", err)
		}
	default:
		return Schema{}, fmt.Errorf("unsupported schema format: %q", ext)
	}
	return schema, nil
}
