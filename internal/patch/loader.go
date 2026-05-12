package patch

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadFromFile reads a patch definition from a JSON or YAML file.
func LoadFromFile(path string) (*Patch, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading patch file: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(path))
	var p Patch
	switch ext {
	case ".json":
		if err := json.Unmarshal(data, &p); err != nil {
			return nil, fmt.Errorf("parsing JSON patch: %w", err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &p); err != nil {
			return nil, fmt.Errorf("parsing YAML patch: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported file extension %q (use .json, .yaml, or .yml)", ext)
	}

	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("validating patch from %s: %w", path, err)
	}
	return &p, nil
}

// LoadDir loads all patch files from a directory, sorted by filename.
func LoadDir(dir string) ([]*Patch, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading patch directory: %w", err)
	}

	var patches []*Patch
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !isSupportedExt(entry.Name()) {
			continue
		}
		p, err := LoadFromFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("loading %s: %w", entry.Name(), err)
		}
		patches = append(patches, p)
	}
	return patches, nil
}

// isSupportedExt reports whether the filename has a supported patch file extension.
func isSupportedExt(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".json" || ext == ".yaml" || ext == ".yml"
}
