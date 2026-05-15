package patch

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// SignaturePath returns the path to the signature file for a given config.
func SignaturePath(configPath string) string {
	return configPath + ".sig"
}

// SignPatch computes an HMAC-SHA256 signature over the canonical JSON of the
// patch operations using the provided secret key.
func SignPatch(p *Patch, secret []byte) (string, error) {
	if len(secret) == 0 {
		return "", errors.New("signature: secret key must not be empty")
	}
	data, err := json.Marshal(p.Ops)
	if err != nil {
		return "", fmt.Errorf("signature: marshal ops: %w", err)
	}
	mac := hmac.New(sha256.New, secret)
	mac.Write(data)
	return hex.EncodeToString(mac.Sum(nil)), nil
}

// VerifySignature checks that the stored signature for a patch matches the
// expected HMAC-SHA256 computed with secret.
func VerifySignature(p *Patch, secret []byte) error {
	if p.Signature == "" {
		return fmt.Errorf("signature: patch %q has no signature", p.ID)
	}
	expected, err := SignPatch(p, secret)
	if err != nil {
		return err
	}
	if !hmac.Equal([]byte(p.Signature), []byte(expected)) {
		return fmt.Errorf("signature: patch %q signature mismatch", p.ID)
	}
	return nil
}

// SaveSignature writes the signature string to a .sig file beside the patch.
func SaveSignature(patchPath, sig string) error {
	return os.WriteFile(SignaturePath(patchPath), []byte(sig), 0o644)
}

// LoadSignature reads the signature from the .sig file beside the patch.
func LoadSignature(patchPath string) (string, error) {
	data, err := os.ReadFile(SignaturePath(patchPath))
	if err != nil {
		return "", fmt.Errorf("signature: read sig file: %w", err)
	}
	return string(data), nil
}
