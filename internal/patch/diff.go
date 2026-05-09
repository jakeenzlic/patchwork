package patch

import (
	"fmt"
	"strings"
)

// Diff computes the list of patch operations needed to transform src into dst.
// Both src and dst should be unmarshalled JSON/YAML documents (map[string]interface{}).
func Diff(src, dst map[string]interface{}) ([]Operation, error) {
	var ops []Operation
	err := diffValues(src, dst, "", &ops)
	if err != nil {
		return nil, err
	}
	return ops, nil
}

func diffValues(src, dst interface{}, path string, ops *[]Operation) error {
	srcMap, srcIsMap := src.(map[string]interface{})
	dstMap, dstIsMap := dst.(map[string]interface{})

	if srcIsMap && dstIsMap {
		// Find keys removed or changed
		for k, sv := range srcMap {
			child := joinPath(path, k)
			dv, exists := dstMap[k]
			if !exists {
				*ops = append(*ops, Operation{Op: "delete", Path: child})
				continue
			}
			if err := diffValues(sv, dv, child, ops); err != nil {
				return err
			}
		}
		// Find keys added
		for k, dv := range dstMap {
			child := joinPath(path, k)
			if _, exists := srcMap[k]; !exists {
				*ops = append(*ops, Operation{Op: "add", Path: child, Value: dv})
			}
		}
		return nil
	}

	if path == "" {
		return fmt.Errorf("diff: top-level values must be objects")
	}

	if !equal(src, dst) {
		*ops = append(*ops, Operation{Op: "replace", Path: path, Value: dst})
	}
	return nil
}

func joinPath(base, key string) string {
	if base == "" {
		return key
	}
	return strings.Join([]string{base, key}, ".")
}

func equal(a, b interface{}) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}
