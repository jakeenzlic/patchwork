package patch

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Format represents the output serialization format.
type Format string

const (
	FormatJSON Format = "json"
	FormatYAML Format = "yaml"
)

// Export writes the given config map to a file in the specified format.
// The format is inferred from the file extension if fmt is empty.
func Export(config map[string]any, destPath string, fmt Format) error {
	if fmt == "" {
		fmt = inferFormat(destPath)
	}

	var data []byte
	var err error

	switch fmt {
	case FormatJSON:
		data, err = json.MarshalIndent(config, "", "  ")
		if err != nil {
			return fmt.Errorf("export: failed to marshal JSON: %w", err)
		}
		data = append(data, '\n')
	case FormatYAML:
		data, err = yaml.Marshal(config)
		if err != nil {
			return fmt.Errorf("export: failed to marshal YAML: %w", err)
		}
	default:
		return fmt.Errorf("export: unsupported format %q", fmt)
	}

	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return fmt.Errorf("export: failed to create directories: %w", err)
	}

	if err := os.WriteFile(destPath, data, 0o644); err != nil {
		return fmt.Errorf("export: failed to write file: %w", err)
	}

	return nil
}

func inferFormat(path string) Format {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".yaml", ".yml":
		return FormatYAML
	default:
		return FormatJSON
	}
}
