package patch

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CompressedBundlePath returns the path for a compressed patch bundle.
func CompressedBundlePath(dir, name string) string {
	return filepath.Join(dir, name+".patch.gz")
}

// Bundle compresses a slice of patches into a gzip-encoded JSON file.
func Bundle(patches []Patch, dest string) error {
	data, err := json.Marshal(patches)
	if err != nil {
		return fmt.Errorf("bundle: marshal: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return fmt.Errorf("bundle: mkdir: %w", err)
	}

	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("bundle: create: %w", err)
	}
	defer f.Close()

	gw := gzip.NewWriter(f)
	defer gw.Close()

	if _, err := io.Copy(gw, bytes.NewReader(data)); err != nil {
		return fmt.Errorf("bundle: write: %w", err)
	}
	return nil
}

// Unbundle reads a gzip-encoded JSON bundle and returns the patches.
func Unbundle(src string) ([]Patch, error) {
	f, err := os.Open(src)
	if err != nil {
		return nil, fmt.Errorf("unbundle: open: %w", err)
	}
	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("unbundle: gzip reader: %w", err)
	}
	defer gr.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, gr); err != nil {
		return nil, fmt.Errorf("unbundle: read: %w", err)
	}

	var patches []Patch
	if err := json.Unmarshal(buf.Bytes(), &patches); err != nil {
		return nil, fmt.Errorf("unbundle: unmarshal: %w", err)
	}
	return patches, nil
}

// IsBundle reports whether the given path looks like a compressed bundle.
func IsBundle(path string) bool {
	return strings.HasSuffix(path, ".patch.gz")
}
